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

		eol := p.GetPrecedence(lexer.EOL)
		Expect(eol).To(BeNumerically(">=", 0))
		Expect(p.GetPrecedence(lexer.Whitespace)).To(BeNumerically("==", eol))

		addition := p.GetPrecedence(lexer.Addition)
		Expect(addition).To(BeNumerically(">", eol))
		Expect(p.GetPrecedence(lexer.Substraction)).To(BeNumerically("==", addition))

		multiplication := p.GetPrecedence(lexer.Multiplication)
		Expect(multiplication).To(BeNumerically(">", addition))
		Expect(p.GetPrecedence(lexer.Division)).To(BeNumerically("==", multiplication))
		Expect(p.GetPrecedence(lexer.FloorDiv)).To(BeNumerically("==", multiplication))
		Expect(p.GetPrecedence(lexer.Modulus)).To(BeNumerically("==", multiplication))

		unaryAddition := p.GetPrecedence(lexer.UnaryAddition)
		Expect(unaryAddition).To(BeNumerically(">", multiplication))
		Expect(p.GetPrecedence(lexer.UnarySubstraction)).To(BeNumerically("==", unaryAddition))

		exponent := p.GetPrecedence(lexer.Exponent)
		Expect(exponent).To(BeNumerically(">", unaryAddition))

		lpar := p.GetPrecedence(lexer.LPar)
		Expect(lpar).To(BeNumerically(">", exponent))
		Expect(p.GetPrecedence(lexer.RPar)).To(BeNumerically("==", lpar))

		identifier := p.GetPrecedence(lexer.Identifier)
		Expect(identifier).To(BeNumerically(">", lpar))
		Expect(p.GetPrecedence(lexer.Number)).To(BeNumerically("==", identifier))
	})

	It("Returns default values when token meta are not set", func() {
		p := parser.TokenPriorities{}
		p[lexer.Addition] = parser.TokenMeta{Precedence: 20, Associativity: parser.RightAssociativity}

		Expect(p.GetAssociativity(lexer.Addition)).To(Equal(parser.RightAssociativity))
		Expect(p.GetPrecedence(lexer.Addition)).To(Equal(20))

		Expect(p.GetAssociativity(lexer.Exponent)).To(Equal(parser.LeftAssociativity))
		Expect(p.GetPrecedence(lexer.Exponent)).To(Equal(0))
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
