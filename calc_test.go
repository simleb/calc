package calc

import (
	"fmt"
	"testing"
)

// print prints a list of tokens inline separated by spaces.
// Useful for debugging.
func print(tokens []token) {
	for _, v := range tokens {
		switch v := v.(type) {
		case unary:
			fmt.Printf("%c̲ ", v)
		case operator:
			fmt.Printf("%c ", v)
		case parenOpen:
			fmt.Printf("( ")
		case parenClose:
			fmt.Printf(") ")
		case separator:
			fmt.Printf(", ")
		default:
			fmt.Printf("%v ", v)
		}
	}
	fmt.Printf("\n")
}

// good is a list of test cases that should succeed.
var good = []struct {
	exp string
	res float64
}{
	{"0", 0},
	{"-1", -1},
	{"+1", 1},
	{"1", 1},
	{" 1", 1},
	{"1 ", 1},
	{"  1	 ", 1},
	{"1.", 1},
	{"1.0", 1},
	{"1.01", 1.01},
	{".9", 0.9},
	{"1.23456", 1.23456},
	{"2e1", 20},
	{"2E2", 200},
	{"2e-1", 0.2},
	{"2e+1", 20},
	{"2.e1", 20},
	{".2e1", 2},
	{"1.2e1", 12},
	{"(1)", 1},
	{"(-1)", -1},
	{"-(1)", -1},
	{"+(1)", 1},
	{"((1))", 1},
	{"( 1 )", 1},
	{"1+2/1", 3},
	{"10/-1*-2", 20},
	{"-2^2", -4},
	{"2^-2", 0.25},
	{"1 +2", 3},
	{"1+	2", 3},
	{"1 + 2", 3},
	{"(1+2)", 3},
	{"(1)+2", 3},
	{"((1)+2)", 3},
	{"1-2", -1},
	{"2^3^2", 512},
	{"1+2*2^3*4/(6-1)", 13.8},
	{"--1", 1},
	{"---1", -1},
	{"++-+-+++-+1", -1},
	{"10%3", 1},
	{"-10%3", -1},
	{"2*4+1%3", 9},
	{"(2*4+1)%3", 0},
	{"2(3+1)", 8},
	{"12/2(1+1)", 3},
}

// bas is a list of test cases that should fail.
var bad = []string{
	"",
	"+",
	"-",
	"%",
	"1+",
	"/1",
	".",
	",",
	"1+.",
	"1+1e+",
	"()",
	"(1",
	"2)",
	")",
	"(3)/2)",
	"(3+(1)",
	"2^(1-)",
	"π",
	"1 2",
	"10/3 5-9",
	",",
	"(,)",
	"f(,)",
	"(1,2)",
}

// vars is a list of variables and functions used in some tests.
var vars = map[string]interface{}{
	"π":    3.14,
	"life": 42,
	"inc": func(x float64) float64 {
		return x + 1
	},
	"sqdist": func(x, y float64) float64 {
		return x*x + y*y
	},
}

// vgood is a list of test cases that should succeed when the variables are provided.
var vgood = []struct {
	exp string
	res float64
}{
	{"π", 3.14},
	{"2*π", 6.28},
	{"2π", 6.28},
	{"1/2π", 1 / 6.28},
	{"life/4", 10.5},
	{"inc(1)", 2},
	{"inc(-2+2)", 1},
	{"inc(2*inc(1-2)^3)", 1},
	{"1+inc(inc(1))", 4},
	{"inc(2)+2", 5},
	{"sqdist(3, 4)", 25},
	{"sqdist(1+2*3-4, 2*2)", 25},
}

// vbad is a list of test cases that should fail even when the variables are provided.
var vbad = []string{
	"foo",
	"foo()",
	"inc()",
	"inc(,)",
	"inc(2,)",
	"inc(,3)",
	"inc(2,3)",
	"sqdist(2)",
	"sqdist(2,3,4)",
	"sqdist(2,)",
	"sqdist(,2)",
}

func TestEvalFloat(t *testing.T) {
	for _, v := range good {
		t.Logf("%q\n", v.exp)
		x, err := EvalFloat(v.exp, nil)
		if err != nil {
			t.Fatal(err)
		}
		if x != v.res {
			t.Fatalf("expected %f, got %f", v.res, x)
		}
	}
	for _, exp := range bad {
		t.Logf("%q\n", exp)
		if _, err := EvalFloat(exp, nil); err == nil {
			t.Fatalf("error: expression %q should fail", exp)
		}
	}
	for _, v := range vgood {
		t.Logf("vars: %q\n", v.exp)
		x, err := EvalFloat(v.exp, vars)
		if err != nil {
			t.Fatal(err)
		}
		if x != v.res {
			t.Fatalf("expected %f, got %f", v.res, x)
		}
	}
	for _, v := range vgood {
		t.Logf("no vars: %q\n", v.exp)
		_, err := EvalFloat(v.exp, nil)
		if err == nil {
			t.Fatalf("error: expression %q should fail", v.exp)
		}
	}
	for _, exp := range vbad {
		t.Logf("%q\n", exp)
		_, err := EvalFloat(exp, vars)
		if err == nil {
			t.Fatalf("error: expression %q should fail", exp)
		}
	}
}

func TestEvalInt(t *testing.T) {
	n, err := EvalInt("1+2*2^3*4/(6-1)", nil)
	if err != nil {
		t.Fatal(err)
	}
	if n != 13 {
		t.Fatalf("expected 13, got %d", n)
	}
}

// testDebug does what EvalFloat does but prints the tokens.
func testDebug(t *testing.T) {
	tokens, err := tokenize("2(3)")
	if err != nil {
		t.Fatal(err)
	}
	print(tokens)
	rpn, err := postfix(tokens, vars)
	if err != nil {
		t.Fatal(err)
	}
	print(rpn)
	_, err = eval(rpn, vars)
	if err != nil {
		t.Fatal(err)
	}
}

func ExampleEvalFloat_variables() {
	x, err := EvalFloat("2x^2-3x+1.5", map[string]interface{}{"x": 4})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(x)
	// Output: 21.5
}

func ExampleEvalFloat_functions() {
	vars := map[string]interface{}{
		"inc": func(x float64) float64 {
			return x + 1
		},
	}
	x, err := EvalFloat("1 + inc(1.5 + inc(2))", vars)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(x)
	// Output: 6.5
}
