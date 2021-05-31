..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

SQL examples
====================

.. _sql-output-examples:

SQL-output examples
^^^^^^^^^^^^^^^^^^^

I like to produce SQL-query output with header-column and tab delimiter: this is CSV but with a tab instead of a comma, also known as TSV. Then I post-process with ``mlr --tsv`` or ``mlr --tsvlite``.  This means I can do some (or all, or none) of my data processing within SQL queries, and some (or none, or all) of my data processing using Miller -- whichever is most convenient for my needs at the moment.

For example, using default output formatting in ``mysql`` we get formatting like Miller's ``--opprint --barred``:

.. code-block:: none
   :emphasize-lines: 1,1

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

Using ``mysql``'s ``-B`` we get TSV output:

.. code-block:: none
   :emphasize-lines: 1,1

    $ mysql --database=mydb -B -e 'show columns in mytable' | mlr --itsvlite --opprint cat
    Field            Type         Null Key Default Extra
    id               bigint(20)   NO  MUL NULL -
    category         varchar(256) NO  -   NULL -
    is_permanent     tinyint(1)   NO  -   NULL -
    assigned_to      bigint(20)   YES -   NULL -
    last_update_time int(11)      YES -   NULL -

Since Miller handles TSV output, we can do as much or as little processing as we want in the SQL query, then send the rest on to Miller. This includes outputting as JSON, doing further selects/joins in Miller, doing stats, etc.  etc.:

.. code-block:: none
   :emphasize-lines: 1,1

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

.. code-block:: none
   :emphasize-lines: 1,1

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

.. _sql-input-examples:

SQL-input examples
^^^^^^^^^^^^^^^^^^

One use of NIDX (value-only, no keys) format is for loading up SQL tables.

Create and load SQL table:

.. code-block:: none

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

Aggregate counts within SQL:

.. code-block:: none
   :emphasize-lines: 1,1

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

Aggregate counts within Miller:

.. code-block:: none
   :emphasize-lines: 1,1

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

Pipe SQL output to aggregate counts within Miller:

.. code-block:: none
   :emphasize-lines: 1,1

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
