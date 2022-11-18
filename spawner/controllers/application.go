package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"spawner/utils"
	"time"
)

// when we use the container NAME for external requests
// if you want to use the container ID use the provisiong APIs
func RestartContainer(containerName string) error {

	if containerName == "" {
		return errors.New("container not specified")
		// c.JSON(http.StatusBadRequest, gin.H{"error": "no container specified"})
	}

	containerList, err := utils.ListContainers()

	if err != nil {
		return err
	}
	// check if we have the requested container
	container_found := false
	container := utils.ContainerSummary{}
	for _, cnt := range containerList {

		// check if the requested container has been created by us and is stopped.
		if cnt.ContainerName == containerName {
			container = cnt
			container_found = true
			break
		}
	}

	if !container_found {
		return errors.New("container not found")
	}

	err = utils.MatchContainerLabel(container.ContainerID, "origin", "spawner")

	if err != nil {
		return err
	}

	if container.ContainerStatus == "exited" {
		err := utils.StartContainer(container.ContainerID)

		if err != nil {
			return err
		}

		endpointUrl, err := utils.GetContainerLabel(container.ContainerID, "healthcheck")

		if err != nil {
			log.Println("healthcheck label not found using sleep method")
			utils.StartContainer(container.ContainerID)
			time.Sleep(time.Second * time.Duration(utils.RedirectTimeout))
			return nil
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
				return nil
			} else {
				return errors.New("endpoint unreachable")
			}
		}
	}

	if container.ContainerStatus == "running" {
		// container maybe started but you dont have a valid cookie so the route doesn't match
		// let's refresh that
		return nil
	}

	return errors.New("container unexpected state")

}

func PathRouting(c *gin.Context) {

	containerName := c.Param("containerName")

	err := RestartContainer(containerName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(301, c.Request.URL.String())
}

func CookieRouting(c *gin.Context) {
	containerName := c.Param("containerName")

	err := RestartContainer(containerName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	// remove every path from the request
	redirectUrl := fmt.Sprintf("%s://%s", scheme, c.Request.Host)
	// set a custom cookie use by traefik
	c.SetCookie(utils.CookieKey, containerName, utils.CookieMaxAge, "/", c.Request.URL.Hostname(), utils.CookieSecure, utils.CookieHttpOnly)
	c.Redirect(http.StatusFound, redirectUrl)

}
