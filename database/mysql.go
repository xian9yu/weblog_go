package database

import (
	"log"
	"os"
	"time"
	"weblog/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitMySQL 初始化sql
func InitMySQL(dsn string) *gorm.DB {
	client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//NamingStrategy: schema.NamingStrategy{
		// TablePrefix:   "db_", // 表名前缀，`User` 的表名应该是 `db_users`
		//	SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `we_user`
		//},
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				//SlowThreshold: time.Nanosecond, // 慢 SQL 阈值
				LogLevel: logger.Info, // Log level
				Colorful: false,       // 彩色打印
			},
		),
	})
	if err != nil {
		log.Fatalf("❌ 数据库连接失败: %v", err)
	}

	// 设置连接池优化（防止本地开发时频繁断开）
	sqlDB, _ := client.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动同步数据库结构
	if err = client.AutoMigrate(
		&models.Article{},
		&models.User{},
		&models.SystemConfig{},
	); err != nil {
		log.Println("同步数据库结构失败: ", err.Error())
	}

	log.Println("✅ MySQL 数据库连接并初始化成功！")
	return client
}
