package models

type User struct { // 结构体首字母大写, 和数据库表名对应, 默认访问数据表users, 可以设置访问数据表的方法
	Id       int
	Username string
	Age      int
	Email    string
	AddTime  int // 大驼峰命名
}

// 配置数据库操作的表名称
func (User) TableName() string {
	return "user"
}
