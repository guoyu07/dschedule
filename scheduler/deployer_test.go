package scheduler

import (
	"github.com/weibocom/dschedule/structs"
	"testing"
)

func TestDeployer(t *testing.T) {
	rm, err := NewResourceManager()
	if err != nil {
		t.Fatalf("NewResourceManager failed: %v", err)
		return
	}
	nodeId, err := rm.AddMeta(&structs.NodeMeta{
		Name:     "local-machine",
		IP:       "10.229.88.91",
		CPU:      2,
		MemoryMB: 10240,
		DiskMB:   102400,
	})
	if err != nil {
		t.Fatalf("AddMeta failed: %v", err)
		return
	}
	node, err := rm.GetNode(nodeId)
	if err != nil {
		t.Fatalf("GetNode by nodeId %d failed: %v", nodeId, err)
		return
	}

	dockerPort := 4243
	container := &structs.Container{
		Type: "docker",

		Image: "docker.io/redis:2.8",

		//环境变量
		Env: map[string]string{
			"Name": "Chine",
		},

		// 文件挂载，key为container中的目录，value为host中的目录
		Volumes: map[string]string{
			"/data0/docker": "/data1/",
		},

		// 容器使用的网络模式，可选值HOST，BRIDGE，NONE，CONTAINER:NAME
		//Network: "HOST",
		Network: "BRIDGE",
		// 端口映射，key为container的端口，value为host中的端口
		PortMapping: map[int]int{
			6379: 6379,
		},

		// 容器启动时需要执行的命令
		Command: "redis-server",
	}

	deployer, err := NewDeployer(node, dockerPort, container)
	if err != nil {
		t.Fatalf("NewDeployer failed: %v", err)
		return
	}
	err = deployer.Start()
	if err != nil {
		t.Fatalf("Deployer start failed: %v", err)
		return
	}
	err = deployer.Stop()
	if err != nil {
		t.Fatalf("Deployer stop failed: %v", err)
		return
	}

}
