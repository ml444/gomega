package scheduler

import (
	"fmt"
	log "github.com/ml444/glog"
	"math/rand"
	"time"

	"go.uber.org/atomic"
)

type RetryItem struct {
	item    *mfile.Item
	delayMs uint32
}
type BarrierConsumer struct {
	fg           *withFileGroup
	concurrent   int
	exitChan     chan bool
	exited       atomic.Bool
	reqCh        chan *BarrierWorker
	pendingItems []*mfile.Item
	workers      []*BarrierWorker
}

func NewBarrierConsumer(fg *withFileGroup, concurrent int, old *BarrierConsumer) *BarrierConsumer {
	if concurrent <= 0 || concurrent > 10*10000 {
		panic(fmt.Sprintf("invalid concurrent %d", concurrent))
	}
	p := &BarrierConsumer{
		fg:         fg,
		concurrent: concurrent,
		reqCh:      make(chan *BarrierWorker, concurrent),
		exitChan:   make(chan bool, 1),
	}
	if old != nil {
		for _, v := range old.workers {
			if v.item != nil {
				p.pendingItems = append(p.pendingItems, v.item)
			}
		}
	}
	return p
}

func (p *BarrierConsumer) notifyExit() {
	for {
		select {
		case w := <-p.reqCh:
			log.Infof("abort req for worker %d", w.workerId)
			w.item = nil
			w.wg.Done()
		default:
			goto OUT
		}
	}
OUT:
	log.Infof("notify consumer")
	select {
	case p.exitChan <- true:
	default:
	}
}

type barrierConsumerStat struct {
	name string

	statAt uint32

	workerReqCnt           int
	workerReqGotItemCnt    int
	workerReqMissedItemCnt int

	assignWorkerCallCnt   int
	assignWorkerAssignCnt int
}

func (p *barrierConsumerStat) dumpAndReset(dt uint32) {
	log.Infof("%s: dt %d workerReqCnt %d workerReqGotItemCnt %d workerReqMissedItemCnt %d "+
		"assignWorkerCallCnt %d assignWorkerAssignCnt %d",
		p.name, dt, p.workerReqCnt, p.workerReqGotItemCnt, p.workerReqMissedItemCnt,
		p.assignWorkerCallCnt, p.assignWorkerAssignCnt,
	)
	p.workerReqCnt = 0
	p.workerReqGotItemCnt = 0
	p.workerReqMissedItemCnt = 0

	p.assignWorkerCallCnt = 0
	p.assignWorkerAssignCnt = 0
}

type asyncMsg struct {
	el *utils.MinHeapElement

	seq   uint64
	hash  uint32
	msgId string
}

type asyncMsgMgr struct {
	mh *utils.MinHeap
	h  map[string]*asyncMsg
}

func newAsyncMsgMgr() *asyncMsgMgr {
	return &asyncMsgMgr{
		mh: utils.NewMinHeap(),
		h:  map[string]*asyncMsg{},
	}
}

func (p *asyncMsgMgr) add(msgId string, seq uint64, hash uint32, maxAsyncWaitMs uint32) (existed bool) {
	if p.h[msgId] != nil {
		existed = true
		return
	}
	el := &utils.MinHeapElement{
		Value:    msgId,
		Priority: int64(utils.NowMs()) + int64(maxAsyncWaitMs),
	}
	p.mh.PushEl(el)
	p.h[msgId] = &asyncMsg{
		el:    el,
		hash:  hash,
		msgId: msgId,
		seq:   seq,
	}
	existed = false
	return
}

func (p *asyncMsgMgr) del(msgId string) *asyncMsg {
	m := p.h[msgId]
	if m == nil {
		return nil
	}
	p.mh.RemoveEl(m.el)
	delete(p.h, msgId)
	return m
}

func (p *asyncMsgMgr) clearTimeout(cb func(m *asyncMsg)) {
	now := int64(utils.NowMs())
	for p.mh.Len() > 0 {
		f := p.mh.PeekEl()
		if f.Priority > now {
			break
		}

		_ = p.mh.PopEl()
		msgId := f.Value.(string)
		m := p.h[msgId]
		if m != nil {
			delete(p.h, msgId)
			if cb != nil {
				cb(m)
			}
		}
	}
}

