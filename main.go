package main

import (
	glog "github.com/ml444/glog"
	"github.com/ml444/glog/config"
	"github.com/ml444/glog/level"

	"github.com/ml444/gomega/broker"
	"github.com/ml444/gomega/log"
	"github.com/ml444/gomega/transport/grpc"
)

func main() {
	glog.InitLog(config.SetLevel2Logger(level.DebugLevel))
	log.SetLogger(glog.GetLogger())
	// get server config
	// start some broker
	broker.InitBroker()
	log.Info("running...")
	err := grpc.RunServer()
	if err != nil {
		log.Error(err.Error())
	}
}
