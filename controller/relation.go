package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mod/model"
	"go.mod/pkg"
	"go.mod/utils"
	"gorm.io/gorm"
)

type UserListResponse struct {
	Response
	UserList []Res_User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	// TODO 验证token有效性
	if token == "" {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
	}

	// 解析token获取user_id
	parseToken, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
	}
	userId := parseToken.UserId

	toUserId := c.Query("to_user_id")
	actionType := c.Query("action_type")
	if toUserId == "" {
		c.JSON(http.StatusOK, pkg.ParamErr)
	}

	var temp model.Follow
	res := model.Mysql.Model(&model.Follow{}).Where("user_id=? and follow_id=?", userId, toUserId).Find(&temp)
	fmt.Println(res)
	switch actionType {
	case "1": //关注
		if res.RowsAffected == 0 {
			follow := model.Follow{
				Id:         int64(uuid.New().ID()),
				UserId:     userId,
				FollowId:   utils.Str2int64(toUserId),
				CreateTime: time.Now(),
			}
			model.Mysql.Create(&follow)

			// 更新userId关注数、toUserId粉丝数
			model.Mysql.Model(&model.User{}).Where("id = ?", userId).UpdateColumn("follow_count", gorm.Expr("follow_count + ?", 1))
			model.Mysql.Model(&model.User{}).Where("id = ?", toUserId).UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1))

			c.JSON(http.StatusOK, pkg.Success)
		} else {
			c.JSON(http.StatusOK, pkg.RecordAlreadyExistErr)
		}
	case "2": //取消关注
		if res.RowsAffected == 0 {
			c.JSON(http.StatusOK, pkg.RecordNotExistErr)
		} else {
			model.Mysql.Delete(&temp)

			// 更新userId关注数、toUserId粉丝数
			model.Mysql.Model(&model.User{}).Where("id = ?", userId).UpdateColumn("follow_count", gorm.Expr("follow_count - ?", 1))
			model.Mysql.Model(&model.User{}).Where("id = ?", toUserId).UpdateColumn("follower_count", gorm.Expr("follower_count - ?", 1))

			c.JSON(http.StatusOK, pkg.Success)
		}
	default:
		c.JSON(http.StatusOK, pkg.ParamErr)
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	token := c.Query("token")
	// TODO 验证token有效性
	if token == "" {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
	}
	// // 解析token获取user_id
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
	}

	userId := c.Query("user_id")
	if userId == "" {
		c.JSON(http.StatusOK, pkg.ParamErr)
	}

	var followId_list []int64
	model.Mysql.Model(&model.Follow{}).Select("follow_id").Where("user_id=?", userId).Find(&followId_list)
	var rsp_user_list []Res_User
	if len(followId_list) > 0 {
		var user_list []model.User
		model.Mysql.Model(&model.User{}).Where("id in (?)", followId_list).Find(&user_list)
		for _, user := range user_list {
			// 改动
			follow := model.Follow{}
			isfollow := true
			res := model.Mysql.Table("tb_follow").Where("user_id = ? and follow_id = ?", claims.UserId, user.Id).Find(&follow)
			if res.RowsAffected == 0 {
				isfollow = false
			}
			//
			var rsp_user Res_User = Res_User{
				Id:            user.Id,
				Name:          user.Name,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      isfollow,
			}
			rsp_user_list = append(rsp_user_list, rsp_user)
		}
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: pkg.SuccessCode,
		},
		UserList: rsp_user_list,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	// TODO 验证token有效性
	if token == "" {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
	}
	// 解析token获取user_id
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
	}

	userId := c.Query("user_id")
	if userId == "" {
		c.JSON(http.StatusOK, pkg.ParamErr)
	}

	var followerId_list []int64
	model.Mysql.Model(&model.Follow{}).Select("user_id").Where("follow_id=?", userId).Find(&followerId_list)
	var rsp_user_list []Res_User
	if len(followerId_list) > 0 {
		var user_list []model.User
		model.Mysql.Model(&model.User{}).Where("id in (?)", followerId_list).Find(&user_list)
		for _, user := range user_list {
			// 改动
			follow := model.Follow{}
			isfollow := true
			if claims.UserId != user.Id {
				res := model.Mysql.Table("tb_follow").Where("user_id = ? and follow_id = ?", claims.UserId, user.Id).Find(&follow)
				if res.RowsAffected == 0 {
					isfollow = false
				}
			}
			//
			var rsp_user Res_User = Res_User{
				Id:            user.Id,
				Name:          user.Name,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      isfollow,
			}
			rsp_user_list = append(rsp_user_list, rsp_user)
		}
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: pkg.SuccessCode,
		},
		UserList: rsp_user_list,
	})
}
