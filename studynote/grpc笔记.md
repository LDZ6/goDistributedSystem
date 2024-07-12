## 一、gRPC框架介绍

### 简介

gRPC 是一个高性能、开源和通用的 RPC 框架，面向移动端和 HTTP/2 设计。目前提供 C、Java 和 Go 语言版本，分别是：grpc, grpc-java, grpc-go，其中 C 版本支持 C, C++, Node.js, Python, Ruby, Objective-C, PHP 和 C#。

### gRPC 特点

1. **多语言支持**：提供几乎所有主流语言的实现，打破语言隔阂。
2. **HTTP/2 标准设计**：带来双向流、流控、头部压缩、单 TCP 连接上的多复用请求等优点。
3. **高性能**：默认使用 Protocol Buffers 序列化，性能相较于 RESTful JSON 更好。
4. **工具链成熟**：代码生成便捷，开箱即用。
5. **双向流式请求和响应**：对批量处理、低延时场景友好。

这些特性使得 gRPC 在移动设备上表现更好，更省电和节省空间占用。有了 gRPC，可以一次性在一个 .proto 文件中定义服务，并使用任何支持它的语言去实现客户端和服务端。

gRPC 默认使用 Protocol Buffers 作为序列化机制（当然也可以使用其他数据格式如 JSON），可以用 proto 文件创建 gRPC 服务，用 Protocol Buffers 消息类型来定义方法参数和返回类型。

在 gRPC 客户端，可以直接调用不同服务器上的远程程序，就像调用本地程序一样，很容易构建分布式应用和服务。客户端和服务器可以分别使用 gRPC 支持的不同语言实现。

### 参考资料

