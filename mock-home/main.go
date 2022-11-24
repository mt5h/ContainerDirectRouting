package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mock-home/build"
	"net/http"
)

func homeCtrl(c *gin.Context){
		c.JSON(200, gin.H{"message": "You are at Home"})
}

func instanceCtrl(c *gin.Context){

		containerName := c.Param("containerName")

		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}

		// remove every path from the request
		redirectUrl := fmt.Sprintf("%s://%s", scheme, c.Request.Host)
		// set a custom cookie use by traefik
		c.SetCookie("instance", containerName, 3600, "/", c.Request.URL.Hostname(), false, false)
		c.Redirect(http.StatusFound, redirectUrl)

}

func main() {

	build.LoadInfo()
	r := gin.Default()
  home := r.Group("/")
  {
    home.GET("home", homeCtrl)
    home.GET("home/", homeCtrl)
	  home.GET("/home/:containerName", instanceCtrl)
  }
	r.Run(":8000")

}
