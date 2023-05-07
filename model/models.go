package model

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Admin struct {
	AdminID       int       `gorm:"primaryKey;"`
	AdminName     string    `gorm:"not null;Unique"`
	Password      string    `gorm:"not null"`
	LatestLogTime time.Time `gorm:"not null;column:LatestLogtime"`
}

type User struct {
	UserID        int       `gorm:"primaryKey"`
	UserName      string    `gorm:"not null;Unique"`
	Password      string    `gorm:"not null"`
	LatestLogTime time.Time `gorm:"not null;column:LatestLogtime"`
	rank          int       `gorm:"default:1"`
}

// MyClaims 自定义token的结构
type MyClaims struct {
	UserID   int    `json:"userid"`
	UserName string `json:"userName"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// LoginUser 存储前端提交的登录信息
type LoginUser struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// UserKeyInfo 从数据库获取的登录用户的关键信息
type UserKeyInfo struct {
	UserID   int
	UserName string
	Password string
}

// AdminKeyInfo 从数据库获取的登录管理员的关键信息
type AdminKeyInfo struct {
	AdminID   int
	AdminName string
	Password  string
}

func (Admin) TableName() string {
	return "admin"
}
func (User) TableName() string {
	return "user"
}
