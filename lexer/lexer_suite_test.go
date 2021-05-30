package lexer_test

import (
	"fmt"
	"testing"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"

	"github.com/arxeiss/go-expression-calculator/lexer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLexer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lexer Suite")
}

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

func (matcher *TokenMatcher) message(actual interface{}, negated bool) (message string) {
	if token, ok := actual.(lexer.Token); ok {
		notStr := ""
		if negated {
			notStr = "not "
		}
		return fmt.Sprintf("Expected {Type: %s, Value: %f, Identifier: '%s', Start pos: %d, End pos: %d}"+
			" %sto be equal to {Type: %s, Value: %f, Identifier: '%s', Start pos: %d, End pos: %d}",
			token.Type().String(), token.Value(), token.Identifier(), token.StartPosition(), token.EndPosition(),
			notStr,
			matcher.tType.String(), matcher.value, matcher.idName, matcher.startPos, matcher.endPos,
		)
	}
	return fmt.Sprintf("Expected %s to be equal to %s", format.Object(matcher, 0), format.Object(actual, 0))
}

func (matcher *TokenMatcher) FailureMessage(actual interface{}) (message string) {
	return matcher.message(actual, false)
}

func (matcher *TokenMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return matcher.message(actual, true)
}
