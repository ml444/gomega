package brokers

import (
	"fmt"
	"github.com/ml444/scheduler/publish"
	"github.com/ml444/scheduler/util"
	"os"
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
	Write(req *publish.PubReq, data []byte)
	SetFinish(item *Item)
}

var brokerMap map[string]*Broker

func GetBrokerByTopicName(namespace, topic string) *Broker {
	key := fmt.Sprintf("%s:%s", namespace, topic)
	broker, ok := brokerMap[key]
	if !ok {
		return nil
	}
	return broker
}

type Broker struct {
	util.Report
	topic     string
	partition int
	sequence  uint64

	idxFile  *os.File
	dataFile *os.File

	itemChan chan *Item
	exitChan chan bool
}

func (b *Broker) getNextSequence() uint64 {
	// TODO Lock
	b.sequence++
	return b.sequence
}

func (b *Broker) getIdxFile() (*os.File, error) {
	return nil, nil
}

func (b *Broker) getDataFile() (*os.File, error) {
	return nil, nil
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

	{
		if b.idxFile == nil {
			b.idxFile, err = b.getDataFile()
			if err != nil {
				b.ReportError(err)
				return
			}
		}
		var n, m int
		n, err = b.dataFile.Write(dataBuf)
		if err != nil {
			b.ReportError(err)
			return
		}
		for n < dataBufSize {
			m, err = b.dataFile.Write(dataBuf[n:])
			if err != nil {
				b.ReportError(err)
				return
			}
			n += m
		}
	}
	{
		if b.idxFile == nil {
			b.idxFile, err = b.getIdxFile()
			if err != nil {
				b.ReportError(err)
				return
			}
		}
		var n, m int
		n, err = b.idxFile.Write(idxBuf)
		if err != nil {
			b.ReportError(err)
			return
		}
		for n < idxBufSize {
			m, err = b.idxFile.Write(idxBuf[n:])
			if err != nil {
				b.ReportError(err)
				return
			}
			n += m
		}
	}
}
