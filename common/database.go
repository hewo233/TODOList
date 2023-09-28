package common

import (
	"TODOList/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dbuser, err1 := gorm.Open(sqlite.Open("C:\\learn\\CS\\codes\\go\\TODOList\\db\\user.db"), &gorm.Config{})
	if err1 != nil {
		panic("failed to connect database")
	}
	_ = dbuser.AutoMigrate(&model.User{})
	DB = dbuser
	return dbuser
}

func GetDB() *gorm.DB {
	return DB
}
