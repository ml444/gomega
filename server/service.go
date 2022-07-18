package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/ml444/scheduler/brokers"
	"github.com/ml444/scheduler/config"
	"github.com/ml444/scheduler/pb"
	"github.com/ml444/scheduler/publish"
	"github.com/ml444/scheduler/subscribe"
	"github.com/ml444/scheduler/topic"
	"google.golang.org/grpc"
	"io"
	"net"
	"time"
)

type OmegaServer struct {
	pb.UnimplementedOmegaServiceServer
}

func (s *OmegaServer) Pub(ctx context.Context, req *pb.PubReq) (*pb.PubRsp, error) {
	var rsp pb.PubRsp
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

func (s *OmegaServer) Sub(ctx context.Context, req *pb.SubReq) (*pb.SubRsp, error) {
	var rsp pb.SubRsp
	fmt.Println("====> Subscribe", req.ClientId, "===")
	//cfg := subscribe.SubConfig{
	//	Namespace:           req.Namespace,
	//	Topic:               req.Topic,
	//	GroupName:           req.Group,
	//	MaxRetryCount:       req.MaxRetryCount,
	//	MaxTimeout:          req.MaxTimeout,
	//	RetryIntervalMs:     req.RetryIntervalMs,
	//	ItemLifetimeInQueue: req.ItemLifetimeIn_Queue,
	//}
	//cfg.Init()
	topicIns, err := topic.Mgr.GetTopic(req.Namespace, req.Topic)
	if err != nil {
		return nil, err
	}
	if c := topicIns.GetSubCfg(req.Group); c == nil {
		cfg := subscribe.SubConfig{
			Type:                req.Policy,
			Namespace:           req.Namespace,
			Topic:               req.Topic,
			GroupName:           req.Group,
			MaxRetryCount:       config.DefaultMaxRetryCount,
			MaxTimeout:          config.DefaultMaxTimeout,
			RetryIntervalMs:     config.DefaultRetryIntervalMs,
			ItemLifetimeInQueue: config.DefaultItemLifetimeInQueue,
		}
		topicIns.AddSubCfg(&cfg)
	}
	rsp.Status = 10000
	rsp.Token, err = topicIns.GetToken(req.Group)
	if err != nil {
		return nil, err
	}
	//{
	//	// test
	//	item := &brokers.Item{
	//		Sequence:   1,
	//		HashCode:   1234,
	//		CreatedAt:  123,
	//		Partition:  0,
	//		Offset:     1,
	//		Size:       0,
	//		RetryCount: 0,
	//	}
	//	payload := brokers.MsgPayload{
	//		MsgId:      "124560",
	//		Data:       []byte(`{"Latitude": 409146138, "Longitude": -746188906}`),
	//		DelayType:  0,
	//		DelayValue: 0,
	//	}
	//	consumeRsp := call.ConsumeRsp{}
	//	worker := subscribe.Worker{S: &cfg, GlobalCfg: &subscribe.Config{MaxExecTimeSeconds: 123}}
	//	err := worker.ConsumeMsg(item, &payload, &consumeRsp)
	//	if err != nil {
	//		println(err)
	//	}
	//}
	return &rsp, nil
}

func (s *OmegaServer) Consume(stream pb.OmegaService_ConsumeServer) error {
	// get worker
	req, err := stream.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	if req.Token != "" {
		fmt.Println(req.Token)
		topicIns, err := topic.Mgr.GetTopic(req.Namespace, req.Topic)
		if err != nil {
			return errors.New("not found topic")
		}
		worker, err := topicIns.GetWorker(req.Topic)
		if err != nil {
			return err
		}
		return worker.Run(req, stream)
	}
	return errors.New("token must be not null")
}

func RunServer() error {
	listen, err := net.Listen("tcp", ":12345")
	if err != nil {
		return err
	}
	svr := grpc.NewServer()
	pb.RegisterOmegaServiceServer(svr, &OmegaServer{})
	err = svr.Serve(listen)
	if err != nil {
		return err
	}
	return nil
}
