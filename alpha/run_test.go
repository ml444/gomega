package alpha

import (
	"context"
	"fmt"
	log "github.com/ml444/glog"
	"io"
	"testing"
	"time"
)

func Run() {
	// add topic
	_, err := AddTopic(context.TODO(), &Topic{
		Namespace:  "testNamespace",
		TopicName:  "testTopic",
		Partitions: 4,
		Priority:   0,
	})
	if err != nil {
		println(err.Error())
		return
	}

	// subscribe
	subRsp, err := Sub(context.TODO(), &SubReq{
		ClientId:  "test_client_id",
		Namespace: "testNamespace",
		Topic:     "testTopic",
		Group:     "testGroup",
		Policy:    0,
	})
	if err != nil {
		println(err.Error())
		return
	}

	// pub
	//var wg sync.WaitGroup
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	for i := 0; i < 100; i++ {
	//		for j := 0; j < 1000; j++ {
	//			rsp, err := Pub(context.TODO(), &PubReq{
	//				Namespace: "testNamespace",
	//				TopicName: "testTopic",
	//				//Partition:            uint32(j),
	//				HashCode: uint64(j),
	//				Data:     []byte(`{"name": "foo", "age": 22}`),
	//			})
	//			if err != nil {
	//				fmt.Println(err)
	//				return
	//			}
	//			fmt.Println(rsp.Status,i, j)
	//		}
	//		time.Sleep(1 * time.Second)
	//	}
	//}()

	fmt.Println("===> token: ", subRsp.Token)
	// sub
	for {
		err = loopConsume(subRsp.Token, handler)
		if err != nil {
			println(err.Error())
			return
		}
	}
	//wg.Wait()
}

func loopConsume(token string, handler func(r *ConsumeRsp) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	stream, err := client.Consume(ctx)
	if err != nil {
		return err
	}
	var rsp *ConsumeRsp
	req := ConsumeReq{
		Namespace:        "testNamespace",
		Topic:            "testTopic",
		Token:            token,
		Sequence:         1,
		IsRetry:          false,
		IgnoreRetryCount: false,
		RetryIntervalMs:  0,
	}
	for {
		if err = stream.Send(&req); err != nil {
			return err
		}
		rsp, err = stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Error(err)
			return err
		}
		err = handler(rsp)
		if err != nil {
			log.Error(err)
			return err
		}
		req.Sequence++
	}
}

func TestRun(t *testing.T) {
	Run()
}
