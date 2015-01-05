// Copyright 2013-2015 Oliver Eilhard.
// Use of this source code is governed by the MIT LICENSE that
// can be found in the MIT-LICENSE file included in the project.

package mruby

import (
	"testing"
)

func TestNewContext(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}
}

func TestNewContextWithOptions(t *testing.T) {
	ctx := NewContext(SetFilename("test.rb"), SetNoExec(true))
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}
	if ctx.filename != "test.rb" {
		t.Errorf("expected filename %q; got: %q", "test.rb", ctx.filename)
	}
	if ctx.noExec != true {
		t.Errorf("expected noExec = %v; got: %v", true, ctx.noExec)
	}
}

func TestLoadString(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	val, err := ctx.LoadString("'Hello world'")
	if err != nil {
		t.Fatal(err)
	}
	if !val.IsString() {
		t.Fatalf("expected value to be a string; got: %v", val.Type())
	}
	s, err := val.ToString()
	if err != nil {
		t.Fatal(err)
	}
	if s != "Hello world" {
		t.Errorf("expected %q; got: %q", "Hello world", s)
	}
}

func TestLoadStringResult(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	s, err := ctx.LoadStringResult("'Hello world'")
	if err != nil {
		t.Fatal(err)
	}
	if s != "Hello world" {
		t.Errorf("expected %q; got: %q", "Hello world", s)
	}
}

func TestLoadStringWithArgs(t *testing.T) {
	ctx := NewContext()

	var res interface{}
	var err error

	res, err = ctx.LoadStringResult("ARGV[0]", "Oliver", "Sandra")
	if err != nil {
		t.Fatal(err)
	}
	if res != "Oliver" {
		t.Errorf("expected %v, got %v", "Oliver", res)
	}

	res, err = ctx.LoadStringResult("ARGV[1]", "Oliver", "Sandra")
	if err != nil {
		t.Fatal(err)
	}
	if res != "Sandra" {
		t.Errorf("expected %v, got %v", "Sandra", res)
	}

	res, err = ctx.LoadStringResult("ARGV.inject { |x,y| x+y }", 1, 2, 3.5)
	if err != nil {
		t.Fatal(err)
	}
	if res != float64(6.5) {
		t.Errorf("expected %v, got %v", float64(6.5), res)
	}

	dict := map[string]interface{}{
		"name": "Oliver",
		"age":  23,
	}
	res, err = ctx.LoadStringResult(`ARGV[0]["age"]`, dict)
	if err != nil {
		t.Fatal(err)
	}
	if res != 23 {
		t.Errorf("expected %v, got %v", 23, res)
	}
}

func TestNoExec(t *testing.T) {
	ctx := NewContext(SetNoExec(true))
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	proc, err := ctx.LoadString("'Hello world'")
	if err != nil {
		t.Fatal(err)
	}
	if !proc.IsProc() {
		t.Fatalf("expected value to be a Proc; got: %v", proc.Type())
	}
	val, err := proc.Run()
	if err != nil {
		t.Fatal(err)
	}
	s, err := val.ToString()
	if err != nil {
		t.Fatal(err)
	}
	if s != "Hello world" {
		t.Errorf("expected %q; got: %q", "Hello world", s)
	}
}

func TestNoExecError(t *testing.T) {
	ctx := NewContext(SetNoExec(true))
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	proc, err := ctx.LoadString("raise 'kaboom'")
	if err != nil {
		t.Fatal(err)
	}
	if !proc.IsProc() {
		t.Fatalf("expected value to be a Proc; got: %v", proc.Type())
	}
	_, err = proc.Run()
	if err == nil {
		t.Fatalf("expected error")
	}
	e, ok := err.(*RunError)
	if !ok {
		t.Fatal("expected RunError")
	}
	expected := "kaboom"
	if e.Message != expected {
		t.Errorf("expected %q; got: %q", expected, e.Message)
	}
	if e.Error() != expected {
		t.Errorf("expected %q; got: %q", expected, e.Error())
	}
}

func TestIntReturnsFixnum(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	f32 := float32(11.5)
	f64 := float64(12.3)

	tests := []struct {
		in   interface{}
		want interface{}
	}{
		{f32, 11.5},
		{f64, 12.3},
		{&f32, 11.5},
		{&f64, 12.3},
		{nil, nil},
	}

	for _, test := range tests {
		val, err := ctx.ToValue(test.in)
		if err != nil {
			t.Fatal(err)
		}
		if test.want != nil && !val.IsFloat() {
			t.Errorf("expected %q; got: %v", "Float", val.Type())
		}
		got, err := ctx.LoadStringResult("ARGV[0]", test.in)
		if err != nil {
			t.Fatal(err)
		}
		if got != test.want {
			t.Errorf("expected %v; got: %v", test.want, got)
		}
	}
}

func TestFloatReturnsFloat(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	i := int(1)
	i8 := int8(2)
	i16 := int16(3)
	i32 := int32(4)
	i64 := int64(5)
	ui := uint(6)
	ui8 := uint8(7)
	ui16 := uint16(8)
	ui32 := uint32(9)
	ui64 := uint64(10)

	tests := []struct {
		in   interface{}
		want interface{}
	}{
		{i, 1},
		{i8, 2},
		{i16, 3},
		{i32, 4},
		{i64, 5},
		{ui, 6},
		{ui8, 7},
		{ui16, 8},
		{ui32, 9},
		{ui64, 10},
		{&i, 1},
		{&i8, 2},
		{&i16, 3},
		{&i32, 4},
		{&i64, 5},
		{&ui, 6},
		{&ui8, 7},
		{&ui16, 8},
		{&ui32, 9},
		{&ui64, 10},
		{nil, nil},
	}

	for _, test := range tests {
		val, err := ctx.ToValue(test.in)
		if err != nil {
			t.Fatal(err)
		}
		if test.want != nil && !val.IsFixnum() {
			t.Errorf("expected %q; got: %v", "Fixnum", val.Type())
		}
		got, err := ctx.LoadStringResult("ARGV[0]", test.in)
		if err != nil {
			t.Fatal(err)
		}
		if got != test.want {
			t.Errorf("expected %v; got: %v", test.want, got)
		}
	}
}

