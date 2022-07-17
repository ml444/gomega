package subscribe

import (
	"fmt"
	"github.com/ml444/scheduler/brokers"
	"github.com/ml444/scheduler/pb"
	"io"
	"sync"
	"time"
)

const (
	WorkerTypeConcurrent = 1
	WorkerTypeSerial     = 2
	WorkerTypeAsync      = 3
)

type Worker struct {
	name       string
	wg         *sync.WaitGroup
	exitChan   chan int
	msgChan    chan *brokers.Item
	finishChan chan *brokers.Item
	retryList  *brokers.MinHeap
	tk         *time.Ticker
	Cfg        *SubConfig
	//futureList *structure.Tree // async
	blockLimit int
	isSync     bool
}

const defaultTimeout = time.Millisecond * 10

const (
	defaultConsumeMaxExecTimeSeconds = 60
	defaultConsumeMaxRetryCount      = 5
)

func NewConsumeWorker(name string, wg *sync.WaitGroup, msgChan *chan *brokers.Item, finishChan *chan *brokers.Item) *Worker {
	return &Worker{
		name:       name,
		wg:         wg,
		exitChan:   make(chan int, 1),
		msgChan:    *msgChan,
		finishChan: *finishChan,
		retryList:  brokers.NewMinHeap(),
		tk:         time.NewTicker(defaultTimeout),
	}
}
func (w *Worker) Init() {
	w.exitChan = make(chan int, 1)
	w.retryList = brokers.NewMinHeap()
	w.tk = time.NewTicker(defaultTimeout)
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
		//case <-w.tk.C:
		//	//TODO no tk
		//	return nil
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

func (w *Worker) ConsumeMsg() {}

func (w *Worker) Run(firstReq *pb.ConsumeReq, stream pb.OmegaService_ConsumeServer) error {
	defer w.wg.Done()
	if w.isSync {
		return w.syncRun(firstReq, stream)
	} else {
		return w.asyncRun(stream)
	}
}
func (w *Worker) syncRun(firstReq *pb.ConsumeReq, stream pb.OmegaService_ConsumeServer) error {

	var err error
	var req *pb.ConsumeReq
	var prvItem *brokers.Item
	for {
		if firstReq != nil {
			req = firstReq
			firstReq = nil
		} else {
			req, err = stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
		}
		if req == nil {
			return nil
		}
		if prvItem != nil {
			if req.IsRetry && prvItem.Sequence == req.Sequence {
				maxRetryCount := w.getMaxRetryCount()
				if prvItem.RetryCount >= maxRetryCount {
					w.setFinish(prvItem)
				} else {
					var waitMs int64
					if req.RetryIntervalMs > 0 {
						waitMs = req.RetryIntervalMs
					} else {
						waitMs = w.getNextRetryWait(prvItem.RetryCount) * 1000
					}
					execAt := time.Now().UnixMilli() + waitMs
					w.retryList.PushEl(&brokers.MinHeapElement{
						Value:    prvItem,
						Priority: execAt,
					})
				}
			} else {
				//w.setFinish(prvItem)
			}
		}

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
		fmt.Println("===> msg:", msg.Sequence)
		err = stream.Send(&pb.ConsumeRsp{
			Partition:  msg.Partition,
			Sequence:   msg.Sequence,
			Data:       msg.Data,
			RetryCount: msg.RetryCount,
		})
		if err != nil {
			return err
		}
		prvItem = msg
		// TODO Skip blackList msg
		//var consumeRsp call.ConsumeRsp
		//_ = w.ConsumeMsg(msg, payload, &consumeRsp)
	}
	return nil
}

func (w *Worker) asyncRun(stream pb.OmegaService_ConsumeServer) error {
	var err error
	var isExist bool
	go func(stream pb.OmegaService_ConsumeServer) {
		var in *pb.ConsumeReq
		for {
			in, err = stream.Recv()
			if err == io.EOF {
				isExist = true
				return
			}
			if err != nil {
				isExist = true
				return
			}
			setFinish(in)
		}
	}(stream)

	for !isExist {
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
		fmt.Println("===> msg:", msg.Sequence)
		err = stream.Send(&pb.ConsumeRsp{
			Partition:  msg.Partition,
			Sequence:   msg.Sequence,
			Data:       msg.Data,
			RetryCount: msg.RetryCount,
		})
		if err != nil {
			return err
		}
		//prvItem = msg
	}
	return err
}

//type retryItem struct {
//	item       *brokers.Item
//	nextExecAt int64
//}

//func (w *Worker) ConsumeMsg(item *brokers.Item, payload *brokers.MsgPayload, consumeRsp *call.ConsumeRsp) error {
//	// TODO getRoute(checkRoute())
//	ctx := context.TODO()
//	s := w.S
//	var timeoutSeconds = w.Cfg.MaxExecTimeSeconds
//	if timeoutSeconds == 0 {
//		timeoutSeconds = defaultConsumeMaxExecTimeSeconds
//	}
//	meta := &call.MsgMeta{
//		CreatedAt: item.CreatedAt,
//		RetryCnt:  item.RetryCount,
//		Data:      payload.Data,
//		MsgId:     payload.MsgId,
//	}
//	if s.BeforeProcess != nil {
//		s.BeforeProcess(ctx, meta)
//	}
//	in, err := s.UnMarshalRequest(payload.Data)
//	if err != nil {
//		log.Error(err)
//		return err
//	}
//	out := s.NewResponse()
//	err = call.Call(ctx, w.S.Route, &in, &out, timeoutSeconds)
//	if err != nil {
//		log.Error(err)
//		consumeRsp.Retry = true
//		item.RetryCount++
//		return err
//	}
//	if s.AfterProcess != nil {
//		isRetry, isIgnoreRetryCount := s.AfterProcess(ctx, meta, &in, &out)
//		if isRetry {
//			if consumeRsp != nil {
//				consumeRsp.Retry = true
//			}
//		}
//		if !isIgnoreRetryCount {
//			item.RetryCount++
//		}
//	}
//	return nil
//}

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
		//if s.RetryIntervalStep > 0 {
		//	return retryMs * s.RetryIntervalStep
		//}
		return int64(retryMs)
	}
	if retryCnt == 0 {
		retryCnt = 1
	}
	return int64(retryCnt) * 5
}
