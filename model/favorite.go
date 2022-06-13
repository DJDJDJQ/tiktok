package model

import (
	"time"

	"github.com/google/uuid"
)

func CreateFavorite(userId int64, videoId int64) int {
	favorite := Favorite{
		Id:         int64(uuid.New().ID()),
		UserId:     userId,
		VideoId:    videoId,
		CreateTime: time.Now(),
	}
	res := Mysql.Create(&favorite)
	if res.RowsAffected != 1 {
		//特殊情况导致插入失败
		return -1
	}
	// 更新点赞视频数和获赞数
	UpdateFavoriteCount(userId, 1)
	var authorId int64
	Mysql.Model(&Video{}).Select("user_id").Where("id = ?", videoId).Find(&authorId)
	UpdateTotalFavoritedCount(authorId, 1)
	return 0
}

func DeleteFavorite(temp Favorite) int {
	res := Mysql.Delete(&temp)
	if res.RowsAffected != 1 {
		//特殊情况导致删除失败
		return -1
	}
	// 更新点赞视频数和获赞数
	UpdateFavoriteCount(temp.UserId, -1)
	var authorId int64
	Mysql.Model(&Video{}).Select("user_id").Where("id = ?", temp.VideoId).Find(&authorId)
	UpdateTotalFavoritedCount(authorId, 1)
	return 0
}

func SearchIsFavorite(userId int64, videoId int64) bool {
	isfavor := false
	res := Mysql.Table("tb_favorite").Where("user_id = ? AND video_id = ?", userId, videoId).Find(&Favorite{})
	if res.RowsAffected != 0 {
		isfavor = true
	}
	return isfavor
}
