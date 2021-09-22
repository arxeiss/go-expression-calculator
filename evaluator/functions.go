package evaluator

import (
	"fmt"
	"math"
	"math/rand"
)

func MathFunctions() map[string]FunctionHandler {
	return map[string]FunctionHandler{
		"abs": {
			Description:  "Returns the absolute value of x.",
			Handler:      func(x ...float64) (float64, error) { return math.Abs(x[0]), nil },
			MinArguments: 1, MaxArguments: 1,
			ArgsNames: []string{"x"},
		},
		"acos": {
			Description:  "Returns the arccosine, in radians, of x.",
			Handler:      func(x ...float64) (float64, error) { return math.Acos(x[0]), nil },
			MinArguments: 1, MaxArguments: 1,
			ArgsNames: []string{"x"},
		},
		"asin": {
			Description:  "Returns the arcsine, in radians, of x.",
			Handler:      func(x ...float64) (float64, error) { return math.Asin(x[0]), nil },
			MinArguments: 1, MaxArguments: 1,
			ArgsNames: []string{"x"},
		},
		"atan": {
			Description:  "Returns the arctangent, in radians, of x.",
			Handler:      func(x ...float64) (float64, error) { return math.Atan(x[0]), nil },
			MinArguments: 1, MaxArguments: 1,
			ArgsNames: []string{"x"},
		},
		"ceil": {
			Description:  "Returns the least integer value greater than or equal to x.",
			Handler:      func(x ...float64) (float64, error) { return math.Ceil(x[0]), nil },
			MinArguments: 1, MaxArguments: 1,
			ArgsNames: []string{"x"},
		},
		"cos": {
			Description:  "Returns the cosine of the radian argument x.",
			Handler:      func(x ...float64) (float64, error) { return math.Cos(x[0]), nil },
			MinArguments: 1, MaxArguments: 1,
			ArgsNames: []string{"x"},
		},
		"floor": {
			Description:  "Returns the greatest integer value less than or equal to x.",
			Handler:      func(x ...float64) (float64, error) { return math.Floor(x[0]), nil },
			MinArguments: 1, MaxArguments: 1,
			ArgsNames: []string{"x"},
		},
		"sin": {
			Description:  "Returns the sine of the radian argument x.",
			Handler:      func(x ...float64) (float64, error) { return math.Sin(x[0]), nil },
			MinArguments: 1, MaxArguments: 1,
			ArgsNames: []string{"x"},
		},
		"sqrt": {
			Description:  "Returns the square root of x.",
			Handler:      func(x ...float64) (float64, error) { return math.Sqrt(x[0]), nil },
			MinArguments: 1, MaxArguments: 1,
			ArgsNames: []string{"x"},
		},
		"tan": {
			Description:  "Returns the tangent of the radian argument x.",
			Handler:      func(x ...float64) (float64, error) { return math.Tan(x[0]), nil },
			MinArguments: 1, MaxArguments: 1,
			ArgsNames: []string{"x"},
		},
		"deg2rad": {
			Description:  "Convert x from degrees into radians.",
			Handler:      func(x ...float64) (float64, error) { return x[0] * (math.Pi / 180), nil },
			MinArguments: 1, MaxArguments: 1,
			ArgsNames: []string{"x"},
		},
		"rad2deg": {
			Description:  "Convert x from radians into degrees.",
			Handler:      func(x ...float64) (float64, error) { return x[0] * (180 / math.Pi), nil },
			MinArguments: 1, MaxArguments: 1,
			ArgsNames: []string{"x"},
		},
	}
}

func MathFunctionsWithVarArgs() map[string]FunctionHandler {
	return map[string]FunctionHandler{
		"pi": {
			Description:  "Returns Pi value.",
			Handler:      func(x ...float64) (float64, error) { return math.Pi, nil },
			MinArguments: 0, MaxArguments: 0,
		},
		"e": {
			Description:  "Returns e value (base of natural logarithm).",
			Handler:      func(x ...float64) (float64, error) { return math.E, nil },
			MinArguments: 0, MaxArguments: 0,
		},
		"phi": {
			Description:  "Returns Phi value.",
			Handler:      func(x ...float64) (float64, error) { return math.Phi, nil },
			MinArguments: 0, MaxArguments: 0,
		},
		"log": {
			Description: "Returns log of value n with given base.",
			// Log_10(20) == ln(20) / ln(10)
			Handler:      func(x ...float64) (float64, error) { return math.Log(x[0]) / math.Log(x[1]), nil },
			MinArguments: 2, MaxArguments: 2,
			ArgsNames: []string{"n", "base"},
		},
		"max": {
			Description: "Returns maximum of provided numbers.",
			Handler: func(x ...float64) (float64, error) {
				c := x[0]
				for i := 1; i < len(x); i++ {
					c = math.Max(c, x[i])
				}
				return c, nil
			},
			MinArguments: 1, MaxArguments: 0,
			ArgsNames: []string{"a", "b"},
		},
		"min": {
			Description: "Returns minimum of provided numbers.",
			Handler: func(x ...float64) (float64, error) {
				c := x[0]
				for i := 1; i < len(x); i++ {
					c = math.Min(c, x[i])
				}
				return c, nil
			},
			MinArguments: 1, MaxArguments: 0,
			ArgsNames: []string{"a", "b"},
		},
		"rand_f": {
			Description: "Returns random float number in range <0;1).",
			Handler: func(x ...float64) (float64, error) {
				return rand.Float64(), nil
			},
			MinArguments: 0, MaxArguments: 0,
		},
		"rand_i": {
			Description: "Returns random decimal number in range <0, a) or <a, b) if b is provided.",
			Handler: func(x ...float64) (float64, error) {
				if len(x) == 1 {
					x = append(x, x[0])
					x[0] = 0
				}
				min, max := int64(x[0]), int64(x[1])
				if min >= max {
					return 0, fmt.Errorf("number %d (min) cannot be higher or equal to %d (max)", min, max)
				}
				return float64(rand.Int63n(max-min) + min), nil
			},
			MinArguments: 1, MaxArguments: 2,
			ArgsNames: []string{"a", "b"},
		},
		"nth_root": {
			Description: "Returns n-th root of a.",
			Handler: func(p ...float64) (float64, error) {
				a := p[0]
				n := p[1]
				if a < 0 {
					return 0, fmt.Errorf("number a cannot be negative")
				}
				if n <= 0 {
					return 0, fmt.Errorf("number n cannot be 0 or negative")
				}

				// If n has fractional part
				if n-math.Floor(n) > 1e-15 {
					return math.Pow(a, 1/n), nil
				}

				// https://rosettacode.org/wiki/Nth_root#Go
				n1 := int64(n) - 1
				n1f, rn := float64(n1), 1/math.Floor(n) // n is not int64 but already float64
				x, x0 := 1.0, 0.0
				for {
					potx, t2 := 1/x, a
					for b := n1; b > 0; b >>= 1 {
						if b&1 == 1 {
							t2 *= potx
						}
						potx *= potx
					}
					x0, x = x, rn*(n1f*x+t2)
					if math.Abs(x-x0)*1e15 < x {
						break
					}
				}
				return x, nil
			},
			MinArguments: 2, MaxArguments: 2,
			ArgsNames: []string{"a", "n"},
		},
	}
}
