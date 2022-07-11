package subscribe

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/ml444/scheduler/subscribe/call"
	"reflect"
)

type Subscriber struct {
	Namespace string
	Name      string
	Topic     string
	Route     string
	Addrs     []string
	Cfg       *Config

	Request  proto.Message
	Response proto.Message

	BeforeProcess func(ctx context.Context, meta *call.MsgMeta)
	AfterProcess  func(ctx context.Context, meta *call.MsgMeta, req, rsp interface{}) (isRetry bool, ignoreRetryCount bool)
}

func NewSubscriber(namespace, topicName string, subCfg *Config) *Subscriber {
	return &Subscriber{
		Cfg: subCfg,
	}
}

func (s *Subscriber) UnMarshalRequest(data []byte) (proto.Message, error) {
	inT := reflect.TypeOf(s.Request).Elem()
	in := reflect.New(inT).Interface().(proto.Message)
	err := proto.Unmarshal(data, in)
	if err != nil {
		return nil, err
	}
	return in, nil
}

func (s Subscriber) NewResponse() proto.Message {
	T := reflect.TypeOf(s.Response).Elem()
	return reflect.New(T).Interface().(proto.Message)
}
