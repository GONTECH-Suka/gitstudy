package database

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDataBase(config *viper.Viper) {

	URL := fmt.Sprintf("%s:%d@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GetString("mysql.user"),
		config.GetInt("mysql.password"),
		config.GetString("mysql.host"),
		config.GetInt("mysql.port"),
		config.GetString("mysql.db"))

	var err error
	DB, err = gorm.Open(mysql.Open(URL), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 创建表时默认使用单数表名
		},
	})
	if err != nil {
		color.Red("Error: " + err.Error())
	}
}
