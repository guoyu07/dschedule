package scheduler

import (
	"fmt"
	log "github.com/omidnikta/logrus"
	"github.com/weibocom/dschedule/structs"
	"sync"
	"time"
)

const (
	NoResourceWaitTime = 10 * time.Second
)

type ScheduleService struct {
	service        *structs.Service
	beforeServices []*structs.Service
	afterServices  []*structs.Service
	deployers      []*Deployer
	mutex          sync.Mutex
}

type Scheduler struct {
	services         map[string]*ScheduleService
	priorityServices []map[string]*ScheduleService
	//index    int

	resourceManager *ResourceManager
	dockerPort      int
}

func NewScheduler(resourceManager *ResourceManager, dockerPort int) (*Scheduler, error) {
	// TODO ConstructFromConsul for  services and priorityServices
	return &Scheduler{
		services:         make(map[string]*ScheduleService),
		priorityServices: make([]map[string]*ScheduleService, structs.MaxPriority+1),
		//index:           0,
		resourceManager: resourceManager,
		dockerPort:      dockerPort,
	}, nil
}

func (this *Scheduler) Register(service *structs.Service) (bool, error) {
	if service.Priority < structs.MinPriority || service.Priority > structs.MaxPriority {
		return false, fmt.Errorf("service.Priority %d is should in [%d, %d]",
			service.Priority, structs.MinPriority, structs.MaxPriority)
	}
	var beforeServices, afterServices []*structs.Service
	for _, serviceId := range service.BeforeServiceIds {
		if schedulerService, ok := this.services[serviceId]; ok {
			beforeServices = append(beforeServices, schedulerService.service)
		} else {
			return false, fmt.Errorf("beforeServiceId:%d not register", serviceId)
		}
	}
	for _, serviceId := range service.AfterServiceIds {
		if schedulerService, ok := this.services[serviceId]; ok {
			afterServices = append(afterServices, schedulerService.service)
		} else {
			return false, fmt.Errorf("afterServiceId:%d not register", serviceId)
		}
	}

	scheduleService := &ScheduleService{
		service:        service,
		beforeServices: beforeServices,
		afterServices:  afterServices,
	}
	this.services[service.ServiceId] = scheduleService
	// TODO StoreToConsul

	// insert into priority map, checked service.Priority at the beginning of this func
	if this.priorityServices[service.Priority] == nil {
		this.priorityServices[service.Priority] = make(map[string]*ScheduleService)
	}
	this.priorityServices[service.Priority][service.ServiceId] = scheduleService

	return true, nil
}

// TODO
func (this *Scheduler) Keep(serviceId string, num int) (bool, error) {
	return false, nil
}

