// Copyright 2013 Oliver Eilhard.
// Use of this source code is governed by the MIT LICENSE that
// can be found in the MIT-LICENSE file included in the project.
package mruby

/*
#cgo linux LDFLAGS: -lm
#cgo LDFLAGS: -lmruby
#include <stdlib.h>
#include <string.h>

#include <mruby.h>
#include <mruby/array.h>
#include <mruby/hash.h>
#include <mruby/proc.h>
#include <mruby/data.h>
#include <mruby/compile.h>
#include <mruby/string.h>
#include <mruby/value.h>

int
my_type(mrb_value v) {
	return mrb_type(v);
}

mrb_value
any_to_str(mrb_state *mrb, mrb_value v) {
	return mrb_any_to_s(mrb, v);
}

mrb_float
get_float(mrb_value v) {
	return mrb_float(v);
}

mrb_int
get_fixnum(mrb_value v) {
	return mrb_fixnum(v);
}

mrb_sym
get_symbol(mrb_value v) {
	return mrb_symbol(v);
}

int
is_fixnum(mrb_value v) {
	return mrb_fixnum_p(v);
}

int
is_float(mrb_value v) {
	return mrb_float_p(v);
}

int
is_undef(mrb_value v) {
	return mrb_undef_p(v);
}

int
is_nil(mrb_value v) {
	return mrb_nil_p(v);
}

int
is_symbol(mrb_value v) {
	return mrb_symbol_p(v);
}

int
is_array(mrb_value v) {
	return mrb_array_p(v);
}

int
is_string(mrb_value v) {
	return mrb_string_p(v);
}

int
is_hash(mrb_value v) {
	return mrb_hash_p(v);
}

int
is_voidp(mrb_value v) {
	return mrb_voidp_p(v);
}

int
is_bool(mrb_value v) {
	return mrb_bool(v) != 0;
}

mrb_value
get_ary_entry(mrb_value ary, int index) {
	return mrb_ary_entry(ary, ((mrb_int)index));
}
*/
import "C"

import (
	"fmt"
	_ "log"
	"reflect"
	"unsafe"
)

