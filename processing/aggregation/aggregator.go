package aggregation

import (
	"time"
)

type GroupingType int

const (
	TYPE_GROUPING_BY_TIME GroupingType = iota
	TYPE_GROUPING_BY_USER
	TYPE_GROUPING_BY_IMPRESSION

	KEY_TIME string = "time"

	KEY_WEEK string = "week"
	KEY_YEAR string = "year"
	KEY_ACTOR string = "actor"
	KEY_SUBJECT string = "subject"
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

func (self *Aggregator) AggregateBy(groupingType GroupingType, data interface{}, aggregationEndPoint map[GroupingType]interface{}) {

	switch groupingType {
	case TYPE_GROUPING_BY_TIME:
		self.aggregateByDay(data, aggregationEndPoint)
	case TYPE_GROUPING_BY_IMPRESSION:
		self.aggregateByImpression(data, aggregationEndPoint)
	case TYPE_GROUPING_BY_USER:
		self.aggregateByUser(data, aggregationEndPoint)
	}
}

func (*Aggregator) aggregateByDay(data interface{}, aggregationEndPoint map[GroupingType]interface{}) {
	//day := time.Parse()
	//aggregationEndPoint[day]

	inputData := data.(map[string]interface{})
	actionTime := inputData[KEY_TIME].(time.Time)
	day := actionTime.Format("2006-01-02")

	if aggregationEndPoint[TYPE_GROUPING_BY_TIME] == nil {
		aggregationEndPoint[TYPE_GROUPING_BY_TIME] = make(map[string][]string)
	}

	storedData := aggregationEndPoint[TYPE_GROUPING_BY_TIME].(map[string][]string)

	storedData[day] = append(storedData[day], inputData[KEY_ACTOR].(string))

	aggregationEndPoint[TYPE_GROUPING_BY_TIME] = storedData
}

func (*Aggregator) aggregateByImpression(data interface{}, aggregationEndPoint map[GroupingType]interface{}) {

}

func (*Aggregator) aggregateByUser(data interface{}, aggregationEndPoint interface{}) {

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