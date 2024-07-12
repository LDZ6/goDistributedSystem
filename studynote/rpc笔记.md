# 微服务和 RPC 框架简介

在学习微服务框架之前，我们需要先了解一下 RPC 架构。通过 RPC，可以更形象地理解微服务的工作流程。

## 一、简介

### RPC的概念

RPC (Remote Procedure Call Protocol)，即远程过程调用协议。通俗地说，就是调用远程的一个函数。相对应的是本地函数调用。

例如，当我们写下如下代码时：

```go
result := Add(1, 2)
```

传入了1和2两个参数，调用了本地代码中的一个 `Add` 函数，得到 `result` 这个返回值。这时，参数、返回值、代码段都在一个进程空间内，这就是本地函数调用。

那么有没有办法调用一个跨进程（例如这个进程部署在另一台服务器上）的函数呢？这就是 RPC 主要实现的功能，也是微服务的主要功能。

### RPC的优势

使用微服务化的一个好处就是：

1. 不限定服务的提供方使用什么技术选型，能够实现公司跨团队的技术解耦。
2. 每个服务都被封装成进程，彼此独立。
3. 使用微服务可以跨进程通信。

### RPC与IPC的区别

- IPC（Inter-Process Communication）：进程间通信
- RPC：远程进程调用 —— 应用层协议（和 HTTP 协议同层），底层使用 TCP 实现

在 Go 语言中实现 RPC 非常简单，有封装好的官方库和一些第三方库提供支持。Go RPC 可以利用 TCP 或 HTTP 来传递数据，可以对要传递的数据使用多种类型的编解码方式。Go 官方的 `net/rpc` 库使用 `encoding/gob` 进行编解码，支持 TCP 或 HTTP 数据传输方式，由于其他语言不支持 gob 编解码方式，所以使用 `net/rpc` 库实现的 RPC 方法无法进行跨语言调用。

Go 官方还提供了 `net/rpc/jsonrpc` 库实现 RPC 方法，JSON RPC 采用 JSON 进行数据编解码，因而支持跨语言调用。但目前的 jsonrpc 库是基于 TCP 协议实现的，暂时不支持使用 HTTP 进行数据传输。

除了 Go 官方提供的 RPC 库，还有许多第三方库为在 Go 中实现 RPC 提供支持，大部分第三方 RPC 库的实现都是使用 protobuf 进行数据编解码，根据 protobuf 声明文件自动生成 RPC 方法定义与服务注册代码，在 Go 中可以很方便地进行 RPC 服务调用。

## 二、`net/rpc` 库实现远程调用

### 使用 HTTP 作为 RPC 的载体实现远程调用

以下演示如何使用 Go 官方的 `net/rpc` 库实现 RPC 方法，使用 HTTP 作为 RPC 的载体，通过 `net/http` 包监听客户端连接请求。

#### 创建 RPC 微服务端

新建 `server/main.go` 文件：

```go
package main

import (
    "fmt"
    "log"
    "net"
    "net/http"
    "net/rpc"
    "os"
)

// 定义类对象
type World struct {}

// 绑定类方法
func (this *World) HelloWorld(req string, res *string) error {
    *res = req + " 你好!"
    return nil
}

// 绑定类方法
func (this *World) Print(req string, res *string) error {
    *res = req + " this is Print!"
    return nil
}

func main() {
    // 1. 注册RPC服务
    rpc.Register(new(World)) // 注册rpc服务
    rpc.HandleHTTP() // 采用http协议作为rpc载体

    // 2. 设置监听
    lis, err := net.Listen("tcp", "127.0.0.1:8800")
    if err != nil {
        log.Fatalln("fatal error: ", err)
    }
    fmt.Fprintf(os.Stdout, "%s", "start connection")

    // 3. 建立链接
    http.Serve(lis, nil)
}
```

注意：以上 `World` 结构体的方法必须满足 Go 语言的 RPC 规则：

1. 方法只能有两个可序列化的参数，其中第二个参数是指针类型，参数的类型不能是 channel（通道）、complex（复数类型）、func（函数），因为它们不能进行序列化。
2. 方法要返回一个 `error` 类型，同时必须是公开的方法。

#### 创建 RPC 客户端

客户端可以是 Go Web 也可以是一个 Go 应用，新建 `client/main.go` 文件：

