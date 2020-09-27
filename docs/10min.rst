Miller in 10 minutes
====================

CSV-file examples
^^^^^^^^^^^^^^^^^

Suppose you have this CSV data file::

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

``mlr cat`` is like cat -- it passes the data through unmodified::

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

but it can also do format conversion (here, you can pretty-print in tabular format)::

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

``mlr head`` and ``mlr tail`` count records rather than lines. Whethere you're getting the first few records or the last few, the CSV header is included either way::

    $ mlr --csv head -n 4 example.csv
    color,shape,flag,index,quantity,rate
    yellow,triangle,1,11,43.6498,9.8870
    red,square,1,15,79.2778,0.0130
    red,circle,1,16,13.8103,2.9010
    red,square,0,48,77.5542,7.4670
    $ mlr --csv tail -n 4 example.csv
    color,shape,flag,index,quantity,rate
    purple,triangle,0,65,80.1405,5.8240
    yellow,circle,1,73,63.9785,4.2370
    yellow,circle,1,87,63.5058,8.3350
    purple,square,0,91,72.3735,8.2430

You can sort primarily alphabetically on one field, then secondarily numerically descending on another field::

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

You can use ``cut`` to retain only specified fields, in the same order they appeared in the input data::

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

You can also use ``cut -o`` to retain only specified fields in your preferred order::

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

You can use ``cut -x`` to omit fields you don't care about::

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

You can use ``filter`` to keep only records you care about::

    $ mlr --icsv --opprint filter '$color == "red"' example.csv
    color shape  flag index quantity rate
    red   square 1    15    79.2778  0.0130
    red   circle 1    16    13.8103  2.9010
    red   square 0    48    77.5542  7.4670
    red   square 0    64    77.1991  9.5310
    $ mlr --icsv --opprint filter '$color == "red" && $flag == 1' example.csv
    color shape  flag index quantity rate
    red   square 1    15    79.2778  0.0130
    red   circle 1    16    13.8103  2.9010

You can use ``put`` to create new fields which are computed from other fields::

    $ mlr --icsv --opprint put '$ratio = $quantity / $rate; $color_shape = $color . "_" . $shape' example.csv
    color  shape    flag index quantity rate   ratio       color_shape
    yellow triangle 1    11    43.6498  9.8870 4.414868    yellow_triangle
    red    square   1    15    79.2778  0.0130 6098.292308 red_square
    red    circle   1    16    13.8103  2.9010 4.760531    red_circle
    red    square   0    48    77.5542  7.4670 10.386260   red_square
    purple triangle 0    51    81.2290  8.5910 9.455127    purple_triangle
    red    square   0    64    77.1991  9.5310 8.099790    red_square
    purple triangle 0    65    80.1405  5.8240 13.760388   purple_triangle
    yellow circle   1    73    63.9785  4.2370 15.099953   yellow_circle
    yellow circle   1    87    63.5058  8.3350 7.619172    yellow_circle
    purple square   0    91    72.3735  8.2430 8.779995    purple_square

Even though Miller's main selling point is name-indexing, sometimes you really want to refer to a field name by its positional index. Use ``$[[3]]`` to access the name of field 3 or ``$[[[3]]]`` to access the value of field 3::

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

