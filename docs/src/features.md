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
# Features

Miller is like awk, sed, cut, join, and sort for **name-indexed data, such as
CSV, TSV, JSON, and JSON Lines**. You get to work with your data using named
fields, without needing to count positional column indices.

This is something the Unix toolkit always could have done, and arguably always
should have done.  It operates on key-value-pair data while the familiar Unix
tools operate on integer-indexed fields: if the natural data structure for the
latter is the array, then Miller's natural data structure is the
insertion-ordered hash map.  This encompasses a **variety of data formats**,
including but not limited to the familiar CSV, TSV, JSON, and JSON Lines.
(Miller can handle **positionally-indexed data** as a special case.)

* Miller is **multi-purpose**: it's useful for **data cleaning**, **data reduction**, **statistical reporting**, **devops**, **system administration**, **log-file processing**, **format conversion**, and **database-query post-processing**.

* You can use Miller to snarf and munge **log-file data**, including selecting out relevant substreams, then produce CSV format and load that into all-in-memory/data-frame utilities for further statistical and/or graphical processing.

* Miller complements **data-analysis tools** such as **R**, **pandas**, etc.: you can use Miller to **clean** and **prepare** your data. While you can do **basic statistics** entirely in Miller, its streaming-data feature and single-pass algorithms enable you to **reduce very large data sets**.

* Miller complements SQL **databases**: you can slice, dice, and reformat data on the client side on its way into or out of a database.  (See [SQL Examples](sql-examples.md).) You can also reap some of the benefits of databases for quick, setup-free one-off tasks when you just need to query some data in disk files in a hurry.

* Miller also goes beyond the classic Unix tools by stepping fully into our modern, **no-SQL** world: its essential record-heterogeneity property allows Miller to operate on data where records with different schemas (field names) are interleaved.

* Miller is **streaming**: most operations need only a single record in memory at a time, rather than ingesting all input before producing any output.  For those operations that require deeper retention (`sort`, `tac`, `stats1`), Miller retains only as much data as needed.  This means that whenever functionally possible, you can operate on files that are larger than your system's available RAM, and you can use Miller in **tail -f** contexts.

* Miller is **pipe-friendly** and interoperates with the Unix toolkit

* Miller's I/O formats include **tabular pretty-printing**, **positionally indexed** (Unix-toolkit style), CSV, JSON, and others

* Miller does **conversion** between formats

* Miller's **processing is format-aware**: e.g., CSV `sort` and `tac` keep header lines first

* Miller has high-throughput **performance** on par with the Unix toolkit

* Not unlike [jq](https://stedolan.github.io/jq/) (for JSON), Miller is written in Go, which is a portable, modern language, and Miller has no runtime dependencies.  You can download or compile a single binary, `scp` it to a faraway machine, and expect it to work.

Releases and release notes: [https://github.com/johnkerl/miller/releases](https://github.com/johnkerl/miller/releases).
