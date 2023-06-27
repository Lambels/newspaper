package interpreter

import (
	"bytes"
	"io"
	"log"
	"strconv"
)

var token_funcs = []struct {
	prefix []byte
	delims [][]byte
	gen    func(*Lexer) Element
}{
	{[]byte("- []"), nil, (*Lexer).todo},
	{[]byte("- [x]"), nil, (*Lexer).ctodo},
	{[]byte("- "), nil, (*Lexer).list},
	{[]byte("->[]"), [][]byte{[]byte("->[]"), []byte("->[x]")}, (*Lexer).chain},
	{[]byte("->[x]"), [][]byte{[]byte("->[]"), []byte("->[x]")}, (*Lexer).cchain},
	{[]byte("~"), nil, (*Lexer).text},
}

func NewLexer(buf []byte) *Lexer {
	return &Lexer{
		start: 0,
		end:   0,
		line:  0,
		err:   nil,
		buf:   buf,
	}
}

type Token struct {
	Marker  int
	Element Element
	line    int
	startx  int
	endx    int
}

type Lexer struct {
	start int
	end   int
	line  int

	err error
	buf []byte
}

// Possible Tokens:
// Text:
// This is example of text. (4)
// This is also example of text, it is terminated by \n.
//
// Todo:
// - [] This is an example of a todo.
// - [] Same here. (3)
//
// List:
// - This is an example of a list. (3)
// - Same here.
//
// Chain:
// ->[] This is the first element (2) ->[] This is the second element ->[] This is the third element.
//
// Script:
// ~This is a script~
// ~THis is not a script~ because it is contained in text ~this is a script~
func (l *Lexer) Next() (*Token, error) {
	if l.err != nil {
		return &Token{}, l.err
	}

	gen := (*Lexer).text
	var delims [][]byte
	for _, t := range token_funcs {
		if bytes.HasPrefix(l.buf[l.end:], t.prefix) {
			gen = t.gen
			delims = t.delims
		}
	}

	token := &Token{
		Marker: -1,
	}

	// move the end pointer to the first delimeter.
	l.moveBefore(delims...)
	switch {
	case l.peek() == '~' && l.buf[l.start] == '~': // fastpath: we are parsing a script, no need to backtrack.
	default: // we need to backtrack the script, it will be parsed on the next call to next.
		l.backtrackScript()
		token.Marker = l.parseMarker()
	}

	if l.peek() == '\n' {
		l.line++
	}

	// the end pointer is now in position, generate the element.
	elem := gen(l)
	token.startx = l.start
	token.endx = l.end
	token.line = l.line
	token.Element = elem

	return token, nil
}

func (l *Lexer) script() Element {
	return nil
}

func (l *Lexer) text() Element {
	return &Text{
		buf: bytes.NewBuffer(l.buf[l.start:l.end]),
	}
}

func (l *Lexer) list() Element {
	return nil
}

func (l *Lexer) todo() Element {
	return nil
}

func (l *Lexer) ctodo() Element {
	return nil
}

func (l *Lexer) chain() Element {
	return nil
}

func (l *Lexer) cchain() Element {
	return nil
}

func (l *Lexer) parseMarker() int {
	// fastpath: no closing braket, no marker
	if l.peek() != ')' {
		return -1
	}

	index := bytes.LastIndexByte(l.buf[l.start:l.end], '(')
	if index == -1 {
		return -1
	}

	marker, err := strconv.Atoi(string(l.buf[index:l.end]))
	if err != nil {
		return -1
	}
	return marker
}

func (l *Lexer) moveBefore(delims ...[]byte) {
	index := bytes.IndexByte(l.buf[l.start:], '\n')
	if index == -1 {
		index = len(l.buf[l.start:])
	}

	for _, delim := range delims {
		if v := bytes.Index(l.buf[l.start:index], delim); v != -1 && v < index {
			index = v
		}
	}

	l.end = index + l.start
}

func (l *Lexer) backtrackScript() {
	// fastpath1: no script tag imediately before end.
	if l.peek() != '~' {
		return
	}

	// push the end pointer back to the first script tag from left to right.
	index := bytes.LastIndexByte(l.buf[l.start:l.end], '~')
	if index != -1 {
		l.end = index
	}
}

func (l *Lexer) peek() byte {
	if l.atEnd() {
		return '\x00'
	}

    log.Println(l.end)
	return l.buf[l.end]
}

func (l *Lexer) atEnd() bool {
    if l.end >= len(l.buf) {
		l.err = io.EOF
		return true
	}

	return false
}
