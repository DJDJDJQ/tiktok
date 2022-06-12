package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateVideo(userId int64, play_url string, cover_url string, title string) int {
	video := Video{
		Id:          int64(uuid.New().ID()),
		UserId:      userId,
		PlayUrl:     play_url,
		CoverUrl:    cover_url,
		Title:       title,
		PublishTime: time.Now(),
		Status:      0,
	}
	res := Mysql.Model(&Video{}).Create(&video)
	if res.RowsAffected == 0 {
		return -1
	}
	return 0
}

func GetVideosByUserId(userId int64) []Video {
	var video_list []Video
	Mysql.Model(&Video{}).Where("user_id=?", userId).Find(&video_list)
	return video_list
}

func UpdataVideoFavoriteCount(videoId int64, cnt int) {
	if cnt > 0 {
		Mysql.Model(&Video{}).Where("id = ?", videoId).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", cnt))
	} else if cnt < 0 {
		Mysql.Model(&Video{}).Where("id = ?", videoId).UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 0-cnt))
	}
	return
}
