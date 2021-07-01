..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Quick examples
================================================================

Column select:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --csv cut -f hostname,uptime mydata.csv

Add new columns as function of other columns:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --nidx put '$sum = $7 < 0.0 ? 3.5 : $7 + 2.1*$8' *.dat

Row filter:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --csv filter '$status != "down" && $upsec >= 10000' *.csv

Apply column labels and pretty-print:

.. code-block:: none
   :emphasize-lines: 1-1

    grep -v '^#' /etc/group | mlr --ifs : --nidx --opprint label group,pass,gid,member then sort -f group

Join multiple data sources on key columns:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr join -j account_id -f accounts.dat then group-by account_name balances.dat

Mulltiple formats including JSON:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --json put '$attr = sub($attr, "([0-9]+)_([0-9]+)_.*", "\1:\2")' data/*.json

Aggregate per-column statistics:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr stats1 -a min,mean,max,p10,p50,p90 -f flag,u,v data/*

Linear regression:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr stats2 -a linreg-pca -f u,v -g shape data/*

Aggregate custom per-column statistics:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr put -q '@sum[$a][$b] += $x; end {emit @sum, "a", "b"}' data/*

Iterate over data using DSL expressions:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --from estimates.tbl put '
      for (k,v in $*) {
        if (is_numeric(v) && k =~ "^[t-z].*$") {
          $sum += v; $count += 1
        }
      }
      $mean = $sum / $count # no assignment if count unset
    '

Run DSL expressions from a script file:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --from infile.dat put -f analyze.mlr

Split/reduce output to multiple filenames:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --from infile.dat put 'tee > "./taps/data-".$a."-".$b, $*'

Compressed I/O:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --from infile.dat put 'tee | "gzip > ./taps/data-".$a."-".$b.".gz", $*'

Interoperate with other data-processing tools using standard pipes:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --from infile.dat put -q '@v=$*; dump | "jq .[]"'

Tap/trace:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --from infile.dat put  '(NR % 1000 == 0) { print > stderr, "Checkpoint ".NR}'
