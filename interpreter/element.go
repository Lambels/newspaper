package interpreter

import "io"

type Element interface {
	String() string
	io.Reader
}

type EvalElement interface {
	Element
	Evaluate() bool
}

type IndexElement interface {
	Element
	io.Writer
	io.WriterTo
}

// Some Text:
//      - dsafgsgs
//      - Idk


