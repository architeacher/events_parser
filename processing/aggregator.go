package processing

import (
	//"time"
	//"sync"
)

type GroupingType int

const (
	TYPE_GROUPING_BY_DAY GroupingType = iota
	TYPE_GROUPING_BY_IMPRESSION
	TYPE_GROUPING_BY_USER
)

var (
	// A buffered channel that we can merge fetched dates on.
	DailyActiveUsers chan string
	// A boolean chanel to indicate that merge is done, and the aggregation process should start.
	isMergeDone chan bool
	// A channel to serve aggregated data.
	AggregationQueue chan map[string]int
	aggregatedData map[string]int
)

type Aggregator struct {
	dailyActiveUsers map[string] int
	top10WeeklyViewedUsers map[string][]string
	averageUserSession int
}

func NewAggregator() *Aggregator {

	// Todo: Dependencies should be injected
	//DailyActiveUsers = make(chan string, 50000)
	//isMergeDone = make(chan bool)
	//AggregationQueue = make(chan map[string]int, 5)
	//aggregatedData = map[string]int{}

	return &Aggregator{}
}

func (self *Aggregator) AggregateBy(groupingType GroupingType, data interface{}) interface{} {

	var output int
	switch groupingType {
	case TYPE_GROUPING_BY_DAY:
		return self.aggregateByDay(data)
	case TYPE_GROUPING_BY_IMPRESSION:
		return self.aggregateByImpression(data)
	case TYPE_GROUPING_BY_USER:
		return self.aggregateByUser(data)
	}

	return output
}

func (*Aggregator) aggregateByDay(data interface{}) map[string]interface{} {
	var output map[string]interface{}
	return output
}

func (*Aggregator) aggregateByImpression(data interface{}) map[string]interface{} {
	var output map[string]interface{}
	return output
}

func (*Aggregator) aggregateByUser(data interface{}) map[string]interface{} {
	var output map[string]interface{}
	return output
}

//func (*Aggregator) MonitorNewData() {
//
//	var lock sync.RWMutex
//
//	for {
//		select {
//		case date := <-DailyActiveUsers:
//			if _, ok := aggregatedData[date]; !ok {
//				aggregatedData[date] = 0
//			}
//
//			lock.Lock()
//			aggregatedData[date]++
//			aggregatedData["total"]++
//			lock.Unlock()
//		case <-time.After(time.Millisecond * time.Duration(5000)):
//			isMergeDone <- true
//		}
//	}
//}
//
//func (*Aggregator) Aggregate(operator *Operator) {
//
//	for isMergeDone := range isMergeDone {
//
//		if (isMergeDone) {
//
//			processedData := operator.Operate(aggregatedData)
//
//			AggregationQueue <- processedData
//		}
//	}
//}