..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

FAQ
=========

No output at all
----------------------------------------------------------------

Try ``od -xcv`` and/or ``cat -e`` on your file to check for non-printable characters.

If you're using Miller version less than 5.0.0 (try ``mlr --version`` on your system to find out), when the line-ending-autodetect feature was introduced, please see http://johnkerl.org/miller-releases/miller-4.5.0/doc/index.html.

Fields not selected
----------------------------------------------------------------

Check the field-separators of the data, e.g. with the command-line ``head`` program. Example: for CSV, Miller's default record separator is comma; if your data is tab-delimited, e.g. ``aTABbTABc``, then Miller won't find three fields named ``a``, ``b``, and ``c`` but rather just one named ``aTABbTABc``.  Solution in this case: ``mlr --fs tab {remaining arguments ...}``.

Also try ``od -xcv`` and/or ``cat -e`` on your file to check for non-printable characters.

Diagnosing delimiter specifications
----------------------------------------------------------------

.. code-block:: none

    # Use the `file` command to see if there are CR/LF terminators (in this case,
    # there are not):
    $ file data/colours.csv 
    data/colours.csv: UTF-8 Unicode text
    
    # Look at the file to find names of fields
    $ cat data/colours.csv 
    KEY;DE;EN;ES;FI;FR;IT;NL;PL;RO;TR
    masterdata_colourcode_1;Weiß;White;Blanco;Valkoinen;Blanc;Bianco;Wit;Biały;Alb;Beyaz
    masterdata_colourcode_2;Schwarz;Black;Negro;Musta;Noir;Nero;Zwart;Czarny;Negru;Siyah
    
    # Extract a few fields:
    $ mlr --csv cut -f KEY,PL,RO data/colours.csv 
    (only blank lines appear)
    
    # Use XTAB output format to get a sharper picture of where records/fields
    # are being split:
    $ mlr --icsv --oxtab cat data/colours.csv 
    KEY;DE;EN;ES;FI;FR;IT;NL;PL;RO;TR masterdata_colourcode_1;Weiß;White;Blanco;Valkoinen;Blanc;Bianco;Wit;Biały;Alb;Beyaz
    
    KEY;DE;EN;ES;FI;FR;IT;NL;PL;RO;TR masterdata_colourcode_2;Schwarz;Black;Negro;Musta;Noir;Nero;Zwart;Czarny;Negru;Siyah
    
    # Using XTAB output format makes it clearer that KEY;DE;...;RO;TR is being
    # treated as a single field name in the CSV header, and likewise each
    # subsequent line is being treated as a single field value. This is because
    # the default field separator is a comma but we have semicolons here.
    # Use XTAB again with different field separator (--fs semicolon):
     mlr --icsv --ifs semicolon --oxtab cat data/colours.csv 
    KEY masterdata_colourcode_1
    DE  Weiß
    EN  White
    ES  Blanco
    FI  Valkoinen
    FR  Blanc
    IT  Bianco
    NL  Wit
    PL  Biały
    RO  Alb
    TR  Beyaz
    
    KEY masterdata_colourcode_2
    DE  Schwarz
    EN  Black
    ES  Negro
    FI  Musta
    FR  Noir
    IT  Nero
    NL  Zwart
    PL  Czarny
    RO  Negru
    TR  Siyah
    
    # Using the new field-separator, retry the cut:
     mlr --csv --fs semicolon cut -f KEY,PL,RO data/colours.csv 
    KEY;PL;RO
    masterdata_colourcode_1;Biały;Alb
    masterdata_colourcode_2;Czarny;Negru

How do I suppress numeric conversion?
----------------------------------------------------------------

**TL;DR use put -S**.

Within ``mlr put`` and ``mlr filter``, the default behavior for scanning input records is to parse them as integer, if possible, then as float, if possible, else leave them as string:

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat data/scan-example-1.tbl
    value
    1
    2.0
    3x
    hello

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --pprint put '$copy = $value; $type = typeof($value)' data/scan-example-1.tbl
    value copy  type
    1     1     int
    2.0   2.0   float
    3x    3x    string
    hello hello string

The numeric-conversion rule is simple:

* Try to scan as integer (``"1"`` should be int);
* If that doesn't succeed, try to scan as float (``"1.0"`` should be float);
* If that doesn't succeed, leave the value as a string (``"1x"`` is string).

