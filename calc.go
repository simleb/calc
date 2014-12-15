// Package calc parses and evaluates mathematical expressions.
//
// Expression are first tokenized (lexical analysis), then rearranged
// into reverse polish notation (RPN) using the shunting-yard algorithm.
// Finally, they are evaluated by treating all numbers as floating point.
// Expressions might contain custom functions and variables.
//
// This package is designed for short expressions,
// hence the input is a string and not an io.Reader.
package calc

// EvalFloat parses and evaluates a mathematical expression and returns a float64.
// All numbers in the expression are treated as floating point numbers.
// A map of variables (int or float64) and functions (taking any number
// of arguments and returning a float64) can be provided.
func EvalFloat(exp string, vars map[string]interface{}) (float64, error) {
	tokens, err := tokenize(exp)
	if err != nil {
		return 0, err
	}
	rpn, err := postfix(tokens, vars)
	if err != nil {
		return 0, err
	}
	return eval(rpn, vars)
}

// EvalInt parses and evaluates a mathematical expression and returns an int.
// The evaluation uses floating point numbers but the result is truncated.
func EvalInt(exp string, vars map[string]interface{}) (int, error) {
	x, err := EvalFloat(exp, vars)
	return int(x), err
}
