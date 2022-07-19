package topics

import (
	"errors"
	"fmt"
	"github.com/ml444/scheduler/config"
	"github.com/ml444/scheduler/subscribe"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"sync"
	"syscall"
)

type Topic struct {
	Namespace  string
	Name       string
	Partitions uint32
	Priority   uint32

	subCfgMap map[string]*subscribe.SubConfig
	groupMap  map[string]subscribe.IGroup
	workerMap map[string]*subscribe.Worker
	mu        sync.RWMutex
}

func NewTopic(namespace, topicName string, partitions uint32, priority uint32) (*Topic, error) {
	var err error
	err = checkNaming(namespace)
	if err != nil {
		return nil, err
	}
	err = checkNaming(topicName)
	if err != nil {
		return nil, err
	}
	for i:=uint32(0); i< partitions; i++ {
		err = makeDir(namespace, topicName, i)
		if err != nil {
			return nil, err
		}
	}
	return &Topic{
		Namespace:  namespace,
		Name:       topicName,
		Partitions: partitions,
		Priority:   priority,
		subCfgMap:  map[string]*subscribe.SubConfig{},
		groupMap:   map[string]subscribe.IGroup{},
		workerMap:  map[string]*subscribe.Worker{},
		mu:         sync.RWMutex{},
	}, nil
}

func makeDir(namespace, topic string, partition uint32) error {
	syscall.Umask(0)
	return os.MkdirAll(filepath.Join(config.GlobalCfg.Broker.BasePath,namespace, topic, strconv.FormatUint(uint64(partition), 10)), 0777)
}

func checkNaming(name string) error {
	l := len(name)
	if l == 0 {
		return errors.New("name is empty")
	}
	if l > 64 {
		return fmt.Errorf("the length of the name[%s] is greater than 64", name)
	}
	for i := 0; i < l; i++ {
		if name[i] == '.' || name[i] == '/' {
			return fmt.Errorf("name[%s] has invalid char", name)
		}
	}
	return nil
}
func (t *Topic) GetToken(groupName string) (token string, err error) {
	cfg, ok := t.subCfgMap[groupName]
	if !ok {
		return "", errors.New("not found subscriber config")
	}
	group, ok := t.groupMap[groupName]
	if !ok {
		group, err = subscribe.NewConsumeGroup(cfg.Type, cfg, t.Partitions)
		if err != nil {
			return "", err
		}
		t.groupMap[groupName] = group
	}
	token = t.generateToken(cfg)
	worker := group.AddWorker(token)
	t.workerMap[token] = worker
	go group.Start()
	return token, nil
}

func (t *Topic) generateToken(cfg *subscribe.SubConfig) string {
	return "test_token"
}
func (t *Topic) UpdateSubCfg(groupName string, newCfg *subscribe.SubConfig) {
	t.mu.Lock()
	defer t.mu.Unlock()
	cfg, ok := t.subCfgMap[groupName]
	if ok {
		if !reflect.DeepEqual(cfg, newCfg) {
			if n := newCfg.MaxRetryCount; n != 0 && n != cfg.MaxRetryCount {
				cfg.MaxRetryCount = n
			}
			if n := newCfg.MaxTimeout; n != 0 && n != cfg.MaxTimeout {
				cfg.MaxTimeout = n
			}
			if n := newCfg.RetryIntervalMs; n != 0 && n != cfg.RetryIntervalMs {
				cfg.RetryIntervalMs = n
			}
			if n := newCfg.ItemLifetimeInQueue; n != 0 && n != cfg.ItemLifetimeInQueue {
				cfg.ItemLifetimeInQueue = n
			}
		}
	} else {
		t.AddSubCfg(newCfg)
	}
}

func (t *Topic) AddSubCfg(sub *subscribe.SubConfig) {
	//if t.subCfgMap == nil {
	//	t.subCfgMap = map[string]*SubConfig{}
	//}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.subCfgMap[sub.GroupName] = sub
}

func (t *Topic) GetSubCfg(group string) *subscribe.SubConfig {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if sub, ok := t.subCfgMap[group]; ok {
		return sub
	} else {
		return nil
	}
}

func (t *Topic) GetWorker(token string) (*subscribe.Worker, error) {
	if worker, ok := t.workerMap[token]; ok {
		return worker, nil
	}
	return nil, errors.New("not found worker")
}
