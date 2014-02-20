# Ruby for Google Go

The [mruby-go](https://github.com/olivere/mruby-go) package enables you
to use the embedded Ruby interpreter [mruby](https://github.com/mruby/mruby)
inside Google Go projects.

## Prerequisites

We're using mruby-go with Google Go 1.2, but it should work with 1.1+ and
tip as well. The mruby interpreter went
[1.0.0](http://www.mruby.org/releases/2014/02/09/mruby-1.0.0-released.html)
on 9th Feb 2014, and we have tested it with that version successfully.

We're using mruby-go in production and are pretty confident that it behaves
as a good citizen. But we don't give any guarantees.

Read the tests to get a feel of what works and what doesn't.

## Installation

Building mruby-go depends on [pkg-config](http://www.freedesktop.org/wiki/Software/pkg-config/).
Make sure to have a valid pkg-config file for mruby in your
`PKG_CONFIG_PATH`. I've compiled mruby locally, so here's my `mruby.pc`:

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

## Getting started

See the examples directory or the tests for example usage.

# Support

Feel free to send pull requests.

# License

This package has a [MIT-LICENSE](https://github.com/olivere/mruby-go/MIT-LICENSE).
