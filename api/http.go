package api

import (
	"encoding/json"
	"fmt"
	log "github.com/omidnikta/logrus"
	sched "github.com/weibocom/dschedule/scheduler"
	"github.com/weibocom/dschedule/util"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
	"time"
)

// HTTPServer is used to wrap an Agent and expose various API's
// in a RESTful manner
type HTTPServer struct {
	addr            string
	port            int
	uiDir           string
	mux             *http.ServeMux
	listener        net.Listener
	resourceManager *sched.ResourceManager
}

// NewHTTPServers starts new HTTP servers to provide an interface to
// the agent.
func NewHTTPServer(ip string, port int, uiDir string, enableDebug bool, resourceManager *sched.ResourceManager) (*HTTPServer, error) {

	httpAddr, err := util.ClientListener(ip, port)
	if err != nil {
		return nil, fmt.Errorf("Failed to get ClientListener address %s:%d -> %v", ip, port, err)
	}

	// Error if we are trying to bind a domain socket to an existing path
	socketPath, isSocket := util.UnixSocketAddr(ip)
	if isSocket {
		if _, err := os.Stat(socketPath); !os.IsNotExist(err) {
			log.Warnf("http: Replacing socket %q", socketPath)
		}
		if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("error removing socket file: %s", err)
		}
	}

	ln, err := net.Listen(httpAddr.Network(), httpAddr.String())
	if err != nil {
		return nil, fmt.Errorf("Failed to get Listen on %s: %v", httpAddr.String(), err)
	}

	var list net.Listener
	if isSocket {
		// Set up ownership/permission bits on the socket file
		/* TODO open
		if err := util.SetFilePermissions(socketPath, config.UnixSockets); err != nil {
			return nil, fmt.Errorf("Failed setting up HTTP socket: %s", err)
		}
		*/
		list = ln
	} else {
		list = tcpKeepAliveListener{ln.(*net.TCPListener)}
	}

	// Create the mux
	mux := http.NewServeMux()

	// Create the server
	srv := &HTTPServer{
		uiDir: uiDir,
		addr:  httpAddr.String(),

		mux:             mux,
		listener:        list,
		resourceManager: resourceManager,
	}
	srv.registerHandlers(enableDebug)

	return srv, nil
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by NewHttpServer so
// dead TCP connections eventually go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(30 * time.Second)
	return tc, nil
}

func (s *HTTPServer) Start() {
	log.Infof("Http Server listening at addr %s", s.addr)
	// Start the server
	http.Serve(s.listener, s.mux)
}

// Shutdown is used to shutdown the HTTP server
func (s *HTTPServer) Shutdown() {
	if s != nil {
		log.Infoln("http: Shutting down http server (%v)", s.addr)
		s.listener.Close()
	}
}

// registerHandlers is used to attach our handlers to the mux
func (s *HTTPServer) registerHandlers(enableDebug bool) {
	s.mux.HandleFunc("/", s.Index)
	s.mux.HandleFunc("/ping", s.Ping)

	// s.mux.HandleFunc("/status", s.wrap(s.Status))

	// s.mux.HandleFunc("/service/", s.wrap(s.ServiceEndpoint))
	// s.mux.HandleFunc("/repair/", s.wrap(s.RepairEndpoint))
	// s.mux.HandleFunc("/exec/", s.wrap(s.ExecEndpoint))
	// s.mux.HandleFunc("/machine/", s.wrap(s.MachineEndpoint))

	s.mux.HandleFunc("/node/", s.wrap(s.NodeEndpoint))
	s.mux.HandleFunc("/service/", s.wrap(s.ServiceEndpoint))

	if enableDebug {
		s.mux.HandleFunc("/debug/pprof/", pprof.Index)
		s.mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		s.mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		s.mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	}

	// Enable the UI + special endpoints
	if s.uiDir != "" {
		// Static file serving done from /ui/
		s.mux.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir(s.uiDir))))
	}
}

// wrap is used to wrap functions to make them more convenient
func (s *HTTPServer) wrap(handler func(resp http.ResponseWriter, req *http.Request) (interface{}, error)) func(resp http.ResponseWriter, req *http.Request) {
	f := func(resp http.ResponseWriter, req *http.Request) {
		// TODO open
		//setHeaders(resp, s.agent.config.HTTPAPIResponseHeaders)

		/*
			// Obfuscate any tokens from appearing in the logs
			formVals, err := url.ParseQuery(req.URL.RawQuery)
			if err != nil {
				log.Error("http: Failed to decode query: %s", err)
				resp.WriteHeader(500)
				return
			}
		*/

		// Invoke the handler
		start := time.Now()
		defer func() {
			log.Debugf("http: Request %v (%v)", req.URL.String(), time.Now().Sub(start))
		}()
		obj, err := handler(resp, req)

		// Check for an error
	HAS_ERR:
		if err != nil {
			log.Errorf("http: Request %v, error: %v", req.URL.String(), err)
			code := 500
			errMsg := err.Error()
			// TODO change
			if strings.Contains(errMsg, "Permission denied") || strings.Contains(errMsg, "ACL not found") {
				code = 403
			}
			resp.WriteHeader(code)
			resp.Write([]byte(err.Error()))
			return
		}

		//默认使用pretty输出
		prettyPrint := true
		// if _, ok := req.URL.Query()["pretty"]; ok {
		// 	prettyPrint = true
		// }

		// Write out the JSON object
		// return null for no data
		//if obj != nil {
		var buf []byte
		if prettyPrint {
			buf, err = json.MarshalIndent(obj, "", "    ")
		} else {
			buf, err = json.Marshal(obj)
		}
		if err != nil {
			goto HAS_ERR
		}
		resp.Header().Set("Content-Type", "application/json")
		resp.Write(buf)
		//}
	}
	return f
}

// Renders a simple index page
func (s *HTTPServer) Index(resp http.ResponseWriter, req *http.Request) {
	// Check if this is a non-index path
	if req.URL.Path != "/" {
		resp.WriteHeader(404)
		return
	}

	// Check if we have no UI configured
	if s.uiDir == "" {
		resp.Write([]byte("dScheduler"))
		return
	}

	// Redirect to the UI endpoint
	http.Redirect(resp, req, "/ui/", 301)
}

func (s *HTTPServer) Ping(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("PONG"))
}
