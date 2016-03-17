package scheduler

import (
	"fmt"
	//log "github.com/omidnikta/logrus"
	"github.com/samalba/dockerclient"
	"github.com/weibocom/dschedule/structs"
	"strconv"
	"strings"
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
		exposedPorts[key] = struct{}{}
		portBindings[key] = []dockerclient.PortBinding{
			dockerclient.PortBinding{
				//HostIp:   "0.0.0.0",
				HostPort: strconv.Itoa(hostPort),
			},
		}
	}

	var envs []string
	for key, val := range this.container.Env {
		envs = append(envs, fmt.Sprintf("%s=%s", key, val))
	}

	volumes := make(map[string]struct{})
	var binds []string
	for constainerFile, hostFile := range this.container.Volumes {
		volumes[constainerFile] = struct{}{}
		binds = append(binds, fmt.Sprintf("%s:%s", hostFile, constainerFile))
		//volumes[constainerFile] = struct{}{hostFile}
	}

	// Create a container
	containerConfig := &dockerclient.ContainerConfig{
		Image:        this.container.Image,
		Env:          envs,
		Volumes:      volumes,
		ExposedPorts: exposedPorts,
		Cmd:          []string{this.container.Command},

		AttachStdin: true,
		Tty:         true,
	}
	//log.Infof("deployer create container: %v", containerConfig)
	containerId, err := this.docker.CreateContainer(containerConfig, "", nil)
	if err != nil {
		return err
	}
	//log.Infof("deployer created container containerId: %v", containerId)
	// Start the container
	hostConfig := &dockerclient.HostConfig{
		Binds:        binds,
		PortBindings: portBindings,
		NetworkMode:  strings.ToLower(this.container.Network),
	}
	err = this.docker.StartContainer(containerId, hostConfig)
	if err != nil {
		//log.Errorf("deployer start container failed, cause: %v", err)
		return err
	}
	//log.Infof("deployer started container containerId: %v", containerId)
	this.containerId = containerId
	return nil
}

func (this *Deployer) Stop() error {
	// Stop the container (with 5 seconds timeout)
	this.docker.StopContainer(this.containerId, 5) // 5 -> timeout
	//log.Infof("deployer stopped container:%v", this.containerId)
	return nil
}
