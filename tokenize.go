package calc

import (
	"regexp"
	"strconv"
)

// A token is a unit of lexical analysis.
type token interface{}

// The following types are specialized tokens:

// A number is a floating point number token.
type number float64

// An ident is an identifier token (variable or function).
type ident string

// An operator is a binary operator token.
type operator rune

// A unary is a unary operator token.
type unary rune

// A parenOpen is a left parenthese token.
type parenOpen struct{}

// A parenClose is a right parenthese token.
type parenClose struct{}

// A separator is a token used to separate function arguments (comma).
type separator struct{}

// These regexps match number and ident tokens respectively.
var (
	numberPattern = regexp.MustCompile(`^(\d+(\.\d*)?|\.\d+?)([eE][-+]?\d+)?`)
	identPattern  = regexp.MustCompile(`^[a-zA-Z][\w.-]*`)
)

// tokenize splits an expression into tokens.
func tokenize(exp string) ([]token, error) {
	skip := 0
	var tokens []token
	for i, r := range exp {
		// Skip previously scanned runes
		if skip > 0 {
			skip--
			continue
		}

		// Skip whitespace
		if isSpace(r) {
			continue
		}

		// Scan identifiers
		if isLetter(r) {
			m := identPattern.FindString(exp[i:])
			tokens = append(tokens, ident(m))
			skip = len([]rune(m)) - 1
			continue
		}

		// Scan numbers
		if isDigit(r) || r == '.' {
			m := numberPattern.FindString(exp[i:])
			x, err := strconv.ParseFloat(m, 64)
			if err != nil {
				return nil, makeError(exp, i, "bad number (%v)", err)
			}
			tokens = append(tokens, number(x))
			skip = len([]rune(m)) - 1
			continue
		}

		// Scan operators
		if _, found := operators[operator(r)]; found {
			if r == '+' || r == '-' {
				if len(tokens) == 0 {
					tokens = append(tokens, unary(r))
					continue
				} else {
					switch tokens[len(tokens)-1].(type) {
					case parenOpen, separator, operator, unary:
						tokens = append(tokens, unary(r))
						continue
					}
				}
			}
			tokens = append(tokens, operator(r))
			continue
		}

		// Scan parentheses and separators
		switch r {
		case '(':
			tokens = append(tokens, parenOpen{})
			continue
		case ')':
			tokens = append(tokens, parenClose{})
			continue
		case ',':
			tokens = append(tokens, separator{})
			continue
		default:
			return nil, makeError(exp, i, "bad character '%c'", r)
		}
	}
	return tokens, nil
}

// isSpace check if the rune is an ASCII whitespace.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isLetter check if the rune is an ASCII letter.
func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// isDigit check if the rune is an ASCII digit.
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
