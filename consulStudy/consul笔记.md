# Consul 简介

在服务发现的框架中，常见的有：

- Zookeeper
- Eureka
- etcd
- Consul

本文重点介绍 Consul。

## 什么是 Consul？

Consul 是一个提供服务发现功能的工具。它具有以下特点：

- **分布式**：Consul 可以在多个节点上运行，确保高可用性。
- **高可用**：即使部分节点出现故障，Consul 仍能正常工作。
- **横向扩展**：可以根据需要增加更多的节点来提升性能和可靠性。

Consul 提供的关键特性包括：

1. **Service Discovery**：
    - Consul 通过 DNS 或 HTTP 接口使服务注册和发现变得简单。
    - 外部服务（例如 SaaS 提供的服务）也可以注册到 Consul 中。

2. **Health Checking**：
    - 健康检测功能使 Consul 可以快速告警集群中的操作。
    - 与服务发现功能集成，防止将请求转发到故障服务。

3. **Key/Value Storage**：
    - 用于存储动态配置的系统。
    - 提供简单的 HTTP 接口，可以在任何地方操作。

4. **Multi-Datacenter**：
    - 支持任意数量的数据中心，配置简单。

本文将介绍服务发现、健康检查以及基本的 KV 存储。多数据中心的支持将在另一篇文章中详细讨论。

## Consul 的几个概念

Consul 集群由多个 Server 和 Client 组成。不论是 Server 还是 Client，都是 Consul 的节点，所有服务都可以注册到这些节点上，通过这些节点实现服务注册信息的共享。

### CLIENT

- **CLIENT** 表示 Consul 的客户端模式。
- 在这种模式下，所有注册到当前节点的服务会被转发到 Server，本身不持久化这些信息。

### SERVER

- **SERVER** 表示 Consul 的服务端模式。
- 与 Client 模式功能相同，但会将所有信息持久化到本地。
- 在故障时，信息可以被保留。

### SERVER-LEADER

- **SERVER-LEADER** 是 Server 中的领导者，负责同步注册信息到其他 Server，并负责各节点的健康监测。

### 其它信息

- Consul 节点之间通过通信协议和算法确保数据同步和实时性。
- 详细信息可以参考官方文档。

## Consul 角色

- **Client**：客户端，无状态，将 HTTP 和 DNS 接口请求转发给局域网内的服务端集群。
- **Server**：服务端，保存配置信息，高可用集群，与本地客户端通讯，通过广域网与其他数据中心通讯。推荐每个数据中心的 Server 数量为 3 或 5 个。

### Server 的选举

两个 Server 竞选 Leader，选举过程由 Raft 协议完成。

## Consul 命令示例

### Server 模式

```bash
consul agent -server -bootstrap-expect 2 -data-dir /tmp/consul -node=cn1 -bind=192.168.188.128 -ui -config-dir /etc/consul.d -rejoin -join 198.168.188.128 -client=0.0.0.0
```

- `-server`：表示当前节点为 Server 模式
- `-bootstrap-expect`：期望启动的 Server 数量
- `-data-dir`：数据存放目录
- `-node`：节点名称
- `-bind`：绑定 IP
- `-ui`：启动 Web 界面
- `-config-dir`：配置文件目录
- `-rejoin`：自动重连
- `-join`：加入哪个 Server 集群
- `-client`：监听的客户端地址

### Client 模式

```bash
consul agent -data-dir /data/consul0 -node=cn1 -bind=192.168.1.202 -config-dir /etc/consul.d -rejoin -join 198.168.188.128
```

- `-data-dir`：数据存放目录
- `-node`：节点名称
- `-bind`：绑定 IP
- `-config-dir`：配置文件目录
- `-rejoin`：自动重连
- `-join`：加入哪个 Server 集群

`consul members` 命令可以查看当前集群的成员信息。

`consul leave` 命令可以主动离开集群。

创建配置文件目录：

```bash
mkdir /etc/consul.d
```

编辑服务配置文件 `/etc/consul.d/app.json`：

```json
{
  "service": {
    "name": "app",
    "tags": ["rails"],
    "address": "127.0.0.1",
    "port": 8080,
    "checks": [
      {
        "name": "HTTP on port 8080",
        "http": "http://localhost:8080/health",
        "interval": "10s"
      }
    ]
  }
}
```

简单的健康检查示例：

```go
package main

import (
  "fmt"
  "net/http"
)

func main() {
  http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "OK")
  })
  http.ListenAndServe(":8080", nil)
}
```

使用 `dig` 命令查看服务：

```bash
dig @127.0.0.1 -p 8600 app.service.consul
```

## Consul 架构图

![consul架构图](https://ucc.alicdn.com/pic/developer-ecology/9b396dc233f0408cbb6f37b248fc2995.png?x-oss-process=image/resize,w_1400/format,webp)

## 参考资料

- [Consul 官方文档](https://www.consul.io/docs)
- [Raft 协议简介](https://raft.github.io/)

---