- gRPC 官方文档中文版：[http://doc.oschina.net/grpc?t=60133](http://doc.oschina.net/grpc?t=60133)
- gRPC 官网：[https://grpc.io](https://grpc.io)

## 二、gRPC的使用

### 安装 gRPC 包

gRPC 使用的包是 `google.golang.org/grpc`。可以在项目中使用以下命令下载包：

```sh
go get -u -v google.golang.org/grpc
```

也可以使用 `go mod tidy` 下载包。

### 案例1：实现基本的 gRPC 微服务

#### 1. 实现服务端

##### (1). 创建文件 main.go

在 `micro/grpc_demo/server/hello` 下创建文件 `main.go`：

```go
package main

import "fmt"

// RPC 远程调用的接口，需要实现 hello.proto 中定义的 Hello 接口，以及里面的方法

func main() {
    fmt.Println("hello")
}
```

然后通过命令 `go mod init hello` 初始化项目，并使用 `go mod tidy` 下载依赖。

```sh
go mod init hello
go mod tidy
```

测试是否成功：

```sh
go run .\main.go
```

如果打印出 `hello`，说明成功。

##### (2). 创建 .proto 文件

在 `micro/grpc_demo/server/hello/proto` 下创建 `hello.proto` 文件：

```proto
syntax = "proto3";  // proto 版本
option go_package = "./helloService"; // 表示在目录 helloService 下生成 hello.pb.go，以及对应的包名

// 通过 service 创建一个 RPC 服务，生成一个 Hello 接口
service Hello {
    // 通过 rpc 来指定远程调用的方法：
    // SayHello 方法，这个方法里面实现对传入的参数 HelloReq，以及返回的参数 HelloRes 进行约束
    rpc SayHello(HelloReq) returns (HelloRes);
}

// message 为传入的参数进行定义消息：结构体类型，这样就要求客户端传入一个结构体，结构体有一个字符串类型的 name 参数
message HelloReq {
    string name = 1;
}

// message 为返回的参数进行定义消息：结构体类型，这样就要求服务端返回一个结构体，结构体有一个字符串类型的 message 参数
message HelloRes {
    string message = 1;
}
```

##### (3). 编译 .proto 文件

使用命令 `protoc --go_out=plugins=grpc:. *.proto` 对 .proto 文件进行编译，生成 .pb.go 的服务。

##### (4). 生成的 .pb.go 文件重要方法/结构体讲解

1. **HelloServer**

```go
// HelloServer is the server API for Hello service.
type HelloServer interface {
    // SayHello 方法，实现对传入的参数 HelloReq，以及返回的参数 HelloRes 进行约束
    SayHello(context.Context, *HelloReq) (*HelloRes, error)
}
```

2. **RegisterHelloServer**

```go
// RegisterHelloServer 注册服务
func RegisterHelloServer(s *grpc.Server, srv HelloServer) {
    s.RegisterService(&_Hello_serviceDesc, srv)
}
```

3. **NewHelloClient**

```go
// NewHelloClient 注册客户端
func NewHelloClient(cc grpc.ClientConnInterface) HelloClient {
    return &helloClient{cc}
}
```

##### (5). 在 main.go 目录下运行命令 `go mod tidy` 加载生成的 .pb.go 中引入的包。

##### (6). 在 main.go 中编写服务端代码

```go
package main

import (
    "context"
    "fmt"
    "go_code/micro/grpc_demo/server/hello/proto/helloService"
    "net"
    "google.golang.org/grpc"
)

// grpc 远程调用的接口，需要实现 hello.proto 中定义的 Hello 服务接口，以及里面的方法
type Hello struct{}

// SayHello 方法参考 hello.pb.go 中的接口
func (this Hello) SayHello(c context.Context, req *helloService.HelloReq) (*helloService.HelloRes, error) {
    fmt.Println(req)
    return &helloService.HelloRes{
        Message: "你好 " + req.Name,
    }, nil
}

func main() {
    // 1. 初始化一个 grpc 对象
    grpcServer := grpc.NewServer()
    // 2. 注册服务
    helloService.RegisterHelloServer(grpcServer, new(Hello))
    // 3. 设置监听，指定 IP、port
    listener, err := net.Listen("tcp", "127.0.0.1:8080")
    if err != nil {
        fmt.Println(err)
    }
    // 4. 退出关闭监听
    defer listener.Close()
    // 5. 启动服务
    grpcServer.Serve(listener)
}
```

测试服务端是否启动成功：

```sh
go run .\main.go
```

如果没有错误，说明启动成功。

#### 2. 实现客户端

##### (1). 创建文件 main.go

在 `micro/grpc_demo/client/hello` 下创建文件 `main.go`。初始化项目并下载依赖：

```sh
go mod init hello
go mod tidy
```

测试是否成功：

```sh
go run .\main.go
```

如果打印出 `hello`，说明成功。

##### (2). 创建 .proto 文件

在 `micro/grpc_demo/client/proto` 下创建 `hello.proto` 文件，内容与服务端一致。

##### (3). 编译 .proto 文件

使用命令 `protoc --go_out=plugins=grpc:. *.proto` 对 .proto 文件进行编译，生成 .pb.go 的服务。

##### (4). 生成的 .pb.go 文件重要方法/结构体讲解

参考服务端步骤。

##### (5). 在 main.go 目录下运行命令 `go mod tidy` 加载生成的 .pb.go 中引入的包。

##### (6). 在 main.go 中编写客户端代码

```go
package main

// grpc 客户端代码

import (
    "clientHello/proto/helloService"
    "context"
    "fmt"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    // 1. 连接服务器
    grpcClient, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        fmt.Println(err)
    }

    // 2. 注册客户端
    client := helloService.NewHelloClient(grpcClient)
    // 3. 调用服务端函数，实现 HelloClient 接口：SayHello()
    res, err1 := client.SayHello(context.Background(), &helloService.HelloReq{
        Name: "张三",
    })
    if err1 != nil {
        fmt.Printf("调用服务端代码失败: %s", err1)
        return
    }

    fmt.Printf("%#v\r\n", res)
    fmt.Printf("调用成功: %s", res.Message)
}
```

#### 3. 测试客户端调用服务端微服务是否成功

##### (1). 启动服务端

在 `server/hello` 下运行命令：

```sh
go run .\main.go
```

##### (2). 启动客户端，请求服务端方法

在 `client/hello` 下运行命令：

```sh
go run .\main.go
```

如果返回了服务端的数据，说明微服务操作完成。