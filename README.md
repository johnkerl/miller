**Miller is like sed, awk, cut, join, and sort for name-indexed data such as CSV.**

[![Build Status](https://travis-ci.org/johnkerl/miller.svg?branch=master)](https://travis-ci.org/johnkerl/miller) [![Docs](https://img.shields.io/badge/docs-here-yellow.svg)](http://johnkerl.org/miller/doc)

With Miller, you get to use named fields without needing to count positional
indices.  Examples:

```
% mlr --csv cut -f hostname,uptime mydata.csv
% mlr --csv --rs lf filter '$status != "down" && $upsec >= 10000' *.csv
% mlr --nidx put '$sum = $7 < 0.0 ? 3.5 : $7 + 2.1*$8' *.
% grep -v '^#' /etc/group | mlr --ifs : --nidx --opprint label group,pass,gid,member then sort -f group
% mlr join -j account_id -f accounts.dat then group-by account_name balances.dat
% mlr put '$attr = sub($attr, "([0-9]+)_([0-9]+)_.*", "\1:\2")' data/*
% mlr stats1 -a min,mean,max,p10,p50,p90 -f flag,u,v data/*
% mlr stats2 -a linreg-pca -f u,v -g shape data/*
```

This is something the Unix toolkit always could have done, and arguably always
should have done.  It operates on **key-value-pair data** while the familiar
Unix tools operate on integer-indexed fields: if the natural data structure for
the latter is the array, then Miller's natural data structure is the
insertion-ordered hash map.  This encompasses a **variety of data formats**,
including but not limited to the familiar **CSV**.  (Miller can handle
positionally-indexed data as a special case.)

Features:

* I/O formats including **tabular pretty-printing** and **positionally indexed** (Unix-toolkit style)

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

* Miller complements SQL **databases**: you can slice, dice, and reformat data on
the client side on its way into or out of a database. You can also reap some of
the benefits of databases for quick, setup-free one-off tasks when just need to
query some data in disk files in a hurry.

* Miller complements **data-analysis tools** such as **R**, **pandas**, etc.:
you can use Miller to **clean** and **prepare** your data. While you can do
basic statistics entirely in Miller, its streaming-data feature and single-pass
algorithms enable you to reduce very large data sets.  You can snarf and munge
log-file data, including selecting out relevant substreams, then produce CSV
format and load that into all-in-memory/data-frame utilities for further
statistical and/or graphical processing.

* Miller also goes beyond the classic Unix tools by stepping into our modern,
**no-SQL** world: its essential record-heterogeneity property allows Miller to
operate on data where records with different schema (field names) are
interleaved.

* Not unlike `jq` (http://stedolan.github.io/jq/) for JSON, Miller is written
in modern C, and it has **zero runtime dependencies**. You can download or
compile a single binary, `scp` it to a faraway machine, and expect it to work.

Documentation:

* **For full documentation, please visit** http://johnkerl.org/miller/doc
* Miller's license is two-clause BSD: https://github.com/johnkerl/miller/blob/master/LICENSE.txt
* Build information including dependencies: http://johnkerl.org/miller/doc/build.html
* Notes about issue-labeling in the Github repo: https://github.com/johnkerl/miller/wiki/Issue-labeling
* See [here] (https://github.com/johnkerl/miller/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc) for active issues.
