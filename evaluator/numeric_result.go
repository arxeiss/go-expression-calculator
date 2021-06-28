package evaluator

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/arxeiss/go-expression-calculator/ast"
)

type FunctionHandler func(x ...float64) (float64, error)

type NumericEvaluator struct {
	variables map[string]float64
	functions map[string]FunctionHandler
}

func NewNumericEvaluator(vars map[string]float64, funcs map[string]FunctionHandler) (*NumericEvaluator, error) {
	variables := make(map[string]float64)
	{
		varNames := make(map[string]string)
		for kcs, v := range vars {
			k := strings.ToLower(kcs)
			if pn, has := varNames[k]; has {
				return nil, fmt.Errorf(
					"variable with name '%s' was defined as '%s' before, variables are case insensitive", kcs, pn)
			}
			varNames[k] = kcs
			variables[k] = v
		}
	}

	functions := make(map[string]FunctionHandler)
	{
		funcsNames := make(map[string]string)
		for kcs, v := range funcs {
			k := strings.ToLower(kcs)
			if pn, has := funcsNames[k]; has {
				return nil, fmt.Errorf(
					"function named '%s' was defined as '%s' before, function names are case insensitive", kcs, pn)
			}
			funcsNames[k] = kcs
			functions[k] = v
		}
	}

	return &NumericEvaluator{
		variables: variables,
		functions: functions,
	}, nil
}

func (e *NumericEvaluator) Eval(rootNode ast.Node) (float64, error) {
	switch n := rootNode.(type) {
	case *ast.BinaryNode:
		return e.handleBinary(n)
	case *ast.UnaryNode:
		return e.handleUnary(n)
	case *ast.FunctionNode:
		return e.handleFunction(n)
	case *ast.VariableNode:
		if v, has := e.variables[strings.ToLower(n.Name())]; has {
			return v, nil
		}
		return 0, EvalError(n.GetToken(), fmt.Errorf("undefined variable '%s'", n.Name()))
	case *ast.NumericNode:
		return n.Value(), nil
	}
	return 0, EvalError(rootNode.GetToken(), fmt.Errorf("unimplemented node type %T", e))
}

func (e *NumericEvaluator) handleUnary(n *ast.UnaryNode) (float64, error) {
	val, err := e.Eval(n.Next())
	if err != nil {
		return 0, err
	}

	switch n.Operator() {
	case ast.Substraction:
		return -val, nil
	case ast.Addition:
		return val, nil
	}

	return 0, EvalError(n.GetToken(), errors.New("unary node supports only Addition and Substraction operator"))
}

func (e *NumericEvaluator) handleBinary(n *ast.BinaryNode) (float64, error) {
	l, err := e.Eval(n.Left())
	if err != nil {
		return 0, err
	}
	r, err := e.Eval(n.Right())
	if err != nil {
		return 0, err
	}

	switch n.Operator() {
	case ast.Addition:
		return l + r, nil
	case ast.Substraction:
		return l - r, nil
	case ast.Multiplication:
		return l * r, nil
	case ast.Division:
		return l / r, nil
	case ast.FloorDiv:
		return math.Floor(l / r), nil
	case ast.Exponent:
		return math.Pow(l, r), nil
	case ast.Modulus:
		return math.Mod(l, r), nil
	}

	return 0, EvalError(n.GetToken(), fmt.Errorf("unimplemented operator %s", n.Operator()))
}

func (e *NumericEvaluator) handleFunction(n *ast.FunctionNode) (float64, error) {
	f, has := e.functions[strings.ToLower(n.Name())]
	if !has {
		return 0, EvalError(n.GetToken(), fmt.Errorf("undefined function '%s'", n.Name()))
	}
	v, err := e.Eval(n.Param())
	if err != nil {
		return 0, err
	}

	return f(v)
}
