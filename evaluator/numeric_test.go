package evaluator_test

import (
	"errors"
	"math"
	"strings"

	"github.com/arxeiss/go-expression-calculator/ast"
	"github.com/arxeiss/go-expression-calculator/evaluator"
	"github.com/arxeiss/go-expression-calculator/lexer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
)

var _ = Describe("Evaluator", func() {
	It("Check value", func() {
		ev, err := evaluator.NewNumericEvaluator(nil)
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

	It("Check assign to variable", func() {
		ev, err := evaluator.NewNumericEvaluator(map[string]float64{"a": 10}, nil)
		Expect(err).To(Succeed())
		Expect(ev.VariableList()).To(HaveLen(1))

		res, err := ev.Eval(ast.NewAssignNode(
			ast.NewVariableNode("myVar", nil),
			ast.NewBinaryNode(
				ast.Addition,
				ast.NewVariableNode("a", nil),
				ast.NewNumericNode(55, nil),
				nil,
			),
			nil,
		))
		Expect(ev.VariableList()).To(ConsistOf(
			MatchAllFields(Fields{"Name": Equal("a"), "Value": BeEquivalentTo(10)}),
			MatchAllFields(Fields{"Name": Equal("myvar"), "Value": BeEquivalentTo(65)}),
		))
		Expect(res).To(BeEquivalentTo(65))
		Expect(err).To(Succeed())
	})

	It("Check undefined variable", func() {
		ev, err := evaluator.NewNumericEvaluator(map[string]float64{"myVar": 89}, nil)
		Expect(err).To(Succeed())
		_, err = ev.Eval(ast.NewVariableNode("anotherVar", lexer.NewToken(lexer.Addition, 0, "", 10, 12)))
		evalErr, ok := err.(*evaluator.Error)
		Expect(ok).To(BeTrue())
		Expect(evalErr.Position()).To(Equal(10))
		Expect(evalErr.Error()).To(ContainSubstring("undefined variable 'anotherVar'"))
	})

	Describe("Handle unary", func() {
		It("Check addition", func() {
			ev, err := evaluator.NewNumericEvaluator(nil)
			Expect(err).To(Succeed())
			res, err := ev.Eval(ast.NewUnaryNode(ast.Addition, ast.NewNumericNode(33, nil), nil))
			Expect(res).To(BeEquivalentTo(33))
			Expect(err).To(Succeed())
		})
		It("Check substraction", func() {
			ev, err := evaluator.NewNumericEvaluator(nil)
			Expect(err).To(Succeed())
			res, err := ev.Eval(ast.NewUnaryNode(ast.Substraction, ast.NewNumericNode(33, nil), nil))
			Expect(res).To(BeEquivalentTo(-33))
			Expect(err).To(Succeed())
		})
		It("Check error", func() {
			ev, err := evaluator.NewNumericEvaluator(nil)
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
			ast.NewBinaryNode(ast.Multiplication,
				ast.NewNumericNode(3.87, nil), ast.NewVariableNode("myVar", nil), nil),
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

	DescribeTable(
		"Handle function",
		func(rootNode ast.Node, errStr string) {
			f := func(x ...float64) (float64, error) { return 0, errors.New("just some error") }
			ev, err := evaluator.NewNumericEvaluator(nil, map[string]evaluator.FunctionHandler{
				"a": {MinArguments: 0, MaxArguments: 0, Handler: func(x ...float64) (float64, error) { return 0, nil }},
				"b": {MinArguments: 3, MaxArguments: 0, Handler: f},
				"c": {MinArguments: 2, MaxArguments: 2, Handler: f},
				"d": {MinArguments: 2, MaxArguments: 3, Handler: f},
			})
			Expect(err).To(Succeed())
			_, err = ev.Eval(rootNode)
			evalErr, ok := err.(*evaluator.Error)
			Expect(ok).To(BeTrue())
			Expect(evalErr.Error()).To(ContainSubstring(errStr))
		},
		Entry("Send args", ast.NewFunctionNode("a", []ast.Node{
			ast.NewNumericNode(20, nil),
			ast.NewNumericNode(10, nil),
		}, nil), "function 'a' require 0 arguments, got 2"),
		Entry("Not enough args", ast.NewFunctionNode("b", []ast.Node{
			ast.NewNumericNode(20, nil),
			ast.NewNumericNode(10, nil),
		}, nil), "function 'b' require at least 3 arguments, got 2"),
		Entry("Too many args", ast.NewFunctionNode("c", []ast.Node{
			ast.NewNumericNode(20, nil),
			ast.NewNumericNode(10, nil),
			ast.NewNumericNode(10, nil),
		}, nil), "function 'c' require 2 arguments, got 3"),
		Entry("Not enough var args", ast.NewFunctionNode("d", []ast.Node{
			ast.NewNumericNode(20, nil),
		}, nil), "function 'd' require between 2 and 3 arguments, got 1"),
		Entry("Too many var args", ast.NewFunctionNode("d", []ast.Node{
			ast.NewNumericNode(20, nil),
			ast.NewNumericNode(20, nil),
			ast.NewNumericNode(20, nil),
			ast.NewNumericNode(20, nil),
		}, nil), "function 'd' require between 2 and 3 arguments, got 4"),
		Entry("Internal error function", ast.NewFunctionNode("d", []ast.Node{
			ast.NewNumericNode(20, nil),
			ast.NewNumericNode(20, nil),
		}, lexer.NewToken(lexer.Number, 123, "", 3, 6)), "just some error in function 'd' at position 3"),
		Entry("Error in function param", ast.NewFunctionNode("d", []ast.Node{
			ast.NewFunctionNode("f", nil, lexer.NewToken(lexer.Identifier, 0, "f", 7, 8)),
			ast.NewNumericNode(20, nil),
		}, lexer.NewToken(lexer.Number, 123, "", 3, 6)), "undefined function 'f' at position 7"),
	)

	Describe("Handle function", func() {
		It("Check substraction", func() {
			ev, err := evaluator.NewNumericEvaluator(nil)
			Expect(err).To(Succeed())
			res, err := ev.Eval(ast.NewUnaryNode(ast.Substraction, ast.NewNumericNode(33, nil), nil))
			Expect(res).To(BeEquivalentTo(-33))
			Expect(err).To(Succeed())
		})
		It("Check error", func() {
			ev, err := evaluator.NewNumericEvaluator(nil)
			Expect(err).To(Succeed())
			_, err = ev.Eval(ast.NewUnaryNode(ast.Multiplication, ast.NewNumericNode(33, nil), nil))
			Expect(err).To(MatchError(ContainSubstring("unary node supports only Addition and Substraction operator")))
		})
	})

	It("Check function", func() {
		ev, err := evaluator.NewNumericEvaluator(nil, map[string]evaluator.FunctionHandler{
			"myFunc": {
				Handler:      func(x ...float64) (float64, error) { return x[0] + 2, nil },
				MinArguments: 1, MaxArguments: 1,
			},
		})
		Expect(err).To(Succeed())
		res, err := ev.Eval(ast.NewFunctionNode("myFunc", []ast.Node{ast.NewNumericNode(7, nil)}, nil))
		Expect(res).To(BeEquivalentTo(9))
		Expect(err).To(Succeed())
	})

	It("Check undefined function", func() {
		ev, err := evaluator.NewNumericEvaluator(nil)
		Expect(err).To(Succeed())
		_, err = ev.Eval(ast.NewFunctionNode("myFunc", []ast.Node{ast.NewNumericNode(7, nil)}, nil))
		Expect(err).To(MatchError(ContainSubstring("undefined function 'myFunc'")))
	})

	It("Handle complex tree", func() {
		ev, err := evaluator.NewNumericEvaluator(
			map[string]float64{"X": 13.8, "Y": 8.9, "Z": 3},
			map[string]evaluator.FunctionHandler{
				"AddTwo": {
					Handler:      func(x ...float64) (float64, error) { return x[0] + 2, nil },
					MinArguments: 1, MaxArguments: 1,
				},
				"Abs": {
					Handler:      func(x ...float64) (float64, error) { return math.Abs(x[0]), nil },
					MinArguments: 1, MaxArguments: 1,
				},
				"Ceil": {
					Handler:      func(x ...float64) (float64, error) { return math.Ceil(x[0]), nil },
					MinArguments: 1, MaxArguments: 1,
				},
			},
		)
		Expect(err).To(Succeed())
		t := ast.NewBinaryNode(
			ast.Addition,
			ast.NewUnaryNode(ast.Substraction, ast.NewFunctionNode(
				"AddTwo",
				[]ast.Node{
					ast.NewBinaryNode(
						ast.Modulus,
						ast.NewNumericNode(1246.67, nil),
						ast.NewVariableNode("Y", nil),
						nil,
					),
				},
				nil,
			), nil),
			ast.NewBinaryNode(
				ast.Multiplication,
				ast.NewFunctionNode("Abs", []ast.Node{ast.NewBinaryNode(
					ast.Substraction,
					ast.NewNumericNode(5, nil),
					ast.NewVariableNode("X", nil),
					nil,
				)}, nil),
				ast.NewFunctionNode(
					"Ceil",
					[]ast.Node{
						ast.NewBinaryNode(
							ast.Exponent,
							ast.NewNumericNode(2.2, nil),
							ast.NewVariableNode("Z", nil),
							nil,
						),
					},
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

	It("Check error when defining same function with different case sensitivity", func() {
		_, err := evaluator.NewNumericEvaluator(nil, evaluator.MathFunctions(), map[string]evaluator.FunctionHandler{
			"aBS": {Description: "test", Handler: func(x ...float64) (float64, error) { return 0, nil }},
		})
		// order in map is non-deterministic, so names in error can also be in different order
		Expect(strings.ToLower(err.Error())).To(
			ContainSubstring("function named 'abs' was defined as 'abs' before"),
		)
	})

	It("Variable list", func() {
		ev, err := evaluator.NewNumericEvaluator(map[string]float64{
			"c": 10,
			"a": 5,
			"b": 12,
		})
		Expect(err).To(Succeed())
		sortedVarList := ev.VariableList()
		Expect(sortedVarList).To(HaveLen(3))
		Expect(sortedVarList[0].Name).To(Equal("a"))
		Expect(sortedVarList[1].Name).To(Equal("b"))
		Expect(sortedVarList[2].Name).To(Equal("c"))
	})

	It("Function list", func() {
		ev, err := evaluator.NewNumericEvaluator(nil, map[string]evaluator.FunctionHandler{
			"d": {},
			"k": {},
			"j": {},
			"c": {},
		})
		Expect(err).To(Succeed())
		sortedFuncList := ev.FunctionList()
		Expect(sortedFuncList).To(HaveLen(4))
		Expect(sortedFuncList[0].Name).To(Equal("c"))
		Expect(sortedFuncList[1].Name).To(Equal("d"))
		Expect(sortedFuncList[2].Name).To(Equal("j"))
		Expect(sortedFuncList[3].Name).To(Equal("k"))
	})
})
