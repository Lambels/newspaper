package newspaper

import (
	"context"
	"os"
)

// Script models the current state of a script. It describes how the script should advance and how many hops
// the script has left till it evaluates to its action.
//
// Each script must evaluate to an action in the furture and must have have a string representation regardless of its state.
type PreCondition interface {
    // Evaluated indicates the action to which the pre condition will evaluate to after TTL passes.
	Evaluated() Action

    // String returns the string representation of the pre condition.
	String() string

    // TTL indicates the time to live for the pre condition before the action occurs.
	TTL() int

    // Advance advances the state of the pre condition by n.
	Advance(n int) error
}

// Action represents a PreCondition which takes places now and needs to be ran.
type Action interface {
    PreCondition

    // Runs the action.
	Run(ctx context.Context, next *os.File, parent Element) error
}
