package common

import (
	"TODOList/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var UserDB *gorm.DB
var TestDB *gorm.DB

// User
func InitUserDB() *gorm.DB {
	dbuser, err1 := gorm.Open(sqlite.Open("C:\\learn\\CS\\codes\\go\\TODOList\\db\\user.db"), &gorm.Config{})
	if err1 != nil {
		panic("failed to connect database")
	}
	_ = dbuser.AutoMigrate(&model.User{})
	UserDB = dbuser
	return dbuser
}

func GetUserDB() *gorm.DB {
	return UserDB
}

// TODOList
func InitTestDB() *gorm.DB {
	dbtest, err1 := gorm.Open(sqlite.Open("C:\\learn\\CS\\codes\\go\\TODOList\\db\\test.db"), &gorm.Config{})
	if err1 != nil {
		panic("failed to connect database")
	}
	_ = dbtest.AutoMigrate(&model.TODO{})
	TestDB = dbtest
	return dbtest
}

func GetTestDB() *gorm.DB {
	return TestDB
}
