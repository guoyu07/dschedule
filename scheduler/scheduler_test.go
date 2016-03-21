package scheduler

import (
	log "github.com/omidnikta/logrus"
	"github.com/weibocom/dschedule/structs"
	"testing"
	"time"
)

func TestScheduler(t *testing.T) {
	rm, err := NewResourceManager()
	if err != nil {
		t.Fatalf("NewResourceManager failed: %v", err)
		return
	}
	scheduler, err := NewScheduler(rm, 4243)
	if err != nil {
		t.Fatalf("NewScheduler failed: %v", err)
		return
	}

	// test in same ip
	meta := &structs.NodeMeta{
		Name:     "local-machine",
		IP:       "10.229.88.91",
		CPU:      2,
		MemoryMB: 10240,
		DiskMB:   102400,
	}
	nodeId1, err := rm.AddMeta(meta)
	if err != nil {
		t.Fatalf("AddMeta 1 failed: %v", err)
		return
	}
	nodeId2, err := rm.AddMeta(meta)
	if err != nil {
		t.Fatalf("AddMeta 2 failed: %v", err)
		return
	}
	/*
			nodeId3, err := rm.AddMeta(meta)
			if err != nil {
				t.Fatalf("AddMeta 3 failed: %v", err)
				return
			}
		t.Fatalf("node1=%s, node2=%s, node3=%s", nodeId1, nodeId2, nodeId3)
	*/
	log.Infof("node1=%s, node2=%s", nodeId1, nodeId2)

	serviceId1 := "redis-bridge-priority1"
	serviceId2 := "redis-bridge-priority2"
	serviceId3 := "redis-host-priority5"
	// DEBUGGED: refer Dedicated, Elastic
	service1 := structs.Service{
		ServiceId: serviceId1,
		Priority:  1,
		Dedicated: 1,
		Elastic:   0,
		Container: &structs.Container{
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
			PortMapping: map[string]string{
				"6379": "6378",
			},

			// 容器启动时需要执行的命令
			Command: "redis-server",
		},
	}

	_, err = scheduler.Register(&service1)
	if err != nil {
		t.Fatalf("Register serviceId1 %s failed: %v", serviceId1, err)
		return
	}
	_, err = scheduler.Add(serviceId1, 1)
	if err != nil {
		t.Fatalf("Add serviceId1 %s failed: %v", serviceId1, err)
		return
	}
	log.Infof("finished serviceId1 %s", serviceId1)
	time.Sleep(20 * time.Second)

	// set higher priority redis but Dedicated 0
	service2 := service1
	service2.ServiceId = serviceId2
	service2.Priority = 3
	service2.Dedicated = 0
	service2.Elastic = 1
	service2.Container.PortMapping = map[string]string{
		"6379": "6377",
	}
	_, err = scheduler.Register(&service2)
	if err != nil {
		t.Fatalf("Register serviceId2 %s failed: %v", serviceId2, err)
		return
	}
	_, err = scheduler.Add(serviceId2, 1)
	if err != nil {
		t.Fatalf("Add serviceId2 %s failed: %v", serviceId2, err)
		return
	}

	log.Infof("finished serviceId2 %s", serviceId2)
	time.Sleep(20 * time.Second)
	// set higher priority redis
	service3 := service2
	service3.ServiceId = serviceId3
	service3.Priority = 5
	service3.Container.Network = "HOST"
	service3.Container.PortMapping = nil
	_, err = scheduler.Register(&service3)
	if err != nil {
		t.Fatalf("Register serviceId3 %s failed: %v", serviceId3, err)
		return
	}
	_, err = scheduler.Add(serviceId3, 1)
	if err != nil {
		t.Fatalf("Add serviceId3 %s failed: %v", serviceId3, err)
		return
	}

	log.Infof("finished serviceId3 %s", serviceId3)
	time.Sleep(20 * time.Second)
}
