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
	// TODO add functions with multiple args and equal operator
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
			/* sin */ lexer.NewToken(lexer.Identifier, 0, "sin", 0, 0),
			/* ( */ lexer.NewToken(lexer.LPar, 0, "", 0, 0),
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
					MatchFunctionNode("sin",
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

	// TODO: test assignment, function with multiple params
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

	Entry("Two number tokens in a row", []*lexer.Token{
		lexer.NewToken(lexer.Number, 20, "", 0, 2),
		lexer.NewToken(lexer.Whitespace, 0, " ", 2, 3),
		lexer.NewToken(lexer.Number, 20, "", 3, 5),
	}, Equal(3), ContainSubstring("expected 'EOL' type, got 'Number'; found Number token at position 3")),

	Entry("First token must be operand",
		[]*lexer.Token{
			lexer.NewToken(lexer.Whitespace, 0, " ", 0, 1),
			lexer.NewToken(lexer.Division, 0, "", 1, 2),
		},
		Equal(1),
		ContainSubstring("expected number, identifier or left parenthesis; found Division token at position 1"),
	),
	Entry("Number followed by identifier", []*lexer.Token{
		lexer.NewToken(lexer.Number, 20, "", 0, 2),
		lexer.NewToken(lexer.Whitespace, 0, "  ", 3, 4),
		lexer.NewToken(lexer.Identifier, 0, "VariableName", 4, 16),
	}, Equal(4), ContainSubstring("expected 'EOL' type, got 'Identifier'; found Identifier token at position 4")),

	FEntry("Multiple unary operators in a row",
		[]*lexer.Token{
			lexer.NewToken(lexer.Substraction, 0, "", 0, 1),
			lexer.NewToken(lexer.Identifier, 0, "var", 2, 5),
			lexer.NewToken(lexer.Addition, 0, "", 5, 6),
			lexer.NewToken(lexer.Substraction, 0, "", 6, 7),
			lexer.NewToken(lexer.Substraction, 0, "", 7, 8),
		},
		Equal(5),
		ContainSubstring("expected number, identifier or left parenthesis; found Addition token at position 5"),
	),
	XEntry("Unexpected operator after another operator - not unary",
		[]*lexer.Token{
			lexer.NewToken(lexer.Number, 20, "", 0, 2),
			lexer.NewToken(lexer.Addition, 0, "", 2, 3),
			lexer.NewToken(lexer.Division, 0, "", 3, 4),
		},
		Equal(3),
		ContainSubstring("expected number, identifier or left parenthesis; found Division token at position 3"),
	),
	XEntry("Found left parenthesis, expecting operator",
		[]*lexer.Token{
			lexer.NewToken(lexer.Number, 20, "", 0, 2),
			lexer.NewToken(lexer.LPar, 0, "", 2, 3),
		},
		Equal(2),
		ContainSubstring("expected operator or right parenthesis; found LPar token at position 2"),
	),
	XEntry("Found right parenthesis, expecting operand",
		[]*lexer.Token{
			lexer.NewToken(lexer.LPar, 0, "", 0, 1),
			lexer.NewToken(lexer.Number, 20, "", 1, 3),
			lexer.NewToken(lexer.Exponent, 0, "", 3, 5),
			lexer.NewToken(lexer.RPar, 0, "", 5, 6),
		},
		Equal(5),
		ContainSubstring("expected number, identifier or left parenthesis; found RPar token at position 5"),
	),
	XEntry("Found EOL, expecting operand",
		[]*lexer.Token{
			lexer.NewToken(lexer.Number, 20, "", 0, 2),
			lexer.NewToken(lexer.FloorDiv, 0, "", 2, 3),
			lexer.NewToken(lexer.EOL, 0, "", 3, 3),
		},
		Equal(3),
		ContainSubstring("unexpected end of input; found EOL token at position 3"),
	),
	XEntry("Extra left parenthesis",
		[]*lexer.Token{
			lexer.NewToken(lexer.Whitespace, 0, "    ", 0, 4),
			lexer.NewToken(lexer.LPar, 0, "", 4, 5),
			lexer.NewToken(lexer.LPar, 0, "", 5, 6),
			lexer.NewToken(lexer.Number, 123, "", 6, 9),
			lexer.NewToken(lexer.RPar, 0, "", 9, 10),
		},
		Equal(4),
		ContainSubstring("cannot find matching right parenthesis; found LPar token at position 4"),
	),
	XEntry("Extra right parenthesis",
		[]*lexer.Token{
			lexer.NewToken(lexer.Whitespace, 0, "  ", 0, 2),
			lexer.NewToken(lexer.LPar, 0, "", 2, 3),
			lexer.NewToken(lexer.Number, 123, "", 3, 6),
			lexer.NewToken(lexer.RPar, 0, "", 7, 8),
			lexer.NewToken(lexer.RPar, 0, "", 8, 9),
		},
		Equal(8),
		ContainSubstring("cannot find matching left parenthesis; found RPar token at position 8"),
	),
)
