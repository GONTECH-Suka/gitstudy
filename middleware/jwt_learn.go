package middleware

import (
	"Gin_EdMaSys/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gookit/color"
	"time"
)

// LoginToken 获取登录者的token
var LoginToken string

// MySecret 从配置文件获取密钥
var MySecret = []byte(Config.GetString("jwt.mySecret"))

func GetToken(userInfo model.UserKeyInfo) (string, error) {
	c := model.MyClaims{
		UserID:   userInfo.UserID,
		UserName: userInfo.UserName,
		Password: userInfo.Password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer:    "Jesus",
		},
	}
	// 使用指定的签名方法创建签名对象;以HS256方式加密，并设置token格式
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的mySecret密钥获得完整的编码后的字符串token
	tokenString, err := token.SignedString(MySecret)

	return tokenString, err
}

func InitJwt() {
	u := model.UserKeyInfo{
		UserID:   1,
		UserName: "114",
		Password: "514",
	}
	tokenString, _ := GetToken(u)
	// 解析完整的编码后的字符串即传入的tokenString
	tokenTarget, _ := jwt.ParseWithClaims(tokenString, &model.MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		// 返回密钥的作用是告诉jwt解析函数解析token时使用的密钥是什么。
		return MySecret, nil
	})
	color.Greenln(tokenTarget)
	// 可以用这个格式取出任意想要的数据
	color.Greenln(tokenTarget.Claims.(*model.MyClaims).UserName)
}
