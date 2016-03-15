package scheduler

import (
	"fmt"
	log "github.com/omidnikta/logrus"
	"github.com/weibocom/dschedule/structs"
	"sync"
)

type ResourceManager struct {
	allNodes  []*structs.Node //map[string]*structs.Node
	freeNodes []*structs.Node
	mutex     sync.Mutex
}

func NewResourceManager() (*ResourceManager, error) {
	return &ResourceManager{
	//usedNodes: make(map[string]*Node),
	}, nil
}

func (this *ResourceManager) AddMeta(meta *structs.NodeMeta) (string, error) {
	node := &structs.Node{
		NodeId:    fmt.Sprintf("%s-random-something-else", meta.IP),
		Used:      false,
		Reachable: false,
		Meta:      meta,
	}
	this.allNodes = append(this.allNodes, node)
	this.freeNodes = append(this.freeNodes, node)
	// TODO store
	log.Infof("node info is : ", node)
	return node.NodeId, nil
}

func (this *ResourceManager) ModifyMeta(nodeId string, meta *structs.NodeMeta) error {
	for _, node := range this.allNodes {
		if node.NodeId == nodeId {
			node.Meta = meta
			// TODO store
			return nil
		}
	}
	return fmt.Errorf("No nodeId: %d", nodeId)
}

func (this *ResourceManager) RemoveNode(nodeId string) error {
	idx := -1
	for i, node := range this.allNodes {
		if node.NodeId == nodeId {
			idx = i
		}
	}
	if idx < 0 {
		return fmt.Errorf("No nodeId: %d", nodeId)
	}

	tNodes := this.allNodes
	this.allNodes = tNodes[:idx]
	if idx+1 < len(tNodes) {
		this.allNodes = append(this.allNodes, tNodes[idx+1:]...)
	}

	// TODO store
	return nil
}

func (this *ResourceManager) GetNode(nodeId string) (*structs.Node, error) {
	for _, node := range this.allNodes {
		if node.NodeId == nodeId {
			return node, nil
		}
	}
	return nil, nil
}

func (this *ResourceManager) GetNodeList() ([]*structs.Node, error) {
	return this.allNodes, nil
}

///////////////////////////////////////////////////////////////

func (this *ResourceManager) AllocNodes(num int /*, rules*/) ([]*structs.Node, error) {

	var allocs []*structs.Node

	if len(this.freeNodes) > num {
		allocs = this.freeNodes[:num]
		this.freeNodes = this.freeNodes[num:]
	} else {
		allocs = this.freeNodes
		this.freeNodes = this.freeNodes[:0]
	}

	for _, node := range allocs {
		node.Used = true
	}
	// TODO store

	return allocs, nil
}

func (this *ResourceManager) ReturnNodes(nodes []*structs.Node) error {
	for _, node := range nodes {
		node.Used = false
	}
	this.freeNodes = append(this.freeNodes, nodes...)
	// TODO store
	return nil
}
