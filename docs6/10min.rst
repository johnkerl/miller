..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Miller in 10 minutes
====================

CSV-file examples
^^^^^^^^^^^^^^^^^

Suppose you have this CSV data file:

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat example.csv
    color,shape,flag,index,quantity,rate
    yellow,triangle,1,11,43.6498,9.8870
    red,square,1,15,79.2778,0.0130
    red,circle,1,16,13.8103,2.9010
    red,square,0,48,77.5542,7.4670
    purple,triangle,0,51,81.2290,8.5910
    red,square,0,64,77.1991,9.5310
    purple,triangle,0,65,80.1405,5.8240
    yellow,circle,1,73,63.9785,4.2370
    yellow,circle,1,87,63.5058,8.3350
    purple,square,0,91,72.3735,8.2430

``mlr cat`` is like cat -- it passes the data through unmodified:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --csv cat example.csv
    color,shape,flag,index,quantity,rate
    yellow,triangle,1,11,43.6498,9.8870
    red,square,1,15,79.2778,0.0130
    red,circle,1,16,13.8103,2.9010
    red,square,0,48,77.5542,7.4670
    purple,triangle,0,51,81.2290,8.5910
    red,square,0,64,77.1991,9.5310
    purple,triangle,0,65,80.1405,5.8240
    yellow,circle,1,73,63.9785,4.2370
    yellow,circle,1,87,63.5058,8.3350
    purple,square,0,91,72.3735,8.2430

but it can also do format conversion (here, you can pretty-print in tabular format):

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --opprint cat example.csv
    color  shape    flag index quantity rate
    yellow triangle 1    11    43.6498  9.8870
    red    square   1    15    79.2778  0.0130
    red    circle   1    16    13.8103  2.9010
    red    square   0    48    77.5542  7.4670
    purple triangle 0    51    81.2290  8.5910
    red    square   0    64    77.1991  9.5310
    purple triangle 0    65    80.1405  5.8240
    yellow circle   1    73    63.9785  4.2370
    yellow circle   1    87    63.5058  8.3350
    purple square   0    91    72.3735  8.2430

``mlr head`` and ``mlr tail`` count records rather than lines. Whether you're getting the first few records or the last few, the CSV header is included either way:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --csv head -n 4 example.csv
    color,shape,flag,index,quantity,rate
    yellow,triangle,1,11,43.6498,9.8870
    red,square,1,15,79.2778,0.0130
    red,circle,1,16,13.8103,2.9010
    red,square,0,48,77.5542,7.4670

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --csv tail -n 4 example.csv
    color,shape,flag,index,quantity,rate
    purple,triangle,0,65,80.1405,5.8240
    yellow,circle,1,73,63.9785,4.2370
    yellow,circle,1,87,63.5058,8.3350
    purple,square,0,91,72.3735,8.2430

You can sort primarily alphabetically on one field, then secondarily numerically descending on another field:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --opprint sort -f shape -nr index example.csv
    color  shape    flag index quantity rate
    yellow circle   1    87    63.5058  8.3350
    yellow circle   1    73    63.9785  4.2370
    red    circle   1    16    13.8103  2.9010
    purple square   0    91    72.3735  8.2430
    red    square   0    64    77.1991  9.5310
    red    square   0    48    77.5542  7.4670
    red    square   1    15    79.2778  0.0130
    purple triangle 0    65    80.1405  5.8240
    purple triangle 0    51    81.2290  8.5910
    yellow triangle 1    11    43.6498  9.8870

You can use ``cut`` to retain only specified fields, in the same order they appeared in the input data:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --opprint cut -f flag,shape example.csv
    shape    flag
    triangle 1
    square   1
    circle   1
    square   0
    triangle 0
    square   0
    triangle 0
    circle   1
    circle   1
    square   0

You can also use ``cut -o`` to retain only specified fields in your preferred order:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --opprint cut -o -f flag,shape example.csv
    flag shape
    1    triangle
    1    square
    1    circle
    0    square
    0    triangle
    0    square
    0    triangle
    1    circle
    1    circle
    0    square

You can use ``cut -x`` to omit fields you don't care about:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --opprint cut -x -f flag,shape example.csv
    color  index quantity rate
    yellow 11    43.6498  9.8870
    red    15    79.2778  0.0130
    red    16    13.8103  2.9010
    red    48    77.5542  7.4670
    purple 51    81.2290  8.5910
    red    64    77.1991  9.5310
    purple 65    80.1405  5.8240
    yellow 73    63.9785  4.2370
    yellow 87    63.5058  8.3350
    purple 91    72.3735  8.2430

