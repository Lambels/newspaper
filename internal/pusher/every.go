package pusher

import (
	"encoding/binary"
	"errors"
)

// Could be implemented using a Forverer(PushIn(n))
type Every struct {
	diff PushIn
	n    int
}

func NewEvery(n int) (*Every, error) {
	if n <= 0 {
		return nil, errors.New("Invalid value for n, n has to be > 0")
	}

	return &Every{PushIn{n}, n}, nil
}

func NewDelayedEvery(diff int, cycle int) (*Chain, error) {
	if diff < 0 || cycle <= 0 {
		return nil, errors.New("Invalid arguments provided, diff > 0 and cycle >= 0")
	}

	return &Chain{&PushIn{diff}, &Every{PushIn{cycle}, cycle}}, nil
}

func (p *Every) Advance() error {
	if err := p.diff.Advance(); errors.Is(err, ErrFinished) {
		p.diff = PushIn{p.n}
		return ErrPushed
	}

	return nil
}

func (p *Every) AdvanceN(n int) (int, error) {
	if n, err := p.diff.AdvanceN(n); errors.Is(err, ErrFinished) {
		p.diff = PushIn{p.n}
		return n, ErrPushed
	}

	return n, nil
}

func (p Every) AppendBinary(buf []byte) ([]byte, error) {
	if p.n != p.diff.n {
		buf, _ = p.diff.AppendBinary(buf)
	}

	buf = append(buf, EveryCode)
	buf = binary.AppendUvarint(buf, uint64(p.n))
	return buf, nil
}
