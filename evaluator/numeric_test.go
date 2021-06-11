package evaluator_test

import (
	"math"

	"github.com/arxeiss/go-expression-calculator/ast"
	"github.com/arxeiss/go-expression-calculator/evaluator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("Evaluator", func() {
	It("Check value", func() {
		ev := evaluator.NewNumericEvaluator(nil, nil)
		res, err := ev.Eval(ast.NumericNode(17))
		Expect(res).To(BeEquivalentTo(17))
		Expect(err).To(Succeed())
	})

	It("Check variable", func() {
		ev := evaluator.NewNumericEvaluator(map[string]float64{"myVar": 89}, nil)
		res, err := ev.Eval(ast.VariableNode("myVar"))
		Expect(res).To(BeEquivalentTo(89))
		Expect(err).To(Succeed())
	})

	It("Check undefined variable", func() {
		ev := evaluator.NewNumericEvaluator(map[string]float64{"myVar": 89}, nil)
		_, err := ev.Eval(ast.VariableNode("anotherVar"))
		Expect(err).To(MatchError("undefined variable 'anotherVar'"))
	})

	Describe("Handle unary", func() {
		It("Check addition", func() {
			ev := evaluator.NewNumericEvaluator(nil, nil)
			res, err := ev.Eval(ast.NewUnaryNode(ast.Addition, ast.NumericNode(33)))
			Expect(res).To(BeEquivalentTo(33))
			Expect(err).To(Succeed())
		})
		It("Check substraction", func() {
			ev := evaluator.NewNumericEvaluator(nil, nil)
			res, err := ev.Eval(ast.NewUnaryNode(ast.Substraction, ast.NumericNode(33)))
			Expect(res).To(BeEquivalentTo(-33))
			Expect(err).To(Succeed())
		})
		It("Check error", func() {
			ev := evaluator.NewNumericEvaluator(nil, nil)
			_, err := ev.Eval(ast.NewUnaryNode(ast.Multiplication, ast.NumericNode(33)))
			Expect(err).To(MatchError("unary node supports only Addition and Substraction operator"))
		})
	})

	DescribeTable("Handle binary",
		func(rn ast.Node, expRes float64, errMatcher types.GomegaMatcher) {
			ev := evaluator.NewNumericEvaluator(map[string]float64{"myVar": 2.89, "intVar": 3}, nil)
			res, err := ev.Eval(rn)
			Expect(res).To(Equal(expRes))
			Expect(err).To(errMatcher)
		},
		Entry("Addition",
			ast.NewBinaryNode(ast.Addition, ast.NumericNode(3.87), ast.VariableNode("myVar")),
			3.87+2.89,
			Succeed()),
		Entry("Substraction",
			ast.NewBinaryNode(ast.Substraction, ast.NumericNode(3.87), ast.VariableNode("myVar")),
			3.87-2.89,
			Succeed()),
		Entry("Multiplication",
			ast.NewBinaryNode(ast.Multiplication, ast.NumericNode(3.87), ast.VariableNode("myVar")),
			3.87*2.89,
			Succeed()),
		Entry("Division",
			ast.NewBinaryNode(ast.Division, ast.NumericNode(3.87), ast.VariableNode("myVar")),
			3.87/2.89,
			Succeed()),
		Entry("Exponent",
			ast.NewBinaryNode(ast.Exponent, ast.NumericNode(3.75), ast.VariableNode("intVar")),
			52.734375,
			Succeed()),
		Entry("FloorDiv",
			ast.NewBinaryNode(ast.FloorDiv, ast.NumericNode(38.7), ast.VariableNode("myVar")),
			13.0,
			Succeed()),
		Entry("Modulus",
			ast.NewBinaryNode(ast.Modulus, ast.NumericNode(39), ast.NumericNode(2.5)),
			1.5,
			Succeed()),
		Entry("Error operation",
			ast.NewBinaryNode(ast.Invalid, ast.NumericNode(39), ast.NumericNode(2.5)),
			0.0,
			MatchError("unimplemented operator Invalid")),
	)

	It("Check function", func() {
		ev := evaluator.NewNumericEvaluator(nil, map[string]evaluator.FunctionHandler{
			"myFunc": func(x ...float64) (float64, error) { return x[0] + 2, nil },
		})
		res, err := ev.Eval(ast.NewFunctionNode("myFunc", ast.NumericNode(7)))
		Expect(res).To(BeEquivalentTo(9))
		Expect(err).To(Succeed())
	})

	It("Check undefined function", func() {
		ev := evaluator.NewNumericEvaluator(nil, nil)
		_, err := ev.Eval(ast.NewFunctionNode("myFunc", ast.NumericNode(7)))
		Expect(err).To(MatchError("undefined function 'myFunc'"))
	})

	It("Handle complex tree", func() {
		ev := evaluator.NewNumericEvaluator(
			map[string]float64{"X": 13.8, "Y": 8.9, "Z": 3},
			map[string]evaluator.FunctionHandler{
				"AddTwo": func(x ...float64) (float64, error) { return x[0] + 2, nil },
				"Abs":    func(x ...float64) (float64, error) { return math.Abs(x[0]), nil },
				"Ceil":   func(x ...float64) (float64, error) { return math.Ceil(x[0]), nil },
			},
		)
		t := ast.NewBinaryNode(
			ast.Addition,
			ast.NewUnaryNode(ast.Substraction, ast.NewFunctionNode(
				"AddTwo",
				ast.NewBinaryNode(ast.Modulus, ast.NumericNode(1246.67), ast.VariableNode("Y")),
			)),
			ast.NewBinaryNode(
				ast.Multiplication,
				ast.NewFunctionNode("Abs", ast.NewBinaryNode(
					ast.Substraction,
					ast.NumericNode(5),
					ast.VariableNode("X"),
				)),
				ast.NewFunctionNode(
					"Ceil",
					ast.NewBinaryNode(ast.Exponent, ast.NumericNode(2.2), ast.VariableNode("Z")),
				),
			),
		)
		res, err := ev.Eval(t)
		Expect(res).To(Equal(94.13))
		Expect(err).To(Succeed())
	})
})
