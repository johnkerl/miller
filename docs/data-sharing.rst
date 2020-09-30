..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Mixing with other languages
================================================================

As discussed in the section on :doc:`file-formats`, Miller supports several different file formats. Different tools are good at different things, so it's important to be able to move data into and out of other languages. **CSV** and **JSON** are well-known, of course; here are some examples using **DKVP** format, with **Ruby** and **Python**. Last, we show how to use arbitrary **shell commands** to extend functionality beyond Miller's domain-specific language.

DKVP I/O in Python
----------------------------------------------------------------

Here are the I/O routines:

::

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

And here is an example using them:

::

    $ cat polyglot-dkvp-io/example.py
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

Run as-is:

::

    $ python polyglot-dkvp-io/example.py < data/small
    a=pan,b=pan,i=1,y=0.7268028627434533,ab=panpan,iy=1.7268028627434533,ta=str,tb=str,ti=int,ty=float,tab=str,tiy=float
    a=eks,b=pan,i=2,y=0.5221511083334797,ab=ekspan,iy=2.5221511083334796,ta=str,tb=str,ti=int,ty=float,tab=str,tiy=float
    a=wye,b=wye,i=3,y=0.33831852551664776,ab=wyewye,iy=3.3383185255166477,ta=str,tb=str,ti=int,ty=float,tab=str,tiy=float
    a=eks,b=wye,i=4,y=0.13418874328430463,ab=ekswye,iy=4.134188743284304,ta=str,tb=str,ti=int,ty=float,tab=str,tiy=float
    a=wye,b=pan,i=5,y=0.8636244699032729,ab=wyepan,iy=5.863624469903273,ta=str,tb=str,ti=int,ty=float,tab=str,tiy=float

Run as-is, then pipe to Miller for pretty-printing:

::

