package main

import (
	"context"
	pb "helloworld-client/proto"
	"time"

	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
)

var (
	service = "helloworld" //需要和微服务服务端对应的service名统一,这样才能调用该微服务
	version = "latest"
)

func main() {
	//集成consul
	consulReg := consul.NewRegistry(
		//指定微服务的ip:  选择注册服务器地址,默认为本机,也可以选择consul集群中的client
		registry.Addrs("127.0.0.1:8500"),
	)
	// Create service
	srv := micro.NewService(
		//注册consul
		micro.Registry(consulReg),
	)
	srv.Init()

	// 创建客户端实例
	c := pb.NewHelloworldService(service, srv.Client())
	for {
		// Call service: CallRequest就是.proto中的
		rsp, err := c.Call(context.Background(), &pb.CallRequest{Name: "张三"})
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info(rsp)
		//每隔一段时间请求
		time.Sleep(2 * time.Second) // 每隔2秒请求
	}
}
