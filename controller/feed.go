package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mod/model"
	"go.mod/utils"
)

type FeedResponse struct {
	Response
	VideoList []Res_Video `json:"video_list,omitempty"`
	NextTime  int64       `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	// c.JSON(http.StatusOK, FeedResponse{
	// 	Response:  Response{StatusCode: 0},
	// 	VideoList: DemoVideos,
	// 	NextTime:  time.Now().Unix(),
	// })
	token := c.Query("token")
	latest_time := c.Query("latest_time")
	timeobj, _ := time.Parse("2006-01-02 15:04:05", latest_time)
	if timeobj.IsZero() || time.Since(timeobj) < 0 {
		// 说明传入时间为空或者传入时间大于当前时间，那么赋值为当前时间
		timeobj = time.Now()
	}
	/*	latest_time, err := strconv.Atoi(c.Query("latest_time"))
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Timestamp parsing error",
			})
			return
		}

	var timeobj time.Time*/
	var nexttime int64 //本次返回的视频中发布最早的时间，作为下次请求的latest_time
	Videos := []model.Video{}

	/*	if latest_time != 0 { //可选参数，不填默认当前时间
			timeobj = time.Unix(int64(latest_time), 0)
		} else {
			timeobj = time.Now()
		}*/

	model.Mysql.Table("tb_video").Limit(30).Where("publish_time < ?", timeobj).Order("publish_time desc").Find(&Videos)

	if len(Videos) > 0 {
		lastvideo := Videos[len(Videos)-1]
		nexttime = lastvideo.PublishTime.Unix()
	} else { //如果查询不到符合要求的视频，nexttime设置为当前时间，不确定逻辑是否正确
		nexttime = time.Now().Unix()
	}

	sendVideos := []Res_Video{}
	if token == "" { // 未登录
		for _, value := range Videos {
			author := model.User{}
			model.Mysql.Table("tb_user").Where("id = ?", value.UserId).Find(&author)
			res_author := Res_User{
				Id:            author.Id,
				Name:          author.Name,
				FollowCount:   author.FollowCount,
				FollowerCount: author.FollowerCount,
				IsFollow:      false,
				//IsFollow:      author.IsFollow,
			}

			temp := Res_Video{
				Id:            value.Id,
				Author:        res_author,
				Title:         value.Title,
				PlayUrl:       value.PlayUrl,
				CoverUrl:      value.CoverUrl,
				FavoriteCount: value.FavoriteCount,
				CommentCount:  value.CommentCount,
				IsFavorite:    false,
			}
			sendVideos = append(sendVideos, temp)
		}
	} else { // 已登录
		claims, _ := utils.ParseToken(token)
		for _, value := range Videos {
			author := model.User{}
			// 改动
			follow := model.Follow{}
			isfollow := true
			if claims.UserId != value.UserId {
				res := model.Mysql.Table("tb_follow").Where("user_id = ? and follow_id = ?", claims.UserId, value.UserId).Find(&follow)
				if res.RowsAffected == 0 {
					isfollow = false
				}
			}
			//
			model.Mysql.Table("tb_user").Where("id = ?", value.UserId).Find(&author)
			res_author := Res_User{
				Id:            author.Id,
				Name:          author.Name,
				FollowCount:   author.FollowCount,
				FollowerCount: author.FollowerCount,
				IsFollow:      isfollow, //
			}
			favorite := model.Favorite{}
			isfavor := false
			res := model.Mysql.Table("tb_favorite").Where("user_id = ? AND video_id = ?", claims.UserId, value.Id).Find(&favorite)
			if res.RowsAffected != 0 {
				isfavor = true
			}
			temp := Res_Video{
				Id:            value.Id,
				Author:        res_author,
				Title:         value.Title,
				PlayUrl:       value.PlayUrl,
				CoverUrl:      value.CoverUrl,
				FavoriteCount: value.FavoriteCount,
				CommentCount:  value.CommentCount,
				IsFavorite:    isfavor,
			}
			sendVideos = append(sendVideos, temp)
		}
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: sendVideos,
		NextTime:  nexttime,
	})
}
