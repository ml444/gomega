package scheduler

import (
	"fmt"
	"sync"
)

type BarrierWorker struct {
	workerId   int
	reqCh      chan *BarrierWorker
	wg         sync.WaitGroup
	fg         *withFileGroup
	item       *mfile.Item
	notifyExit atomic.Bool
	exited     atomic.Bool

	// 任务结果控制
	lastHash       uint32
	isAsync        bool
	maxAsyncWaitMs uint32
	msgId          string
	seq            uint64

	// 重试
	retry        bool
	retryItem    *mfile.Item
	retryDelayMs uint32
}
func NewBarrierWorker(fg *withFileGroup, reqCh chan *BarrierWorker, workerId int) *BarrierWorker {
	return &BarrierWorker{
		fg:       fg,
		reqCh:    reqCh,
		workerId: workerId,
	}
}
func (p *BarrierWorker) start() {
	p.exited.Store(false)
	go func() {
		defer func() {
			log.Warnf("worker %d exit", p.workerId)
			p.exited.Store(true)
		}()
		for !p.notifyExit.Load() {
			p.wg.Add(1)
			p.reqCh <- p
			p.wg.Wait()
			if p.item == nil {
				log.Infof("wait item return nil, exit")
				break
			}
			msg := p.item
			fg := p.fg
			payload, err := decodePayload(msg.Data)
			if err != nil {
				log.Errorf("%s: decode payload err:%v", fg.logName, err)
				warning.ReportMsg(nil, fmt.Sprintf("smq: %s: decode payload err:%v", fg.logName, err))
				fg.qr.SetFinish(msg)
				fg.onFinalFail(msg, payload)
				continue
			}
			var consumeRsp smq.ConsumeRsp
			fg.consumeMsg(msg, payload, &consumeRsp)
			if consumeRsp.Retry {
				maxRetryCount := fg.getMaxRetryCount()
				if msg.RetryCount >= maxRetryCount {
					fg.qr.SetFinish(msg)
					fg.onFinalFail(msg, payload)
				} else {
					var waitMs uint32
					if consumeRsp.RetryIntervalSeconds > 0 {
						waitMs = consumeRsp.RetryIntervalSeconds * 1000
					} else if consumeRsp.RetryIntervalMs > 0 {
						waitMs = consumeRsp.RetryIntervalMs
					} else {
						waitMs = fg.getNextRetryWait(msg.RetryCount) * 1000
					}
					log.Warnf("%s: wait %d ms to retry %s cpu:%f memory:%f", fg.name, waitMs, payload.MsgId, consumeRsp.CpuPercent, consumeRsp.MemoryPercent)
					p.retry = true
					p.retryDelayMs = waitMs
					p.retryItem = msg
				}
			} else {
				// async 仅 retry == false 有效
				p.isAsync = consumeRsp.IsAsync
				p.maxAsyncWaitMs = consumeRsp.MaxAsyncWaitMs
				p.msgId = payload.MsgId
			}
		}
	}()
}
