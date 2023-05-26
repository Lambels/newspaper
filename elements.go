package newspaper

import "io"

// An element is an individual piece of text which is scriptable by the use of a pre condition.
//
// Element must implement Writer to describe how another element can be written into it and WriterTo to describe how it is written into a file or
// indexed into another element.
type Element interface {
	io.Writer
	io.WriterTo

    // PreCondition returns the PreCondition for the Element, it always returns non nil PreConditions since newspaper will only look at elements
    // with pre conditions.
	PreCondition() PreCondition
    // String represents the string representation of the Element along side any pre conditions.
	String() string
    // MorphInto downgrades or upgrades the element to the passed element.
	MorphInto(target Element) Element
    // Completed indicates wether the current element is marked as complete, if there is no complete state for the element it defaults to false.
    Completed() bool

    // Parent returns the parent element of the element.
    Parent() Element
    // Children returns the children elements of the element.
    Children() []Element
}
