package logger

type FormatterInterface interface {
	format(message interface{}, context ...Context) string
}

type Formatter struct {
	FormatterInterface
}

func NewFormatter() *Formatter {
	return &Formatter{}
}

func (self *Formatter) format(context ...Context) interface{} {
	return context
}
