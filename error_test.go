package calc

import (
	"testing"
)

func TestError(t *testing.T) {
	exp := "1:5"
	s := `calc: bad character ':'
      1:5
       ^`
	_, err := EvalFloat(exp, nil)
	if err == nil {
		t.Fatalf("%q should fail", exp)
	}
	if err.Error() != s {
		t.Fatalf("bad error formatting. Expected:\n%s\nGot:\n%s\n", s, err)
	}
}
