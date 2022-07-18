package topic

import (
	"errors"
	"strings"
	"sync"
)

var Mgr *Manager
func init() {
	Mgr = &Manager{
		topicMap: map[string]*Topic{},
		mu:       sync.RWMutex{},
	}
}
type Manager struct {
	topicMap map[string]*Topic
	//workerMap map[string]*subscribe.Worker
	mu       sync.RWMutex
}




func (m *Manager) GetTopic(namespace, topicName string) (*Topic, error) {
	if t, ok := m.topicMap[strings.Join([]string{namespace, topicName}, ".")]; ok {
		return t, nil
	}
	return nil, errors.New("not found topic")
}

//func (m *Manager) GetWorker(token string) (*subscribe.Worker, error){
//	if worker, ok := m.workerMap[token]; ok {
//		return worker, nil
//	}
//	return nil, errors.New("not found worker")
//}

func TODO() {
	// 根据Topic本身的优先级，选择一个合适的服务节点，更新etcd
	// sdk：topic
}