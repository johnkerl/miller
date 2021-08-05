<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# Data-diving examples

## flins data

The [flins.csv](data/flins.csv) file is some sample data obtained from [https://support.spatialkey.com/spatialkey-sample-csv-data](https://support.spatialkey.com/spatialkey-sample-csv-data).

Vertical-tabular format is good for a quick look at CSV data layout -- seeing what columns you have to work with:

<pre class="pre-highlight">
<b>head -n 2 data/flins.csv | mlr --icsv --oxtab cat</b>
</pre>
<pre class="pre-non-highlight">
county   Seminole
tiv_2011 22890.55
tiv_2012 20848.71
line     Residential
</pre>

A few simple queries:

<pre class="pre-highlight">
<b>mlr --from data/flins.csv --icsv --opprint count-distinct -f county | head</b>
</pre>
<pre class="pre-non-highlight">
county     count
Seminole   1
Miami Dade 2
Palm Beach 1
Highlands  2
Duval      1
St. Johns  1
</pre>

<pre class="pre-highlight">
<b>mlr --from data/flins.csv --icsv --opprint count-distinct -f construction,line</b>
</pre>

Categorization of total insured value:

<pre class="pre-highlight">
<b>mlr --from data/flins.csv --icsv --opprint stats1 -a min,mean,max -f tiv_2012</b>
</pre>
<pre class="pre-non-highlight">
tiv_2012_min tiv_2012_mean      tiv_2012_max
19757.91     1061531.4637499999 2785551.63
</pre>

<pre class="pre-highlight">
<b>mlr --from data/flins.csv --icsv --opprint \</b>
<b>  stats1 -a min,mean,max -f tiv_2012 -g construction,line</b>
</pre>

<pre class="pre-highlight">
<b>mlr --from data/flins.csv --icsv --oxtab \</b>
<b>  stats1 -a p0,p10,p50,p90,p95,p99,p100 -f hu_site_deductible</b>
</pre>

<pre class="pre-highlight">
<b>mlr --from data/flins.csv --icsv --opprint \</b>
<b>  stats1 -a p95,p99,p100 -f hu_site_deductible -g county \</b>
<b>  then sort -f county | head</b>
</pre>
<pre class="pre-non-highlight">
county
Duval
Highlands
Miami Dade
Palm Beach
Seminole
St. Johns
</pre>

<pre class="pre-highlight">
<b>mlr --from data/flins.csv --icsv --oxtab \</b>
<b>  stats2 -a corr,linreg-ols,r2 -f tiv_2011,tiv_2012</b>
</pre>
<pre class="pre-non-highlight">
tiv_2011_tiv_2012_corr  0.9353629581411828
tiv_2011_tiv_2012_ols_m 1.0890905877734807
tiv_2011_tiv_2012_ols_b 103095.52335638746
tiv_2011_tiv_2012_ols_n 8
tiv_2011_tiv_2012_r2    0.8749038634626236
</pre>

<pre class="pre-highlight">
<b>mlr --from data/flins.csv --icsv --opprint \</b>
<b>  stats2 -a corr,linreg-ols,r2 -f tiv_2011,tiv_2012 -g county</b>
</pre>
<pre class="pre-non-highlight">
county     tiv_2011_tiv_2012_corr tiv_2011_tiv_2012_ols_m tiv_2011_tiv_2012_ols_b tiv_2011_tiv_2012_ols_n tiv_2011_tiv_2012_r2
Seminole   -                      -                       -                       1                       -
Miami Dade 1                      0.9306426512386247      -2311.1543275160047     2                       0.9999999999999999
Palm Beach -                      -                       -                       1                       -
Highlands  0.9999999999999997     1.055692910750992       -4529.7939388307705     2                       0.9999999999999992
Duval      -                      -                       -                       1                       -
St. Johns  -                      -                       -                       1                       -
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

<pre class="pre-highlight">
<b>wc -l data/colored-shapes.dkvp</b>
</pre>
<pre class="pre-non-highlight">
   10078 data/colored-shapes.dkvp
</pre>

<pre class="pre-highlight">
<b>head -n 6 data/colored-shapes.dkvp | mlr --opprint cat</b>
</pre>
<pre class="pre-non-highlight">
color  shape    flag i  u                   v                    w                   x
yellow triangle 1    11 0.6321695890307647  0.9887207810889004   0.4364983936735774  5.7981881667050565
red    square   1    15 0.21966833570651523 0.001257332190235938 0.7927778364718627  2.944117399716207
red    circle   1    16 0.20901671281497636 0.29005231936593445  0.13810280912907674 5.065034003400998
red    square   0    48 0.9562743938458542  0.7467203085342884   0.7755423050923582  7.117831369597269
purple triangle 0    51 0.4355354501763202  0.8591292672156728   0.8122903963006748  5.753094629505863
red    square   0    64 0.2015510269821953  0.9531098083420033   0.7719912015786777  5.612050466474166
</pre>

Look at uncategorized stats (using [creach](https://github.com/johnkerl/scripts/blob/master/fundam/creach) for spacing).

Here it looks reasonable that `u` is unit-uniform; something's up with `v` but we can't yet see what:

<pre class="pre-highlight">
<b>mlr --oxtab stats1 -a min,mean,max -f flag,u,v data/colored-shapes.dkvp | creach 3</b>
</pre>
<pre class="pre-non-highlight">
flag_min  0
flag_mean 0.39888866838658465
flag_max  1

u_min     0.000043912454007477564
u_mean    0.4983263438118866
u_max     0.9999687954968421

v_min     -0.09270905318501277
v_mean    0.49778696527477023
v_max     1.0724998185026013
</pre>

The histogram shows the different distribution of 0/1 flags:

<pre class="pre-highlight">
<b>mlr --opprint histogram -f flag,u,v --lo -0.1 --hi 1.1 --nbins 12 data/colored-shapes.dkvp</b>
</pre>
<pre class="pre-non-highlight">
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

Look at univariate stats by color and shape. In particular, color-dependent flag probabilities pop out, aligning with their original Bernoulli probablities from the data-generator script:

<pre class="pre-highlight">
<b>mlr --opprint stats1 -a min,mean,max -f flag,u,v -g color \</b>
<b>  then sort -f color \</b>
<b>  data/colored-shapes.dkvp</b>
</pre>
<pre class="pre-non-highlight">
color  flag_min flag_mean           flag_max u_min                   u_mean              u_max              v_min                 v_mean              v_max
blue   0        0.5843537414965987  1        0.000043912454007477564 0.517717155039078   0.9999687954968421 0.0014886830387470518 0.49105642841387653 0.9995761761685742
green  0        0.20919747520288548 1        0.00048750676198217047  0.5048610622924616  0.9999361779701204 0.0005012669003675585 0.49908475928072205 0.9996764373885353
orange 0        0.5214521452145214  1        0.00123537823160913     0.49053241689014415 0.9988853487546249 0.0024486660337188493 0.4877637745987629  0.998475130432018
purple 0        0.09019264448336252 1        0.0002655214518428872   0.4940049543793683  0.9996465731736793 0.0003641137096487279 0.497050699948439   0.9999751864255598
red    0        0.3031674208144796  1        0.0006711367180041172   0.49255964831571375 0.9998822102016469 -0.09270905318501277  0.4965350959465078  1.0724998185026013
yellow 0        0.8924274593064402  1        0.001300228762057487    0.49712912165196765 0.99992313390574   0.0007109695568577878 0.510626599360317   0.9999189897724752
</pre>

<pre class="pre-highlight">
<b>mlr --opprint stats1 -a min,mean,max -f flag,u,v -g shape \</b>
<b>  then sort -f shape \</b>
<b>  data/colored-shapes.dkvp</b>
</pre>
<pre class="pre-non-highlight">
shape    flag_min flag_mean           flag_max u_min                   u_mean              u_max              v_min                  v_mean              v_max
circle   0        0.3998456194519491  1        0.000043912454007477564 0.49855450951394115 0.99992313390574   -0.09270905318501277   0.49552415740048406 1.0724998185026013
square   0        0.39611178614823817 1        0.0001881939925673093   0.499385458061097   0.9999687954968421 0.00008930277299445954 0.49653825501903986 0.9999751864255598
triangle 0        0.4015421115065243  1        0.000881025170573424    0.4968585405884252  0.9996614910922645 0.000716883409890845   0.501049532862137   0.9999946837499262
</pre>

Look at bivariate stats by color and shape. In particular, `u,v` pairwise correlation for red circles pops out:

<pre class="pre-highlight">
<b>mlr --opprint --right stats2 -a corr -f u,v,w,x data/colored-shapes.dkvp</b>
</pre>
<pre class="pre-non-highlight">
           u_v_corr              w_x_corr 
0.13341803768384553 -0.011319938208638764 
</pre>

<pre class="pre-highlight">
<b>mlr --opprint --right \</b>
<b>  stats2 -a corr -f u,v,w,x -g color,shape then sort -nr u_v_corr \</b>
<b>  data/colored-shapes.dkvp</b>
</pre>
<pre class="pre-non-highlight">
 color    shape              u_v_corr               w_x_corr 
   red   circle    0.9807984157534667  -0.018565046320623148 
orange   square   0.17685846147882145   -0.07104374629148885 
 green   circle   0.05764430126828069   0.011795210176784067 
   red   square  0.055744791559722166 -0.0006802175149145207 
yellow triangle   0.04457267106380469    0.02460476240108526 
yellow   square   0.04379171794446621   -0.04462267239937856 
purple   circle   0.03587354791796681    0.13411247530136805 
  blue   square   0.03241156493114544   -0.05350791240143263 
  blue triangle  0.015356295190464324 -0.0006084778850362686 
orange   circle   0.01051866723398945    -0.1627949723421722 
   red triangle   0.00809781003735548   0.012485753551391776 
purple triangle  0.005155038421780437   -0.04505792148014131 
purple   square  -0.02568020549187632    0.05769444883779078 
 green   square -0.025775985300150128  -0.003265248022084335 
orange triangle -0.030456930370361554     -0.131870019629393 
yellow   circle  -0.06477338560056926    0.07369474300245252 
  blue   circle   -0.1023476302678634  -0.030529007506883508 
 green triangle  -0.10901830007460846    -0.0484881707807228 
</pre>
