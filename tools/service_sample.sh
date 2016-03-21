#!/bin/sh

ip="localhost"
port=11989


#新增离线服务
curl -X POST -d '{
	"serviceId"		: 	"hadoop-job-1",
	"serviceType"	:	"nonprod",
	"strategyName" : "crontab",
	"strategyConfig":[
		{
			"time":"23:10:00",
			"instanceNum":2
		}
	],
	"priority"		:	1,
	"dedicated"		: 	1,
	"elastic"		:   1,
	"container" : {
		"type" : "docker",
		"network":"host",
			"image": "registry.intra.weibo.com/tongwei/node_manager:cdh5.3.2",
    	"env":{
   	    	"YARN_DIRS":"/data0/yarn:/data1/yarn:/data2/yarn:/data3/yarn"
	    },
	    "portMapping":{
	      "8088":"8088",
	      "8031":"8031",
	      "8031":"8031",
	      "8032":"8032",
	      "8033":"8033",
	      "57829":"57829"
	    },
	    "volumes":{
	      "/data0/yarn":"/data0/yarn",
	      "/data1/yarn":"/data1/yarn",
	      "/data2/yarn":"/data2/yarn",
	      "/data3/yarn":"/data3/yarn",
	      "/etc/hadoop/conf":"/etc/hadoop/conf"
	    }
	}
}' "http://$ip:$port/service/"


#变更离线服务
curl -X PUT -d '{
	"serviceId"		: 	"hadoop-job-1",
	"serviceType"	:	"nonprod",
	"strategyName" : "crontab",
	"strategyConfig":[
		{
			"time":"23:10:00",
			"instanceNum":1
		}
	],
	"priority"		:	1,
	"dedicated"		: 	1,
	"elastic"		:   1,
	"container" : {
		"type" : "docker",
		"network":"host",
			"image": "registry.intra.weibo.com/tongwei/node_manager:cdh5.3.2",
    	"env":{
   	    	"YARN_DIRS":"/data0/yarn:/data1/yarn:/data2/yarn:/data3/yarn"
	    },
	    "portMapping":{
	      "8088":"8088",
	      "8031":"8031",
	      "8031":"8031",
	      "8032":"8032",
	      "8033":"8033",
	      "57829":"57829"
	    },
	    "volumes":{
	      "/data0/yarn":"/data0/yarn",
	      "/data1/yarn":"/data1/yarn",
	      "/data2/yarn":"/data2/yarn",
	      "/data3/yarn":"/data3/yarn",
	      "/etc/hadoop/conf":"/etc/hadoop/conf"
	    }
	}
}' "http://$ip:$port/service/"

#新增在线服务
curl -X POST -d '{
	"serviceId"		: 	"remind-web",
	"serviceType"	:	"prod",
	"strategyName" : "crontab",
	"strategyConfig":[
		{
			"time":"23:10:00",
			"instanceNum":1
		}
	],
	"priority"		:	5,
	"dedicated"		: 	1,
	"elastic"		:   1,
	"container" : {
		"type" : "docker",
		"network":"host",
		"name" : "blossom",
		"image": "registry.intra.weibo.com/weibo_rd_if/remind-web:remind-web_RELEASE_V3.49",
    	"env":{
   	    	"NAME_CONF":"openapi_remind-tc-inner=/data1/weibo"
	    },
	    "volumes":{
	      "/data1/mblog/logs/"	:	"/data1/weibo/logs/",
	      "/data1/mblog/gclogs/":	"/data1/weibo/gclogs"
	    },
	    "command":"/docker_init.sh"
	}
}' "http://$ip:$port/service/"

#更新在线服务
curl -X PUT -d '{
	"serviceId"		: 	"remind-web",
	"serviceType"	:	"prod",
	"strategyName" : "crontab",
	"strategyConfig":[
		{
			"time":"23:10:00",
			"instanceNum":2
		}
	],
	"priority"		:	5,
	"dedicated"		: 	1,
	"elastic"		:   1,
	"container" : {
		"type" : "docker",
		"network":"host",
		"name" : "blossom",
		"image": "registry.intra.weibo.com/weibo_rd_if/remind-web:remind-web_RELEASE_V3.49",
    	"env":{
   	    	"NAME_CONF":"openapi_remind-tc-inner=/data1/weibo"
	    },
	    "volumes":{
	      "/data1/mblog/logs/"	:	"/data1/weibo/logs/",
	      "/data1/mblog/gclogs/":	"/data1/weibo/gclogs"
	    },
	    "command":"/docker_init.sh"
	}
}' "http://$ip:$port/service/"

#查询服务
curl "http://$ip:$port/service/{serviceId}"

#查询所有服务
curl "http://$ip:$port/service/"