package broker

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ml444/gomega/config"
	"github.com/ml444/gomega/log"
)

type WorkerGroup struct {
	cfg           *config.SubCfg
	mu            sync.RWMutex
	workerListIdx int
	workerList    []*Worker
	workerMap     map[string]*Worker
	//workerMsgChanList    []chan *Item
	//workerFinishChan chan *Item
	finishChan chan uint64

	latestFile uint64

	done bool

	//minFinishSequence uint64
	finishFileSequence uint64
	topic              string
	fileDir            string
	idxFile            *os.File
	dataFile           *os.File
	needRotate         bool

	itemList         []*Item
	maxFillItemCount int
	sequenceCursor   uint64 // sequenceCursor of offset sequence
	openFileTk       *time.Ticker
	sleepFillTk      *time.Ticker
}

func NewWorkerGroup(topic string, cfg *config.SubCfg) *WorkerGroup {
	return &WorkerGroup{
		cfg:                cfg,
		topic:              topic,
		workerMap:          map[string]*Worker{},
		mu:                 sync.RWMutex{},
		fileDir:            filepath.Join(config.GetRootPath(), topic),
		finishFileSequence: 1,
		sequenceCursor:     cfg.LastSequence,
		maxFillItemCount:   1000,
		openFileTk:         time.NewTicker(time.Second * 1),
		sleepFillTk:        time.NewTicker(time.Millisecond * 100),
	}
}

func (g *WorkerGroup) GetOrCreateWorker(workerId string) *Worker {
	g.mu.Lock()
	defer g.mu.Unlock()
	worker, ok := g.workerMap[workerId]
	if !ok {
		worker = NewWorker(&g.finishChan)
		g.workerMap[workerId] = worker
		g.workerList = append(g.workerList, worker)
		//g.workerMsgChanList = append(g.workerMsgChanList, ch)
	}
	return worker
}

func (g *WorkerGroup) RemoveWorker(workerId string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	worker, ok := g.workerMap[workerId]
	if ok {
		delete(g.workerMap, workerId)
		for i, w := range g.workerList {
			if w == worker {
				g.workerList = append(g.workerList[:i], g.workerList[i+1:]...)
				break
			}
		}
	}
	if len(g.workerMap) == 0 {
		g.done = true
	}
}

func (g *WorkerGroup) Stop() {
	g.done = true
	g.openFileTk.Stop()
	g.sleepFillTk.Stop()
	for ok := false; !ok; {
		if len(g.itemList) == 0 {
			ok = true
		} else {
			time.Sleep(time.Millisecond * 10)
		}
	}
	for _, w := range g.workerList {
		w.Stop()
	}
	g.cfg.LastSequence = g.sequenceCursor
}

func (g *WorkerGroup) SetLatestFileSequence(sequence uint64) {
	g.latestFile = sequence
}

func (g *WorkerGroup) openIndexFile() (*os.File, error) {
	path := filepath.Join(g.fileDir, fmt.Sprintf("%d%s", g.finishFileSequence, config.IdxFileSuffix))
	return os.OpenFile(path, os.O_RDONLY, 0666)
}
func (g *WorkerGroup) openDataFile() (*os.File, error) {
	path := filepath.Join(g.fileDir, fmt.Sprintf("%d%s", g.finishFileSequence, config.DataFileSuffix))
	return os.OpenFile(path, os.O_RDONLY, 0666)
}
func (g *WorkerGroup) openNextFile() error {
	if g.idxFile != nil {
		_ = g.idxFile.Close()
		g.idxFile = nil
	}
	if g.dataFile != nil {
		_ = g.dataFile.Close()
		g.dataFile = nil
	}
AgainIdx:
	idxFile, err := g.openIndexFile()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			select {
			case <-g.openFileTk.C:
				log.Debugf("retry open %d idx file", g.finishFileSequence)
				goto AgainIdx
			}
		}
		log.Errorf("err: %v", err)
		return err
	}
	g.idxFile = idxFile
