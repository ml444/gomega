package subscribe

import (
	log "github.com/ml444/glog"
)

type Subscriber struct {
	Name  string
	Topic string
	Route string
	Addrs []string
	Cfg   *Config
}

func NewSubscriber(namespace, topicName string, subCfg *Config) *Subscriber {
	log.Infof("init subCfg %+v", subCfg)
	return &Subscriber{
		Cfg: subCfg,

	}
}