This is a sensible default: you should be able to put ``'$z = $x + $y'`` without having to write ``'$z = int($x) + float($y)'``.  Also note that default output format for floating-point numbers created by ``put`` (and other verbs such as ``stats1``) is six decimal places; you can override this using ``mlr --ofmt``.  Also note that Miller uses your system's Go library functions whenever possible: e.g. ``sscanf`` for converting strings to integer or floating-point.

But now suppose you have data like these:

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat data/scan-example-2.tbl
    value
    0001
    0002
    0005
    0005WA
    0006
    0007
    0007WA
    0008
    0009
    0010

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --pprint put '$copy = $value; $type = typeof($value)' data/scan-example-2.tbl
    value  copy   type
    0001   0001   int
    0002   0002   int
    0005   0005   int
    0005WA 0005WA string
    0006   0006   int
    0007   0007   int
    0007WA 0007WA string
    0008   0008   float
    0009   0009   float
    0010   0010   int

The same conversion rules as above are being used. Namely:

* By default field values are inferred to int, else float, else string;

* leading zeroes indicate octal for integers (``sscanf`` semantics);

* since ``0008`` doesn't scan as integer (leading 0 requests octal but 8 isn't a valid octal digit), the float scan is tried next and it succeeds;

* default floating-point output format is 6 decimal places (override with ``mlr --ofmt``).

Taken individually the rules make sense; taken collectively they produce a mishmash of types here.

