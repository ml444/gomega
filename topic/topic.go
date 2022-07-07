package topic

import (
	log "github.com/ml444/glog"
	"github.com/ml444/scheduler/subscribe"
	"sync"
)

func init() {
	topicMgr = &Manager{
		m:  map[string]*Topic{},
		mu: sync.RWMutex{},
	}
}

type Config struct {
	Name      string
	NameSpace string
}

type Topic struct {
	Namespace string
	Name      string
	Partition int
	//cfg           *Config
	subscribers   map[string]*subscribe.Subscriber
	subscribersMu sync.RWMutex

	Priority uint32
}

func NewTopic(t *Config) (*Topic, error) {
	log.Infof("new topic %+v", t)

	x := &Topic{

	}

	return x, nil
}

func (p *Topic) getSubscriber(name string) *subscribe.Subscriber {
	p.subscribersMu.RLock()
	x := p.subscribers[name]
	p.subscribersMu.RUnlock()
	return x
}
