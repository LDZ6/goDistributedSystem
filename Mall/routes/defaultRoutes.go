package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 设置前台后台路由
func DefaultRoutersInit(r *gin.Engine) {
	defaultRouters := r.Group("/")
	{
		defaultRouters.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "default/index.html", gin.H{
				"msg": "前台首页",
			})
		})
		defaultRouters.GET("/news", func(c *gin.Context) {
			c.String(200, "新闻")
		})
	}
}
