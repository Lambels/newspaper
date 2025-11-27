package pusher

import (
	"encoding/binary"
	"errors"
)

type Repeat struct {
	p Pusher
	f factory
	n int // the number of repetitions left.
}

func NewRepeat(n int, f factory) Repeat {
	return Repeat{f(), f, n-1}
}

func (p *Repeat) refresh() error {
	if p.n < 1 {
		return ErrFinished
	}

	p.n -= 1
	p.p = p.f()
	return ErrPushed
}

func (p *Repeat) Advance() error {
	err := p.p.Advance()

	if errors.Is(err, ErrFinished) {
		return p.refresh()
	}

	return err
}

func (p *Repeat) AdvanceN(n int) (int, error) {
	n, err := p.p.AdvanceN(n)

	if errors.Is(err, ErrFinished) {
		return n, p.refresh()
	}

	return n, err
}

func (p Repeat) AppendBinary(buf []byte) ([]byte, error) {
	if p.n < 1 {
		return buf, ErrFinished
	}

	var err error
	// we first write the current pusher.
	buf, err = p.p.AppendBinary(buf)
	// static pushers have no reason to fail but we still check for consistency issues (only possible error 
	// should be ErrFinished)
	if err != nil && !errors.Is(err, ErrFinished) {
		return nil, err
	}

	buf = append(buf, RepeatCode) // start of repeat sequence
	buf = binary.AppendUvarint(buf, uint64(p.n))

	buf, err = p.f().AppendBinary(buf)
	if err != nil {
		return nil, err
	}

	buf = append(buf, EndSeqCode) // end of repeat sequence
	return buf, nil
}

type Forever struct {
	p Pusher
	f factory
}


func NewForever(f factory) Forever {
}
