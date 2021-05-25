..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Main reference
================================================================


.. _reference-command-overview:

Command overview
----------------------------------------------------------------

Whereas the Unix toolkit is made of the separate executables ``cat``, ``tail``, ``cut``,
``sort``, etc., Miller has subcommands, or **verbs**, invoked as follows:

::

    mlr tac *.dat
    mlr cut --complement -f os_version *.dat
    mlr sort -f hostname,uptime *.dat

These fall into categories as follows:

* Analogs of their Unix-toolkit namesakes, discussed below as well as in :doc:`feature-comparison`: :ref:`reference-verbs-cat` :ref:`reference-verbs-cut` :ref:`reference-verbs-grep` :ref:`reference-verbs-head` :ref:`reference-verbs-join` :ref:`reference-verbs-sort` :ref:`reference-verbs-tac` :ref:`reference-verbs-tail` :ref:`reference-verbs-top` :ref:`reference-verbs-uniq`.

* ``awk``-like functionality: :ref:`reference-verbs-filter` :ref:`reference-verbs-put` :ref:`reference-verbs-sec2gmt` :ref:`reference-verbs-sec2gmtdate` :ref:`reference-verbs-step` :ref:`reference-verbs-tee`.

* Statistically oriented: :ref:`reference-verbs-bar` :ref:`reference-verbs-bootstrap` :ref:`reference-verbs-decimate` :ref:`reference-verbs-histogram` :ref:`reference-verbs-least-frequent` :ref:`reference-verbs-most-frequent` :ref:`reference-verbs-sample` :ref:`reference-verbs-shuffle` :ref:`reference-verbs-stats1` :ref:`reference-verbs-stats2`.

* Particularly oriented toward :doc:`record-heterogeneity`, although all Miller commands can handle heterogeneous records: :ref:`reference-verbs-group-by` :ref:`reference-verbs-group-like` :ref:`reference-verbs-having-fields`.

