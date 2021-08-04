..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Data-cleaning examples
================================================================

Here are some ways to use the type-checking options as described in :ref:`reference-dsl-type-tests-and-assertions` Suppose you have the following data file, with inconsistent typing for boolean. (Also imagine that, for the sake of discussion, we have a million-line file rather than a four-line file, so we can't see it all at once and some automation is called for.)

.. code-block:: none
   :emphasize-lines: 1-1

    cat data/het-bool.csv
    name,reachable
    barney,false
    betty,true
    fred,true
    wilma,1

One option is to coerce everything to boolean, or integer:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --opprint put '$reachable = boolean($reachable)' data/het-bool.csv
    name   reachable
    barney false
    betty  true
    fred   true
    wilma  true

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --opprint put '$reachable = int(boolean($reachable))' data/het-bool.csv
    name   reachable
    barney 0
    betty  1
    fred   1
    wilma  1

A second option is to flag badly formatted data within the output stream:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --opprint put '$format_ok = is_string($reachable)' data/het-bool.csv
    name   reachable format_ok
    barney false     false
    betty  true      false
    fred   true      false
    wilma  1         false

Or perhaps to flag badly formatted data outside the output stream:

.. code-block:: none
   :emphasize-lines: 1-3

    mlr --icsv --opprint put '
      if (!is_string($reachable)) {eprint "Malformed at NR=".NR}
    ' data/het-bool.csv
    Malformed at NR=1
    Malformed at NR=2
    Malformed at NR=3
    Malformed at NR=4
    name   reachable
    barney false
    betty  true
    fred   true
    wilma  1

A third way is to abort the process on first instance of bad data:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --csv put '$reachable = asserting_string($reachable)' data/het-bool.csv
    Miller: is_string type-assertion failed at NR=1 FNR=1 FILENAME=data/het-bool.csv
