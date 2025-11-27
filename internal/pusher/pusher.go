package pusher

import (
	"encoding"
)

const (
	// These are the primitive atoms in which all the Packed types should break into. A packed type
	// should wrap one or multiple of these primitive atoms to provide extra functionality.

	// The code for the PushIn atom expects an unsigned integer as per UVarint encoding in encoding/binary.
	PushInCode byte = iota
	// The code for the push atom.
	PushCode
	// The code for the PushN atom expects an unsigned integer as per UVarint encoding in encoding/binary.
	PushNCode

	// These are the Packed types. They provide extra functionality over the standard and limited primitive atoms.

	// The code for the Every pack expects an unsigned integer as per UVarint encoding in encoding/binary.
	EveryCode
	// The code for the Repeat pack expects an EndSeqCode and .an unsigned integer as per UVarint encoding in encoding/binary.  
	RepeatCode
	// The code for the Forever pack expects an EndSeqCode.
	ForeverCode

	// EndSeqCode represents the end of the arguments of a multi argument code. It can be nested.
	EndSeqCode
)

type Pusher interface {
	Advance() error
	AdvanceN(int) (int, error)
	encoding.BinaryAppender
}

type factory func() Pusher
