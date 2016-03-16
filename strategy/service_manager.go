package strategy

import (
	"fmt"
	log "github.com/omidnikta/logrus"
	"github.com/weibocom/dschedule/scheduler"
	"github.com/weibocom/dschedule/structs"
)

type ServiceManager struct {
	scheduler       *scheduler.Scheduler
	strategyName    string
	strategy        Strategy
	services        map[string]*structs.Service
	resourceManager *scheduler.ResourceManager
}

func NewServiceManager(strategyName string, resourceManager *scheduler.ResourceManager) (*ServiceManager, error) {
	// TODO scheduler dockerPort
	scheduler, _ := scheduler.NewScheduler(resourceManager, 4243)
	strategy, _ := NewStrategy(strategyName)

	return &ServiceManager{
		scheduler:       scheduler,
		strategy:        strategy,
		services:        make(map[string]*structs.Service),
		resourceManager: resourceManager,
	}, nil
}

func (serviceManager *ServiceManager) GetService(serviceId string) (*structs.Service, error) {
	// serviceManager.scheduler.
	return nil, nil
}

func (serviceManager *ServiceManager) AddService(service *structs.Service) (string, error) {
	serviceManager.services[service.ServiceId] = service

	serviceManager.setServiceDefaultProperties(service)
	ok, err := serviceManager.scheduler.Register(service)
	if !ok {
		return "", fmt.Errorf("scheduler register a service failed, cause: %v", err)
	}
	// TODO use strategy
	errApp := serviceManager.strategy.Applying(service, serviceManager.scheduler)
	if errApp != nil {
		return "", fmt.Errorf("strategy applyling service failed, cause : %v", errApp)
	}
	// TODO return
	return service.ServiceId, nil
}

func (serviceManager *ServiceManager) ModifyService(serviceId string, service *structs.Service) (bool, error) {
	serviceManager.setServiceDefaultProperties(service)
	ok, err := serviceManager.scheduler.Register(service)
	if !ok {
		return false, fmt.Errorf("scheduler register a service failed, cause: %v", err)
	}
	if _, ok := serviceManager.services[serviceId]; ok {
		serviceManager.services[serviceId] = service
		err := serviceManager.strategy.Applying(service, serviceManager.scheduler)
		if err != nil {
			return false, fmt.Errorf("strategy applyling service failed, cause : %v", err)
		}
	}
	log.Infof("modify service, service info: %v", service)
	return true, nil
}

func (serviceManager *ServiceManager) DeleteService(serviceId string) (string, error) {
	delete(serviceManager.services, serviceId)
	return "", nil
}

func (serviceManager *ServiceManager) GetServiceList() ([]*structs.Service, error) {
	return nil, nil
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
