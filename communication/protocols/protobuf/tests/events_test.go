package protobuf

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"splash/client"
	"splash/communication/protocols/protobuf"
	"strconv"
	httpProtocol "splash/communication/protocols/http"
	"net/http"
)

func TestShouldBuildPayloadCorrectly (t *testing.T) {

	assert := assert.New(t)

	data := [][]string{
		{"1439003289214921956", "signup", "1241687973565112848"},
		{"1439003289214921956", "follow", "5577006791947779410", "1241687973565112848"},
	}

	for _, eventData := range data {

		payload, err := client.NewClient(httpProtocol.NewProtocol(&http.Client{})).BuildPayload(eventData)
		assert.Nil(err)

		time, err := strconv.ParseInt(eventData[0], 10, 0)

		var eventType protobuf.Event_EventType

		switch eventData[1] {
		case "signup":
			eventType = protobuf.Event_SIGNUP
			break
		case "follow":
			eventType = protobuf.Event_FOLLOW
			break
		case "viorama":
			eventType = protobuf.Event_SPLASH_CREATION
			break
		case "view":
			eventType = protobuf.Event_IMPRESSION
			break
		}

		actorId, err := strconv.ParseInt(eventData[2], 10, 0)
		assert.Nil(err)

		eventPayload := protobuf.Event_Payload{
			Time: &time,
			EventType: &eventType,
			ActorId: &actorId,
		}

		if len(eventData) > 3 {
			subjectId, err := strconv.ParseInt(eventData[3], 10, 0)
			assert.Nil(err)
			eventPayload.SubjectId = &subjectId
		}

		assert.Equal(payload, &eventPayload)
	}
}

func TestShouldMarshalEventsCorrectly (t *testing.T) {

	assert := assert.New(t)

	assert.Equal(123, 123, "they should be equal")
}