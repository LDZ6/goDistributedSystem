## go mod 依赖管理
```golang
module gomod

go 1.16

require (
// 依赖包
)

exclude (
// 排除依赖包
)

replace (
// 替换依赖包
// 例如：
// replace example.com/foo => example.com/bar v1.2.3
)

retract (
// 撤销依赖包
// 例如：
// retract example.com/foo v1.2.3
)
```

---

***go mod 命令行***

```go mod init gomod```<br/>
```
go mod download
//仅仅下载指定的包到mod下,不安装相关依赖
```  
```
go mod tidy
//依赖对齐
```
```
go mod edit 
//自动编辑go.mod文件
go mod edit -require="包名@版本号"//添加require
go mod edit -replace="包名1@版本号=包名1@版本号"//添加require
go mod edit -exclude="包名1@版本号"//添加exclude
go mod edit -retract="v版本号"//添加retract
go mod edit -dropretract="v版本号"//删除retract
go mod vendor
go mod verify//确定mod是否被改变
go mod why 包名
```
---

```
go install //安装包的可执行文件
go get //-u 获取最新版本
go clean
```