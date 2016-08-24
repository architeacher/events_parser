package http

import (
	"bytes"
	"net/http"
	"splash/communication"
)

type Protocol struct {
	client *http.Transport
}

func NewProtocol(client *http.Transport) *Protocol {
	return &Protocol{client}
}

func (self *Protocol) Client() *http.Transport {
	return self.client
}

func (self *Protocol) Send(req *communication.Request) (*communication.Response, error) {

	httpRequest, err := NewHttpRequest(req)

	if err != nil {
		return nil, err
	}

	nativeRequest, err := http.NewRequest(
		httpRequest.Method(),
		httpRequest.Protocol()+httpRequest.Host()+httpRequest.path,
		bytes.NewBuffer(httpRequest.Body()),
	)

	if err != nil {
		return nil, err
	}

	// Setting headers
	for key, value := range req.Meta() {
		nativeRequest.Header.Add(key, value)
	}

	nativeResponse, err := self.Client().RoundTrip(nativeRequest)

	if err != nil {
		return nil, err
	}

	httpResponse, err := NewHttpResponseFromNative(nativeResponse)

	if err != nil {
		return nil, err
	}

	return httpResponse.toBaseResponse(), nil
}
