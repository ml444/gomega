package topics

import (
	"errors"
	"strings"
	"sync"
)

var mgr *Manager
func init() {
	mgr = &Manager{
		topicMap: map[string]*Topic{},
		mu:       sync.RWMutex{},
	}
}
type Manager struct {
	topicMap map[string]*Topic
	//workerMap map[string]*subscribe.Worker
	mu       sync.RWMutex
}

func generateTopicKey(namespace, topicName string) string {
	return strings.Join([]string{namespace, topicName}, ".")
}

func AddTopic(t *Topic) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	mgr.topicMap[generateTopicKey(t.Namespace, t.Name)] = t
}


func GetTopic(namespace, topicName string) (*Topic, error) {
	if t, ok := mgr.topicMap[generateTopicKey(namespace, topicName)]; ok {
		return t, nil
	}
	return nil, errors.New("not found topics")
}


func GetTopicPartitionCount(namespace, topicName string) (uint32, error) {
	if t, ok := mgr.topicMap[generateTopicKey(namespace, topicName)]; ok {
		return t.Partitions, nil
	}
	return 0, errors.New("not found topic")
}

//func (m *Manager) GetWorker(token string) (*subscribe.Worker, error){
//	if worker, ok := m.workerMap[token]; ok {
//		return worker, nil
//	}
//	return nil, errors.New("not found worker")
//}

func TODO() {
	// 根据Topic本身的优先级，选择一个合适的服务节点，更新etcd
	// sdk：topics
}