OK, CSV and pretty-print are fine. But Miller can also convert between a few other formats -- let's take a look at JSON output::

    $ mlr --icsv --ojson put '$ratio = $quantity/$rate; $shape = toupper($shape)' example.csv
    { "color": "yellow", "shape": "TRIANGLE", "flag": 1, "index": 11, "quantity": 43.6498, "rate": 9.8870, "ratio": 4.414868 }
    { "color": "red", "shape": "SQUARE", "flag": 1, "index": 15, "quantity": 79.2778, "rate": 0.0130, "ratio": 6098.292308 }
    { "color": "red", "shape": "CIRCLE", "flag": 1, "index": 16, "quantity": 13.8103, "rate": 2.9010, "ratio": 4.760531 }
    { "color": "red", "shape": "SQUARE", "flag": 0, "index": 48, "quantity": 77.5542, "rate": 7.4670, "ratio": 10.386260 }
    { "color": "purple", "shape": "TRIANGLE", "flag": 0, "index": 51, "quantity": 81.2290, "rate": 8.5910, "ratio": 9.455127 }
    { "color": "red", "shape": "SQUARE", "flag": 0, "index": 64, "quantity": 77.1991, "rate": 9.5310, "ratio": 8.099790 }
    { "color": "purple", "shape": "TRIANGLE", "flag": 0, "index": 65, "quantity": 80.1405, "rate": 5.8240, "ratio": 13.760388 }
    { "color": "yellow", "shape": "CIRCLE", "flag": 1, "index": 73, "quantity": 63.9785, "rate": 4.2370, "ratio": 15.099953 }
    { "color": "yellow", "shape": "CIRCLE", "flag": 1, "index": 87, "quantity": 63.5058, "rate": 8.3350, "ratio": 7.619172 }
    { "color": "purple", "shape": "SQUARE", "flag": 0, "index": 91, "quantity": 72.3735, "rate": 8.2430, "ratio": 8.779995 }

Or, JSON output with vertical-formatting flags::

    $ mlr --icsv --ojsonx tail -n 2 example.csv
    {
      "color": "yellow",
      "shape": "circle",
      "flag": 1,
      "index": 87,
      "quantity": 63.5058,
      "rate": 8.3350
    }
    {
      "color": "purple",
      "shape": "square",
      "flag": 0,
      "index": 91,
      "quantity": 72.3735,
      "rate": 8.2430
    }

Sorts and stats
^^^^^^^^^^^^^^^

Now suppose you want to sort the data on a given column, *and then* take the top few in that ordering. You can use Miller's ``then`` feature to pipe commands together.  

Here are the records with the top three ``index`` values::

    $ mlr --icsv --opprint sort -f shape -nr index then head -n 3 example.csv
    color  shape  flag index quantity rate
    yellow circle 1    87    63.5058  8.3350
    yellow circle 1    73    63.9785  4.2370
    red    circle 1    16    13.8103  2.9010

Lots of Miller commands take a ``-g`` option for group-by: here, ``head -n 1 -g shape`` outputs the first record for each distinct value of the ``shape`` field. This means we're finding the record with highest ``index`` field for each distinct ``shape`` field::

    $ mlr --icsv --opprint sort -f shape -nr index then head -n 1 -g shape example.csv
    color  shape    flag index quantity rate
    yellow circle   1    87    63.5058  8.3350
    purple square   0    91    72.3735  8.2430
    purple triangle 0    65    80.1405  5.8240

Statistics can be computed with or without group-by field(s)::

    $ mlr --icsv --opprint --from example.csv stats1 -a count,min,mean,max -f quantity -g shape
    shape    quantity_count quantity_min quantity_mean quantity_max
    triangle 3              43.649800    68.339767     81.229000
    square   4              72.373500    76.601150     79.277800
    circle   3              13.810300    47.098200     63.978500
    $ mlr --icsv --opprint --from example.csv stats1 -a count,min,mean,max -f quantity -g shape,color
    shape    color  quantity_count quantity_min quantity_mean quantity_max
    triangle yellow 1              43.649800    43.649800     43.649800
    square   red    3              77.199100    78.010367     79.277800
    circle   red    1              13.810300    13.810300     13.810300
    triangle purple 2              80.140500    80.684750     81.229000
    circle   yellow 2              63.505800    63.742150     63.978500
    square   purple 1              72.373500    72.373500     72.373500

If your output has a lot of columns, you can use XTAB format to line things up vertically for you instead::

    $ mlr --icsv --oxtab --from example.csv stats1 -a p0,p10,p25,p50,p75,p90,p99,p100 -f rate
    rate_p0   0.013000
    rate_p10  2.901000
    rate_p25  4.237000
    rate_p50  8.243000
    rate_p75  8.591000
    rate_p90  9.887000
    rate_p99  9.887000
    rate_p100 9.887000

