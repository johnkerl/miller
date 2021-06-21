..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Miller in 10 minutes
====================

Obtaining Miller
^^^^^^^^^^^^^^^^

You can install Miller for various platforms as follows:

* Linux: ``yum install miller`` or ``apt-get install miller`` depending on your flavor of Linux
* MacOS: ``brew install miller`` or ``port install miller`` depending on your preference of `Homebrew <https://brew.sh>`_ or `MacPorts <https://macports.org>`_.
* Windows: ``choco install miller``  using `Chocolatey <https://chocolatey.org>`_.
* You can get latest builds for Linux, MacOS, and Windows by visiting https://github.com/johnkerl/miller/actions, selecting the latest build, and clicking _Artifacts_. (These are retained for 5 days after each commit.)
* See also :doc:`build` if you prefer -- in particular, if your platform's package manager doesn't have the latest release.

As a first check, you should be able to run ``mlr --version`` at your system's command prompt and see something like the following:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --version
    Miller v6.0.0-dev

As a second check, given (`example.csv <./example.csv>`_) you should be able to do

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --csv cat example.csv
    color,shape,flag,index,quantity,rate
    yellow,triangle,true,11,43.6498,9.8870
    red,square,true,15,79.2778,0.0130
    red,circle,true,16,13.8103,2.9010
    red,square,false,48,77.5542,7.4670
    purple,triangle,false,51,81.2290,8.5910
    red,square,false,64,77.1991,9.5310
    purple,triangle,false,65,80.1405,5.8240
    yellow,circle,true,73,63.9785,4.2370
    yellow,circle,true,87,63.5058,8.3350
    purple,square,false,91,72.3735,8.2430

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --opprint cat example.csv
    color  shape    flag  index quantity rate
    yellow triangle true  11    43.6498  9.8870
    red    square   true  15    79.2778  0.0130
    red    circle   true  16    13.8103  2.9010
    red    square   false 48    77.5542  7.4670
    purple triangle false 51    81.2290  8.5910
    red    square   false 64    77.1991  9.5310
    purple triangle false 65    80.1405  5.8240
    yellow circle   true  73    63.9785  4.2370
    yellow circle   true  87    63.5058  8.3350
    purple square   false 91    72.3735  8.2430

If you run into issues on these checks, please check out the resources on the :doc:`community` page for help.

Miller verbs
^^^^^^^^^^^^

Let's take a quick look at some of the most useful Miller verbs -- file-format-aware, name-index-empowered equivalents of standard system commands.

``mlr cat`` is like system ``cat`` (or ``type`` on Windows) -- it passes the data through unmodified:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --csv cat example.csv
    color,shape,flag,index,quantity,rate
    yellow,triangle,true,11,43.6498,9.8870
    red,square,true,15,79.2778,0.0130
    red,circle,true,16,13.8103,2.9010
    red,square,false,48,77.5542,7.4670
    purple,triangle,false,51,81.2290,8.5910
    red,square,false,64,77.1991,9.5310
    purple,triangle,false,65,80.1405,5.8240
    yellow,circle,true,73,63.9785,4.2370
    yellow,circle,true,87,63.5058,8.3350
    purple,square,false,91,72.3735,8.2430

But ``mlr cat`` can also do format conversion -- for example, you can pretty-print in tabular format:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --opprint cat example.csv
    color  shape    flag  index quantity rate
    yellow triangle true  11    43.6498  9.8870
    red    square   true  15    79.2778  0.0130
    red    circle   true  16    13.8103  2.9010
    red    square   false 48    77.5542  7.4670
    purple triangle false 51    81.2290  8.5910
    red    square   false 64    77.1991  9.5310
    purple triangle false 65    80.1405  5.8240
    yellow circle   true  73    63.9785  4.2370
    yellow circle   true  87    63.5058  8.3350
    purple square   false 91    72.3735  8.2430

