// Copyright 2013-2015 Oliver Eilhard.
// Use of this source code is governed by the MIT LICENSE that
// can be found in the MIT-LICENSE file included in the project.

/*
Package mruby embeds the Ruby programming language into your Go project.
mruby is a lightweight implementation of the Ruby language complying to
(part of) the ISO standard. mruby can be linked and embedded within an
application. This package makes mruby available in your Go projects,
effectively enabling Ruby as a scripting language within your Go code.

mruby went 1.0.0 on 9th Feb 2014. The current version of mruby is 1.1.0,
and mruby-go is tested with it.

You can find the mruby source code at https://github.com/mruby/mruby.
Tarballs are available from https://github.com/mruby/mruby/releases.

Example:

	ctx := mruby.NewContext()
	ctx.LoadString("p 'Hello world'")

If successful, these two lines of code will print "Hello world!" to stdout.
*/
package mruby
