package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"

	"github.com/arxeiss/go-expression-calculator/lexer"
)

const PrettyPrintErrorOffset = 15

func main() {
	l := lexer.NewLexer("  2.58e0 + (sin(5^4)) *  55  +   77 // 8  ")
	expr, err := l.Tokenize()
	prettyPrintError(l, err)
	_ = expr
}

func prettyPrintError(l *lexer.Lexer, err error) {
	if err == nil {
		return
	}
	lexerErr := &lexer.LexerError{}
	if !errors.As(err, &lexerErr) {
		fmt.Println(err.Error())
		return
	}
	pos := lexerErr.Position()
	if pos < 0 {
		fmt.Println(err.Error())
		return
	}
	expr := l.Expression()
	start := pos - PrettyPrintErrorOffset
	end := pos + PrettyPrintErrorOffset
	if start < 0 {
		start = 0
	}
	if end > len(expr) {
		end = len(expr)
	}

	fmt.Print(color.RedString("\n  Error: "), color.HiRedString(lexerErr.Error()))
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
