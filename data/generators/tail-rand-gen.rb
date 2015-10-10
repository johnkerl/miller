#!/usr/bin/env ruby
# For playing with stats1/2 -s
$stdout.sync = true
while true
  x = rand()
  y = rand()
  xy = x*y
  puts "x=#{x},y=#{y},xy=#{xy}"
  sleep 0.1
end
