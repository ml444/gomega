package main

import (
	"reflect"

	"github.com/golang/protobuf/proto"

	"github.com/ml444/gomega/config"
)

type Subscriber struct {
	Name  string
	Topic string
	Route string
	Addrs []string
	Cfg   *config.Config

	Request  proto.Message
	Response proto.Message
}

func NewSubscriber(subCfg *config.Config) *Subscriber {
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

func (s *Subscriber) NewResponse() proto.Message {
	T := reflect.TypeOf(s.Response).Elem()
	return reflect.New(T).Interface().(proto.Message)
}