// ruby2go converts a value in mruby to an interface{} in Go.
func ruby2go(ctx *Context, v C.mrb_value) interface{} {
	switch C.my_type(v) {
	case C.MRB_TT_FALSE:
		if C.is_nil(v) != 0 {
			return nil
		}
		return false
	case C.MRB_TT_FREE:
		return nil
	case C.MRB_TT_TRUE:
		return true
	case C.MRB_TT_FIXNUM:
		return int(C.get_fixnum(v))
	case C.MRB_TT_SYMBOL:
		// Return symbol as string
		return C.GoString(C.mrb_sym2name(ctx.mrb, C.get_symbol(v)))
	case C.MRB_TT_UNDEF:
		return nil
	case C.MRB_TT_FLOAT:
		return float64(C.get_float(v))
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
		return rubyArray2Slice(ctx, v)
	case C.MRB_TT_HASH:
		return rubyHash2Map(ctx, v)
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

func go2ruby(ctx *Context, v interface{}) C.mrb_value {
	kv := reflect.ValueOf(v)
	switch kv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return C.mrb_fixnum_value(C.mrb_int(kv.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return C.mrb_fixnum_value(C.mrb_int(kv.Int()))
	case reflect.Float32, reflect.Float64:
		return C.mrb_float_value(C.mrb_float(kv.Float()))
	case reflect.String:
		cs := C.CString(kv.String())
		defer C.free(unsafe.Pointer(cs))
		return C.mrb_str_new_cstr(ctx.mrb, cs)
	case reflect.Bool:
		if kv.Bool() {
			return C.mrb_true_value()
		}
		return C.mrb_false_value()
	case reflect.Array, reflect.Slice:
		ary := C.mrb_ary_new(ctx.mrb)
		for i := 0; i < kv.Len(); i++ {
			C.mrb_ary_push(ctx.mrb, ary, go2ruby(ctx, kv.Index(i).Interface()))
		}
		return ary
	case reflect.Map:
		hsh := C.mrb_hash_new(ctx.mrb)
		for _, key := range kv.MapKeys() {
			value := kv.MapIndex(key)
			C.mrb_hash_set(ctx.mrb, hsh,
				go2ruby(ctx, key.String()),
				go2ruby(ctx, value.Interface()))
		}
		return hsh
	case reflect.Interface:
		return go2ruby(ctx, kv.Elem().Interface())
	}
	return C.mrb_nil_value()
}

// rubyType encapsulates information about a mrb_value
// in a Go struct.
type rubyType struct {
	TypeName  string
	ClassName string
	Value     interface{}
}

// String returns a representation of the rubyType as a string.
func (typ rubyType) String() string {
	return fmt.Sprintf("<RubyType{Type:%s,Class:%s,Value:%#v}>",
		typ.TypeName, typ.ClassName, typ.Value)
}

// rubyTypeOf converts a value in mruby to an interface{} in Go.
func rubyTypeOf(ctx *Context, v C.mrb_value) *rubyType {
	typ := &rubyType{TypeName: "<unknown>"}
	typ.ClassName = C.GoString(C.mrb_obj_classname(ctx.mrb, v))
	typ.Value = ruby2go(ctx, v)

	switch C.my_type(v) {
	case C.MRB_TT_FALSE:
		typ.TypeName = "MRB_TT_FALSE"
	case C.MRB_TT_FREE:
		typ.TypeName = "MRB_TT_FREE"
	case C.MRB_TT_TRUE:
		typ.TypeName = "MRB_TT_TRUE"
	case C.MRB_TT_FIXNUM:
		typ.TypeName = "MRB_TT_FIXNUM"
	case C.MRB_TT_SYMBOL:
		typ.TypeName = "MRB_TT_SYMBOL"
	case C.MRB_TT_UNDEF:
		typ.TypeName = "MRB_TT_UNDEF"
	case C.MRB_TT_FLOAT:
		typ.TypeName = "MRB_TT_FLOAT"
	case C.MRB_TT_VOIDP:
		typ.TypeName = "MRB_TT_VOIDP"
	case C.MRB_TT_OBJECT:
		typ.TypeName = "MRB_TT_OBJECT"
	case C.MRB_TT_CLASS:
		typ.TypeName = "MRB_TT_CLASS"
	case C.MRB_TT_MODULE:
		typ.TypeName = "MRB_TT_MODULE"
	case C.MRB_TT_ICLASS:
		typ.TypeName = "MRB_TT_ICLASS"
	case C.MRB_TT_SCLASS:
		typ.TypeName = "MRB_TT_SCLASS"
	case C.MRB_TT_PROC:
		typ.TypeName = "MRB_TT_PROC"
	case C.MRB_TT_ARRAY:
		typ.TypeName = "MRB_TT_ARRAY"
	case C.MRB_TT_HASH:
		typ.TypeName = "MRB_TT_HASH"
	case C.MRB_TT_STRING:
		typ.TypeName = "MRB_TT_STRING"
	case C.MRB_TT_RANGE:
		typ.TypeName = "MRB_TT_RANGE"
	case C.MRB_TT_EXCEPTION:
		typ.TypeName = "MRB_TT_EXCEPTION"
	case C.MRB_TT_FILE:
		typ.TypeName = "MRB_TT_FILE"
	case C.MRB_TT_ENV:
		typ.TypeName = "MRB_TT_ENV"
	case C.MRB_TT_DATA:
		typ.TypeName = "MRB_TT_DATA"
	case C.MRB_TT_MAXDEFINE:
		typ.TypeName = "MRB_TT_MAXDEFINE"
	}
	return typ
}

// rubyArray2Slice takes a Ruby array and turns it into a Go slice.
func rubyArray2Slice(ctx *Context, mary C.mrb_value) []interface{} {
	ary := make([]interface{}, 0)

	for i := 0; i < int(C.mrb_ary_len(ctx.mrb, mary)); i++ {
		v := C.get_ary_entry(mary, C.int(i))
		value := ruby2go(ctx, v)
		ary = append(ary, value)
	}

	return ary
}

// rubyHash2Map converts a Ruby hash to a Go map.
func rubyHash2Map(ctx *Context, mhash C.mrb_value) map[string]interface{} {
	hsh := make(map[string]interface{})

	mkeys := C.mrb_hash_keys(ctx.mrb, mhash)

	for i := 0; i < int(C.mrb_ary_len(ctx.mrb, mkeys)); i++ {
		mkey := C.get_ary_entry(mkeys, C.int(i))
		mvalue := C.mrb_hash_get(ctx.mrb, mhash, mkey)

		key := ""
		if C.is_symbol(mkey) != 0 {
			key = C.GoString(C.mrb_string_value_ptr(ctx.mrb, mkey))
		} else {
			key = C.GoString(C.mrb_string_value_ptr(ctx.mrb, mkey))
		}

		hsh[key] = ruby2go(ctx, mvalue)
	}

	//log.Printf(C.GoString(C.mrb_string_value_ptr(ctx.mrb, C.any_to_str(ctx.mrb, v))))

	return hsh
}
