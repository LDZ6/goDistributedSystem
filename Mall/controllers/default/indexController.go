package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type IndexController struct {
}

func (con IndexController) Index(c *gin.Context) {
	//获取中间件中设置的username值,数据共享
	username, _ := c.Get("username")
	fmt.Println(username)
	//username是一个空接口类型,故要使用则需要用类型断言转换username
	v, ok := username.(string)
	if ok != true {
		c.String(200, "后台首页--获取用户名失败")
	} else {
		c.String(200, "后台首页,用户名:"+v)
	}

}
