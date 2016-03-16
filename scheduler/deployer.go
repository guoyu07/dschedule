package scheduler

import (
	"fmt"
	"github.com/samalba/dockerclient"
	"github.com/weibocom/dschedule/structs"
	"strconv"
)

type Deployer struct {
	docker      *dockerclient.DockerClient
	node        *structs.Node
	container   *structs.Container
	containerId string
}

func NewDeployer(node *structs.Node, dockerPort int, container *structs.Container) (*Deployer, error) {
	host := fmt.Sprintf("%s:%d", node.Meta.IP, dockerPort)
	docker, err := dockerclient.NewDockerClient(host, nil)
	if err != nil {
		return nil, err
	}
	return &Deployer{
		docker:    docker,
		node:      node,
		container: container,
	}, nil
}

func (this *Deployer) Start() error {
	// Network: Expose frist and binding after: https://github.com/docker/docker/issues/2785
	//          https://docker-py.readthedocs.org/en/latest/port-bindings/
	exposedPorts := make(map[string]struct{})
	portBindings := make(map[string][]dockerclient.PortBinding)
	for containerPort, hostPort := range this.container.PortMapping {
		key := fmt.Sprintf("%d/tcp", containerPort)
		exposedPorts[key] = struct{}{} // TODO may not ok
		portBindings[key] = []dockerclient.PortBinding{
			dockerclient.PortBinding{
				//HostIp:   "0.0.0.0",
				HostPort: strconv.Itoa(hostPort),
			},
		}
	}
	// Create a container
	containerConfig := &dockerclient.ContainerConfig{
		Image:        this.container.Image,
		Cmd:          []string{this.container.Command},
		ExposedPorts: exposedPorts,

		AttachStdin: true,
		Tty:         true,
	}
	containerId, err := this.docker.CreateContainer(containerConfig, "", nil)
	if err != nil {
		return err
	}

	// Start the container
	hostConfig := &dockerclient.HostConfig{
		NetworkMode:  this.container.Network,
		PortBindings: portBindings,
	}
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
