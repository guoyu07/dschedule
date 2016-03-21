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
		//Name:     "local-machine",
		//IP:       "10.229.88.91",
		Name:     "online-machine",
		IP:       "10.73.88.41",
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
	/*
		container := &structs.Container{
			Type:  "docker",
			Image: "registry.intra.weibo.com/weibo_rd_if/remind-web:remind-web_RELEASE_V3.49",
			Env: map[string]string{
				"NAME_CONF": "openapi_remind-tc-inner=/data1/weibo",
			},
			Volumes: map[string]string{
				// DEBUGGED: remove -v when host network in dockerclient, but ok in cmdLine
				//"/etc/resolv.conf":     "/etc/resolv.conf",
				"/data1/mblog/logs/":   "/data1/weibo/logs/",
				"/data1/mblog/gclogs/": "/data1/weibo/gclogs",
				"/data0/docker":        "/data1/",
			},
			Network: "HOST",
			Command: "/docker_init.sh", //"redis-server",
		}
	*/
	container := &structs.Container{
		Type: "docker",

		Image: "registry.api.weibo.com/liubin8/nginx:latest",

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
			80: 100,
		},

		// 容器启动时需要执行的命令
		Command: "", //"redis-server",
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
	//err = deployer.Stop()
	if err != nil {
		t.Fatalf("Deployer stop failed: %v", err)
		return
	}

}
