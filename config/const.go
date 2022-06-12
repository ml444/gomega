package config

const (
	MaxDataSize = 1024 * 1024
	MaxFileSize = 1024 * 1024 * 1024

	DefaultQueueMaxSize uint64 = 1024
)

const (
	IdxFileSuffix  = ".idx"
	DataFileSuffix = ".dat"
)

const (
	DefaultMaxTimeout          = 60
	DefaultMaxRetryCount       = 3
	DefaultRetryIntervalMs     = 5
	DefaultItemLifetimeInQueue = 600
)
