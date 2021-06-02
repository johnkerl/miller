..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

DSL reference
================================================================

Overview
----------------------------------------------------------------

Here's comparison of verbs and ``put``/``filter`` DSL expressions:

Example:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr stats1 -a sum -f x -g a data/small
    a=pan,x_sum=0.3467901443380824
    a=eks,x_sum=1.1400793586611044
    a=wye,x_sum=0.7778922255683036

* Verbs are coded in Go
* They run a bit faster
* They take fewer keystrokes
* There is less to learn
* Their customization is limited to each verb's options

Example:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr  put -q '@x_sum[$a] += $x; end{emit @x_sum, "a"}' data/small
    a=pan,x_sum=0.3467901443380824
    a=eks,x_sum=1.1400793586611044
    a=wye,x_sum=0.7778922255683036

* You get to write your own DSL expressions
* They run a bit slower
* They take more keystrokes
* There is more to learn
* They are highly customizable

Please see :doc:`reference-verbs` for information on verbs other than ``put`` and ``filter``.

The essential usages of ``mlr filter`` and ``mlr put`` are for record-selection and record-updating expressions, respectively. For example, given the following input data:

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729

you might retain only the records whose ``a`` field has value ``eks``:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr filter '$a == "eks"' data/small
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463

or you might add a new field which is a function of existing fields:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put '$ab = $a . "_" . $b ' data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533,ab=pan_pan
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797,ab=eks_pan
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776,ab=wye_wye
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463,ab=eks_wye
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,ab=wye_pan

The two verbs ``mlr filter`` and ``mlr put`` are essentially the same. The only differences are:

* Expressions sent to ``mlr filter`` must end with a boolean expression, which is the filtering criterion;

* ``mlr filter`` expressions may not reference the ``filter`` keyword within them; and

* ``mlr filter`` expressions may not use ``tee``, ``emit``, ``emitp``, or ``emitf``.

All the rest is the same: in particular, you can define and invoke functions and subroutines to help produce the final boolean statement, and record fields may be assigned to in the statements preceding the final boolean statement.

There are more details and more choices, of course, as detailed in the following sections.

Syntax
----------------------------------------------------------------

Expression formatting
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Multiple expressions may be given, separated by semicolons, and each may refer to the ones before:

.. code-block:: none
   :emphasize-lines: 1,1

    $ ruby -e '10.times{|i|puts "i=#{i}"}' | mlr --opprint put '$j = $i + 1; $k = $i +$j'
    i j  k
    0 1  1
    1 2  3
    2 3  5
    3 4  7
    4 5  9
    5 6  11
    6 7  13
    7 8  15
    8 9  17
    9 10 19

Newlines within the expression are ignored, which can help increase legibility of complex expressions:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint put '
      $nf       = NF;
      $nr       = NR;
      $fnr      = FNR;
      $filenum  = FILENUM;
      $filename = FILENAME
    ' data/small data/small2
    a   b   i     x                    y                    nf nr fnr filenum filename
    pan pan 1     0.3467901443380824   0.7268028627434533   5  1  1   1       data/small
    eks pan 2     0.7586799647899636   0.5221511083334797   5  2  2   1       data/small
    wye wye 3     0.20460330576630303  0.33831852551664776  5  3  3   1       data/small
    eks wye 4     0.38139939387114097  0.13418874328430463  5  4  4   1       data/small
    wye pan 5     0.5732889198020006   0.8636244699032729   5  5  5   1       data/small
    pan eks 9999  0.267481232652199086 0.557077185510228001 5  6  1   2       data/small2
    wye eks 10000 0.734806020620654365 0.884788571337605134 5  7  2   2       data/small2
    pan wye 10001 0.870530722602517626 0.009854780514656930 5  8  3   2       data/small2
    hat wye 10002 0.321507044286237609 0.568893318795083758 5  9  4   2       data/small2
    pan zee 10003 0.272054845593895200 0.425789896597056627 5  10 5   2       data/small2

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint filter '($x > 0.5 && $y < 0.5) || ($x < 0.5 && $y > 0.5)' then stats2 -a corr -f x,y data/medium
    x_y_corr
    -0.7479940285189345

.. _reference-dsl-expressions-from-files:

Expressions from files
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

The simplest way to enter expressions for ``put`` and ``filter`` is between single quotes on the command line, e.g.

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --from data/small put '$xy = sqrt($x**2 + $y**2)'
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533,xy=0.8052985815845617
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797,xy=0.9209978658539777
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776,xy=0.3953756915115773
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463,xy=0.40431685157744135
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,xy=1.036584492737304

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --from data/small put 'func f(a, b) { return sqrt(a**2 + b**2) } $xy = f($x, $y)'
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533,xy=0.8052985815845617
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797,xy=0.9209978658539777
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776,xy=0.3953756915115773
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463,xy=0.40431685157744135
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,xy=1.036584492737304

You may, though, find it convenient to put expressions into files for reuse, and read them
**using the -f option**. For example:

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat data/fe-example-3.mlr
    func f(a, b) {
      return sqrt(a**2 + b**2)
    }
    $xy = f($x, $y)

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --from data/small put -f data/fe-example-3.mlr
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533,xy=0.8052985815845617
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797,xy=0.9209978658539777
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776,xy=0.3953756915115773
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463,xy=0.40431685157744135
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,xy=1.036584492737304

If you have some of the logic in a file and you want to write the rest on the command line, you can **use the -f and -e options together**:

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat data/fe-example-4.mlr
    func f(a, b) {
      return sqrt(a**2 + b**2)
    }

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --from data/small put -f data/fe-example-4.mlr -e '$xy = f($x, $y)'
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533,xy=0.8052985815845617
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797,xy=0.9209978658539777
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776,xy=0.3953756915115773
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463,xy=0.40431685157744135
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,xy=1.036584492737304

A suggested use-case here is defining functions in files, and calling them from command-line expressions.

Another suggested use-case is putting default parameter values in files, e.g. using ``begin{@count=is_present(@count)?@count:10}`` in the file, where you can precede that using ``begin{@count=40}`` using ``-e``.

