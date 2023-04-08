package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
  "spawner/utils"
  "spawner/sessions"
	"spawner/build"
	"spawner/config"
	"spawner/controllers"
)


func main() {

	config.LoadFlags()
	build.LoadInfo()

	if config.PathRouting == config.CookieRouting {
		log.Fatal("Set only one type of routing at the time")
	}

	management := gin.Default()
  /////////////////////////////////////////////////////////////	

  sessions.TokenCache = utils.NewTokenSessions()
  sessions.TokenCache.SetValidity(config.TokenExpireTime)
  sessions.TokenCache.StartMaintenance(config.TokenCleanUpLoop)

  sessions.LsDB = utils.NewLoginStore()
	sessions.LsDB.ReadPasswordsFile(config.UsersPassFile)

	public := management.Group("/login")
  {
    public.POST("", controllers.Login)

  }

  /////////////////////////////////////////////////////////////
	// management endpoint
	provisioning := management.Group("/deploy")
  provisioning.Use(utils.SimpleAuth(sessions.TokenCache, config.EnableMgMtAuth))
	{
    
		provisioning.POST("/", controllers.CreateContainer)
		provisioning.GET("/", controllers.ListContainers)
		provisioning.PUT("/:container-id", controllers.StartContainer)
		provisioning.DELETE("/:container-id", controllers.DeleteContainer)
	}

	users := gin.Default()
	// Path routing
	if config.PathRouting {
		instances := users.Group(fmt.Sprintf("/%s", config.ContainerPrefix))
		{
			instances.GET("/:container-name/*any", controllers.PathRouting)
			instances.GET("/:container-name", controllers.PathRouting)
		}
	}

	// We have to handle:
	// - request with a cookie:
	//    - is valid -> container exists:
	//      - is started -> handled by traefick
	//      - is stopped -> container is started and the client is redirected to the root url <- spawner
	//    - is invalid -> container doesn't exist -> redirect to the fallback site <- spawner
	// - request has no cookie:
	//   - redirect to the fallback url <- spawner
	if config.CookieRouting {
		instances := users.Group("/")
		{
      //instances.GET("", controllers.CookieRouting)
			instances.GET("*any", controllers.CookieRouting)
		}
	}

	go management.Run(":8008")
	users.Run(":8000")

}
