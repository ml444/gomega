package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/ml444/scheduler/brokers"
	"github.com/ml444/scheduler/config"
	"github.com/ml444/scheduler/publish"
	"github.com/ml444/scheduler/subscribe"
	"github.com/ml444/scheduler/subscribe/call"
	"google.golang.org/grpc"
	"net"
	"time"
)

type OmegaServer struct {
}

func (s *OmegaServer) Pub(ctx context.Context, req *PubRequest) (*PubResponse, error) {
	var rsp PubResponse
	fmt.Println(req)

	dataSize := len(req.Data)
	if dataSize > config.MaxDataSize {
		return nil, errors.New("data length exceeds limit")
	}
	item := brokers.Item{
		CreatedAt: uint32(time.Now().Unix()),
		HashCode:  req.HashCode,
		Partition: req.Partition,
		Size:      uint32(dataSize),
		//DelayType:  req.DelayType,
		//DelayValue: req.DelayValue,
		//Priority:   req.Priority,
		Data:      req.Data,
		Namespace: req.Namespace,
	}

	err := publish.Pub(&ctx, req.Namespace, req.TopicName, &item)
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}

func (s *OmegaServer) Subscribe(ctx context.Context, req *SubscribeReq) (*SubscribeRsp, error) {
	var rsp SubscribeRsp
	fmt.Println("====> Subscribe")
	sub := subscribe.Subscriber{
		Namespace:     req.Namespace,
		Name:          req.Name,
		Topic:         req.Topic,
		Route:         req.Route,
		Addrs:         nil,
		Cfg:           nil,
		Request:       req.Request,
		Response:      req.Response,
		BeforeProcess: nil,
		AfterProcess:  nil,
	}
	subscribe.SubMgr.SetSubscriber(&sub)
	rsp.Status = 10000
	rsp.Message = "successfully"

	{
		// test
		item := &brokers.Item{
			Sequence:   1,
			HashCode:   1234,
			CreatedAt:  123,
			Partition:  0,
			Offset:     1,
			Size:       0,
			RetryCount: 0,
		}
		payload := brokers.MsgPayload{
			MsgId:      "124560",
			Data:       []byte(`{"Latitude": 409146138, "Longitude"": -746188906}`),
			DelayType:  0,
			DelayValue: 0,
		}
		consumeRsp := call.ConsumeRsp{}
		worker := subscribe.Worker{S: &sub, Cfg: &subscribe.Config{MaxExecTimeSeconds: 123}}
		err := worker.ConsumeMsg(item, &payload, &consumeRsp)
		if err != nil {
			println(err)
		}
	}
	return &rsp, nil
}

func RunServer() error {
	listen, err := net.Listen("tcp", ":12345")
	if err != nil {
		return err
	}
	svr := grpc.NewServer()
	RegisterOmegaServiceServer(svr, &OmegaServer{})
	err = svr.Serve(listen)
	if err != nil {
		return err
	}
	return nil
}