``mlr head`` and ``mlr tail`` count records rather than lines. Whether you're getting the first few records or the last few, the CSV header is included either way:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --csv head -n 4 example.csv
    color,shape,flag,index,quantity,rate
    yellow,triangle,true,11,43.6498,9.8870
    red,square,true,15,79.2778,0.0130
    red,circle,true,16,13.8103,2.9010
    red,square,false,48,77.5542,7.4670

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --csv tail -n 4 example.csv
    color,shape,flag,index,quantity,rate
    purple,triangle,false,65,80.1405,5.8240
    yellow,circle,true,73,63.9785,4.2370
    yellow,circle,true,87,63.5058,8.3350
    purple,square,false,91,72.3735,8.2430

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --ojson tail -n 2 example.csv
    {
      "color": "yellow",
      "shape": "circle",
      "flag": true,
      "index": 87,
      "quantity": 63.5058,
      "rate": 8.3350
    }
    {
      "color": "purple",
      "shape": "square",
      "flag": false,
      "index": 91,
      "quantity": 72.3735,
      "rate": 8.2430
    }

You can sort on a single field:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --opprint sort -f shape example.csv
    color  shape    flag  index quantity rate
    red    circle   true  16    13.8103  2.9010
    yellow circle   true  73    63.9785  4.2370
    yellow circle   true  87    63.5058  8.3350
    red    square   true  15    79.2778  0.0130
    red    square   false 48    77.5542  7.4670
    red    square   false 64    77.1991  9.5310
    purple square   false 91    72.3735  8.2430
    yellow triangle true  11    43.6498  9.8870
    purple triangle false 51    81.2290  8.5910
    purple triangle false 65    80.1405  5.8240

Or, you can sort primarily alphabetically on one field, then secondarily numerically descending on another field, and so on:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --opprint sort -f shape -nr index example.csv
    color  shape    flag  index quantity rate
    yellow circle   true  87    63.5058  8.3350
    yellow circle   true  73    63.9785  4.2370
    red    circle   true  16    13.8103  2.9010
    purple square   false 91    72.3735  8.2430
    red    square   false 64    77.1991  9.5310
    red    square   false 48    77.5542  7.4670
    red    square   true  15    79.2778  0.0130
    purple triangle false 65    80.1405  5.8240
    purple triangle false 51    81.2290  8.5910
    yellow triangle true  11    43.6498  9.8870

If there are fields you don't want to see in your data, you can use ``cut`` to keep only the ones you want, in the same order they appeared in the input data:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --opprint cut -f flag,shape example.csv
    shape    flag
    triangle true
    square   true
    circle   true
    square   false
    triangle false
    square   false
    triangle false
    circle   true
    circle   true
    square   false

You can also use ``cut -o`` to keep specified fields, but in your preferred order:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --opprint cut -o -f flag,shape example.csv
    flag  shape
    true  triangle
    true  square
    true  circle
    false square
    false triangle
    false square
    false triangle
    true  circle
    true  circle
    false square

You can use ``cut -x`` to omit fields you don't care about:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --opprint cut -x -f flag,shape example.csv
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
   :emphasize-lines: 1-1

    mlr --icsv --opprint filter '$color == "red"' example.csv
    color shape  flag  index quantity rate
    red   square true  15    79.2778  0.0130
    red   circle true  16    13.8103  2.9010
    red   square false 48    77.5542  7.4670
    red   square false 64    77.1991  9.5310

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --opprint filter '$color == "red" && $flag == true' example.csv
    color shape  flag index quantity rate
    red   square true 15    79.2778  0.0130
    red   circle true 16    13.8103  2.9010

You can use ``put`` to create new fields which are computed from other fields:

