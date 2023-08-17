package alpha

import (
	"context"
	"io"

	"github.com/ml444/gomega/alpha/omega"
	"github.com/ml444/gomega/log"
)

type consumeFunc func(ctx context.Context, msg *omega.ConsumeRsp) error

func Sub(ctx context.Context, req *omega.SubReq) (*omega.SubRsp, error) {
	//var err error
	//var rsp omega.SubRsp
	return cli.Sub(ctx, req)
}

func LoopConsume(ctx context.Context, topic, group, clientId string, handler consumeFunc) error {
	subRsp, err := Sub(ctx, &omega.SubReq{
		ClientId: clientId,
		Topic:    topic,
		Group:    group,
	})
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	var rsp *omega.ConsumeRsp
	req := omega.ConsumeReq{
		Token:    subRsp.Token,
		Sequence: subRsp.Sequence,
	}

	stream, err := cli.Consume(ctx, &req)
	if err != nil {
		log.Error(err)
		return err
	}
	for {
		rsp, err = stream.Recv()
		if err == io.EOF {
			log.Error(err)
			break
		}
		if err != nil {
			log.Error(err)
			return err
		}
		err = handler(ctx, rsp)
		if err != nil {
			log.Error(err)
			return err
		}
		req.Sequence++
	}
	log.Warnf("%s LoopConsume exit", clientId)
	return nil
}
