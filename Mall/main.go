package main

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Article struct {
	Title   string `json:"title"`
	Content string `json:"
"`
}

// 时间戳转换成日期函数
func UnixToTime(timestamp int) string {
	t := time.Unix(int64(timestamp), 0)
	return t.Format("2006-01-02 15:04:05")
}

func Println(str1 string, str2 string) string {
	return str1 + str2
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
	//前台路由
	r.GET("/", func(c *gin.Context) {
		//渲染模板文件
		c.HTML(http.StatusOK, "default/index.html", gin.H{
			"title": "首页",
			"score": 88,
			"hobby": []string{"吃饭", "睡觉", "打豆豆"}, // 切片
			"newList": []interface{}{ // 接口
				&Article{
					Title:   "新闻标题1",
					Content: "新闻内容1",
				},
				&Article{
					Title:   "新闻标题2",
					Content: "新闻内容2",
				},
			},
			"testSlice": []string{}, // 空数组/空切片
			"news": &Article{ // 结构体
				Title:   "新闻标题3",
				Content: "新闻内容3",
			},
			"date": 1672648334,
		})
	})
	r.GET("/news", func(c *gin.Context) {
		news := &Article{
			Title:   "新闻标题",
			Content: "新闻内容",
		}
		c.HTML(http.StatusOK, "default/news.html", gin.H{
			"title": "新闻详情",
			"news":  news,
		})
	})

	//后台路由
	r.GET("/admin", func(c *gin.Context) {
		//渲染模板文件
		c.HTML(http.StatusOK, "admin/index.html", gin.H{
			"title": "后台首页",
		})
	})
	r.GET("/admin/news", func(c *gin.Context) {
		news := &Article{
			Title:   "后台新闻标题",
			Content: "后台新闻内容",
		}
		c.HTML(http.StatusOK, "admin/news.html", gin.H{
			"title": "后台新闻详情",
			"news":  news,
		})
	})
	r.Run() // 启动一个web服务
}
