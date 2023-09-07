package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple-douyin/db"
	"simple-douyin/model"
	"strconv"
	"time"
)

var (
	maxVideoNum = 30
)

func Feed(c *gin.Context) {
	currentUserID, _ := c.Get("user_id")
	lastTimeStr := c.Query("latest_time")

	lastTime, _ := strconv.ParseInt(lastTimeStr, 10, 64)

	videoArr := make([]model.Video, 30)
	if currentUserID == nil {
		query := db.DB.Table("video").Preload("Author").
			Order("video.create_at DESC").
			Limit(maxVideoNum)
		if lastTime > 0 {
			query = query.Where("video.create_at < ?", lastTime)
		}
		if err := query.Find(&videoArr).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
			return
		}
	} else {
		query := db.DB.Table("video").
			Preload("Author").
			Preload("Author.Followers", "follower = ? ", currentUserID). // 加载用户的粉丝
			Preload("FavoriteUser", "user_id=?", currentUserID).         // 加载用户喜欢视频字段
			Select("video.*").
			Order("video.create_at DESC").
			Limit(maxVideoNum)
		if lastTime > 0 {
			query = query.Where("video.create_at > ?", lastTime)
		}

		if err := query.
			Find(&videoArr).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
			return
		}
	}

	for i := range videoArr {
		video := &videoArr[i]
		video.PlayURL = staticBaseUrl + video.PlayURL
		if video.FavoriteUser != nil { // 用户是否喜欢这个视频
			video.IsFavorite = video.FavoriteUser.Flag
		}
		if video.Author != nil && len(video.Author.Followers) > 0 { // 用户是否关注作者
			video.Author.IsFollow = video.Author.Followers[0].Flag
		}
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videoArr,
		NextTime:  time.Now().Unix(),
	})
}
