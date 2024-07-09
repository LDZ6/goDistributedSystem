# RPC 服务详解

## 1. ProtoBuf 示例

### 1.1 服务定义

```proto
syntax = "proto3";

package helloworld;

// 定义 Greeter 服务，包括 SayHello 和 SayName 方法
service Greeter {
  rpc SayHello (HelloRequest) returns (HelloReply);  // rpc 方法定义，接收 HelloRequest，返回 HelloReply
  rpc SayName (NameRequest) returns (NameReply);    // rpc 方法定义，接收 NameRequest，返回 NameReply
}

// 定义 HelloRequest 消息类型，包含一个字符串字段 name
message HelloRequest {
  string name = 1;  // 请求参数定义
}

// 定义 HelloReply 消息类型，包含一个字符串字段 message
message HelloReply {
  string message = 1;  // 响应参数定义
}

// 定义 NameRequest 消息类型，包含一个字符串字段 name
message NameRequest {
  string name = 1;
}

// 定义 NameReply 消息类型，包含一个字符串字段 message
message NameReply {
  string message = 1;
}
```

### 1.2 生成 Go 代码

使用 `protoc` 工具生成 Go 代码：

```bash
protoc --go_out=. helloworld.proto           # 生成普通的 Go 文件
protoc --go-grpc_out=. helloworld.proto      # 生成 gRPC 相关的 Go 文件
```

## 2. Go 中的 RPC

Go 标准包中提供了对 RPC 的支持，且支持三个级别的 RPC：TCP、HTTP、JSONRPC。但 Go 的 RPC 包独特，它只支持 Go 语言开发的服务器与客户端之间的交互，因为内部使用了 Gob 编码。

### 2.1 RPC 函数要求

1. 函数必须是导出的（首字母大写）。
2. 必须有两个导出类型的参数。
   - 第一个参数是接收的参数。
   - 第二个参数是返回给客户端的参数，且必须是指针类型。
3. 函数必须有一个返回值 `error`。

### 2.2 示例：RPC 函数

```go
func (t *T) MethodName(argType T1, replyType *T2) error
```

- `T` 是结构体类型。
- `T1` 是结构体类型。
- `T2` 是结构体指针类型。

### 2.3 服务端示例

```go
package main

import (
    "errors"
    "fmt"
    "net/http"
    "net/rpc"
)

// 定义参数结构体
type Args struct {
    A, B int
}

// 定义商和余数结构体
type Quotient struct {
    Quo, Rem int
}

// 定义 Arith 类型
type Arith int

// Multiply 方法：乘法操作，符合 RPC 函数要求
func (t *Arith) Multiply(args *Args, reply *int) error {
    *reply = args.A * args.B
    return nil
}

// Divide 方法：除法操作，符合 RPC 函数要求
func (t *Arith) Divide(args *Args, quo *Quotient) error {
    if args.B == 0 {
        return errors.New("divide by zero")  // 除零错误
    }
    quo.Quo = args.A / args.B
    quo.Rem = args.A % args.B
    return nil
}

func main() {
    arith := new(Arith)  // 创建 Arith 实例
    rpc.Register(arith)  // 注册 RPC 服务
    rpc.HandleHTTP()     // 绑定到 HTTP 协议

    err := http.ListenAndServe(":1234", nil)  // 启动 HTTP 服务器
    if err != nil {
        fmt.Println(err.Error())
    }
}
```

### 2.4 客户端示例

```go
package main

import (
    "fmt"
    "log"
    "net/rpc"
    "os"
)

// 定义参数结构体
type Args struct {
    A, B int
}

// 定义商和余数结构体
type Quotient struct {
    Quo, Rem int
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: ", os.Args[0], "server")
        os.Exit(1)
    }
    serverAddress := os.Args[1]

    // 连接到 RPC 服务器
    client, err := rpc.DialHTTP("tcp", serverAddress+":1234")
    if err != nil {
        log.Fatal("dialing:", err)
    }

    // 同步调用 Multiply 方法
    args := Args{17, 8}
    var reply int
    err = client.Call("Arith.Multiply", args, &reply)
    if err != nil {
        log.Fatal("arith error:", err)
    }
    fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

    // 同步调用 Divide 方法
    var quot Quotient
    err = client.Call("Arith.Divide", args, &quot)
    if err != nil {
        log.Fatal("arith error:", err)
    }
    fmt.Printf("Arith: %d/%d=%d remainder %d\n", args.A, args.B, quot.Quo, quot.Rem)
}
```

## 3. gRPC 示例

### 3.1 服务器端

```go
package main

import (
    "context"
    "fmt"
    "net"
    "google.golang.org/grpc"
    pd "myproject/proto"  // 引入生成的 protobuf 文件
)

// 定义服务端结构体
type server struct{}

// 实现 SayHello 方法
func (s *server) SayHello(ctx context.Context, in *pd.HelloRequest) (*pd.HelloReply, error) {
    return &pd.HelloReply{Message: "Hello " + in.Name}, nil
}

// 实现 SayName 方法
func (s *server) SayName(ctx context.Context, in *pd.NameRequest) (*pd.NameReply, error) {
    return &pd.NameReply{Message: "My name is " + in.Name}, nil
}

func main() {
    ln, err := net.Listen("tcp", ":1234")  // 监听 TCP 连接
    if err != nil {
        fmt.Println(err.Error())
    }
    srv := grpc.NewServer()  // 创建 gRPC 服务器
    pd.RegisterGreeterServer(srv, &server{})  // 注册 Greeter 服务
    err = srv.Serve(ln)  // 启动服务
    if err != nil {
        fmt.Println(err.Error())
    }
}
```

### 3.2 客户端

```go
package main

import (
    "context"
    "fmt"
    "log"
    "google.golang.org/grpc"
    pd "myproject/proto"  // 引入生成的 protobuf 文件
    "os"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: ", os.Args[0], "server")
        os.Exit(1)
    }
    serverAddress := os.Args[1]

    // 连接到 gRPC 服务器
    conn, err := grpc.Dial(serverAddress+":1234", grpc.WithInsecure())
    if err != nil {
        log.Fatal("dialing:", err)
    }
    defer conn.Close()

    // 创建 Greeter 客户端
    client := pd.NewGreeterClient(conn)
    name := "world"
    // 调用 SayHello 方法
    reply, err := client.SayHello(context.Background(), &pd.HelloRequest{Name: name})
    if err != nil {
        log.Fatal("say hello error:", err)
    }
    fmt.Printf("Greeter: Hello %s, %s\n", name, reply.Message)
}
```

---

### 注解说明

- **ProtoBuf 服务定义**：使用 `.proto` 文件定义服务和消息类型，生成 Go 代码后可用于 gRPC 通信。
- **RPC 函数要求**：Go 的 RPC 函数必须是导出的，具有两个参数（一个值参数，一个指针参数）和一个 `error` 返回值。
- **服务端和客户端示例**：展示如何在 Go 中实现并使用 RPC 服务，包括注册服务、处理请求和调用远程方法。