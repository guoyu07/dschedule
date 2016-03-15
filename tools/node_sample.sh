#!/bin/sh

#节点新增，删除，查询接口脚本
host="localhost"
ip=""

# 添加新的节点
curl -X POST -d '{
	"name"		: 	"",
	"ip"		:	"",
	"cpu"		:	"12",
	"memoryMB"	: 	"10240",
	"diskMB"	:	"51200",
	"attibutes" : {
		"hostname" : "",
		"os" 	   : ""
	}
}' "http://$host:$ip/node"


# 变更节点信息
curl -X PUT -d '{
	"name"		: 	"",
	"ip"		:	"",
	"cpu"		:	"12",
	"memoryMB"	: 	"10240",
	"diskMB"	:	"51200",
	"attibutes" : {
		"hostname" : "",
		"os" 	   : ""
	}
}' "http://$host:$ip/node/{nodeId}"


# 查询已存在的节点信息
curl "http://$host:$ip/node/{nodeId}"

# 查询所有存在的节点信息
curl "http://$host:$ip/node/"

# 删除某一个节点
curl -X DELETE "http://$host:$ip/node/{nodeId}"
