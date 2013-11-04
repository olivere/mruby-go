package main

import (
	"fmt"
	"io/ioutil"
	"os"

	mruby "github.com/olivere/mruby-go"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: run <rubyfile>\n")
		os.Exit(2)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	script, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	ctx := mruby.NewContext()
	res, err := ctx.LoadString(string(script))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if res != nil {
	    fmt.Fprintf(os.Stdout, "%v\n", res)
	}
}
