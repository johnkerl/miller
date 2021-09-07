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
# DSL output statements

You can **output** variable-values or expressions in **five ways**:

* **Assign** them to stream-record fields. For example, `$cumulative_sum = @sum`. For another example, `$nr = NR` adds a field named `nr` to each output record, containing the value of the built-in variable `NR` as of when that record was ingested.

* Use the **print** or **eprint** keywords which immediately print an expression *directly to standard output or standard error*, respectively. Note that `dump`, `edump`, `print`, and `eprint` don't output records which participate in `then`-chaining; rather, they're just immediate prints to stdout/stderr. The `printn` and `eprintn` keywords are the same except that they don't print final newlines. Additionally, you can print to a specified file instead of stdout/stderr.

* Use the **dump** or **edump** keywords, which *immediately print all out-of-stream variables as a JSON data structure to the standard output or standard error* (respectively).

* Use **tee** which formats the current stream record (not just an arbitrary string as with **print**) to a specific file.

* Use **emit**/**emitp**/**emitf** to send out-of-stream variables' current values to the output record stream, e.g.  `@sum += $x; emit @sum` which produces an extra record such as `sum=3.1648382`. These records, just like records from input file(s), participate in downstream [then-chaining](reference-main-then-chaining.md) to other verbs.

For the first two options you are populating the output-records stream which feeds into the next verb in a `then`-chain (if any), or which otherwise is formatted for output using `--o...` flags.

For the last three options you are sending output directly to standard output, standard error, or a file.

## Print statements

The `print` statement is perhaps self-explanatory, but with a few light caveats:

* There are four variants: `print` goes to stdout with final newline, `printn` goes to stdout without final newline (you can include one using "\n" in your output string), `eprint` goes to stderr with final newline, and `eprintn` goes to stderr without final newline.

* Output goes directly to stdout/stderr, respectively: data produced this way do not go downstream to the next verb in a `then`-chain. (Use `emit` for that.)

* Print statements are for strings (`print "hello"`), or things which can be made into strings: numbers (`print 3`, `print $a + $b`), or concatenations thereof (`print "a + b = " . ($a + $b)`). Maps (in `$*`, map-valued out-of-stream or local variables, and map literals) as well as arrays are printed as JSON.

* You can redirect print output to a file:

<pre class="pre-highlight-non-pair">
<b>mlr --from myfile.dat put 'print > "tap.txt", $x'</b>
</pre>

* You can redirect print output to multiple files, split by values present in various records:

<pre class="pre-highlight-non-pair">
<b>mlr --from myfile.dat put 'print > $a.".txt", $x'</b>
</pre>