* These draw from other sources (see also :doc:`originality`): :ref:`reference-verbs-count-distinct` is SQL-ish, and :ref:`reference-verbs-rename` can be done by ``sed`` (which does it faster: see :doc:`performance`. Verbs: :ref:`reference-verbs-check` :ref:`reference-verbs-count-distinct` :ref:`reference-verbs-label` :ref:`reference-verbs-merge-fields` :ref:`reference-verbs-nest` :ref:`reference-verbs-nothing` :ref:`reference-verbs-regularize` :ref:`reference-verbs-rename` :ref:`reference-verbs-reorder` :ref:`reference-verbs-reshape` :ref:`reference-verbs-seqgen`.

I/O options
----------------------------------------------------------------

Formats
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Options:

::

    --dkvp    --idkvp    --odkvp
    --nidx    --inidx    --onidx
    --csv     --icsv     --ocsv
    --csvlite --icsvlite --ocsvlite
    --pprint  --ipprint  --opprint  --right
    --xtab    --ixtab    --oxtab
    --json    --ijson    --ojson

These are as discussed in :doc:`file-formats`, with the exception of ``--right`` which makes pretty-printed output right-aligned:

::

    $ mlr --opprint cat data/small
    a   b   i x                   y
    pan pan 1 0.3467901443380824  0.7268028627434533
    eks pan 2 0.7586799647899636  0.5221511083334797
    wye wye 3 0.20460330576630303 0.33831852551664776
    eks wye 4 0.38139939387114097 0.13418874328430463
    wye pan 5 0.5732889198020006  0.8636244699032729

::

    $ mlr --opprint --right cat data/small
      a   b i                   x                   y
    pan pan 1  0.3467901443380824  0.7268028627434533
    eks pan 2  0.7586799647899636  0.5221511083334797
    wye wye 3 0.20460330576630303 0.33831852551664776
    eks wye 4 0.38139939387114097 0.13418874328430463
    wye pan 5  0.5732889198020006  0.8636244699032729

Additional notes:

* Use ``--csv``, ``--pprint``, etc. when the input and output formats are the same.

* Use ``--icsv --opprint``, etc. when you want format conversion as part of what Miller does to your data.

* DKVP (key-value-pair) format is the default for input and output. So, ``--oxtab`` is the same as ``--idkvp --oxtab``.

**Pro-tip:** Please use either **--format1**, or **--iformat1 --oformat2**.  If you use **--format1 --oformat2** then what happens is that flags are set up for input *and* output for format1, some of which are overwritten for output in format2. For technical reasons, having ``--oformat2`` clobber all the output-related effects of ``--format1`` also removes some flexibility from the command-line interface. See also https://github.com/johnkerl/miller/issues/180 and https://github.com/johnkerl/miller/issues/199.

In-place mode
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Use the ``mlr -I`` flag to process files in-place. For example, ``mlr -I --csv cut -x -f unwanted_column_name mydata/*.csv`` will remove ``unwanted_column_name`` from all your ``*.csv`` files in your ``mydata/`` subdirectory.

By default, Miller output goes to the screen (or you can redirect a file using ``>`` or to another process using ``|``). With ``-I``, for each file name on the command line, output is written to a temporary file in the same directory. Miller writes its output into that temp file, which is then renamed over the original.  Then, processing continues on the next file. Each file is processed in isolation: if the output format is CSV, CSV headers will be present in each output file; statistics are only over each file's own records; and so on.

Please see :ref:`10min-choices-for-printing-to-files` for examples.

Compression
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Options:

::

    --prepipe {command}


The prepipe command is anything which reads from standard input and produces data acceptable to Miller. Nominally this allows you to use whichever decompression utilities you have installed on your system, on a per-file basis. If the command has flags, quote them: e.g. ``mlr --prepipe 'zcat -cf'``. Examples:

::

    # These two produce the same output:
    $ gunzip < myfile1.csv.gz | mlr cut -f hostname,uptime
    $ mlr --prepipe gunzip cut -f hostname,uptime myfile1.csv.gz
    # With multiple input files you need --prepipe:
    $ mlr --prepipe gunzip cut -f hostname,uptime myfile1.csv.gz myfile2.csv.gz
    $ mlr --prepipe gunzip --idkvp --oxtab cut -f hostname,uptime myfile1.dat.gz myfile2.dat.gz

::

    # Similar to the above, but with compressed output as well as input:
    $ gunzip < myfile1.csv.gz | mlr cut -f hostname,uptime | gzip > outfile.csv.gz
    $ mlr --prepipe gunzip cut -f hostname,uptime myfile1.csv.gz | gzip > outfile.csv.gz
    $ mlr --prepipe gunzip cut -f hostname,uptime myfile1.csv.gz myfile2.csv.gz | gzip > outfile.csv.gz

::

    # Similar to the above, but with different compression tools for input and output:
    $ gunzip < myfile1.csv.gz | mlr cut -f hostname,uptime | xz -z > outfile.csv.xz
    $ xz -cd < myfile1.csv.xz | mlr cut -f hostname,uptime | gzip > outfile.csv.xz
    $ mlr --prepipe 'xz -cd' cut -f hostname,uptime myfile1.csv.xz myfile2.csv.xz | xz -z > outfile.csv.xz

.. _reference-separators:

Record/field/pair separators
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Miller has record separators ``IRS`` and ``ORS``, field separators ``IFS`` and ``OFS``, and pair separators ``IPS`` and ``OPS``.  For example, in the DKVP line ``a=1,b=2,c=3``, the record separator is newline, field separator is comma, and pair separator is the equals sign. These are the default values.

Options:

::

    --rs --irs --ors
    --fs --ifs --ofs --repifs
    --ps --ips --ops

* You can change a separator from input to output via e.g. ``--ifs = --ofs :``. Or, you can specify that the same separator is to be used for input and output via e.g. ``--fs :``.

* The pair separator is only relevant to DKVP format.

* Pretty-print and xtab formats ignore the separator arguments altogether.

* The ``--repifs`` means that multiple successive occurrences of the field separator count as one.  For example, in CSV data we often signify nulls by empty strings, e.g. ``2,9,,,,,6,5,4``. On the other hand, if the field separator is a space, it might be more natural to parse ``2 4    5`` the same as ``2 4 5``: ``--repifs --ifs ' '`` lets this happen.  In fact, the ``--ipprint`` option above is internally implemented in terms of ``--repifs``.

* Just write out the desired separator, e.g. ``--ofs '|'``. But you may use the symbolic names ``newline``, ``space``, ``tab``, ``pipe``, or ``semicolon`` if you like.

Number formatting
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

The command-line option ``--ofmt {format string}`` is the global number format for commands which generate numeric output, e.g. ``stats1``, ``stats2``, ``histogram``, and ``step``, as well as ``mlr put``. Examples:

::

    --ofmt %.9le  --ofmt %.6lf  --ofmt %.0lf

These are just familiar ``printf`` formats applied to double-precision numbers.  Please don't use ``%s`` or ``%d``. Additionally, if you use leading width (e.g. ``%18.12lf``) then the output will contain embedded whitespace, which may not be what you want if you pipe the output to something else, particularly CSV. I use Miller's pretty-print format (``mlr --opprint``) to column-align numerical data.

To apply formatting to a single field, overriding the global ``ofmt``, use ``fmtnum`` function within ``mlr put``. For example:

::

    $ echo 'x=3.1,y=4.3' | mlr put '$z=fmtnum($x*$y,"%08lf")'
    x=3.1,y=4.3,z=13.330000

::

    $ echo 'x=0xffff,y=0xff' | mlr put '$z=fmtnum(int($x*$y),"%08llx")'
    x=0xffff,y=0xff,z=00feff01

Input conversion from hexadecimal is done automatically on fields handled by ``mlr put`` and ``mlr filter`` as long as the field value begins with "0x".  To apply output conversion to hexadecimal on a single column, you may use ``fmtnum``, or the keystroke-saving ``hexfmt`` function. Example:

::

    $ echo 'x=0xffff,y=0xff' | mlr put '$z=hexfmt($x*$y)'
    x=0xffff,y=0xff,z=0xfeff01

Data transformations (verbs)
----------------------------------------------------------------

Please see the separate page :doc:`reference-verbs`.

Expression language for filter and put
----------------------------------------------------------------

Please see the separate page :doc:`reference-dsl`.

then-chaining
----------------------------------------------------------------

In accord with the `Unix philosophy <http://en.wikipedia.org/wiki/Unix_philosophy>`_, you can pipe data into or out of Miller. For example:

::

    mlr cut --complement -f os_version *.dat | mlr sort -f hostname,uptime

You can, if you like, instead simply chain commands together using the ``then`` keyword:

::

    mlr cut --complement -f os_version then sort -f hostname,uptime *.dat

(You can precede the very first verb with ``then``, if you like, for symmetry.)

Here's a performance comparison:

::

    % cat piped.sh
    mlr cut -x -f i,y data/big | mlr sort -n y > /dev/null
    
    % time sh piped.sh
    real 0m2.828s
    user 0m3.183s
    sys  0m0.137s
    
    
    % cat chained.sh
    mlr cut -x -f i,y then sort -n y data/big > /dev/null
    
    % time sh chained.sh
    real 0m2.082s
    user 0m1.933s
    sys  0m0.137s

There are two reasons to use then-chaining: one is for performance, although I don't expect this to be a win in all cases.  Using then-chaining avoids redundant string-parsing and string-formatting at each pipeline step: instead input records are parsed once, they are fed through each pipeline stage in memory, and then output records are formatted once. On the other hand, Miller is single-threaded, while modern systems are usually multi-processor, and when streaming-data programs operate through pipes, each one can use a CPU.  Rest assured you get the same results either way.

The other reason to use then-chaining is for simplicity: you don't have re-type formatting flags (e.g. ``--csv --fs tab``) at every pipeline stage.

Auxiliary commands
----------------------------------------------------------------

There are a few nearly-standalone programs which have nothing to do with the rest of Miller, do not participate in record streams, and do not deal with file formats. They might as well be little standalone executables but they're delivered within the main Miller executable for convenience.

::

    $ mlr aux-list
    Available subcommands:
      aux-list
      lecat
      termcvt
      hex
      unhex
      netbsd-strptime
    For more information, please invoke mlr {subcommand} --help

::

    $ mlr lecat --help
    Usage: mlr lecat [options] {zero or more file names}
    Simply echoes input, but flags CR characters in red and LF characters in green.
    If zero file names are supplied, standard input is read.
    Options:
    --mono: don't try to colorize the output
    -h or --help: print this message

::

    $ mlr termcvt --help
    Usage: mlr termcvt [option] {zero or more file names}
    Option (exactly one is required):
    --cr2crlf
    --lf2crlf
    --crlf2cr
    --crlf2lf
    --cr2lf
    --lf2cr
    -I in-place processing (default is to write to stdout)
    -h or --help: print this message
    Zero file names means read from standard input.
    Output is always to standard output; files are not written in-place.

::

    $ mlr hex --help
    Usage: mlr hex [options] {zero or more file names}
    Simple hex-dump.
    If zero file names are supplied, standard input is read.
    Options:
    -r: print only raw hex without leading offset indicators or trailing ASCII dump.
    -h or --help: print this message

::

    $ mlr unhex --help
    Usage: mlr unhex [option] {zero or more file names}
    Options:
    -h or --help: print this message
    Zero file names means read from standard input.
    Output is always to standard output; files are not written in-place.

Examples:

::

    $ echo 'Hello, world!' | mlr lecat --mono
    Hello, world![LF]

::

    $ echo 'Hello, world!' | mlr termcvt --lf2crlf | mlr lecat --mono
    Hello, world![CR][LF]

::

    $ mlr hex data/budget.csv
    00000000: 23 20 41 73  61 6e 61 20  2d 2d 20 68  65 72 65 20 |# Asana -- here |
    00000010: 61 72 65 20  74 68 65 20  62 75 64 67  65 74 20 66 |are the budget f|
    00000020: 69 67 75 72  65 73 20 79  6f 75 20 61  73 6b 65 64 |igures you asked|
    00000030: 20 66 6f 72  21 0a 74 79  70 65 2c 71  75 61 6e 74 | for!.type,quant|
    00000040: 69 74 79 0a  70 75 72 70  6c 65 2c 34  35 36 2e 37 |ity.purple,456.7|
    00000050: 38 0a 67 72  65 65 6e 2c  36 37 38 2e  31 32 0a 6f |8.green,678.12.o|
    00000060: 72 61 6e 67  65 2c 31 32  33 2e 34 35  0a          |range,123.45.|

::

    $ mlr hex -r data/budget.csv
    23 20 41 73  61 6e 61 20  2d 2d 20 68  65 72 65 20 
    61 72 65 20  74 68 65 20  62 75 64 67  65 74 20 66 
    69 67 75 72  65 73 20 79  6f 75 20 61  73 6b 65 64 
    20 66 6f 72  21 0a 74 79  70 65 2c 71  75 61 6e 74 
    69 74 79 0a  70 75 72 70  6c 65 2c 34  35 36 2e 37 
    38 0a 67 72  65 65 6e 2c  36 37 38 2e  31 32 0a 6f 
    72 61 6e 67  65 2c 31 32  33 2e 34 35  0a          

::

    $ mlr hex -r data/budget.csv | sed 's/20/2a/g' | mlr unhex
    #*Asana*--*here*are*the*budget*figures*you*asked*for!
    type,quantity
    purple,456.78
    green,678.12
    orange,123.45

Data types
----------------------------------------------------------------

Miller's input and output are all string-oriented: there is (as of August 2015 anyway) no support for binary record packing. In this sense, everything is a string in and out of Miller.  During processing, field names are always strings, even if they have names like "3"; field values are usually strings.  Field values' ability to be interpreted as a non-string type only has meaning when comparison or function operations are done on them.  And it is an error condition if Miller encounters non-numeric (or otherwise mistyped) data in a field in which it has been asked to do numeric (or otherwise type-specific) operations.

Field values are treated as numeric for the following:

* Numeric sort: ``mlr sort -n``, ``mlr sort -nr``.
* Statistics: ``mlr histogram``, ``mlr stats1``, ``mlr stats2``.
* Cross-record arithmetic: ``mlr step``.

For ``mlr put`` and ``mlr filter``:

* Miller's types for function processing are **empty-null** (empty string), **absent-null** (reads of unset right-hand sides, or fall-through non-explicit return values from user-defined functions), **error**, **string**, **float** (double-precision), **int** (64-bit signed), and **boolean**.

* On input, string values representable as numbers, e.g. "3" or "3.1", are treated as int or float, respectively. If a record has ``x=1,y=2`` then ``mlr put '$z=$x+$y'`` will produce ``x=1,y=2,z=3``, and ``mlr put '$z=$x.$y'`` does not give an error simply because the dot operator has been generalized to stringify non-strings.  To coerce back to string for processing, use the ``string`` function: ``mlr put '$z=string($x).string($y)'`` will produce ``x=1,y=2,z=12``.

* On input, string values representable as boolean  (e.g. ``"true"``, ``"false"``) are *not* automatically treated as boolean.  (This is because ``"true"`` and ``"false"`` are ordinary words, and auto string-to-boolean on a column consisting of words would result in some strings mixed with some booleans.) Use the ``boolean`` function to coerce: e.g. giving the record ``x=1,y=2,w=false`` to ``mlr put '$z=($x<$y) || boolean($w)'``.

* Functions take types as described in ``mlr --help-all-functions``: for example, ``log10`` takes float input and produces float output, ``gmt2sec`` maps string to int, and ``sec2gmt`` maps int to string.

* All math functions described in ``mlr --help-all-functions`` take integer as well as float input.

.. _reference-null-data:

Null data: empty and absent
----------------------------------------------------------------

One of Miller's key features is its support for **heterogeneous** data.  For example, take ``mlr sort``: if you try to sort on field ``hostname`` when not all records in the data stream *have* a field named ``hostname``, it is not an error (although you could pre-filter the data stream using ``mlr having-fields --at-least hostname then sort ...``).  Rather, records lacking one or more sort keys are simply output contiguously by ``mlr sort``.

Miller has two kinds of null data:

* **Empty (key present, value empty)**: a field name is present in a record (or in an out-of-stream variable) with empty value: e.g. ``x=,y=2`` in the data input stream, or assignment ``$x=""`` or ``@x=""`` in ``mlr put``.

* **Absent (key not present)**: a field name is not present, e.g. input record is ``x=1,y=2`` and a ``put`` or ``filter`` expression refers to ``$z``. Or, reading an out-of-stream variable which hasn't been assigned a value yet, e.g.  ``mlr put -q '@sum += $x; end{emit @sum}'`` or ``mlr put -q '@sum[$a][$b] += $x; end{emit @sum, "a", "b"}'``.

You can test these programatically using the functions ``is_empty``/``is_not_empty``, ``is_absent``/``is_present``, and ``is_null``/``is_not_null``. For the last pair, note that null means either empty or absent.

Rules for null-handling:

* Records with one or more empty sort-field values sort after records with all sort-field values present:

::

    $ mlr cat data/sort-null.dat
    a=3,b=2
    a=1,b=8
    a=,b=4
    x=9,b=10
    a=5,b=7

::

    $ mlr sort -n  a data/sort-null.dat
    a=1,b=8
    a=3,b=2
    a=5,b=7
    a=,b=4
    x=9,b=10

::

    $ mlr sort -nr a data/sort-null.dat
    a=,b=4
    a=5,b=7
    a=3,b=2
    a=1,b=8
    x=9,b=10

* Functions/operators which have one or more *empty* arguments produce empty output: e.g.

::

    $ echo 'x=2,y=3' | mlr put '$a=$x+$y'
    x=2,y=3,a=5

::

    $ echo 'x=,y=3' | mlr put '$a=$x+$y'
    x=,y=3,a=

::

    $ echo 'x=,y=3' | mlr put '$a=log($x);$b=log($y)'
    x=,y=3,a=,b=1.098612

with the exception that the ``min`` and ``max`` functions are special: if one argument is non-null, it wins:

::

    $ echo 'x=,y=3' | mlr put '$a=min($x,$y);$b=max($x,$y)'
    x=,y=3,a=3,b=3

* Functions of *absent* variables (e.g. ``mlr put '$y = log10($nonesuch)'``) evaluate to absent, and arithmetic/bitwise/boolean operators with both operands being absent evaluate to absent. Arithmetic operators with one absent operand return the other operand. More specifically, absent values act like zero for addition/subtraction, and one for multiplication: Furthermore, **any expression which evaluates to absent is not stored in the left-hand side of an assignment statement**:

::

    $ echo 'x=2,y=3' | mlr put '$a=$u+$v; $b=$u+$y; $c=$x+$y'
    x=2,y=3,b=3,c=5

::

    $ echo 'x=2,y=3' | mlr put '$a=min($x,$v);$b=max($u,$y);$c=min($u,$v)'
    x=2,y=3,a=2,b=3

* Likewise, for assignment to maps, **absent-valued keys or values result in a skipped assignment**.

The reasoning is as follows:

* Empty values are explicit in the data so they should explicitly affect accumulations: ``mlr put '@sum += $x'`` should accumulate numeric ``x`` values into the sum but an empty ``x``, when encountered in the input data stream, should make the sum non-numeric. To work around this you can use the ``is_not_null`` function as follows: ``mlr put 'is_not_null($x) { @sum += $x }'``

* Absent stream-record values should not break accumulations, since Miller by design handles heterogenous data: the running ``@sum`` in ``mlr put '@sum += $x'`` should not be invalidated for records which have no ``x``.

* Absent out-of-stream-variable values are precisely what allow you to write ``mlr put '@sum += $x'``. Otherwise you would have to write ``mlr put 'begin{@sum = 0}; @sum += $x'`` -- which is tolerable -- but for ``mlr put 'begin{...}; @sum[$a][$b] += $x'`` you'd have to pre-initialize ``@sum`` for all values of ``$a`` and ``$b`` in your input data stream, which is intolerable.

* The penalty for the absent feature is that misspelled variables can be hard to find: e.g. in ``mlr put 'begin{@sumx = 10}; ...; update @sumx somehow per-record; ...; end {@something = @sum * 2}'`` the accumulator is spelt ``@sumx`` in the begin-block but ``@sum`` in the end-block, where since it is absent, ``@sum*2`` evaluates to 2. See also the section on :ref:`reference-dsl-errors-and-transparency`.

Since absent plus absent is absent (and likewise for other operators), accumulations such as ``@sum += $x`` work correctly on heterogenous data, as do within-record formulas if both operands are absent. If one operand is present, you may get behavior you don't desire.  To work around this -- namely, to set an output field only for records which have all the inputs present -- you can use a pattern-action block with ``is_present``:

::

    $ mlr cat data/het.dkvp
    resource=/path/to/file,loadsec=0.45,ok=true
    record_count=100,resource=/path/to/file
    resource=/path/to/second/file,loadsec=0.32,ok=true
    record_count=150,resource=/path/to/second/file
    resource=/some/other/path,loadsec=0.97,ok=false

::

    $ mlr put 'is_present($loadsec) { $loadmillis = $loadsec * 1000 }' data/het.dkvp
    resource=/path/to/file,loadsec=0.45,ok=true,loadmillis=450.000000
    record_count=100,resource=/path/to/file
    resource=/path/to/second/file,loadsec=0.32,ok=true,loadmillis=320.000000
    record_count=150,resource=/path/to/second/file
    resource=/some/other/path,loadsec=0.97,ok=false,loadmillis=970.000000

::

    $ mlr put '$loadmillis = (is_present($loadsec) ? $loadsec : 0.0) * 1000' data/het.dkvp
    resource=/path/to/file,loadsec=0.45,ok=true,loadmillis=450.000000
    record_count=100,resource=/path/to/file,loadmillis=0.000000
    resource=/path/to/second/file,loadsec=0.32,ok=true,loadmillis=320.000000
    record_count=150,resource=/path/to/second/file,loadmillis=0.000000
    resource=/some/other/path,loadsec=0.97,ok=false,loadmillis=970.000000

If you're interested in a formal description of how empty and absent fields participate in arithmetic, here's a table for plus (other arithmetic/boolean/bitwise operators are similar):

::

    $ mlr --print-type-arithmetic-info
    (+)    | error  absent empty  string int    float  bool  
    ------ + ------ ------ ------ ------ ------ ------ ------
    error  | error  error  error  error  error  error  error 
    absent | error  absent absent error  int    float  error 
    empty  | error  absent empty  error  empty  empty  error 
    string | error  error  error  error  error  error  error 
    int    | error  int    empty  error  int    float  error 
    float  | error  float  empty  error  float  float  error 
    bool   | error  error  error  error  error  error  error 

String literals
----------------------------------------------------------------

You can use the following backslash escapes for strings such as between the double quotes in contexts such as ``mlr filter '$name =~ "..."'``, ``mlr put '$name = $othername . "..."'``, ``mlr put '$name = sub($name, "...", "...")``, etc.:

* ``\a``: ASCII code 0x07 (alarm/bell)
* ``\b``: ASCII code 0x08 (backspace)
* ``\f``: ASCII code 0x0c (formfeed)
* ``\n``: ASCII code 0x0a (LF/linefeed/newline)
* ``\r``: ASCII code 0x0d (CR/carriage return)
* ``\t``: ASCII code 0x09 (tab)
* ``\v``: ASCII code 0x0b (vertical tab)
* ``\\``: backslash
* ``\"``: double quote
* ``\123``: Octal 123, etc. for ``\000`` up to ``\377``
* ``\x7f``: Hexadecimal 7f, etc. for ``\x00`` up to ``\xff``

See also https://en.wikipedia.org/wiki/Escape_sequences_in_C.

These replacements apply only to strings you key in for the DSL expressions for ``filter`` and ``put``: that is, if you type ``\t`` in a string literal for a ``filter``/``put`` expression, it will be turned into a tab character. If you want a backslash followed by a ``t``, then please type ``\\t``.

However, these replacements are not done automatically within your data stream. If you wish to make these replacements, you can do, for example, for a field named ``field``, ``mlr put '$field = gsub($field, "\\t", "\t")'``. If you need to make such a replacement for all fields in your data, you should probably simply use the system ``sed`` command.

Regular expressions
----------------------------------------------------------------

Miller lets you use regular expressions (of type POSIX.2) in the following contexts:

* In ``mlr filter`` with ``=~`` or ``!=~``, e.g. ``mlr filter '$url =~ "http.*com"'``

* In ``mlr put`` with ``sub`` or ``gsub``, e.g. ``mlr put '$url = sub($url, "http.*com", "")'``

* In ``mlr having-fields``, e.g. ``mlr having-fields --any-matching '^sda[0-9]'``

* In ``mlr cut``, e.g. ``mlr cut -r -f '^status$,^sda[0-9]'``

* In ``mlr rename``, e.g. ``mlr rename -r '^(sda[0-9]).*$,dev/\1'``

* In ``mlr grep``, e.g. ``mlr --csv grep 00188555487 myfiles*.csv``

Points demonstrated by the above examples:

* There are no implicit start-of-string or end-of-string anchors; please use ``^`` and/or ``$`` explicitly.

* Miller regexes are wrapped with double quotes rather than slashes.

* The ``i`` after the ending double quote indicates a case-insensitive regex.

* Capture groups are wrapped with ``(...)`` rather than ``\(...\)``; use ``\(`` and ``\)`` to match against parentheses.

For ``filter`` and ``put``, if the regular expression is a string literal (the normal case), it is precompiled at process start and reused thereafter, which is efficient. If the regular expression is a more complex expression, including string concatenation using ``.``, or a column name (in which case you can take regular expressions from input data!), then regexes are compiled on each record which works but is less efficient. As well, in this case there is no way to specify case-insensitive matching.

Example:

::

    $ cat data/regex-in-data.dat
    name=jane,regex=^j.*e$
    name=bill,regex=^b[ou]ll$
    name=bull,regex=^b[ou]ll$

::

    $ mlr filter '$name =~ $regex' data/regex-in-data.dat
    name=jane,regex=^j.*e$
    name=bull,regex=^b[ou]ll$

Regex captures
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Regex captures of the form ``\0`` through ``\9`` are supported as

* Captures have in-function context for ``sub`` and ``gsub``. For example, the first ``\1,\2`` pair belong to the first ``sub`` and the second ``\1,\2`` pair belong to the second ``sub``:

::

    mlr put '$b = sub($a, "(..)_(...)", "\2-\1"); $c = sub($a, "(..)_(.)(..)", ":\1:\2:\3")'

* Captures endure for the entirety of a ``put`` for the ``=~`` and ``!=~`` operators. For example, here the ``\1,\2`` are set by the ``=~`` operator and are used by both subsequent assignment statements:

::

    mlr put '$a =~ "(..)_(....); $b = "left_\1"; $c = "right_\2"'

* The captures are not retained across multiple puts. For example, here the ``\1,\2`` won't be expanded from the regex capture:

::

    mlr put '$a =~ "(..)_(....)' then {... something else ...} then put '$b = "left_\1"; $c = "right_\2"'

* Captures are ignored in ``filter`` for the ``=~`` and ``!=~`` operators. For example, there is no mechanism provided to refer to the first ``(..)`` as ``\1`` or to the second ``(....)`` as ``\2`` in the following filter statement:

::

    mlr filter '$a =~ "(..)_(....)'

* Up to nine matches are supported: ``\1`` through ``\9``, while ``\0`` is the entire match string; ``\15`` is treated as ``\1`` followed by an unrelated ``5``.

Arithmetic
----------------------------------------------------------------

Input scanning
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Numbers in Miller are double-precision float or 64-bit signed integers. Anything scannable as int, e.g ``123`` or ``0xabcd``, is treated as an integer; otherwise, input scannable as float (``4.56`` or ``8e9``) is treated as float; everything else is a string.

If you want all numbers to be treated as floats, then you may use ``float()`` in your filter/put expressions (e.g. replacing ``$c = $a * $b`` with ``$c = float($a) * float($b)``) -- or, more simply, use ``mlr filter -F`` and ``mlr put -F`` which forces all numeric input, whether from expression literals or field values, to float. Likewise ``mlr stats1 -F`` and ``mlr step -F`` force integerable accumulators (such as ``count``) to be done in floating-point.

Conversion by math routines
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

For most math functions, integers are cast to float on input, and produce float output: e.g. ``exp(0) = 1.0`` rather than ``1``.  The following, however, produce integer output if their inputs are integers: ``+`` ``-`` ``*`` ``/`` ``//`` ``%`` ``abs`` ``ceil`` ``floor`` ``max`` ``min`` ``round`` ``roundm`` ``sgn``. As well, ``stats1 -a min``, ``stats1 -a max``, ``stats1 -a sum``, ``step -a delta``, and ``step -a rsum`` produce integer output if their inputs are integers.

Conversion by arithmetic operators
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

The sum, difference, and product of integers is again integer, except for when that would overflow a 64-bit integer at which point Miller converts the result to float.

The short of it is that Miller does this transparently for you so you needn't think about it.

Implementation details of this, for the interested: integer adds and subtracts overflow by at most one bit so it suffices to check sign-changes. Thus, Miller allows you to add and subtract arbitrary 64-bit signed integers, converting only to float precisely when the result is less than -2\ :sup:`63` or greater than 2\ :sup:`63`\ -1.  Multiplies, on the other hand, can overflow by a word size and a sign-change technique does not suffice to detect overflow. Instead Miller tests whether the floating-point product exceeds the representable integer range. Now, 64-bit integers have 64-bit precision while IEEE-doubles have only 52-bit mantissas -- so, there are 53 bits including implicit leading one.  The following experiment explicitly demonstrates the resolution at this range:

::

    64-bit integer     64-bit integer     Casted to double           Back to 64-bit
    in hex           in decimal                                    integer
    0x7ffffffffffff9ff 9223372036854774271 9223372036854773760.000000 0x7ffffffffffff800
    0x7ffffffffffffa00 9223372036854774272 9223372036854773760.000000 0x7ffffffffffff800
    0x7ffffffffffffbff 9223372036854774783 9223372036854774784.000000 0x7ffffffffffffc00
    0x7ffffffffffffc00 9223372036854774784 9223372036854774784.000000 0x7ffffffffffffc00
    0x7ffffffffffffdff 9223372036854775295 9223372036854774784.000000 0x7ffffffffffffc00
    0x7ffffffffffffe00 9223372036854775296 9223372036854775808.000000 0x8000000000000000
    0x7ffffffffffffffe 9223372036854775806 9223372036854775808.000000 0x8000000000000000
    0x7fffffffffffffff 9223372036854775807 9223372036854775808.000000 0x8000000000000000

That is, one cannot check an integer product to see if it is precisely greater than 2\ :sup:`63`\ -1 or less than -2\ :sup:`63` using either integer arithmetic (it may have already overflowed) or using double-precision (due to granularity).  Instead Miller checks for overflow in 64-bit integer multiplication by seeing whether the absolute value of the double-precision product exceeds the largest representable IEEE double less than 2\ :sup:`63`, which we see from the listing above is 9223372036854774784. (An alternative would be to do all integer multiplies using handcrafted multi-word 128-bit arithmetic.  This approach is not taken.)

Pythonic division
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Division and remainder are `pythonic <http://python-history.blogspot.com/2010/08/why-pythons-integer-division-floors.html>`_:

* Quotient of integers is floating-point: ``7/2`` is ``3.5``.
* Integer division is done with ``//``: ``7//2`` is ``3``.  This rounds toward the negative.
* Remainders are non-negative.

On-line help
----------------------------------------------------------------

Examples:

::

    $ mlr --help
    Usage: mlr [I/O options] {verb} [verb-dependent options ...] {zero or more file names}
    
    Command-line-syntax examples:
      mlr --csv cut -f hostname,uptime mydata.csv
      mlr --tsv --rs lf filter '$status != "down" && $upsec >= 10000' *.tsv
      mlr --nidx put '$sum = $7 < 0.0 ? 3.5 : $7 + 2.1*$8' *.dat
      grep -v '^#' /etc/group | mlr --ifs : --nidx --opprint label group,pass,gid,member then sort -f group
      mlr join -j account_id -f accounts.dat then group-by account_name balances.dat
      mlr --json put '$attr = sub($attr, "([0-9]+)_([0-9]+)_.*", "\1:\2")' data/*.json
      mlr stats1 -a min,mean,max,p10,p50,p90 -f flag,u,v data/*
      mlr stats2 -a linreg-pca -f u,v -g shape data/*
      mlr put -q '@sum[$a][$b] += $x; end {emit @sum, "a", "b"}' data/*
      mlr --from estimates.tbl put '
      for (k,v in $*) {
        if (is_numeric(v) && k =~ "^[t-z].*$") {
          $sum += v; $count += 1
        }
      }
      $mean = $sum / $count # no assignment if count unset'
      mlr --from infile.dat put -f analyze.mlr
      mlr --from infile.dat put 'tee > "./taps/data-".$a."-".$b, $*'
      mlr --from infile.dat put 'tee | "gzip > ./taps/data-".$a."-".$b.".gz", $*'
      mlr --from infile.dat put -q '@v=$*; dump | "jq .[]"'
      mlr --from infile.dat put  '(NR % 1000 == 0) { print > stderr, "Checkpoint ".NR}'
    
    Data-format examples:
      DKVP: delimited key-value pairs (Miller default format)
      +---------------------+
      | apple=1,bat=2,cog=3 | Record 1: "apple" => "1", "bat" => "2", "cog" => "3"
      | dish=7,egg=8,flint  | Record 2: "dish" => "7", "egg" => "8", "3" => "flint"
      +---------------------+
    
      NIDX: implicitly numerically indexed (Unix-toolkit style)
      +---------------------+
      | the quick brown     | Record 1: "1" => "the", "2" => "quick", "3" => "brown"
      | fox jumped          | Record 2: "1" => "fox", "2" => "jumped"
      +---------------------+
    
      CSV/CSV-lite: comma-separated values with separate header line
      +---------------------+
      | apple,bat,cog       |
      | 1,2,3               | Record 1: "apple => "1", "bat" => "2", "cog" => "3"
      | 4,5,6               | Record 2: "apple" => "4", "bat" => "5", "cog" => "6"
      +---------------------+
    
      Tabular JSON: nested objects are supported, although arrays within them are not:
      +---------------------+
      | {                   |
      |  "apple": 1,        | Record 1: "apple" => "1", "bat" => "2", "cog" => "3"
      |  "bat": 2,          |
      |  "cog": 3           |
      | }                   |
      | {                   |
      |   "dish": {         | Record 2: "dish:egg" => "7", "dish:flint" => "8", "garlic" => ""
      |     "egg": 7,       |
      |     "flint": 8      |
      |   },                |
      |   "garlic": ""      |
      | }                   |
      +---------------------+
    
      PPRINT: pretty-printed tabular
      +---------------------+
      | apple bat cog       |
      | 1     2   3         | Record 1: "apple => "1", "bat" => "2", "cog" => "3"
      | 4     5   6         | Record 2: "apple" => "4", "bat" => "5", "cog" => "6"
      +---------------------+
    
      XTAB: pretty-printed transposed tabular
      +---------------------+
      | apple 1             | Record 1: "apple" => "1", "bat" => "2", "cog" => "3"
      | bat   2             |
      | cog   3             |
      |                     |
      | dish 7              | Record 2: "dish" => "7", "egg" => "8"
      | egg  8              |
      +---------------------+
    
      Markdown tabular (supported for output only):
      +-----------------------+
      | | apple | bat | cog | |
      | | ---   | --- | --- | |
      | | 1     | 2   | 3   | | Record 1: "apple => "1", "bat" => "2", "cog" => "3"
      | | 4     | 5   | 6   | | Record 2: "apple" => "4", "bat" => "5", "cog" => "6"
      +-----------------------+
    
    Help options:
      -h or --help                 Show this message.
      --version                    Show the software version.
      {verb name} --help           Show verb-specific help.
      --help-all-verbs             Show help on all verbs.
      -l or --list-all-verbs       List only verb names.
      -L                           List only verb names, one per line.
      -f or --help-all-functions   Show help on all built-in functions.
      -F                           Show a bare listing of built-in functions by name.
      -k or --help-all-keywords    Show help on all keywords.
      -K                           Show a bare listing of keywords by name.
    
    Customization via .mlrrc:
    You can set up personal defaults via a $HOME/.mlrrc and/or ./.mlrrc.
    For example, if you usually process CSV, then you can put "--csv" in your .mlrrc file
    and that will be the default input/output format unless otherwise specified on the command line.
    
    The .mlrrc file format is one "--flag" or "--option value" per line, with the leading "--" optional.
    Hash-style comments and blank lines are ignored.
    
    Sample .mlrrc:
    # Input and output formats are CSV by default (unless otherwise specified
    # on the mlr command line):
    csv
    # These are no-ops for CSV, but when I do use JSON output, I want these
    # pretty-printing options to be used:
    jvstack
    jlistwrap
    
    How to specify location of .mlrrc:
    * If $MLRRC is set:
      o If its value is "__none__" then no .mlrrc files are processed.
      o Otherwise, its value (as a filename) is loaded and processed. If there are syntax
        errors, they abort mlr with a usage message (as if you had mistyped something on the
        command line). If the file can't be loaded at all, though, it is silently skipped.
      o Any .mlrrc in your home directory or current directory is ignored whenever $MLRRC is
        set in the environment.
    * Otherwise:
      o If $HOME/.mlrrc exists, it's then processed as above.
      o If ./.mlrrc exists, it's then also processed as above.
      (I.e. current-directory .mlrrc defaults are stacked over home-directory .mlrrc defaults.)
    
    See also:
    https://johnkerl.org/miller/doc/customization.html
    
    Verbs:
       altkv bar bootstrap cat check clean-whitespace count count-distinct
       count-similar cut decimate fill-down fill-empty filter format-values
       fraction grep group-by group-like having-fields head histogram join label
       least-frequent merge-fields most-frequent nest nothing put regularize
       remove-empty-columns rename reorder repeat reshape sample sec2gmt
       sec2gmtdate seqgen shuffle skip-trivial-records sort sort-within-records
       stats1 stats2 step tac tail tee top uniq unsparsify
    
    Functions for the filter and put verbs:
       + + - - * / // .+ .+ .- .- .* ./ .// % ** | ^ & ~ << >> bitcount == != =~
       !=~ > >= < <= && || ^^ ! ? : . gsub regextract regextract_or_else strlen sub
       ssub substr tolower toupper truncate capitalize lstrip rstrip strip
       collapse_whitespace clean_whitespace system abs acos acosh asin asinh atan
       atan2 atanh cbrt ceil cos cosh erf erfc exp expm1 floor invqnorm log log10
       log1p logifit madd max mexp min mmul msub pow qnorm round roundm sgn sin
       sinh sqrt tan tanh urand urandrange urand32 urandint dhms2fsec dhms2sec
       fsec2dhms fsec2hms gmt2sec localtime2sec hms2fsec hms2sec sec2dhms sec2gmt
       sec2gmt sec2gmtdate sec2localtime sec2localtime sec2localdate sec2hms
       strftime strftime_local strptime strptime_local systime is_absent is_bool
       is_boolean is_empty is_empty_map is_float is_int is_map is_nonempty_map
       is_not_empty is_not_map is_not_null is_null is_numeric is_present is_string
       asserting_absent asserting_bool asserting_boolean asserting_empty
       asserting_empty_map asserting_float asserting_int asserting_map
       asserting_nonempty_map asserting_not_empty asserting_not_map
       asserting_not_null asserting_null asserting_numeric asserting_present
       asserting_string boolean float fmtnum hexfmt int string typeof depth haskey
       joink joinkv joinv leafcount length mapdiff mapexcept mapselect mapsum
       splitkv splitkvx splitnv splitnvx
    
    Please use "mlr --help-function {function name}" for function-specific help.
    
    Data-format options, for input, output, or both:
      --idkvp   --odkvp   --dkvp      Delimited key-value pairs, e.g "a=1,b=2"
                                      (this is Miller's default format).
    
      --inidx   --onidx   --nidx      Implicitly-integer-indexed fields
                                      (Unix-toolkit style).
      -T                              Synonymous with "--nidx --fs tab".
    
      --icsv    --ocsv    --csv       Comma-separated value (or tab-separated
                                      with --fs tab, etc.)
    
      --itsv    --otsv    --tsv       Keystroke-savers for "--icsv --ifs tab",
                                      "--ocsv --ofs tab", "--csv --fs tab".
      --iasv    --oasv    --asv       Similar but using ASCII FS 0x1f and RS 0x1e
      --iusv    --ousv    --usv       Similar but using Unicode FS U+241F (UTF-8 0xe2909f)
                                      and RS U+241E (UTF-8 0xe2909e)
    
      --icsvlite --ocsvlite --csvlite Comma-separated value (or tab-separated
                                      with --fs tab, etc.). The 'lite' CSV does not handle
                                      RFC-CSV double-quoting rules; is slightly faster;
                                      and handles heterogeneity in the input stream via
                                      empty newline followed by new header line. See also
                                      http://johnkerl.org/miller/doc/file-formats.html#CSV/TSV/etc.
    
      --itsvlite --otsvlite --tsvlite Keystroke-savers for "--icsvlite --ifs tab",
                                      "--ocsvlite --ofs tab", "--csvlite --fs tab".
      -t                              Synonymous with --tsvlite.
      --iasvlite --oasvlite --asvlite Similar to --itsvlite et al. but using ASCII FS 0x1f and RS 0x1e
      --iusvlite --ousvlite --usvlite Similar to --itsvlite et al. but using Unicode FS U+241F (UTF-8 0xe2909f)
                                      and RS U+241E (UTF-8 0xe2909e)
    
      --ipprint --opprint --pprint    Pretty-printed tabular (produces no
                                      output until all input is in).
                          --right     Right-justifies all fields for PPRINT output.
                          --barred    Prints a border around PPRINT output
                                      (only available for output).
    
                --omd                 Markdown-tabular (only available for output).
    
      --ixtab   --oxtab   --xtab      Pretty-printed vertical-tabular.
                          --xvright   Right-justifies values for XTAB format.
    
      --ijson   --ojson   --json      JSON tabular: sequence or list of one-level
                                      maps: {...}{...} or [{...},{...}].
        --json-map-arrays-on-input    JSON arrays are unmillerable. --json-map-arrays-on-input
        --json-skip-arrays-on-input   is the default: arrays are converted to integer-indexed
        --json-fatal-arrays-on-input  maps. The other two options cause them to be skipped, or
                                      to be treated as errors.  Please use the jq tool for full
                                      JSON (pre)processing.
                          --jvstack   Put one key-value pair per line for JSON
                                      output.
                    --jsonx --ojsonx  Keystroke-savers for --json --jvstack
                    --jsonx --ojsonx  and --ojson --jvstack, respectively.
                          --jlistwrap Wrap JSON output in outermost [ ].
                        --jknquoteint Do not quote non-string map keys in JSON output.
                         --jvquoteall Quote map values in JSON output, even if they're
                                      numeric.
                  --jflatsep {string} Separator for flattening multi-level JSON keys,
                                      e.g. '{"a":{"b":3}}' becomes a:b => 3 for
                                      non-JSON formats. Defaults to :.
    
      -p is a keystroke-saver for --nidx --fs space --repifs
    
      Examples: --csv for CSV-formatted input and output; --idkvp --opprint for
      DKVP-formatted input and pretty-printed output.
    
      Please use --iformat1 --oformat2 rather than --format1 --oformat2.
      The latter sets up input and output flags for format1, not all of which
      are overridden in all cases by setting output format to format2.
    
    Comments in data:
      --skip-comments                 Ignore commented lines (prefixed by "#")
                                      within the input.
      --skip-comments-with {string}   Ignore commented lines within input, with
                                      specified prefix.
      --pass-comments                 Immediately print commented lines (prefixed by "#")
                                      within the input.
      --pass-comments-with {string}   Immediately print commented lines within input, with
                                      specified prefix.
    Notes:
    * Comments are only honored at the start of a line.
    * In the absence of any of the above four options, comments are data like
      any other text.
    * When pass-comments is used, comment lines are written to standard output
      immediately upon being read; they are not part of the record stream.
      Results may be counterintuitive. A suggestion is to place comments at the
      start of data files.
    
    Format-conversion keystroke-saver options, for input, output, or both:
    As keystroke-savers for format-conversion you may use the following:
            --c2t --c2d --c2n --c2j --c2x --c2p --c2m
      --t2c       --t2d --t2n --t2j --t2x --t2p --t2m
      --d2c --d2t       --d2n --d2j --d2x --d2p --d2m
      --n2c --n2t --n2d       --n2j --n2x --n2p --n2m
      --j2c --j2t --j2d --j2n       --j2x --j2p --j2m
      --x2c --x2t --x2d --x2n --x2j       --x2p --x2m
      --p2c --p2t --p2d --p2n --p2j --p2x       --p2m
    The letters c t d n j x p m refer to formats CSV, TSV, DKVP, NIDX, JSON, XTAB,
    PPRINT, and markdown, respectively. Note that markdown format is available for
    output only.
    
    Compressed-data options:
      --prepipe {command} This allows Miller to handle compressed inputs. You can do
      without this for single input files, e.g. "gunzip < myfile.csv.gz | mlr ...".
    
      However, when multiple input files are present, between-file separations are
      lost; also, the FILENAME variable doesn't iterate. Using --prepipe you can
      specify an action to be taken on each input file. This pre-pipe command must
      be able to read from standard input; it will be invoked with
        {command} < {filename}.
      Examples:
        mlr --prepipe 'gunzip'
        mlr --prepipe 'zcat -cf'
        mlr --prepipe 'xz -cd'
        mlr --prepipe cat
        mlr --prepipe-gunzip
        mlr --prepipe-zcat
      Note that this feature is quite general and is not limited to decompression
      utilities. You can use it to apply per-file filters of your choice.
      For output compression (or other) utilities, simply pipe the output:
        mlr ... | {your compression command}
    
      There are shorthands --prepipe-zcat and --prepipe-gunzip which are
      valid in .mlrrc files. The --prepipe flag is not valid in .mlrrc
      files since that would put execution of the prepipe command under 
      control of the .mlrrc file.
    
    Separator options, for input, output, or both:
      --rs     --irs     --ors              Record separators, e.g. 'lf' or '\r\n'
      --fs     --ifs     --ofs  --repifs    Field separators, e.g. comma
      --ps     --ips     --ops              Pair separators, e.g. equals sign
    
      Notes about line endings:
      * Default line endings (--irs and --ors) are "auto" which means autodetect from
        the input file format, as long as the input file(s) have lines ending in either
        LF (also known as linefeed, '\n', 0x0a, Unix-style) or CRLF (also known as
        carriage-return/linefeed pairs, '\r\n', 0x0d 0x0a, Windows style).
      * If both irs and ors are auto (which is the default) then LF input will lead to LF
        output and CRLF input will lead to CRLF output, regardless of the platform you're
        running on.
      * The line-ending autodetector triggers on the first line ending detected in the input
        stream. E.g. if you specify a CRLF-terminated file on the command line followed by an
        LF-terminated file then autodetected line endings will be CRLF.
      * If you use --ors {something else} with (default or explicitly specified) --irs auto
        then line endings are autodetected on input and set to what you specify on output.
      * If you use --irs {something else} with (default or explicitly specified) --ors auto
        then the output line endings used are LF on Unix/Linux/BSD/MacOSX, and CRLF on Windows.
    
      Notes about all other separators:
      * IPS/OPS are only used for DKVP and XTAB formats, since only in these formats
        do key-value pairs appear juxtaposed.
      * IRS/ORS are ignored for XTAB format. Nominally IFS and OFS are newlines;
        XTAB records are separated by two or more consecutive IFS/OFS -- i.e.
        a blank line. Everything above about --irs/--ors/--rs auto becomes --ifs/--ofs/--fs
        auto for XTAB format. (XTAB's default IFS/OFS are "auto".)
      * OFS must be single-character for PPRINT format. This is because it is used
        with repetition for alignment; multi-character separators would make
        alignment impossible.
      * OPS may be multi-character for XTAB format, in which case alignment is
        disabled.
      * TSV is simply CSV using tab as field separator ("--fs tab").
      * FS/PS are ignored for markdown format; RS is used.
      * All FS and PS options are ignored for JSON format, since they are not relevant
        to the JSON format.
      * You can specify separators in any of the following ways, shown by example:
        - Type them out, quoting as necessary for shell escapes, e.g.
          "--fs '|' --ips :"
        - C-style escape sequences, e.g. "--rs '\r\n' --fs '\t'".
        - To avoid backslashing, you can use any of the following names:
          cr crcr newline lf lflf crlf crlfcrlf tab space comma pipe slash colon semicolon equals
      * Default separators by format:
          File format  RS       FS       PS
          gen          N/A      (N/A)    (N/A)
          dkvp         auto     ,        =
          json         auto     (N/A)    (N/A)
          nidx         auto     space    (N/A)
          csv          auto     ,        (N/A)
          csvlite      auto     ,        (N/A)
          markdown     auto     (N/A)    (N/A)
          pprint       auto     space    (N/A)
          xtab         (N/A)    auto     space
    
    Relevant to CSV/CSV-lite input only:
      --implicit-csv-header Use 1,2,3,... as field labels, rather than from line 1
                         of input files. Tip: combine with "label" to recreate
                         missing headers.
      --no-implicit-csv-header Do not use --implicit-csv-header. This is the default
                         anyway -- the main use is for the flags to 'mlr join' if you have
                         main file(s) which are headerless but you want to join in on
                         a file which does have a CSV header. Then you could use
                         'mlr --csv --implicit-csv-header join --no-implicit-csv-header
                         -l your-join-in-with-header.csv ... your-headerless.csv'
      --allow-ragged-csv-input|--ragged If a data line has fewer fields than the header line,
                         fill remaining keys with empty string. If a data line has more
                         fields than the header line, use integer field labels as in
                         the implicit-header case.
      --headerless-csv-output   Print only CSV data lines.
      -N                 Keystroke-saver for --implicit-csv-header --headerless-csv-output.
    
    Double-quoting for CSV output:
      --quote-all        Wrap all fields in double quotes
      --quote-none       Do not wrap any fields in double quotes, even if they have
                         OFS or ORS in them
      --quote-minimal    Wrap fields in double quotes only if they have OFS or ORS
                         in them (default)
      --quote-numeric    Wrap fields in double quotes only if they have numbers
                         in them
      --quote-original   Wrap fields in double quotes if and only if they were
                         quoted on input. This isn't sticky for computed fields:
                         e.g. if fields a and b were quoted on input and you do
                         "put '$c = $a . $b'" then field c won't inherit a or b's
                         was-quoted-on-input flag.
    
    Numerical formatting:
      --ofmt {format}    E.g. %.18lf, %.0lf. Please use sprintf-style codes for
                         double-precision. Applies to verbs which compute new
                         values, e.g. put, stats1, stats2. See also the fmtnum
                         function within mlr put (mlr --help-all-functions).
                         Defaults to %lf.
    
    Other options:
      --seed {n} with n of the form 12345678 or 0xcafefeed. For put/filter
                         urand()/urandint()/urand32().
      --nr-progress-mod {m}, with m a positive integer: print filename and record
                         count to stderr every m input records.
      --from {filename}  Use this to specify an input file before the verb(s),
                         rather than after. May be used more than once. Example:
                         "mlr --from a.dat --from b.dat cat" is the same as
                         "mlr cat a.dat b.dat".
      -n                 Process no input files, nor standard input either. Useful
                         for mlr put with begin/end statements only. (Same as --from
                         /dev/null.) Also useful in "mlr -n put -v '...'" for
                         analyzing abstract syntax trees (if that's your thing).
      -I                 Process files in-place. For each file name on the command
                         line, output is written to a temp file in the same
                         directory, which is then renamed over the original. Each
                         file is processed in isolation: if the output format is
                         CSV, CSV headers will be present in each output file;
                         statistics are only over each file's own records; and so on.
    
    Then-chaining:
    Output of one verb may be chained as input to another using "then", e.g.
      mlr stats1 -a min,mean,max -f flag,u,v -g color then sort -f color
    
    Auxiliary commands:
    Miller has a few otherwise-standalone executables packaged within it.
    They do not participate in any other parts of Miller.
    Available subcommands:
      aux-list
      lecat
      termcvt
      hex
      unhex
      netbsd-strptime
    For more information, please invoke mlr {subcommand} --help
    
    For more information please see http://johnkerl.org/miller/doc and/or
    http://github.com/johnkerl/miller. This is Miller version v5.10.2-dev.

::

    $ mlr sort --help
    Usage: mlr sort {flags}
    Flags:
      -f  {comma-separated field names}  Lexical ascending
      -n  {comma-separated field names}  Numerical ascending; nulls sort last
      -nf {comma-separated field names}  Same as -n
      -r  {comma-separated field names}  Lexical descending
      -nr {comma-separated field names}  Numerical descending; nulls sort first
    Sorts records primarily by the first specified field, secondarily by the second
    field, and so on.  (Any records not having all specified sort keys will appear
    at the end of the output, in the order they were encountered, regardless of the
    specified sort order.) The sort is stable: records that compare equal will sort
    in the order they were encountered in the input record stream.
    
    Example:
      mlr sort -f a,b -nr x,y,z
    which is the same as:
      mlr sort -f a -f b -nr x -nr y -nr z
