package broker

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/ml444/gomega/config"
	"github.com/ml444/gomega/log"
)

type IBackendWriter interface {
	Write(item *Item)
	SetFinish(item *Item)
}

var brokerMap map[string]*Broker
var workerMap map[string]*Worker

var mutex = &sync.RWMutex{}

func InitBroker() {
	brokerMap = map[string]*Broker{}
	workerMap = map[string]*Worker{}
}

func GetWorker(workerId string) *Worker {
	mutex.RLock()
	defer mutex.RUnlock()
	return workerMap[workerId]
}
func SetWorker(workerId string, worker *Worker) {
	mutex.Lock()
	defer mutex.Unlock()
	workerMap[workerId] = worker
}

func GetOrCreateBroker(topic string) *Broker {
	//mutex.RLock()
	//defer mutex.RUnlock()
	broker, ok := brokerMap[topic]
	if !ok {
		broker = NewBroker(topic)
		brokerMap[topic] = broker
		go broker.ioLoop()
	}
	return broker
}

type Broker struct {
	cfg *config.BrokerCfg

	topic    string
	sequence uint64

	itemChan chan *Item
	exitChan chan bool

	fileDir  string
	idxFile  *os.File
	dataFile *os.File

	mu             *sync.Mutex
	workerGroupMap map[string]*WorkerGroup
	//groupSubCfgMap map[string]*config.SubCfg

}

func NewBroker(topic string) *Broker {
	itemChSize := config.DefaultQueueMaxSize
	brokerCfg, ok := config.GetBrokerCfg(topic)
	if !ok {
		brokerCfg = &config.BrokerCfg{
			QueueMaxSize:   config.DefaultQueueMaxSize,
			PubCfg:         nil,
			GroupSubCfgMap: map[string]*config.SubCfg{},
		}
		config.SetBrokerCfg(topic, brokerCfg)
	} else {
		itemChSize = brokerCfg.QueueMaxSize
	}
	fileDir := filepath.Join(config.GlobalCfg.RootPath, topic)
	_ = os.MkdirAll(fileDir, 0755)
	return &Broker{
		cfg:            brokerCfg,
		topic:          topic,
		sequence:       0,
		fileDir:        fileDir,
		itemChan:       make(chan *Item, itemChSize),
		exitChan:       make(chan bool, 1),
		mu:             &sync.Mutex{},
		workerGroupMap: map[string]*WorkerGroup{},
	}
}

func (b *Broker) GetConsumeWorker(topic, group string, workerId string) *Worker {
	b.mu.Lock()
	defer b.mu.Unlock()
	workerGroup, ok := b.workerGroupMap[group]
	if !ok {
		cfg, ok := b.cfg.GetGroupSubCfg(group)
		if !ok {
			cfg = &config.SubCfg{
				Topic:          topic,
				Group:          group,
				Version:        1,
				LastSequence:   0,
				IsHashDispatch: false,
			}
			b.cfg.SetGroupSubCfg(group, cfg)
		}
		workerGroup = NewWorkerGroup(topic, cfg)
		b.workerGroupMap[group] = workerGroup
		go workerGroup.Start() // TODO error
	}
	return workerGroup.GetOrCreateWorker(workerId)
}

func (b *Broker) getNextSequence() uint64 {
	b.sequence++
	return b.sequence
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
			// check if there are items in the channel
			var items []*Item
			var size int
			for done := false; !done; {
				select {
				case item := <-b.itemChan:
					items = append(items, item)
					size += len(item.Data)
				default:
					b.Flush(items, size)
					done = true
				}
			}
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
	firstItem := itemList[0]
	firstSeq := firstItem.Sequence
	{
		// ROTATE file
		var size int64
		size, err = b.dataFile.Seek(0, io.SeekEnd)
		if err == nil && size+int64(len(dataBuf)) >= config.GlobalCfg.FileRotatorSize {
			b.dataFile.Close()
			b.dataFile = nil
			b.idxFile.Close()
			b.idxFile = nil
			b.dataFile, err = getDataFile(b.fileDir, firstSeq)
			if err != nil {
				log.Error(err)
				return
			}
			b.idxFile, err = getIdxFile(b.fileDir, firstSeq)
			if err != nil {
				log.Error(err)
				return
			}
		}

		if b.dataFile == nil {
			b.dataFile, err = getDataFile(b.fileDir, firstSeq)
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
			b.idxFile, err = getIdxFile(b.fileDir, firstSeq)
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

func (b *Broker) Stop() {
	b.exitChan <- true
	for _, workerGroup := range b.workerGroupMap {
		workerGroup.Stop()
	}
	// save config
	if b.cfg.PubCfg == nil {
		b.cfg.PubCfg = &config.PubCfg{Topic: b.topic}
	}
	b.cfg.PubCfg.LastSequence = b.sequence
	b.cfg.PubCfg.LastFileSequence = getIdxFileSequence(b.idxFile.Name())
	_ = config.FlushConfig()
}

func getIdxFileSequence(filename string) uint64 {
	s := strings.TrimSuffix(filename, config.IdxFileSuffix)
	seq, _ := strconv.ParseUint(s, 10, 64)
	return seq
}