Choices for printing to files
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Often we want to print output <span class="boldmaroon">to the screen</span>. Miller does this by default, as we've seen in the previous examples.

Sometimes we want to print output to another file: <span class="boldmaroon">just use '> outputfilenamegoeshere'</span> at the end of your command::

    % mlr --icsv --opprint cat example.csv > newfile.csv
    # Output goes to the new file;
    # nothing is printed to the screen.
    </pre> </div>
    </td><td>
    <div class="pokipanel"> <pre>
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

Other times we just want our files to be changed in-place: <span class="boldmaroon">just use 'mlr -I'</span>.::

    % cp example.csv newfile.txt

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

    % mlr -I --icsv --opprint cat newfile.txt

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

Also using ``mlr -I`` you can bulk-operate on lots of files: e.g.::

    mlr -I --csv cut -x -f unwanted_column_name *.csv

If you like, you can first copy off your original data somewhere else, before doing in-place operations.

Lastly, using ``tee`` within ``put``, you can split your input data into separate files per one or more field names::

    $ mlr --csv --from example.csv put -q 'tee > $shape.".csv", $*'

    $ cat circle.csv
    color,shape,flag,index,quantity,rate
    red,circle,1,16,13.8103,2.9010
    yellow,circle,1,73,63.9785,4.2370
    yellow,circle,1,87,63.5058,8.3350

    $ cat square.csv
    color,shape,flag,index,quantity,rate
    red,square,1,15,79.2778,0.0130
    red,square,0,48,77.5542,7.4670
    red,square,0,64,77.1991,9.5310
    purple,square,0,91,72.3735,8.2430

    $ cat triangle.csv
    color,shape,flag,index,quantity,rate
    yellow,triangle,1,11,43.6498,9.8870
    purple,triangle,0,51,81.2290,8.5910
    purple,triangle,0,65,80.1405,5.8240


Other-format examples
^^^^^^^^^^^^^^^^^^^^^

What's a CSV file, really? It's an array of rows, or *records*, each being a list of key-value pairs, or *fields*: for CSV it so happens that all the keys are shared in the header line and the values vary data line by data line.

For example, if you have::

    shape,flag,index
    circle,1,24
    square,0,36
    </pre>

then that's a way of saying::

    shape=circle,flag=1,index=24
    shape=square,flag=0,index=36

Data written this way are called <span class="boldmaroon">DKVP</span>, for *delimited key-value pairs*.

We've also already seen other ways to write the same data::

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

SQL-output examples
^^^^^^^^^^^^^^^^^^^

I like to produce SQL-query output with header-column and tab delimiter: this is CSV but with a tab instead of a comma, also known as TSV. Then I post-process with ``mlr --tsv`` or ``mlr --tsvlite``.  This means I can do some (or all, or none) of my data processing within SQL queries, and some (or none, or all) of my data processing using Miller -- whichever is most convenient for my needs at the moment.

For example, using default output formatting in ``mysql`` we get formatting like Miller's ``--opprint --barred``::

    $ mysql --database=mydb -e 'show columns in mytable'
    +------------------+--------------+------+-----+---------+-------+
    | Field            | Type         | Null | Key | Default | Extra |
    +------------------+--------------+------+-----+---------+-------+
    | id               | bigint(20)   | NO   | MUL | NULL    |       |
    | category         | varchar(256) | NO   |     | NULL    |       |
    | is_permanent     | tinyint(1)   | NO   |     | NULL    |       |
    | assigned_to      | bigint(20)   | YES  |     | NULL    |       |
    | last_update_time | int(11)      | YES  |     | NULL    |       |
    +------------------+--------------+------+-----+---------+-------+

