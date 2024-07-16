package routers

import "github.com/gin-gonic/gin"

//设置api路由
func ApiRoutersInit(r *gin.Engine) {
	apiRouters := r.Group("/api")
	{
		apiRouters.GET("/", func(c *gin.Context) {
			c.String(200, "api首页")
		})
		apiRouters.GET("/userlist", func(c *gin.Context) {
			c.String(200, "一个userlist接口")
		})
	}
}
