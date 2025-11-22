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

// join is used to join multiple pushers to chain format. This is usually how internal state is represented for each pusher
//
// reps1 -> reps2 --> ... -> repsn
func join(reps ...Pusher) []byte {
	return ""
}

// joinRec is used to join recursive definitions without producing an infinite byte array (that would require allot of memory!)
// 
// to use joinRec put your recursive identifier as the first argument and the current state under reps variadic arguments.
// the following will be produced:
// 	- if reps in empty/produces empty string: rec
// 	- else: reps1 -> reps2 -> reps3 -> ... -> rec ( -> reps1 -> reps2 -> ... -> rec -> ...)
// the part in () isnt actually represented but is implied by the rec.
func joinRec(rec string, reps ...Pusher) []byte {
	// makes use of join
	return ""
}

// join joins the string representation of multiple pushers into a chain representation.
//func join(reps ...string) string {
//	return ""
//}
