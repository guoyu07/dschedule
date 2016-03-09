package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/omidnikta/logrus"
	sched "github.com/weibocom/dschedule/scheduler"
	"io"
	"net/http"
	"strings"
)

func (s *HTTPServer) NodeEndpoint(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	nodeId := strings.TrimPrefix(req.URL.Path, "/node/")
	// Switch on the method
	switch req.Method {
	case "GET":
		if nodeId != "" {
			return s.getNode(resp, req, nodeId)
		}
		return s.listNode(resp, req)
	case "PUT":
		fallthrough
	case "POST":
		if nodeId != "" {
			return s.modifyNode(resp, req)
		}
		return s.addNode(resp, req)
	case "DELETE":
		if nodeId != "" {
			return s.deleteNode(resp, req, nodeId)
		} else {
			return nil, fmt.Errorf("Url '%s' with DELETE is illeage, should like '/node/{nodeId}'", req.URL.Path)
		}

	default:
		resp.WriteHeader(405)
		return nil, nil
	}
}

func (s *HTTPServer) getNode(resp http.ResponseWriter, req *http.Request, nodeId string) (*sched.Node, error) {
	node, err := s.resourceManager.GetNode(nodeId)
	if err != nil {
		return nil, fmt.Errorf("nodeId:%v is not exist.", nodeId)
	}
	return node, nil
}

func (s *HTTPServer) listNode(resp http.ResponseWriter, req *http.Request) ([]*sched.Node, error) {
	nodes, err := s.resourceManager.GetNodeList()
	if err != nil {
		return nil, fmt.Errorf("list node failed, cause %v", err)
	}
	return nodes, nil
}

func (s *HTTPServer) modifyNode(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	// todo
	return nil, nil
}

func (s *HTTPServer) addNode(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, req.Body); err != nil {
		return nil, err
	}
	var meta *sched.NodeMeta
	if err := json.Unmarshal(buf.Bytes(), meta); err != nil {
		return nil, fmt.Errorf("json unmarshal node meta failed: %v", err)
	}

	err := s.resourceManager.AddMeta(meta)
	if err != nil {
		log.Warnf("RM add node meta failed, cause:%v", err)
		return nil, fmt.Errorf("RM add node meta failed, cause: %v", err)
	}
	return nil, nil
}

func (s *HTTPServer) deleteNode(resp http.ResponseWriter, req *http.Request, nodeId string) (interface{}, error) {
	// todo
	return nil, nil
}

/*
func (s *HTTPServer) addNode(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	// Copy the value
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, req.Body); err != nil {
		return nil, err
	}
	iplist := strings.Split(string(buf.Bytes()), ",")

	machineDao, err := s.storeManager.GetMachineDaoForPool()
	if err != nil {
		return nil, err
	}
	defer s.storeManager.PutMachineDaoForPool(machineDao)

	for _, ip := range iplist {
		m := &store.Machine{
			Ip:         ip,
			CreateTime: time.Now(),
		}
		if err := machineDao.Set(m); err != nil {
			return nil, err
		}
	}
	return map[string]interface{}{
		"Success": len(iplist),
	}, nil
}
*/