Moreover, you can have one or more ``-f`` expressions (maybe one function per file, for example) and one or more ``-e`` expressions on the command line.  If you mix ``-f`` and ``-e`` then the expressions are evaluated in the order encountered. (Since the expressions are all simply concatenated together in order, don't forget intervening semicolons: e.g. not ``mlr put -e '$x=1' -e '$y=2 ...'`` but rather ``mlr put -e '$x=1;' -e '$y=2' ...``.)

Semicolons, commas, newlines, and curly braces
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Miller uses **semicolons as statement separators**, not statement terminators. This means you can write:

.. code-block:: none

    mlr put 'x=1'
    mlr put 'x=1;$y=2'
    mlr put 'x=1;$y=2;'
    mlr put 'x=1;;;;$y=2;'

Semicolons are optional after closing curly braces (which close conditionals and loops as discussed below).

.. code-block:: none
   :emphasize-lines: 1,1

    $ echo x=1,y=2 | mlr put 'while (NF < 10) { $[NF+1] = ""}  $foo = "bar"'
    x=1,y=2,3=,4=,5=,6=,7=,8=,9=,10=,foo=bar

.. code-block:: none
   :emphasize-lines: 1,1

    $ echo x=1,y=2 | mlr put 'while (NF < 10) { $[NF+1] = ""}; $foo = "bar"'
    x=1,y=2,3=,4=,5=,6=,7=,8=,9=,10=,foo=bar

Semicolons are required between statements even if those statements are on separate lines.  **Newlines** are for your convenience but have no syntactic meaning: line endings do not terminate statements. For example, adjacent assignment statements must be separated by semicolons even if those statements are on separate lines:

.. code-block:: none

    mlr put '
      $x = 1
      $y = 2 # Syntax error
    '
    
    mlr put '
      $x = 1;
      $y = 2 # This is OK
    '

**Trailing commas** are allowed in function/subroutine definitions, function/subroutine callsites, and map literals. This is intended for (although not restricted to) the multi-line case:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --csvlite --from data/a.csv put '
      func f(
        num a,
        num b,
      ): num {
        return a**2 + b**2;
      }
      $* = {
        "s": $a + $b,
        "t": $a - $b,
        "u": f(
          $a,
          $b,
        ),
        "v": NR,
      }
    '
    s,t,u,v
    3,-1,5,1
    9,-1,41,2

Bodies for all compound statements must be enclosed in **curly braces**, even if the body is a single statement:

.. code-block:: none
   :emphasize-lines: 1,1

    mlr put 'if ($x == 1) $y = 2' # Syntax error

.. code-block:: none
   :emphasize-lines: 1,1

    mlr put 'if ($x == 1) { $y = 2 }' # This is OK

Bodies for compound statements may be empty:

.. code-block:: none
   :emphasize-lines: 1,1

    mlr put 'if ($x == 1) { }' # This no-op is syntactically acceptable

Variables
----------------------------------------------------------------

Miller has the following kinds of variables:

**Built-in variables** such as ``NF``, ``NF``, ``FILENAME``, ``M_PI``, and ``M_E``.  These are all capital letters and are read-only (although some of them change value from one record to another).

**Fields of stream records**, accessed using the ``$`` prefix. These refer to fields of the current data-stream record. For example, in ``echo x=1,y=2 | mlr put '$z = $x + $y'``, ``$x`` and ``$y`` refer to input fields, and ``$z`` refers to a new, computed output field. In a few contexts, presented below, you can refer to the entire record as ``$*``.

**Out-of-stream variables** accessed using the ``@`` prefix. These refer to data which persist from one record to the next, including in ``begin`` and ``end`` blocks (which execute before/after the record stream is consumed, respectively). You use them to remember values across records, such as sums, differences, counters, and so on.  In a few contexts, presented below, you can refer to the entire out-of-stream-variables collection as ``@*``.

**Local variables** are limited in scope and extent to the current statements being executed: these include function arguments, bound variables in for loops, and explicitly declared local variables.

**Keywords** are not variables, but since their names are reserved, you cannot use these names for local variables.

Built-in variables
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

These are written all in capital letters, such as ``NR``, ``NF``, ``FILENAME``, and only a small, specific set of them is defined by Miller.

Namely, Miller supports the following five built-in variables for :doc:`filter and put <reference-dsl>`, all ``awk``-inspired: ``NF``, ``NR``, ``FNR``, ``FILENUM``, and ``FILENAME``, as well as the mathematical constants ``M_PI`` and ``M_E``.  Lastly, the ``ENV`` hashmap allows read access to environment variables, e.g.  ``ENV["HOME"]`` or ``ENV["foo_".$hostname]``.

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr filter 'FNR == 2' data/small*
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    1=pan,2=pan,3=1,4=0.3467901443380824,5=0.7268028627434533
    a=wye,b=eks,i=10000,x=0.734806020620654365,y=0.884788571337605134

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put '$fnr = FNR' data/small*
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533,fnr=1
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797,fnr=2
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776,fnr=3
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463,fnr=4
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,fnr=5
    1=a,2=b,3=i,4=x,5=y,fnr=1
    1=pan,2=pan,3=1,4=0.3467901443380824,5=0.7268028627434533,fnr=2
    1=eks,2=pan,3=2,4=0.7586799647899636,5=0.5221511083334797,fnr=3
    1=wye,2=wye,3=3,4=0.20460330576630303,5=0.33831852551664776,fnr=4
    1=eks,2=wye,3=4,4=0.38139939387114097,5=0.13418874328430463,fnr=5
    1=wye,2=pan,3=5,4=0.5732889198020006,5=0.8636244699032729,fnr=6
    a=pan,b=eks,i=9999,x=0.267481232652199086,y=0.557077185510228001,fnr=1
    a=wye,b=eks,i=10000,x=0.734806020620654365,y=0.884788571337605134,fnr=2
    a=pan,b=wye,i=10001,x=0.870530722602517626,y=0.009854780514656930,fnr=3
    a=hat,b=wye,i=10002,x=0.321507044286237609,y=0.568893318795083758,fnr=4
    a=pan,b=zee,i=10003,x=0.272054845593895200,y=0.425789896597056627,fnr=5

Their values of ``NF``, ``NR``, ``FNR``, ``FILENUM``, and ``FILENAME`` change from one record to the next as Miller scans through your input data stream. The mathematical constants, of course, do not change; ``ENV`` is populated from the system environment variables at the time Miller starts and is read-only for the remainder of program execution.

Their **scope is global**: you can refer to them in any ``filter`` or ``put`` statement. Their values are assigned by the input-record reader:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --csv put '$nr = NR' data/a.csv
    a,b,c,nr
    1,2,3,1
    4,5,6,2

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --csv repeat -n 3 then put '$nr = NR' data/a.csv
    a,b,c,nr
    1,2,3,1
    1,2,3,1
    1,2,3,1
    4,5,6,2
    4,5,6,2
    4,5,6,2

The **extent** is for the duration of the put/filter: in a ``begin`` statement (which executes before the first input record is consumed) you will find ``NR=1`` and in an ``end`` statement (which is executed after the last input record is consumed) you will find ``NR`` to be the total number of records ingested.

These are all **read-only** for the ``mlr put`` and ``mlr filter`` DSLs: they may be assigned from, e.g. ``$nr=NR``, but they may not be assigned to: ``NR=100`` is a syntax error.

Field names
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Names of fields within stream records must be specified using a ``$`` in :doc:`filter and put expressions <reference-dsl>`, even though the dollar signs don't appear in the data stream itself. For integer-indexed data, this looks like ``awk``'s ``$1,$2,$3``, except that Miller allows non-numeric names such as ``$quantity`` or ``$hostname``.  Likewise, enclose string literals in double quotes in ``filter`` expressions even though they don't appear in file data.  In particular, ``mlr filter '$x=="abc"'`` passes through the record ``x=abc``.

If field names have **special characters** such as ``.`` then you can use braces, e.g. ``'${field.name}'``.

You may also use a **computed field name** in square brackets, e.g.

.. code-block:: none
   :emphasize-lines: 1,1

    $ echo a=3,b=4 | mlr filter '$["x"] < 0.5'

.. code-block:: none
   :emphasize-lines: 1,1

    $ echo s=green,t=blue,a=3,b=4 | mlr put '$[$s."_".$t] = $a * $b'
    s=green,t=blue,a=3,b=4,green_blue=12

Notes:

The names of record fields depend on the contents of your input data stream, and their values change from one record to the next as Miller scans through your input data stream.

Their **extent** is limited to the current record; their **scope** is the ``filter`` or ``put`` command in which they appear.

These are **read-write**: you can do ``$y=2*$x``, ``$x=$x+1``, etc.

Records are Miller's output: field names present in the input stream are passed through to output (written to standard output) unless fields are removed with ``cut``, or records are excluded with ``filter`` or ``put -q``, etc. Simply assign a value to a field and it will be output.

Positional field names
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Even though Miller's main selling point is name-indexing, sometimes you really want to refer to a field name by its positional index (starting from 1).

Use ``$[[3]]`` to access the name of field 3.  More generally, any expression evaluating to an integer can go between ``$[[`` and ``]]``.

Then using a computed field name, ``$[ $[[3]] ]`` is the value in the third field. This has the shorter equivalent notation ``$[[[3]]]``.

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr cat data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put '$[[3]] = "NEW"' data/small
    a=pan,b=pan,NEW=1,x=0.3467901443380824,y=0.7268028627434533
    a=eks,b=pan,NEW=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,NEW=3,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=wye,NEW=4,x=0.38139939387114097,y=0.13418874328430463
    a=wye,b=pan,NEW=5,x=0.5732889198020006,y=0.8636244699032729

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put '$[[[3]]] = "NEW"' data/small
    a=pan,b=pan,i=NEW,x=0.3467901443380824,y=0.7268028627434533
    a=eks,b=pan,i=NEW,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=NEW,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=wye,i=NEW,x=0.38139939387114097,y=0.13418874328430463
    a=wye,b=pan,i=NEW,x=0.5732889198020006,y=0.8636244699032729

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put '$NEW = $[[NR]]' data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533,NEW=a
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797,NEW=b
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776,NEW=i
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463,NEW=x
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,NEW=y

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put '$NEW = $[[[NR]]]' data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533,NEW=pan
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797,NEW=pan
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776,NEW=3
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463,NEW=0.38139939387114097
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,NEW=0.8636244699032729

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put '$[[[NR]]] = "NEW"' data/small
    a=NEW,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
    a=eks,b=NEW,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=NEW,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=wye,i=4,x=NEW,y=0.13418874328430463
    a=wye,b=pan,i=5,x=0.5732889198020006,y=NEW

Right-hand side accesses to non-existent fields -- i.e. with index less than 1 or greater than ``NF`` -- return an absent value. Likewise, left-hand side accesses only refer to fields which already exist. For example, if a field has 5 records then assigning the name or value of the 6th (or 600th) field results in a no-op.

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put '$[[6]] = "NEW"' data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put '$[[[6]]] = "NEW"' data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729

Out-of-stream variables
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

These are prefixed with an at-sign, e.g. ``@sum``.  Furthermore, unlike built-in variables and stream-record fields, they are maintained in an arbitrarily nested hashmap: you can do ``@sum += $quanity``, or ``@sum[$color] += $quanity``, or ``@sum[$color][$shape] += $quanity``. The keys for the multi-level hashmap can be any expression which evaluates to string or integer: e.g.  ``@sum[NR] = $a + $b``, ``@sum[$a."-".$b] = $x``, etc.

Their names and their values are entirely under your control; they change only when you assign to them.

Just as for field names in stream records, if you want to define out-of-stream variables with **special characters** such as ``.`` then you can use braces, e.g. ``'@{variable.name}["index"]'``.

You may use a **computed key** in square brackets, e.g.

.. code-block:: none
   :emphasize-lines: 1,1

    $ echo s=green,t=blue,a=3,b=4 | mlr put -q '@[$s."_".$t] = $a * $b; emit all'
    green_blue=12

Out-of-stream variables are **scoped** to the ``put`` command in which they appear.  In particular, if you have two or more ``put`` commands separated by ``then``, each put will have its own set of out-of-stream variables:

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat data/a.dkvp
    a=1,b=2,c=3
    a=4,b=5,c=6

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put '@sum += $a; end {emit @sum}' then put 'is_present($a) {$a=10*$a; @sum += $a}; end {emit @sum}' data/a.dkvp
    a=10,b=2,c=3
    a=40,b=5,c=6
    sum=5
    sum=50

Out-of-stream variables' **extent** is from the start to the end of the record stream, i.e. every time the ``put`` or ``filter`` statement referring to them is executed.

Out-of-stream variables are **read-write**: you can do ``$sum=@sum``, ``@sum=$sum``, etc.

Indexed out-of-stream variables
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Using an index on the ``@count`` and ``@sum`` variables, we get the benefit of the ``-g`` (group-by) option which ``mlr stats1`` and various other Miller commands have:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put -q '
      @x_count[$a] += 1;
      @x_sum[$a] += $x;
      end {
        emit @x_count, "a";
        emit @x_sum, "a";
      }
    ' ../data/small
    a=pan,x_count=2
    a=eks,x_count=3
    a=wye,x_count=2
    a=zee,x_count=2
    a=hat,x_count=1
    a=pan,x_sum=0.8494161498792961
    a=eks,x_sum=1.75186341922895
    a=wye,x_sum=0.7778922255683036
    a=zee,x_sum=1.1256801691982772
    a=hat,x_sum=0.03144187646093577

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr stats1 -a count,sum -f x -g a ../data/small
    a=pan,x_count=2,x_sum=0.8494161498792961
    a=eks,x_count=3,x_sum=1.75186341922895
    a=wye,x_count=2,x_sum=0.7778922255683036
    a=zee,x_count=2,x_sum=1.1256801691982772
    a=hat,x_count=1,x_sum=0.03144187646093577

Indices can be arbitrarily deep -- here there are two or more of them:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --from data/medium put -q '
      @x_count[$a][$b] += 1;
      @x_sum[$a][$b] += $x;
      end {
        emit (@x_count, @x_sum), "a", "b";
      }
    '
    a=pan,b=pan,x_count=427,x_sum=219.1851288316854
    a=pan,b=wye,x_count=395,x_sum=198.43293070748447
    a=pan,b=eks,x_count=429,x_sum=216.07522773165525
    a=pan,b=hat,x_count=417,x_sum=205.22277621488686
    a=pan,b=zee,x_count=413,x_sum=205.09751802331917
    a=eks,b=pan,x_count=371,x_sum=179.96303047250723
    a=eks,b=wye,x_count=407,x_sum=196.9452860713734
    a=eks,b=zee,x_count=357,x_sum=176.8803651584733
    a=eks,b=eks,x_count=413,x_sum=215.91609712937984
    a=eks,b=hat,x_count=417,x_sum=208.783170520597
    a=wye,b=wye,x_count=377,x_sum=185.29584980261419
    a=wye,b=pan,x_count=392,x_sum=195.84790012056564
    a=wye,b=hat,x_count=426,x_sum=212.0331829346132
    a=wye,b=zee,x_count=385,x_sum=194.77404756708714
    a=wye,b=eks,x_count=386,x_sum=204.8129608356315
    a=zee,b=pan,x_count=389,x_sum=202.21380378504267
    a=zee,b=wye,x_count=455,x_sum=233.9913939194868
    a=zee,b=eks,x_count=391,x_sum=190.9617780631925
    a=zee,b=zee,x_count=403,x_sum=206.64063510417319
    a=zee,b=hat,x_count=409,x_sum=191.30000620900935
    a=hat,b=wye,x_count=423,x_sum=208.8830097609959
    a=hat,b=zee,x_count=385,x_sum=196.3494502965293
    a=hat,b=eks,x_count=389,x_sum=189.0067933716193
    a=hat,b=hat,x_count=381,x_sum=182.8535323148762
    a=hat,b=pan,x_count=363,x_sum=168.5538067327806

The idea is that ``stats1``, and other Miller verbs, encapsulate frequently-used patterns with a minimum of keystroking (and run a little faster), whereas using out-of-stream variables you have more flexibility and control in what you do.

Begin/end blocks can be mixed with pattern/action blocks. For example:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put '
      begin {
        @num_total = 0;
        @num_positive = 0;
      };
      @num_total += 1;
      $x > 0.0 {
        @num_positive += 1;
        $y = log10($x); $z = sqrt($y)
      };
      end {
        emitf @num_total, @num_positive
      }
    ' data/put-gating-example-1.dkvp
    x=-1
    x=0
    x=1,y=0,z=0
    x=2,y=0.3010299956639812,z=0.5486620049392715
    x=3,y=0.4771212547196624,z=0.6907396432228734
    num_total=5,num_positive=3

.. _reference-dsl-local-variables:

Local variables
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Local variables are similar to out-of-stream variables, except that their extent is limited to the expressions in which they appear (and their basenames can't be computed using square brackets). There are three kinds of local variables: **arguments** to functions/subroutines, **variables bound within for-loops**, and **locals** defined within control blocks. They may be untyped using ``var``, or typed using ``num``, ``int``, ``float``, ``str``, ``bool``, and ``map``.

For example:

.. code-block:: none
   :emphasize-lines: 1,1

    $ # Here I'm using a specified random-number seed so this example always
    # produces the same output for this web document: in everyday practice we
    # would leave off the --seed 12345 part.
    mlr --seed 12345 seqgen --start 1 --stop 10 then put '
      func f(a, b) {                          # function arguments a and b
          r = 0.0;                            # local r scoped to the function
          for (int i = 0; i < 6; i += 1) {    # local i scoped to the for-loop
              num u = urand();                # local u scoped to the for-loop
              r += u;                         # updates r from the enclosing scope
          }
          r /= 6;
          return a + (b - a) * r;
      }
      num o = f(10, 20);                      # local to the top-level scope
      $o = o;
    '
    i=1,o=15.952526011537227
    i=2,o=12.782237754999116
    i=3,o=15.126606630220966
    i=4,o=14.794357488895775
    i=5,o=15.168665974047421
    i=6,o=16.20662783079942
    i=7,o=13.966128063060479
    i=8,o=13.99248245928659
    i=9,o=15.784270485515197
    i=10,o=15.37686787628025

Things which are completely unsurprising, resembling many other languages:

* Parameter names are bound to their arguments but can be reassigned, e.g. if there is a parameter named ``a`` then you can reassign the value of ``a`` to be something else within the function if you like.

* However, you cannot redeclare the *type* of an argument or a local: ``var a=1; var a=2`` is an error but ``var a=1;  a=2`` is OK.

* All argument-passing is positional rather than by name; arguments are passed by value, not by reference. (This is also true for map-valued variables: they are not, and cannot be, passed by reference)

* You can define locals (using ``var``, ``num``, etc.) at any scope (if-statements, else-statements, while-loops, for-loops, or the top-level scope), and nested scopes will have access (more details on scope in the next section).  If you define a local variable with the same name inside an inner scope, then a new variable is created with the narrower scope.

* If you assign to a local variable for the first time in a scope without declaring it as ``var``, ``num``, etc. then: if it exists in an outer scope, that outer-scope variable will be updated; if not, it will be defined in the current scope as if ``var`` had been used. (See also :ref:`reference-dsl-type-checking` for an example.) I recommend always declaring variables explicitly to make the intended scoping clear.

* Functions and subroutines never have access to locals from their callee (unless passed by value as arguments).

Things which are perhaps surprising compared to other languages:

* Type declarations using ``var``, or typed using ``num``, ``int``, ``float``, ``str``, and ``bool`` are necessary to declare local variables.  Function arguments and variables bound in for-loops over stream records and out-of-stream variables are *implicitly* declared using ``var``. (Some examples are shown below.)

* Type-checking is done at assignment time. For example, ``float f = 0`` is an error (since ``0`` is an integer), as is ``float f = 0.0; f = 1``. For this reason I prefer to use ``num`` over ``float`` in most contexts since ``num`` encompasses integer and floating-point values. More information about type-checking is at :ref:`reference-dsl-type-checking`.

* Bound variables in for-loops over stream records and out-of-stream variables are implicitly local to that block. E.g. in ``for (k, v in $*) { ... }`` ``for ((k1, k2), v in @*) { ... }`` if there are ``k``, ``v``, etc. in the enclosing scope then those will be masked by the loop-local bound variables in the loop, and moreover the values of the loop-local bound variables are not available after the end of the loop.

* For C-style triple-for loops, if a for-loop variable is defined using ``var``, ``int``, etc. then it is scoped to that for-loop. E.g. ``for (i = 0; i < 10; i += 1) { ... }`` and ``for (int i = 0; i < 10; i += 1) { ... }``. (This is unsurprising.). If there is no typedecl and an outer-scope variable of that name exists, then it is used. (This is also unsurprising.) But of there is no outer-scope variable of that name then the variable is scoped to the for-loop only.

The following example demonstrates the scope rules:

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat data/scope-example.mlr
    func f(a) {      # argument is local to the function
      var b = 100;   # local to the function
      c = 100;       # local to the function; does not overwrite outer c
      return a + 1;
    }
    var a = 10;      # local at top level
    var b = 20;      # local at top level
    c = 30;          # local at top level; there is no more-outer-scope c
    if (NR == 3) {
      var a = 40;    # scoped to the if-statement; doesn't overwrite outer a
      b = 50;        # not scoped to the if-statement; overwrites outer b
      c = 60;        # not scoped to the if-statement; overwrites outer c
      d = 70;        # there is no outer d so a local d is created here
    
      $inner_a = a;
      $inner_b = b;
      $inner_c = c;
      $inner_d = d;
    }
    $outer_a = a;
    $outer_b = b;
    $outer_c = c;
    $outer_d = d;    # there is no outer d defined so no assignment happens

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat data/scope-example.dat
    n=1,x=123
    n=2,x=456
    n=3,x=789

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --oxtab --from data/scope-example.dat put -f data/scope-example.mlr
    n       1
    x       123
    outer_a 10
    outer_b 20
    outer_c 30
    
    n       2
    x       456
    outer_a 10
    outer_b 20
    outer_c 30
    
    n       3
    x       789
    inner_a 40
    inner_b 50
    inner_c 60
    inner_d 70
    outer_a 10
    outer_b 50
    outer_c 60

And this example demonstrates the type-declaration rules:

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat data/type-decl-example.mlr
    subr s(a, str b, int c) {                         # a is implicitly var (untyped).
                                                      # b is explicitly str.
                                                      # c is explicitly int.
                                                      # The type-checking is done at the callsite
                                                      # when arguments are bound to parameters.
                                                      #
        var b = 100;     # error                      # Re-declaration in the same scope is disallowed.
        int n = 10;                                   # Declaration of variable local to the subroutine.
        n = 20;                                       # Assignment is OK.
        int n = 30;      # error                      # Re-declaration in the same scope is disallowed.
        str n = "abc";   # error                      # Re-declaration in the same scope is disallowed.
                                                      #
        float f1 = 1;    # error                      # 1 is an int, not a float.
        float f2 = 2.0;                               # 2.0 is a float.
        num f3 = 3;                                   # 3 is a num.
        num f4 = 4.0;                                 # 4.0 is a num.
    }                                                 #
                                                      #
    call s(1, 2, 3);                                  # Type-assertion '3 is int' is done here at the callsite.
                                                      #
    k = "def";                                        # Top-level variable k.
                                                      #
    for (str k, v in $*) {                            # k and v are bound here, masking outer k.
      print k . ":" . v;                              # k is explicitly str; v is implicitly var.
    }                                                 #
                                                      #
    print "k is".k;                                   # k at this scope level is still "def".
    print "v is".v;                                   # v is undefined in this scope.
                                                      #
    i = -1;                                           #
    for (i = 1, int j = 2; i <= 10; i += 1, j *= 2) { # C-style triple-for variables use enclosing scope, unless
                                                      # declared local: i is outer, j is local to the loop.
      print "inner i =" . i;                          #
      print "inner j =" . j;                          #
    }                                                 #
    print "outer i =" . i;                            # i has been modified by the loop.
    print "outer j =" . j;                            # j is undefined in this scope.

Map literals
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Miller's ``put``/``filter`` DSL has four kinds of hashmaps. **Stream records** are (single-level) maps from name to value. **Out-of-stream variables** and **local variables** can also be maps, although they can be multi-level hashmaps (e.g. ``@sum[$x][$y]``).  The fourth kind is **map literals**. These cannot be on the left-hand side of assignment expressions. Syntactically they look like JSON, although Miller allows string and integer keys in its map literals while JSON allows only string keys (e.g. ``"3"`` rather than ``3``).

For example, the following swaps the input stream's ``a`` and ``i`` fields, modifies ``y``, and drops the rest:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint put '
      $* = {
        "a": $i,
        "i": $a,
        "y": $y * 10,
      }
    ' data/small
    a i   y
    1 pan 7.268028627434533
    2 eks 5.221511083334796
    3 wye 3.3831852551664774
    4 eks 1.3418874328430463
    5 wye 8.63624469903273

Likewise, you can assign map literals to out-of-stream variables or local variables; pass them as arguments to user-defined functions, return them from functions, and so on:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --from data/small put '
      func f(map m): map {
        m["x"] *= 200;
        return m;
      }
      $* = f({"a": $a, "x": $x});
    '
    a=pan,x=69.35802886761648
    a=eks,x=151.73599295799272
    a=wye,x=40.92066115326061
    a=eks,x=76.2798787742282
    a=wye,x=114.65778396040011

Like out-of-stream and local variables, map literals can be multi-level:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --from data/small put -q '
      begin {
        @o = {
          "nrec": 0,
          "nkey": {"numeric":0, "non-numeric":0},
        };
      }
      @o["nrec"] += 1;
      for (k, v in $*) {
        if (is_numeric(v)) {
          @o["nkey"]["numeric"] += 1;
        } else {
          @o["nkey"]["non-numeric"] += 1;
        }
      }
      end {
        dump @o;
      }
    '
    {
      "nrec": 5,
      "nkey": {
        "numeric": 15,
        "non-numeric": 10
      }
    }

By default, map-valued expressions are dumped using JSON formatting. If you use ``dump`` to print a hashmap with integer keys and you don't want them double-quoted (JSON-style) then you can use ``mlr put --jknquoteint``. See also ``mlr put --help``.

.. _reference-dsl-type-checking:

Type-checking
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Miller's ``put``/``filter`` DSLs support two optional kinds of type-checking.  One is inline **type-tests** and **type-assertions** within expressions.  The other is **type declarations** for assignments to local variables, binding of arguments to user-defined functions, and return values from user-defined functions, These are discussed in the following subsections.

Use of type-checking is entirely up to you: omit it if you want flexibility with heterogeneous data; use it if you want to help catch misspellings in your DSL code or unexpected irregularities in your input data.

.. _reference-dsl-type-tests-and-assertions:

Type-test and type-assertion expressions
................................................................

The following ``is...`` functions take a value and return a boolean indicating whether the argument is of the indicated type. The ``assert_...`` functions return their argument if it is of the specified type, and cause a fatal error otherwise:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr -F | grep ^is
    is_absent
    is_array
    is_bool
    is_boolean
    is_empty
    is_empty_map
    is_error
    is_float
    is_int
    is_map
    is_nonempty_map
    is_not_empty
    is_not_map
    is_not_array
    is_not_null
    is_null
    is_numeric
    is_present
    is_string

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr -F | grep ^assert
    asserting_absent
    asserting_array
    asserting_bool
    asserting_boolean
    asserting_error
    asserting_empty
    asserting_empty_map
    asserting_float
    asserting_int
    asserting_map
    asserting_nonempty_map
    asserting_not_empty
    asserting_not_map
    asserting_not_array
    asserting_not_null
    asserting_null
    asserting_numeric
    asserting_present
    asserting_string

Please see :ref:`cookbook-data-cleaning-examples` for examples of how to use these.

Type-declarations for local variables, function parameter, and function return values
...............................................................................................

Local variables can be defined either untyped as in ``x = 1``, or typed as in ``int x = 1``. Types include **var** (explicitly untyped), **int**, **float**, **num** (int or float), **str**, **bool**, and **map**. These optional type declarations are enforced at the time values are assigned to variables: whether at the initial value assignment as in ``int x = 1`` or in any subsequent assignments to the same variable farther down in the scope.

The reason for ``num`` is that ``int`` and ``float`` typedecls are very precise:

.. code-block:: none

    float a = 0;   # Runtime error since 0 is int not float
    int   b = 1.0; # Runtime error since 1.0 is float not int
    num   c = 0;   # OK
    num   d = 1.0; # OK

A suggestion is to use ``num`` for general use when you want numeric content, and use ``int`` when you genuinely want integer-only values, e.g. in loop indices or map keys (since Miller map keys can only be strings or ints).

The ``var`` type declaration indicates no type restrictions, e.g. ``var x = 1`` has the same type restrictions on ``x`` as ``x = 1``. The difference is in intentional shadowing: if you have ``x = 1`` in outer scope and ``x = 2`` in inner scope (e.g. within a for-loop or an if-statement) then outer-scope ``x`` has value 2 after the second assignment.  But if you have ``var x = 2`` in the inner scope, then you are declaring a variable scoped to the inner block.) For example:

.. code-block:: none

    x = 1;
    if (NR == 4) {
      x = 2; # Refers to outer-scope x: value changes from 1 to 2.
    }
    print x; # Value of x is now two

.. code-block:: none

    x = 1;
    if (NR == 4) {
      var x = 2; # Defines a new inner-scope x with value 2
    }
    print x;     # Value of this x is still 1

Likewise function arguments can optionally be typed, with type enforced when the function is called:

.. code-block:: none

    func f(map m, int i) {
      ...
    }
    $a = f({1:2, 3:4}, 5);     # OK
    $b = f({1:2, 3:4}, "abc"); # Runtime error
    $c = f({1:2, 3:4}, $x);    # Runtime error for records with non-integer field named x
    if (NR == 4) {
      var x = 2; # Defines a new inner-scope x with value 2
    }
    print x;     # Value of this x is still 1

Thirdly, function return values can be type-checked at the point of ``return`` using ``:`` and a typedecl after the parameter list:

.. code-block:: none

    func f(map m, int i): bool {
      ...
      ...
      if (...) {
        return "false"; # Runtime error if this branch is taken
      }
      ...
      ...
      if (...) {
        return retval; # Runtime error if this function doesn't have an in-scope
        # boolean-valued variable named retval
      }
      ...
      ...
      # In Miller if your functions don't explicitly return a value, they return absent-null.
      # So it would also be a runtime error on reaching the end of this function without
      # an explicit return statement.
    }

Null data: empty and absent
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Please see :ref:`reference-null-data`.

Aggregate variable assignments
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

There are three remaining kinds of variable assignment using out-of-stream variables, the last two of which use the ``$*`` syntax:

* Recursive copy of out-of-stream variables
* Out-of-stream variable assigned to full stream record
* Full stream record assigned to an out-of-stream variable

Example recursive copy of out-of-stream variables:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint put -q '@v["sum"] += $x; @v["count"] += 1; end{dump; @w = @v; dump}' data/small
    {
      "v": {
        "sum": 2.264761728567491,
        "count": 5
      }
    }
    {
      "v": {
        "sum": 2.264761728567491,
        "count": 5
      },
      "w": {
        "sum": 2.264761728567491,
        "count": 5
      }
    }

Example of out-of-stream variable assigned to full stream record, where the 2nd record is stashed, and the 4th record is overwritten with that:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put 'NR == 2 {@keep = $*}; NR == 4 {$* = @keep}' data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729

Example of full stream record assigned to an out-of-stream variable, finding the record for which the ``x`` field has the largest value in the input stream:

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

    $ mlr --opprint put -q 'is_null(@xmax) || $x > @xmax {@xmax=$x; @recmax=$*}; end {emit @recmax}' data/small
    a   b   i x                  y
    eks pan 2 0.7586799647899636 0.5221511083334797

Keywords for filter and put
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --help-all-keywords
    TOD: port mlr_dsl_list_all_keywords

Operator precedence
----------------------------------------------------------------

Operators are listed in order of decreasing precedence, highest first.

.. code-block:: none

    Operators              Associativity
    ---------              -------------
    ()                     left to right
    **                    right to left
    ! ~ unary+ unary- &    right to left
    binary* / // %         left to right
    binary+ binary- .      left to right
    << >>                  left to right
    &                      left to right
    ^                      left to right
    |                      left to right
    < <= > >=              left to right
    == != =~ !=~           left to right
    &&                     left to right
    ^^                     left to right
    ||                     left to right
    ? :                    right to left
    =                      N/A for Miller (there is no $a=$b=$c)

Operator and function semantics
----------------------------------------------------------------

* Functions are often pass-throughs straight to the system-standard Go libraries.

* The ``min`` and ``max`` functions are different from other multi-argument functions which return null if any of their inputs are null: for ``min`` and ``max``, by contrast, if one argument is absent-null, the other is returned. Empty-null loses min or max against numeric or boolean; empty-null is less than any other string.

* Symmetrically with respect to the bitwise OR, XOR, and AND operators ``|``, ``^``, ``&``, Miller has logical operators ``||``, ``^^``, ``&&``: the logical XOR not existing in Go.

* The exponentiation operator ``**`` is familiar from many languages.

* The regex-match and regex-not-match operators ``=~`` and ``!=~`` are similar to those in Ruby and Perl.

Control structures
----------------------------------------------------------------

Pattern-action blocks
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

These are reminiscent of ``awk`` syntax.  They can be used to allow assignments to be done only when appropriate -- e.g. for math-function domain restrictions, regex-matching, and so on:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr cat data/put-gating-example-1.dkvp
    x=-1
    x=0
    x=1
    x=2
    x=3

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put '$x > 0.0 { $y = log10($x); $z = sqrt($y) }' data/put-gating-example-1.dkvp
    x=-1
    x=0
    x=1,y=0,z=0
    x=2,y=0.3010299956639812,z=0.5486620049392715
    x=3,y=0.4771212547196624,z=0.6907396432228734

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr cat data/put-gating-example-2.dkvp
    a=abc_123
    a=some other name
    a=xyz_789

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put '$a =~ "([a-z]+)_([0-9]+)" { $b = "left_\1"; $c = "right_\2" }' data/put-gating-example-2.dkvp
    a=abc_123,b=left_\1,c=right_\2
    a=some other name
    a=xyz_789,b=left_\1,c=right_\2

This produces heteregenous output which Miller, of course, has no problems with (see :doc:`record-heterogeneity`).  But if you want homogeneous output, the curly braces can be replaced with a semicolon between the expression and the body statements.  This causes ``put`` to evaluate the boolean expression (along with any side effects, namely, regex-captures ``\1``, ``\2``, etc.) but doesn't use it as a criterion for whether subsequent assignments should be executed. Instead, subsequent assignments are done unconditionally:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put '$x > 0.0; $y = log10($x); $z = sqrt($y)' data/put-gating-example-1.dkvp
    x=1,y=0,z=0
    x=2,y=0.3010299956639812,z=0.5486620049392715
    x=3,y=0.4771212547196624,z=0.6907396432228734

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put '$a =~ "([a-z]+)_([0-9]+)"; $b = "left_\1"; $c = "right_\2"' data/put-gating-example-2.dkvp
    a=abc_123,b=left_\1,c=right_\2
    a=xyz_789,b=left_\1,c=right_\2

If-statements
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

These are again reminiscent of ``awk``. Pattern-action blocks are a special case of ``if`` with no ``elif`` or ``else`` blocks, no ``if`` keyword, and parentheses optional around the boolean expression:

.. code-block:: none

    mlr put 'NR == 4 {$foo = "bar"}'

.. code-block:: none

    mlr put 'if (NR == 4) {$foo = "bar"}'

Compound statements use ``elif`` (rather than ``elsif`` or ``else if``):

.. code-block:: none

    mlr put '
      if (NR == 2) {
        ...
      } elif (NR ==4) {
        ...
      } elif (NR ==6) {
        ...
      } else {
        ...
      }
    '

While and do-while loops
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Miller's ``while`` and ``do-while`` are unsurprising in comparison to various languages, as are ``break`` and ``continue``:

.. code-block:: none
   :emphasize-lines: 1,1

    $ echo x=1,y=2 | mlr put '
      while (NF < 10) {
        $[NF+1] = ""
      }
      $foo = "bar"
    '
    x=1,y=2,3=,4=,5=,6=,7=,8=,9=,10=,foo=bar

.. code-block:: none
   :emphasize-lines: 1,1

    $ echo x=1,y=2 | mlr put '
      do {
        $[NF+1] = "";
        if (NF == 5) {
          break
        }
      } while (NF < 10);
      $foo = "bar"
    '
    x=1,y=2,3=,4=,5=,foo=bar

A ``break`` or ``continue`` within nested conditional blocks or if-statements will, of course, propagate to the innermost loop enclosing them, if any. A ``break`` or ``continue`` outside a loop is a syntax error that will be flagged as soon as the expression is parsed, before any input records are ingested.
The existence of ``while``, ``do-while``, and ``for`` loops in Miller's DSL means that you can create infinite-loop scenarios inadvertently.  In particular, please recall that DSL statements are executed once if in ``begin`` or ``end`` blocks, and once *per record* otherwise. For example, **while (NR < 10) will never terminate as NR is only incremented between records**.

For-loops
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

While Miller's ``while`` and ``do-while`` statements are much as in many other languages, ``for`` loops are more idiosyncratic to Miller. They are loops over key-value pairs, whether in stream records, out-of-stream variables, local variables, or map-literals: more reminiscent of ``foreach``, as in (for example) PHP. There are **for-loops over map keys** and **for-loops over key-value tuples**.  Additionally, Miller has a **C-style triple-for loop** with initialize, test, and update statements.

As with ``while`` and ``do-while``, a ``break`` or ``continue`` within nested control structures will propagate to the innermost loop enclosing them, if any, and a ``break`` or ``continue`` outside a loop is a syntax error that will be flagged as soon as the expression is parsed, before any input records are ingested.

Key-only for-loops
................................................................

The ``key`` variable is always bound to the *key* of key-value pairs:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --from data/small put '
      print "NR = ".NR;
      for (key in $*) {
        value = $[key];
        print "  key:" . key . "  value:".value;
      }
    
    '
    NR = 1
      key:a  value:pan
      key:b  value:pan
      key:i  value:1
      key:x  value:0.3467901443380824
      key:y  value:0.7268028627434533
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
    NR = 2
      key:a  value:eks
      key:b  value:pan
      key:i  value:2
      key:x  value:0.7586799647899636
      key:y  value:0.5221511083334797
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    NR = 3
      key:a  value:wye
      key:b  value:wye
      key:i  value:3
      key:x  value:0.20460330576630303
      key:y  value:0.33831852551664776
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    NR = 4
      key:a  value:eks
      key:b  value:wye
      key:i  value:4
      key:x  value:0.38139939387114097
      key:y  value:0.13418874328430463
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
    NR = 5
      key:a  value:wye
      key:b  value:pan
      key:i  value:5
      key:x  value:0.5732889198020006
      key:y  value:0.8636244699032729
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr -n put '
      end {
        o = {1:2, 3:{4:5}};
        for (key in o) {
          print "  key:" . key . "  valuetype:" . typeof(o[key]);
        }
      }
    '
      key:1  valuetype:int
      key:3  valuetype:map

Note that the value corresponding to a given key may be gotten as through a **computed field name** using square brackets as in ``$[key]`` for stream records, or by indexing the looped-over variable using square brackets.

Key-value for-loops
................................................................

Single-level keys may be gotten at using either ``for(k,v)`` or ``for((k),v)``; multi-level keys may be gotten at using ``for((k1,k2,k3),v)`` and so on.  The ``v`` variable will be bound to to a scalar value (a string or a number) if the map stops at that level, or to a map-valued variable if the map goes deeper. If the map isn't deep enough then the loop body won't be executed.

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat data/for-srec-example.tbl
    label1 label2 f1  f2  f3
    blue   green  100 240 350
    red    green  120 11  195
    yellow blue   140 0   240

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --pprint --from data/for-srec-example.tbl put '
      $sum1 = $f1 + $f2 + $f3;
      $sum2 = 0;
      $sum3 = 0;
      for (key, value in $*) {
        if (key =~ "^f[0-9]+") {
          $sum2 += value;
          $sum3 += $[key];
        }
      }
    '
    label1 label2 f1  f2  f3  sum1 sum2 sum3
    blue   green  100 240 350 690  690  690
    red    green  120 11  195 326  326  326
    yellow blue   140 0   240 380  380  380

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --from data/small --opprint put 'for (k,v in $*) { $[k."_type"] = typeof(v) }'
    a   b   i x                   y                   a_type b_type i_type x_type y_type
    pan pan 1 0.3467901443380824  0.7268028627434533  string string int    float  float
    eks pan 2 0.7586799647899636  0.5221511083334797  string string int    float  float
    wye wye 3 0.20460330576630303 0.33831852551664776 string string int    float  float
    eks wye 4 0.38139939387114097 0.13418874328430463 string string int    float  float
    wye pan 5 0.5732889198020006  0.8636244699032729  string string int    float  float

Note that the value of the current field in the for-loop can be gotten either using the bound variable ``value``, or through a **computed field name** using square brackets as in ``$[key]``.

Important note: to avoid inconsistent looping behavior in case you're setting new fields (and/or unsetting existing ones) while looping over the record, **Miller makes a copy of the record before the loop: loop variables are bound from the copy and all other reads/writes involve the record itself**:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --from data/small --opprint put '
      $sum1 = 0;
      $sum2 = 0;
      for (k,v in $*) {
        if (is_numeric(v)) {
          $sum1 +=v;
          $sum2 += $[k];
        }
      }
    '
    a   b   i x                   y                   sum1               sum2
    pan pan 1 0.3467901443380824  0.7268028627434533  2.0735930070815356 8.294372028326142
    eks pan 2 0.7586799647899636  0.5221511083334797  3.280831073123443  13.123324292493772
    wye wye 3 0.20460330576630303 0.33831852551664776 3.5429218312829507 14.171687325131803
    eks wye 4 0.38139939387114097 0.13418874328430463 4.515588137155445  18.06235254862178
    wye pan 5 0.5732889198020006  0.8636244699032729  6.436913389705273  25.747653558821092

It can be confusing to modify the stream record while iterating over a copy of it, so instead you might find it simpler to use a local variable in the loop and only update the stream record after the loop:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --from data/small --opprint put '
      sum = 0;
      for (k,v in $*) {
        if (is_numeric(v)) {
          sum += $[k];
        }
      }
      $sum = sum
    '
    a   b   i x                   y                   sum
    pan pan 1 0.3467901443380824  0.7268028627434533  2.0735930070815356
    eks pan 2 0.7586799647899636  0.5221511083334797  3.280831073123443
    wye wye 3 0.20460330576630303 0.33831852551664776 3.5429218312829507
    eks wye 4 0.38139939387114097 0.13418874328430463 4.515588137155445
    wye pan 5 0.5732889198020006  0.8636244699032729  6.436913389705273

You can also start iterating on sub-hashmaps of an out-of-stream or local variable; you can loop over nested keys; you can loop over all out-of-stream variables.  The bound variables are bound to a copy of the sub-hashmap as it was before the loop started.  The sub-hashmap is specified by square-bracketed indices after ``in``, and additional deeper indices are bound to loop key-variables. The terminal values are bound to the loop value-variable whenever the keys are not too shallow. The value-variable may refer to a terminal (string, number) or it may be map-valued if the map goes deeper. Example indexing is as follows:

.. code-block:: none

    # Parentheses are optional for single key:
    for (k1,           v in @a["b"]["c"]) { ... }
    for ((k1),         v in @a["b"]["c"]) { ... }
    # Parentheses are required for multiple keys:
    for ((k1, k2),     v in @a["b"]["c"]) { ... } # Loop over subhashmap of a variable
    for ((k1, k2, k3), v in @a["b"]["c"]) { ... } # Ditto
    for ((k1, k2, k3), v in @a { ... }            # Loop over variable starting from basename
    for ((k1, k2, k3), v in @* { ... }            # Loop over all variables (k1 is bound to basename)

That's confusing in the abstract, so a concrete example is in order. Suppose the out-of-stream variable ``@myvar`` is populated as follows:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr -n put --jknquoteint -q '
      begin {
        @myvar = {
          1: 2,
          3: { 4 : 5 },
          6: { 7: { 8: 9 } }
        }
      }
      end { dump }
    '
    {
      "myvar": {
        "1": 2,
        "3": {
          "4": 5
        },
        "6": {
          "7": {
            "8": 9
          }
        }
      }
    }

Then we can get at various values as follows:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr -n put --jknquoteint -q '
      begin {
        @myvar = {
          1: 2,
          3: { 4 : 5 },
          6: { 7: { 8: 9 } }
        }
      }
      end {
        for (k, v in @myvar) {
          print
            "key=" . k .
            ",valuetype=" . typeof(v);
        }
      }
    '
    key=1,valuetype=int
    key=3,valuetype=map
    key=6,valuetype=map

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr -n put --jknquoteint -q '
      begin {
        @myvar = {
          1: 2,
          3: { 4 : 5 },
          6: { 7: { 8: 9 } }
        }
      }
      end {
        for ((k1, k2), v in @myvar) {
          print
            "key1=" . k1 .
            ",key2=" . k2 .
            ",valuetype=" . typeof(v);
        }
      }
    '
    key1=3,key2=4,valuetype=int
    key1=6,key2=7,valuetype=map

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr -n put --jknquoteint -q '
      begin {
        @myvar = {
          1: 2,
          3: { 4 : 5 },
          6: { 7: { 8: 9 } }
        }
      }
      end {
        for ((k1, k2), v in @myvar[6]) {
          print
            "key1=" . k1 .
            ",key2=" . k2 .
            ",valuetype=" . typeof(v);
        }
      }
    '
    key1=7,key2=8,valuetype=int

C-style triple-for loops
................................................................

These are supported as follows:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --from data/small --opprint put '
      num suma = 0;
      for (a = 1; a <= NR; a += 1) {
        suma += a;
      }
      $suma = suma;
    '
    a   b   i x                   y                   suma
    pan pan 1 0.3467901443380824  0.7268028627434533  1
    eks pan 2 0.7586799647899636  0.5221511083334797  3
    wye wye 3 0.20460330576630303 0.33831852551664776 6
    eks wye 4 0.38139939387114097 0.13418874328430463 10
    wye pan 5 0.5732889198020006  0.8636244699032729  15

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --from data/small --opprint put '
      num suma = 0;
      num sumb = 0;
      for (num a = 1, num b = 1; a <= NR; a += 1, b *= 2) {
        suma += a;
        sumb += b;
      }
      $suma = suma;
      $sumb = sumb;
    '
    a   b   i x                   y                   suma sumb
    pan pan 1 0.3467901443380824  0.7268028627434533  1    1
    eks pan 2 0.7586799647899636  0.5221511083334797  3    3
    wye wye 3 0.20460330576630303 0.33831852551664776 6    7
    eks wye 4 0.38139939387114097 0.13418874328430463 10   15
    wye pan 5 0.5732889198020006  0.8636244699032729  15   31

Notes:

* In ``for (start; continuation; update) { body }``, the start, continuation, and update statements may be empty, single statements, or multiple comma-separated statements. If the continuation is empty (e.g. ``for(i=1;;i+=1)``) it defaults to true.

* In particular, you may use ``$``-variables and/or ``@``-variables in the start, continuation, and/or update steps (as well as the body, of course).

* The typedecls such as ``int`` or ``num`` are optional.  If a typedecl is provided (for a local variable), it binds a variable scoped to the for-loop regardless of whether a same-name variable is present in outer scope. If a typedecl is not provided, then the variable is scoped to the for-loop if no same-name variable is present in outer scope, or if a same-name variable is present in outer scope then it is modified.

* Miller has no ``++`` or ``--`` operators.

* As with all for/if/while statements in Miller, the curly braces are required even if the body is a single statement, or empty.

Begin/end blocks
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Miller supports an ``awk``-like ``begin/end`` syntax.  The statements in the ``begin`` block are executed before any input records are read; the statements in the ``end`` block are executed after the last input record is read.  (If you want to execute some statement at the start of each file, not at the start of the first file as with ``begin``, you might use a pattern/action block of the form ``FNR == 1 { ... }``.) All statements outside of ``begin`` or ``end`` are, of course, executed on every input record. Semicolons separate statements inside or outside of begin/end blocks; semicolons are required between begin/end block bodies and any subsequent statement.  For example:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put '
      begin { @sum = 0 };
      @x_sum += $x;
      end { emit @x_sum }
    ' ../data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729
    a=zee,b=pan,i=6,x=0.5271261600918548,y=0.49322128674835697
    a=eks,b=zee,i=7,x=0.6117840605678454,y=0.1878849191181694
    a=zee,b=wye,i=8,x=0.5985540091064224,y=0.976181385699006
    a=hat,b=wye,i=9,x=0.03144187646093577,y=0.7495507603507059
    a=pan,b=wye,i=10,x=0.5026260055412137,y=0.9526183602969864
    x_sum=4.536293840335763

Since uninitialized out-of-stream variables default to 0 for addition/substraction and 1 for multiplication when they appear on expression right-hand sides (not quite as in ``awk``, where they'd default to 0 either way), the above can be written more succinctly as

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put '
      @x_sum += $x;
      end { emit @x_sum }
    ' ../data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729
    a=zee,b=pan,i=6,x=0.5271261600918548,y=0.49322128674835697
    a=eks,b=zee,i=7,x=0.6117840605678454,y=0.1878849191181694
    a=zee,b=wye,i=8,x=0.5985540091064224,y=0.976181385699006
    a=hat,b=wye,i=9,x=0.03144187646093577,y=0.7495507603507059
    a=pan,b=wye,i=10,x=0.5026260055412137,y=0.9526183602969864
    x_sum=4.536293840335763

The **put -q** option is a shorthand which suppresses printing of each output record, with only ``emit`` statements being output. So to get only summary outputs, one could write

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put -q '
      @x_sum += $x;
      end { emit @x_sum }
    ' ../data/small
    x_sum=4.536293840335763

We can do similarly with multiple out-of-stream variables:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put -q '
      @x_count += 1;
      @x_sum += $x;
      end {
        emit @x_count;
        emit @x_sum;
      }
    ' ../data/small
    x_count=10
    x_sum=4.536293840335763

This is of course not much different than

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr stats1 -a count,sum -f x ../data/small
    x_count=10,x_sum=4.536293840335763

Note that it's a syntax error for begin/end blocks to refer to field names (beginning with ``$``), since these execute outside the context of input records.

Output statements
----------------------------------------------------------------

You can **output** variable-values or expressions in **five ways**:

* **Assign** them to stream-record fields. For example, ``$cumulative_sum = @sum``. For another example, ``$nr = NR`` adds a field named ``nr`` to each output record, containing the value of the built-in variable ``NR`` as of when that record was ingested.

* Use the **print** or **eprint** keywords which immediately print an expression *directly to standard output or standard error*, respectively. Note that ``dump``, ``edump``, ``print``, and ``eprint`` don't output records which participate in ``then``-chaining; rather, they're just immediate prints to stdout/stderr. The ``printn`` and ``eprintn`` keywords are the same except that they don't print final newlines. Additionally, you can print to a specified file instead of stdout/stderr.

* Use the **dump** or **edump** keywords, which *immediately print all out-of-stream variables as a JSON data structure to the standard output or standard error* (respectively).

* Use **tee** which formats the current stream record (not just an arbitrary string as with **print**) to a specific file.

* Use **emit**/**emitp**/**emitf** to send out-of-stream variables' current values to the output record stream, e.g.  ``@sum += $x; emit @sum`` which produces an extra output record such as ``sum=3.1648382``.

For the first two options you are populating the output-records stream which feeds into the next verb in a ``then``-chain (if any), or which otherwise is formatted for output using ``--o...`` flags.

For the last three options you are sending output directly to standard output, standard error, or a file.

.. _reference-dsl-print-statements:

Print statements
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

The ``print`` statement is perhaps self-explanatory, but with a few light caveats:

* There are four variants: ``print`` goes to stdout with final newline, ``printn`` goes to stdout without final newline (you can include one using "\n" in your output string), ``eprint`` goes to stderr with final newline, and ``eprintn`` goes to stderr without final newline.

* Output goes directly to stdout/stderr, respectively: data produced this way do not go downstream to the next verb in a ``then``-chain. (Use ``emit`` for that.)

* Print statements are for strings (``print "hello"``), or things which can be made into strings: numbers (``print 3``, ``print $a + $b``, or concatenations thereof (``print "a + b = " . ($a + $b)``). Maps (in ``$*``, map-valued out-of-stream or local variables, and map literals) aren't convertible into strings. If you print a map, you get ``{is-a-map}`` as output. Please use ``dump`` to print maps.

* You can redirect print output to a file: ``mlr --from myfile.dat put 'print > "tap.txt", $x'`` ``mlr --from myfile.dat put 'o=$*; print > $a.".txt", $x'``.

* See also :ref:`reference-dsl-redirected-output-statements` for examples.

.. _reference-dsl-dump-statements:

Dump statements
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

The ``dump`` statement is for printing expressions, including maps, directly to stdout/stderr, respectively:

* There are two variants: ``dump`` prints to stdout; ``edump`` prints to stderr.

* Output goes directly to stdout/stderr, respectively: data produced this way do not go downstream to the next verb in a ``then``-chain. (Use ``emit`` for that.)

* You can use ``dump`` to output single strings, numbers, or expressions including map-valued data. Map-valued data are printed as JSON. Miller allows string and integer keys in its map literals while JSON allows only string keys, so use ``mlr put --jknquoteint`` if you want integer-valued map keys not double-quoted.

* If you use ``dump`` (or ``edump``) with no arguments, you get a JSON structure representing the current values of all out-of-stream variables.

* As with ``print``, you can redirect output to files.

* See also :ref:`reference-dsl-redirected-output-statements` for examples.

Tee statements
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Records produced by a ``mlr put`` go downstream to the next verb in your ``then``-chain, if any, or otherwise to standard output.  If you want to additionally copy out records to files, you can do that using ``tee``.

The syntax is, by example, ``mlr --from myfile.dat put 'tee > "tap.dat", $*' then sort -n index``.  First is ``tee >``, then the filename expression (which can be an expression such as ``"tap.".$a.".dat"``), then a comma, then ``$*``. (Nothing else but ``$*`` is teeable.)

See also :ref:`reference-dsl-redirected-output-statements` for examples.

.. _reference-dsl-redirected-output-statements:

Redirected-output statements
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

The **print**, **dump** **tee**, **emitf**, **emit**, and **emitp** keywords all allow you to redirect output to one or more files or pipe-to commands. The filenames/commands are strings which can be constructed using record-dependent values, so you can do things like splitting a table into multiple files, one for each account ID, and so on.

Details:

* The ``print`` and ``dump`` keywords produce output immediately to standard output, or to specified file(s) or pipe-to command if present.

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --help-keyword print
    TOD: port mlr_dsl_keyword_usage

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --help-keyword dump
    TOD: port mlr_dsl_keyword_usage

* ``mlr put`` sends the current record (possibly modified by the ``put`` expression) to the output record stream. Records are then input to the following verb in a ``then``-chain (if any), else printed to standard output (unless ``put -q``). The **tee** keyword *additionally* writes the output record to specified file(s) or pipe-to command, or immediately to ``stdout``/``stderr``.

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --help-keyword tee
    TOD: port mlr_dsl_keyword_usage

* ``mlr put``'s ``emitf``, ``emitp``, and ``emit`` send out-of-stream variables to the output record stream. These are then input to the following verb in a ``then``-chain (if any), else printed to standard output. When redirected with ``>``, ``>>``, or ``|``, they *instead* write the out-of-stream variable(s) to specified file(s) or pipe-to command, or immediately to ``stdout``/``stderr``.

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --help-keyword emitf
    TOD: port mlr_dsl_keyword_usage

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --help-keyword emitp
    TOD: port mlr_dsl_keyword_usage

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --help-keyword emit
    TOD: port mlr_dsl_keyword_usage

.. _reference-dsl-emit-statements:

Emit statements
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

There are three variants: ``emitf``, ``emit``, and ``emitp``. Keep in mind that out-of-stream variables are a nested, multi-level hashmap (directly viewable as JSON using ``dump``), whereas Miller output records are lists of single-level key-value pairs. The three emit variants allow you to control how the multilevel hashmaps are flatten down to output records. You can emit any map-valued expression, including ``$*``, map-valued out-of-stream variables, the entire out-of-stream-variable collection ``@*``, map-valued local variables, map literals, or map-valued function return values.

Use **emitf** to output several out-of-stream variables side-by-side in the same output record. For ``emitf`` these mustn't have indexing using ``@name[...]``. Example:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put -q '@count += 1; @x_sum += $x; @y_sum += $y; end { emitf @count, @x_sum, @y_sum}' data/small
    count=5,x_sum=2.264761728567491,y_sum=2.585085709781158

Use **emit** to output an out-of-stream variable. If it's non-indexed you'll get a simple key-value pair:

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

    $ mlr put -q '@sum += $x; end { dump }' data/small
    {
      "sum": 2.264761728567491
    }

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put -q '@sum += $x; end { emit @sum }' data/small
    sum=2.264761728567491

If it's indexed then use as many names after ``emit`` as there are indices:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put -q '@sum[$a] += $x; end { dump }' data/small
    {
      "sum": {
        "pan": 0.3467901443380824,
        "eks": 1.1400793586611044,
        "wye": 0.7778922255683036
      }
    }

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put -q '@sum[$a] += $x; end { emit @sum, "a" }' data/small
    a=pan,sum=0.3467901443380824
    a=eks,sum=1.1400793586611044
    a=wye,sum=0.7778922255683036

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put -q '@sum[$a][$b] += $x; end { dump }' data/small
    {
      "sum": {
        "pan": {
          "pan": 0.3467901443380824
        },
        "eks": {
          "pan": 0.7586799647899636,
          "wye": 0.38139939387114097
        },
        "wye": {
          "wye": 0.20460330576630303,
          "pan": 0.5732889198020006
        }
      }
    }

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put -q '@sum[$a][$b] += $x; end { emit @sum, "a", "b" }' data/small
    a=pan,b=pan,sum=0.3467901443380824
    a=eks,b=pan,sum=0.7586799647899636
    a=eks,b=wye,sum=0.38139939387114097
    a=wye,b=wye,sum=0.20460330576630303
    a=wye,b=pan,sum=0.5732889198020006

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put -q '@sum[$a][$b][$i] += $x; end { dump }' data/small
    {
      "sum": {
        "pan": {
          "pan": {
            "1": 0.3467901443380824
          }
        },
        "eks": {
          "pan": {
            "2": 0.7586799647899636
          },
          "wye": {
            "4": 0.38139939387114097
          }
        },
        "wye": {
          "wye": {
            "3": 0.20460330576630303
          },
          "pan": {
            "5": 0.5732889198020006
          }
        }
      }
    }

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put -q '@sum[$a][$b][$i] += $x; end { emit @sum, "a", "b", "i" }' data/small
    a=pan,b=pan,i=1,sum=0.3467901443380824
    a=eks,b=pan,i=2,sum=0.7586799647899636
    a=eks,b=wye,i=4,sum=0.38139939387114097
    a=wye,b=wye,i=3,sum=0.20460330576630303
    a=wye,b=pan,i=5,sum=0.5732889198020006

Now for **emitp**: if you have as many names following ``emit`` as there are levels in the out-of-stream variable's hashmap, then ``emit`` and ``emitp`` do the same thing. Where they differ is when you don't specify as many names as there are hashmap levels. In this case, Miller needs to flatten multiple map indices down to output-record keys: ``emitp`` includes full prefixing (hence the ``p`` in ``emitp``) while ``emit`` takes the deepest hashmap key as the output-record key:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put -q '@sum[$a][$b] += $x; end { dump }' data/small
    {
      "sum": {
        "pan": {
          "pan": 0.3467901443380824
        },
        "eks": {
          "pan": 0.7586799647899636,
          "wye": 0.38139939387114097
        },
        "wye": {
          "wye": 0.20460330576630303,
          "pan": 0.5732889198020006
        }
      }
    }

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put -q '@sum[$a][$b] += $x; end { emit @sum, "a" }' data/small
    a=pan,pan=0.3467901443380824
    a=eks,pan=0.7586799647899636,wye=0.38139939387114097
    a=wye,wye=0.20460330576630303,pan=0.5732889198020006

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put -q '@sum[$a][$b] += $x; end { emit @sum }' data/small
    pan.pan=0.3467901443380824,eks.pan=0.7586799647899636,eks.wye=0.38139939387114097,wye.wye=0.20460330576630303,wye.pan=0.5732889198020006

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put -q '@sum[$a][$b] += $x; end { emitp @sum, "a" }' data/small
    a=pan,sum.pan=0.3467901443380824
    a=eks,sum.pan=0.7586799647899636,sum.wye=0.38139939387114097
    a=wye,sum.wye=0.20460330576630303,sum.pan=0.5732889198020006

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put -q '@sum[$a][$b] += $x; end { emitp @sum }' data/small
    sum.pan.pan=0.3467901443380824,sum.eks.pan=0.7586799647899636,sum.eks.wye=0.38139939387114097,sum.wye.wye=0.20460330576630303,sum.wye.pan=0.5732889198020006

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --oxtab put -q '@sum[$a][$b] += $x; end { emitp @sum }' data/small
    sum.pan.pan 0.3467901443380824
    sum.eks.pan 0.7586799647899636
    sum.eks.wye 0.38139939387114097
    sum.wye.wye 0.20460330576630303
    sum.wye.pan 0.5732889198020006

Use **--oflatsep** to specify the character which joins multilevel
keys for ``emitp`` (it defaults to a colon):

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put -q --oflatsep / '@sum[$a][$b] += $x; end { emitp @sum, "a" }' data/small
    a=pan,sum.pan=0.3467901443380824
    a=eks,sum.pan=0.7586799647899636,sum.wye=0.38139939387114097
    a=wye,sum.wye=0.20460330576630303,sum.pan=0.5732889198020006

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put -q --oflatsep / '@sum[$a][$b] += $x; end { emitp @sum }' data/small
    sum.pan.pan=0.3467901443380824,sum.eks.pan=0.7586799647899636,sum.eks.wye=0.38139939387114097,sum.wye.wye=0.20460330576630303,sum.wye.pan=0.5732889198020006

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --oxtab put -q --oflatsep / '@sum[$a][$b] += $x; end { emitp @sum }' data/small
    sum.pan.pan 0.3467901443380824
    sum.eks.pan 0.7586799647899636
    sum.eks.wye 0.38139939387114097
    sum.wye.wye 0.20460330576630303
    sum.wye.pan 0.5732889198020006

Multi-emit statements
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

You can emit **multiple map-valued expressions side-by-side** by
including their names in parentheses:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --from data/medium --opprint put -q '
      @x_count[$a][$b] += 1;
      @x_sum[$a][$b] += $x;
      end {
          for ((a, b), _ in @x_count) {
              @x_mean[a][b] = @x_sum[a][b] / @x_count[a][b]
          }
          emit (@x_sum, @x_count, @x_mean), "a", "b"
      }
    '
    a   b   x_sum              x_count x_mean
    pan pan 219.1851288316854  427     0.5133141190437597
    pan wye 198.43293070748447 395     0.5023618498923658
    pan eks 216.07522773165525 429     0.5036718595143479
    pan hat 205.22277621488686 417     0.492140950155604
    pan zee 205.09751802331917 413     0.4966041598627583
    eks pan 179.96303047250723 371     0.48507555383425127
    eks wye 196.9452860713734  407     0.4838950517724162
    eks zee 176.8803651584733  357     0.49546320772681596
    eks eks 215.91609712937984 413     0.5227992666570941
    eks hat 208.783170520597   417     0.5006790659966355
    wye wye 185.29584980261419 377     0.49150092785839306
    wye pan 195.84790012056564 392     0.4996119901034838
    wye hat 212.0331829346132  426     0.4977304763723314
    wye zee 194.77404756708714 385     0.5059066170573692
    wye eks 204.8129608356315  386     0.5306035254809106
    zee pan 202.21380378504267 389     0.5198298297816007
    zee wye 233.9913939194868  455     0.5142667998230479
    zee eks 190.9617780631925  391     0.4883932942792647
    zee zee 206.64063510417319 403     0.5127559183726382
    zee hat 191.30000620900935 409     0.46772617655014515
    hat wye 208.8830097609959  423     0.49381326184632596
    hat zee 196.3494502965293  385     0.5099985721987774
    hat eks 189.0067933716193  389     0.48587864619953547
    hat hat 182.8535323148762  381     0.47993053101017374
    hat pan 168.5538067327806  363     0.4643355557376876

What this does is walk through the first out-of-stream variable (``@x_sum`` in this example) as usual, then for each keylist found (e.g. ``pan,wye``), include the values for the remaining out-of-stream variables (here, ``@x_count`` and ``@x_mean``). You should use this when all out-of-stream variables in the emit statement have **the same shape and the same keylists**.

Emit-all statements
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Use **emit all** (or ``emit @*`` which is synonymous) to output all out-of-stream variables. You can use the following idiom to get various accumulators output side-by-side (reminiscent of ``mlr stats1``):

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --from data/small --opprint put -q '@v[$a][$b]["sum"] += $x; @v[$a][$b]["count"] += 1; end{emit @*,"a","b"}'
    a b   pan.sum            pan.count
    v pan 0.3467901443380824 1
    
    a b   pan.sum            pan.count wye.sum             wye.count
    v eks 0.7586799647899636 1         0.38139939387114097 1
    
    a b   wye.sum             wye.count pan.sum            pan.count
    v wye 0.20460330576630303 1         0.5732889198020006 1

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --from data/small --opprint put -q '@sum[$a][$b] += $x; @count[$a][$b] += 1; end{emit @*,"a","b"}'
    a   b   pan
    sum pan 0.3467901443380824
    
    a   b   pan                wye
    sum eks 0.7586799647899636 0.38139939387114097
    
    a   b   wye                 pan
    sum wye 0.20460330576630303 0.5732889198020006
    
    a     b   pan
    count pan 1
    
    a     b   pan wye
    count eks 1   1
    
    a     b   wye pan
    count wye 1   1

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --from data/small --opprint put -q '@sum[$a][$b] += $x; @count[$a][$b] += 1; end{emit (@sum, @count),"a","b"}'
    a   b   sum                 count
    pan pan 0.3467901443380824  1
    eks pan 0.7586799647899636  1
    eks wye 0.38139939387114097 1
    wye wye 0.20460330576630303 1
    wye pan 0.5732889198020006  1

Unset statements
----------------------------------------------------------------

You can clear a map key by assigning the empty string as its value: ``$x=""`` or ``@x=""``. Using ``unset`` you can remove the key entirely. Examples:

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

    $ mlr put 'unset $x, $a' data/small
    b=pan,i=1,y=0.7268028627434533
    b=pan,i=2,y=0.5221511083334797
    b=wye,i=3,y=0.33831852551664776
    b=wye,i=4,y=0.13418874328430463
    b=pan,i=5,y=0.8636244699032729

This can also be done, of course, using ``mlr cut -x``. You can also clear out-of-stream or local variables, at the base name level, or at an indexed sublevel:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put -q '@sum[$a][$b] += $x; end { dump; unset @sum; dump }' data/small
    {
      "sum": {
        "pan": {
          "pan": 0.3467901443380824
        },
        "eks": {
          "pan": 0.7586799647899636,
          "wye": 0.38139939387114097
        },
        "wye": {
          "wye": 0.20460330576630303,
          "pan": 0.5732889198020006
        }
      }
    }
    {}

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put -q '@sum[$a][$b] += $x; end { dump; unset @sum["eks"]; dump }' data/small
    {
      "sum": {
        "pan": {
          "pan": 0.3467901443380824
        },
        "eks": {
          "pan": 0.7586799647899636,
          "wye": 0.38139939387114097
        },
        "wye": {
          "wye": 0.20460330576630303,
          "pan": 0.5732889198020006
        }
      }
    }
    {
      "sum": {
        "pan": {
          "pan": 0.3467901443380824
        },
        "wye": {
          "wye": 0.20460330576630303,
          "pan": 0.5732889198020006
        }
      }
    }

If you use ``unset all`` (or ``unset @*`` which is synonymous), that will unset all out-of-stream variables which have been defined up to that point.

Filter statements
----------------------------------------------------------------

You can use ``filter`` within ``put``. In fact, the following two are synonymous:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr filter 'NR==2 || NR==3' data/small
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put 'filter NR==2 || NR==3' data/small
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776

The former, of course, is much easier to type. But the latter allows you to define more complex expressions for the filter, and/or do other things in addition to the filter:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put '@running_sum += $x; filter @running_sum > 1.3' data/small
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr put '$z = $x * $y; filter $z > 0.3' data/small
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797,z=0.3961455844854848
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,z=0.4951063394654227

Built-in functions for filter and put, summary
----------------------------------------------------------------

mlr: option "--list-all-functions-as-table" not recognized.
Please run "mlr --help" for usage information.

Built-in functions for filter and put
----------------------------------------------------------------

Each function takes a specific number of arguments, as shown below, except for functions marked as variadic such as ``min`` and ``max``. (The latter compute min and max of any number of numerical arguments.) There is no notion of optional or default-on-absent arguments. All argument-passing is positional rather than by name; arguments are passed by value, not by reference.

You can get a list of all functions using **mlr -F**.


.. _reference-dsl-plus:

\+
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    +  (class=arithmetic #args=1,2) Addition as binary operator; unary plus operator.



.. _reference-dsl-minus:

\-
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    -  (class=arithmetic #args=1,2) Subtraction as binary operator; unary negation operator.



.. _reference-dsl-times:

\*
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    *  (class=arithmetic #args=2) Multiplication, with integer*integer overflow to float.



.. _reference-dsl-/:

/
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    /  (class=arithmetic #args=2) Division. Integer / integer is floating-point.



.. _reference-dsl-//:

//
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    //  (class=arithmetic #args=2) Pythonic integer division, rounding toward negative.



.. _reference-dsl-exponentiation:

\**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    **  (class=arithmetic #args=2) Exponentiation. Same as pow, but as an infix operator.



.. _reference-dsl-.+:

.+
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    .+  (class=arithmetic #args=2) Addition, with integer-to-integer overflow.



.. _reference-dsl-.-:

.-
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    .-  (class=arithmetic #args=2) Subtraction, with integer-to-integer overflow.



.. _reference-dsl-.*:

.*
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    .*  (class=arithmetic #args=2) Multiplication, with integer-to-integer overflow.



.. _reference-dsl-./:

./
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    ./  (class=arithmetic #args=2) Integer division; not pythonic.



.. _reference-dsl-%:

%
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    %  (class=arithmetic #args=2) Remainder; never negative-valued (pythonic).



.. _reference-dsl-~:

~
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    ~  (class=arithmetic #args=1) Bitwise NOT. Beware '$y=~$x' since =~ is the
    regex-match operator: try '$y = ~$x'.



.. _reference-dsl-&:

&
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    &  (class=arithmetic #args=2) Bitwise AND.



.. _reference-dsl-bitwise-or:

\|
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    |  (class=arithmetic #args=2) Bitwise OR.



.. _reference-dsl-^:

^
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    ^  (class=arithmetic #args=2) Bitwise XOR.



.. _reference-dsl-<<:

<<
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    <<  (class=arithmetic #args=2) Bitwise left-shift.



.. _reference-dsl->>:

>>
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    >>  (class=arithmetic #args=2) Bitwise signed right-shift.



.. _reference-dsl->>>:

>>>
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    >>>  (class=arithmetic #args=2) Bitwise unsigned right-shift.



.. _reference-dsl-!:

!
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    !  (class=boolean #args=1) Logical negation.



.. _reference-dsl-==:

==
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    ==  (class=boolean #args=2) String/numeric equality. Mixing number and string results in string compare.



.. _reference-dsl-!=:

!=
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    !=  (class=boolean #args=2) String/numeric inequality. Mixing number and string results in string compare.



.. _reference-dsl->:

>
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    >  (class=boolean #args=2) String/numeric greater-than. Mixing number and string results in string compare.



.. _reference-dsl->=:

>=
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    >=  (class=boolean #args=2) String/numeric greater-than-or-equals. Mixing number and string results in string compare.



.. _reference-dsl-<:

<
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    <  (class=boolean #args=2) String/numeric less-than. Mixing number and string results in string compare.



.. _reference-dsl-<=:

<=
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    <=  (class=boolean #args=2) String/numeric less-than-or-equals. Mixing number and string results in string compare.



.. _reference-dsl-=~:

=~
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    =~  (class=boolean #args=2) String (left-hand side) matches regex (right-hand side), e.g. '$name =~ "^a.*b$"'.



.. _reference-dsl-!=~:

!=~
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    !=~  (class=boolean #args=2) String (left-hand side) does not match regex (right-hand side), e.g. '$name !=~ "^a.*b$"'.



.. _reference-dsl-&&:

&&
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    &&  (class=boolean #args=2) Logical AND.



.. _reference-dsl-||:

||
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    ||  (class=boolean #args=2) Logical OR.



.. _reference-dsl-^^:

^^
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    ^^  (class=boolean #args=2) Logical XOR.



.. _reference-dsl-??:

??
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    ??  (class=boolean #args=2) Absent-coalesce operator. $a ?? 1 evaluates to 1 if $a isn't defined in the current record.



.. _reference-dsl-???:

???
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    ???  (class=boolean #args=2) Absent-coalesce operator. $a ?? 1 evaluates to 1 if $a isn't defined in the current record, or has empty value.



.. _reference-dsl-question-mark-colon:

\?
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    ?:  (class=boolean #args=3) Standard ternary operator.



.. _reference-dsl-.:

.
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    .  (class=string #args=2) String concatenation.



.. _reference-dsl-abs:

abs
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    abs  (class=math #args=1) Absolute value.



.. _reference-dsl-acos:

acos
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    acos  (class=math #args=1) Inverse trigonometric cosine.



.. _reference-dsl-acosh:

acosh
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    acosh  (class=math #args=1) Inverse hyperbolic cosine.



.. _reference-dsl-append:

append
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    append  (class=maps/arrays #args=2) Appends second argument to end of first argument, which must be an array.



.. _reference-dsl-arrayify:

arrayify
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    arrayify  (class=maps/arrays #args=1) Walks through a nested map/array, converting any map with consecutive keys
    "1", "2", ... into an array. Useful to wrap the output of unflatten.



.. _reference-dsl-asin:

asin
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asin  (class=math #args=1) Inverse trigonometric sine.



.. _reference-dsl-asinh:

asinh
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asinh  (class=math #args=1) Inverse hyperbolic sine.



.. _reference-dsl-asserting_absent:

asserting_absent
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asserting_absent  (class=typing #args=1) Aborts with an error if is_absent on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_array:

asserting_array
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asserting_array  (class=typing #args=1) Aborts with an error if is_array on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_bool:

asserting_bool
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asserting_bool  (class=typing #args=1) Aborts with an error if is_bool on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_boolean:

asserting_boolean
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asserting_boolean  (class=typing #args=1) Aborts with an error if is_boolean on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_empty:

asserting_empty
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asserting_empty  (class=typing #args=1) Aborts with an error if is_empty on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_empty_map:

asserting_empty_map
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asserting_empty_map  (class=typing #args=1) Aborts with an error if is_empty_map on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_error:

asserting_error
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asserting_error  (class=typing #args=1) Aborts with an error if is_error on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_float:

asserting_float
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asserting_float  (class=typing #args=1) Aborts with an error if is_float on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_int:

asserting_int
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asserting_int  (class=typing #args=1) Aborts with an error if is_int on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_map:

asserting_map
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asserting_map  (class=typing #args=1) Aborts with an error if is_map on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_nonempty_map:

asserting_nonempty_map
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asserting_nonempty_map  (class=typing #args=1) Aborts with an error if is_nonempty_map on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_not_array:

asserting_not_array
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asserting_not_array  (class=typing #args=1) Aborts with an error if is_not_array on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_not_empty:

asserting_not_empty
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asserting_not_empty  (class=typing #args=1) Aborts with an error if is_not_empty on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_not_map:

asserting_not_map
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asserting_not_map  (class=typing #args=1) Aborts with an error if is_not_map on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_not_null:

asserting_not_null
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asserting_not_null  (class=typing #args=1) Aborts with an error if is_not_null on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_null:

asserting_null
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asserting_null  (class=typing #args=1) Aborts with an error if is_null on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_numeric:

asserting_numeric
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asserting_numeric  (class=typing #args=1) Aborts with an error if is_numeric on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_present:

asserting_present
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asserting_present  (class=typing #args=1) Aborts with an error if is_present on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_string:

asserting_string
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    asserting_string  (class=typing #args=1) Aborts with an error if is_string on the argument returns false,
    else returns its argument.



.. _reference-dsl-atan:

atan
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    atan  (class=math #args=1) One-argument arctangent.



.. _reference-dsl-atan2:

atan2
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    atan2  (class=math #args=2) Two-argument arctangent.



.. _reference-dsl-atanh:

atanh
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    atanh  (class=math #args=1) Inverse hyperbolic tangent.



.. _reference-dsl-bitcount:

bitcount
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    bitcount  (class=arithmetic #args=1) Count of 1-bits.



.. _reference-dsl-boolean:

boolean
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    boolean  (class=conversion #args=1) Convert int/float/bool/string to boolean.



.. _reference-dsl-capitalize:

capitalize
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    capitalize  (class=string #args=1) Convert string's first character to uppercase.



.. _reference-dsl-cbrt:

cbrt
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    cbrt  (class=math #args=1) Cube root.



.. _reference-dsl-ceil:

ceil
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    ceil  (class=math #args=1) Ceiling: nearest integer at or above.



.. _reference-dsl-clean_whitespace:

clean_whitespace
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    clean_whitespace  (class=string #args=1) Same as collapse_whitespace and strip.



.. _reference-dsl-collapse_whitespace:

collapse_whitespace
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    collapse_whitespace  (class=string #args=1) Strip repeated whitespace from string.



.. _reference-dsl-cos:

cos
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    cos  (class=math #args=1) Trigonometric cosine.



.. _reference-dsl-cosh:

cosh
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    cosh  (class=math #args=1) Hyperbolic cosine.



.. _reference-dsl-depth:

depth
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    depth  (class=maps/arrays #args=1) Prints maximum depth of map/array. Scalars have depth 0.



.. _reference-dsl-dhms2fsec:

dhms2fsec
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    dhms2fsec  (class=time #args=1) Recovers floating-point seconds as in dhms2fsec("5d18h53m20.250000s") = 500000.250000
    



.. _reference-dsl-dhms2sec:

dhms2sec
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    dhms2sec  (class=time #args=1) Recovers integer seconds as in dhms2sec("5d18h53m20s") = 500000
    



.. _reference-dsl-erf:

erf
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    erf  (class=math #args=1) Error function.



.. _reference-dsl-erfc:

erfc
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    erfc  (class=math #args=1) Complementary error function.



.. _reference-dsl-exp:

exp
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    exp  (class=math #args=1) Exponential function e**x.



.. _reference-dsl-expm1:

expm1
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    expm1  (class=math #args=1) e**x - 1.



.. _reference-dsl-flatten:

flatten
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    flatten  (class=maps/arrays #args=3) Flattens multi-level maps to single-level ones. Examples:
    flatten("a", ".", {"b": { "c": 4 }}) is {"a.b.c" : 4}.
    flatten("", ".", {"a": { "b": 3 }}) is {"a.b" : 3}.
    Two-argument version: flatten($*, ".") is the same as flatten("", ".", $*).
    Useful for nested JSON-like structures for non-JSON file formats like CSV.



.. _reference-dsl-float:

float
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    float  (class=conversion #args=1) Convert int/float/bool/string to float.



.. _reference-dsl-floor:

floor
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    floor  (class=math #args=1) Floor: nearest integer at or below.



.. _reference-dsl-fmtnum:

fmtnum
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    fmtnum  (class=conversion #args=2) Convert int/float/bool to string using
    printf-style format string, e.g. '$s = fmtnum($n, "%06lld")'.



.. _reference-dsl-fsec2dhms:

fsec2dhms
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    fsec2dhms  (class=time #args=1) Formats floating-point seconds as in fsec2dhms(500000.25) = "5d18h53m20.250000s"
    



.. _reference-dsl-fsec2hms:

fsec2hms
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    fsec2hms  (class=time #args=1) Formats floating-point seconds as in fsec2hms(5000.25) = "01:23:20.250000"
    



.. _reference-dsl-get_keys:

get_keys
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    get_keys  (class=maps/arrays #args=1) Returns array of keys of map or array



.. _reference-dsl-get_values:

get_values
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    get_values  (class=maps/arrays #args=1) Returns array of keys of map or array -- in the latter case, returns a copy of the array



.. _reference-dsl-gmt2sec:

gmt2sec
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    gmt2sec  (class=time #args=1) Parses GMT timestamp as integer seconds since the epoch.



.. _reference-dsl-gsub:

gsub
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    gsub  (class=string #args=3) Example: '$name=gsub($name, "old", "new")' (replace all).



.. _reference-dsl-haskey:

haskey
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    haskey  (class=maps/arrays #args=2) True/false if map has/hasn't key, e.g. 'haskey($*, "a")' or
    'haskey(mymap, mykey)', or true/false if array index is in bounds / out of bounds.
    Error if 1st argument is not a map or array. Note -n..-1 alias to 1..n in Miller arrays.



.. _reference-dsl-hexfmt:

hexfmt
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    hexfmt  (class=conversion #args=1) Convert int to hex string, e.g. 255 to "0xff".



.. _reference-dsl-hms2fsec:

hms2fsec
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    hms2fsec  (class=time #args=1) Recovers floating-point seconds as in hms2fsec("01:23:20.250000") = 5000.250000
    



.. _reference-dsl-hms2sec:

hms2sec
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    hms2sec  (class=time #args=1) Recovers integer seconds as in hms2sec("01:23:20") = 5000
    



.. _reference-dsl-hostname:

hostname
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    hostname  (class=system #args=0) Returns the hostname as a string.



.. _reference-dsl-int:

int
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    int  (class=conversion #args=1) Convert int/float/bool/string to int.



.. _reference-dsl-invqnorm:

invqnorm
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    invqnorm  (class=math #args=1) Inverse of normal cumulative distribution function.
    Note that invqorm(urand()) is normally distributed.



.. _reference-dsl-is_absent:

is_absent
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    is_absent  (class=typing #args=1) False if field is present in input, true otherwise



.. _reference-dsl-is_array:

is_array
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    is_array  (class=typing #args=1) True if argument is an array.



.. _reference-dsl-is_bool:

is_bool
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    is_bool  (class=typing #args=1) True if field is present with boolean value. Synonymous with is_boolean.



.. _reference-dsl-is_boolean:

is_boolean
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    is_boolean  (class=typing #args=1) True if field is present with boolean value. Synonymous with is_bool.



.. _reference-dsl-is_empty:

is_empty
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    is_empty  (class=typing #args=1) True if field is present in input with empty string value, false otherwise.



.. _reference-dsl-is_empty_map:

is_empty_map
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    is_empty_map  (class=typing #args=1) True if argument is a map which is empty.



.. _reference-dsl-is_error:

is_error
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    is_error  (class=typing #args=1) True if if argument is an error, such as taking string length of an integer.



.. _reference-dsl-is_float:

is_float
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    is_float  (class=typing #args=1) True if field is present with value inferred to be float



.. _reference-dsl-is_int:

is_int
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    is_int  (class=typing #args=1) True if field is present with value inferred to be int



.. _reference-dsl-is_map:

is_map
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    is_map  (class=typing #args=1) True if argument is a map.



.. _reference-dsl-is_nonempty_map:

is_nonempty_map
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    is_nonempty_map  (class=typing #args=1) True if argument is a map which is non-empty.



.. _reference-dsl-is_not_array:

is_not_array
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    is_not_array  (class=typing #args=1) True if argument is not an array.



.. _reference-dsl-is_not_empty:

is_not_empty
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    is_not_empty  (class=typing #args=1) False if field is present in input with empty value, true otherwise



.. _reference-dsl-is_not_map:

is_not_map
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    is_not_map  (class=typing #args=1) True if argument is not a map.



.. _reference-dsl-is_not_null:

is_not_null
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    is_not_null  (class=typing #args=1) False if argument is null (empty or absent), true otherwise.



.. _reference-dsl-is_null:

is_null
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    is_null  (class=typing #args=1) True if argument is null (empty or absent), false otherwise.



.. _reference-dsl-is_numeric:

is_numeric
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    is_numeric  (class=typing #args=1) True if field is present with value inferred to be int or float



.. _reference-dsl-is_present:

is_present
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    is_present  (class=typing #args=1) True if field is present in input, false otherwise.



.. _reference-dsl-is_string:

is_string
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    is_string  (class=typing #args=1) True if field is present with string (including empty-string) value



.. _reference-dsl-joink:

joink
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    joink  (class=conversion #args=2) Makes string from map/array keys. Examples:
    joink({"a":3,"b":4,"c":5}, ",") = "a,b,c"
    joink([1,2,3], ",") = "1,2,3".



.. _reference-dsl-joinkv:

joinkv
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    joinkv  (class=conversion #args=3) Makes string from map/array key-value pairs. Examples:
    joinkv([3,4,5], "=", ",") = "1=3,2=4,3=5"
    joinkv({"a":3,"b":4,"c":5}, "=", ",") = "a=3,b=4,c=5"



.. _reference-dsl-joinv:

joinv
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    joinv  (class=conversion #args=2) Makes string from map/array values.
    joinv([3,4,5], ",") = "3,4,5"
    joinv({"a":3,"b":4,"c":5}, ",") = "3,4,5"



.. _reference-dsl-json_parse:

json_parse
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    json_parse  (class=maps/arrays #args=1) Converts value from JSON-formatted string.



.. _reference-dsl-json_stringify:

json_stringify
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    json_stringify  (class=maps/arrays #args=1,2) Converts value to JSON-formatted string. Default output is single-line.
    With optional second boolean argument set to true, produces multiline output.



.. _reference-dsl-leafcount:

leafcount
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    leafcount  (class=maps/arrays #args=1) Counts total number of terminal values in map/array. For single-level
    map/array, same as length.



.. _reference-dsl-length:

length
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    length  (class=maps/arrays #args=1) Counts number of top-level entries in array/map. Scalars have length 1.



.. _reference-dsl-log:

log
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    log  (class=math #args=1) Natural (base-e) logarithm.



.. _reference-dsl-log10:

log10
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    log10  (class=math #args=1) Base-10 logarithm.



.. _reference-dsl-log1p:

log1p
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    log1p  (class=math #args=1) log(1-x).



.. _reference-dsl-logifit:

logifit
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    logifit  (class=math #args=3)  Given m and b from logistic regression, compute fit:
    $yhat=logifit($x,$m,$b).



.. _reference-dsl-lstrip:

lstrip
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    lstrip  (class=string #args=1) Strip leading whitespace from string.



.. _reference-dsl-madd:

madd
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    madd  (class=arithmetic #args=3) a + b mod m (integers)



.. _reference-dsl-mapdiff:

mapdiff
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    mapdiff  (class=maps/arrays #args=variadic) With 0 args, returns empty map. With 1 arg, returns copy of arg.
    With 2 or more, returns copy of arg 1 with all keys from any of remaining
    argument maps removed.



.. _reference-dsl-mapexcept:

mapexcept
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    mapexcept  (class=maps/arrays #args=variadic) Returns a map with keys from remaining arguments, if any, unset.
    Remaining arguments can be strings or arrays of string.
    E.g. 'mapexcept({1:2,3:4,5:6}, 1, 5, 7)' is '{3:4}'
    and  'mapexcept({1:2,3:4,5:6}, [1, 5, 7])' is '{3:4}'.



.. _reference-dsl-mapselect:

mapselect
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    mapselect  (class=maps/arrays #args=variadic) Returns a map with only keys from remaining arguments set.
    Remaining arguments can be strings or arrays of string.
    E.g. 'mapselect({1:2,3:4,5:6}, 1, 5, 7)' is '{1:2,5:6}'
    and  'mapselect({1:2,3:4,5:6}, [1, 5, 7])' is '{1:2,5:6}'.



.. _reference-dsl-mapsum:

mapsum
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    mapsum  (class=maps/arrays #args=variadic) With 0 args, returns empty map. With >= 1 arg, returns a map with
    key-value pairs from all arguments. Rightmost collisions win, e.g.
    'mapsum({1:2,3:4},{1:5})' is '{1:5,3:4}'.



.. _reference-dsl-max:

max
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    max  (class=math #args=variadic) Max of n numbers; null loses.



.. _reference-dsl-md5:

md5
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    md5  (class=hashing #args=1) MD5 hash.



.. _reference-dsl-mexp:

mexp
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    mexp  (class=arithmetic #args=3) a ** b mod m (integers)



.. _reference-dsl-min:

min
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    min  (class=math #args=variadic) Min of n numbers; null loses.



.. _reference-dsl-mmul:

mmul
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    mmul  (class=arithmetic #args=3) a * b mod m (integers)



.. _reference-dsl-msub:

msub
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    msub  (class=arithmetic #args=3) a - b mod m (integers)



.. _reference-dsl-os:

os
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    os  (class=system #args=0) Returns the operating-system name as a string.



.. _reference-dsl-pow:

pow
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    pow  (class=arithmetic #args=2) Exponentiation. Same as **, but as a function.



.. _reference-dsl-qnorm:

qnorm
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    qnorm  (class=math #args=1) Normal cumulative distribution function.



.. _reference-dsl-regextract:

regextract
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    regextract  (class=string #args=2) Example: '$name=regextract($name, "[A-Z]{3}[0-9]{2}")'



.. _reference-dsl-regextract_or_else:

regextract_or_else
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    regextract_or_else  (class=string #args=3) Example: '$name=regextract_or_else($name, "[A-Z]{3}[0-9]{2}", "default")'



.. _reference-dsl-round:

round
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    round  (class=math #args=1) Round to nearest integer.



.. _reference-dsl-roundm:

roundm
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    roundm  (class=math #args=2) Round to nearest multiple of m: roundm($x,$m) is
    the same as round($x/$m)*$m.



.. _reference-dsl-rstrip:

rstrip
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    rstrip  (class=string #args=1) Strip trailing whitespace from string.



.. _reference-dsl-sec2dhms:

sec2dhms
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    sec2dhms  (class=time #args=1) Formats integer seconds as in sec2dhms(500000) = "5d18h53m20s"
    



.. _reference-dsl-sec2gmt:

sec2gmt
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    sec2gmt  (class=time #args=1,2) Formats seconds since epoch (integer part)
    as GMT timestamp, e.g. sec2gmt(1440768801.7) = "2015-08-28T13:33:21Z".
    Leaves non-numbers as-is. With second integer argument n, includes n decimal places
    for the seconds part



.. _reference-dsl-sec2gmtdate:

sec2gmtdate
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    sec2gmtdate  (class=time #args=1) Formats seconds since epoch (integer part)
    as GMT timestamp with year-month-date, e.g. sec2gmtdate(1440768801.7) = "2015-08-28".
    Leaves non-numbers as-is.
    



.. _reference-dsl-sec2hms:

sec2hms
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    sec2hms  (class=time #args=1) Formats integer seconds as in sec2hms(5000) = "01:23:20"
    



.. _reference-dsl-sgn:

sgn
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    sgn  (class=math #args=1)  +1, 0, -1 for positive, zero, negative input respectively.



.. _reference-dsl-sha1:

sha1
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    sha1  (class=hashing #args=1) SHA1 hash.



.. _reference-dsl-sha256:

sha256
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    sha256  (class=hashing #args=1) SHA256 hash.



.. _reference-dsl-sha512:

sha512
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    sha512  (class=hashing #args=1) SHA512 hash.



.. _reference-dsl-sin:

sin
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    sin  (class=math #args=1) Trigonometric sine.



.. _reference-dsl-sinh:

sinh
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    sinh  (class=math #args=1) Hyperbolic sine.



.. _reference-dsl-splita:

splita
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    splita  (class=conversion #args=2) Splits string into array with type inference. Example:
    splita("3,4,5", ",") = [3,4,5]



.. _reference-dsl-splitax:

splitax
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    splitax  (class=conversion #args=2) Splits string into array without type inference. Example:
    splita("3,4,5", ",") = ["3","4","5"]



.. _reference-dsl-splitkv:

splitkv
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    splitkv  (class=conversion #args=3) Splits string by separators into map with type inference. Example:
    splitkv("a=3,b=4,c=5", "=", ",") = {"a":3,"b":4,"c":5}



.. _reference-dsl-splitkvx:

splitkvx
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    splitkvx  (class=conversion #args=3) Splits string by separators into map without type inference (keys and
    values are strings). Example:
    splitkvx("a=3,b=4,c=5", "=", ",") = {"a":"3","b":"4","c":"5"}



.. _reference-dsl-splitnv:

splitnv
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    splitnv  (class=conversion #args=2) Splits string by separator into integer-indexed map with type inference. Example:
    splitnv("a,b,c", ",") = {"1":"a","2":"b","3":"c"}



.. _reference-dsl-splitnvx:

splitnvx
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    splitnvx  (class=conversion #args=2) Splits string by separator into integer-indexed map without type
    inference (values are strings). Example:
    splitnvx("3,4,5", ",") = {"1":"3","2":"4","3":"5"}



.. _reference-dsl-sqrt:

sqrt
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    sqrt  (class=math #args=1) Square root.



.. _reference-dsl-ssub:

ssub
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    ssub  (class=string #args=3) Like sub but does no regexing. No characters are special.



.. _reference-dsl-strftime:

strftime
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    strftime  (class=time #args=2)  Formats seconds since the epoch as timestamp, e.g.
    	strftime(1440768801.7,"%Y-%m-%dT%H:%M:%SZ") = "2015-08-28T13:33:21Z", and
    	strftime(1440768801.7,"%Y-%m-%dT%H:%M:%3SZ") = "2015-08-28T13:33:21.700Z".
    	Format strings are as in the C library (please see "man strftime" on your system),
    	with the Miller-specific addition of "%1S" through "%9S" which format the seconds
    	with 1 through 9 decimal places, respectively. ("%S" uses no decimal places.)
    	See also strftime_local.
    



.. _reference-dsl-string:

string
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    string  (class=conversion #args=1) Convert int/float/bool/string/array/map to string.



.. _reference-dsl-strip:

strip
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    strip  (class=string #args=1) Strip leading and trailing whitespace from string.



.. _reference-dsl-strlen:

strlen
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    strlen  (class=string #args=1) String length.



.. _reference-dsl-strptime:

strptime
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    strptime  (class=time #args=2) strptime: Parses timestamp as floating-point seconds since the epoch,
    	e.g. strptime("2015-08-28T13:33:21Z","%Y-%m-%dT%H:%M:%SZ") = 1440768801.000000,
    	and  strptime("2015-08-28T13:33:21.345Z","%Y-%m-%dT%H:%M:%SZ") = 1440768801.345000.
    	See also strptime_local.
    



.. _reference-dsl-sub:

sub
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    sub  (class=string #args=3) Example: '$name=sub($name, "old", "new")' (replace once).



.. _reference-dsl-substr:

substr
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    substr  (class=string #args=3) substr is an alias for substr0. See also substr1. Miller is generally 1-up
    with all array indices, but, this is a backward-compatibility issue with Miller 5 and below.
    Arrays are new in Miller 6; the substr function is older.



.. _reference-dsl-substr0:

substr0
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    substr0  (class=string #args=3) substr0(s,m,n) gives substring of s from 0-up position m to n
    inclusive. Negative indices -len .. -1 alias to 0 .. len-1.



.. _reference-dsl-substr1:

substr1
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    substr1  (class=string #args=3) substr1(s,m,n) gives substring of s from 1-up position m to n
    inclusive. Negative indices -len .. -1 alias to 1 .. len.



.. _reference-dsl-system:

system
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    system  (class=system #args=1) Run command string, yielding its stdout minus final carriage return.



.. _reference-dsl-systime:

systime
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    systime  (class=time #args=0) help string will go here



.. _reference-dsl-systimeint:

systimeint
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    systimeint  (class=time #args=0) help string will go here



.. _reference-dsl-tan:

tan
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    tan  (class=math #args=1) Trigonometric tangent.



.. _reference-dsl-tanh:

tanh
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    tanh  (class=math #args=1) Hyperbolic tangent.



.. _reference-dsl-tolower:

tolower
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    tolower  (class=string #args=1) Convert string to lowercase.



.. _reference-dsl-toupper:

toupper
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    toupper  (class=string #args=1) Convert string to uppercase.



.. _reference-dsl-truncate:

truncate
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    truncate  (class=string #args=2) Truncates string first argument to max length of int second argument.



.. _reference-dsl-typeof:

typeof
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    typeof  (class=typing #args=1) Convert argument to type of argument (e.g. "str"). For debug.



.. _reference-dsl-unflatten:

unflatten
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    unflatten  (class=maps/arrays #args=2) Reverses flatten. Example:
    unflatten({"a.b.c" : 4}, ".") is {"a": "b": { "c": 4 }}}.
    Useful for nested JSON-like structures for non-JSON file formats like CSV.
    See also arrayify.



.. _reference-dsl-uptime:

uptime
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    uptime  (class=time #args=0) help string will go here



.. _reference-dsl-urand:

urand
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    urand  (class=math #args=0) Floating-point numbers uniformly distributed on the unit interval.
    Int-valued example: '$n=floor(20+urand()*11)'.



.. _reference-dsl-urand32:

urand32
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    urand32  (class=math #args=0) Integer uniformly distributed 0 and 2**32-1 inclusive.



.. _reference-dsl-urandint:

urandint
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    urandint  (class=math #args=2) Integer uniformly distributed between inclusive integer endpoints.



.. _reference-dsl-urandrange:

urandrange
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    urandrange  (class=math #args=2) Floating-point numbers uniformly distributed on the interval [a, b).



.. _reference-dsl-version:

version
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    version  (class=system #args=0) Returns the Miller version as a string.



User-defined functions and subroutines
----------------------------------------------------------------

As of Miller 5.0.0 you can define your own functions, as well as subroutines.

User-defined functions
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Here's the obligatory example of a recursive function to compute the factorial function:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint --from data/small put '
        func f(n) {
            if (is_numeric(n)) {
                if (n > 0) {
                    return n * f(n-1);
                } else {
                    return 1;
                }
            }
            # implicitly return absent-null if non-numeric
        }
        $ox = f($x + NR);
        $oi = f($i);
    '
    a   b   i x                   y                   ox                  oi
    pan pan 1 0.3467901443380824  0.7268028627434533  0.46705354854811026 1
    eks pan 2 0.7586799647899636  0.5221511083334797  3.680838410072862   2
    wye wye 3 0.20460330576630303 0.33831852551664776 1.7412511955594865  6
    eks wye 4 0.38139939387114097 0.13418874328430463 18.588348778962008  24
    wye pan 5 0.5732889198020006  0.8636244699032729  211.38730958519247  120

Properties of user-defined functions:

* Function bodies start with ``func`` and a parameter list, defined outside of ``begin``, ``end``, or other ``func`` or ``subr`` blocks. (I.e. the Miller DSL has no nested functions.)

* A function (uniqified by its name) may not be redefined: either by redefining a user-defined function, or by redefining a built-in function. However, functions and subroutines have separate namespaces: you can define a subroutine ``log`` which does not clash with the mathematical ``log`` function.

* Functions may be defined either before or after use (there is an object-binding/linkage step at startup).  More specifically, functions may be either recursive or mutually recursive. Functions may not call subroutines.

* Functions may be defined and called either within ``mlr put`` or ``mlr put``.

* Functions have read access to ``$``-variables and ``@``-variables but may not modify them. See also :ref:`cookbook-memoization-with-oosvars` for an example.

* Argument values may be reassigned: they are not read-only.

* When a return value is not implicitly returned, this results in a return value of absent-null. (In the example above, if there were records for which the argument to ``f`` is non-numeric, the assignments would be skipped.) See also the section on :ref:`reference-null-data`.

* See the section on :ref:`reference-dsl-local-variables` for information on scope and extent of arguments, as well as for information on the use of local variables within functions.

* See the section on :ref:`reference-dsl-expressions-from-files` for information on the use of ``-f`` and ``-e`` flags.

User-defined subroutines
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Example:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint --from data/small put -q '
      begin {
        @call_count = 0;
      }
      subr s(n) {
        @call_count += 1;
        if (is_numeric(n)) {
          if (n > 1) {
            call s(n-1);
          } else {
            print "numcalls=" . @call_count;
          }
        }
      }
      print "NR=" . NR;
      call s(NR);
    '
    NR=1
    numcalls=1
    NR=2
    numcalls=3
    NR=3
    numcalls=6
    NR=4
    numcalls=10
    NR=5
    numcalls=15

Properties of user-defined subroutines:

* Subroutine bodies start with ``subr`` and a parameter list, defined outside of ``begin``, ``end``, or other ``func`` or ``subr`` blocks. (I.e. the Miller DSL has no nested subroutines.)

* A subroutine (uniqified by its name) may not be redefined. However, functions and subroutines have separate namespaces: you can define a subroutine ``log`` which does not clash with the mathematical ``log`` function.

* Subroutines may be defined either before or after use (there is an object-binding/linkage step at startup).  More specifically, subroutines may be either recursive or mutually recursive. Subroutines may call functions.

* Subroutines may be defined and called either within ``mlr put`` or ``mlr put``.

* Subroutines have read/write access to ``$``-variables and ``@``-variables.

* Argument values may be reassigned: they are not read-only.

* See the section on :ref:`reference-dsl-local-variables` for information on scope and extent of arguments, as well as for information on the use of local variables within functions.

* See the section on :ref:`reference-dsl-expressions-from-files` for information on the use of ``-f`` and ``-e`` flags.

.. _reference-dsl-errors-and-transparency:

Errors and transparency
----------------------------------------------------------------

As soon as you have a programming language, you start having the problem *What is my code doing, and why?* This includes getting syntax errors -- which are always annoying -- as well as the even more annoying problem of a program which parses without syntax error but doesn't do what you expect.

The ``syntax error`` message is cryptic: it says ``syntax error at `` followed by the next symbol it couldn't parse. This is good, but (as of 5.0.0) it doesn't say things like ``syntax error at line 17, character 22``. Here are some common causes of syntax errors:

* Don't forget ``;`` at end of line, before another statement on the next line.

* Miller's DSL lacks the ``++`` and ``--`` operators.

* Curly braces are required for the bodies of ``if``/``while``/``for`` blocks, even when the body is a single statement.

Now for transparency:

* As in any language, you can do (see :ref:`reference-dsl-print-statements`) ``print`` (or ``eprint`` to print to stderr). See also :ref:`reference-dsl-dump-statements` and :ref:`reference-dsl-emit-statements`.

* The ``-v`` option to ``mlr put`` and ``mlr filter`` prints abstract syntax trees for your code. While not all details here will be of interest to everyone, certainly this makes questions such as operator precedence completely unambiguous.

* The ``-T`` option prints a trace of each statement executed.

* The ``-t`` and ``-a`` options show low-level details for the parsing process and for stack-variable-index allocation, respectively. These will likely be of interest to people who enjoy compilers, and probably less useful for a more general audience.

* Please see :ref:`reference-dsl-type-checking` for type declarations and type-assertions you can use to make sure expressions and the data flowing them are evaluating as you expect.  I made them optional because one of Miller's important use-cases is being able to say simple things like ``mlr put '$y = $x + 1' myfile.dat`` with a minimum of punctuational bric-a-brac -- but for programs over a few lines I generally find that the more type-specification, the better.

A note on the complexity of Miller's expression language
----------------------------------------------------------------

One of Miller's strengths is its brevity: it's much quicker -- and less error-prone -- to type ``mlr stats1 -a sum -f x,y -g a,b`` than having to track summation variables as in ``awk``, or using Miller's out-of-stream variables. And the more language features Miller's put-DSL has (for-loops, if-statements, nested control structures, user-defined functions, etc.) then the *less* powerful it begins to seem: because of the other programming-language features it *doesn't* have (classes, execptions, and so on).

When I was originally prototyping Miller in 2015, the decision I had was whether to hand-code in a low-level language like C or Rust, with my own hand-rolled DSL, or whether to use a higher-level language (like Python or Lua or Nim) and let the ``put`` statements be handled by the implementation language's own ``eval``: the implementation language would take the place of a DSL. Multiple performance experiments showed me I could get better throughput using the former, by a wide margin. So Miller is Go under the hood with a hand-rolled DSL.

I do want to keep focusing on what Miller is good at -- concise notation, low latency, and high throughput -- and not add too much in terms of high-level-language features to the DSL.  That said, some sort of customizability is a basic thing to want. As of 4.1.0 we have recursive for/while/if structures on about the same complexity level as ``awk``; as of 5.0.0 we have user-defined functions and map-valued variables, again on about the same complexity level as ``awk`` along with optional type-declaration syntax.  While I'm excited by these powerful language features, I hope to keep new features beyond 5.0.0 focused on Miller's sweet spot which is speed plus simplicity.

