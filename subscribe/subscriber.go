package subscribe

type Subscriber struct {
	Name  string
	Topic string
	Route string
	Addrs []string
	Cfg   *Config
}

func NewSubscriber(namespace, topicName string, subCfg *Config) *Subscriber {
	return &Subscriber{
		Cfg: subCfg,
	}
}
