package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"net/http"
)


type Instance struct {
	Name    string            `from: "Name", json:"name" binding:"required"`
	Network string            `from: "Network", json:"network" binding:"required"`
	Image   string            `from: "Image", json:"image" binding:"required"`
	Labels  map[string]string `from: "Labels", json:"labels" binding:"required"`
	Envs    map[string]string `from: "Envs", json:"envs" binding:"required"`
}

type containerSummary struct {
	ContainerID     string
	ContainerName   string
	ContainerLabels map[string]string
	ContainerImage  string
	ContainerStatus string
	ContainerState  string
}

type InstanceStates struct {
	Instances []Instance
}

func (containerInstance *Instance) envDockerFormat() []string {

	envNum := len(containerInstance.Envs)
	res := make([]string, envNum)
	counter := 0
	for key, value := range containerInstance.Envs {
		res[counter] = fmt.Sprintf("%s=%s", key, value)
		counter++
	}
	return res

}


func createContainer(containerInstance Instance) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer cli.Close()

	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{},
	}

	networkConfig.EndpointsConfig[containerInstance.Network] = &network.EndpointSettings{}

	containerConfig := &container.Config{
		Hostname: containerInstance.Name,
		Env:      containerInstance.envDockerFormat(),
		Labels:   containerInstance.Labels,
		Image:    containerInstance.Image,
	}

	resp, err := cli.ContainerCreate(
		ctx,
		containerConfig,
		nil,
		networkConfig,
		nil,
		containerInstance.Name,
	)

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}
	fmt.Printf("%s\n", resp.ID)
	return resp.ID, nil

}

func listContainers() ([]containerSummary, error) {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	listOptions := types.ContainerListOptions{
		All: true,
	}
	containers, err := cli.ContainerList(ctx, listOptions)
	if err != nil {
		return nil, err
	}

	containerList := make([]containerSummary, len(containers))

	for cindex, container := range containers {
		containerList[cindex] = containerSummary{
			ContainerID:     container.ID,
			ContainerImage:  container.Image,
			ContainerName:   container.Names[0][1:],
			ContainerLabels: container.Labels,
			ContainerStatus: container.Status,
			ContainerState:  container.State,
		}
	}

	return containerList, nil
}

func deleteContainer(containerId string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	if err := cli.ContainerStop(ctx, containerId, nil); err != nil {
		return err
	}


	cli.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{})
	fmt.Println("Success")
	return nil
}

func deleteContainerWrapper(c *gin.Context) {

	containerId := c.Param("containerid")
	err := deleteContainer(containerId)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"result": "ok"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func listContainersWrapper(c *gin.Context) {
	res, err := listContainers()

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"res": res})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func createContainerWrapper(c *gin.Context) {

	var instance Instance

	if err := c.ShouldBindJSON(&instance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	containerID, err := createContainer(instance)

	if err == nil {
		fmt.Println("no errors")
		c.JSON(http.StatusOK, gin.H{"id": containerID})
	} else {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}

func main() {

	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.POST("/", createContainerWrapper)
		v1.GET("/", listContainersWrapper)
		v1.DELETE("/:containerid", deleteContainerWrapper)
	}

	router.Run(":8008")
}

/*
  labels := map[string]string{}
  labels["traefik.enable"] = "true"
  labels[fmt.Sprintf("traefik.http.routers.%s.rule",containerName)] = fmt.Sprintf("Path(`/%s`)",containerName)
  labels[fmt.Sprintf("traefik.http.routers.%s.entrypoints",containerName)] = "web"
*/

//env := []string{fmt.Sprintf("INSTANCE=%s",containerName)}
/*
	labels := map[string]string{}
	labels["traefik.enable"] = "true"
	labels[fmt.Sprintf("traefik.http.routers.%s.rule", containerName)] = fmt.Sprintf("Path(`/%s`)", containerName)
	labels[fmt.Sprintf("traefik.http.routers.%s.entrypoints", containerName)] = "web"

	envs := map[string]string{
		"INSTANCE": containerName,
		"TEST":     "123",
	}

	instance := Instance{
		Name:    containerName,
		Network: "example",
		Labels:  labels,
		Envs:    envs,
		Image:   "app-instance",
	}
*/
