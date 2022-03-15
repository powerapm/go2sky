package tracer

import (
	"log"

	"github.com/powerapm/go2sky"
	"github.com/powerapm/go2sky/reporter"
)

/*
    创建GRPC的Tracer对象，用于上报采集数据（生产环境使用）
	oapAddr: 上报服务端地址，示例 127.0.0.1:11800
	serviceName: 上报服务名称，用于区分同一类应用
	instanceUUID: 上报应用实例的唯一ID，传递空是将生成随机UUID，服务端展示为服务名-pid:PID@IP 示例  app-pid:2839@127.0.0.1
	              如果传递为NAME:ip:port则注册到服务端的为ip:port，示例 NAME:127.0.0.1:8080
*/
func CreateGRPCTracer(oapAddr string, serviceName string, instanceUUID string) (tracer *go2sky.Tracer, err error) {
	// Use gRPC reporter for production
	reporter, err := reporter.NewGRPCReporter(oapAddr)
	if err != nil {
		log.Fatalf("new grpc reporter error %v \n", err)
	}

	tracer, err = go2sky.NewTracer(serviceName, go2sky.WithReporter(reporter), go2sky.WithInstance(instanceUUID))
	if err != nil {
		log.Fatalf("create grpc reporter tracer error %v \n", err)
	}

	return tracer, err

}

/*
    创建Log的Tracer对象，用于上报采集数据（测试环境使用）
	serviceName: 上报服务名称，用于区分同一类应用
*/
func CreateLogTracer(serviceName string) (tracer *go2sky.Tracer, err error) {
	// Use gRPC reporter for production
	reporter, err := reporter.NewLogReporter()
	if err != nil {
		log.Fatalf("new log reporter error %v \n", err)
	}

	tracer, err = go2sky.NewTracer(serviceName, go2sky.WithReporter(reporter))
	if err != nil {
		log.Fatalf("create log reporter tracer error %v \n", err)
	}
	//等待注册完成
	tracer.WaitUntilRegister()
	return tracer, err

}
