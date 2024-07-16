package admin

import (
	"net/http"
	"os"
	"path"
	"strconv"

	"example.com/mall/models"
	"github.com/gin-gonic/gin"
)

// 定义一个UserController结构体,可以实例化结构体访问里面的方法
type UserController struct {
	BaseController // 继承基础控制器
}

func (con UserController) Index(c *gin.Context) {
	con.Success(c, "登录成功", "/admin")
}

func (con UserController) Add(c *gin.Context) {
	c.String(200, "用户添加2")
}

func (con UserController) Edit(c *gin.Context) {
	c.String(200, "用户编辑3")
}

// 添加:按日期单文件上传
func (con UserController) AddByDate(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/user/add-by-date.html", gin.H{
		"title": "添加:按日期单文件上传",
	})
}

//按日期单文件上传
/**
1.获取上传文件
2.获取后缀名,判断后缀是否正确: .jpg,.png,.gif,.jpeg
3.创建图片保存目录 ./static/upload/20230203
4.生成文件名称和文件保存目录
5.执行上传
*/
func (con UserController) DoUploadByDate(c *gin.Context) {
	//获取表单中提交的username
	username := c.PostForm("username")
	//1.获取上传文件
	file, err := c.FormFile("face")
	//判断上传文件上否存在
	if err != nil { //说明上传文件不存在
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": "fail",
			"message": err.Error(),
		})
		return
	}
	//2.获取后缀名,判断后缀是否正确: .jpg,.png,.gif,.jpeg
	extName := path.Ext(file.Filename)
	//设置后缀map
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
	}
	//判断后缀是否合法
	if _, ok := allowExtMap[extName]; !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": "fail",
			"message": "上传文件后缀不合法",
		})
		return
	}
	//3.创建图片保存目录 ./static/upload/20230203
	//获取日期
	day := models.GetDay()
	//拼接目录
	dir := "./static/upload/" + day
	//创建目录:MkdirAll 目录不存在,会一次性创建多层
	err = os.MkdirAll(dir, 0666)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": "fail",
			"message": "创建目录失败",
		})
		return
	}
	//4.生成文件名称和文件保存目录: models.GetUnix() 获取时间戳(int64); strconv.FormatInt() 把时间戳(int64)转换成字符串
	filename := strconv.FormatInt(models.GetUnix(), 10) + extName
	//5.执行上传
	dst := path.Join(dir, filename)
	//上传文件到指定目录
	c.SaveUploadedFile(file, dst)
	c.JSON(http.StatusOK, gin.H{
		"success":  "true",
		"username": username,
		"dst":      dst,
	})
}

// 添加:不同名字的多文件上传
func (con UserController) AddNotFileNameUpload(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/user/addnotfilenameupload.html", gin.H{
		"title": "添加不同名字的多文件上传",
	})
}

// 不同文件名的多文件上传
func (con UserController) DoNotFileNameUploads(c *gin.Context) {
	//获取表单中提交的username
	username := c.PostForm("username")
	//获取文件
	file1, err1 := c.FormFile("face1")
	file2, err2 := c.FormFile("face2")
	//判断上传文件上否存在
	if err1 != nil { //说明上传文件不存在
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": "fail",
			"message": err1.Error(),
		})
		return
	}
	//判断上传文件上否存在
	if err2 != nil { //说明上传文件不存在
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": "fail",
			"message": err2.Error(),
		})
		return
	}

	//设置需要上传的文件目录 file.Filename 获取文件名; "./static/upload" 是基于main.go文件路由的
	dst1 := path.Join("./static/upload", file1.Filename)
	//上传文件到指定目录
	c.SaveUploadedFile(file1, dst1)

	//设置需要上传的文件目录 file.Filename 获取文件名; "./static/upload" 是基于main.go文件路由的
	dst2 := path.Join("./static/upload", file2.Filename)
	//上传文件到指定目录
	c.SaveUploadedFile(file2, dst2)

	c.JSON(http.StatusOK, gin.H{
		"success":  "true",
		"username": username,
		"dst1":     dst1,
		"dst2":     dst2,
	})
}

// 单文件上传
func (con UserController) DoUpload(c *gin.Context) {
	//获取表单中提交的username
	username := c.PostForm("username")
	//获取文件
	file, err := c.FormFile("face")
	//判断上传文件上否存在
	if err != nil { //说明上传文件不存在
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": "fail",
			"message": err.Error(),
		})
		return
	}
	//设置需要上传的文件目录 file.Filename 获取文件名; "./static/upload" 是基于main.go文件路由的
	dst := path.Join("./static/upload", file.Filename)
	//上传文件到指定目录
	c.SaveUploadedFile(file, dst)
	c.JSON(http.StatusOK, gin.H{
		"success":  "true",
		"username": username,
		"dst":      dst,
	})
}

func (con UserController) DoCommonFileNameUploads(c *gin.Context) {
	//获取表单中提交的username
	username := c.PostForm("username")
	//获取form
	form, _ := c.MultipartForm()
	//获取多文件
	files := form.File["face[]"]
	//遍历文件,并上传
	for _, file := range files {
		//设置需要上传的文件目录 file.Filename 获取文件名; "./static/upload" 是基于main.go文件路由的
		dst := path.Join("./static/upload", file.Filename)
		//上传文件到指定目录
		c.SaveUploadedFile(file, dst)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  "true",
		"username": username,
		"message":  "文件上传成功",
	})
}

func (con UserController) AddCommonFileNameUpload(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/user/addcommonfilenameupload.html", gin.H{
		"title": "相同名字的多文件上传",
	})
}
