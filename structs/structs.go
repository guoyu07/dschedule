package structs

import (
// "strings"
)

const (
	ContainerTypeDocker = "DOCKER"
	ContainerTypeQume   = "QEMU"

	ContainerNetworkHost      = "HOST"
	ContainerNetworkBridge    = "BRIDGE"
	ContainerNetworkNone      = "NONE"
	ContainerNetworkContainer = "CONTAINER"

	ServiceTypeProd      = "PROD"
	ServiceTypeNonProd   = "NON-PROD"
	ServiceTypeAuxiliary = "AUXILIARY"

	ServiceStrategyAuto    = "AUTO"
	ServiceStrategyCrontab = "CRONTAB"
	ServiceStrategyStable  = "STABLE"

	MaxPriority = 5
	MinPriority = 1
)

type Node struct {
	NodeId    string `json:"nodeId"` // uniq
	Used      bool   `json:"used"`
	Reachable bool   `json:"reachable"`
	Failed    int    `json:"failed"` // TODO more strategy

	Meta *NodeMeta `json:"meta"`
}

type NodeMeta struct {
	Name       string            `json:"name"`
	IP         string            `json:"ip"`
	CPU        int               `json:"cpu"`
	MemoryMB   int               `json:"memoryMB"`
	DiskMB     int               `json:"diskMB"`
	Attributes map[string]string `json:"attributes"` //other
}

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

	//容器名
	Name string

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

type Service struct {
	ServiceId string

	//任务类型，可选值prod，non-prod, auxiliary
	ServiceType string

	//任务执行的策略，可选值AUTO，CRONTAB，STABLE
	StrategyName string

	StrategyConfig interface{}

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
