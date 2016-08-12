package processing

import (
	"time"
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
}

func NewAggregator() *Aggregator {

	DailyActiveUsers = make(chan string)
	isMergeDone = make(chan bool)
	AggregationQueue = make(chan map[string]int)
	aggregatedData = map[string]int{}

	return &Aggregator{}
}

func (*Aggregator) MonitorNewData() {

	for {
		select {
		case date := <-DailyActiveUsers:
			if _, ok := aggregatedData[date]; !ok {
				aggregatedData[date] = 0
			}

			aggregatedData[date]++
			aggregatedData["total"]++
		case <-time.After(time.Millisecond * time.Duration(5000)):
			isMergeDone <- true
		}
	}
}

func (*Aggregator) Aggregate(operator *Operator) {

	for isMergeDone := range isMergeDone {

		if (isMergeDone) {

			processedData := operator.Operate(aggregatedData)

			AggregationQueue <- processedData
		}
	}
}