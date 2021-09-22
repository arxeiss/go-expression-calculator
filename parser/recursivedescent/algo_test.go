package recursivedescent_test

import (
	"github.com/onsi/gomega/types"

	"github.com/arxeiss/go-expression-calculator/ast"
	. "github.com/arxeiss/go-expression-calculator/ast/astutils"
	"github.com/arxeiss/go-expression-calculator/lexer"
	"github.com/arxeiss/go-expression-calculator/parser"
	"github.com/arxeiss/go-expression-calculator/parser/recursivedescent"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Successful parsing", func() {
	It("All supported token types", func() {
		p, err := recursivedescent.NewParser(parser.DefaultTokenPriorities())
		Expect(err).To(Succeed())
		input := []*lexer.Token{
			/* n */ lexer.NewToken(lexer.Identifier, 0, "n", 0, 0),
			/* = */ lexer.NewToken(lexer.Equal, 0, "", 0, 0),
			/* ( */ lexer.NewToken(lexer.LPar, 0, "", 0, 0),
			/* 30 */ lexer.NewToken(lexer.Number, 30, "", 0, 0),
			/* + */ lexer.NewToken(lexer.Addition, 0, "", 0, 0),
			/* 10 */ lexer.NewToken(lexer.Number, 10, "", 0, 0),
			/* * */ lexer.NewToken(lexer.Multiplication, 0, "", 0, 0),
			/* abc */ lexer.NewToken(lexer.Identifier, 0, "abc", 0, 0),
			/* ) */ lexer.NewToken(lexer.RPar, 0, "", 0, 0),
			/* - */ lexer.NewToken(lexer.Substraction, 0, "", 0, 0),
			/* 10 */ lexer.NewToken(lexer.Number, 10, "", 0, 0),
			/* / */ lexer.NewToken(lexer.Division, 0, "", 0, 0),
			/* 3 */ lexer.NewToken(lexer.Number, 3, "", 0, 0),
			/* // */ lexer.NewToken(lexer.FloorDiv, 0, "", 0, 0),
			/* 2 */ lexer.NewToken(lexer.Number, 2, "", 0, 0),
			/* + */ lexer.NewToken(lexer.Addition, 0, "", 0, 0),
			/* - */ lexer.NewToken(lexer.Substraction, 0, "", 0, 0), // Will be changed to unary
			/* max */ lexer.NewToken(lexer.Identifier, 0, "max", 0, 0),
			/* ( */ lexer.NewToken(lexer.LPar, 0, "", 0, 0),
			/* 5 */ lexer.NewToken(lexer.Number, 5, "", 0, 0),
			/* , */ lexer.NewToken(lexer.Comma, 0, "", 0, 0),
			/* 17 */ lexer.NewToken(lexer.Number, 17, "", 0, 0),
			/* % */ lexer.NewToken(lexer.Modulus, 0, "", 0, 0),
			/* 7 */ lexer.NewToken(lexer.Number, 7, "", 0, 0),
			/* ) */ lexer.NewToken(lexer.RPar, 0, "", 0, 0),

			lexer.NewToken(lexer.EOL, 0, "", 0, 0),
		}
		rootNode, err := p.Parse(input)
		Expect(err).To(Succeed())
		Expect(rootNode).NotTo(BeNil())

		Expect(rootNode).To(MatchAssignNode(
			MatchVariableNode("n"),
			MatchBinaryNode(
				ast.Addition,
				MatchBinaryNode(
					ast.Substraction,
					MatchBinaryNode(
						ast.Addition,
						MatchNumericNode(30),
						MatchBinaryNode(
							ast.Multiplication,
							MatchNumericNode(10),
							MatchVariableNode("abc"),
						),
					),
					MatchBinaryNode(
						ast.FloorDiv,
						MatchBinaryNode(
							ast.Division,
							MatchNumericNode(10),
							MatchNumericNode(3),
						),
						MatchNumericNode(2),
					),
				),
				MatchUnaryNode(
					ast.Substraction,
					MatchFunctionNode("max",
						MatchNumericNode(5),
						MatchBinaryNode(
							ast.Modulus,
							MatchNumericNode(17),
							MatchNumericNode(7),
						),
					),
				),
			),
		))
		// Just test it will not panic
		Expect(ast.ToTreeDrawer(rootNode)).NotTo(BeNil())
	})

	It("Simple expression with parenthesis", func() {
		p, err := recursivedescent.NewParser(parser.DefaultTokenPriorities())
		Expect(err).To(Succeed())
		input := []*lexer.Token{
			lexer.NewToken(lexer.LPar, 0, "", 0, 0),
			lexer.NewToken(lexer.Number, 23, "", 0, 0),
			lexer.NewToken(lexer.Addition, 0, "", 0, 0),
			lexer.NewToken(lexer.Number, 11, "", 0, 0),
			lexer.NewToken(lexer.RPar, 0, "", 0, 0),
			lexer.NewToken(lexer.Division, 0, "", 0, 0),
			lexer.NewToken(lexer.Number, 4, "", 0, 0),
			lexer.NewToken(lexer.Substraction, 0, "", 0, 0),
			lexer.NewToken(lexer.Number, 25, "", 0, 0),
			lexer.NewToken(lexer.Multiplication, 0, "", 0, 0),
			lexer.NewToken(lexer.Number, 10, "", 0, 0),
			lexer.NewToken(lexer.EOL, 0, "", 0, 0),
		}
		rootNode, err := p.Parse(input)
		Expect(err).To(Succeed())
		Expect(rootNode).NotTo(BeNil())

		Expect(rootNode).To(MatchBinaryNode(
			ast.Substraction,
			MatchBinaryNode(
				ast.Division,
				MatchBinaryNode(
					ast.Addition,
					MatchNumericNode(23),
					MatchNumericNode(11),
				),
				MatchNumericNode(4),
			),
			MatchBinaryNode(
				ast.Multiplication,
				MatchNumericNode(25),
				MatchNumericNode(10),
			),
		))
	})

	It("Starts with unary substraction and number", func() {
		p, err := recursivedescent.NewParser(parser.DefaultTokenPriorities())
		Expect(err).To(Succeed())
		input := []*lexer.Token{
			lexer.NewToken(lexer.Substraction, 0, "", 0, 0),
			lexer.NewToken(lexer.Number, 46, "", 0, 0),
			lexer.NewToken(lexer.EOL, 0, "", 0, 0),
		}
		rootNode, err := p.Parse(input)
		Expect(err).To(Succeed())
		Expect(rootNode).NotTo(BeNil())

		Expect(rootNode).To(MatchUnaryNode(
			ast.Substraction,
			MatchNumericNode(46),
		))
	})

	It("Simple multiple unary expressions", func() {
		p, err := recursivedescent.NewParser(parser.DefaultTokenPriorities())
		Expect(err).To(Succeed())
		input := []*lexer.Token{
			lexer.NewToken(lexer.Number, 20, "", 0, 0),
			lexer.NewToken(lexer.Substraction, 0, "", 0, 0),
			lexer.NewToken(lexer.Substraction, 0, "", 0, 0),
			lexer.NewToken(lexer.Substraction, 0, "", 0, 0),
			lexer.NewToken(lexer.Substraction, 0, "", 0, 0),
			lexer.NewToken(lexer.Number, 5, "", 0, 0),
			lexer.NewToken(lexer.Addition, 0, "", 0, 0),
			lexer.NewToken(lexer.Addition, 0, "", 0, 0),
			lexer.NewToken(lexer.Number, 7, "", 0, 0),
			lexer.NewToken(lexer.EOL, 0, "", 0, 0),
		}
		rootNode, err := p.Parse(input)
		Expect(err).To(Succeed())
		Expect(rootNode).NotTo(BeNil())

		Expect(rootNode).To(MatchBinaryNode(
			ast.Addition,
			MatchBinaryNode(
				ast.Substraction,
				MatchNumericNode(20),
				MatchUnaryNode(
					ast.Substraction,
					MatchUnaryNode(
						ast.Substraction,
						MatchUnaryNode(
							ast.Substraction,
							MatchNumericNode(5),
						),
					),
				),
			),
			MatchUnaryNode(ast.Addition, MatchNumericNode(7)),
		))
	})

	It("Simple unary with higher precedence operator", func() {
		p, err := recursivedescent.NewParser(parser.DefaultTokenPriorities())
		Expect(err).To(Succeed())
		input := []*lexer.Token{
			lexer.NewToken(lexer.Substraction, 0, "", 0, 0),
			lexer.NewToken(lexer.Number, 5, "", 0, 0),
			lexer.NewToken(lexer.Exponent, 0, "", 0, 0),
			lexer.NewToken(lexer.Substraction, 0, "", 0, 0),
			lexer.NewToken(lexer.Number, 2, "", 0, 0),
			lexer.NewToken(lexer.EOL, 0, "", 0, 0),
		}
		rootNode, err := p.Parse(input)
		Expect(err).To(Succeed())
		Expect(rootNode).NotTo(BeNil())

		Expect(rootNode).To(MatchUnaryNode(
			ast.Substraction,
			MatchBinaryNode(
				ast.Exponent,
				MatchNumericNode(5),
				MatchUnaryNode(ast.Substraction, MatchNumericNode(2)),
			),
		))
	})

	It("Nested unary operators in parenthesis", func() {
		p, err := recursivedescent.NewParser(parser.DefaultTokenPriorities())
		Expect(err).To(Succeed())
		input := []*lexer.Token{
			lexer.NewToken(lexer.Addition, 0, "", 0, 0),
			lexer.NewToken(lexer.LPar, 0, "", 0, 0),
			lexer.NewToken(lexer.Substraction, 0, "", 0, 0),
			lexer.NewToken(lexer.Identifier, 0, "abs", 0, 0),
			lexer.NewToken(lexer.LPar, 0, "", 0, 0),
			lexer.NewToken(lexer.Substraction, 0, "", 0, 0),
			lexer.NewToken(lexer.Identifier, 0, "jkl", 0, 0),
			lexer.NewToken(lexer.RPar, 0, "", 0, 0),
			lexer.NewToken(lexer.RPar, 0, "", 0, 0),
			lexer.NewToken(lexer.EOL, 0, "", 0, 0),
		}
		rootNode, err := p.Parse(input)
		Expect(err).To(Succeed())
		Expect(rootNode).NotTo(BeNil())

		Expect(rootNode).To(MatchUnaryNode(
			ast.Addition,
			MatchUnaryNode(
				ast.Substraction,
				MatchFunctionNode("abs",
					MatchUnaryNode(
						ast.Substraction,
						MatchVariableNode("jkl"),
					),
				),
			),
		))
	})

	It("Handles unary operators in a row", func() {
		p, err := recursivedescent.NewParser(parser.DefaultTokenPriorities())
		Expect(err).To(Succeed())
		input := []*lexer.Token{
			lexer.NewToken(lexer.Number, 11, "", 0, 0),
			lexer.NewToken(lexer.Addition, 0, "", 0, 0),
			lexer.NewToken(lexer.Addition, 0, "", 0, 0),
			lexer.NewToken(lexer.LPar, 0, "", 0, 0),
			lexer.NewToken(lexer.Identifier, 0, "x", 0, 0),
			lexer.NewToken(lexer.Substraction, 0, "", 0, 0),
			lexer.NewToken(lexer.Substraction, 0, "", 0, 0),
			lexer.NewToken(lexer.Number, 150, "", 0, 0),
			lexer.NewToken(lexer.RPar, 0, "", 0, 0),
			lexer.NewToken(lexer.Substraction, 0, "", 0, 0),
			lexer.NewToken(lexer.Addition, 0, "", 0, 0),
			lexer.NewToken(lexer.Identifier, 0, "y", 0, 0),

			lexer.NewToken(lexer.EOL, 0, "", 0, 0),
		}
		rootNode, err := p.Parse(input)
		Expect(err).To(Succeed())
		Expect(rootNode).NotTo(BeNil())

		Expect(rootNode).To(MatchBinaryNode(
			ast.Substraction,
			MatchBinaryNode(
				ast.Addition,
				MatchNumericNode(11),
				MatchUnaryNode(
					ast.Addition,
					MatchBinaryNode(
						ast.Substraction,
						MatchVariableNode("x"),
						MatchUnaryNode(ast.Substraction, MatchNumericNode(150)),
					),
				),
			),
			MatchUnaryNode(ast.Addition, MatchVariableNode("y")),
		))
	})

	It("Correctly handles right associativity", func() {
		p, err := recursivedescent.NewParser(parser.DefaultTokenPriorities())
		Expect(err).To(Succeed())
		input := []*lexer.Token{
			lexer.NewToken(lexer.Identifier, 0, "i", 0, 0),
			lexer.NewToken(lexer.Exponent, 0, "", 0, 0),
			lexer.NewToken(lexer.LPar, 0, "", 0, 0),
			lexer.NewToken(lexer.Identifier, 0, "x", 0, 0),
			lexer.NewToken(lexer.Exponent, 0, "", 0, 0),
			lexer.NewToken(lexer.Identifier, 0, "y", 0, 0),
			lexer.NewToken(lexer.RPar, 0, "", 0, 0),
			lexer.NewToken(lexer.Exponent, 0, "", 0, 0),
			lexer.NewToken(lexer.Identifier, 0, "j", 0, 0),
			lexer.NewToken(lexer.Exponent, 0, "", 0, 0),
			lexer.NewToken(lexer.Identifier, 0, "k", 0, 0),
			lexer.NewToken(lexer.EOL, 0, "", 0, 0),
		}
		rootNode, err := p.Parse(input)
		Expect(err).To(Succeed())
		Expect(rootNode).NotTo(BeNil())

		Expect(rootNode).To(MatchBinaryNode(
			ast.Exponent,
			MatchVariableNode("i"),
			MatchBinaryNode(
				ast.Exponent,
				MatchBinaryNode(
					ast.Exponent,
					MatchVariableNode("x"),
					MatchVariableNode("y"),
				), MatchBinaryNode(
					ast.Exponent,
					MatchVariableNode("j"),
					MatchVariableNode("k"),
				),
			),
		))
	})

	It("Support assign only number", func() {
		p, err := recursivedescent.NewParser(parser.DefaultTokenPriorities())
		Expect(err).To(Succeed())
		input := []*lexer.Token{
			lexer.NewToken(lexer.Identifier, 0, "a", 0, 0),
			lexer.NewToken(lexer.Equal, 0, "", 0, 0),
			lexer.NewToken(lexer.Number, 46, "", 0, 0),
			lexer.NewToken(lexer.EOL, 0, "", 0, 0),
		}
		rootNode, err := p.Parse(input)
		Expect(err).To(Succeed())
		Expect(rootNode).NotTo(BeNil())

		Expect(rootNode).To(MatchAssignNode(
			MatchVariableNode("a"),
			MatchNumericNode(46),
		))
	})

	It("Support assign value of another variable", func() {
		p, err := recursivedescent.NewParser(parser.DefaultTokenPriorities())
		Expect(err).To(Succeed())
		input := []*lexer.Token{
			lexer.NewToken(lexer.Identifier, 0, "a", 0, 0),
			lexer.NewToken(lexer.Equal, 0, "", 0, 0),
			lexer.NewToken(lexer.Identifier, 0, "b", 0, 0),
			lexer.NewToken(lexer.EOL, 0, "", 0, 0),
		}
		rootNode, err := p.Parse(input)
		Expect(err).To(Succeed())
		Expect(rootNode).NotTo(BeNil())

		Expect(rootNode).To(MatchAssignNode(
			MatchVariableNode("a"),
			MatchVariableNode("b"),
		))
	})

	It("Support functions without arguments", func() {
		p, err := recursivedescent.NewParser(parser.DefaultTokenPriorities())
		Expect(err).To(Succeed())
		input := []*lexer.Token{
			lexer.NewToken(lexer.Identifier, 0, "rand", 0, 0),
			lexer.NewToken(lexer.LPar, 0, "", 0, 0),
			lexer.NewToken(lexer.RPar, 46, "", 0, 0),
			lexer.NewToken(lexer.EOL, 0, "", 0, 0),
		}
		rootNode, err := p.Parse(input)
		Expect(err).To(Succeed())
		Expect(rootNode).NotTo(BeNil())

		Expect(rootNode).To(MatchFunctionNode("rand"))
	})

	It("Support functions with multiple arguments", func() {
		p, err := recursivedescent.NewParser(parser.DefaultTokenPriorities())
		Expect(err).To(Succeed())
		input := []*lexer.Token{
			// max(a, 123, b, -55---31^2)
			lexer.NewToken(lexer.Identifier, 0, "max", 0, 0),
			lexer.NewToken(lexer.LPar, 0, "", 0, 0),
			lexer.NewToken(lexer.Identifier, 0, "a", 0, 0),
			lexer.NewToken(lexer.Comma, 0, "", 0, 0),
			lexer.NewToken(lexer.Number, 123, "", 0, 0),
			lexer.NewToken(lexer.Comma, 0, "", 0, 0),
			lexer.NewToken(lexer.Identifier, 0, "b", 0, 0),
			lexer.NewToken(lexer.Comma, 0, "", 0, 0),
			lexer.NewToken(lexer.Substraction, 0, "", 0, 0),
			lexer.NewToken(lexer.Number, 55, "", 0, 0),
			lexer.NewToken(lexer.Substraction, 0, "", 0, 0),
			lexer.NewToken(lexer.Substraction, 0, "", 0, 0),
			lexer.NewToken(lexer.Substraction, 0, "", 0, 0),
			lexer.NewToken(lexer.Number, 31, "", 0, 0),
			lexer.NewToken(lexer.Exponent, 0, "", 0, 0),
			lexer.NewToken(lexer.Number, 2, "", 0, 0),
			lexer.NewToken(lexer.RPar, 0, "", 0, 0),
			lexer.NewToken(lexer.EOL, 0, "", 0, 0),
		}
		rootNode, err := p.Parse(input)
		Expect(err).To(Succeed())
		Expect(rootNode).NotTo(BeNil())

		Expect(rootNode).To(MatchFunctionNode(
			"max",
			MatchVariableNode("a"),
			MatchNumericNode(123),
			MatchVariableNode("b"),
			MatchBinaryNode(
				ast.Substraction,
				MatchUnaryNode(ast.Substraction, MatchNumericNode(55)),
				MatchUnaryNode(
					ast.Substraction,
					MatchUnaryNode(
						ast.Substraction,
						MatchBinaryNode(
							ast.Exponent,
							MatchNumericNode(31),
							MatchNumericNode(2),
						),
					),
				),
			),
		))
	})

	It("Support functions inside functions", func() {
		p, err := recursivedescent.NewParser(parser.DefaultTokenPriorities())
		Expect(err).To(Succeed())
		input := []*lexer.Token{
			/* c */ lexer.NewToken(lexer.Identifier, 0, "c", 0, 0),
			/* = */ lexer.NewToken(lexer.Equal, 0, "", 0, 0),
			/* min */ lexer.NewToken(lexer.Identifier, 0, "min", 0, 0),
			/* ( */ lexer.NewToken(lexer.LPar, 0, "", 0, 0),
			/* sin */ lexer.NewToken(lexer.Identifier, 0, "sin", 0, 0),
			/* ( */ lexer.NewToken(lexer.LPar, 0, "", 0, 0),
			/* 4 */ lexer.NewToken(lexer.Number, 4, "", 0, 0),
			/* * */ lexer.NewToken(lexer.Multiplication, 0, "", 0, 0),
			/* 7 */ lexer.NewToken(lexer.Number, 7, "", 0, 0),
			/* ) */ lexer.NewToken(lexer.RPar, 0, "", 0, 0),
			/* , */ lexer.NewToken(lexer.Comma, 0, "", 0, 0),
			/* cos */ lexer.NewToken(lexer.Identifier, 0, "cos", 0, 0),
			/* ( */ lexer.NewToken(lexer.LPar, 0, "", 0, 0),
			/* max */ lexer.NewToken(lexer.Identifier, 0, "max", 0, 0),
			/* ( */ lexer.NewToken(lexer.LPar, 0, "", 0, 0),
			/* a */ lexer.NewToken(lexer.Identifier, 0, "a", 0, 0),
			/* , */ lexer.NewToken(lexer.Comma, 0, "", 0, 0),
			/* sqrt */ lexer.NewToken(lexer.Identifier, 0, "sqrt", 0, 0),
			/* ( */ lexer.NewToken(lexer.LPar, 0, "", 0, 0),
			/* 18 */ lexer.NewToken(lexer.Number, 18, "", 0, 0),
			/* ) */ lexer.NewToken(lexer.RPar, 0, "", 0, 0),
			/* , */ lexer.NewToken(lexer.Comma, 0, "", 0, 0),
			/* sqrt */ lexer.NewToken(lexer.Identifier, 0, "sqrt", 0, 0),
			/* ( */ lexer.NewToken(lexer.LPar, 0, "", 0, 0),
			/* 25 */ lexer.NewToken(lexer.Number, 25, "", 0, 0),
			/* ) */ lexer.NewToken(lexer.RPar, 0, "", 0, 0),
			/* , */ lexer.NewToken(lexer.Comma, 0, "", 0, 0),
			/* c */ lexer.NewToken(lexer.Identifier, 0, "c", 0, 0),
			/* ) */ lexer.NewToken(lexer.RPar, 0, "", 0, 0),
			/* ) */ lexer.NewToken(lexer.RPar, 0, "", 0, 0),
			/* ) */ lexer.NewToken(lexer.RPar, 0, "", 0, 0),
			lexer.NewToken(lexer.EOL, 0, "", 0, 0),
		}

		rootNode, err := p.Parse(input)
		Expect(err).To(Succeed())
		Expect(rootNode).NotTo(BeNil())

		Expect(rootNode).To(MatchAssignNode(
			MatchVariableNode("c"),
			MatchFunctionNode(
				"min",
				MatchFunctionNode("sin", MatchBinaryNode(ast.Multiplication, MatchNumericNode(4), MatchNumericNode(7))),
				MatchFunctionNode("cos", MatchFunctionNode(
					"max",
					MatchVariableNode("a"),
					MatchFunctionNode("sqrt", MatchNumericNode(18)),
					MatchFunctionNode("sqrt", MatchNumericNode(25)),
					MatchVariableNode("c"),
				)),
			),
		))
	})
})