func asyncMsgMgrTest() {
	//mgr := newAsyncMsgMgr()
	//mgr.add("msgId", 1, 10)
	//mgr.add("msgId2", 2, 10)
	//mgr.add("msgId3", 3, 2000)
	//v := mgr.del("msgId")
	//log.Infof("ret %+v", v)
	//
	//time.Sleep(time.Second)
	//mgr.clearTimeout(func(m *asyncMsg) {
	//	log.Infof("timeout %+v", m)
	//})
	//
	//v = mgr.del("msgId3")
	//log.Infof("ret 3 %+v", v)
	//
	//v = mgr.del("msgId4")
	//log.Infof("ret 4 %+v", v)
	//
	//mgr.clearTimeout(func(m *asyncMsg) {
	//	log.Infof("timeout 2 %+v", m)
	//})
}

func (p *BarrierConsumer) start() {
	p.exited.Store(false)
	for i := 0; i < p.concurrent; i++ {
		w := NewBarrierWorker(p.fg, p.reqCh, i)
		p.workers = append(p.workers, w)
		w.start()
	}
	log.Infof("%s: start %d workers", p.fg.logName, p.concurrent)

	go func() {
		var st barrierConsumerStat
		st.statAt = utils.Now()
		st.name = p.fg.logName

		am := newAsyncMsgMgr()

		var pendingWorkers []*BarrierWorker
		defer func() {
			log.Warnf("consumer exit")
			for _, v := range pendingWorkers {
				log.Warnf("abort worker %d", v.workerId)
				v.item = nil
				v.wg.Done()
			}
			p.exited.Store(true)
		}()
		q := p.fg.qr.(*mfile.BarrierQueueGroupReader)
		barrierCount := 0
		if p.fg.subCfg != nil {
			barrierCount = int(p.fg.subCfg.BarrierCount)
		}
		moreCh := q.GetNotifyWriteChan4Consumer()
		asyncConfirmCh := q.GetAsyncMsgConfirmChan()
		var skipHashCounter = map[uint32]int{}
		dispatchItem := func(worker *BarrierWorker, item *mfile.Item) {
			worker.item = item
			worker.isAsync = false
			worker.maxAsyncWaitMs = 0
			worker.msgId = ""
			worker.lastHash = item.Hash
			worker.seq = item.Seq
			worker.retry = false
			worker.retryItem = nil
			worker.retryDelayMs = 0
			c := skipHashCounter[item.Hash]
			c++
			skipHashCounter[item.Hash] = c
			worker.wg.Done()
		}
		popSkipHash := func() (*mfile.Item, error) {
			if len(p.pendingItems) == 0 {
				i, err := q.PopSkipHash(skipHashCounter, barrierCount)
				if err != nil || i == nil {
					return i, err
				}
				if i.Hash == 0 {
					i.Hash = i.CorpId
				}
				if i.Hash == 0 {
					i.Hash = uint32(rand.Int() % p.concurrent)
					log.Warnf("used rand hash %d", i.Hash)
				}
				return i, nil
			}
			i := p.pendingItems[0]
			p.pendingItems = p.pendingItems[1:]
			return i, nil
		}
		assignPendingWorker := func() {
			st.assignWorkerCallCnt++

			for len(pendingWorkers) > 0 {
				item, err := popSkipHash()
				if err != nil {
					log.Errorf("err:%v", err)
					return
				}
				if item == nil {
					break
				}
				worker := pendingWorkers[0]
				pendingWorkers = pendingWorkers[1:]
				dispatchItem(worker, item)

				st.assignWorkerAssignCnt++
			}
		}

		confirmHash := func(hash uint32) (cleared bool) {
			cleared = true
			c := skipHashCounter[hash]
			if c > 0 {
				c--
			}
			if c == 0 {
				delete(skipHashCounter, hash)
			} else {
				skipHashCounter[hash] = c
			}
			return
		}

		lastRepAt := uint32(0)

		tk := time.NewTicker(time.Second)
		defer tk.Stop()

		log.Infof("%s: into main loop", p.fg.logName)
		for {
			if p.fg.subCfg != nil {
				barrierCount = int(p.fg.subCfg.BarrierCount)
			}
			select {
			case worker := <-p.reqCh:
				hash := worker.lastHash

				if worker.retry {
					if worker.retryItem == nil {
						reportWarnMsg(p.fg.topicName, "worker retry item nil")
					} else {
						q.Retry(worker.retryItem, worker.retryDelayMs)
					}

					worker.retryItem = nil
					worker.retry = false
					worker.retryDelayMs = 0

					// 挂起 worker
					confirmHash(hash)
					pendingWorkers = append(pendingWorkers, worker)
					break
				}

				if worker.isAsync {
					wait := worker.maxAsyncWaitMs
					if wait == 0 {
						wait = 30 * 1000
					} else if wait > 60*10*1000 {
						// 异步任务的等待时间太长
						reportWarnMsg(p.fg.name, "async msg: invalid wait time")
						goto doConfirm
					}
					if worker.msgId == "" {
						reportWarnMsg(p.fg.name, "async msg: invalid msgId empty")
						goto doConfirm
					}
					existed := am.add(worker.msgId, worker.seq, worker.lastHash, wait)
					if existed {
						reportWarnMsg(p.fg.name, "async msg: msgId existed")
						goto doConfirm
					}

					// 异步任务，先挂起 worker
					pendingWorkers = append(pendingWorkers, worker)

					// 这个 break 是跳出 select ... case 不是 for
					// https://blog.csdn.net/c359719435/article/details/79624227
					break
				}

			doConfirm:
				st.workerReqCnt++

				confirmHash(hash)
				item, err := popSkipHash()
				if err != nil {
					log.Errorf("err:%v", err)
				}
				if item == nil {
					st.workerReqMissedItemCnt++
					pendingWorkers = append(pendingWorkers, worker)
				} else {
					st.workerReqGotItemCnt++
					dispatchItem(worker, item)
				}
			case <-moreCh:
				assignPendingWorker()
			case <-tk.C:
				am.clearTimeout(func(m *asyncMsg) {
					confirmHash(m.hash)
					q.SetLastFinishTs(m.seq, m.hash)

					log.Warnf("%s: async msg %s timeout", p.fg.name, m.msgId)
					reportAsyncMsgTimeout(p.fg.topicName, p.fg.channelName)
				})

				assignPendingWorker()

				now := utils.Now()
				if st.statAt+10 < now {
					dt := now - st.statAt
					st.statAt = now
					st.dumpAndReset(dt)
				}

				if lastRepAt+5 < now {
					lastRepAt = now
					reportQueueLen(p.fg.topicName, p.fg.channelName, q.GetQueueLen())
				}
			case msgId := <-asyncConfirmCh:
				m := am.del(msgId)
				if m != nil {
					log.Infof("%s: confirm async msg %s ok", p.fg.name, msgId)
					if confirmHash(m.hash) {
						q.SetLastFinishTs(m.seq, m.hash)
						assignPendingWorker()
					}
				} else {
					log.Warnf("%s: confirm async msg %s but not found", p.fg.name, msgId)
				}
			case <-p.exitChan:
				return
			}
		}
	}()
}

func (p *BarrierConsumer) stop() {
	log.Infof("notify %d workers", len(p.workers))
	for _, w := range p.workers {
		w.notifyExit.Store(true)
	}
	p.notifyExit()
	// 确定所有 worker 停了
	for {
		allWorkerExited := true
		for _, w := range p.workers {
			if !w.exited.Load() {
				allWorkerExited = false
				break
			}
		}
		if allWorkerExited {
			break
		}
		time.Sleep(time.Second)
	}
	for {
		if p.exited.Load() {
			break
		}
		time.Sleep(time.Second)
	}
}
