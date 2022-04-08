// Licensed to SkyAPM org under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. SkyAPM org licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package go2sky

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/powerapm/go2sky/internal/tool"
	"github.com/powerapm/go2sky/propagation"
)

const (
	errParameter = tool.Error("parameter are nil")
	EmptyTraceID = "N/A"
	NoopTraceID  = "[Ignored Trace]"
	// -1 represent the object doesn't exist.
	Inexistence = -1
)

// Tracer is go2sky tracer implementation.
type Tracer struct {
	service  string
	instance string
	reporter Reporter
	// 0 not init 1 init
	initFlag   int32
	serviceID  int32
	instanceID int32
	wg         *sync.WaitGroup
}

// TracerOption allows for functional options to adjust behaviour
// of a Tracer to be created by NewTracer
type TracerOption func(t *Tracer)

// NewTracer return a new go2sky Tracer
func NewTracer(service string, opts ...TracerOption) (tracer *Tracer, err error) {
	if service == "" {
		return nil, errParameter
	}
	t := &Tracer{
		service:    service,
		initFlag:   0,
		serviceID:  0,
		instanceID: 0,
	}
	for _, opt := range opts {
		opt(t)
	}
	if t.reporter != nil {
		if t.instance == "" {
			id, err := uuid.NewUUID()
			if err != nil {
				return nil, err
			}
			t.instance = id.String()
		}
		t.wg = &sync.WaitGroup{}
		t.wg.Add(1)
		go func() {
			defer t.wg.Done()
			for {
				serviceID, instanceID, err := t.reporter.Register(t.service, t.instance)
				if err != nil {
					time.Sleep(5 * time.Second)
					continue
				}
				if atomic.SwapInt32(&t.serviceID, serviceID) == 0 && atomic.SwapInt32(&t.instanceID, instanceID) == 0 {
					atomic.SwapInt32(&t.initFlag, 1)
					break
				}
			}
		}()
	}
	return t, nil
}

//WaitUntilRegister is a tool helps user to wait until register process has finished
func (t *Tracer) WaitUntilRegister() {
	if t.wg != nil {
		t.wg.Wait()
	}
}

// CreateEntrySpan creates and starts an entry span for incoming request
func (t *Tracer) CreateEntrySpan(ctx context.Context, operationName string, extractor propagation.Extractor) (s Span, nCtx context.Context, err error) {
	if ctx == nil || operationName == "" || extractor == nil {
		return nil, nil, errParameter
	}
	if s, nCtx = t.createNoop(ctx); s != nil {
		return
	}
	header, err := extractor()
	if err != nil {
		return
	}
	var refSc *propagation.SpanContext
	if header != "" {
		refSc = &propagation.SpanContext{}
		err = refSc.DecodeSW6(header)
		if err != nil {
			return
		}
	}
	s, nCtx, err = t.CreateLocalSpan(ctx, WithContext(refSc), WithSpanType(SpanTypeEntry))
	if err != nil {
		return
	}
	s.SetOperationName(operationName)
	ref, ok := nCtx.Value(refKeyInstance).(*propagation.SpanContext)
	if ok && ref != nil {
		return
	}
	sc := &propagation.SpanContext{
		Sample:                 1,
		ParentEndpoint:         operationName,
		EntryEndpoint:          operationName,
		EntryServiceInstanceID: t.instanceID,
	}
	if refSc != nil {
		sc.Sample = refSc.Sample
		if refSc.EntryEndpoint != "" {
			sc.EntryEndpoint = refSc.EntryEndpoint
		}
		sc.EntryEndpointID = refSc.EntryEndpointID
		sc.EntryServiceInstanceID = refSc.EntryServiceInstanceID
	}
	nCtx = context.WithValue(nCtx, refKeyInstance, sc)
	return
}

// CreateLocalSpan creates and starts a span for local usage
func (t *Tracer) CreateLocalSpan(ctx context.Context, opts ...SpanOption) (s Span, c context.Context, err error) {
	if ctx == nil {
		return nil, nil, errParameter
	}
	if s, c = t.createNoop(ctx); s != nil {
		return
	}
	ds := newLocalSpan(t)
	for _, opt := range opts {
		opt(ds)
	}
	parentSpan, ok := ctx.Value(ctxKeyInstance).(segmentSpan)
	if !ok {
		parentSpan = nil
	}
	s = newSegmentSpan(ds, parentSpan)
	return s, context.WithValue(ctx, ctxKeyInstance, s), nil
}

