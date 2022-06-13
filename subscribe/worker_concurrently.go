package subscribe

import (
	"fmt"
	log "github.com/ml444/glog"
	"github.com/ml444/scheduler/backend"
	"github.com/ml444/scheduler/mfile"
	"strings"
	"sync"
	"time"
)

const (
	WorkerTypeConcurrent = 1
	WorkerTypeSerial     = 2
	WorkerTypeAsync      = 3
)

type consumeWorker struct {
	typ        int8
	wg         *sync.WaitGroup
	idx        int
	exitChan   chan int
	msgChan    chan *mfile.Item
	backend    backend.IBackendReader
	retryList  *mfile.MinHeap
	tk         *time.Ticker
	cfg        *Config
	futureList *mfile.Tree // async
	blockLimit int
}

const defaultTimeout = time.Millisecond * 10

const (
	defaultConsumeMaxExecTimeSeconds = 60
	defaultConsumeMaxRetryCount      = 5
)

func NewConsumeWorker(wg *sync.WaitGroup, idx int, msgChan chan *mfile.Item, backend backend.IBackendReader) *consumeWorker {
	return &consumeWorker{
		wg:        wg,
		idx:       idx,
		exitChan:  make(chan int, 1),
		msgChan:   msgChan,
		backend:   backend,
		retryList: mfile.NewMinHeap(),
		tk:        time.NewTicker(defaultTimeout),
	}
}
func (p *consumeWorker) notifyExit() {
	select {
	case p.exitChan <- 1:
	default:
	}
}
func (p *consumeWorker) getNextMsg(exit *bool) *mfile.Item {

	select {
	case msg := <-p.msgChan:
		return msg
	case <-p.exitChan:
		*exit = true
		return nil
	case <-p.tk.C:
		return nil
	}
}

func (p *consumeWorker) tryRetryMsg() *mfile.Item {
	top := p.retryList.PeekEl()
	now := time.Now().UnixMilli()
	if top.Priority <= now {
		return p.retryList.PopEl().Value.(*mfile.Item)
	}
	return nil
}

func (p *consumeWorker) Run() {
	defer p.wg.Done()
	fg := p.backend
	for {
		var exit bool
		var msg *mfile.Item
		var blockCount int
		blockCount = p.retryList.Len()
		if blockCount > 0  {
			msg = p.tryRetryMsg()
		}
		if msg == nil {
			msg = p.getNextMsg(&exit)
		}
		if exit {
			break
		}
		if msg == nil {
			continue
		}
		// TODO Skip blackList msg
		payload, err := decodePayload(msg.Data)
		if err != nil {
			// TODO report err
			fg.qr.SetFinish(msg)
			p.onFinalFail(msg, payload)
			continue
		}
		var consumeRsp ConsumeRsp
		p.consumeMsg(msg, payload, &consumeRsp)
		if consumeRsp.Retry {
			maxRetryCount := p.getMaxRetryCount()
			if msg.RetryCount >= maxRetryCount {
				fg.qr.SetFinish(msg)
				p.onFinalFail(msg, payload)
			} else {
				var waitMs int64
				if consumeRsp.RetryIntervalSeconds > 0 {
					waitMs = consumeRsp.RetryIntervalSeconds * 1000
				} else if consumeRsp.RetryIntervalMs > 0 {
					waitMs = consumeRsp.RetryIntervalMs
				} else {
					waitMs = int64(p.getNextRetryWait(msg.RetryCount)) * 1000
				}
				execAt := time.Now().UnixMilli() + waitMs
				p.retryList.PushEl(&mfile.MinHeapElement{
					Value: msg,
					Priority: execAt,
				})
			}
		} else {
			fg.qr.SetFinish(msg)
		}
	}
}
//type retryItem struct {
//	item       *mfile.Item
//	nextExecAt int64
//}

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