.. code-block:: none
   :emphasize-lines: 1-4

    mlr --icsv --opprint put '
      $ratio = $quantity / $rate;
      $color_shape = $color . "_" . $shape
    ' example.csv
    color  shape    flag  index quantity rate   ratio              color_shape
    yellow triangle true  11    43.6498  9.8870 4.414868008496004  yellow_triangle
    red    square   true  15    79.2778  0.0130 6098.292307692308  red_square
    red    circle   true  16    13.8103  2.9010 4.760530851430541  red_circle
    red    square   false 48    77.5542  7.4670 10.386259541984733 red_square
    purple triangle false 51    81.2290  8.5910 9.455127458968688  purple_triangle
    red    square   false 64    77.1991  9.5310 8.099790158430384  red_square
    purple triangle false 65    80.1405  5.8240 13.760388049450551 purple_triangle
    yellow circle   true  73    63.9785  4.2370 15.09995279679018  yellow_circle
    yellow circle   true  87    63.5058  8.3350 7.619172165566886  yellow_circle
    purple square   false 91    72.3735  8.2430 8.779995147397793  purple_square

Even though Miller's main selling point is name-indexing, sometimes you really want to refer to a field name by its positional index. Use ``$[[3]]`` to access the name of field 3 or ``$[[[3]]]`` to access the value of field 3:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --opprint put '$[[3]] = "NEW"' example.csv
    color  shape    NEW   index quantity rate
    yellow triangle true  11    43.6498  9.8870
    red    square   true  15    79.2778  0.0130
    red    circle   true  16    13.8103  2.9010
    red    square   false 48    77.5542  7.4670
    purple triangle false 51    81.2290  8.5910
    red    square   false 64    77.1991  9.5310
    purple triangle false 65    80.1405  5.8240
    yellow circle   true  73    63.9785  4.2370
    yellow circle   true  87    63.5058  8.3350
    purple square   false 91    72.3735  8.2430

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --opprint put '$[[[3]]] = "NEW"' example.csv
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

You can find the full list of verbs at the :doc:`reference-verbs` page.

Multiple input files
^^^^^^^^^^^^^^^^^^^^

Miller takes all the files from the command line as an input stream. But it's format-aware, so it doesn't repeat CSV header lines. For example, with input files (`data/a.csv <data/a.csv>`_) and (`data/b.csv <data/b.csv>`_), the system ``cat`` command will repeat header lines:

.. code-block:: none
   :emphasize-lines: 1-1

    cat data/a.csv
    a,b,c
    1,2,3
    4,5,6

.. code-block:: none
   :emphasize-lines: 1-1

    cat data/b.csv
    a,b,c
    7,8,9

.. code-block:: none
   :emphasize-lines: 1-1

    cat data/a.csv data/b.csv
    a,b,c
    1,2,3
    4,5,6
    a,b,c
    7,8,9

However, ``mlr cat`` will not:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --csv cat data/a.csv data/b.csv
    a,b,c
    1,2,3
    4,5,6
    7,8,9

Chaining verbs together
^^^^^^^^^^^^^^^^^^^^^^^

Often we want to chain queries together -- for example, sorting by a field and taking the top few values. We can do this using pipes:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --csv sort -nr index example.csv | mlr --icsv --opprint head -n 3
    color  shape  flag  index quantity rate
    purple square false 91    72.3735  8.2430
    yellow circle true  87    63.5058  8.3350
    yellow circle true  73    63.9785  4.2370

This works fine -- but Miller also lets you chain verbs together using the word ``then``. Think of this as a Miller-internal pipe that lets you use fewer keystrokes:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --opprint sort -nr index then head -n 3 example.csv
    color  shape  flag  index quantity rate
    purple square false 91    72.3735  8.2430
    yellow circle true  87    63.5058  8.3350
    yellow circle true  73    63.9785  4.2370

As another convenience, you can put the filename first using ``--from``. When you're interacting with your data at the command line, this makes it easier to up-arrow and append to the previous command:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --opprint --from example.csv sort -nr index then head -n 3
    color  shape  flag  index quantity rate
    purple square false 91    72.3735  8.2430
    yellow circle true  87    63.5058  8.3350
    yellow circle true  73    63.9785  4.2370

