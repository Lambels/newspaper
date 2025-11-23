package pusher

import (
	"encoding"
)

const (
	PushInCode byte = iota
	PushCode
	PushNCode
	EveryCode
)

type Pusher interface {
	Advance() error
	AdvanceN(int) (int, error)
	encoding.BinaryMarshaler
}
