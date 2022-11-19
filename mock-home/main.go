package main

import (
  "mock-home/build"
	"github.com/gin-gonic/gin"
)

func main() {

	build.LoadInfo()


  home := gin.Default()
  home.GET("/home", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "You are at Home",
          })
        })
  home.Run(":8000")
    
  }
