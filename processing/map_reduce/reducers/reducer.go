package reducers

import (
	"splash/processing/aggregation"
)

func Reducer(input chan interface{}, output chan interface{}) {
	results := map[aggregation.GroupingType]interface{}{}
	aggregator := aggregation.NewAggregator()

	for matches := range input {
		for key, value := range matches.(map[aggregation.GroupingType]interface{}) {
			//	_, exists := results[key]
			//	if !exists {
			//		results[key] = value
			//	} else {
			//		results[key] = results[key] + value
			//	}
			//for index, event := range value {

				aggregator.AggregateBy(key, value, results)

			//}
		}
	}
	output <- results
}
