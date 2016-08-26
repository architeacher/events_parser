package reducers

import (
	"splash/processing/aggregation"
)

func Reducer(input chan interface{}, output chan interface{}) {
	results := aggregation.AggregationType{}
	aggregator := aggregation.NewAggregator()

	for matches := range input {

		for key, value := range matches.(map[aggregation.GroupingType]interface {}) {
			aggregator.AggregateBy(key, value, results)
			// Todo: Remove this.
			//fmt.Println(key, value, len(results[aggregation.TYPE_GROUPING_BY_TIME].(map[string][]string)["2015-11-18"]))
		}
	}
	output <- results
}
