// Copyright 2013-2015 Oliver Eilhard.
// Use of this source code is governed by the MIT LICENSE that
// can be found in the MIT-LICENSE file included in the project.

package mruby

import (
	"testing"
)

func TestValueType(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	tests := []struct {
		Input      string
		ShouldFail bool
		Type       string
		Class      string
	}{
		{`true`, false, "MRB_TT_TRUE", "TrueClass"},
		{`false`, false, "MRB_TT_FALSE", "FalseClass"},
		{`''`, false, "MRB_TT_STRING", "String"},
		{`""`, false, "MRB_TT_STRING", "String"},
		{`"abc"`, false, "MRB_TT_STRING", "String"},
		{`:abc`, false, "MRB_TT_SYMBOL", "Symbol"},
		{`1`, false, "MRB_TT_FIXNUM", "Fixnum"},
		{`1.5`, false, "MRB_TT_FLOAT", "Float"},
		{`nil`, false, "MRB_TT_FALSE", "NilClass"},
		{"['Oliver', 2, 42.3, true, nil]", false, "MRB_TT_ARRAY", "Array"},
		{"{:name => 'Oliver', :age => 21}", false, "MRB_TT_HASH", "Hash"},
		{"raise 'kaboom'", true, "MRB_TT_EXCEPTION", "Exception"},
	}

	for _, test := range tests {
		val, err := ctx.LoadString(test.Input)
		if err != nil {
			// Should it fail?
			if !test.ShouldFail {
				t.Fatal(err)
			}
		} else {
			// Should succeed
			typ := val.Type()
			if typ.class != test.Class {
				t.Errorf("expected class %q; got: %q", test.Class, typ.class)
			}
			if typ.typ != test.Type {
				t.Errorf("expected type %q; got: %q", test.Type, typ.typ)
			}
		}
	}
}
