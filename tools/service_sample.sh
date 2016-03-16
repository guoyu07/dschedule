#!/bin/sh

ip="localhost"
port=


#新增服务
curl -X POST -d '{
	"serviceId"		: 	"test-redis",
	"serviceType"		:	"prod",
	"strategyName" : "crontab",
	"strategyConfig":[{
		"time":"",
		"instanceNum":1
	}],
	"priority"		:	5,
	"container" : {
		"type" : "docker",
		"image": "docker.io/redis:2.8",
		"command":"redis-server"
	}

}' "http://$host:$ip/service/"


#变更服务
curl -X POST -d '{
	"serviceId"		: 	"test-redis",
	"serviceType"		:	"prod",
	"strategyName" : "crontab",
	"strategyConfig":[{
		"time":"",
		"instanceNum":1
	}],
	"priority"		:	5,
	"container" : {
		"type" : "docker",
		"image": "docker.io/redis:2.8",
		"command":"redis-server"
	}

}' "http://$host:$ip/service/"