func TestPtrValues(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	i := int(1)
	i8 := int8(2)
	i16 := int16(3)
	i32 := int32(4)
	i64 := int64(5)
	ui := uint(6)
	ui8 := uint8(7)
	ui16 := uint16(8)
	ui32 := uint32(9)
	ui64 := uint64(10)
	f32 := float32(11.5)
	f64 := float64(12.3)
	b := true
	s := "oliver"

	tests := []struct {
		in        interface{}
		want      interface{}
		wantClass string
	}{
		{&i, 1, "Fixnum"},
		{&i8, 2, "Fixnum"},
		{&i16, 3, "Fixnum"},
		{&i32, 4, "Fixnum"},
		{&i64, 5, "Fixnum"},
		{&ui, 6, "Fixnum"},
		{&ui8, 7, "Fixnum"},
		{&ui16, 8, "Fixnum"},
		{&ui32, 9, "Fixnum"},
		{&ui64, 10, "Fixnum"},
		{&f32, 11.5, "Float"},
		{&f64, 12.3, "Float"},
		{&b, true, "TrueClass"},
		{&s, "oliver", "String"},
		{nil, nil, "NilClass"},
	}

	for _, test := range tests {
		val, err := ctx.ToValue(test.in)
		if err != nil {
			t.Fatal(err)
		}
		if val.Type().class != test.wantClass {
			t.Errorf("expected value type %q; got: %v", test.wantClass, val.Type())
		}
		got, err := ctx.LoadStringResult("ARGV[0]", test.in)
		if err != nil {
			t.Fatal(err)
		}
		if got != test.want {
			t.Errorf("expected %v; got: %v", test.want, got)
		}
	}
}

func TestPtrValuesInHash(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	i := int(1)
	i8 := int8(2)
	i16 := int16(3)
	i32 := int32(4)
	i64 := int64(5)
	ui := uint(6)
	ui8 := uint8(7)
	ui16 := uint16(8)
	ui32 := uint32(9)
	ui64 := uint64(10)
	f32 := float32(11.5)
	f64 := float64(12.3)
	b := true
	s := "Oliver"

	in := map[string]interface{}{
		"i":    &i,
		"i8":   &i8,
		"i16":  &i16,
		"i32":  &i32,
		"i64":  &i64,
		"ui":   &ui,
		"ui8":  &ui8,
		"ui16": &ui16,
		"ui32": &ui32,
		"ui64": &ui64,
		"f32":  &f32,
		"f64":  &f64,
		"b":    &b,
		"s":    &s,
		"nil":  nil,
	}

	val, err := ctx.ToValue(in)
	if err != nil {
		t.Fatal(err)
	}
	if val.Type().class != "Hash" {
		t.Errorf("expected value type %q; got: %v", "Hash", val.Type())
	}
	res, err := ctx.LoadStringResult("ARGV[0]", in)
	if err != nil {
		t.Fatal(err)
	}
	hsh, ok := res.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map; got: %v", hsh)
	}
	// Flattens return types to int, for Ruby does only have Fixnum
	if got, ok := hsh["i"].(int); !ok || got != i {
		t.Errorf("expected i=%v; got: %v", i, got)
	}
	if got, ok := hsh["i8"].(int); !ok || got != int(i8) {
		t.Errorf("expected i8=%v; got: %v", i8, got)
	}
	if got, ok := hsh["i16"].(int); !ok || got != int(i16) {
		t.Errorf("expected i16=%v; got: %v", i16, got)
	}
	if got, ok := hsh["i32"].(int); !ok || got != int(i32) {
		t.Errorf("expected i32=%v; got: %v", i32, got)
	}
	if got, ok := hsh["i64"].(int); !ok || got != int(i64) {
		t.Errorf("expected i64=%v; got: %v", i64, got)
	}
	if got, ok := hsh["ui"].(int); !ok || got != int(ui) {
		t.Errorf("expected ui=%v; got: %v", ui, got)
	}
	if got, ok := hsh["ui8"].(int); !ok || got != int(ui8) {
		t.Errorf("expected ui8=%v; got: %v", ui8, got)
	}
	if got, ok := hsh["ui16"].(int); !ok || got != int(ui16) {
		t.Errorf("expected ui16=%v; got: %v", ui16, got)
	}
	if got, ok := hsh["ui32"].(int); !ok || got != int(ui32) {
		t.Errorf("expected ui32=%v; got: %v", ui32, got)
	}
	if got, ok := hsh["ui64"].(int); !ok || got != int(ui64) {
		t.Errorf("expected ui64=%v; got: %v", ui64, got)
	}
	// Flattens return types to float64, for Ruby does only have Float
	if got, ok := hsh["f32"].(float64); !ok || got != float64(f32) {
		t.Errorf("expected f32=%v; got: %v", f32, got)
	}
	if got, ok := hsh["f64"].(float64); !ok || got != float64(f64) {
		t.Errorf("expected f64=%v; got: %v", f64, got)
	}
	if got, ok := hsh["b"].(bool); !ok || got != b {
		t.Errorf("expected b=%v; got: %v", b, got)
	}
	if got, ok := hsh["s"].(string); !ok || got != s {
		t.Errorf("expected s=%v; got: %v", s, got)
	}
	if got := hsh["nil"]; got != nil {
		t.Errorf("expected got=%v; got: %v", nil, got)
	}
}