AgainData:
	dataFile, err := g.openDataFile()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			select {
			case <-g.openFileTk.C:
				log.Debugf("retry open %d data file", g.finishFileSequence)
				goto AgainData
			}
		}
		log.Errorf("err: %v", err)
		return err
	}
	g.dataFile = dataFile
	return nil
}
func (g *WorkerGroup) push(item *Item) error {
	g.itemList = append(g.itemList, item)
	return nil
}
func (g *WorkerGroup) pop() (*Item, error) {
	if len(g.itemList) == 0 {
		return nil, nil
	}
	item := g.itemList[0]
	err := g.fillData(item)
	if err != nil {
		return nil, err
	}
	g.itemList = g.itemList[1:]
	return item, nil
}
func (g *WorkerGroup) fillIndex() error {
	if g.idxFile == nil {
		err := g.openNextFile()
		if err != nil {
			return err
		}
	}
	size := g.maxFillItemCount * indexItemSize
	var buf = make([]byte, size)
	var off int64
	var maybeRotate bool
	if g.sequenceCursor > g.finishFileSequence {
		off = int64(g.sequenceCursor-g.finishFileSequence+1) * indexItemSize
	}
	n, err := g.idxFile.ReadAt(buf, off)
	if err != nil {
		if err == io.EOF {
			log.Debugf("read %d idx file eof", g.finishFileSequence)
			maybeRotate = true
		} else {
			log.Error(err)
			return err
		}
	}

	if n <= 0 {
		return nil
	}
	count := n / indexItemSize
	var lastItem Item
	for i := 0; i < count; i++ {
		begin := i * indexItemSize
		var item Item
		err = item.FillIndex(buf[begin : begin+indexItemSize])
		if err != nil {
			log.Error(err)
			return err
		}
		err = g.push(&item)
		if err != nil {
			log.Error(err)
			return err
		}
		lastItem = item
	}
	if maybeRotate {
		if lastItem.Sequence >= g.latestFile {
			// file is writing
			return nil
		} else {
			g.needRotate = true
		}
	}
	return nil
}
func (g *WorkerGroup) fillData(item *Item) error {
	d := g.dataFile
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
func (g *WorkerGroup) Start() error {
	// read idx file and data file
	searchSeq := g.finishFileSequence
	if g.finishFileSequence == 1 && g.sequenceCursor > g.finishFileSequence {
		searchSeq = g.sequenceCursor
	}

	_, fileFirstSeq, err := searchIdxFileWithSequence(g.fileDir, searchSeq)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	if fileFirstSeq != 0 {
		g.finishFileSequence = fileFirstSeq
	}
	err = g.openNextFile()
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}

	for !g.done {
		if g.needRotate {
			log.Debugf("rotate file, finishFileSequence:%d", g.finishFileSequence)
			g.finishFileSequence = g.sequenceCursor + 1
			err = g.openNextFile()
			if err != nil {
				log.Errorf("err: %v", err)
				return err
			}
			g.needRotate = false
		}
		// start fill index
		err = g.fillIndex()
		if err != nil {
			log.Errorf("err: %v", err)
			return err
		}
		if len(g.itemList) == 0 {
			select {
			case <-g.sleepFillTk.C:
				continue
			}
		}
		workerLen := len(g.workerList)

		for {
			var item *Item
			item, err = g.pop()
			if err != nil {
				log.Errorf("err: %v", err)
				return err
			}
			if item == nil {
				break
			}
			err = g.fillData(item)
			if err != nil {
				log.Errorf("err: %v", err)
				return err
			}
			if g.cfg.IsHashDispatch {
				g.workerListIdx = int(item.HashCode) % workerLen
			} else {
				g.workerListIdx = (g.workerListIdx + 1) % workerLen
			}
			worker := g.workerList[g.workerListIdx]
			worker.ReceiveItem(item)
			g.sequenceCursor++
		}
	}

	return nil
}
