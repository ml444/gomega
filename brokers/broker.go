package brokers

import (
	"errors"
	"fmt"
	log "github.com/ml444/glog"
	"github.com/ml444/scheduler/config"
	"os"
	"path/filepath"
	"strconv"
)

/*
{
	"127.0.0.1:9091": {
		"topic1": {
			"partition1": "broker1",
			"partition2": "broker2"
		}
	},
	"127.0.0.1:9092": {
		"topic2": {
			"partition1": "broker3",
			"partition2": "broker4"
		}
	}
}
*/
type IBackendReader interface {
	Read()
}

type IBackendWriter interface {
	Write(item *Item)
	SetFinish(item *Item)
}

var brokerMap map[string]*Broker

func InitBroker() {
	if brokerMap == nil {
		brokerMap = map[string]*Broker{}
	}
}

func GetBroker(namespace, topic string, partition uint32) (*Broker, error) {
	key := fmt.Sprintf("%s:%s:%d", namespace, topic, partition)
	broker, ok := brokerMap[key]
	if !ok {
		broker = NewBroker(namespace, topic, partition)
		// TODO goroutine
		go broker.ioLoop()
		if broker != nil {
			brokerMap[key] = broker
			return broker, nil
		}
		return nil, errors.New("not found broker")
	}
	return broker, nil
}

type Broker struct {
	//namespace string
	//topic     string
	//partition uint32
	sequence  uint64

	fileDir  string
	idxFile  *os.File
	dataFile *os.File

	itemChan chan *Item
	exitChan chan bool
}

func NewBroker(namespace, topic string, partition uint32) *Broker {
	return &Broker{
		//partition: partition,
		sequence:  0,
		fileDir:   filepath.Join(config.GlobalCfg.Broker.BasePath, namespace, topic, strconv.FormatUint(uint64(partition), 10)),
		//idxFile:   nil,
		//dataFile:  nil,
		itemChan:  make(chan *Item, 1024),
		exitChan:  make(chan bool, 1),

	}
}

func (b *Broker) getNextSequence() uint64 {
	// TODO Lock
	b.sequence++
	return b.sequence
}

func (b *Broker) getIdxFile() (*os.File, error) {
	path := filepath.Join(b.fileDir, fmt.Sprintf("%d.idx", b.sequence))
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (b *Broker) getDataFile() (*os.File, error) {
	path := filepath.Join(b.fileDir, fmt.Sprintf("%d.dat", b.sequence))
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}
func (b *Broker) Send(item *Item) error {
	item.Sequence = b.getNextSequence()
	//item.Offset = b.
	b.itemChan <- item
	return nil
}

func (b *Broker) ioLoop() {
	for {
		select {
		case item0 := <-b.itemChan:
			var items []*Item
			size := len(item0.Data)
			items = append(items, item0)
			for {
				select {
				case item := <-b.itemChan:
					items = append(items, item)
					size += len(item.Data)
					if size >= 5*1024*1024 {
						goto FLUSH
					}
				default:
					goto FLUSH
				}
			}
		FLUSH:
			b.Flush(items, size)
		case <-b.exitChan:
			return
		}
	}
}

func (b *Broker) Flush(itemList []*Item, dataSize int) {
	var err error
	itemsLen := len(itemList)
	idxBufSize := itemsLen * indexItemSize
	idxBuf := make([]byte, idxBufSize)
	dataBufSize := itemsLen*dataItemExtraSize + dataSize
	dataBuf := make([]byte, dataBufSize)
	var begin int
	var end int
	for i, item := range itemList {
		idxBufSlice := idxBuf[i*indexItemSize : (i+1)*indexItemSize]
		item.Marshal2Index(idxBufSlice)
		end = begin + dataItemExtraSize + len(item.Data)
		dataBufSlice := dataBuf[begin:end]
		item.Marshal2Data(dataBufSlice)
		begin = end
	}
	fmt.Println("===> ", b.fileDir)
	{
		if b.dataFile == nil {
			b.dataFile, err = b.getDataFile()
			if err != nil {
				log.Error(err)
				return
			}
		}
		var n, m int
		n, err = b.dataFile.Write(dataBuf)
		if err != nil {
			log.Error(err)
			return
		}
		for n < dataBufSize {
			m, err = b.dataFile.Write(dataBuf[n:])
			if err != nil {
				log.Error(err)
				return
			}
			n += m
		}
	}
	{
		if b.idxFile == nil {
			b.idxFile, err = b.getIdxFile()
			if err != nil {
				log.Error(err)
				return
			}
		}
		var n, m int
		n, err = b.idxFile.Write(idxBuf)
		if err != nil {
			log.Error(err)
			return
		}
		for n < idxBufSize {
			m, err = b.idxFile.Write(idxBuf[n:])
			if err != nil {
				log.Error(err)
				return
			}
			n += m
		}
	}
}
