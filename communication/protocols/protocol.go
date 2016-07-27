package protocols

import (
	"splash/communication"
)

type Protocol interface {
	Send(*communication.Request) (*communication.Response, error)
}
