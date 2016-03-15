package strategy

import (
	// "fmt"
	log "github.com/omidnikta/logrus"
	"github.com/robfig/cron"
	"github.com/weibocom/dschedule/scheduler"
	"github.com/weibocom/dschedule/structs"
	"reflect"
	"time"
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

		// TODO handle registerId
		// 传入的InstanceNum应该经过计算，现在测试链路直接传入原值
		crontabStrategy.cronObject.AddFunc(expression, func() {
			scheduler.Add(service.ServiceId, config.InstanceNum)
		})
	}
	time.Sleep(time.Second * 10)
	return nil
}
