package scheduler

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-multierror"
	log "github.com/omidnikta/logrus"
	"github.com/weibocom/dschedule/storage"
	"github.com/weibocom/dschedule/structs"
	"github.com/weibocom/dschedule/util"
	"sync"
)

const (
	MaxFailed = 3
)

const (
	NodeStoragePrefix = "node"
)

type ResourceManager struct {
	allNodes    []*structs.Node //map[string]*structs.Node
	freeNodes   []*structs.Node
	failedNodes []*structs.Node
	mutex       sync.Mutex
	storage     *storage.Storage
}

func NewResourceManager(storage *storage.Storage) (*ResourceManager, error) {
	// TODO ConstructFromConsul for allNodes, freeNodes, failedNodes
	return &ResourceManager{
		//usedNodes: make(map[string]*Node),
		storage: storage,
	}, nil
}

func (this *ResourceManager) AddMeta(meta *structs.NodeMeta) (string, error) {
	node := &structs.Node{
		NodeId:    fmt.Sprintf("%s-%s", meta.IP, util.GenerateUUID()),
		Used:      false,
		Reachable: false,
		Meta:      meta,
	}
	this.allNodes = append(this.allNodes, node)
	this.freeNodes = append(this.freeNodes, node)
	err := this.StoreNode(node)
	if err != nil {
		return "", fmt.Errorf("add node failed, cause: %v", err)
	}
	//log.Infof("node info is : ", node)
	return node.NodeId, nil
}

func (this *ResourceManager) ModifyMeta(nodeId string, meta *structs.NodeMeta) error {
	for _, node := range this.allNodes {
		if node.NodeId == nodeId {
			node.Meta = meta

			err := this.StoreNode(node)
			if err != nil {
				return fmt.Errorf("modify node failed, cause: %v", err)
			}
			return nil
		}
	}
	return fmt.Errorf("No nodeId: %s", nodeId)
}

func (this *ResourceManager) DeleteNode(nodeId string) error {
	idx := -1
	for i, node := range this.allNodes {
		if node.NodeId == nodeId {
			idx = i
		}
	}
	if idx < 0 {
		return fmt.Errorf("No nodeId: %s", nodeId)
	}

	tNodes := this.allNodes
	this.allNodes = tNodes[:idx]
	if idx+1 < len(tNodes) {
		this.allNodes = append(this.allNodes, tNodes[idx+1:]...)
	}

	err := this.RemoveNode(nodeId)
	if err != nil {
		return fmt.Errorf("remove node failed, cause: %v", err)
	}
	return nil
}

func (this *ResourceManager) GetNode(nodeId string) (*structs.Node, error) {
	for _, node := range this.allNodes {
		if node.NodeId == nodeId {
			return node, nil
		}
	}
	nodes, err := this.RetriveNode(nodeId)
	if err != nil {
		return nil, fmt.Errorf("reosurce manager get node failed, cause: %v", err)
	}
	if len(nodes) > 0 {
		return nodes[0], nil
	}
	return nil, nil
}

func (this *ResourceManager) GetNodeList() ([]*structs.Node, error) {
	if len(this.allNodes) == 0 {
		nodes, err := this.RetriveNode("")
		if err != nil {
			return nil, fmt.Errorf("resource manager get node list faield, cause : %v", err)
		}
		if len(nodes) != 0 {
			this.allNodes = nodes
		}
	}
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
	var mErr multierror.Error
	for _, node := range allocs {
		node.Used = true
		err := this.StoreNode(node)
		if err != nil {
			mErr.Errors = append(mErr.Errors, err)
		}
	}

	return allocs, mErr.ErrorOrNil()
}

func (this *ResourceManager) ReturnNodes(nodes []*structs.Node) error {
	var mErr multierror.Error
	for _, node := range nodes {
		node.Used = false
		if node.Failed < MaxFailed {
			this.freeNodes = append(this.freeNodes, node)
		} else {
			log.Warnf("Node '%s' IP '%s' Failed have reached %d, insert into failed queue", node.NodeId, node.Meta.IP, node.Failed)
			this.failedNodes = append(this.failedNodes, node)
		}
		err := this.StoreNode(node)
		if err != nil {
			mErr.Errors = append(mErr.Errors, err)
		}
	}
	return mErr.ErrorOrNil()
}

/**
* save node to storage
 */
func (this *ResourceManager) StoreNode(node *structs.Node) error {

	value, err := json.Marshal(node)
	if err != nil {
		return fmt.Errorf("json marshal node failed, cause: %v", err)
	}
	key := fmt.Sprintf("%s/%s", NodeStoragePrefix, node.NodeId)
	log.Infof("store node key:%v, value:%v", key, string(value))
	err = this.storage.Client.Put(key, []byte(value), nil)
	if err != nil {
		return fmt.Errorf("Error trying to put value at key: %v", key)
	}
	return nil
}

/**
*	remove node from storage
 */
func (this *ResourceManager) RemoveNode(nodeId string) error {
	key := fmt.Sprintf("%s/%s", NodeStoragePrefix, nodeId)
	err := this.storage.Client.Delete(key)
	if err != nil {
		return fmt.Errorf("Error trying to delete key: %v", key)
	}
	return nil
}

/**
* retrive node message from storage
 */
func (this *ResourceManager) RetriveNode(nodeId string) ([]*structs.Node, error) {
	var nodes []*structs.Node
	key := fmt.Sprintf("%s/%s", NodeStoragePrefix, nodeId)
	if nodeId == "" { // retrive list
		entries, err := this.storage.Client.List(key)
		if err != nil {
			return nil, fmt.Errorf("storage get node list failed, cause: %v", err)
		}
		for _, pair := range entries {
			var node *structs.Node
			if err := json.Unmarshal(pair.Value, &node); err != nil {
				log.Errorf("json unmarshal node failed, cause: %v", err)
				continue
			}
			log.Infof("retrive node : %v", node)
			nodes = append(nodes, node)
		}
	} else {
		pair, err := this.storage.Client.Get(key)
		if err != nil {
			return nil, fmt.Errorf("Error trying to get key: %v, cause: %v", key, err)
		}
		var node *structs.Node
		if err := json.Unmarshal(pair.Value, &node); err != nil {
			return nil, fmt.Errorf("json unmarshal node failed, cause: %v", err)
		}
		log.Infof("retrive node : %v", node)
		nodes = append(nodes, node)
	}

	return nodes, nil
}
