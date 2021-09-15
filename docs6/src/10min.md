<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-main-flag-list/index.html">Flag list</a>
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
# Miller in 10 minutes

## Miller verbs

Let's take a quick look at some of the most useful Miller verbs -- file-format-aware, name-index-empowered equivalents of standard system commands.

`mlr cat` is like system `cat` (or `type` on Windows) -- it passes the data through unmodified:

<pre class="pre-highlight-in-pair">
<b>mlr --csv cat example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
yellow,triangle,true,1,11,43.6498,9.8870
red,square,true,2,15,79.2778,0.0130
red,circle,true,3,16,13.8103,2.9010
red,square,false,4,48,77.5542,7.4670
purple,triangle,false,5,51,81.2290,8.5910
red,square,false,6,64,77.1991,9.5310
purple,triangle,false,7,65,80.1405,5.8240
yellow,circle,true,8,73,63.9785,4.2370
yellow,circle,true,9,87,63.5058,8.3350
purple,square,false,10,91,72.3735,8.2430
</pre>

But `mlr cat` can also do format conversion -- for example, you can pretty-print in tabular format:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint cat example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate
yellow triangle true  1  11    43.6498  9.8870
red    square   true  2  15    79.2778  0.0130
red    circle   true  3  16    13.8103  2.9010
red    square   false 4  48    77.5542  7.4670
purple triangle false 5  51    81.2290  8.5910
red    square   false 6  64    77.1991  9.5310
purple triangle false 7  65    80.1405  5.8240
yellow circle   true  8  73    63.9785  4.2370
yellow circle   true  9  87    63.5058  8.3350
purple square   false 10 91    72.3735  8.2430
</pre>

`mlr head` and `mlr tail` count records rather than lines. Whether you're getting the first few records or the last few, the CSV header is included either way:

<pre class="pre-highlight-in-pair">
<b>mlr --csv head -n 4 example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
yellow,triangle,true,1,11,43.6498,9.8870
red,square,true,2,15,79.2778,0.0130
red,circle,true,3,16,13.8103,2.9010
red,square,false,4,48,77.5542,7.4670
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv tail -n 4 example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
purple,triangle,false,7,65,80.1405,5.8240
yellow,circle,true,8,73,63.9785,4.2370
yellow,circle,true,9,87,63.5058,8.3350
purple,square,false,10,91,72.3735,8.2430
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson tail -n 2 example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "color": "yellow",
  "shape": "circle",
  "flag": "true",
  "k": 9,
  "index": 87,
  "quantity": 63.5058,
  "rate": 8.3350
}
{
  "color": "purple",
  "shape": "square",
  "flag": "false",
  "k": 10,
  "index": 91,
  "quantity": 72.3735,
  "rate": 8.2430
}
</pre>

You can sort on a single field:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint sort -f shape example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate
red    circle   true  3  16    13.8103  2.9010
yellow circle   true  8  73    63.9785  4.2370
yellow circle   true  9  87    63.5058  8.3350
red    square   true  2  15    79.2778  0.0130
red    square   false 4  48    77.5542  7.4670
red    square   false 6  64    77.1991  9.5310
purple square   false 10 91    72.3735  8.2430
yellow triangle true  1  11    43.6498  9.8870
purple triangle false 5  51    81.2290  8.5910
purple triangle false 7  65    80.1405  5.8240
</pre>

Or, you can sort primarily alphabetically on one field, then secondarily numerically descending on another field, and so on:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint sort -f shape -nr index example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate
yellow circle   true  9  87    63.5058  8.3350
yellow circle   true  8  73    63.9785  4.2370
red    circle   true  3  16    13.8103  2.9010
purple square   false 10 91    72.3735  8.2430
red    square   false 6  64    77.1991  9.5310
red    square   false 4  48    77.5542  7.4670
red    square   true  2  15    79.2778  0.0130
purple triangle false 7  65    80.1405  5.8240
purple triangle false 5  51    81.2290  8.5910
yellow triangle true  1  11    43.6498  9.8870
</pre>

If there are fields you don't want to see in your data, you can use `cut` to keep only the ones you want, in the same order they appeared in the input data:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint cut -f flag,shape example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
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
</pre>

