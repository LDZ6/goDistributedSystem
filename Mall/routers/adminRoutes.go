package routers

import (
	"example.com/mall/controllers/admin"
	"example.com/mall/middlewares"
	"github.com/gin-gonic/gin"
)

// 设置admin后台路由
func AdminRoutersInit(r *gin.Engine) {
	//路由分组: 配置全局中间件:middlewares.InitMiddleware
	adminRouters := r.Group("/admin", middlewares.InitMiddleware)
	{
		//登录页面
		adminRouters.GET("/login", admin.LoginController{}.Index) // 实例化控制器,并访问其中方法
		adminRouters.POST("/doLogin", admin.LoginController{}.DoIndex)
		adminRouters.GET("/loginOut", admin.LoginController{}.LoginOut)

		//验证码
		adminRouters.GET("/captcha", admin.LoginController{}.Captcha)
	}
}
