package mruby_test

import (
	mruby "github.com/olivere/mruby-go"
	"testing"
)

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
