package evaluator_test

import (
	"math"
	"strings"

	"github.com/arxeiss/go-expression-calculator/ast"
	"github.com/arxeiss/go-expression-calculator/evaluator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("Evaluator", func() {
	It("Check value", func() {
		ev, err := evaluator.NewNumericEvaluator(nil, nil)
		Expect(err).To(Succeed())
		res, err := ev.Eval(ast.NewNumericNode(17, nil))
		Expect(res).To(BeEquivalentTo(17))
		Expect(err).To(Succeed())
	})

	It("Check variable", func() {
		ev, err := evaluator.NewNumericEvaluator(map[string]float64{"myVar": 89}, nil)
		Expect(err).To(Succeed())
		res, err := ev.Eval(ast.NewVariableNode("myVar", nil))
		Expect(res).To(BeEquivalentTo(89))
		Expect(err).To(Succeed())
	})

	It("Check undefined variable", func() {
		ev, err := evaluator.NewNumericEvaluator(map[string]float64{"myVar": 89}, nil)
		Expect(err).To(Succeed())
		_, err = ev.Eval(ast.NewVariableNode("anotherVar", nil))
		Expect(err).To(MatchError(ContainSubstring("undefined variable 'anotherVar'")))
	})

	Describe("Handle unary", func() {
		It("Check addition", func() {
			ev, err := evaluator.NewNumericEvaluator(nil, nil)
			Expect(err).To(Succeed())
			res, err := ev.Eval(ast.NewUnaryNode(ast.Addition, ast.NewNumericNode(33, nil), nil))
			Expect(res).To(BeEquivalentTo(33))
			Expect(err).To(Succeed())
		})
		It("Check substraction", func() {
			ev, err := evaluator.NewNumericEvaluator(nil, nil)
			Expect(err).To(Succeed())
			res, err := ev.Eval(ast.NewUnaryNode(ast.Substraction, ast.NewNumericNode(33, nil), nil))
			Expect(res).To(BeEquivalentTo(-33))
			Expect(err).To(Succeed())
		})
		It("Check error", func() {
			ev, err := evaluator.NewNumericEvaluator(nil, nil)
			Expect(err).To(Succeed())
			_, err = ev.Eval(ast.NewUnaryNode(ast.Multiplication, ast.NewNumericNode(33, nil), nil))
			Expect(err).To(MatchError(ContainSubstring("unary node supports only Addition and Substraction operator")))
		})
	})

	DescribeTable("Handle binary",
		func(rn ast.Node, expRes float64, errMatcher types.GomegaMatcher) {
			ev, err := evaluator.NewNumericEvaluator(map[string]float64{"myVar": 2.89, "intVar": 3}, nil)
			Expect(err).To(Succeed())
			res, err := ev.Eval(rn)
			Expect(res).To(Equal(expRes))
			Expect(err).To(errMatcher)
		},
		Entry("Addition",
			ast.NewBinaryNode(ast.Addition, ast.NewNumericNode(3.87, nil), ast.NewVariableNode("myVar", nil), nil),
			3.87+2.89,
			Succeed()),
		Entry("Substraction",
			ast.NewBinaryNode(ast.Substraction, ast.NewNumericNode(3.87, nil), ast.NewVariableNode("myVar", nil), nil),
			3.87-2.89,
			Succeed()),
		Entry("Multiplication",
			ast.NewBinaryNode(ast.Multiplication, ast.NewNumericNode(3.87, nil), ast.NewVariableNode("myVar", nil), nil),
			3.87*2.89,
			Succeed()),
		Entry("Division",
			ast.NewBinaryNode(ast.Division, ast.NewNumericNode(3.87, nil), ast.NewVariableNode("myVar", nil), nil),
			3.87/2.89,
			Succeed()),
		Entry("Exponent",
			ast.NewBinaryNode(ast.Exponent, ast.NewNumericNode(3.75, nil), ast.NewVariableNode("intVar", nil), nil),
			52.734375,
			Succeed()),
		Entry("FloorDiv",
			ast.NewBinaryNode(ast.FloorDiv, ast.NewNumericNode(38.7, nil), ast.NewVariableNode("myVar", nil), nil),
			13.0,
			Succeed()),
		Entry("Modulus",
			ast.NewBinaryNode(ast.Modulus, ast.NewNumericNode(39, nil), ast.NewNumericNode(2.5, nil), nil),
			1.5,
			Succeed()),
		Entry("Error operation",
			ast.NewBinaryNode(ast.Invalid, ast.NewNumericNode(39, nil), ast.NewNumericNode(2.5, nil), nil),
			0.0,
			MatchError(ContainSubstring("unimplemented operator Invalid"))),
	)

	It("Check function", func() {
		ev, err := evaluator.NewNumericEvaluator(nil, map[string]evaluator.FunctionHandler{
			"myFunc": func(x ...float64) (float64, error) { return x[0] + 2, nil },
		})
		Expect(err).To(Succeed())
		res, err := ev.Eval(ast.NewFunctionNode("myFunc", ast.NewNumericNode(7, nil), nil))
		Expect(res).To(BeEquivalentTo(9))
		Expect(err).To(Succeed())
	})

	It("Check undefined function", func() {
		ev, err := evaluator.NewNumericEvaluator(nil, nil)
		Expect(err).To(Succeed())
		_, err = ev.Eval(ast.NewFunctionNode("myFunc", ast.NewNumericNode(7, nil), nil))
		Expect(err).To(MatchError(ContainSubstring("undefined function 'myFunc'")))
	})

	It("Handle complex tree", func() {
		ev, err := evaluator.NewNumericEvaluator(
			map[string]float64{"X": 13.8, "Y": 8.9, "Z": 3},
			map[string]evaluator.FunctionHandler{
				"AddTwo": func(x ...float64) (float64, error) { return x[0] + 2, nil },
				"Abs":    func(x ...float64) (float64, error) { return math.Abs(x[0]), nil },
				"Ceil":   func(x ...float64) (float64, error) { return math.Ceil(x[0]), nil },
			},
		)
		Expect(err).To(Succeed())
		t := ast.NewBinaryNode(
			ast.Addition,
			ast.NewUnaryNode(ast.Substraction, ast.NewFunctionNode(
				"AddTwo",
				ast.NewBinaryNode(ast.Modulus, ast.NewNumericNode(1246.67, nil), ast.NewVariableNode("Y", nil), nil),
				nil,
			), nil),
			ast.NewBinaryNode(
				ast.Multiplication,
				ast.NewFunctionNode("Abs", ast.NewBinaryNode(
					ast.Substraction,
					ast.NewNumericNode(5, nil),
					ast.NewVariableNode("X", nil),
					nil,
				), nil),
				ast.NewFunctionNode(
					"Ceil",
					ast.NewBinaryNode(ast.Exponent, ast.NewNumericNode(2.2, nil), ast.NewVariableNode("Z", nil), nil),
					nil,
				),
				nil,
			),
			nil,
		)
		res, err := ev.Eval(t)
		Expect(res).To(Equal(94.13))
		Expect(err).To(Succeed())
	})

	It("Check error when defining same variable with different case sensitivity", func() {
		_, err := evaluator.NewNumericEvaluator(map[string]float64{
			"my_variable": 123,
			"my_Variable": 123,
		}, nil)
		// order in map is non-deterministic, so names in error can also be in different order
		Expect(strings.ToLower(err.Error())).To(
			ContainSubstring("variable with name 'my_variable' was defined as 'my_variable'"),
		)
	})
})
