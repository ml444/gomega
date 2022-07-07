package subscribe

import (
	log "github.com/ml444/glog"
	"github.com/ml444/scheduler/brokers"
	"github.com/ml444/scheduler/subscribe/call"
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
	msgChan    chan *brokers.Item
	finishChan chan *brokers.Item
	retryList  *brokers.MinHeap
	tk         *time.Ticker
	cfg        *Config
	//futureList *structure.Tree // async
	blockLimit int
}

const defaultTimeout = time.Millisecond * 10

const (
	defaultConsumeMaxExecTimeSeconds = 60
	defaultConsumeMaxRetryCount      = 5
)

func NewConsumeWorker(wg *sync.WaitGroup, idx int, msgChan chan *brokers.Item, finishChan chan *brokers.Item) *consumeWorker {
	return &consumeWorker{
		wg:         wg,
		idx:        idx,
		exitChan:   make(chan int, 1),
		msgChan:    msgChan,
		finishChan: finishChan,
		retryList:  brokers.NewMinHeap(),
		tk:         time.NewTicker(defaultTimeout),
	}
}
func (p *consumeWorker) notifyExit() {
	select {
	case p.exitChan <- 1:
	default:
	}
}
func (p *consumeWorker) NextMsg(exit *bool) *brokers.Item {

	select {
	case msg := <-p.msgChan:
		return msg
	case <-p.exitChan:
		*exit = true
		return nil
	case <-p.tk.C:
		//TODO no tk
		return nil
	}
}

func (p *consumeWorker) tryRetryMsg() *brokers.Item {
	top := p.retryList.PeekEl()
	now := time.Now().UnixMilli()
	if top.Priority <= now {
		return p.retryList.PopEl().Value.(*brokers.Item)
	}
	return nil
}

func (p *consumeWorker) setFinish(msg *brokers.Item) {
	p.finishChan <- msg
}
func (p *consumeWorker) Run() {
	defer p.wg.Done()

	cfg := p.cfg
	if cfg == nil {
		panic("cfg is nil")
	}

	if cfg.ServicePath == "" || cfg.ServiceName == "" {
		panic("cfg hasn't config ServicePath or ServiceName")
	}
	for {
		var exit bool
		var msg *brokers.Item
		var blockCount int
		blockCount = p.retryList.Len()
		if blockCount > 0 {
			msg = p.tryRetryMsg()
		}
		if msg == nil {
			msg = p.NextMsg(&exit)
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
			p.setFinish(msg)
			p.onFinalFail(msg, payload)
			continue
		}
		var consumeRsp call.ConsumeRsp
		p.consumeMsg(msg, payload, &consumeRsp)
		if consumeRsp.Retry {
			maxRetryCount := p.getMaxRetryCount()
			if msg.RetryCount >= maxRetryCount {
				p.setFinish(msg)
				p.onFinalFail(msg, payload)
			} else {
				var waitMs int64
				if consumeRsp.RetryIntervalMs > 0 {
					waitMs = consumeRsp.RetryIntervalMs
				} else {
					waitMs = p.getNextRetryWait(msg.RetryCount) * 1000
				}
				execAt := time.Now().UnixMilli() + waitMs
				p.retryList.PushEl(&brokers.MinHeapElement{
					Value:    msg,
					Priority: execAt,
				})
			}
		} else {
			p.setFinish(msg)
		}
	}
}

//type retryItem struct {
//	item       *brokers.Item
//	nextExecAt int64
//}

func (p *consumeWorker) consumeMsg(item *brokers.Item, payload *MsgPayload, consumeRsp *call.ConsumeRsp) {
	// TODO getRoute(checkRoute())

	var timeoutSeconds = p.cfg.MaxExecTimeSeconds
	if timeoutSeconds == 0 {
		timeoutSeconds = defaultConsumeMaxExecTimeSeconds
	}
	req := &call.ConsumeReq{
		CreatedAt:   item.CreatedAt,
		RetryCnt:    item.RetryCount,
		Data:        payload.Data,
		MsgId:       payload.MsgId,
		MaxRetryCnt: p.getMaxRetryCount(),
		Timeout:     timeoutSeconds,
	}
	consumeRsp = call.Call(req, &item.RetryCount)
}

func (p *consumeWorker) getMaxRetryCount() uint32 {
	s := p.cfg
	maxRetryCount := uint32(defaultConsumeMaxRetryCount)
	if s != nil && s.MaxRetryCount > 0 {
		maxRetryCount = s.MaxRetryCount
	}
	return maxRetryCount
}

func (p *consumeWorker) getNextRetryWait(retryCnt uint32) int64 {
	s := p.cfg
	retryMs := s.RetryIntervalMs
	if s != nil && retryMs > 0 {
		if s.RetryIntervalStep > 0 {
			return retryMs * s.RetryIntervalStep
		}
		return retryMs
	}
	if retryCnt == 0 {
		retryCnt = 1
	}
	return int64(retryCnt) * 5
}

func (p *consumeWorker) onFinalFail(item *brokers.Item, payload *MsgPayload) {
	log.Warnf("%s: msg %s touch max retry count, drop",
		p, payload.MsgId)

}

func decodePayload(buf []byte) (*MsgPayload, error) {
	var p MsgPayload
	//err := proto.Unmarshal(buf, &p)
	//if err != nil {
	//	log.Errorf("err:%v", err)
	//	return nil, err
	//}
	return &p, nil
}

type MsgPayload struct {
	//Ctx                  *MsgPayload_Context `protobuf:"bytes,1,opt,name=ctx" json:"ctx,omitempty"`
	MsgId string `protobuf:"bytes,2,opt,name=msg_id,json=msgId" json:"msg_id,omitempty"`
	Data  []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}
