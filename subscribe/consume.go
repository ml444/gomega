package subscribe

import (
	"fmt"
	log "github.com/ml444/glog"
	"github.com/ml444/scheduler/brokers"
	"sync"
)

const defaultConsumeConcurrentCount = 100

type ConcurrentConsume struct {
	cfg        *Config
	wg         sync.WaitGroup
	retryList  *brokers.MinHeap
	workers    []*Worker
	msgChan    chan *brokers.Item
	finishChan chan *brokers.Item
	writer     brokers.IBackendWriter
	queueGroup brokers.IQueueGroup
	isExist    bool
}

func NewConcurrentConsume(cfg *Config) *ConcurrentConsume {
	return &ConcurrentConsume{cfg: cfg}
}

func (c *ConcurrentConsume) init() {
	if c.cfg.ConcurrentCount == 0 {
		c.cfg.ConcurrentCount = defaultConsumeConcurrentCount
	}
	c.msgChan = make(chan *brokers.Item, 1024)
	c.finishChan = make(chan *brokers.Item, 1024)
}

func (c *ConcurrentConsume) Start() {
	c.init()
	go func() {
		for !c.isExist {
			msg := <-c.finishChan
			c.writer.SetFinish(msg)
		}
	}()
	for i := 0; i < int(c.cfg.ConcurrentCount); i++ {
		w := NewConsumeWorker(&c.wg, i, c.msgChan, c.finishChan)
		c.workers = append(c.workers, w)
		c.wg.Add(1)
		//go w.Run()
	}
	if c.retryList.Len() > 0 {
		for {
			msg := c.retryList.PopEl()
			if msg == nil {
				break
			}
			c.msgChan <- msg.Value.(*brokers.Item)
		}
	}
	for {
		item, err := c.queueGroup.SequentialRead()
		if err != nil {
			log.Error(err)
			continue
		}
		if item == nil {
			continue
		}
		c.msgChan <- item
	}
}
func (c *ConcurrentConsume) Stop() {
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

type SerialConsume struct {
	cfg         *Config
	wg          sync.WaitGroup
	workerMap   map[int]*Worker
	heapMap     map[int]*brokers.MinHeap
	msgChanMap  map[int]chan *brokers.Item
	finishChan  chan *brokers.Item
	retryList   []*brokers.Item
	writer      brokers.IBackendWriter
	queueGroup  brokers.IQueueGroup
	workerCount uint32
	isExist     bool
}

func NewSerialConsume(cfg *Config) *SerialConsume {
	return &SerialConsume{
		cfg:         cfg,
		wg:          sync.WaitGroup{},
		workerMap:   nil,
		heapMap:     nil,
		msgChanMap:  nil,
		finishChan:  nil,
		retryList:   nil,
		workerCount: 0,
	}
}

func (c *SerialConsume) init() {
	if c.cfg == nil {
		panic("Cfg is nil")
	}
	if c.cfg.ConcurrentCount == 0 {
		c.workerCount = defaultConsumeConcurrentCount
	} else {
		c.workerCount = c.cfg.ConcurrentCount
	}
	c.workerMap = map[int]*Worker{}
	c.msgChanMap = map[int]chan *brokers.Item{}
	c.finishChan = make(chan *brokers.Item, 1024)
}

func (c *SerialConsume) Start() {
	c.init()
	go func() {
		for !c.isExist {
			msg := <-c.finishChan
			c.writer.SetFinish(msg)
		}
	}()
	for i := 0; i < int(c.workerCount); i++ {
		// TODO chan
		ch := make(chan *brokers.Item, 1024)
		c.msgChanMap[i] = ch
		w := NewConsumeWorker(&c.wg, i, ch, c.finishChan)
		c.workerMap[i] = w
		c.wg.Add(1)
		//go w.Run()
	}

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
func (c *SerialConsume) specifyPartitionSend(partition int) {
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
func (c *SerialConsume) selectEmit(msg *brokers.Item) {
	hashSize := msg.HashCode
	idx := int(hashSize / uint64(c.workerCount))
	c.msgChanMap[idx] <- msg
}
func (c *SerialConsume) Stop() {
	for _, w := range c.workerMap {
		w.notifyExit()
	}
	c.wg.Wait()
	for _, w := range c.workerMap {
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

type AsyncConsume struct {
	maxAsyncWaitMs uint32
}
type TimingConsume struct {
	// 时间轮算法
}
