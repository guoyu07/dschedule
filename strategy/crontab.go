package strategy

import (
	"fmt"
	log "github.com/omidnikta/logrus"
	"strings"
)

const (
	CRONTAB_EXPRESSION_FORMAT = "* * * * * *" //秒 分 时 天 月 周
	TIME_OF_DAY_SEPERATOR     = ":"

	// SECONDS_OF_MINUTE = 60
	// MINUTE_OF_HOUR    = 60
	// HOUR_OF_DAY       = 24
)

/**
* 将时间转换为crontab的格式，时间格式为：HH:MM:SS 或者MM:SS
 */
func ParseTime(time string) (string, error) {
	var expression string = ""
	if time == "" {
		return "", fmt.Errorf("input time is not valid, time format like HH:MM:SS or MM:SS, input time: %v", time)
	}
	// 对于类似@hourly, @weekly的时间配置，直接返回。
	if strings.HasPrefix(time, "@") {
		log.Infof("input time is %v", time)
		return time, nil
	}
	tod := strings.Split(time, TIME_OF_DAY_SEPERATOR)

	cronExpSectionSize := len(strings.Split(CRONTAB_EXPRESSION_FORMAT, " "))
	if len(tod) > 3 || len(tod) < 2 {
		return "", fmt.Errorf("input time is not valid, time format like HH:MM:SS or MM:SS, input time: %v", time)
	}
	for i := len(tod) - 1; i >= 0; i-- {
		expression += tod[i] + " "
	}
	for i := 0; i < cronExpSectionSize-len(tod); i++ {
		expression += "* "
	}
	log.Infof("parse time: %v to cronatab expression: %v", time, expression)
	return strings.TrimSpace(expression), nil
}
