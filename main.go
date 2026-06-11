package main

import (
	"log"
	"time"
	"weblog/ctrls"
	"weblog/database"
	"weblog/middleware"
	"weblog/models"
	"weblog/router"
	"weblog/utils"

	"github.com/gin-gonic/gin"
)

func initRouter(ac *ctrls.Controllers, r *gin.Engine) *gin.Engine {
	//  挂载全局标准中间件
	r.Use(gin.Logger())      // 打印请求日志
	r.Use(gin.Recovery())    // 崩溃恢复，防止 panic 导致进程挂掉
	r.Use(middleware.Cors()) // 如果有跨域需求，可以在这里挂载跨域中间件
	// 限制上传文件大小
	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	// 注册前台公开路由（传 controller 进去）
	router.InitBlogRoutes(r, ac)

	// 注册后台管理路由（传 controller 进去）
	router.InitAdminRoutes(r, ac)

	return r

}

func initApp(conf *utils.Config) (*gin.Engine, string) {
	r := gin.Default() // 初始化router

	// 初始化数据库连接
	db := database.InitMySQL(conf.Database.Mysql)

	// 初始化 Repository
	repos := models.NewRepositories(db)
	ac := ctrls.NewControllers(repos)

	// 载入路由系统，并把 database 等依赖注入进去
	r = initRouter(ac, r)

	// 模式为 release
	gin.SetMode(gin.ReleaseMode)

	return r, conf.App.Address

}

func main() {
	// 启动时一次性加载配置
	conf, err := utils.InitConfig("config.yaml")
	if err != nil {
		log.Fatalf("❌ 项目启动失败, 配置加载异常: %v\n", err)
	}

	// 将 Go 运行时的默认本地时区强行指向上海
	time.Local, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("强制设置 Go 运行时时区失败: %v", err)
	}

	r, addr := initApp(conf)
	if err := r.Run(addr); err != nil {
		log.Fatalln("服务启动失败 ：", err)
	}
}
