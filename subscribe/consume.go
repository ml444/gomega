package subscribe

import (
	"errors"
	"fmt"
	log "github.com/ml444/glog"
	"github.com/ml444/scheduler/backend"
	"github.com/ml444/scheduler/mfile"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type ConsumeReq struct {
	CreatedAt uint32 `protobuf:"varint,1,opt,name=created_at,json=createdAt" json:"created_at,omitempty"`
	RetryCnt  uint32 `protobuf:"varint,2,opt,name=retry_cnt,json=retryCnt" json:"retry_cnt,omitempty"`
	Data      []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	CorpId    uint32 `protobuf:"varint,4,opt,name=corp_id,json=corpId" json:"corp_id,omitempty"`
	AppId     uint32 `protobuf:"varint,5,opt,name=app_id,json=appId" json:"app_id,omitempty"`
	MsgId     string `protobuf:"bytes,6,opt,name=msg_id,json=msgId" json:"msg_id,omitempty"`
	// @desc: 服务端最大重试次数，业务可以根据这个来做最后失败逻辑
	MaxRetryCnt          uint32 `protobuf:"varint,7,opt,name=max_retry_cnt,json=maxRetryCnt" json:"max_retry_cnt,omitempty"`
	Timeout              uint32
	XXX_NoUnkeyedLiteral struct{} `json:"-" bson:"-"`
	XXX_unrecognized     []byte   `json:"-" bson:"-"`
	XXX_sizecache        int32    `json:"-" bson:"-"`
}
type ConsumeRsp struct {
	// @desc: 是否需要重试
	Retry bool `protobuf:"varint,1,opt,name=retry" json:"retry,omitempty"`
	// @desc: 重试间隔时间，可以这里指定，不指定的话，由 smq 来定 (单位:秒)
	RetryIntervalSeconds int64 `protobuf:"varint,2,opt,name=retry_interval_seconds,json=retryIntervalSeconds" json:"retry_interval_seconds,omitempty"`
	// @desc: 为 true 时，不要累加重试次数
	SkipIncRetryCount bool `protobuf:"varint,3,opt,name=skip_inc_retry_count,json=skipIncRetryCount" json:"skip_inc_retry_count,omitempty"`
	// @desc: 支持毫秒级重试时间
	RetryIntervalMs int64 `protobuf:"varint,4,opt,name=retry_interval_ms,json=retryIntervalMs" json:"retry_interval_ms,omitempty"`
	// @desc: 是否异步消费
	IsAsync bool `protobuf:"varint,5,opt,name=is_async,json=isAsync" json:"is_async,omitempty"`
	// @desc: 异步消费最大等待确认时间
	MaxAsyncWaitMs int64 `protobuf:"varint,6,opt,name=max_async_wait_ms,json=maxAsyncWaitMs" json:"max_async_wait_ms,omitempty"`
	// @desc: 消费端服务cpu使用率
	CpuPercent float32 `protobuf:"fixed32,7,opt,name=cpu_percent,json=cpuPercent" json:"cpu_percent,omitempty"`
	// @desc: 消费端服务内存使用率
	MemoryPercent        float32 `protobuf:"fixed32,8,opt,name=memory_percent,json=memoryPercent" json:"memory_percent,omitempty"`
	Took                 int64
	XXX_NoUnkeyedLiteral struct{} `json:"-" bson:"-"`
	XXX_unrecognized     []byte   `json:"-" bson:"-"`
	XXX_sizecache        int32    `json:"-" bson:"-"`
}
type retryItem struct {
	item       *mfile.Item
	nextExecAt int64
}

func (p *consumeWorker) RunAsync() {

}

func (p *consumeWorker) consumeMsg(item *mfile.Item, payload *MsgPayload, consumeRsp *ConsumeRsp) {
	cfg := p.cfg
	if cfg == nil {
		log.Warnf("sub cfg is nil, skip")
		consumeRsp.Retry = true
		return
	}

	if cfg.ServicePath == "" || cfg.ServiceName == "" {
		consumeRsp.Retry = true
		return
	}

	// TODO getRoute(checkRoute())


	var timeoutSeconds = cfg.MaxExecTimeSeconds
	if timeoutSeconds == 0 {
		timeoutSeconds = defaultConsumeMaxExecTimeSeconds
	}
	req := &ConsumeReq{
		CreatedAt:   item.CreatedAt,
		RetryCnt:    item.RetryCount,
		Data:        payload.Data,
		MsgId:       payload.MsgId,
		MaxRetryCnt: p.getMaxRetryCount(),
		Timeout:     timeoutSeconds,
	}

	beg := time.Now()
	callRsp := Call(req)
	consumeRsp.Took = int64(time.Since(beg) / time.Millisecond)

	if callRsp.retry {
		consumeRsp.Retry = true
		if !callRsp.skipIncRetryCount {
			item.RetryCount++
		}
	}
}

func getRoute(key string) (string, error) {
	return "", nil
}

func Call(req *ConsumeReq) *CallRsp {
	var rsp CallRsp


	addr, err := dispatch.Route(cfg.ServiceName, rand.Int(), 0, xBrickBeta)
	if err != nil {
		rsp.retry = true
		rsp.skipIncRetryCount = true
		rsp.err = err
		return &rsp
	}
	rsp.addr = addr

	cell := &SysCall{
		CmdPath:    cfg.ServicePath,
		CmdService: cfg.ServiceName,
		ModuleName: "smq",
		CallId:     fmt.Sprintf("%s.%s.%s", p.topicName, p.channelName, payload.MsgId),
		ErrCode:    0,
		ErrMsg:     "",
		CallAt:     0,
		ReqId:      ctx.GetReqId(),
		BaseReqId:  getBaseReqId(ctx.GetReqId()),
		Duration:   0,
	}
	defer func() {
		if cell.Duration > 500 || cell.ErrCode != 0 {
			sendMsg("sysCall", cell)
		}
	}()


	cell.CallAt = time.Now().UnixMilli()

	err = rpc.ClientCallSpecAddressWithTimeout(
		ctx, cfg.ServiceName, addr, cfg.ServicePath, req, consumeRsp,
		time.Duration(timeoutSeconds)*time.Second)

	cell.Duration = time.Now().UnixMilli() - cell.CallAt

	if err != nil {
		rsp.err = err

		var errCode int32
		if x, ok := err.(*rpc.ErrMsg); ok {
			errCode = x.ErrCode
			cell.ErrCode = x.ErrCode
			cell.ErrMsg = x.ErrMsg
		} else {
			errCode = -1
			cell.ErrCode = -1
			cell.ErrMsg = err.Error()
		}

		if errCode < 0 {
			rsp.retry = true
			rsp.skipIncRetryCount = false
		} else {
			rsp.retry = consumeRsp.Retry
			rsp.skipIncRetryCount = consumeRsp.SkipIncRetryCount
		}
	} else {
		rsp.retry = consumeRsp.Retry
		rsp.skipIncRetryCount = consumeRsp.SkipIncRetryCount
	}

	return &rsp

OUT:
	rsp.retry = true
	rsp.skipIncRetryCount = true
	rsp.err = errors.New("cpu load too high, retry")
	log.Warnf("cpu load too high smq: %s: call %s %s msgId:%s cml:%d addr:%s, retry", p.logName, cfg.ServiceName, cfg.ServicePath, payload.MsgId, cml, addr)
	return &rsp
}

type CallRsp struct {
	addr  string
	retry bool
	took  int64
	err   error

	skipIncRetryCount bool
}

func (p *consumeWorker) getMaxRetryCount() uint32 {
	s := p.cfg
	maxRetryCount := uint32(defaultConsumeMaxRetryCount)
	if s != nil && s.MaxRetryCount > 0 {
		maxRetryCount = s.MaxRetryCount
	}
	return maxRetryCount
}

func (p *consumeWorker) getNextRetryWait(retryCnt uint32) uint32 {
	s := p.cfg
	if s != nil && s.RetryIntervalMax > 0 {
		min := s.RetryIntervalMin
		max := s.RetryIntervalMax
		if min >= max {
			return max
		}
		diff := max - min
		cnt := p.getMaxRetryCount()
		var step uint32
		if s.RetryIntervalStep > 0 {
			step = s.RetryIntervalStep
		} else {
			step = diff / cnt
			if step == 0 {
				step = 1
			}
		}
		if retryCnt > 0 {
			retryCnt--
		}
		res := min + (retryCnt * step)
		if res >= max {
			res = max
		}
		return res
	}
	if retryCnt == 0 {
		retryCnt = 1
	}
	return retryCnt * 5
}

func genNewLogCtx(oldCtx string, id string) string {
	if oldCtx == "" {
		return id
	}
	pos := strings.IndexByte(oldCtx, '.')
	var s string
	if pos > 0 {
		s = oldCtx[:pos]
	} else {
		s = oldCtx
	}
	return fmt.Sprintf("%s.%s", s, id)
}

func (p *consumeWorker) onFinalFail(item *mfile.Item, payload *MsgPayload) {
	log.Warnf("%s: msg %s touch max retry count, drop",
		p.logName, payload.MsgId)
	warning.ReportMsg(nil, "smq: %s: final fail", p.logName)
	obj := &smq.SosObjectPayload{
		CreatedAt:   item.CreatedAt,
		Hash:        item.Hash,
		Data:        payload.Data,
		TopicName:   p.topicName,
		ChannelName: p.channelName,
	}
	pb, err := proto.Marshal(obj)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	id, err := sos.SetObject(nil, &sos.SetReq{
		PersistentType: finalFailObjectPersistentType,
		Buf:            pb,
		CorpId:         item.CorpId,
		AppId:          item.AppId,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		warning.ReportMsg(nil, "sos.SetObject err %v", err)
	}

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

func getBaseReqId(str string) string {
	idx := strings.Index(str, ".")
	if idx > 0 {
		return utils.ShortStr(str, idx)
	}
	return str
}

func decodePayload(buf []byte) (*MsgPayload, error) {
	var p smq.MsgPayload
	err := proto.Unmarshal(buf, &p)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &p, nil
}

type MsgPayload struct {
	//Ctx                  *MsgPayload_Context `protobuf:"bytes,1,opt,name=ctx" json:"ctx,omitempty"`
	MsgId                string   `protobuf:"bytes,2,opt,name=msg_id,json=msgId" json:"msg_id,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-" bson:"-"`
	XXX_unrecognized     []byte   `json:"-" bson:"-"`
	XXX_sizecache        int32    `json:"-" bson:"-"`
}
