package api

import (
	"fmt"
	"net/http"
	// "strconv"
	"strings"
)

const (
	ContainerTypeDocker = "DOCKER"
	ContainerTypeQume   = "QEMU"

	ContainerNetworkHost      = "HOST"
	ContainerNetworkBridge    = "BRIDGE"
	ContainerNetworkNone      = "NONE"
	ContainerNetworkContainer = "CONTAINER"

	TaskTypeProd      = "PROD"
	TaskTypeNonProd   = "NON-PROD"
	TaskTypeAuxiliary = "AUXILIARY"

	TaskStrategyAuto    = "AUTO"
	TaskStrategyCrontab = "CRONTAB"
	TaskStrategyStable  = "STABLE"
)

type Constraint struct {
	// 左值
	LeftTarget string
	// 右值
	RightTarget string
	// 操作
	Operand string
}

type Container struct {
	//容器类型，可选值DOCKER，QUMU，现在支持DOCKER
	Type string

	//镜像名
	Image string

	// 容器使用的网络模式，可选值HOST，BRIDGE，NONE，CONTAINER:NAME
	Network string

	//环境变量
	Env map[string]string

	// 文件挂载，key为container中的目录，value为host中的目录
	Volumes map[string]string

	// 端口映射，key为container的端口，value为host中的端口
	PortMapping map[string]string

	// 容器启动时需要执行的命令
	Command string
}

type Task struct {
	TaskId string

	//任务类型，可选值prod，non-prod, auxiliary
	TaskType string

	//任务执行的策略，可选值AUTO，CRONTAB，STABLE
	Strategy string

	//任务优先级，1-5，数值越大优先级越高
	Priority int

	//专用实例数
	Dedicated int

	//可伸缩实例数
	Elastic int

	//限制条件，可能是操作系统的版本，网卡个数，磁盘大小等
	Constraints []*Constraint

	//容器
	Container *Container
}

func (s *HTTPServer) TaskEndpoint(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	taskId := strings.TrimPrefix(req.URL.Path, "/task/")
	// Switch on the method
	switch req.Method {
	case "GET":
		if taskId != "" {
			return s.getTask(resp, req, taskId)
		}
		return s.listTask(resp, req)
	case "PUT":
		fallthrough
	case "POST":
		if taskId != "" {
			return s.modifyTask(resp, req, taskId)
		}
		return s.addTask(resp, req)
	case "DELETE":
		// Pull out the node id,
		if taskId != "" {
			return s.deleteTask(resp, req, taskId)
		} else {
			return nil, fmt.Errorf("Url '%s' with DELETE is illeage, should like '/task/{taskId}'", req.URL.Path)
		}

	default:
		resp.WriteHeader(405)
		return nil, nil
	}
}

func (s *HTTPServer) getTask(resp http.ResponseWriter, req *http.Request, taskId string) (interface{}, error) {
	// todo
	return nil, nil
}

func (s *HTTPServer) listTask(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	// todo
	return nil, nil
}

func (s *HTTPServer) modifyTask(resp http.ResponseWriter, req *http.Request, taskId string) (interface{}, error) {
	//todo
	return nil, nil
}

func (s *HTTPServer) addTask(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	//todo
	return nil, nil
}

func (s *HTTPServer) deleteTask(resp http.ResponseWriter, req *http.Request, taskId string) (interface{}, error) {
	//todo
	return nil, nil
}