The solution is to **use the -S flag** for ``mlr put`` and/or ``mlr filter``. Then all field values are left as string. You can type-coerce on demand using syntax like ``'$z = int($x) + float($y)'``. (See also :doc:`reference-dsl`; see also https://github.com/johnkerl/miller/issues/150.)

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --pprint put -S '$copy = $value; $type = typeof($value)' data/scan-example-2.tbl
    value  copy   type
    0001   0001   int
    0002   0002   int
    0005   0005   int
    0005WA 0005WA string
    0006   0006   int
    0007   0007   int
    0007WA 0007WA string
    0008   0008   float
    0009   0009   float
    0010   0010   int

How do I examine then-chaining?
----------------------------------------------------------------

Then-chaining found in Miller is intended to function the same as Unix pipes, but with less keystroking. You can print your data one pipeline step at a time, to see what intermediate output at one step becomes the input to the next step.

First, look at the input data:

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat data/then-example.csv
    Status,Payment_Type,Amount
    paid,cash,10.00
    pending,debit,20.00
    paid,cash,50.00
    pending,credit,40.00
    paid,debit,30.00

Next, run the first step of your command, omitting anything from the first ``then`` onward:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --opprint count-distinct -f Status,Payment_Type data/then-example.csv
    Status  Payment_Type count
    paid    cash         2
    pending debit        1
    pending credit       1
    paid    debit        1

After that, run it with the next ``then`` step included:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --opprint count-distinct -f Status,Payment_Type then sort -nr count data/then-example.csv
    Status  Payment_Type count
    paid    cash         2
    pending debit        1
    pending credit       1
    paid    debit        1

Now if you use ``then`` to include another verb after that, the columns ``Status``, ``Payment_Type``, and ``count`` will be the input to that verb.

Note, by the way, that you'll get the same results using pipes:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --csv count-distinct -f Status,Payment_Type data/then-example.csv | mlr --icsv --opprint sort -nr count
    Status  Payment_Type count
    paid    cash         2
    pending debit        1
    pending credit       1
    paid    debit        1

I assigned $9 and it's not 9th
----------------------------------------------------------------

Miller records are ordered lists of key-value pairs. For NIDX format, DKVP format when keys are missing, or CSV/CSV-lite format with ``--implicit-csv-header``, Miller will sequentially assign keys of the form ``1``, ``2``, etc. But these are not integer array indices: they're just field names taken from the initial field ordering in the input data.

.. code-block:: none
   :emphasize-lines: 1,1

    $ echo x,y,z | mlr --dkvp cat
    1=x,2=y,3=z

.. code-block:: none
   :emphasize-lines: 1,1

    $ echo x,y,z | mlr --dkvp put '$6="a";$4="b";$55="cde"'
    1=x,2=y,3=z,6=a,4=b,55=cde

.. code-block:: none
   :emphasize-lines: 1,1

    $ echo x,y,z | mlr --nidx cat
    x,y,z

.. code-block:: none
   :emphasize-lines: 1,1

    $ echo x,y,z | mlr --csv --implicit-csv-header cat
    1,2,3
    x,y,z

.. code-block:: none
   :emphasize-lines: 1,1

    $ echo x,y,z | mlr --dkvp rename 2,999
    1=x,999=y,3=z

.. code-block:: none
   :emphasize-lines: 1,1

    $ echo x,y,z | mlr --dkvp rename 2,newname
    1=x,newname=y,3=z

.. code-block:: none
   :emphasize-lines: 1,1

    $ echo x,y,z | mlr --csv --implicit-csv-header reorder -f 3,1,2
    3,1,2
    z,x,y

How can I filter by date?
----------------------------------------------------------------

Given input like

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat dates.csv
    date,event
    2018-02-03,initialization
    2018-03-07,discovery
    2018-02-03,allocation

we can use ``strptime`` to parse the date field into seconds-since-epoch and then do numeric comparisons.  Simply match your input dataset's date-formatting to the :ref:`reference-dsl-strptime` format-string.  For example:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --csv filter 'strptime($date, "%Y-%m-%d") > strptime("2018-03-03", "%Y-%m-%d")' dates.csv
    date,event
    2018-03-07,discovery

Caveat: localtime-handling in timezones with DST is still a work in progress; see https://github.com/johnkerl/miller/issues/170. See also https://github.com/johnkerl/miller/issues/208 -- thanks @aborruso!

How can I handle commas-as-data in various formats?
----------------------------------------------------------------

:doc:`CSV <file-formats>` handles this well and by design:

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat commas.csv
    Name,Role
    "Xiao, Lin",administrator
    "Khavari, Darius",tester

Likewise :ref:`file-formats-json`:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --ojson cat commas.csv
    {
      "Name": "Xiao, Lin",
      "Role": "administrator"
    }
    {
      "Name": "Khavari, Darius",
      "Role": "tester"
    }

For Miller's :ref:`vertical-tabular format <file-formats-xtab>` there is no escaping for carriage returns, but commas work fine:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --oxtab cat commas.csv
    Name Xiao, Lin
    Role administrator
    
    Name Khavari, Darius
    Role tester

But for :ref:`Key-value_pairs <file-formats-dkvp>` and :ref:`index-numbered <file-formats-nidx>`, commas are the default field separator. And -- as of Miller 5.4.0 anyway -- there is no CSV-style double-quote-handling like there is for CSV. So commas within the data look like delimiters:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --odkvp cat commas.csv
    Name=Xiao, Lin,Role=administrator
    Name=Khavari, Darius,Role=tester

One solution is to use a different delimiter, such as a pipe character:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --odkvp --ofs pipe cat commas.csv
    Name=Xiao, Lin|Role=administrator
    Name=Khavari, Darius|Role=tester

To be extra-sure to avoid data/delimiter clashes, you can also use control
characters as delimiters -- here, control-A:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --odkvp --ofs '\001'  cat commas.csv | cat -v
    Name=Xiao, Lin\001Role=administrator
    Name=Khavari, Darius\001Role=tester

How can I handle field names with special symbols in them?
----------------------------------------------------------------

Simply surround the field names with curly braces:

.. code-block:: none
   :emphasize-lines: 1,1

    $ echo 'x.a=3,y:b=4,z/c=5' | mlr put '${product.all} = ${x.a} * ${y:b} * ${z/c}'
    x.a=3,y:b=4,z/c=5,product.all=60

How to escape '?' in regexes?
----------------------------------------------------------------

One way is to use square brackets; an alternative is to use simple string-substitution rather than a regular expression.

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat data/question.dat
    a=is it?,b=it is!
.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --oxtab put '$c = gsub($a, "[?]"," ...")' data/question.dat
    a is it?
    b it is!
    c is it ...
.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --oxtab put '$c = ssub($a, "?"," ...")' data/question.dat
    a is it?
    b it is!
    c is it ...

The ``ssub`` function exists precisely for this reason: so you don't have to escape anything.

How can I put single-quotes into strings?
----------------------------------------------------------------

This is a little tricky due to the shell's handling of quotes. For simplicity, let's first put an update script into a file:

.. code-block:: none

    $a = "It's OK, I said, then 'for now'."

.. code-block:: none
   :emphasize-lines: 1,1

    $ echo a=bcd | mlr put -f data/single-quote-example.mlr
    a=It's OK, I said, then 'for now'.

So, it's simple: Miller's DSL uses double quotes for strings, and you can put single quotes (or backslash-escaped double-quotes) inside strings, no problem.

Without putting the update expression in a file, it's messier:

.. code-block:: none
   :emphasize-lines: 1,1

    $ echo a=bcd | mlr put '$a="It'\''s OK, I said, '\''for now'\''."'
    a=It's OK, I said, 'for now'.

The idea is that the outermost single-quotes are to protect the ``put`` expression from the shell, and the double quotes within them are for Miller. To get a single quote in the middle there, you need to actually put it *outside* the single-quoting for the shell. The pieces are the following, all concatenated together:

* ``$a="It``
* ``\'``
* ``s OK, I said,``
* ``\'``
* ``for now``
* ``\'``
* ``.``

Why doesn't mlr cut put fields in the order I want?
----------------------------------------------------------------

Example: columns ``x,i,a`` were requested but they appear here in the order ``a,i,x``:

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr cut -f x,i,a data/small
    a=pan,i=1,x=0.3467901443380824
    a=eks,i=2,x=0.7586799647899636
    a=wye,i=3,x=0.20460330576630303
    a=eks,i=4,x=0.38139939387114097
    a=wye,i=5,x=0.5732889198020006

The issue is that Miller's ``cut``, by default, outputs cut fields in the order they appear in the input data. This design decision was made intentionally to parallel the Unix/Linux system ``cut`` command, which has the same semantics.

The solution is to use the ``-o`` option:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr cut -o -f x,i,a data/small
    x=0.3467901443380824,i=1,a=pan
    x=0.7586799647899636,i=2,a=eks
    x=0.20460330576630303,i=3,a=wye
    x=0.38139939387114097,i=4,a=eks
    x=0.5732889198020006,i=5,a=wye

NR is not consecutive after then-chaining
----------------------------------------------------------------

Given this input data:

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729

why don't I see ``NR=1`` and ``NR=2`` here??

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr filter '$x > 0.5' then put '$NR = NR' data/small
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797,NR=2
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,NR=5

The reason is that ``NR`` is computed for the original input records and isn't dynamically updated. By contrast, ``NF`` is dynamically updated: it's the number of fields in the current record, and if you add/remove a field, the value of ``NF`` will change:

.. code-block:: none
   :emphasize-lines: 1,1

    $ echo x=1,y=2,z=3 | mlr put '$nf1 = NF; $u = 4; $nf2 = NF; unset $x,$y,$z; $nf3 = NF'
    nf1=3,u=4,nf2=5,nf3=3

``NR``, by contrast (and ``FNR`` as well), retains the value from the original input stream, and records may be dropped by a ``filter`` within a ``then``-chain. To recover consecutive record numbers, you can use out-of-stream variables as follows:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint --from data/small put '
      begin{ @nr1 = 0 }
      @nr1 += 1;
      $nr1 = @nr1
    ' \
    then filter '$x>0.5' \
    then put '
      begin{ @nr2 = 0 }
      @nr2 += 1;
      $nr2 = @nr2
    '
    a   b   i x                  y                  nr1 nr2
    eks pan 2 0.7586799647899636 0.5221511083334797 2   1
    wye pan 5 0.5732889198020006 0.8636244699032729 5   2

Or, simply use ``mlr cat -n``:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr filter '$x > 0.5' then cat -n data/small
    n=1,a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    n=2,a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729

Why am I not seeing all possible joins occur?
----------------------------------------------------------------

**This section describes behavior before Miller 5.1.0. As of 5.1.0, -u is the default.**

For example, the right file here has nine records, and the left file should add in the ``hostname`` column -- so the join output should also have 9 records:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsvlite --opprint cat data/join-u-left.csv
    hostname              ipaddr
    nadir.east.our.org    10.3.1.18
    zenith.west.our.org   10.3.1.27
    apoapsis.east.our.org 10.4.5.94

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsvlite --opprint cat data/join-u-right.csv
    ipaddr    timestamp  bytes
    10.3.1.27 1448762579 4568
    10.3.1.18 1448762578 8729
    10.4.5.94 1448762579 17445
    10.3.1.27 1448762589 12
    10.3.1.18 1448762588 44558
    10.4.5.94 1448762589 8899
    10.3.1.27 1448762599 0
    10.3.1.18 1448762598 73425
    10.4.5.94 1448762599 12200

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsvlite --opprint join -s -j ipaddr -f data/join-u-left.csv data/join-u-right.csv
    ipaddr    hostname              timestamp  bytes
    10.3.1.27 zenith.west.our.org   1448762579 4568
    10.4.5.94 apoapsis.east.our.org 1448762579 17445
    10.4.5.94 apoapsis.east.our.org 1448762589 8899
    10.4.5.94 apoapsis.east.our.org 1448762599 12200

The issue is that Miller's ``join``, by default (before 5.1.0), took input sorted (lexically ascending) by the sort keys on both the left and right files.  This design decision was made intentionally to parallel the Unix/Linux system ``join`` command, which has the same semantics. The benefit of this default is that the joiner program can stream through the left and right files, needing to load neither entirely into memory. The drawback, of course, is that is requires sorted input.

The solution (besides pre-sorting the input files on the join keys) is to simply use **mlr join -u** (which is now the default). This loads the left file entirely into memory (while the right file is still streamed one line at a time) and does all possible joins without requiring sorted input:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsvlite --opprint join -u -j ipaddr -f data/join-u-left.csv data/join-u-right.csv
    ipaddr    hostname              timestamp  bytes
    10.3.1.27 zenith.west.our.org   1448762579 4568
    10.3.1.18 nadir.east.our.org    1448762578 8729
    10.4.5.94 apoapsis.east.our.org 1448762579 17445
    10.3.1.27 zenith.west.our.org   1448762589 12
    10.3.1.18 nadir.east.our.org    1448762588 44558
    10.4.5.94 apoapsis.east.our.org 1448762589 8899
    10.3.1.27 zenith.west.our.org   1448762599 0
    10.3.1.18 nadir.east.our.org    1448762598 73425
    10.4.5.94 apoapsis.east.our.org 1448762599 12200

General advice is to make sure the left-file is relatively small, e.g. containing name-to-number mappings, while saving large amounts of data for the right file.

How to rectangularize after joins with unpaired?
----------------------------------------------------------------

Suppose you have the following two data files:

.. code-block:: none

    id,code
    3,0000ff
    2,00ff00
    4,ff0000

.. code-block:: none

    id,color
    4,red
    2,green

Joining on color the results are as expected:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --csv join -j id -f data/color-codes.csv data/color-names.csv
    id,code,color
    4,ff0000,red
    2,00ff00,green

However, if we ask for left-unpaireds, since there's no ``color`` column, we get a row not having the same column names as the other:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --csv join --ul -j id -f data/color-codes.csv data/color-names.csv
    id,code,color
    4,ff0000,red
    2,00ff00,green
    
    id,code
    3,0000ff

To fix this, we can use **unsparsify**:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --csv join --ul -j id -f data/color-codes.csv then unsparsify --fill-with "" data/color-names.csv
    id,code,color
    4,ff0000,red
    2,00ff00,green
    3,0000ff,

Thanks to @aborruso for the tip!

What about XML or JSON file formats?
----------------------------------------------------------------

Miller handles **tabular data**, which is a list of records each having fields which are key-value pairs. Miller also doesn't require that each record have the same field names (see also :doc:`record-heterogeneity`). Regardless, tabular data is a **non-recursive data structure**.

XML, JSON, etc. are, by contrast, all **recursive** or **nested** data structures. For example, in JSON you can represent a hash map whose values are lists of lists.

Now, you can put tabular data into these formats -- since list-of-key-value-pairs is one of the things representable in XML or JSON. Example:

.. code-block:: none

    # DKVP
    x=1,y=2
    z=3

.. code-block:: none

    # XML
    <table>
      <record>
        <field>
          <key> x </key> <value> 1 </value>
        </field>
        <field>
          <key> y </key> <value> 2 </value>
        </field>
      </record>
      <record>
        <field>
          <key> z </key> <value> 3 </value>
        </field>
      </record>
    </table>

.. code-block:: none

    # JSON
    [{"x":1,"y":2},{"z":3}]

However, a tool like Miller which handles non-recursive data is never going to be able to handle full XML/JSON semantics -- only a small subset.  If tabular data represented in XML/JSON/etc are sufficiently well-structured, it may be easy to grep/sed out the data into a simpler text form -- this is a general text-processing problem.

Miller does support tabular data represented in JSON: please see :doc:`file-formats`.  See also `jq <https://stedolan.github.io/jq/>`_ for a truly powerful, JSON-specific tool.

For XML, my suggestion is to use a tool like `ff-extractor <http://ff-extractor.sourceforge.net>`_ to do format conversion.
