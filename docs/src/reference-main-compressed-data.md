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
# Compressed data

As of [Miller 6](new-in-miller-6.md), Miller supports reading GZIP, BZIP2, and
ZLIB formats transparently, and in-process. And (as before Miller 6) you have a
more general `--prepipe` option to support other decompression programs.

## Automatic detection on input

If your files end in `.gz`, `.bz2`, or `.z` then Miller will autodetect by file extension:

<pre class="pre-highlight-in-pair">
<b>file gz-example.csv.gz</b>
</pre>
<pre class="pre-non-highlight-in-pair">
gz-example.csv.gz: gzip compressed data, was "gz-example.csv", last modified: Mon Aug 23 02:04:34 2021, from Unix, original size modulo 2^32 429
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv sort -f color gz-example.csv.gz</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
purple,triangle,false,5,51,81.2290,8.5910
purple,triangle,false,7,65,80.1405,5.8240
purple,square,false,10,91,72.3735,8.2430
red,square,true,2,15,79.2778,0.0130
red,circle,true,3,16,13.8103,2.9010
red,square,false,4,48,77.5542,7.4670
red,square,false,6,64,77.1991,9.5310
yellow,triangle,true,1,11,43.6498,9.8870
yellow,circle,true,8,73,63.9785,4.2370
yellow,circle,true,9,87,63.5058,8.3350
</pre>

This will decompress the input data on the fly, while leaving the disk file unmodified. This helps you save disk space, at the cost of some additional runtime CPU usage to decompress the data.

## Manual detection on input

If the filename doesn't in in `.gz`, `.bz2`, or `.z` then you can use the flags `--gzin`, `--bz2in`, or `--zin` to let Miller know:

<pre class="pre-highlight-non-pair">
<b>mlr --csv --gzin sort -f color myfile.bin # myfile.bin has gzip contents</b>
</pre>

## External decompressors on input

Using the `--prepipe` flag, you can provide the name of any decompression
program in your `$PATH` and Miller will run it on each input file, effectively
piping the standard output of that program to Miller's standard input.

You can, of course, already do without this for single input files, for example:

<pre class="pre-highlight-in-pair">
<b>gunzip < gz-example.csv.gz | mlr --csv sort -f color</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
purple,triangle,false,5,51,81.2290,8.5910
purple,triangle,false,7,65,80.1405,5.8240
purple,square,false,10,91,72.3735,8.2430
red,square,true,2,15,79.2778,0.0130
red,circle,true,3,16,13.8103,2.9010
red,square,false,4,48,77.5542,7.4670
red,square,false,6,64,77.1991,9.5310
yellow,triangle,true,1,11,43.6498,9.8870
yellow,circle,true,8,73,63.9785,4.2370
yellow,circle,true,9,87,63.5058,8.3350
</pre>

The benefit of `--prepipe` is that Miller will run the specified program once per
file, respecting file boundaries.

The prepipe command can be anything which reads from standard input and produces
data acceptable to Miller. Nominally this allows you to use whichever
decompression utilities you have installed on your system, on a per-file basis.

If the command has flags, quote them: e.g. `mlr --prepipe 'zcat -cf'`.

In your [.mlrrc file](customization.md), `--prepipe` and `--prepipex` are not
allowed as they could be used for unexpected code execution. You can use
`--prepipe-bz2`, `--prepipe-gunzip`, and `--prepipe-zcat` in `.mlrrc`, though.

Note that this feature is quite general and is not limited to decompression
utilities. You can use it to apply per-file filters of your choice: e.g. `mlr
--prepipe 'head -n 10' ...`, if you like.

There is a `--prepipe` and a `--prepipex`:

* If the command normally runs with `nameofprogram < filename.ext` (such as `gunzip` or `zcat -cf` or `xz -cd`) then use `--prepipe`.
* If the command normally runs with `nameofprogram filename.ext` (such as `unzip -qc`) then use `--prepipex`.

Lastly, note that if `--prepipe` or `--prepipex` is specified on the Miller
command line, it replaces any autodetect decisions that might have been made
based on the filename extension. Likewise, `--gzin`/`--bz2in`/`--zin` are ignored if
`--prepipe` or `--prepipex` is also specified.

## Compressed output

Everything said so far on this page has to do with compressed input.

For compressed output:

* Normally Miller output is to stdout, so you can pipe the output: `mlr sort -n quantity foo.csv | gzip > sorted.csv.gz`.

* For [`tee` statements](reference-dsl-output-statements.md#tee-statements), which write output to files rather than stdout, use `tee`'s redirect syntax:

<!---
  gzip by default puts timestamp into the header, so every regen of these *.md.in files makes a modified
  file, which is annoying for version control. That can be suppressed by using 'gzip -n' but then that's
  confusing for the reader, who has no need for -n. We handle this by making this code sample non-live.
--->
<pre class="pre-highlight-non-pair">
<b>mlr --from example.csv --csv put -q '</b>
<b>  filename = $color.".csv.gz";</b>
<b>  tee | "gzip > ".filename, $*</b>
<b>'</b>
</pre>

<pre class="pre-highlight-in-pair">
<b>file red.csv.gz purple.csv.gz yellow.csv.gz</b>
</pre>
<pre class="pre-non-highlight-in-pair">
red.csv.gz:    gzip compressed data, last modified: Mon Aug 23 02:34:05 2021, from Unix, original size modulo 2^32 185
purple.csv.gz: gzip compressed data, last modified: Mon Aug 23 02:34:05 2021, from Unix, original size modulo 2^32 164
yellow.csv.gz: gzip compressed data, last modified: Mon Aug 23 02:34:05 2021, from Unix, original size modulo 2^32 158
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv cat yellow.csv.gz</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
yellow,triangle,true,1,11,43.6498,9.8870
yellow,circle,true,8,73,63.9785,4.2370
yellow,circle,true,9,87,63.5058,8.3350
</pre>

* Using the [in-place flag](reference-main-in-place-processing.md) `-I`, the overwritten file will
be compressed when possible. See the [page on in-place mode](reference-main-in-place-processing.md) for details.
