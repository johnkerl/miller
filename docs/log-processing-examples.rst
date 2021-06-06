..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Log-processing examples
----------------------------------------------------------------

Another of my favorite use-cases for Miller is doing ad-hoc processing of log-file data.  Here's where DKVP format really shines: one, since the field names and field values are present on every line, every line stands on its own. That means you can ``grep`` or what have you. Also it means not every line needs to have the same list of field names ("schema").

Again, all the examples in the CSV section apply here -- just change the input-format flags. But there's more you can do when not all the records have the same shape.

Writing a program -- in any language whatsoever -- you can have it print out log lines as it goes along, with items for various events jumbled together. After the program has finished running you can sort it all out, filter it, analyze it, and learn from it.

Suppose your program has printed something like this (`log.txt <./log.txt>`_):

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat log.txt
    op=enter,time=1472819681
    op=cache,type=A9,hit=0
    op=cache,type=A4,hit=1
    time=1472819690,batch_size=100,num_filtered=237
    op=cache,type=A1,hit=1
    op=cache,type=A9,hit=0
    op=cache,type=A1,hit=1
    op=cache,type=A9,hit=0
    op=cache,type=A9,hit=0
    op=cache,type=A1,hit=1
    time=1472819705,batch_size=100,num_filtered=348
    op=cache,type=A4,hit=1
    op=cache,type=A9,hit=0
    op=cache,type=A9,hit=0
    op=cache,type=A9,hit=0
    op=cache,type=A9,hit=0
    op=cache,type=A4,hit=1
    time=1472819713,batch_size=100,num_filtered=493
    op=cache,type=A9,hit=1
    op=cache,type=A1,hit=1
    op=cache,type=A9,hit=0
    op=cache,type=A9,hit=0
    op=cache,type=A9,hit=1
    time=1472819720,batch_size=100,num_filtered=554
    op=cache,type=A1,hit=0
    op=cache,type=A4,hit=1
    op=cache,type=A9,hit=0
    op=cache,type=A9,hit=0
    op=cache,type=A9,hit=0
    op=cache,type=A4,hit=0
    op=cache,type=A4,hit=0
    op=cache,type=A9,hit=0
    time=1472819736,batch_size=100,num_filtered=612
    op=cache,type=A1,hit=1
    op=cache,type=A9,hit=0
    op=cache,type=A9,hit=0
    op=cache,type=A9,hit=0
    op=cache,type=A9,hit=0
    op=cache,type=A4,hit=1
    op=cache,type=A1,hit=1
    op=cache,type=A9,hit=0
    op=cache,type=A9,hit=0
    time=1472819742,batch_size=100,num_filtered=728

Each print statement simply contains local information: the current timestamp, whether a particular cache was hit or not, etc. Then using either the system ``grep`` command, or Miller's ``having-fields``, or ``is_present``, we can pick out the parts we want and analyze them:

.. code-block:: none
   :emphasize-lines: 1,1

    $ grep op=cache log.txt \
      | mlr --idkvp --opprint stats1 -a mean -f hit -g type then sort -f type
    type hit_mean
    A1   0.857143
    A4   0.714286
    A9   0.090909

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --from log.txt --opprint \
      filter 'is_present($batch_size)' \
      then step -a delta -f time,num_filtered \
      then sec2gmt time
    time                 batch_size num_filtered time_delta num_filtered_delta
    2016-09-02T12:34:50Z 100        237          0          0
    2016-09-02T12:35:05Z 100        348          15         111
    2016-09-02T12:35:13Z 100        493          8          145
    2016-09-02T12:35:20Z 100        554          7          61
    2016-09-02T12:35:36Z 100        612          16         58
    2016-09-02T12:35:42Z 100        728          6          116

Alternatively, we can simply group the similar data for a better look:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint group-like log.txt
    op    time
    enter 1472819681
    
    op    type hit
    cache A9   0
    cache A4   1
    cache A1   1
    cache A9   0
    cache A1   1
    cache A9   0
    cache A9   0
    cache A1   1
    cache A4   1
    cache A9   0
    cache A9   0
    cache A9   0
    cache A9   0
    cache A4   1
    cache A9   1
    cache A1   1
    cache A9   0
    cache A9   0
    cache A9   1
    cache A1   0
    cache A4   1
    cache A9   0
    cache A9   0
    cache A9   0
    cache A4   0
    cache A4   0
    cache A9   0
    cache A1   1
    cache A9   0
    cache A9   0
    cache A9   0
    cache A9   0
    cache A4   1
    cache A1   1
    cache A9   0
    cache A9   0
    
    time       batch_size num_filtered
    1472819690 100        237
    1472819705 100        348
    1472819713 100        493
    1472819720 100        554
    1472819736 100        612
    1472819742 100        728

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint group-like then sec2gmt time log.txt
    op    time
    enter 2016-09-02T12:34:41Z
    
    op    type hit
    cache A9   0
    cache A4   1
    cache A1   1
    cache A9   0
    cache A1   1
    cache A9   0
    cache A9   0
    cache A1   1
    cache A4   1
    cache A9   0
    cache A9   0
    cache A9   0
    cache A9   0
    cache A4   1
    cache A9   1
    cache A1   1
    cache A9   0
    cache A9   0
    cache A9   1
    cache A1   0
    cache A4   1
    cache A9   0
    cache A9   0
    cache A9   0
    cache A4   0
    cache A4   0
    cache A9   0
    cache A1   1
    cache A9   0
    cache A9   0
    cache A9   0
    cache A9   0
    cache A4   1
    cache A1   1
    cache A9   0
    cache A9   0
    
    time                 batch_size num_filtered
    2016-09-02T12:34:50Z 100        237
    2016-09-02T12:35:05Z 100        348
    2016-09-02T12:35:13Z 100        493
    2016-09-02T12:35:20Z 100        554
    2016-09-02T12:35:36Z 100        612
    2016-09-02T12:35:42Z 100        728
