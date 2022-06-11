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
	user_id := parseToken.UserId                     //用户id
	video_id := utils.Str2int64(c.Query("video_id")) //视频id
	action_type := c.Query("action_type")            //1-发布评论，2-删除评论
	comment_text := c.Query("comment_text")          //可选，用户填写的评论内容
	comment_id := c.Query("comment_id")              //可选，要删除的评论id

	//验证参数合法性
	//TODO检验user_id合法性
	if action_type != "1" && action_type != "2" {
		c.JSON(http.StatusOK, pkg.ParamErr)
		return
	}
	var comment_temp model.Comment
	resSearch := model.Mysql.Model(&model.Comment{}).Where("id=?", comment_id).Find(&comment_temp)

	if action_type == "1" { //新建评论
		//TODO 如果comment_id不为空可能导致下述情况
		if resSearch.RowsAffected != 0 { //评论已存在
			c.JSON(http.StatusOK, pkg.RecordAlreadyExistErr)
			return
		} else {
			CommitIDNew := int64(uuid.New().ID()) //生成随机CommitID
			CommitTimeNew := time.Now()           //生成评论时间
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
					},
					Content:    comment_text,
					CreateDate: CommitTimeNew.Format("01-02"),
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
			resDelete := model.Mysql.Where("id=?", comment_id).Delete(&comment)
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
	token := c.Query("token")                        //用户鉴权token
	video_id := utils.Str2int64(c.Query("video_id")) //视频id
	//判断token
	//println(token)
	//不需要登录也可以查看评论

	// if token == "" {
	// 	c.JSON(http.StatusOK, pkg.TokenInvalidErr)
	// 	return
	// }
	var video model.Video
	resSearch := model.Mysql.Model(&model.Video{}).Where("id=?", video_id).First(&video)
	if resSearch.RowsAffected != 1 { //未找到视频
		c.JSON(http.StatusOK, pkg.RecordNotExistErrCode)
		return
	}
	//搜索到视频
	var comment_list []model.Comment  //评论
	var rescomment_list []Res_Comment //返回的评论列表
	model.Mysql.Model(&model.Comment{}).Where("video_id=?", video_id).Find(&comment_list)

	//model.Comment To Res_Comment
	for _, comment_temp := range comment_list {
		var commenter model.User //该评论的用户信息
		resSearchUser := model.Mysql.Model(&model.User{}).Where("id=?", comment_temp.UserId).First(&commenter)
		if resSearchUser.RowsAffected != 1 {
			c.JSON(http.StatusOK, pkg.RecordNotExistErrCode)
			return
		}
		// 改动 IsFollow，判断评论者是否关注了视频作者
		isfollow := false
		if token != "" {
			follow := model.Follow{}
			res := model.Mysql.Table("tb_follow").Where("user_id = ? and follow_id = ?", commenter.Id, video.UserId).Find(&follow)
			if res.RowsAffected == 0 {
				isfollow = true
			}
		}
		//

		var res_commenter Res_User = Res_User{
			Id:            commenter.Id,
			Name:          commenter.Name,
			FollowCount:   commenter.FollowCount,
			FollowerCount: commenter.FollowerCount,
			IsFollow:      isfollow,
		}

		rescomment_list = append(rescomment_list,
			Res_Comment{
				Id:         comment_temp.Id,
				User:       res_commenter,
				Content:    comment_temp.Content,
				CreateDate: comment_temp.CreateDate.Format("01-02"),
			})
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: rescomment_list,
	})

}
