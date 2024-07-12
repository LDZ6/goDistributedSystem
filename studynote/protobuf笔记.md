## Protocol Buffers (protobuf) 使用笔记

### 参考文档
[Protocol Buffers Documentation](https://developers.google.com/protocol-buffers/docs/proto3)

### 基本语法

#### 案例引入

创建一个后缀为 `.proto` 的文件，如下：

```proto
syntax = "proto3";  // 指定版本信息，不指定会报错，默认是proto2

option go_package = "./proto;helloworld"; // 分号前表示路径，后面表示包名

message Person {
    string name = 1; // 名字
    int32 age = 2; // 年龄
    repeated string hobby = 3; // 爱好，类似于go中的切片
}
```

#### 总结

- **版本声明**：`syntax = "proto3";` 是必须的。
- **包声明**：`option go_package` 指定生成的 Go 包路径和名称。
- **消息类型**：使用 `message` 定义类似 Go 语言结构体的数据结构。
- **字段编号**：字段编号用于标识字段，在同一消息中不能重复。

#### 消息的格式说明

消息由至少一个字段组成，格式如下：

```proto
// 注释格式
（字段修饰符）数据类型 字段名称 = 唯一的编号标签值;
```

- **字段修饰符**：`singular`、`repeated`
- **数据类型**：如 `string`、`int32` 等
- **编号标签**：用于唯一标识字段

### proto2 与 proto3 区别

- **字段规则**：proto3 移除了 `required`，`optional` 改为 `singular`。
- **默认值**：proto3 字段的默认值由系统决定，不再提供 `default` 选项。
- **零值**：proto3 必须有一个零值，以便使用 0 作为数字默认值。
- **扩展支持**：proto3 移除了扩展，新增了 `Any` 和 JSON 映射。

### 高级用法

#### message 嵌套

```proto
syntax = "proto3";
option go_package = "./proto;helloworld";

message Person {
    string name = 1;
    int32 age = 2;

    message PhoneNumber {
        string number = 1;
        int64 type = 2;
    }
    PhoneNumber phone = 3;
}
```

#### repeated 关键字

```proto
syntax = "proto3";
option go_package = "./proto;helloworld";

message Person {
    string name = 1;
    int32 age = 2;

    message PhoneNumber {
        string number = 1;
        int64 type = 2;
    }
    repeated PhoneNumber phone = 3;
}
```

#### 默认值

- **string**：空字符串
- **bytes**：空 bytes
- **bool**：false
- **numeric**：0
- **enums**：第一个定义的枚举值
- **repeated**：空列表
- **message**：空对象

#### enum 关键字

```proto
syntax = "proto3";
package pb;

message Person {
    string name = 1;
    int32 age = 2;

    message PhoneNumber {
        string number = 1;
        PhoneType type = 2;
    }
    repeated PhoneNumber phone = 3;

    enum Corpus {
        UNIVERSAL = 0;
        WEB = 1;
        IMAGES = 2;
        LOCAL = 3;
        NEWS = 4;
        PRODUCTS = 5;
        VIDEO = 6;
    }
    Corpus corpus = 4;
}

enum PhoneType {
    MOBILE = 0;
    HOME = 1;
    WORK = 2;
}
```

### 定义 RPC 服务

```proto
syntax = "proto3";
option go_package = "./sayService";

service sayService {
    rpc SayHello (HelloRequest) returns (HelloRes);
}

message HelloRequest {
    string name = 1;
}

message HelloRes {
    string message = 1;
}
```

生成 Go 文件命令：

```sh
protoc --go_out=plugins=grpc:. *.proto
```

### 案例总结

```proto
syntax = "proto3";
package tutorial;

message Student {
    uint64 id = 1;
    string name = 2;
    string email = 3;

    enum PhoneType {
        MOBILE = 0;
        HOME = 1;
    }

    message PhoneNumber {
        string number = 1;
        PhoneType type = 2;
    }
    repeated PhoneNumber phone = 4;
}
```

### protobuf 编译

生成命令如下：

```sh
protoc --proto_path=IMPORT_PATH --go_out=DST_DIR path/to/file.proto
```

- **`--proto_path`**：指定 .proto 文件路径
- **`--go_out`**：指定生成的 Go 文件目录

一般使用：

```sh
protoc --go_out=./ *.proto
```

### protobuf 序列化与反序列化

#### 定义 proto 文件

```proto
syntax = "proto3";
option go_package = "./protoService";

message Userinfo {
    string name = 1;
    int32 age = 2;
    repeated string hobby = 3;
    PhoneType phone = 4;
}

enum PhoneType {
    MOBILE = 0;
    HOME = 1;
    WORK = 2;
}
```

#### 生成 .go 文件

```sh
protoc --go_out=./ *.proto  // 一般情况下使用这个命令
protoc --go_out=plugins=grpc:. *.proto  // 有 RPC 服务的情况下使用这个命令
```

#### 序列化与反序列化示例

```go
package main

import (
    "fmt"
    "go_code/micro/protoc/protoService"
    "google.golang.org/protobuf/proto"
)

func main() {
    // 初始化并赋值
    u := &protoService.Userinfo{
        Name: "zhangsan",
        Age: 20,
        Hobby: []string{"吃饭", "睡觉", "写代码"},
    }
    fmt.Println(u.GetHobby())
    
    // 序列化
    data, err := proto.Marshal(u)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(data)
    
    // 反序列化
    info := protoService.Userinfo{}
    err = proto.Unmarshal(data, &info)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Printf("%#v", info)
    fmt.Println(info.GetHobby())
}
```