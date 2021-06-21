..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Special symbols and formatting
==============================

How can I handle commas-as-data in various formats?
----------------------------------------------------------------

:doc:`CSV <file-formats>` handles this well and by design:

.. code-block:: none
   :emphasize-lines: 1-1

    cat commas.csv
    Name,Role
    "Xiao, Lin",administrator
    "Khavari, Darius",tester

Likewise :ref:`file-formats-json`:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --ojson cat commas.csv
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
   :emphasize-lines: 1-1

    mlr --icsv --oxtab cat commas.csv
    Name Xiao, Lin
    Role administrator
    
    Name Khavari, Darius
    Role tester

But for :ref:`Key-value_pairs <file-formats-dkvp>` and :ref:`index-numbered <file-formats-nidx>`, commas are the default field separator. And -- as of Miller 5.4.0 anyway -- there is no CSV-style double-quote-handling like there is for CSV. So commas within the data look like delimiters:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --odkvp cat commas.csv
    Name=Xiao, Lin,Role=administrator
    Name=Khavari, Darius,Role=tester

One solution is to use a different delimiter, such as a pipe character:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --odkvp --ofs pipe cat commas.csv
    Name=Xiao, Lin|Role=administrator
    Name=Khavari, Darius|Role=tester

To be extra-sure to avoid data/delimiter clashes, you can also use control
characters as delimiters -- here, control-A:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --odkvp --ofs '\001'  cat commas.csv | cat -v
    Name=Xiao, Lin\001Role=administrator
    Name=Khavari, Darius\001Role=tester

How can I handle field names with special symbols in them?
----------------------------------------------------------------

Simply surround the field names with curly braces:

.. code-block:: none
   :emphasize-lines: 1-1

    echo 'x.a=3,y:b=4,z/c=5' | mlr put '${product.all} = ${x.a} * ${y:b} * ${z/c}'
    x.a=3,y:b=4,z/c=5,product.all=60

How can I put single-quotes into strings?
----------------------------------------------------------------

This is a little tricky due to the shell's handling of quotes. For simplicity, let's first put an update script into a file:

.. code-block:: none

    $a = "It's OK, I said, then 'for now'."

.. code-block:: none
   :emphasize-lines: 1-1

    echo a=bcd | mlr put -f data/single-quote-example.mlr
    a=It's OK, I said, then 'for now'.

So, it's simple: Miller's DSL uses double quotes for strings, and you can put single quotes (or backslash-escaped double-quotes) inside strings, no problem.

Without putting the update expression in a file, it's messier:

.. code-block:: none
   :emphasize-lines: 1-1

    echo a=bcd | mlr put '$a="It'\''s OK, I said, '\''for now'\''."'
    a=It's OK, I said, 'for now'.

The idea is that the outermost single-quotes are to protect the ``put`` expression from the shell, and the double quotes within them are for Miller. To get a single quote in the middle there, you need to actually put it *outside* the single-quoting for the shell. The pieces are the following, all concatenated together:

* ``$a="It``
* ``\'``
* ``s OK, I said,``
* ``\'``
* ``for now``
* ``\'``
* ``.``

How to escape '?' in regexes?
----------------------------------------------------------------

One way is to use square brackets; an alternative is to use simple string-substitution rather than a regular expression.

.. code-block:: none
   :emphasize-lines: 1-1

    cat data/question.dat
    a=is it?,b=it is!
.. code-block:: none
   :emphasize-lines: 1-1

    mlr --oxtab put '$c = gsub($a, "[?]"," ...")' data/question.dat
    a is it?
    b it is!
    c is it ...
.. code-block:: none
   :emphasize-lines: 1-1

    mlr --oxtab put '$c = ssub($a, "?"," ...")' data/question.dat
    a is it?
    b it is!
    c is it ...

The ``ssub`` function exists precisely for this reason: so you don't have to escape anything.

How do I suppress numeric conversion?
----------------------------------------------------------------

**TL;DR use put -S**.

**TODO: probably remove this for Miller 6**

Within ``mlr put`` and ``mlr filter``, the default behavior for scanning input records is to parse them as integer, if possible, then as float, if possible, else leave them as string:

.. code-block:: none
   :emphasize-lines: 1-1

    cat data/scan-example-1.tbl
    value
    1
    2.0
    3x
    hello

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --pprint put '$copy = $value; $type = typeof($value)' data/scan-example-1.tbl
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
   :emphasize-lines: 1-1

    cat data/scan-example-2.tbl
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
   :emphasize-lines: 1-1

    mlr --pprint put '$copy = $value; $type = typeof($value)' data/scan-example-2.tbl
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
   :emphasize-lines: 1-1

    mlr --pprint put -S '$copy = $value; $type = typeof($value)' data/scan-example-2.tbl
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
