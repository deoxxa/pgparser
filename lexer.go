package pgparser

import (
	"strings"

	"go.bmatsuo.co/go-lexer"
)

const (
	itemLeftBrace lexer.ItemType = iota
	itemRightBrace
	itemLeftParen
	itemRightParen
	itemComma
	itemQuotedString
	itemBareString
)

func stateBegin(l *lexer.Lexer) lexer.StateFn {
	r, n := l.Advance()

	switch r {
	case '{':
		l.Emit(itemLeftBrace)
		return stateBegin
	case '}':
		l.Emit(itemRightBrace)
		return stateBegin
	case '(':
		l.Emit(itemLeftParen)
		return stateBegin
	case ')':
		l.Emit(itemRightParen)
		return stateBegin
	case ',':
		l.Emit(itemComma)
		return stateBegin
	case '"':
		l.Backup()
		return stateQuotedString
	default:
		if n > 0 {
			l.Backup()
			return stateBareString
		}
	}

	return nil
}

func stateQuotedString(l *lexer.Lexer) lexer.StateFn {
	if !l.Accept("\"") {
		l.Errorf("expected an opening quote")
	}

	for {
		if l.Accept("\\") {
			l.Advance()
		}

		n := l.AcceptRunFunc(func(r rune) bool {
			return r != '\\' && r != '"'
		})

		if n == 0 {
			break
		}
	}

	if !l.Accept("\"") {
		l.Errorf("expected a closing quote")
	}

	l.Emit(itemQuotedString)

	return stateBegin
}

func stateBareString(l *lexer.Lexer) lexer.StateFn {
	l.AcceptRunFunc(func(r rune) bool {
		return !strings.ContainsRune(`{}(),"`, r)
	})

	l.Emit(itemBareString)

	return stateBegin
}
