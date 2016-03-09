/* *
 * Date: 2016.03.08
 * */

package main

import (
	"flag"
	//"fmt"
	"github.com/weibocom/dschedule/scheduler"
	"log"
	//"time"
)

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

	log.Println("Dschedule start with: storage=%v, prefix=%s, listen-port=%d, debug=%v, ui-dir=%v .",
		*storage, storageKeyPrefix, *port, *debug, *uiDir)

	// for compile test
	_, _ = scheduler.NewResourceManager()

	log.Fatalln("Dschedule stopped!")
}
