package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mod/model"
	"go.mod/pkg"
	"go.mod/utils"
	"gorm.io/gorm"
)

type CommentListResponse struct {
	Response
	CommentList []Res_Comment `json:"comment_list,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token") //用户鉴权token
	//判断token
	if token == "" {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
		return
	}
	// 解析token获取user_id
	parseToken, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
	}
	user_id := parseToken.UserId                                 //用户id
	video_id, _ := strconv.ParseInt(c.Query("video_id"), 10, 64) //视频id
	action_type := c.Query("action_type")                        //1-发布评论，2-删除评论
	comment_text := c.Query("comment_text")                      //可选，用户填写的评论内容
	comment_id := c.Query("comment_id")                          //可选，要删除的评论id

	//验证参数合法性
	if action_type != "1" && action_type != "2" {
		c.JSON(http.StatusOK, pkg.ParamErr)
		return
	}
	if action_type == "1" { //新建评论
		comment := model.Comment{
			Id:      int64(uuid.New().ID()),
			UserId:  user_id,
			VideoId: int64(video_id),
			Content: comment_text}
		model.Mysql.Create(&comment)
		c.JSON(http.StatusOK, pkg.Success)
		//更新video评论数
		model.Mysql.Model(&model.Video{}).Where("id = ?", video_id).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))

	} else { //删除评论
		var comment model.Comment
		model.Mysql.Where("Id=?", comment_id).Delete(&comment)
		c.JSON(http.StatusOK, pkg.Success)
		//更新video评论数
		model.Mysql.Model(&model.Video{}).Where("id = ?", video_id).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	token := c.Query("token")                                    //用户鉴权token
	video_id, _ := strconv.ParseInt(c.Query("video_id"), 10, 64) //视频id
	//判断token
	if token == "" {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
		return
	}
	var comment []Res_Comment
	model.Mysql.Model(&model.Comment{}).Where("VideoId=?", video_id).Find(comment)
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: comment,
	})
}
