package api

import (
	log "github.com/omidnikta/logrus"
	sched "github.com/weibocom/dschedule/scheduler"
	// "net/http"
	// "net/http/httptest"
	"testing"
)

func MakeHTTPServer(t *testing.T) *HTTPServer {
	resourceManager := &sched.ResourceManager{}
	server, err := NewHTTPServer("0.0.0.0", 11989, "", true, resourceManager)
	if err != nil {
		log.Errorf("create http server failed, cause: %v", err)
		t.Fatalf("create http server failed, cause: %v", err)
	}
	go server.Start()
	log.Infoln("http server started.")
	return server
}
