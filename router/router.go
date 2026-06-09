package router

import (
	"weblog/ctrls"
	"weblog/ctrls/article"
	"weblog/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Article(db *gorm.DB, articleCtrl *article.Controller, r *gin.RouterGroup) {
	// 🌲 1. 前台公开路由组（无需登录）
	guestGroup := r.Group("/article")
	{
		guestGroup.GET("/details/:id", articleCtrl.GetInfoById) // 简化：/article/details/10 直接拿 :id
		guestGroup.GET("/list", articleCtrl.GuestList)          // 查全量、查分类一箭双雕：/article/list
	}
	// 🔒 2. 后台管理路由组（必须鉴权）
	adminGroup := r.Group("/admin/article").Use(middleware.JWTAuth())
	{
		adminGroup.POST("/create", articleCtrl.Create)
		adminGroup.POST("/edit", articleCtrl.Update)
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
	//	u.GET("/details/id/:id", user.GetInfoById)
	//	u.GET("/details/email/:email", user.GetInfoByEmail)
	//	//u.GET("/details/account/:account", user.GetInfoByAccount)
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
	//	s.GET("/details/id/:id", systemConfig.GetInfoById)
	//	s.GET("/details/name/:name", systemConfig.GetInfoByName)
	//}
}

func InitBlogRoutes(r *gin.Engine, ac *ctrls.Controllers) {
	// 建立前台根路由组
	blogGroup := r.Group("/api/v1")
	{
		// ================= 📝 文章前台接口 =================
		// 游客看文章列表（只看已发布）
		blogGroup.GET("/articles", ac.Article.GuestList)
		// 浏览单篇文章详情（通过 URL 参数 :id 传参）
		blogGroup.GET("/article/detail/:id", ac.Article.GetInfoById)

		// ================= ⚙️ 系统配置前台接口 =================
		// 前台获取单项配置（比如网站标题、ICP备案号）
		blogGroup.GET("/config/value", ac.SystemConfig.GetConfigByKey)

		// ================= 🔑 认证/登录接口 =================
		// 🌟 登录绝对不能加 JWT 拦截器，否则没人能登录成功！
		// 登录成功后，后端会吐出一个 Token 给前端保存
		blogGroup.POST("/login", ac.User.SignIn)
	}
}

func InitAdminRoutes(r *gin.Engine, ac *ctrls.Controllers) {
	// 建立后台管理路由组，并注入 AuthJWT 中间件拦截器
	// 只有 Header 携带了合法 Token 的请求才能走进去
	adminGroup := r.Group("/api/v1/admin", middleware.JWTAuth())
	{
		// ================= ⚙️ 系统配置项管理 =================
		adminGroup.POST("/config", ac.SystemConfig.CreateConfig)            // 新增配置
		adminGroup.PUT("/config", ac.SystemConfig.UpdateConfigByKey)        // 修改配置
		adminGroup.DELETE("/config", ac.SystemConfig.DeleteConfig)          // 批量/单个删除配置
		adminGroup.GET("/config", ac.SystemConfig.GetConfigList)            // 获取配置分页列表
		adminGroup.GET("/config/detail/:id", ac.SystemConfig.GetConfigById) // 通过ID反显配置详情

		// ================= 📝 文章管理 =================
		adminGroup.POST("/article", ac.Article.Create)        // 发表/保存文章
		adminGroup.PUT("/article", ac.Article.Update)         // 编辑更新文章
		adminGroup.DELETE("/article", ac.Article.BatchDelete) // 批量软删除文章
		adminGroup.GET("/articles", ac.Article.AdminList)

		// ================= 👤 个人中心/密码管理 =================
		// 管理员登录后台后，获取自己的基本信息（如：头像、昵称）
		adminGroup.GET("/user/info", ac.User.GetInfoById)
		// 定期修改后台登录密码
		adminGroup.PUT("/user/password", ac.User.UpdatePassword)
	}
}
