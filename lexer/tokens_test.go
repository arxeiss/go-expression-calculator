package lexer_test

import (
	"fmt"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"

	"github.com/arxeiss/go-expression-calculator/lexer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

type TokenMatcher struct {
	tType            lexer.TokenType
	value            float64
	idName           string
	startPos, endPos int
}

func MatchToken(tType lexer.TokenType, value float64, idName string, startPos, endPos int) types.GomegaMatcher {
	return &TokenMatcher{
		tType:    tType,
		value:    value,
		idName:   idName,
		startPos: startPos,
		endPos:   endPos,
	}
}

func (matcher *TokenMatcher) Match(actual interface{}) (success bool, err error) {
	if token, ok := actual.(lexer.Token); ok {
		return token.Type() == matcher.tType &&
			token.Value() == matcher.value &&
			token.Identifier() == matcher.idName &&
			token.StartPosition() == matcher.startPos &&
			token.EndPosition() == matcher.endPos, nil
	}
	return false, fmt.Errorf("matcher MatchToken expects a `lexer.Token` Got:\n%s", format.Object(actual, 1))
}

func (matcher *TokenMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected %s to be equal to %s", format.Object(actual, 0), format.Object(actual, 0))
}

func (matcher *TokenMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected %s not to be equal to %s", format.Object(actual, 0), format.Object(actual, 0))
}

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
