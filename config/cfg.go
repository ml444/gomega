package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/ml444/gomega/log"
)

var GlobalCfg *Config
var rwMux = &sync.RWMutex{}

func init() {
	GlobalCfg = NewDefaultConfig()
}

type Config struct {
	RootPath        string                `json:"root_path" yaml:"root_path"`
	MaxMessageSize  int                   `json:"max_message_size" yaml:"max_message_size"`
	FileRotatorSize int64                 `json:"file_rotator_size" yaml:"file_rotator_size"`
	BrokerCfgMap    map[string]*BrokerCfg `json:"broker_cfg_map" yaml:"broker_cfg_map"`
}

func SetBrokerCfg(topic string, brokerCfg *BrokerCfg) {
	rwMux.Lock()
	defer rwMux.Unlock()
	GlobalCfg.BrokerCfgMap[topic] = brokerCfg
}

func GetBrokerCfg(topic string) (*BrokerCfg, bool) {
	rwMux.RLock()
	defer rwMux.RUnlock()
	cfg, ok := GlobalCfg.BrokerCfgMap[topic]
	return cfg, ok
}

func NewDefaultConfig() *Config {
	return &Config{
		RootPath:        ".",
		MaxMessageSize:  1024 * 1024,
		FileRotatorSize: 1024 * 1024 * 256,
		BrokerCfgMap:    map[string]*BrokerCfg{},
	}
}

// FlushConfig flushes the config to file.
func FlushConfig() error {
	filePath := filepath.Join(GlobalCfg.RootPath, "config.json")
	f, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Errorf("err: %v", err)
			return
		}
	}()
	buf, err := json.Marshal(GlobalCfg)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	_, err = f.Write(buf)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	return nil
}
