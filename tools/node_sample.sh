#!/bin/sh

#节点新增，删除，查询接口脚本
ip="localhost"
port="11989"

# 添加新的节点
curl -X POST -d '{
	"name"		: 	"10.73.88.40",
	"ip"		:	"10.73.88.40",
	"cpu"		:	12,
	"memoryMB"	: 	10240,
	"diskMB"	:	51200
}' "http://$ip:$port/node/"


# 新增节点信息
curl -X POST -d '{
	"name"		: 	"10.73.88.41",
	"ip"		:	"10.73.88.41",
	"cpu"		:	12,
	"memoryMB"	: 	10240,
	"diskMB"	:	51200
}' "http://$ip:$port/node/"

# 变更节点信息
curl -X PUT -d '{
	"name"		: 	"10.73.88.41",
	"ip"		:	"10.73.88.41",
	"cpu"		:	12,
	"memoryMB"	: 	10240,
	"diskMB"	:	51200
}' "http://$ip:$port/node/{nodeId}"


# 查询已存在的节点信息
curl "http://$ip:$port/node/{nodeId}"

# 查询所有存在的节点信息
curl "http://$ip:$port/node/"

# 删除某一个节点
curl -X DELETE "http://$ip:$port/node/{nodeId}"
