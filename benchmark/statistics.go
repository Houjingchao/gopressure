package benchmark

import (
	"github.com/Houjingchao/gopressure/data"
	"fmt"
	"log"
)

type countInfo struct {
	Name         string
	Count        int64
	CostSum      int64
	CostMin      int64
	CostMax      int64
	SuccessCount int64
}

func (c *countInfo) add(info *countInfo) {
	c.Name = info.Name
	c.Count += info.Count
	c.CostSum += info.CostSum
	if info.CostMin < c.CostMin || c.CostMin <= 0 {
		c.CostMin = info.CostMin
	}
	if info.CostMax > c.CostMax {
		c.CostMax = info.CostMax
	}
	c.SuccessCount += info.SuccessCount
}

func (c *countInfo) statistics(report data.ReqResult, cost int64) {
	c.Name = report.Name
	if cost > c.CostMax {
		c.CostMax = cost
	}

	if c.CostMin == 0 {
		c.CostMin = cost
	}

	if cost < c.CostMin {
		c.CostMin = cost
	}

	c.CostSum += cost
	if report.Success {
		c.SuccessCount++
	}
	c.Count++
}

type statInfo struct {
	linkCount *countInfo
	txCount   map[string]*countInfo
}

func (s *statInfo) Add(info *statInfo) {
	s.linkCount.add(info.linkCount)
	for k, v := range info.txCount {
		if s.txCount[k] == nil {
			s.txCount[k] = v
		} else {
			s.txCount[k].add(v)
		}
	}
}

func genReport(info *statInfo, last int64) (reportLog string) {
	linkInfos := info.linkCount
	if linkInfos.Count <= 0 {
		return
	}
	reportLog += fmt.Sprintf("link:%s, link_count: %d, success_count: %d, success_rate: %f, average_cost: %d, min_cost: %d, max_cost: %d qps: %d \n",
		linkInfos.Name,
		linkInfos.Count,
		linkInfos.SuccessCount,
		float64(linkInfos.SuccessCount)/float64(linkInfos.Count),
		linkInfos.CostSum/linkInfos.Count,
		linkInfos.CostMin,
		linkInfos.CostMax,
		linkInfos.Count/last,
	)
	for k, v := range info.txCount {
		reportLog += fmt.Sprintf("tx:%s, count: %d, success_count: %d, success_rate: %f, average_cost: %d, min_cost: %d, max_cost: %d qps: %d \n",
			k,
			v.Count,
			v.SuccessCount,
			float64(v.SuccessCount)/float64(v.Count),
			v.CostSum/v.Count,
			v.CostMin,
			v.CostMax,
			v.Count/last,
		)
	}
	log.Println(reportLog)
	return
}