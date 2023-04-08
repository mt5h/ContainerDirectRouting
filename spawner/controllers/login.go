package controllers

import(
  "spawner/sessions"
 	"github.com/gin-gonic/gin"
  "net/http"
)

type loginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}


func Login(c *gin.Context) {
    li := loginInput{}

		if err := c.ShouldBindJSON(&li); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
    if sessions.LsDB.CheckCredentials(li.Username, li.Password) {
			generatedToken := sessions.TokenCache.Add()
			c.JSON(http.StatusOK, gin.H{"token": generatedToken})
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"message": "login invalid"})
}


