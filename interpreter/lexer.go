package interpreter

import (
	"bytes"
	"strconv"
	"unicode"
)

// token_func maps each of the tokens prefix to their respective generator function.
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

// A token represents a literal wrapped by the type Element with aditional information such as the location and marker.
type Token struct {
	Marker  int
	Element Element
	Line    int
	Start   int
	End     int
}

// NewLexer creates a new lexer which scans over buf for tokens.
func NewLexer(buf []byte, reporter *ErrorReporter) *Lexer {
	return &Lexer{
		start:    0,
		end:      0,
		line:     0,
		buf:      buf,
		reporter: reporter,
	}
}

type Lexer struct {
	// start represents the current start of the token in the file buf.
	start int
	// end represents the current end of the token in token buf.
	end int
	// line represents the current line.
	line   int
	indent int

	// parsedScript is a flag which signals wether the last element parsed was a script to prevent illegal advances.
	parsedScript bool

	// token represents the last token parsed.
	tkn *Token

	reporter *ErrorReporter
	buf      []byte
}

// Advance advances the next token in buf. It also reports wether the advance was legal and similarlly reports any errors
// to the provided error reporter.
//
// The possible tokens are:
// Script: ~SCRIPT~
// Text: Some text
// List: - Some list
// Todo: - [] Some todo
// Completed Todo: - [x] Some completed todo
// Chain: ->[] Some chain
// Completed chain: ->[x] Some completed chain
func (l *Lexer) Advance() bool {
	if l.atEnd() {
		return false
	}

	var lexeme Element
	marker := -1

Outer:
	for {
		c := l.advance()

		switch {
		case c == '\x00':
			return false
		case c == '~': // script tag, emit script lexeme.
			l.parsedScript = true
			lexeme = l.script()
			break Outer
		case c == '\n': // new line, reset indent counter and increment line.
			l.indent = 0
			l.start++
			l.line++
			l.parsedScript = false
			// for each consecutive match of a tab increase the indent and consume it.
			for l.match('\t') {
				l.indent++
			}

			for once := true; once || l.matchSpace(); once = false {
				if l.peek() == '\n' {
					l.tkn = nil
					l.start = l.end
					return true
				}
			}
		case unicode.IsSpace(rune(c)): // skip out of order space.
			for once := true; once || l.matchSpace(); once = false {
				if l.peek() == '\n' {
					l.tkn = nil
					l.start = l.end
					return true
				}
			}
		default: // we need to generate a more complex lexeme.
			gen := (*Lexer).text
			for _, v := range token_funcs {
				if bytes.HasPrefix(l.buf[l.end-1:], v.prefix) {
					gen = v.gen
				}
			}
			lexeme = gen(l)
			marker = l.parseMarker()
			break Outer
		}
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
	l.start = l.end

	return true
}

// Token returns the previously parsed token.
//
// Token also reports the last parsed token even after multiple illegal advancements.
func (l *Lexer) Token() *Token {
	return l.tkn
}

// TODO: figure this out.
func indexBeforeClosing(buf []byte, delims ...[]byte) int {
	// everything is ended by either a newline or a script tag.
	var upperBound int
	if i := bytes.IndexByte(buf, '\n'); i != -1 {
		upperBound = i
	} else {
		upperBound = len(buf) - 1
	}

	// check if there is a script tag before the newline.
	if i := bytes.IndexByte(buf[:upperBound], '~'); i != -1 && i < upperBound {
		upperBound = i
	}

	// now check for the delims.
	for _, delim := range delims {
		if i := bytes.Index(buf[:upperBound], delim); i != -1 && i < upperBound {
			upperBound = i
		}
	}

	return upperBound
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

	index := bytes.LastIndexByte(l.buf[l.start:l.end], '(')
	if index == -1 {
		return -1
	} else {
		index += l.start
	}

	marker, err := strconv.Atoi(string(l.buf[index+1 : end]))
	if marker < 0 && err == nil {
		l.reporter.AddWarning(nil)
		return -1
	} else if err != nil {
		return -1
	}

	return marker
}

// advance advances the end pointer and reports the byte before advancing, this is to allow
// first iteration advancing.
//
// it returns a null terminator byte if currently at end.
func (l *Lexer) advance() byte {
	if l.atEnd() {
		return '\x00'
	}

	l.end++
	return l.buf[l.end-1]
}

// peek reports the current byte under the end pointer.
//
// it returns a null terminator byte if currently at end.
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

// match conditionally consumes s if the current character matches it and reports the advancment.
func (l *Lexer) match(s byte) bool {
	if l.peek() != s {
		return false
	}

	l.advance()
	return true
}

// matchSpace conditionally consumes a space character (ignoring \n) and reports the advancement.
func (l *Lexer) matchSpace() bool {
	c := l.peek()
	if c == ' ' || c == '\r' {
		l.advance()
		return true
	}

	return false
}

// script parses a script starting from the first character of the script.
//
// It also consumes the closing script tag and reports:
//
// - Any errors while parsing the script.
//
// - Any unterminated scripts.
func (l *Lexer) script() Element {
	// illegal to parse 2 scripts one after another.
	if l.parsedScript {
		l.reporter.AddError(nil)
		return nil
	}

	for l.peek() != '~' && !l.atEnd() {
		if l.peek() == '\n' {
			l.reporter.AddError(nil)
			return nil
		}
		l.advance()
	}

	if l.atEnd() {
		l.reporter.AddError(nil)
		return nil
	}

	// consume closing "~"
	l.advance()

	script, err := parseScript(l.buf[l.start:l.end])
	if err != nil {
		l.reporter.AddError(err)
		return nil
	}

	return script
}

func (l *Lexer) text() Element {
	// illegal to parse text after script on the same line.
	if l.parsedScript {
		l.reporter.AddError(nil)
		return nil
	}

	l.end += indexBeforeClosing(l.buf[l.end:])

	buf := bytes.NewBuffer(l.buf[l.start : l.end+1])
	return &Text{
		buf: buf,
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
