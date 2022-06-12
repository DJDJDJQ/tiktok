package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mod/model"
	"go.mod/pkg"
	"go.mod/utils"
)

type VideoListResponse struct {
	Response
	VideoList []Res_Video `json:"video_list"`
}

type CommentActionResponse struct {
	Response
	Comment Res_Comment `json:"comment,omitempty"`
}

type CommentListResponse struct {
	Response
	CommentList []Res_Comment `json:"comment_list,omitempty"`
}

type PublishVideo struct {
	Token string `json:"token,omitempty"`
	Data  []byte `json:"data,omitempty"`
	Title string `json:"title,omitempty"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	// TODO 验证token有效性
	if token == "" {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
		return
	}

	// 解析token获取user_id
	parseToken, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
	}
	userId := parseToken.UserId

	// 检查用户是否存在
	var user model.User
	result := model.Mysql.Where(" id = ?", userId).Find(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, pkg.RecordNotExistErr.WithMessage("User doesn't exist"))
		return
	}

	title := c.PostForm("title")

	// 上传视频，获取视频play_url

	play_url, cover_url, err := stroage_upload(user.Id, c)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// //存储在本地
	// TODO 截取封面，获取cover_url 目前是写死的

	// 视频封面写死
	// snapShotName, err := service.GetSnapshot(finalName, finalName, 1)
	// if err != nil {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 7002, StatusMsg: "视频封面截取失败"})
	// }
	// cover_url := "http://" + host + ":8080/?url=" + snapShotName
	// cover_url := "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg"

	// 添加video
	var video model.Video = model.Video{
		Id:          int64(uuid.New().ID()),
		UserId:      parseToken.UserId,
		PlayUrl:     play_url,
		CoverUrl:    cover_url,
		Title:       title,
		PublishTime: time.Now(),
		Status:      0,
	}
	res := model.Mysql.Model(&model.Video{}).Create(&video)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusOK, pkg.ServiceErr)
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")
	// 验证参数有效性
	if token == "" {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
		return
	}
	// // 解析token获取user_id
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
		return
	}

	// 查找user_id发布的视频
	var video_list []model.Video
	model.Mysql.Model(&model.Video{}).Where("user_id=?", user_id).Find(&video_list)
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
	model.Mysql.Model(&model.User{}).Where("id=?", user_id).Find(&user)
	var author Res_User = Res_User{
		Id:            user.Id,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      isfollow,
	}
	// 转换为Res_Video
	var rsp_video_list []Res_Video
	for _, video := range video_list {
		favorite := model.Favorite{}
		isfavor := false
		res := model.Mysql.Table("tb_favorite").Where("user_id = ? AND video_id = ?", claims.UserId, video.Id).Find(&favorite)
		if res.RowsAffected != 0 {
			isfavor = true
		}
		var rsp_video Res_Video = Res_Video{
			Id:            video.Id,
			Author:        author,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    isfavor,
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