var _ = DescribeTable("Handle errors",
	func(input []*lexer.Token, posMatcher, errMatcher types.GomegaMatcher) {
		p, err := recursivedescent.NewParser(parser.DefaultTokenPriorities())
		Expect(err).To(Succeed())

		rootNode, err := p.Parse(input)
		Expect(rootNode).To(BeNil())

		parseErr, ok := err.(*parser.Error)
		Expect(ok).To(BeTrue())
		Expect(parseErr.Position()).To(posMatcher)
		Expect(parseErr.Error()).To(errMatcher)
	},
	Entry("Empty token list", []*lexer.Token{}, Equal(-1), ContainSubstring(recursivedescent.ErrEmptyInput.Error())),
	Entry("Empty token list", []*lexer.Token{
		lexer.NewToken(lexer.EOL, 0, "", 0, 0),
	}, Equal(0), ContainSubstring(recursivedescent.ErrEmptyInput.Error())),

	Entry("Not EOL at the end", []*lexer.Token{
		lexer.NewToken(lexer.Identifier, 0, "abc", 0, 3),
	}, Equal(0), ContainSubstring(recursivedescent.ErrExpectedEOL.Error())),

	Entry("Two number tokens in a row", []*lexer.Token{
		lexer.NewToken(lexer.Number, 20, "", 0, 2),
		lexer.NewToken(lexer.Whitespace, 0, " ", 2, 3),
		lexer.NewToken(lexer.Number, 20, "", 3, 5),
		lexer.NewToken(lexer.EOL, 0, "", 5, 5),
	}, Equal(3), ContainSubstring("unexpected token; found Number token at position 3")),

	Entry("First token must be operand",
		[]*lexer.Token{
			lexer.NewToken(lexer.Whitespace, 0, " ", 0, 1),
			lexer.NewToken(lexer.Division, 0, "", 1, 2),
			lexer.NewToken(lexer.EOL, 0, "", 2, 2),
		},
		Equal(1),
		ContainSubstring("expected number, identifier or left parenthesis; found Division token at position 1"),
	),
	Entry("Number followed by identifier", []*lexer.Token{
		lexer.NewToken(lexer.Number, 20, "", 0, 2),
		lexer.NewToken(lexer.Whitespace, 0, "  ", 3, 4),
		lexer.NewToken(lexer.Identifier, 0, "VariableName", 4, 16),
		lexer.NewToken(lexer.EOL, 0, "", 16, 16),
	}, Equal(4), ContainSubstring("unexpected token; found Identifier token at position 4")),

	Entry("Multiple unary operators in a row",
		[]*lexer.Token{
			lexer.NewToken(lexer.Substraction, 0, "", 0, 1),
			lexer.NewToken(lexer.Identifier, 0, "var", 2, 5),
			lexer.NewToken(lexer.Addition, 0, "", 5, 6),
			lexer.NewToken(lexer.Substraction, 0, "", 6, 7),
			lexer.NewToken(lexer.Substraction, 0, "", 7, 8),
			lexer.NewToken(lexer.EOL, 0, "", 8, 8),
		},
		Equal(8),
		ContainSubstring("expected number, identifier or left parenthesis; found EOL token at position 8"),
	),
	Entry("Only unary operators in a row",
		[]*lexer.Token{
			lexer.NewToken(lexer.Substraction, 0, "", 0, 1),
			lexer.NewToken(lexer.Substraction, 0, "", 1, 2),
			lexer.NewToken(lexer.EOL, 0, "", 2, 2),
		},
		Equal(2),
		ContainSubstring("expected number, identifier or left parenthesis; found EOL token at position 2"),
	),
	Entry("Missing right operand after unary",
		[]*lexer.Token{
			lexer.NewToken(lexer.Number, 10, "", 0, 2),
			lexer.NewToken(lexer.Substraction, 0, "", 2, 3),
			lexer.NewToken(lexer.Substraction, 0, "", 3, 4),
			lexer.NewToken(lexer.EOL, 0, "", 4, 4),
		},
		Equal(4),
		ContainSubstring("expected number, identifier or left parenthesis; found EOL token at position 4"),
	),
	Entry("Assign to number is not valid",
		[]*lexer.Token{
			lexer.NewToken(lexer.Number, 20, "", 0, 2),
			lexer.NewToken(lexer.Addition, 0, "", 2, 3),
			lexer.NewToken(lexer.Number, 77, "", 3, 5),
			lexer.NewToken(lexer.Equal, 0, "", 5, 6),
			lexer.NewToken(lexer.Number, 99, "", 3, 5),
			lexer.NewToken(lexer.EOL, 0, "", 6, 6),
		},
		Equal(5),
		ContainSubstring("expected one of ['Addition', 'Substraction', 'Multiplication', 'Division', 'FloorDiv', 'Modulus', 'Exponent'] types, got 'Equal'"), //nolint:lll
	),
	Entry("Assign to variable inside expression is not valid",
		[]*lexer.Token{
			lexer.NewToken(lexer.Number, 20, "", 0, 2),
			lexer.NewToken(lexer.Addition, 0, "", 2, 3),
			lexer.NewToken(lexer.LPar, 0, "", 3, 4),
			lexer.NewToken(lexer.Identifier, 0, "abc", 4, 7),
			lexer.NewToken(lexer.Equal, 0, "", 7, 8),
			lexer.NewToken(lexer.Number, 99, "", 8, 10),
			lexer.NewToken(lexer.RPar, 0, "", 10, 11),
			lexer.NewToken(lexer.EOL, 0, "", 11, 11),
		},
		Equal(7),
		ContainSubstring("expected one of ['Addition', 'Substraction', 'Multiplication', 'Division', 'FloorDiv', 'Modulus', 'Exponent'] types, got 'Equal'"), //nolint:lll
	),
	Entry("Unexpected operator after another operator - not unary",
		[]*lexer.Token{
			lexer.NewToken(lexer.Number, 20, "", 0, 2),
			lexer.NewToken(lexer.Addition, 0, "", 2, 3),
			lexer.NewToken(lexer.Division, 0, "", 3, 4),
			lexer.NewToken(lexer.Number, 20, "", 0, 2),
			lexer.NewToken(lexer.EOL, 0, "", 4, 4),
		},
		Equal(3),
		ContainSubstring("expected number, identifier or left parenthesis; found Division token at position 3"),
	),
	Entry("Found left parenthesis, expecting operator",
		[]*lexer.Token{
			lexer.NewToken(lexer.Number, 20, "", 0, 2),
			lexer.NewToken(lexer.LPar, 0, "", 2, 3),
			lexer.NewToken(lexer.RPar, 0, "", 3, 4),
			lexer.NewToken(lexer.EOL, 0, "", 4, 4),
		},
		Equal(2),
		ContainSubstring("unexpected token; found LPar token at position 2"),
	),
	Entry("Found right parenthesis, expecting operand",
		[]*lexer.Token{
			lexer.NewToken(lexer.LPar, 0, "", 0, 1),
			lexer.NewToken(lexer.Number, 20, "", 1, 3),
			lexer.NewToken(lexer.Exponent, 0, "", 3, 5),
			lexer.NewToken(lexer.RPar, 0, "", 5, 6),
			lexer.NewToken(lexer.EOL, 0, "", 6, 6),
		},
		Equal(5),
		ContainSubstring("expected number, identifier or left parenthesis; found RPar token at position 5"),
	),
	Entry("Found EOL, expecting operand",
		[]*lexer.Token{
			lexer.NewToken(lexer.Number, 20, "", 0, 2),
			lexer.NewToken(lexer.FloorDiv, 0, "", 2, 3),
			lexer.NewToken(lexer.EOL, 0, "", 3, 3),
		},
		Equal(3),
		ContainSubstring("expected number, identifier or left parenthesis; found EOL token at position 3"),
	),
	Entry("Extra left parenthesis",
		[]*lexer.Token{
			lexer.NewToken(lexer.Whitespace, 0, "    ", 0, 4),
			lexer.NewToken(lexer.LPar, 0, "", 4, 5),
			lexer.NewToken(lexer.LPar, 0, "", 5, 6),
			lexer.NewToken(lexer.Number, 123, "", 6, 9),
			lexer.NewToken(lexer.RPar, 0, "", 9, 10),
			lexer.NewToken(lexer.EOL, 0, "", 10, 10),
		},
		Equal(10),
		ContainSubstring("expected 'RPar' type, got 'EOL'; found EOL token at position 10"),
	),
	Entry("Extra right parenthesis",
		[]*lexer.Token{
			lexer.NewToken(lexer.Whitespace, 0, "  ", 0, 2),
			lexer.NewToken(lexer.LPar, 0, "", 2, 3),
			lexer.NewToken(lexer.Number, 123, "", 3, 6),
			lexer.NewToken(lexer.RPar, 0, "", 7, 8),
			lexer.NewToken(lexer.RPar, 0, "", 8, 9),
			lexer.NewToken(lexer.EOL, 0, "", 9, 9),
		},
		Equal(8),
		ContainSubstring("unexpected token; found RPar token at position 8"),
	),
)
