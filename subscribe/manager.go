package subscribe

import (
	"reflect"
	"sync"
)

var SubMgr *Manager

func init() {
	SubMgr = NewManager()
}

type Manager struct {
	subCfgMap map[string]*SubConfig
	groupMap  map[string]IGroup
	workerMap map[string]*Worker
	mu        sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		subCfgMap: map[string]*SubConfig{},
		groupMap: map[string]IGroup{},
		workerMap: map[string]*Worker{},
		mu:        sync.RWMutex{},
	}
}

func (m *Manager) GetToken(key string) string {
	cfg, ok := m.subCfgMap[key]
	if !ok {
		return ""
	}
	group, ok := m.groupMap[key]
	if !ok {
		group = GetGroup(cfg.Type, cfg)
		m.groupMap[key] = group
	}
	token := m.generateToken(cfg)
	worker := group.AddWorker(token)
	m.workerMap[token] = worker
	go group.Start()
	return token
}

func (m *Manager) generateToken(cfg *SubConfig) string {
	return "test_token"
}
func (m *Manager) UpdateSubCfg(key string, newCfg *SubConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cfg, ok := m.subCfgMap[key]
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
		m.AddSubCfg(newCfg)
	}
}

func (m *Manager) AddSubCfg(sub *SubConfig) {
	//if m.subCfgMap == nil {
	//	m.subCfgMap = map[string]*SubConfig{}
	//}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subCfgMap[sub.Name()] = sub
}

func (m *Manager) GetSubCfg(key string) *SubConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if sub, ok := m.subCfgMap[key]; ok {
		return sub
	} else {
		return nil
	}
}

func (m *Manager) GetWorker(token string) *Worker {
	return m.workerMap[token]
}
