package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mod/model"
	"go.mod/pkg"
	"go.mod/utils"
	"gorm.io/gorm"
)

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
	//TODO检验user_id合法性
	if action_type != "1" && action_type != "2" {
		c.JSON(http.StatusOK, pkg.ParamErr)
		return
	}
	var comment_temp model.Comment
	resSearch := model.Mysql.Model(&model.Comment{}).Where("id=?", user_id).First(&comment_temp)

	if action_type == "1" { //新建评论
		if resSearch.RowsAffected != 0 { //评论已存在
			c.JSON(http.StatusOK, pkg.RecordAlreadyExistErr)
			return
		} else {
			CommitIDNew := int64(uuid.New().ID())       //生成随机CommitID
			CommitTimeNew := time.Now().Format("01-02") //生成评论时间
			comment := model.Comment{
				Id:         CommitIDNew,
				UserId:     user_id,
				VideoId:    video_id,
				Content:    comment_text,
				CreateDate: CommitTimeNew,
			}
			resCreate := model.Mysql.Create(&comment)
			if resCreate.RowsAffected != 1 {
				//特殊情况导致插入失败
				c.JSON(http.StatusOK, pkg.ServiceErrCode)
				return
			}
			var user_temp model.User
			model.Mysql.Model(&model.User{}).Where("id=?", user_id).First(&user_temp)
			//查询成功
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 0},
				Comment: Res_Comment{
					Id: CommitIDNew,
					User: Res_User{
						Id:            user_temp.Id,
						Name:          user_temp.Name,
						FollowCount:   user_temp.FollowCount,
						FollowerCount: user_temp.FollowerCount,
						IsFollow:      user_temp.IsFollow,
					},
					Content:    comment_text,
					CreateDate: CommitTimeNew,
				}})
			//更新video评论数
			model.Mysql.Model(&model.Video{}).Where("id = ?", video_id).
				UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))
		}
	} else if action_type == "2" { //删除评论
		if resSearch.RowsAffected != 1 { //评论不存在
			c.JSON(http.StatusOK, pkg.RecordNotExistErrCode)
			return
		} else {
			var comment model.Comment
			resDelete := model.Mysql.Where("Id=?", comment_id).Delete(&comment)
			if resDelete.RowsAffected != 1 {
				//特殊情况导致删除失败
				c.JSON(http.StatusOK, pkg.ServiceErrCode)
				return
			}
			c.JSON(http.StatusOK, pkg.Success)
			//更新video评论数
			model.Mysql.Model(&model.Video{}).Where("id = ?", video_id).
				UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))
		}
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
