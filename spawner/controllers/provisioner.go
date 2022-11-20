package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"spawner/utils"
)

// this uses the container id
func DeleteContainer(c *gin.Context) {

	containerId := c.Param("container-id")

	err := utils.MatchContainerLabel(containerId, "origin", "spawner")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = utils.DeleteContainer(containerId)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"result": "ok"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func ListContainers(c *gin.Context) {
	localContainers, err := utils.ListContainers()

	instances := []utils.ContainerSummary{}

	for _, cnt := range localContainers {
		err := utils.MatchContainerLabel(cnt.ContainerID, "origin", "spawner")
		if err == nil {
			instances = append(instances, cnt)
		}
	}

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"instances": instances})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func CreateContainer(c *gin.Context) {

	var instance utils.Instance

	if err := c.ShouldBindJSON(&instance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// we add our label to make sure we created it
	instance.Labels["origin"] = "spawner"

	containerID, err := utils.CreateContainer(instance)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"id": containerID})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}

func StartContainer(c *gin.Context) {
	containerId := c.Param("container-id")
	if containerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "container id missing"})
		return
	}

	cnt, err := utils.GetContainer(containerId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = utils.MatchContainerLabel(cnt.ContainerID, "origin", "spawner")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if cnt.ContainerStatus == "exited" {
		err = utils.StartContainer(cnt.ContainerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": cnt.ContainerID})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("container is %s", cnt.ContainerStatus)})
		return
	}
}
