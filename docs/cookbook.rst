..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Cookbook part 1: common patterns
================================================================

Headerless CSV on input or output
----------------------------------------------------------------

Sometimes we get CSV files which lack a header. For example:

::

    $ cat data/headerless.csv
    John,23,present
    Fred,34,present
    Alice,56,missing
    Carol,45,present

You can use Miller to add a header. The ``--implicit-csv-header`` applies positionally indexed labels:

::

    $ mlr --csv --implicit-csv-header cat data/headerless.csv
    1,2,3
    John,23,present
    Fred,34,present
    Alice,56,missing
    Carol,45,present

Following that, you can rename the positionally indexed labels to names with meaning for your context.  For example:

::

    $ mlr --csv --implicit-csv-header label name,age,status data/headerless.csv
    name,age,status
    John,23,present
    Fred,34,present
    Alice,56,missing
    Carol,45,present

Likewise, if you need to produce CSV which is lacking its header, you can pipe Miller's output to the system command ``sed 1d``, or you can use Miller's ``--headerless-csv-output`` option:

::

    $ head -5 data/colored-shapes.dkvp | mlr --ocsv cat
    color,shape,flag,i,u,v,w,x
    yellow,triangle,1,11,0.6321695890307647,0.9887207810889004,0.4364983936735774,5.7981881667050565
    red,square,1,15,0.21966833570651523,0.001257332190235938,0.7927778364718627,2.944117399716207
    red,circle,1,16,0.20901671281497636,0.29005231936593445,0.13810280912907674,5.065034003400998
    red,square,0,48,0.9562743938458542,0.7467203085342884,0.7755423050923582,7.117831369597269
    purple,triangle,0,51,0.4355354501763202,0.8591292672156728,0.8122903963006748,5.753094629505863

::

    $ head -5 data/colored-shapes.dkvp | mlr --ocsv --headerless-csv-output cat
    yellow,triangle,1,11,0.6321695890307647,0.9887207810889004,0.4364983936735774,5.7981881667050565
    red,square,1,15,0.21966833570651523,0.001257332190235938,0.7927778364718627,2.944117399716207
    red,circle,1,16,0.20901671281497636,0.29005231936593445,0.13810280912907674,5.065034003400998
    red,square,0,48,0.9562743938458542,0.7467203085342884,0.7755423050923582,7.117831369597269
    purple,triangle,0,51,0.4355354501763202,0.8591292672156728,0.8122903963006748,5.753094629505863

Lastly, often we say "CSV" or "TSV" when we have positionally indexed data in columns which are separated by commas or tabs, respectively. In this case it's perhaps simpler to **just use NIDX format** which was designed for this purpose. (See also :doc:`file-formats`.) For example: 

::

    $ mlr --inidx --ifs comma --oxtab cut -f 1,3 data/headerless.csv
    1 John
    3 present
    
    1 Fred
    3 present
    
    1 Alice
    3 missing
    
    1 Carol
    3 present

Doing multiple joins
----------------------------------------------------------------

Suppose we have the following data:

::

    $ cat multi-join/input.csv
    id,task
    10,chop
    20,puree
    20,wash
    30,fold
    10,bake
    20,mix
    10,knead
    30,clean

And we want to augment the ``id`` column with lookups from the following data files:

::

    $ cat multi-join/name-lookup.csv
    id,name
    30,Alice
    10,Bob
    20,Carol

::

    $ cat multi-join/status-lookup.csv
    id,status
    30,occupied
    10,idle
    20,idle

We can run the input file through multiple ``join`` commands in a ``then``-chain:

::

    $ mlr --icsv --opprint join -f multi-join/name-lookup.csv -j id then join -f multi-join/status-lookup.csv -j id multi-join/input.csv
    id status   name  task
    10 idle     Bob   chop
    20 idle     Carol puree
    20 idle     Carol wash
    30 occupied Alice fold
    10 idle     Bob   bake
    20 idle     Carol mix
    10 idle     Bob   knead
    30 occupied Alice clean

Bulk rename of fields
----------------------------------------------------------------

Suppose you want to replace spaces with underscores in your column names:

::

    $ cat data/spaces.csv
    a b c,def,g h i
    123,4567,890
    2468,1357,3579
    9987,3312,4543

The simplest way is to use ``mlr rename`` with ``-g`` (for global replace, not just first occurrence of space within each field) and ``-r`` for pattern-matching (rather than explicit single-column renames):

