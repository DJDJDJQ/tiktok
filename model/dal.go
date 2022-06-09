package model

import (
	"database/sql/driver"
	"time"
)

type MyBool bool

func (b MyBool) Value() (driver.Value, error) {
	result := make([]byte, 1)
	if b {
		result[0] = byte(1)
	} else {
		result[0] = 0
	}
	return result, nil
}
func (b MyBool) Scan(v interface{}) error {
	bytes := v.([]byte)

	if bytes[0] == 0 {
		b = false
	} else {
		b = true
	}
	return nil
}

type User struct {
	Id            int64     `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	FollowCount   int64     `json:"follow_count,omitempty"`
	FollowerCount int64     `json:"follower_count,omitempty"`
	IsFollow      MyBool    `json:"is_follow,omitempty"`
	RegisterTime  time.Time `json:"register_time,omitempty"`
}

type Video struct {
	Id            int64     `json:"id,omitempty"`
	UserId        int64     `json:"user_id,omitempty"`
	PlayUrl       string    `json:"play_url,omitempty"`
	CoverUrl      string    `json:"cover_url,omitempty"`
	FavoriteCount int64     `json:"favorite_count,omitempty"`
	CommentCount  int64     `json:"comment_count,omitempty"`
	IsFavorite    MyBool    `json:"is_favorite,omitempty"`
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
	CreateDate string    `json:"create_date,omitempty"`
	IsDelete   bool      `json:"is_delete,omitempty"`
	DeleteTime time.Time `json:"delete_time,omitempty"`
}

type Follow struct {
	Id         int64     `json:"id,omitempty"`
	UserId     int64     `json:"user_id,omitempty"`
	FollowerId int64     `json:"follower_id,omitempty"`
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
