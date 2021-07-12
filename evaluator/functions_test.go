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

	type funcArg struct {
		value         float64
		resultMatcher types.GomegaMatcher
	}

	DescribeTable("Math Functions",
		func(name string, descriptionMatcher types.GomegaMatcher, testArgs []funcArg) {
			f, has := mathFunctions[name]
			Expect(has).To(BeTrue())
			Expect(f.Description).To(descriptionMatcher)
			for i, v := range testArgs {
				Expect(f.Handler(v.value)).To(v.resultMatcher, fmt.Sprintf("%d. match - value = %f", i+1, v.value))
			}
		},
		Entry("Abs", "Abs", ContainSubstring("Abs returns the absolute value"), []funcArg{
			{value: 155, resultMatcher: BeEquivalentTo(155)},
			{value: -15, resultMatcher: BeEquivalentTo(15)},
			{value: 0, resultMatcher: BeEquivalentTo(0)},
		}),
		Entry("Acos", "Acos", ContainSubstring("Acos returns the arccosine, in radians"), []funcArg{
			{value: 0, resultMatcher: BeEquivalentTo(0.5 * math.Pi)},
			{value: 1, resultMatcher: BeEquivalentTo(0)},
			{value: 0.5, resultMatcher: BeNumerically("~", math.Pi/3)},
		}),
		Entry("Asin", "Asin", ContainSubstring("Asin returns the arcsine, in radians"), []funcArg{
			{value: 0, resultMatcher: BeEquivalentTo(0)},
			{value: 1, resultMatcher: BeEquivalentTo(0.5 * math.Pi)},
			{value: 0.5, resultMatcher: BeNumerically("~", math.Pi/6)},
		}),
		Entry("Atan", "Atan", ContainSubstring("Atan returns the arctangent, in radians"), []funcArg{
			{value: 0, resultMatcher: BeEquivalentTo(0)},
			{value: 1, resultMatcher: BeEquivalentTo(0.25 * math.Pi)},
			{value: -1, resultMatcher: BeEquivalentTo(-0.25 * math.Pi)},
		}),
		Entry("Ceil", "Ceil", ContainSubstring("Ceil returns the least integer value greater than or equal"), []funcArg{
			{value: 2.3, resultMatcher: BeEquivalentTo(3)},
			{value: 2.9, resultMatcher: BeEquivalentTo(3)},
			{value: -2.5, resultMatcher: BeEquivalentTo(-2)},
		}),
		Entry("Cos", "Cos", ContainSubstring("Cos returns the cosine of the radian"), []funcArg{
			{value: 0, resultMatcher: BeEquivalentTo(1)},
			{value: math.Pi, resultMatcher: BeEquivalentTo(-1)},
			{value: math.Pi * 0.5, resultMatcher: BeEquivalentTo(0)},
		}),
		Entry("Floor", "Floor", ContainSubstring("returns the greatest integer value less than or equal"), []funcArg{
			{value: 2.3, resultMatcher: BeEquivalentTo(2)},
			{value: 2.9, resultMatcher: BeEquivalentTo(2)},
			{value: -2.5, resultMatcher: BeEquivalentTo(-3)},
		}),
		Entry("Sin", "Sin", ContainSubstring("Sin returns the sine of the radian"), []funcArg{
			{value: 0, resultMatcher: BeEquivalentTo(0)},
			{value: math.Pi, resultMatcher: BeNumerically("~", 0)},
			{value: math.Pi * 0.5, resultMatcher: BeEquivalentTo(1)},
		}),
		Entry("Sqrt", "Sqrt", ContainSubstring("Sqrt returns the square root"), []funcArg{
			{value: 0, resultMatcher: BeEquivalentTo(0)},
			{value: 4, resultMatcher: BeEquivalentTo(2)},
			{value: 25, resultMatcher: BeEquivalentTo(5)},
			{value: 2, resultMatcher: BeNumerically("~", 1.414213562)},
		}),
		Entry("Tan", "Tan", ContainSubstring("Tan returns the tangent of the radian"), []funcArg{
			{value: 0, resultMatcher: BeEquivalentTo(0)},
			{value: math.Pi, resultMatcher: BeEquivalentTo(0)},
			{value: math.Pi * 0.25, resultMatcher: BeEquivalentTo(1)},
		}),
	)
})