You can use ``filter`` to keep only records you care about:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --opprint filter '$color == "red"' example.csv
    color shape  flag index quantity rate
    red   square 1    15    79.2778  0.0130
    red   circle 1    16    13.8103  2.9010
    red   square 0    48    77.5542  7.4670
    red   square 0    64    77.1991  9.5310

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --opprint filter '$color == "red" && $flag == 1' example.csv
    color shape  flag index quantity rate
    red   square 1    15    79.2778  0.0130
    red   circle 1    16    13.8103  2.9010

You can use ``put`` to create new fields which are computed from other fields:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --opprint put '$ratio = $quantity / $rate; $color_shape = $color . "_" . $shape' example.csv
    color  shape    flag index quantity rate   ratio              color_shape
    yellow triangle 1    11    43.6498  9.8870 4.414868008496004  yellow_triangle
    red    square   1    15    79.2778  0.0130 6098.292307692308  red_square
    red    circle   1    16    13.8103  2.9010 4.760530851430541  red_circle
    red    square   0    48    77.5542  7.4670 10.386259541984733 red_square
    purple triangle 0    51    81.2290  8.5910 9.455127458968688  purple_triangle
    red    square   0    64    77.1991  9.5310 8.099790158430384  red_square
    purple triangle 0    65    80.1405  5.8240 13.760388049450551 purple_triangle
    yellow circle   1    73    63.9785  4.2370 15.09995279679018  yellow_circle
    yellow circle   1    87    63.5058  8.3350 7.619172165566886  yellow_circle
    purple square   0    91    72.3735  8.2430 8.779995147397793  purple_square

Even though Miller's main selling point is name-indexing, sometimes you really want to refer to a field name by its positional index. Use ``$[[3]]`` to access the name of field 3 or ``$[[[3]]]`` to access the value of field 3:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --opprint put '$[[3]] = "NEW"' example.csv
    color  shape    NEW index quantity rate
    yellow triangle 1   11    43.6498  9.8870
    red    square   1   15    79.2778  0.0130
    red    circle   1   16    13.8103  2.9010
    red    square   0   48    77.5542  7.4670
    purple triangle 0   51    81.2290  8.5910
    red    square   0   64    77.1991  9.5310
    purple triangle 0   65    80.1405  5.8240
    yellow circle   1   73    63.9785  4.2370
    yellow circle   1   87    63.5058  8.3350
    purple square   0   91    72.3735  8.2430

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --opprint put '$[[[3]]] = "NEW"' example.csv
    color  shape    flag index quantity rate
    yellow triangle NEW  11    43.6498  9.8870
    red    square   NEW  15    79.2778  0.0130
    red    circle   NEW  16    13.8103  2.9010
    red    square   NEW  48    77.5542  7.4670
    purple triangle NEW  51    81.2290  8.5910
    red    square   NEW  64    77.1991  9.5310
    purple triangle NEW  65    80.1405  5.8240
    yellow circle   NEW  73    63.9785  4.2370
    yellow circle   NEW  87    63.5058  8.3350
    purple square   NEW  91    72.3735  8.2430

JSON-file examples
^^^^^^^^^^^^^^^^^^

