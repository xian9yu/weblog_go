package utils

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

// ================= 定义配置结构体 =================

// Config 对应配置文件的最外层大骨架
type Config struct {
	App         AppConfig         `mapstructure:"app"`
	Database    DatabaseConfig    `mapstructure:"database"`
	SecurityKey SecurityKeyConfig `mapstructure:"securitykey"` // 注意：YAML里是驼峰或全小写，Viper默认不区分大小写，写成 securitykey 最稳妥
}

// AppConfig 对应 app 节点
type AppConfig struct {
	Address string `mapstructure:"address"`
}

// DatabaseConfig 对应 db 节点
type DatabaseConfig struct {
	Mysql string `mapstructure:"mysql"` // 直接用 string 接收整条 DSN 链接字符串
}

// SecurityKeyConfig 对应 securityKey 节点
type SecurityKeyConfig struct {
	Token    string `mapstructure:"token"`
	Password string `mapstructure:"password"`
}

// ================= 单例模式与加载逻辑 =================

var (
	globalConfig *Config
	configOnce   sync.Once
)

// InitConfig 初始化全局配置
func InitConfig(filePath string) (*Config, error) {
	var err error

	configOnce.Do(func() {
		v := viper.New()
		v.SetConfigFile(filePath)
		v.SetConfigType("yaml")

		if readErr := v.ReadInConfig(); readErr != nil {
			err = fmt.Errorf("读取配置文件失败: %w", readErr)
			return // 跳出匿名函数，此时 err 不为空
		}

		var conf Config
		if unmarshalErr := v.Unmarshal(&conf); unmarshalErr != nil {
			err = fmt.Errorf("解析配置文件失败: %w", unmarshalErr)
			return // 跳出匿名函数，此时 err 不为空
		}

		// 🎉 只有完全成功，才赋值给全局单例
		globalConfig = &conf
	})

	// 🛡️ 防御升级：如果单例内部由于报错没赋值成功，强制将 Once 重置，或者直接返回错误
	if err != nil {
		// 很多框架在这里直接选择 panic，因为配置加载失败，程序启动没有任何意义
		// panic(err)
		return nil, err
	}

	return globalConfig, nil
}

// GetConfig 获取配置单例
func GetConfig() *Config {
	if globalConfig == nil {
		panic("【致命错误】: 配置未初始化，请检查 main.go 是否在最前方案行了 utils.InitConfig()")
	}
	return globalConfig
}
