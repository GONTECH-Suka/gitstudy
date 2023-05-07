package main

import (
	"Gin_EdMaSys/database"
	"Gin_EdMaSys/middleware"
	"Gin_EdMaSys/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
)

func main() {
	// 读取配置文件信息
	middleware.InitConfig("config", "ini", "configs/")
	// 读取数据库
	database.InitDataBase(middleware.Config)
	// 让输出的日志带颜色
	gin.ForceConsoleColor()

	// 初始化Casbin
	middleware.InitCasbinEnforcer()
	// 设置路由
	ginServer := router.SetRouters()

	// 启动网页
	address := fmt.Sprint("localhost:", middleware.Config.GetInt("server.port"))
	color.Greenln(fmt.Sprintf("Starting HTTP service on %s...", address))
	// 设置端口
	ginError := ginServer.Run(address)
	if ginError != nil {
		color.Redln("Failed! Address may be occupied!")
	}

}
