#!/bin/sh

ip="localhost"
port=8500

#启动consul agent
# consul agent -server -bootstrap -dc yf -data-dir /tmp/consul -client=0.0.0.0 -ui-dir=/tmp/consul_ui/ 

#查询某个节点信息
curl "http://$ip:$port/v1/kv/node/{nodeId}"

#查询所有节点信息
curl "http://$ip:$port/v1/kv/node?recurse"

#查询某服务信息
curl "http://$ip:$port/v1/kv/service/{serviceId}"

#查询所有服务信息
curl "http://$ip:$port/v1/kv/service/"