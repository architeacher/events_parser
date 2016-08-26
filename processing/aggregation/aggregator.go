package aggregation

import (
	"time"
)

type GroupingType int
type AggregationType map[GroupingType]interface{}
type AggregationTypeChannel chan AggregationType

type Aggregator struct {
	dailyActiveUsers       map[string]int
	top10WeeklyViewedUsers map[string][]string
	averageUserSession     int
}

const (
	TYPE_GROUPING_BY_TIME GroupingType = iota
	TYPE_GROUPING_BY_USER
	TYPE_GROUPING_BY_IMPRESSION

	KEY_TIME string = "time"

	KEY_WEEK    string = "week"
	KEY_YEAR    string = "year"
	KEY_ACTOR   string = "actor"
	KEY_SUBJECT string = "subject"
)

var (
	AggregationQueue AggregationTypeChannel
)

func NewAggregator() *Aggregator {

	// Todo: Dependencies should be injected
	return &Aggregator{}
}

func (self *Aggregator) AggregateBy(groupingType GroupingType, data interface{}, aggregationEndPoint AggregationType) {

	switch groupingType {
	case TYPE_GROUPING_BY_TIME:
		self.aggregateByDay(data, aggregationEndPoint)
	case TYPE_GROUPING_BY_IMPRESSION:
		self.aggregateByImpression(data, aggregationEndPoint)
	case TYPE_GROUPING_BY_USER:
		self.aggregateByUser(data, aggregationEndPoint)
	}
}

func (*Aggregator) aggregateByDay(data interface{}, aggregationEndPoint AggregationType) {

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

func (*Aggregator) aggregateByImpression(data interface{}, aggregationEndPoint AggregationType) {

}

func (*Aggregator) aggregateByUser(data interface{}, aggregationEndPoint AggregationType) {

}
