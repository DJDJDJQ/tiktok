package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mod/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Status Code 说明
var Success int32 = 0
var PostDataError int32 = 1
var ParamsError int32 = 2
var NoDataError int32 = 3

type favorite_action_request struct {
	Userid     int64  `json:"user_id,omitempty"`
	Token      string `json:"token,omitempty"`
	Videoid    int64  `json:"video_id,omitempty"`
	Actiontype int32  `json:"action_type,omitempty"`
}

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	// 获取post请求中的参数
	var req favorite_action_request
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: PostDataError,
			StatusMsg:  "Post data parameter errors"})
	}

	// req.Userid = c.Query("user_id")
	// req.Token = c.Query("token")
	// req.Videoid = c.Query("video_id")
	// req.Actiontype = c.Query("action_type")

	// 验证参数合法性
	if req.Token == "" {
		c.JSON(http.StatusOK, Response{
			StatusCode: ParamsError,
			StatusMsg:  "Token is invalid"})
	} else if req.Actiontype != 1 && req.Actiontype != 2 {
		c.JSON(http.StatusOK, Response{
			StatusCode: ParamsError,
			StatusMsg:  "Action type is invalid"})
	}

	db, err := gorm.Open(
		mysql.Open("root:leelee@tcp(127.0.0.1:3306)/douyin"),
	)
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	var temp model.Favorite
	db.Model(&model.Favorite{}).Where("user_id=? and video_id=? and is_cancel=0", req.Userid, req.Videoid).First(&temp)
	// 判断赞操作action_type
	if req.Actiontype == 1 {
		if temp.Id != 0 {
			favorite := model.Favorite{
				Id:         int64(uuid.New().ID()),
				UserId:     req.Userid,
				VideoId:    req.Videoid,
				CreateTime: time.Now(),
				IsCancel:   false,
			}
			db.Create(&favorite)

			var video model.Video
			db.Model(&model.Video{}).Where("id=?", req.Videoid).First(&video)
			video.FavoriteCount += 1
			db.Save(&video)
		}
	} else if req.Actiontype == 2 {
		if temp.Id == 0 {
			c.JSON(http.StatusOK, Response{
				StatusCode: NoDataError,
				StatusMsg:  "Can't find this data",
			})
		} else {
			db.Model(&temp).Updates(map[string]interface{}{"is_cancel": true, "cancel_time": time.Now()})

			var video model.Video
			db.Model(&model.Video{}).Where("id=?", req.Videoid).First(&video)
			video.FavoriteCount -= 1
			db.Save(&video)
		}
	}
	c.JSON(http.StatusOK, Response{StatusCode: Success})
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")

	// 验证token有效性
	if token == "" {
		c.JSON(http.StatusOK, Response{
			StatusCode: ParamsError,
			StatusMsg:  "Token is invalid",
		})
	}

	// 连接数据库
	db, err := gorm.Open(
		mysql.Open("root:leelee@tcp(127.0.0.1:3306)/douyin"),
	)
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	var videoId_list []int64
	var rsp_video_list []Res_Video
	db.Model(&model.Favorite{}).Select("video_id").Where("user_id=? and is_cancel=0", user_id).Find(&videoId_list)
	if len(videoId_list) != 0 {
		var video_list []model.Video
		db.Model(&model.Video{}).Where("id in (?)", videoId_list).Find(&video_list)
		for _, video := range video_list {
			// 查找视频作者
			var user model.User
			db.Where("id=?", video.UserId).Find(&user)
			var author Res_User = Res_User{
				Id:            user.Id,
				Name:          user.Name,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      bool(user.IsFollow),
			}
			// 转换为Res_Video
			var rsp_video Res_Video = Res_Video{
				Id:            video.Id,
				Author:        author,
				PlayUrl:       video.PlayUrl,
				CoverUrl:      video.CoverUrl,
				FavoriteCount: video.FavoriteCount,
				CommentCount:  video.CommentCount,
				IsFavorite:    bool(video.IsFavorite),
				Title:         video.Title,
			}
			rsp_video_list = append(rsp_video_list, rsp_video)
		}
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: Success,
		},
		VideoList: rsp_video_list,
	})
}
