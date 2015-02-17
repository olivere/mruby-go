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

// Class represents a class in Ruby.
type Class struct {
	ctx   *Context
	class *C.struct_RClass
}

// NewClass defines a new class with the given name and super-class
// in the context.
// If super is nil, ObjectClass is used by default.
func NewClass(ctx *Context, name string, super *Class) (*Class, error) {
	if super == nil {
		super = ctx.ObjectClass()
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	class := C.mrb_define_class(ctx.mrb, cname, super.class)
	return &Class{ctx: ctx, class: class}, nil
}

// DefineClass defines a new class with the given name and super-class
// in the context. If super is nil, ObjectClass is used by default.
func (ctx *Context) DefineClass(name string, super *Class) (*Class, error) {
	return NewClass(ctx, name, super)
}

// HasClass tests if the context has a class with the given name.
func (ctx *Context) HasClass(name string, outer *Class) bool {
	/*
		cname := C.CString(name)
		defer C.free(unsafe.Pointer(cname))
		flag := C.my_mrb_const_defined_at(ctx.mrb, cname)
		return flag == C.mrb_bool(1)
	*/
	_, found := ctx.GetClass(name, outer)
	return found
}

// GetClass returns the given class.
func (ctx *Context) GetClass(name string, outer *Class) (*Class, bool) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	if outer == nil {
		outer = ctx.ObjectClass()
	}
	class := C.mrb_class_get_under(ctx.mrb, outer.class, cname)
	if C.has_exception(ctx.mrb) != 0 {
		return nil, false
	}
	if class == nil {
		return nil, false
	}
	return &Class{ctx: ctx, class: class}, true
}

// DefineMethod registers an instance method with the name in the class.
// The function is called when executed in Ruby. The args value specifies
// the number of required and optional arguments (if any) of f.
func (c *Class) DefineMethod(name string, f Function, args Args) {
	c.ctx.addMethod(c.class, name, f)

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	C.mrb_define_method(
		c.ctx.mrb,
		c.class,
		cname,
		C.my_mrb_func_call_t(),
		C.mrb_aspec(args))
}

// DefineMethod registers a class method with the name in the class.
// The function is called when executed in Ruby. The args value specifies
// the number of required and optional arguments (if any) of f.
func (c *Class) DefineClassMethod(name string, f Function, args Args) {
	c.ctx.addMethod(c.class.c, name, f)

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	C.mrb_define_class_method(
		c.ctx.mrb,
		c.class,
		cname,
		C.my_mrb_func_call_t(),
		C.mrb_aspec(args))
}
