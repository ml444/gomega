package publish

import (
	"context"
	"errors"
	"github.com/ml444/scheduler/brokers"
	"time"
)

const (
	maxDataSize = 1024 * 1024
	defaultNamespace= "default"
)

type PubReq struct {
	Namespace string
	TopicName string // @v: required
	Partition uint32
	Data  []byte // @v: required
	Hash  uint64
	// @desc: delay 仅对 barrier mode 生效
	//  delay_type == DelayTypeRelate 时，delay_value 的单位是 ms
	DelayType  uint32
	DelayValue uint32
	// @desc: [仅串行化模式有效] 范围是 0-3，数值越大，越优先执行
	Priority uint32 `protobuf:"varint,32,opt,name=priority" json:"priority,omitempty"`
}

func Pub(ctx *context.Context, req *PubReq) error {
	if req.Namespace == "" {
		req.Namespace = defaultNamespace
	}
	dataSize := len(req.Data)
	if dataSize > maxDataSize {
		return errors.New("data length exceeds limit")
	}
	// file sequence offset
	bk:= brokers.GetBrokerByTopicName(req.Namespace, req.TopicName)
	item := brokers.Item{
		CreatedAt:  uint32(time.Now().Unix()),
		HashCode:   req.Hash,
		Partition:  req.Partition,
		Size:       uint32(dataSize),
		DelayType:  req.DelayType,
		DelayValue: req.DelayValue,
		Priority:   req.Priority,
		Data:       req.Data,
		Namespace:  req.Namespace,
	}
	err := bk.Send(&item)
	if err != nil {
		return err
	}
	return nil
}