.. code-block:: none
   :emphasize-lines: 1-4

    mlr --icsv --opprint --from example.csv \
      sort -nr index \
      then head -n 3 \
      then cut -f shape,quantity
    shape  quantity
    square 72.3735
    circle 63.5058
    circle 63.9785

Sorts and stats
^^^^^^^^^^^^^^^

Now suppose you want to sort the data on a given column, *and then* take the top few in that ordering. You can use Miller's ``then`` feature to pipe commands together.

Here are the records with the top three ``index`` values:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --opprint sort -nr index then head -n 3 example.csv
    color  shape  flag  index quantity rate
    purple square false 91    72.3735  8.2430
    yellow circle true  87    63.5058  8.3350
    yellow circle true  73    63.9785  4.2370

Lots of Miller commands take a ``-g`` option for group-by: here, ``head -n 1 -g shape`` outputs the first record for each distinct value of the ``shape`` field. This means we're finding the record with highest ``index`` field for each distinct ``shape`` field:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --opprint sort -f shape -nr index then head -n 1 -g shape example.csv
    color  shape    flag  index quantity rate
    yellow circle   true  87    63.5058  8.3350
    purple square   false 91    72.3735  8.2430
    purple triangle false 65    80.1405  5.8240

Statistics can be computed with or without group-by field(s):

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --opprint --from example.csv stats1 -a count,min,mean,max -f quantity -g shape
    shape    quantity_count quantity_min quantity_mean     quantity_max
    triangle 3              43.6498      68.33976666666666 81.229
    square   4              72.3735      76.60114999999999 79.2778
    circle   3              13.8103      47.0982           63.9785

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --opprint --from example.csv stats1 -a count,min,mean,max -f quantity -g shape,color
    shape    color  quantity_count quantity_min quantity_mean      quantity_max
    triangle yellow 1              43.6498      43.6498            43.6498
    square   red    3              77.1991      78.01036666666666  79.2778
    circle   red    1              13.8103      13.8103            13.8103
    triangle purple 2              80.1405      80.68475000000001  81.229
    circle   yellow 2              63.5058      63.742149999999995 63.9785
    square   purple 1              72.3735      72.3735            72.3735

If your output has a lot of columns, you can use XTAB format to line things up vertically for you instead:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --icsv --oxtab --from example.csv stats1 -a p0,p10,p25,p50,p75,p90,p99,p100 -f rate
    rate_p0   0.0130
    rate_p10  2.9010
    rate_p25  4.2370
    rate_p50  8.2430
    rate_p75  8.5910
    rate_p90  9.8870
    rate_p99  9.8870
    rate_p100 9.8870


File formats and format conversion
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Miller supports the following formats:

* CSV (comma-separared values)
* TSV (tab-separated values)
* JSON (JavaScript Object Notation)
* PPRINT (pretty-printed tabular)
* XTAB (vertical-tabular or sideways-tabular)
* NIDX (numerically indexed, label-free, with implicit labels ``"1"``, ``"2"``, etc.)
* DKVP (delimited key-value pairs).

What's a CSV file, really? It's an array of rows, or *records*, each being a list of key-value pairs, or *fields*: for CSV it so happens that all the keys are shared in the header line and the values vary from one data line to another.

For example, if you have:

.. code-block:: none

    shape,flag,index
    circle,1,24
    square,0,36

then that's a way of saying:

.. code-block:: none

    shape=circle,flag=1,index=24
    shape=square,flag=0,index=36

Other ways to write the same data:

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

How to specify these to Miller:

* If you use ``--csv`` or ``--json`` or ``--pprint``, etc., then Miller will use that format for input and output.
* If you use ``--icsv`` and ``--ojson`` (note the extra ``i`` and ``o``) then Miller will use CSV for input and JSON for output, etc.  See also :doc:`keystroke-savers` for even shorter options like ``--c2j``.

You can read more about this at the :doc:`file-formats` page.

Choices for printing to files
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Often we want to print output to the screen. Miller does this by default, as we've seen in the previous examples.

