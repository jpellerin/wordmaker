package wordmaker

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func Lex(name, input string) (*lexer, chan item) {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	go l.run()
	return l, l.items
}

var eof rune // intentionally nil

type stateFn func(*lexer) stateFn

type lexer struct {
	name       string
	input      string
	start      int
	pos        int
	width      int
	parenDepth int
	items      chan item
}

type item struct {
	typ itemType
	val string
}

type itemType int

const (
	itemError itemType = iota

	itemClass
	itemColon
	itemSlash
	itemChoice
	itemLeftParen
	itemRightParen
	itemPattern
	itemEOL
	itemEOF
)

func (i item) String() string {
	switch i.typ {
	case itemError:
		return i.val
	case itemEOF:
		return "EOF"
	}
	return fmt.Sprintf("%q", i.val)
}

func (l *lexer) run() {
	for state := lexLine; state != nil; {
		state = state(l)
	}
	defer close(l.items)
}

func (l *lexer) emit(t itemType) {
	// fmt.Printf("Emit %q %q@%v\n", t, l.input[l.start:l.pos], l.start)
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func lexLine(l *lexer) stateFn {
	r := l.peek()
	if isClass(r) {
		return lexClass
	} else if r == 'r' {
		return lexPattern
	} else if r == 'n' {
		return lexNumber
	}
	l.emit(itemError)
	return nil
}

func lexClass(l *lexer) stateFn {
	l.next()
	if l.peek() == ':' {
		l.emit(itemClass)
		l.next()
		l.emit(itemColon)
		return lexInChoices
	}
	l.emit(itemError)
	return nil
}

func lexPattern(l *lexer) stateFn {
	l.next()
	if l.peek() == ':' {
		l.emit(itemPattern)
		l.next()
		l.emit(itemColon)
		return lexInPattern
	}
	l.emit(itemError)
	return nil
}

func lexInChoices(l *lexer) (next stateFn) {
	return choiceLexer(lexChoice)(l)
}

func choiceLexer(choiceFn stateFn) stateFn {
	return func(l *lexer) stateFn {
		switch r := l.next(); {
		case r == eof || isEOL(r):
			return nil
		case r == '/':
			l.emit(itemSlash)
			return choiceFn
		case r == '(':
			l.emit(itemLeftParen)
			l.parenDepth++
			return choiceFn
		case r == ')':
			l.emit(itemRightParen)
			l.parenDepth--
			if l.parenDepth < 0 {
				l.emit(itemError)
				return nil
			}
			// blanks are not allowed after closing parens
			return choiceLexer(choiceFn)
		default:
			l.backup()
			return choiceFn
		}
	}
}

func lexChoice(l *lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case isAlpha(r):
			// absorb
			// fmt.Print(" absorb\n")
		default:
			// fmt.Print("Emit due to %v\n", r)
			l.backup()
			l.emit(itemChoice)
			break Loop
		}
	}
	return lexInChoices
}

func lexInPattern(l *lexer) stateFn {
	return choiceLexer(lexPatternItem)(l)
}

func lexPatternItem(l *lexer) stateFn {
	emit := true
Loop:
	for {
		switch r := l.next(); {
		case isDelim(r):
			if r == '(' {
				// don't emit a blank item between ((
				emit = false
			}
			l.backup()
			break Loop
		case r == eof:
			break Loop
		}
	}
	if emit || l.start < l.pos {
		l.emit(itemChoice)
	}
	return lexInPattern
}

func lexNumber(l *lexer) stateFn {
	return nil
}

func isEOL(r rune) bool {
	return r == '\r' || r == '\n'
}

func isAlpha(r rune) bool {
	return unicode.IsLetter(r)
}

func isClass(r rune) bool {
	return isAlpha(r) && unicode.IsUpper(r)
}

func isDelim(r rune) bool {
	return (r == '/' || r == '(' || r == ')')
}
