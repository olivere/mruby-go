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
	"reflect"
	"runtime"
	"unsafe"
)

// Context serves as the entry point for all communication with mruby.
type Context struct {
	mrb *C.mrb_state
	ctx *C.mrbc_context

	captureErrors bool   // indicates whether script errors are captured
	noExec        bool   // automatically "run" the scripts given to the context
	filename      string // filename used internally
}

// NewContext creates a new mruby context. Use the options to handle
// configuration.
//
// Examples:
//   ctx := mruby.NewContext()
//   ctx := mruby.NewContext(mruby.SetNoExec(true), mruby.SetFilename("simple.rb"))
func NewContext(options ...func(*Context)) *Context {
	ctx := &Context{
		captureErrors: true,
		noExec:        false,
		filename:      "(mruby-go)",
	}

	// Run configuration handlers
	for _, option := range options {
		option(ctx)
	}

	// Finalize setup of MRB context
	cfilename := C.CString(ctx.filename)
	defer C.free(unsafe.Pointer(cfilename))

	captureErrors := C.mrb_bool(0)
	if ctx.captureErrors {
		captureErrors = C.mrb_bool(1)
	}
	noExec := C.mrb_bool(0)
	if ctx.noExec {
		noExec = C.mrb_bool(1)
	}

	ctx.mrb = C.mrb_open()
	ctx.ctx = C.my_context_new(ctx.mrb, cfilename, captureErrors, noExec)

	runtime.SetFinalizer(ctx, func(x *Context) {
		C.mrbc_context_free(x.mrb, x.ctx)
		C.mrb_close(x.mrb)
		x.mrb = nil
		x.ctx = nil
	})

	return ctx
}

// SetCaptureErrors indicates whether script errors are captured (default: true).
// It is used for configuring a Context (see NewContext for details).
func SetCaptureErrors(captureErrors bool) func(*Context) {
	return func(ctx *Context) {
		ctx.captureErrors = captureErrors
	}
}

// SetNoExec indicates whether scripts given to this context, e.g. via
// LoadString, are automatically run once loaded and/or parsed.
// It is used for configuring a Context (see NewContext for details).
func SetNoExec(noExec bool) func(*Context) {
	return func(ctx *Context) {
		ctx.noExec = noExec
	}
}

// SetFilename sets the filename to be used in the context (default: "(mruby-go)").
// It is used for configuring a Context (see NewContext for details).
func SetFilename(filename string) func(*Context) {
	return func(ctx *Context) {
		ctx.filename = filename
	}
}

// GC runs the full MRuby garbage collector.
func (ctx *Context) GC() {
	C.mrb_full_gc(ctx.mrb)
}

// GC runs the incremental MRuby garbage collector.
func (ctx *Context) IncrementalGC() {
	C.mrb_incremental_gc(ctx.mrb)
}

// LoadString loads a snippet of Ruby code and returns its output.
// An error is returned if the interpreter failes or the Ruby code
// raises an exception of type RunError.
func (ctx *Context) LoadString(code string, args ...interface{}) (Value, error) {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))

	ai := C.mrb_gc_arena_save(ctx.mrb)
	defer C.mrb_gc_arena_restore(ctx.mrb, ai)

	// Create ARGV global variable and push the args into it
	argv := C.CString("ARGV")
	defer C.free(unsafe.Pointer(argv))
	argvAry := C.mrb_ary_new_capa(ctx.mrb, C.mrb_int(len(args)))
	for i := 0; i < len(args); i++ {
		val, err := ctx.ToValue(args[i])
		if err != nil {
			return NilValue(ctx), err
		}
		C.mrb_ary_push(ctx.mrb, argvAry, val.v)
	}
	C.mrb_define_global_const(ctx.mrb, argv, argvAry)

	result := C.mrb_load_string_cxt(ctx.mrb, ccode, ctx.ctx)

	if C.has_exception(ctx.mrb) != 0 {
		return NilValue(ctx), newRunError(ctx, true)
	}

	//log.Printf("mruby result type: %s\n", rubyTypeOf(ctx, result))

	return Value{ctx: ctx, v: result}, nil
}

// LoadStringResult invokes LoadString and returns the Go value immediately.
// Use this method to skip testing the returned Value.
func (ctx *Context) LoadStringResult(code string, args ...interface{}) (interface{}, error) {
	val, err := ctx.LoadString(code, args...)
	if err != nil {
		return NilValue(ctx), err
	}
	return val.ToInterface()
}

// ToValue stores the given value for encoding/decoding from/to Go and MRuby.
func (ctx *Context) ToValue(value interface{}) (Value, error) {
	valof := reflect.ValueOf(value)
	switch valof.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Value{ctx: ctx, v: C.mrb_fixnum_value(C.mrb_int(valof.Int()))}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return Value{ctx: ctx, v: C.mrb_fixnum_value(C.mrb_int(valof.Uint()))}, nil
	case reflect.Float32, reflect.Float64:
		return Value{ctx: ctx, v: C.get_float_value(ctx.mrb, C.mrb_float(valof.Float()))}, nil
	case reflect.String:
		cs := C.CString(valof.String())
		defer C.free(unsafe.Pointer(cs))
		return Value{ctx: ctx, v: C.mrb_str_new_cstr(ctx.mrb, cs)}, nil
	case reflect.Bool:
		if valof.Bool() {
			return Value{ctx: ctx, v: C.mrb_true_value()}, nil
		} else {
			return Value{ctx: ctx, v: C.mrb_false_value()}, nil
		}
	case reflect.Array, reflect.Slice:
		ary := C.mrb_ary_new(ctx.mrb)
		for i := 0; i < valof.Len(); i++ {
			elem, err := ctx.ToValue(valof.Index(i).Interface())
			if err != nil {
				return NilValue(ctx), err
			}
			C.mrb_ary_push(ctx.mrb, ary, elem.v)
		}
		return Value{ctx: ctx, v: ary}, nil
	case reflect.Map:
		hsh := C.mrb_hash_new(ctx.mrb)
		for _, key := range valof.MapKeys() {
			mapvalue := valof.MapIndex(key)
			keyv, err := ctx.ToValue(key.String())
			if err != nil {
				return NilValue(ctx), err
			}
			valv, err := ctx.ToValue(mapvalue.Interface())
			if err != nil {
				return NilValue(ctx), err
			}
			C.mrb_hash_set(ctx.mrb, hsh, keyv.v, valv.v)
		}
		return Value{ctx: ctx, v: hsh}, nil
	case reflect.Interface:
		return ctx.ToValue(valof.Elem().Interface())
	case reflect.Ptr:
		// Drill down.
		if valof.IsNil() {
			return NilValue(ctx), nil
		} else {
			valof = valof.Elem()
			return ctx.ToValue(valof.Interface())
		}
	}
	return NilValue(ctx), nil
}