Sometimes, though, we want to print output to another file. Just use **> outputfilenamegoeshere** at the end of your command:

.. code-block:: none
   :emphasize-lines: 1,1

    mlr --icsv --opprint cat example.csv > newfile.csv
    # Output goes to the new file;
    # nothing is printed to the screen.

.. code-block:: none
   :emphasize-lines: 1,1

    cat newfile.csv
    color  shape    flag     index quantity rate
    yellow triangle true     11    43.6498  9.8870
    red    square   true     15    79.2778  0.0130
    red    circle   true     16    13.8103  2.9010
    red    square   false    48    77.5542  7.4670
    purple triangle false    51    81.2290  8.5910
    red    square   false    64    77.1991  9.5310
    purple triangle false    65    80.1405  5.8240
    yellow circle   true     73    63.9785  4.2370
    yellow circle   true     87    63.5058  8.3350
    purple square   false    91    72.3735  8.2430

Other times we just want our files to be **changed in-place**: just use **mlr -I**:

.. code-block:: none
   :emphasize-lines: 1,1

    cp example.csv newfile.txt

.. code-block:: none
   :emphasize-lines: 1,1

    cat newfile.txt
    color,shape,flag,index,quantity,rate
    yellow,triangle,true,11,43.6498,9.8870
    red,square,true,15,79.2778,0.0130
    red,circle,true,16,13.8103,2.9010
    red,square,false,48,77.5542,7.4670
    purple,triangle,false,51,81.2290,8.5910
    red,square,false,64,77.1991,9.5310
    purple,triangle,false,65,80.1405,5.8240
    yellow,circle,true,73,63.9785,4.2370
    yellow,circle,true,87,63.5058,8.3350
    purple,square,false,91,72.3735,8.2430

.. code-block:: none
   :emphasize-lines: 1,1

    mlr -I --csv sort -f shape newfile.txt

.. code-block:: none
   :emphasize-lines: 1,1

    cat newfile.txt
    color,shape,flag,index,quantity,rate
    red,circle,true,16,13.8103,2.9010
    yellow,circle,true,73,63.9785,4.2370
    yellow,circle,true,87,63.5058,8.3350
    red,square,true,15,79.2778,0.0130
    red,square,false,48,77.5542,7.4670
    red,square,false,64,77.1991,9.5310
    purple,square,false,91,72.3735,8.2430
    yellow,triangle,true,11,43.6498,9.8870
    purple,triangle,false,51,81.2290,8.5910
    purple,triangle,false,65,80.1405,5.8240

Also using ``mlr -I`` you can bulk-operate on lots of files: e.g.:

.. code-block:: none
   :emphasize-lines: 1,1

    mlr -I --csv cut -x -f unwanted_column_name *.csv

If you like, you can first copy off your original data somewhere else, before doing in-place operations.

Lastly, using ``tee`` within ``put``, you can split your input data into separate files per one or more field names:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --csv --from example.csv put -q 'tee > $shape.".csv", $*'

.. code-block:: none
   :emphasize-lines: 1-1

    cat circle.csv
    color,shape,flag,index,quantity,rate
    red,circle,true,16,13.8103,2.9010
    yellow,circle,true,73,63.9785,4.2370
    yellow,circle,true,87,63.5058,8.3350

.. code-block:: none
   :emphasize-lines: 1-1

    cat square.csv
    color,shape,flag,index,quantity,rate
    red,square,true,15,79.2778,0.0130
    red,square,false,48,77.5542,7.4670
    red,square,false,64,77.1991,9.5310
    purple,square,false,91,72.3735,8.2430

.. code-block:: none
   :emphasize-lines: 1-1

    cat triangle.csv
    color,shape,flag,index,quantity,rate
    yellow,triangle,true,11,43.6498,9.8870
    purple,triangle,false,51,81.2290,8.5910
    purple,triangle,false,65,80.1405,5.8240
