// Copyright 2013-2015 Oliver Eilhard.
// Use of this source code is governed by the MIT LICENSE that
// can be found in the MIT-LICENSE file included in the project.

package mruby

import (
	"testing"
)

func TestParse(t *testing.T) {
	ctx := NewContext()

	parser, err := ctx.Parse("'Parsed'")
	if err != nil {
		t.Fatal(err)
	}
	val, err := parser.Run()
	if err != nil {
		t.Fatal(err)
	}
	res, err := val.ToInterface()
	if err != nil {
		t.Fatal(err)
	}
	s, ok := res.(string)
	if !ok {
		t.Errorf("expected string; got: %v", res)
	}
	if s != "Parsed" {
		t.Errorf("expected %q, got %q", "Parsed", s)
	}
}

func TestParseError(t *testing.T) {
	ctx := NewContext()

	rubycode := ".fail here!"

	_, err := ctx.Parse(rubycode)
	if err == nil {
		t.Fatal("expected parse error")
	}
	parseErr, ok := err.(*ParseError)
	if !ok {
		t.Fatal("expected ParseError")
	}
	if parseErr.Line != 1 {
		t.Errorf("expected error in line %d; got: %d", 1, parseErr.Line)
	}
	expected := "syntax error, unexpected '.'"
	if parseErr.Message != expected {
		t.Errorf("expected error message %q; got: %q", expected, parseErr.Message)
	}
	expected = "parse error: line 1: syntax error, unexpected '.'"
	if parseErr.Error() != expected {
		t.Errorf("expected error %q; got: %q", expected, parseErr.Error())
	}
}

func TestExceptionOnRun(t *testing.T) {
	ctx := NewContext()

	parser, err := ctx.Parse("raise 'kaboom'")
	if err != nil {
		t.Fatal(err)
	}
	val, err := parser.Run()
	if err == nil {
		t.Fatal("expected exception message as error, got nil")
	}
	if err.Error() != "kaboom" {
		t.Errorf("expected exception message 'kaboom', got %v", err.Error())
	}
	if !val.IsNil() {
		t.Fatal("expected result to be nil")
	}
}

func BenchmarkParse(b *testing.B) {
	ctx := NewContext()

	code := `
def concat(a, b)
	a + b
end

concat "Hello", "World"
`

	parser, err := ctx.Parse(code)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		val, err := parser.Run()
		if err != nil {
			b.Fatal(err)
		}
		res, err := val.ToInterface()
		if err != nil {
			b.Fatal(err)
		}
		s, ok := res.(string)
		if !ok {
			b.Errorf("run %d: expected string, got: %v", i, val)
		}
		if s != "HelloWorld" {
			b.Errorf("run %d: expected %q; got: %q", i, "HelloWorld", s)
		}
	}
}
