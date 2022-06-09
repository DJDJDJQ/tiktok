package controller

import "time"

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type User_Token struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type Res_Video struct {
	Id            int64    `json:"id"`
	Author        Res_User `json:"author"`
	PlayUrl       string   `json:"play_url"`
	CoverUrl      string   `json:"cover_url"`
	FavoriteCount int64    `json:"favorite_count"`
	CommentCount  int64    `json:"comment_count"`
	IsFavorite    bool     `json:"is_favorite"`
	Title         string   `json:"title"`
}

type Res_Comment struct {
	Id         int64     `json:"id"`
	User       Res_User  `json:"user"`
	Content    string    `json:"content"`
	CreateDate time.Time `json:"create_date"`
}

type Res_User struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow"`
}
