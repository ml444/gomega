package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/ml444/gomega/log"
)

const (
	EnvGomegaConfigPath = "GOMEGA_CONFIG_PATH"
	GomegaConfigName    = "gomega_config.json"
)

var configFilePath string
var globalCfg *Config
var rwMux = &sync.RWMutex{}

func InitConfig() error {
	globalCfg = NewDefaultConfig()
	// load config from file
	configFilePath = filepath.Join(os.Getenv(EnvGomegaConfigPath), GomegaConfigName)
	f, err := os.Open(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warnf("config file not exist, use default config")
			return nil
		}
		log.Errorf("err: %v", err)
		return err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(globalCfg)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	return nil
}

type Config struct {
	GRPCServerNetwork string `json:"grpc_server_network" yaml:"grpc_server_network"`
	GRPCServerAddress string `json:"grpc_server_address" yaml:"grpc_server_address"`

	RootPath        string                `json:"root_path" yaml:"root_path"`
	MaxMessageSize  int                   `json:"max_message_size" yaml:"max_message_size"`
	FileRotatorSize int64                 `json:"file_rotator_size" yaml:"file_rotator_size"`
	BrokerCfgMap    map[string]*BrokerCfg `json:"broker_cfg_map" yaml:"broker_cfg_map"`
}

func SetBrokerCfg(topic string, brokerCfg *BrokerCfg) {
	rwMux.Lock()
	defer rwMux.Unlock()
	globalCfg.BrokerCfgMap[topic] = brokerCfg
}

func GetBrokerCfg(topic string) (*BrokerCfg, bool) {
	rwMux.RLock()
	defer rwMux.RUnlock()
	cfg, ok := globalCfg.BrokerCfgMap[topic]
	return cfg, ok
}

func GetRootPath() string {
	return globalCfg.RootPath
}

func GetMaxMessageSize() int {
	return globalCfg.MaxMessageSize
}

func GetFileRotatorSize() int64 {
	return globalCfg.FileRotatorSize
}

func GetGRPCServerNetworkAndAddress() (string, string) {
	return globalCfg.GRPCServerNetwork, globalCfg.GRPCServerAddress
}

func NewDefaultConfig() *Config {
	return &Config{
		GRPCServerNetwork: "tcp",
		GRPCServerAddress: ":12345",

		RootPath:        ".",
		MaxMessageSize:  1024 * 1024,
		FileRotatorSize: 1024 * 1024 * 256,
		BrokerCfgMap:    map[string]*BrokerCfg{},
	}
}

// FlushConfig flushes the config to file.
func FlushConfig() error {
	f, err := os.OpenFile(configFilePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
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
	buf, err := json.Marshal(globalCfg)
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
