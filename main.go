/* *
 * Date: 2016.03.08
 * */

package main

import (
	"flag"
	//"fmt"
	log "github.com/omidnikta/logrus"
	"github.com/weibocom/dschedule/api"
	"github.com/weibocom/dschedule/scheduler"
	"github.com/weibocom/dschedule/strategy"
	//"time"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	textFormatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(textFormatter)
	//log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	// log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

}

// run consul first:
// consul agent -server -bootstrap -data-dir /tmp/consul -client=0.0.0.0 -ui-dir=/data0/consul_ui/
func main() {
	var debug = flag.Bool("debug", true, "enable debug")
	//var cluster = flag.Bool("cluster", false, "for online job server")
	var port = flag.Int("port", 11989, "listen port")
	var uiDir = flag.String("ui-dir", "", "ui directory")
	var storage = flag.String("storage", "consul://localhost:8500", "backend to store meta data")
	var storageKeyPrefix = flag.String("storage-key-prefix", "dschedule", "storage key prefix")
	flag.Parse()

	log.Infof("Dschedule start with: storage=%v, prefix=%s, listen-port=%d, debug=%v, ui-dir=%v .",
		*storage, storageKeyPrefix, *port, *debug, *uiDir)

	// for compile test
	resourceManager, _ := scheduler.NewResourceManager()

	serviceManager, _ := strategy.NewServiceManager("CRONTAB", resourceManager)

	server, err := api.NewHTTPServer("0.0.0.0", *port, *uiDir, *debug, resourceManager, serviceManager)
	if err != nil {
		log.Errorf("[main.NewHTTPServer] falied: %v", err)
		return
	}

	server.Start()

	log.Fatalln("Dschedule stopped!")
}
