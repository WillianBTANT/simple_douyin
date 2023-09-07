package model

import (
	"gorm.io/gorm"
)
import dbPackage "simple-douyin/db"

var db *gorm.DB

func init() {
	db = dbPackage.DB
	err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
		AutoMigrate(&User{}, &Comment{}, &Video{}, &UserUser{}, &UserVideo{})
	if err != nil {
		panic(err)
	}
}

func CreateTables() error {
	return db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
		AutoMigrate(&User{}, &Comment{}, &Video{}, &UserUser{}, &UserVideo{})
}

func SelectUserByID(userID int64) (*User, error) {
	var user User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
