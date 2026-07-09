<!--  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. -->
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

## Doing SQL-style left, right, inner, and full-outer joins

Miller's `join` verb is defined in terms of _paired_ and _unpaired_ records, rather than SQL-database terminology -- but you can get SQL-style joins using the `--ul` and `--ur` flags (which emit unpaired left-file and right-file records, respectively), along with [`unsparsify`](reference-verbs.md#unsparsify) to fill in empty cells for non-matches.

Suppose you have the following two data files, where we want to join on the left file's `a` field matching the right file's `e` field:

<pre class="pre-non-highlight-non-pair">
a,b,c
a,t,1
b,u,2
c,v,3
</pre>

<pre class="pre-non-highlight-non-pair">
e,f,g
a,t,3
b,u,2
d,w,1
</pre>

In all the following examples, the `-f` file (`data/join-x.csv`) is the left file, and the file in the main input stream (`data/join-y.csv`) is the right file. The flags `-j a -r e` say that the left file's `a` field is matched against the right file's `e` field, with the output join column named `a`.

**Inner join** -- only matching records -- is what Miller's `join` does by default, since only paired records are emitted:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ocsv join -j a -r e -f data/join-x.csv data/join-y.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c,f,g
a,t,1,t,3
b,u,2,u,2
</pre>

**Left join** keeps all records from the left file, with empty cells where the right file has no match. Use `--ul` to also emit unpaired left-file records, then `unsparsify` to square up the output:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ocsv join --ul -j a -r e -f data/join-x.csv \</b>
<b>  then unsparsify --fill-with "" \</b>
<b>  data/join-y.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c,f,g
a,t,1,t,3
b,u,2,u,2
c,v,3,,
</pre>

**Right join** keeps all records from the right file. Use `--ur` to also emit unpaired right-file records:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ocsv join --ur -j a -r e -f data/join-x.csv \</b>
<b>  then unsparsify --fill-with "" \</b>
<b>  data/join-y.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c,f,g
a,t,1,t,3
b,u,2,u,2
d,,,w,1
</pre>

**Full outer join** keeps all records from both files. Use both `--ul` and `--ur`:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ocsv join --ul --ur -j a -r e -f data/join-x.csv \</b>
<b>  then unsparsify --fill-with "" \</b>
<b>  data/join-y.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c,f,g
a,t,1,t,3
b,u,2,u,2
d,,,w,1
c,v,3,,
</pre>

