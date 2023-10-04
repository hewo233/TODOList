package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"varchar(20);not null"`
	Email    string `gorm:"varchar(20);not null;unique"`
	Password string `gorm:"size:255;not null"`
}

type MailUser struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Source   string `form:"source" json:"source" uri:"source" xml:"source" binding:"required"`
	Contacts string `form:"contacts" json:"contacts" uri:"contacts" xml:"contacts" binding:"required"`
	Subject  string `form:"subject" json:"subject" uri:"subject" xml:"subject" binding:"required"`
	Content  string `form:"content" json:"content" uri:"content" xml:"content" binding:"required"`
}
