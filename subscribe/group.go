package subscribe

import (
	"errors"
	"fmt"
	log "github.com/ml444/glog"
	"github.com/ml444/scheduler/brokers"
	"github.com/ml444/scheduler/pb"
	"sync"
)

const defaultConsumeConcurrentCount = 100

type IGroup interface {
	Start()
	Stop()
	AddWorker(name string) *Worker
}

func NewConsumeGroup(policy pb.Policy, cfg *SubConfig, partitions uint32) (IGroup, error) {
	switch policy {
	case pb.Policy_PolicyConcurrence:
		return NewConcurrentConsume(cfg, partitions)
	case pb.Policy_PolicySerial:
		return NewSerialConsume(cfg, partitions)
	default:
		return NewConcurrentConsume(cfg, partitions)
	}
}

type ConcurrentGroup struct {
	cfg         *SubConfig
	wg          sync.WaitGroup
	retryList   *brokers.MinHeap
	queueGroup  brokers.IQueueGroup
	writer      brokers.IBackendWriter
	workers     []*Worker
	workerCount int
	msgChan     chan *brokers.Item
	finishChan  chan *brokers.Item
	isExist     bool
}

func NewConcurrentConsume(cfg *SubConfig, partitions uint32) (*ConcurrentGroup, error) {
	if cfg == nil {
		return nil, errors.New("subscriber config is nil")
	}
	var err error
	c := ConcurrentGroup{cfg: cfg}
	c.queueGroup, err = brokers.NewQueueGroup(cfg.Namespace, cfg.Topic, partitions)
	if err != nil {
		return nil, err
	}
	c.init()
	return &c, nil
}

func (c *ConcurrentGroup) init() {
	c.msgChan = make(chan *brokers.Item, 1024)
	c.finishChan = make(chan *brokers.Item, 1024)
	c.retryList = brokers.NewMinHeap()
}

func (c *ConcurrentGroup) Start() {
	go func() {
		for !c.isExist {
			msg := <-c.finishChan
			c.writer.SetFinish(msg)
		}
	}()
	//for i := 0; i < int(c.cfg.ConcurrentCount); i++ {
	//	w := NewConsumeWorker(&c.wg, i, c.msgChan, c.finishChan)
	//	c.workers = append(c.workers, w)
	//	c.wg.Add(1)
	//	//go w.Run()
	//}
	if c.retryList.Len() > 0 {
		for {
			msg := c.retryList.PopEl()
			if msg == nil {
				break
			}
			c.msgChan <- msg.Value.(*brokers.Item)
		}
	}
	//var testSeq uint64
	for {
		item, err := c.queueGroup.SequentialRead()
		if err != nil {
			log.Error(err)
			continue
		}
		if item == nil {
			continue
		}
		//testSeq++
		//item := &brokers.Item{
		//	Sequence:   testSeq,
		//	HashCode:   testSeq,
		//	CreatedAt:  0,
		//	Offset:     0,
		//	Size:       0,
		//	RetryCount: 0,
		//	DelayType:  0,
		//	DelayValue: 0,
		//	Priority:   0,
		//	Data:       []byte("wahaha"),
		//}
		fmt.Println("===>", item.Sequence)
		c.msgChan <- item
	}
}
func (c *ConcurrentGroup) AddWorker(name string) *Worker {
	w := NewConsumeWorker(name, &c.wg, &c.msgChan, &c.finishChan)
	c.workers = append(c.workers, w)
	c.wg.Add(1)
	c.workerCount++
	// TODO reBalance
	return w
}
func (c *ConcurrentGroup) Stop() {
	for _, w := range c.workers {
		w.notifyExit()
	}
	c.wg.Wait()
	for _, w := range c.workers {
		if w.retryList != nil {
			for {
				el := w.retryList.PopEl()
				if el == nil {
					break
				}
				c.retryList.PushEl(el)
				//c.retryList = append(c.retryList, el.)
			}
		}
	}
}

type SerialGroup struct {
	cfg         *SubConfig
	wg          sync.WaitGroup
	workers     []*Worker
	heapMap     map[int]*brokers.MinHeap
	msgChanMap  map[int]chan *brokers.Item
	finishChan  chan *brokers.Item
	retryList   []*brokers.Item
	writer      brokers.IBackendWriter
	queueGroup  brokers.IQueueGroup
	workerCount uint32
	isExist     bool
}

func NewSerialConsume(cfg *SubConfig, partitions uint32) (*SerialGroup, error) {
	g := SerialGroup{
		cfg:         cfg,
		wg:          sync.WaitGroup{},
		heapMap:     nil,
		msgChanMap:  nil,
		finishChan:  nil,
		retryList:   nil,
		workerCount: 0,
	}
	var err error
	g.queueGroup, err = brokers.NewQueueGroup(cfg.Namespace, cfg.Topic, partitions)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (c *SerialGroup) init() {
	c.msgChanMap = map[int]chan *brokers.Item{}
	c.finishChan = make(chan *brokers.Item, 1024)
}

func (c *SerialGroup) Start() {
	c.init()
	go func() {
		for !c.isExist {
			msg := <-c.finishChan
			c.writer.SetFinish(msg)
		}
	}()
	//for i := 0; i < int(c.workerCount); i++ {
	//	// TODO chan
	//	ch := make(chan *brokers.Item, 1024)
	//	c.msgChanMap[i] = ch
	//	w := NewConsumeWorker(&c.wg, i, ch, c.finishChan)
	//	c.workerMap[i] = w
	//	c.wg.Add(1)
	//	//go w.Run()
	//}

	if len(c.retryList) > 0 {
		for {
			msg := c.retryList[0]
			c.retryList = c.retryList[1:]
			if msg == nil {
				break
			}
			c.selectEmit(msg)
		}
	}
	for i := 0; i < int(c.workerCount); i++ {
		go c.specifyPartitionSend(i)
	}
}

func (c *SerialGroup) AddWorker(name string) *Worker {
	// TODO chan
	ch := make(chan *brokers.Item, 1024)
	w := NewConsumeWorker(name, &c.wg, &ch, &c.finishChan)
	c.msgChanMap[int(c.workerCount)] = ch
	c.workers = append(c.workers, w)
	c.wg.Add(1)
	c.workerCount++
	// TODO reBalance
	return w
}

func (c *SerialGroup) specifyPartitionSend(partition int) {
	ch, ok := c.msgChanMap[partition]
	if !ok {
		log.Error(fmt.Errorf("not found msgChan with partition %d", partition))
		return
	}
	for {
		// TODO: get queue to msg
		item, err := c.queueGroup.SpecifyRead(partition)
		if err != nil {
			log.Error(err)
			continue
		}
		if item == nil {
			continue
		}
		ch <- item
	}
}
func (c *SerialGroup) selectEmit(msg *brokers.Item) {
	hashSize := msg.HashCode
	idx := int(hashSize / uint64(c.workerCount))
	c.msgChanMap[idx] <- msg
}
func (c *SerialGroup) Stop() {
	for _, w := range c.workers {
		w.notifyExit()
	}
	c.wg.Wait()
	for _, w := range c.workers {
		if w.retryList != nil {
			for {
				el := w.retryList.PopEl()
				if el == nil {
					break
				}
				//c.retryList.PushEl(el)
				c.retryList = append(c.retryList, el.Value.(*brokers.Item))
			}
		}
	}
}

type AsyncGroup struct {
	maxAsyncWaitMs uint32
}
type TimingGroup struct {
	// 时间轮算法
}
