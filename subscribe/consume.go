package subscribe

import (
	"github.com/ml444/scheduler/backend"
	"github.com/ml444/scheduler/structure"
	"sync"
)

type ConsumeReq struct {
	CreatedAt uint32 `protobuf:"varint,1,opt,name=created_at,json=createdAt" json:"created_at,omitempty"`
	RetryCnt  uint32 `protobuf:"varint,2,opt,name=retry_cnt,json=retryCnt" json:"retry_cnt,omitempty"`
	Data      []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	CorpId    uint32 `protobuf:"varint,4,opt,name=corp_id,json=corpId" json:"corp_id,omitempty"`
	AppId     uint32 `protobuf:"varint,5,opt,name=app_id,json=appId" json:"app_id,omitempty"`
	MsgId     string `protobuf:"bytes,6,opt,name=msg_id,json=msgId" json:"msg_id,omitempty"`
	// @desc: 服务端最大重试次数，业务可以根据这个来做最后失败逻辑
	MaxRetryCnt          uint32 `protobuf:"varint,7,opt,name=max_retry_cnt,json=maxRetryCnt" json:"max_retry_cnt,omitempty"`
	Timeout              uint32
	XXX_NoUnkeyedLiteral struct{} `json:"-" bson:"-"`
	XXX_unrecognized     []byte   `json:"-" bson:"-"`
	XXX_sizecache        int32    `json:"-" bson:"-"`
}
type ConsumeRsp struct {
	// @desc: 是否需要重试
	Retry bool `protobuf:"varint,1,opt,name=retry" json:"retry,omitempty"`
	// @desc: 重试间隔时间，可以这里指定，不指定的话，由 smq 来定 (单位:秒)
	RetryIntervalSeconds int64 `protobuf:"varint,2,opt,name=retry_interval_seconds,json=retryIntervalSeconds" json:"retry_interval_seconds,omitempty"`
	// @desc: 为 true 时，不要累加重试次数
	SkipIncRetryCount bool `protobuf:"varint,3,opt,name=skip_inc_retry_count,json=skipIncRetryCount" json:"skip_inc_retry_count,omitempty"`
	// @desc: 支持毫秒级重试时间
	RetryIntervalMs int64 `protobuf:"varint,4,opt,name=retry_interval_ms,json=retryIntervalMs" json:"retry_interval_ms,omitempty"`
	// @desc: 是否异步消费
	IsAsync bool `protobuf:"varint,5,opt,name=is_async,json=isAsync" json:"is_async,omitempty"`
	// @desc: 异步消费最大等待确认时间
	MaxAsyncWaitMs int64 `protobuf:"varint,6,opt,name=max_async_wait_ms,json=maxAsyncWaitMs" json:"max_async_wait_ms,omitempty"`
	// @desc: 消费端服务cpu使用率
	CpuPercent float32 `protobuf:"fixed32,7,opt,name=cpu_percent,json=cpuPercent" json:"cpu_percent,omitempty"`
	// @desc: 消费端服务内存使用率
	MemoryPercent        float32 `protobuf:"fixed32,8,opt,name=memory_percent,json=memoryPercent" json:"memory_percent,omitempty"`
	Took                 int64
	XXX_NoUnkeyedLiteral struct{} `json:"-" bson:"-"`
	XXX_unrecognized     []byte   `json:"-" bson:"-"`
	XXX_sizecache        int32    `json:"-" bson:"-"`
}

const defaultConsumeConcurrentCount = 100

type ConcurrentConsume struct {
	cfg       *Config
	wg        sync.WaitGroup
	retryList *structure.MinHeap
	workers   []*consumeWorker
	msgChan   chan *backend.Item
	backend   backend.IBackendReader
}

func (c *ConcurrentConsume) Start() {
	var concurrentCount uint32
	cfg := c.cfg
	if cfg == nil || cfg.ConcurrentCount == 0 {
		concurrentCount = defaultConsumeConcurrentCount
	}
	for i := 0; i < int(concurrentCount); i++ {
		w := NewConsumeWorker(&c.wg, i, c.msgChan, c.backend)
		c.workers = append(c.workers, w)
		c.wg.Add(1)
		go w.Run()
	}
	if c.retryList.Len() > 0 {
		for {
			msg := c.retryList.PopEl()
			if msg == nil {
				break
			}
			c.msgChan <- msg.Value.(*backend.Item)
		}
	}
	for false {
		// TODO read one file
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
	workerMap   map[int]*consumeWorker
	heapMap     map[int]*structure.MinHeap
	msgChanMap  map[int]chan *backend.Item
	backend     backend.IBackendReader
	retryList   []*backend.Item
	workerCount uint32
}

func (c *SerialConsume) Start() {
	var concurrentCount uint32
	cfg := c.cfg
	if cfg == nil || cfg.ConcurrentCount == 0 {
		concurrentCount = defaultConsumeConcurrentCount
	}
	c.workerCount = concurrentCount
	if c.workerMap == nil {
		c.workerMap = map[int]*consumeWorker{}
	}
	if c.msgChanMap == nil {
		c.msgChanMap = map[int]chan *backend.Item{}
	}
	for i := 0; i < int(concurrentCount); i++ {
		// TODO chan
		ch := make(chan *backend.Item, 1024)
		c.msgChanMap[i] = ch
		w := NewConsumeWorker(&c.wg, i, ch, c.backend)
		c.workerMap[i] = w
		c.wg.Add(1)
		go w.Run()
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
	for false {
		// Read N file

	}
}
func (c *SerialConsume) selectEmit(msg *backend.Item) {
	hashSize := msg.Hash
	idx := int(hashSize / uint64(c.workerCount))
	ch := c.msgChanMap[idx]
	ch <- msg
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
				c.retryList = append(c.retryList, el.Value.(*backend.Item))
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
