package config

type SubCfg struct {
	Topic          string
	Group          string
	Version        uint32
	LastSequence   uint64
	IsHashDispatch bool
}
