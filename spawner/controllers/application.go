package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"spawner/config"
	"spawner/utils"
	"strings"
	"time"
)

type healthCheckMethod int64

const (
	NoHC healthCheckMethod = iota
	HttpHC
	DockerHC
)

func CheckTraefickRoute(containerId string) error {
  // get the container status
  instance, err := utils.GetContainer(containerId)
		if err != nil {
			return err
	}

  log.Println("-->", instance.ContainerName)
  
  endpointUrl := fmt.Sprintf("%s/api/http/routers/%s@%s", config.TraefikBaseUrl, instance.ContainerName, config.TraefikPlatform)
  // if the route for containerName has not been found we do expect a 404
  retries:=0
  for {
    if found := utils.HttpHealthCheck(endpointUrl, 10*time.Second, 200); found {
      return nil
    }
    log.Println(instance.ContainerName, "traefick route not ready")
    time.Sleep(500*time.Millisecond)
    retries+=1
    if retries > 20{
      // try to stop the container and  exit
      if err := utils.StopContainer(instance.ContainerID); err != nil {
        log.Println("Error stopping", instance.ContainerName, ":", err.Error())
      }
      // at this point we dont care if we can stop the container or not
      return errors.New("Traefik route not detected. Container Stopped")
    }
  }
}

func CheckContainerHealth(containerId string) error {
		
  instance, err := utils.GetContainer(containerId)
		if err != nil {
			return err
		}

		usedHealthCheck := NoHC

		endpointUrl, err := utils.GetContainerLabel(instance.ContainerID, "health-check")

		if err == nil {
			usedHealthCheck = HttpHC
		}

		if instance.ContainerHealth != "" {
			usedHealthCheck = DockerHC
		}

		if usedHealthCheck == DockerHC {
			// Status is one of Starting, Healthy or Unhealthy
			log.Println("Using docker healthcheck")
			for {
				instance, err = utils.GetContainer(instance.ContainerID)
				if err != nil {
					return err
				}
				// we dont check healthcheck prop here we have already checked above
				if strings.ToLower(instance.ContainerHealth) == "starting" {
          log.Println(instance.ContainerName, "is starting. Fail streak:", instance.ContainerFailStreak)
				} else if strings.ToLower(instance.ContainerHealth) == "healthy" {
					log.Println(instance.ContainerName, "is healthy")
					// app is up nothing to do
					return nil
				} else if strings.ToLower(instance.ContainerHealth) == "unhealthy" {
          log.Println(instance.ContainerName, "is unhealthy. Fail streak:", instance.ContainerFailStreak)
					// max retries exceeded stop the container and return an error
					utils.StopContainer(instance.ContainerID)
					return errors.New(fmt.Sprintf("%s is unhealthy. Container stopped", instance.ContainerName))
				} else {
					log.Println("Invalid container healthcheck status")
					usedHealthCheck = NoHC
					break
				}
				time.Sleep(2 * time.Second)
			}
		}

		if usedHealthCheck == HttpHC {
			log.Println("Using HTTP health check")
			retries := 1
			for {

				if retries > config.HealthCheckRetries {
					// we exceeded the maximun number of retries
					utils.StopContainer(instance.ContainerID)
					return errors.New("http health-check failed")
				}

				log.Println(retries, "HTTP healthcheck on", endpointUrl)

				httpHealthCheckResult := utils.HttpHealthCheck(endpointUrl, config.HttpProbeTimeOut, config.HttpProbeOkStatus)

				if httpHealthCheckResult {
					return nil
				} else {
					log.Println(instance.ContainerName, "HTTP healthcheck failed")
					time.Sleep(config.HealthCheckInterval)
				}
				retries += 1
			}
		}

		log.Println("No valid health-check method found just using a timeout")
		time.Sleep(time.Duration(config.RedirectTimeout))
		return nil
}

// when we use the container NAME for external requests
// if you want to use the container ID use the provisiong APIs
func CheckContainerStatus(containerName string) (string, error) {

	if containerName == "" {
		return "", errors.New("container not specified")
	}

	containerList, err := utils.ListContainers()

	if err != nil {
		return "", err
	}
	// check if we have the requested container
	containerFound := false
	instance := utils.ContainerSummary{}
	for _, cnt := range containerList {
		// check if the requested container has been created by us and is stopped.
		if cnt.ContainerName == containerName {
			instance = cnt
			containerFound = true
			break
		}
	}

	if !containerFound {
		return "", errors.New("container not found")
	}

	if err = utils.MatchContainerLabel(instance.ContainerID, "origin", "spawner"); err != nil {
		return "", err
	}

	if strings.ToLower(instance.ContainerStatus) == "exited" {
		// start the container
		if err := utils.StartContainer(instance.ContainerID); err != nil {
			return "", err
		}
  } else if strings.ToLower(instance.ContainerStatus) == "running" {
		// get the current status
    log.Println(instance.ContainerName, "is already running")
  } else {
	  return "", errors.New("container unexpected state")
  }

  return instance.ContainerID, nil

}

func PathRouting(c *gin.Context) {

	containerName := c.Param("containerName")

	containerId, err := CheckContainerStatus(containerName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

  if err:=CheckContainerHealth(containerId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
  }

	c.Redirect(301, c.Request.URL.String())

}

func CookieRouting(c *gin.Context) {
	containerName, err := c.Cookie(config.CookieKey)
	// invalid or non existent cookie
	if err != nil {
		log.Printf("Request from %s, user agent: %s has the following error:%s\n", c.Request.RemoteAddr, c.Request.UserAgent(), err.Error())
		c.Redirect(http.StatusFound, config.CookieFallBackUrl)
		return
	}

	log.Println("Checking", containerName, "status")


	// if the correct cookie exists check the container status
  containerId, err := CheckContainerStatus(containerName)
  if err != nil {
		c.Redirect(http.StatusFound, config.CookieFallBackUrl)
  }

  err = CheckContainerHealth(containerId)

  if err != nil {
		log.Printf("Request from %s, user agent: %s has the following error:%s\n", c.Request.RemoteAddr, c.Request.UserAgent(), err.Error())
		c.Redirect(http.StatusFound, config.CookieFallBackUrl)
		return
	}

  if config.TraefikCheckEnabled {
    if err := CheckTraefickRoute(containerId); err != nil {
		  log.Printf("Request from %s, user agent: %s has the following error:%s\n", c.Request.RemoteAddr, c.Request.UserAgent(), err.Error())
		  c.Redirect(http.StatusFound, config.CookieFallBackUrl)
		  return
	  }
  } else {
    // always way 1s to ensure traefik route is added and avoid redirects loops
    time.Sleep(1*time.Second)
  }

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	// remove every path from the request
	redirectUrl := fmt.Sprintf("%s://%s", scheme, c.Request.Host)
	// set a custom cookie use by traefik
	// c.SetCookie(utils.CookieKey, containerName, utils.CookieMaxAge, "/", c.Request.URL.Hostname(), utils.CookieSecure, utils.CookieHttpOnly)
	c.Redirect(http.StatusFound, redirectUrl)

}
