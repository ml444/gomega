package alpha

import "context"

func Subscribe(ctx context.Context, req *SubscribeReq) (*SubscribeRsp, error) {
	var err error
	var rsp SubscribeRsp
	Init()
	//return cli.Pub(ctx, req)
	err = conn.Invoke(ctx, "/omega.OmegaService/Subscribe", req, &rsp)
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}
