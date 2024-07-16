package admin

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/mall/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginController struct {
	BaseController
}

// 进入登录页面
func (con LoginController) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/login/login.html", gin.H{})
}

// 执行登录操作
func (con LoginController) DoIndex(c *gin.Context) {
	//获取表单中的数据
	captchaId := c.PostForm("captchaId")     // 验证码id
	verifyValue := c.PostForm("verifyValue") //验证码的值
	//获取用户名以及密码
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 1.判断验证码是否验证成功
	if flag := models.VerifyCaptcha(captchaId, verifyValue); flag {
		//2.查询数据库,判断用户以及密码是否正确
		userinfo := []models.Manager{}
		password = models.Md5(password)
		models.DB.Where("username = ? and password = ? ", username, password).Find(&userinfo)
		if len(userinfo) > 0 {
			//3.执行登录,保存用户信息,执行跳转操作
			session := sessions.Default(c)
			//注意: session.Set没法保存结构体对应的切片,所以需要把结构体转换成json字符串
			userinfoSlice, _ := json.Marshal(userinfo)
			session.Set("userinfo", string(userinfoSlice))
			session.Save()
			con.Success(c, "登录成功", "/admin")
		} else {
			con.Error(c, "用户名或密码错误", "/admin/login")
		}
	} else {
		con.Error(c, "验证码验证失败", "/admin/login")
	}
}

// 获取验证码,验证验证码
func (con LoginController) Captcha(c *gin.Context) {
	id, b64s, err := models.MakeCaptcha()
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"captchaId":    id,
		"captchaImage": b64s,
	})
}

func (con LoginController) LoginOut(c *gin.Context) {
	//1.销毁session中用户信息
	session := sessions.Default(c)
	session.Delete("userinfo")
	session.Save()
	con.Success(c, "退出登录成功", "/admin/login")
}
