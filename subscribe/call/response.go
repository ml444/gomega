package call
type ConsumeRsp struct {
	// @desc: 是否需要重试
	Retry bool `protobuf:"varint,1,opt,name=retry" json:"retry,omitempty"`
	// @desc: 重试间隔时间，可以这里指定，不指定的话，由 smq 来定 (单位:秒)
	RetryIntervalSeconds int64 `protobuf:"varint,2,opt,name=retry_interval_seconds,json=retryIntervalSeconds" json:"retry_interval_seconds,omitempty"`
	// @desc: 为 true 时，不要累加重试次数
	SkipIncRetryCount bool `protobuf:"varint,3,opt,name=skip_inc_retry_count,json=skipIncRetryCount" json:"skip_inc_retry_count,omitempty"`
	// @desc: 支持毫秒级重试时间
	RetryIntervalMs int64 `protobuf:"varint,4,opt,name=retry_interval_ms,json=retryIntervalMs" json:"retry_interval_ms,omitempty"`
	// @desc: 是否异步消费
	IsAsync bool `protobuf:"varint,5,opt,name=is_async,json=isAsync" json:"is_async,omitempty"`
	// @desc: 异步消费最大等待确认时间
	MaxAsyncWaitMs int64 `protobuf:"varint,6,opt,name=max_async_wait_ms,json=maxAsyncWaitMs" json:"max_async_wait_ms,omitempty"`
	// @desc: 消费端服务cpu使用率
	CpuPercent float32 `protobuf:"fixed32,7,opt,name=cpu_percent,json=cpuPercent" json:"cpu_percent,omitempty"`
	// @desc: 消费端服务内存使用率
	MemoryPercent        float32 `protobuf:"fixed32,8,opt,name=memory_percent,json=memoryPercent" json:"memory_percent,omitempty"`
	Took                 int64
}