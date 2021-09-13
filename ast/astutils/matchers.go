package astutils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	errorsutil "github.com/onsi/gomega/gstruct/errors"
	"github.com/onsi/gomega/types"

	"github.com/arxeiss/go-expression-calculator/ast"
)

type binaryMatcher struct {
	operation ast.Operation
	left      types.GomegaMatcher
	right     types.GomegaMatcher
	failures  []error
}

type assignMatcher struct {
	left     types.GomegaMatcher
	right    types.GomegaMatcher
	failures []error
}

type unaryMatcher struct {
	operation ast.Operation
	next      types.GomegaMatcher
	failures  []error
}

type numericMatcher struct {
	value    interface{}
	failures []error
}

type variableMatcher struct {
	name     interface{}
	failures []error
}

type functionMatcher struct {
	name     interface{}
	param    types.GomegaMatcher
	failures []error
}

func MatchBinaryNode(operation ast.Operation, left types.GomegaMatcher, right types.GomegaMatcher) types.GomegaMatcher {
	return &binaryMatcher{
		operation: operation,
		left:      left,
		right:     right,
	}
}
func MatchAssignNode(left types.GomegaMatcher, right types.GomegaMatcher) types.GomegaMatcher {
	return &assignMatcher{
		left:  left,
		right: right,
	}
}
func MatchUnaryNode(operation ast.Operation, node types.GomegaMatcher) types.GomegaMatcher {
	return &unaryMatcher{
		operation: operation,
		next:      node,
	}
}

// MatchNumericNode expects types.GomegaMatcher or passed value will be compared with gomega.Equal
func MatchNumericNode(value interface{}) types.GomegaMatcher {
	return &numericMatcher{
		value: value,
	}
}

// MatchVariableNode expects types.GomegaMatcher or passed name will be compared with gomega.Equal
func MatchVariableNode(name interface{}) types.GomegaMatcher {
	return &variableMatcher{
		name: name,
	}
}

func MatchFunctionNode(name interface{}, param types.GomegaMatcher) types.GomegaMatcher {
	return &functionMatcher{
		name:  name,
		param: param,
	}
}

func matchNode(matcher types.GomegaMatcher, current interface{}, nestStr string, failures []error) []error {
	m, err := matcher.Match(current)
	if err != nil {
		failures = append(failures, errorsutil.Nest(nestStr, err))
	} else if !m {
		failures = append(
			failures,
			errorsutil.Nest(nestStr, errors.New(matcher.FailureMessage(current))),
		)
	}
	return failures
}

func matchOperation(nodeOp, matcherOp ast.Operation, failures []error) []error {
	if nodeOp != matcherOp {
		failures = append(failures, fmt.Errorf(" -> operator %s, got %s", matcherOp.String(), nodeOp.String()))
	}
	return failures
}

func formatFailureMessage(failures []error, actual interface{}) string {
	strFailures := make([]string, len(failures))
	for i := range failures {
		strFailures[i] = failures[i].Error()
	}

	t := reflect.TypeOf(actual)
	n := t.Name()
	if n == "" && t.Kind() == reflect.Ptr {
		n = t.Elem().Name()
	}

	return fmt.Sprintf("Expected to %s to match: \n%v\n", n, strings.Join(strFailures, "\n"))
}

func (matcher *binaryMatcher) Match(actual interface{}) (success bool, err error) {
	if node, ok := actual.(*ast.BinaryNode); ok {
		matcher.failures = matchOperation(node.Operator(), matcher.operation, matcher.failures)
		matcher.failures = matchNode(matcher.left, node.Left(), " -> Left", matcher.failures)
		matcher.failures = matchNode(matcher.right, node.Right(), " -> Right", matcher.failures)

		return len(matcher.failures) == 0, nil
	}
	return false, fmt.Errorf("matcher MatchBinaryNode expects a `*ast.BinaryNode` Got:\n%s", format.Object(actual, 1))
}

func (matcher *binaryMatcher) FailureMessage(actual interface{}) (message string) {
	return formatFailureMessage(matcher.failures, actual)
}

func (matcher *binaryMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("not to match node %s", format.Object(actual, 0))
}

