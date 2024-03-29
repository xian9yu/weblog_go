package dbConn

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"os"
	"time"
	"weblog/models"
	"weblog/utils"

	"log"
)

var v = utils.InitConfig()

func MySQL() *gorm.DB {
	return initSQL()
}

// InitSQL 初始化sql
func initSQL() *gorm.DB {

	dsn := v.GetString("db.mysql.dsn")

	client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   "db_", // 表名前缀，`User` 的表名应该是 `db_users`
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `we_user`
		},
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Nanosecond, // 慢 SQL 阈值
				LogLevel:      logger.Info,     // Log level
				Colorful:      true,            // 彩色打印
			},
		),
	})
	if err != nil {
		log.Fatalln("MySQL连接异常, 请检查后再试")
	}

	// 自动同步数据库结构
	if err = client.AutoMigrate(
		&models.Article{},
		&models.Category{},
		&models.Links{},
		&models.User{},
		&models.Auth{},
		&models.Setting{},
		&models.Upload{},
	); err != nil {
		log.Println("同步数据库结构失败: ", err.Error())

	}

	return client
}