See also [Redirected-output statements](reference-dsl-output-statements.md#redirected-output-statements) for examples.

## Dump statements

The `dump` statement is for printing expressions, including maps, directly to stdout/stderr, respectively:

* There are two variants: `dump` prints to stdout; `edump` prints to stderr.

* Output goes directly to stdout/stderr, respectively: data produced this way do not go downstream to the next verb in a `then`-chain. (Use `emit` for that.)

* You can use `dump` to output single strings, numbers, or expressions including map-valued data. Map-valued data are printed as JSON.

* If you use `dump` (or `edump`) with no arguments, you get a JSON structure representing the current values of all out-of-stream variables.

* As with `print`, you can redirect output to files.

* See also [Redirected-output statements](reference-dsl-output-statements.md#redirected-output-statements) for examples.

## Tee statements

Records produced by a `mlr put` go downstream to the next verb in your `then`-chain, if any, or otherwise to standard output.  If you want to additionally copy out records to files, you can do that using `tee`.

The syntax is, by example:

<pre class="pre-highlight-non-pair">
<b>mlr --from myfile.dat put 'tee > "tap.dat", $*' then sort -n index</b>
</pre>

First is `tee >`, then the filename expression (which can be an expression such as `"tap.".$a.".dat"`), then a comma, then `$*`. (Nothing else but `$*` is teeable.)

You can also write to a variable file name -- for example, you can split a
single file into multiple ones on field names:

<pre class="pre-highlight-in-pair">
<b>mlr --csv cat example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
yellow,triangle,true,1,11,43.6498,9.8870
red,square,true,2,15,79.2778,0.0130
red,circle,true,3,16,13.8103,2.9010
red,square,false,4,48,77.5542,7.4670
purple,triangle,false,5,51,81.2290,8.5910
red,square,false,6,64,77.1991,9.5310
purple,triangle,false,7,65,80.1405,5.8240
yellow,circle,true,8,73,63.9785,4.2370
yellow,circle,true,9,87,63.5058,8.3350
purple,square,false,10,91,72.3735,8.2430
</pre>

<pre class="pre-highlight-non-pair">
<b>mlr --csv --from example.csv put -q 'tee > $shape.".csv", $*'</b>
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv cat circle.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
red,circle,true,3,16,13.8103,2.9010
yellow,circle,true,8,73,63.9785,4.2370
yellow,circle,true,9,87,63.5058,8.3350
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv cat square.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
red,square,true,2,15,79.2778,0.0130
red,square,false,4,48,77.5542,7.4670
red,square,false,6,64,77.1991,9.5310
purple,square,false,10,91,72.3735,8.2430
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv cat triangle.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
yellow,triangle,true,1,11,43.6498,9.8870
purple,triangle,false,5,51,81.2290,8.5910
purple,triangle,false,7,65,80.1405,5.8240
</pre>

See also [Redirected-output statements](reference-dsl-output-statements.md#redirected-output-statements) for examples.

## Redirected-output statements

The **print**, **dump** **tee**, **emitf**, **emit**, and **emitp** keywords all allow you to redirect output to one or more files or pipe-to commands. The filenames/commands are strings which can be constructed using record-dependent values, so you can do things like splitting a table into multiple files, one for each account ID, and so on.

Details:

* The `print` and `dump` keywords produce output immediately to standard output, or to specified file(s) or pipe-to command if present.

<pre class="pre-highlight-in-pair">
<b>mlr help keyword print</b>
</pre>
<pre class="pre-non-highlight-in-pair">
print: prints expression immediately to stdout.

  Example: mlr --from f.dat put -q 'print "The sum of x and y is ".($x+$y)'
  Example: mlr --from f.dat put -q 'for (k, v in $*) { print k . " => " . v }'
  Example: mlr --from f.dat put  '(NR % 1000 == 0) { print > stderr, "Checkpoint ".NR}'
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr help keyword dump</b>
</pre>
<pre class="pre-non-highlight-in-pair">
dump: prints all currently defined out-of-stream variables immediately
to stdout as JSON.

With >, >>, or |, the data do not become part of the output record stream but
are instead redirected.

The > and >> are for write and append, as in the shell, but (as with awk) the
file-overwrite for > is on first write, not per record. The | is for piping to
a process which will process the data. There will be one open file for each
distinct file name (for > and >>) or one subordinate process for each distinct
value of the piped-to command (for |). Output-formatting flags are taken from
the main command line.

  Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump }'
  Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump >  "mytap.dat"}'
  Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump >> "mytap.dat"}'
  Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump | "jq .[]"}'
</pre>

* `mlr put` sends the current record (possibly modified by the `put` expression) to the output record stream. Records are then input to the following verb in a `then`-chain (if any), else printed to standard output (unless `put -q`). The **tee** keyword *additionally* writes the output record to specified file(s) or pipe-to command, or immediately to `stdout`/`stderr`.

<pre class="pre-highlight-in-pair">
<b>mlr help keyword tee</b>
</pre>
<pre class="pre-non-highlight-in-pair">
tee: prints the current record to specified file.
This is an immediate print to the specified file (except for pprint format
which of course waits until the end of the input stream to format all output).

The > and >> are for write and append, as in the shell, but (as with awk) the
file-overwrite for > is on first write, not per record. The | is for piping to
a process which will process the data. There will be one open file for each
distinct file name (for > and >>) or one subordinate process for each distinct
value of the piped-to command (for |). Output-formatting flags are taken from
the main command line.

