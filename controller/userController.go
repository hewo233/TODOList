package controller

import (
	"TODOList/common"
	"TODOList/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func Register(c *gin.Context) {
	dbuser := common.GetDB()
	name := c.PostForm("name")
	password := c.PostForm("password")
	telephone := c.PostForm("telephone") //注意都是string

	if len(telephone) != 11 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    42200,
			"message": "错误的手机号",
		})
		return
	}
	if len(name) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    42200,
			"message": "你没名字啊？",
		})
		return
	}
	if len(password) < 6 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    42200,
			"message": "密码不能少于6位",
		})
		return
	}

	var user model.User
	dbuser.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    42200,
			"message": "手机号已被注册",
		})
		return
	}

	//创号
	HashPassword, err1 := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"message": "密码加密寄了",
		})
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(HashPassword),
	}
	err2 := dbuser.Create(&newUser).Error
	if err2 != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    50000,
			"message": "i don't know",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"message": "注册成功",
	})

}

// 登录
func Login(c *gin.Context) {
	dbuser := common.GetDB()
	var requestUser model.User
	c.Bind(&requestUser)
	telephone := requestUser.Telephone
	password := requestUser.Password

	//name := c.PostForm("name")

	if len(telephone) != 11 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    42200,
			"message": "手机号错误",
		})
		return
	}

	if len(password) < 6 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    42200,
			"message": "密码不能少于6位",
		})
		return
	}

	var user model.User
	dbuser.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    42200,
			"message": "用户不存在",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    42200,
			"message": "密码错误",
		})
		return
	}

	token, err := common.ReleaseToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"message": "系统错误",
		})
		log.Printf("Token generate error: %v", err) //打印错误日志
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": gin.H{
			"token": token,
		},
		"message": "登录成功",
	})

}

func Info(c *gin.Context) {
	user, _ := c.Get("user") //上文所说，返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user": user,
		},
	})
}
