package scheduler

import (
	"github.com/samalba/dockerclient"
	"github.com/weibocom/dschedule/structs"
)

type Deployer struct {
	docker      *dockerclient.DockerClient
	container   *structs.Container
	containerId string
}

func NewDeployer(host string, container *structs.Container) (*Deployer, error) {
	docker, err := dockerclient.NewDockerClient(host, nil)
	if err != nil {
		return nil, err
	}
	return &Deployer{
		docker:    docker,
		container: container,
	}, nil
}

func (this *Deployer) Start() error {
	// Create a container
	containerConfig := &dockerclient.ContainerConfig{
		Image:       this.container.Image,
		Cmd:         []string{this.container.Command},
		AttachStdin: true,
		Tty:         true,
	}
	containerId, err := this.docker.CreateContainer(containerConfig, "", nil)
	if err != nil {
		return err
	}

	// Start the container
	hostConfig := &dockerclient.HostConfig{}
	err = this.docker.StartContainer(containerId, hostConfig)
	if err != nil {
		return err
	}
	this.containerId = containerId
	return nil
}

func (this *Deployer) Stop() error {
	// Stop the container (with 5 seconds timeout)
	this.docker.StopContainer(this.containerId, 5) // 5 -> timeout
	return nil
}
