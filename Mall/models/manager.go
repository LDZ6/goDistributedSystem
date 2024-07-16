package models

//管理员表

type Manager struct { // 结构体首字母大写, 和数据库表名对应, 默认访问数据表users, 可以设置访问数据表的方法
	Id       int
	Username string
	Password string
	Mobile   string
	Email    string
	Status   int
	RoleId   int
	AddTime  int
	IsSuper  int
	Role     Role `gorm:"foreignKey:RoleId;references:Id"` // 配置关联关系
}

//配置数据库操作的表名称
func (Manager) TableName() string {
	return "manager"
}
