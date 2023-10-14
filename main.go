package main

import (
	"TODOList/common"
	"TODOList/routes"
	"github.com/gin-gonic/gin"
	//"strconv"
)

func main() {

	engine := gin.Default()

	_ = common.InitUserDB() //init
	_ = common.InitTestDB()

	routes.SetupRouter(engine)
	engine.Run(":8080")
}
