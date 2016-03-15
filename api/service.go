package api

import (
	"fmt"
	"net/http"
	// "strconv"
	"strings"
)

func (s *HTTPServer) ServiceEndpoint(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	serviceId := strings.TrimPrefix(req.URL.Path, "/service/")
	// Switch on the method
	switch req.Method {
	case "GET":
		if serviceId != "" {
			return s.getService(resp, req, serviceId)
		}
		return s.listService(resp, req)
	case "PUT":
		fallthrough
	case "POST":
		if serviceId != "" {
			return s.modifyService(resp, req, serviceId)
		}
		return s.addService(resp, req)
	case "DELETE":
		// Pull out the node id,
		if serviceId != "" {
			return s.deleteService(resp, req, serviceId)
		} else {
			return nil, fmt.Errorf("Url '%s' with DELETE is illeage, should like '/service/{serviceId}'", req.URL.Path)
		}

	default:
		resp.WriteHeader(405)
		return nil, nil
	}
}

func (s *HTTPServer) getService(resp http.ResponseWriter, req *http.Request, serviceId string) (interface{}, error) {
	// todo
	return nil, nil
}

func (s *HTTPServer) listService(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	// todo
	return nil, nil
}

func (s *HTTPServer) modifyService(resp http.ResponseWriter, req *http.Request, serviceId string) (interface{}, error) {
	//todo
	return nil, nil
}

func (s *HTTPServer) addService(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	//todo
	return nil, nil
}

func (s *HTTPServer) deleteService(resp http.ResponseWriter, req *http.Request, serviceId string) (interface{}, error) {
	//todo
	return nil, nil
}
