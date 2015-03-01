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

func TestHasClassAndGetClassFailWhenMissing(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	found := ctx.HasClass("MissingClass", nil)
	if found {
		t.Errorf("expected to not find class %q; got: %v", "MissingClass", found)
	}

	class, found := ctx.GetClass("MissingClass", nil)
	if found {
		t.Errorf("expected to not find class %q; got: %v", "MissingClass", found)
	}
	if class != nil {
		t.Fatalf("expected to return nil for missing class; got: %v", class)
	}
}

func TestDefineClassOnTopLevel(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	class, err := ctx.DefineClass("MyClass", nil)
	if err != nil {
		t.Fatalf("expected no error; got: %v", err)
	}
	if class == nil {
		t.Fatalf("expected class; got: %v", class)
	}
	_, found := ctx.GetClass("MyClass", nil)
	if !found {
		t.Errorf("expected to find class %q; got: %v", "MyClass", found)
	}
	found = ctx.HasClass("MyClass", nil)
	if !found {
		t.Errorf("expected to find class %q; got: %v", "MyClass", found)
	}
}

func TestDefineClassesInModule(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	module, err := ctx.DefineModule("MyModule", nil)
	if err != nil {
		t.Fatalf("expected to define module; got: %v", err)
	}
	if module == nil {
		t.Fatalf("expected module; got: %v", module)
	}

	// This defines a class MyModule::MyClass
	class, err := ctx.DefineClassUnder("MyClass", nil, module)
	if err != nil {
		t.Fatalf("expected no error; got: %v", err)
	}
	if class == nil {
		t.Fatalf("expected class; got: %v", class)
	}
	_, found := ctx.GetClass("MyClass", module)
	if !found {
		t.Errorf("expected to find class %q; got: %v", "MyModule::MyClass", found)
	}
	found = ctx.HasClass("MissingClass", module)
	if found {
		t.Errorf("expected to not find class %q; got: %v", "MyModule::MissingClass", found)
	}

	// This defines a class MyModule::MyDerivedClass, derived from MyClass
	derived, err := ctx.DefineClassUnder("MyDerivedClass", class, module)
	if err != nil {
		t.Fatalf("expected no error; got: %v", err)
	}
	if derived == nil {
		t.Fatalf("expected class; got: %v", class)
	}
	_, found = ctx.GetClass("MyDerivedClass", module)
	if !found {
		t.Errorf("expected to find class %q; got: %v", "MyModule::MyDerivedClass", found)
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
