package interpreter

import (
	"bufio"
	"io"
)

type Lexer struct {
	line int
	col  int

	reader bufio.Reader
}

type lexError struct {
    line int
    col int

    msg string
}

func NewLexer(source io.Reader) *Lexer {
	return &Lexer{
		reader: *bufio.NewReader(source),
	}
}

func (l *Lexer) Next(_) (_, error) {

}
