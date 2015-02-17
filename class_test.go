// Copyright 2013-2015 Oliver Eilhard.
// Use of this source code is governed by the MIT LICENSE that
// can be found in the MIT-LICENSE file included in the project.

package mruby

import (
	"html"
	"testing"
)

func TestNewClass(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	class, err := NewClass(ctx, "MyClass", nil)
	if err != nil {
		t.Fatalf("expected no error; got: %v", err)
	}
	if class == nil {
		t.Errorf("expected class; got: %v", class)
	}

	if !ctx.HasClass("MyClass", nil) {
		t.Fatalf("expected to find class %q", "MyClass")
	}
}

func TestDefineClass(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	class, err := ctx.DefineClass("MyClass", nil)
	if err != nil {
		t.Fatalf("expected no error; got: %v", err)
	}
	if class == nil {
		t.Errorf("expected class; got: %v", class)
	}

	if !ctx.HasClass("MyClass", nil) {
		t.Fatalf("expected to find class %q", "MyClass")
	}
}

func TestDefineClassScoping(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	var found bool
	/*
		This will raise an error which must be captured.
		_, found = ctx.GetModule("MissingClass", nil)
		if found {
			t.Fatalf("expected to not find class %q; got: %v", "MissingClass", found)
		}
	*/

	outer, err := ctx.DefineClass("Outer", nil)
	if err != nil {
		t.Fatalf("expected no error; got: %v", err)
	}
	if outer == nil {
		t.Errorf("expected outer class; got: %v", outer)
	}
	_, found = ctx.GetClass("Outer", nil)
	if !found {
		t.Fatalf("expected to find class %q; got: %v", "Outer", found)
	}
	if !ctx.HasClass("Outer", nil) {
		t.Fatalf("expected to find class %q", "Outer")
	}

	inner, err := ctx.DefineClass("Inner", outer)
	if err != nil {
		t.Fatalf("expected no error; got: %v", err)
	}
	if inner == nil {
		t.Errorf("expected inner class; got: %v", inner)
	}
	_, found = ctx.GetClass("Inner", outer)
	if !found {
		t.Fatalf("expected to find class %q; got: %v", "Outer::Inner", found)
	}
	/*
		_, found = ctx.GetClass("Inner", nil)
		if found {
			t.Fatalf("expected to not find class %q; got: %v", "::Inner", found)
		}
	*/
	if !ctx.HasClass("Inner", outer) {
		t.Fatalf("expected to find class %q", "Outer::Inner")
	}
}

func TestClassMethodWithNoArgs(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	class, err := ctx.DefineClass("MyClass", nil)
	if err != nil {
		t.Fatalf("expected no error; got: %v", err)
	}
	if class == nil {
		t.Errorf("expected classule; got: %v", class)
	}

	helloWorld := func(ctx *Context, self Value) (Value, error) {
		return ctx.ToValue("Hello world")
	}

	class.DefineClassMethod("hello", helloWorld, ArgsNone())

	s, err := ctx.LoadStringResult("MyClass.hello()")
	if err != nil {
		t.Fatal(err)
	}
	if s != "Hello world" {
		t.Errorf("expected %q; got: %q", "Hello world", s)
	}
}

func TestClassMethodWithRequiredArg(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	class, err := ctx.DefineClass("MyClass", nil)
	if err != nil {
		t.Fatalf("expected no error; got: %v", err)
	}
	if class == nil {
		t.Errorf("expected class; got: %v", class)
	}

	escapeHtml := func(ctx *Context, self Value) (output Value, err error) {
		// We expect a string here.
		sv, err := ctx.GetArgs("o", self)
		if err != nil {
			return NilValue(ctx), err
		}
		s, err := sv.ToString()
		if err == nil {
			s = html.EscapeString(s)
		}
		return ctx.ToValue(s)
	}

	class.DefineClassMethod("escape_html", escapeHtml, ArgsRequired(1))

	input := "<esc&ped>"
	expected := html.EscapeString(input)

	got, err := ctx.LoadStringResult("MyClass.escape_html(ARGV[0])", input)
	if err != nil {
		t.Fatal(err)
	}
	if got != expected {
		t.Errorf("expected %q; got: %q", expected, got)
	}
}
