package recursivedescent

import (
	"errors"
	"fmt"
	"strings"

	"github.com/arxeiss/go-expression-calculator/ast"
	"github.com/arxeiss/go-expression-calculator/lexer"
	"github.com/arxeiss/go-expression-calculator/parser"
)

var (
	ErrEmptyInput      = errors.New("there are no tokens to parse")
	ErrExpectedOperand = errors.New("expected number, identifier or left parenthesis")
	ErrExpectedEOL     = errors.New("last token is expected to be the end of input")
	ErrUnexpectedToken = errors.New("unexpected token")

	binaryOperators = []lexer.TokenType{
		lexer.Addition, lexer.Substraction,
		lexer.Multiplication, lexer.Division, lexer.FloorDiv, lexer.Modulus,
		lexer.Exponent,
	}
)

type Parser struct {
	priorities parser.TokenPriorities
}

func NewParser(priorities parser.TokenPriorities) (parser.Parser, error) {
	if err := priorities.Normalize(); err != nil {
		return nil, err
	}
	return &Parser{
		priorities: priorities,
	}, nil
}

type parserInstance struct {
	tokenList     []*lexer.Token
	i             int
	parser        *Parser
	maxPrecedence parser.TokenPrecedence
}

// Parse uses Recursive Descent parser.
func (p *Parser) Parse(tokenList []*lexer.Token) (ast.Node, error) {
	noWhiteSpaceList := make([]*lexer.Token, 0)
	for _, v := range tokenList {
		if v.Type() != lexer.Whitespace {
			noWhiteSpaceList = append(noWhiteSpaceList, v)
		}
	}
	if len(noWhiteSpaceList) == 0 {
		return nil, parser.ParseError(nil, ErrEmptyInput)
	}
	if lastToken := noWhiteSpaceList[len(noWhiteSpaceList)-1]; lastToken.Type() != lexer.EOL {
		return nil, parser.ParseError(lastToken, ErrExpectedEOL)
	}
	n, err := (&parserInstance{
		tokenList:     noWhiteSpaceList,
		i:             0,
		parser:        p,
		maxPrecedence: p.priorities.MaxPrecedence(),
	}).parseBlock()

	if err == nil && n == nil {
		return nil, parser.ParseError(noWhiteSpaceList[0], ErrEmptyInput)
	}
	return n, err
}

func (p *parserInstance) getPrecedence(tokenType lexer.TokenType) parser.TokenPrecedence {
	return p.parser.priorities.GetPrecedence(tokenType)
}
func (p *parserInstance) getAssociativity(tokenType lexer.TokenType) parser.TokenAssociativity {
	return p.parser.priorities.GetAssociativity(tokenType)
}

func (p *parserInstance) parseBlock() (ast.Node, error) {
	// Parse assignment
	var node ast.Node
	var err error

	if p.hasNth(0, lexer.Identifier) && p.hasNth(1, lexer.Equal) {
		variable, _ := p.expect()
		equalOp, _ := p.expect()

		node, err = p.parseExpression(p.getPrecedence(equalOp.Type()))
		if err == nil {
			node = ast.NewAssignNode(ast.NewVariableNode(variable.Identifier(), variable), node, equalOp)
		}
	} else {
		node, err = p.parseExpression(p.parser.priorities.MinPrecedence())
	}
	if err != nil {
		return nil, err
	}
	if !p.has(lexer.EOL) {
		return nil, parser.ParseError(p.current(), ErrUnexpectedToken)
	}
	return node, nil
}

func (p *parserInstance) parseExpression(currentPrecedence parser.TokenPrecedence) (ast.Node, error) {
	var node ast.Node
	var err error

	// Always nest to deepest precedence as "normal recursive descent" would do by calling methods like factor and term
	if currentPrecedence < p.maxPrecedence {
		if node, err = p.parseExpression(p.parser.priorities.NextPrecedence(currentPrecedence)); err != nil {
			return nil, err
		}
	}

	// If there is no node returned, we should expect either term or unary operators
	if node == nil {
		switch {
		case p.has(lexer.LPar, lexer.Identifier, lexer.Number):
			node, err = p.parseTerm()
		// Unary operators are checked only if their precedence match current one
		case p.has(lexer.Addition) && currentPrecedence == p.getPrecedence(lexer.UnaryAddition),
			p.has(lexer.Substraction) && currentPrecedence == p.getPrecedence(lexer.UnarySubstraction):
			node, err = p.handleUnary()
		}
	}
	if err != nil {
		return nil, err
	}

	// Iterate until there are operators with same precedence
	for p.getPrecedence(p.current().Type()) == currentPrecedence {
		node, err = p.handleSamePrecedenceTokens(currentPrecedence, node)
		if err != nil {
			return nil, err
		}
	}

	return node, err
}

