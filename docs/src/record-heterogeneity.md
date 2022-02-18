<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-main-flag-list/index.html">Flags</a>
&nbsp;
<a class="quicklink" href="../reference-verbs/index.html">Verbs</a>
&nbsp;
<a class="quicklink" href="../reference-dsl-builtin-functions/index.html">Functions</a>
&nbsp;
<a class="quicklink" href="../glossary/index.html">Glossary</a>
&nbsp;
<a class="quicklink" href="../release-docs/index.html">Release docs</a>
</span>
</div>
# Record-heterogeneity

We think of CSV tables as rectangular: if there are 17 columns in the header
then there are 17 columns for every row, else the data have a formatting error.

But heterogeneous data abound -- log-file entries, JSON documents, no-SQL
databases such as MongoDB, etc. -- not to mention **data-cleaning
opportunities** we'll look at in this page. Miller offers several ways to
handle data heterogeneity.

## Terminology, examples, and solutions

Different kinds of heterogeneous data include _ragged_, _irregular_, and _sparse_.

### Homogeneous/rectangular data

A **homogeneous** list of records is one in which all records have _the same keys, in the same order_.
For example, here is a well-formed [CSV file](file-formats.md#csvtsvasvusvetc):

<pre class="pre-highlight-in-pair">
<b>mlr --csv cat data/het/hom.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c
1,2,3
4,5,6
7,8,9
</pre>

It has three records (written here using JSON Lines formatting):

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojsonl cat data/het/hom.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{"a": 1, "b": 2, "c": 3}
{"a": 4, "b": 5, "c": 6}
{"a": 7, "b": 8, "c": 9}
</pre>

Here every row has the same keys, in the same order: `a,b,c`.

These are also sometimes called **rectangular** since if we pretty-print them we get a nice rectangle:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint cat data/het/hom.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a b c
1 2 3
4 5 6
7 8 9
</pre>

### Fillable data

A second example has some empty cells which could be **filled**:

<pre class="pre-highlight-in-pair">
<b>mlr --csv cat data/het/fillable.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c
1,2,3
4,,6
,8,9
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojsonl cat data/het/fillable.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{"a": 1, "b": 2, "c": 3}
{"a": 4, "b": "", "c": 6}
{"a": "", "b": 8, "c": 9}
</pre>

This example is still homogeneous, though: every row has the same keys, in the same order: `a,b,c`.
Empty values don't make the data heterogeneous.

Note however that we can use the [`fill-empty`](reference-verbs.md#fill-empty) verb to make these
values non-empty, if we like:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint fill-empty -v filler data/het/fillable.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a      b      c
1      2      3
4      filler 6
filler 8      9
</pre>

### Ragged data

Next let's look at non-well-formed CSV files. For a third example:

<pre class="pre-highlight-in-pair">
<b>cat data/het/ragged.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c
1,2,3
4,5
7,8,9,10
</pre>

If you `mlr --csv cat` this, you'll get an error message:

<pre class="pre-highlight-in-pair">
<b>mlr --csv cat data/het/ragged.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
mlr :  mlr: CSV header/data length mismatch 3 != 2 at filename data/het/ragged.csv row 3.

</pre>

There are two kinds of raggedness here. Since CSVs form records by zipping the
keys from the header line together with the values from each data line, the
second record has a missing value for key `c` (which ought to be fillable),
while the third record has a value `10` with no key for it.

Using the [`--allow-ragged-csv-input` flag](reference-main-flag-list.md#csv-only-flags)
we can fill values in too-short rows, and provide a key (column number starting
with 1) for too-long rows:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --allow-ragged-csv-input cat data/het/ragged.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "a": 1,
  "b": 2,
  "c": 3
},
{
  "a": 4,
  "b": 5,
  "c": ""
},
{
  "a": 7,
  "b": 8,
  "c": 9,
  "4": 10
}
]
</pre>

### Irregular data

Here's another situation -- this file has, in some sense, the "same" data as
our `ragged.csv` example above:

<pre class="pre-highlight-in-pair">
<b>cat data/het/irregular.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{"a": 1, "b": 2, "c": 3}
{"c": 6, "a": 4, "b": 5}
{"b": 8, "c": 9, "a": 7}
</pre>

For example, on the second record, `a` is 4, `b` is 5, `c` is 6. But this data
is heterogeneous because the keys `a,b,c` aren't in the same order in each
record.

This kind of data arises often in practice. One reason is that, while many
programming languages (including the Miller DSL) [preserve insertion
order](reference-main-maps.md#insertion-order-is-preserved) in maps; others do
not. So someone might have written `{"a":4,"b":5,"c":6}` in the source code,
but the data may not have printed that way into a given data file.

We can use the [`regularize`](reference-verbs.md#regularize) or
[`sort-within-records`](reference-verbs.md#sort-within-records) verb to order
the keys:

<pre class="pre-highlight-in-pair">
<b>mlr --jsonl regularize data/het/irregular.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{"a": 1, "b": 2, "c": 3}
{"a": 4, "b": 5, "c": 6}
{"a": 7, "b": 8, "c": 9}
</pre>

The `regularize` verb tries to re-order subsequent rows to look like the first
(whatever order that is); the `sort-within-records` verb simply uses
alphabetical order (which is the same in the above example where the first
record has keys in the order `a,b,c`).

### Sparse data

Here's another frequently occurring situation -- quite often, systems will log
data for items which are present, but won't log data for items which aren't.

<pre class="pre-highlight-in-pair">
<b>mlr --json cat data/het/sparse.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "host": "xy01.east",
  "status": "running",
  "volume": "/dev/sda1"
},
{
  "host": "xy92.west",
  "status": "running"
},
{
  "purpose": "failover",
  "host": "xy55.east",
  "volume": "/dev/sda1",
  "reimaged": true
}
]
</pre>

This data is called **sparse** (from the [data-storage term](https://en.wikipedia.org/wiki/Sparse_matrix)).

We can use the [`unsparsify`](reference-verbs.md#unsparsify) verb to make sure
every record has the same keys:

<pre class="pre-highlight-in-pair">
<b>mlr --json unsparsify data/het/sparse.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "host": "xy01.east",
  "status": "running",
  "volume": "/dev/sda1",
  "purpose": "",
  "reimaged": ""
},
{
  "host": "xy92.west",
  "status": "running",
  "volume": "",
  "purpose": "",
  "reimaged": ""
},
{
  "host": "xy55.east",
  "status": "",
  "volume": "/dev/sda1",
  "purpose": "failover",
  "reimaged": true
}
]
</pre>

Since this data is now homogeneous (rectangular), it pretty-prints nicely:

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --opprint unsparsify data/het/sparse.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
host      status  volume    purpose  reimaged
xy01.east running /dev/sda1 -        -
xy92.west running -         -        -
xy55.east -       /dev/sda1 failover true
</pre>

## Reading and writing heterogeneous data

In the previous sections we saw different kinds of data heterogeneity, and ways
to transform the data to make it homogeneous.

### Non-rectangular file formats: JSON, XTAB, NIDX, DKVP

For these formats, record-heterogeneity comes naturally:

<pre class="pre-highlight-in-pair">
<b>cat data/het/sparse.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "host": "xy01.east",
  "status": "running",
  "volume": "/dev/sda1"
}
{
  "host": "xy92.west",
  "status": "running"
}
{
  "purpose": "failover",
  "host": "xy55.east",
  "volume": "/dev/sda1",
  "reimaged": true
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --onidx --ofs ' ' cat data/het/sparse.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
xy01.east running /dev/sda1
xy92.west running
failover xy55.east /dev/sda1 true
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --oxtab cat data/het/sparse.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
host   xy01.east
status running
volume /dev/sda1

host   xy92.west
status running

purpose  failover
host     xy55.east
volume   /dev/sda1
reimaged true
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --odkvp cat data/het/sparse.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
host=xy01.east,status=running,volume=/dev/sda1
host=xy92.west,status=running
purpose=failover,host=xy55.east,volume=/dev/sda1,reimaged=true
</pre>

Even then, we may wish to put like with like, using the [`group-like`](reference-verbs.md#group-like) verb:

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --odkvp cat data/het.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
resource=/path/to/file,loadsec=0.45,ok=true
record_count=100,resource=/path/to/file
resource=/path/to/second/file,loadsec=0.32,ok=true
record_count=150,resource=/path/to/second/file
resource=/some/other/path,loadsec=0.97,ok=false
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --odkvp group-like data/het.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
resource=/path/to/file,loadsec=0.45,ok=true
resource=/path/to/second/file,loadsec=0.32,ok=true
resource=/some/other/path,loadsec=0.97,ok=false
record_count=100,resource=/path/to/file
record_count=150,resource=/path/to/second/file
</pre>

### Rectangular file formats: CSV and pretty-print

CSV and pretty-print formats expect rectangular structure. But Miller lets you
process non-rectangular using CSV and pretty-print.

Miller simply prints a newline and a new header when there is a schema change
-- where by _schema_ we mean simply the list of record keys in the order they
are encountered. When there is no schema change, you get CSV per se as a
special case. Likewise, Miller reads heterogeneous CSV or pretty-print input
the same way. The difference between CSV and CSV-lite is that the former is
[RFC-4180-compliant](file-formats.md#csvtsvasvusvetc), while the latter readily
handles heterogeneous data (which is non-compliant). For example:

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

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --opprint group-like data/het.json</b>
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

Miller handles explicit header changes as just shown. If your CSV input contains ragged data -- if there are implicit header changes (no intervening blank line and new header line) as seen above -- you can use `--allow-ragged-csv-input` (or keystroke-saver `--ragged`).

<pre class="pre-highlight-in-pair">
<b>mlr --csv --ragged cat data/het/ragged.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c
1,2,3
4,5,

a,b,c,4
7,8,9,10
</pre>

## Processing heterogeneous data

Above we saw how to make heterogeneous data homogeneous, and then how to print heterogeneous data.
As for other processing, record-heterogeneity is not a problem for Miller.

Miller operates on specified fields and takes the rest along: for example, if
you are sorting on the `count` field then all records in the input stream must
have a `count` field but the other fields can vary, and moreover the sorted-on
field name(s) don't need to be in the same position on each line:

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
