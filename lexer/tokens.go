package lexer

var (
	tokenTypeStr = []string{
		"EOL", "Whitespace", "Identifier", "LPar", "RPar", "Exponent", "Multiplication", "Division", "FloorDiv",
		"Modulus", "Addition", "Substraction", "Number", "Equal", "Comma", "UnaryAddition", "UnarySubstraction"}
)

type TokenType uint8

const (
	EOL TokenType = iota
	Whitespace

	Identifier
	LPar
	RPar
	Exponent
	Multiplication
	Division
	FloorDiv
	Modulus
	Addition
	Substraction
	Number
	Equal
	Comma

	// Unary operators cannot be recognized by lexer, but are prepared for parsers

	UnaryAddition
	UnarySubstraction
)

func (tt TokenType) String() string {
	return tokenTypeStr[tt]
}

type Token struct {
	tType            TokenType
	value            float64
	idName           string
	startPos, endPos int
}

func NewToken(tType TokenType, value float64, idName string, startPos, endPos int) *Token {
	return &Token{
		tType:    tType,
		value:    value,
		idName:   idName,
		startPos: startPos,
		endPos:   endPos,
	}
}

func (t *Token) Type() TokenType {
	if t == nil {
		return EOL
	}
	return t.tType
}

func (t *Token) Value() float64 {
	if t == nil {
		return 0
	}
	return t.value
}

func (t *Token) Identifier() string {
	if t == nil {
		return ""
	}
	return t.idName
}

func (t *Token) StartPosition() int {
	if t == nil {
		return 0
	}
	return t.startPos
}

func (t *Token) EndPosition() int {
	if t == nil {
		return 0
	}
	return t.endPos
}

func (t *Token) ChangeToUnary() error {
	if t == nil {
		return ErrInvalidUnary
	}
	switch t.tType {
	case Addition, UnaryAddition:
		t.tType = UnaryAddition
	case Substraction, UnarySubstraction:
		t.tType = UnarySubstraction
	default:
		return ErrInvalidUnary
	}
	return nil
}
