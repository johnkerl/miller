#!/usr/bin/env ruby

5000000000.times do
  u = rand < 0.5 ? 0 : 1
  v = rand < 0.1 ? 0 : 1
  puts "#{u} #{v}"
end
puts "3 5"
