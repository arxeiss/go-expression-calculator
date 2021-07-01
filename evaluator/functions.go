package evaluator

import "math"

func MathFunctions() map[string]FunctionHandler {
	return map[string]FunctionHandler{
		"Abs": {
			Description: "Abs returns the absolute value of x.",
			Handler:     func(x ...float64) (float64, error) { return math.Abs(x[0]), nil },
		},
		"Acos": {
			Description: "Acos returns the arccosine, in radians, of x.",
			Handler:     func(x ...float64) (float64, error) { return math.Acos(x[0]), nil },
		},
		"Asin": {
			Description: "Asin returns the arcsine, in radians, of x.",
			Handler:     func(x ...float64) (float64, error) { return math.Asin(x[0]), nil },
		},
		"Atan": {
			Description: "Atan returns the arctangent, in radians, of x.",
			Handler:     func(x ...float64) (float64, error) { return math.Atan(x[0]), nil },
		},
		"Ceil": {
			Description: "Ceil returns the least integer value greater than or equal to x.",
			Handler:     func(x ...float64) (float64, error) { return math.Ceil(x[0]), nil },
		},
		"Cos": {
			Description: "Cos returns the cosine of the radian argument x.",
			Handler:     func(x ...float64) (float64, error) { return math.Cos(x[0]), nil },
		},
		"Floor": {
			Description: "Floor returns the greatest integer value less than or equal to x.",
			Handler:     func(x ...float64) (float64, error) { return math.Floor(x[0]), nil },
		},
		"Sin": {
			Description: "Sin returns the sine of the radian argument x.",
			Handler:     func(x ...float64) (float64, error) { return math.Sin(x[0]), nil },
		},
		"Sqrt": {
			Description: "Sqrt returns the square root of x.",
			Handler:     func(x ...float64) (float64, error) { return math.Sqrt(x[0]), nil },
		},
		"Tan": {
			Description: "Tan returns the tangent of the radian argument x.",
			Handler:     func(x ...float64) (float64, error) { return math.Tan(x[0]), nil },
		},
	}
}
