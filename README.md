# Calc

[![GoDoc](https://godoc.org/github.com/simleb/calc?status.svg)](http://godoc.org/github.com/simleb/calc)
[![Coverage Status](https://img.shields.io/coveralls/simleb/calc.svg)](https://coveralls.io/r/simleb/calc)
[![Build Status](https://drone.io/github.com/simleb/calc/status.png)](https://drone.io/github.com/simleb/calc/latest)

Package `calc` parses and evaluates mathematical expressions.

Expression are first tokenized (lexical analysis), then rearranged into reverse polish notation (RPN) using the shunting-yard algorithm. Finally, they are evaluated by treating all numbers as floating point.

Expressions might contain custom functions and variables.

This package is designed for short expressions, hence the input is a `string` and not an `io.Reader`.

## Example

```go
vars := map[string]interface{}{
	"x": 4,
	"inv": func(x float64) float64 {
		return 1 / x
	},
}
x, err := EvalFloat("2x^2-3x+inv(1+inv(x))", vars)
if err != nil {
	// handle error
}
fmt.Println(x)
```

## Todo

- [x] add more tests
- [x] allow unary minus
- [x] add modulus operator (`%`)
- [x] allow `2x` notation for `2*x` (and `2(…)` for `2*(…)`)
- [x] allow unicode identifiers
- [ ] add typical math functions (sqrt, trigo, abs, min, max…)

## License

The MIT License (MIT)

	Copyright (c) 2014 Simon Leblanc
	
	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:
	
	The above copyright notice and this permission notice shall be included in
	all copies or substantial portions of the Software.
	
	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE.
