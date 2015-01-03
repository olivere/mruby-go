// Copyright 2013-2015 Oliver Eilhard.
// Use of this source code is governed by the MIT LICENSE that
// can be found in the MIT-LICENSE file included in the project.

package mruby

import (
	"reflect"
	"testing"
)

func TestToValue(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	tests := []struct {
		Input    string
		Expected interface{}
		Failed   bool
		Error    string
	}{
		{`true`, true, false, ""},
		{`false`, false, false, ""},
		{`''`, string(""), false, ""},
		{`""`, string(""), false, ""},
		{`"abc"`, string("abc"), false, ""},
		{`:abc`, string("abc"), false, ""},
		{`1`, int(1), false, ""},
		{`1.5`, float64(1.5), false, ""},
		{`nil`, nil, false, ""},
		{"['Oliver', 2, 42.3, true, nil]", []interface{}{"Oliver", 2, float64(42.3), true, nil}, false, ""},
		{"{:name => 'Oliver', :age => 21}",
			map[string]interface{}{
				"name": "Oliver",
				"age":  21,
			},
			false,
			""},
		{"{:name => 'Oliver', 'age' => 21, address: {city: 'Munich'}}",
			map[string]interface{}{
				"name": "Oliver",
				"age":  21,
				"address": map[string]interface{}{
					"city": "Munich",
				},
			},
			false,
			""},
		{"raise 'kaboom'", nil, true, "kaboom"},
	}

	for _, test := range tests {
		val, err := ctx.LoadString(test.Input)
		if err != nil {
			// Should it fail?
			if !test.Failed {
				t.Fatal(err)
			} else if test.Error != err.Error() {
				t.Errorf("expected error %q; got: %q", test.Error, err.Error())
			}
		} else {
			// Should succeed
			got, err := val.ToInterface()
			if err != nil {
				t.Fatal(err)
			}
			if reflect.TypeOf(got) != reflect.TypeOf(test.Expected) {
				t.Errorf("expected %v; got: %v", reflect.TypeOf(test.Expected), reflect.TypeOf(got))
			}
			if !reflect.DeepEqual(got, test.Expected) {
				t.Errorf("expected %v; got: %v", test.Expected, got)
			}
		}
	}
}

func TestValueStringer(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	val, err := ctx.LoadString("'Hello'")
	if err != nil {
		t.Fatal(err)
	}
	expected := "Value{Value:\"Hello\",ValueType{Type:MRB_TT_STRING,Class:String}}"
	if val.String() != expected {
		t.Errorf("expected %q; got: %q", expected, val.String())
	}
}

func TestNilType(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	val, err := ctx.LoadString("nil")
	if err != nil {
		t.Fatal(err)
	}
	if !val.IsNil() {
		t.Errorf("expected type NilClass; got: %v", val.Type())
	}

	res, err := val.ToInterface()
	if err != nil {
		t.Fatal(err)
	}
	if res != nil {
		t.Errorf("expected %v; got: %v", nil, res)
	}
}

func TestTrueType(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	val, err := ctx.LoadString("1 == 1")
	if err != nil {
		t.Fatal(err)
	}
	if !val.IsBool() {
		t.Errorf("expected type True; got: %v", val.Type())
	}

	flag, err := val.ToBool()
	if err != nil {
		t.Fatal(err)
	}
	if !flag {
		t.Errorf("expected %v; got: %v", true, flag)
	}
}

func TestFalseType(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	val, err := ctx.LoadString("1 != 1")
	if err != nil {
		t.Fatal(err)
	}
	if !val.IsBool() {
		t.Errorf("expected type False; got: %v", val.Type())
	}

	flag, err := val.ToBool()
	if err != nil {
		t.Fatal(err)
	}
	if flag {
		t.Errorf("expected %v; got: %v", false, flag)
	}
}

func TestFixnumType(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	val, err := ctx.LoadString("1+2")
	if err != nil {
		t.Fatal(err)
	}
	if !val.IsFixnum() {
		t.Errorf("expected type Fixnum; got: %v", val.Type())
	}
	i, err := val.ToInt()
	if err != nil {
		t.Fatal(err)
	}
	if i != 3 {
		t.Errorf("expected %d; got: %d", 3, i)
	}

	i8, err := val.ToInt8()
	if err != nil {
		t.Fatal(err)
	}
	if i8 != 3 {
		t.Errorf("expected %d; got: %d", 3, i8)
	}

	i16, err := val.ToInt()
	if err != nil {
		t.Fatal(err)
	}
	if i16 != 3 {
		t.Errorf("expected %d; got: %d", 3, i16)
	}

	i32, err := val.ToInt()
	if err != nil {
		t.Fatal(err)
	}
	if i32 != 3 {
		t.Errorf("expected %d; got: %d", 3, i32)
	}

	i64, err := val.ToInt()
	if err != nil {
		t.Fatal(err)
	}
	if i64 != 3 {
		t.Errorf("expected %d; got: %d", 3, i64)
	}
}

