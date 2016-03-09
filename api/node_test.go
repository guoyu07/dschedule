package api

import (
	"bytes"
	"encoding/json"
	"github.com/weibocom/dschedule/scheduler"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNodeEndpoint_ADD(t *testing.T) {
	srv := MakeHTTPServer(t)
	defer srv.Shutdown()
	// if err != nil {
	// 	t.Fatalf("create httpserver failed, cause:%v", err)
	// }
	meta := &scheduler.NodeMeta{
		Name:     "yf-remind-1",
		IP:       "10.75.24.12",
		CPU:      24,
		MemoryMB: 10240,
		DiskMB:   102400,
	}
	metaJson, err := json.Marshal(meta)
	t.logf("meta json: %v", string(metaJson))
	if err != nil {
		t.Fatalf("json marshal meta failed, cause:%v", err)
	}
	req, err := http.NewRequest("POST", "/node", bytes.NewBuffer(metaJson))
	if err != nil {
		t.Fatalf("create http request failed, cause:%v", err)
	}
	resp := httptest.NewRecorder()
	info, err := srv.NodeEndpoint(resp, req)
	if err != nil {
		t.Fatalf("add node failed, cause:%v", err)
	}
	t.Logf("add node success, info: %v", info)
}

func TestNodeEndpoint_GET(t *testing.T) {
	srv := MakeHTTPServer(t)
	defer srv.Shutdown()

	req, err := http.NewRequest("GET", "/node/<nodeId>", nil)
	if err != nil {
		t.Fatalf("create http request failed, cause:%v", err)
	}
	resp := httptest.NewRecorder()
	ret, err := srv.NodeEndpoint(resp, req)
	if err != nil {
		t.Fatalf("add node failed, cause:%v", err)
	}
	t.Logf("add node success, ret: %v", ret)
}
