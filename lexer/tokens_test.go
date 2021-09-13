package lexer_test

import (
	"github.com/arxeiss/go-expression-calculator/lexer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Token", func() {
	It("Getters and matcher", func() {
		t := lexer.NewToken(lexer.Addition, 10.123, "identifier", 10, 12)
		Expect(t.Type()).To(Equal(lexer.Addition))
		Expect(t.Value()).To(Equal(10.123))
		Expect(t.Identifier()).To(Equal("identifier"))
		Expect(t.StartPosition()).To(Equal(10))
		Expect(t.EndPosition()).To(Equal(12))

		Expect(t).To(PointTo(MatchToken(lexer.Addition, 10.123, "identifier", 10, 12)))
	})
})

var _ = DescribeTable("TokenType stringer",
	func(tt lexer.TokenType, expected string) {
		Expect(tt.String()).To(Equal(expected))
	},
	Entry("EOL", lexer.EOL, "EOL"),
	Entry("Whitespace", lexer.Whitespace, "Whitespace"),
	Entry("Identifier", lexer.Identifier, "Identifier"),
	Entry("LPar", lexer.LPar, "LPar"),
	Entry("RPar", lexer.RPar, "RPar"),
	Entry("Exponent", lexer.Exponent, "Exponent"),
	Entry("Multiplication", lexer.Multiplication, "Multiplication"),
	Entry("Division", lexer.Division, "Division"),
	Entry("FloorDiv", lexer.FloorDiv, "FloorDiv"),
	Entry("Modulus", lexer.Modulus, "Modulus"),
	Entry("Addition", lexer.Addition, "Addition"),
	Entry("Substraction", lexer.Substraction, "Substraction"),
	Entry("Number", lexer.Number, "Number"),
	Entry("Equal", lexer.Equal, "Equal"),
	Entry("Comma", lexer.Comma, "Comma"),
	Entry("UnaryAddition", lexer.UnaryAddition, "UnaryAddition"),
	Entry("UnarySubstraction", lexer.UnarySubstraction, "UnarySubstraction"),
)
