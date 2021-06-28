..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Mixing with other languages
================================================================

As discussed in the section on :doc:`file-formats`, Miller supports several different file formats. Different tools are good at different things, so it's important to be able to move data into and out of other languages. **CSV** and **JSON** are well-known, of course; here are some examples using **DKVP** format, with **Ruby** and **Python**. Last, we show how to use arbitrary **shell commands** to extend functionality beyond Miller's domain-specific language.

DKVP I/O in Python
----------------------------------------------------------------

Here are the I/O routines:

.. code-block:: none

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

.. code-block:: none
   :emphasize-lines: 1,1

    cat polyglot-dkvp-io/example.py
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

.. code-block:: none
   :emphasize-lines: 1,1

    python polyglot-dkvp-io/example.py < data/small
    a=pan,b=pan,i=1,y=0.7268028627434533,ab=panpan,iy=1.7268028627434533,ta=str,tb=str,ti=int,ty=float,tab=str,tiy=float
    a=eks,b=pan,i=2,y=0.5221511083334797,ab=ekspan,iy=2.5221511083334796,ta=str,tb=str,ti=int,ty=float,tab=str,tiy=float
    a=wye,b=wye,i=3,y=0.33831852551664776,ab=wyewye,iy=3.3383185255166477,ta=str,tb=str,ti=int,ty=float,tab=str,tiy=float
    a=eks,b=wye,i=4,y=0.13418874328430463,ab=ekswye,iy=4.134188743284304,ta=str,tb=str,ti=int,ty=float,tab=str,tiy=float
    a=wye,b=pan,i=5,y=0.8636244699032729,ab=wyepan,iy=5.863624469903273,ta=str,tb=str,ti=int,ty=float,tab=str,tiy=float

Run as-is, then pipe to Miller for pretty-printing:

.. code-block:: none
   :emphasize-lines: 1,1

    python polyglot-dkvp-io/example.py < data/small | mlr --opprint cat
    a   b   i y                   ab     iy                 ta  tb  ti  ty    tab tiy
    pan pan 1 0.7268028627434533  panpan 1.7268028627434533 str str int float str float
    eks pan 2 0.5221511083334797  ekspan 2.5221511083334796 str str int float str float
    wye wye 3 0.33831852551664776 wyewye 3.3383185255166477 str str int float str float
    eks wye 4 0.13418874328430463 ekswye 4.134188743284304  str str int float str float
    wye pan 5 0.8636244699032729  wyepan 5.863624469903273  str str int float str float

DKVP I/O in Ruby
----------------------------------------------------------------

Here are the I/O routines:

.. code-block:: none

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

And here is an example using them:

.. code-block:: none
   :emphasize-lines: 1,1

    cat polyglot-dkvp-io/example.rb
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

Run as-is:

.. code-block:: none
   :emphasize-lines: 1,1

    ruby -I./polyglot-dkvp-io polyglot-dkvp-io/example.rb data/small
    a=pan,b=pan,i=1,y=0.7268028627434533,ab=panpan,iy=1.7268028627434533,ta=String,tb=String,ti=Integer,ty=Float,tab=String,tiy=Float
    a=eks,b=pan,i=2,y=0.5221511083334797,ab=ekspan,iy=2.5221511083334796,ta=String,tb=String,ti=Integer,ty=Float,tab=String,tiy=Float
    a=wye,b=wye,i=3,y=0.33831852551664776,ab=wyewye,iy=3.3383185255166477,ta=String,tb=String,ti=Integer,ty=Float,tab=String,tiy=Float
    a=eks,b=wye,i=4,y=0.13418874328430463,ab=ekswye,iy=4.134188743284304,ta=String,tb=String,ti=Integer,ty=Float,tab=String,tiy=Float
    a=wye,b=pan,i=5,y=0.8636244699032729,ab=wyepan,iy=5.863624469903273,ta=String,tb=String,ti=Integer,ty=Float,tab=String,tiy=Float

Run as-is, then pipe to Miller for pretty-printing:

.. code-block:: none
   :emphasize-lines: 1,1

    ruby -I./polyglot-dkvp-io polyglot-dkvp-io/example.rb data/small | mlr --opprint cat
    a   b   i y                   ab     iy                 ta     tb     ti      ty    tab    tiy
    pan pan 1 0.7268028627434533  panpan 1.7268028627434533 String String Integer Float String Float
    eks pan 2 0.5221511083334797  ekspan 2.5221511083334796 String String Integer Float String Float
    wye wye 3 0.33831852551664776 wyewye 3.3383185255166477 String String Integer Float String Float
    eks wye 4 0.13418874328430463 ekswye 4.134188743284304  String String Integer Float String Float
    wye pan 5 0.8636244699032729  wyepan 5.863624469903273  String String Integer Float String Float

SQL-output examples
----------------------------------------------------------------

Please see :ref:`sql-output-examples`.

SQL-input examples
----------------------------------------------------------------

Please see :ref:`sql-input-examples`.

Running shell commands
----------------------------------------------------------------

The :ref:`reference-dsl-system` DSL function allows you to run a specific shell command and put its output -- minus the final newline -- into a record field. The command itself is any string, either a literal string, or a concatenation of strings, perhaps including other field values or what have you.

.. code-block:: none
   :emphasize-lines: 1,1

    mlr --opprint put '$o = system("echo hello world")' data/small
    a   b   i x                   y                   o
    pan pan 1 0.3467901443380824  0.7268028627434533  hello world
    eks pan 2 0.7586799647899636  0.5221511083334797  hello world
    wye wye 3 0.20460330576630303 0.33831852551664776 hello world
    eks wye 4 0.38139939387114097 0.13418874328430463 hello world
    wye pan 5 0.5732889198020006  0.8636244699032729  hello world

.. code-block:: none
   :emphasize-lines: 1,1

    mlr --opprint put '$o = system("echo {" . NR . "}")' data/small
    a   b   i x                   y                   o
    pan pan 1 0.3467901443380824  0.7268028627434533  {1}
    eks pan 2 0.7586799647899636  0.5221511083334797  {2}
    wye wye 3 0.20460330576630303 0.33831852551664776 {3}
    eks wye 4 0.38139939387114097 0.13418874328430463 {4}
    wye pan 5 0.5732889198020006  0.8636244699032729  {5}

.. code-block:: none
   :emphasize-lines: 1,1

    mlr --opprint put '$o = system("echo -n ".$a."| sha1sum")' data/small
    a   b   i x                   y                   o
    pan pan 1 0.3467901443380824  0.7268028627434533  f29c748220331c273ef16d5115f6ecd799947f13  -
    eks pan 2 0.7586799647899636  0.5221511083334797  456d988ecb3bf1b75f057fc6e9fe70db464e9388  -
    wye wye 3 0.20460330576630303 0.33831852551664776 eab0de043d67f441c7fd1e335f0ca38708e6ebf7  -
    eks wye 4 0.38139939387114097 0.13418874328430463 456d988ecb3bf1b75f057fc6e9fe70db464e9388  -
    wye pan 5 0.5732889198020006  0.8636244699032729  eab0de043d67f441c7fd1e335f0ca38708e6ebf7  -

Note that running a subprocess on every record takes a non-trivial amount of time. Comparing asking the system ``date`` command for the current time in nanoseconds versus computing it in process:

..
    hard-coded, not live-code, since %N doesn't exist on all platforms

.. code-block:: none

    $ mlr --opprint put '$t=system("date +%s.%N")' then step -a delta -f t data/small
    a   b   i x                   y                   t                    t_delta
    pan pan 1 0.3467901443380824  0.7268028627434533  1568774318.513903817 0
    eks pan 2 0.7586799647899636  0.5221511083334797  1568774318.514722876 0.000819
    wye wye 3 0.20460330576630303 0.33831852551664776 1568774318.515618046 0.000895
    eks wye 4 0.38139939387114097 0.13418874328430463 1568774318.516547441 0.000929
    wye pan 5 0.5732889198020006  0.8636244699032729  1568774318.517518828 0.000971

.. code-block:: none

    $ mlr --opprint put '$t=systime()' then step -a delta -f t data/small
    a   b   i x                   y                   t                 t_delta
    pan pan 1 0.3467901443380824  0.7268028627434533  1568774318.518699 0
    eks pan 2 0.7586799647899636  0.5221511083334797  1568774318.518717 0.000018
    wye wye 3 0.20460330576630303 0.33831852551664776 1568774318.518723 0.000006
    eks wye 4 0.38139939387114097 0.13418874328430463 1568774318.518727 0.000004
    wye pan 5 0.5732889198020006  0.8636244699032729  1568774318.518730 0.000003
