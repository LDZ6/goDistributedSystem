package main

import (
	"encoding/xml"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 时间戳转换成日期函数
func UnixToTime(timestamp int) string {
	t := time.Unix(int64(timestamp), 0)
	return t.Format("2006-01-02 15:04:05")
}

func Println(str1 string, str2 string) string {
	return str1 + str2
}

type UserInfo struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type Article struct {
	Title   string `json:"title" xml:"title"`
	Content string `json:"content" xml:"content"`
}

func main() {
	//初始化路由
	r := gin.Default()
	//自定义模板函数,必须在r.LoadHTMLGlob前面
	r.SetFuncMap(template.FuncMap{
		"UnixToTime": UnixToTime, //注册模板函数
		"Println":    Println,
	})
	//加载templates中所有模板文件, 使用不同目录下名称相同的模板,注意:一定要放在配置路由之前才得行
	r.LoadHTMLGlob("templates/**/*")
	//配置静态web目录 第一个参数表示路由,第二个参数表示映射的目录
	r.Static("/static", "./static")
	//配置路由
	//GET传值
	r.GET("/", func(c *gin.Context) {
		username := c.Query("username")
		age := c.Query("age")
		page := c.DefaultQuery("page", "1")
		c.JSON(http.StatusOK, gin.H{
			"username": username,
			"age":      age,
			"page":     page,
		})
	})
	//GET传值 获取文章id
	r.GET("/article", func(c *gin.Context) {
		id := c.DefaultQuery("id", "1")
		c.JSON(http.StatusOK, gin.H{
			"id":     id,
			"conent": "文章详情",
		})
	})

	//post演示
	r.GET("/user", func(c *gin.Context) {
		//渲染模板
		c.HTML(http.StatusOK, "default/user.html", gin.H{})
	})
	//获取表单post过来的数据
	r.POST("/doAddUser", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		age := c.DefaultPostForm("age", "20")
		c.JSON(http.StatusOK, gin.H{
			"username": username,
			"password": password,
			"age":      age,
		})
	})
	//获取GET POST传递的数据绑定到结构体
	//r.POST("/getUser", func(c *gin.Context) {
	r.GET("/getUser", func(c *gin.Context) {
		user := &UserInfo{}
		//绑定到对应的结构体
		if err := c.ShouldBind(&user); err == nil {
			c.JSON(http.StatusOK, user)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"err": err.Error(),
			})
		}
	})
	//获取 POST xml数据绑定到结构体
	r.POST("/xml", func(c *gin.Context) {
		article := &Article{}
		//获取c.Request.Body读取请求数据, 返回的是一个xml切片
		xmlSliceData, _ := c.GetRawData()
		//解析xml切片
		if err := xml.Unmarshal(xmlSliceData, &article); err == nil {
			c.JSON(http.StatusOK, article)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
		}
	})

	//动态路由传值
	//http://127.0.0.1:8080/list/2221  http://127.0.0.1:8080/list/2
	r.GET("/list/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	})
	r.Run() // 启动一个web服务
}