::

    $ mlr --csv rename -g -r ' ,_'  data/spaces.csv
    a_b_c,def,g_h_i
    123,4567,890
    2468,1357,3579
    9987,3312,4543

::

    $ mlr --csv --opprint rename -g -r ' ,_'  data/spaces.csv
    a_b_c def  g_h_i
    123   4567 890
    2468  1357 3579
    9987  3312 4543

You can also do this with a for-loop:

::

    $ cat data/bulk-rename-for-loop.mlr
    map newrec = {};
    for (oldk, v in $*) {
        newrec[gsub(oldk, " ", "_")] = v;
    }
    $* = newrec

::

    $ mlr --icsv --opprint put -f data/bulk-rename-for-loop.mlr data/spaces.csv
    a_b_c def  g_h_i
    123   4567 890
    2468  1357 3579
    9987  3312 4543

Search-and-replace over all fields
----------------------------------------------------------------

How to do ``$name = gsub($name, "old", "new")`` for all fields?

::

    $ cat data/sar.csv
    a,b,c
    the quick,brown fox,jumped
    over,the,lazy dogs

::

    $ cat data/sar.mlr
      for (k in $*) {
        $[k] = gsub($[k], "e", "X");
      }

::

    $ mlr --csv put -f data/sar.mlr data/sar.csv
    a,b,c
    thX quick,brown fox,jumpXd
    ovXr,thX,lazy dogs

Full field renames and reassigns
----------------------------------------------------------------

Using Miller 5.0.0's map literals and assigning to ``$*``, you can fully generalize :ref:`mlr rename <reference-verbs-rename>`, :ref:`mlr reorder <reference-verbs-reorder>`, etc.

::

    $ cat data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729

::

    $ mlr put '
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
    z=0.346790,KEYFIELD=pan,i=1,b=pan,y=0.346790,x=0.726803
    z=0.758680,KEYFIELD=eks,i=3,b=pan,y=0.758680,x=0.522151
    z=0.204603,KEYFIELD=wye,i=6,b=wye,y=0.204603,x=0.338319
    z=0.381399,KEYFIELD=eks,i=10,b=wye,y=0.381399,x=0.134189
    z=0.573289,KEYFIELD=wye,i=15,b=pan,y=0.573289,x=0.863624

Numbering and renumbering records
----------------------------------------------------------------

The ``awk``-like built-in variable ``NR`` is incremented for each input record:

::

    $ cat data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729

::

    $ mlr put '$nr = NR' data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533,nr=1
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797,nr=2
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776,nr=3
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463,nr=4
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,nr=5

However, this is the record number within the original input stream -- not after any filtering you may have done:

::

    $ mlr filter '$a == "wye"' then put '$nr = NR' data/small
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776,nr=3
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,nr=5

There are two good options here. One is to use the ``cat`` verb with ``-n``:

::

    $ mlr filter '$a == "wye"' then cat -n data/small
    n=1,a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    n=2,a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729

The other is to keep your own counter within the ``put`` DSL:

::

    $ mlr filter '$a == "wye"' then put 'begin {@n = 1} $n = @n; @n += 1' data/small
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776,n=1
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,n=2

The difference is a matter of taste (although ``mlr cat -n`` puts the counter first).

Options for dealing with duplicate rows
----------------------------------------------------------------

If your data has records appearing multiple times, you can use :ref:`mlr uniq <reference-verbs-uniq>` to show and/or count the unique records.

If you want to look at partial uniqueness -- for example, show only the first record for each unique combination of the ``account_id`` and ``account_status`` fields -- you might use ``mlr head -n 1 -g account_id,account_status``. Please also see :ref:`mlr head <reference-verbs-head>`.

.. _cookbook-data-cleaning-examples:

Data-cleaning examples
----------------------------------------------------------------

