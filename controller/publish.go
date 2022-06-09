package controller

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.mod/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type VideoListResponse struct {
	Response
	VideoList []Res_Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")

	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	user := usersLoginInfo[token]
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")

	// 验证参数有效性
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

	// 查找user_id发布的视频
	var video_list []model.Video
	db.Model(&model.Video{}).Where("user_id=?", user_id).Find(&video_list)
	// 查找视频作者
	var user model.User
	db.Model(&model.User{}).Where("id=?", user_id).Find(&user)
	var author Res_User = Res_User{
		Id:            user.Id,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      bool(user.IsFollow),
	}
	// 转换为Res_Video
	var rsp_video_list []Res_Video
	for _, video := range video_list {
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

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: Success,
		},
		VideoList: rsp_video_list,
	})
}
