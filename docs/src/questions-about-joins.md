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
3,0000ff,
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

See also the [record-heterogeneity page](record-heterogeneity.md).

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

## Updating a database file with new values from another

Suppose you have a durable "database" file, and you periodically receive a download
containing updates for some of its records, on some of its columns. You'd like to
update matching records in place, and pass records without updates through unchanged --
without introducing any new records or columns. For example:

<pre class="pre-non-highlight-non-pair">
id,name,color,size,shape
1,alice,red,10,circle
2,bob,green,20,square
3,carol,blue,30,triangle
4,dave,yellow,40,hexagon
</pre>

<pre class="pre-non-highlight-non-pair">
id,color,size
1,crimson,11
3,navy,
</pre>

The tool for this is `join`, and the key fact is its collision rule: **when a
non-join field is present on both sides of a paired record, the value from the
right file overwrites the value from the left file.** So put the database file
on the left (with `-f`) and the download on the right -- then the downloaded
values win. Adding `--ul` also emits unpaired left records, i.e., database
records for which no update arrived. A pleasant side effect of this ordering is
that the output columns are exactly the database file's columns, in the
database file's order:

<pre class="pre-highlight-in-pair">
<b>mlr --csv join --ul -j id -f data/update-db.csv data/update-download.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
id,name,color,size,shape
1,alice,crimson,11,circle
3,carol,navy,,triangle
2,bob,green,20,square
4,dave,yellow,40,hexagon
</pre>

Records 1 and 3 have been updated from the download, while records 2 and 4 are
passed through unchanged. Two things to note:

* Paired records are emitted first (in download order, since the right file is
the one being streamed), then the unpaired database records. If you want the
original ordering back, append a `sort`:

<pre class="pre-highlight-in-pair">
<b>mlr --csv join --ul -j id -f data/update-db.csv \</b>
<b>  then sort -t id \</b>
<b>  data/update-download.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
id,name,color,size,shape
1,alice,crimson,11,circle
2,bob,green,20,square
3,carol,navy,,triangle
4,dave,yellow,40,hexagon
</pre>

* The download's _empty_ `size` cell for record 3 overwrote the database's
`30`. If empty cells in the download should mean "leave the database value
alone", give the download's fields a prefix using `--rp` -- so they no longer
collide -- and then copy over only the non-empty ones:

<pre class="pre-highlight-in-pair">
<b>mlr --csv join --ul --rp upd_ -j id -f data/update-db.csv \</b>
<b>  then put '</b>
<b>    for (k, v in $*) {</b>
<b>      if (k =~ "^upd_") {</b>
<b>        if (!is_empty(v)) {</b>
<b>          $[sub(k, "^upd_", "")] = v;</b>
<b>        }</b>
<b>        unset $[k];</b>
<b>      }</b>
<b>    }</b>
<b>  ' \</b>
<b>  then sort -t id \</b>
<b>  data/update-download.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
id,name,color,size,shape
1,alice,crimson,11,circle
2,bob,green,20,square
3,carol,navy,30,triangle
4,dave,yellow,40,hexagon
</pre>

Here record 3's `size` keeps its database value `30`, while all the non-empty
downloaded values are applied as before.

See also [issue 826](https://github.com/johnkerl/miller/issues/826).
