package subscribe

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/ml444/scheduler/pb"
)

type Subscriber struct {
	IGroup
	Type                pb.Policy
	//Version             string
	Id                  string
	Namespace           string
	Topic               string
	GroupName           string
	MaxRetryCount       uint32
	MaxTimeout          uint32
	//SerialCount         uint32
	RetryIntervalMs     uint32
	ItemLifetimeInQueue uint32
	//Cfg   *Config

}

func (s *Subscriber) Init() {
	s.Id = s.GenerateId()
	s.IGroup = GetGroup(s.Type)
}

func (s *Subscriber) GetToken() string {
	// TODO
	return "token"
}

func (s *Subscriber) GenerateId() string {
	str := fmt.Sprintf("%d:%s:%s:%s:%d:%d:%d:%d",
		s.Type, s.Namespace, s.Topic, s.GroupName,
		s.MaxRetryCount, s.MaxTimeout,
		s.RetryIntervalMs, s.ItemLifetimeInQueue)
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//func (s *Subscriber) UnMarshalRequest(data []byte) (proto.Message, error) {
//	inT := reflect.TypeOf(s.Request).Elem()
//	in := reflect.New(inT).Interface().(proto.Message)
//	err := proto.Unmarshal(data, in)
//	if err != nil {
//		return nil, err
//	}
//	return in, nil
//}
//
//func (s Subscriber) NewResponse() proto.Message {
//	T := reflect.TypeOf(s.Response).Elem()
//	return reflect.New(T).Interface().(proto.Message)
//}
