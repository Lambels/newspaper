package interpreter

import (
	"bytes"
	"io"
)

const whitespace string = " \t\v"

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
	//TODO: maybe need to trim start here.
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
	// start represents the current start of the token in the file buf.
	start int
	// end represents the current end of the token in token buf.
	end int
	// line represents the current line.
	line int
	// token holds the bytes of an actual token, marker and script.
	token []byte

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
    //TODO: check for newline character and increment line, only once.
    if l.peek() == '\n' {
        l.line++
        l.start++
    }

	if l.err != nil {
		return &Token{}, l.err
	}

	gen := (*Lexer).text
	switch len(l.token) {
	case 0: // we need to get a token.
		var delims [][]byte
		for _, t := range token_funcs {
			if hasPrefixWithSpace(l.buf[l.start:], t.prefix) {
				gen = t.gen
				delims = t.delims
			}
		}

		index := indexFirst(l.buf[l.start:], delims...)
        l.setToken(l.buf[l.start:index])

        // now that we have a token set, we need to move the end offset before the script.
        l.end = len(l.token)
        if l.token[l.end-1] == '~' {
            if v := bytes.LastIndexByte(l.token[:l.end-1], '~'); v != -1 {
                l.end = v
            }
        }

	default: // fastpath: if there is already a token set, it must be a script, truncate whitespace and generate it.
		gen = (*Lexer).script
        l.setToken(bytes.TrimLeft(l.token, whitespace))
	}

    elem := gen(l)
    if elem == nil {
        return nil, nil
    }

    return nil, nil
}

func (l *Lexer) setToken(token []byte) {
	if len(token) == 0 {
		l.start += l.end
		l.token = nil
		return
	}

	if l.token != nil {
		l.start += len(l.token) - len(token)
	}

	l.token = token
}

func (l *Lexer) script() Element {
	return nil
}

func (l *Lexer) text() Element {
    clean := bytes.TrimLeft(l.token, whitespace)
    l.setToken(clean)
    if len(clean) == 0 {
        return nil
    }

    return &Text{
        buf: bytes.NewBuffer(bytes.TrimSpace(l.token[:l.end])), 
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

func indexFirst(buf []byte, delims ...[]byte) int {
	index := bytes.IndexByte(buf, '\n')
	if index == -1 {
		index = len(buf)
	}

	for _, delim := range delims {
		var offset int
		if hasPrefixWithSpace(buf, delim) {
			offset = bytes.Index(buf, delim)
		}

		if v := bytes.Index(buf[offset:index], delim); v != -1 && v+offset < index {
			index = v + offset
		}
	}

	return index
}

func (l *Lexer) peek() byte {
	if l.atEnd() {
		return '\x00'
	}

	return l.buf[l.end]
}

func (l *Lexer) atEnd() bool {
	if l.end >= len(l.buf) {
		l.err = io.EOF
		return true
	}

	return false
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

func leadingTab(s []byte) int {
	var count int
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c != '\t' {
			return count
		}
		count++
	}

	return count
}

func hasPrefixWithSpace(s []byte, prefix []byte) bool {
	return bytes.HasPrefix(bytes.TrimLeft(s, " \t\v"), prefix)
}
