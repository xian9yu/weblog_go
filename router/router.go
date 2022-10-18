package router

import (
	"github.com/gin-gonic/gin"
	"weblog/ctrls/article"
	"weblog/ctrls/category"
	"weblog/ctrls/setting"
	"weblog/ctrls/upload"
	"weblog/ctrls/user"
	"weblog/middleware"
)

func User(r *gin.Engine) {
	u := r.Group("/user/")
	// 公开接口
	{
		u.POST("/sign_in", user.SignIn)
		u.POST("/sign_up", user.SignUp)
	}

	u.Use(middleware.Auth()) // 鉴权
	{
		u.GET("/details/id/:id", user.GetDetailsById)
		u.GET("/details/mail/:mail", user.GetDetailsByMail)
		u.GET("/details/account/:account", user.GetDetailsByAccount)
		u.POST("/edit", user.Edit)
		u.POST("/delete", user.Delete)
		u.POST("/sign_out", user.SignOut)

	}
	// 仅管理员可操作
	u.Use(middleware.AdminAuth()) // 鉴权
	{
		u.GET("/list", user.List)
	}
}

func Article(r *gin.Engine) {
	a := r.Group("/article/")
	// 公开接口
	{
		a.GET("/details/id/:id", article.GetDetailsById)
		a.GET("/details/title/:title", article.GetDetailsByTitle)
		a.GET("/list", article.List)
		a.GET("/category_list", article.CategoryList)
	}

	a.Use(middleware.Auth()) // 鉴权
	{
		a.POST("/add", article.Add)
		a.POST("/edit", article.Edit)
		a.POST("/delete", article.Delete)
	}
}

func Category(r *gin.Engine) {
	c := r.Group("/category/")
	c.Use(middleware.Auth()) // 鉴权
	{
		c.GET("/list", category.List)
		c.GET("/details/id/:id", category.GetDetailsById)
		c.GET("/details/name/:name", category.GetDetailsByName)
	}

	// 仅管理员可操作
	c.Use(middleware.AdminAuth()) // 鉴权
	{
		c.POST("/add", category.Add)
		c.POST("/edit", category.Edit)
		c.POST("/delete", category.Delete)
	}
}

func Setting(r *gin.Engine) {
	s := r.Group("/setting/")

	// 仅管理员可操作
	s.Use(middleware.AdminAuth()) // 鉴权
	{
		s.POST("/add", setting.Add)
		s.POST("/edit", setting.Edit)
		s.POST("/delete", setting.Delete)
		s.GET("/list", setting.List)
		s.GET("/details/id/:id", setting.GetDetailsById)
		s.GET("/details/name/:name", setting.GetDetailsByName)
	}
}
func Upload(r *gin.Engine) {
	s := r.Group("/upload/")

	s.Use(middleware.Auth()) // 鉴权
	{
		s.POST("/image", upload.Image)
		s.POST("/delete", upload.Delete)

	}
}
