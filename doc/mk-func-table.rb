#!/usr/bin/env ruby

lines = `mlr --list-all-functions-as-table`
counter = 0
puts '<table border=1>'
lines.split("\n").each do |line|
  counter = counter + 1
  if line =~ /^? :.*/ # has an extra whitespace which throws off the naive split
    ig, nore, function_class, nargs = line.split(/\s+/)
    name = '? :'
  else
    name, function_class, nargs = line.split(/\s+/)
  end

  if counter == 1
    puts "<tr class=\"mlrbg\">"
    puts "<th>#{name}</th> <th>#{function_class}</th> <th>#{nargs}</th>"
  else
    puts "<tr>"
    puts "<td><a href=\"##{name}\">#{name}</a></td> <td>#{function_class}</td> <td>#{nargs}</td>"
  end
  puts '</tr>'
end
puts '</table>'
