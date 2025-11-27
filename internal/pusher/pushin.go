package pusher

import (
	"encoding/binary"
	"errors"
)

type PushIn struct {
	n int
}

func NewPushIn(n int) (PushIn, error) {
	if n <= 0 {
		return PushIn{}, errors.New("Invalid value for n, n has to be > 0")
	}

	return PushIn{n}, nil
}

func (p *PushIn) Advance() error {
	if p.n < 1 {
		return ErrFinished
	}

	p.n -= 1

	return nil
}

func (p *PushIn) AdvanceN(n int) (int, error) {
	p.n -= n
	if p.n < 1 {
		// return p.n value if overflowed, this value makes
		// sense up to first ErrFinished returned afterwards
		// the value no longer makes sense.
		return p.n + n, ErrFinished
	}

	return n, nil
}

func (p PushIn) AppendBinary(buf []byte) ([]byte, error) {
	if p.n < 1 {
		return buf, ErrFinished	
	}

	buf = append(buf, PushInCode)
	buf = binary.AppendUvarint(buf, uint64(p.n))
	return buf, nil
}
