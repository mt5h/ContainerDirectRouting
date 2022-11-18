package main

import (
	"fmt"
	"log"
	"spawner/build"
	"spawner/controllers"
	"spawner/utils"

	"github.com/gin-gonic/gin"
)

func main() {

	utils.LoadFlags()
	build.LoadInfo()

	if utils.PathRouting == utils.CookieRouting {
		log.Fatal("Set only one type of routing at the time")
	}

	router := gin.Default()

	// management endpoint
	provisioning_v1 := router.Group("/deploy")
	{
		provisioning_v1.POST("/", controllers.CreateContainer)
		provisioning_v1.GET("/", controllers.ListContainers)
		provisioning_v1.PUT("/:container-id", controllers.StartContainer)
		provisioning_v1.DELETE("/:container-id", controllers.DeleteContainer)
	}

	// Path routing
	if utils.PathRouting {
		spawned := router.Group(fmt.Sprintf("/%s", utils.ContainerPrefix))
		{
			spawned.GET("/:container-name/*other", controllers.PathRouting)
			spawned.GET("/:container-name", controllers.PathRouting)
		}
	}

	if utils.CookieRouting {
		spawner := router.Group("/")
		{
			spawner.GET("/:containerName", controllers.CookieRouting)
		}
	}

	router.Run(":8008")
}
