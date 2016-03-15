package crontab

import (
	"fmt"
	log "github.com/omidnikta/logrus"
	"github.com/robfig/cron"
	"github.com/weibocom/dschedule/scheduler"
	"github.com/weibocom/dschedule/structs"
)

type CrontabStrategy struct {
	configs    []*Config
	cronObject cron.Cron
}

type Config struct {
	Time        string
	InstanceNum int
}

func NewCrontabStrategy(service *structs.Service) *CrontabStrategy {
	configs, _ := service.StrategyConfig.([]*Config)
	cn := cron.New()
	cn.Start()
	return &CrontabStrategy{
		configs:    configs,
		cronObject: cn,
	}
}

func (crontabStrategy *CrontabStrategy) Applying(srevice *structs.Service, scheduler *scheduler.Scheduler) error {

	for _, config := range crontabStrategy.configs {
		expression, err := ParseTime(config.Time)
		if err != nil {
			log.Errorf("parse config failed, cause : %v", err)
			continue
		}

		// TODO handle registerId
		// 传入的InstanceNum应该经过计算，现在测试链路直接传入原值
		crontabStrategy.cronObject.AddFunc(expression, func() {
			scheduler.Add("registerId", config.InstanceNum)
		})
	}
	return nil
}
