package models

//连接数据库核心代码

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局使用DB,就需要把DB定义成公有的
var DB *gorm.DB
var err error

// 自动初始化数据库
func init() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:123456@tcp(127.0.0.1:3306)/gin?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
}