func (p *parserInstance) handleSamePrecedenceTokens(
	currentPrecedence parser.TokenPrecedence,
	leftNode ast.Node,
) (ast.Node, error) {
	// Left part is matched, we always need to find operator now
	current := p.current()
	operatorToken, err := p.expect(binaryOperators...)
	if err != nil {
		return nil, err
	}
	// If there is operator, but nothing on the left side, it is always an error
	if leftNode == nil {
		return nil, parser.ParseError(current, ErrExpectedOperand)
	}

	var rightNode ast.Node
	current = p.current()

	// Has another opearator after operator, it must be unary
	if p.has(lexer.Addition, lexer.Substraction) {
		rightNode, err = p.handleUnary()
	} else {
		// If operator is RightPrecedence, keep same precedence as right parts should be lower in AST
		// So next iteration with same precedence will parse it first
		nextPrecedence := currentPrecedence
		if p.getAssociativity(operatorToken.Type()) == parser.LeftAssociativity {
			nextPrecedence = p.parser.priorities.NextPrecedence(currentPrecedence)
		}
		rightNode, err = p.parseExpression(nextPrecedence)
	}
	if err != nil {
		return nil, err
	}
	// If there is no part on the right side of operator, it must be error
	if rightNode == nil {
		return nil, parser.ParseError(current, ErrExpectedOperand)
	}

	return ast.NewBinaryNode(tokenTypeToOperation(operatorToken.Type()), leftNode, rightNode, operatorToken), nil
}

func (p *parserInstance) handleUnary() (ast.Node, error) {
	token, err := p.expect(lexer.Addition, lexer.Substraction)
	if err != nil {
		return nil, err
	}
	_ = token.ChangeToUnary()
	current := p.current()

	// Recurse again to match more unary operators or operators with higher precedence
	node, err := p.parseExpression(p.getPrecedence(token.Type()))
	if err != nil {
		return nil, err
	}
	// If there is unary operator and nothing returned (not either term) it is error always
	if node == nil {
		return nil, parser.ParseError(current, ErrExpectedOperand)
	}
	return ast.NewUnaryNode(tokenTypeToOperation(token.Type()), node, token), nil
}

func (p *parserInstance) parseTerm() (ast.Node, error) {
	token, err := p.expect(lexer.LPar, lexer.Identifier, lexer.Number)
	if err != nil {
		return nil, err
	}
	var node ast.Node
	switch token.Type() {
	case lexer.LPar:
		// If there is left parenthesis, just nest with lowest priority to handle sub-expression
		node, err = p.parseExpression(p.parser.priorities.MinPrecedence())
		if err != nil {
			return nil, err
		}
		// Last must be right parenthesis again
		if _, err := p.expect(lexer.RPar); err != nil {
			return nil, err
		}
	case lexer.Identifier:
		if p.has(lexer.LPar) {
			_, _ = p.expect() // Just pop out if it is function
			args := []ast.Node{}
			for {
				// Parse sub-expresion as argument
				node, err = p.parseExpression(p.parser.priorities.MinPrecedence())
				if err != nil {
					return nil, err
				}
				if node == nil { // When there is nothing, it is most like end
					break
				}
				args = append(args, node)
				if !p.has(lexer.Comma) {
					break
				}
				_, _ = p.expect() // If pop out the comma
			}
			// Function call must end with right parenthesis
			if _, err := p.expect(lexer.RPar); err != nil {
				return nil, err
			}
			node = ast.NewFunctionNode(token.Identifier(), args, token)
		} else {
			node = ast.NewVariableNode(token.Identifier(), token)
		}
	case lexer.Number:
		node = ast.NewNumericNode(token.Value(), token)
	}

	return node, nil
}

func (p *parserInstance) moveForward() {
	p.i++
}

func (p *parserInstance) has(expectedTypes ...lexer.TokenType) bool {
	return p.hasNth(0, expectedTypes...)
}

func (p *parserInstance) hasNth(nth int, expectedTypes ...lexer.TokenType) bool {
	nthToken := p.nextNth(nth)
	for _, expType := range expectedTypes {
		if expType == nthToken.Type() {
			return true
		}
	}
	return false
}

func (p *parserInstance) expect(expectedTypes ...lexer.TokenType) (*lexer.Token, error) {
	defer p.moveForward()

	if len(expectedTypes) == 0 {
		return p.tokenList[p.i], nil
	}

	anyOf := []string{}
	current := p.current()
	for _, expType := range expectedTypes {
		if current.Type() == expType {
			return current, nil
		}
		anyOf = append(anyOf, expType.String())
	}
	var err error
	if len(anyOf) > 1 {
		err = fmt.Errorf("expected one of ['%s'] types, got '%s'", strings.Join(anyOf, "', '"), current.Type())
	} else {
		err = fmt.Errorf("expected '%s' type, got '%s'", anyOf[0], current.Type())
	}
	return nil, parser.ParseError(p.current(), err)
}

func (p *parserInstance) current() *lexer.Token {
	return p.nextNth(0)
}

func (p *parserInstance) nextNth(nth int) *lexer.Token {
	if p.i+nth < len(p.tokenList) {
		return p.tokenList[p.i+nth]
	}
	return nil
}

func tokenTypeToOperation(tt lexer.TokenType) ast.Operation {
	switch tt {
	case lexer.UnaryAddition, lexer.Addition:
		return ast.Addition
	case lexer.UnarySubstraction, lexer.Substraction:
		return ast.Substraction
	case lexer.Multiplication:
		return ast.Multiplication
	case lexer.Division:
		return ast.Division
	case lexer.Exponent:
		return ast.Exponent
	case lexer.FloorDiv:
		return ast.FloorDiv
	case lexer.Modulus:
		return ast.Modulus
	}
	return ast.Invalid
}
