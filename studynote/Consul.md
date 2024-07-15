# Consul集群详解及应用

## 简介
Consul 是 HashiCorp 公司推出的开源工具，用于实现分布式系统的服务发现与配置。Consul 是分布式的、高可用的、可横向扩展的。完成 Consul 的安装后，必须运行 agent 代理，agent 可以运行为 server 模式或 client 模式。

- **服务模式 (server 模式)**: 主要参与维护集群状态，响应 RPC 查询，与其他数据中心交换 WAN gossip ，以及向上级或远程数据中心转发查询，并且会将数据持久化。推荐使用 3 到 5 台机器。
- **客户模式 (client 模式)**: 客户模式下 Consul Agent 是一个非常轻量级的进程，它消耗最小的资源开销和少量的网络带宽，提供注册服务，运行健康检查，并将查询转发给服务器。客户端是相对无状态的，不负责数据的持久化，必须要有一个服务模式的 Consul。

一个数据中心至少必须拥有一台 server，建议在一个集群中有 3 或者 5 个 server。单一 server 在出现失败时，会不可避免地出现数据丢失。

## 架构
Consul 的架构分为上下两个部分，通过 WAN GOSSIP 进行报文交互。单个 datacenter 中，节点被划分为 server 和 client，他们之间通过 GRPC 进行通信。此外，Client 和 Server 之间通过 LAN Gossip 进行通信。

Consul 的 client 是一个非常轻量级的进程，用于注册服务、运行健康检查和转发对 server 的查询。每个数据中心至少必须拥有一个 server。Agent 必须在集群中的每个主机上运行。

## Server 功能
- 参与共识仲裁（Raft）
- 存储机器状态（日志存储）
- 处理查询
- 维护周边（LAN/WAN) 节点之间的关系

## Client 功能
- 负责通过该节点注册到 Consul 微服务的健康检查
- 将客户端的注册请求和查询转换为 server 的 RPC 请求
- 维护周边各节点（LAN/WAN) 的关系

## 服务端口
- **8300**: 只存在于 server 模式，选取 Leader 节点（Raft 协议），为 Leader 节点和 Client 节点提供 RPC 调用。
- **8301**: LAN 网中集群数据同步的通信端口（Gossip 协议），也是加入集群的通信端口。
- **8302**: 只存在于 server 模式，WAN 网中集群数据同步的通信端口（Gossip 协议），主要支持数据中心与数据中心之间通过 WAN 交互。
- **8500**: 提供 Http 服务（或 web 界面）。
- **8600**: 提供 DNS 服务端口。

## 实现原理
Consul 主要基于两个协议来实现其核心功能：
- **Gossip 协议**: 在集群内消息传递。
- **Raft 协议**: 保障日志的一致性。

## 案例讲解
以下是通过 Consul 实现微服务注册与发现的示例。

### 集群配置
准备四个虚拟机，三个作为服务端，一个作为客户端。安装好 Consul 等相关软件。假设 IP 地址如下：
- 服务端：192.168.1.129, 192.168.1.130, 192.168.1.131
- 客户端：192.168.1.132

### 步骤
1. **启动服务端**
    ```sh
    consul agent -server -bootstrap-expect 3 -node=server_01 -bind=192.168.1.129 -ui -data-dir=/root/usr/local/consul/data -client 0.0.0.0
    consul agent -server -bootstrap-expect 3 -node=server_02 -bind=192.168.1.130 -ui -data-dir=/root/usr/local/consul/data -client 0.0.0.0
    consul agent -server -bootstrap-expect 3 -node=server_03 -bind=192.168.1.131 -ui -data-dir=/root/usr/local/consul/data -client 0.0.0.0
    ```

2. **启动客户端**
    ```sh
    consul agent -data-dir=/root/usr/local/consul/data -node=client-01 -bind=192.168.1.132 -ui -client 0.0.0.0
    ```

3. **关联集群**
    在 server-02、server-03、client-01 节点上运行以下命令建立集群关系：
    ```sh
    consul join 192.168.1.129
    ```

4. **查看 Consul 成员和集群状态**
    ```sh
    consul members
    ```

### 服务注册与发现
通过修改示例代码，将服务注册到 Consul，并通过 Consul 客户端发现并调用服务。

- **服务端代码修改**:
    ```go
    package main

    import (
        "context"
        "fmt"
        "net"
        "google.golang.org/grpc"
        "github.com/hashicorp/consul/api"
        "serverHello/proto/helloService"
    )

    type Hello struct{}

    func (this Hello) SayHello(c context.Context, req *helloService.HelloReq) (*helloService.HelloRes, error) {
        fmt.Println(req)
        return &helloService.HelloRes{
            Message: "你好" + req.Name,
        }, nil
    }

    func main() {
        consulConfig := api.DefaultConfig()
        consulConfig.Address = "192.168.1.132:8500"
        consulClient, _ := api.NewClient(consulConfig)
        agentService := api.AgentServiceRegistration{
            ID:      "1",
            Tags:    []string{"test"},
            Name:    "HelloService",
            Port:    8082,
            Address: "192.168.1.111",
            Check: &api.AgentServiceCheck{
                TCP:      "192.168.1.111:8082",
                Timeout:  "5s",
                Interval: "30s",
            },
        }
        consulClient.Agent().ServiceRegister(&agentService)
        grpcServer := grpc.NewServer()
        helloService.RegisterHelloServer(grpcServer, new(Hello))
        listener, err := net.Listen("tcp", "192.168.1.111:8082")
        if err != nil {
            fmt.Println(err)
        }
        defer listener.Close()
        grpcServer.Serve(listener)
    }
    ```

