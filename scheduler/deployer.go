package scheduler

import (
	"fmt"
	log "github.com/omidnikta/logrus"
	"github.com/samalba/dockerclient"
	"github.com/weibocom/dschedule/structs"
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
	log.Infof("deployer start container:%v", this.containerId)
	// Create a container
	containerConfig := &dockerclient.ContainerConfig{
		Image:       this.container.Image,
		Cmd:         []string{this.container.Command},
		AttachStdin: true,
		Tty:         true,
	}
	log.Infof("deployer create container: %v", containerConfig)
	containerId, err := this.docker.CreateContainer(containerConfig, "", nil)
	if err != nil {
		return err
	}
	log.Infof("create container containerId: %v", containerId)
	// Start the container
	hostConfig := &dockerclient.HostConfig{}
	err = this.docker.StartContainer(containerId, hostConfig)
	if err != nil {
		log.Errorf("deployer start container failed, cause: %v", err)
		return err
	}
	this.containerId = containerId
	return nil
}

func (this *Deployer) Stop() error {
	log.Infoln("deployer stop container:%v", this.containerId)
	// Stop the container (with 5 seconds timeout)
	this.docker.StopContainer(this.containerId, 5) // 5 -> timeout
	return nil
}
