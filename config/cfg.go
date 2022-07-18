package config

var GlobalCfg *Config

type Config struct {
	Broker struct {
		BasePath     string
		QueueMaxSize uint64
	}
}

func NewConfig() *Config {
	return &Config{}
}
