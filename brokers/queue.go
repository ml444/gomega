package brokers

import (
	"errors"
	"fmt"
	log "github.com/ml444/glog"
	"github.com/ml444/scheduler/util"
	"io"
	"os"
	"path/filepath"
	"time"
)

type IQueue interface {
	Push(item *Item)
	Pop() *Item
}

type IQueueGroup interface {
	SequentialRead() (*Item, error)
	SpecifyRead(partition int) (*Item, error)
}

type FifoQueue struct {
	util.Report
	Name      string
	Partition int
	cursor    uint64
	offset    uint64
	beginSeq  uint64
	endSeq    uint64
	items     []*Item
	indexFile *os.File
	dataFile  *os.File
	Len       uint64

	MaxItemCount uint32
	MaxDataBytes uint32
}

func NewFifoQueue() *FifoQueue {
	return &FifoQueue{
		cursor:       0,
		offset:       0,
		beginSeq:     0,
		endSeq:       0,
		items:        nil,
		dataFile:     nil,
		Len:          0,
		MaxItemCount: 8192,
		MaxDataBytes: 64 * 1024 * 1024,
	}
}

func (q *FifoQueue) Init() {

}

func (q *FifoQueue) Push(item *Item) error {
	return nil
}

func (q *FifoQueue) Pop() (*Item, error) {
	if len(q.items) == 0 {
		return nil, nil
	}
	item := q.items[0]
	err := q.FillData(item)
	if err != nil {
		return nil, err
	}
	q.items = q.items[1:]
	return item, nil
}

func (q *FifoQueue) PopWait() (*Item, error) {
	item, err := q.Pop()
	if err != nil {
		return nil, err
	}
	if item == nil {
		time.Sleep(time.Millisecond * 10)
		return q.Pop()
	}
	return item, nil
}

func (q *FifoQueue) openIndexFile() (*os.File, error) {
	path := filepath.Join(IndexFilePath, q.Name)
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
}

func (q *FifoQueue) FillIndex(fillCount uint64) error {
	if fillCount == 0 {
		return nil
	}
	if q.indexFile == nil {
		idxFile, err := q.openIndexFile()
		if err != nil {
			q.ReportError(err)
			return err
		}
		q.indexFile = idxFile
	}
	size := fillCount * indexItemSize
	var buf = make([]byte, size)
	n, err := q.indexFile.ReadAt(buf, int64(q.offset)*indexItemSize)
	if err != nil && err != io.EOF {
		q.ReportError(err)
		return err
	}
	if n <= 0 {
		return nil
	}
	need := n / indexItemSize
	for i := 0; i < need; i++ {
		begin := i * indexItemSize
		var item Item
		err = item.FillIndex(buf[begin : begin+indexItemSize])
		if err != nil {
			q.ReportError(err)
			return err
		}
		err = q.Push(&item)
		if err != nil {
			q.ReportError(err)
			return err
		}
	}
	return nil
}
func (q *FifoQueue) FillData(item *Item) error {
	d := q.dataFile
	if d == nil {
		return errors.New("file not opened")
	}
	dataBuf := make([]byte, item.Size)
	var read uint32
	for read < item.Size {
		n, err := d.ReadAt(dataBuf[read:], int64(item.Offset)+int64(read))
		if err != nil {
			if err == io.EOF {
				return errors.New("data file truncated")
			}
			log.Errorf("err:%v", err)
			return err
		} else if n > 0 {
			read += uint32(n)
		} else {
			return errors.New("data file truncated")
		}
	}
	return item.FillData(dataBuf)
}

func (q *FifoQueue) Close() {

}

type QueueGroup struct {
	util.Report
	partitionNum int
	queueMaxSize uint64
	consumeIdx   int
	queueList    []*FifoQueue
	tk           *time.Ticker
	closeChan    chan bool
}

func (g *QueueGroup) Init(queueCount int) {
	g.tk = time.NewTicker(time.Second * 5)
	for i := 0; i < queueCount; i++ {
		q := NewFifoQueue()
		g.queueList = append(g.queueList, q)
	}

}

func (g *QueueGroup) IoLoop() {
	tkc := g.tk.C
	for {
		select {
		case <-tkc:
			for _, queue := range g.queueList {
				fillCount := g.queueMaxSize - (queue.offset - queue.cursor)
				err := queue.FillIndex(fillCount)
				if err != nil {
					g.ReportError(err)
				}
			}
		case <-g.closeChan:
			return

		}
	}
}

func (g *QueueGroup) Close() {
	g.closeChan <- true
	for _, queue := range g.queueList {
		queue.Close()
	}
}

func (g *QueueGroup) SequentialRead() (*Item, error) {
	defer func() {
		if g.consumeIdx+1 < g.partitionNum {
			g.consumeIdx++
		} else {
			g.consumeIdx = 0
		}
	}()
	q := g.queueList[g.consumeIdx]
	return q.PopWait()
}

func (g *QueueGroup) SpecifyRead(partition int) (*Item, error) {
	if partition+1 > g.partitionNum {
		return nil, fmt.Errorf("partition [%d] out of range", partition)
	}
	q := g.queueList[partition]
	return q.Pop()
}
