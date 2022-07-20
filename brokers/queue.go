package brokers

import (
	"errors"
	"fmt"
	log "github.com/ml444/glog"
	"github.com/ml444/scheduler/config"
	"io"
	"os"
	"path/filepath"
	"strconv"
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

type Queue struct {
	filePrefix string
	Partition  int
	cursor     uint64
	offset     uint64
	beginSeq   uint64
	endSeq     uint64
	items      []*Item
	indexFile  *os.File
	dataFile   *os.File
	//Len       uint64

	MaxItemCount uint32
	MaxDataBytes uint32
}

func NewQueue(baseDir string, parentPath string) (*Queue, error) {
	q := Queue{
		filePrefix: filepath.Join(baseDir, parentPath),
		cursor:   0,
		offset:   0,
		beginSeq: 1,
		endSeq:   0,
		items:    nil,
		//indexFile: nil,
		//dataFile: nil,
		//Len:          0,
		MaxItemCount: 1024,
		MaxDataBytes: 64 * 1024 * 1024,
	}
	err := q.Init()
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func (q *Queue) Init() error {
	var err error
	q.indexFile, err = q.openIndexFile()
	if err != nil {
		log.Error(err)
		return err
	}
	q.dataFile, err = q.openDataFile()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (q *Queue) Push(item *Item) error {
	q.items = append(q.items, item)
	return nil
}

func (q *Queue) Pop() (*Item, error) {
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

func (q *Queue) PopWait() (*Item, error) {
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

func (q *Queue) openIndexFile() (*os.File, error) {
	path := filepath.Join(q.filePrefix, fmt.Sprintf("%d.idx", q.beginSeq))
	return os.OpenFile(path, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
}

func (q *Queue) openDataFile() (*os.File, error) {
	path := filepath.Join(q.filePrefix, fmt.Sprintf("%d.dat", q.beginSeq))
	return os.OpenFile(path, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
}

func (q *Queue) FillIndex(fillCount uint64) error {
	if fillCount == 0 {
		return nil
	}
	if q.indexFile == nil {
		idxFile, err := q.openIndexFile()
		if err != nil {
			log.Error(err)
			return err
		}
		q.indexFile = idxFile
	}
	size := fillCount * indexItemSize
	var buf = make([]byte, size)
	n, err := q.indexFile.ReadAt(buf, int64(q.offset)*indexItemSize)
	if err != nil && err != io.EOF {
		log.Error(err)
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
			log.Error(err)
			return err
		}
		err = q.Push(&item)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}
func (q *Queue) FillData(item *Item) error {
	d := q.dataFile
	if d == nil {
		return errors.New("file not opened")
	}
	dataBuf := make([]byte, item.Size+4)
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

func (q *Queue) Close() {

}

type QueueGroup struct {
	partitions   uint32
	queueCursor  uint32
	queueMaxSize uint64
	queueList    []*Queue
	tk           *time.Ticker
	closeChan    chan struct{}
	baseDir      string
}

func NewQueueGroup(namespace, topic string, partitions uint32) (*QueueGroup, error) {
	dir := filepath.Join(config.GlobalCfg.Broker.BasePath, namespace, topic)
	group := QueueGroup{
		baseDir:      dir,
		partitions:   partitions,
		queueMaxSize: config.GlobalCfg.Broker.QueueMaxSize,
		queueCursor:  0,
	}
	err := group.Init(int(partitions))
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (g *QueueGroup) Init(queueCount int) error {
	g.closeChan = make(chan struct{}, 1)
	g.tk = time.NewTicker(time.Second * 5)
	for i := 0; i < queueCount; i++ {
		q, err := NewQueue(g.baseDir, strconv.FormatInt(int64(i), 10))
		if err != nil {
			return err
		}
		g.queueList = append(g.queueList, q)
	}
	go g.IoLoop()
	return nil
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
					log.Error(err)
				}
			}
		case <-g.closeChan:
			return
		}
	}
}

func (g *QueueGroup) Close() {
	g.closeChan <- struct{}{}
	for _, queue := range g.queueList {
		queue.Close()
	}
}

func (g *QueueGroup) SequentialRead() (*Item, error) {
	defer func() {
		if g.queueCursor+1 < g.partitions {
			g.queueCursor++
		} else {
			g.queueCursor = 0
		}
	}()
	q := g.queueList[g.queueCursor]
	return q.PopWait()
}

func (g *QueueGroup) SpecifyRead(partition int) (*Item, error) {
	if partition+1 > int(g.partitions) {
		return nil, fmt.Errorf("partition [%d] out of range", partition)
	}
	q := g.queueList[partition]
	return q.Pop()
}
