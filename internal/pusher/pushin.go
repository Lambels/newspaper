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
	p.n -= 1

	if p.n < 1 {
		return ErrFinished
	}

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

func (p PushIn) MarshalBinary() ([]byte, error) {
	if p.n < 1 {
		return nil, ErrFinished	
	}

	buf := make([]byte, 9)
	buf[0] = PushInCode
	n := binary.PutUvarint(buf[1:], uint64(p.n))
	return buf[:n+1], nil
}