You can also use `cut -o` to keep specified fields, but in your preferred order:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint cut -o -f flag,shape example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
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
</pre>

You can use `cut -x` to omit fields you don't care about:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint cut -x -f flag,shape example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  k  index quantity rate
yellow 1  11    43.6498  9.8870
red    2  15    79.2778  0.0130
red    3  16    13.8103  2.9010
red    4  48    77.5542  7.4670
purple 5  51    81.2290  8.5910
red    6  64    77.1991  9.5310
purple 7  65    80.1405  5.8240
yellow 8  73    63.9785  4.2370
yellow 9  87    63.5058  8.3350
purple 10 91    72.3735  8.2430
</pre>

Even though Miller's main selling point is name-indexing, sometimes you really want to refer to a field name by its positional index. Use `$[[3]]` to access the name of field 3 or `$[[[3]]]` to access the value of field 3:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint put '$[[3]] = "NEW"' example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    NEW   k  index quantity rate
yellow triangle true  1  11    43.6498  9.8870
red    square   true  2  15    79.2778  0.0130
red    circle   true  3  16    13.8103  2.9010
red    square   false 4  48    77.5542  7.4670
purple triangle false 5  51    81.2290  8.5910
red    square   false 6  64    77.1991  9.5310
purple triangle false 7  65    80.1405  5.8240
yellow circle   true  8  73    63.9785  4.2370
yellow circle   true  9  87    63.5058  8.3350
purple square   false 10 91    72.3735  8.2430
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint put '$[[[3]]] = "NEW"' example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag k  index quantity rate
yellow triangle NEW  1  11    43.6498  9.8870
red    square   NEW  2  15    79.2778  0.0130
red    circle   NEW  3  16    13.8103  2.9010
red    square   NEW  4  48    77.5542  7.4670
purple triangle NEW  5  51    81.2290  8.5910
red    square   NEW  6  64    77.1991  9.5310
purple triangle NEW  7  65    80.1405  5.8240
yellow circle   NEW  8  73    63.9785  4.2370
yellow circle   NEW  9  87    63.5058  8.3350
purple square   NEW  10 91    72.3735  8.2430
</pre>

You can find the full list of verbs at the [Verbs Reference](reference-verbs.md) page.

## Filtering

You can use `filter` to keep only records you care about:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint filter '$color == "red"' example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color shape  flag  k index quantity rate
red   square true  2 15    79.2778  0.0130
red   circle true  3 16    13.8103  2.9010
red   square false 4 48    77.5542  7.4670
red   square false 6 64    77.1991  9.5310
</pre>

<pre class="pre-highlight-non-pair">
<b>mlr --icsv --opprint filter '$color == "red" && $flag == true' example.csv</b>
</pre>

## Computing new fields

You can use `put` to create new fields which are computed from other fields:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint put '</b>
<b>  $ratio = $quantity / $rate;</b>
<b>  $color_shape = $color . "_" . $shape</b>
<b>' example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate   ratio              color_shape
yellow triangle true  1  11    43.6498  9.8870 4.414868008496004  yellow_triangle
red    square   true  2  15    79.2778  0.0130 6098.292307692308  red_square
red    circle   true  3  16    13.8103  2.9010 4.760530851430541  red_circle
red    square   false 4  48    77.5542  7.4670 10.386259541984733 red_square
purple triangle false 5  51    81.2290  8.5910 9.455127458968688  purple_triangle
red    square   false 6  64    77.1991  9.5310 8.099790158430384  red_square
purple triangle false 7  65    80.1405  5.8240 13.760388049450551 purple_triangle
yellow circle   true  8  73    63.9785  4.2370 15.09995279679018  yellow_circle
yellow circle   true  9  87    63.5058  8.3350 7.619172165566886  yellow_circle
purple square   false 10 91    72.3735  8.2430 8.779995147397793  purple_square
</pre>

