package mappers

import (
	"splash/queue/jobs"
	"splash/processing/aggregation"
	"splash/services"
	"splash/communication/protocols/protobuf"
)

func Mapper(input interface{}, output chan interface{}) error {

	serviceLocator := services.NewLocator()

	results := map[aggregation.GroupingType]interface{}{}

	job := input.(*jobs.Job)

	payload := job.GetPayload().(*jobs.Payload)
	eventType := payload.GetType()

	time := serviceLocator.GetAsTimestamp(payload.GetTime())

	year, week := time.ISOWeek()

	results[aggregation.TYPE_GROUPING_BY_TIME] = map[string]interface{}{
		aggregation.KEY_TIME : time,
		aggregation.KEY_ACTOR: payload.GetActorId(),
	}

	switch eventType {
	case protobuf.Event_IMPRESSION:
		results[aggregation.TYPE_GROUPING_BY_IMPRESSION] = map[string]interface{}{
			aggregation.KEY_WEEK: week,
			aggregation.KEY_YEAR: year,
			aggregation.KEY_SUBJECT: payload.GetSubjectId(),
		}
		break
	}

	output <- results

	return nil
}