Using ``mysql``'s ``-B`` we get TSV output::

    $ mysql --database=mydb -B -e 'show columns in mytable' | mlr --itsvlite --opprint cat
    Field            Type         Null Key Default Extra
    id               bigint(20)   NO  MUL NULL -
    category         varchar(256) NO  -   NULL -
    is_permanent     tinyint(1)   NO  -   NULL -
    assigned_to      bigint(20)   YES -   NULL -
    last_update_time int(11)      YES -   NULL -

Since Miller handles TSV output, we can do as much or as little processing as we want in the SQL query, then send the rest on to Miller. This includes outputting as JSON, doing further selects/joins in Miller, doing stats, etc.  etc.::

    $ mysql --database=mydb -B -e 'show columns in mytable' | mlr --itsvlite --ojson --jlistwrap --jvstack cat 
    [
      {
        "Field": "id",
        "Type": "bigint(20)",
        "Null": "NO",
        "Key": "MUL",
        "Default": "NULL",
        "Extra": ""
      },
      {
        "Field": "category",
        "Type": "varchar(256)",
        "Null": "NO",
        "Key": "",
        "Default": "NULL",
        "Extra": ""
      },
      {
        "Field": "is_permanent",
        "Type": "tinyint(1)",
        "Null": "NO",
        "Key": "",
        "Default": "NULL",
        "Extra": ""
      },
      {
        "Field": "assigned_to",
        "Type": "bigint(20)",
        "Null": "YES",
        "Key": "",
        "Default": "NULL",
        "Extra": ""
      },
      {
        "Field": "last_update_time",
        "Type": "int(11)",
        "Null": "YES",
        "Key": "",
        "Default": "NULL",
        "Extra": ""
      }
    ]

    $ mysql --database=mydb -B -e 'select * from mytable' > query.tsv

    $ mlr --from query.tsv --t2p stats1 -a count -f id -g category,assigned_to
    category assigned_to id_count
    special  10000978    207
    special  10003924    385
    special  10009872    168
    standard 10000978    524
    standard 10003924    392
    standard 10009872    108
    ...

Again, all the examples in the CSV section apply here -- just change the input-format flags.

SQL-input examples
^^^^^^^^^^^^^^^^^^

One use of NIDX (value-only, no keys) format is for loading up SQL tables.

Create and load SQL table::

    mysql> CREATE TABLE abixy(
      a VARCHAR(32),
      b VARCHAR(32),
      i BIGINT(10),
      x DOUBLE,
      y DOUBLE
    );
    Query OK, 0 rows affected (0.01 sec)

    bash$ mlr --onidx --fs comma cat data/medium > medium.nidx

    mysql> LOAD DATA LOCAL INFILE 'medium.nidx' REPLACE INTO TABLE abixy FIELDS TERMINATED BY ',' ;
    Query OK, 10000 rows affected (0.07 sec)
    Records: 10000  Deleted: 0  Skipped: 0  Warnings: 0

    mysql> SELECT COUNT(*) AS count FROM abixy;
    +-------+
    | count |
    +-------+
    | 10000 |
    +-------+
    1 row in set (0.00 sec)

    mysql> SELECT * FROM abixy LIMIT 10;
    +------+------+------+---------------------+---------------------+
    | a    | b    | i    | x                   | y                   |
    +------+------+------+---------------------+---------------------+
    | pan  | pan  |    1 |  0.3467901443380824 |  0.7268028627434533 |
    | eks  | pan  |    2 |  0.7586799647899636 |  0.5221511083334797 |
    | wye  | wye  |    3 | 0.20460330576630303 | 0.33831852551664776 |
    | eks  | wye  |    4 | 0.38139939387114097 | 0.13418874328430463 |
    | wye  | pan  |    5 |  0.5732889198020006 |  0.8636244699032729 |
    | zee  | pan  |    6 |  0.5271261600918548 | 0.49322128674835697 |
    | eks  | zee  |    7 |  0.6117840605678454 |  0.1878849191181694 |
    | zee  | wye  |    8 |  0.5985540091064224 |   0.976181385699006 |
    | hat  | wye  |    9 | 0.03144187646093577 |  0.7495507603507059 |
    | pan  | wye  |   10 |  0.5026260055412137 |  0.9526183602969864 |
    +------+------+------+---------------------+---------------------+

