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

// Module is a Ruby module.
type Module struct {
	ctx    *Context
	module *C.struct_RClass
}

// NewModule defines a new module with the given name under outer.
// If outer is nil, the registered module is a top-level module.
func NewModule(ctx *Context, name string, outer RClass) (*Module, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	if outer == nil {
		outer = ctx.ObjectModule()
	}
	module := C.mrb_define_module_under(ctx.mrb, outer.RClass(), cname)
	return &Module{ctx: ctx, module: module}, nil
}

// RClass returns the MRuby object of type *struct RClass.
func (m *Module) RClass() *C.struct_RClass {
	return m.module
}

// DefineModule defines a new module with the given name under outer.
// If outer is nil, the registered module is a top-level module.
func (ctx *Context) DefineModule(name string, outer RClass) (*Module, error) {
	return NewModule(ctx, name, outer)
}

// HasModule tests if the context has a module with the given name.
func (ctx *Context) HasModule(name string, outer RClass) bool {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var klass *C.struct_RClass
	if outer != nil {
		klass = outer.RClass()
	}
	flag := C.my_mrb_has_module(ctx.mrb, klass, cname)
	return flag == C.mrb_bool(1)
}

// GetModule returns the given module.
func (ctx *Context) GetModule(name string, outer RClass) (*Module, bool) {
	if !ctx.HasModule(name, outer) {
		return nil, false
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	if outer == nil {
		outer = ctx.ObjectModule()
	}
	module := C.mrb_module_get_under(ctx.mrb, outer.RClass(), cname)
	if C.has_exception(ctx.mrb) != 0 {
		return nil, false
	}
	if module == nil {
		return nil, false
	}
	return &Module{ctx: ctx, module: module}, true
}

// DefineMethod registers a method with the name in the module.
// The function is called when executed in Ruby. The args value specifies
// the number of required and optional arguments (if any) of f.
func (m *Module) DefineMethod(name string, f Function, args Args) {
	m.ctx.addMethod(m.module, name, f)

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	C.mrb_define_method(
		m.ctx.mrb,
		m.module,
		cname,
		C.my_mrb_func_call_t(),
		C.mrb_aspec(args))
}

// DefineClassMethod registers a class method with the name in the module.
// The function is called when executed in Ruby. The args value specifies
// the number of required and optional arguments (if any) of f.
func (m *Module) DefineClassMethod(name string, f Function, args Args) {
	m.ctx.addMethod(m.module.c, name, f)

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	// Note: Use mrb_define_method instead of mrb_define_class_method here.
	C.mrb_define_method(
		m.ctx.mrb,
		m.module.c,
		cname,
		C.my_mrb_func_call_t(),
		C.mrb_aspec(args))
}
