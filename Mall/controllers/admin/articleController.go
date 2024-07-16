package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArticleController struct {
}

func (con ArticleController) Index(c *gin.Context) {
	c.String(http.StatusOK, "文章列表1")
}

func (con ArticleController) Add(c *gin.Context) {
	c.String(http.StatusOK, "文章添加2")
}

func (con ArticleController) Edit(c *gin.Context) {
	c.String(http.StatusOK, "文章编辑3")
}
