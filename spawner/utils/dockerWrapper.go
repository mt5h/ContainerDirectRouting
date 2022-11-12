package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type Instance struct {
	Name    string            `from: "Name", json:"name" binding:"required"`
	Network string            `from: "Network", json:"network" binding:"required"`
	Image   string            `from: "Image", json:"image" binding:"required"`
	Labels  map[string]string `from: "Labels", json:"labels" binding:"required"`
	Envs    map[string]string `from: "Envs", json:"envs" binding:"required"`
}

type ContainerSummary struct {
	ContainerID     string
	ContainerName   string
	ContainerLabels map[string]string
	ContainerImage  string
	ContainerStatus string
}

type InstanceStates struct {
	Instances []Instance
}

func (containerInstance *Instance) EnvDockerFormat() []string {

	envNum := len(containerInstance.Envs)
	res := make([]string, envNum)
	counter := 0
	for key, value := range containerInstance.Envs {
		res[counter] = fmt.Sprintf("%s=%s", key, value)
		counter++
	}
	return res
}

func GetContainerLabel(containerId, label_key string) (string, error) {
	container, err := GetContainer(containerId)
	if err != nil {
		return "", err
	}
	value, found := container.ContainerLabels[label_key]
	if found {
		return value, nil
	} else {
		return "", errors.New("Container tag not found")
	}
}

// TODO write the reverse from docker to our form

func MatchContainerLabel(containerId, label_key, label_value string) error {
	container, err := GetContainer(containerId)
	if err != nil {
		return err
	}
	value, found := container.ContainerLabels[label_key]
	if found {
		if value == label_value {
			return nil
		} else {
			return errors.New("Container tag not match")
		}
	}
	return errors.New("Container tag not found")

}

func CreateContainer(containerInstance Instance) (string, error) {
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
		Env:      containerInstance.EnvDockerFormat(),
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
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}
	return resp.ID, nil

}

func ListContainers() ([]ContainerSummary, error) {

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

	containerList := make([]ContainerSummary, len(containers))

	for cindex, container := range containers {

		containerList[cindex] = ContainerSummary{
			ContainerID:     container.ID,
			ContainerImage:  container.Image,
			ContainerName:   container.Names[0][1:],
			ContainerLabels: container.Labels,
			ContainerStatus: container.State,
		}
	}

	return containerList, nil
}

func StartContainer(containerId string) error {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	defer cli.Close()
	if err := cli.ContainerStart(ctx, containerId, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return err

}

func GetContainer(containerId string) (ContainerSummary, error) {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return ContainerSummary{}, err
	}
	defer cli.Close()

	containerInfo, err := cli.ContainerInspect(ctx, containerId)
	if err != nil {
		return ContainerSummary{}, err
	}

	cnt := ContainerSummary{
		ContainerID:     containerInfo.ID,
		ContainerName:   containerInfo.Name,
		ContainerImage:  containerInfo.Image,
		ContainerLabels: containerInfo.Config.Labels,
		ContainerStatus: containerInfo.State.Status,
	}

	return cnt, nil
}

func DeleteContainer(containerId string) error {
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
	return nil
}
