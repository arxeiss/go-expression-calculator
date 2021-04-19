package lexer_test

import (
	"github.com/onsi/gomega/types"

	"github.com/arxeiss/go-expression-calculator/lexer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Tokenize", func() {
	It("Handle all tokens", func() {
		allTokensInput := "(**^s0m3_Identifier)//  %+-*159.4587e-5/" // #nosec G101
		l := lexer.NewLexer(allTokensInput)
		Expect(l.Expression()).To(Equal(allTokensInput))
		tokens, err := l.Tokenize()
		Expect(err).To(Succeed())
		Expect(tokens).To(MatchAllElementsWithIndex(IndexIdentity, Elements{
			"0":  PointTo(MatchToken(lexer.LPar, 0, "", 0, 1)),
			"1":  PointTo(MatchToken(lexer.Exponent, 0, "", 1, 3)),
			"2":  PointTo(MatchToken(lexer.Exponent, 0, "", 3, 4)),
			"3":  PointTo(MatchToken(lexer.Identifier, 0, "s0m3_Identifier", 4, 19)),
			"4":  PointTo(MatchToken(lexer.RPar, 0, "", 19, 20)),
			"5":  PointTo(MatchToken(lexer.FloorDiv, 0, "", 20, 22)),
			"6":  PointTo(MatchToken(lexer.Whitespace, 0, "", 22, 24)),
			"7":  PointTo(MatchToken(lexer.Modulus, 0, "", 24, 25)),
			"8":  PointTo(MatchToken(lexer.Addition, 0, "", 25, 26)),
			"9":  PointTo(MatchToken(lexer.Substraction, 0, "", 26, 27)),
			"10": PointTo(MatchToken(lexer.Multiplication, 0, "", 27, 28)),
			"11": PointTo(MatchToken(lexer.Number, 159.4587e-5, "", 28, 39)),
			"12": PointTo(MatchToken(lexer.Division, 0, "", 39, 40)),
		}))
	})

	DescribeTable("Handle valid numbers",
		func(expr string, valueMatcher types.GomegaMatcher) {
			l := lexer.NewLexer(expr)
			tokens, err := l.Tokenize()
			Expect(err).To(Succeed())
			Expect(tokens).To(HaveLen(1))
			Expect(tokens[0].Value()).To(valueMatcher)
		},
		Entry("Whole numbers only", "123", BeEquivalentTo(123)),
		Entry("Decimal", "876.93", BeEquivalentTo(876.93)),
		Entry("Fraction part only", ".72", BeEquivalentTo(0.72)),

		Entry("Whole numbers only with exponent", "123e2", BeEquivalentTo(12300)),
		Entry("Decimal with exponent", "876.93e3", BeEquivalentTo(876930)),
		Entry("Fraction part only with exponent", ".72e5", BeEquivalentTo(72000)),

		Entry("Decimal with positive exponent", ".047e+5", BeEquivalentTo(4700)),
		Entry("Decimal with negative exponent", "4.7e-5", BeEquivalentTo(0.000047)),
	)

	DescribeTable("Handle invalid character error",
		func(expr string, pos int, errStr string, wrapperErr error) {
			l := lexer.NewLexer(expr)
			tokens, err := l.Tokenize()
			Expect(tokens).To(BeNil())
			lexErr := err.(*lexer.Error)
			Expect(lexErr.Position()).To(Equal(pos))
			Expect(lexErr.Error()).To(Equal(errStr))
			Expect(lexErr.Unwrap()).To(Equal(wrapperErr))
		},
		Entry("At the begining", "? 123", 0, "unexpected character at position 0", lexer.ErrUnexpectedChar),
		Entry("In the middle", "+ < 123", 2, "unexpected character at position 2", lexer.ErrUnexpectedChar),
		Entry("At the end", "123.", 3, "unexpected character at position 3", lexer.ErrUnexpectedChar),
	)

	It("Handle empty error", func() {
		err := lexer.Error{}
		Expect(err.Position()).To(Equal(-1))
		Expect(err.Error()).To(Equal("unexpected error"))
		Expect(err.Unwrap()).To(BeNil())
	})
})
