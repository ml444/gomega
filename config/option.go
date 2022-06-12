package config

import "github.com/ml444/gomega/log"

type Option func(c *Config)

func WithLogger(logger log.Logger) Option {
	return func(c *Config) {
		log.SetLogger(logger)
	}
}
