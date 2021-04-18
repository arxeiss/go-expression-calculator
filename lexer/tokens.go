package lexer

var (
	tokenTypeStr = []string{
		"EOL", "Unknown", "LPar", "RPar", "Exponent", "Multiplication", "Division", "FloorDiv", "Modulus",
		"Addition", "Substraction", "Number", "Identifier", "Whitespace"}
)

type TokenType uint8

const (
	EOL TokenType = iota
	Unknown
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
	Identifier
	Whitespace
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

func (t Token) Type() TokenType {
	return t.tType
}

func (t Token) Value() float64 {
	return t.value
}

func (t Token) Identifier() string {
	return t.idName
}

func (t Token) StartPosition() int {
	return t.startPos
}

func (t Token) EndPosition() int {
	return t.endPos
}
