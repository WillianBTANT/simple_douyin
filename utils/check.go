package utils

import "regexp"

const (
	usernamePattern = "^[a-zA-Z0-9]{5,64}$" // 英文数字长度为 5-64
	passwordPattern = "^[\x21-\x7E]{5,64}$" // ascii字符长度为 5-64
)

// IsValidUsername 检查用户名是否合法
func IsValidUsername(username string) bool {
	validUsernamePattern := regexp.MustCompile(usernamePattern)
	return validUsernamePattern.MatchString(username)
}

// IsValidPassword 检查用户密码是否合法
func IsValidPassword(password string) bool {
	validPasswordPattern := regexp.MustCompile(passwordPattern)
	return validPasswordPattern.MatchString(password)
}
