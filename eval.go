package calc

import (
	"fmt"
	"math"
	"reflect"
)

// A stackf is a dead simple stack of float64s.
type stackf []float64

// Push pushes a float64 on the stack.
func (s *stackf) Push(x float64) {
	*s = append(*s, x)
}

// Pop pops an element from the stack.
func (s *stackf) Pop() float64 {
	n := len(*s)
	x := (*s)[n-1]
	*s = (*s)[:n-1]
	return x
}

// eval evaluates a list of tokens in RPN.
func eval(rpn []token, vars map[string]interface{}) (float64, error) {
	stk := make(stackf, 0)
	for _, t := range rpn {
		switch t := t.(type) {
		case number:
			stk.Push(float64(t))
		case unary:
			if len(stk) < 1 {
				return 0, fmt.Errorf("calc: invalid expression")
			}
			a := stk.Pop()
			switch t {
			case '+':
				stk.Push(a)
			case '-':
				stk.Push(-a)
			}
		case operator:
			if len(stk) < 2 {
				return 0, fmt.Errorf("calc: invalid expression")
			}
			b, a := stk.Pop(), stk.Pop()
			switch t {
			case '+':
				stk.Push(a + b)
			case '-':
				stk.Push(a - b)
			case '*', ' ':
				stk.Push(a * b)
			case '/':
				stk.Push(a / b)
			case '%':
				stk.Push(math.Mod(a, b))
			case '^':
				stk.Push(math.Pow(a, b))
			}
		case ident:
			switch v := vars[string(t)].(type) {
			case int:
				stk.Push(float64(v))
			case float64:
				stk.Push(v)
			default: // func
				f := reflect.TypeOf(v)
				n := f.NumIn()
				args := make([]reflect.Value, n)
				for i := n - 1; i >= 0; i-- {
					if len(stk) < 1 {
						return 0, fmt.Errorf("calc: invalid expression")
					}
					args[i] = reflect.ValueOf(stk.Pop())
				}
				r := reflect.ValueOf(v).Call(args)
				stk.Push(r[0].Interface().(float64))
			}
		}
	}
	if len(stk) != 1 {
		return 0, fmt.Errorf("calc: invalid expression")
	}
	return stk.Pop(), nil
}
