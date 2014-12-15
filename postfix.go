package calc

import (
	"fmt"
	"reflect"
)

// Associativity.
const (
	left = iota
	right
)

// Precedence and associativity of operators.
var operators = map[token]struct {
	prec  int
	assoc int
}{
	operator('^'): {5, right},
	unary('+'):    {4, right},
	unary('-'):    {4, right},
	operator(' '): {3, left}, // implicit multiplication (U+202F NARROW NO-BREAK SPACE)
	operator('*'): {2, left},
	operator('/'): {2, left},
	operator('%'): {2, left},
	operator('+'): {1, left},
	operator('-'): {1, left},
}

// A stack is a dead simple stack of tokens.
type stack []token

// Push pushes a token on the stack.
func (s *stack) Push(t token) {
	*s = append(*s, t)
}

// Pop pops an element from the stack.
func (s *stack) Pop() token {
	n := len(*s)
	if n == 0 {
		return nil
	}
	t := (*s)[n-1]
	*s = (*s)[:n-1]
	return t
}

// Top returns the top element from the stack without popping.
func (s stack) Top() token {
	n := len(s)
	if n == 0 {
		return nil
	}
	return s[n-1]
}

// postfix rearranges tokens into RPN using the shunting-yard algorithm.
func postfix(tokens []token, vars map[string]interface{}) ([]token, error) {
	var out, stk stack
	for i, t := range tokens {
		switch t := t.(type) {
		case number:
			out.Push(t)
		case ident:
			if i > 0 {
				if _, ok := tokens[i-1].(number); ok {
					handleOperator(operator(' '), &stk, &out)
				}
			}
			v, found := vars[string(t)]
			if !found {
				return nil, fmt.Errorf("calc: '%v' not provided", t)
			}
			switch v.(type) {
			case int, float64:
				out.Push(t)
			default:
				stk.Push(t)
			}
		case separator:
			for {
				o := stk.Top()
				if o == nil {
					return nil, fmt.Errorf("calc: bad comma")
				}
				if _, ok := o.(parenOpen); ok {
					switch tokens[i-1].(type) {
					case separator, parenOpen:
						return nil, fmt.Errorf("calc: bad comma")
					}
					break
				}
				out.Push(stk.Pop())
			}
		case unary:
			stk.Push(t)
		case operator:
			handleOperator(t, &stk, &out)
		case parenOpen:
			if i > 0 {
				if _, ok := tokens[i-1].(number); ok {
					handleOperator(operator(' '), &stk, &out)
				}
			}
			stk.Push(t)
		case parenClose:
			switch tokens[i-1].(type) {
			case separator:
				return nil, fmt.Errorf("calc: bad comma")
			}
			for {
				o := stk.Pop()
				if o == nil {
					return nil, fmt.Errorf("calc: mismatched parentheses")
				}
				if _, ok := o.(parenOpen); ok {
					break
				}
				out.Push(o)
			}
			// Check if top of stack contains a function
			o := stk.Top()
			if o == nil {
				continue
			}
			tk, ok := o.(ident)
			if !ok {
				continue
			}
			if reflect.TypeOf(vars[string(tk)]).Kind() == reflect.Func {
				out.Push(stk.Pop())
			}
		}
	}
	for {
		o := stk.Pop()
		if o == nil {
			break
		}
		if _, ok := o.(parenOpen); ok {
			return nil, fmt.Errorf("calc: mismatched parentheses")
		}
		out.Push(o)
	}
	return out, nil
}

// handleOperator updates the stacks with a new operator.
func handleOperator(t operator, stk, out *stack) {
	for {
		o := stk.Top()
		if o == nil {
			break
		}
		switch o.(type) {
		case operator, unary:
			break
		}
		if (operators[t].assoc == right || operators[t].prec > operators[o].prec) && (operators[t].assoc == left || operators[t].prec >= operators[o].prec) {
			break
		}
		out.Push(stk.Pop())
	}
	stk.Push(t)
}
