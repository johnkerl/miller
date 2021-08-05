#!/usr/bin/env python

# ================================================================
# Example of DKVP I/O using Python.
#
# Key point: Use Miller for what it's good at; pass data into/out of tools in
# other languages to do what they're good at.
#
#   bash$ python -i dkvp_io.py
#
#   # READ
#   >>> map = dkvpline2map('x=1,y=2', '=', ',')
#   >>> map
#   OrderedDict([('x', '1'), ('y', '2')])
#
#   # MODIFY
#   >>> map['z'] = map['x'] + map['y']
#   >>> map
#   OrderedDict([('x', '1'), ('y', '2'), ('z', 3)])
#
#   # WRITE
#   >>> line = map2dkvpline(map, '=', ',')
#   >>> line
#   'x=1,y=2,z=3'
#
# ================================================================

import re
import collections

# ----------------------------------------------------------------
# ips and ifs (input pair separator and input field separator) are nominally '=' and ','.
def dkvpline2map(line, ips, ifs):
	pairs = re.split(ifs, line)
	map = collections.OrderedDict()
	for pair in pairs:
		key, value = re.split(ips, pair, 1)

		# Type inference:
		try:
			value = int(value)
		except:
			try:
				value = float(value)
			except:
				pass

		map[key] = value
	return map

# ----------------------------------------------------------------
# ops and ofs (output pair separator and output field separator) are nominally '=' and ','.
def map2dkvpline(map , ops, ofs):
	line = ''
	pairs = []
	for key in map:
		pairs.append(str(key) + ops + str(map[key]))
	return str.join(ofs, pairs)
