package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"spawner/controllers"
	"spawner/utils"
)

func main() {

	utils.LoadFlags()

	router := gin.Default()
	provisioning_v1 := router.Group("/deploy")
	{
		provisioning_v1.POST("/", controllers.CreateContainer)
		provisioning_v1.GET("/", controllers.ListContainers)
		provisioning_v1.PUT("/:container-id", controllers.StartContainer)
		provisioning_v1.DELETE("/:container-id", controllers.DeleteContainer)
	}

	spawned := router.Group(fmt.Sprintf("/%s", utils.ContainerPrefix))
	{
		spawned.GET("/:container-name/*other", controllers.RestartContainer)
	}
	router.Run(":8008")
}
