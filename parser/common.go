package parser

import (
	"github.com/arxeiss/go-expression-calculator/ast"
	"github.com/arxeiss/go-expression-calculator/lexer"
)

type Parser interface {
	Parse(tokenList []*lexer.Token) (ast.Node, error)
}
