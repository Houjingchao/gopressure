package benchmark

import (
	"github.com/Houjingchao/gopressure/interfaces"
	"github.com/getlantern/errors"
	"log"
	"time"
)

const DOC_INTERVAL = 5 // 每5秒打一次点

type BenchMark struct {
	factory     interfaces.Factory
	concurrency int
	last        int64
	debug       bool // 是否是测试
	occupyData  bool
}

func (b *BenchMark) checkFactory() error {
	provider, err := b.factory.GetDataProvider()
	if err != nil {
		log.Println("get dataprovider err ", err)
		return err
	}
	capacity := provider.GetCapacity()
	if capacity < b.concurrency {
		return errors.New("data provider's capacity is not enough. The capacity is %d", capacity)
	}
	return nil
}

func NewbenchMark(factory interfaces.Factory, concurrency int, last int64, debug bool, occupyData bool) *BenchMark {
	return &BenchMark{
		factory:     factory,
		concurrency: concurrency,
		last:        last,
		debug:       debug,
		occupyData:  occupyData,
	}
}

func (b *BenchMark) Run() (reportStr string, err error) {
	if b.occupyData {
		if err = b.checkFactory(); err != nil {
			return
		}
	}
	reportChan := make(chan *statInfo, b.concurrency*2)
	summaryReport := newStatInfo()
	docReport := newStatInfo()

	benchTimer := NewTimer().Start()
	for {
		tmpTimer := NewTimer().Start()
		for i := 0; i < b.concurrency; i++ {
			go newWorker(b.factory, DOC_INTERVAL*time.Second, b.occupyData).start(reportChan)
		}
		for p := 0; p < b.concurrency; p++ {
			ch := <-reportChan
			docReport.Add(ch)
		}
		tmpLast := tmpTimer.CostSecond()
		genReport(docReport, tmpLast)

		summaryReport.Add(docReport)
		docReport = newStatInfo()
		if benchTimer.CostSecond() >= b.last { //跳出
			break
		}
	}
	last := benchTimer.CostSecond()
	reportStr = genReport(summaryReport, last)

	return
}

func newStatInfo() *statInfo {
	return &statInfo{
		linkCount: &countInfo{},
		txCount:   make(map[string]*countInfo),
	}
}
