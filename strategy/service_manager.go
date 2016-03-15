package strategy

import (
	"fmt"
	"github.com/weibocom/dschedule/scheduler"
	"github.com/weibocom/dschedule/structs"
)

type ServiceManager struct {
	scheduler    *scheduler.Scheduler
	strategyName string
	strategy     *Strategy
}

func NewServiceManager(strategyName string) {
	// TODO new scheduler
	scheduler := &scheduler.Scheduler{}

	return &ServiceManager{
		scheduler: scheduler,
		strategy:  NewStrategy(strategyName),
	}
}

func (serviceManager *ServiceManager) getService(serviceId string) (*structs.Service, error) {
	// serviceManager.scheduler.
	return nil, nil
}

func (serviceManager *ServiceManager) addService(service *structs.Service) (string, error) {
	registerId, err := serviceManager.scheduler.Register(service)
	if err != nil {
		return "", fmt.Errorf("scheduler register a service failed, cause: %v", err)
	}
	// TODO use strategy
	errApp := serviceManager.strategy.Applying(service)
	if errApp != nil {
		return "", fmt.Errorf("strategy applyling service failed, cause : %v", errApp)
	}
	// TODO return
	return "", nil
}

func (serviceManager *ServiceManager) modifyService(serviceId string, service *structs.Service) (*structs.Service, error) {
	return nil, nil
}

func (serviceManager *ServiceManager) deleteService(serviceId string) (string, error) {
	return nil, nil
}

func (serviceManager *ServiceManager) getServiceList() ([]*structs.Service, error) {
	return nil, nil
}
