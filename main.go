package main

import (
	"os"
	"os/signal"
	"syscall"

	glog "github.com/ml444/glog"
	gcfg "github.com/ml444/glog/config"
	"github.com/ml444/glog/level"

	"github.com/ml444/gomega/broker"
	"github.com/ml444/gomega/config"
	"github.com/ml444/gomega/log"
	"github.com/ml444/gomega/transport/grpc"
)

func main() {
	// init log
	glog.InitLog(gcfg.SetLevel2Logger(level.DebugLevel))
	log.SetLogger(glog.GetLogger())

	// init config
	err := config.InitConfig()
	if err != nil {
		log.Error(err.Error())
		return
	}
	// start some broker
	broker.InitBroker()

	go func() {
		log.Info("grpc server running...")
		err := grpc.RunServer()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	// signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	log.Debugf("signal: %s", <-c)
	broker.ExitBrokers()

}
