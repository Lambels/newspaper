package interpreter

import (
	"bytes"
	"io"
	"strconv"
	"unicode"
)

const whitespace string = " \t\v"

var token_funcs = []struct {
	prefix []byte
	gen    func(*Lexer) Element
}{
	{[]byte("- []"), (*Lexer).todo},
	{[]byte("- [x]"), (*Lexer).ctodo},
	{[]byte("- "), (*Lexer).list},
	{[]byte("->[]"), (*Lexer).chain},
	{[]byte("->[x]"), (*Lexer).cchain},
}

func NewLexer(buf []byte) *Lexer {
	return &Lexer{
		start: 0,
		end:   0,
		line:  0,
		buf:   buf,
	}
}

type Token struct {
	Marker  int
	Element Element
	Line    int
	Start   int
	End     int
}

type Lexer struct {
	// start represents the current start of the token in the file buf.
	start int
	// end represents the current end of the token in token buf.
	end int
	// line represents the current line.
	line   int
	indent int

	tkn *Token

	reporter *ErrorReporter
	buf      []byte
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
func (l *Lexer) Advance() bool {
	if l.atEnd() {
		return false
	}

	var lexeme Element
	marker := -1
	c := l.advance()

	switch {
	case c == '~': // script tag, emit script lexeme.
		lexeme = l.script()
	case c == '\n': // new line, reset indent counter and increment line.
		l.indent = 0
		l.line++
		// for each consecutive match of a tab increase the indent and consume it.
		for l.match('\t') {
			l.indent++
		}
	case unicode.IsSpace(rune(c)): // skip out of order space.
	default: // we need to generate a more complex lexeme.
		gen := (*Lexer).text
		for _, v := range token_funcs {
			if bytes.HasPrefix(l.buf[l.start:], v.prefix) {
				gen = v.gen
			}
		}
		lexeme = gen(l)
		marker = l.parseMarker()
	}

	// if parsing the current lexeme had an error, report an illegal advance.
	if l.reporter.HadError() {
		return false
	}

	l.tkn = &Token{
		Marker:  marker,
		Element: lexeme,
		Line:    l.line,
		Start:   l.start,
		End:     l.end,
	}

	return true
}

func (l *Lexer) Token() *Token {
    return l.tkn
}

func indexBeforeClosing(delims ...[]byte) int {

}

func (l *Lexer) parseMarker() int {
	// try to parse marker.
	var end int
	for i := l.end; i > l.start; i-- {
		if !unicode.IsSpace(rune(l.buf[i])) && l.buf[i] != ')' {
			return -1
		} else if l.buf[i] == ')' {
			end = i
			break
		}
	}

	index := bytes.LastIndexByte(l.buf[l.start:], '(')
	if index == -1 {
		return -1
	}

	marker, err := strconv.Atoi(string(l.buf[index+1 : end]))
	if marker < 0 && err == nil {
		//TODO: report error.
		return -1
	} else if err != nil {
		return -1
	}

	return marker
}

func (l *Lexer) advance() byte {
	if l.end >= len(l.buf)-1 {
		return '\x00'
	}

	l.end++
	return l.buf[l.end-1]
}

func (l *Lexer) peek() byte {
	if l.atEnd() {
		return '\x00'
	}

	return l.buf[l.end]
}

func (l *Lexer) atEnd() bool {
	if l.end >= len(l.buf) {
		return true
	}

	return false
}

func (l *Lexer) match(s byte) bool {
	if l.peek() != s {
		return false
	}

	l.advance()
	return true
}

func (l *Lexer) script() Element {
	return nil
}

func (l *Lexer) text() Element {
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
