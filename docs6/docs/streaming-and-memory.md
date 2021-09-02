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
# Streaming processing, and memory usage

## What does streaming mean?

When we say that Miller is **streaming**, we mean that most operations need
only a single record in memory at a time, rather than ingesting all input
before producing any output.

This is contrast to, say, the dataframes approach where you ingest all data,
wait for end of file, then start manipulating the data.

Both approaches have their advantages: the dataframes approach requires that
all data fit in system memory (which, as hardware gets larger over time, is
less and less of a constraint); the streaming approach requires that you
sometimes need to accumulate results on records (rows) as they arrive rather
than looping through them explicitly.

Since Miller takes the streaming approach when possible (see below for
exceptions), you can often operate on files which are larger than your system's
memory . It also means you can do `tail -f some-file | mlr --some-flags` and
Miller will operate on records as they arrive one at a time.  You don't have
to wait for and end-of-file marker (which never arrives with `tail-f`) to
start seeing partial results. This also means if you pipe Miller's output
to other streaming tools (like `cat`, `grep`, `sed`, and so on), they
will also output partial results as data arrives.

The statements in the [Miller programming language](programming-language.md)
(outside of optional `begin`/`end` blocks which execute before and after all
records have been read, respectively) are implicit callbacks which are executed
once per record.

If you do wish to accumulate all records into memory and loop over them
explicitly, you can do so -- see the page on [operating on all
records](operating-on-all-records.md).

## Streaming and non-streaming verbs

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

Yet other verbs, such as
[`stats1`](reference-verbs.md#stats1) and
[`stats2`](reference-verbs.md#stats2), retain only summary arithmetic on the
records they visit.

| Verb      | Description                          |
| ----------- | ------------------------------------ |
| `foo`       | TODO  |
| `bar`       | TODO |
| `sort`      | Non-streaming: retains all records, then emits sorted data after end of input stream |
| `tac`       | Non-streaming: retains all records, then emits reversed data after end of input stream |

TODO
