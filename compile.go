package mruby

/*
#cgo CFLAGS: -I./include
#cgo darwin LDFLAGS: -L./lib/darwin_amd64
#cgo linux LDFLAGS: -L./lib/linux_amd64
#cgo LDFLAGS: -lmruby
#include <stdlib.h>
#include <string.h>

#include <mruby.h>
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
my_type(mrb_value v) {
	return mrb_type(v);
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

type Parser struct {
	ctx    *Context
	parser *C.struct_mrb_parser_state
	n      C.int
}

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

func (p *Parser) Run() (interface{}, error) {

	ai := C.mrb_gc_arena_save(p.ctx.mrb)
	defer C.mrb_gc_arena_restore(p.ctx.mrb, ai)

	result := C.my_run(p.ctx.mrb, p.n)

	if C.has_exception(p.ctx.mrb) != 0 {
		msg := C.GoString(C.get_exception_message(p.ctx.mrb))
		C.reset_exception(p.ctx.mrb)
		return nil, errors.New(fmt.Sprintf("%s", msg))
	}

	return ruby2go(p.ctx, result), nil
}

/*
func (p *Parser) Parse(code string) (interface{}, error) {
	//var result C.mrb_value

	ai := C.mrb_gc_arena_save(p.ctx.mrb)
	defer C.mrb_gc_arena_restore(p.ctx.mrb, ai)

	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))

	parser := C.my_parse(p.ctx.mrb, p.ctx.ctx, ccode)
	defer C.mrb_parser_free(parser)

	if parser.nerr > 0 {
		lineno := parser.error_buffer[0].lineno
		msg := C.GoString(parser.error_buffer[0].message)
		return nil, errors.New(fmt.Sprintf("error: line %d: %s", lineno, msg))
	}

	n := C.mrb_generate_code(p.ctx.mrb, parser)

	result := C.my_run(p.ctx.mrb, n)

	return ruby2go(p.ctx, result), nil
}
*/

func ruby2go(ctx *Context, v C.mrb_value) interface{} {
	switch C.my_type(v) {
	case C.MRB_TT_FALSE:
		return false
	case C.MRB_TT_FREE:
		return nil
	case C.MRB_TT_TRUE:
		return true
	case C.MRB_TT_FIXNUM:
		return nil
	case C.MRB_TT_SYMBOL:
		return nil
	case C.MRB_TT_UNDEF:
		return nil
	case C.MRB_TT_FLOAT:
		return nil
	case C.MRB_TT_VOIDP:
		return nil
	case C.MRB_TT_OBJECT:
		return nil
	case C.MRB_TT_CLASS:
		return nil
	case C.MRB_TT_MODULE:
		return nil
	case C.MRB_TT_ICLASS:
		return nil
	case C.MRB_TT_SCLASS:
		return nil
	case C.MRB_TT_PROC:
		return nil
	case C.MRB_TT_ARRAY:
		return nil
	case C.MRB_TT_HASH:
		return nil
	case C.MRB_TT_STRING:
		return C.GoString(C.mrb_string_value_ptr(ctx.mrb, v))
	case C.MRB_TT_RANGE:
		return nil
	case C.MRB_TT_EXCEPTION:
		return nil
	case C.MRB_TT_FILE:
		return nil
	case C.MRB_TT_ENV:
		return nil
	case C.MRB_TT_DATA:
		return nil
	case C.MRB_TT_MAXDEFINE:
		return nil
	}
	return nil
}
