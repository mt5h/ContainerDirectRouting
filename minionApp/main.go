package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	instance_env := os.Getenv("CONNSTR")
	hostname, _ := os.Hostname()
	r := gin.Default()

  instance :=	r.Group(hostname)
	{
		instance.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message":           "pong",
				"connection string": instance_env,
			})
		})
	}
  r.Run(":9000")
}
