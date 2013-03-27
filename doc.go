// Copyright 2013 Oliver Eilhard.
// Use of this source code is governed by the MIT LICENSE that
// can be found in the MIT-LICENSE file included in the project.

/*
Package mruby embeds the Ruby programming language into your Go project.
mruby is a lightweight implementation of the Ruby language complying to
(part of) the ISO standard. mruby can be linked and embedded within an
application. This package makes mruby available in your Go projects,
effectively enabling Ruby as a scripting language within your Go code.

Notice that mruby is a work in progress. But with the backing of the
original authors of the Ruby programming language, it is likely to succeed
in the long run as a way to embed Ruby within larger applications.

You can find the mruby source code at https://github.com/mruby/mruby.

Introduction

Let us all start with a very simple example.

	ctx := mruby.NewContext()
	ctx.LoadString("p 'Hello world'")

If successful, these two lines of code will print "Hello world!" on stdout.
*/
package mruby
