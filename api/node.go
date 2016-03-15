package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/omidnikta/logrus"
	"github.com/weibocom/dschedule/structs"
	"io"
	"net/http"
	"strings"
)

func (s *HTTPServer) NodeEndpoint(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	log.Infof("http request url: %v, method: %v", req.URL.String(), req.Method)
	nodeId := strings.TrimPrefix(req.URL.Path, "/node/")
	// Switch on the method
	log.Infoln(nodeId)
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
			return s.modifyNode(resp, req, nodeId)
		}
		log.Infoln("add node")
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

func (s *HTTPServer) getNode(resp http.ResponseWriter, req *http.Request, nodeId string) (*structs.Node, error) {
	node, err := s.resourceManager.GetNode(nodeId)
	if err != nil {
		return nil, fmt.Errorf("nodeId:%v is not exist.", nodeId)
	}
	return node, nil
}

func (s *HTTPServer) listNode(resp http.ResponseWriter, req *http.Request) ([]*structs.Node, error) {
	nodes, err := s.resourceManager.GetNodeList()
	if err != nil {
		return nil, fmt.Errorf("list node failed, cause %v", err)
	}
	return nodes, nil
}

func (s *HTTPServer) modifyNode(resp http.ResponseWriter, req *http.Request, nodeId string) (interface{}, error) {
	meta, err := parseNodeMeta(req)
	if err != nil {
		return nil, fmt.Errorf("parse NodeMeta from request failed, cause: %v", err)
	}

	errModify := s.resourceManager.ModifyMeta(nodeId, meta)
	if errModify != nil {
		return nil, fmt.Errorf("RM modify node meta failed, cause: %v", errModify)
	}
	return nil, nil
}

func (s *HTTPServer) addNode(resp http.ResponseWriter, req *http.Request) (string, error) {
	meta, err := parseNodeMeta(req)
	if err != nil {
		return nil, fmt.Errorf("parse NodeMeta from request failed, cause: %v", err)
	}
	log.Debugf("resource :%v", s.resourceManager)
	nodeId, err := s.resourceManager.AddMeta(meta)
	if err != nil {
		log.Warnf("RM add node meta failed, cause:%v", err)
		return nil, fmt.Errorf("RM add node meta failed, cause: %v", err)
	}
	log.Infof("node id: %v", nodeId)
	return nodeId, nil
}

func (s *HTTPServer) deleteNode(resp http.ResponseWriter, req *http.Request, nodeId string) (string, error) {
	err := s.resourceManager.RemoveNode(nodeId)
	if err != nil {
		log.Warnf("RM remove node failed, cause:%v", err)
		return nil, fmt.Errorf("RM remove node failed, cause: %v", err)
	}
	return "SUCCESS", nil
}

func parseNodeMeta(req *http.Request) (*structs.NodeMeta, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, req.Body); err != nil {
		return nil, err
	}
	var meta *structs.NodeMeta
	if err := json.Unmarshal(buf.Bytes(), &meta); err != nil {
		return nil, fmt.Errorf("json unmarshal node meta failed: %v", err)
	}
	return meta, nil
}
