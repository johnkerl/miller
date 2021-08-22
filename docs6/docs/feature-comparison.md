<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# Unix-toolkit context

How does Miller fit within the Unix toolkit (`grep`, `sed`, `awk`, etc.)?

## File-format awareness

Miller respects CSV headers. If you do `mlr --csv cat *.csv` then the header line is written once:

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
<b>mlr --csv cat data/a.csv data/b.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c
1,2,3
4,5,6
7,8,9
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv sort -nr b data/a.csv data/b.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c
7,8,9
4,5,6
1,2,3
</pre>

Likewise with `mlr sort`, `mlr tac`, and so on.

## awk-like features: mlr filter and mlr put

* `mlr filter` includes/excludes records based on a filter expression, e.g. `mlr filter '$count > 10'`.

* `mlr put` adds a new field as a function of others, e.g. `mlr put '$xy = $x * $y'` or `mlr put '$counter = NR'`.

* The `$name` syntax is straight from `awk`'s `$1 $2 $3` (adapted to name-based indexing), as are the variables `FS`, `OFS`, `RS`, `ORS`, `NF`, `NR`, and `FILENAME`. The `ENV[...]` syntax is from Ruby.

* While `awk` functions are record-based, Miller subcommands (or *verbs*) are stream-based: each of them maps a stream of records into another stream of records.

* Like `awk`, Miller (as of v5.0.0) allows you to define new functions within its `put` and `filter` expression language.  Further programmability comes from chaining with `then`.

* As with `awk`, `$`-variables are stream variables and all verbs (such as `cut`, `stats1`, `put`, etc.) as well as `put`/`filter` statements operate on streams.  This means that you define actions to be done on each record and then stream your data through those actions.  The built-in variables `NF`, `NR`, etc.  change from one record to another, `$x` is a label for field `x` in the current record, and the input to `sqrt($x)` changes from one record to the next.  The expression language for the `put` and `filter` verbs additionally allows you to define `begin {...}` and `end {...}` blocks for actions to be taken before and after records are processed, respectively.

* As with `awk`, Miller's `put`/`filter` language lets you set `@sum=0` before records are read, then update that sum on each record, then print its value at the end.  Unlike `awk`, Miller makes syntactically explicit the difference between variables with extent across all records (names starting with `@`, such as `@sum`) and variables which are local to the current expression invocation (names starting without `@`, such as `sum`).

## See also

See [Verbs Reference](reference-verbs.md) for more on Miller's subcommands `cat`, `cut`, `head`, `sort`, `tac`, `tail`, `top`, and `uniq`, as well as [DSL reference](reference-dsl.md) for more on the awk-like `mlr filter` and `mlr put`.
