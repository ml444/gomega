package subscribe

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/ml444/scheduler/pb"
)

//type Config struct {
//
//	// @desc: 并发消费数
//	//ConcurrentCount    uint32 `protobuf:"varint,1,opt,Id=concurrent_count,json=concurrentCount" json:"concurrent_count,omitempty"`
//	MaxRetryCount      uint32 `protobuf:"varint,2,opt,Id=max_retry_count,json=maxRetryCount" json:"max_retry_count,omitempty"`
//	MaxExecTimeSeconds uint32 `protobuf:"varint,3,opt,Id=max_exec_time_seconds,json=maxExecTimeSeconds" json:"max_exec_time_seconds,omitempty"`
//	// @desc: 消费者的服务名和path
//	//ServiceName string `protobuf:"bytes,4,opt,Id=service_name,json=serviceName" json:"service_name,omitempty"`
//	//ServicePath string `protobuf:"bytes,5,opt,Id=service_path,json=servicePath" json:"service_path,omitempty"`
//	// @desc: 在队列中最大存活时间，超时会抛弃
//	//  默认: 一直有效
//	MaxInQueueTimeSeconds uint32 `protobuf:"varint,6,opt,Id=max_in_queue_time_seconds,json=maxInQueueTimeSeconds" json:"max_in_queue_time_seconds,omitempty"`
//	RetryIntervalMs       int64  `protobuf:"varint,8,opt,Id=retry_interval_max,json=retryIntervalMax" json:"retry_interval_max,omitempty"`
//	RetryIntervalStep     int64  `protobuf:"varint,9,opt,Id=retry_interval_step,json=retryIntervalStep" json:"retry_interval_step,omitempty"`
//	// @desc: 串行模式
//	BarrierMode bool `protobuf:"varint,20,opt,Id=barrier_mode,json=barrierMode" json:"barrier_mode,omitempty"`
//	// @desc: 串行数，串行模式下，默认串行数是1,可以配置大于1
//	BarrierCount uint32 `protobuf:"varint,21,opt,Id=barrier_count,json=barrierCount" json:"barrier_count,omitempty"`
//}

type SubConfig struct {
	//IGroup
	Type          pb.Policy
	Version       string
	Namespace     string
	Topic         string
	GroupName     string
	MaxRetryCount uint32
	MaxTimeout    uint32
	//SerialCount         uint32
	RetryIntervalMs     uint32
	ItemLifetimeInQueue uint32
	//Cfg   *Config

}

func (s *SubConfig) Init() {
	//s.Id = s.GenerateId()
	//s.IGroup = NewConsumeGroup(s.Type)
}

func (s *SubConfig) GenerateId() string {
	str := fmt.Sprintf("%s:%s:%s:%d:%d:%d:%d",
		s.Namespace, s.Topic, s.GroupName,
		s.MaxRetryCount, s.MaxTimeout,
		s.RetryIntervalMs, s.ItemLifetimeInQueue)
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *SubConfig) Name() string {
	return GetSubCfgName(s.Namespace, s.Topic, s.GroupName)
}

func GetSubCfgName(namespace, topic, group string) string {
	return fmt.Sprintf("%s:%s:%s", namespace, topic, group)
}
