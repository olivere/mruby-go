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
	"errors"
	"fmt"
)

// Value is used for encoding/decoding data types from MRuby to Go and vice versa.
type Value struct {
	ctx *Context
	v   C.mrb_value
}

// Type returns type information.
func (v Value) Type() ValueType {
	return newValueType(v)
}

// String returns a string representation of the value.
func (v Value) String() string {
	goval, _ := v.ToInterface()
	return fmt.Sprintf("Value{Value:%#v,%v}", goval, v.Type())
}

// NilValue returns a mruby nil value.
func NilValue(ctx *Context) Value {
	return Value{ctx: ctx, v: C.mrb_nil_value()}
}

// -- Detectors: Is... --

// IsNil indicates whether this value is nil.
func (v Value) IsNil() bool {
	return C.is_nil(v.v) != 0
}

// IsBool indicates whether this value is a bool.
//
// Notice: In Ruby, everything that is not false is a boolean.
// However, this method respects the Go type bool that has two distinct
// values: true or false.
func (v Value) IsBool() bool {
	switch C.my_type(v.v) {
	case C.MRB_TT_FALSE, C.MRB_TT_TRUE:
		return true
	default:
		return false
	}
}

// IsFixnum indicates whether this value is a fixed number.
//
// Notice: Ruby only has fixed numbers and floats. So there is no equivalent
// to e.g. int8 or uint16.
func (v Value) IsFixnum() bool {
	switch C.my_type(v.v) {
	case C.MRB_TT_FIXNUM:
		return true
	default:
		return false
	}
}

// IsFloat indicates whether this value is a float.
//
// Notice: Ruby only has fixed numbers and floats. So this method returns
// true if the Ruby type is a float and therefore can be converted to float32.
func (v Value) IsFloat() bool {
	switch C.my_type(v.v) {
	case C.MRB_TT_FLOAT:
		return true
	default:
		return false
	}
}

// IsString indicates whether this value can be converted to a string.
func (v Value) IsString() bool {
	switch C.my_type(v.v) {
	case C.MRB_TT_STRING:
		return true
	default:
		return false
	}
}

// IsSymbol indicates whether this value is a Ruby symbol.
func (v Value) IsSymbol() bool {
	switch C.my_type(v.v) {
	case C.MRB_TT_SYMBOL:
		return true
	default:
		return false
	}
}

// IsArray indicates whether this value is an array.
func (v Value) IsArray() bool {
	switch C.my_type(v.v) {
	case C.MRB_TT_ARRAY:
		return true
	default:
		return false
	}
}

// IsHash indicates whether this value is a hash.
func (v Value) IsHash() bool {
	switch C.my_type(v.v) {
	case C.MRB_TT_HASH:
		return true
	default:
		return false
	}
}

/*
// IsException indicates whether this value is an exception.
func (v Value) IsException() bool {
	switch C.my_type(v.v) {
	case C.MRB_TT_EXCEPTION:
		return true
	default:
		return false
	}
}
*/

// IsProc indicates whether this value is a Proc.
func (v Value) IsProc() bool {
	switch typ := C.my_type(v.v); typ {
	case C.MRB_TT_PROC:
		return true
	default:
		return false
	}
}

// -- Conversion to Go --

// ToBool treats this value as a bool and returns its value.
// If the value is not a bool, an error is returned.
func (v Value) ToBool() (bool, error) {
	switch typ := C.my_type(v.v); typ {
	case C.MRB_TT_FALSE:
		return false, nil
	case C.MRB_TT_TRUE:
		return true, nil
	default:
		return false, fmt.Errorf("value is not a bool but %v", v.Type())
	}
}

// ToInt treats this value as an int and returns its value.
// If the value is not a Ruby Fixnum, an error is returned.
func (v Value) ToInt() (int, error) {
	switch typ := C.my_type(v.v); typ {
	case C.MRB_TT_FIXNUM:
		return int(C.get_fixnum(v.v)), nil
	default:
		return 0, fmt.Errorf("value is not an int but %v", v.Type())
	}
}

// ToInt8 treats this value as an int8 and returns its value.
// If the value is not a Ruby Fixnum, an error is returned.
func (v Value) ToInt8() (int8, error) {
	switch typ := C.my_type(v.v); typ {
	case C.MRB_TT_FIXNUM:
		return int8(C.get_fixnum(v.v)), nil
	default:
		return 0, fmt.Errorf("value is not an int8 but %v", v.Type())
	}
}

// ToInt16 treats this value as an int16 and returns its value.
// If the value is not a Ruby Fixnum, an error is returned.
func (v Value) ToInt16() (int16, error) {
	switch typ := C.my_type(v.v); typ {
	case C.MRB_TT_FIXNUM:
		return int16(C.get_fixnum(v.v)), nil
	default:
		return 0, fmt.Errorf("value is not an int16 but %v", v.Type())
	}
}

