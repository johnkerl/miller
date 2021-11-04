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
# Statistics examples

## Computing interquartile ranges

For one or more specified field names, simply compute p25 and p75, then write the IQR as the difference of p75 and p25:

<pre class="pre-highlight-in-pair">
<b>mlr --oxtab stats1 -f x -a p25,p75 \</b>
<b>    then put '$x_iqr = $x_p75 - $x_p25' \</b>
<b>    data/medium </b>
</pre>
<pre class="pre-non-highlight-in-pair">
x_p25 0.24667037823231752
x_p75 0.7481860062358446
x_iqr 0.5015156280035271
</pre>

For wildcarded field names, first compute p25 and p75, then loop over field names with `p25` in them:

<pre class="pre-highlight-in-pair">
<b>mlr --oxtab stats1 --fr '[i-z]' -a p25,p75 \</b>
<b>    then put 'for (k,v in $*) {</b>
<b>      if (k =~ "(.*)_p25") {</b>
<b>        $["\1_iqr"] = $["\1_p75"] - $["\1_p25"]</b>
<b>      }</b>
<b>    }' \</b>
<b>    data/medium </b>
</pre>
<pre class="pre-non-highlight-in-pair">
i_p25 2501
i_p75 7501
x_p25 0.24667037823231752
x_p75 0.7481860062358446
y_p25 0.25213670524015686
y_p75 0.7640028449996572
i_iqr 5000
x_iqr 0.5015156280035271
y_iqr 0.5118661397595003
</pre>

## Computing weighted means

This might be more elegantly implemented as an option within the `stats1` verb. Meanwhile, it's expressible within the DSL:

<pre class="pre-highlight-in-pair">
<b>mlr --from data/medium put -q '</b>
<b>  # Using the y field for weighting in this example</b>
<b>  weight = $y;</b>
<b></b>
<b>  # Using the a field for weighted aggregation in this example</b>
<b>  @sumwx[$a] += weight * $i;</b>
<b>  @sumw[$a] += weight;</b>
<b></b>
<b>  @sumx[$a] += $i;</b>
<b>  @sumn[$a] += 1;</b>
<b></b>
<b>  end {</b>
<b>    map wmean = {};</b>
<b>    map mean  = {};</b>
<b>    for (a in @sumwx) {</b>
<b>      wmean[a] = @sumwx[a] / @sumw[a]</b>
<b>    }</b>
<b>    for (a in @sumx) {</b>
<b>      mean[a] = @sumx[a] / @sumn[a]</b>
<b>    }</b>
<b>    #emit wmean, "a";</b>
<b>    #emit mean, "a";</b>
<b>    emit (wmean, mean), "a";</b>
<b>  }'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,wmean=4979.563722208067,mean=5028.259010091302
a=eks,wmean=4890.3815931472145,mean=4956.2900763358775
a=wye,wmean=4946.987746229947,mean=4920.001017293998
a=zee,wmean=5164.719684856538,mean=5123.092330239375
a=hat,wmean=4925.533162478552,mean=4967.743946419371
</pre>
