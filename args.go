// Copyright 2013-2015 Oliver Eilhard.
// Use of this source code is governed by the MIT LICENSE that
// can be found in the MIT-LICENSE file included in the project.

package mruby

import (
	"sync"
)

/*
#cgo pkg-config: mruby
#include "mruby_go.h"
*/
import "C"

// Args is used to specify the number of arguments a Go extension method
// will use. Use one of the ArgsAny, ArgsNone, ArgsRequired, ArgsOptional,
// or ArgsArg helper functions to initialize it.
type Args C.mrb_aspec

// ArgsAny specifies a function with any number of arguments.
func ArgsAny() Args {
	return Args(C.args_any())
}

// ArgsNone specifies a function that accepts no arguments.
func ArgsNone() Args {
	return Args(C.args_none())
}

// ArgsRequired specifies a function that accepts a number of required arguments.
func ArgsRequired(required int) Args {
	return Args(C.args_req(C.int(required)))
}

// ArgsOptional specifies a function that accepts a number of options arguments.
func ArgsOptional(optional int) Args {
	return Args(C.args_opt(C.int(optional)))
}

// ArgsArg specifies a function that a number of required and optional arguments.
func ArgsArg(required, optional int) Args {
	return Args(C.args_arg(C.int(required), C.int(optional)))
}

var getArgAccumulator []*C.mrb_value
var getArgLock sync.Mutex

//export go_get_arg_append
func go_get_arg_append(v *C.mrb_value) {
	getArgAccumulator = append(getArgAccumulator, v)
}
