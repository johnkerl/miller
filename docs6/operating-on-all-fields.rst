..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Operating on all fields
=======================

Bulk rename of fields
----------------------------------------------------------------

Suppose you want to replace spaces with underscores in your column names:

.. code-block:: none
   :emphasize-lines: 1-1

    cat data/spaces.csv
    a b c,def,g h i
    123,4567,890
    2468,1357,3579
    9987,3312,4543

The simplest way is to use ``mlr rename`` with ``-g`` (for global replace, not just first occurrence of space within each field) and ``-r`` for pattern-matching (rather than explicit single-column renames):

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --csv rename -g -r ' ,_'  data/spaces.csv
    a_b_c,def,g_h_i
    123,4567,890
    2468,1357,3579
    9987,3312,4543

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --csv --opprint rename -g -r ' ,_'  data/spaces.csv
    a_b_c def  g_h_i
    123   4567 890
    2468  1357 3579
    9987  3312 4543

You can also do this with a for-loop:

.. code-block:: none
   :emphasize-lines: 1-1

    cat data/bulk-rename-for-loop.mlr
    map newrec = {};
    for (oldk, v in $*) {
        newrec[gsub(oldk, " ", "_")] = v;
    }
    $* = newrec

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --opprint put -f data/bulk-rename-for-loop.mlr data/spaces.csv
    a_b_c def  g_h_i
    123   4567 890
    2468  1357 3579
    9987  3312 4543

Search-and-replace over all fields
----------------------------------------------------------------

How to do ``$name = gsub($name, "old", "new")`` for all fields?

.. code-block:: none
   :emphasize-lines: 1-1

    cat data/sar.csv
    a,b,c
    the quick,brown fox,jumped
    over,the,lazy dogs

.. code-block:: none
   :emphasize-lines: 1-1

    cat data/sar.mlr
      for (k in $*) {
        $[k] = gsub($[k], "e", "X");
      }

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --csv put -f data/sar.mlr data/sar.csv
    a,b,c
    thX quick,brown fox,jumpXd
    ovXr,thX,lazy dogs

Full field renames and reassigns
----------------------------------------------------------------

Using Miller 5.0.0's map literals and assigning to ``$*``, you can fully generalize :ref:`mlr rename <reference-verbs-rename>`, :ref:`mlr reorder <reference-verbs-reorder>`, etc.

.. code-block:: none
   :emphasize-lines: 1-1

    cat data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729

.. code-block:: none
   :emphasize-lines: 1-15

    mlr put '
      begin {
        @i_cumu = 0;
      }
    
      @i_cumu += $i;
      $* = {
        "z": $x + y,
        "KEYFIELD": $a,
        "i": @i_cumu,
        "b": $b,
        "y": $x,
        "x": $y,
      };
    ' data/small
    z=0.3467901443380824,KEYFIELD=pan,i=1,b=pan,y=0.3467901443380824,x=0.7268028627434533
    z=0.7586799647899636,KEYFIELD=eks,i=3,b=pan,y=0.7586799647899636,x=0.5221511083334797
    z=0.20460330576630303,KEYFIELD=wye,i=6,b=wye,y=0.20460330576630303,x=0.33831852551664776
    z=0.38139939387114097,KEYFIELD=eks,i=10,b=wye,y=0.38139939387114097,x=0.13418874328430463
    z=0.5732889198020006,KEYFIELD=wye,i=15,b=pan,y=0.5732889198020006,x=0.8636244699032729