When you create a new field, it can immediately be used in subsequent statements:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint --from example.csv put '</b>
<b>  $y = $index + 1;</b>
<b>  $z = $y**2 + $k;</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate   y  z
yellow triangle true  1  11    43.6498  9.8870 12 145
red    square   true  2  15    79.2778  0.0130 16 258
red    circle   true  3  16    13.8103  2.9010 17 292
red    square   false 4  48    77.5542  7.4670 49 2405
purple triangle false 5  51    81.2290  8.5910 52 2709
red    square   false 6  64    77.1991  9.5310 65 4231
purple triangle false 7  65    80.1405  5.8240 66 4363
yellow circle   true  8  73    63.9785  4.2370 74 5484
yellow circle   true  9  87    63.5058  8.3350 88 7753
purple square   false 10 91    72.3735  8.2430 92 8474
</pre>

For `put` and `filter` we were able to type out expressions using a programming-language syntax.
See the [Miller programming language page](miller-programming-language.md) for more information.

## Multiple input files

Miller takes all the files from the command line as an input stream. But it's format-aware, so it doesn't repeat CSV header lines. For example, with input files [data/a.csv](data/a.csv) and [data/b.csv](data/b.csv), the system `cat` command will repeat header lines:

<pre class="pre-highlight-in-pair">
<b>cat data/a.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c
1,2,3
4,5,6
</pre>

<pre class="pre-highlight-in-pair">
<b>cat data/b.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c
7,8,9
</pre>

<pre class="pre-highlight-in-pair">
<b>cat data/a.csv data/b.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c
1,2,3
4,5,6
a,b,c
7,8,9
</pre>

However, `mlr cat` will not:

<pre class="pre-highlight-in-pair">
<b>mlr --csv cat data/a.csv data/b.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c
1,2,3
4,5,6
7,8,9
</pre>

## Chaining verbs together

Often we want to chain queries together -- for example, sorting by a field and taking the top few values. We can do this using pipes:

<pre class="pre-highlight-in-pair">
<b>mlr --csv sort -nr index example.csv | mlr --icsv --opprint head -n 3</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape  flag  k  index quantity rate
purple square false 10 91    72.3735  8.2430
yellow circle true  9  87    63.5058  8.3350
yellow circle true  8  73    63.9785  4.2370
</pre>

This works fine -- but Miller also lets you chain verbs together using the word `then`. Think of this as a Miller-internal pipe that lets you use fewer keystrokes:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint sort -nr index then head -n 3 example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape  flag  k  index quantity rate
purple square false 10 91    72.3735  8.2430
yellow circle true  9  87    63.5058  8.3350
yellow circle true  8  73    63.9785  4.2370
</pre>

As another convenience, you can put the filename first using `--from`. When you're interacting with your data at the command line, this makes it easier to up-arrow and append to the previous command:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint --from example.csv sort -nr index then head -n 3</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape  flag  k  index quantity rate
purple square false 10 91    72.3735  8.2430
yellow circle true  9  87    63.5058  8.3350
yellow circle true  8  73    63.9785  4.2370
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint --from example.csv \</b>
<b>  sort -nr index \</b>
<b>  then head -n 3 \</b>
<b>  then cut -f shape,quantity</b>
</pre>
<pre class="pre-non-highlight-in-pair">
shape  quantity
square 72.3735
circle 63.5058
circle 63.9785
</pre>

## Sorts and stats

Now suppose you want to sort the data on a given column, *and then* take the top few in that ordering. You can use Miller's `then` feature to pipe commands together.

Here are the records with the top three `index` values:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint sort -nr index then head -n 3 example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape  flag  k  index quantity rate
purple square false 10 91    72.3735  8.2430
yellow circle true  9  87    63.5058  8.3350
yellow circle true  8  73    63.9785  4.2370
</pre>

Lots of Miller commands take a `-g` option for group-by: here, `head -n 1 -g shape` outputs the first record for each distinct value of the `shape` field. This means we're finding the record with highest `index` field for each distinct `shape` field:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint sort -f shape -nr index then head -n 1 -g shape example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate
yellow circle   true  9  87    63.5058  8.3350
purple square   false 10 91    72.3735  8.2430
purple triangle false 7  65    80.1405  5.8240
</pre>

