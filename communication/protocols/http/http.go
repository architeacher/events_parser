package http

import (
	"net/http"
	"bytes"
	"splash/communication"
)

type Protocol struct {
	client *http.Client
}

func NewProtocol(client *http.Client) *Protocol {
	return &Protocol{client}
}

func (self *Protocol) Client() *http.Client {
	return self.client
}

func (self *Protocol) Send(req *communication.Request) (*communication.Response, error){

	httpRequest, err := NewHttpRequest(req)

	if err != nil {
		return nil, err
	}

	nativeRequest, err := http.NewRequest(
		httpRequest.Method(),
		httpRequest.Protocol() + httpRequest.Host() + httpRequest.path,
		bytes.NewBuffer(httpRequest.Body()),
	)

	// Setting headers
	for key, value := range req.Meta() {

		nativeRequest.Header.Add(key, value)
	}

	if err != nil {
		return nil, err
	}

	nativeResponse, err := self.Client().Do(nativeRequest)

	if err != nil {
		return nil, err
	}

	httpResponse, err := NewHttpResponseFromNative(nativeResponse)

	if err != nil {
		return nil, err
	}

	return httpResponse.toBaseResponse(), nil
}