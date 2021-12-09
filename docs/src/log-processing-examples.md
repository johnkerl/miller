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
# Log-processing examples

Another of my favorite use-cases for Miller is doing ad-hoc processing of log-file data.  Here's where DKVP format really shines: one, since the field names and field values are present on every line, every line stands on its own. That means you can `grep` or what have you. Also it means not every line needs to have the same list of field names ("schema").

## Generating and aggregating log-file output

Again, all the examples in the CSV section apply here -- just change the input-format flags. But there's more you can do when not all the records have the same shape.

Writing a program -- in any language whatsoever -- you can have it print out log lines as it goes along, with items for various events jumbled together. After the program has finished running you can sort it all out, filter it, analyze it, and learn from it.

Suppose your program has printed something like this [log.txt](./log.txt):

<pre class="pre-highlight-in-pair">
<b>cat log.txt</b>
</pre>
<pre class="pre-non-highlight-in-pair">
op=enter,time=1472819681
op=cache,type=A9,hit=0
op=cache,type=A4,hit=1
time=1472819690,batch_size=100,num_filtered=237
op=cache,type=A1,hit=1
op=cache,type=A9,hit=0
op=cache,type=A1,hit=1
op=cache,type=A9,hit=0
op=cache,type=A9,hit=0
op=cache,type=A1,hit=1
time=1472819705,batch_size=100,num_filtered=348
op=cache,type=A4,hit=1
op=cache,type=A9,hit=0
op=cache,type=A9,hit=0
op=cache,type=A9,hit=0
op=cache,type=A9,hit=0
op=cache,type=A4,hit=1
time=1472819713,batch_size=100,num_filtered=493
op=cache,type=A9,hit=1
op=cache,type=A1,hit=1
op=cache,type=A9,hit=0
op=cache,type=A9,hit=0
op=cache,type=A9,hit=1
time=1472819720,batch_size=100,num_filtered=554
op=cache,type=A1,hit=0
op=cache,type=A4,hit=1
op=cache,type=A9,hit=0
op=cache,type=A9,hit=0
op=cache,type=A9,hit=0
op=cache,type=A4,hit=0
op=cache,type=A4,hit=0
op=cache,type=A9,hit=0
time=1472819736,batch_size=100,num_filtered=612
op=cache,type=A1,hit=1
op=cache,type=A9,hit=0
op=cache,type=A9,hit=0
op=cache,type=A9,hit=0
op=cache,type=A9,hit=0
op=cache,type=A4,hit=1
op=cache,type=A1,hit=1
op=cache,type=A9,hit=0
op=cache,type=A9,hit=0
time=1472819742,batch_size=100,num_filtered=728
</pre>

Each print statement simply contains local information: the current timestamp, whether a particular cache was hit or not, etc. Then using either the system `grep` command, or Miller's [having-fields verb](reference-verbs.md#having-fields), or the [is_present DSL function](reference-dsl-builtin-functions.md#is_present), we can pick out the parts we want and analyze them:

<pre class="pre-highlight-in-pair">
<b>grep op=cache log.txt \</b>
<b>  | mlr --idkvp --opprint stats1 -a mean -f hit -g type then sort -f type</b>
</pre>
<pre class="pre-non-highlight-in-pair">
type hit_mean
A1   0.8571428571428571
A4   0.7142857142857143
A9   0.09090909090909091
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --from log.txt --opprint \</b>
<b>  filter 'is_present($batch_size)' \</b>
<b>  then step -a delta -f time,num_filtered \</b>
<b>  then sec2gmt time</b>
</pre>
<pre class="pre-non-highlight-in-pair">
time                 batch_size num_filtered time_delta num_filtered_delta
2016-09-02T12:34:50Z 100        237          0          0
2016-09-02T12:35:05Z 100        348          15         111
2016-09-02T12:35:13Z 100        493          8          145
2016-09-02T12:35:20Z 100        554          7          61
2016-09-02T12:35:36Z 100        612          16         58
2016-09-02T12:35:42Z 100        728          6          116
</pre>

Alternatively, we can simply group the similar data for a better look:

<pre class="pre-highlight-in-pair">
<b>mlr --opprint group-like log.txt</b>
</pre>
<pre class="pre-non-highlight-in-pair">
op    time
enter 1472819681

op    type hit
cache A9   0
cache A4   1
cache A1   1
cache A9   0
cache A1   1
cache A9   0
cache A9   0
cache A1   1
cache A4   1
cache A9   0
cache A9   0
cache A9   0
cache A9   0
cache A4   1
cache A9   1
cache A1   1
cache A9   0
cache A9   0
cache A9   1
cache A1   0
cache A4   1
cache A9   0
cache A9   0
cache A9   0
cache A4   0
cache A4   0
cache A9   0
cache A1   1
cache A9   0
cache A9   0
cache A9   0
cache A9   0
cache A4   1
cache A1   1
cache A9   0
cache A9   0

time       batch_size num_filtered
1472819690 100        237
1472819705 100        348
1472819713 100        493
1472819720 100        554
1472819736 100        612
1472819742 100        728
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --opprint group-like then sec2gmt time log.txt</b>
</pre>
<pre class="pre-non-highlight-in-pair">
op    time
enter 2016-09-02T12:34:41Z

op    type hit
cache A9   0
cache A4   1
cache A1   1
cache A9   0
cache A1   1
cache A9   0
cache A9   0
cache A1   1
cache A4   1
cache A9   0
cache A9   0
cache A9   0
cache A9   0
cache A4   1
cache A9   1
cache A1   1
cache A9   0
cache A9   0
cache A9   1
cache A1   0
cache A4   1
cache A9   0
cache A9   0
cache A9   0
cache A4   0
cache A4   0
cache A9   0
cache A1   1
cache A9   0
cache A9   0
cache A9   0
cache A9   0
cache A4   1
cache A1   1
cache A9   0
cache A9   0

time                 batch_size num_filtered
2016-09-02T12:34:50Z 100        237
2016-09-02T12:35:05Z 100        348
2016-09-02T12:35:13Z 100        493
2016-09-02T12:35:20Z 100        554
2016-09-02T12:35:36Z 100        612
2016-09-02T12:35:42Z 100        728
</pre>

## Parsing log-file output

This, of course, depends highly on what's in your log files. But, as an example, suppose you have log-file lines such as

<pre class="pre-non-highlight-non-pair">
2015-10-08 08:29:09,445 INFO com.company.path.to.ClassName @ [sometext] various/sorts/of data {& punctuation} hits=1 status=0 time=2.378
</pre>

I prefer to pre-filter with `grep` and/or `sed` to extract the structured text, then hand that to Miller. Example:

<pre class="pre-highlight-in-pair">
<b>grep 'various sorts' *.log \</b>
<b>  | sed 's/.*} //' \</b>
<b>  | mlr --fs space --repifs --oxtab stats1 -a min,p10,p50,p90,max -f time -g status</b>
</pre>
<pre class="pre-non-highlight-in-pair">
... output here ...
</pre>
