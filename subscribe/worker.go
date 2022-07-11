package subscribe

import (
	"context"
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

type Worker struct {
	typ        int8
	wg         *sync.WaitGroup
	idx        int
	exitChan   chan int
	msgChan    chan *brokers.Item
	finishChan chan *brokers.Item
	retryList  *brokers.MinHeap
	tk         *time.Ticker
	Cfg        *Config
	//futureList *structure.Tree // async
	blockLimit int

	S *Subscriber
}

const defaultTimeout = time.Millisecond * 10

const (
	defaultConsumeMaxExecTimeSeconds = 60
	defaultConsumeMaxRetryCount      = 5
)

func NewConsumeWorker(wg *sync.WaitGroup, idx int, msgChan chan *brokers.Item, finishChan chan *brokers.Item) *Worker {
	return &Worker{
		wg:         wg,
		idx:        idx,
		exitChan:   make(chan int, 1),
		msgChan:    msgChan,
		finishChan: finishChan,
		retryList:  brokers.NewMinHeap(),
		tk:         time.NewTicker(defaultTimeout),
	}
}
func (w *Worker) notifyExit() {
	select {
	case w.exitChan <- 1:
	default:
	}
}
func (w *Worker) NextMsg(exit *bool) *brokers.Item {

	select {
	case msg := <-w.msgChan:
		return msg
	case <-w.exitChan:
		*exit = true
		return nil
	case <-w.tk.C:
		//TODO no tk
		return nil
	}
}

func (w *Worker) tryRetryMsg() *brokers.Item {
	top := w.retryList.PeekEl()
	now := time.Now().UnixMilli()
	if top.Priority <= now {
		return w.retryList.PopEl().Value.(*brokers.Item)
	}
	return nil
}

func (w *Worker) setFinish(msg *brokers.Item) {
	w.finishChan <- msg
}
func (w *Worker) Run() {
	defer w.wg.Done()

	cfg := w.Cfg
	if cfg == nil {
		panic("Cfg is nil")
	}

	if cfg.ServicePath == "" || cfg.ServiceName == "" {
		panic("Cfg hasn't config ServicePath or ServiceName")
	}
	for {
		var exit bool
		var msg *brokers.Item
		var blockCount int
		blockCount = w.retryList.Len()
		if blockCount > 0 {
			msg = w.tryRetryMsg()
		}
		if msg == nil {
			msg = w.NextMsg(&exit)
		}
		if exit {
			break
		}
		if msg == nil {
			continue
		}
		// TODO Skip blackList msg
		payload, err := brokers.DecodeMsgPayload(msg.Data)
		if err != nil {
			// TODO report err
			w.setFinish(msg)
			continue
		}
		var consumeRsp call.ConsumeRsp
		_ = w.ConsumeMsg(msg, payload, &consumeRsp)
		if consumeRsp.Retry {
			maxRetryCount := w.getMaxRetryCount()
			if msg.RetryCount >= maxRetryCount {
				w.setFinish(msg)
			} else {
				var waitMs int64
				if consumeRsp.RetryIntervalMs > 0 {
					waitMs = consumeRsp.RetryIntervalMs
				} else {
					waitMs = w.getNextRetryWait(msg.RetryCount) * 1000
				}
				execAt := time.Now().UnixMilli() + waitMs
				w.retryList.PushEl(&brokers.MinHeapElement{
					Value:    msg,
					Priority: execAt,
				})
			}
		} else {
			w.setFinish(msg)
		}
	}
}

//type retryItem struct {
//	item       *brokers.Item
//	nextExecAt int64
//}

func (w *Worker) ConsumeMsg(item *brokers.Item, payload *brokers.MsgPayload, consumeRsp *call.ConsumeRsp) error {
	// TODO getRoute(checkRoute())
	ctx := context.TODO()
	s := w.S
	var timeoutSeconds = w.Cfg.MaxExecTimeSeconds
	if timeoutSeconds == 0 {
		timeoutSeconds = defaultConsumeMaxExecTimeSeconds
	}
	meta := &call.MsgMeta{
		CreatedAt:   item.CreatedAt,
		RetryCnt:    item.RetryCount,
		Data:        payload.Data,
		MsgId:       payload.MsgId,
	}
	if s.BeforeProcess != nil {
		s.BeforeProcess(ctx, meta)
	}
	in, err := s.UnMarshalRequest(payload.Data)
	if err != nil {
		log.Error(err)
		return err
	}
	out := s.NewResponse()
	err = call.Call(ctx, w.S.Route, &in, &out, timeoutSeconds)
	if err != nil {
		log.Error(err)
		consumeRsp.Retry = true
		item.RetryCount++
		return err
	}
	if s.AfterProcess != nil {
		isRetry, isIgnoreRetryCount := s.AfterProcess(ctx, meta, &in, &out)
		if isRetry {
			if consumeRsp != nil {
				consumeRsp.Retry = true
			}
		}
		if !isIgnoreRetryCount {
			item.RetryCount++
		}
	}
	return nil
}

func (w *Worker) getMaxRetryCount() uint32 {
	s := w.Cfg
	maxRetryCount := uint32(defaultConsumeMaxRetryCount)
	if s != nil && s.MaxRetryCount > 0 {
		maxRetryCount = s.MaxRetryCount
	}
	return maxRetryCount
}

func (w *Worker) getNextRetryWait(retryCnt uint32) int64 {
	s := w.Cfg
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



