package evaluator

import "math"

func MathFunctions() map[string]FunctionHandler {
	return map[string]FunctionHandler{
		"Abs":   func(x ...float64) (float64, error) { return math.Abs(x[0]), nil },
		"Acos":  func(x ...float64) (float64, error) { return math.Acos(x[0]), nil },
		"Asin":  func(x ...float64) (float64, error) { return math.Asin(x[0]), nil },
		"Atan":  func(x ...float64) (float64, error) { return math.Atan(x[0]), nil },
		"Ceil":  func(x ...float64) (float64, error) { return math.Ceil(x[0]), nil },
		"Cos":   func(x ...float64) (float64, error) { return math.Cos(x[0]), nil },
		"Floor": func(x ...float64) (float64, error) { return math.Floor(x[0]), nil },
		"Sin":   func(x ...float64) (float64, error) { return math.Sin(x[0]), nil },
		"Sqrt":  func(x ...float64) (float64, error) { return math.Sqrt(x[0]), nil },
		"Tan":   func(x ...float64) (float64, error) { return math.Tan(x[0]), nil },
	}
}
