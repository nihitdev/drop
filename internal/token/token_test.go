package token

import "testing"

func TestGenerate(t *testing.T) {
	a, err := Generate()
	if err != nil {
		t.Fatal(err)
	}
	b, err := Generate()
	if err != nil {
		t.Fatal(err)
	}
	if len(a) != 12 || len(b) != 12 || a == b {
		t.Fatalf("expected two distinct 12-character tokens, got %q and %q", a, b)
	}
}