// Providing a non-nil stopCh can be used to stop continuing waiting avilable resource
func (this *Scheduler) Add(serviceId string, num int, stopCh <-chan struct{}) (bool, error) {
	log.Infof("invoke scheduler Add....... serviceId:%v, num:%v", serviceId, num)
	// check if > max
	scheduleService := this.services[serviceId]
	if scheduleService == nil {
		return false, fmt.Errorf("serviceId %d not Register before", serviceId)
	}

	service := scheduleService.service
	if service.Dedicated+service.Elastic < len(scheduleService.deployers)+num {
		return false, fmt.Errorf("Add num %d(existing %d) has been larger than settting by register(%d + %d)",
			num, len(scheduleService.deployers), service.Dedicated, service.Elastic)
	}

	var containers []*structs.Container
	for _, service := range scheduleService.beforeServices {
		containers = append(containers, service.Container)
	}
	containers = append(containers, service.Container)
	for _, service := range scheduleService.afterServices {
		containers = append(containers, service.Container)
	}

	// asynchronous
	//go func() {
	var wg sync.WaitGroup
	for needRequestNum := num; needRequestNum > 0; {
		select {
		case <-stopCh:
			// CARE: maybe need break for, if there are codes after outer for
			return false, fmt.Errorf("Got a stopCh, left %d containers need to expand", needRequestNum)
		default: // stopCh == nil
		}

		wg.Wait()
		allocNodes, err := this.resourceManager.AllocNodes(needRequestNum)
		if err != nil {
			return false, fmt.Errorf("AllocNodes failed: %v, needRequestNum:%d", err, needRequestNum)
		}
		// remove low priority service nodes
		if needKill := needRequestNum - len(allocNodes); needKill > 0 {
			// TODO TODO TODO do wait priority  -> channel
			wg.Add(1)
			go func() {
				defer wg.Done()
				expectKillNum := needKill
			LOOP:
				// search low priority queue and stop them
				for i := structs.MinPriority; i < service.Priority; i++ {
					for serviceId, _ := range this.priorityServices[i] {
						//fmt.Printf("Key: %s  Value: %s\n", key, value)
						// TODO asynchronous
						removedNum, err := this.Remove(serviceId, needKill)
						if err != nil {
							log.Errorf("Remove low priority serviceId %s failed : %v", serviceId, err)
							continue
						}
						needKill -= removedNum
						if needKill <= 0 {
							break LOOP
						}
					}
				}
				if expectKillNum == needKill { // no node killed
					time.Sleep(NoResourceWaitTime)
				}
			}()
		}

		// TODO continue last deploying when process restart, ???tcc????
		// deploy
		var deployers []*Deployer
		for _, node := range allocNodes {
			deployer, err := NewDeployer(node, this.dockerPort, containers)
			if err == nil {
				err = deployer.Start()
			}
			if err != nil {
				node.Failed++
				log.Errorf("NewDeployer or Start container IP '%s' failed %d times: %v",
					node.Meta.IP, node.Failed, err)
				err := this.resourceManager.ReturnNodes([]*structs.Node{node})
				if err != nil {
					log.Errorf("ReturnNodes %s IP '%s' failed: %v", node.NodeId, node.Meta.IP, err)
				}
				continue
			}

			// TODO StoreToConsul for new node-container in this service
			// TODO tell resourceManager, scheduler have used the resource, esle rm will reback the resource
			log.Infof("Started serviceId:%s, nodeId:%s, containerIds:%v", serviceId, deployer.node.NodeId, deployer.containerIds)
			node.Failed = 0
			deployers = append(deployers, deployer)
		}
		scheduleService.mutex.Lock()
		scheduleService.deployers = append(scheduleService.deployers, deployers...)
		scheduleService.mutex.Unlock()

		log.Infof("needRequestNum=%d, available deployer num=%d", needRequestNum, len(deployers))
		needRequestNum -= len(deployers)
	}
	//}()
	return true, nil
}

func (this *Scheduler) Remove(serviceId string, num int) (int, error) {
	log.Infof("deploy remove ...... serviceId:%v, num:%v", serviceId, num)
	scheduleService := this.services[serviceId]
	if scheduleService == nil {
		return 0, fmt.Errorf("serviceId %d not Register before", serviceId)
	}

	if len(scheduleService.deployers) <= scheduleService.service.Dedicated {
		return 0, fmt.Errorf("scheduler remove failed, cause deployers number < Dedicated num, serviceId:%v", serviceId)
	}

	scheduleService.mutex.Lock()
	reduceNum := num
	if num < 0 { // remove all when negative num
		reduceNum = len(scheduleService.deployers)
	}

	log.Debugf("deployers:%v, dedicated:%v, reduceNum:%v", len(scheduleService.deployers), scheduleService.service.Dedicated, reduceNum)
	elasticNum := len(scheduleService.deployers) - scheduleService.service.Dedicated
	if reduceNum > elasticNum {
		reduceNum = elasticNum // make sure the Dedicated
	}
	if reduceNum == len(scheduleService.deployers) {
		scheduleService.deployers = scheduleService.deployers[:0]
	} else {
		scheduleService.deployers = scheduleService.deployers[reduceNum:]
	}
	reduceDeployers := scheduleService.deployers[:reduceNum]
	scheduleService.mutex.Unlock()

	var returnNodes []*structs.Node
	for _, deployer := range reduceDeployers {
		err := deployer.Stop()
		if err != nil {
			log.Errorf("Stop container failed, serviceId:%s: %v", serviceId, err)
		}
		// TODO StoreToConsul for killing node-container in this service
		log.Infof("Stopped serviceId:%s, nodeId:%s, containerIds:%v", serviceId, deployer.node.NodeId, deployer.containerIds)
		// TODO deal with the failed node
		returnNodes = append(returnNodes, deployer.node)
	}

	err := this.resourceManager.ReturnNodes(returnNodes)
	if err != nil {
		log.Errorf("ReturnNodes len=%d failed: %v", len(returnNodes), err)
	}

	return reduceNum, nil
}

func (this *Scheduler) Status(serviceId string) (*structs.Service, int, error) {
	scheduleService := this.services[serviceId]
	if scheduleService == nil {
		return nil, 0, fmt.Errorf("serviceId %d not Register before", serviceId)
	}
	return scheduleService.service, len(scheduleService.deployers), nil
}