Aggregate counts within SQL::

    mysql> SELECT a, b, COUNT(*) AS count FROM abixy GROUP BY a, b ORDER BY COUNT DESC;
    +------+------+-------+
    | a    | b    | count |
    +------+------+-------+
    | zee  | wye  |   455 |
    | pan  | eks  |   429 |
    | pan  | pan  |   427 |
    | wye  | hat  |   426 |
    | hat  | wye  |   423 |
    | pan  | hat  |   417 |
    | eks  | hat  |   417 |
    | pan  | zee  |   413 |
    | eks  | eks  |   413 |
    | zee  | hat  |   409 |
    | eks  | wye  |   407 |
    | zee  | zee  |   403 |
    | pan  | wye  |   395 |
    | wye  | pan  |   392 |
    | zee  | eks  |   391 |
    | zee  | pan  |   389 |
    | hat  | eks  |   389 |
    | wye  | eks  |   386 |
    | wye  | zee  |   385 |
    | hat  | zee  |   385 |
    | hat  | hat  |   381 |
    | wye  | wye  |   377 |
    | eks  | pan  |   371 |
    | hat  | pan  |   363 |
    | eks  | zee  |   357 |
    +------+------+-------+
    25 rows in set (0.01 sec)

Aggregate counts within Miller::

    $ mlr --opprint uniq -c -g a,b then sort -nr count data/medium
    a   b   count
    zee wye 455
    pan eks 429
    pan pan 427
    wye hat 426
    hat wye 423
    pan hat 417
    eks hat 417
    eks eks 413
    pan zee 413
    zee hat 409
    eks wye 407
    zee zee 403
    pan wye 395
    hat pan 363
    eks zee 357

Pipe SQL output to aggregate counts within Miller::

    $ mysql -D miller -B -e 'select * from abixy' | mlr --itsv --opprint uniq -c -g a,b then sort -nr count
    a   b   count
    zee wye 455
    pan eks 429
    pan pan 427
    wye hat 426
    hat wye 423
    pan hat 417
    eks hat 417
    eks eks 413
    pan zee 413
    zee hat 409
    eks wye 407
    zee zee 403
    pan wye 395
    wye pan 392
    zee eks 391
    zee pan 389
    hat eks 389
    wye eks 386
    hat zee 385
    wye zee 385
    hat hat 381
    wye wye 377
    eks pan 371
    hat pan 363
    eks zee 357

Log-processing examples
^^^^^^^^^^^^^^^^^^^^^^^

Another of my favorite use-cases for Miller is doing ad-hoc processing of log-file data.  Here's where DKVP format really shines: one, since the field names and field values are present on every line, every line stands on its own. That means you can ``grep`` or what have you. Also it means not every line needs to have the same list of field names ("schema").

Again, all the examples in the CSV section apply here -- just change the input-format flags. But there's more you can do when not all the records have the same shape.

Writing a program -- in any language whatsoever -- you can have it print out log lines as it goes along, with items for various events jumbled together. After the program has finished running you can sort it all out, filter it, analyze it, and learn from it.

Suppose your program has printed something like this::

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

Each print statement simply contains local information: the current timestamp, whether a particular cache was hit or not, etc. Then using either the system ``grep`` command, or Miller's ``having-fields``, or ``is_present``, we can pick out the parts we want and analyze them::

    $ grep op=cache log.txt \
  | mlr --idkvp --opprint stats1 -a mean -f hit -g type then sort -f type
    type hit_mean
    A1   0.857143
    A4   0.714286
    A9   0.090909
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

Alternatively, we can simply group the similar data for a better look::

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

More
^^^^

Please see the <a href="reference.html">reference</a> for complete information, as well as the <a href="faq.html">FAQ</a> and the <a href="cookbook.html">cookbook</a> for more tips.
