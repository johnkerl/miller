..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

DSL reference: control structures
================================================================

Pattern-action blocks
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

These are reminiscent of ``awk`` syntax.  They can be used to allow assignments to be done only when appropriate -- e.g. for math-function domain restrictions, regex-matching, and so on:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr cat data/put-gating-example-1.dkvp
    x=-1
    x=0
    x=1
    x=2
    x=3

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr put '$x > 0.0 { $y = log10($x); $z = sqrt($y) }' data/put-gating-example-1.dkvp
    x=-1
    x=0
    x=1,y=0,z=0
    x=2,y=0.3010299956639812,z=0.5486620049392715
    x=3,y=0.4771212547196624,z=0.6907396432228734

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr cat data/put-gating-example-2.dkvp
    a=abc_123
    a=some other name
    a=xyz_789

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr put '$a =~ "([a-z]+)_([0-9]+)" { $b = "left_\1"; $c = "right_\2" }' data/put-gating-example-2.dkvp
    a=abc_123,b=left_\1,c=right_\2
    a=some other name
    a=xyz_789,b=left_\1,c=right_\2

This produces heteregenous output which Miller, of course, has no problems with (see :doc:`record-heterogeneity`).  But if you want homogeneous output, the curly braces can be replaced with a semicolon between the expression and the body statements.  This causes ``put`` to evaluate the boolean expression (along with any side effects, namely, regex-captures ``\1``, ``\2``, etc.) but doesn't use it as a criterion for whether subsequent assignments should be executed. Instead, subsequent assignments are done unconditionally:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr put '$x > 0.0; $y = log10($x); $z = sqrt($y)' data/put-gating-example-1.dkvp
    x=1,y=0,z=0
    x=2,y=0.3010299956639812,z=0.5486620049392715
    x=3,y=0.4771212547196624,z=0.6907396432228734

.. code-block:: none
   :emphasize-lines: 1-1

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
   :emphasize-lines: 1-6

    $ echo x=1,y=2 | mlr put '
      while (NF < 10) {
        $[NF+1] = ""
      }
      $foo = "bar"
    '
    x=1,y=2,3=,4=,5=,6=,7=,8=,9=,10=,foo=bar

.. code-block:: none
   :emphasize-lines: 1-9

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
   :emphasize-lines: 1-8

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
   :emphasize-lines: 1-8

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
   :emphasize-lines: 1-1

    $ cat data/for-srec-example.tbl
    label1 label2 f1  f2  f3
    blue   green  100 240 350
    red    green  120 11  195
    yellow blue   140 0   240

.. code-block:: none
   :emphasize-lines: 1-11

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
   :emphasize-lines: 1-1

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
   :emphasize-lines: 1-10

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
   :emphasize-lines: 1-9

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
   :emphasize-lines: 1-10

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
   :emphasize-lines: 1-16

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
   :emphasize-lines: 1-17

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
   :emphasize-lines: 1-17

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
   :emphasize-lines: 1-7

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
   :emphasize-lines: 1-10

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
   :emphasize-lines: 1-5

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
   :emphasize-lines: 1-4

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
   :emphasize-lines: 1-4

    $ mlr put -q '
      @x_sum += $x;
      end { emit @x_sum }
    ' ../data/small
    x_sum=4.536293840335763

We can do similarly with multiple out-of-stream variables:

.. code-block:: none
   :emphasize-lines: 1-8

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
   :emphasize-lines: 1-1

    $ mlr stats1 -a count,sum -f x ../data/small
    x_count=10,x_sum=4.536293840335763

Note that it's a syntax error for begin/end blocks to refer to field names (beginning with ``$``), since these execute outside the context of input records.

