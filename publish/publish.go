package publish

import (
	"context"
	"github.com/ml444/scheduler/brokers"
	"github.com/ml444/scheduler/config"
)


type PubReq struct {
	Namespace string
	TopicName string // @v: required
	Partition uint32
	Data      []byte // @v: required
	HashCode  uint64
	// @desc: delay 仅对 barrier mode 生效
	//  delay_type == DelayTypeRelate 时，delay_value 的单位是 ms
	DelayType  uint32
	DelayValue uint32
	// @desc: [仅串行化模式有效] 范围是 0-3，数值越大，越优先执行
	Priority uint32 `protobuf:"varint,32,opt,name=priority" json:"priority,omitempty"`
}

func Pub(ctx context.Context, namespace, topic string, partition uint32, item *brokers.Item) error {
	if namespace == "" {
		namespace = config.DefaultNamespace
	}
	// file sequence offset
	bk, err := brokers.GetBroker(namespace, topic, partition)
	if err != nil {
		return err
	}

	err = bk.Send(item)
	if err != nil {
		return err
	}
	return nil
}