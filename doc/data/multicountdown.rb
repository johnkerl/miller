#!/usr/bin/ruby
$stdout.sync = true

# Simulates some multi-threaded program making progress over time on some set of tasks
counts = {'red' => 3842, 'green' => 1224, 'blue' => 2697, 'purple' => 979}
colors = counts.keys
n = colors.length

start_time = Time::now.to_f

upsec = 0.0
while counts.any?{|color,count| count > 0}
  upsec += 0.001 * rand(250)
  color = colors[rand(n)]
  count = counts[color]
  delta = rand(40)
  count -= delta
  count = [0, count].max
  counts[color] = count
  puts "upsec=#{upsec},color=#{color},count=#{counts[color]}"
end