// CreateExitSpanWithContext creates and starts an exit span for client with context
func (t *Tracer) CreateExitSpanWithContext(ctx context.Context, operationName string, peer string,
	injector propagation.Injector) (s Span, nCtx context.Context, err error) {

	if ctx == nil || operationName == "" || peer == "" || injector == nil {
		return nil, nil, errParameter
	}
	if s, nCtx := t.createNoop(ctx); s != nil {
		return s, nCtx, nil
	}
	s, mCtx, err := t.CreateLocalSpan(ctx, WithSpanType(SpanTypeExit), WithOperationName(operationName))
	if err != nil {
		return s, mCtx, err
	}
	s.SetOperationName(operationName)
	s.SetPeer(peer)
	spanContext := &propagation.SpanContext{}
	span, ok := s.(ReportedSpan)
	if !ok {
		return nil, nCtx, errors.New("span type is wrong")
	}
	spanContext.Sample = 1
	spanContext.TraceID = span.Context().TraceID
	spanContext.ParentSpanID = span.Context().SpanID
	spanContext.ParentSegmentID = span.Context().SegmentID
	spanContext.NetworkAddress = peer
	spanContext.ParentServiceInstanceID = t.instanceID

	// Since 6.6.0, if first span is not entry span, then this is an internal segment(no RPC), rather than an endpoint.
	// EntryEndpoint
	firstSpan := span.Context().FirstSpan
	firstSpanOperationName := firstSpan.GetOperationName()
	ref, ok := ctx.Value(refKeyInstance).(*propagation.SpanContext)
	var entryEndpoint = ""
	var entryServiceInstanceID int32
	if ok && ref != nil {
		spanContext.Sample = ref.Sample
		entryEndpoint = ref.EntryEndpoint
		entryServiceInstanceID = ref.EntryServiceInstanceID
	} else {
		if firstSpan.IsEntry() {
			entryEndpoint = firstSpanOperationName
		} else {
			spanContext.EntryEndpointID = Inexistence
		}
		entryServiceInstanceID = t.instanceID
	}
	spanContext.EntryServiceInstanceID = entryServiceInstanceID
	if entryEndpoint != "" {
		spanContext.EntryEndpoint = entryEndpoint
	}
	// ParentEndpoint
	if firstSpan.IsEntry() && firstSpanOperationName != "" {
		spanContext.ParentEndpoint = firstSpanOperationName
	} else {
		spanContext.ParentEndpointID = Inexistence
	}

	err = injector(spanContext.EncodeSW6())
	if err != nil {
		return nil, nCtx, err
	}
	return s, nCtx, nil
}

// CreateExitSpan creates and starts an exit span for client
func (t *Tracer) CreateExitSpan(ctx context.Context, operationName string, peer string, injector propagation.Injector) (Span, error) {
	// if ctx == nil || operationName == "" || peer == "" || injector == nil {
	// 	return nil, errParameter
	// }
	// if s, _ := t.createNoop(ctx); s != nil {
	// 	return s, nil
	// }
	// s, _, err := t.CreateLocalSpan(ctx, WithSpanType(SpanTypeExit))
	// if err != nil {
	// 	return nil, err
	// }
	// s.SetOperationName(operationName)
	// s.SetPeer(peer)
	// spanContext := &propagation.SpanContext{}
	// span, ok := s.(ReportedSpan)
	// if !ok {
	// 	return nil, errors.New("span type is wrong")
	// }
	// spanContext.Sample = 1
	// spanContext.TraceID = span.Context().TraceID
	// spanContext.ParentSpanID = span.Context().SpanID
	// spanContext.ParentSegmentID = span.Context().SegmentID
	// spanContext.NetworkAddress = peer
	// spanContext.ParentServiceInstanceID = t.instanceID

	// // Since 6.6.0, if first span is not entry span, then this is an internal segment(no RPC), rather than an endpoint.
	// // EntryEndpoint
	// firstSpan := span.Context().FirstSpan
	// firstSpanOperationName := firstSpan.GetOperationName()
	// ref, ok := ctx.Value(refKeyInstance).(*propagation.SpanContext)
	// var entryEndpoint = ""
	// var entryServiceInstanceID int32
	// if ok && ref != nil {
	// 	spanContext.Sample = ref.Sample
	// 	entryEndpoint = ref.EntryEndpoint
	// 	entryServiceInstanceID = ref.EntryServiceInstanceID
	// } else {
	// 	if firstSpan.IsEntry() {
	// 		entryEndpoint = firstSpanOperationName
	// 	} else {
	// 		spanContext.EntryEndpointID = Inexistence
	// 	}
	// 	entryServiceInstanceID = t.instanceID
	// }
	// spanContext.EntryServiceInstanceID = entryServiceInstanceID
	// if entryEndpoint != "" {
	// 	spanContext.EntryEndpoint = entryEndpoint
	// }
	// // ParentEndpoint
	// if firstSpan.IsEntry() && firstSpanOperationName != "" {
	// 	spanContext.ParentEndpoint = firstSpanOperationName
	// } else {
	// 	spanContext.ParentEndpointID = Inexistence
	// }

	// err = injector(spanContext.EncodeSW6())
	// if err != nil {
	// 	return nil, err
	// }
	// return s, nil
	var s, _, err = t.CreateExitSpanWithContext(ctx, operationName, peer, injector)
	return s, err
}

