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
	stream, err := cli.Consume(ctx)
	if err != nil {
		return err
	}
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

	for {
		if err = stream.Send(&req); err != nil {
			log.Error(err)
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
		err = handler(ctx, rsp)
		if err != nil {
			log.Error(err)
			return err
		}
		req.Sequence++
	}
}
