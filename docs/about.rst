About
=====

Miller is like awk, sed, cut, join, and sort for **name-indexed data such as
CSV, TSV, and tabular JSON**. You get to work with your data using named
fields, without needing to count positional column indices.

This is something the Unix toolkit always could have done, and arguably
always should have done.  It operates on key-value-pair data while the familiar
Unix tools operate on integer-indexed fields: if the natural data structure for
the latter is the array, then Miller's natural data structure is the
insertion-ordered hash map.  This encompasses a **variety of data formats**,
including but not limited to the familiar CSV, TSV, and JSON.  (Miller can handle
**positionally-indexed data** as a special case.)

Features
^^^^^^^^

* Miller is **multi-purpose**: it's useful for **data cleaning**, **data reduction**, **statistical reporting**, **devops**, **system administration**, **log-file processing**, **format conversion**, and **database-query post-processing**.

* You can use Miller to snarf and munge **log-file data**, including selecting out relevant substreams, then produce CSV format and load that into all-in-memory/data-frame utilities for further statistical and/or graphical processing.

* Miller complements **data-analysis tools** such as **R**, **pandas**, etc.: you can use Miller to **clean** and **prepare** your data. While you can do **basic statistics** entirely in Miller, its streaming-data feature and single-pass algorithms enable you to **reduce very large data sets**.

* Miller complements SQL **databases**: you can slice, dice, and reformat data on the client side on its way into or out of a database.  (Examples <a href="10-min.html#SQL-input_examples">here</a> and <a href="10-min.html#SQL-output_examples">here</a>). You can also reap some of the benefits of databases for quick, setup-free one-off tasks when you just need to query some data in disk files in a hurry.

* Miller also goes beyond the classic Unix tools by stepping fully into our modern, **no-SQL** world: its essential record-heterogeneity property allows Miller to operate on data where records with different schema (field names) are interleaved.

* Miller is **streaming**: most operations need only a single record in memory at a time, rather than ingesting all input before producing any output.  For those operations which require deeper retention (``sort``, ``tac``, ``stats1``), Miller retains only as much data as needed.  This means that whenever functionally possible, you can operate on files which are larger than your system's available RAM, and you can use Miller in **tail -f** contexts.

* Miller is **pipe-friendly** and interoperates with the Unix toolkit

* Miller's I/O formats include **tabular pretty-printing**, **positionally indexed** (Unix-toolkit style), CSV, JSON, and others

* Miller does **conversion** between formats

* Miller's **processing is format-aware**: e.g. CSV ``sort`` and ``tac`` keep header lines first

* Miller has high-throughput **performance** on par with the Unix toolkit

* Not unlike <a href="http://stedolan.github.io/jq/">jq</a> (for JSON), Miller is written in portable, modern C, with **zero runtime dependencies**.  You can download or compile a single binary, ``scp`` it to a faraway machine, and expect it to work.

Releases and release notes: <a href="https://github.com/johnkerl/miller/releases">https://github.com/johnkerl/miller/releases</a>.

Examples
^^^^^^^^

Column select::

    % mlr --csv cut -f hostname,uptime mydata.csv

Add new columns as function of other columns::

    % mlr --nidx put '$sum = $7 < 0.0 ? 3.5 : $7 + 2.1*$8' *.dat

Row filter::

    % mlr --csv filter '$status != "down" && $upsec >= 10000' *.csv

Apply column labels and pretty-print::

    % grep -v '^#' /etc/group | mlr --ifs : --nidx --opprint label group,pass,gid,member then sort -f group

Join multiple data sources on key columns::

    % mlr join -j account_id -f accounts.dat then group-by account_name balances.dat

Multiple formats including JSON::

    % mlr --json put '$attr = sub($attr, "([0-9]+)_([0-9]+)_.*", "\1:\2")' data/*.json

Aggregate per-column statistics::

    % mlr stats1 -a min,mean,max,p10,p50,p90 -f flag,u,v data/*

Linear regression::

    % mlr stats2 -a linreg-pca -f u,v -g shape data/*

Aggregate custom per-column statistics::

    % mlr put -q '@sum[$a][$b] += $x; end {emit @sum, "a", "b"}' data/*

Iterate over data using DSL expressions::

    % mlr --from estimates.tbl put '
      for (k,v in $*) {
        if (is_numeric(v) && k =~ "^[t-z].*$") {
          $sum += v; $count += 1
        }
      }
      $mean = $sum / $count # no assignment if count unset
    '

Run DSL expressions from a script file::

    % mlr --from infile.dat put -f analyze.mlr

Split/reduce output to multiple filenames::

    % mlr --from infile.dat put 'tee > "./taps/data-".$a."-".$b, $*'

Compressed I/O::

    % mlr --from infile.dat put 'tee | "gzip > ./taps/data-".$a."-".$b.".gz", $*'

Interoperate with other data-processing tools using standard pipes::

    % mlr --from infile.dat put -q '@v=$*; dump | "jq .[]"'

Tap/trace::

    % mlr --from infile.dat put  '(NR % 1000 == 0) { print > stderr, "Checkpoint ".NR}'
