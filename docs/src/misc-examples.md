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
# Miscellaneous examples

Column select:

<pre class="pre-highlight-non-pair">
<b>mlr --csv cut -f hostname,uptime mydata.csv</b>
</pre>

Add new columns as function of other columns:

<pre class="pre-highlight-non-pair">
<b>mlr --nidx put '$sum = $7 < 0.0 ? 3.5 : $7 + 2.1*$8' *.dat</b>
</pre>

Row filter:

<pre class="pre-highlight-non-pair">
<b>mlr --csv filter '$status != "down" && $upsec >= 10000' *.csv</b>
</pre>

Apply column labels and pretty-print:

<pre class="pre-highlight-non-pair">
<b>grep -v '^#' /etc/group | mlr --ifs : --nidx --opprint label group,pass,gid,member then sort -f group</b>
</pre>

Join multiple data sources on key columns:

<pre class="pre-highlight-non-pair">
<b>mlr join -j account_id -f accounts.dat then group-by account_name balances.dat</b>
</pre>

Mulltiple formats including JSON:

<pre class="pre-highlight-non-pair">
<b>mlr --json put '$attr = sub($attr, "([0-9]+)_([0-9]+)_.*", "\1:\2")' data/*.json</b>
</pre>

Aggregate per-column statistics:

<pre class="pre-highlight-non-pair">
<b>mlr stats1 -a min,mean,max,p10,p50,p90 -f flag,u,v data/*</b>
</pre>

Linear regression:

<pre class="pre-highlight-non-pair">
<b>mlr stats2 -a linreg-pca -f u,v -g shape data/*</b>
</pre>

Aggregate custom per-column statistics:

<pre class="pre-highlight-non-pair">
<b>mlr put -q '@sum[$a][$b] += $x; end {emit @sum, "a", "b"}' data/*</b>
</pre>

Iterate over data using DSL expressions:

<pre class="pre-highlight-non-pair">
<b>mlr --from estimates.tbl put '</b>
<b>  for (k,v in $*) {</b>
<b>    if (is_numeric(v) && k =~ "^[t-z].*$") {</b>
<b>      $sum += v; $count += 1</b>
<b>    }</b>
<b>  }</b>
<b>  $mean = $sum / $count # no assignment if count unset</b>
<b>'</b>
</pre>

Run DSL expressions from a script file:

<pre class="pre-highlight-non-pair">
<b>mlr --from infile.dat put -f analyze.mlr</b>
</pre>

Split/reduce output to multiple filenames:

<pre class="pre-highlight-non-pair">
<b>mlr --from infile.dat put 'tee > "./taps/data-".$a."-".$b, $*'</b>
</pre>

Compressed I/O:

<pre class="pre-highlight-non-pair">
<b>mlr --from infile.dat put 'tee | "gzip > ./taps/data-".$a."-".$b.".gz", $*'</b>
</pre>

Interoperate with other data-processing tools using standard pipes:

<pre class="pre-highlight-non-pair">
<b>mlr --from infile.dat put -q '@v=$*; dump | "jq .[]"'</b>
</pre>

Tap/trace:

<pre class="pre-highlight-non-pair">
<b>mlr --from infile.dat put  '(NR % 1000 == 0) { print > stderr, "Checkpoint ".NR}'</b>
</pre>

## Program timing

This admittedly artificial example demonstrates using Miller time and stats functions to introspectively acquire some information about Miller's own runtime. The `delta` function computes the difference between successive timestamps.

<pre class="pre-non-highlight-non-pair">
$ ruby -e '10000.times{|i|puts "i=#{i+1}"}' &gt; lines.txt

$ head -n 5 lines.txt
i=1
i=2
i=3
i=4
i=5

mlr --ofmt '%.9le' --opprint put '$t=systime()' then step -a delta -f t lines.txt | head -n 7
i     t                 t_delta
1     1430603027.018016 1.430603027e+09
2     1430603027.018043 2.694129944e-05
3     1430603027.018048 5.006790161e-06
4     1430603027.018052 4.053115845e-06
5     1430603027.018055 2.861022949e-06
6     1430603027.018058 3.099441528e-06

mlr --ofmt '%.9le' --oxtab \
  put '$t=systime()' then \
  step -a delta -f t then \
  filter '$i&gt;1' then \
  stats1 -a min,mean,max -f t_delta \
  lines.txt
t_delta_min  2.861022949e-06
t_delta_mean 4.077508505e-06
t_delta_max  5.388259888e-05
</pre>

## Showing differences between successive queries

Suppose you have a database query which you run at one point in time, producing the output on the left, then again later producing the output on the right:

<pre class="pre-highlight-in-pair">
<b>cat data/previous_counters.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,count
red,3472
blue,6838
orange,694
purple,12
</pre>

<pre class="pre-highlight-in-pair">
<b>cat data/current_counters.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,count
red,3467
orange,670
yellow,27
blue,6944
</pre>

And, suppose you want to compute the differences in the counters between adjacent keys. Since the color names aren't all in the same order, nor are they all present on both sides, we can't just paste the two files side-by-side and do some column-four-minus-column-two arithmetic.

First, rename counter columns to make them distinct:

<pre class="pre-highlight-in-pair">
<b>mlr --csv rename count,previous_count data/previous_counters.csv > data/prevtemp.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

<pre class="pre-highlight-in-pair">
<b>cat data/prevtemp.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,previous_count
red,3472
blue,6838
orange,694
purple,12
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv rename count,current_count data/current_counters.csv > data/currtemp.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

<pre class="pre-highlight-in-pair">
<b>cat data/currtemp.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,current_count
red,3467
orange,670
yellow,27
blue,6944
</pre>

Then, join on the key field(s), and use unsparsify to zero-fill counters absent on one side but present on the other. Use `--ul` and `--ur` to emit unpaired records (namely, purple on the left and yellow on the right):

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint \</b>
<b>  join -j color --ul --ur -f data/prevtemp.csv \</b>
<b>  then unsparsify --fill-with 0 \</b>
<b>  then put '$count_delta = $current_count - $previous_count' \</b>
<b>  data/currtemp.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  previous_count current_count count_delta
red    3472           3467          -5
orange 694            670           -24
yellow 0              27            (error)
blue   6838           6944          106
purple 12             0             (error)
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

See also the [record-heterogeneity page](record-heterogeneity.md).

## Memoization with out-of-stream variables

The recursive function for the Fibonacci sequence is famous for its computational complexity.  Namely, using f(0)=1, f(1)=1, f(n)=f(n-1)+f(n-2) for n>=2, the evaluation tree branches left as well as right at each non-trivial level, resulting in millions or more paths to the root 0/1 nodes for larger n. This program

<pre class="pre-non-highlight-non-pair">
mlr --ofmt '%.9lf' --opprint seqgen --start 1 --stop 28 then put '
  func f(n) {
      @fcount += 1;              # count number of calls to the function
      if (n &lt; 2) {
          return 1
      } else {
          return f(n-1) + f(n-2) # recurse
      }
  }

  @fcount = 0;
  $o = f($i);
  $fcount = @fcount;

' then put '$seconds=systime()' then step -a delta -f seconds then cut -x -f seconds
</pre>

produces output like this:

<pre class="pre-non-highlight-non-pair">
i  o      fcount  seconds_delta
1  1      1       0
2  2      3       0.000039101
3  3      5       0.000015974
4  5      9       0.000019073
5  8      15      0.000026941
6  13     25      0.000036955
7  21     41      0.000056028
8  34     67      0.000086069
9  55     109     0.000134945
10 89     177     0.000217915
11 144    287     0.000355959
12 233    465     0.000506163
13 377    753     0.000811815
14 610    1219    0.001297235
15 987    1973    0.001960993
16 1597   3193    0.003417969
17 2584   5167    0.006215811
18 4181   8361    0.008294106
19 6765   13529   0.012095928
20 10946  21891   0.019592047
21 17711  35421   0.031193972
22 28657  57313   0.057254076
23 46368  92735   0.080307961
24 75025  150049  0.129482031
25 121393 242785  0.213325977
26 196418 392835  0.334423065
27 317811 635621  0.605969906
28 514229 1028457 0.971235037
</pre>

Note that the time it takes to evaluate the function is blowing up exponentially as the input argument increases. Using `@`-variables, which persist across records, we can cache and reuse the results of previous computations:

<pre class="pre-non-highlight-non-pair">
mlr --ofmt '%.9lf' --opprint seqgen --start 1 --stop 28 then put '
  func f(n) {
    @fcount += 1;                 # count number of calls to the function
    if (is_present(@fcache[n])) { # cache hit
      return @fcache[n]
    } else {                      # cache miss
      num rv = 1;
      if (n &gt;= 2) {
        rv = f(n-1) + f(n-2)      # recurse
      }
      @fcache[n] = rv;
      return rv
    }
  }
  @fcount = 0;
  $o = f($i);
  $fcount = @fcount;
' then put '$seconds=systime()' then step -a delta -f seconds then cut -x -f seconds
</pre>

with output like this:

<pre class="pre-non-highlight-non-pair">
i  o      fcount seconds_delta
1  1      1      0
2  2      3      0.000053883
3  3      3      0.000035048
4  5      3      0.000045061
5  8      3      0.000014067
6  13     3      0.000028849
7  21     3      0.000028133
8  34     3      0.000027895
9  55     3      0.000014067
10 89     3      0.000015020
11 144    3      0.000012875
12 233    3      0.000033140
13 377    3      0.000014067
14 610    3      0.000012875
15 987    3      0.000029087
16 1597   3      0.000013828
17 2584   3      0.000013113
18 4181   3      0.000012875
19 6765   3      0.000013113
20 10946  3      0.000012875
21 17711  3      0.000013113
22 28657  3      0.000013113
23 46368  3      0.000015974
24 75025  3      0.000012875
25 121393 3      0.000013113
26 196418 3      0.000012875
27 317811 3      0.000013113
28 514229 3      0.000012875
</pre>
