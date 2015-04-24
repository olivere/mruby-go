// Copyright 2013-2015 Oliver Eilhard.
// Use of this source code is governed by the MIT LICENSE that
// can be found in the MIT-LICENSE file included in the project.

package mruby_test

import (
	"html"
	"testing"

	"github.com/olivere/mruby-go"
)

func BenchmarkFunctionCalls(b *testing.B) {
	// Create a new context, and set some options
	ctx := mruby.NewContext()
	if ctx == nil {
		b.Fatal("cannot create context")
	}

	// EscapeHtml as a non-trivial helper method.
	escapeHtml := func(ctx *mruby.Context) (output mruby.Value, err error) {
		// We expect a string here.
		args, err := ctx.GetArgs()
		if err != nil {
			return mruby.NilValue(ctx), err
		}
		if len(args) == 1 {
			s, err := args[0].ToString()
			if err != nil {
				return mruby.NilValue(ctx), err
			}
			s = html.EscapeString(s)
			return ctx.ToValue(s)
		}
		return mruby.NilValue(ctx), nil
	}

	// We create a new module called Helpers that will hold our extension method.
	module, err := ctx.DefineModule("Helpers", nil)
	if err != nil {
		b.Fatal(err)
	}
	if module == nil {
		b.Fatal(err)
	}

	module.DefineClassMethod("escape_html", escapeHtml)

	input := "<test&go>"
	expected := html.EscapeString(input)

	for i := 0; i < b.N; i++ {
		got, err := ctx.LoadStringResult("Helpers.escape_html(ARGV[0])", input)
		if err != nil {
			b.Fatal(err)
		}
		if expected != got {
			b.Fatalf("expected %q; got: %q", expected, got)
		}
	}
}

func BenchmarkFunctionCallsInParallel(b *testing.B) {
	// EscapeHtml as a non-trivial helper method.
	escapeHtml := func(ctx *mruby.Context) (output mruby.Value, err error) {
		// We expect a string here.
		args, err := ctx.GetArgs()
		if err != nil {
			return mruby.NilValue(ctx), err
		}
		if len(args) == 1 {
			s, err := args[0].ToString()
			if err != nil {
				return mruby.NilValue(ctx), err
			}
			s = html.EscapeString(s)
			return ctx.ToValue(s)
		}
		return mruby.NilValue(ctx), nil
	}

	input := "<test&go>"
	expected := html.EscapeString(input)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Create a new context, and set some options
			ctx := mruby.NewContext()
			if ctx == nil {
				b.Fatal("cannot create context")
			}

			// We create a new module called Helpers that will hold our extension method.
			module, err := ctx.DefineModule("Helpers", nil)
			if err != nil {
				b.Fatal(err)
			}
			if module == nil {
				b.Fatal(err)
			}

			module.DefineClassMethod("escape_html", escapeHtml)

			got, err := ctx.LoadStringResult("Helpers.escape_html(ARGV[0])", input)
			if err != nil {
				b.Fatal(err)
			}
			if expected != got {
				b.Fatalf("expected %q; got: %q", expected, got)
			}
		}
	})
}
