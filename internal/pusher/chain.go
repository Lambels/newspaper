package pusher

import "errors"

// A chain is a pusher made of multiple pushers, this is the state most pushers will be in since they all break down
// into chains.
type Chain []Pusher

func (c *Chain) Advance() error {
	if len(*c) == 0 {
		return ErrFinished
	}

	err := (*c)[0].Advance()
	// specifically check for ErrPushed and not use error wrapping.
	if err == nil || err == ErrPushed {
		return err
	}

	if errors.Is(err, ErrFinished) {
		*c = (*c)[1:]
		return ErrPushed
	}

	return err
}

func (c *Chain) AdvanceN(n int) (int, error) {
	if len(*c) == 0 {
		return 0, ErrFinished
	}

	n, err := (*c)[0].AdvanceN(n)
	if err == nil || err == ErrPushed {
		return n, err
	}

	if errors.Is(err, ErrFinished) {
		*c = (*c)[1:]
		return n, ErrPushed
	}

	return n, err
}

// No need for a start sequence since eveything is implicitly a chain under the hood.
func (p Chain) AppendBinary(buf []byte) ([]byte, error) {
	var err error
	for _, push := range p {
		buf, err = push.AppendBinary(buf)
		if err != nil && !errors.Is(ErrFinished, err) {
			// error not related to pusher state, cannot marshal chain.
			return nil, err
		}
	}

	buf = append(buf, EndSeqCode)
	return buf, nil
}

// Add end sequence and add start sequence also.
// Eagerly execute all pushers.
// return err finished if no pushers left.
// return err pushed if all pushed.
type AndChain []Pusher

func (c *AndChain) Advance() error {
	if len(*c) == 0 {
		return ErrFinished
	}

	var shifted int
	rerr := ErrPushed
	for i, pusher := range *c {
		err := pusher.Advance()
		if !errors.Is(err, ErrPushed) {
			rerr = nil
		}

		if errors.Is(err, ErrFinished) {
			removeAndReslice((*[]Pusher)(c), i-shifted)
			if len(*c) == 0 {
				return ErrFinished
			}
		}
	}

	return rerr
}

// TODO: implement somehow, choose good representation
func (c *AndChain) AdvanceN(n int) (int, error) {
	if len(*c) == 0 {
		return 0, ErrFinished
	}

}

type OrChain []Pusher

func removeAndReslice(ps *[]Pusher, i int) {
	
}
