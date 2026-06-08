package utils

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

// ================= 1. 定义配置结构体 =================

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

// ================= 2. 单例模式与加载逻辑 =================

var (
	globalConfig *Config
	configOnce   sync.Once
)

func InitConfig(filePath string) (*Config, error) {
	var err error
	configOnce.Do(func() {
		v := viper.New()
		v.SetConfigFile(filePath)
		v.SetConfigType("yaml")

		if readErr := v.ReadInConfig(); readErr != nil {
			err = fmt.Errorf("读取配置文件失败: %w", readErr)
			return
		}

		var conf Config
		// Viper 会根据 mapstructure 标签把数据塞进 conf 结构体
		if unmarshalErr := v.Unmarshal(&conf); unmarshalErr != nil {
			err = fmt.Errorf("解析配置文件失败: %w", unmarshalErr)
			return
		}

		globalConfig = &conf
	})

	return globalConfig, err
}

func GetConfig() *Config {
	if globalConfig == nil {
		panic("配置未初始化，请先在程序启动时调用 InitConfig()")
	}
	return globalConfig
}
