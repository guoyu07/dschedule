package scheduler

import (
	"fmt"
	log "github.com/omidnikta/logrus"
	//"github.com/samalba/dockerclient"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	PROTOCOL_TCP       = "tcp"
	PROTOCOL_HTTP      = "http"
	PROTOCOL_CONTAINER = "container"
)

type HealthConfig struct {
	protocol string // tcp , http, container
	ip       string // TODO should check in every localhost
	port     int
	path     string // only for protocol container(id), http(url path)
	httpTag  string // should contain this string when http check

	interval int // second
	timeout  int

	startWait  int
	maxFailure int

	//startedNotify chan struct{}
	//unhealthyFunc func() bool // bool -> need continue health check
}

func StartHealthCheck(config *HealthConfig, startedNotify chan<- struct{},
	unhealthyFunc func(*HealthConfig) bool) {
	ticker := time.NewTicker(time.Second * 1)
	var isFirstHealth = true
	var unhealthyCount int
	for _ = range ticker.C {
		var isHealth bool
		switch config.protocol {
		case PROTOCOL_TCP:
			isHealth = checkTcp(config.ip, config.port)
		case PROTOCOL_HTTP:
			isHealth = checkHttp(config.ip, config.port, config.path, config.httpTag)
		case PROTOCOL_CONTAINER:
			isHealth = checkContainer(config.ip, config.port, config.path)
		default:
			log.Errorf("config.protocol=%s not support", config.protocol)
			return
		}
		if isHealth {
			if isFirstHealth {
				if startedNotify != nil { // block when send to nil channel
					startedNotify <- struct{}{}
				}
				isFirstHealth = false
			}
			unhealthyCount = 0
		} else {
			if unhealthyCount++; unhealthyCount >= config.maxFailure &&
				!(isFirstHealth && unhealthyCount < config.startWait) {
				needContinue := unhealthyFunc(config)
				if !needContinue {
					break
				}
			}
		}
	}
}

func checkTcp(ip string, port int) bool {
	service := fmt.Sprintf("%s:%d", ip, port)
	_, err := net.Dial("tcp", service)
	return err == nil
}

func checkHttp(ip string, port int, path, httpTag string) bool {
	// Refer: http://www.01happy.com/golang-http-client-get-and-post/
	service := fmt.Sprintf("http://%s:%d/%s", ip, port, path)
	resp, err := http.Get(service)
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false
		}
		//fmt.Println("body", string(body))
		if httpTag != "" {
			return strings.Contains(string(body), httpTag)
		}
		return true
	} else {
		//fmt.Println("err", err)
		return false
	}
}

// TODO /containers/(id)/stats now accepts stream  only >= v1.19 API
func checkContainer(ip string, port int, containerId string) bool {
	return true
}
