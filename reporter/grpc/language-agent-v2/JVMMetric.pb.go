// Code generated by protoc-gen-go. DO NOT EDIT.
// source: language-agent-v2/JVMMetric.proto

package language_agent_v2

import (
	context "context"
	fmt "fmt"
	common "github.com/powerapm/go2sky/reporter/grpc/common"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type JVMMetricCollection struct {
	Metrics              []*common.JVMMetric `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"`
	ServiceInstanceId    int32               `protobuf:"varint,2,opt,name=serviceInstanceId,proto3" json:"serviceInstanceId,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *JVMMetricCollection) Reset()         { *m = JVMMetricCollection{} }
func (m *JVMMetricCollection) String() string { return proto.CompactTextString(m) }
func (*JVMMetricCollection) ProtoMessage()    {}
func (*JVMMetricCollection) Descriptor() ([]byte, []int) {
	return fileDescriptor_a5bd38fe036677f3, []int{0}
}

func (m *JVMMetricCollection) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JVMMetricCollection.Unmarshal(m, b)
}
func (m *JVMMetricCollection) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JVMMetricCollection.Marshal(b, m, deterministic)
}
func (m *JVMMetricCollection) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JVMMetricCollection.Merge(m, src)
}
func (m *JVMMetricCollection) XXX_Size() int {
	return xxx_messageInfo_JVMMetricCollection.Size(m)
}
func (m *JVMMetricCollection) XXX_DiscardUnknown() {
	xxx_messageInfo_JVMMetricCollection.DiscardUnknown(m)
}

var xxx_messageInfo_JVMMetricCollection proto.InternalMessageInfo

func (m *JVMMetricCollection) GetMetrics() []*common.JVMMetric {
	if m != nil {
		return m.Metrics
	}
	return nil
}

func (m *JVMMetricCollection) GetServiceInstanceId() int32 {
	if m != nil {
		return m.ServiceInstanceId
	}
	return 0
}

func init() {
	proto.RegisterType((*JVMMetricCollection)(nil), "JVMMetricCollection")
}

func init() { proto.RegisterFile("language-agent-v2/JVMMetric.proto", fileDescriptor_a5bd38fe036677f3) }

var fileDescriptor_a5bd38fe036677f3 = []byte{
	// 288 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0x41, 0x4b, 0xc3, 0x30,
	0x14, 0xc7, 0xed, 0x44, 0x87, 0xf1, 0xa2, 0x9d, 0xc8, 0xe8, 0x69, 0x0e, 0x0f, 0x3b, 0x6c, 0x09,
	0x74, 0x17, 0xaf, 0x3a, 0x10, 0x36, 0xa8, 0x94, 0x15, 0x26, 0x78, 0xcb, 0xb2, 0x90, 0x85, 0x34,
	0x79, 0x25, 0xc9, 0x3a, 0xfa, 0x25, 0xfc, 0x20, 0x7e, 0x4a, 0x59, 0xdb, 0xcd, 0xc3, 0x3c, 0x05,
	0x7e, 0xef, 0xe5, 0xfd, 0xff, 0xfc, 0xd0, 0x53, 0x4e, 0x8d, 0xd8, 0x51, 0xc1, 0x27, 0x54, 0x70,
	0xe3, 0x27, 0x65, 0x4c, 0x16, 0xab, 0x24, 0xe1, 0xde, 0x4a, 0x86, 0x0b, 0x0b, 0x1e, 0xa2, 0x1e,
	0x03, 0xad, 0xc1, 0x90, 0xe6, 0x69, 0xe1, 0x5d, 0x0b, 0x17, 0xab, 0xa4, 0x21, 0x43, 0x89, 0x7a,
	0xa7, 0x9f, 0x33, 0xc8, 0x73, 0xce, 0xbc, 0x04, 0x13, 0x3e, 0xa3, 0xae, 0xae, 0x99, 0xeb, 0x07,
	0x83, 0xcb, 0xd1, 0x6d, 0x8c, 0xf0, 0x69, 0x6d, 0x79, 0x1c, 0x85, 0x63, 0x74, 0xef, 0xb8, 0x2d,
	0x25, 0xe3, 0x73, 0xe3, 0x3c, 0x35, 0x8c, 0xcf, 0x37, 0xfd, 0xce, 0x20, 0x18, 0x5d, 0x2d, 0xcf,
	0x07, 0xf1, 0x3b, 0x7a, 0xfc, 0xbb, 0xc1, 0x0b, 0xb0, 0x3e, 0x6b, 0x76, 0xc2, 0x31, 0xea, 0xb2,
	0x26, 0x3b, 0x7c, 0xc0, 0xff, 0xd4, 0x89, 0x6e, 0xf0, 0x0c, 0xb4, 0xa6, 0x66, 0xe3, 0x86, 0x17,
	0x6f, 0xdf, 0x01, 0x9a, 0x82, 0x15, 0x98, 0x16, 0x94, 0x6d, 0x39, 0x76, 0xaa, 0xda, 0xd3, 0x5c,
	0x49, 0x73, 0x20, 0x1a, 0x1b, 0xee, 0xf7, 0x60, 0x15, 0x3e, 0x1a, 0xc2, 0xb5, 0x21, 0x5c, 0xc6,
	0x69, 0xf0, 0xf5, 0x22, 0xa4, 0xdf, 0xee, 0xd6, 0x98, 0x81, 0x26, 0x99, 0xaa, 0x5e, 0xd3, 0x84,
	0x08, 0x88, 0x9d, 0xaa, 0x88, 0xad, 0xfb, 0x70, 0x4b, 0x84, 0x2d, 0x18, 0x39, 0xb3, 0xfb, 0xd3,
	0x89, 0x32, 0x55, 0x7d, 0xb6, 0x31, 0x1f, 0x4d, 0x44, 0x7a, 0x10, 0xc8, 0x20, 0x5f, 0x5f, 0xd7,
	0x2a, 0xa7, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xc1, 0x26, 0x7e, 0x2f, 0x96, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// JVMMetricReportServiceClient is the client API for JVMMetricReportService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type JVMMetricReportServiceClient interface {
	Collect(ctx context.Context, in *JVMMetricCollection, opts ...grpc.CallOption) (*common.Commands, error)
}

type jVMMetricReportServiceClient struct {
	cc *grpc.ClientConn
}

func NewJVMMetricReportServiceClient(cc *grpc.ClientConn) JVMMetricReportServiceClient {
	return &jVMMetricReportServiceClient{cc}
}

func (c *jVMMetricReportServiceClient) Collect(ctx context.Context, in *JVMMetricCollection, opts ...grpc.CallOption) (*common.Commands, error) {
	out := new(common.Commands)
	err := c.cc.Invoke(ctx, "/JVMMetricReportService/collect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// JVMMetricReportServiceServer is the server API for JVMMetricReportService service.
type JVMMetricReportServiceServer interface {
	Collect(context.Context, *JVMMetricCollection) (*common.Commands, error)
}

func RegisterJVMMetricReportServiceServer(s *grpc.Server, srv JVMMetricReportServiceServer) {
	s.RegisterService(&_JVMMetricReportService_serviceDesc, srv)
}

func _JVMMetricReportService_Collect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JVMMetricCollection)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JVMMetricReportServiceServer).Collect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/JVMMetricReportService/Collect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JVMMetricReportServiceServer).Collect(ctx, req.(*JVMMetricCollection))
	}
	return interceptor(ctx, in, info, handler)
}

var _JVMMetricReportService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "JVMMetricReportService",
	HandlerType: (*JVMMetricReportServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "collect",
			Handler:    _JVMMetricReportService_Collect_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "language-agent-v2/JVMMetric.proto",
}
