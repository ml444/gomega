package subscribe

type Config struct {

	// @desc: 并发消费数
	ConcurrentCount    uint32 `protobuf:"varint,1,opt,name=concurrent_count,json=concurrentCount" json:"concurrent_count,omitempty"`
	MaxRetryCount      uint32 `protobuf:"varint,2,opt,name=max_retry_count,json=maxRetryCount" json:"max_retry_count,omitempty"`
	MaxExecTimeSeconds uint32 `protobuf:"varint,3,opt,name=max_exec_time_seconds,json=maxExecTimeSeconds" json:"max_exec_time_seconds,omitempty"`
	// @desc: 消费者的服务名和path
	ServiceName string `protobuf:"bytes,4,opt,name=service_name,json=serviceName" json:"service_name,omitempty"`
	ServicePath string `protobuf:"bytes,5,opt,name=service_path,json=servicePath" json:"service_path,omitempty"`
	// @desc: 在队列中最大存活时间，超时会抛弃
	//  默认: 一直有效
	MaxInQueueTimeSeconds uint32 `protobuf:"varint,6,opt,name=max_in_queue_time_seconds,json=maxInQueueTimeSeconds" json:"max_in_queue_time_seconds,omitempty"`
	RetryIntervalMs  int64 `protobuf:"varint,8,opt,name=retry_interval_max,json=retryIntervalMax" json:"retry_interval_max,omitempty"`
	RetryIntervalStep int64 `protobuf:"varint,9,opt,name=retry_interval_step,json=retryIntervalStep" json:"retry_interval_step,omitempty"`
	// @desc: 串行模式
	BarrierMode bool `protobuf:"varint,20,opt,name=barrier_mode,json=barrierMode" json:"barrier_mode,omitempty"`
	// @desc: 串行数，串行模式下，默认串行数是1,可以配置大于1
	BarrierCount         uint32   `protobuf:"varint,21,opt,name=barrier_count,json=barrierCount" json:"barrier_count,omitempty"`
}
