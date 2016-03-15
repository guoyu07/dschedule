package strategy

import (
	"fmt"
	// log "github.com/omidnikta/logrus"
	"github.com/weibocom/dschedule/scheduler"
	"github.com/weibocom/dschedule/structs"
)

var BuiltinStrategy = map[string]Factory{
	structs.ServiceStrategyCrontab: NewCrontabStrategy,
	// structs.ServiceStrategyAuto:    NewAutoStrategy,
	// structs.ServiceStrategyStable:  NewStableStrategy,
}

type Factory func() Strategy

type Strategy interface {
	Applying(service *structs.Service, scheduler *scheduler.Scheduler) error
}

func NewStrategy(name string) (Strategy, error) {
	factory, ok := BuiltinStrategy[name]
	if !ok {
		return nil, fmt.Errorf(" unsupported strategy: %v", name)
	}
	strategy := factory()
	return strategy, nil
}
