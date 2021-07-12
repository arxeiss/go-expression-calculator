package lexer

import (
	"errors"
	"regexp"
	"strconv"
)

var (
	//nolint:lll
	tokenRegexp = regexp.MustCompile(
		`\(|\)|\*\*|\^|//|%|\+|\-|\*|/|(?P<num>(?:[0-9]+(?:\.[0-9]+)?|\.[0-9]+)(?:e[+-]?[0-9]+)?)|(?P<id>(?i)[a-z_][a-z0-9_]*)|(?P<ws>\s+)`,
	)
)

type Lexer struct {
	expr string
}

func NewLexer(expression string) *Lexer {
	return &Lexer{expr: expression}
}

func (l *Lexer) Expression() string {
	return l.expr
}

// Tokenize converts input expresion into the list of tokens
func (l *Lexer) Tokenize() ([]*Token, error) {
	expr := make([]*Token, 0)
	subMatchNames := tokenRegexp.SubexpNames()

	lastIndex := 0
	for _, indexes := range tokenRegexp.FindAllStringSubmatchIndex(l.expr, -1) {
		t := &Token{startPos: indexes[0], endPos: indexes[1]}
		// If current token does not start where previous ended, there is something unexpected
		if t.startPos != lastIndex {
			return nil, PositionError(lastIndex, ErrUnexpectedChar)
		}
		lastIndex = t.endPos

		// submatch contains numbers, identifier and whitespace
		if handled, err := l.handleSubMatches(t, indexes, subMatchNames); err != nil {
			return nil, err
		} else if !handled {
			// tokens are not part of submatch
			t.tType = operatorTokenType(l.expr[t.startPos:t.endPos])
		}
		// Returned EOL means some internal error, ie unhandled characters
		if t.tType == EOL {
			return nil, PositionError(lastIndex, ErrUnexpectedChar)
		}

		expr = append(expr, t)
	}
	// If all regex matches are processed, but there is still some text
	if lastIndex != len(l.expr) {
		return nil, PositionError(lastIndex, ErrUnexpectedChar)
	}
	// Always add EOL for easier handling in parsers
	expr = append(expr, &Token{tType: EOL, startPos: lastIndex, endPos: lastIndex})
	return expr, nil
}

func (l *Lexer) handleSubMatches(t *Token, indexes []int, subMatchNames []string) (bool, error) {
	for i := 1; i < len(subMatchNames); i++ {
		// There are always begin and end index for each submatch
		// So if the value at index is 0, submatch was not found
		if indexes[i*2] < 0 {
			continue
		}

		switch subMatchNames[i] {
		case "num":
			t.tType = Number
			var err error
			if t.value, err = strconv.ParseFloat(l.expr[t.startPos:t.endPos], 64); err != nil {
				if errors.Is(err, strconv.ErrRange) {
					return false, TokenError(t, ErrNumberOutOfRange)
				}
				return false, TokenError(t, ErrInvalidNumber)
			}
			return true, nil
		case "id":
			t.tType = Identifier
			t.idName = l.expr[t.startPos:t.endPos]
			return true, nil
		case "ws":
			t.tType = Whitespace
			return true, nil
		}
	}
	return false, nil
}

func operatorTokenType(operator string) TokenType {
	switch operator {
	case "(":
		return LPar
	case ")":
		return RPar
	case "^", "**":
		return Exponent
	case "*":
		return Multiplication
	case "/":
		return Division
	case "//":
		return FloorDiv
	case "%":
		return Modulus
	case "+":
		return Addition
	case "-":
		return Substraction
	}
	return EOL
}
