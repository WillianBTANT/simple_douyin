package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"time"
)

func InitGinLogWriter() {
	logDir := "./log"
	logFilePath := logDir + "/http.log"

	_, err := os.Stat(logFilePath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			panic(err)
		}
		create, err := os.Create(logFilePath)
		if err != nil {
			panic(err)
		}
		gin.DefaultWriter = create
	} else {
		file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			panic(err)
		}
		// gin.DefaultWriter = file  // todo 只打印到文件中，在程序允许的时候添加参数设置
		gin.DefaultWriter = io.MultiWriter(os.Stdout, file)
	}

}

func LogToFileFormatter() gin.LogFormatter {
	return func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[http] %s - [%s] \"%s %s %s | %s | %s | %s %d %s | %s \" || %s || \"%s \"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
			param.BodySize,
			param.Keys,
			param.Request.Header.Get("X-Request-Id"),
			param.Request.Header.Get("X-Forwarded-For"),
		)
	}
}
