<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-main-flag-list/index.html">Flag list</a>
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
# Questions about joins

## Why am I not seeing all possible joins occur?

**This section describes behavior before Miller 5.1.0. As of 5.1.0, -u is the default.**

For example, the right file here has nine records, and the left file should add in the `hostname` column -- so the join output should also have 9 records:

<pre class="pre-highlight-in-pair">
<b>mlr --icsvlite --opprint cat data/join-u-left.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
hostname              ipaddr
nadir.east.our.org    10.3.1.18
zenith.west.our.org   10.3.1.27
apoapsis.east.our.org 10.4.5.94
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsvlite --opprint cat data/join-u-right.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
ipaddr    timestamp  bytes
10.3.1.27 1448762579 4568
10.3.1.18 1448762578 8729
10.4.5.94 1448762579 17445
10.3.1.27 1448762589 12
10.3.1.18 1448762588 44558
10.4.5.94 1448762589 8899
10.3.1.27 1448762599 0
10.3.1.18 1448762598 73425
10.4.5.94 1448762599 12200
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsvlite --opprint join -s -j ipaddr -f data/join-u-left.csv data/join-u-right.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
ipaddr    hostname              timestamp  bytes
10.3.1.27 zenith.west.our.org   1448762579 4568
10.4.5.94 apoapsis.east.our.org 1448762579 17445
10.4.5.94 apoapsis.east.our.org 1448762589 8899
10.4.5.94 apoapsis.east.our.org 1448762599 12200
</pre>

The issue is that Miller's `join`, by default (before 5.1.0), took input sorted (lexically ascending) by the sort keys on both the left and right files.  This design decision was made intentionally to parallel the Unix/Linux system `join` command, which has the same semantics. The benefit of this default is that the joiner program can stream through the left and right files, needing to load neither entirely into memory. The drawback, of course, is that is requires sorted input.

The solution (besides pre-sorting the input files on the join keys) is to simply use **mlr join -u** (which is now the default). This loads the left file entirely into memory (while the right file is still streamed one line at a time) and does all possible joins without requiring sorted input:

<pre class="pre-highlight-in-pair">
<b>mlr --icsvlite --opprint join -u -j ipaddr -f data/join-u-left.csv data/join-u-right.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
ipaddr    hostname              timestamp  bytes
10.3.1.27 zenith.west.our.org   1448762579 4568
10.3.1.18 nadir.east.our.org    1448762578 8729
10.4.5.94 apoapsis.east.our.org 1448762579 17445
10.3.1.27 zenith.west.our.org   1448762589 12
10.3.1.18 nadir.east.our.org    1448762588 44558
10.4.5.94 apoapsis.east.our.org 1448762589 8899
10.3.1.27 zenith.west.our.org   1448762599 0
10.3.1.18 nadir.east.our.org    1448762598 73425
10.4.5.94 apoapsis.east.our.org 1448762599 12200
</pre>

General advice is to make sure the left-file is relatively small, e.g. containing name-to-number mappings, while saving large amounts of data for the right file.

## How to rectangularize after joins with unpaired?

Suppose you have the following two data files:

<pre class="pre-non-highlight-non-pair">
id,code
3,0000ff
2,00ff00
4,ff0000
</pre>

<pre class="pre-non-highlight-non-pair">
id,color
4,red
2,green
</pre>

Joining on color the results are as expected:

<pre class="pre-highlight-in-pair">
<b>mlr --csv join -j id -f data/color-codes.csv data/color-names.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
id,code,color
4,ff0000,red
2,00ff00,green
</pre>

However, if we ask for left-unpaireds, since there's no `color` column, we get a row not having the same column names as the other:

<pre class="pre-highlight-in-pair">
<b>mlr --csv join --ul -j id -f data/color-codes.csv data/color-names.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
id,code,color
4,ff0000,red
2,00ff00,green

id,code
3,0000ff
</pre>

To fix this, we can use **unsparsify**:

<pre class="pre-highlight-in-pair">
<b>mlr --csv join --ul -j id -f data/color-codes.csv \</b>
<b>  then unsparsify --fill-with "" \</b>
<b>  data/color-names.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
id,code,color
4,ff0000,red
2,00ff00,green
3,0000ff,
</pre>

Thanks to @aborruso for the tip!

## Doing multiple joins

Suppose we have the following data:

<pre class="pre-highlight-in-pair">
<b>cat multi-join/input.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
id,task
10,chop
20,puree
20,wash
30,fold
10,bake
20,mix
10,knead
30,clean
</pre>

And we want to augment the `id` column with lookups from the following data files:

<pre class="pre-highlight-in-pair">
<b>cat multi-join/name-lookup.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
id,name
30,Alice
10,Bob
20,Carol
</pre>

<pre class="pre-highlight-in-pair">
<b>cat multi-join/status-lookup.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
id,status
30,occupied
10,idle
20,idle
</pre>

We can run the input file through multiple `join` commands in a `then`-chain:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint join -f multi-join/name-lookup.csv -j id \</b>
<b>  then join -f multi-join/status-lookup.csv -j id \</b>
<b>  multi-join/input.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
id status   name  task
10 idle     Bob   chop
20 idle     Carol puree
20 idle     Carol wash
30 occupied Alice fold
10 idle     Bob   bake
20 idle     Carol mix
10 idle     Bob   knead
30 occupied Alice clean
</pre>
