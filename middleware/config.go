package middleware

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// Config 创建config对象
var Config = viper.New()

// InitConfig 读取配置文件
func InitConfig(configName string, configType string, configPath string) {

	Config.SetConfigName(configName)
	Config.SetConfigType(configType)
	// path to look for the config file in
	Config.AddConfigPath(configPath)

	// Find and read the config file
	err := Config.ReadInConfig()
	if err != nil {
		color.Red("Fatal error config file: " + err.Error())
	}
	// 监控配置文件并且热加载程序，不重启就可以加载新的配置文件
	Config.WatchConfig()

	// 注意ini里的分级，每出现一个[xxx]便多加一个xxx.否则无法提取配置信息
	fmt.Println(Config.GetInt("server.port"))
	// 获取所有配置
	fmt.Println(Config.AllSettings())

}
