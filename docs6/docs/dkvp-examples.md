<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# DKVP I/O examples

## DKVP I/O in Python

Here are the I/O routines:

<pre class="pre-non-highlight-non-pair">
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
</pre>

And here is an example using them:

<pre class="pre-highlight-in-pair">
<b>cat polyglot-dkvp-io/example.py</b>
</pre>
<pre class="pre-non-highlight-in-pair">
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
</pre>

Run as-is:

<pre class="pre-highlight-in-pair">
<b>python polyglot-dkvp-io/example.py < data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=1,y=0.7268028627434533,ab=panpan,iy=1.7268028627434533,ta=str,tb=str,ti=int,ty=float,tab=str,tiy=float
a=eks,b=pan,i=2,y=0.5221511083334797,ab=ekspan,iy=2.5221511083334796,ta=str,tb=str,ti=int,ty=float,tab=str,tiy=float
a=wye,b=wye,i=3,y=0.33831852551664776,ab=wyewye,iy=3.3383185255166477,ta=str,tb=str,ti=int,ty=float,tab=str,tiy=float
a=eks,b=wye,i=4,y=0.13418874328430463,ab=ekswye,iy=4.134188743284304,ta=str,tb=str,ti=int,ty=float,tab=str,tiy=float
a=wye,b=pan,i=5,y=0.8636244699032729,ab=wyepan,iy=5.863624469903273,ta=str,tb=str,ti=int,ty=float,tab=str,tiy=float
</pre>

Run as-is, then pipe to Miller for pretty-printing:

<pre class="pre-highlight-in-pair">
<b>python polyglot-dkvp-io/example.py < data/small | mlr --opprint cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i y                   ab     iy                 ta  tb  ti  ty    tab tiy
pan pan 1 0.7268028627434533  panpan 1.7268028627434533 str str int float str float
eks pan 2 0.5221511083334797  ekspan 2.5221511083334796 str str int float str float
wye wye 3 0.33831852551664776 wyewye 3.3383185255166477 str str int float str float
eks wye 4 0.13418874328430463 ekswye 4.134188743284304  str str int float str float
wye pan 5 0.8636244699032729  wyepan 5.863624469903273  str str int float str float
</pre>

## DKVP I/O in Ruby

Here are the I/O routines:

<pre class="pre-non-highlight-non-pair">
#!/usr/bin/env ruby

# ================================================================
# Example of DKVP I/O using Ruby.
#
# Key point: Use Miller for what it's good at; pass data into/out of tools in
# other languages to do what they're good at.
#
#   bash$ irb -I. -r dkvp_io.rb
#
#   # READ
#   irb(main):001:0> map = dkvpline2map('x=1,y=2', '=', ',')
#   => {"x"=>"1", "y"=>"2"}
#
#   # MODIFY
#   irb(main):001:0> map['z'] = map['x'] + map['y']
#   => 3
#
#   # WRITE
#   irb(main):002:0> line = map2dkvpline(map, '=', ',')
#   => "x=1,y=2,z=3"
#
# ================================================================

# ----------------------------------------------------------------
# ips and ifs (input pair separator and input field separator) are nominally '=' and ','.
def dkvpline2map(line, ips, ifs)
  map = {}
  line.split(ifs).each do |pair|
    (k, v) = pair.split(ips, 2)

    # Type inference:
    begin
      v = Integer(v)
    rescue ArgumentError
      begin
        v = Float(v)
      rescue ArgumentError
        # Leave as string
      end
    end

    map[k] = v
  end
  map
end

# ----------------------------------------------------------------
# ops and ofs (output pair separator and output field separator) are nominally '=' and ','.
def map2dkvpline(map, ops, ofs)
  map.collect{|k,v| k.to_s + ops + v.to_s}.join(ofs)
end
</pre>

And here is an example using them:

<pre class="pre-highlight-in-pair">
<b>cat polyglot-dkvp-io/example.rb</b>
</pre>
<pre class="pre-non-highlight-in-pair">
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
</pre>

Run as-is:

<pre class="pre-highlight-in-pair">
<b>ruby -I./polyglot-dkvp-io polyglot-dkvp-io/example.rb data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=1,y=0.7268028627434533,ab=panpan,iy=1.7268028627434533,ta=String,tb=String,ti=Integer,ty=Float,tab=String,tiy=Float
a=eks,b=pan,i=2,y=0.5221511083334797,ab=ekspan,iy=2.5221511083334796,ta=String,tb=String,ti=Integer,ty=Float,tab=String,tiy=Float
a=wye,b=wye,i=3,y=0.33831852551664776,ab=wyewye,iy=3.3383185255166477,ta=String,tb=String,ti=Integer,ty=Float,tab=String,tiy=Float
a=eks,b=wye,i=4,y=0.13418874328430463,ab=ekswye,iy=4.134188743284304,ta=String,tb=String,ti=Integer,ty=Float,tab=String,tiy=Float
a=wye,b=pan,i=5,y=0.8636244699032729,ab=wyepan,iy=5.863624469903273,ta=String,tb=String,ti=Integer,ty=Float,tab=String,tiy=Float
</pre>

Run as-is, then pipe to Miller for pretty-printing:

<pre class="pre-highlight-in-pair">
<b>ruby -I./polyglot-dkvp-io polyglot-dkvp-io/example.rb data/small | mlr --opprint cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i y                   ab     iy                 ta     tb     ti      ty    tab    tiy
pan pan 1 0.7268028627434533  panpan 1.7268028627434533 String String Integer Float String Float
eks pan 2 0.5221511083334797  ekspan 2.5221511083334796 String String Integer Float String Float
wye wye 3 0.33831852551664776 wyewye 3.3383185255166477 String String Integer Float String Float
eks wye 4 0.13418874328430463 ekswye 4.134188743284304  String String Integer Float String Float
wye pan 5 0.8636244699032729  wyepan 5.863624469903273  String String Integer Float String Float
</pre>
