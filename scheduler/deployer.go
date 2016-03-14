package scheduler

import (
	"github.com/samalba/dockerclient"
)

type Deployer struct {
	docker *dockerclient.DockerClient
}

func NewDeployer(host string) (*Deployer, error) {
	docker, err := dockerclient.NewDockerClient(host, nil)
	if err != nil {
		return nil, err
	}
	return &Deployer{
		docker: docker,
	}, nil
}

func (this *Deployer) Start(image, cmd, name string) (string, error) {
	// Create a container
	containerConfig := &dockerclient.ContainerConfig{
		Image:       image,
		Cmd:         []string{cmd},
		AttachStdin: true,
		Tty:         true,
	}
	containerId, err := this.docker.CreateContainer(containerConfig, name, nil)
	if err != nil {
		return "", err
	}

	// Start the container
	hostConfig := &dockerclient.HostConfig{}
	err = this.docker.StartContainer(containerId, hostConfig)
	if err != nil {
		return "", err
	}
	return containerId, nil
}

func (this *Deployer) Stop(containerId string) error {
	// Stop the container (with 5 seconds timeout)
	this.docker.StopContainer(containerId, 5) // 5 -> timeout
	return nil
}