func (matcher *assignMatcher) Match(actual interface{}) (success bool, err error) {
	if node, ok := actual.(*ast.AssignNode); ok {
		matcher.failures = matchNode(matcher.left, node.Left(), " -> Left", matcher.failures)
		matcher.failures = matchNode(matcher.right, node.Right(), " -> Right", matcher.failures)

		return len(matcher.failures) == 0, nil
	}
	return false, fmt.Errorf("matcher MatchAssignNode expects a `*ast.AssignNode` Got:\n%s", format.Object(actual, 1))
}

func (matcher *assignMatcher) FailureMessage(actual interface{}) (message string) {
	return formatFailureMessage(matcher.failures, actual)
}

func (matcher *assignMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("not to match node %s", format.Object(actual, 0))
}

func (matcher *unaryMatcher) Match(actual interface{}) (success bool, err error) {
	if node, ok := actual.(*ast.UnaryNode); ok {
		matcher.failures = matchOperation(node.Operator(), matcher.operation, matcher.failures)
		matcher.failures = matchNode(matcher.next, node.Next(), " -> Next", matcher.failures)

		return len(matcher.failures) == 0, nil
	}
	return false, fmt.Errorf("matcher MatchUnaryNode expects a `*ast.UnaryNode` Got:\n%s", format.Object(actual, 1))
}

func (matcher *unaryMatcher) FailureMessage(actual interface{}) (message string) {
	return formatFailureMessage(matcher.failures, actual)
}

func (matcher *unaryMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("not to match node %s", format.Object(actual, 0))
}

func (matcher *numericMatcher) Match(actual interface{}) (success bool, err error) {
	if val, ok := actual.(*ast.NumericNode); ok {
		var valMatcher types.GomegaMatcher
		switch vm := matcher.value.(type) {
		case types.GomegaMatcher:
			valMatcher = vm
		case float64:
			valMatcher = gomega.Equal(vm)
		case float32, int, int32, int64:
			valMatcher = gomega.BeEquivalentTo(vm)
		default:
			valMatcher = gomega.Equal(matcher.value)
		}
		matcher.failures = matchNode(valMatcher, val.Value(), " -> Value", matcher.failures)

		return len(matcher.failures) == 0, nil
	}
	return false, fmt.Errorf("matcher MatchNumericNode expects a `ast.NumericNode` Got:\n%s", format.Object(actual, 1))
}

func (matcher *numericMatcher) FailureMessage(actual interface{}) (message string) {
	return formatFailureMessage(matcher.failures, actual)
}

func (matcher *numericMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("not to match value %s", format.Object(actual, 0))
}

func (matcher *variableMatcher) Match(actual interface{}) (success bool, err error) {
	if val, ok := actual.(*ast.VariableNode); ok {
		var valMatcher types.GomegaMatcher
		if vm, ok := matcher.name.(types.GomegaMatcher); ok {
			valMatcher = vm
		} else {
			valMatcher = gomega.Equal(matcher.name)
		}
		matcher.failures = matchNode(valMatcher, val.Name(), " -> Name", matcher.failures)

		return len(matcher.failures) == 0, nil
	}
	return false, fmt.Errorf(
		"matcher MatchVariableNode expects a `ast.VariableNode` Got:\n%s", format.Object(actual, 1))
}

func (matcher *variableMatcher) FailureMessage(actual interface{}) (message string) {
	return formatFailureMessage(matcher.failures, actual)
}

func (matcher *variableMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("not to match value %s", format.Object(actual, 0))
}

func (matcher *functionMatcher) Match(actual interface{}) (success bool, err error) {
	if node, ok := actual.(*ast.FunctionNode); ok {
		var valMatcher types.GomegaMatcher
		if vm, ok := matcher.name.(types.GomegaMatcher); ok {
			valMatcher = vm
		} else {
			valMatcher = gomega.Equal(matcher.name)
		}
		matcher.failures = matchNode(valMatcher, node.Name(), " -> Name", matcher.failures)
		matcher.failures = matchNode(matcher.param, node.Param(), " -> Param", matcher.failures)

		return len(matcher.failures) == 0, nil
	}
	return false, fmt.Errorf(
		"matcher MatchFunctionNode expects a `*ast.FunctionNode` Got:\n%s", format.Object(actual, 1))
}

func (matcher *functionMatcher) FailureMessage(actual interface{}) (message string) {
	return formatFailureMessage(matcher.failures, actual)
}

func (matcher *functionMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("not to match value %s", format.Object(actual, 0))
}
