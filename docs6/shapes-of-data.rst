..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Shapes of data
==============

No output at all
----------------------------------------------------------------

Try ``od -xcv`` and/or ``cat -e`` on your file to check for non-printable characters.

If you're using Miller version less than 5.0.0 (try ``mlr --version`` on your system to find out), when the line-ending-autodetect feature was introduced, please see http://johnkerl.org/miller-releases/miller-4.5.0/doc/index.html.

Fields not selected
----------------------------------------------------------------

Check the field-separators of the data, e.g. with the command-line ``head`` program. Example: for CSV, Miller's default record separator is comma; if your data is tab-delimited, e.g. ``aTABbTABc``, then Miller won't find three fields named ``a``, ``b``, and ``c`` but rather just one named ``aTABbTABc``.  Solution in this case: ``mlr --fs tab {remaining arguments ...}``.

Also try ``od -xcv`` and/or ``cat -e`` on your file to check for non-printable characters.

Diagnosing delimiter specifications
----------------------------------------------------------------

.. code-block:: none

    # Use the `file` command to see if there are CR/LF terminators (in this case,
    # there are not):
    $ file data/colours.csv 
    data/colours.csv: UTF-8 Unicode text
    
    # Look at the file to find names of fields
    $ cat data/colours.csv 
    KEY;DE;EN;ES;FI;FR;IT;NL;PL;RO;TR
    masterdata_colourcode_1;Weiß;White;Blanco;Valkoinen;Blanc;Bianco;Wit;Biały;Alb;Beyaz
    masterdata_colourcode_2;Schwarz;Black;Negro;Musta;Noir;Nero;Zwart;Czarny;Negru;Siyah
    
    # Extract a few fields:
    $ mlr --csv cut -f KEY,PL,RO data/colours.csv 
    (only blank lines appear)
    
    # Use XTAB output format to get a sharper picture of where records/fields
    # are being split:
    $ mlr --icsv --oxtab cat data/colours.csv 
    KEY;DE;EN;ES;FI;FR;IT;NL;PL;RO;TR masterdata_colourcode_1;Weiß;White;Blanco;Valkoinen;Blanc;Bianco;Wit;Biały;Alb;Beyaz
    
    KEY;DE;EN;ES;FI;FR;IT;NL;PL;RO;TR masterdata_colourcode_2;Schwarz;Black;Negro;Musta;Noir;Nero;Zwart;Czarny;Negru;Siyah
    
    # Using XTAB output format makes it clearer that KEY;DE;...;RO;TR is being
    # treated as a single field name in the CSV header, and likewise each
    # subsequent line is being treated as a single field value. This is because
    # the default field separator is a comma but we have semicolons here.
    # Use XTAB again with different field separator (--fs semicolon):
     mlr --icsv --ifs semicolon --oxtab cat data/colours.csv 
    KEY masterdata_colourcode_1
    DE  Weiß
    EN  White
    ES  Blanco
    FI  Valkoinen
    FR  Blanc
    IT  Bianco
    NL  Wit
    PL  Biały
    RO  Alb
    TR  Beyaz
    
    KEY masterdata_colourcode_2
    DE  Schwarz
    EN  Black
    ES  Negro
    FI  Musta
    FR  Noir
    IT  Nero
    NL  Zwart
    PL  Czarny
    RO  Negru
    TR  Siyah
    
    # Using the new field-separator, retry the cut:
     mlr --csv --fs semicolon cut -f KEY,PL,RO data/colours.csv 
    KEY;PL;RO
    masterdata_colourcode_1;Biały;Alb
    masterdata_colourcode_2;Czarny;Negru

I assigned $9 and it's not 9th
----------------------------------------------------------------

Miller records are ordered lists of key-value pairs. For NIDX format, DKVP format when keys are missing, or CSV/CSV-lite format with ``--implicit-csv-header``, Miller will sequentially assign keys of the form ``1``, ``2``, etc. But these are not integer array indices: they're just field names taken from the initial field ordering in the input data.

.. code-block:: none
   :emphasize-lines: 1-1

    $ echo x,y,z | mlr --dkvp cat
    1=x,2=y,3=z

.. code-block:: none
   :emphasize-lines: 1-1

    $ echo x,y,z | mlr --dkvp put '$6="a";$4="b";$55="cde"'
    1=x,2=y,3=z,6=a,4=b,55=cde

.. code-block:: none
   :emphasize-lines: 1-1

    $ echo x,y,z | mlr --nidx cat
    x,y,z

.. code-block:: none
   :emphasize-lines: 1-1

    $ echo x,y,z | mlr --csv --implicit-csv-header cat
    1,2,3
    x,y,z

.. code-block:: none
   :emphasize-lines: 1-1

    $ echo x,y,z | mlr --dkvp rename 2,999
    1=x,999=y,3=z

.. code-block:: none
   :emphasize-lines: 1-1

    $ echo x,y,z | mlr --dkvp rename 2,newname
    1=x,newname=y,3=z

.. code-block:: none
   :emphasize-lines: 1-1

    $ echo x,y,z | mlr --csv --implicit-csv-header reorder -f 3,1,2
    3,1,2
    z,x,y

