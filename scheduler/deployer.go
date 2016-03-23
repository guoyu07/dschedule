package scheduler

import (
	"fmt"
	//log "github.com/omidnikta/logrus"
	"github.com/samalba/dockerclient"
	"github.com/weibocom/dschedule/structs"
	"strings"
)

type Deployer struct {
	docker   *dockerclient.DockerClient
	node     *structs.Node
	nodeEnvs []string

	containers   []*structs.Container
	containerIds []string
}

func NewDeployer(node *structs.Node, dockerPort int,
	containers []*structs.Container) (*Deployer, error) {

	host := fmt.Sprintf("%s:%d", node.Meta.IP, dockerPort)
	docker, err := dockerclient.NewDockerClient(host, nil)
	if err != nil {
		return nil, err
	}

	var nodeEnvs []string
	nodeEnvs = append(nodeEnvs, fmt.Sprintf("NODE_META_NAME=%s", node.Meta.Name))
	nodeEnvs = append(nodeEnvs, fmt.Sprintf("NODE_META_IP=%s", node.Meta.IP))
	nodeEnvs = append(nodeEnvs, fmt.Sprintf("NODE_META_CPU=%s", node.Meta.CPU))
	nodeEnvs = append(nodeEnvs, fmt.Sprintf("NODE_META_MEMORY_MB=%s", node.Meta.MemoryMB))
	nodeEnvs = append(nodeEnvs, fmt.Sprintf("NODE_META_DISK_MB=%s", node.Meta.DiskMB))
	nodeEnvs = append(nodeEnvs, fmt.Sprintf("NODE_META_DISK_DIRS=%s", strings.Join(node.Meta.DiskDirs, ":")))
	for key, val := range node.Meta.Attributes {
		nodeEnvs = append(nodeEnvs, fmt.Sprintf("NODE_META_DISK_ATTRIBUTE_%s=%s", strings.ToUpper(key), val))
	}

	return &Deployer{
		docker:     docker,
		node:       node,
		nodeEnvs:   nodeEnvs,
		containers: containers,
	}, nil
}

func (this *Deployer) Start() error {
	for _, container := range this.containers {
		containerId, err := this.run(container)
		if err != nil {
			return err
		}
		this.containerIds = append(this.containerIds, containerId)
	}
	return nil
}

func (this *Deployer) Stop() error {
	for idx := len(this.containers) - 1; idx >= 0; idx-- {
		// Stop the container (with 5 seconds timeout)
		this.docker.StopContainer(this.containerIds[idx], 5) // 5 -> timeout
		//log.Infof("deployer stopped container:%v", containerId)
	}
	return nil
}

func (this *Deployer) run(container *structs.Container) (string, error) {
	// Network: Expose frist and binding after: https://github.com/docker/docker/issues/2785
	//          https://docker-py.readthedocs.org/en/latest/port-bindings/
	exposedPorts := make(map[string]struct{})
	portBindings := make(map[string][]dockerclient.PortBinding)
	for containerPort, hostPort := range container.PortMapping {
		key := fmt.Sprintf("%s/tcp", containerPort)
		exposedPorts[key] = struct{}{}
		portBindings[key] = []dockerclient.PortBinding{
			dockerclient.PortBinding{
				//HostIp:   "0.0.0.0",
				HostPort: hostPort,
			},
		}
	}

	envs := this.nodeEnvs[:]
	for key, val := range container.Env {
		envs = append(envs, fmt.Sprintf("%s=%s", key, val))
	}

	volumes := make(map[string]struct{})
	var binds []string
	for constainerFile, hostFile := range container.Volumes {
		volumes[constainerFile] = struct{}{}
		binds = append(binds, fmt.Sprintf("%s:%s", hostFile, constainerFile))
		//volumes[constainerFile] = struct{}{hostFile}
	}

	// Create a container
	containerConfig := &dockerclient.ContainerConfig{
		Image:        container.Image,
		Env:          envs,
		Volumes:      volumes,
		ExposedPorts: exposedPorts,

		AttachStdin: true,
		Tty:         true,
	}
	// DEBUGGED: System error: exec: "": executable file not found in $PATH
	if container.Command != "" {
		containerConfig.Cmd = []string{container.Command}
	}
	//images, err := this.docker.SearchImages(containerConfig.Image, "", nil)
	//if err != nil {
	//	return err
	//}
	//fmt.Println("images:", images)
	//if len(images) == 0 {
	//fmt.Println("pulling ", containerConfig.Image)

	// DEBUGGED: "Image not found"
	err := this.docker.PullImage(containerConfig.Image, nil)
	if err != nil {
		return "", fmt.Errorf("Pull image(%s) failed: %v", containerConfig.Image, err)
	}
	//}
	//log.Infof("deployer create container: %v", containerConfig)
	containerId, err := this.docker.CreateContainer(containerConfig, "", nil)
	if err != nil {
		return "", fmt.Errorf("Create container failed: %v", err)
	}
	//log.Infof("deployer created container containerId: %v", containerId)
	// Start the container
	hostConfig := &dockerclient.HostConfig{
		Binds:        binds,
		PortBindings: portBindings,
		NetworkMode:  strings.ToLower(container.Network),
	}
	err = this.docker.StartContainer(containerId, hostConfig)
	if err != nil {
		//log.Errorf("deployer start container failed, cause: %v", err)
		return containerId, fmt.Errorf("Start container(%s) failed: %v", err)
	}
	//log.Infof("deployer started container containerId: %v", containerId)
	return containerId, nil
}