You can use any of the output-format command-line flags, e.g. --ocsv, --ofs,
etc., to control the format of the output. See also mlr -h.

emit with redirect and tee with redirect are identical, except tee can only
output $*.

  Example: mlr --from f.dat put 'tee >  "/tmp/data-".$a, $*'
  Example: mlr --from f.dat put 'tee >> "/tmp/data-".$a.$b, $*'
  Example: mlr --from f.dat put 'tee >  stderr, $*'
  Example: mlr --from f.dat put -q 'tee | "tr \[a-z\\] \[A-Z\\]", $*'
  Example: mlr --from f.dat put -q 'tee | "tr \[a-z\\] \[A-Z\\] > /tmp/data-".$a, $*'
  Example: mlr --from f.dat put -q 'tee | "gzip > /tmp/data-".$a.".gz", $*'
  Example: mlr --from f.dat put -q --ojson 'tee | "gzip > /tmp/data-".$a.".gz", $*'
</pre>

* `mlr put`'s `emitf`, `emitp`, and `emit` send out-of-stream variables to the output record stream. These are then input to the following verb in a `then`-chain (if any), else printed to standard output. When redirected with `>`, `>>`, or `|`, they *instead* write the out-of-stream variable(s) to specified file(s) or pipe-to command, or immediately to `stdout`/`stderr`.

<pre class="pre-highlight-in-pair">
<b>mlr help keyword emitf</b>
</pre>
<pre class="pre-non-highlight-in-pair">
emitf: inserts non-indexed out-of-stream variable(s) side-by-side into the
output record stream.

With >, >>, or |, the data do not become part of the output record stream but
are instead redirected.

The > and >> are for write and append, as in the shell, but (as with awk) the
file-overwrite for > is on first write, not per record. The | is for piping to
a process which will process the data. There will be one open file for each
distinct file name (for > and >>) or one subordinate process for each distinct
value of the piped-to command (for |). Output-formatting flags are taken from
the main command line.

You can use any of the output-format command-line flags, e.g. --ocsv, --ofs,
etc., to control the format of the output if the output is redirected. See also mlr -h.

  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf @a'
  Example: mlr --from f.dat put --oxtab '@a=$i;@b+=$x;@c+=$y; emitf > "tap-".$i.".dat", @a'
  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf @a, @b, @c'
  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf > "mytap.dat", @a, @b, @c'
  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf >> "mytap.dat", @a, @b, @c'
  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf > stderr, @a, @b, @c'
  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf | "grep somepattern", @a, @b, @c'
  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf | "grep somepattern > mytap.dat", @a, @b, @c'

Please see https://johnkerl.org/miller6://johnkerl.org/miller/doc for more information.
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr help keyword emitp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
emitp: inserts an out-of-stream variable into the output record stream.
Hashmap indices present in the data but not slotted by emitp arguments are
output concatenated with ":".

With >, >>, or |, the data do not become part of the output record stream but
are instead redirected.

The > and >> are for write and append, as in the shell, but (as with awk) the
file-overwrite for > is on first write, not per record. The | is for piping to
a process which will process the data. There will be one open file for each
distinct file name (for > and >>) or one subordinate process for each distinct
value of the piped-to command (for |). Output-formatting flags are taken from
the main command line.

You can use any of the output-format command-line flags, e.g. --ocsv, --ofs,
etc., to control the format of the output if the output is redirected. See also mlr -h.

  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp @sums'
  Example: mlr --from f.dat put --opprint '@sums[$a][$b]+=$x; emitp > "tap-".$a.$b.".dat", @sums'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp @sums, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp >  "mytap.dat", @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp >> "mytap.dat", @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp | "gzip > mytap.dat.gz", @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp > stderr, @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp | "grep somepattern", @*, "index1", "index2"'

