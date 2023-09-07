package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"simple-douyin/config"
)

var DB *gorm.DB

func init() {
	dsn := config.GetConfig().DatabaseURI
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{PrepareStmt: true})
	if err != nil {
		panic(err)
	}
	// todo 出现错误的时候不再打印日志，后期考虑添加到应用程序日志
	db.Logger = logger.Default.LogMode(logger.Silent)
	DB = db
}