- **客户端代码修改**:
    ```go
    package main

    import (
        "context"
        "fmt"
        "google.golang.org/grpc"
        "github.com/hashicorp/consul/api"
        "serverHello/proto/helloService"
    )

    func main() {
        consulConfig := api.DefaultConfig()
        consulClient, _ := api.NewClient(consulConfig)
        services, _ := consulClient.Agent().Services()
        helloServiceData := services["HelloService"]
        helloServiceAddress := helloServiceData.Address + ":" + fmt.Sprintf("%d", helloServiceData.Port)
        grpcConn, _ := grpc.Dial(helloServiceAddress, grpc.WithInsecure())
        defer grpcConn.Close()
        grpcClient := helloService.NewHelloClient(grpcConn)
        res, err := grpcClient.SayHello(context.Background(), &helloService.HelloReq{
            Name: "小名",
        })
        if err != nil {
            fmt.Println(err)
        }
        fmt.Println(res)
    }
    ```

运行客户端代码，控制台将打印出 "你好小名"。

### 集群退出
可以使用 `Ctrl-C` 优雅地关闭 Agent，或运行以下命令优雅退出集群：
```sh
consul leave
```
## GRPC 微服务集群概念
Consul 集群通过多台服务器的协作，当一台服务器宕机时，集群会从另一台服务器中获取服务，从而保证客户端访问 Consul 集群的负载均衡性。然而，如果终端对应的微服务宕机，Consul 集群服务器将无法访问该微服务。这时，我们需要引入 GRPC 微服务集群来解决这个问题。

### 什么是 GRPC 微服务集群

GRPC 微服务集群是指将一个 GRPC 微服务部署到多台不同服务器上。当其中一个微服务宕机时，Consul 会访问另一台服务器上的对应微服务，从而实现微服务的负载均衡。这样可以提高系统的可靠性和可用性。

GRPC 微服务集群主要实现的是微服务的负载均衡，通过将同样的微服务部署在不同的服务器上，结合 Consul 可以非常简单地搭建 GRPC 微服务集群：

- 同一个微服务的不同应用使用同样的注册名
- 同一个微服务的不同应用注册服务时使用不同的 ID

### 代码展示

#### 注销相关服务

首先，检查 Consul UI 上是否存在 `goods` 服务。如果存在，则先注销之前的 `goods` 服务。

#### 部署微服务集群

以 `goods` 微服务为例，将 `server/goods` 部署到两台服务器上，并修改 `main.go` 中的代码，以区分不同服务器上的同一个微服务。

##### 部署到服务器 192.168.1.111

```go
package main

import (
    "context"
    "fmt"
    "github.com/hashicorp/consul/api"
    "goods/proto/goodsService"
    "net"
    "google.golang.org/grpc"
    "strconv"
)

type Goods struct{}

func (g Goods) AddGoods(c context.Context, req *goodsService.AddGoodsReq) (*goodsService.AddGoodsRes, error) {
    fmt.Println(req)
    return &goodsService.AddGoodsRes{
        Message: "第一个goods微服务-增加数据成功",
        Success: true,
    }, nil
}

func (g Goods) GetGoods(c context.Context, req *goodsService.GetGoodsReq) (*goodsService.GetGoodsRes, error) {
    var tempList []*goodsService.GoodsModel
    for i := 0; i < 10; i++ {
        tempList = append(tempList, &goodsService.GoodsModel{
            Title:   "商品" + strconv.Itoa(i),
            Price:   float64(i),
            Content: "测试商品内容" + strconv.Itoa(i),
        })
    }
    return &goodsService.GetGoodsRes{
        GoodsList: tempList,
    }, nil
}

func main() {
    consulConfig := api.DefaultConfig()
    consulConfig.Address = "192.168.1.132:8500"
    consulClient, _ := api.NewClient(consulConfig)
    agentService := api.AgentServiceRegistration{
        ID:      "1",
        Tags:    []string{"goods"},
        Name:    "GoodsService",
        Port:    8080,
        Address: "192.168.1.111",
        Check: &api.AgentServiceCheck{
            TCP:      "192.168.1.111:8080",
            Timeout:  "5s",
            Interval: "30s",
        },
    }
    consulClient.Agent().ServiceRegister(&agentService)

    grpcServer := grpc.NewServer()
    goodsService.RegisterGoodsServer(grpcServer, new(Goods))
    listener, err := net.Listen("tcp", "192.168.1.111:8080")
    if err != nil {
        fmt.Println(err)
    }
    defer listener.Close()
    grpcServer.Serve(listener)
}
```

