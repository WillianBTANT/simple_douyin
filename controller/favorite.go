package controller

import (
	"github.com/gin-gonic/gin"
)

func FavoriteAction(c *gin.Context) {
	// todo 用户喜欢某个视频的操作
	// todo 视频表中 喜欢数+-1
	// todo user_video表中添加记录 否则FLag->false

}

func FavoriteList(c *gin.Context) {
	// todo 某个用户喜欢的视频列表
}