func TestFloatType(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	val, err := ctx.LoadString("1.5+2.25")
	if err != nil {
		t.Fatal(err)
	}
	if !val.IsFloat() {
		t.Errorf("expected type Float; got: %v", val.Type())
	}
	f32, err := val.ToFloat32()
	if err != nil {
		t.Fatal(err)
	}
	if f32 != 3.75 {
		t.Errorf("expected %v; got: %v", 3.75, f32)
	}

	f64, err := val.ToFloat64()
	if err != nil {
		t.Fatal(err)
	}
	if f64 != 3.75 {
		t.Errorf("expected %v; got: %v", 3.75, f64)
	}
}

func TestStringType(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	val, err := ctx.LoadString("'Hello'")
	if err != nil {
		t.Fatal(err)
	}
	if !val.IsString() {
		t.Errorf("expected type String; got: %v", val.Type())
	}
	s, err := val.ToString()
	if err != nil {
		t.Fatal(err)
	}
	if s != "Hello" {
		t.Errorf("expected %q; got: %q", "Hello", s)
	}
}

func TestSymbolType(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	val, err := ctx.LoadString(":Hello")
	if err != nil {
		t.Fatal(err)
	}
	if !val.IsSymbol() {
		t.Errorf("expected type Symbol; got: %v", val.Type())
	}
	s, err := val.ToString()
	if err != nil {
		t.Fatal(err)
	}
	if s != "Hello" {
		t.Errorf("expected %q; got: %q", "Hello", s)
	}
}

func TestArrayType(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	val, err := ctx.LoadString("[1,2,'Oliver']")
	if err != nil {
		t.Fatal(err)
	}
	if !val.IsArray() {
		t.Errorf("expected type Array; got: %v", val.Type())
	}

	got, err := val.ToArray()
	if err != nil {
		t.Fatal(err)
	}
	expected := []interface{}{1, 2, "Oliver"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v; got: %v", expected, got)
	}
}

func TestMapType(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	val, err := ctx.LoadString("{city: 'Munich', :name => 'Oliver', age: 21}")
	if err != nil {
		t.Fatal(err)
	}
	if !val.IsHash() {
		t.Errorf("expected type Hash; got: %v", val.Type())
	}

	got, err := val.ToMap()
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]interface{}{"city": "Munich", "name": "Oliver", "age": 21}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v; got: %v", expected, got)
	}
}

func TestExceptionType(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	_, err := ctx.LoadString("raise 'bang bang'")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestArrayOfHashes(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	in := []map[string]interface{}{
		{
			"a": 1,
			"b": 2,
		},
		{
			"a": 17,
			"b": 4,
		},
	}

	script := `
ARGV[0].each do |hsh|
	hsh[:c] = hsh["a"] + hsh["b"]
end
`

	val, err := ctx.LoadString(script, in)
	if err != nil {
		t.Fatal(err)
	}
	if !val.IsArray() {
		t.Errorf("expected type Array; got: %v", val.Type())
	}

	got, err := val.ToArray()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("expected %d entries; got: %d", 2, len(got))
	}

	ent, ok := got[0].(map[string]interface{})
	if !ok {
		t.Fatal("expected entry to be a map")
	}
	a, found := ent["a"]
	if !found {
		t.Errorf("expected entry %q", "a")
	}
	if a != 1 {
		t.Errorf("expected entry %q = %d; got: %d", "a", 1, a)
	}
	b, found := ent["b"]
	if !found {
		t.Errorf("expected entry %q", "b")
	}
	if b != 2 {
		t.Errorf("expected entry %q = %d; got: %d", "b", 2, b)
	}
	c, found := ent["c"]
	if !found {
		t.Errorf("expected entry %q", "c")
	}
	if c != 3 {
		t.Errorf("expected entry %q = %d; got: %d", "c", 3, c)
	}

	ent, ok = got[1].(map[string]interface{})
	if !ok {
		t.Fatal("expected entry to be a map")
	}
	a, found = ent["a"]
	if !found {
		t.Errorf("expected entry %q", "a")
	}
	if a != 17 {
		t.Errorf("expected entry %q = %d; got: %d", "a", 17, a)
	}
	b, found = ent["b"]
	if !found {
		t.Errorf("expected entry %q", "b")
	}
	if b != 4 {
		t.Errorf("expected entry %q = %d; got: %d", "b", 4, b)
	}
	c, found = ent["c"]
	if !found {
		t.Errorf("expected entry %q", "c")
	}
	if c != 21 {
		t.Errorf("expected entry %q = %d; got: %d", "c", 21, c)
	}
}