```go
package main

import (
    "fmt"
    "net/rpc"
)

func main() {
    // 1. 用 rpc 链接服务器 --Dial()
    conn, err := rpc.DialHTTP("tcp", "127.0.0.1:8800")
    if err != nil {
        fmt.Println("Dial err:", err)
        return
    }
    defer conn.Close()

    // 2. 调用远程函数
    var reply1 string // 接受返回值 --- 传出参数
    err1 := conn.Call("World.HelloWorld", "张三", &reply1)
    if err1 != nil {
        fmt.Println("Call:", err1)
        return
    }
    fmt.Println(reply1)

    var reply2 string // 接受返回值 --- 传出参数
    err2 := conn.Call("World.Print", "李四", &reply2)
    if err2 != nil {
        fmt.Println("Call:", err2)
        return
    }
    fmt.Println(reply2)
}
```

### 使用 TCP 作为 RPC 的载体实现远程调用

#### 创建 RPC 微服务端

新建 `server/main.go` 文件：

```go
package main

import (
    "fmt"
    "net"
    "net/rpc"
)

// 定义类对象
type World struct {}

// 绑定类方法
func (this *World) HelloWorld(req string, res *string) error {
    *res = req + " 你好!"
    return nil
}

func main() {
    // 1. 注册RPC服务
    err := rpc.RegisterName("hello", new(World))
    if err != nil {
        fmt.Println("注册 rpc 服务失败!", err)
        return
    }

    // 2. 设置监听
    listener, err := net.Listen("tcp", "127.0.0.1:8800")
    if err != nil {
        fmt.Println("net.Listen err:", err)
        return
    }
    defer listener.Close()
    fmt.Println("开始监听 ...")

    // 3. 建立链接
    for {
        //接收连接
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Accept() err:", err)
            return
        }
        // 4. 绑定服务
        go rpc.ServeConn(conn)
    }
}
```

注意：以上 `World` 结构体的方法必须满足 Go 语言的 RPC 规则：

1. 方法只能有两个可序列化的参数，其中第二个参数是指针类型，参数的类型不能是 channel（通道）、complex（复数类型）、func（函数），因为它们不能进行序列化。
2. 方法要返回一个 `error` 类型，同时必须是公开的方法。

#### 创建 RPC 客户端

新建 `client/main.go` 文件：

```go
package main

import (
    "fmt"
    "net/rpc"
)

func main() {
    // 1. 用 rpc 链接服务器 --Dial()
    conn, err := rpc.Dial("tcp", "127.0.0.1:8800")
    if err != nil {
        fmt.Println("Dial err:", err)
        return
    }
    defer conn.Close()

    // 2. 调用远程函数
    var reply string // 接受返回值 --- 传出参数
    err = conn.Call("hello.HelloWorld", "张三", &reply)
    if err != nil {
        fmt.Println("Call:", err)
        return
    }
    fmt.Println(reply)
}
```

说明：

首先是通过 `rpc.Dial` 拨号 RPC 服务，然后通过 `client.Call` 调用具体的 RPC 方法。在调用 `client.Call` 时，第一个参数是用点号链接的 RPC 服务名字和方法名字，第二和第三个参数分别定义 RPC 方法的两个参数。

## 三、使用 TCP 作为 RPC 的载体实现远程调用具体案例

### 案例 1：简单使用

1. 创建一个 `hello` 微服务端，编写微服务端 RPC 代码，完成后启动该微服务端。
2. 创建一个 `hello` 客户端，编写客户端 RPC 代码，完成后启动该客户端，访问微服务端 RPC 功能，并返回相关数据。

#### 创建 `hello` 微服务端

新建 `mirco/server/hello/main.go` 文件，并编写代码：

```go
package main

import (
    "fmt"
    "net"
    "net/rpc"
)

// 定义类对象
type Hello struct{}

// 绑定类方法
func (this Hello) SayHello(req string, res *string) error {
    fmt.Println("请求的参数:", req)
    *res = "你好" + req
    return nil
}

func main() {
    // 1. 注册RPC服务
    err1 := rpc

.RegisterName("hello", new(Hello))
    if err1 != nil {
        fmt.Println("注册 rpc 服务失败!", err1)
        return
    }

    // 2. 设置监听
    listener, err2 := net.Listen("tcp", "127.0.0.1:8800")
    if err2 != nil {
        fmt.Println("net.Listen err:", err2)
        return
    }
    defer listener.Close()
    fmt.Println("开始监听 ...")

    // 3. 建立链接
    for {
        //接收连接
        conn, err3 := listener.Accept()
        if err3 != nil {
            fmt.Println("Accept() err:", err3)
            return
        }
        // 4. 绑定服务
        go rpc.ServeConn(conn)
    }
}
```

