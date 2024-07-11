# Protocol Buffers (Protobuf) Message 嵌套详解

## Protobuf 基础语法

在编写 Protobuf 文件时，我们需要指定版本信息，包名，并定义消息类型（message）。下面是一个简单的示例：

```proto
syntax = "proto3"; // 指定版本信息，不指定会报错
package pb;        // 后期生成 Go 文件的包名

// message为关键字，作用为定义一种消息类型
message Person {
    string name = 1;  // 名字
    int32  age = 2 ;  // 年龄
    
    // 定义一个嵌套的message
    message PhoneNumber {
        string number = 1;
        int64 type = 2;
    }
    PhoneNumber phone = 3;
}
```

在上述示例中，`Person.PhoneNumber` 是一个嵌套的 message 类型，可以直接在 `Person` 中使用。

## 嵌套的 Message 示例

我们可以在 `Person` 中嵌套定义多个 `PhoneNumber`，并使用 `repeated` 关键字表示重复字段：

```proto
syntax = "proto3"; // 指定版本信息，不指定会报错
package pb;        // 后期生成 Go 文件的包名

// message为关键字，作用为定义一种消息类型
message Person {
    string name = 1;   // 名字
    int32  age = 2 ;   // 年龄
    
    // 定义一个message
    message PhoneNumber {
        string number = 1;
        int64 type = 2;
    }
    repeated PhoneNumber phone = 3;
}

// enum为关键字，作用为定义一种枚举类型
enum PhoneType {
    MOBILE = 0;
    HOME = 1;
    WORK = 2;
}
```

在这个例子中，`Person` 包含了一个重复的 `PhoneNumber` 字段，并且定义了一个枚举类型 `PhoneType`。

## 使用枚举类型的嵌套 Message 示例

我们可以将 `PhoneType` 枚举类型应用到嵌套的 `PhoneNumber` 中：

```proto
syntax = "proto3"; // 指定版本信息，不指定会报错
package pb;        // 后期生成 Go 文件的包名

// message为关键字，作用为定义一种消息类型
message Person {
    string name = 1;   // 名字
    int32  age = 2 ;   // 年龄
    
    // 定义一个message
    message PhoneNumber {
        string number = 1;
        PhoneType type = 2;
    }
    repeated PhoneNumber phone = 3;
}

// enum为关键字，作用为定义一种枚举类型
enum PhoneType {
    // 如果不设置将报错
    option allow_alias = true;
    MOBILE = 0;
    HOME = 1;
    WORK = 2;
    Personal = 2;
}
```

## Protobuf 编译

Protobuf 编译是通过编译器 `protoc` 进行的。通过这个编译器，我们可以把 `.proto` 文件生成 Go、Java、Python、C++、Ruby 或者 C# 代码。

### 编译命令

以下命令用于将 `.proto` 文件生成 Go 代码（以及 gRPC 代码）：

```sh
// 将当前目录中的所有 .proto 文件进行编译生成 Go 代码
protoc --go_out=./ --go_opt=paths=source_relative *.proto
```

Protobuf 编译器会把 `.proto` 文件编译成 `.pb.go` 文件。

### 参数说明

- **--go_out 参数**：指定 Go 代码生成的基本路径。编译器会将生成的 Go 代码输出到该路径。
  
  ```sh
  protoc --go_out=./ --go_opt=paths=source_relative *.proto
  ```

- **--go_opt 参数**：`protoc-gen-go` 提供的参数，用于指定生成文件的路径选项，可以设置多个。

  - `paths=import`：生成的文件会按 `go_package` 路径来生成，默认输出模式。
  - `paths=source_relative`：输出文件与输入文件放在相同的目录中。
  - `module=$PREFIX`：输出文件放在以 Go 包的导入路径命名的目录中，但从输出文件名中删除指定的目录前缀。

- **--proto_path 参数**：指定 `.proto` 文件所在的路径，如果忽略则默认当前目录。如果有多个目录可以多次调用 `--proto_path`，它们将会顺序地被访问并执行导入。

### 示例

```sh
protoc --proto_path=src --go_out=out --go_opt=paths=source_relative foo.proto bar/baz.proto
```

编译器将从 `src` 目录中读取输入文件 `foo.proto` 和 `bar/baz.proto`，并将输出文件 `foo.pb.go` 和 `bar/baz.pb.go` 写入 `out` 目录。如果需要，编译器会自动创建嵌套的输出子目录，但不会创建输出目录本身。