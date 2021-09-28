package parser

import (
	"errors"
	"math"

	"github.com/arxeiss/go-expression-calculator/lexer"
)

type TokenAssociativity uint16
type TokenPrecedence uint16

const (
	LeftAssociativity TokenAssociativity = iota
	RightAssociativity
)

type TokenMeta struct {
	// Precedence sets operator priority
	Precedence TokenPrecedence
	// Associativity tels if token is executed Left-to-Right or Right-to-Left
	Associativity TokenAssociativity
}
type TokenPriorities map[lexer.TokenType]TokenMeta

var (
	ErrZeroPrecedenceSet = errors.New("precedence must be greater than 0")
)

func DefaultTokenPriorities() TokenPriorities {
	// It make sense to have precedence and associativity set only to operators
	return TokenPriorities{
		lexer.Equal: TokenMeta{Precedence: 10, Associativity: RightAssociativity},

		lexer.Addition:     TokenMeta{Precedence: 20},
		lexer.Substraction: TokenMeta{Precedence: 20},

		lexer.Multiplication: TokenMeta{Precedence: 40},
		lexer.Division:       TokenMeta{Precedence: 40},
		lexer.FloorDiv:       TokenMeta{Precedence: 40},
		lexer.Modulus:        TokenMeta{Precedence: 40},

		lexer.UnaryAddition:     TokenMeta{Precedence: 60},
		lexer.UnarySubstraction: TokenMeta{Precedence: 60},

		lexer.Exponent: TokenMeta{Precedence: 80, Associativity: RightAssociativity},
	}
}

func (tp TokenPriorities) Normalize() error {
	for k := range tp {
		switch k {
		case lexer.Equal, lexer.Addition, lexer.Substraction, lexer.Multiplication, lexer.Division, lexer.FloorDiv,
			lexer.Modulus, lexer.UnaryAddition, lexer.UnarySubstraction, lexer.Exponent:

		default:
			delete(tp, k)
		}
	}
	if tp.MinPrecedence() == 0 {
		return ErrZeroPrecedenceSet
	}
	return nil
}

func (tp TokenPriorities) GetMeta(tokenType lexer.TokenType) TokenMeta {
	if meta, ok := tp[tokenType]; ok {
		return meta
	}
	return TokenMeta{}
}

func (tp TokenPriorities) GetPrecedence(tokenType lexer.TokenType) TokenPrecedence {
	return tp.GetMeta(tokenType).Precedence
}

func (tp TokenPriorities) GetAssociativity(tokenType lexer.TokenType) TokenAssociativity {
	return tp.GetMeta(tokenType).Associativity
}

func (tp TokenPriorities) MaxPrecedence() TokenPrecedence {
	max := TokenPrecedence(0)
	for _, v := range tp {
		if v.Precedence > max {
			max = v.Precedence
		}
	}
	return max
}

func (tp TokenPriorities) MinPrecedence() TokenPrecedence {
	min := TokenPrecedence(math.MaxUint16)
	for _, v := range tp {
		if v.Precedence < min {
			min = v.Precedence
		}
	}
	return min
}

func (tp TokenPriorities) NextPrecedence(current TokenPrecedence) TokenPrecedence {
	next := tp.MaxPrecedence()
	for _, v := range tp {
		if v.Precedence > current && v.Precedence < next {
			next = v.Precedence
		}
	}
	return next
}
