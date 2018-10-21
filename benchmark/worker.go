package benchmark

import (
	"sync"
	"github.com/Houjingchao/gopressure/interfaces"
	"time"
	"log"
)

type worker struct {
	factory        interfaces.Factory
	last           time.Duration
	occupyData     bool
	startTime      time.Time
	lastDocTime    time.Time
	wg             sync.WaitGroup
	statisticsInfo *statInfo
}

func newWorker(factory interfaces.Factory, last time.Duration, occupyData bool) *worker {
	return &worker{
		factory:        factory,
		last:           last,
		occupyData:     occupyData,
		statisticsInfo: newStatInfo(),
	}
}

func (w *worker) start(reportChan chan *statInfo) {
	w.startTime = time.Now()
	for {
		w.runLink()
		if w.needStop() {
			reportChan <- w.genReport()
			return
		}
	}
}

func (w *worker) needStop() bool {
	return w.startTime.Add(w.last).Before(time.Now())
}

func (w *worker) genReport() *statInfo {
	return w.statisticsInfo
}

func (w *worker) runLink() {
	ctx, err := w.factory.GetContext()
	if err != nil {
		log.Println(err)
		return
	}

	provider, err := w.factory.GetDataProvider()
	if err != nil {
		log.Println(err)
		return
	}
	data := provider.Get()
	if w.occupyData {
		defer provider.Put(data)
	} else {
		provider.Put(data)
	}

	link, err := w.factory.GetLink()
	if err != nil {
		log.Println(err)
		return
	}
	link.Setup(ctx)
	linkTimer := NewTimer().Start()
	for {
		txs, err := link.Next(ctx)
		if len(txs) == 0 || err != nil { // 当
			break
		}
		for _, tx := range txs {
			tmpTx := tx
			w.wg.Add(1)
			go func() {
				defer w.wg.Done()
				t := NewTimer().Start()
				tmpTx.Request(ctx, data)
				txCost := t.CostMillisecond()
				report := tmpTx.Report()
				if w.statisticsInfo.txCount[report.Name] == nil {
					w.statisticsInfo.txCount[report.Name] = &countInfo{}
				}
				w.statisticsInfo.txCount[report.Name].statistics(report, txCost)
			}()
			w.wg.Wait()
		}

	}
	cost := linkTimer.CostMillisecond()
	w.statisticsInfo.linkCount.statistics(link.Report(), cost)
}