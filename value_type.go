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
	"fmt"
)

// ValueType encapsulates information about a Value.
type ValueType struct {
	typ   string
	class string
}

// String returns a representation of the rubyType as a string.
func (vt ValueType) String() string {
	return fmt.Sprintf("ValueType{Type:%s,Class:%s}", vt.typ, vt.class)
}

// newValueType returns information about a value.
func newValueType(v Value) ValueType {
	vt := ValueType{}
	vt.class = C.GoString(C.mrb_obj_classname(v.ctx.mrb, v.v))
	switch C.my_type(v.v) {
	case C.MRB_TT_FALSE:
		vt.typ = "MRB_TT_FALSE"
	case C.MRB_TT_FREE:
		vt.typ = "MRB_TT_FREE"
	case C.MRB_TT_TRUE:
		vt.typ = "MRB_TT_TRUE"
	case C.MRB_TT_FIXNUM:
		vt.typ = "MRB_TT_FIXNUM"
	case C.MRB_TT_SYMBOL:
		vt.typ = "MRB_TT_SYMBOL"
	case C.MRB_TT_UNDEF:
		vt.typ = "MRB_TT_UNDEF"
	case C.MRB_TT_FLOAT:
		vt.typ = "MRB_TT_FLOAT"
	case C.MRB_TT_CPTR:
		vt.typ = "MRB_TT_CPTR"
	case C.MRB_TT_OBJECT:
		vt.typ = "MRB_TT_OBJECT"
	case C.MRB_TT_CLASS:
		vt.typ = "MRB_TT_CLASS"
	case C.MRB_TT_MODULE:
		vt.typ = "MRB_TT_MODULE"
	case C.MRB_TT_ICLASS:
		vt.typ = "MRB_TT_ICLASS"
	case C.MRB_TT_SCLASS:
		vt.typ = "MRB_TT_SCLASS"
	case C.MRB_TT_PROC:
		vt.typ = "MRB_TT_PROC"
	case C.MRB_TT_ARRAY:
		vt.typ = "MRB_TT_ARRAY"
	case C.MRB_TT_HASH:
		vt.typ = "MRB_TT_HASH"
	case C.MRB_TT_STRING:
		vt.typ = "MRB_TT_STRING"
	case C.MRB_TT_RANGE:
		vt.typ = "MRB_TT_RANGE"
	case C.MRB_TT_EXCEPTION:
		vt.typ = "MRB_TT_EXCEPTION"
	case C.MRB_TT_FILE:
		vt.typ = "MRB_TT_FILE"
	case C.MRB_TT_ENV:
		vt.typ = "MRB_TT_ENV"
	case C.MRB_TT_DATA:
		vt.typ = "MRB_TT_DATA"
	case C.MRB_TT_FIBER:
		vt.typ = "MRB_TT_FIBER"
	case C.MRB_TT_MAXDEFINE:
		vt.typ = "MRB_TT_MAXDEFINE"
	default:
		vt.typ = "<unknown>"
	}
	return vt
}
