<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# CSV, with and without headers

## Headerless CSV on input or output

Sometimes we get CSV files which lack a header. For example (`data/headerless.csv <./data/headerless.csv>`_):

<pre>
<b>cat data/headerless.csv</b>
John,23,present
Fred,34,present
Alice,56,missing
Carol,45,present
</pre>

You can use Miller to add a header. The ``--implicit-csv-header`` applies positionally indexed labels:

<pre>
<b>mlr --csv --implicit-csv-header cat data/headerless.csv</b>
1,2,3
John,23,present
Fred,34,present
Alice,56,missing
Carol,45,present
</pre>

Following that, you can rename the positionally indexed labels to names with meaning for your context.  For example:

<pre>
<b>mlr --csv --implicit-csv-header label name,age,status data/headerless.csv</b>
name,age,status
John,23,present
Fred,34,present
Alice,56,missing
Carol,45,present
</pre>

Likewise, if you need to produce CSV which is lacking its header, you can pipe Miller's output to the system command ``sed 1d``, or you can use Miller's ``--headerless-csv-output`` option:

<pre>
<b>head -5 data/colored-shapes.dkvp | mlr --ocsv cat</b>
color,shape,flag,i,u,v,w,x
yellow,triangle,1,11,0.6321695890307647,0.9887207810889004,0.4364983936735774,5.7981881667050565
red,square,1,15,0.21966833570651523,0.001257332190235938,0.7927778364718627,2.944117399716207
red,circle,1,16,0.20901671281497636,0.29005231936593445,0.13810280912907674,5.065034003400998
red,square,0,48,0.9562743938458542,0.7467203085342884,0.7755423050923582,7.117831369597269
purple,triangle,0,51,0.4355354501763202,0.8591292672156728,0.8122903963006748,5.753094629505863
</pre>

<pre>
<b>head -5 data/colored-shapes.dkvp | mlr --ocsv --headerless-csv-output cat</b>
yellow,triangle,1,11,0.6321695890307647,0.9887207810889004,0.4364983936735774,5.7981881667050565
red,square,1,15,0.21966833570651523,0.001257332190235938,0.7927778364718627,2.944117399716207
red,circle,1,16,0.20901671281497636,0.29005231936593445,0.13810280912907674,5.065034003400998
red,square,0,48,0.9562743938458542,0.7467203085342884,0.7755423050923582,7.117831369597269
purple,triangle,0,51,0.4355354501763202,0.8591292672156728,0.8122903963006748,5.753094629505863
</pre>

Lastly, often we say "CSV" or "TSV" when we have positionally indexed data in columns which are separated by commas or tabs, respectively. In this case it's perhaps simpler to **just use NIDX format** which was designed for this purpose. (See also :doc:`file-formats`.) For example:

<pre>
<b>mlr --inidx --ifs comma --oxtab cut -f 1,3 data/headerless.csv</b>
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

Miller is (by central design) a mapping from name to value, rather than integer position to value as in most tools in the Unix toolkit such as ``sort``, ``cut``, ``awk``, etc. So given input ``Yea=1,Yea=2`` on the same input line, first ``Yea=1`` is stored, then updated with ``Yea=2``. This is in the input-parser and the value ``Yea=1`` is unavailable to any further processing. The following example line comes from a headerless CSV file and includes 5 times the string (value) ``'NA'``:

<pre>
<b>ag '0.9' nas.csv | head -1</b>
2:-349801.10097848,4537221.43295653,2,1,NA,NA,NA,NA,NA
</pre>

The repeated ``'NA'`` strings (values) in the same line will be treated as fields (columns) with same name, thus only one is kept in the output.

This can be worked around by telling ``mlr`` that there is no header row by using ``--implicit-csv-header`` or changing the input format by using ``nidx`` like so:

<pre>
ag '0.9' nas.csv | mlr --n2c --fs "," label xsn,ysn,x,y,t,a,e29,e31,e32 then head
</pre>

## Regularizing ragged CSV

Miller handles compliant CSV: in particular, it's an error if the number of data fields in a given data line don't match the number of header lines. But in the event that you have a CSV file in which some lines have less than the full number of fields, you can use Miller to pad them out. The trick is to use NIDX format, for which each line stands on its own without respect to a header line.

<pre>
<b>cat data/ragged.csv</b>
a,b,c
1,2,3
4,5
6,7,8,9
</pre>

<pre>
<b>mlr --from data/ragged.csv --fs comma --nidx put '</b>
<b>  @maxnf = max(@maxnf, NF);</b>
<b>  @nf = NF;</b>
<b>  while(@nf < @maxnf) {</b>
<b>    @nf += 1;</b>
<b>    $[@nf] = ""</b>
<b>  }</b>
<b>'</b>
a,b,c
1,2,3
4,5
6,7,8,9
</pre>

or, more simply,

<pre>
<b>mlr --from data/ragged.csv --fs comma --nidx put '</b>
<b>  @maxnf = max(@maxnf, NF);</b>
<b>  while(NF < @maxnf) {</b>
<b>    $[NF+1] = "";</b>
<b>  }</b>
<b>'</b>
a,b,c
1,2,3
4,5
6,7,8,9
</pre>