Why doesn't mlr cut put fields in the order I want?
----------------------------------------------------------------

Example: columns ``x,i,a`` were requested but they appear here in the order ``a,i,x``:

.. code-block:: none
   :emphasize-lines: 1-1

    $ cat data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr cut -f x,i,a data/small
    a=pan,i=1,x=0.3467901443380824
    a=eks,i=2,x=0.7586799647899636
    a=wye,i=3,x=0.20460330576630303
    a=eks,i=4,x=0.38139939387114097
    a=wye,i=5,x=0.5732889198020006

The issue is that Miller's ``cut``, by default, outputs cut fields in the order they appear in the input data. This design decision was made intentionally to parallel the Unix/Linux system ``cut`` command, which has the same semantics.

The solution is to use the ``-o`` option:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr cut -o -f x,i,a data/small
    x=0.3467901443380824,i=1,a=pan
    x=0.7586799647899636,i=2,a=eks
    x=0.20460330576630303,i=3,a=wye
    x=0.38139939387114097,i=4,a=eks
    x=0.5732889198020006,i=5,a=wye

Numbering and renumbering records
----------------------------------------------------------------

The ``awk``-like built-in variable ``NR`` is incremented for each input record:

.. code-block:: none
   :emphasize-lines: 1-1

    $ cat data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr put '$nr = NR' data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533,nr=1
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797,nr=2
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776,nr=3
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463,nr=4
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,nr=5

However, this is the record number within the original input stream -- not after any filtering you may have done:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr filter '$a == "wye"' then put '$nr = NR' data/small
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776,nr=3
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,nr=5

There are two good options here. One is to use the ``cat`` verb with ``-n``:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr filter '$a == "wye"' then cat -n data/small
    n=1,a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    n=2,a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729

The other is to keep your own counter within the ``put`` DSL:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr filter '$a == "wye"' then put 'begin {@n = 1} $n = @n; @n += 1' data/small
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776,n=1
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,n=2

The difference is a matter of taste (although ``mlr cat -n`` puts the counter first).

Splitting nested fields
----------------------------------------------------------------

Suppose you have a TSV file like this:

.. code-block:: none

    a	b
    x	z
    s	u:v:w

The simplest option is to use :ref:`mlr nest <reference-verbs-nest>`:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --tsv nest --explode --values --across-records -f b --nested-fs : data/nested.tsv
    a	b
    x	z
    s	u
    s	v
    s	w

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --tsv nest --explode --values --across-fields  -f b --nested-fs : data/nested.tsv
    a	b_1
    x	z
    
    a	b_1	b_2	b_3
    s	u	v	w

While ``mlr nest`` is simplest, let's also take a look at a few ways to do this using the ``put`` DSL.

One option to split out the colon-delimited values in the ``b`` column is to use ``splitnv`` to create an integer-indexed map and loop over it, adding new fields to the current record:

.. code-block:: none
   :emphasize-lines: 1-6

    $ mlr --from data/nested.tsv --itsv --oxtab put '
      o = splitnv($b, ":");
      for (k,v in o) {
        $["p".k]=v
      }
    '
    a  x
    b  z
    p1 z
    
    a  s
    b  u:v:w
    p1 u
    p2 v
    p3 w

while another is to loop over the same map from ``splitnv`` and use it (with ``put -q`` to suppress printing the original record) to produce multiple records:

.. code-block:: none
   :emphasize-lines: 1-7

    $ mlr --from data/nested.tsv --itsv --oxtab put -q '
      o = splitnv($b, ":");
      for (k,v in o) {
        x = mapsum($*, {"b":v});
        emit x
      }
    '
    a x
    b z
    
    a s
    b u
    
    a s
    b v
    
    a s
    b w

.. code-block:: none
   :emphasize-lines: 1-6

    $ mlr --from data/nested.tsv --tsv put -q '
      o = splitnv($b, ":");
      for (k,v in o) {
        x = mapsum($*, {"b":v}); emit x
      }
    '
    a	b
    x	z
    s	u
    s	v
    s	w

Options for dealing with duplicate rows
----------------------------------------------------------------

If your data has records appearing multiple times, you can use :ref:`mlr uniq <reference-verbs-uniq>` to show and/or count the unique records.

If you want to look at partial uniqueness -- for example, show only the first record for each unique combination of the ``account_id`` and ``account_status`` fields -- you might use ``mlr head -n 1 -g account_id,account_status``. Please also see :ref:`mlr head <reference-verbs-head>`.

Rectangularizing data
----------------------------------------------------------------

Suppose you have a method (in whatever language) which is printing things of the form

.. code-block:: none

    outer=1
    outer=2
    outer=3

and then calls another method which prints things of the form

.. code-block:: none

    middle=10
    middle=11
    middle=12
    middle=20
    middle=21
    middle=30
    middle=31

and then, perhaps, that second method calls a third method which prints things of the form

.. code-block:: none

    inner1=100,inner2=101
    inner1=120,inner2=121
    inner1=200,inner2=201
    inner1=210,inner2=211
    inner1=300,inner2=301
    inner1=312
    inner1=313,inner2=314

with the result that your program's output is

.. code-block:: none

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

.. code-block:: none
   :emphasize-lines: 1-10

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
