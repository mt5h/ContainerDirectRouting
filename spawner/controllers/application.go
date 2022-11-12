package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"spawner/utils"
	"time"
)

// when we use the container NAME for external requests
// if you want to use the container ID use the provisiong APIs

func RestartContainer(c *gin.Context) {

	containerName := c.Param("container-name")
	container := utils.ContainerSummary{}

	if containerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no container specified"})
	}

	containerList, err := utils.ListContainers()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// check if we have the requested container
	container_found := false
	for _, cnt := range containerList {
		
    // check if the requested container has been created by us and is stopped.
		if cnt.ContainerName == containerName {
			container = cnt
			container_found = true
			break
		}
	}

	if !container_found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "container not found"})
		return
	}

	err = utils.MatchContainerLabel(container.ContainerID, "origin", "spawner")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if container.ContainerStatus == "exited" {
		err := utils.StartContainer(container.ContainerID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		endpointUrl, err := utils.GetContainerLabel(container.ContainerID, "healthcheck")
    
		if err != nil {
			log.Println("healthcheck label not found using sleep method")
			utils.StartContainer(container.ContainerID)
			time.Sleep(time.Second * time.Duration(utils.RedirectTimeout))
			c.Redirect(301, c.Request.URL.String())
			return
		} else {

			healthcheckResult := false
			retries := 3
			for {
				retries -= 1
				log.Println(retries, "checking", endpointUrl)

				healthcheckResult = utils.HttpHealthCheck(endpointUrl, time.Duration(20*time.Second), 200)

				if healthcheckResult || retries == 0 {
					break
				} else {
					time.Sleep(time.Second * time.Duration(utils.RedirectHealthcheckTimeout))
				}

			}
			if healthcheckResult {
				c.Redirect(301, c.Request.URL.String())
				return
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "endpoint unreachable"})
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "container unexpected state"})
	}

}
