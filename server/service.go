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
	"github.com/ml444/scheduler/topics"
	"google.golang.org/grpc"
	"io"
	"math/rand"
	"net"
	"time"
)

type OmegaServer struct {
	pb.UnimplementedOmegaServiceServer
}

func (s *OmegaServer) AddTopic(ctx context.Context, req *pb.Topic) (*pb.Response, error) {
	var rsp pb.Response
	topic, err := topics.NewTopic(req.Namespace, req.TopicName, req.Partitions, req.Priority)
	if err != nil {
		return nil, err
	}
	topics.AddTopic(topic)
	return &rsp, nil
}

func generateRandomHashCode() uint64 {
	return rand.Uint64()
}

func (s *OmegaServer) Pub(ctx context.Context, req *pb.PubReq) (*pb.PubRsp, error) {
	var rsp pb.PubRsp
	//fmt.Println(req)

	dataSize := len(req.Data)
	if dataSize > config.MaxDataSize {
		return nil, errors.New("data length exceeds limit")
	}
	partition := req.Partition
	hashCode := req.HashCode
	if hashCode > 0 {
		partitionCount, err := topics.GetTopicPartitionCount(req.Namespace, req.TopicName)
		if err != nil {
			return nil, err
		}
		// TODO
		partition = uint32(hashCode % uint64(partitionCount))
	} else {
		// TODO generate hashCode
		hashCode = generateRandomHashCode()
	}
	item := brokers.Item{
		CreatedAt: uint32(time.Now().Unix()),
		HashCode:  hashCode,
		Partition: partition,
		Size:      uint32(dataSize),
		//DelayType:  req.DelayType,
		//DelayValue: req.DelayValue,
		//Priority:   req.Priority,
		Data:      req.Data,
		//Namespace: req.Namespace,
	}

	err := publish.Pub(ctx, req.Namespace, req.TopicName, partition, &item)
	if err != nil {
		return nil, err
	}
	rsp.Status = 10000
	return &rsp, nil
}

func (s *OmegaServer) Sub(ctx context.Context, req *pb.SubReq) (*pb.SubRsp, error) {
	var rsp pb.SubRsp
	fmt.Println("====> Subscribe", req.ClientId, "===")
	topicIns, err := topics.GetTopic(req.Namespace, req.Topic)
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
		topicIns, err := topics.GetTopic(req.Namespace, req.Topic)
		if err != nil {
			return errors.New("not found topics")
		}
		worker, err := topicIns.GetWorker(req.Token)
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
