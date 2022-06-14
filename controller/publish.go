package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mod/model"
	"go.mod/pkg"
	"go.mod/service"
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
	// 解析token获取user_id
	parseToken, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
		return
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
	play_url, cover_url, err := service.Stroage_upload(user.Id, c)
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
	code := model.CreateVideo(parseToken.UserId, play_url, cover_url, title)
	if code == -1 {
		c.JSON(http.StatusOK, pkg.DataBaseErr.WithMessage("Publish video failed"))
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
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
	// 解析token获取user_id
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
		return
	}

	// 查找user_id发布的视频
	var video_list []model.Video = model.GetVideosByUserId(user_id)

	// 查找视频作者
	var user model.User = model.GetUserById(user_id)
	isfollow := model.SearchIsFollow(claims.UserId, user_id)
	var author Res_User = Convert2ResUser(user, isfollow)
	// 转换为Res_Video
	var rsp_video_list []Res_Video
	for _, video := range video_list {
		isfavor := model.SearchIsFavorite(claims.UserId, video.Id)
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
