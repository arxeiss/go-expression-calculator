package lexer

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrUnexpectedChar = errors.New("unexpected character")
	ErrInvalidNumber  = errors.New("cannot parse number")
)

type LexerError struct {
	token    *Token
	position *int
	err      error
}

func (e *LexerError) Position() int {
	if e.token != nil {
		return e.token.startPos
	}
	if e.position != nil {
		return *e.position
	}
	return -1
}

func (e *LexerError) Unwrap() error {
	return e.err
}

func (e *LexerError) Error() string {
	b := strings.Builder{}
	pos := e.Position()

	if pos < 0 {
		b.WriteString("unexpected error")
		if e.err != nil {
			b.WriteByte(' ')
			b.WriteString(e.err.Error())
		}
		return b.String()
	}

	if e.err != nil {
		b.WriteString(e.err.Error())
	} else {
		b.WriteString("error")
	}
	b.WriteString(" at position ")
	b.WriteString(strconv.Itoa(pos))

	if e.token != nil {
		b.WriteString(" found ")
		b.WriteString(e.token.tType.String())
		b.WriteString(" token")
	}

	return b.String()
}

func PositionError(pos int, err error) *LexerError {
	return &LexerError{
		position: &pos,
		err:      err,
	}
}

func TokenError(token *Token, err error) *LexerError {
	return &LexerError{
		token: token,
		err:   err,
	}
}
