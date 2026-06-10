package router

import (
	"weblog/ctrls"
	"weblog/middleware"

	"github.com/gin-gonic/gin"
)

func InitBlogRoutes(r *gin.Engine, ac *ctrls.Controllers) {
	// 建立前台根路由组
	blogGroup := r.Group("/api/v1")
	{
		// ================= 📝 文章前台接口 =================
		// 游客看文章列表（只看已发布）
		blogGroup.GET("/articles", ac.Article.GuestList)
		// 浏览单篇文章详情
		blogGroup.GET("/article/info", ac.Article.GetInfoById)

		// ================= ⚙️ 系统配置前台接口 =================
		// 前台获取单项配置（比如网站标题、ICP备案号）
		blogGroup.GET("/config/batch", ac.SystemConfig.GetBatchValues)

		// ================= 🔑 认证/登录接口 =================
		// 🌟 登录绝对不能加 JWT 拦截器，否则没人能登录成功！
		// 登录成功后，后端会吐出一个 Token 给前端保存
		blogGroup.POST("/user/login", ac.User.Login)
		blogGroup.POST("/user/register", ac.User.Register) // 注册新账号
	}
}

func InitAdminRoutes(r *gin.Engine, ac *ctrls.Controllers) {
	// 建立后台管理路由组，并注入 AuthJWT 中间件拦截器
	// 只有 Header 携带了合法 Token 的请求才能走进去
	adminGroup := r.Group("/api/v1/admin", middleware.JWTAuth())
	{
		// ================= ⚙️ 系统配置项管理 =================
		adminGroup.POST("/config", ac.SystemConfig.CreateConfig)      // 新增配置
		adminGroup.PUT("/config", ac.SystemConfig.UpdateConfigByKey)  // 修改配置
		adminGroup.DELETE("/config", ac.SystemConfig.DeleteConfig)    // 批量/单个删除配置
		adminGroup.GET("/config", ac.SystemConfig.GetConfigList)      // 获取配置分页列表
		adminGroup.GET("/config/info", ac.SystemConfig.GetConfigById) // 通过ID反显配置详情

		// ================= 📝 文章管理 =================
		adminGroup.POST("/article/create", ac.Article.Create)        // 发表/保存文章
		adminGroup.PUT("/article/update", ac.Article.Update)         // 编辑更新文章
		adminGroup.DELETE("/article/delete", ac.Article.BatchDelete) // 批量软删除文章
		adminGroup.GET("/articles", ac.Article.AdminList)

		// ================= 👤 个人中心/密码管理 =================

		adminGroup.GET("/user/info", ac.User.GetInfoById) // 管理员登录后台后，获取自己的基本信息（如：头像、昵称）
		adminGroup.GET("/user/list", ac.User.List)        // 管理员登录后台后，获取自己的基本信息（如：头像、昵称）
		adminGroup.PUT("/user/password", ac.User.UpdatePassword)
		adminGroup.POST("/user/logout", ac.User.Logout) // 退出登录 解析出 Token 里的用户 ID 并将其拉黑或让其失效
		adminGroup.DELETE("/user/delete", ac.User.DeleteAccount)
	}
}
