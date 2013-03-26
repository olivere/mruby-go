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

func TestLoadString(t *testing.T) {
	ctx := mruby.NewContext()

	if ctx == nil {
		t.Fatal("expected NewContext() to be != nil")
	}
	ctx.LoadString("p 'Hello world'")
}
