// Copyright 2013-2015 Oliver Eilhard.
// Use of this source code is governed by the MIT LICENSE that
// can be found in the MIT-LICENSE file included in the project.

package mruby

import (
	"testing"
)

func TestNewContext(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}
}

func TestNewContextWithOptions(t *testing.T) {
	ctx := NewContext(SetFilename("test.rb"), SetNoExec(true))
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}
	if ctx.filename != "test.rb" {
		t.Errorf("expected filename %q; got: %q", "test.rb", ctx.filename)
	}
	if ctx.noExec != true {
		t.Errorf("expected noExec = %v; got: %v", true, ctx.noExec)
	}
}

func TestLoadString(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	val, err := ctx.LoadString("'Hello world'")
	if err != nil {
		t.Fatal(err)
	}
	if !val.IsString() {
		t.Fatalf("expected value to be a string; got: %v", val.Type())
	}
	s, err := val.ToString()
	if err != nil {
		t.Fatal(err)
	}
	if s != "Hello world" {
		t.Errorf("expected %q; got: %q", "Hello world", s)
	}
}

func TestLoadStringResult(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	s, err := ctx.LoadStringResult("'Hello world'")
	if err != nil {
		t.Fatal(err)
	}
	if s != "Hello world" {
		t.Errorf("expected %q; got: %q", "Hello world", s)
	}
}

func TestLoadStringWithArgs(t *testing.T) {
	ctx := NewContext()

	var res interface{}
	var err error

	res, err = ctx.LoadStringResult("ARGV[0]", "Oliver", "Sandra")
	if err != nil {
		t.Fatal(err)
	}
	if res != "Oliver" {
		t.Errorf("expected %v, got %v", "Oliver", res)
	}

	res, err = ctx.LoadStringResult("ARGV[1]", "Oliver", "Sandra")
	if err != nil {
		t.Fatal(err)
	}
	if res != "Sandra" {
		t.Errorf("expected %v, got %v", "Sandra", res)
	}

	res, err = ctx.LoadStringResult("ARGV.inject { |x,y| x+y }", 1, 2, 3.5)
	if err != nil {
		t.Fatal(err)
	}
	if res != float64(6.5) {
		t.Errorf("expected %v, got %v", float64(6.5), res)
	}

	dict := map[string]interface{}{
		"name": "Oliver",
		"age":  23,
	}
	res, err = ctx.LoadStringResult(`ARGV[0]["age"]`, dict)
	if err != nil {
		t.Fatal(err)
	}
	if res != 23 {
		t.Errorf("expected %v, got %v", 23, res)
	}
}

func TestNoExec(t *testing.T) {
	ctx := NewContext(SetNoExec(true))
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	proc, err := ctx.LoadString("'Hello world'")
	if err != nil {
		t.Fatal(err)
	}
	if !proc.IsProc() {
		t.Fatalf("expected value to be a Proc; got: %v", proc.Type())
	}
	val, err := proc.Run()
	if err != nil {
		t.Fatal(err)
	}
	s, err := val.ToString()
	if err != nil {
		t.Fatal(err)
	}
	if s != "Hello world" {
		t.Errorf("expected %q; got: %q", "Hello world", s)
	}
}

func TestNoExecError(t *testing.T) {
	ctx := NewContext(SetNoExec(true))
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	proc, err := ctx.LoadString("raise 'kaboom'")
	if err != nil {
		t.Fatal(err)
	}
	if !proc.IsProc() {
		t.Fatalf("expected value to be a Proc; got: %v", proc.Type())
	}
	_, err = proc.Run()
	if err == nil {
		t.Fatalf("expected error")
	}
	e, ok := err.(*RunError)
	if !ok {
		t.Fatal("expected RunError")
	}
	expected := "kaboom"
	if e.Message != expected {
		t.Errorf("expected %q; got: %q", expected, e.Message)
	}
	if e.Error() != expected {
		t.Errorf("expected %q; got: %q", expected, e.Error())
	}
}
