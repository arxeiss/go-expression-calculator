package parser

import (
	"github.com/arxeiss/go-expression-calculator/lexer"
)

type TokenAssociativity int

const (
	LeftAssociativity TokenAssociativity = iota
	RightAssociativity
)

type TokenMeta struct {
	// Precedence sets operator priority
	Precedence int
	// Associativity tels if token is executed Left-to-Right or Right-to-Left
	Associativity TokenAssociativity
}

type TokenPriorities map[lexer.TokenType]TokenMeta

func DefaultTokenPriorities() TokenPriorities {
	return TokenPriorities{
		lexer.EOL:        TokenMeta{Precedence: 0},
		lexer.Whitespace: TokenMeta{Precedence: 0},

		lexer.Addition:     TokenMeta{Precedence: 20},
		lexer.Substraction: TokenMeta{Precedence: 20},

		lexer.Multiplication: TokenMeta{Precedence: 40},
		lexer.Division:       TokenMeta{Precedence: 40},
		lexer.FloorDiv:       TokenMeta{Precedence: 40},
		lexer.Modulus:        TokenMeta{Precedence: 40},

		lexer.UnaryAddition:     TokenMeta{Precedence: 60},
		lexer.UnarySubstraction: TokenMeta{Precedence: 60},

		lexer.Exponent: TokenMeta{Precedence: 80, Associativity: RightAssociativity},

		lexer.LPar: TokenMeta{Precedence: 100},
		lexer.RPar: TokenMeta{Precedence: 100},

		lexer.Identifier: TokenMeta{Precedence: 120},
		lexer.Number:     TokenMeta{Precedence: 120},
	}
}

func (tp TokenPriorities) GetPrecedence(tokenType lexer.TokenType) int {
	if meta, ok := tp[tokenType]; ok {
		return meta.Precedence
	}
	return 0
}

func (tp TokenPriorities) GetAssociativity(tokenType lexer.TokenType) TokenAssociativity {
	if meta, ok := tp[tokenType]; ok {
		return meta.Associativity
	}
	return LeftAssociativity
}
