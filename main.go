package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"weblog/middleware"
	"weblog/models"
	"weblog/router"
	"weblog/utils"
)

func initRouter(r *gin.Engine) {
	router.User(r)
	router.Article(r)
	router.Category(r)
	router.Setting(r)
	router.Upload(r)

}

func main() {
	v := utils.InitConfig() // 初始化 配置文件

	models.InitSQL() // 初始化 sql
	//models.InitRedis() // 初始化 redis

	r := gin.Default() // 初始化router
	initRouter(r)

	// 限制上传文件大小
	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	// pprof.Register(r)
	// 使用gin自带的异常恢复中间件，避免出现异常时程序退出
	r.Use(gin.Recovery())
	// 跨域中间件
	r.Use(middleware.Cors())

	// 模式为 release
	gin.SetMode(gin.ReleaseMode)

	addr := v.GetString("app.addr")
	port := v.GetString("app.port")
	err := r.Run(addr + ":" + port)
	if err != nil {
		log.Fatalln("服务启动失败 ：", err)
	}
}
