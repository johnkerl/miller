<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-verbs/index.html">Verb list</a>
&nbsp;
<a class="quicklink" href="../reference-dsl-builtin-functions/index.html">Function list</a>
&nbsp;
<a class="quicklink" href="../glossary/index.html">Glossary</a>
&nbsp;
<a class="quicklink" href="https://github.com/johnkerl/miller" target="_blank">Repository ↗</a>
</span>
</div>
# Record-heterogeneity

We think of CSV tables as rectangular: if there are 17 columns in the header
then there are 17 columns for every row, else the data have a formatting error.

But heterogeneous data abound -- log-file entries, JSON documents, today's
no-SQL databases, etc. Miller handles heterogeneity in multiple ways.

## Examples and terminology

### Homogeneous/rectangular

A **homogeneoust** list of records is one in which all records have the _same keys, in the same order_.
For example, if a well-formed [CSV file](file-formats.md#csvtsvasvusvetc) looks like

<pre class="pre-non-highlight-non-pair">
a,b,c
1,2,3
4,5,6
7,8,9
</pre>

then there are three records (written here using JSON formatting):

<pre class="pre-non-highlight-non-pair">
{"a": 1, "b": 2, "c": 3}
{"a": 4, "b": 5, "c": 6}
{"a": 7, "b": 8, "c": 9}
</pre>

And every row has the same keys, in the same order: `a,b,c`.

These are also sometimes called **rectangular** since if we pretty-print them we get a nice rectangle:

<pre class="pre-non-highlight-non-pair">
a b c
1 2 3
4 5 6
7 8 9
</pre>

### Fillable

A second example has some empty cells which could be **filled**:

<pre class="pre-non-highlight-non-pair">
a,b,c
1,2,3
4,,6
,8,9
</pre>

<pre class="pre-non-highlight-non-pair">
{"a": 1, "b": 2, "c": 3}
{"a": 4, "b": "", "c": 6}
{"a": "", "b": 8, "c": 9}
</pre>

This example is still homogeneous, though: every row has the same keys, in the same order: `a,b,c`.

TODO: link to fill-down etc (here or below)

### Ragged

Next let's look at non-well-formed CSV files. For a third example:

<pre class="pre-non-highlight-non-pair">
a,b,c
1,2,3
4,5
7,8,9,10
</pre>

TODO: named data/x for all these.

If you `mlr csv cat` this, you'll get an error message like `CSV header/data
length mismatch 3 != 2 at filename (stdin) row 3.` This kind of data is referred to as **ragged**.

### Regular/irregular

Next let's look at some JSON data:

<pre class="pre-non-highlight-non-pair">
{"a": 1, "b": 2, "c": 3}
{"c": 6, "a": 4, "b": 5}
{"b": 8, "c": 9, "a": 7}
</pre>

xxx same. xxx regularize.

xxx regularize then unsparsify.

xxx refer to data-cleaning prominently, maybe in the page title
xxx link data-cleaning examples <--> here, both ways

xxx on how sparse originates -- log files / etc from software which removes (or
never populates) keys with empty values.

## TBF

```
  heterogeneity ragged/rectangular sparse/unsparse regularize
  - heterogeneity 'all the same' = same keys in same order
  - rectangular = homoegeneous
  - ragged = varying number of keys; maybe 'sparse'
  - sparse / unsparsify = not rect due to either empty values or absent keys
  - irregular / regularize = sort keys -- a,b,c / c,a,b problem

  record-heterogeneity.md.in

  shapes-of-data.md: instead, use the following
  csv-with-and-without-headers.md.in:## Regularizing ragged CSV
  misc-examples.md.in:Then, join on the key field(s), and use unsparsify to zero-fill counters
  questions-about-joins.md.in:To fix this, we can use **unsparsify**:
  two-pass-algorithms.md.in:There is a keystroke-saving verb for this: [unsparsify](reference-verbs.md#unsparsify)
  programming-language.md.in:easy for you to handle non-heterogeneous data
  two-pass-algorithms.md.in:Suppose you have some heterogeneous data like this:
```

## Reading and writing heterogeneous data

### Rectangular file formats: CSV and pretty-print

Miller simply prints a newline and a new header when there is a schema change. When there is no schema change, you get CSV per se as a special case. Likewise, Miller reads heterogeneous CSV or pretty-print input the same way. The difference between CSV and CSV-lite is that the former is RFC4180-compliant, while the latter readily handles heterogeneous data (which is non-compliant). For example:

<pre class="pre-highlight-in-pair">
<b>cat data/het.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "resource": "/path/to/file",
  "loadsec": 0.45,
  "ok": true
}
{
  "record_count": 100,
  "resource": "/path/to/file"
}
{
  "resource": "/path/to/second/file",
  "loadsec": 0.32,
  "ok": true
}
{
  "record_count": 150,
  "resource": "/path/to/second/file"
}
{
  "resource": "/some/other/path",
  "loadsec": 0.97,
  "ok": false
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --opprint cat data/het.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
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

Miller handles explicit header changes as just shown. If your CSV input contains ragged data -- if there are implicit header changes (no intervening blank line and new header line) -- you can use `--allow-ragged-csv-input` (or keystroke-saver `--ragged`). For too-short data lines, values are filled with empty string; for too-long data lines, missing field names are replaced with positional indices (counting up from 1, not 0), as follows:

<pre class="pre-highlight-in-pair">
<b>cat data/ragged.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c
1,2,3
4,5
6,7,8,9
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --oxtab --allow-ragged-csv-input cat data/ragged.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
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

You may also find Miller's `group-like` feature handy (see also [Verbs Reference](reference-verbs.md)):

<pre class="pre-highlight-in-pair">
<b>mlr --j2p cat data/het.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
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

<pre class="pre-highlight-in-pair">
<b>mlr --j2p group-like data/het.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
resource             loadsec ok
/path/to/file        0.45    true
/path/to/second/file 0.32    true
/some/other/path     0.97    false

record_count resource
100          /path/to/file
150          /path/to/second/file
</pre>

### Non-rectangular file formats: JSON, XTAB, NIDX, DKVP

For these formats, record-heterogeneity comes naturally:

<pre class="pre-highlight-in-pair">
<b>cat data/het.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "resource": "/path/to/file",
  "loadsec": 0.45,
  "ok": true
}
{
  "record_count": 100,
  "resource": "/path/to/file"
}
{
  "resource": "/path/to/second/file",
  "loadsec": 0.32,
  "ok": true
}
{
  "record_count": 150,
  "resource": "/path/to/second/file"
}
{
  "resource": "/some/other/path",
  "loadsec": 0.97,
  "ok": false
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --onidx --ofs ' ' cat data/het.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
/path/to/file 0.45 true
100 /path/to/file
/path/to/second/file 0.32 true
150 /path/to/second/file
/some/other/path 0.97 false
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --oxtab cat data/het.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
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

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --oxtab group-like data/het.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
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

## Processing heterogeneous data

Miller operates on specified fields and takes the rest along: for example, if you are sorting on the `count` field then all records in the input stream must have a `count` field but the other fields can vary, and moreover the sorted-on field name(s) don't need to be in the same position on each line:

<pre class="pre-highlight-in-pair">
<b>cat data/sort-het.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
count=500,color=green
count=600
status=ok,count=250,hours=0.22
status=ok,count=200,hours=3.4
count=300,color=blue
count=100,color=green
count=450
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr sort -n count data/sort-het.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
count=100,color=green
status=ok,count=200,hours=3.4
status=ok,count=250,hours=0.22
count=300,color=blue
count=450
count=500,color=green
count=600
</pre>

## Making data more heterogeneous

TODO
