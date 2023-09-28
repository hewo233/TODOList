package main

import (
	"TODOList/common"
	controller "TODOList/controller"
	"TODOList/middleware"
	"TODOList/model"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	//"strconv"
)

// todolist的结构体
type TODO struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Done      bool   `json:"done"`
	Exist     bool   `json:"exist"`
	Telephone string `json:"telephone"` //对应每个客户
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

		if user, ok := user.(model.User); ok { //硬抄,我也不知道为什么断言
			telephone := user.Telephone

			var todos []TODO
			result := db.Where("telephone = ?", telephone).Find(&todos).Error
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

		if user, ok := user.(model.User); ok { //硬抄,我也不知道为什么断言
			telephone := user.Telephone

			var todos []TODO
			result := db.Where("telephone = ? AND id = ?", telephone, id).Find(&todos).Error
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

		if user, ok := user.(model.User); ok { //硬抄,我也不知道为什么断言
			telephone := user.Telephone

			request.Telephone = telephone

			Newid := request.Id
			var single TODO
			db.First(&single, "id = ? AND telephone = ?", Newid, telephone)

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
	//是否需要考虑没有此id的时候直接新建？目前的选择，同时还有就是报错的选择被我注释掉了
	engine.PUT("/todo", middleware.AuthMiddleware(), func(c *gin.Context) {
		request := TODO{}
		c.BindJSON(&request)

		user, err2 := c.Get("user")
		if !err2 {
			c.JSON(http.StatusTemporaryRedirect, gin.H{"error": user})
			return
		}

		if user, ok := user.(model.User); ok { //硬抄,我也不知道为什么断言
			telephone := user.Telephone

			request.Telephone = telephone

			Newid := request.Id
			var single TODO
			db.First(&single, "id = ? AND telephone = ?", Newid, telephone)

			if single.Exist == false {
				//没有就添加
				err := db.Create(&request).Error
				if err != nil {
					c.JSON(http.StatusBadGateway, Resp{
						Code: 50200,
					})
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

		if user, ok := user.(model.User); ok { //硬抄,我也不知道为什么断言
			telephone := user.Telephone

			var single TODO
			result := db.First(&single, "id = ? AND telephone = ?", id, telephone).Error
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

	//以 title 为索引查询
	engine.GET("/todo/title/:title", middleware.AuthMiddleware(), func(c *gin.Context) {
		title := c.Param("title") //注意，返回值是 string

		user, err2 := c.Get("user")
		if !err2 {
			c.JSON(http.StatusNotFound, Resp{
				Code: 40400,
			})
			return
		}

		if user, ok := user.(model.User); ok { //硬抄,我也不知道为什么断言
			telephone := user.Telephone

			var todos []TODO
			result := db.Where("telephone = ? AND title = ?", telephone, title).Find(&todos).Error
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
			telephone := user.Telephone

			var todos []TODO
			result := db.Where("telephone = ?", telephone).Find(&todos).Error
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
			result2 := db.Where("telephone = ?", telephone).Find(&todos2).Error
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

	engine.Run(":8080")

}