func (t *Tracer) createNoop(ctx context.Context) (s Span, nCtx context.Context) {
	if ns, ok := ctx.Value(ctxKeyInstance).(*NoopSpan); ok {
		nCtx = ctx
		s = ns
		return
	}
	if t.initFlag == 0 {
		s = &NoopSpan{}
		nCtx = context.WithValue(ctx, ctxKeyInstance, s)
		return
	}
	return
}

type ctxKey struct{}

type refKey struct{}

var ctxKeyInstance = ctxKey{}

var refKeyInstance = refKey{}

//Reporter is a data transit specification
type Reporter interface {
	Register(service string, instance string) (int32, int32, error)
	Send(spans []ReportedSpan)
	Close()
}

func getFuncInfo() (info string) {
	// 获取上层调用者PC，文件名，所在行
	pc, codePath, codeLine, ok := runtime.Caller(2)
	if !ok {
		// 不ok，函数栈用尽了
		info = "-.-()"
	} else {
		// 拼接文件名与所在行,根据PC获取函数名
		info = fmt.Sprintf("%s.%s(%d)", codePath, runtime.FuncForPC(pc).Name(), codeLine)
	}
	return

}

/**
2022-04-07 黄尧 针对协程的处理,在协程函数调用一开始

eg:

//ctx为CallGoroutineWithContext()调用之后返回的上下文

 go ExcuteAll( ctx,tracer)

//被并发执行的函数，ctx为必传

 func ExcuteAll(ctx context.Context,tracer *go2sky.Tracer ) {
	callSpan, ctx, err := tracer.CallGoroutineWithContext(ctx)
	if err != nil {
		log.Fatalf("CallGoroutineWithContext失败 %v \n", err)
	}
	err = nil
	defer func() {
		if err != nil {
			err = fmt.Errorf("协程执行失败: %s", err.Error())
			callSpan.Error(time.Now(), err.Error())
		}
		callSpan.End()
	}()

	//doSomething
 }
**/
func (t *Tracer) CallGoroutineWithContext(ctx context.Context) (s Span, nCtx context.Context, err error) {
	if ctx == nil {
		return nil, nil, errParameter
	}
	operationName := "goroutine#" + getFuncInfo()
	var extractor propagation.Extractor
	extractor = func() (string, error) {
		return ctx.Value(propagation.Header).(string), nil
	}

	if s, nCtx = t.createNoop(ctx); s != nil {
		return
	}
	header, err := extractor()
	if err != nil {
		return
	}
	var refSc *propagation.SpanContext
	if header != "" {
		refSc = &propagation.SpanContext{}
		err = refSc.DecodeSW6(header)
		if err != nil {
			return
		}
	}
	s, nCtx, err = t.CreateLocalSpan(ctx, WithContext(refSc), WithSpanType(SpanTypeLocal))
	if err != nil {
		return
	}
	s.SetOperationName(operationName)
	ref, ok := nCtx.Value(refKeyInstance).(*propagation.SpanContext)
	if ok && ref != nil {
		return
	}
	sc := &propagation.SpanContext{
		Sample:                 1,
		ParentEndpoint:         operationName,
		EntryEndpoint:          operationName,
		EntryServiceInstanceID: t.instanceID,
	}
	if refSc != nil {
		sc.Sample = refSc.Sample
		if refSc.EntryEndpoint != "" {
			sc.EntryEndpoint = refSc.EntryEndpoint
		}
		sc.EntryEndpointID = refSc.EntryEndpointID
		sc.EntryServiceInstanceID = refSc.EntryServiceInstanceID
	}
	nCtx = context.WithValue(nCtx, refKeyInstance, sc)
	return
}

