package api

import (
	"fmt"
	"net/http"
	// "strconv"
	"bytes"
	"encoding/json"
	log "github.com/omidnikta/logrus"
	"github.com/weibocom/dschedule/structs"
	"io"
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
		log.Infoln("post, add service.....")
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

func (s *HTTPServer) getService(resp http.ResponseWriter, req *http.Request, serviceId string) (*structs.Service, error) {
	service, err := s.serviceManager.GetService(serviceId)
	if err != nil {
		return nil, fmt.Errorf("ServiceManager get service failed, cause: %v", err)
	}
	return service, nil
}

func (s *HTTPServer) listService(resp http.ResponseWriter, req *http.Request) ([]*structs.Service, error) {
	services, err := s.serviceManager.GetServiceList()
	if err != nil {
		return nil, fmt.Errorf("ServiceManager list service failed, cause: %v", err)
	}
	return services, nil
}

func (s *HTTPServer) modifyService(resp http.ResponseWriter, req *http.Request, serviceId string) (*structs.Service, error) {
	service, err := parseService(req)
	if err != nil {
		return nil, fmt.Errorf("parse service from request failed, cause: %v", err)
	}

	// TODO modify service
	s.serviceManager.ModifyService(serviceId, service)
	return nil, nil
}

func (s *HTTPServer) addService(resp http.ResponseWriter, req *http.Request) (string, error) {
	service, err := parseService(req)
	if err != nil {
		return "", fmt.Errorf("parse service from request failed, cause: %v", err)
	}
	log.Infof("add service :%v", service)
	mfService, err := s.serviceManager.AddService(service)
	if err != nil {
		return "", fmt.Errorf("ServiceManager add service failed, cause: %v", err)
	}
	return mfService, nil
}

func (s *HTTPServer) deleteService(resp http.ResponseWriter, req *http.Request, serviceId string) (string, error) {
	ok, err := s.serviceManager.DeleteService(serviceId)
	if err != nil {
		return "", fmt.Errorf("ServiceManager delete service failed, cause: %v", err)
	}
	return ok, nil
}

func parseService(req *http.Request) (*structs.Service, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, req.Body); err != nil {
		return nil, err
	}
	var service *structs.Service
	if err := json.Unmarshal(buf.Bytes(), &service); err != nil {
		return nil, fmt.Errorf("json unmarshal service failed: %v", err)
	}
	return service, nil
}
