..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Dates and times
===============

How can I filter by date?
----------------------------------------------------------------

Given input like

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat dates.csv
    date,event
    2018-02-03,initialization
    2018-03-07,discovery
    2018-02-03,allocation

we can use ``strptime`` to parse the date field into seconds-since-epoch and then do numeric comparisons.  Simply match your input dataset's date-formatting to the :ref:`reference-dsl-strptime` format-string.  For example:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --csv filter 'strptime($date, "%Y-%m-%d") > strptime("2018-03-03", "%Y-%m-%d")' dates.csv
    date,event
    2018-03-07,discovery

Caveat: localtime-handling in timezones with DST is still a work in progress; see https://github.com/johnkerl/miller/issues/170. See also https://github.com/johnkerl/miller/issues/208 -- thanks @aborruso!

Finding missing dates
----------------------------------------------------------------

Suppose you have some date-stamped data which may (or may not) be missing entries for one or more dates:

.. code-block:: none
   :emphasize-lines: 1,1

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

.. code-block:: none
   :emphasize-lines: 1,1

    $ wc -l data/miss-date.csv
        1372 data/miss-date.csv

Since there are 1372 lines in the data file, some automation is called for. To find the missing dates, you can convert the dates to seconds since the epoch using ``strptime``, then compute adjacent differences (the ``cat -n`` simply inserts record-counters):

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --from data/miss-date.csv --icsv \
      cat -n \
      then put '$datestamp = strptime($date, "%Y-%m-%d")' \
      then step -a delta -f datestamp \
    | head
    n=1,date=2012-03-05,qoh=10055,datestamp=1.3309056e+09,datestamp_delta=0
    n=2,date=2012-03-06,qoh=10486,datestamp=1.330992e+09,datestamp_delta=86400
    n=3,date=2012-03-07,qoh=10430,datestamp=1.3310784e+09,datestamp_delta=86400
    n=4,date=2012-03-08,qoh=10674,datestamp=1.3311648e+09,datestamp_delta=86400
    n=5,date=2012-03-09,qoh=10880,datestamp=1.3312512e+09,datestamp_delta=86400
    n=6,date=2012-03-10,qoh=10718,datestamp=1.3313376e+09,datestamp_delta=86400
    n=7,date=2012-03-11,qoh=10795,datestamp=1.331424e+09,datestamp_delta=86400
    n=8,date=2012-03-12,qoh=11043,datestamp=1.3315104e+09,datestamp_delta=86400
    n=9,date=2012-03-13,qoh=11177,datestamp=1.3315968e+09,datestamp_delta=86400
    n=10,date=2012-03-14,qoh=11498,datestamp=1.3316832e+09,datestamp_delta=86400

Then, filter for adjacent difference not being 86400 (the number of seconds in a day):

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --from data/miss-date.csv --icsv \
      cat -n \
      then put '$datestamp = strptime($date, "%Y-%m-%d")' \
      then step -a delta -f datestamp \
      then filter '$datestamp_delta != 86400 && $n != 1'
    n=774,date=2014-04-19,qoh=130140,datestamp=1.3978656e+09,datestamp_delta=259200
    n=1119,date=2015-03-31,qoh=181625,datestamp=1.42776e+09,datestamp_delta=172800

Given this, it's now easy to see where the gaps are:

.. code-block:: none
   :emphasize-lines: 1,1

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

.. code-block:: none
   :emphasize-lines: 1,1

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
