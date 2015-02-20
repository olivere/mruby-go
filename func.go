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
	"log"
	"unsafe"
)

var _ = log.Print

// methodMap maps a Ruby symbol to a function.
type methodMap map[C.mrb_sym]Function

// Function defines the signature of a Go function that can be called
// from within a MRuby script. The self parameter can be used to extract
// the arguments passed to the extension method.
type Function func(ctx *Context, self Value) (Value, error)

//export my_mrb_func_call
func my_mrb_func_call(mrb *C.mrb_state, v C.mrb_value) C.mrb_value {
	contextsFu.Lock()
	defer contextsFu.Unlock()

	// Find the context by mrb.
	ctx, found := contexts[mrb]
	if !found {
		return C.mrb_nil_value()
	}

	// Find the function in the ctx.
	ctx.methodsMu.Lock()
	defer ctx.methodsMu.Unlock()

	if ctx.methodsByRClass == nil {
		return C.mrb_nil_value()
	}

	input := Value{ctx: ctx, v: v}

	callinfo := mrb.c.ci

	methods, found := ctx.methodsByRClass[callinfo.proc.target_class]
	if !found {
		return C.mrb_nil_value()
	}

	method, found := methods[callinfo.mid]
	if !found {
		return C.mrb_nil_value()
	}

	output, err := method(ctx, input)
	if err != nil {
		return C.mrb_nil_value()
	}
	return output.v
}

// addMethod inserts a method to the given class.
func (ctx *Context) addMethod(class *C.struct_RClass, name string, f Function) {
	ctx.methodsMu.Lock()
	defer ctx.methodsMu.Unlock()

	// Register the function.
	if ctx.methodsByRClass == nil {
		ctx.methodsByRClass = make(map[*C.struct_RClass]methodMap)
	}

	methods, found := ctx.methodsByRClass[class]
	if !found {
		methods = make(map[C.mrb_sym]Function)
		ctx.methodsByRClass[class] = methods
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	sym := C.mrb_intern_cstr(ctx.mrb, cname)
	methods[sym] = f
}
