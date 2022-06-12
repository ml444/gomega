package subscribe

import (
	"github.com/ml444/scheduler/backend"
	"github.com/ml444/scheduler/mfile"
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
	index      int
	exitChan   chan int
	msgChan    chan *mfile.Item
	backend    backend.IBackendReader
	retryList  *mfile.MinHeap
	tk         *time.Ticker
	cfg        *Config
	hashMap    map[uint32]mfile.MinHeap
	futureList *mfile.Tree // async
	blockLimit int
}

const defaultTimeout = time.Millisecond * 10

const (
	defaultConsumeMaxExecTimeSeconds = 60
	defaultConsumeMaxRetryCount      = 5
)

func NewConsumeWorker(wg *sync.WaitGroup, index int, msgChan chan *mfile.Item, fg backend.IBackendReader) *consumeWorker {
	return &consumeWorker{
		wg:        wg,
		index:     index,
		exitChan:  make(chan int, 1),
		msgChan:   msgChan,
		backend:   fg,
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
		return p.retryList.PopEl().Value.(*retryItem).item
	}
	return nil
}

func (p *consumeWorker) RunConcurrently() {
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
					Value: &retryItem{
						item:       msg,
						nextExecAt: execAt,
					},
					Priority: execAt,
				})

			}
		} else {
			fg.qr.SetFinish(msg)
		}
	}
}
