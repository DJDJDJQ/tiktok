package model

import "gorm.io/gorm"

func GetUserById(userId int64) User {
	var user User
	Mysql.Model(&User{}).Where("id=?", userId).Find(&user)
	return user
}

func UpdataUserFollowCount(userId int64, cnt int) {
	if cnt > 0 {
		Mysql.Model(&User{}).Where("id = ?", userId).UpdateColumn("follow_count", gorm.Expr("follow_count + ?", cnt))

	} else if cnt < 0 {
		Mysql.Model(&User{}).Where("id = ?", userId).UpdateColumn("follow_count", gorm.Expr("follow_count - ?", 0-cnt))
	}
	return
}

func UpdataUserFollowerCount(userId int64, cnt int) {
	if cnt > 0 {
		Mysql.Model(&User{}).Where("id = ?", userId).UpdateColumn("follower_count", gorm.Expr("follower_count + ?", cnt))

	} else if cnt < 0 {
		Mysql.Model(&User{}).Where("id = ?", userId).UpdateColumn("follower_count", gorm.Expr("follower_count - ?", 0-cnt))
	}
	return
}