// ToInt32 treats this value as an int32 and returns its value.
// If the value is not a Ruby Fixnum, an error is returned.
func (v Value) ToInt32() (int32, error) {
	switch typ := C.my_type(v.v); typ {
	case C.MRB_TT_FIXNUM:
		return int32(C.get_fixnum(v.v)), nil
	default:
		return 0, fmt.Errorf("value is not an int32 but %v", v.Type())
	}
}

// ToInt64 treats this value as an int64 and returns its value.
// If the value is not a Ruby Fixnum, an error is returned.
func (v Value) ToInt64() (int64, error) {
	switch typ := C.my_type(v.v); typ {
	case C.MRB_TT_FIXNUM:
		return int64(C.get_fixnum(v.v)), nil
	default:
		return 0, fmt.Errorf("value is not an int64 but %v", v.Type())
	}
}

// ToUint treats this value as an uint and returns its value.
// If the value is not a Ruby Fixnum, an error is returned.
func (v Value) ToUint() (uint, error) {
	switch typ := C.my_type(v.v); typ {
	case C.MRB_TT_FIXNUM:
		return uint(C.get_fixnum(v.v)), nil
	default:
		return 0, fmt.Errorf("value is not a uint but %v", v.Type())
	}
}

// ToUint8 treats this value as an uint8 and returns its value.
// If the value is not a Ruby Fixnum, an error is returned.
func (v Value) ToUint8() (uint8, error) {
	switch typ := C.my_type(v.v); typ {
	case C.MRB_TT_FIXNUM:
		return uint8(C.get_fixnum(v.v)), nil
	default:
		return 0, fmt.Errorf("value is not a uint8 but %v", v.Type())
	}
}

// ToUint16 treats this value as an uint16 and returns its value.
// If the value is not a Ruby Fixnum, an error is returned.
func (v Value) ToUint16() (uint16, error) {
	switch typ := C.my_type(v.v); typ {
	case C.MRB_TT_FIXNUM:
		return uint16(C.get_fixnum(v.v)), nil
	default:
		return 0, fmt.Errorf("value is not a uint16 but %v", v.Type())
	}
}

// ToUint32 treats this value as an uint32 and returns its value.
// If the value is not a Ruby Fixnum, an error is returned.
func (v Value) ToUint32() (uint32, error) {
	switch typ := C.my_type(v.v); typ {
	case C.MRB_TT_FIXNUM:
		return uint32(C.get_fixnum(v.v)), nil
	default:
		return 0, fmt.Errorf("value is not a uint32 but %v", v.Type())
	}
}

// ToUint64 treats this value as an uint64 and returns its value.
// If the value is not a Ruby Fixnum, an error is returned.
func (v Value) ToUint64() (uint64, error) {
	switch typ := C.my_type(v.v); typ {
	case C.MRB_TT_FIXNUM:
		return uint64(C.get_fixnum(v.v)), nil
	default:
		return 0, fmt.Errorf("value is not a uint64 but %v", v.Type())
	}
}

// ToFloat32 treats this value as a float32 and returns its value.
// If the value is not a Ruby Float, an error is returned.
func (v Value) ToFloat32() (float32, error) {
	switch typ := C.my_type(v.v); typ {
	case C.MRB_TT_FLOAT:
		return float32(C.get_float(v.v)), nil
	default:
		return 0.0, fmt.Errorf("value is not a float32 but %v", v.Type())
	}
}

// ToFloat64 treats this value as a float64 and returns its value.
// If the value is not a Ruby Float, an error is returned.
func (v Value) ToFloat64() (float64, error) {
	switch typ := C.my_type(v.v); typ {
	case C.MRB_TT_FLOAT:
		return float64(C.get_float(v.v)), nil
	default:
		return 0.0, fmt.Errorf("value is not a float64 but %v", v.Type())
	}
}

// ToString returns a string for MRuby types String and Symbol.
// If the value is not a Ruby String or Symbol, an error is returned.
func (v Value) ToString() (string, error) {
	switch typ := C.my_type(v.v); typ {
	case C.MRB_TT_STRING:
		return C.GoString(C.mrb_string_value_ptr(v.ctx.mrb, v.v)), nil
	case C.MRB_TT_SYMBOL:
		return C.GoString(C.mrb_sym2name(v.ctx.mrb, C.get_symbol(v.v))), nil
	default:
		return "", fmt.Errorf("value is not a string but %v", v.Type())
	}
}

// ToArray treats this value as an array and returns its values.
// If the value is not a Ruby Array, an error is returned.
func (v Value) ToArray() ([]interface{}, error) {
	switch typ := C.my_type(v.v); typ {
	case C.MRB_TT_ARRAY:
		return v.mrbArrayToSlice()
	default:
		return nil, fmt.Errorf("value is not an array but %v", v.Type())
	}
}

