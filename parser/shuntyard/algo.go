package shuntyard

import (
	"errors"
	"fmt"

	"github.com/arxeiss/go-expression-calculator/ast"
	"github.com/arxeiss/go-expression-calculator/lexer"
	"github.com/arxeiss/go-expression-calculator/parser"
)

type expectState uint8

const (
	operandToken expectState = iota
	operatorToken
)

var (
	ErrEmptyInput       = errors.New("there are no tokens to parse")
	ErrExpectedOperand  = errors.New("expected number, identifier or left parenthesis")
	ErrMultipleUnary    = errors.New("too many unary operators in a row")
	ErrExpectedOperator = errors.New("expected operator or right parenthesis")
	ErrUnexpectedEOL    = errors.New("unexpected end of input")
	ErrMissingLPar      = errors.New("cannot find matching left parenthesis")
	ErrMissingRPar      = errors.New("cannot find matching right parenthesis")
)

type Parser struct {
	priorities parser.TokenPriorities
}

func NewParser(priorities parser.TokenPriorities) parser.Parser {
	return &Parser{
		priorities: priorities,
	}
}

// Parse uses Shunting Yard algorithm to parse the input.
// Algorithm is based on Wiki page https://en.wikipedia.org/wiki/Shunting-yard_algorithm
// with some improvements discussed on StackOverflow https://stackoverflow.com/a/29652095/1513087
// and modified to produce Abstract Syntax Tree rather than RPN
func (p *Parser) Parse(tokenList []*lexer.Token) (ast.Node, error) {
	expect := operandToken
	output := make([]ast.Node, 0)
	opStack := make([]*lexer.Token, 0)
	var err error

	if len(tokenList) == 0 {
		return nil, parser.ParseError(nil, ErrEmptyInput)
	}

tokenLoop:
	for i := 0; i < len(tokenList); i++ {
		curToken := tokenList[i]

		switch curToken.Type() {
		case lexer.Whitespace:
			continue

		case lexer.Number:
			expect, output, err = p.handleNumber(expect, curToken, output)

		case lexer.Identifier:
			expect, opStack, output, err = p.handleIdentifier(expect, tokenList, i, opStack, output)

		case lexer.UnaryAddition, lexer.UnarySubstraction, lexer.Addition, lexer.Substraction:
			// If the received token is Addition or Substraction and we expect operand, it is probably unary operator
			if expect == operandToken {
				opStack, err = p.handleUnary(curToken, opStack)
				break
			}
			// If operator is expected, fallthrough to handle operator
			fallthrough
		case lexer.Exponent, lexer.Multiplication, lexer.Division, lexer.FloorDiv, lexer.Modulus:
			expect, opStack, output, err = p.handleOperator(expect, curToken, opStack, output)

		case lexer.LPar:
			expect, opStack, err = p.handleLPar(expect, curToken, opStack)

		case lexer.RPar:
			expect, opStack, output, err = p.handleRPar(expect, curToken, opStack, output)

		case lexer.EOL:
			if expect == operandToken {
				return nil, parser.ParseError(curToken, ErrUnexpectedEOL)
			}
			break tokenLoop
		default:
			return nil, fmt.Errorf("parser internal error, unhandled '%s' type", curToken.Type())
		}

		if err != nil {
			return nil, err
		}
	}

	return p.clearOpStack(opStack, output)
}

func (p *Parser) clearOpStack(opStack []*lexer.Token, output []ast.Node) (ast.Node, error) {
	var err error
	for len(opStack) > 0 {
		topStackEl := opStack[len(opStack)-1]
		if topStackEl.Type() == lexer.LPar {
			return nil, parser.ParseError(topStackEl, ErrMissingRPar)
		}
		output, err = p.addToOutput(output, topStackEl)
		if err != nil {
			return nil, err
		}
		opStack = opStack[:len(opStack)-1]
	}

	if len(output) != 1 {
		return nil, fmt.Errorf("internal error, at the end output must contain single node")
	}

	return output[0], err
}

