#!mruby

# This example only works if you compiled mruby with the sandbox gem
# from mattn: https://github.com/mattn/mruby-sandbox.
# Here's an example of how to use the Sandbox gem:
# https://github.com/mattn/mruby-lingrbot/blob/master/mruby-lingrbot.rb

begin
	@sb = Sandbox.new
rescue => ex
	puts "Please install the Sandbox gem from https://github.com/mattn/mruby-sandbox to run this example"
	exit # you need the mruby-exit gem for Kernel#exit
end

puts @sb.eval "'Hello World'"
