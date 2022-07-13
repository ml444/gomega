package subscribe

import "sync"

var SubMgr *Manager

func init() {
	SubMgr = NewManager()
}

type Manager struct {
	subMap    map[string]*Subscriber
	workerMap map[string]*Worker
	mu        sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		subMap:    map[string]*Subscriber{},
		workerMap: map[string]*Worker{},
		mu:        sync.RWMutex{},
	}
}

func (m *Manager) ReplaceSubscriber(oldKey string, newSub *Subscriber) {
	//if m.subMap == nil {
	//	m.subMap = map[string]*Subscriber{}
	//}
	m.mu.Lock()
	defer m.mu.Unlock()
	if sub, ok := m.subMap[oldKey]; ok {
		sub.Stop()
		delete(m.subMap, oldKey)
	}
	if newSub.Id == "" {
		newSub.Id = newSub.GenerateId()
	}
	m.subMap[newSub.Id] = newSub
}

func (m *Manager) AddSubscriber(sub *Subscriber) {
	//if m.subMap == nil {
	//	m.subMap = map[string]*Subscriber{}
	//}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subMap[sub.Id] = sub
	go sub.Start()
}

func (m *Manager) GetSubscriber(key string) *Subscriber {
	if sub, ok := m.subMap[key]; ok {
		return sub
	} else {
		return nil
	}
}