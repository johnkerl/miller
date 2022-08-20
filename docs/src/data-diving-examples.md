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
# Data-diving examples

## flins data

The [flins.csv](data/flins.csv) file is some sample data obtained from [https://support.spatialkey.com/spatialkey-sample-csv-data](https://support.spatialkey.com/spatialkey-sample-csv-data).

Vertical-tabular format is good for a quick look at CSV data layout -- seeing what columns you have to work with, as this is a file big enough that we can't just see it on a single screenful:

<pre class="pre-highlight-in-pair">
<b>wc -l data/flins.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
   36635 data/flins.csv
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --c2x --from data/flins.csv head -n 2</b>
</pre>
<pre class="pre-non-highlight-in-pair">
policyID           119736
statecode          FL
county             CLAY COUNTY
eq_site_limit      498960
hu_site_limit      498960
fl_site_limit      498960
fr_site_limit      498960
tiv_2011           498960
tiv_2012           792148.9
eq_site_deductible 0
hu_site_deductible 9979.2
fl_site_deductible 0
fr_site_deductible 0
point_latitude     30.102261
point_longitude    -81.711777
line               Residential
construction       Masonry
point_granularity  1

policyID           448094
statecode          FL
county             CLAY COUNTY
eq_site_limit      1322376.3
hu_site_limit      1322376.3
fl_site_limit      1322376.3
fr_site_limit      1322376.3
tiv_2011           1322376.3
tiv_2012           1438163.57
eq_site_deductible 0
hu_site_deductible 0
fl_site_deductible 0
fr_site_deductible 0
point_latitude     30.063936
point_longitude    -81.707664
line               Residential
construction       Masonry
point_granularity  3
</pre>

A few simple queries:

<pre class="pre-highlight-in-pair">
<b>mlr --c2p --from data/flins.csv count-distinct -f county | head</b>
</pre>
<pre class="pre-non-highlight-in-pair">
county              count
CLAY COUNTY         363
SUWANNEE COUNTY     154
NASSAU COUNTY       135
COLUMBIA COUNTY     125
ST  JOHNS COUNTY    657
BAKER COUNTY        70
BRADFORD COUNTY     31
HAMILTON COUNTY     35
UNION COUNTY        15
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --c2p --from data/flins.csv count-distinct -f line</b>
</pre>
<pre class="pre-non-highlight-in-pair">
line        count
Residential 30838
Commercial  5796
</pre>

Categorization of total insured value:

<pre class="pre-highlight-in-pair">
<b>mlr --c2x --from data/flins.csv stats1 -a min,mean,max -f tiv_2012</b>
</pre>
<pre class="pre-non-highlight-in-pair">
tiv_2012_min  73.37
tiv_2012_mean 2571004.0973420837
tiv_2012_max  1701000000
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --c2p --from data/flins.csv \</b>
<b>  stats1 -a min,mean,max -f tiv_2012 -g construction,line</b>
</pre>
<pre class="pre-non-highlight-in-pair">
construction        line        tiv_2012_min tiv_2012_mean      tiv_2012_max
Masonry             Residential 261168.07    1041986.1292168079 3234970.92
Wood                Residential 73.37        113493.01704925536 649046.12
Reinforced Concrete Commercial  6416016.01   20212428.681839883 60570000
Reinforced Masonry  Commercial  1287817.34   4621372.981117158  16650000
Steel Frame         Commercial  29790000     133492500          1701000000
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --c2x --from data/flins.csv \</b>
<b>  stats1 -a p0,p10,p50,p90,p95,p99,p100 -f hu_site_deductible</b>
</pre>
<pre class="pre-non-highlight-in-pair">
hu_site_deductible_p0   0
hu_site_deductible_p10  0
hu_site_deductible_p50  0
hu_site_deductible_p90  76.5
hu_site_deductible_p95  6829.2
hu_site_deductible_p99  126270
hu_site_deductible_p100 7380000
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --c2p --from data/flins.csv \</b>
<b>  stats1 -a p95,p99,p100 -f hu_site_deductible -g county \</b>
<b>  then sort -f county | head</b>
</pre>
<pre class="pre-non-highlight-in-pair">
county              hu_site_deductible_p95 hu_site_deductible_p99 hu_site_deductible_p100
ALACHUA COUNTY      30630.6                107312.4               1641375
BAKER COUNTY        0                      0                      0
BAY COUNTY          26131.5                181912.5               630000
BRADFORD COUNTY     3355.2                 8163                   8163
BREVARD COUNTY      5360.4                 78975                  1973461.5
BROWARD COUNTY      0                      148500                 3258900
CALHOUN COUNTY      0                      33339.6                33339.6
CHARLOTTE COUNTY    5400                   52650                  250994.7
CITRUS COUNTY       1332.9                 79974.9                483785.1
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --c2x --from data/flins.csv \</b>
<b>  stats2 -a corr,linreg-ols,r2 -f tiv_2011,tiv_2012</b>
</pre>
<pre class="pre-non-highlight-in-pair">
tiv_2011_tiv_2012_corr  0.9730497632351692
tiv_2011_tiv_2012_ols_m 0.9835583980337723
tiv_2011_tiv_2012_ols_b 433854.6428968317
tiv_2011_tiv_2012_ols_n 36634
tiv_2011_tiv_2012_r2    0.9468258417320189
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --c2x --from data/flins.csv --ofmt '%.4f' \</b>
<b>  stats2 -a corr,linreg-ols,r2 -f tiv_2011,tiv_2012 -g county \</b>
<b>  then head -n 5</b>
</pre>
<pre class="pre-non-highlight-in-pair">
county                  CLAY COUNTY
tiv_2011_tiv_2012_corr  0.9627
tiv_2011_tiv_2012_ols_m 1.0901
tiv_2011_tiv_2012_ols_b 46450.5313
tiv_2011_tiv_2012_ols_n 363
tiv_2011_tiv_2012_r2    0.9268

county                  SUWANNEE COUNTY
tiv_2011_tiv_2012_corr  0.9892
tiv_2011_tiv_2012_ols_m 1.0747
tiv_2011_tiv_2012_ols_b 36253.0032
tiv_2011_tiv_2012_ols_n 154
tiv_2011_tiv_2012_r2    0.9785

county                  NASSAU COUNTY
tiv_2011_tiv_2012_corr  0.9731
tiv_2011_tiv_2012_ols_m 1.2963
tiv_2011_tiv_2012_ols_b -45369.2427
tiv_2011_tiv_2012_ols_n 135
tiv_2011_tiv_2012_r2    0.9470

county                  COLUMBIA COUNTY
tiv_2011_tiv_2012_corr  0.9995
tiv_2011_tiv_2012_ols_m 0.9314
tiv_2011_tiv_2012_ols_b 117183.5484
tiv_2011_tiv_2012_ols_n 125
tiv_2011_tiv_2012_r2    0.9990

county                  ST  JOHNS COUNTY
tiv_2011_tiv_2012_corr  0.9662
tiv_2011_tiv_2012_ols_m 1.2301
tiv_2011_tiv_2012_ols_b -596.6239
tiv_2011_tiv_2012_ols_n 657
tiv_2011_tiv_2012_r2    0.9335
</pre>

## Color/shape data

The [data/colored-shapes.dkvp](data/colored-shapes.dkvp) file is some sample data produced by the [mkdat2](../data/mkdat2) script. The idea is:

* Produce some data with known distributions and correlations, and verify that Miller recovers those properties empirically.
* Each record is labeled with one of a few colors and one of a few shapes.
* The `flag` field is 0 or 1, with probability dependent on color
* The `u` field is plain uniform on the unit interval.
* The `v` field is the same, except tightly correlated with `u` for red circles.
* The `w` field is autocorrelated for each color/shape pair.
* The `x` field is boring Gaussian with mean 5 and standard deviation about 1.2, with no dependence on color or shape.

Peek at the data:

<pre class="pre-highlight-in-pair">
<b>wc -l data/colored-shapes.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
   10078 data/colored-shapes.dkvp
</pre>

<pre class="pre-highlight-in-pair">
<b>head -n 6 data/colored-shapes.dkvp | mlr --opprint cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag i   u        v        w        x
yellow triangle 1    56  0.632170 0.988721 0.436498 5.798188
red    square   1    80  0.219668 0.001257 0.792778 2.944117
red    circle   1    84  0.209017 0.290052 0.138103 5.065034
red    square   0    243 0.956274 0.746720 0.775542 7.117831
purple triangle 0    257 0.435535 0.859129 0.812290 5.753095
red    square   0    322 0.201551 0.953110 0.771991 5.612050
</pre>

Look at uncategorized stats (using [creach](https://github.com/johnkerl/scripts/blob/master/fundam/creach) for spacing).

Here it looks reasonable that `u` is unit-uniform; something's up with `v` but we can't yet see what:

<pre class="pre-highlight-in-pair">
<b>mlr --oxtab stats1 -a min,mean,max -f flag,u,v data/colored-shapes.dkvp | creach 3</b>
</pre>
<pre class="pre-non-highlight-in-pair">
flag_min  0
flag_mean 0.39888866838658465
flag_max  1

u_min     0.000044
u_mean    0.49832634262750525
u_max     0.999969

v_min     -0.092709
v_mean    0.49778696586624427
v_max     1.0725

</pre>

The histogram shows the different distribution of 0/1 flags:

<pre class="pre-highlight-in-pair">
<b>mlr --opprint histogram -f flag,u,v --lo -0.1 --hi 1.1 --nbins 12 data/colored-shapes.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
bin_lo                bin_hi              flag_count u_count v_count
-0.010000000000000002 0.09000000000000002 6058       0       36
0.09000000000000002   0.19000000000000003 0          1062    988
0.19000000000000003   0.29000000000000004 0          985     1003
0.29000000000000004   0.39000000000000007 0          1024    1014
0.39000000000000007   0.4900000000000001  0          1002    991
0.4900000000000001    0.5900000000000002  0          989     1041
0.5900000000000002    0.6900000000000002  0          1001    1016
0.6900000000000002    0.7900000000000001  0          972     962
0.7900000000000001    0.8900000000000002  0          1035    1070
0.8900000000000002    0.9900000000000002  0          995     993
0.9900000000000002    1.0900000000000003  4020       1013    939
1.0900000000000003    1.1900000000000002  0          0       25
</pre>

Look at univariate stats by color and shape. In particular, color-dependent flag probabilities pop out, aligning with their original Bernoulli probabilities from the data-generator script:

<pre class="pre-highlight-in-pair">
<b>mlr --opprint stats1 -a min,mean,max -f flag,u,v -g color \</b>
<b>  then sort -f color \</b>
<b>  data/colored-shapes.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  flag_min flag_mean           flag_max u_min    u_mean              u_max    v_min     v_mean              v_max
blue   0        0.5843537414965987  1        0.000044 0.5177171537414964  0.999969 0.001489  0.4910564278911574  0.999576
green  0        0.20919747520288548 1        0.000488 0.5048610595130744  0.999936 0.000501  0.49908475924256035 0.999676
orange 0        0.5214521452145214  1        0.001235 0.49053241584158375 0.998885 0.002449  0.4877637788778878  0.998475
purple 0        0.09019264448336252 1        0.000266 0.49400496322241666 0.999647 0.000364  0.4970507127845888  0.999975
red    0        0.3031674208144796  1        0.000671 0.49255964641241273 0.999882 -0.092709 0.4965350941607402  1.0725
yellow 0        0.8924274593064402  1        0.0013   0.4971291160651098  0.999923 0.000711  0.5106265987261144  0.999919
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --opprint stats1 -a min,mean,max -f flag,u,v -g shape \</b>
<b>  then sort -f shape \</b>
<b>  data/colored-shapes.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
shape    flag_min flag_mean           flag_max u_min    u_mean              u_max    v_min     v_mean              v_max
circle   0        0.3998456194519491  1        0.000044 0.498554505982246   0.999923 -0.092709 0.49552416171362396 1.0725
square   0        0.39611178614823817 1        0.000188 0.4993854558930749  0.999969 0.000089  0.49653825929526124 0.999975
triangle 0        0.4015421115065243  1        0.000881 0.49685854240806604 0.999661 0.000717  0.5010495260972719  0.999995
</pre>

Look at bivariate stats by color and shape. In particular, `u,v` pairwise correlation for red circles pops out:

<pre class="pre-highlight-in-pair">
<b>mlr --opprint --right stats2 -a corr -f u,v,w,x data/colored-shapes.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
          u_v_corr              w_x_corr
0.1334180491027861 -0.011319841199866178
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --opprint --right \</b>
<b>  stats2 -a corr -f u,v,w,x -g color,shape then sort -nr u_v_corr \</b>
<b>  data/colored-shapes.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
 color    shape              u_v_corr               w_x_corr
   red   circle    0.9807984401887236   -0.01856553658708754
orange   square   0.17685855992752927   -0.07104431573806054
 green   circle   0.05764419437577255    0.01179572988801509
   red   square   0.05574477124893523 -0.0006801456507510942
yellow triangle   0.04457273771962798   0.024604310103081825
yellow   square   0.04379172927296089   -0.04462197201631237
purple   circle   0.03587354936895086     0.1341133954140899
  blue   square   0.03241153095761164  -0.053507648119643196
  blue triangle  0.015356427073158766 -0.0006089997461435399
orange   circle  0.010518953877704048   -0.16279397329279383
   red triangle   0.00809782571528034   0.012486621357942596
purple triangle  0.005155190909099334  -0.045057909256220656
purple   square -0.025680276963377404    0.05769429647930396
 green   square   -0.0257760734502851  -0.003265173252087127
orange triangle -0.030456661186085785    -0.1318699981926352
yellow   circle  -0.06477331572781474    0.07369449819706045
  blue   circle  -0.10234761901929677  -0.030528539069837757
 green triangle  -0.10901825107358765   -0.04848782060162929
</pre>
