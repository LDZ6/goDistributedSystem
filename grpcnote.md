# gRPC 服务与客户端实现笔记

## gRPC 服务端实现

### 包声明和导入
```go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)
```

### 命令行标志
```go
var (
	port = flag.Int("port", 50051, "The server port")
)
```

### Greeter 服务实现
```go
type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
```

### 主函数
```go
func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```

## gRPC 客户端实现

### 包声明和导入
```go
package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)
```

### 命令行标志
```go
var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", "world", "Name to greet")
)
```

### 主函数
```go
func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
```

## 总结
以上代码分别实现了一个简单的 gRPC 服务器和客户端，通过 Protocol Buffers 定义了 Greeter 服务的通信协议，并生成了对应的自动生成代码进行通信。服务器接收客户端的问候，并返回相应的问候消息。

# 使用 gRPC 插件生成代码的笔记

## 添加 gRPC 插件

### 在 .proto 文件中添加 gRPC 插件声明
```protobuf
syntax = "proto3";

package helloworld;

option go_package = "helloworldpb";

service Greeter {
    rpc SayHello (HelloRequest) returns (HelloReply) {}
}

message HelloRequest {
    string name = 1;
}

message HelloReply {
    string message = 1;
}
```

### 使用 protoc 编译 .proto 文件并生成 gRPC 代码
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative helloworld.proto
```

## 生成的代码说明

- `--go_out` 标志用于生成 Protocol Buffers 的 Go 语言代码。
- `--go-grpc_out` 标志用于生成 gRPC 的 Go 语

言代码。
- `--go_opt=paths=source_relative` 和 `--go-grpc_opt=paths=source_relative` 标志指定生成的代码相对于源文件的路径。
- 生成的 Go 语言代码将包含在指定的包中，并且会根据服务和消息的名称生成相应的结构体、接口和方法。
