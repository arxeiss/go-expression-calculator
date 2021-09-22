package evaluator_test

import (
	"fmt"
	"math"

	"github.com/arxeiss/go-expression-calculator/evaluator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("Functions", func() {
	mathFunctions := evaluator.MathFunctions()
	mathFunctionsVarArgs := evaluator.MathFunctionsWithVarArgs()
	precission := 0.00000000001

	type funcArg struct {
		value                     []float64
		resultMatcher, errMatcher types.GomegaMatcher
	}

	DescribeTable("Math Functions",
		func(name string, descriptionMatcher types.GomegaMatcher, testArgs []funcArg) {
			f, has := mathFunctions[name]
			Expect(has).To(BeTrue())
			Expect(f.Description).To(descriptionMatcher)
			Expect(f.MinArguments).To(Equal(1))
			Expect(f.MaxArguments).To(Equal(1))
			for i, v := range testArgs {
				Expect(f.Handler(v.value...)).To(v.resultMatcher, fmt.Sprintf("%d. match - value = %f", i+1, v.value))
			}
		},
		Entry("abs", "abs", ContainSubstring("Returns the absolute value"), []funcArg{
			{value: []float64{155}, resultMatcher: BeEquivalentTo(155)},
			{value: []float64{-15}, resultMatcher: BeEquivalentTo(15)},
			{value: []float64{0}, resultMatcher: BeEquivalentTo(0)},
		}),
		Entry("acos", "acos", ContainSubstring("Returns the arccosine, in radians"), []funcArg{
			{value: []float64{0}, resultMatcher: BeEquivalentTo(0.5 * math.Pi)},
			{value: []float64{1}, resultMatcher: BeEquivalentTo(0)},
			{value: []float64{0.5}, resultMatcher: BeNumerically("~", math.Pi/3)},
		}),
		Entry("asin", "asin", ContainSubstring("Returns the arcsine, in radians"), []funcArg{
			{value: []float64{0}, resultMatcher: BeEquivalentTo(0)},
			{value: []float64{1}, resultMatcher: BeEquivalentTo(0.5 * math.Pi)},
			{value: []float64{0.5}, resultMatcher: BeNumerically("~", math.Pi/6)},
		}),
		Entry("atan", "atan", ContainSubstring("Returns the arctangent, in radians"), []funcArg{
			{value: []float64{0}, resultMatcher: BeEquivalentTo(0)},
			{value: []float64{1}, resultMatcher: BeEquivalentTo(0.25 * math.Pi)},
			{value: []float64{-1}, resultMatcher: BeEquivalentTo(-0.25 * math.Pi)},
		}),
		Entry("ceil", "ceil", ContainSubstring("Returns the least integer value greater than or equal"), []funcArg{
			{value: []float64{2.3}, resultMatcher: BeEquivalentTo(3)},
			{value: []float64{2.9}, resultMatcher: BeEquivalentTo(3)},
			{value: []float64{-2.5}, resultMatcher: BeEquivalentTo(-2)},
		}),
		Entry("cos", "cos", ContainSubstring("Returns the cosine of the radian"), []funcArg{
			{value: []float64{0}, resultMatcher: BeEquivalentTo(1)},
			{value: []float64{math.Pi}, resultMatcher: BeEquivalentTo(-1)},
			{value: []float64{math.Pi * 0.5}, resultMatcher: BeEquivalentTo(0)},
		}),
		Entry("floor", "floor", ContainSubstring("Returns the greatest integer value less than or equal"), []funcArg{
			{value: []float64{2.3}, resultMatcher: BeEquivalentTo(2)},
			{value: []float64{2.9}, resultMatcher: BeEquivalentTo(2)},
			{value: []float64{-2.5}, resultMatcher: BeEquivalentTo(-3)},
		}),
		Entry("sin", "sin", ContainSubstring("Returns the sine of the radian"), []funcArg{
			{value: []float64{0}, resultMatcher: BeEquivalentTo(0)},
			{value: []float64{math.Pi}, resultMatcher: BeNumerically("~", 0)},
			{value: []float64{math.Pi * 0.5}, resultMatcher: BeEquivalentTo(1)},
		}),
		Entry("sqrt", "sqrt", ContainSubstring("Returns the square root"), []funcArg{
			{value: []float64{0}, resultMatcher: BeEquivalentTo(0)},
			{value: []float64{4}, resultMatcher: BeEquivalentTo(2)},
			{value: []float64{25}, resultMatcher: BeEquivalentTo(5)},
			{value: []float64{2}, resultMatcher: BeNumerically("~", 1.414213562)},
		}),
		Entry("tan", "tan", ContainSubstring("Returns the tangent of the radian"), []funcArg{
			{value: []float64{0}, resultMatcher: BeEquivalentTo(0)},
			{value: []float64{math.Pi}, resultMatcher: BeEquivalentTo(0)},
			{value: []float64{math.Pi * 0.25}, resultMatcher: BeEquivalentTo(1)},
		}),
	)

	DescribeTable("Math Functions with VarArgs",
		func(name string, descMatcher types.GomegaMatcher, minArgs, maxArgs int, testArgs []funcArg) {
			f, has := mathFunctionsVarArgs[name]
			Expect(has).To(BeTrue())
			Expect(f.Description).To(descMatcher)
			Expect(f.MinArguments).To(Equal(minArgs))
			Expect(f.MaxArguments).To(Equal(maxArgs))
			for i, v := range testArgs {
				res, err := f.Handler(v.value...)
				Expect(err).To(v.errMatcher, fmt.Sprintf("%d. match", i+1))
				Expect(res).To(v.resultMatcher, fmt.Sprintf("%d. match - value = %f", i+1, v.value))
			}
		},
		Entry("pi", "pi", ContainSubstring("Returns Pi value."), 0, 0, []funcArg{
			{value: nil, resultMatcher: Equal(math.Pi), errMatcher: Succeed()},
		}),
		Entry("e", "e", ContainSubstring("Returns e value (base of natural logarithm)."), 0, 0, []funcArg{
			{value: nil, resultMatcher: Equal(math.E), errMatcher: Succeed()},
		}),
		Entry("phi", "phi", ContainSubstring("Returns Phi value."), 0, 0, []funcArg{
			{value: nil, resultMatcher: Equal(math.Phi), errMatcher: Succeed()},
		}),
		Entry("log", "log",
			ContainSubstring("Returns log of value n with given base."), 2, 2, []funcArg{
				{value: []float64{10, 2}, resultMatcher: BeNumerically("~", 3.32192809489, precission),
					errMatcher: Succeed()},
				{value: []float64{10, math.E}, resultMatcher: BeNumerically("~", 2.30258509299, precission),
					errMatcher: Succeed()},
				{value: []float64{10, 10}, resultMatcher: BeNumerically("==", 1), errMatcher: Succeed()},
			}),
		Entry("max", "max", ContainSubstring("Returns maximum of provided numbers."), 1, 0, []funcArg{
			{value: []float64{12.4, 12.2, 13.1, 13.00000001}, resultMatcher: Equal(13.1), errMatcher: Succeed()},
			{value: []float64{14.2}, resultMatcher: Equal(14.2), errMatcher: Succeed()},
			{value: []float64{-14.2, -6.3242}, resultMatcher: Equal(-6.3242), errMatcher: Succeed()},
		}),
		Entry("min", "min", ContainSubstring("Returns minimum of provided numbers."), 1, 0, []funcArg{
			{value: []float64{12.4, 12.2, 13.1, 13.00000001}, resultMatcher: Equal(12.2), errMatcher: Succeed()},
			{value: []float64{14.2}, resultMatcher: Equal(14.2), errMatcher: Succeed()},
			{value: []float64{-14.2, -6.3242}, resultMatcher: Equal(-14.2), errMatcher: Succeed()},
		}),
		Entry("rand_i", "rand_i",
			ContainSubstring("Returns random decimal number in range <0, a) or <a, b) if b is provided."), 1, 2,
			[]funcArg{
				{value: []float64{-20, -40}, resultMatcher: BeEquivalentTo(0),
					errMatcher: MatchError("number -20 (min) cannot be higher or equal to -40 (max)")},
				{value: []float64{20, 20}, resultMatcher: BeEquivalentTo(0),
					errMatcher: MatchError("number 20 (min) cannot be higher or equal to 20 (max)")},
				{value: []float64{0}, resultMatcher: BeEquivalentTo(0),
					errMatcher: MatchError("number 0 (min) cannot be higher or equal to 0 (max)")},
			}),
		Entry("nth_root", "nth_root", ContainSubstring("Returns n-th root of a."), 2, 2,
			[]funcArg{
				{value: []float64{-25, 5}, resultMatcher: BeEquivalentTo(0),
					errMatcher: MatchError("number a cannot be negative")},
				{value: []float64{25, -5}, resultMatcher: BeEquivalentTo(0),
					errMatcher: MatchError("number n cannot be 0 or negative")},
				{value: []float64{25, 2}, resultMatcher: BeNumerically("~", 5, precission), errMatcher: Succeed()},
				{value: []float64{25, 5}, resultMatcher: BeNumerically("~", 1.903653938715, precission),
					errMatcher: Succeed()},
				{value: []float64{25.99, 5.9}, resultMatcher: BeNumerically("~", 1.736991425077, precission),
					errMatcher: Succeed()},
			}),
	)
})
