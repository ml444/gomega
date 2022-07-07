package call

import (
	"time"
)

func Call(req *ConsumeReq, retryCount *uint32) *ConsumeRsp {
	var rsp ConsumeRsp
	beg := time.Now()
	callRsp := Request(req)
	rsp.Took = int64(time.Since(beg) / time.Millisecond)

	if callRsp.retry {
		rsp.Retry = true
		if !callRsp.skipIncRetryCount {
			*retryCount++
		}
	}
	return &rsp
}

func getRoute(key string) (string, error) {
	//dispatch.Route(cfg.ServiceName, rand.Int(), 0, xBrickBeta)
	return "", nil
}

func Request(req *ConsumeReq) *CallRsp {
	var rsp CallRsp

//	addr, err := getRoute()
//	if err != nil {
//		rsp.retry = true
//		rsp.skipIncRetryCount = true
//		rsp.err = err
//		return &rsp
//	}
//	rsp.addr = addr
//
//	cell := &SysCall{
//		CmdPath:    cfg.ServicePath,
//		CmdService: cfg.ServiceName,
//		ModuleName: "smq",
//		CallId:     fmt.Sprintf("%s.%s.%s", p.topicName, p.channelName, payload.MsgId),
//		ErrCode:    0,
//		ErrMsg:     "",
//		CallAt:     0,
//		ReqId:      ctx.GetReqId(),
//		BaseReqId:  getBaseReqId(ctx.GetReqId()),
//		Duration:   0,
//	}
//	defer func() {
//		if cell.Duration > 500 || cell.ErrCode != 0 {
//			sendMsg("sysCall", cell)
//		}
//	}()
//
//	cell.CallAt = time.Now().UnixMilli()
//
//	err = rpc.ClientCallSpecAddressWithTimeout(
//		ctx, cfg.ServiceName, addr, cfg.ServicePath, req, consumeRsp,
//		time.Duration(timeoutSeconds)*time.Second)
//
//	cell.Duration = time.Now().UnixMilli() - cell.CallAt
//
//	if err != nil {
//		rsp.err = err
//
//		var errCode int32
//		if x, ok := err.(*rpc.ErrMsg); ok {
//			errCode = x.ErrCode
//			cell.ErrCode = x.ErrCode
//			cell.ErrMsg = x.ErrMsg
//		} else {
//			errCode = -1
//			cell.ErrCode = -1
//			cell.ErrMsg = err.Error()
//		}
//
//		if errCode < 0 {
//			rsp.retry = true
//			rsp.skipIncRetryCount = false
//		} else {
//			rsp.retry = consumeRsp.Retry
//			rsp.skipIncRetryCount = consumeRsp.SkipIncRetryCount
//		}
//	} else {
//		rsp.retry = consumeRsp.Retry
//		rsp.skipIncRetryCount = consumeRsp.SkipIncRetryCount
//	}
//
//	return &rsp
//
//OUT:
//	rsp.retry = true
//	rsp.skipIncRetryCount = true
//	rsp.err = errors.New("cpu load too high, retry")
//	log.Warnf("cpu load too high smq: %s: call %s %s msgId:%s cml:%d addr:%s, retry", p.logName, cfg.ServiceName, cfg.ServicePath, payload.MsgId, cml, addr)
	return &rsp
}

type CallRsp struct {
	addr  string
	retry bool
	took  int64
	err   error

	skipIncRetryCount bool
}
type SysCall struct {
	CmdPath    string `json:"cmd_path,omitempty"`
	CmdService string `json:"cmd_service,omitempty"`

	ModuleName string `json:"module_name,omitempty"`
	// 任务id、msgid等
	CallId string `json:"call_id,omitempty"`

	ErrCode int32  `json:"err_code,omitempty"`
	ErrMsg  string `json:"err_msg,omitempty"`

	CallAt int64 `json:"call_at,omitempty"`

	ReqId     string `json:"req_id,omitempty"`
	BaseReqId string `json:"base_req_id,omitempty"`

	Duration int64 `json:"duration,omitempty"`
}
