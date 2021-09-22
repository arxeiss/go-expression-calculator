package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/arxeiss/go-expression-calculator/ast"
	"github.com/arxeiss/go-expression-calculator/evaluator"
	"github.com/arxeiss/go-expression-calculator/lexer"
	"github.com/arxeiss/go-expression-calculator/parser"
)

func parseLine(expr string, numEvaluator *evaluator.NumericEvaluator, p parser.Parser) {
	expr = strings.TrimSpace(expr)
	switch expr {
	case "help":
		fmt.Printf(
			"%s\n   %s - %s\n   %s - %s\n   %s - %s\n   %s - %s\n   %s - %s\n",
			"Write directly any expression to evaluate, or one of those commands:",
			color.HiYellowString("functions  "), "Show all available functions",
			color.HiYellowString("variables  "), "Prints all variables with values",
			color.HiYellowString("help       "), "Show this help",
			color.HiYellowString("tree {expr}"), "Write tree and then expression to print AST tree",
			color.HiYellowString("exit       "), "Quit this REPL",
		)
	case "func", "funcs", "functions":
		funcs := numEvaluator.FunctionList()
		if len(funcs) == 0 {
			fmt.Println(color.YellowString("There are no defined functions"))
			return
		}
		prettyPrintFunctions(funcs)
	case "vars", "variables":
		vars := numEvaluator.VariableList()
		if len(vars) == 0 {
			fmt.Println(color.YellowString("There are no variables now"))
			return
		}
		prettyPrintVariables(vars)
	default:
		parseExpression(numEvaluator, p, expr)
	}
}

func parseExpression(numEval *evaluator.NumericEvaluator, p parser.Parser, expr string) {
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
	rootNode, err := p.Parse(tokenized)
	if err != nil {
		prettyPrintError(l.Expression(), err)
		return
	}
	value, err := numEval.Eval(rootNode)
	if err != nil {
		prettyPrintError(l.Expression(), err)
		return
	}
	fmt.Printf("%s %.8f\n", color.HiBlackString("<-"), value)
	if printTree {
		fmt.Println(color.HiCyanString(" Here comes the AST Tree: "))
		fmt.Print(ast.ToTreeDrawer(rootNode))
	}
}
