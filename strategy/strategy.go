package strategy

import (
	"fmt"
	log "github.com/omidnikta/logrus"
	"github.com/weibocom/dschedule/scheduler"
	"github.com/weibocom/dschedule/structs"
)

var BuiltinStrategy = map[string]Factory{
	structs.ServiceStrategyCrontab: NewCrontabStrategy,
	structs.ServiceStrategyAuto:    NewAutoStrategy,
	structs.ServiceStrategyStable:  NewStableStrategy,
}

type Factory func(service *structs.Service) Strategy

type Strategy interface {
	Applying(service *structs.Service, scheduler *scheduler.Scheduler) error
}

func NewStrategy(name string, service *structs.Service) (Strategy, error) {
	factory, ok := BuiltinStrategy[name]
	if !ok {
		return nil, fmt.Errorf(" unsupported strategy: %v", name)
	}
	strategy := factory(service)
	return strategy, nil
}
