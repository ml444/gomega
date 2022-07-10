package main

import (
	"fmt"
	log "github.com/ml444/glog"
	"github.com/ml444/scheduler/brokers"
	"github.com/ml444/scheduler/server"
)

func main() {
	// get server config
	// start some broker
	brokers.InitBroker()
	fmt.Println("running...")
	err := server.RunServer()
	if err != nil {
		fmt.Println(err)
		log.Error(err)
	}
}
