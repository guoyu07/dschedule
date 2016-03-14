package scheduler

import (
	"fmt"
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

func (this *ResourceManager) AddMeta(meta *structs.NodeMeta) error {
	node := &structs.Node{
		NodeId:    fmt.Sprintf("%s-random-something-else", meta.IP),
		Used:      false,
		Reachable: false,
		Meta:      meta,
	}
	this.allNodes = append(this.allNodes, node)
	this.freeNodes = append(this.freeNodes, node)
	// TODO store
	return nil
}

func (this *ResourceManager) ModifyMeta(nodeId string, meta *structs.NodeMeta) error {
	// TODO modify
	// TODO store
	return nil
}

func (this *ResourceManager) RemoveNode(nodeId string) error {
	// TODO remove
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
