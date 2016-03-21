package strategy

import (
	"fmt"
	log "github.com/omidnikta/logrus"
	"github.com/robfig/cron"
	"github.com/weibocom/dschedule/scheduler"
	"github.com/weibocom/dschedule/structs"
	"reflect"
)

type CrontabStrategy struct {
	configs    []*Config
	cronObject *cron.Cron
}

type Config struct {
	Time        string
	InstanceNum int
}

func NewCrontabStrategy() Strategy {
	cn := cron.New()
	cn.Start()
	return &CrontabStrategy{
		cronObject: cn,
	}
}

func (crontabStrategy *CrontabStrategy) Applying(service *structs.Service, scheduler *scheduler.Scheduler) error {
	var configs []*Config

	switch reflect.TypeOf(service.StrategyConfig).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(service.StrategyConfig)

		for i := 0; i < s.Len(); i++ {
			c, ok := s.Index(i).Interface().(map[string]interface{})
			if !ok {
				log.Errorf("service strategyConfig assertion failed, cause : %v", c)
				continue
			}
			log.Infof("config: %v", c)

			config := &Config{
				Time:        c["time"].(string),
				InstanceNum: (int)(c["instanceNum"].(float64)),
			}
			configs = append(configs, config)
		}
	case reflect.Map:
		s := reflect.ValueOf(service.StrategyConfig)
		c, ok := s.Interface().(map[string]interface{})
		if !ok {
			log.Errorf("service strategyConfig assertion failed, cause : %v", c)
		}

		config := &Config{
			Time:        c["time"].(string),
			InstanceNum: (int)(c["instanceNum"].(float64)),
		}
		configs = append(configs, config)
	default:
		log.Errorf("service StrategyConfig: %v is not the valid object", service.StrategyConfig)
	}

	crontabStrategy.configs = configs
	for _, config := range configs {
		expression, err := ParseTime(config.Time)
		if err != nil {
			log.Errorf("parse config failed, cause : %v", err)
			continue
		}
		log.Infof("run right now .....")
		// just run it immediately
		crontabStrategy.executableFunc(service, config, scheduler)
		log.Infof("run right now already: %v", config)

		crontabStrategy.cronObject.AddFunc(expression, crontabStrategy.crontabFunc(service, config, scheduler))

	}
	return nil
}

func (crontabStrategy *CrontabStrategy) crontabFunc(service *structs.Service, config *Config, scheduler *scheduler.Scheduler) func() {
	return func() {
		err := crontabStrategy.executableFunc(service, config, scheduler)
		if err != nil {
			log.Errorf("cronatab strategy execute function failed, cause : %v", err)
		}
	}
}

func (crontabStrategy *CrontabStrategy) executableFunc(service *structs.Service, config *Config, scheduler *scheduler.Scheduler) error {
	log.Infoln("start run cron job....")
	onlineNum, err := crontabStrategy.getServiceOnlineInstanceNum(service.ServiceId, scheduler)
	if err != nil {
		return fmt.Errorf("scheduler get service:%v status faield, cause: %v", service.ServiceId, err)
	}
	log.Infof("onlineNum:%v, config.InstanceNum:%v", onlineNum, config.InstanceNum)
	if onlineNum > config.InstanceNum {
		num, err := scheduler.Remove(service.ServiceId, onlineNum-config.InstanceNum)
		if err != nil {
			return fmt.Errorf("scheduler remove service:%v failed, cause: %v", service.ServiceId, err)
		}
		log.Infof("scheduler remove success, remove instance num:%v", num)
	} else if onlineNum < config.InstanceNum {
		_, err := scheduler.Add(service.ServiceId, config.InstanceNum-onlineNum)
		if err != nil {
			return fmt.Errorf("scheduler add service:%v failed, cause: %v", service.ServiceId, err)
		}
		log.Infof("scheduler add service:%v success, add instance num:%v, online instance num:%v", service.ServiceId,
			config.InstanceNum-onlineNum, config.InstanceNum)
	} else {
		log.Warnf("crontab strategy check online instance num equals config.InstanceNum, onlineInstanceNum:%v", onlineNum)
	}
	return nil
}

func (crontabStrategy *CrontabStrategy) getServiceOnlineInstanceNum(serviceId string, scheduler *scheduler.Scheduler) (int, error) {
	_, num, err := scheduler.Status(serviceId)
	if err != nil {
		return -1, fmt.Errorf("crontab strategy get service:%v online instance num failed, cause: %v", serviceId, err)
	}
	return num, nil
}