/**
2022-04-07 黄尧 针对协程的处理,在调用协程前进行调用

eg:

//入参：当前上下文，如果没有则传递context.Background()

 ctx, err := tracer.StartGoroutineWithContext(ginCtx.Request.Context())
 if err != nil {
	log.Fatalf("Start Goroutine error: %v", err)
 }

//执行协程，传递相关参数以及返回的上下文

 go ExcuteAll(tracer, ctx)

**/
func (t *Tracer) StartGoroutineWithContext(ctx context.Context) (headCtx context.Context, err error) {
	operationName := "call-goroutine"
	peer := "goroutine"
	var injector propagation.Injector
	headCtx = context.Background()
	injector = func(header string) error {
		headCtx = context.WithValue(headCtx, propagation.Header, header)
		return nil
	}
	if ctx == nil || operationName == "" || injector == nil {
		return nil, errParameter
	}
	var s Span
	defer func() {
		if s != nil {
			s.End()
		}
	}()
	if s, nCtx := t.createNoop(ctx); s != nil {
		return nCtx, nil
	}
	s, mCtx, err := t.CreateLocalSpan(ctx, WithSpanType(SpanTypeLocal), WithOperationName(operationName))
	if err != nil {
		return mCtx, err
	}
	s.SetOperationName(operationName)
	s.SetPeer(peer)
	spanContext := &propagation.SpanContext{}
	span, ok := s.(ReportedSpan)
	if !ok {
		return headCtx, errors.New("span type is wrong")
	}
	spanContext.Sample = 1
	spanContext.TraceID = span.Context().TraceID
	spanContext.ParentSpanID = span.Context().SpanID
	spanContext.ParentSegmentID = span.Context().SegmentID
	spanContext.NetworkAddress = peer
	spanContext.ParentServiceInstanceID = t.instanceID

	// Since 6.6.0, if first span is not entry span, then this is an internal segment(no RPC), rather than an endpoint.
	// EntryEndpoint
	firstSpan := span.Context().FirstSpan
	firstSpanOperationName := firstSpan.GetOperationName()
	ref, ok := ctx.Value(refKeyInstance).(*propagation.SpanContext)
	var entryEndpoint = ""
	var entryServiceInstanceID int32
	if ok && ref != nil {
		spanContext.Sample = ref.Sample
		entryEndpoint = ref.EntryEndpoint
		entryServiceInstanceID = ref.EntryServiceInstanceID
	} else {
		if firstSpan.IsEntry() {
			entryEndpoint = firstSpanOperationName
		} else {
			spanContext.EntryEndpointID = Inexistence
		}
		entryServiceInstanceID = t.instanceID
	}
	spanContext.EntryServiceInstanceID = entryServiceInstanceID
	if entryEndpoint != "" {
		spanContext.EntryEndpoint = entryEndpoint
	}
	// ParentEndpoint
	if firstSpan.IsEntry() && firstSpanOperationName != "" {
		spanContext.ParentEndpoint = firstSpanOperationName
	} else {
		spanContext.ParentEndpointID = Inexistence
	}

	err = injector(spanContext.EncodeSW6())
	if err != nil {
		return headCtx, err
	}
	return headCtx, nil
}

func TraceID(ctx context.Context) string {
	activeSpan := ctx.Value(ctxKeyInstance)
	if activeSpan == nil {
		return EmptyTraceID
	}
	span, ok := activeSpan.(segmentSpan)
	if ok {
		return span.context().GetReadableGlobalTraceID()
	}
	return NoopTraceID
}
