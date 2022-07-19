package config

var GlobalCfg *Config

type Config struct {
	Broker struct {
		BasePath     string
		QueueMaxSize uint64
	}
}

func NewConfig() *Config {
	return &Config{
		Broker: struct {
			BasePath     string
			QueueMaxSize uint64
		}{
			BasePath:     ".",
			QueueMaxSize: 1024,
		},
	}
}

func init() {
	GlobalCfg = NewConfig()
}
