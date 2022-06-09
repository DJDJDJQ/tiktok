package model

import (
	"time"
)

type User struct {
	Id            int64     `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	Password      string    `json:"password,omitempty"`
	FollowCount   int64     `json:"follow_count,omitempty"`
	FollowerCount int64     `json:"follower_count,omitempty"`
	IsFollow      bool      `json:"is_follow,omitempty"`
	RegisterTime  time.Time `json:"register_time,omitempty"`
}

type Video struct {
	Id            int64     `json:"id,omitempty"`
	UserId        int64     `json:"user_id,omitempty"`
	PlayUrl       string    `json:"play_url,omitempty"`
	CoverUrl      string    `json:"cover_url,omitempty"`
	FavoriteCount int64     `json:"favorite_count,omitempty"`
	CommentCount  int64     `json:"comment_count,omitempty"`
	IsFavorite    bool      `json:"is_favorite,omitempty"`
	Title         string    `json:"title,omitempty"`
	PublishTime   time.Time `json:"publish_time,omitempty"`
	Status        int8      `json:"status,omitempty"`
}

type Favorite struct {
	Id         int64     `json:"id,omitempty"`
	UserId     int64     `json:"user_id,omitempty"`
	VideoId    int64     `json:"video_id,omitempty"`
	CreateTime time.Time `json:"create_time,omitempty"`
	IsCancel   bool      `json:"is_cancel,omitempty"` //DEFAULT NULL,
	CancelTime time.Time `json:"cancel_time,omitempty"`
}

type Comment struct {
	Id         int64     `json:"id,omitempty"`
	UserId     int64     `json:"user_id,omitempty"`
	VideoId    int64     `json:"video_id,omitempty"`
	Content    string    `json:"content,omitempty"`
	CreateDate time.Time `json:"create_date,omitempty"`
	IsDelete   bool      `json:"is_delete,omitempty"`
	DeleteTime time.Time `json:"delete_time,omitempty"`
}

type Follow struct {
	Id         int64     `json:"id,omitempty"`
	UserId     int64     `json:"user_id,omitempty"`
	FollowId   int64     `json:"follow_id,omitempty"`
	CreateTime time.Time `json:"create_time,omitempty"`
	IsCancel   bool      `json:"is_cancel,omitempty"`
	CancelTime time.Time `json:"cancel_time,omitempty"`
}

//设置表名，可以通过给User struct类型定义 TableName函数，返回一个字符串作为表名
//不重构TableName方法，gorm会默认将struct类型名加s作为表名 例如User结构体 查询users表而不会查tb_user表
func (v User) TableName() string {
	return "tb_user"
}
func (v Video) TableName() string {
	return "tb_video"
}
func (v Favorite) TableName() string {
	return "tb_favorite"
}
func (v Comment) TableName() string {
	return "tb_comment"
}
func (v Follow) TableName() string {
	return "tb_follow"
}
