package config

type BrokerCfg struct {
	QueueMaxSize   uint64
	PubCfg         *PubCfg
	GroupSubCfgMap map[string]*SubCfg
}

func (b *BrokerCfg) SetGroupSubCfg(group string, cfg *SubCfg) {
	rwMux.Lock()
	defer rwMux.Unlock()
	b.GroupSubCfgMap[group] = cfg
	// save config to file
	_ = FlushConfig()
}

func (b *BrokerCfg) GetGroupSubCfg(group string) (*SubCfg, bool) {
	rwMux.RLock()
	defer rwMux.RUnlock()
	cfg, ok := b.GroupSubCfgMap[group]
	return cfg, ok
}
