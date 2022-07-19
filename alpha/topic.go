package alpha

import "context"

func AddTopic(ctx context.Context, req *Topic) (*Response, error) {
	var err error
	var rsp Response
	Init()
	//return cli.Pub(ctx, req)
	err = conn.Invoke(ctx, "/omega.OmegaService/AddTopic", req, &rsp)
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}
