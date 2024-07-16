package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 基础控制器
type BaseController struct {
}

// 公共成功函数
func (con BaseController) Success(c *gin.Context, message string, redirectUrl string) {
	c.HTML(http.StatusOK, "admin/public/success.html", gin.H{
		"message":     message,
		"redirectUrl": redirectUrl,
	})
}

// 公共失败函数
func (con BaseController) Error(c *gin.Context, message string, redirectUrl string) {
	c.HTML(http.StatusOK, "admin/public/error.html", gin.H{
		"message":     message,
		"redirectUrl": redirectUrl,
	})
}
