package utils

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// 获取项目执行路径
func getExecPath() string {
	p, err := os.Executable()
	if err != nil {
		panic("read config file path failed -->" + err.Error())
	}
	return filepath.Dir(p)
}

// 获取命令执行路径
func getWdPath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic("read config file path failed -->" + err.Error())
	}
	return dir
}

func Config() *viper.Viper {
	return initConfig()
}

// InitConfig 初始化配置文件
func initConfig() *viper.Viper {
	vip := viper.New()
	vip.SetConfigFile("./config.yaml")
	vip.AddConfigPath(getExecPath())

	// 读取配置文件
	if err := vip.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			vip.AddConfigPath(getWdPath())
		} else {
			panic("read config file failed -->" + err.Error())
		}
	}

	// 实时监测配置文件修改
	vip.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed -->", e.Name)
	})
	vip.WatchConfig()

	return vip
}
