package model

import (
	"time"

	"github.com/google/uuid"
)

func CreateFollow(userId int64, followId int64) int {
	follow := Follow{
		Id:         int64(uuid.New().ID()),
		UserId:     userId,
		FollowId:   followId,
		CreateTime: time.Now(),
	}
	res := Mysql.Create(&follow)
	if res.RowsAffected != 1 {
		//特殊情况导致插入失败
		return -1
	}
	return 0
}

func DeleteFollow(temp Follow) int {
	res := Mysql.Delete(&temp)
	if res.RowsAffected != 1 {
		//特殊情况导致删除失败
		return -1
	}
	return 0
}

func SearchIsFollow(userId int64, followId int64) bool {
	isfollow := true
	if userId != followId {
		res := Mysql.Table("tb_follow").Where("user_id = ? and follow_id = ?", userId, followId).Find(&Follow{})
		if res.RowsAffected == 0 {
			isfollow = false
		}
	}
	return isfollow
}