// handleNumber parse number if expected token in operand, otherwise error is returned
func (*Parser) handleNumber(
	expect expectState,
	curToken *lexer.Token,
	output []ast.Node,
) (expectState, []ast.Node, error) {
	if expect == operatorToken {
		return expect, nil, parser.ParseError(curToken, ErrExpectedOperator)
	}
	output = append(output, ast.NewNumericNode(curToken.Value(), curToken))
	expect = operatorToken

	return expect, output, nil
}

// handleIdentifier parse incoming identifier and decide whenever it is function or variable
// if operator is expected, returns error
func (*Parser) handleIdentifier(
	expect expectState,
	tokenList []*lexer.Token,
	currIndex int,
	opStack []*lexer.Token,
	output []ast.Node,
) (expectState, []*lexer.Token, []ast.Node, error) {
	if expect == operatorToken {
		return expect, opStack, output, parser.ParseError(tokenList[currIndex], ErrExpectedOperator)
	}
	// Find next token to decide if identifier is function
	nextTokenType := lexer.EOL
	for n := currIndex + 1; n < len(tokenList); n++ {
		if tokenList[n].Type() != lexer.Whitespace {
			nextTokenType = tokenList[n].Type()
			break
		}
	}
	if nextTokenType == lexer.LPar {
		// Expecting the identifier is function name, because is followed by (
		opStack = append(opStack, tokenList[currIndex])
		expect = operandToken
	} else {
		// Identifier is variable name
		output = append(output, ast.NewVariableNode(tokenList[currIndex].Identifier(), tokenList[currIndex]))
		expect = operatorToken
	}

	return expect, opStack, output, nil
}

// handleUnary convert current Addition or Substraction into unary token
func (*Parser) handleUnary(
	curToken *lexer.Token,
	opStack []*lexer.Token,
) ([]*lexer.Token, error) {
	if len(opStack) > 0 {
		// When unary operator is found, the expect flag is not changing
		// It is required to check last element of operator stack. If is unary, throw error
		lastOp := opStack[len(opStack)-1]
		if lastOp.Type() == lexer.UnaryAddition || lastOp.Type() == lexer.UnarySubstraction {
			return nil, parser.ParseError(curToken, ErrMultipleUnary)
		}
	}
	if err := curToken.ChangeToUnary(); err != nil {
		return nil, err
	}
	opStack = append(opStack, curToken)

	return opStack, nil
}

// handleOperator parse token as binary operator or return error if operand is expected
func (p *Parser) handleOperator(
	expect expectState,
	curToken *lexer.Token,
	opStack []*lexer.Token,
	output []ast.Node,
) (expectState, []*lexer.Token, []ast.Node, error) {
	if expect == operandToken {
		return expect, nil, nil, parser.ParseError(curToken, ErrExpectedOperand)
	}
	for len(opStack) > 0 {
		topStackEl := opStack[len(opStack)-1]
		// Continue only, if the operator at the top of the operator stack is not a left parenthesis
		if topStackEl.Type() == lexer.LPar {
			break
		}
		// If the operator at the top of the operator stack has greater precedence
		// OR
		// The operator at the top of the operator stack has equal precedence and the token is left associative
		if p.priorities.GetPrecedence(topStackEl.Type()) > p.priorities.GetPrecedence(curToken.Type()) ||
			(p.priorities.GetPrecedence(topStackEl.Type()) == p.priorities.GetPrecedence(curToken.Type()) &&
				p.priorities.GetAssociativity(curToken.Type()) == parser.LeftAssociativity) {
			var err error
			output, err = p.addToOutput(output, topStackEl)
			if err != nil {
				return expect, nil, nil, err
			}
			opStack = opStack[:len(opStack)-1]
			continue
		}
		// If no operation was done, exit loop
		break
	}
	opStack = append(opStack, curToken)
	expect = operandToken

	return expect, opStack, output, nil
}

// handleLPar parse left parenthesis or return error, if operator is expected
func (*Parser) handleLPar(
	expect expectState,
	curToken *lexer.Token,
	opStack []*lexer.Token,
) (expectState, []*lexer.Token, error) {
	if expect == operatorToken {
		return expect, nil, parser.ParseError(curToken, ErrExpectedOperator)
	}
	opStack = append(opStack, curToken)
	expect = operandToken

	return expect, opStack, nil
}

