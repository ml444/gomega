package subscribe

import (
	"context"
	"encoding/json"
	"github.com/ml444/scheduler/subscribe/call"
	"reflect"
)

type Subscriber struct {
	Name  string
	Topic string
	Route string
	Addrs []string
	Cfg   *Config

	Request interface{}
	Response interface{}

	BeforeProcess func(ctx context.Context, meta *call.MsgMeta)
	AfterProcess func(ctx context.Context, meta *call.MsgMeta, req, rsp interface{}) (isRetry bool, ignoreRetryCount bool)
}

func NewSubscriber(namespace, topicName string, subCfg *Config) *Subscriber {
	return &Subscriber{
		Cfg: subCfg,
	}
}

func (s *Subscriber) UnMarshalRequest(data []byte) (interface{}, error) {
	var in = s.Request
	err := json.Unmarshal(data, &in)
	if err != nil {
		return nil, err
	}
	return in, nil
}

func (s Subscriber) NewResponse() interface{} {
	T := reflect.TypeOf(s.Response)
	return reflect.New(T).Interface()
}