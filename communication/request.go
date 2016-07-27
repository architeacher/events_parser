package communication

type Request struct {
	body         []byte
	meta         map[string]string
}

func NewRequest(body []byte, meta map[string]string) *Request {
	return &Request{body, meta}
}

func (self *Request) Body() []byte {
	return self.body
}

func (self *Request) Meta() map[string]string {
	return self.meta
}