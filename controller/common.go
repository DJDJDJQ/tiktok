package controller

import "go.mod/model"

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
	Id         int64    `json:"id"`
	User       Res_User `json:"user"`
	Content    string   `json:"content"`
	CreateDate string   `json:"create_date"`
}

type Res_User struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow"`
}

// 将user转换成Res_User
func Convert2ResUser(user model.User, isfollow bool) Res_User {
	var author Res_User = Res_User{
		Id:            user.Id,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      isfollow,
	}
	return author
}

// 将video转换成Res_Video
func Convert2ResVideo(video model.Video, user Res_User, isfavor bool) Res_Video {
	var rsp_video Res_Video = Res_Video{
		Id:            video.Id,
		Author:        user,
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
		IsFavorite:    isfavor,
		Title:         video.Title,
	}
	return rsp_video
}

// 将comment转换成Res_Comment
func Convert2Comment(comment model.Comment, user Res_User) Res_Comment {
	var rsp_comment Res_Comment = Res_Comment{
		Id:         comment.Id,
		User:       user,
		Content:    comment.Content,
		CreateDate: comment.CreateDate.Format("01-02"),
	}
	return rsp_comment
}
