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
# CSV, with and without headers

## Headerless CSV on input or output

Sometimes we get CSV files which lack a header. For example, [data/headerless.csv](./data/headerless.csv):

<pre class="pre-highlight-in-pair">
<b>cat data/headerless.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
John,23,present
Fred,34,present
Alice,56,missing
Carol,45,present
</pre>

You can use Miller to add a header. The `--implicit-csv-header` applies positionally indexed labels:

<pre class="pre-highlight-in-pair">
<b>mlr --csv --implicit-csv-header cat data/headerless.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1,2,3
John,23,present
Fred,34,present
Alice,56,missing
Carol,45,present
</pre>

Following that, you can rename the positionally indexed labels to names with meaning for your context.  For example:

<pre class="pre-highlight-in-pair">
<b>mlr --csv --implicit-csv-header label name,age,status data/headerless.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
name,age,status
John,23,present
Fred,34,present
Alice,56,missing
Carol,45,present
</pre>

Likewise, if you need to produce CSV which is lacking its header, you can pipe Miller's output to the system command `sed 1d`, or you can use Miller's `--headerless-csv-output` option:

<pre class="pre-highlight-in-pair">
<b>head -5 data/colored-shapes.dkvp | mlr --ocsv cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,i,u,v,w,x
yellow,triangle,1,56,0.632170,0.988721,0.436498,5.798188
red,square,1,80,0.219668,0.001257,0.792778,2.944117
red,circle,1,84,0.209017,0.290052,0.138103,5.065034
red,square,0,243,0.956274,0.746720,0.775542,7.117831
purple,triangle,0,257,0.435535,0.859129,0.812290,5.753095
</pre>

<pre class="pre-highlight-in-pair">
<b>head -5 data/colored-shapes.dkvp | mlr --ocsv --headerless-csv-output cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
yellow,triangle,1,56,0.632170,0.988721,0.436498,5.798188
red,square,1,80,0.219668,0.001257,0.792778,2.944117
red,circle,1,84,0.209017,0.290052,0.138103,5.065034
red,square,0,243,0.956274,0.746720,0.775542,7.117831
purple,triangle,0,257,0.435535,0.859129,0.812290,5.753095
</pre>

Lastly, often we say "CSV" or "TSV" when we have positionally indexed data in columns which are separated by commas or tabs, respectively. In this case it's perhaps simpler to **just use NIDX format** which was designed for this purpose. (See also [File Formats](file-formats.md).) For example:

<pre class="pre-highlight-in-pair">
<b>mlr --inidx --ifs comma --oxtab cut -f 1,3 data/headerless.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1 John
3 present

1 Fred
3 present

1 Alice
3 missing

1 Carol
3 present
</pre>

## Headerless CSV with duplicate field values

Miller is (by central design) a mapping from name to value, rather than integer position to value as in most tools in the Unix toolkit such as `sort`, `cut`, `awk`, etc. So given input `Yea=1,Yea=2` on the same input line, first `Yea=1` is stored, then updated with `Yea=2`. This is in the input-parser and the value `Yea=1` is unavailable to any further processing. The following example line comes from a headerless CSV file and includes 5 times the string (value) `'NA'`:

<pre class="pre-highlight-in-pair">
<b>ag '0.9' nas.csv | head -1</b>
</pre>
<pre class="pre-non-highlight-in-pair">
2:-349801.10097848,4537221.43295653,2,1,NA,NA,NA,NA,NA
</pre>

The repeated `'NA'` strings (values) in the same line will be treated as fields (columns) with same name, thus only one is kept in the output.

This can be worked around by telling `mlr` that there is no header row by using `--implicit-csv-header` or changing the input format by using `nidx` like so:

<pre class="pre-non-highlight-non-pair">
ag '0.9' nas.csv | mlr --n2c --fs "," label xsn,ysn,x,y,t,a,e29,e31,e32 then head
</pre>

## Regularizing ragged CSV

Miller handles compliant CSV: in particular, it's an error if the number of data fields in a given data line don't match the number of header lines. But in the event that you have a CSV file in which some lines have less than the full number of fields, you can use Miller to pad them out. The trick is to use NIDX format, for which each line stands on its own without respect to a header line.

<pre class="pre-highlight-in-pair">
<b>cat data/ragged.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c
1,2,3
4,5
6,7,8,9
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --from data/ragged.csv --fs comma --nidx put '</b>
<b>  @maxnf = max(@maxnf, NF);</b>
<b>  @nf = NF;</b>
<b>  while(@nf < @maxnf) {</b>
<b>    @nf += 1;</b>
<b>    $[@nf] = ""</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c
1,2,3
4,5
6,7,8,9
</pre>

or, more simply,

<pre class="pre-highlight-in-pair">
<b>mlr --from data/ragged.csv --fs comma --nidx put '</b>
<b>  @maxnf = max(@maxnf, NF);</b>
<b>  while(NF < @maxnf) {</b>
<b>    $[NF+1] = "";</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c
1,2,3
4,5
6,7,8,9
</pre>

See also the [record-heterogeneity page](record-heterogeneity.md).