// handleRPar parse right parenthesis, checks all matching left parenthesis or return error if operand is expected
func (p *Parser) handleRPar(
	expect expectState,
	curToken *lexer.Token,
	opStack []*lexer.Token,
	output []ast.Node,
) (expectState, []*lexer.Token, []ast.Node, error) {
	if expect == operandToken {
		return expect, nil, nil, parser.ParseError(curToken, ErrExpectedOperand)
	}
	// Repeat until top of operator stack is no Left parenthesis
	for len(opStack) > 0 {
		topStackEl := opStack[len(opStack)-1]
		if topStackEl.Type() == lexer.LPar {
			// When left parenthesis is found, do not remove it
			break
		}
		var err error
		output, err = p.addToOutput(output, topStackEl)
		if err != nil {
			return expect, nil, nil, err
		}
		opStack = opStack[:len(opStack)-1]
	}
	// If operator stack is empty, there is no matching left parenthesis
	if len(opStack) == 0 {
		return expect, nil, nil, parser.ParseError(curToken, ErrMissingLPar)
	}
	// If not empty, there must be left parenthesis, just remove it
	opStack = opStack[:len(opStack)-1]

	// Check if left parenthesis was there because of function call
	// If identifier is found, it must be function. Variables are never added to operator stack
	if len(opStack) > 0 && opStack[len(opStack)-1].Type() == lexer.Identifier {
		// Remove it and add to output
		var err error
		output, err = p.addToOutput(output, opStack[len(opStack)-1])
		if err != nil {
			return expect, nil, nil, err
		}
		opStack = opStack[:len(opStack)-1]
	}

	expect = operatorToken
	return expect, opStack, output, nil
}

func (p *Parser) addToOutput(output []ast.Node, token *lexer.Token) ([]ast.Node, error) {
	var err error
	var op ast.Operation
	switch t := token.Type(); t {
	case lexer.UnaryAddition, lexer.UnarySubstraction:
		if len(output) < 1 {
			return nil, errors.New("internal error, missing value for unary operator")
		}
		if op, err = tokenTypeToOperation(t); err != nil {
			return nil, err
		}
		output[len(output)-1] = ast.NewUnaryNode(op, output[len(output)-1], token)
	case lexer.Addition, lexer.Substraction, lexer.Multiplication, lexer.Division, lexer.Exponent,
		lexer.FloorDiv, lexer.Modulus:

		if len(output) < 2 {
			return nil, errors.New("internal error, missing values for binary operator")
		}
		r := output[len(output)-1]
		l := output[len(output)-2]
		output = output[:len(output)-1]

		if op, err = tokenTypeToOperation(t); err != nil {
			return nil, err
		}
		output[len(output)-1] = ast.NewBinaryNode(op, l, r, token)
	case lexer.Identifier:
		// Current functions support only 1 parameter
		if len(output) < 1 {
			return nil, errors.New("internal error, missing value for function")
		}
		output[len(output)-1] = ast.NewFunctionNode(token.Identifier(), output[len(output)-1], token)
	default:
		return nil, fmt.Errorf("unexpected token '%s' received to add to output", t.String())
	}
	return output, err
}

func tokenTypeToOperation(tt lexer.TokenType) (ast.Operation, error) {
	switch tt {
	case lexer.UnaryAddition, lexer.Addition:
		return ast.Addition, nil
	case lexer.UnarySubstraction, lexer.Substraction:
		return ast.Substraction, nil
	case lexer.Multiplication:
		return ast.Multiplication, nil
	case lexer.Division:
		return ast.Division, nil
	case lexer.Exponent:
		return ast.Exponent, nil
	case lexer.FloorDiv:
		return ast.FloorDiv, nil
	case lexer.Modulus:
		return ast.Modulus, nil
	}
	return ast.Invalid, fmt.Errorf("missing convertion of %s to AST operation", tt.String())
}
