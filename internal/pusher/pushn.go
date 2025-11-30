package pusher

import "encoding/binary"

// could be a chain of push n times but due to efficiency we will provide native support.
type PushN struct {
	n int
}

func NewPushN(n int) *PushN {
	return &PushN{n}
}

func (p *PushN) Advance() error {
	if p.n < 1 {
		return ErrFinished
	}

	p.n--
	return ErrPushed
}

func (p *PushN) AdvanceN(_ int) (int, error) {
	return 1, p.Advance()
}

func (p PushN) AppendBinary(buf []byte) ([]byte, error) {
	if p.n < 1 {
		return buf, ErrFinished
	}

	buf = append(buf, PushNCode)
	return binary.AppendUvarint(buf, uint64(p.n)), nil
}