注意：以上 `Hello` 结构体的方法必须满足 Go 语言的 RPC 规则：

1. 方法只能有两个可序列化的参数，其中第二个参数是指针类型，参数的类型不能是 channel（通道）、complex（复数类型）、func（函数），因为它们不能进行序列化。
2. 方法要返回一个 `error` 类型，同时必须是公开的方法。

#### 创建 `hello` 客户端

新建 `mirco/client/hello/main.go` 文件，并编写代码：

```go
package main

import (
    "fmt"
    "net/rpc"
)

func main() {
    // 1. 用 rpc 链接服务器 --Dial()
    conn, err := rpc.Dial("tcp", "127.0.0.1:8800")
    if err != nil {
        fmt.Println("Dial err:", err)
        return
    }
    defer conn.Close()

    // 2. 调用远程函数
    var reply string // 接受返回值 --- 传出参数
    err = conn.Call("hello.SayHello", "张三", &reply)
    if err != nil {
        fmt.Println("Call:", err)
        return
    }
    fmt.Println(reply)
}
```

说明：

首先是通过 `rpc.Dial` 拨号 RPC 服务，然后通过 `client.Call` 调用具体的 RPC 方法。在调用 `client.Call` 时，第一个参数是用点号链接的 RPC 服务名字和方法名字，第二和第三个参数分别定义 RPC 方法的两个参数。

### 案例 2：综合使用

1. 创建一个 `hello` 微服务端，编写微服务端 RPC 代码，完成后启动该微服务端。
2. 创建一个 `hello` 客户端，编写客户端 RPC 代码，完成后启动该客户端，访问微服务端 RPC 功能，并返回相关数据。

#### 创建 `hello` 微服务端

新建 `mirco/server/hello/main.go` 文件，并编写代码：

```go
package main

import (
    "fmt"
    "net"
    "net/rpc"
)

// 定义类对象
type Hello struct{}

// 绑定类方法
func (this Hello) SayHello(req string, res *string) error {
    fmt.Println("请求的参数:", req)
    *res = "你好" + req
    return nil
}

func main() {
    // 1. 注册RPC服务
    err1 := rpc.RegisterName("hello", new(Hello))
    if err1 != nil {
        fmt.Println("注册 rpc 服务失败!", err1)
        return
    }

    // 2. 设置监听
    listener, err2 := net.Listen("tcp", "127.0.0.1:8800")
    if err2 != nil {
        fmt.Println("net.Listen err:", err2)
        return
    }
    defer listener.Close()
    fmt.Println("开始监听 ...")

    // 3. 建立链接
    for {
        //接收连接
        conn, err3 := listener.Accept()
        if err3 != nil {
            fmt.Println("Accept() err:", err3)
            return
        }
        // 4. 绑定服务
        go rpc.ServeConn(conn)
    }
}
```

注意：以上 `Hello` 结构体的方法必须满足 Go 语言的 RPC 规则：

1. 方法只能有两个可序列化的参数，其中第二个参数是指针类型，参数的类型不能是 channel（通道）、complex（复数类型）、func（函数），因为它们不能进行序列化。
2. 方法要返回一个 `error` 类型，同时必须是公开的方法。

#### 创建 `hello` 客户端

继续新建 `mirco/client/hello/main.go` 文件，并编写以下代码：

```go
package main

import (
    "fmt"
    "net/rpc"
)

func main() {
    // 1. 用 rpc 链接服务器 --Dial()
    conn, err := rpc.Dial("tcp", "127.0.0.1:8800")
    if err != nil {
        fmt.Println("Dial err:", err)
        return
    }
    defer conn.Close()

    // 2. 调用远程函数
    var reply string // 接受返回值 --- 传出参数
    err = conn.Call("hello.SayHello", "张三", &reply)
    if err != nil {
        fmt.Println("Call:", err)
        return
    }
    fmt.Println(reply)
}
```

#### 启动与运行

1. 启动微服务端：在终端中，导航到 `mirco/server/hello` 目录，然后执行以下命令：

    ```sh
    go run main.go
    ```

2. 启动客户端：在另一个终端中，导航到 `mirco/client/hello` 目录，然后执行以下命令：

    ```sh
    go run main.go
    ```

3. 客户端将会调用微服务端的 `SayHello` 方法，并输出类似如下的结果：

    ```
    你好张三
    ```

