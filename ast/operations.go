package ast

var (
	operationsStr = []string{"Invalid", "+", "-", "*", "/", "^", "//", "%", "="}
)

type Operation uint8

const (
	Invalid Operation = iota

	Addition
	Substraction
	Multiplication
	Division

	Exponent
	FloorDiv
	Modulus
	Assign
)

func (o Operation) String() string {
	return operationsStr[o]
}
