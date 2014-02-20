// Copyright 2013-2014 Oliver Eilhard.
// Use of this source code is governed by the MIT LICENSE that
// can be found in the MIT-LICENSE file included in the project.
package mruby_test

import (
	mruby "github.com/olivere/mruby-go"
	"testing"
)

func TestNewContext(t *testing.T) {
	ctx := mruby.NewContext()

	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}
}
