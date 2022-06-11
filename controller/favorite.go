package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mod/model"
	"go.mod/pkg"
	"go.mod/utils"
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
	// var req favorite_action_request
	// err := c.BindJSON(&req)
	// if err != nil {
	// 	c.JSON(http.StatusOK, pkg.ParamErr)
	// }

	token := c.Query("token")
	// TODO 验证token有效性
	if token == "" {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
		return
	}

	// 解析token获取user_id
	parseToken, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
		return
	}
	userId := parseToken.UserId

	// 验证参数合法性
	//userId := c.Query("user_id")
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")
	if videoId == "" {
		c.JSON(http.StatusOK, pkg.ParamErr)
		return
	} else if actionType != "1" && actionType != "2" {
		c.JSON(http.StatusOK, pkg.ParamErr)
		return
	}

	var temp model.Favorite
	res := model.Mysql.Model(&model.Favorite{}).Where("user_id=? and video_id=?", userId, videoId).Find(&temp)
	// 判断赞操作action_type
	if actionType == "1" {
		if res.RowsAffected == 0 {
			favorite := model.Favorite{
				Id:         int64(uuid.New().ID()),
				UserId:     userId,
				VideoId:    utils.Str2int64(videoId),
				CreateTime: time.Now(),
				IsCancel:   false,
			}
			resCreate := model.Mysql.Create(&favorite)
			if resCreate.RowsAffected != 1 {
				//特殊情况导致插入失败
				c.JSON(http.StatusOK, pkg.ServiceErrCode)
				return
			}
			// 更新video点赞数
			model.Mysql.Model(&model.Video{}).Where("id = ?", videoId).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1))

			c.JSON(http.StatusOK, pkg.Success)
		} else {
			c.JSON(http.StatusOK, pkg.RecordAlreadyExistErr)
		}
	} else if actionType == "2" {
		if res.RowsAffected == 0 {
			c.JSON(http.StatusOK, pkg.RecordNotExistErr)
		} else {
			resDelete := model.Mysql.Delete(&temp)
			if resDelete.RowsAffected != 1 {
				//特殊情况导致删除失败
				c.JSON(http.StatusOK, pkg.ServiceErrCode)
				return
			}
			// 更新video点赞数
			model.Mysql.Model(&model.Video{}).Where("id = ?", videoId).UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1))

			c.JSON(http.StatusOK, pkg.Success)
		}
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")
	// 验证token有效性
	if token == "" {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
	}
	// // 解析token获取user_id
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
	}

	var videoId_list []int64
	var rsp_video_list []Res_Video
	model.Mysql.Model(&model.Favorite{}).Select("video_id").Where("user_id=? and is_cancel=0", user_id).Find(&videoId_list)
	if len(videoId_list) != 0 {
		var video_list []model.Video
		model.Mysql.Model(&model.Video{}).Where("id in (?)", videoId_list).Find(&video_list)
		for _, video := range video_list {
			// 查找视频作者
			// 改动
			follow := model.Follow{}
			isfollow := true
			if claims.UserId != utils.Str2int64(user_id) {
				res := model.Mysql.Table("tb_follow").Where("user_id = ? and follow_id = ?", claims.UserId, user_id).Find(&follow)
				if res.RowsAffected == 0 {
					isfollow = false
				}
			}
			//
			var user model.User
			model.Mysql.Where("id=?", video.UserId).Find(&user)
			var author Res_User = Res_User{
				Id:            user.Id,
				Name:          user.Name,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      isfollow,
			}
			// 转换为Res_Video
			var rsp_video Res_Video = Res_Video{
				Id:            video.Id,
				Author:        author,
				PlayUrl:       video.PlayUrl,
				CoverUrl:      video.CoverUrl,
				FavoriteCount: video.FavoriteCount,
				CommentCount:  video.CommentCount,
				IsFavorite:    true,
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
