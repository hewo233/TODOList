package main

import (
	"TODOList/common"
	controller "TODOList/controller"
	"TODOList/middleware"
	"TODOList/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"strings"
	//"strconv"
)

// todolist的结构体
type TODO struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Done     bool   `json:"done"`
	Exist    bool   `json:"exist"`
	Email    string `json:"email"` //对应每个客户
	Tag      string `json:"tag"`
	Deadline string `json:"deadline"`
}

//var todos []*TODO

// 标准化响应
type Resp struct {
	Code int  `json:"code"`
	Data TODO `json:"data"`
}

func main() {

	_ = common.InitDB()

	db, err2 := gorm.Open(sqlite.Open("C:\\learn\\CS\\codes\\go\\TODOList\\db\\test.db"), &gorm.Config{})
	if err2 != nil {
		panic("failed to connect database")
	}
	_ = db.AutoMigrate(&TODO{})

	engine := gin.Default()

	//注册用户
	engine.POST("/register", controller.Register)

	//用户登录
	engine.POST("/login", controller.Login)

	//返回用户信息
	engine.GET("/userinfo", middleware.AuthMiddleware(), controller.Info)

	//查询全部的todolist
	engine.GET("/todo/list", middleware.AuthMiddleware(), func(c *gin.Context) {

		user, err := c.Get("user")
		if !err {
			c.JSON(http.StatusNotFound, Resp{
				Code: 40400,
			})
			return
		}

		if user, ok := user.(model.User); ok {
			email := user.Email

			var todos []TODO
			result := db.Where("email = ?", email).Find(&todos).Error
			if result != nil {
				c.JSON(http.StatusNotFound, Resp{
					Code: 40400,
				})
			}

			c.JSON(http.StatusOK, gin.H{
				"code":   20000,
				"result": todos,
			})

		} else {
			c.JSON(http.StatusBadGateway, Resp{
				Code: 50200,
			})
		}
	})

	//查询给定的id
	engine.GET("/todo/id/:id", middleware.AuthMiddleware(), func(c *gin.Context) {

		id := c.Param("id") //注意，返回值是 string

		user, err2 := c.Get("user")
		if !err2 {
			c.JSON(http.StatusNotFound, Resp{
				Code: 40400,
			})
			return
		}

		if user, ok := user.(model.User); ok {
			email := user.Email

			var todos []TODO
			result := db.Where("email = ? AND id = ?", email, id).Find(&todos).Error //从数据库中找出来
			if result != nil {
				c.JSON(http.StatusNotFound, Resp{
					Code: 40400,
				})
			}

			c.JSON(http.StatusOK, gin.H{
				"Code":   20000,
				"result": todos,
			})

		} else {
			c.JSON(http.StatusBadGateway, Resp{
				Code: 50200,
			})
		}

	})

	//上传
	engine.POST("/todo", middleware.AuthMiddleware(), func(c *gin.Context) {

		request := TODO{}
		c.BindJSON(&request)

		user, err2 := c.Get("user")
		if !err2 {
			c.JSON(http.StatusTemporaryRedirect, gin.H{"error": user})
			return
		}

		if user, ok := user.(model.User); ok {
			email := user.Email

			request.Email = email

			Newid := request.Id
			var single TODO
			db.First(&single, "id = ? AND email = ?", Newid, email)

			//如果此处已经存了东西，返回报错和以前的东西
			if single.Exist == true {
				c.JSON(http.StatusBadGateway, Resp{
					Code: 50200,
					Data: single,
				})
				return
			}

			//否则就添加

			err := db.Create(&request).Error
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{
					"code":  50200,
					"error": err,
				})
				return
			}

			c.JSON(http.StatusOK, Resp{
				Code: 20000,
				Data: request,
			})

		} else {
			c.JSON(http.StatusBadGateway, Resp{
				Code: 50200,
			})
		}

	})

	//以id为索引更新
	engine.PUT("/todo", middleware.AuthMiddleware(), func(c *gin.Context) {
		request := TODO{}
		c.BindJSON(&request)

		user, err2 := c.Get("user")
		if !err2 {
			c.JSON(http.StatusTemporaryRedirect, gin.H{"error": user})
			return
		}

		if user, ok := user.(model.User); ok {
			email := user.Email

			request.Email = email

			Newid := request.Id
			var single TODO
			db.First(&single, "id = ? AND email = ?", Newid, email)

			if single.Exist == false {
				//没有就添加
				err := db.Create(&request).Error
				if err != nil {
					c.JSON(http.StatusBadGateway, Resp{
						Code: 50200,
					})
					return
				}

				c.JSON(http.StatusOK, Resp{
					Code: 20000,
					Data: request,
				})
				return
			}

			single = request

			db.Save(single)

			c.JSON(http.StatusOK, Resp{
				Code: 20000,
				Data: request,
			})

		} else {
			c.JSON(http.StatusBadGateway, Resp{
				Code: 50200,
			})
		}
	})

	//以id为索引删除
	engine.DELETE("todo/:id", middleware.AuthMiddleware(), func(c *gin.Context) {

		id := c.Param("id") //注意，返回值是 string

		user, err2 := c.Get("user")
		if !err2 {
			c.JSON(http.StatusNotFound, Resp{
				Code: 40400,
			})
			return
		}

		if user, ok := user.(model.User); ok {
			email := user.Email

			var single TODO
			result := db.First(&single, "id = ? AND email = ?", id, email).Error
			if result != nil {
				c.JSON(http.StatusNotFound, Resp{
					Code: 40400,
				})
				return
			}

			db.Delete(single)
			c.JSON(http.StatusOK, Resp{
				Code: 20000,
			})

		} else {
			c.JSON(http.StatusBadGateway, Resp{
				Code: 50200,
			})
		}

	})

	//以 tag 为索引查询
	engine.GET("/todo/tag/:tag", middleware.AuthMiddleware(), func(c *gin.Context) {
		tag := c.Param("tag") //注意，返回值是 string

		user, err2 := c.Get("user")
		if !err2 {
			c.JSON(http.StatusNotFound, Resp{
				Code: 40400,
			})
			return
		}

		if user, ok := user.(model.User); ok {
			email := user.Email

			var todos []TODO
			result := db.Where("email = ? AND tag = ?", email, tag).Find(&todos).Error
			if result != nil {
				c.JSON(http.StatusNotFound, Resp{
					Code: 40400,
				})
			}

			c.JSON(http.StatusOK, gin.H{
				"Code":   20000,
				"result": todos,
			})

		} else {
			c.JSON(http.StatusBadGateway, Resp{
				Code: 50200,
			})
		}
	})

	//清理 已完成的
	engine.POST("/todo/clear", middleware.AuthMiddleware(), func(c *gin.Context) {

		user, err := c.Get("user")
		if !err {
			c.JSON(http.StatusNotFound, Resp{
				Code: 40400,
			})
			return
		}

		if user, ok := user.(model.User); ok { //硬抄,我也不知道为什么断言
			email := user.Email

			var todos []TODO
			result := db.Where("email = ?", email).Find(&todos).Error
			if result != nil {
				c.JSON(http.StatusNotFound, Resp{
					Code: 40400,
				})
			}

			for _, todo := range todos {
				if todo.Done == true {
					db.Delete(todo)
				}
			}

			var todos2 []TODO
			result2 := db.Where("email = ?", email).Find(&todos2).Error
			if result2 != nil {
				c.JSON(http.StatusNotFound, Resp{
					Code: 40400,
				})
			}

			c.JSON(http.StatusOK, gin.H{
				"code":   20000,
				"result": todos2,
			})

		} else {
			c.JSON(http.StatusBadGateway, Resp{
				Code: 50200,
			})
		}
	})

	engine.GET("/todo/send", middleware.AuthMiddleware(), func(c *gin.Context) {

		user, err := c.Get("user")
		if !err {
			c.JSON(http.StatusNotFound, Resp{
				Code: 40400,
			})
			return
		}

		if user, ok := user.(model.User); ok {
			email := user.Email
			var objc model.MailUser

			objc.Source = "heian"
			objc.Contacts = email
			objc.Subject = "TODOList-Warning"

			var flag bool

			var todos []TODO
			result := db.Where("email = ?", email).Find(&todos).Error
			if result != nil {
				c.JSON(http.StatusNotFound, Resp{
					Code: 40400,
				})
			}
			for _, todo := range todos {
				if todo.Done == false {
					flag = true
					break
				}
			}
			if flag == true {
				objc.Content = "你还有任务没做，速速完成"
			} else {
				objc.Content = "已完成所有任务，芜湖！"
			}

			//fmt.Println(json.Content, json.Contacts)
			//c.JSON(http.StatusOK, gin.H{"status": &json})
			user := "2842137843@qq.com"
			password := "zqntpnqgmrxcdgea"
			host := "smtp.qq.com:25"
			source := objc.Source
			if source != "heian" {
				fmt.Println("Send mail error!,source 认证失败")
				c.JSON(http.StatusOK, gin.H{
					"error": "Send mail error!,source 认证失败",
				})
				return
			}
			to := objc.Contacts
			if strings.TrimSpace(to) == "" {
				fmt.Println("Send mail error!,发送人为空")
				c.JSON(http.StatusOK, gin.H{
					"error": "Send mail error!,发送人为空",
				})
				return
			}
			subject := objc.Subject
			if strings.TrimSpace(subject) == "" {
				fmt.Println("Send mail error!标题为空")
				c.JSON(http.StatusOK, gin.H{
					"error": "Send mail error!,标题为空",
				})
				return
			}
			body := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="iso-8859-15">
			<title>MMOGA POWER</title>
		</head>
		<body>
			` + fmt.Sprintf(objc.Content) +
				`</body>
		</html>`
			//log.Printf("接收人：", json.Contacts+"\n"+"标题:", json.Subject+"\n", "发送内容：", json.Content+"\n")
			fmt.Printf("接收人:%s \n 标题: %s \n 内容: %s \n", objc.Contacts, objc.Subject, objc.Content)
			sendUserName := "TODOList" //发送邮件的人名称
			fmt.Println("send email")
			err := controller.SendToMail(user, sendUserName, password, host, to, subject, body, "html")
			if err != nil {
				fmt.Println("Send mail error!")
				c.JSON(http.StatusOK, gin.H{
					"error":   "Send mail error! !",
					"message": err,
				})
				//fmt.Println(err)
			} else {
				fmt.Println("Send mail success!")
				c.JSON(http.StatusOK, gin.H{
					"success": "Send mail success! !",
				})
			}

		} else {
			c.JSON(http.StatusBadGateway, Resp{
				Code: 50200,
			})
		}

	})

	engine.Run(":8080")

}
