#!/usr/bin/env ruby

lines = `mlr --list-all-functions-as-table`
counter = 0
lines.split("\n").each do |line|
  counter = counter + 1
  if line =~ /^? :.*/ # has an extra whitespace which throws off the naive split
    ig, nore, function_class, nargs = line.split(/\s+/)
    name = '? :'
  else
    name, function_class, nargs = line.split(/\s+/)
  end

  if counter == 1
    puts "+--------------+-----------+-------+"
    puts "| #{name} | #{function_class} | #{nargs} |"
    puts "+==============+===========+=======+"
  else
    puts "| a href=\"##{name}\" #{name} | #{function_class} | #{nargs} |"
    puts "+--------------+-----------+-------+"
  end
end
