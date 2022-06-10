package model

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Mysql *gorm.DB

func init() {
	var err error
	//dsn := "root:leelee@tcp(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=true&loc=Local"       //本地运行测试用
	//dsn := "user_tiktok:123@tcp(127.0.0.1:3306)/douyin?charset=utf8&parseTime=true&loc=Local"      //服务器运行及测试测试用
	dsn := "user_tiktok:123@tcp(150.158.97.105:3306)/douyin?charset=utf8&parseTime=true&loc=Local" //本地运行服务器测试用

	logrus.Info("初始化数据库···")

	Mysql, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//err = Mysql.AutoMigrate(User{}, Comment{}, Like{}, Video{})
	//if err != nil {
	//	logrus.Errorln("表生成出错", err)
	//	return
	//}
}