Statistics can be computed with or without group-by field(s):

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint --from example.csv \</b>
<b>  stats1 -a count,min,mean,max -f quantity -g shape</b>
</pre>
<pre class="pre-non-highlight-in-pair">
shape    quantity_count quantity_min quantity_mean     quantity_max
triangle 3              43.6498      68.33976666666666 81.229
square   4              72.3735      76.60114999999999 79.2778
circle   3              13.8103      47.0982           63.9785
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint --from example.csv \</b>
<b>  stats1 -a count,min,mean,max -f quantity -g shape,color</b>
</pre>
<pre class="pre-non-highlight-in-pair">
shape    color  quantity_count quantity_min quantity_mean      quantity_max
triangle yellow 1              43.6498      43.6498            43.6498
square   red    3              77.1991      78.01036666666666  79.2778
circle   red    1              13.8103      13.8103            13.8103
triangle purple 2              80.1405      80.68475000000001  81.229
circle   yellow 2              63.5058      63.742149999999995 63.9785
square   purple 1              72.3735      72.3735            72.3735
</pre>

If your output has a lot of columns, you can use XTAB format to line things up vertically for you instead:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --oxtab --from example.csv \</b>
<b>  stats1 -a p0,p10,p25,p50,p75,p90,p99,p100 -f rate</b>
</pre>
<pre class="pre-non-highlight-in-pair">
rate_p0   0.0130
rate_p10  2.9010
rate_p25  4.2370
rate_p50  8.2430
rate_p75  8.5910
rate_p90  9.8870
rate_p99  9.8870
rate_p100 9.8870
</pre>


## File formats and format conversion

Miller supports the following formats:

* CSV (comma-separared values)
* TSV (tab-separated values)
* JSON (JavaScript Object Notation)
* PPRINT (pretty-printed tabular)
* XTAB (vertical-tabular or sideways-tabular)
* NIDX (numerically indexed, label-free, with implicit labels `"1"`, `"2"`, etc.)
* DKVP (delimited key-value pairs).

What's a CSV file, really? It's an array of rows, or *records*, each being a list of key-value pairs, or *fields*: for CSV it so happens that all the keys are shared in the header line and the values vary from one data line to another.

For example, if you have:

<pre class="pre-non-highlight-non-pair">
shape,flag,index
circle,1,24
square,0,36
</pre>

then that's a way of saying:

<pre class="pre-non-highlight-non-pair">
shape=circle,flag=1,index=24
shape=square,flag=0,index=36
</pre>

Other ways to write the same data:

<pre class="pre-non-highlight-non-pair">
CSV                   PPRINT
shape,flag,index      shape  flag index
circle,1,24           circle 1    24
square,0,36           square 0    36

JSON                  XTAB
{                     shape circle
  "shape": "circle",  flag  1
  "flag": 1,          index 24
  "index": 24         .
}                     shape square
{                     flag  0
  "shape": "square",  index 36
  "flag": 0,
  "index": 36
}

DKVP
shape=circle,flag=1,index=24
shape=square,flag=0,index=36
</pre>

Anything we can do with CSV input data, we can do with any other format input data.  And you can read from one format, do any record-processing, and output to the same format as the input, or to a different output format.

How to specify these to Miller:

* If you use `--csv` or `--json` or `--pprint`, etc., then Miller will use that format for input and output.
* If you use `--icsv` and `--ojson` (note the extra `i` and `o`) then Miller will use CSV for input and JSON for output, etc.  See also [Keystroke Savers](keystroke-savers.md) for even shorter options like `--c2j`.

You can read more about this at the [File Formats](file-formats.md) page.

If all record values are numbers, strings, etc., then converting back and forth between CSV and JSON is
a matter of specifying input-format and output-format flags:

<pre class="pre-highlight-in-pair">
<b>mlr --json cat example.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "color": "yellow",
  "shape": "triangle",
  "flag": "true",
  "k": 1,
  "index": 11,
  "quantity": 43.6498,
  "rate": 9.8870
}
{
  "color": "red",
  "shape": "square",
  "flag": "true",
  "k": 2,
  "index": 15,
  "quantity": 79.2778,
  "rate": 0.0130
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --ocsv cat example.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
yellow,triangle,true,1,11,43.6498,9.8870
red,square,true,2,15,79.2778,0.0130
</pre>

However, if JSON data has map-valued or array-valued fields, Miller gives you choices on how to
convert these to CSV columns. For example, here's some JSON data with map-valued fields:

<pre class="pre-highlight-in-pair">
<b>cat data/server-log.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "hostname": "localhost",
  "pid": 12345,
  "req": {
    "id": 6789,
    "method": "GET",
    "path": "api/check",
    "host": "foo.bar",
    "headers": {
      "host": "bar.baz",
      "user-agent": "browser"
    }
  },
  "res": {
    "status_code": 200,
    "header": {
      "content-type": "text",
      "content-encoding": "plain"
    }
  }
}
</pre>

