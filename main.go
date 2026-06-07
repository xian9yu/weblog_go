package main

import (
	"log"
	"time"

	"weblog/models"
	"weblog/router"
	"weblog/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	router.User(r)
	router.Article(r)
	router.Category(r)
	router.Setting(r)
}

func initApp(r *gin.Engine) string {
	v := utils.Config() // 配置文件

	models.MySQL() // sql
	// models.InitRedis() // redis

	initRouter(r)

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

	return v.GetString("app.address")

}

func main() {
	r := gin.Default() // 初始化router
	addr := initApp(r)
	if err := r.Run(addr); err != nil {
		log.Fatalln("服务启动失败 ：", err)
	}
}
