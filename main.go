package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"

	"github.com/arxeiss/go-expression-calculator/ast"
	"github.com/arxeiss/go-expression-calculator/lexer"
	"github.com/arxeiss/go-expression-calculator/parser"
	"github.com/arxeiss/go-expression-calculator/parser/shuntyard"
)

const PrettyPrintErrorOffset = 15

func main() {
	l := lexer.NewLexer("  2.58e0 + (sin(5^4)) *  y  +  77 // some_variable  ")
	expr, err := l.Tokenize()
	if err != nil {
		prettyPrintError(l.Expression(), err)
		return
	}
	rootNode, err := shuntyard.NewParser(parser.DefaultTokenPriorities()).Parse(expr)
	if err != nil {
		prettyPrintError(l.Expression(), err)
		return
	}
	fmt.Print(ast.ToTreeDrawer(rootNode))
}

func prettyPrintError(expr string, err error) {
	if err == nil {
		return
	}
	pos := -1
	prefix := "Error"
	lexerErr := &lexer.Error{}
	if errors.As(err, &lexerErr) {
		pos = lexerErr.Position()
		prefix = "Lexer error"
	}
	parserErr := &parser.Error{}
	if errors.As(err, &parserErr) {
		pos = parserErr.Position()
		prefix = "Parser error"
	}

	if pos < 0 {
		fmt.Println(err.Error())
		return
	}
	start := pos - PrettyPrintErrorOffset
	end := pos + PrettyPrintErrorOffset
	if start < 0 {
		start = 0
	}
	if end > len(expr) {
		end = len(expr)
	}

	fmt.Print(color.RedString("\n  %s: ", prefix), color.HiRedString(err.Error()))
	fmt.Printf(
		"\n   | %s\n   | %s^\n\n",
		colorizeCode(expr[start:end]),
		color.HiBlackString(strings.Repeat(".", pos-start)),
	)
}

func colorizeCode(text string) string {
	r := regexp.MustCompile(`(\S+)|\s+`)
	b := strings.Builder{}
	for _, v := range r.FindAllStringSubmatch(text, -1) {
		if v[1] != "" {
			b.WriteString(color.HiYellowString(v[1]))
		} else {
			b.WriteString(color.BlackString(strings.Repeat(".", len(v[0]))))
		}
	}
	return b.String()
}
