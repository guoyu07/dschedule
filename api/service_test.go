package api

import (
	"bytes"
	"encoding/json"
	"github.com/weibocom/dschedule/structs"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestService_ADD(t *testing.T) {
	TestNodeEndpoint_ADD(t)
	srv := MakeHTTPServer(t)
	defer srv.Shutdown()

	service := map[string]interface{}{
		"serviceId":    "feed-1",
		"serviceType":  structs.ServiceTypeProd,
		"strategyName": structs.ServiceStrategyCrontab,
		"strategyConfig": []map[string]interface{}{
			map[string]interface{}{
				"time":        "12:10:00",
				"instanceNum": 0,
			},
			// map[string]interface{}{
			// 	"time":        "@every 2s",
			// 	"instanceNum": 3,
			// },
		},
		"priority": 2,
		"container": map[string]interface{}{
			"type":    "docker",
			"image":   "docker.io/redis:2.8",
			"command": "redis-server",
		},
		"dedicated": 0,
		"elastic":   2,
	}

	data, err := json.Marshal(service)
	if err != nil {
		t.Fatalf("json marshal service failed, cause: %v", err)
	}
	t.Log(string(data))
	req, err := http.NewRequest("POST", "/service/", bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	t.Logf("url: %v", req.URL.String())
	resp := httptest.NewRecorder()
	info, err := srv.ServiceEndpoint(resp, req)
	if err != nil {
		t.Fatalf("add node failed, cause:%v", err)
	}
	t.Logf("add service success, info: %v", info)
	time.Sleep(time.Second * 90)
}
