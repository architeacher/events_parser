package client

import (
	"os"
	"encoding/csv"
	"bufio"
	"strconv"
	"github.com/golang/protobuf/proto"
	"splash/communication/protocols/protobuf"
	"splash/communication"
	httpProtocol "splash/communication/protocols/http"
)

type Client struct {
	protocol *httpProtocol.Protocol
}

func NewClient(protocol *httpProtocol.Protocol) *Client {
	return &Client{
		protocol,
	}
}

func (self *Client) LoadCSVFile(path *string) ([][]byte, error){

	file, err := os.Open(*path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	// Enforce the Reader not to check the number of fields
	reader.FieldsPerRecord = -1

	rows, err := reader.ReadAll()

	if err != nil {
		return nil, err
	}

	payloadCollection := make([]*protobuf.Event_Payload, 0)
	eventsBuffers := make([][]byte, 0)

	totalRecords, rowsSizes,  k := len(rows), 0, 0

	requiredItemsInPayloadCollection := 1

	for index, row := range rows {

		payload, err := self.BuildPayload(row)

		if nil != err {
			return nil, err
		}

		rowsSizes += self.getRowSize(row)

		k++

		payloadCollection = append(payloadCollection, payload)

		if (requiredItemsInPayloadCollection == k){

			event := new(protobuf.Event)
			event.PayloadCollection = payloadCollection

			buffer, err := proto.Marshal(event)

			if err != nil {
				return nil, err
			}

			eventsBuffers = append(eventsBuffers, buffer)
			payloadCollection = make([]*protobuf.Event_Payload, 0)

			requiredItemsInPayloadCollection++

			remainingRecords := totalRecords - index - 1

			if requiredItemsInPayloadCollection > remainingRecords{
				requiredItemsInPayloadCollection = remainingRecords
			}

			k, rowsSizes = 0, 0
		}
	}

	return eventsBuffers, nil
}

func (self *Client) BuildPayload(row []string) (*protobuf.Event_Payload, error){

	payload := &protobuf.Event_Payload{}

	for key, value := range row {

		switch key {
		case 0:

			entry, err := strconv.ParseInt(value, 10, 0)

			if nil != err {
				return nil, err
			}

			payload.Time = &entry
			break
		case 1:

			var entry protobuf.Event_EventType

			switch value {
			case "signup":
				entry = protobuf.Event_SIGNUP
				break
			case "follow":
				entry = protobuf.Event_FOLLOW
				break
			case "viorama":
				entry = protobuf.Event_SPLASH_CREATION
				break
			case "view":
				entry = protobuf.Event_IMPRESSION
				break
			}

			payload.EventType = &entry
			break
		case 2:

			entry, err := strconv.ParseInt(value, 10, 0)

			if nil != err {
				return nil, err
			}

			payload.ActorId = &entry
		case 3:

			entry, err := strconv.ParseInt(value, 10, 0)

			if nil != err {
				return nil, err
			}

			payload.SubjectId = &entry
			break
		}
	}

	return payload, nil
}

func (self *Client) getRowSize(row []string) int {

	size := 0
	for _, record := range row  {

		size += len(record)
	}

	return size
}

func (self *Client) SendData(data *[]byte, host, path *string) (*communication.Response, error) {

	request := communication.NewRequest(*data, map[string]string{
		"method": "Post",
		"protocol": "http://",
		"host": *host,
		"path": *path,
		"Content-Type": "application/x-protobuf",
	})

	protocolResponse, err := self.protocol.Send(request)

	if err != nil {
		return protocolResponse, err
	}

	return protocolResponse, nil
}
