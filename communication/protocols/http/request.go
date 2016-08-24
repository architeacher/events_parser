package http

import (
	"errors"
	"splash/communication"
)

type Request struct {
	headers     map[string]string
	method      string
	protocol    string
	host        string
	path        string
	body        []byte
	baseRequest *communication.Request
}

func NewHttpRequest(req *communication.Request) (*Request, error) {
	var method, protocol, host, path string

	for key, value := range req.Meta() {
		switch key {
		case "method":
			method = value
			delete(req.Meta(), key)
			break
		case "protocol":
			protocol = value
			delete(req.Meta(), key)
			break
		case "host":
			host = value
			delete(req.Meta(), key)
			break
		case "path":
			path = value
			delete(req.Meta(), key)
			break
		}
	}

	if "" == method {
		return nil, errors.New("Can not build http request, \"method\" is missing.")
	}

	if "" == protocol {
		protocol = "http://"
	}

	if "" == host {
		return nil, errors.New("Can not build http request, \"host\" is missing.")
	}

	if "" == path {
		path = "/"
	}

	return &Request{
		req.Meta(),
		method,
		protocol,
		host,
		path,
		req.Body(),
		req,
	}, nil
}

func (self *Request) Headers() map[string]string {
	return self.headers
}

func (self *Request) Method() string {
	return self.method
}

func (self *Request) Protocol() string {
	return self.protocol
}

func (self *Request) Host() string {
	return self.host
}

func (self *Request) Path() string {
	return self.path
}

func (self *Request) Body() []byte {
	return self.body
}

func (self *Request) toBaseRequest() *communication.Request {
	return self.baseRequest
}
