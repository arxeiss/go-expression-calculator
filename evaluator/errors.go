package evaluator

import (
	"strconv"
	"strings"

	"github.com/arxeiss/go-expression-calculator/lexer"
)

type Error struct {
	token *lexer.Token
	err   error
}

// Position returns 0 based index of error in original input expression
// If position is not set, -1 is returned
func (e *Error) Position() int {
	if e.token != nil {
		return e.token.StartPosition()
	}
	return -1
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Error() string {
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

		if e.token != nil {
			b.WriteString("; found ")
			b.WriteString(e.token.Type().String())
			b.WriteString(" token")
		}
	}

	b.WriteString(" at position ")
	b.WriteString(strconv.Itoa(pos))

	return b.String()
}

func EvalError(token *lexer.Token, err error) *Error {
	return &Error{
		token: token,
		err:   err,
	}
}
