..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

DSL reference: syntax
================================================================

Expression formatting
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Multiple expressions may be given, separated by semicolons, and each may refer to the ones before:

.. code-block:: none
   :emphasize-lines: 1-1

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
   :emphasize-lines: 1-7

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
   :emphasize-lines: 1-1

    $ mlr --opprint filter '($x > 0.5 && $y < 0.5) || ($x < 0.5 && $y > 0.5)' then stats2 -a corr -f x,y data/medium
    x_y_corr
    -0.7479940285189345

.. _reference-dsl-expressions-from-files:

Expressions from files
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

The simplest way to enter expressions for ``put`` and ``filter`` is between single quotes on the command line, e.g.

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --from data/small put '$xy = sqrt($x**2 + $y**2)'
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533,xy=0.8052985815845617
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797,xy=0.9209978658539777
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776,xy=0.3953756915115773
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463,xy=0.40431685157744135
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,xy=1.036584492737304

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --from data/small put 'func f(a, b) { return sqrt(a**2 + b**2) } $xy = f($x, $y)'
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533,xy=0.8052985815845617
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797,xy=0.9209978658539777
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776,xy=0.3953756915115773
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463,xy=0.40431685157744135
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,xy=1.036584492737304

You may, though, find it convenient to put expressions into files for reuse, and read them
**using the -f option**. For example:

.. code-block:: none
   :emphasize-lines: 1-1

    $ cat data/fe-example-3.mlr
    func f(a, b) {
      return sqrt(a**2 + b**2)
    }
    $xy = f($x, $y)

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --from data/small put -f data/fe-example-3.mlr
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533,xy=0.8052985815845617
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797,xy=0.9209978658539777
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776,xy=0.3953756915115773
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463,xy=0.40431685157744135
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,xy=1.036584492737304

If you have some of the logic in a file and you want to write the rest on the command line, you can **use the -f and -e options together**:

.. code-block:: none
   :emphasize-lines: 1-1

    $ cat data/fe-example-4.mlr
    func f(a, b) {
      return sqrt(a**2 + b**2)
    }

.. code-block:: none
   :emphasize-lines: 1-1

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
   :emphasize-lines: 1-1

    $ echo x=1,y=2 | mlr put 'while (NF < 10) { $[NF+1] = ""}  $foo = "bar"'
    x=1,y=2,3=,4=,5=,6=,7=,8=,9=,10=,foo=bar

.. code-block:: none
   :emphasize-lines: 1-1

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
   :emphasize-lines: 1-17

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

