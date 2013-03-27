// Copyright 2013 Oliver Eilhard.
// Use of this source code is governed by the MIT LICENSE that
// can be found in the MIT-LICENSE file included in the project.
package mruby_test

import (
	mruby "github.com/olivere/mruby-go"
	"reflect"
	"testing"
)

func TestLoadString(t *testing.T) {
	ctx := mruby.NewContext()
	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}

	res, err := ctx.LoadString("'Hello world'")
	if err != nil {
		t.Fatal(err)
	}
	if res != "Hello world" {
		t.Errorf("expected 'Hello world', got %v", res)
	}
}

func TestLoadStringWithDifferentResults(t *testing.T) {
	ctx := mruby.NewContext()

	tests := []struct {
		Name     string
		Code     string
		Expected interface{}
		Failure  bool
	}{
		{"nil", "nil", nil, false},
		{"string", "'Oliver'", "Oliver", false},
		{"string2", "\"Oliver\"", "Oliver", false},
		{"zero int", "0", 0, false},
		{"zero float", "0.0", float64(0.0), false},
		{"int", "42", 42, false},
		{"float", "42.3", float64(42.3), false},
		{"symbol", ":oliver", "oliver", false},
		{"array", "['Oliver', 2, 42.3, true, nil]", []interface{}{
			"Oliver", 2, float64(42.3), true, nil,
		}, false},
		{"hash", "{:name => 'Oliver', :age => 21}", map[string]interface{}{
			"name": "Oliver",
			"age":  21,
		}, false},
		{"complex hash", "{:name => 'Oliver', 'age' => 21, address: {city: 'Munich'}}", map[string]interface{}{
			"name": "Oliver",
			"age":  21,
			"address": map[string]interface{}{
				"city": "Munich",
			},
		}, false},
		{"exception", "raise 'kaboom'", nil, true},
	}

	for _, test := range tests {
		res, err := ctx.LoadString(test.Code)
		if test.Failure {
			// test should have failed
			if err == nil {
				t.Fatalf("test %s:\n  should have failed", test.Name)
			}
			if reflect.TypeOf(res) != reflect.TypeOf(test.Expected) {
				t.Errorf("test %s:\n  expected type %v, got type %v",
					test.Name, reflect.TypeOf(test.Expected), reflect.TypeOf(res))
			}
		} else {
			// test should have succeeded
			if err != nil {
				t.Fatalf("test %s:\n  %v", test.Name, err)
			}
			if reflect.TypeOf(res) != reflect.TypeOf(test.Expected) {
				t.Errorf("test %s:\n  expected type %v, got type %v",
					test.Name, reflect.TypeOf(test.Expected), reflect.TypeOf(res))
			}
			if !reflect.DeepEqual(res, test.Expected) {
				t.Errorf("test %s:\n  expected %v, got %v", test.Name, test.Expected, res)
			}
		}
	}
}

func TestLoadStringWithArgs(t *testing.T) {
	ctx := mruby.NewContext()

	var res interface{}
	var err error

	res, err = ctx.LoadString("ARGV[0]", "Oliver", "Sandra")
	if err != nil {
		t.Fatal(err)
	}
	if res != "Oliver" {
		t.Errorf("expected %v, got %v", "Oliver", res)
	}

	res, err = ctx.LoadString("ARGV[1]", "Oliver", "Sandra")
	if err != nil {
		t.Fatal(err)
	}
	if res != "Sandra" {
		t.Errorf("expected %v, got %v", "Sandra", res)
	}

	res, err = ctx.LoadString("ARGV.inject { |x,y| x+y }", 1, 2, 3.5)
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
	res, err = ctx.LoadString(`ARGV[0]["age"]`, dict)
	if err != nil {
		t.Fatal(err)
	}
	if res != 23 {
		t.Errorf("expected %v, got %v", 23, res)
	}
}

func TestParse(t *testing.T) {
	ctx := mruby.NewContext()

	parser, err := ctx.Parse("'Parsed'")
	if err != nil {
		t.Fatal(err)
	}
	res, err := parser.Run()
	if err != nil {
		t.Fatal(err)
	}
	if res != "Parsed" {
		t.Errorf("expected 'Parsed', got %v", res)
	}
}

func TestExceptionOnRun(t *testing.T) {
	ctx := mruby.NewContext()

	parser, err := ctx.Parse("raise 'kaboom'")
	if err != nil {
		t.Fatal(err)
	}
	res, err := parser.Run()
	if err == nil {
		t.Fatal("expected exception message as error, got nil")
	}
	if err.Error() != "kaboom" {
		t.Errorf("expected exception message 'kaboom', got %v", err.Error())
	}
	if res != nil {
		t.Fatal("expected result to be nil, got %v", res)
	}
}

func BenchmarkParse(b *testing.B) {
	ctx := mruby.NewContext()

	code := `
def concat(a, b)
	a + b
end

concat "Hello", "World"
`

	parser, err := ctx.Parse(code)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		res, err := parser.Run()
		if err != nil {
			b.Fatal(err)
		}
		if res != "HelloWorld" {
			b.Errorf("run %d: expected 'HelloWorld', got %v", i, res)
		}
	}
}
