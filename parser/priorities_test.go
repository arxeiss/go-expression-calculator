package parser_test

import (
	"github.com/arxeiss/go-expression-calculator/lexer"
	"github.com/arxeiss/go-expression-calculator/parser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Priorities", func() {
	It("Default priorities - Precedence", func() {

		p := parser.DefaultTokenPriorities()
		Expect(p.Normalize()).To(Succeed())

		equal := p.GetPrecedence(lexer.Equal)
		Expect(equal).To(BeNumerically(">", 0))

		addition := p.GetPrecedence(lexer.Addition)
		Expect(addition).To(BeNumerically(">", equal))
		Expect(p.GetPrecedence(lexer.Substraction)).To(BeNumerically("==", addition))
		Expect(p.NextPrecedence(equal)).To(Equal(addition))

		multiplication := p.GetPrecedence(lexer.Multiplication)
		Expect(multiplication).To(BeNumerically(">", addition))
		Expect(p.GetPrecedence(lexer.Division)).To(BeNumerically("==", multiplication))
		Expect(p.GetPrecedence(lexer.FloorDiv)).To(BeNumerically("==", multiplication))
		Expect(p.GetPrecedence(lexer.Modulus)).To(BeNumerically("==", multiplication))
		Expect(p.NextPrecedence(addition)).To(Equal(multiplication))

		unaryAddition := p.GetPrecedence(lexer.UnaryAddition)
		Expect(unaryAddition).To(BeNumerically(">", multiplication))
		Expect(p.GetPrecedence(lexer.UnarySubstraction)).To(BeNumerically("==", unaryAddition))
		Expect(p.NextPrecedence(multiplication)).To(Equal(unaryAddition))

		exponent := p.GetPrecedence(lexer.Exponent)
		Expect(exponent).To(BeNumerically(">", unaryAddition))
		Expect(p.NextPrecedence(unaryAddition)).To(Equal(exponent))

		Expect(p.MaxPrecedence()).To(Equal(exponent))
	})

	It("Returns default values when token meta are not set", func() {
		p := parser.TokenPriorities{}
		p[lexer.Addition] = parser.TokenMeta{Precedence: 20, Associativity: parser.RightAssociativity}

		Expect(p.GetAssociativity(lexer.Addition)).To(Equal(parser.RightAssociativity))
		Expect(p.GetPrecedence(lexer.Addition)).To(BeNumerically("==", 20))

		Expect(p.GetAssociativity(lexer.Exponent)).To(Equal(parser.LeftAssociativity))
		Expect(p.GetPrecedence(lexer.Exponent)).To(BeNumerically("==", 0))
	})

	It("Normalize - fails on zero precedence", func() {
		p := parser.TokenPriorities{}
		p[lexer.Addition] = parser.TokenMeta{Precedence: 0, Associativity: parser.RightAssociativity}
		Expect(p.Normalize()).To(MatchError(parser.ErrZeroPrecedenceSet))
	})

	It("Normalize - removes non operation priorities", func() {
		p := parser.DefaultTokenPriorities()

		p[lexer.LPar] = parser.TokenMeta{Precedence: 100}
		p[lexer.RPar] = parser.TokenMeta{Precedence: 100}
		p[lexer.Number] = parser.TokenMeta{Precedence: 100}
		p[lexer.Identifier] = parser.TokenMeta{Precedence: 100}
		p[lexer.Comma] = parser.TokenMeta{Precedence: 100}
		p[lexer.EOL] = parser.TokenMeta{Precedence: 100}
		p[lexer.Whitespace] = parser.TokenMeta{Precedence: 100}

		Expect(p).To(HaveLen(17))
		Expect(p.Normalize()).To(Succeed())
		Expect(p).To(HaveLen(10))
	})
})

var _ = DescribeTable("Check default associativity",
	func(tt lexer.TokenType, ta parser.TokenAssociativity) {
		Expect(parser.DefaultTokenPriorities().GetAssociativity(tt)).To(Equal(ta))
	},
	Entry("EOL", lexer.EOL, parser.LeftAssociativity),
	Entry("Whitespace", lexer.Whitespace, parser.LeftAssociativity),
	Entry("Addition", lexer.Addition, parser.LeftAssociativity),
	Entry("Substraction", lexer.Substraction, parser.LeftAssociativity),
	Entry("Multiplication", lexer.Multiplication, parser.LeftAssociativity),
	Entry("Division", lexer.Division, parser.LeftAssociativity),
	Entry("FloorDiv", lexer.FloorDiv, parser.LeftAssociativity),
	Entry("Modulus", lexer.Modulus, parser.LeftAssociativity),
	Entry("UnaryAddition", lexer.UnaryAddition, parser.LeftAssociativity),
	Entry("UnarySubstraction", lexer.UnarySubstraction, parser.LeftAssociativity),
	Entry("Exponent", lexer.Exponent, parser.RightAssociativity),
	Entry("LPar", lexer.LPar, parser.LeftAssociativity),
	Entry("RPar", lexer.RPar, parser.LeftAssociativity),
	Entry("Identifier", lexer.Identifier, parser.LeftAssociativity),
	Entry("Number", lexer.Number, parser.LeftAssociativity),
)
