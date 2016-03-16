package strategy

import (
	"fmt"
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

func (serviceManager *ServiceManager) ModifyService(serviceId string, service *structs.Service) (*structs.Service, error) {
	return nil, nil
}

func (serviceManager *ServiceManager) DeleteService(serviceId string) (string, error) {
	return "", nil
}

func (serviceManager *ServiceManager) GetServiceList() ([]*structs.Service, error) {
	return nil, nil
}
