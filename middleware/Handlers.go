package middleware

import (
	"Gin_EdMaSys/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthCheckHandler 进行拦截，不登录无法访问除主页以及登录以外的网页
func AuthCheckHandler(tokenKeeper string) gin.HandlerFunc {

	return func(c *gin.Context) {

		// 从请求头中获取 token
		cookie, err := c.Request.Cookie(tokenKeeper)
		if err != nil {
			// 处理获取 Cookie 失败的情况
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "Cookie has something wrong",
			})
		}
		tokenString := cookie.Value

		// 返回错误信息
		returnErrorMsg := func(c *gin.Context, msg string) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": msg,
			})
			c.Abort()
			return
		}

		// 如果没有 token 返回错误信息
		if tokenString == "" {
			returnErrorMsg(c, "Authorization token not found")
		}

		// 解析前端传来的token
		fontToken, err1 := jwt.ParseWithClaims(tokenString, &model.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法，检查token中的签名方法是否为HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("Unexpected signing method", 401)
			}

			return MySecret, nil
		})
		// 如果解析出错，返回错误信息
		if err1 != nil {
			returnErrorMsg(c, "Authorization token can not be decoded properly")
		}

		// 解析后端保存的token
		savedToken, err2 := jwt.ParseWithClaims(LoginToken, &model.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法，检查token中的签名方法是否为HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("Unexpected signing method", 401)
			}

			return MySecret, nil
		})
		if err2 != nil {
			returnErrorMsg(c, "Authorization token can not be decoded properly")
		}

		// 解析出的结果与保存的相同则验证成功,可以继续访问
		if savedToken.Raw == fontToken.Raw {
			c.Next()
		} else {
			returnErrorMsg(c, "Authorization token is not valid")
		}
	}

}
