package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mod/model"
	"go.mod/pkg"
	"go.mod/utils"
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
	token := c.Query("token")
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
	videoId, err := utils.Str2int64(c.Query("video_id"))
	if err != nil {
		c.JSON(http.StatusOK, pkg.ParamErr)
		return
	}
	actionType := c.Query("action_type")
	if actionType != "1" && actionType != "2" {
		c.JSON(http.StatusOK, pkg.ParamErr)
		return
	}

	var temp model.Favorite
	res := model.Mysql.Model(&model.Favorite{}).Where("user_id=? and video_id=?", userId, videoId).Find(&temp)
	// 判断赞操作action_type
	if actionType == "1" {
		if res.RowsAffected == 0 {
			if code := model.CreateFavorite(userId, videoId); code == -1 {
				c.JSON(http.StatusOK, pkg.DataBaseErr.WithMessage("Favorite action failed"))
				return
			}
			// 点赞完成，更新video点赞数+1
			model.UpdataVideoFavoriteCount(videoId, 1)

			c.JSON(http.StatusOK, pkg.Success)
			return
		} else {
			c.JSON(http.StatusOK, pkg.RecordAlreadyExistErr)
		}
	} else if actionType == "2" {
		if res.RowsAffected == 0 {
			c.JSON(http.StatusOK, pkg.RecordNotExistErr)
		} else {
			code := model.DeleteFavorite(temp)
			if code == -1 {
				c.JSON(http.StatusOK, pkg.DataBaseErr.WithMessage("Cancel favorite action failed"))
				return
			}
			// 取消点赞完成，更新video点赞数-1
			model.UpdataVideoFavoriteCount(videoId, -1)

			c.JSON(http.StatusOK, pkg.Success)
			return
		}
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	user_id, err := utils.Str2int64(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusOK, pkg.ParamErr)
		return
	}
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
		return
	}
	// 验证token有效性
	if token == "" {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
	}
	// 解析token获取user_id
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
	}

	var videoId_list []int64
	var rsp_video_list []Res_Video
	// 查询登录用户点赞视频的ID
	model.Mysql.Model(&model.Favorite{}).Select("video_id").Where("user_id=?", user_id).Find(&videoId_list)
	if len(videoId_list) != 0 {
		var video_list []model.Video
		model.Mysql.Model(&model.Video{}).Where("id in (?)", videoId_list).Find(&video_list)
		for _, video := range video_list {
			// 查看登录用户是否了关注视频作者
			isfollow := model.SearchIsFollow(claims.UserId, video.UserId)
			// 查找视频作者
			var user model.User
			model.Mysql.Where("id=?", video.UserId).Find(&user)
			var author Res_User = Convert2ResUser(user, isfollow)
			// 转换为Res_Video
			var rsp_video Res_Video = Convert2ResVideo(video, author, true)
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
