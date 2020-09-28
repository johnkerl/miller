#!/usr/bin/env python

import sys
import re
import copy
import dkvp_io

while True:
	# Read the original record:
	line = sys.stdin.readline().strip()
	if line == '':
		break
	map = dkvp_io.dkvpline2map(line, '=', ',')

	# Drop a field:
	map.pop('x')

	# Compute some new fields:
	map['ab'] = map['a'] + map['b']
	map['iy'] = map['i'] + map['y']

	# Add new fields which show type of each already-existing field:
	omap = copy.copy(map) # since otherwise the for-loop will modify what it loops over
	keys = omap.keys()
	for key in keys:
		# Convert "<type 'int'>" to just "int", etc.:
		type_string = str(map[key].__class__)
		type_string = re.sub("<type '", "", type_string) # python2
		type_string = re.sub("<class '", "", type_string) # python3
		type_string = re.sub("'>", "", type_string)
		map['t'+key] = type_string

	# Write the modified record:
	print(dkvp_io.map2dkvpline(map, '=', ','))
