package pusher

type Push struct {
	called bool
}

func NewPush() Push {
	return Push{}
}

func (p *Push) Advance() error {
	p.called = true
	return ErrFinished
}

func (p *Push) AdvanceN(_ int) (int, error) {
	p.called = true
	return 1, ErrFinished
}

func (p Push) MarshalBinary() ([]byte, error) {
	if !p.called {
		return []byte{PushCode}, nil
	}

	return nil, ErrFinished
}
