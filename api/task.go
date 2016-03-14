package api

import (
	"fmt"
	"net/http"
	// "strconv"
	"strings"
)

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