We can convert this to CSV, or other tabular formats:

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --ocsv cat data/server-log.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
hostname,pid,req.id,req.method,req.path,req.host,req.headers.host,req.headers.user-agent,res.status_code,res.header.content-type,res.header.content-encoding
localhost,12345,6789,GET,api/check,foo.bar,bar.baz,browser,200,text,plain
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --oxtab cat data/server-log.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
hostname                    localhost
pid                         12345
req.id                      6789
req.method                  GET
req.path                    api/check
req.host                    foo.bar
req.headers.host            bar.baz
req.headers.user-agent      browser
res.status_code             200
res.header.content-type     text
res.header.content-encoding plain
</pre>

These transformations are reversible:

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --oxtab cat data/server-log.json | mlr --ixtab --ojson cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "hostname": "localhost",
  "pid": 12345,
  "req": {
    "id": 6789,
    "method": "GET",
    "path": "api/check",
    "host": "foo.bar",
    "headers": {
      "host": "bar.baz",
      "user-agent": "browser"
    }
  },
  "res": {
    "status_code": 200,
    "header": {
      "content-type": "text",
      "content-encoding": "plain"
    }
  }
}
</pre>

See the [flatten/unflatten page](flatten-unflatten.md) for more information.

## Choices for printing to files

Often we want to print output to the screen. Miller does this by default, as we've seen in the previous examples.

Sometimes, though, we want to print output to another file. Just use `> outputfilenamegoeshere` at the end of your command:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint cat example.csv > newfile.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
# Output goes to the new file;
# nothing is printed to the screen.
</pre>

<pre class="pre-highlight-in-pair">
<b>cat newfile.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
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
</pre>

Other times we just want our files to be **changed in-place**: just use `mlr -I`:

<pre class="pre-highlight-non-pair">
<b>cp example.csv newfile.txt</b>
</pre>

<pre class="pre-highlight-in-pair">
<b>cat newfile.txt</b>
</pre>
<pre class="pre-non-highlight-in-pair">
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
</pre>

<pre class="pre-highlight-non-pair">
<b>mlr -I --csv sort -f shape newfile.txt</b>
</pre>

<pre class="pre-highlight-in-pair">
<b>cat newfile.txt</b>
</pre>
<pre class="pre-non-highlight-in-pair">
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
</pre>

Also using `mlr -I` you can bulk-operate on lots of files: e.g.:

<pre class="pre-highlight-non-pair">
<b>mlr -I --csv cut -x -f unwanted_column_name *.csv</b>
</pre>

If you like, you can first copy off your original data somewhere else, before doing in-place operations.

Lastly, using `tee` within `put`, you can split your input data into separate files per one or more field names:

<pre class="pre-highlight-non-pair">
<b>mlr --csv --from example.csv put -q 'tee > $shape.".csv", $*'</b>
</pre>

<pre class="pre-highlight-in-pair">
<b>cat circle.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
red,circle,true,3,16,13.8103,2.9010
yellow,circle,true,8,73,63.9785,4.2370
yellow,circle,true,9,87,63.5058,8.3350
</pre>

<pre class="pre-highlight-in-pair">
<b>cat square.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
red,square,true,2,15,79.2778,0.0130
red,square,false,4,48,77.5542,7.4670
red,square,false,6,64,77.1991,9.5310
purple,square,false,10,91,72.3735,8.2430
</pre>

<pre class="pre-highlight-in-pair">
<b>cat triangle.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
yellow,triangle,true,1,11,43.6498,9.8870
purple,triangle,false,5,51,81.2290,8.5910
purple,triangle,false,7,65,80.1405,5.8240
</pre>
