**Miller is like sed, awk, cut, join, and sort for name-indexed data such as CSV.**

With Miller, you get to use named fields without needing to count positional
indices.

For example:

```
% mlr --csv cut -f hostname,uptime mydata.csv
% mlr --csv sort -f hostname,uptime mydata.csv
% mlr --csv put '$z = $x + 2.7*$y' mydata.csv
% mlr --csv filter '$status != "down"' mydata.csv
```

This is something the Unix toolkit always could have done, and arguably always
should have done.  It operates on **key-value-pair data** while the familiar
Unix tools operate on integer-indexed fields: if the natural data structure for
the latter is the array, then Miller's natural data structure is the
insertion-ordered hash map.  This encompasses a **variety of data formats**,
including but not limited to the familiar **CSV**.  (Miller can handle
positionally-indexed data as a special case.) 

Features:

* I/O formats including **tabular pretty-printing**

* **Conversion** between formats

* **Format-aware processing**: e.g. CSV `sort` and `tac` keep header lines first

* High-throughput **performance** on par with the Unix toolkit

* Miller is **pipe-friendly** and interoperates with Unix toolkit

* Miller is **streaming**: most operations need only a single record in
memory at a time, rather than ingesting all input before producing any output.
For those operations which require deeper retention (`sort`, `tac`, `stats1`),
Miller retains only as much data as needed. This means that whenever
functionally possible, you can operate on files which are larger than your
system&rsquo;s available RAM, and you can use Miller in `tail -f`
contexts.

* It complements SQL **databases**: you can slice, dice, and reformat data on
the client side on its way into or out of a database. You can also reap some of
the benefits of databases for quick, setup-free one-off tasks when just need to
query some data in disk files in a hurry.

* Miller also goes beyond the classic Unix tools by stepping into our modern,
**no-SQL** world: its essential record-heterogeneity property allows Miller to
operate on data where records with different schema (field names) are
interleaved.

* Not unlike `jq` (http://stedolan.github.io/jq/) for JSON, Miller is written
in modern C, and it has **zero runtime dependencies**. You can download or
compile a single binary, `scp` it to a faraway machine, and expect it to work.

For more information, please visit http://johnkerl.org/miller/doc

For notes about issue-labeling in the Github repo, please see
https://github.com/johnkerl/miller/wiki/Issue-labeling.
