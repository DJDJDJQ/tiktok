package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mod/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type CommentListResponse struct {
	Response
	CommentList []Res_Comment `json:"comment_list,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	user_id := c.Query("user_id")           //用户id
	token := c.Query("token")               //用户鉴权token
	video_id := c.Query("video_id")         //视频id
	action_type := c.Query("action_type")   //1-发布评论，2-删除评论
	comment_text := c.Query("comment_text") //可选，用户填写的评论内容
	comment_id := c.Query("comment_id")     //可选，要删除的评论id

	//token和数据库连接的代码可以复用
	//判断token
	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "token invalid"})
	}
	//建立数据库连接
	username := "root"  //账号
	password := "123"   //密码
	host := "127.0.0.1" //数据库地址，可以是Ip或者域名
	port := 3306        //数据库端口
	Dbname := "douyin"  //数据库名
	timeout := "10s"    //连接超时，10秒
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s",
		username, password, host, port, Dbname, timeout)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//数据库连接出问题
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "database err"})
	}
	if action_type == "1" { //新建评论
		comment := model.Comment{
			User:    model.User{Id: int64(user_id)},
			VideoId: int64(video_id),
			Content: comment_text}
		db.Create(&comment)

	} else if action_type == "2" { //删除评论
		var comment model.Comment
		db.Where("Id=?", comment_id).Delete(&comment)
	} else { //传入action_type非法
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "action_type invalid"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	token := c.Query("token")       //用户鉴权token
	video_id := c.Query("video_id") //视频id
	//判断token
	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "token invalid"})
	}
	//建立数据库连接
	username := "root"  //账号
	password := "123"   //密码
	host := "127.0.0.1" //数据库地址，可以是Ip或者域名
	port := 3306        //数据库端口
	Dbname := "douyin"  //数据库名
	timeout := "10s"    //连接超时，10秒
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s",
		username, password, host, port, Dbname, timeout)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//数据库连接出问题
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "database err"})
	}
	var comment []Res_Comment
	db.Model(&model.Comment{}).Where("VideoId=?", strconv.ParseInt(video_id, 10, 64)).Find(comment)
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: comment,
	})
}
