package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mod/model"
	"go.mod/pkg"
	"go.mod/service"
	"go.mod/utils"
)

type UserListResponse struct {
	Response
	UserList []Res_User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
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

	toUserId, err := utils.Str2int64(c.Query("to_user_id"))
	if err != nil {
		c.JSON(http.StatusOK, pkg.ParamErr)
		return
	}
	actionType := c.Query("action_type")

	var temp model.Follow
	res := model.Mysql.Model(&model.Follow{}).Where("user_id=? and follow_id=?", userId, toUserId).Find(&temp)
	switch actionType {
	case "1": //关注
		if res.RowsAffected == 0 {
			response := service.Follow(userId, toUserId)
			c.JSON(http.StatusOK, response)
		} else {
			c.JSON(http.StatusOK, pkg.RecordAlreadyExistErr)
		}
	case "2": //取消关注
		if res.RowsAffected == 0 {
			c.JSON(http.StatusOK, pkg.RecordNotExistErr)
		} else {
			response := service.CancelFollow(temp, userId, toUserId)
			c.JSON(http.StatusOK, response)
		}
	default:
		c.JSON(http.StatusOK, pkg.ParamErr)
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
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

	userId, err := utils.Str2int64(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusOK, pkg.ParamErr)
		return
	}

	var followId_list []int64
	model.Mysql.Model(&model.Follow{}).Select("follow_id").Where("user_id=?", userId).Find(&followId_list)
	var rsp_user_list []Res_User
	if len(followId_list) > 0 {
		// 查找关注用户，返回列表
		var user_list []model.User
		model.Mysql.Model(&model.User{}).Where("id in (?)", followId_list).Find(&user_list)
		for _, user := range user_list {
			isfollow := model.SearchIsFollow(claims.UserId, user.Id)
			var rsp_user Res_User = Convert2ResUser(user, isfollow)
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

	userId, err := utils.Str2int64(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusOK, pkg.ParamErr)
		return
	}

	var followerId_list []int64
	model.Mysql.Model(&model.Follow{}).Select("user_id").Where("follow_id=?", userId).Find(&followerId_list)
	var rsp_user_list []Res_User
	if len(followerId_list) > 0 {
		// 查找粉丝用户，返回列表
		var user_list []model.User
		model.Mysql.Model(&model.User{}).Where("id in (?)", followerId_list).Find(&user_list)
		for _, user := range user_list {
			isfollow := model.SearchIsFollow(claims.UserId, user.Id)
			var rsp_user Res_User = Convert2ResUser(user, isfollow)
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
