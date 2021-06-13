..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Cookbook part 1: common patterns
================================================================

Data-cleaning examples
----------------------------------------------------------------

Here are some ways to use the type-checking options as described in :ref:`reference-dsl-type-tests-and-assertions` Suppose you have the following data file, with inconsistent typing for boolean. (Also imagine that, for the sake of discussion, we have a million-line file rather than a four-line file, so we can't see it all at once and some automation is called for.)

.. code-block:: none
   :emphasize-lines: 1-1

    $ cat data/het-bool.csv
    name,reachable
    barney,false
    betty,true
    fred,true
    wilma,1

One option is to coerce everything to boolean, or integer:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --icsv --opprint put '$reachable = boolean($reachable)' data/het-bool.csv
    name   reachable
    barney false
    betty  true
    fred   true
    wilma  true

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --icsv --opprint put '$reachable = int(boolean($reachable))' data/het-bool.csv
    name   reachable
    barney 0
    betty  1
    fred   1
    wilma  1

A second option is to flag badly formatted data within the output stream:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --icsv --opprint put '$format_ok = is_string($reachable)' data/het-bool.csv
    name   reachable format_ok
    barney false     false
    betty  true      false
    fred   true      false
    wilma  1         false

Or perhaps to flag badly formatted data outside the output stream:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --icsv --opprint put 'if (!is_string($reachable)) {eprint "Malformed at NR=".NR} ' data/het-bool.csv
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

    $ mlr --csv put '$reachable = asserting_string($reachable)' data/het-bool.csv
    Miller: is_string type-assertion failed at NR=1 FNR=1 FILENAME=data/het-bool.csv

Showing differences between successive queries
----------------------------------------------------------------

Suppose you have a database query which you run at one point in time, producing the output on the left, then again later producing the output on the right:

.. code-block:: none
   :emphasize-lines: 1-1

    $ cat data/previous_counters.csv
    color,count
    red,3472
    blue,6838
    orange,694
    purple,12

.. code-block:: none
   :emphasize-lines: 1-1

    $ cat data/current_counters.csv
    color,count
    red,3467
    orange,670
    yellow,27
    blue,6944

And, suppose you want to compute the differences in the counters between adjacent keys. Since the color names aren't all in the same order, nor are they all present on both sides, we can't just paste the two files side-by-side and do some column-four-minus-column-two arithmetic.

First, rename counter columns to make them distinct:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --csv rename count,previous_count data/previous_counters.csv > data/prevtemp.csv

.. code-block:: none
   :emphasize-lines: 1-1

    $ cat data/prevtemp.csv
    color,previous_count
    red,3472
    blue,6838
    orange,694
    purple,12

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --csv rename count,current_count data/current_counters.csv > data/currtemp.csv

.. code-block:: none
   :emphasize-lines: 1-1

    $ cat data/currtemp.csv
    color,current_count
    red,3467
    orange,670
    yellow,27
    blue,6944

Then, join on the key field(s), and use unsparsify to zero-fill counters absent on one side but present on the other. Use ``--ul`` and ``--ur`` to emit unpaired records (namely, purple on the left and yellow on the right):

.. code-block:: none
   :emphasize-lines: 1-5

    $ mlr --icsv --opprint \
      join -j color --ul --ur -f data/prevtemp.csv \
      then unsparsify --fill-with 0 \
      then put '$count_delta = $current_count - $previous_count' \
      data/currtemp.csv
    color  previous_count current_count count_delta
    red    3472           3467          -5
    orange 694            670           -24
    yellow 0              27            (error)
    blue   6838           6944          106
    purple 12             0             (error)

Two-pass algorithms
----------------------------------------------------------------

Miller is a streaming record processor; commands are performed once per record. This makes Miller particularly suitable for single-pass algorithms, allowing many of its verbs to process files that are (much) larger than the amount of RAM present in your system. (Of course, Miller verbs such as ``sort``, ``tac``, etc. all must ingest and retain all input records before emitting any output records.) You can also use out-of-stream variables to perform multi-pass computations, at the price of retaining all input records in memory.

Two-pass algorithms: computation of percentages
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

For example, mapping numeric values down a column to the percentage between their min and max values is two-pass: on the first pass you find the min and max values, then on the second, map each record's value to a percentage.

.. code-block:: none
   :emphasize-lines: 1-16

    $ mlr --from data/small --opprint put -q '
      # These are executed once per record, which is the first pass.
      # The key is to use NR to index an out-of-stream variable to
      # retain all the x-field values.
      @x_min = min($x, @x_min);
      @x_max = max($x, @x_max);
      @x[NR] = $x;
    
      # The second pass is in a for-loop in an end-block.
      end {
        for (nr, x in @x) {
          @x_pct[nr] = 100 * (x - @x_min) / (@x_max - @x_min);
        }
        emit (@x, @x_pct), "NR"
      }
    '
    NR x                   x_pct
    1  0.3467901443380824  25.66194338926441
    2  0.7586799647899636  100
    3  0.20460330576630303 0
    4  0.38139939387114097 31.90823602213647
    5  0.5732889198020006  66.54054236562845

Two-pass algorithms: line-number ratios
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Similarly, finding the total record count requires first reading through all the data:

.. code-block:: none
   :emphasize-lines: 1-11

    $ mlr --opprint --from data/small put -q '
      @records[NR] = $*;
      end {
        for((I,k),v in @records) {
          @records[I]["I"] = I;
          @records[I]["N"] = NR;
          @records[I]["PCT"] = 100*I/NR
        }
        emit @records,"I"
      }
    ' then reorder -f I,N,PCT
    I N PCT     a   b   i x                   y
    1 5 (error) pan pan 1 0.3467901443380824  0.7268028627434533
    2 5 (error) eks pan 2 0.7586799647899636  0.5221511083334797
    3 5 (error) wye wye 3 0.20460330576630303 0.33831852551664776
    4 5 (error) eks wye 4 0.38139939387114097 0.13418874328430463
    5 5 (error) wye pan 5 0.5732889198020006  0.8636244699032729

Two-pass algorithms: records having max value
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

The idea is to retain records having the largest value of ``n`` in the following data:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --itsv --opprint cat data/maxrows.tsv
    a      b      n score
    purple red    5 0.743231
    blue   purple 2 0.093710
    red    purple 2 0.802103
    purple red    5 0.389055
    red    purple 2 0.880457
    orange red    2 0.540349
    purple purple 1 0.634451
    orange purple 5 0.257223
    orange purple 5 0.693499
    red    red    4 0.981355
    blue   purple 5 0.157052
    purple purple 1 0.441784
    red    purple 1 0.124912
    orange blue   1 0.921944
    blue   purple 4 0.490909
    purple red    5 0.454779
    green  purple 4 0.198278
    orange blue   5 0.705700
    red    red    3 0.940705
    purple red    5 0.072936
    orange blue   3 0.389463
    orange purple 2 0.664985
    blue   purple 1 0.371813
    red    purple 4 0.984571
    green  purple 5 0.203577
    green  purple 3 0.900873
    purple purple 0 0.965677
    blue   purple 2 0.208785
    purple purple 1 0.455077
    red    purple 4 0.477187
    blue   red    4 0.007487

Of course, the largest value of ``n`` isn't known until after all data have been read. Using an out-of-stream variable we can retain all records as they are read, then filter them at the end:

.. code-block:: none
   :emphasize-lines: 1-1

    $ cat data/maxrows.mlr
    # Retain all records
    @records[NR] = $*;
    # Track max value of n
    @maxn = max(@maxn, $n);
    
    # After all records have been read, loop through retained records
    # and print those with the max n value.
    end {
      for (nr in @records) {
        map record = @records[nr];
        if (record["n"] == @maxn) {
          emit record;
        }
      }
    }

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --itsv --opprint put -q -f data/maxrows.mlr data/maxrows.tsv
    a      b      n score
    purple red    5 0.743231
    purple red    5 0.389055
    orange purple 5 0.257223
    orange purple 5 0.693499
    blue   purple 5 0.157052
    purple red    5 0.454779
    orange blue   5 0.705700
    purple red    5 0.072936
    green  purple 5 0.203577

Feature-counting
----------------------------------------------------------------

Suppose you have some heterogeneous data like this:

.. code-block:: none

    { "qoh": 29874, "rate": 1.68, "latency": 0.02 }
    { "name": "alice", "uid": 572 }
    { "qoh": 1227, "rate": 1.01, "latency": 0.07 }
    { "qoh": 13458, "rate": 1.72, "latency": 0.04 }
    { "qoh": 56782, "rate": 1.64 }
    { "qoh": 23512, "rate": 1.71, "latency": 0.03 }
    { "qoh": 9876, "rate": 1.89, "latency": 0.08 }
    { "name": "bill", "uid": 684 }
    { "name": "chuck", "uid2": 908 }
    { "name": "dottie", "uid": 440 }
    { "qoh": 0, "rate": 0.40, "latency": 0.01 }
    { "qoh": 5438, "rate": 1.56, "latency": 0.17 }

A reasonable question to ask is, how many occurrences of each field are there? And, what percentage of total row count has each of them? Since the denominator of the percentage is not known until the end, this is a two-pass algorithm:

.. code-block:: none

    for (key in $*) {
      @key_counts[key] += 1;
    }
    @record_count += 1;
    
    end {
      for (key in @key_counts) {
          @key_fraction[key] = @key_counts[key] / @record_count
      }
      emit @record_count;
      emit @key_counts, "key";
      emit @key_fraction,"key"
    }

Then

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --json put -q -f data/feature-count.mlr data/features.json
    {
      "record_count": 12
    }
    {
      "key": "qoh",
      "key_counts": 8
    }
    {
      "key": "rate",
      "key_counts": 8
    }
    {
      "key": "latency",
      "key_counts": 7
    }
    {
      "key": "name",
      "key_counts": 4
    }
    {
      "key": "uid",
      "key_counts": 3
    }
    {
      "key": "uid2",
      "key_counts": 1
    }
    {
      "key": "qoh",
      "key_fraction": 0.6666666666666666
    }
    {
      "key": "rate",
      "key_fraction": 0.6666666666666666
    }
    {
      "key": "latency",
      "key_fraction": 0.5833333333333334
    }
    {
      "key": "name",
      "key_fraction": 0.3333333333333333
    }
    {
      "key": "uid",
      "key_fraction": 0.25
    }
    {
      "key": "uid2",
      "key_fraction": 0.08333333333333333
    }

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --ijson --opprint put -q -f data/feature-count.mlr data/features.json
    record_count
    12
    
    key     key_counts
    qoh     8
    rate    8
    latency 7
    name    4
    uid     3
    uid2    1
    
    key     key_fraction
    qoh     0.6666666666666666
    rate    0.6666666666666666
    latency 0.5833333333333334
    name    0.3333333333333333
    uid     0.25
    uid2    0.08333333333333333

Unsparsing
----------------------------------------------------------------

The previous section discussed how to fill out missing data fields within CSV with full header line -- so the list of all field names is present within the header line. Next, let's look at a related problem: we have data where each record has various key names but we want to produce rectangular output having the union of all key names.

For example, suppose you have JSON input like this:

.. code-block:: none
   :emphasize-lines: 1-1

    $ cat data/sparse.json
    {"a":1,"b":2,"v":3}
    {"u":1,"b":2}
    {"a":1,"v":2,"x":3}
    {"v":1,"w":2}

There are field names ``a``, ``b``, ``v``, ``u``, ``x``, ``w`` in the data -- but not all in every record.  Since we don't know the names of all the keys until we've read them all, this needs to be a two-pass algorithm. On the first pass, remember all the unique key names and all the records; on the second pass, loop through the records filling in absent values, then producing output. Use ``put -q`` since we don't want to produce per-record output, only emitting output in the ``end`` block:

.. code-block:: none
   :emphasize-lines: 1-1

    $ cat data/unsparsify.mlr
    # First pass:
    # Remember all unique key names:
    for (k in $*) {
      @all_keys[k] = 1;
    }
    # Remember all input records:
    @records[NR] = $*;
    
    # Second pass:
    end {
      for (nr in @records) {
        # Get the sparsely keyed input record:
        irecord = @records[nr];
        # Fill in missing keys with empty string:
        map orecord = {};
        for (k in @all_keys) {
          if (haskey(irecord, k)) {
            orecord[k] = irecord[k];
          } else {
            orecord[k] = "";
          }
        }
        # Produce the output:
        emit orecord;
      }
    }

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --json put -q -f data/unsparsify.mlr data/sparse.json
    {
      "a": 1,
      "b": 2,
      "v": 3,
      "u": "",
      "x": "",
      "w": ""
    }
    {
      "a": "",
      "b": 2,
      "v": "",
      "u": 1,
      "x": "",
      "w": ""
    }
    {
      "a": 1,
      "b": "",
      "v": 2,
      "u": "",
      "x": 3,
      "w": ""
    }
    {
      "a": "",
      "b": "",
      "v": 1,
      "u": "",
      "x": "",
      "w": 2
    }

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --ijson --ocsv put -q -f data/unsparsify.mlr data/sparse.json
    a,b,v,u,x,w
    1,2,3,,,
    ,2,,1,,
    1,,2,,3,
    ,,1,,,2

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --ijson --opprint put -q -f data/unsparsify.mlr data/sparse.json
    a b v u x w
    1 2 3 - - -
    - 2 - 1 - -
    1 - 2 - 3 -
    - - 1 - - 2

There is a keystroke-saving verb for this: :ref:`mlr unsparsify <reference-verbs-unsparsify>`.

Parsing log-file output
----------------------------------------------------------------

This, of course, depends highly on what's in your log files. But, as an example, suppose you have log-file lines such as

.. code-block:: none

    2015-10-08 08:29:09,445 INFO com.company.path.to.ClassName @ [sometext] various/sorts/of data {& punctuation} hits=1 status=0 time=2.378

I prefer to pre-filter with ``grep`` and/or ``sed`` to extract the structured text, then hand that to Miller. Example:

.. code-block:: none
   :emphasize-lines: 1,1

    grep 'various sorts' *.log | sed 's/.*} //' | mlr --fs space --repifs --oxtab stats1 -a min,p10,p50,p90,max -f time -g status

.. _cookbook-memoization-with-oosvars:

Memoization with out-of-stream variables
----------------------------------------------------------------

The recursive function for the Fibonacci sequence is famous for its computational complexity.  Namely, using f(0)=1, f(1)=1, f(n)=f(n-1)+f(n-2) for n>=2, the evaluation tree branches left as well as right at each non-trivial level, resulting in millions or more paths to the root 0/1 nodes for larger n. This program

.. code-block:: none

    mlr --ofmt '%.9lf' --opprint seqgen --start 1 --stop 28 then put '
      func f(n) {
          @fcount += 1;              # count number of calls to the function
          if (n < 2) {
              return 1
          } else {
              return f(n-1) + f(n-2) # recurse
          }
      }
    
      @fcount = 0;
      $o = f($i);
      $fcount = @fcount;
    
    ' then put '$seconds=systime()' then step -a delta -f seconds then cut -x -f seconds

produces output like this:

.. code-block:: none

    i  o      fcount  seconds_delta
    1  1      1       0
    2  2      3       0.000039101
    3  3      5       0.000015974
    4  5      9       0.000019073
    5  8      15      0.000026941
    6  13     25      0.000036955
    7  21     41      0.000056028
    8  34     67      0.000086069
    9  55     109     0.000134945
    10 89     177     0.000217915
    11 144    287     0.000355959
    12 233    465     0.000506163
    13 377    753     0.000811815
    14 610    1219    0.001297235
    15 987    1973    0.001960993
    16 1597   3193    0.003417969
    17 2584   5167    0.006215811
    18 4181   8361    0.008294106
    19 6765   13529   0.012095928
    20 10946  21891   0.019592047
    21 17711  35421   0.031193972
    22 28657  57313   0.057254076
    23 46368  92735   0.080307961
    24 75025  150049  0.129482031
    25 121393 242785  0.213325977
    26 196418 392835  0.334423065
    27 317811 635621  0.605969906
    28 514229 1028457 0.971235037

Note that the time it takes to evaluate the function is blowing up exponentially as the input argument increases. Using ``@``-variables, which persist across records, we can cache and reuse the results of previous computations:

.. code-block:: none

    mlr --ofmt '%.9lf' --opprint seqgen --start 1 --stop 28 then put '
      func f(n) {
        @fcount += 1;                 # count number of calls to the function
        if (is_present(@fcache[n])) { # cache hit
          return @fcache[n]
        } else {                      # cache miss
          num rv = 1;
          if (n >= 2) {
            rv = f(n-1) + f(n-2)      # recurse
          }
          @fcache[n] = rv;
          return rv
        }
      }
      @fcount = 0;
      $o = f($i);
      $fcount = @fcount;
    ' then put '$seconds=systime()' then step -a delta -f seconds then cut -x -f seconds

with output like this:

.. code-block:: none

    i  o      fcount seconds_delta
    1  1      1      0
    2  2      3      0.000053883
    3  3      3      0.000035048
    4  5      3      0.000045061
    5  8      3      0.000014067
    6  13     3      0.000028849
    7  21     3      0.000028133
    8  34     3      0.000027895
    9  55     3      0.000014067
    10 89     3      0.000015020
    11 144    3      0.000012875
    12 233    3      0.000033140
    13 377    3      0.000014067
    14 610    3      0.000012875
    15 987    3      0.000029087
    16 1597   3      0.000013828
    17 2584   3      0.000013113
    18 4181   3      0.000012875
    19 6765   3      0.000013113
    20 10946  3      0.000012875
    21 17711  3      0.000013113
    22 28657  3      0.000013113
    23 46368  3      0.000015974
    24 75025  3      0.000012875
    25 121393 3      0.000013113
    26 196418 3      0.000012875
    27 317811 3      0.000013113
    28 514229 3      0.000012875
