# Ruby for Google Go

The mruby-go package enables users to use the embedded Ruby interpreter
[mruby](git://github.com/mruby/mruby.git) inside Google Go projects.

## Status

This is a work in progress. Read the tests to get a feel of what works
and what doesn't.

## Installation

Use the latest tip of Google Go to make this work. It all depends
on [issue 4069](https://code.google.com/p/go/issues/detail?id=4069),
which will be included in Go 1.1.

Building mruby-go depends on pkg-config. Make sure to have a valid
pkg-config file for mruby in your `PKG_CONFIG_PATH`. I've compiled
mruby locally, so here's my `mruby.pc`:

    prefix=<home>/ext/mruby
    exec_prefix=${prefix}
    libdir=${exec_prefix}/build/host/lib
    includedir=${prefix}/include
    
    Name: libmruby
    Description: Embedded Ruby (mruby)
    Version: 0.1.0
    Libs: -L${libdir} -lmruby
    Libs.private: -lm
    Cflags: -I${includedir}

Make sure that `pkg-config --list-all` includes `mruby`.

Next, build mruby-go:

    go build

Then run tests with:

    go test

and benchmarks with:

    go test -test.bench .

## mruby

For convenience, the [mruby repository](https://github.com/mruby/mruby.git)
is added as a submodule.

## Getting started

See the examples directory or the tests for example usage.

# License

This package has a [MIT-LICENSE](https://github.com/olivere/mruby-go/MIT-LICENSE).
