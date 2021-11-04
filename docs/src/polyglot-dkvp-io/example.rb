#!/usr/bin/env ruby

require 'dkvp_io'

ARGF.each do |line|
  # Read the original record:
  map = dkvpline2map(line.chomp, '=', ',')

  # Drop a field:
  map.delete('x')

  # Compute some new fields:
  map['ab'] = map['a'] + map['b']
  map['iy'] = map['i'] + map['y']

  # Add new fields which show type of each already-existing field:
  keys = map.keys
  keys.each do |key|
    map['t'+key] = map[key].class
  end

  # Write the modified record:
  puts map2dkvpline(map, '=', ',')
end
