package pusher

import (
	"errors"
	"fmt"
)

// The following errors are returned by Advance methods.

// Represents the pusher pushing: either push or push_in 1 are evaluated.
var ErrPushed = errors.New("Pusher pushed")
// Represents the pusher finishing, also happens after a push or push_in 1, but signals that there will be no more pushes
// afterwards. When a pusher finishes, it always must finish on a push, hence we wrap the ErrPushed.
var ErrFinished = fmt.Errorf("Pusher finished and implicitly %w", ErrPushed)
