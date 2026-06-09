package router

import (
	"weblog/ctrls/article"
	"weblog/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Article(db *gorm.DB, articleCtrl *article.Controller, r *gin.RouterGroup) {
	// 🌲 1. 前台公开路由组（无需登录）
	guestGroup := r.Group("/article")
	{
		guestGroup.GET("/details/:id", articleCtrl.GetDetailsById) // 简化：/article/details/10 直接拿 :id
		guestGroup.GET("/list", articleCtrl.GuestList)             // 查全量、查分类一箭双雕：/article/list
	}
	// 🔒 2. 后台管理路由组（必须鉴权）
	adminGroup := r.Group("/admin/article").Use(middleware.JWTAuth())
	{
		adminGroup.POST("/create", articleCtrl.Create)
		adminGroup.POST("/edit", articleCtrl.Edit)
		adminGroup.POST("/delete", articleCtrl.BatchDelete)
		adminGroup.GET("/list", articleCtrl.AdminList)
	}
}

func User(db *gorm.DB, articleCtrl *article.Controller, r *gin.RouterGroup) {
	//u := r.Group("/user/")
	//// 公开接口
	//{
	//	u.POST("/sign_in", user.SignIn)
	//	u.POST("/sign_up", user.SignUp)
	//}
	//
	////u.Use(middleware.JWTAuth()) // 鉴权
	//{
	//	u.GET("/details/id/:id", user.GetDetailsById)
	//	u.GET("/details/email/:email", user.GetDetailsByEmail)
	//	//u.GET("/details/account/:account", user.GetDetailsByAccount)
	//	u.POST("/edit", user.Edit)
	//	u.POST("/delete", user.Delete)
	//	u.POST("/sign_out", user.SignOut)
	//
	//}
	//// 仅管理员可操作
	////u.Use(middleware.AdminAuth()) // 鉴权
	//{
	//	u.GET("/list", user.List)
	//}
}

func SystemConfig(db *gorm.DB, articleCtrl *article.Controller, r *gin.RouterGroup) {
	//s := r.Group("/systemConfig/")

	// 仅管理员可操作
	//s.Use(middleware.AdminAuth()) // 鉴权
	//{
	//	s.POST("/add", systemConfig.Add)
	//	s.POST("/edit", systemConfig.Edit)
	//	s.POST("/delete", systemConfig.Delete)
	//	s.GET("/list", systemConfig.List)
	//	s.GET("/details/id/:id", systemConfig.GetDetailsById)
	//	s.GET("/details/name/:name", systemConfig.GetDetailsByName)
	//}
}
