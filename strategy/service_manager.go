package strategy

import (
	"encoding/json"
	"fmt"
	log "github.com/omidnikta/logrus"
	"github.com/weibocom/dschedule/scheduler"
	"github.com/weibocom/dschedule/storage"
	"github.com/weibocom/dschedule/structs"
)

const (
	ServiceStoragePrefix = "service"
)

type ServiceManager struct {
	scheduler       *scheduler.Scheduler
	strategyName    string
	strategy        Strategy
	services        map[string]*structs.Service
	resourceManager *scheduler.ResourceManager
	storage         *storage.Storage
}

func NewServiceManager(strategyName string, resourceManager *scheduler.ResourceManager, storage *storage.Storage) (*ServiceManager, error) {
	// TODO scheduler dockerPort
	scheduler, _ := scheduler.NewScheduler(resourceManager, 4243)
	strategy, _ := NewStrategy(strategyName)

	return &ServiceManager{
		scheduler:       scheduler,
		strategy:        strategy,
		services:        make(map[string]*structs.Service),
		resourceManager: resourceManager,
		storage:         storage,
	}, nil
}

func (serviceManager *ServiceManager) AddService(service *structs.Service) (string, error) {
	if _, ok := serviceManager.services[service.ServiceId]; ok {
		return "", fmt.Errorf("add service failed, cause service:%v already exist.", service.ServiceId)
	}
	//set default service properties
	serviceManager.setServiceDefaultProperties(service)
	serviceManager.services[service.ServiceId] = service

	err := serviceManager.StoreService(service)
	if err != nil {
		return "", fmt.Errorf("store service failed, cause: %v", err)
	}

	ok, err := serviceManager.scheduler.Register(service)
	if !ok {
		return "", fmt.Errorf("scheduler register a service failed, cause: %v", err)
	}

	errApp := serviceManager.strategy.Applying(service, serviceManager.scheduler)
	if errApp != nil {
		return "", fmt.Errorf("strategy applyling service failed, cause : %v", errApp)
	}

	return service.ServiceId, nil
}

func (serviceManager *ServiceManager) ModifyService(serviceId string, service *structs.Service) (bool, error) {
	_, ok := serviceManager.services[serviceId]
	if !ok {
		serv, _ := serviceManager.GetService(serviceId)
		if serv == nil {
			return false, fmt.Errorf("modify service failed, cause service:%v is not exist.", serviceId)
		}
	}

	serviceManager.setServiceDefaultProperties(service)
	serviceManager.services[serviceId] = service
	err := serviceManager.StoreService(service)
	if err != nil {
		return false, fmt.Errorf("store service failed, cause: %v", err)
	}

	ok, err = serviceManager.scheduler.Register(service)
	if !ok {
		return false, fmt.Errorf("scheduler register a service failed, cause: %v", err)
	}

	err = serviceManager.strategy.Applying(service, serviceManager.scheduler)
	if err != nil {
		return false, fmt.Errorf("strategy applyling service failed, cause : %v", err)
	}

	log.Infof("modify service, service info: %v", service)
	return true, nil
}

func (serviceManager *ServiceManager) DeleteService(serviceId string) (string, error) {
	delete(serviceManager.services, serviceId)
	err := serviceManager.RemoveService(serviceId)
	if err != nil {
		return "", fmt.Errorf("remove service from storage failed, cause: %v", err)
	}
	return "SUCCESS", nil
}

func (serviceManager *ServiceManager) GetService(serviceId string) (*structs.Service, error) {
	service, ok := serviceManager.services[serviceId]
	if !ok {
		services, err := serviceManager.RetriveService(serviceId)
		if err != nil {
			return nil, fmt.Errorf("service manager get service failed, cause: %v", err)
		}
		if len(services) > 0 {
			serviceManager.services[serviceId] = services[0]
			return services[0], nil
		}
		return nil, fmt.Errorf("get service failed, service is not exist serviceId:%v", serviceId)
	}
	return service, nil
}

func (serviceManager *ServiceManager) GetServiceList() ([]*structs.Service, error) {
	var services []*structs.Service
	for _, service := range serviceManager.services {
		services = append(services, service)
	}
	var err error
	if len(services) == 0 {
		services, err = serviceManager.RetriveService("")
		if err != nil {
			return nil, fmt.Errorf("service manager get service list faield, cause : %v", err)
		}
		for _, service := range services {
			serviceManager.services[service.ServiceId] = service
		}
	}
	return services, nil
}

func (serviceManager *ServiceManager) setServiceDefaultProperties(service *structs.Service) {
	if service.Priority == 0 {
		service.Priority = 3
	}
	container := service.Container
	if container == nil {
		container = &structs.Container{
			Type:    structs.ContainerTypeDocker,
			Network: structs.ContainerNetworkHost,
		}
		service.Container = container
	} else {
		if container.Type == "" {
			container.Type = structs.ContainerTypeDocker
		}
		if container.Network == "" {
			container.Network = structs.ContainerNetworkHost
		}
	}
}

/**
* save service to storage
 */
func (serviceManager *ServiceManager) StoreService(service *structs.Service) error {
	value, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("json marshal service failed, cause: %v", err)
	}
	key := fmt.Sprintf("%s/%s", ServiceStoragePrefix, service.ServiceId)
	err = serviceManager.storage.Client.Put(key, []byte(value), nil)
	if err != nil {
		return fmt.Errorf("Error trying to put value at key: %v", key)
	}
	return nil
}

/**
* remove service from storage
 */
func (serviceManager *ServiceManager) RemoveService(serviceId string) error {
	key := fmt.Sprintf("%s/%s", ServiceStoragePrefix, serviceId)
	err := serviceManager.storage.Client.Delete(key)
	if err != nil {
		return fmt.Errorf("Error trying to delete key: %v", key)
	}
	return nil
}

/**
* retrive service from storage
 */
func (serviceManager *ServiceManager) RetriveService(serviceId string) ([]*structs.Service, error) {
	var services []*structs.Service
	key := fmt.Sprintf("%s/%s", ServiceStoragePrefix, serviceId)
	if serviceId == "" { // retrive list
		entries, err := serviceManager.storage.Client.List(key)
		if err != nil {
			return nil, fmt.Errorf("storage get service list failed, cause: %v", err)
		}
		for _, pair := range entries {
			var service *structs.Service
			if err := json.Unmarshal(pair.Value, &service); err != nil {
				log.Errorf("json unmarshal service failed, cause: %v", err)
				continue
			}
			log.Infof("retrive service : %v", service)
			services = append(services, service)
		}
	} else {
		pair, err := serviceManager.storage.Client.Get(key)
		if err != nil {
			return nil, fmt.Errorf("Error trying to get key: %v, cause: %v", key, err)
		}
		var service *structs.Service
		if err := json.Unmarshal(pair.Value, &service); err != nil {
			return nil, fmt.Errorf("json unmarshal service failed, cause: %v", err)
		}
		log.Infof("retrive service : %v", service)
		services = append(services, service)
	}

	return services, nil
}
