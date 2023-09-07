package controller

import (
	"net/http"
	"simple-douyin/db"
	"simple-douyin/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []model.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment model.Comment `json:"comment,omitempty"`
}

func CommentAction(c *gin.Context) {
	actionType := c.Query("action_type")
	videoID := c.Query("video_id")
	videoId, err := strconv.ParseInt(videoID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "video_id value error")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, "user_id not provided")
		return
	}

	userId, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusBadRequest, "user_id value error")
		return
	}

	switch actionType {
	case "1": // 发布评论
		content := c.Query("comment_text")
		if content == "" {
			c.JSON(http.StatusBadRequest, "comment_text cannot be empty")
			return
		}

		comment := model.Comment{UserID: userId, VideoID: videoId, Content: content}
		if err := db.DB.Create(&comment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0},
			Comment:  comment,
		})

	case "2": // 删除评论
		commentID := c.Query("comment_id")
		commentId, err := strconv.ParseInt(commentID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, "comment_id value error")
			return
		}

		if err := db.DB.Delete(&model.Comment{}, commentId).Error; err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, Response{StatusCode: 0})

	default:
		c.JSON(http.StatusBadRequest, "invalid action_type")
	}
}

func CommentList(c *gin.Context) {
	videoID := c.Query("video_id")
	videoId, err := strconv.ParseInt(videoID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "video_id value error")
		return
	}

	var comments []model.Comment
	if err := db.DB.Preload("User").
		Where("video_id = ?", videoId).
		Order("create_at DESC").
		Debug().
		Find(&comments).
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: comments,
	})
}
