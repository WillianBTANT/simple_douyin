package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple-douyin/model"
	"simple-douyin/utils"
)

func authenticateUser(userID int64) (*model.User, error) {
	return model.SelectUserByID(userID)
}

// TokenAuthMiddleware
//
// authNeeded=true:该接口需要token才能访问;false:该接口的token为非必须参数
func TokenAuthMiddleware(authNeeded bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData struct {
			Token string `json:"token" form:"token"`
		}
		switch c.Request.Method {
		case http.MethodGet:
			requestData.Token = c.Query("token")
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			if err := c.ShouldBind(&requestData); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
				return
			}
		default:
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
			return
		}
		token := requestData.Token

		if authNeeded {
			if token == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing token"})
				c.Abort()
				return
			}

			userID, err := utils.ParseToken(token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			}

			user, err := authenticateUser(userID)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			}

			c.Set("user_id", userID)
			c.Set("user", user)
		} else if token != "" {
			userID, err := utils.ParseToken(token)
			if err == nil {
				if user, err := authenticateUser(userID); err == nil {
					c.Set("user_id", userID)
					c.Set("user", user)
				}
			}
		}

		c.Next()
	}
}