##### 部署到服务器 192.168.1.112

```go
package main

import (
    "context"
    "fmt"
    "github.com/hashicorp/consul/api"
    "goods/proto/goodsService"
    "net"
    "google.golang.org/grpc"
    "strconv"
)

type Goods struct{}

func (g Goods) AddGoods(c context.Context, req *goodsService.AddGoodsReq) (*goodsService.AddGoodsRes, error) {
    fmt.Println(req)
    return &goodsService.AddGoodsRes{
        Message: "第二个goods微服务-增加数据成功",
        Success: true,
    }, nil
}

func (g Goods) GetGoods(c context.Context, req *goodsService.GetGoodsReq) (*goodsService.GetGoodsRes, error) {
    var tempList []*goodsService.GoodsModel
    for i := 0; i < 10; i++ {
        tempList = append(tempList, &goodsService.GoodsModel{
            Title:   "商品" + strconv.Itoa(i),
            Price:   float64(i),
            Content: "测试商品内容" + strconv.Itoa(i),
        })
    }
    return &goodsService.GetGoodsRes{
        GoodsList: tempList,
    }, nil
}

func main() {
    consulConfig := api.DefaultConfig()
    consulConfig.Address = "192.168.1.132:8500"
    consulClient, _ := api.NewClient(consulConfig)
    agentService := api.AgentServiceRegistration{
        ID:      "2",
        Tags:    []string{"goods"},
        Name:    "GoodsService",
        Port:    8080,
        Address: "192.168.1.112",
        Check: &api.AgentServiceCheck{
            TCP:      "192.168.1.112:8080",
            Timeout:  "5s",
            Interval: "30s",
        },
    }
    consulClient.Agent().ServiceRegister(&agentService)

    grpcServer := grpc.NewServer()
    goodsService.RegisterGoodsServer(grpcServer, new(Goods))
    listener, err := net.Listen("tcp", "192.168.1.112:8080")
    if err != nil {
        fmt.Println(err)
    }
    defer listener.Close()
    grpcServer.Serve(listener)
}
```

注意，两台服务器上的代码不同之处在于：

- ID 不同
- Port 为对应服务器上的端口号
- Address 为对应服务器的 IP
- Name 和 Tag 一定要一致

启动服务：

```bash
go run main.go
```

这样就注册 `goods` 微服务到 Consul 服务发现集群中了。

### 客户端进行调度

客户端请求微服务时，需要实现负载均衡。这里有两种方式：

#### 方式一：随机取地址

```go
package main

import (
    "client/proto/goodsService"
    "context"
    "fmt"
    "github.com/hashicorp/consul/api"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "strconv"
)

func main() {
    consulConfig := api.DefaultConfig()
    consulConfig.Address = "192.168.1.132:8500"
    consulClient, _ := api.NewClient(consulConfig)
    serviceGoodsEntry, _, _ := consulClient.Health().Service("GoodsService", "test", false, nil)
    addressGoods := serviceGoodsEntry[0].Service.Address + ":" + strconv.Itoa(serviceGoodsEntry[0].Service.Port)

    grpcGoodsClient, err := grpc.Dial(addressGoods, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        fmt.Println(err)
    }
    clientGoods := goodsService.NewGoodsClient(grpcGoodsClient)
    res1, _ := clientGoods.AddGoods(context.Background(), &goodsService.AddGoodsReq{
        Goods: &goodsService.GoodsModel{
            Title:   "测试商品",
            Price:   20,
            Content: "测试商品的内容",
        },
    })
    fmt.Println(res1.Message)
    fmt.Println(res1.Success)

    res2, _ := clientGoods.GetGoods(context.Background(), &goodsService.GetGoodsReq{})
    fmt.Printf("%#v", res2.GoodsList)
    for i := 0; i < len(res2.GoodsList); i++ {
        fmt.Printf("%#v\r\n", res2.GoodsList[i])
    }
}
```

#### 方式二：使用 grpc-consul-resolver 实现域名解析

```go
package main

import (
    "client/proto/goodsService"
    "context"
    "fmt"
    _ "github.com/mbobakov/grpc-consul-resolver"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    grpcGoodsClient, err := grpc.Dial(
        "consul://192.168.234.132:8500/GoodsService",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
    )
    if err != nil {
        fmt.Println(err)
    }
    clientGoods := goodsService.NewGoodsClient(grpcGoodsClient)
    res1, _ := clientGoods.AddGoods(context.Background(), &goodsService.AddGoodsReq{
        Goods: &goodsService.GoodsModel{
            Title:   "测试商品",
            Price:   20,
            Content: "测试商品的内容",
        },
    })
    fmt.Println(res1.Message)
    fmt.Println(res1.Success)

    res2, _ := clientGoods.GetGoods(context.Background(), &goods

Service.GetGoodsReq{})
    fmt.Printf("%#v", res2.GoodsList)
    for i := 0; i < len(res2.GoodsList); i++ {
        fmt.Printf("%#v\r\n", res2.GoodsList[i])
    }
}
```

注意，使用 `grpc-consul-resolver` 可以自动处理负载均衡，不需要手动配置地址。