Here are some ways to use the type-checking options as described in :ref:`reference-dsl-type-tests-and-assertions` Suppose you have the following data file, with inconsistent typing for boolean. (Also imagine that, for the sake of discussion, we have a million-line file rather than a four-line file, so we can't see it all at once and some automation is called for.) 

::

    $ cat data/het-bool.csv
    name,reachable
    barney,false
    betty,true
    fred,true
    wilma,1

One option is to coerce everything to boolean, or integer:

::

    $ mlr --icsv --opprint put '$reachable = boolean($reachable)' data/het-bool.csv
    name   reachable
    barney false
    betty  true
    fred   true
    wilma  true

::

    $ mlr --icsv --opprint put '$reachable = int(boolean($reachable))' data/het-bool.csv
    name   reachable
    barney 0
    betty  1
    fred   1
    wilma  1

A second option is to flag badly formatted data within the output stream:

::

    $ mlr --icsv --opprint put '$format_ok = is_string($reachable)' data/het-bool.csv
    name   reachable format_ok
    barney false     true
    betty  true      true
    fred   true      true
    wilma  1         false

Or perhaps to flag badly formatted data outside the output stream:

::

    $ mlr --icsv --opprint put 'if (!is_string($reachable)) {eprint "Malformed at NR=".NR} ' data/het-bool.csv
    Malformed at NR=4
    name   reachable
    barney false
    betty  true
    fred   true
    wilma  1

A third way is to abort the process on first instance of bad data:

::

    $ mlr --csv put '$reachable = asserting_string($reachable)' data/het-bool.csv
    mlr: string type-assertion failed at NR=4 FNR=4 FILENAME=data/het-bool.csv
    name,reachable
    barney,false
    betty,true
    fred,true

Splitting nested fields
----------------------------------------------------------------

Suppose you have a TSV file like this:

::

    a	b
    x	z
    s	u:v:w

The simplest option is to use :ref:`mlr nest <reference-verbs-nest>`:

::

    $ mlr --tsv nest --explode --values --across-records -f b --nested-fs : data/nested.tsv
    a	b
    x	z
    s	u
    s	v
    s	w

::

    $ mlr --tsv nest --explode --values --across-fields  -f b --nested-fs : data/nested.tsv
    a	b_1
    x	z
    
    a	b_1	b_2	b_3
    s	u	v	w

While ``mlr nest`` is simplest, let's also take a look at a few ways to do this using the ``put`` DSL.

One option to split out the colon-delimited values in the ``b`` column is to use ``splitnv`` to create an integer-indexed map and loop over it, adding new fields to the current record: 

::

    $ mlr --from data/nested.tsv --itsv --oxtab put 'o=splitnv($b, ":"); for (k,v in o) {$["p".k]=v}'
    a  x
    b  z
    p1 z
    
    a  s
    b  u:v:w
    p1 u
    p2 v
    p3 w

while another is to loop over the same map from ``splitnv`` and use it (with ``put -q`` to suppress printing the original record) to produce multiple records:

::

    $ mlr --from data/nested.tsv --itsv --oxtab put -q 'o=splitnv($b, ":"); for (k,v in o) {emit mapsum($*, {"b":v})}'
    a x
    b z
    
    a s
    b u
    
    a s
    b v
    
    a s
    b w

::

    $ mlr --from data/nested.tsv --tsv put -q 'o=splitnv($b, ":"); for (k,v in o) {emit mapsum($*, {"b":v})}'
    a	b
    x	z
    s	u
    s	v
    s	w

Showing differences between successive queries
----------------------------------------------------------------

Suppose you have a database query which you run at one point in time, producing the output on the left, then again later producing the output on the right:

::

    $ cat data/previous_counters.csv
    color,count
    red,3472
    blue,6838
    orange,694
    purple,12

::

    $ cat data/current_counters.csv
    color,count
    red,3467
    orange,670
    yellow,27
    blue,6944

And, suppose you want to compute the differences in the counters between adjacent keys. Since the color names aren't all in the same order, nor are they all present on both sides, we can't just paste the two files side-by-side and do some column-four-minus-column-two arithmetic.

First, rename counter columns to make them distinct:

::

    $ mlr --csv rename count,previous_count data/previous_counters.csv > data/prevtemp.csv

::

    $ cat data/prevtemp.csv
    color,previous_count
    red,3472
    blue,6838
    orange,694
    purple,12

::

    $ mlr --csv rename count,current_count data/current_counters.csv > data/currtemp.csv

::

    $ cat data/currtemp.csv
    color,current_count
    red,3467
    orange,670
    yellow,27
    blue,6944

Then, join on the key field(s), and use unsparsify to zero-fill counters absent on one side but present on the other. Use ``--ul`` and ``--ur`` to emit unpaired records (namely, purple on the left and yellow on the right): 

::

    $ mlr --icsv --opprint \
      join -j color --ul --ur -f data/prevtemp.csv \
      then unsparsify --fill-with 0 \
      then put '$count_delta = $current_count - $previous_count' \
      data/currtemp.csv
    color  previous_count current_count count_delta
    red    3472           3467          -5
    orange 694            670           -24
    yellow 0              27            27
    blue   6838           6944          106
    purple 12             0             -12

Finding missing dates
----------------------------------------------------------------

Suppose you have some date-stamped data which may (or may not) be missing entries for one or more dates:

::

    $ head -n 10 data/miss-date.csv
    date,qoh
    2012-03-05,10055
    2012-03-06,10486
    2012-03-07,10430
    2012-03-08,10674
    2012-03-09,10880
    2012-03-10,10718
    2012-03-11,10795
    2012-03-12,11043
    2012-03-13,11177

::

    $ wc -l data/miss-date.csv
        1372 data/miss-date.csv

Since there are 1372 lines in the data file, some automation is called for. To find the missing dates, you can convert the dates to seconds since the epoch using ``strptime``, then compute adjacent differences (the ``cat -n`` simply inserts record-counters): 

::

    $ mlr --from data/miss-date.csv --icsv \
      cat -n \
      then put '$datestamp = strptime($date, "%Y-%m-%d")' \
      then step -a delta -f datestamp \
    | head
    n=1,date=2012-03-05,qoh=10055,datestamp=1330905600.000000,datestamp_delta=0
    n=2,date=2012-03-06,qoh=10486,datestamp=1330992000.000000,datestamp_delta=86400.000000
    n=3,date=2012-03-07,qoh=10430,datestamp=1331078400.000000,datestamp_delta=86400.000000
    n=4,date=2012-03-08,qoh=10674,datestamp=1331164800.000000,datestamp_delta=86400.000000
    n=5,date=2012-03-09,qoh=10880,datestamp=1331251200.000000,datestamp_delta=86400.000000
    n=6,date=2012-03-10,qoh=10718,datestamp=1331337600.000000,datestamp_delta=86400.000000
    n=7,date=2012-03-11,qoh=10795,datestamp=1331424000.000000,datestamp_delta=86400.000000
    n=8,date=2012-03-12,qoh=11043,datestamp=1331510400.000000,datestamp_delta=86400.000000
    n=9,date=2012-03-13,qoh=11177,datestamp=1331596800.000000,datestamp_delta=86400.000000
    n=10,date=2012-03-14,qoh=11498,datestamp=1331683200.000000,datestamp_delta=86400.000000

Then, filter for adjacent difference not being 86400 (the number of seconds in a day):

::

    $ mlr --from data/miss-date.csv --icsv \
      cat -n \
      then put '$datestamp = strptime($date, "%Y-%m-%d")' \
      then step -a delta -f datestamp \
      then filter '$datestamp_delta != 86400 && $n != 1'
    n=774,date=2014-04-19,qoh=130140,datestamp=1397865600.000000,datestamp_delta=259200.000000
    n=1119,date=2015-03-31,qoh=181625,datestamp=1427760000.000000,datestamp_delta=172800.000000

Given this, it's now easy to see where the gaps are:

::

    $ mlr cat -n then filter '$n >= 770 && $n <= 780' data/miss-date.csv
    n=770,1=2014-04-12,2=129435
    n=771,1=2014-04-13,2=129868
    n=772,1=2014-04-14,2=129797
    n=773,1=2014-04-15,2=129919
    n=774,1=2014-04-16,2=130181
    n=775,1=2014-04-19,2=130140
    n=776,1=2014-04-20,2=130271
    n=777,1=2014-04-21,2=130368
    n=778,1=2014-04-22,2=130368
    n=779,1=2014-04-23,2=130849
    n=780,1=2014-04-24,2=131026

::

    $ mlr cat -n then filter '$n >= 1115 && $n <= 1125' data/miss-date.csv
    n=1115,1=2015-03-25,2=181006
    n=1116,1=2015-03-26,2=180995
    n=1117,1=2015-03-27,2=181043
    n=1118,1=2015-03-28,2=181112
    n=1119,1=2015-03-29,2=181306
    n=1120,1=2015-03-31,2=181625
    n=1121,1=2015-04-01,2=181494
    n=1122,1=2015-04-02,2=181718
    n=1123,1=2015-04-03,2=181835
    n=1124,1=2015-04-04,2=182104
    n=1125,1=2015-04-05,2=182528

Two-pass algorithms
----------------------------------------------------------------

Miller is a streaming record processor; commands are performed once per record. This makes Miller particularly suitable for single-pass algorithms, allowing many of its verbs to process files that are (much) larger than the amount of RAM present in your system. (Of course, Miller verbs such as ``sort``, ``tac``, etc. all must ingest and retain all input records before emitting any output records.) You can also use out-of-stream variables to perform multi-pass computations, at the price of retaining all input records in memory. 

Two-pass algorithms: computation of percentages
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

For example, mapping numeric values down a column to the percentage between their min and max values is two-pass: on the first pass you find the min and max values, then on the second, map each record's value to a percentage. 

::

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
    NR x        x_pct
    1  0.346790 25.661943
    2  0.758680 100.000000
    3  0.204603 0.000000
    4  0.381399 31.908236
    5  0.573289 66.540542

Two-pass algorithms: line-number ratios
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Similarly, finding the total record count requires first reading through all the data:

::

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
    I N PCT a   b   i x                   y
    1 5 20  pan pan 1 0.3467901443380824  0.7268028627434533
    2 5 40  eks pan 2 0.7586799647899636  0.5221511083334797
    3 5 60  wye wye 3 0.20460330576630303 0.33831852551664776
    4 5 80  eks wye 4 0.38139939387114097 0.13418874328430463
    5 5 100 wye pan 5 0.5732889198020006  0.8636244699032729

Two-pass algorithms: records having max value
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

The idea is to retain records having the largest value of ``n`` in the following data:

::

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

::

    $ cat data/maxrows.mlr
    # Retain all records
    @records[NR] = $*;
    # Track max value of n
    @maxn = max(@maxn, $n);
    
    # After all records have been read, loop through retained records
    # and print those with the max n value.
    end {
      for (int nr in @records) {
        map record = @records[nr];
        if (record["n"] == @maxn) {
          emit record;
        }
      }
    }

::

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

Rectangularizing data
----------------------------------------------------------------

Suppose you have a method (in whatever language) which is printing things of the form

::

    outer=1
    outer=2
    outer=3

and then calls another method which prints things of the form

::

    middle=10
    middle=11
    middle=12
    middle=20
    middle=21
    middle=30
    middle=31

and then, perhaps, that second method calls a third method which prints things of the form

::

    inner1=100,inner2=101
    inner1=120,inner2=121
    inner1=200,inner2=201
    inner1=210,inner2=211
    inner1=300,inner2=301
    inner1=312
    inner1=313,inner2=314

with the result that your program's output is

::

    outer=1
    middle=10
    inner1=100,inner2=101
    middle=11
    middle=12
    inner1=120,inner2=121
    outer=2
    middle=20
    inner1=200,inner2=201
    middle=21
    inner1=210,inner2=211
    outer=3
    middle=30
    inner1=300,inner2=301
    middle=31
    inner1=312
    inner1=313,inner2=314

The idea here is that middles starting with a 1 belong to the outer value of 1, and so on.  (For example, the outer values might be account IDs, the middle values might be invoice IDs, and the inner values might be invoice line-items.) If you want all the middle and inner lines to have the context of which outers they belong to, you can modify your software to pass all those through your methods. Alternatively, don't refactor your code just to handle some ad-hoc log-data formatting -- instead, use the following to rectangularize the data.  The idea is to use an out-of-stream variable to accumulate fields across records. Clear that variable when you see an outer ID; accumulate fields; emit output when you see the inner IDs. 

::

    $ mlr --from data/rect.txt put -q '
      is_present($outer) {
        unset @r
      }
      for (k, v in $*) {
        @r[k] = v
      }
      is_present($inner1) {
        emit @r
      }'
    outer=1,middle=10,inner1=100,inner2=101
    outer=1,middle=12,inner1=120,inner2=121
    outer=2,middle=20,inner1=200,inner2=201
    outer=2,middle=21,inner1=210,inner2=211
    outer=3,middle=30,inner1=300,inner2=301
    outer=3,middle=31,inner1=312,inner2=301
    outer=3,middle=31,inner1=313,inner2=314

Regularizing ragged CSV
----------------------------------------------------------------

Miller handles compliant CSV: in particular, it's an error if the number of data fields in a given data line don't match the number of header lines. But in the event that you have a CSV file in which some lines have less than the full number of fields, you can use Miller to pad them out. The trick is to use NIDX format, for which each line stands on its own without respect to a header line. 

::

    $ cat data/ragged.csv
    a,b,c
    1,2,3
    4,5
    6,7,8,9

::

    $ mlr --from data/ragged.csv --fs comma --nidx put '
      @maxnf = max(@maxnf, NF);
      @nf = NF;
      while(@nf < @maxnf) {
        @nf += 1;
        $[@nf] = ""
      }
    '
    a,b,c
    1,2,3
    4,5,
    6,7,8,9

or, more simply,

::

    $ mlr --from data/ragged.csv --fs comma --nidx put '
      @maxnf = max(@maxnf, NF);
      while(NF < @maxnf) {
        $[NF+1] = "";
      }
    '
    a,b,c
    1,2,3
    4,5,
    6,7,8,9

Feature-counting
----------------------------------------------------------------

Suppose you have some heterogeneous data like this:

::

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

::

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

::

    $ mlr --json put -q -f data/feature-count.mlr data/features.json
    { "record_count": 12 }
    { "key": "qoh", "key_counts": 8 }
    { "key": "rate", "key_counts": 8 }
    { "key": "latency", "key_counts": 7 }
    { "key": "name", "key_counts": 4 }
    { "key": "uid", "key_counts": 3 }
    { "key": "uid2", "key_counts": 1 }
    { "key": "qoh", "key_fraction": 0.666667 }
    { "key": "rate", "key_fraction": 0.666667 }
    { "key": "latency", "key_fraction": 0.583333 }
    { "key": "name", "key_fraction": 0.333333 }
    { "key": "uid", "key_fraction": 0.250000 }
    { "key": "uid2", "key_fraction": 0.083333 }

::

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
    qoh     0.666667
    rate    0.666667
    latency 0.583333
    name    0.333333
    uid     0.250000
    uid2    0.083333

Unsparsing
----------------------------------------------------------------

The previous section discussed how to fill out missing data fields within CSV with full header line -- so the list of all field names is present within the header line. Next, let's look at a related problem: we have data where each record has various key names but we want to produce rectangular output having the union of all key names. 

For example, suppose you have JSON input like this:

::

    $ cat data/sparse.json
    {"a":1,"b":2,"v":3}
    {"u":1,"b":2}
    {"a":1,"v":2,"x":3}
    {"v":1,"w":2}

There are field names ``a``, ``b``, ``v``, ``u``, ``x``, ``w`` in the data -- but not all in every record.  Since we don't know the names of all the keys until we've read them all, this needs to be a two-pass algorithm. On the first pass, remember all the unique key names and all the records; on the second pass, loop through the records filling in absent values, then producing output. Use ``put -q`` since we don't want to produce per-record output, only emitting output in the ``end`` block: 

::

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

::

    $ mlr --json put -q -f data/unsparsify.mlr data/sparse.json
    { "a": 1, "b": 2, "v": 3, "u": "", "x": "", "w": "" }
    { "a": "", "b": 2, "v": "", "u": 1, "x": "", "w": "" }
    { "a": 1, "b": "", "v": 2, "u": "", "x": 3, "w": "" }
    { "a": "", "b": "", "v": 1, "u": "", "x": "", "w": 2 }

::

    $ mlr --ijson --ocsv put -q -f data/unsparsify.mlr data/sparse.json
    a,b,v,u,x,w
    1,2,3,,,
    ,2,,1,,
    1,,2,,3,
    ,,1,,,2

::

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

::

    2015-10-08 08:29:09,445 INFO com.company.path.to.ClassName @ [sometext] various/sorts/of data {& punctuation} hits=1 status=0 time=2.378

I prefer to pre-filter with ``grep`` and/or ``sed`` to extract the structured text, then hand that to Miller. Example:

::

    grep 'various sorts' *.log | sed 's/.*} //' | mlr --fs space --repifs --oxtab stats1 -a min,p10,p50,p90,max -f time -g status

Memoization with out-of-stream variables
----------------------------------------------------------------

The recursive function for the Fibonacci sequence is famous for its computational complexity.  Namely, using *f*(0)=1, *f*(1)=1, *f*(*n*)=*f*(*n*-1)+*f*(*n*-2) for *n*&ge;2, the evaluation tree branches left as well as right at each non-trivial level, resulting in millions or more paths to the root 0/1 nodes for larger *n*. This program

::

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

::

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

::

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

::

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
