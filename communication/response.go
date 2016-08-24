package communication

type Response struct {
	body         string
	meta         map[string]string
	isSuccessful bool
}

func NewResponse(body string, meta map[string]string, isSuccessful bool) *Response {
	return &Response{body, meta, isSuccessful}
}

func (self *Response) Body() string {
	return self.body
}

func (self *Response) Meta() map[string]string {
	return self.meta
}

func (self *Response) IsSuccessful() bool {
	return self.isSuccessful
}
