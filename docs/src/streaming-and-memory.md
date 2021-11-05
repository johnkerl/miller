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
# Streaming processing, and memory usage

## What does streaming mean?

When we say that Miller is **streaming**, we mean that most operations need
only a single record in memory at a time, rather than ingesting all input
before producing any output.

This is contrast to, say, the dataframes approach where you ingest all data,
wait for **end of file**, then start manipulating the data.

Both approaches have their advantages: the dataframes approach requires that
all data fit in **system memory** (which, as hardware gets larger over time, is
less and less of a constraint); the streaming approach requires that you
sometimes need to accumulate results on records (rows) **as they arrive**
rather than looping through them explicitly.

Since Miller takes the streaming approach when possible (see below for
exceptions), you can often operate on files which are larger than your system's
memory . It also means you can do `tail -f some-file | mlr --some-flags` and
Miller will operate on records as they arrive one at a time.  You don't have to
wait for and end-of-file marker (which never arrives with `tail-f`) to start
seeing partial results. This also means if you pipe Miller's output to other
streaming tools (like `cat`, `grep`, `sed`, and so on), they will also output
partial results as data arrives.

The statements in the [Miller programming language](miller-programming-language.md)
(outside of optional `begin`/`end` blocks which execute before and after all
records have been read, respectively) are implicit callbacks which are executed
once per record. For example, using `mlr --csv put '$z = $x + $y' myfile.csv`,
the statement `$z = $x + $y` will be executed 10,000 times if you `myfile.csv`
has 10,000 records.

If you do wish to accumulate all records into memory and loop over them
explicitly, you can do so -- see the page on [operating on all
records](operating-on-all-records.md).

## Streaming and non-streaming verbs

Most verbs, including [`cat`](reference-verbs.md#cat),
[`cut`](reference-verbs.md#cut), etc. operate on each record independently.
They have no state to retain from one record to the next, and are streaming.

For those operations which require deeper retention, Miller retains only as
much data as needed.  For example, the [`sort`](reference-verbs.md#sort) and
[`tac`](reference-verbs.md#tac) (stream-reverse, backward spelling of
[`cat`](reference-verbs.md#cat)) must ingest and retain all records in memory
before emitting any -- the last input record may well end up being the first
one to be emitted.

[`stats1`](reference-verbs.md#stats1) Other verbs, such as
[`tail`](reference-verbs.md#tail) and [`top`](reference-verbs.md#top), need to
retain only a fixed number of records -- 10, perhaps, even if the input data
has a million records.

Yet other verbs, such as [`stats1`](reference-verbs.md#stats1) and
[`stats2`](reference-verbs.md#stats2), retain only summary arithmetic on the
records they visit. These are memory-friendly: memory usage is bounded. However,
they only produce output at the end of the record stream.

## Fully streaming verbs

These don't retain any state from one record to the next.
They are memory-friendly, and they don't wait for end of input to produce their output.

* [altkv](reference-verbs.md#altkv)
* [bar](reference-verbs.md#bar) -- if not auto-mode
* [cat](reference-verbs.md#cat)
* [check](reference-verbs.md#check)
* [clean-whitespace](reference-verbs.md#clean-whitespace)
* [cut](reference-verbs.md#cut)
* [decimate](reference-verbs.md#decimate)
* [fill-down](reference-verbs.md#fill-down)
* [fill-empty](reference-verbs.md#fill-empty)
* [flatten](reference-verbs.md#flatten)
* [format-values](reference-verbs.md#format-values)
* [gap](reference-verbs.md#gap)
* [grep](reference-verbs.md#grep)
* [having-fields](reference-verbs.md#having-fields)
* [head](reference-verbs.md#head)
* [json-parse](reference-verbs.md#json-parse)
* [json-stringify](reference-verbs.md#json-stringify)
* [label](reference-verbs.md#label)
* [merge-fields](reference-verbs.md#merge-fields)
* [nest](reference-verbs.md#nest) -- if not `implode-values-across-records`
* [nothing](reference-verbs.md#nothing)
* [regularize](reference-verbs.md#regularize)
* [rename](reference-verbs.md#rename)
* [reorder](reference-verbs.md#reorder)
* [repeat](reference-verbs.md#repeat)
* [reshape](reference-verbs.md#reshape) -- if not long-to-wide
* [sec2gmt](reference-verbs.md#sec2gmt)
* [sec2gmtdate](reference-verbs.md#sec2gmtdate)
* [seqgen](reference-verbs.md#seqgen)
* [skip-trivial-records](reference-verbs.md#skip-trivial-records)
* [sort-within-records](reference-verbs.md#sort-within-records)
* [step](reference-verbs.md#step)
* [tee](reference-verbs.md#tee)
* [template](reference-verbs.md#template)
* [unflatten](reference-verbs.md#unflatten)
* [unsparsify](reference-verbs.md#unsparsify) if invoked with `-f`

## Non-streaming, retaining all records

These retain all records from one record to the next.
They are memory-unfriendly, and they wait for end of input to produce their output.

* [bar](reference-verbs.md#bar) -- if auto-mode
* [bootstrap](reference-verbs.md#bootstrap)
* [count-similar](reference-verbs.md#count-similar)
* [fraction](reference-verbs.md#fraction)
* [group-by](reference-verbs.md#group-by)
* [group-like](reference-verbs.md#group-like)
* [least-frequent](reference-verbs.md#least-frequent)
* [most-frequent](reference-verbs.md#most-frequent)
* [nest](reference-verbs.md#nest) -- if `implode-values-across-records`
* [remove-empty-columns](reference-verbs.md#remove-empty-columns)
* [reshape](reference-verbs.md#reshape) -- if long-to-wide
* [sample](reference-verbs.md#sample)
* [shuffle](reference-verbs.md#shuffle)
* [sort](reference-verbs.md#sort)
* [tac](reference-verbs.md#tac)
* [uniq](reference-verbs.md#uniq) -- if `mlr uniq -a -c`
* [unsparsify](reference-verbs.md#unsparsify) if invoked without `-f`

## Non-streaming, retaining some records

These retain a bounded number of records from one record to the next.
They are memory-friendly, but they wait for end of input to produce their output.

* [tail](reference-verbs.md#tail)
* [top](reference-verbs.md#top)

## Non-streaming, retaining some state

These retain an amount of state from one record to the next, but less than if
they were to retain all records in memory.  They are variably memory-friendly
-- depending on how many distinct values for the group-by keys exist in the
input data -- and they wait for end of input to produce their output.

* [count-distinct](reference-verbs.md#count-distinct)
* [count](reference-verbs.md#count)
* [histogram](reference-verbs.md#histogram)
* [stats1](reference-verbs.md#stats1) -- except `mlr stats1 -s` for incremental stats before end of stream
* [stats2](reference-verbs.md#stats2)
* [uniq](reference-verbs.md#uniq) -- if not `mlr uniq -a -c`

## Variable

Any `end` blocks you provide will not be executed until end of stream; otherwise these
don't want for end of stream. Similarly, if you write logic to retain all records
(see also the page on [operating on all records](operating-on-all-records.md.in))
these will be memory-unfriendly; otherwhise they are memory-friendly.

Most simple operations such as `mlr put '$z = $x + $y'` are fully streaming.

* [filter](reference-verbs.md#filter)
* [put](reference-verbs.md#put)

## Half-streaming

The main input files are streamed, but the join file (using `-f`) is loaded into memory at the start.