Note that unpaired records are emitted after all paired records, so the output ordering may differ from what a SQL database would produce; you can pipe the output through [`sort`](reference-verbs.md#sort) if you need a particular ordering.

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

## Merging several files on a common key when column names collide

The previous example worked painlessly because the non-key column names -- `name` and `status` -- were different in each lookup file. Now suppose you have several files, each containing measurements of a _different_ quantity, but all with the _same_ column names -- say, one file each for temperature, humidity, and pressure, keyed by timestamp:

<pre class="pre-highlight-in-pair">
<b>cat data/sensor-temp.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
unixTime,minValue,averageValue,maxValue
1740000000,37.2,37.4,37.6
1740000060,37.5,37.6,37.7
1740000120,37.8,37.9,38.0
</pre>

<pre class="pre-highlight-in-pair">
<b>cat data/sensor-humidity.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
unixTime,minValue,averageValue,maxValue
1740000000,50.1,50.3,50.5
1740000120,52.3,52.5,52.7
</pre>

<pre class="pre-highlight-in-pair">
<b>cat data/sensor-pressure.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
unixTime,minValue,averageValue,maxValue
1740000000,1012.2,1012.3,1012.4
1740000060,1011.8,1011.9,1012.0
1740000120,1011.5,1011.6,1011.7
</pre>

(Note that the humidity file is missing a row for the middle timestamp.)

If we merge these with a `then`-chain of `join` commands, as in the previous section, columns are lost: since every file's non-key columns have the same names, each join step overwrites the values from the step before, and only one file's values survive:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint \</b>
<b>  join -j unixTime -f data/sensor-temp.csv \</b>
<b>  then join -j unixTime -f data/sensor-humidity.csv \</b>
<b>  data/sensor-pressure.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
unixTime   minValue averageValue maxValue
1740000000 1012.2   1012.3       1012.4
1740000120 1011.5   1011.6       1011.7
</pre>

The fix is `join`'s `--lp` (left-prefix) option, which renames the non-key columns coming from each `-f` file so that nothing collides:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint \</b>
<b>  join --lp temp: -j unixTime -f data/sensor-temp.csv \</b>
<b>  then join --lp humidity: -j unixTime -f data/sensor-humidity.csv \</b>
<b>  data/sensor-pressure.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
unixTime   humidity:minValue humidity:averageValue humidity:maxValue temp:minValue temp:averageValue temp:maxValue minValue averageValue maxValue
1740000000 50.1              50.3                  50.5              37.2          37.4              37.6          1012.2   1012.3       1012.4
1740000120 52.3              52.5                  52.7              37.8          37.9              38.0          1011.5   1011.6       1011.7
</pre>

The columns from the file at the end of the command line -- here, the pressure file -- keep their unprefixed names.

If you want just one value column per file, you can also use `join`'s `--lk` option to keep only that column from each `-f` file, then use [cut](reference-verbs.md#cut) and [label](reference-verbs.md#label) to arrange and rename the output columns:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint \</b>
<b>  join --lp temp: --lk averageValue -j unixTime -f data/sensor-temp.csv \</b>
<b>  then join --lp humidity: --lk averageValue -j unixTime -f data/sensor-humidity.csv \</b>
<b>  then cut -o -f unixTime,temp:averageValue,humidity:averageValue,averageValue \</b>
<b>  then label unixTime,temperature,humidity,pressure \</b>
<b>  data/sensor-pressure.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
unixTime   temperature humidity pressure
1740000000 37.4        50.3     1012.3
1740000120 37.9        52.5     1011.6
</pre>

Note that `join` gives inner-join semantics by default, so the timestamp missing from the humidity file has been dropped from the output. (This is also why the `paste` command is not a substitute for `join` here: `paste` matches rows by position, so a row missing from one file shifts that file's remaining values onto the wrong rows.) If you'd rather keep those rows, with empty cells where a file has no data, add `--ul --ur` to each join step, and [unsparsify](reference-verbs.md#unsparsify) afterward:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint \</b>
<b>  join --ul --ur --lp temp: --lk averageValue -j unixTime -f data/sensor-temp.csv \</b>
<b>  then join --ul --ur --lp humidity: --lk averageValue -j unixTime -f data/sensor-humidity.csv \</b>
<b>  then unsparsify \</b>
<b>  then cut -o -f unixTime,temp:averageValue,humidity:averageValue,averageValue \</b>
<b>  then label unixTime,temperature,humidity,pressure \</b>
<b>  then sort -t unixTime \</b>
<b>  data/sensor-pressure.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
unixTime   temperature humidity pressure
1740000000 37.4        50.3     1012.3
1740000060 37.6        -        1011.9
1740000120 37.9        52.5     1011.6
</pre>

The `unsparsify` must come before the `cut`-and-`label` step, since `label` renames columns positionally. The `sort` is there because unpaired records may be emitted out of order relative to paired ones.

## How to preprocess the left file of a join?

The left file (the `-f` argument to `join`) is opened by the `join` verb itself, so it doesn't pass through the main record stream: a `then`-chain can preprocess the right file(s), but not the left one.

For example, suppose both of these files have multi-valued `id` fields which need [nest --explode](reference-verbs.md#nest) before joining:

<pre class="pre-highlight-in-pair">
<b>cat data/join-nest-left.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
id,color
1;2,blue
3,green
</pre>

<pre class="pre-highlight-in-pair">
<b>cat data/join-nest-right.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
id,shape
1,circle
2;3,square
</pre>

Using `nest` in a `then`-chain handles the right file, but the left file still has the unexploded `id` value `1;2`, so only `id=3` pairs up:

<pre class="pre-highlight-in-pair">
<b>mlr --csv nest --evar ';' -f id \</b>
<b>  then join -j id -f data/join-nest-left.csv \</b>
<b>  data/join-nest-right.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
id,color,shape
3,green,square
</pre>

One way to preprocess the left file, without creating an intermediate file, is the main-level `--prepipe` flag. The left file inherits the main input options -- including `--prepipe` -- so the specified command is applied to the left file as well as to the right file(s):

<pre class="pre-highlight-in-pair">
<b>mlr --csv --prepipe 'mlr --csv nest --evar ";" -f id' \</b>
<b>  join -j id -f data/join-nest-left.csv \</b>
<b>  data/join-nest-right.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
id,color,shape
1,blue,circle
2,blue,square
3,green,square
</pre>

Note that `--prepipe` applies the same command to _every_ input file -- which is just what's wanted here, since both files need the same `nest`.

Another way, if your shell supports it -- bash, zsh, and ksh do, although plain POSIX `sh` and Windows `cmd` do not -- is [process substitution](https://en.wikipedia.org/wiki/Process_substitution), which lets you preprocess the left file with any command at all, independently of the right file(s):

<pre class="pre-highlight-in-pair">
<b>mlr --csv nest --evar ';' -f id \</b>
<b>  then join -j id -f <(mlr --csv nest --evar ';' -f id data/join-nest-left.csv) \</b>
<b>  data/join-nest-right.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
id,color,shape
1,blue,circle
2,blue,square
3,green,square
</pre>

Note that while `mlr join --help` lists verb-level `--prepipe` and `--prepipex` flags, as of Miller 6 these do not take effect for the left file -- please use one of the recipes above instead.

Thanks to @sonicdoe for the process-substitution tip!

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