OK, CSV and pretty-print are fine. But Miller can also convert between a few other formats -- let's take a look at JSON output:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --ojson put '$ratio = $quantity/$rate; $shape = toupper($shape)' example.csv
    {
      "color": "yellow",
      "shape": "TRIANGLE",
      "flag": 1,
      "index": 11,
      "quantity": 43.6498,
      "rate": 9.8870,
      "ratio": 4.414868008496004
    }
    {
      "color": "red",
      "shape": "SQUARE",
      "flag": 1,
      "index": 15,
      "quantity": 79.2778,
      "rate": 0.0130,
      "ratio": 6098.292307692308
    }
    {
      "color": "red",
      "shape": "CIRCLE",
      "flag": 1,
      "index": 16,
      "quantity": 13.8103,
      "rate": 2.9010,
      "ratio": 4.760530851430541
    }
    {
      "color": "red",
      "shape": "SQUARE",
      "flag": 0,
      "index": 48,
      "quantity": 77.5542,
      "rate": 7.4670,
      "ratio": 10.386259541984733
    }
    {
      "color": "purple",
      "shape": "TRIANGLE",
      "flag": 0,
      "index": 51,
      "quantity": 81.2290,
      "rate": 8.5910,
      "ratio": 9.455127458968688
    }
    {
      "color": "red",
      "shape": "SQUARE",
      "flag": 0,
      "index": 64,
      "quantity": 77.1991,
      "rate": 9.5310,
      "ratio": 8.099790158430384
    }
    {
      "color": "purple",
      "shape": "TRIANGLE",
      "flag": 0,
      "index": 65,
      "quantity": 80.1405,
      "rate": 5.8240,
      "ratio": 13.760388049450551
    }
    {
      "color": "yellow",
      "shape": "CIRCLE",
      "flag": 1,
      "index": 73,
      "quantity": 63.9785,
      "rate": 4.2370,
      "ratio": 15.09995279679018
    }
    {
      "color": "yellow",
      "shape": "CIRCLE",
      "flag": 1,
      "index": 87,
      "quantity": 63.5058,
      "rate": 8.3350,
      "ratio": 7.619172165566886
    }
    {
      "color": "purple",
      "shape": "SQUARE",
      "flag": 0,
      "index": 91,
      "quantity": 72.3735,
      "rate": 8.2430,
      "ratio": 8.779995147397793
    }

Sorts and stats
^^^^^^^^^^^^^^^

Now suppose you want to sort the data on a given column, *and then* take the top few in that ordering. You can use Miller's ``then`` feature to pipe commands together.

Here are the records with the top three ``index`` values:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --opprint sort -f shape -nr index then head -n 3 example.csv
    color  shape  flag index quantity rate
    yellow circle 1    87    63.5058  8.3350
    yellow circle 1    73    63.9785  4.2370
    red    circle 1    16    13.8103  2.9010

Lots of Miller commands take a ``-g`` option for group-by: here, ``head -n 1 -g shape`` outputs the first record for each distinct value of the ``shape`` field. This means we're finding the record with highest ``index`` field for each distinct ``shape`` field:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --opprint sort -f shape -nr index then head -n 1 -g shape example.csv
    color  shape    flag index quantity rate
    yellow circle   1    87    63.5058  8.3350
    purple square   0    91    72.3735  8.2430
    purple triangle 0    65    80.1405  5.8240

Statistics can be computed with or without group-by field(s):

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --opprint --from example.csv stats1 -a count,min,mean,max -f quantity -g shape
    shape    quantity_count quantity_min quantity_mean     quantity_max
    triangle 3              43.6498      68.33976666666666 81.229
    square   4              72.3735      76.60114999999999 79.2778
    circle   3              13.8103      47.0982           63.9785

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --opprint --from example.csv stats1 -a count,min,mean,max -f quantity -g shape,color
    shape    color  quantity_count quantity_min quantity_mean      quantity_max
    triangle yellow 1              43.6498      43.6498            43.6498
    square   red    3              77.1991      78.01036666666666  79.2778
    circle   red    1              13.8103      13.8103            13.8103
    triangle purple 2              80.1405      80.68475000000001  81.229
    circle   yellow 2              63.5058      63.742149999999995 63.9785
    square   purple 1              72.3735      72.3735            72.3735

If your output has a lot of columns, you can use XTAB format to line things up vertically for you instead:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --icsv --oxtab --from example.csv stats1 -a p0,p10,p25,p50,p75,p90,p99,p100 -f rate
    rate_p0   0.0130
    rate_p10  2.9010
    rate_p25  4.2370
    rate_p50  8.2430
    rate_p75  8.5910
    rate_p90  9.8870
    rate_p99  9.8870
    rate_p100 9.8870

.. _10min-choices-for-printing-to-files:

Choices for printing to files
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Often we want to print output to the screen. Miller does this by default, as we've seen in the previous examples.

Sometimes we want to print output to another file: just use **> outputfilenamegoeshere** at the end of your command:

.. code-block:: none
   :emphasize-lines: 1,1

    % mlr --icsv --opprint cat example.csv > newfile.csv
    # Output goes to the new file;
    # nothing is printed to the screen.

