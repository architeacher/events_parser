package processing

import (
	"encoding/json"
	"splash/communication"
)

type ResponsePayloadProcessor struct{
}

func (self * ResponsePayloadProcessor) Process(resp *communication.Response) ([]int, error) {

	var numbers []int

	if resp.IsSuccessful() {

		numbersResponse := map[string][]interface{}{}

		err := json.Unmarshal([]byte(resp.Body()), &numbersResponse)

		if (nil != err) {

			return nil, err
		}

		if nil != numbersResponse["Numbers"] {

			for _, value := range numbersResponse["Numbers"] {

				numbers = append(numbers, int(value.(float64)))
			}
		}
	}

	return numbers, nil
}

