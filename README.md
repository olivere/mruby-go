# Ruby for Google Go

The mruby-go package enables users to use the embedded Ruby interpreter
[mruby](git://github.com/mruby/mruby.git) inside Google Go projects.

## Status

This is a work in progress. Read the tests to get a feel of what works
and what doesn't.

## Installation

The mruby repository at https://github.com/mruby/mruby.git is added
as a submodule. However, you must make sure that include files and
library paths can be resolved.

You can do this either by setting C_INCLUDE_PATH and LIBRARY_PATH 
manually or adding mruby as a package to your system. The latter is
obviously preferred, but probably isn't available as mruby is not 
yet mature.

Start compiling with:

    make compile

## Getting started

See the examples directory or the tests for example usage.

# License

This package has a [MIT-LICENSE](https://github.com/olivere/mruby-go/MIT-LICENSE).
