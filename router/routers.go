package router

import (
	"Gin_EdMaSys/controller"
	"Gin_EdMaSys/middleware"
	"github.com/gin-gonic/gin"
)

func SetRouters() *gin.Engine {

	ginServer := gin.Default()

	// 加载当前目录下以及所有子目录下的模板文件
	ginServer.LoadHTMLGlob("templates/**/*")

	ginServer.GET("/", controller.IndexFunc)

	ginServer.GET("/login", controller.LoginGetFunc)
	ginServer.POST("/login", controller.LoginPostFunc)
	ginServer.GET("/user/home", middleware.AuthCheckHandler("user_token"), controller.GetUserHomePage)
	ginServer.GET("/admin/home", middleware.AuthCheckHandler("admin_token"), controller.GetAdminHomePage)

	return ginServer
}
