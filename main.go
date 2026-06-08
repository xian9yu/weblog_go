package main

import (
	"log"
	"time"
	"weblog/ctrls"
	"weblog/database"
	"weblog/models"
	"weblog/router"
	"weblog/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initRouter(db *gorm.DB, r *gin.Engine) {
	// 初始化 Repository (数据层)，把 database 注入进去
	repos := models.NewRepositories(db)
	// 初始化 Controller (控制层)，把 repo 注入进去
	ct := ctrls.NewControllers(repos)
	api := r.Group("/api/v1")
	{
		router.Article(db, ct.Article, api)
		// 以后加功能只需要在这继续追加：
		// api.POST("/articles", articleCtrl.Create)
		// api.POST("/comments", commentCtrl.Create)
	}

}

func initApp() (*gin.Engine, string) {
	// 启动时一次性加载配置
	conf, err := utils.InitConfig("config.yaml")
	if err != nil {
		log.Fatalf("❌ 项目启动失败, 配置加载异常: %v\n", err)
	}
	r := gin.Default() // 初始化router

	// 初始化核心组件
	db := database.InitMySQL(conf.Database.Mysql)
	// 载入路由系统，并把 database 注入进去
	initRouter(db, r)

	// 限制上传文件大小
	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	// pprof.Register(r)
	// 使用gin自带的异常恢复中间件，避免出现异常时程序退出
	r.Use(gin.Recovery())

	// 添加 CORS 中间件
	r.Use(cors.New(cors.Config{
		// AllowAllOrigins:  true,
		AllowOrigins:     []string{"http://127.0.0.1:8000"}, // 前端地址
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 模式为 release
	gin.SetMode(gin.ReleaseMode)

	return r, conf.App.Address

}

func main() {
	r, addr := initApp()
	if err := r.Run(addr); err != nil {
		log.Fatalln("服务启动失败 ：", err)
	}
}
