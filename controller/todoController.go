package controller

import (
	"TODOList/common"
	"TODOList/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strings"
)

// 查询全部的todolist,GET
func GetAllController(c *gin.Context) {
	dbtest := common.GetTestDB()
	user, err := c.Get("user")
	if !err {
		c.JSON(http.StatusNotFound, model.Resp{
			Code: 50201,
		})
		return
	}

	if user, ok := user.(model.User); ok {
		email := user.Email

		var todos []model.TODO
		result := dbtest.Where("email = ?", email).Find(&todos).Error
		if result != nil {
			c.JSON(http.StatusNotFound, model.Resp{
				Code: 40401,
			})
		}

		sort.Slice(todos, func(i, j int) bool {
			return todos[i].Id < todos[j].Id
		})

		c.JSON(http.StatusOK, gin.H{
			"code":   20000,
			"result": todos,
		})

	} else {
		c.JSON(http.StatusBadGateway, model.Resp{
			Code: 50202,
		})
	}
}

// 查询给定的id,GET
func GetIdController(c *gin.Context) {
	dbtest := common.GetTestDB()
	id := c.Param("id") //注意，返回值是 string

	user, err2 := c.Get("user")
	if !err2 {
		c.JSON(http.StatusNotFound, model.Resp{
			Code: 40401,
		})
		return
	}

	if user, ok := user.(model.User); ok {
		email := user.Email

		var todos []model.TODO
		result := dbtest.Where("email = ? AND id = ?", email, id).Find(&todos).Error //从数据库中找出来
		if result != nil {
			c.JSON(http.StatusNotFound, model.Resp{
				Code: 40402,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"Code":   20000,
			"result": todos,
		})

	} else {
		c.JSON(http.StatusBadGateway, model.Resp{
			Code: 50201,
		})
	}
}

// 上传,POST
func PostController(c *gin.Context) {
	dbtest := common.GetTestDB()

	request := model.TODO{}
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
		var single model.TODO
		dbtest.First(&single, "id = ? AND email = ?", Newid, email)

		//如果此处已经存了东西，返回报错和以前的东西
		if single.Exist == true {
			c.JSON(http.StatusBadGateway, model.Resp{
				Code: 50201,
				Data: single,
			})
			return
		}

		//否则就添加

		err := dbtest.Create(&request).Error
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"code":  50202,
				"error": err,
			})
			return
		}

		c.JSON(http.StatusOK, model.Resp{
			Code: 20000,
			Data: request,
		})

	} else {
		c.JSON(http.StatusBadGateway, model.Resp{
			Code: 50203,
		})
	}
}

// 以id为索引更新
func PutIdController(c *gin.Context) {
	dbtest := common.GetTestDB()

	request := model.TODO{}
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
		var single model.TODO
		dbtest.First(&single, "id = ? AND email = ?", Newid, email)

		if single.Exist == false {
			//没有就添加
			err := dbtest.Create(&request).Error
			if err != nil {
				c.JSON(http.StatusBadGateway, model.Resp{
					Code: 50201,
				})
				return
			}

			c.JSON(http.StatusOK, model.Resp{
				Code: 20000,
				Data: request,
			})
			return
		}

		single = request

		dbtest.Save(single)

		c.JSON(http.StatusOK, model.Resp{
			Code: 20000,
			Data: request,
		})

	} else {
		c.JSON(http.StatusBadGateway, model.Resp{
			Code: 50202,
		})
	}
}

// 以id为索引删除
func DelIdController(c *gin.Context) {
	dbtest := common.GetTestDB()

	id := c.Param("id") //注意，返回值是 string

	user, err2 := c.Get("user")
	if !err2 {
		c.JSON(http.StatusNotFound, model.Resp{
			Code: 40401,
		})
		return
	}

	if user, ok := user.(model.User); ok {
		email := user.Email

		var single model.TODO
		result := dbtest.First(&single, "id = ? AND email = ?", id, email).Error
		if result != nil {
			c.JSON(http.StatusNotFound, model.Resp{
				Code: 40402,
			})
			return
		}

		dbtest.Delete(single)
		c.JSON(http.StatusOK, model.Resp{
			Code: 20000,
		})

	} else {
		c.JSON(http.StatusBadGateway, model.Resp{
			Code: 50201,
		})
	}
}

// 以 tag 为索引查询
func GetTagController(c *gin.Context) {
	dbtest := common.GetTestDB()

	tag := c.Param("tag") //注意，返回值是 string

	user, err2 := c.Get("user")
	if !err2 {
		c.JSON(http.StatusNotFound, model.Resp{
			Code: 40401,
		})
		return
	}

	if user, ok := user.(model.User); ok {
		email := user.Email

		var todos []model.TODO
		result := dbtest.Where("email = ? AND tag = ?", email, tag).Find(&todos).Error
		if result != nil {
			c.JSON(http.StatusNotFound, model.Resp{
				Code: 40402,
			})
		}

		sort.Slice(todos, func(i, j int) bool {
			return todos[i].Id < todos[j].Id
		})

		c.JSON(http.StatusOK, gin.H{
			"Code":   20000,
			"result": todos,
		})

	} else {
		c.JSON(http.StatusBadGateway, model.Resp{
			Code: 50201,
		})
	}
}

// 清理 已完成的
func PostClearController(c *gin.Context) {
	dbtest := common.GetTestDB()

	user, err := c.Get("user")
	if !err {
		c.JSON(http.StatusNotFound, model.Resp{
			Code: 40401,
		})
		return
	}

	if user, ok := user.(model.User); ok { //硬抄,我也不知道为什么断言
		email := user.Email

		var todos []model.TODO
		result := dbtest.Where("email = ?", email).Find(&todos).Error
		if result != nil {
			c.JSON(http.StatusNotFound, model.Resp{
				Code: 40402,
			})
		}

		for _, todo := range todos {
			if todo.Done == true {
				dbtest.Delete(todo)
			}
		}

		var todos2 []model.TODO
		result2 := dbtest.Where("email = ?", email).Find(&todos2).Error
		if result2 != nil {
			c.JSON(http.StatusNotFound, model.Resp{
				Code: 40403,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"code":   20000,
			"result": todos2,
		})

	} else {
		c.JSON(http.StatusBadGateway, model.Resp{
			Code: 50201,
		})
	}
}

// 发邮件(又臭又长)
func GetSendController(c *gin.Context) {
	dbtest := common.GetTestDB()

	user, err := c.Get("user")
	if !err {
		c.JSON(http.StatusNotFound, model.Resp{
			Code: 40401,
		})
		return
	}

	if user, ok := user.(model.User); ok {
		email := user.Email
		var objc model.MailUser

		objc.Source = "1!5!"
		objc.Contacts = email
		objc.Subject = "TODOList-Warning"

		var flag bool

		var todos []model.TODO
		result := dbtest.Where("email = ?", email).Find(&todos).Error
		if result != nil {
			c.JSON(http.StatusNotFound, model.Resp{
				Code: 40402,
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
		if source != "1!5!" {
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
		err := SendToMail(user, sendUserName, password, host, to, subject, body, "html")
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
		c.JSON(http.StatusBadGateway, model.Resp{
			Code: 50201,
		})
	}
}
