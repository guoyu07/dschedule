package api

import (
	log "github.com/omidnikta/logrus"
	sched "github.com/weibocom/dschedule/scheduler"
	// "net/http"
	// "net/http/httptest"
	"testing"
)

var server *HTTPServer

func MakeHTTPServer(t *testing.T) *HTTPServer {
	if server != nil {
		return server
	}
	resourceManager := &sched.ResourceManager{}
	var err error
	server, err = NewHTTPServer("0.0.0.0", 11989, "", true, resourceManager)
	if err != nil {
		log.Errorf("create http server failed, cause: %v", err)
		log.Fatalf("create http server failed, cause: %v", err)
	}
	go server.Start()
	log.Infoln("http server started.")
	return server
}
