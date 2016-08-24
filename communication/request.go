package communication

import "strconv"

type Request struct {
	body []byte
	meta map[string]interface{}
}

func NewRequest(body []byte, meta map[string]interface{}) *Request {
	return &Request{body, meta}
}

func (self *Request) Body() []byte {
	return self.body
}

func (self *Request) Meta() map[string]string {

	copy := make(map[string]string)

	for key, item := range self.meta {

		switch value := item.(type) {
		case string:
			copy[key] = value
		case int:
			copy[key] = strconv.Itoa(value)
		}
	}

	return copy
}