Please see https://johnkerl.org/miller6://johnkerl.org/miller/doc for more information.
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr help keyword emit</b>
</pre>
<pre class="pre-non-highlight-in-pair">
emit: inserts an out-of-stream variable into the output record stream. Hashmap
indices present in the data but not slotted by emit arguments are not output.

With >, >>, or |, the data do not become part of the output record stream but
are instead redirected.

The > and >> are for write and append, as in the shell, but (as with awk) the
file-overwrite for > is on first write, not per record. The | is for piping to
a process which will process the data. There will be one open file for each
distinct file name (for > and >>) or one subordinate process for each distinct
value of the piped-to command (for |). Output-formatting flags are taken from
the main command line.

You can use any of the output-format command-line flags, e.g. --ocsv, --ofs,
etc., to control the format of the output if the output is redirected. See also mlr -h.

  Example: mlr --from f.dat put 'emit >  "/tmp/data-".$a, $*'
  Example: mlr --from f.dat put 'emit >  "/tmp/data-".$a, mapexcept($*, "a")'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit @sums'
  Example: mlr --from f.dat put --ojson '@sums[$a][$b]+=$x; emit > "tap-".$a.$b.".dat", @sums'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit @sums, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit >  "mytap.dat", @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit >> "mytap.dat", @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit | "gzip > mytap.dat.gz", @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit > stderr, @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit | "grep somepattern", @*, "index1", "index2"'

Please see https://johnkerl.org/miller6://johnkerl.org/miller/doc for more information.
</pre>

## Emit statements

There are three variants: `emitf`, `emit`, and `emitp`. Keep in mind that
out-of-stream variables are a nested, multi-level [map](reference-main-maps.md) (directly viewable as
JSON using `dump`), while Miller record values are as well during processing --
but records may be flattened down for output to tabular formats. See the page
[Flatten/unflatten: JSON vs. tabular formats](flatten-unflatten.md) for more
information.

You can emit any map-valued expression, including `$*`, map-valued out-of-stream variables, the entire out-of-stream-variable collection `@*`, map-valued local variables, map literals, or map-valued function return values.

Use **emitf** to output several out-of-stream variables side-by-side in the same output record. For `emitf` these mustn't have indexing using `@name[...]`. Example:

<pre class="pre-highlight-in-pair">
<b>mlr put -q '</b>
<b>  @count += 1;</b>
<b>  @x_sum += $x;</b>
<b>  @y_sum += $y;</b>
<b>  end { emitf @count, @x_sum, @y_sum}</b>
<b>' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
count=5,x_sum=2.26476,y_sum=2.585083
</pre>

Use **emit** to output an out-of-stream variable. If it's non-indexed you'll get a simple key-value pair:

<pre class="pre-highlight-in-pair">
<b>cat data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=1,x=0.346791,y=0.726802
a=eks,b=pan,i=2,x=0.758679,y=0.522151
a=wye,b=wye,i=3,x=0.204603,y=0.338318
a=eks,b=wye,i=4,x=0.381399,y=0.134188
a=wye,b=pan,i=5,x=0.573288,y=0.863624
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put -q '@sum += $x; end { dump }' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "sum": 2.26476
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put -q '@sum += $x; end { emit @sum }' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
sum=2.26476
</pre>

If it's indexed then use as many names after `emit` as there are indices:

<pre class="pre-highlight-in-pair">
<b>mlr put -q '@sum[$a] += $x; end { dump }' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "sum": {
    "pan": 0.346791,
    "eks": 1.140078,
    "wye": 0.777891
  }
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put -q '@sum[$a] += $x; end { emit @sum, "a" }' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,sum=0.346791
a=eks,sum=1.140078
a=wye,sum=0.777891
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put -q '@sum[$a][$b] += $x; end { dump }' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "sum": {
    "pan": {
      "pan": 0.346791
    },
    "eks": {
      "pan": 0.758679,
      "wye": 0.381399
    },
    "wye": {
      "wye": 0.204603,
      "pan": 0.573288
    }
  }
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put -q '@sum[$a][$b] += $x; end { emit @sum, "a", "b" }' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,sum=0.346791
a=eks,b=pan,sum=0.758679
a=eks,b=wye,sum=0.381399
a=wye,b=wye,sum=0.204603
a=wye,b=pan,sum=0.573288
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put -q '@sum[$a][$b][$i] += $x; end { dump }' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "sum": {
    "pan": {
      "pan": {
        "1": 0.346791
      }
    },
    "eks": {
      "pan": {
        "2": 0.758679
      },
      "wye": {
        "4": 0.381399
      }
    },
    "wye": {
      "wye": {
        "3": 0.204603
      },
      "pan": {
        "5": 0.573288
      }
    }
  }
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put -q '</b>
<b>  @sum[$a][$b][$i] += $x;</b>
<b>  end { emit @sum, "a", "b", "i" }</b>
<b>' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=1,sum=0.346791
a=eks,b=pan,i=2,sum=0.758679
a=eks,b=wye,i=4,sum=0.381399
a=wye,b=wye,i=3,sum=0.204603
a=wye,b=pan,i=5,sum=0.573288
</pre>

Now for **emitp**: if you have as many names following `emit` as there are levels in the out-of-stream variable's map, then `emit` and `emitp` do the same thing. Where they differ is when you don't specify as many names as there are map levels. In this case, Miller needs to flatten multiple map indices down to output-record keys: `emitp` includes full prefixing (hence the `p` in `emitp`) while `emit` takes the deepest map key as the output-record key:

<pre class="pre-highlight-in-pair">
<b>mlr put -q '@sum[$a][$b] += $x; end { dump }' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "sum": {
    "pan": {
      "pan": 0.346791
    },
    "eks": {
      "pan": 0.758679,
      "wye": 0.381399
    },
    "wye": {
      "wye": 0.204603,
      "pan": 0.573288
    }
  }
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put -q '@sum[$a][$b] += $x; end { emit @sum, "a" }' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,pan=0.346791
a=eks,pan=0.758679,wye=0.381399
a=wye,wye=0.204603,pan=0.573288
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put -q '@sum[$a][$b] += $x; end { emit @sum }' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
pan=0.346791
pan=0.758679,wye=0.381399
wye=0.204603,pan=0.573288
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put -q '@sum[$a][$b] += $x; end { emitp @sum, "a" }' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,sum.pan=0.346791
a=eks,sum.pan=0.758679,sum.wye=0.381399
a=wye,sum.wye=0.204603,sum.pan=0.573288
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put -q '@sum[$a][$b] += $x; end { emitp @sum }' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
sum.pan.pan=0.346791,sum.eks.pan=0.758679,sum.eks.wye=0.381399,sum.wye.wye=0.204603,sum.wye.pan=0.573288
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --oxtab put -q '@sum[$a][$b] += $x; end { emitp @sum }' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
sum.pan.pan 0.346791
sum.eks.pan 0.758679
sum.eks.wye 0.381399
sum.wye.wye 0.204603
sum.wye.pan 0.573288
</pre>

Use **--oflatsep** to specify the character which joins multilevel
keys for `emitp` (it defaults to a colon):

<pre class="pre-highlight-in-pair">
<b>mlr put -q --oflatsep / '@sum[$a][$b] += $x; end { emitp @sum, "a" }' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,sum.pan=0.346791
a=eks,sum.pan=0.758679,sum.wye=0.381399
a=wye,sum.wye=0.204603,sum.pan=0.573288
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put -q --oflatsep / '@sum[$a][$b] += $x; end { emitp @sum }' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
sum.pan.pan=0.346791,sum.eks.pan=0.758679,sum.eks.wye=0.381399,sum.wye.wye=0.204603,sum.wye.pan=0.573288
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --oxtab put -q --oflatsep / '</b>
<b>  @sum[$a][$b] += $x;</b>
<b>  end { emitp @sum }</b>
<b>' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
sum.pan.pan 0.346791
sum.eks.pan 0.758679
sum.eks.wye 0.381399
sum.wye.wye 0.204603
sum.wye.pan 0.573288
</pre>

