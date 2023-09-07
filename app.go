package main

import (
	"github.com/gin-gonic/gin"
	"simple-douyin/controller"
	"simple-douyin/middleware"
)

func createAndInitEngine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	middleware.InitGinLogWriter() // 设置日志输入到哪里
	//r.Use(gin.LoggerWithFormatter(middleware.LogToFileFormatter()))  // 使用自己的formatter
	r.Use(gin.Logger())
	initRouter(r)
	return r
}

func initRouter(r *gin.Engine) {
	// 公共文件夹文件目录
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// 不需要认证的API添加到 apiRouter
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)

	// token 作为参数但是不是非必须参数的及接口
	aApiRouter := apiRouter.Group("/", middleware.TokenAuthMiddleware(false))
	aApiRouter.GET("/feed/", controller.Feed)

	// 使用中间件做用户认证，需要token的API请添加到protectedApiRouter中
	protectedApiRouter := apiRouter.Group("/", middleware.TokenAuthMiddleware(true))
	protectedApiRouter.GET("/user/", controller.UserInfo)
	protectedApiRouter.GET("/publish/list/", controller.PublishList)
	protectedApiRouter.POST("/publish/action/", controller.Publish)

	//// extra apis - I
	protectedApiRouter.POST("/comment/action/", controller.CommentAction)   // 评论操作
	protectedApiRouter.GET("/comment/list/", controller.CommentList)        // 评论列表
	protectedApiRouter.POST("/favorite/action/", controller.FavoriteAction) // 赞操作
	protectedApiRouter.GET("/favorite/list/", controller.FavoriteList)      // 喜欢列表

	//// extra apis - II
	//apiRouter.POST("/relation/action/", controller.RelationAction)
	//apiRouter.GET("/relation/follow/list/", controller.FollowList)
	//apiRouter.GET("/relation/follower/list/", controller.FollowerList)
	//apiRouter.GET("/relation/friend/list/", controller.FriendList)
	//apiRouter.GET("/message/chat/", controller.MessageChat)
	//apiRouter.POST("/message/action/", controller.MessageAction)
}
