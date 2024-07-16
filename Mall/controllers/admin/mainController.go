package admin

// 进入后台系统首页

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MainController struct {
}

func (con MainController) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/main/index.html", gin.H{
		"username": "超级管理员",
	})
}

func (con MainController) Welcome(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/main/welcome.html", gin.H{})
}
