package middleware

import (
	"TODOList/common"
	"TODOList/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := c.GetHeader("Authorization") //获得

		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    40100,
				"message": "权限不足1",
			})
			return
		}

		tokenString = tokenString[7:] // 去头

		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    40100,
				"message": "权限不足2",
			})
			return
		}

		// 验证通过后取得 userId
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId) //利用数据库开找

		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    40100,
				"message": "权限不足3",
			})
			return
		}

		c.Set("user", user) //以后通过 c.Get(key) 方法获取。
		c.Next()

	}
}
