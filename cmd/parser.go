package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/arxeiss/go-expression-calculator/ast"
	"github.com/arxeiss/go-expression-calculator/evaluator"
	"github.com/arxeiss/go-expression-calculator/lexer"
	"github.com/arxeiss/go-expression-calculator/parser"
	"github.com/arxeiss/go-expression-calculator/parser/shuntyard"
)

func parseExpression(numEval *evaluator.NumericEvaluator, expr string) {
	printTree := false
	if strings.HasPrefix(expr, "tree") {
		expr = strings.TrimSpace(expr[4:])
		printTree = true
	}

	l := lexer.NewLexer(expr)
	tokenized, err := l.Tokenize()
	if err != nil {
		prettyPrintError(expr, err)
		return
	}
	rootNode, err := shuntyard.NewParser(parser.DefaultTokenPriorities()).Parse(tokenized)
	if err != nil {
		prettyPrintError(l.Expression(), err)
		return
	}
	value, err := numEval.Eval(rootNode)
	if err != nil {
		prettyPrintError(l.Expression(), err)
		return
	}
	fmt.Printf("%s %f\n", color.HiBlackString("<-"), value)
	if printTree {
		fmt.Println(color.HiCyanString(" Here comes the AST Tree: "))
		fmt.Print(ast.ToTreeDrawer(rootNode))
	}
}
