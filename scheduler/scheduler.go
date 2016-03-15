package scheduler

import (
	"fmt"
)

type DeployArg struct {
	image   string
	cmd     string
	name    string
	network string // host or bridge ...
	envs    []string
	volumes []string
}

type Service struct {
	min int
	max int

	deployArg DeployArg
}

type Scheduler struct {
	services map[string]interface{}
	//index    int

	resourceManager *ResourceManager
	dockerHost      string
}

func NewScheduler(resourceManager *ResourceManager, dockerHost string) (*Scheduler, error) {
	return &Scheduler{
		services: make(map[string]interface{}),
		//index:           0,
		resourceManager: resourceManager,
		dockerHost:      dockerHost,
	}, nil
}

func (this *Scheduler) Register( /*struct config*/ ) (bool, error) {
	registerId := fmt.Sprintf("registerId-%d", this.index)
	this.index++

	this.services[registerId] = ""
	return true, nil
}

func (this *Scheduler) Add(serviceId string, num int) error {
	// TODO asynchronous
	// TODO check if > max
	// TODO request resource from rm
	// TODO deploy
	return nil
}

func (this *Scheduler) Remove(registerId string, num int) error {
	return nil
}

func (this *Scheduler) Status(registerId string) /*struct runtime status */ error {
	return nil
}
