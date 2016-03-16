package api

import (
	"bytes"
	"encoding/json"
	//	"github.com/weibocom/dschedule/scheduler"
	"github.com/weibocom/dschedule/structs"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TNodeEndpoint_ADD(t *testing.T) {
	srv := MakeHTTPServer(t)
	// defer srv.Shutdown()
	// if err != nil {
	// 	t.Fatalf("create httpserver failed, cause:%v", err)
	// }
	meta := &structs.NodeMeta{
		Name:     "yf-remind-1",
		IP:       "10.229.88.91",
		CPU:      24,
		MemoryMB: 10240,
		DiskMB:   102400,
	}
	metaJson, err := json.Marshal(meta)
	t.Logf("meta json: %v", string(metaJson))
	if err != nil {
		t.Fatalf("json marshal meta failed, cause:%v", err)
	}
	req, err := http.NewRequest("POST", "/node/", bytes.NewBuffer(metaJson))
	if err != nil {
		t.Fatalf("create http request failed, cause:%v", err)
	}
	t.Logf("url: %v", req.URL.String())
	resp := httptest.NewRecorder()
	info, err := srv.NodeEndpoint(resp, req)
	if err != nil {
		t.Fatalf("add node failed, cause:%v", err)
	}
	t.Logf("add node success, info: %v", info)
	time.Sleep(time.Second * 1)
}

/*
func TestNodeEndpoint_GET(t *testing.T) {
	TestNodeEndpoint_ADD(t)
	time.Sleep(time.Second * 1)
	srv := MakeHTTPServer(t)
	defer srv.Shutdown()

	req, err := http.NewRequest("GET", "/node/10.75.24.12-random-something-else", nil)
	if err != nil {
		t.Fatalf("create http request failed, cause:%v", err)
	}
	resp := httptest.NewRecorder()
	ret, err := srv.NodeEndpoint(resp, req)
	if err != nil {
		t.Fatalf("get node failed, cause:%v", err)
	}
	t.Logf("get node success, ret: %v", ret)
}

func TestNodeEndpoint_REMOVE(t *testing.T) {
	TestNodeEndpoint_ADD(t)
	time.Sleep(time.Second * 1)
	srv := MakeHTTPServer(t)
	defer srv.Shutdown()

	req, err := http.NewRequest("DELETE", "/node/10.75.24.12-random-something-else", nil)
	if err != nil {
		t.Fatalf("create http request failed, cause:%v", err)
	}
	resp := httptest.NewRecorder()
	ret, err := srv.NodeEndpoint(resp, req)
	if err != nil {
		t.Fatalf("remove node failed, cause:%v", err)
	}
	t.Logf("remove node success, ret: %v", ret)
}
*/
