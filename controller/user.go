package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"simple-douyin/db"
	"simple-douyin/model"
	"simple-douyin/utils"
	"strconv"
)

type UsrnPwdRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type UserTokenResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token,omitempty"`
}

func Login(c *gin.Context) {
	var loginRequest UsrnPwdRequest
	if err := c.ShouldBindQuery(&loginRequest); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "parameter mismatch"})
		return
	}

	// 用户名和密码的校验
	if !utils.IsValidUsername(loginRequest.Username) {
		c.JSON(http.StatusBadGateway, gin.H{"error": "username format is incorrect"})
		return
	}
	if !utils.IsValidPassword(loginRequest.Password) {
		c.JSON(http.StatusBadGateway, gin.H{"error": "password format is incorrect"})
		return
	}

	var loginUser model.User
	if err := db.DB.Where(&model.User{Name: loginRequest.Username}).First(&loginUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, UserTokenResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query error"})
		return
	}

	if err := utils.CompareHashAndPassword([]byte(loginUser.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "wrong password"})
		return
	}

	token, err := utils.GenerateToken(&loginUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation error"})
		return
	}

	c.JSON(http.StatusOK, UserTokenResponse{
		Response: Response{StatusCode: 0, StatusMsg: "ok"},
		UserId:   loginUser.ID,
		Token:    token,
	})
}

func Register(c *gin.Context) {
	var registerRequest UsrnPwdRequest
	if err := c.ShouldBindQuery(&registerRequest); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "parameter mismatch"})
		return
	}

	// 用户名和密码的校验
	if !utils.IsValidUsername(registerRequest.Username) {
		c.JSON(http.StatusBadGateway, gin.H{"error": "username format is incorrect"})
		return
	}
	if !utils.IsValidPassword(registerRequest.Password) {
		c.JSON(http.StatusBadGateway, gin.H{"error": "password format is incorrect"})
		return
	}

	hashedPassword, err := utils.GenerateFromPassword([]byte(registerRequest.Password))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	newUser := model.User{Name: registerRequest.Username, Password: string(hashedPassword)}
	result := db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&newUser)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "user already exists"})
		return
	}

	token, err := utils.GenerateToken(&newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, UserTokenResponse{
		Response: Response{StatusCode: 0, StatusMsg: "ok"},
		UserId:   newUser.ID,
		Token:    token,
	})
	return
}

type UserInfoResponse struct {
	Response
	User model.User `json:"user,omitempty"`
}

func UserInfo(c *gin.Context) {
	targetUserID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}
	currentUserID := c.GetInt64("user_id")

	var targetUser model.User
	if err := db.DB.
		Preload("Followers", "follower = ? ", currentUserID).
		First(&targetUser, targetUserID).Error; err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User not found"})
		return
	}
	if len(targetUser.Followers) > 0 {
		targetUser.IsFollow = targetUser.Followers[0].Flag
	}

	c.JSON(http.StatusOK, UserInfoResponse{
		Response: Response{StatusCode: 0, StatusMsg: "ok"},
		User:     targetUser,
	})
}
