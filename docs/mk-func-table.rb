#!/usr/bin/env ruby

# Sphinx tables need to be precisely lined up, unlike Markdown tables which
# are more flexible. For example this is fine Markdown:
#
#   | a | b | i | x | y |
#   | --- | --- | --- | --- | --- |
#   | pan | pan | 1 | 0.3467901443380824 | 0.7268028627434533 |
#   | eks | pan | 2 | 0.7586799647899636 | 0.5221511083334797 |
#   | wye | wye | 3 | 0.20460330576630303 | 0.33831852551664776 |
#   | eks | wye | 4 | 0.38139939387114097 | 0.13418874328430463 |
#
# but Sphinx wants
#
#   +-----+-----+---+---------------------+---------------------+
#   | a   | b   | i | x                   | y                   |
#   +=====+=====+===+=====================+=====================+
#   | pan | pan | 1 | 0.3467901443380824  | 0.7268028627434533  |
#   +-----+-----+---+---------------------+---------------------+
#   | eks | pan | 2 | 0.7586799647899636  | 0.5221511083334797  |
#   +-----+-----+---+---------------------+---------------------+
#   | wye | wye | 3 | 0.20460330576630303 | 0.33831852551664776 |
#   +-----+-----+---+---------------------+---------------------+
#   | eks | wye | 4 | 0.38139939387114097 | 0.13418874328430463 |
#   +-----+-----+---+---------------------+---------------------+
#
# So we need to build up a list of tuples, then compute the max length down
# each column, then write a fixed number of -/= per row.

lines = `mlr --list-all-functions-as-table`

# ----------------------------------------------------------------
# Pass 1

max_name_length = 1
max_function_class_length = 1
max_nargs_length = 1

rows = []
lines.split("\n").each do |line|
  if line =~ /^? :.*/ # has an extra whitespace which throws off the naive split
    ig, nore, function_class, nargs = line.split(/\s+/)
    name = '? :'
  else
    name, function_class, nargs = line.split(/\s+/)
  end

  name = '``' + name + '``'

  if max_name_length < name.length
    max_name_length = name.length
  end
  if max_function_class_length < function_class.length
    max_function_class_length = function_class.length
  end
  if max_nargs_length < nargs.length
    max_nargs_length = nargs.length
  end

  # TODO: format name with page-internal link
  row = [name, function_class, nargs]
  rows.append(row)
end

# ----------------------------------------------------------------
# Pass 2

name_dashes = '-' * max_name_length
function_class_dashes = '-' * max_function_class_length
nargs_dashes = '-' * max_nargs_length

name_equals = '=' * max_name_length
function_class_equals = '=' * max_function_class_length
nargs_equals = '=' * max_nargs_length

dashes_line = "+-#{name_dashes}-+-#{function_class_dashes}-+-#{nargs_dashes}-+"
equals_line = "+=#{name_equals}=+=#{function_class_equals}=+=#{nargs_equals}=+"

counter = 0
rows.each do |row|
  name, function_class, nargs = row

  counter = counter + 1
  if counter == 1
    puts dashes_line
    puts "| #{name.ljust(max_name_length)} | #{function_class.ljust(max_function_class_length)} | #{nargs.ljust(max_nargs_length)} |"
    puts equals_line
  else
    puts "| #{name.ljust(max_name_length)} | #{function_class.ljust(max_function_class_length)} | #{nargs.ljust(max_nargs_length)} |"
    puts dashes_line
  end
end

