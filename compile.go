// Copyright 2013 Oliver Eilhard.
// Use of this source code is governed by the MIT LICENSE that
// can be found in the MIT-LICENSE file included in the project.
package mruby

/*
#cgo CFLAGS: -I./include
#cgo darwin LDFLAGS: -L./lib/darwin_x64
#cgo linux LDFLAGS: -L./lib/linux_x64 -lm
#cgo LDFLAGS: -lmruby
#include <stdlib.h>
#include <string.h>

#include <mruby.h>
#include <mruby/array.h>
#include <mruby/proc.h>
#include <mruby/data.h>
#include <mruby/compile.h>
#include <mruby/string.h>
#include <mruby/value.h>

struct mrb_parser_state *
my_parse(mrb_state *mrb, mrbc_context *ctx, char *ruby_code) {
	struct mrb_parser_state *parser = mrb_parser_new(mrb);

	parser->s = ruby_code;
	parser->send = ruby_code + strlen(ruby_code);
	parser->lineno = 1;
	mrb_parser_parse(parser, ctx);

	return parser;
}

mrb_value
my_run(mrb_state *mrb, int n) {
	return mrb_run(mrb,
		mrb_proc_new(mrb, mrb->irep[n]),
		mrb_top_self(mrb)
	);
}

int
has_exception(mrb_state *mrb) {
	return mrb->exc != 0;
}

void
reset_exception(mrb_state *mrb) {
	mrb->exc = 0;
}

char *
get_exception_message(mrb_state *mrb) {
	mrb_value val = mrb_obj_value(mrb->exc);
	return mrb_string_value_ptr(mrb, val);
}
*/
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

// Parser is a parser for Ruby code. It can be used to parse
// Ruby code once and run it multiple times.
type Parser struct {
	ctx    *Context
	parser *C.struct_mrb_parser_state
	n      C.int
}

// Parse parses a string into parsed Ruby code. An error is
// returned if compilation failes.
func (ctx *Context) Parse(code string) (*Parser, error) {
	p := &Parser{ctx: ctx, n: -1}

	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))

	p.parser = C.my_parse(p.ctx.mrb, p.ctx.ctx, ccode)

	if p.parser.nerr > 0 {
		lineno := p.parser.error_buffer[0].lineno
		msg := C.GoString(p.parser.error_buffer[0].message)
		return nil, errors.New(fmt.Sprintf("error: line %d: %s", lineno, msg))
	}

	p.n = C.mrb_generate_code(p.ctx.mrb, p.parser)

	runtime.SetFinalizer(p, func(p *Parser) {
		if p.parser != nil {
			C.mrb_parser_free(p.parser)
		}
	})

	return p, nil
}

// Run runs a previously compiled Ruby code and returns its output.
// An error is returned if the Ruby code raises an exception.
func (p *Parser) Run(args ...interface{}) (interface{}, error) {
	//ai := C.mrb_gc_arena_save(p.ctx.mrb)
	//defer C.mrb_gc_arena_restore(p.ctx.mrb, ai)

	// Create ARGV global variable and push the args into it
	argvAry := C.mrb_ary_new(p.ctx.mrb)
	for i := 0; i < len(args); i++ {
		C.mrb_ary_push(p.ctx.mrb, argvAry, go2ruby(p.ctx, args[i]))
	}
	argv := C.CString("ARGV")
	defer C.free(unsafe.Pointer(argv))
	C.mrb_define_global_const(p.ctx.mrb, argv, argvAry)

	// Run the code
	result := C.my_run(p.ctx.mrb, p.n)

	// Check for exception
	if C.has_exception(p.ctx.mrb) != 0 {
		msg := C.GoString(C.get_exception_message(p.ctx.mrb))
		C.reset_exception(p.ctx.mrb)
		return nil, errors.New(msg)
	}

	return ruby2go(p.ctx, result), nil
}

// LoadString loads a snippet of Ruby code and returns its output.
// An error is returned if the interpreter failes or the Ruby code
// raises an exception.
func (ctx *Context) LoadString(code string, args ...interface{}) (interface{}, error) {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))

	//ai := C.mrb_gc_arena_save(ctx.mrb)
	//defer C.mrb_gc_arena_restore(ctx.mrb, ai)

	// Create ARGV global variable and push the args into it
	argv := C.CString("ARGV")
	defer C.free(unsafe.Pointer(argv))
	argvAry := C.mrb_ary_new(ctx.mrb)
	for i := 0; i < len(args); i++ {
		C.mrb_ary_push(ctx.mrb, argvAry, go2ruby(ctx, args[i]))
	}
	C.mrb_define_global_const(ctx.mrb, argv, argvAry)

	result := C.mrb_load_string_cxt(ctx.mrb, ccode, ctx.ctx)

	if C.has_exception(ctx.mrb) != 0 {
		msg := C.GoString(C.get_exception_message(ctx.mrb))
		C.reset_exception(ctx.mrb)
		return nil, errors.New(msg)
	}

	//log.Printf("mruby result type: %s\n", rubyTypeOf(ctx, result))

	return ruby2go(ctx, result), nil
}
