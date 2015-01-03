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

var (
	// ErrInvalidType is returned when the package cannot convert a Ruby
	// type to the Go equivalent.
	ErrInvalidType = errors.New("invalid type")
)

// RunError is used to indicate errors while running Ruby code.
type RunError struct {
	Message string // Message details
}

func newRunError(ctx *Context, resetException bool) *RunError {
	err := &RunError{}
	err.Message = C.GoString(C.get_exception_message(ctx.mrb))
	if resetException {
		C.reset_exception(ctx.mrb)
	}
	return err
}

// Error returns the error as a string.
func (e *RunError) Error() string {
	return e.Message
}

// ParseError is used to indicate errors while parsing Ruby code.
type ParseError struct {
	Line    int    // Line number
	Message string // Message details
}

// Error returns the error as a string.
func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error: line %d: %s", e.Line, e.Message)
}
