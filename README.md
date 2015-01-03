# MRuby for Google Go

The [mruby-go](https://github.com/olivere/mruby-go) package enables you
to use the embedded Ruby interpreter [mruby](https://github.com/mruby/mruby)
inside Google Go projects.

## Prerequisites

We're using mruby-go with Google Go 1.4, but it should work with 1.1+ and
tip as well. The mruby interpreter went
[1.0.0](http://www.mruby.org/releases/2014/02/09/mruby-1.0.0-released.html)
on 9th Feb 2014, and we have tested it with that version successfully.
You can find a tarball of the latest version of mruby [here](https://github.com/mruby/mruby/releases).

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
    Version: 1.1.0
    Libs: -L${libdir} -lmruby
    Libs.private: -lm
    Cflags: -I${includedir}

Make sure that `pkg-config --list-all` includes `mruby`.

Notice that you configuration of mruby is quite simple, you can (and
probably should) review all settings with regard to your environment.
There's a [whole section on configuring mruby](#mruby-config).

Next, build mruby-go:

    go build

Then run tests with:

    go test

and benchmarks with:

    go test -test.bench .

## Getting started

See the examples directory or the tests for example usage.

# <a name="mruby-config">Configuring mruby</a>

There is a short summary of how to configure mruby in the
[mruby repository](https://github.com/mruby/mruby),
especially in
the [INSTALL](https://github.com/mruby/mruby/blob/master/INSTALL) file.
However, we will try to summarize configuration in our own words.

1. Download mruby from [here](https://github.com/mruby/mruby/releases).
1. Untar the repository and cd into the expanded directory.
1. Run `ruby ./minirake`
1. That should succeed and create some binaries in `bin` and some libs
   in `lib`.

If you want to add some (mruby-specific) gems to your mruby installation,
edit the `./build_config.rb` and add a line to it, e.g.
`conf.gem :git => 'git@github.com:iij/mruby-io.git', :branch => master`.
Then compile again via `ruby ./minirake`. We're using e.g.
[mattn/mruby-json](https://github.com/mattn/mruby-json) and
[mattn/mruby-sandbox](https://github.com/mattn/mruby-sandbox).
The Internet Initiative Japan Inc. has
[quite a lot of mruby gems](https://github.com/iij) on GitHub. Pick what
you need.

For configuring details, we used to edit `./include/mrbconf.h`. We're not
sure if this is the correct (or right) way of doing things, but it works.


# Support

Feel free to send pull requests.

# License

This package has a [MIT-LICENSE](https://github.com/olivere/mruby-go/MIT-LICENSE).