.. code-block:: none
   :emphasize-lines: 1,1

    % cat newfile.csv
    color  shape    flag index quantity rate
    yellow triangle 1    11    43.6498  9.8870
    red    square   1    15    79.2778  0.0130
    red    circle   1    16    13.8103  2.9010
    red    square   0    48    77.5542  7.4670
    purple triangle 0    51    81.2290  8.5910
    red    square   0    64    77.1991  9.5310
    purple triangle 0    65    80.1405  5.8240
    yellow circle   1    73    63.9785  4.2370
    yellow circle   1    87    63.5058  8.3350
    purple square   0    91    72.3735  8.2430

Other times we just want our files to be **changed in-place**: just use **mlr -I**:


.. code-block:: none
   :emphasize-lines: 1,1

    % cp example.csv newfile.txt

.. code-block:: none
   :emphasize-lines: 1,1

    % cat newfile.txt
    color,shape,flag,index,quantity,rate
    yellow,triangle,1,11,43.6498,9.8870
    red,square,1,15,79.2778,0.0130
    red,circle,1,16,13.8103,2.9010
    red,square,0,48,77.5542,7.4670
    purple,triangle,0,51,81.2290,8.5910
    red,square,0,64,77.1991,9.5310
    purple,triangle,0,65,80.1405,5.8240
    yellow,circle,1,73,63.9785,4.2370
    yellow,circle,1,87,63.5058,8.3350
    purple,square,0,91,72.3735,8.2430

.. code-block:: none
   :emphasize-lines: 1,1

    % mlr -I --icsv --opprint cat newfile.txt

.. code-block:: none
   :emphasize-lines: 1,1

    % cat newfile.txt
    color  shape    flag index quantity rate
    yellow triangle 1    11    43.6498  9.8870
    red    square   1    15    79.2778  0.0130
    red    circle   1    16    13.8103  2.9010
    red    square   0    48    77.5542  7.4670
    purple triangle 0    51    81.2290  8.5910
    red    square   0    64    77.1991  9.5310
    purple triangle 0    65    80.1405  5.8240
    yellow circle   1    73    63.9785  4.2370
    yellow circle   1    87    63.5058  8.3350
    purple square   0    91    72.3735  8.2430

Also using ``mlr -I`` you can bulk-operate on lots of files: e.g.:

.. code-block:: none
   :emphasize-lines: 1,1

    mlr -I --csv cut -x -f unwanted_column_name *.csv

If you like, you can first copy off your original data somewhere else, before doing in-place operations.

Lastly, using ``tee`` within ``put``, you can split your input data into separate files per one or more field names:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --csv --from example.csv put -q 'tee > $shape.".csv", $*'

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat circle.csv
    color,shape,flag,index,quantity,rate
    red,circle,1,16,13.8103,2.9010
    yellow,circle,1,73,63.9785,4.2370
    yellow,circle,1,87,63.5058,8.3350

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat square.csv
    color,shape,flag,index,quantity,rate
    red,square,1,15,79.2778,0.0130
    red,square,0,48,77.5542,7.4670
    red,square,0,64,77.1991,9.5310
    purple,square,0,91,72.3735,8.2430

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat triangle.csv
    color,shape,flag,index,quantity,rate
    yellow,triangle,1,11,43.6498,9.8870
    purple,triangle,0,51,81.2290,8.5910
    purple,triangle,0,65,80.1405,5.8240

Other-format examples
^^^^^^^^^^^^^^^^^^^^^

What's a CSV file, really? It's an array of rows, or *records*, each being a list of key-value pairs, or *fields*: for CSV it so happens that all the keys are shared in the header line and the values vary data line by data line.

For example, if you have:

.. code-block:: none

    shape,flag,index
    circle,1,24
    square,0,36

then that's a way of saying:

.. code-block:: none

    shape=circle,flag=1,index=24
    shape=square,flag=0,index=36

Data written this way are called **DKVP**, for *delimited key-value pairs*.

We've also already seen other ways to write the same data:


.. code-block:: none

    CSV                               PPRINT                 JSON
    shape,flag,index                  shape  flag index      [
    circle,1,24                       circle 1    24           {
    square,0,36                       square 0    36             "shape": "circle",
                                                                 "flag": 1,
                                                                 "index": 24
                                                               },
    DKVP                              XTAB                     {
    shape=circle,flag=1,index=24      shape circle               "shape": "square",
    shape=square,flag=0,index=36      flag  1                    "flag": 0,
                                      index 24                   "index": 36
                                                               }
                                      shape square           ]
                                      flag  0
                                      index 36

Anything we can do with CSV input data, we can do with any other format input data.  And you can read from one format, do any record-processing, and output to the same format as the input, or to a different output format.
