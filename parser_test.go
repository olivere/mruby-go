package mruby_test

import (
	mruby "github.com/olivere/mruby-go"
	"testing"
)

func TestParse(t *testing.T) {
	ctx := mruby.NewContext()

	parser := mruby.NewParser(ctx)

	res, err := parser.Parse("p 'Parsed'")
	if err != nil {
		t.Fatal(err)
	}
	if res != "Parsed" {
		t.Errorf("expected 'Parsed', got %v", res)
	}
}
