module Helpers
  def self.hello
    "Hello!"
  end

  def self.escape_html(s)
    s
  end
end

puts Helpers.hello
puts Helpers.escape_html("<esca&e>")