// ToMap treats this value as a hash and returns its values as a map.
// If the value is not a Ruby Hash, an error is returned.
func (v Value) ToMap() (map[string]interface{}, error) {
	switch typ := C.my_type(v.v); typ {
	case C.MRB_TT_HASH:
		return v.mrbHashToMap()
	default:
		return nil, fmt.Errorf("value is not a hash but %v", v.Type())
	}
}

// ToInterface will return the Go-equivalent of the Ruby value.
// It will only handle the following Ruby types: TrueClass, FalseClass,
// NilClass, Fixnum, Float, Symbol, String, Array, and Hash.
// All other Ruby types return nil.
func (v Value) ToInterface() (interface{}, error) {
	switch C.my_type(v.v) {
	case C.MRB_TT_FALSE:
		if C.is_nil(v.v) != 0 {
			return nil, nil
		}
		return false, nil
	case C.MRB_TT_FREE:
		return nil, nil
	case C.MRB_TT_TRUE:
		return true, nil
	case C.MRB_TT_FIXNUM:
		return int(C.get_fixnum(v.v)), nil
	case C.MRB_TT_SYMBOL:
		// Return symbol as string
		return C.GoString(C.mrb_sym2name(v.ctx.mrb, C.get_symbol(v.v))), nil
	case C.MRB_TT_UNDEF:
		return nil, nil
	case C.MRB_TT_FLOAT:
		return float64(C.get_float(v.v)), nil
	case C.MRB_TT_CPTR:
		return nil, nil
	case C.MRB_TT_OBJECT:
		return nil, nil
	case C.MRB_TT_CLASS:
		return nil, nil
	case C.MRB_TT_MODULE:
		return nil, nil
	case C.MRB_TT_ICLASS:
		return nil, nil
	case C.MRB_TT_SCLASS:
		return nil, nil
	case C.MRB_TT_PROC:
		return nil, nil
	case C.MRB_TT_ARRAY:
		return v.mrbArrayToSlice()
	case C.MRB_TT_HASH:
		return v.mrbHashToMap()
	case C.MRB_TT_STRING:
		return C.GoString(C.mrb_string_value_ptr(v.ctx.mrb, v.v)), nil
	case C.MRB_TT_RANGE:
		return nil, nil
	case C.MRB_TT_EXCEPTION:
		return nil, nil
	case C.MRB_TT_FILE:
		return nil, nil
	case C.MRB_TT_ENV:
		return nil, nil
	case C.MRB_TT_DATA:
		return nil, nil
	case C.MRB_TT_FIBER:
		return nil, nil
	case C.MRB_TT_MAXDEFINE:
		return nil, nil
	}
	return nil, ErrInvalidType
}

// mrbArrayToSlice takes the value (which has to be an array) and returns
// the elements of the Ruby array as an array of Go values.
func (v Value) mrbArrayToSlice() ([]interface{}, error) {
	goary := make([]interface{}, 0)
	for i := 0; i < int(C.mrb_ary_len(v.ctx.mrb, v.v)); i++ {
		mrbval := C.get_ary_entry(v.v, C.int(i))
		aryval := Value{ctx: v.ctx, v: mrbval}
		goval, err := aryval.ToInterface()
		if err != nil {
			return nil, err
		}
		goary = append(goary, goval)
	}
	return goary, nil
}

// mrbHashToMap takes the value (which has to be a Ruby Hash) and returns
// the key/value pairs as a Go map|string]interface{}. Ruby keys which are
// symbols are turned into strings.
func (v Value) mrbHashToMap() (map[string]interface{}, error) {
	gomap := make(map[string]interface{})
	mrbkeys := C.mrb_hash_keys(v.ctx.mrb, v.v)
	for i := 0; i < int(C.mrb_ary_len(v.ctx.mrb, mrbkeys)); i++ {
		mrbkey := C.get_ary_entry(mrbkeys, C.int(i))
		mrbvalue := C.mrb_hash_get(v.ctx.mrb, v.v, mrbkey)

		key := ""
		if C.is_symbol(mrbkey) != 0 {
			key = C.GoString(C.mrb_string_value_ptr(v.ctx.mrb, mrbkey))
		} else {
			key = C.GoString(C.mrb_string_value_ptr(v.ctx.mrb, mrbkey))
		}

		val := Value{ctx: v.ctx, v: mrbvalue}
		goval, err := val.ToInterface()
		if err != nil {
			return nil, err
		}
		gomap[key] = goval
	}
	return gomap, nil
}

// Run runs the code given that it is a reference to a Proc.
func (v Value) Run() (Value, error) {
	if !v.IsProc() {
		return NilValue(v.ctx), errors.New("value is not a Proc")
	}
	proc := C.my_mrb_proc_ptr(v.v)
	newv := C.mrb_run(v.ctx.mrb, proc, v.v)
	if C.has_exception(v.ctx.mrb) != 0 {
		return NilValue(v.ctx), newRunError(v.ctx, true)
	}
	return Value{ctx: v.ctx, v: newv}, nil
}
