package routers

import (
	"example.com/mall/controller/admin"
	"example.com/mall/middlewares"
	"github.com/gin-gonic/gin"
)

// 设置admin后台路由
func AdminRoutersInit(r *gin.Engine) {
	//路由分组: 配置全局中间件:middlewares.InitMiddleware
	adminRouters := r.Group("/admin", middlewares.InitMiddleware)
	{
		//单文件上传
		adminRouters.POST("/user/doUpload", admin.UserController{}.DoUpload)

		//不同名字的多文件上传
		adminRouters.GET("/user/addnotfilenameupload", admin.UserController{}.AddNotFileNameUpload)
		//不同文件名的多文件上传
		adminRouters.POST("/user/doNotFileNameUploads", admin.UserController{}.DoNotFileNameUploads)

		//相同名字的多文件上传
		adminRouters.GET("/user/addcommonfilenameupload", admin.UserController{}.AddCommonFileNameUpload)
		//相同文件名的多文件上传
		adminRouters.POST("/user/doCommonFileNameUploads", admin.UserController{}.DoCommonFileNameUploads)
	}

	//添加:按日期单文件上传
	adminRouters.GET("/user/addByDate", admin.UserController{}.AddByDate)
	//添加:按日期单文件上传
	adminRouters.POST("/user/doUploadByDate", admin.UserController{}.DoUploadByDate)
}
