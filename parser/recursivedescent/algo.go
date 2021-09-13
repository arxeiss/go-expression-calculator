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
	return (&parserInstance{
		tokenList:     noWhiteSpaceList,
		i:             0,
		parser:        p,
		maxPrecedence: p.priorities.MaxPrecedence(),
	}).parseBlock()
}

func (p *parserInstance) getPrecedence(tokenType lexer.TokenType) parser.TokenPrecedence {
	return p.parser.priorities.GetPrecedence(tokenType)
}
func (p *parserInstance) getAssociativity(tokenType lexer.TokenType) parser.TokenAssociativity {
	return p.parser.priorities.GetAssociativity(tokenType)
}

func (p *parserInstance) parseBlock() (ast.Node, error) {
	if len(p.tokenList) == 0 {
		return nil, parser.ParseError(nil, ErrEmptyInput)
	}

	// Parse assignment
	if p.hasNth(0, lexer.Identifier) && p.hasNth(1, lexer.Equal) {
		variable, _ := p.expect(lexer.Identifier)
		equalOp, _ := p.expect(lexer.Equal)

		right, err := p.parseExpression(p.getPrecedence(equalOp.Type()))
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.EOL); err != nil {
			return nil, err
		}

		return ast.NewAssignNode(ast.NewVariableNode(variable.Identifier(), variable), right, equalOp), nil
	}

	node, err := p.parseExpression(p.parser.priorities.MinPrecedence())
	if err == nil {
		_, err = p.expect(lexer.EOL)
	}
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parserInstance) parseExpression(currentPrecedence parser.TokenPrecedence) (ast.Node, error) {
	// var unaryToken *lexer.Token
	// if p.GetPrecedence(lastToken.Type()) >= p.GetPrecedence(lexer.UnaryAddition) &&
	// 	p.has(lexer.Addition, lexer.Substraction) {
	// 	unaryToken, _ = p.expect(lexer.Addition, lexer.Substraction)
	// }

	var leftNode ast.Node
	var err error
	if currentPrecedence < p.maxPrecedence {
		if leftNode, err = p.parseExpression(p.parser.priorities.NextPrecedence(currentPrecedence)); err != nil {
			return nil, err
		}
	}

	if leftNode == nil {
		current := p.current()
		switch {
		case p.has(lexer.LPar, lexer.Identifier, lexer.Number):
			leftNode, err = p.parseTerm()
		case p.has(lexer.Addition) && currentPrecedence == p.getPrecedence(lexer.UnaryAddition),
			p.has(lexer.Substraction) && currentPrecedence == p.getPrecedence(lexer.UnarySubstraction):

			token, _ := p.expect(lexer.Addition, lexer.Substraction)
			leftNode, err = p.parseExpression(currentPrecedence)
			if err == nil {
				leftNode = ast.NewUnaryNode(tokenTypeToOperation(token.Type()), leftNode, token)
			}
		}
		if err != nil {
			return nil, err
		}
		if leftNode == nil {
			return nil, parser.ParseError(current, ErrExpectedOperand)
		}
	}

	for p.getPrecedence(p.current().Type()) == currentPrecedence {
		current := p.current()
		operatorToken, err := p.expect(binaryOperators...)
		if err != nil {
			return nil, err
		}

		var rightNode ast.Node

		// Has another opearator after operator, it must be unary
		if p.has(lexer.Addition, lexer.Substraction) {
			token, _ := p.expect(lexer.Addition, lexer.Substraction)
			_ = token.ChangeToUnary()
			rightNode, err = p.parseExpression(p.getPrecedence(token.Type()))
			if err != nil {
				return nil, err
			}
			if rightNode == nil {
				return nil, parser.ParseError(current, ErrExpectedOperand)
			}
			rightNode = ast.NewUnaryNode(tokenTypeToOperation(token.Type()), rightNode, token)
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
		if rightNode == nil {
			return nil, parser.ParseError(current, ErrExpectedOperand)
		}

		leftNode = ast.NewBinaryNode(tokenTypeToOperation(operatorToken.Type()), leftNode, rightNode, operatorToken)
	}
	return leftNode, nil
}

func (p *parserInstance) parseTerm() (ast.Node, error) {
	token, err := p.expect(lexer.LPar, lexer.Identifier, lexer.Number)
	if err != nil {
		return nil, err
	}
	var node ast.Node
	switch token.Type() {
	case lexer.LPar:
		node, err = p.parseExpression(p.parser.priorities.MinPrecedence())
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.RPar); err != nil {
			return nil, err
		}
	case lexer.Identifier:
		if p.has(lexer.LPar) {
			_, _ = p.expect(lexer.LPar) // Just pop out if it is function
			args := []ast.Node{}
			for {
				node, err = p.parseExpression(p.parser.priorities.MinPrecedence())
				if err != nil {
					return nil, err
				}
				args = append(args, node)
				if !p.has(lexer.Comma) {
					break
				}
				_, _ = p.expect(lexer.Comma)
			}
			if _, err := p.expect(lexer.RPar); err != nil {
				return nil, err
			}
			node = ast.NewFunctionNode(token.Identifier(), args[0], token)
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
