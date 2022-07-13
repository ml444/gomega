package subscribe

import "sync"

var SubMgr *Manager

func init() {
	SubMgr = NewManager()
}

type Manager struct {
	subMap map[string]*Subscriber
	mu     sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		subMap: map[string]*Subscriber{},
		mu:     sync.RWMutex{},
	}
}

func (m *Manager) UpdateSubConfig(name string, cfg *Config) {
	if m.subMap == nil {
		m.subMap = map[string]*Subscriber{}
	}
	if sub, ok := m.subMap[name]; ok {
		oldCfg := sub.Cfg
		if oldCfg.BarrierMode != cfg.BarrierMode && cfg.BarrierMode != false {

		}
	}
}

func (m *Manager) SetSubscriber(sub *Subscriber) {
	if m.subMap == nil {
		m.subMap = map[string]*Subscriber{}
	}
	m.subMap[sub.Id] = sub
}
