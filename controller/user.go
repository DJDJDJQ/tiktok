package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mod/model"
	"go.mod/pkg"
	"go.mod/utils"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]Res_User{
	"zhangleidouyin": {
		Id:            2,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User Res_User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if username == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 2,
			"status_msg":  "The account or password is empty",
		})
		return
	}
	//密码必须同时包含字母和数字，且长度为6-16
	expr := `^(?![a-zA-Z]+$)(?![0-9]+$)[0-9A-Za-z]{6,16}$`
	reg, _ := regexp2.Compile(expr, 0)
	m, _ := reg.FindStringMatch(password)
	if m != nil {
		password = m.String()
		log.Println(password)
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 3, StatusMsg: "密码必须同时包含字母和数字,且长度为6-16"},
		})
		return
	}

	userlogin := model.User{}

	res := model.Mysql.Table("tb_user").Where("name = ?", username).Find(&userlogin)

	if res.RowsAffected != 0 { //重复注册
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		newUser := model.User{
			Id:           int64(uuid.New().ID()),
			Name:         username,
			Password:     utils.StrEncrypt(password, username),
			RegisterTime: time.Now(),
		}
		tokenstring, err := utils.GenToken(newUser.Id, username)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status_code": 2,
				"status_msg":  "Failed to generate token",
			})
			return
		}

		model.Mysql.Table("tb_user").Save(&newUser)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "success"},
			UserId:   newUser.Id,
			Token:    tokenstring,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	//验证是否输入为空
	if username == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 2,
			"status_msg":  "账号或密码为空",
		})
		return
	}

	//处理密码，加密后与数据库中的储存相比较
	var user model.User
	res := model.Mysql.Where(" name = ?", username).Find(&user)
	//如果用户名不存在
	if res.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "用户名不存在",
		})
		return
	}
	Match_res := utils.StrMatch(user.Password, password, username)
	if !Match_res {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "密码错误",
		})
		return
	}

	// 获取token
	token, err := utils.GenToken(user.Id, user.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "登录成功",
		"user_id":     user.Id,
		"token":       token,
	})
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	// 解析token获取user_id
	parseToken, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, pkg.TokenInvalidErr)
	}
	userId := parseToken.UserId

	var user model.User
	result := model.Mysql.Where(" id = ?", userId).Find(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 2,
			"status_msg":  "该用户不存在",
		})
		return
	} else {
		var rsp_user Res_User = Res_User{
			Id:             user.Id,
			Name:           user.Name,
			FollowCount:    user.FollowCount,
			FollowerCount:  user.FollowerCount,
			IsFollow:       true,
			TotalFavorited: user.TotalFavorited,
			FavoriteCount:  user.FavoriteCount,
		}
		c.JSON(http.StatusOK, gin.H{
			"status_code": 0,
			"status_msg":  "success",
			"user":        rsp_user,
		})
	}
}
