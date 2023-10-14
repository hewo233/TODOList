package routes

import (
	"TODOList/controller"
	"TODOList/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(engine *gin.Engine) {

	//登录版
	engine.POST("/register", controller.Register)
	engine.POST("/login", controller.Login)
	engine.GET("/userinfo", middleware.AuthMiddleware(), controller.Info)

	//TODOList版
	engine.GET("/todo/list", middleware.AuthMiddleware(), controller.GetAllController)
	engine.GET("/todo/id/:id", middleware.AuthMiddleware(), controller.GetIdController)
	engine.POST("/todo", middleware.AuthMiddleware(), controller.PostController)
	engine.PUT("/todo/", middleware.AuthMiddleware(), controller.PutIdController)
	engine.DELETE("/todo/:id", middleware.AuthMiddleware(), controller.DelIdController)
	engine.GET("/todo/tag/:tag", middleware.AuthMiddleware(), controller.GetTagController)
	engine.POST("/todo/clear", middleware.AuthMiddleware(), controller.PostClearController)

	//邮件
	engine.GET("/todo/send", middleware.AuthMiddleware(), controller.GetSendController)

}
