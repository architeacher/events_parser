package reducers

import (
	"fmt"
	"splash/processing/aggregation"
)

func Reducer(input chan interface{}, output chan interface{}) {
	results := map[aggregation.GroupingType]interface{}{}
	aggregator := aggregation.NewAggregator()

	for {
		select {
		//case <- map_reduce.IsPatchFinished:
		//	output <- results

		case matches := <-input:

			for key, value := range matches.(map[aggregation.GroupingType]interface{}) {
				//	_, exists := results[key]
				//	if !exists {
				//		results[key] = value
				//	} else {
				//		results[key] = results[key] + value
				//	}
				//for index, event := range value {

				aggregator.AggregateBy(key, value, results)
				//fmt.Println(key, value, len(results[aggregation.TYPE_GROUPING_BY_TIME].(map[string][]string)["2015-11-18"]))
				//}
			}
		}
	}
	fmt.Println("will never output")
	output <- results
}
