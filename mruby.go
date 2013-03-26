package mruby

/*
#cgo CFLAGS: -I./include
#cgo darwin LDFLAGS: -L./lib/darwin_amd64
#cgo linux LDFLAGS: -L./lib/linux_amd64
#cgo LDFLAGS: -lmruby
#include <stdlib.h>
#include <stdlib.h>
#include <string.h>

#include <mruby.h>
#include <mruby/proc.h>
#include <mruby/data.h>
#include <mruby/compile.h>
*/
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

type Context struct {
	mrb *C.mrb_state
	ctx *C.mrbc_context
}

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

func (ctx *Context) LoadString(script string) (string, error) {
	cs := C.CString(script)
	defer C.free(unsafe.Pointer(cs))

	C.mrb_load_string(ctx.mrb, cs)

	return "", errors.New("mruby: not implemented")
}
