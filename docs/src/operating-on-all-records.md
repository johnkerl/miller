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
# Operating on all records

As we saw in the DSL-overview page, the Miller programming language has an
[implicit loop over records for main statements](reference-dsl.md#implicit-loop-over-records-for-main-statements).

Miller's feature of [_streaming operation over
records_](streaming-and-memory.md) is implemented by the main statements
(everything outside `begin`/`end`/`func`/`subr`) getting invoked once per
record. You don't explicitly loop over records, as you would in some dataframes
contexts; rather, _Miller loops over records for you_, and it lets you specify
what to do on each record: you write the body of the loop.

That's fine for most simple use-cases, but sometimes you _do_ want to loop over
all records. Here we describe a few options.

## Sums/counters

The first option is to leverage the fact that main DSL statements are already
invoked in a loop over records, and use
[out-of-stream variables](reference-dsl-variables.md#out-of-stream-variables)
to retain sums, counters, etc.

For example, let's look at our short data file [data/short.csv](data/short.csv):

<pre class="pre-highlight-in-pair">
<b>cat data/short.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
word,value
apple,37
ball,28
cat,54
</pre>

We can track count and sum using
[out-of-stream variables](reference-dsl-variables.md#out-of-stream-variables) -- the ones that
start with the `@` sigil -- then
[emit](reference-dsl-output-statements.md#emit-statements) them as a new record
after all the input is read.

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --from data/short.csv put '</b>
<b>  begin {</b>
<b>    @count = 0;</b>
<b>    @sum = 0;</b>
<b>  }</b>
<b>  @count += 1;</b>
<b>  @sum += $value;</b>
<b>  end {</b>
<b>    emit (@count, @sum);</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "word": "apple",
  "value": 37
},
{
  "word": "ball",
  "value": 28
},
{
  "word": "cat",
  "value": 54
},
{
  "count": 3,
  "sum": 119
}
]
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

And if all we want is the final output and not the input data, we can use `put
-q` to not pass through the input records:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --from data/short.csv put -q '</b>
<b>  begin {</b>
<b>    @count = 0;</b>
<b>    @sum = 0;</b>
<b>  }</b>
<b>  @count += 1;</b>
<b>  @sum += $value;</b>
<b>  end {</b>
<b>    emit (@count, @sum);</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "count": 3,
  "sum": 119
}
]
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

As discussed a bit more on the page on [streaming processing and memory
usage](streaming-and-memory.md), this doesn't keep all records in memory, only
the count and sum variables. You can use this on very large files without
running out of memory.

## Retaining records in a map

The second option is to retain entire records in a [map](reference-main-maps.md), then loop over them in an `end` block.

Let's use the same short data file [data/short.csv](data/short.csv):

<pre class="pre-highlight-in-pair">
<b>cat data/short.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
word,value
apple,37
ball,28
cat,54
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --from data/short.csv put -q '</b>
<b>  # map</b>
<b>  begin {</b>
<b>    @records = {};</b>
<b>  }</b>
<b>  @records[NR] = $*;</b>
<b>  end {</b>
<b>    count = length(@records);</b>
<b>    sum = 0;</b>
<b>    for (i = 1; i <= NR; i += 1) {</b>
<b>      sum += @records[i]["value"];</b>
<b>    }</b>
<b>    dump @records; # show the map</b>
<b>    emit (count, sum);</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "1": {
    "word": "apple",
    "value": 37
  },
  "2": {
    "word": "ball",
    "value": 28
  },
  "3": {
    "word": "cat",
    "value": 54
  }
}
[
{
  "count": 3,
  "sum": 119
}
]
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

The downside to this, of course, is that this retains all records (plus data-structure overhead) in memory, so you're limited to processing files that fit in your computer's memory. The upside, though, is that you can do random access over the records using things like

<pre class="pre-non-highlight-non-pair">
    output = 0;
    for (i = 1; i <= NR; i += 1) {
      for (j = 1; j <= NR; j += 1) {
        for (k = 1; k <= NR; k += 1) {
          output += call_some_function_of(@records[i], @records[j], @record[k])
        }
      }
    }
    # do something with the output
</pre>

## Retaining records in an array

The third option is to retain records in an [array](reference-main-arrays.md), then loop over them in an `end` block.

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --from data/short.csv put -q '</b>
<b>  # array</b>
<b>  begin {</b>
<b>    @records = [];</b>
<b>  }</b>
<b>  @records[NR] = $*;</b>
<b>  end {</b>
<b>    count = length(@records);</b>
<b>    sum = 0;</b>
<b>    for (i = 1; i <= NR; i += 1) {</b>
<b>      sum += @records[i]["value"];</b>
<b>    }</b>
<b>    dump @records; # show the array</b>
<b>    emit (count, sum);</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
  {
    "word": "apple",
    "value": 37
  },
  {
    "word": "ball",
    "value": 28
  },
  {
    "word": "cat",
    "value": 54
  }
]
[
{
  "count": 3,
  "sum": 119
}
]
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

Just as with the retain-as-map approach, the downside is the overhead of
retaining all records in memory, and the upside is that you get random access
over records.

## Using maps vs using arrays

Retaining records as a map or as an array is a matter of taste. Some things to note:

If we initialize `@records = {}` in the `begin` block (or, if we don't initialize it at all and just start writing to it in the main statements) then `@records` is a [map](reference-main-maps.md) . If we initialize `@records=[]` then it's an array.

Arrays are, of course, contiguously indexed. (And, in Miller, their indices
start with 1, not 0 as discussed in the [Arrays](reference-main-arrays.md)
page.) This means that if you are only retaining a subset of records then your
array will have [null-gaps](reference-main-arrays.md) in it:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --from data/short.csv put -q '</b>
<b>  begin {</b>
<b>    @records = [];</b>
<b>  }</b>
<b>  if (NR != 2) {</b>
<b>    @records[NR] = $*</b>
<b>  }</b>
<b>  end {</b>
<b>    dump @records;</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
  {
    "word": "apple",
    "value": 37
  },
  null,
  {
    "word": "cat",
    "value": 54
  }
]
[
]
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

You can index `@records` by `@count` rather than `NR` to get a contiguous array:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --from data/short.csv put -q '</b>
<b>  begin {</b>
<b>    @records = [];</b>
<b>    @count = 0;</b>
<b>  }</b>
<b>  # main statement</b>
<b>  if (NR != 2) {</b>
<b>    @count += 1;</b>
<b>    @records[@count] = $*;</b>
<b>  }</b>
<b>  end {</b>
<b>    dump @records;</b>
<b>    count = length(@records);</b>
<b>    sum = 0;</b>
<b>    for (record in @records) {</b>
<b>      sum += record["value"];</b>
<b>    }</b>
<b>    emit (count, sum);</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
  {
    "word": "apple",
    "value": 37
  },
  {
    "word": "cat",
    "value": 54
  }
]
[
{
  "count": 2,
  "sum": 91
}
]
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

If you use a map to retain records, then this is a non-issue: maps can retain whatever values you like:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --from data/short.csv put -q '</b>
<b>  begin {</b>
<b>    @records = {};</b>
<b>  }</b>
<b>  # main statement</b>
<b>  if (NR != 2) {</b>
<b>    @records[NR] = $*;</b>
<b>  }</b>
<b>  end {</b>
<b>    dump @records;</b>
<b>    count = length(@records);</b>
<b>    sum = 0;</b>
<b>    for (key in @records) {</b>
<b>      sum += @records[key]["value"];</b>
<b>    }</b>
<b>    emit (count, sum);</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "1": {
    "word": "apple",
    "value": 37
  },
  "3": {
    "word": "cat",
    "value": 54
  }
}
[
{
  "count": 2,
  "sum": 91
}
]
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

Do note that Miller [maps](reference-main-maps.md) preserve insertion order, so
at the end you're guaranteed to loop over records in the same order you read
them. Also note that when you index a Miller map with an integer key, this
works, but the [key is stringified](reference-main-maps.md).

## Retaining partial records in map or array

If all you need is one or a few attributes out of a record, you don't need to
retain full records. You can retain a map, or array, of just the fields you're
interested in:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --from data/short.csv put -q '</b>
<b>  begin {</b>
<b>    @values = {};</b>
<b>  }</b>
<b>  # main statement</b>
<b>  if (NR != 2) {</b>
<b>    @values[NR] = $value;</b>
<b>  }</b>
<b>  end {</b>
<b>    dump @values;</b>
<b>    count = length(@values);</b>
<b>    sum = 0;</b>
<b>    for (key in @values) {</b>
<b>      sum += @values[key];</b>
<b>    }</b>
<b>    emit (count, sum);</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "1": 37,
  "3": 54
}
[
{
  "count": 2,
  "sum": 91
}
]
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

## Sorting

Please see the [sorting page](sorting.md).

## For more information

Please see the page on [two-pass algorithms](two-pass-algorithms.md); see also
the page on [higher-order functions](reference-dsl-higher-order-functions.md).
