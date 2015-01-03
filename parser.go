// Copyright 2013-2015 Oliver Eilhard.
// Use of this source code is governed by the MIT LICENSE that
// can be found in the MIT-LICENSE file included in the project.

package mruby

/*
#cgo pkg-config: mruby
#include "mruby_go.h"
*/
import "C"

import (
	"unsafe"
)

// Parser is a parser for Ruby code. It can be used to parse
// Ruby code once and run it multiple times.
type Parser struct {
	ctx  *Context
	proc *C.struct_RProc
}

// Parse parses a string into parsed Ruby code. An error is
// returned if compilation failes.
func (ctx *Context) Parse(code string) (*Parser, error) {
	p := &Parser{ctx: ctx}

	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))

	parser := C.my_parse(p.ctx.mrb, p.ctx.ctx, ccode)
	defer C.mrb_parser_free(parser)

	if parser.nerr > 0 {
		line := int(parser.error_buffer[0].lineno)
		msg := C.GoString(parser.error_buffer[0].message)
		return nil, &ParseError{Line: line, Message: msg}
	}

	p.proc = C.mrb_generate_code(p.ctx.mrb, parser)
	return p, nil
}

// Run runs a previously compiled Ruby code and returns its output.
// An error is returned if the Ruby code raises an exception.
func (p *Parser) Run(args ...interface{}) (Value, error) {
	ai := C.mrb_gc_arena_save(p.ctx.mrb)
	defer C.mrb_gc_arena_restore(p.ctx.mrb, ai)

	// Create ARGV global variable and push the args into it
	argvAry := C.mrb_ary_new_capa(p.ctx.mrb, C.mrb_int(len(args)))
	for i := 0; i < len(args); i++ {
		val, err := p.ctx.ToValue(args[i])
		if err != nil {
			return NilValue(p.ctx), err
		}
		C.mrb_ary_push(p.ctx.mrb, argvAry, val.v)
	}
	argv := C.CString("ARGV")
	defer C.free(unsafe.Pointer(argv))
	C.mrb_define_global_const(p.ctx.mrb, argv, argvAry)

	// Run the code
	result := C.my_run(p.ctx.mrb, p.proc)

	// Check for exception
	if C.has_exception(p.ctx.mrb) != 0 {
		return NilValue(p.ctx), newRunError(p.ctx, true)
	}

	return Value{ctx: p.ctx, v: result}, nil
}
