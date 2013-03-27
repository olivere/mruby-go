// Copyright 2013 Oliver Eilhard.
// Use of this source code is governed by the MIT LICENSE that
// can be found in the MIT-LICENSE file included in the project.
package mruby

/*
#cgo LDFLAGS: -lmruby
#include <stdlib.h>
#include <string.h>

#include <mruby.h>
#include <mruby/proc.h>
#include <mruby/data.h>
#include <mruby/compile.h>
*/
import "C"

import (
	"runtime"
)

// Context serves as the entry point for all communication with mruby.
type Context struct {
	mrb *C.mrb_state
	ctx *C.mrbc_context
}

// NewContext creates a new mruby context.
func NewContext() *Context {
	ctx := &Context{}

	ctx.mrb = C.mrb_open()
	ctx.ctx = C.mrbc_context_new(ctx.mrb)
	//ctx.ctx.capture_errors = C.int(1)

	runtime.SetFinalizer(ctx, func(x *Context) {
		C.mrbc_context_free(x.mrb, x.ctx)
		C.mrb_close(x.mrb)
	})

	return ctx
}
