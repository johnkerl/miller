<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# I/O options

## Formats

Options:

<pre class="pre-non-highlight-non-pair">
--dkvp    --idkvp    --odkvp
--nidx    --inidx    --onidx
--csv     --icsv     --ocsv
--csvlite --icsvlite --ocsvlite
--pprint  --ipprint  --opprint  --right
--xtab    --ixtab    --oxtab
--json    --ijson    --ojson
</pre>

These are as discussed in [File Formats](file-formats.md), with the exception of `--right` which makes pretty-printed output right-aligned:

<pre class="pre-highlight-in-pair">
<b>mlr --opprint cat data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x        y
pan pan 1 0.346791 0.726802
eks pan 2 0.758679 0.522151
wye wye 3 0.204603 0.338318
eks wye 4 0.381399 0.134188
wye pan 5 0.573288 0.863624
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --opprint --right cat data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
  a   b i        x        y 
pan pan 1 0.346791 0.726802 
eks pan 2 0.758679 0.522151 
wye wye 3 0.204603 0.338318 
eks wye 4 0.381399 0.134188 
wye pan 5 0.573288 0.863624 
</pre>

Additional notes:

* Use `--csv`, `--pprint`, etc. when the input and output formats are the same.

* Use `--icsv --opprint`, etc. when you want format conversion as part of what Miller does to your data.

* DKVP (key-value-pair) format is the default for input and output. So, `--oxtab` is the same as `--idkvp --oxtab`.

**Pro-tip:** Please use either **--format1**, or **--iformat1 --oformat2**.  If you use **--format1 --oformat2** then what happens is that flags are set up for input *and* output for format1, some of which are overwritten for output in format2. For technical reasons, having `--oformat2` clobber all the output-related effects of `--format1` also removes some flexibility from the command-line interface. See also Miller issues [180](https://github.com/johnkerl/miller/issues/180) and [199](https://github.com/johnkerl/miller/issues/199).

## In-place mode

Use the `mlr -I` flag to process files in-place. For example, `mlr -I --csv cut -x -f unwanted_column_name mydata/*.csv` will remove `unwanted_column_name` from all your `*.csv` files in your `mydata/` subdirectory.

By default, Miller output goes to the screen (or you can redirect a file using `>` or to another process using `|`). With `-I`, for each file name on the command line, output is written to a temporary file in the same directory. Miller writes its output into that temp file, which is then renamed over the original.  Then, processing continues on the next file. Each file is processed in isolation: if the output format is CSV, CSV headers will be present in each output file; statistics are only over each file's own records; and so on.

Please see [Choices for printing to files](10min.md#choices-for-printing-to-files) for examples.

## Compression

See the separate page on [Compressed data](reference-main-compressed-data.md).

## Record/field/pair separators

Miller has record separators `IRS` and `ORS`, field separators `IFS` and `OFS`, and pair separators `IPS` and `OPS`.  For example, in the DKVP line `a=1,b=2,c=3`, the record separator is newline, field separator is comma, and pair separator is the equals sign. These are the default values.

Options:

<pre class="pre-non-highlight-non-pair">
--rs --irs --ors
--fs --ifs --ofs --repifs
--ps --ips --ops
</pre>

* You can change a separator from input to output via e.g. `--ifs = --ofs :`. Or, you can specify that the same separator is to be used for input and output via e.g. `--fs :`.

* The pair separator is only relevant to DKVP format.

* Pretty-print and xtab formats ignore the separator arguments altogether.

* The `--repifs` means that multiple successive occurrences of the field separator count as one.  For example, in CSV data we often signify nulls by empty strings, e.g. `2,9,,,,,6,5,4`. On the other hand, if the field separator is a space, it might be more natural to parse `2 4    5` the same as `2 4 5`: `--repifs --ifs ' '` lets this happen.  In fact, the `--ipprint` option above is internally implemented in terms of `--repifs`.

* Just write out the desired separator, e.g. `--ofs '|'`. But you may use the symbolic names `newline`, `space`, `tab`, `pipe`, or `semicolon` if you like.

## Number formatting

The command-line option `--ofmt {format string}` is the global number format for commands which generate numeric output, e.g. `stats1`, `stats2`, `histogram`, and `step`, as well as `mlr put`. Examples:

<pre class="pre-non-highlight-non-pair">
--ofmt %.9e  --ofmt %.6f  --ofmt %.0f
</pre>

These are just familiar `printf` formats.  (TODO: write about type-checking once that's implemented.) Additionally, if you use leading width (e.g. `%18.12f`) then the output will contain embedded whitespace, which may not be what you want if you pipe the output to something else, particularly CSV. I use Miller's pretty-print format (`mlr --opprint`) to column-align numerical data.

To apply formatting to a single field, overriding the global `ofmt`, use `fmtnum` function within `mlr put`. For example:

<pre class="pre-highlight-in-pair">
<b>echo 'x=3.1,y=4.3' | mlr put '$z=fmtnum($x*$y,"%08f")'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x=3.1,y=4.3,z=13.330000
</pre>

<pre class="pre-highlight-in-pair">
<b>echo 'x=0xffff,y=0xff' | mlr put '$z=fmtnum(int($x*$y),"%08x")'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x=0xffff,y=0xff,z=00feff01
</pre>

Input conversion from hexadecimal is done automatically on fields handled by `mlr put` and `mlr filter` as long as the field value begins with "0x".  To apply output conversion to hexadecimal on a single column, you may use `fmtnum`, or the keystroke-saving `hexfmt` function. Example:

<pre class="pre-highlight-in-pair">
<b>echo 'x=0xffff,y=0xff' | mlr put '$z=hexfmt($x*$y)'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x=0xffff,y=0xff,z=0xfeff01
</pre>
