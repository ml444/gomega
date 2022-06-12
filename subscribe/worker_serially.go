package subscribe

import (
	"fmt"
	log "github.com/ml444/glog"
	"github.com/ml444/scheduler/mfile"
	"go.uber.org/atomic"
	"sync"
	"time"
)

type SerialWorker struct {
	workerId   int
	msgCh      chan *mfile.Item
	wg         sync.WaitGroup
	backend    *withFileGroup
	item       *mfile.Item
	notifyExit atomic.Bool
	exited     atomic.Bool

	// 任务结果控制
	lastHash       uint32
	isAsync        bool
	maxAsyncWaitMs uint32
	msgId          string
	seq            uint64


}
func (p *SerialWorker) Run() {

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
