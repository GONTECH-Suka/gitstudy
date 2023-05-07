package controller

import (
	"Gin_EdMaSys/database"
	"Gin_EdMaSys/middleware"
	"Gin_EdMaSys/model"
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	"net/http"
)

// IndexFunc 加载首页
func IndexFunc(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"code": http.StatusOK,
		"msg":  "我是一个一个首页啊",
	})
}

// LoginGetFunc 加载登录界面
func LoginGetFunc(c *gin.Context) {
	c.HTML(http.StatusOK, "index/login.tmpl", gin.H{})
}

// ConfirmCasbin 验证登录者权限并设置cookie
func ConfirmCasbin(c *gin.Context, res bool, role string, cookieName string, token string) {
	if res == true {
		// 这样设置可以让每次请求的请求头都带上token
		c.SetCookie(cookieName, token, 3600, "/", "", false, false)
		// 权限验证成功返回信息
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusMovedPermanently,
			"msg":  "Welcome",
			"role": role,
		})
		return
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "User does not have permission to visit this part...",
			"role": role,
		})
		return
	}

}

// LoginPostFunc 处理登录请求
func LoginPostFunc(c *gin.Context) {

	var loginUser model.LoginUser
	// 获取axios传来的json数据
	err := c.ShouldBind(&loginUser)
	if err != nil {
		return
	}

	// 搜索数据库的登录者数据(用户表以及管理员表)
	var userInfo model.UserKeyInfo
	database.DB.Table("user").Select([]string{"user_id", "user_name", "password"}).Where("user_name = ?", loginUser.UserName).Find(&userInfo)
	/*
		不知为何提取管理员数据时使用 UserKeyInfo 结构体便无法正确提取数据
		以后提取不同角色的数据最好还是用不同结构体
	*/
	var adminInfo model.AdminKeyInfo
	database.DB.Table("admin").Select([]string{"admin_id", "admin_name", "password"}).Where("admin_name = ?", loginUser.UserName).Find(&adminInfo)

	// 若两个表都搜不到信息说明用户不存在
	if userInfo.UserName == "" && adminInfo.AdminName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Account not exist",
			"role": "guest",
		})
		return
	}

	// 在用户表中查询到了用户信息
	if userInfo.UserName != "" {
		// 确认密码无误
		if loginUser.Password == userInfo.Password {

			// 创建用户的token提取完整的字符串并记录
			token, err := middleware.GetToken(userInfo)
			if err != nil {
				color.Redln("ERROR!", err.Error())
				return
			}
			middleware.LoginToken = token

			// 验证访问权限
			res, _ := middleware.Enforcer.Enforce(loginUser.UserName, "templates/user", "normal")
			ConfirmCasbin(c, res, "User", "user_token", token)

			// 密码错误
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"role": "User",
				"msg":  "Password wrong!",
			})
			return
		}
	}

	// 在管理员表查到密码不为空说明登录的不是用户而是管理员
	if adminInfo.Password != "" {
		// 确认密码无误
		if loginUser.Password == adminInfo.Password {

			// 创建管理员的token并提取完整的字符串保存进上下文
			u := model.UserKeyInfo{
				UserID:   adminInfo.AdminID,
				UserName: adminInfo.AdminName,
				Password: adminInfo.Password,
			}
			token, err := middleware.GetToken(u)
			if err != nil {
				color.Redln("ERROR! ", err.Error())
				return
			}
			middleware.LoginToken = token

			// 验证访问权限
			res, _ := middleware.Enforcer.Enforce(loginUser.UserName, "templates/admin", "supreme")
			ConfirmCasbin(c, res, "Admin", "admin_token", token)

			// 密码错误
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"role": "Admin",
				"msg":  "Password wrong!",
			})
		}
	}
}

// GetUserHomePage 加载用户主页面
func GetUserHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "user/home.tmpl", gin.H{
		"code": http.StatusOK,
		"msg":  "我是一个一个用户",
	})
}

// GetAdminHomePage 加载管理员主页面
func GetAdminHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/home.tmpl", gin.H{
		"code": http.StatusOK,
		"msg":  "我是一个一个管理员",
	})
}
