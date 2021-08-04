<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# Record-heterogeneity

We think of CSV tables as rectangular: if there are 17 columns in the header then there are 17 columns for every row, else the data have a formatting error.

But heterogeneous data abound (today's no-SQL databases for example). Miller handles this.

## For I/O

### CSV and pretty-print

Miller simply prints a newline and a new header when there is a schema change. When there is no schema change, you get CSV per se as a special case. Likewise, Miller reads heterogeneous CSV or pretty-print input the same way. The difference between CSV and CSV-lite is that the former is RFC4180-compliant, while the latter readily handles heterogeneous data (which is non-compliant). For example:

<pre>
<b>cat data/het.dkvp</b>
resource=/path/to/file,loadsec=0.45,ok=true
record_count=100,resource=/path/to/file
resource=/path/to/second/file,loadsec=0.32,ok=true
record_count=150,resource=/path/to/second/file
resource=/some/other/path,loadsec=0.97,ok=false
</pre>

<pre>
<b>mlr --ocsvlite cat data/het.dkvp</b>
resource,loadsec,ok
/path/to/file,0.45,true

record_count,resource
100,/path/to/file

resource,loadsec,ok
/path/to/second/file,0.32,true

record_count,resource
150,/path/to/second/file

resource,loadsec,ok
/some/other/path,0.97,false
</pre>

<pre>
<b>mlr --opprint cat data/het.dkvp</b>
resource      loadsec ok
/path/to/file 0.45    true

record_count resource
100          /path/to/file

resource             loadsec ok
/path/to/second/file 0.32    true

record_count resource
150          /path/to/second/file

resource         loadsec ok
/some/other/path 0.97    false
</pre>

Miller handles explicit header changes as just shown. If your CSV input contains ragged data -- if there are implicit header changes -- you can use ``--allow-ragged-csv-input`` (or keystroke-saver ``--ragged``). For too-short data lines, values are filled with empty string; for too-long data lines, missing field names are replaced with positional indices (counting up from 1, not 0), as follows:

<pre>
<b>cat data/ragged.csv</b>
a,b,c
1,2,3
4,5
6,7,8,9
</pre>

<pre>
<b>mlr --icsv --oxtab --allow-ragged-csv-input cat data/ragged.csv</b>
a 1
b 2
c 3

a 4
b 5
c 

a 6
b 7
c 8
4 9
</pre>

You may also find Miller's ``group-like`` feature handy (see also :doc:`reference-verbs`):

<pre>
<b>mlr --ocsvlite group-like data/het.dkvp</b>
resource,loadsec,ok
/path/to/file,0.45,true
/path/to/second/file,0.32,true
/some/other/path,0.97,false

record_count,resource
100,/path/to/file
150,/path/to/second/file
</pre>

<pre>
<b>mlr --opprint group-like data/het.dkvp</b>
resource             loadsec ok
/path/to/file        0.45    true
/path/to/second/file 0.32    true
/some/other/path     0.97    false

record_count resource
100          /path/to/file
150          /path/to/second/file
</pre>

### Key-value-pair, vertical-tabular, and index-numbered formats

For these formats, record-heterogeneity comes naturally:

<pre>
<b>cat data/het.dkvp</b>
resource=/path/to/file,loadsec=0.45,ok=true
record_count=100,resource=/path/to/file
resource=/path/to/second/file,loadsec=0.32,ok=true
record_count=150,resource=/path/to/second/file
resource=/some/other/path,loadsec=0.97,ok=false
</pre>

<pre>
<b>mlr --onidx --ofs ' ' cat data/het.dkvp</b>
/path/to/file 0.45 true
100 /path/to/file
/path/to/second/file 0.32 true
150 /path/to/second/file
/some/other/path 0.97 false
</pre>

<pre>
<b>mlr --oxtab cat data/het.dkvp</b>
resource /path/to/file
loadsec  0.45
ok       true

record_count 100
resource     /path/to/file

resource /path/to/second/file
loadsec  0.32
ok       true

record_count 150
resource     /path/to/second/file

resource /some/other/path
loadsec  0.97
ok       false
</pre>

<pre>
<b>mlr --oxtab group-like data/het.dkvp</b>
resource /path/to/file
loadsec  0.45
ok       true

resource /path/to/second/file
loadsec  0.32
ok       true

resource /some/other/path
loadsec  0.97
ok       false

record_count 100
resource     /path/to/file

record_count 150
resource     /path/to/second/file
</pre>

## For processing

Miller operates on specified fields and takes the rest along: for example, if you are sorting on the ``count`` field then all records in the input stream must have a ``count`` field but the other fields can vary, and moreover the sorted-on field name(s) don't need to be in the same position on each line:

<pre>
<b>cat data/sort-het.dkvp</b>
count=500,color=green
count=600
status=ok,count=250,hours=0.22
status=ok,count=200,hours=3.4
count=300,color=blue
count=100,color=green
count=450
</pre>

<pre>
<b>mlr sort -n count data/sort-het.dkvp</b>
count=100,color=green
status=ok,count=200,hours=3.4
status=ok,count=250,hours=0.22
count=300,color=blue
count=450
count=500,color=green
count=600
</pre>
