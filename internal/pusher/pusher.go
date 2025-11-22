package pusher

import "encoding"

type Pusher interface {
	// Advance represents 
	Advance() error
	AdvanceN(int) (int, error)
	encoding.BinaryMarshaler
	encoding.BinaryAppender
}
