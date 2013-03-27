# Ruby for Google Go

The mruby-go package enables users to use the embedded Ruby interpreter
[mruby](git://github.com/mruby/mruby.git) inside Google Go projects.

## Status

This is a work in progress. Read the tests to get a feel of what works
and what doesn't.

## Getting started

First you need to compile the mruby package. It has been added as a
submodule of this project. The original repository is located at
https://github.com/mruby/mruby.git. After compiling mruby, just
compile the package.

Here's a shortcut:

    make mruby
    make compile

See the examples directory or the tests for example usage.

# License

This package has a [MIT-LICENSE](https://github.com/olivere/mruby-go/MIT-LICENSE).