## Multi-emit statements

You can emit **multiple map-valued expressions side-by-side** by
including their names in parentheses:

<pre class="pre-highlight-in-pair">
<b>mlr --from data/medium --opprint put -q '</b>
<b>  @x_count[$a][$b] += 1;</b>
<b>  @x_sum[$a][$b] += $x;</b>
<b>  end {</b>
<b>      for ((a, b), _ in @x_count) {</b>
<b>          @x_mean[a][b] = @x_sum[a][b] / @x_count[a][b]</b>
<b>      }</b>
<b>      emit (@x_sum, @x_count, @x_mean), "a", "b"</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   x_sum              x_count x_mean
pan pan 219.1851288316854  427     0.5133141190437597
pan wye 198.43293070748447 395     0.5023618498923658
pan eks 216.07522773165525 429     0.5036718595143479
pan hat 205.22277621488686 417     0.492140950155604
pan zee 205.09751802331917 413     0.4966041598627583
eks pan 179.96303047250723 371     0.48507555383425127
eks wye 196.9452860713734  407     0.4838950517724162
eks zee 176.8803651584733  357     0.49546320772681596
eks eks 215.91609712937984 413     0.5227992666570941
eks hat 208.783170520597   417     0.5006790659966355
wye wye 185.29584980261419 377     0.49150092785839306
wye pan 195.84790012056564 392     0.4996119901034838
wye hat 212.0331829346132  426     0.4977304763723314
wye zee 194.77404756708714 385     0.5059066170573692
wye eks 204.8129608356315  386     0.5306035254809106
zee pan 202.21380378504267 389     0.5198298297816007
zee wye 233.9913939194868  455     0.5142667998230479
zee eks 190.9617780631925  391     0.4883932942792647
zee zee 206.64063510417319 403     0.5127559183726382
zee hat 191.30000620900935 409     0.46772617655014515
hat wye 208.8830097609959  423     0.49381326184632596
hat zee 196.3494502965293  385     0.5099985721987774
hat eks 189.0067933716193  389     0.48587864619953547
hat hat 182.8535323148762  381     0.47993053101017374
hat pan 168.5538067327806  363     0.4643355557376876
</pre>

What this does is walk through the first out-of-stream variable (`@x_sum` in this example) as usual, then for each keylist found (e.g. `pan,wye`), include the values for the remaining out-of-stream variables (here, `@x_count` and `@x_mean`). You should use this when all out-of-stream variables in the emit statement have **the same shape and the same keylists**.

## Emit-all statements

Use **emit all** (or `emit @*` which is synonymous) to output all out-of-stream variables. You can use the following idiom to get various accumulators output side-by-side (reminiscent of `mlr stats1`):

<pre class="pre-highlight-in-pair">
<b>mlr --from data/small --opprint put -q '</b>
<b>  @v[$a][$b]["sum"] += $x;</b>
<b>  @v[$a][$b]["count"] += 1;</b>
<b>  end{emit @*,"a","b"}</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   v.sum    v.count
pan pan 0.346791 1
eks pan 0.758679 1
eks wye 0.381399 1
wye wye 0.204603 1
wye pan 0.573288 1
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --from data/small --opprint put -q '</b>
<b>  @sum[$a][$b] += $x;</b>
<b>  @count[$a][$b] += 1;</b>
<b>  end{emit @*,"a","b"}</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   sum
pan pan 0.346791
eks pan 0.758679
eks wye 0.381399
wye wye 0.204603
wye pan 0.573288

a   b   count
pan pan 1
eks pan 1
eks wye 1
wye wye 1
wye pan 1
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --from data/small --opprint put -q '</b>
<b>  @sum[$a][$b] += $x;</b>
<b>  @count[$a][$b] += 1;</b>
<b>  end{emit (@sum, @count),"a","b"}</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   sum      count
pan pan 0.346791 1
eks pan 0.758679 1
eks wye 0.381399 1
wye wye 0.204603 1
wye pan 0.573288 1
</pre>

