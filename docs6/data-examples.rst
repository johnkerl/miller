..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Data-diving examples
================================================================

flins data
----------------------------------------------------------------

The `flins.csv <data/flins.csv>`_ file is some sample data obtained from https://support.spatialkey.com/spatialkey-sample-csv-data.

Vertical-tabular format is good for a quick look at CSV data layout -- seeing what columns you have to work with:

.. code-block::
   :emphasize-lines: 1,1

    $ head -n 2 data/flins.csv | mlr --icsv --oxtab cat
    county   Seminole
    tiv_2011 22890.55
    tiv_2012 20848.71
    line     Residential

A few simple queries:

.. code-block::
   :emphasize-lines: 1,1

    $ mlr --from data/flins.csv --icsv --opprint count-distinct -f county | head
    county     count
    Seminole   1
    Miami Dade 2
    Palm Beach 1
    Highlands  2
    Duval      1
    St. Johns  1

.. code-block::
   :emphasize-lines: 1,1

    $ mlr --from data/flins.csv --icsv --opprint count-distinct -f construction,line

Categorization of total insured value:

.. code-block::
   :emphasize-lines: 1,1

    $ mlr --from data/flins.csv --icsv --opprint stats1 -a min,mean,max -f tiv_2012
    tiv_2012_min tiv_2012_mean  tiv_2012_max
    19757.910000 1061531.463750 2785551.630000

.. code-block::
   :emphasize-lines: 1,1

    $ mlr --from data/flins.csv --icsv --opprint stats1 -a min,mean,max -f tiv_2012 -g construction,line

.. code-block::
   :emphasize-lines: 1,1

    $ mlr --from data/flins.csv --icsv --oxtab stats1 -a p0,p10,p50,p90,p95,p99,p100 -f hu_site_deductible
    hu_site_deductible_p0   
    hu_site_deductible_p10  
    hu_site_deductible_p50  
    hu_site_deductible_p90  
    hu_site_deductible_p95  
    hu_site_deductible_p99  
    hu_site_deductible_p100 

.. code-block::
   :emphasize-lines: 1,1

    $ mlr --from data/flins.csv --icsv --opprint stats1 -a p95,p99,p100 -f hu_site_deductible -g county then sort -f county | head
    county     hu_site_deductible_p95 hu_site_deductible_p99 hu_site_deductible_p100
    Duval      -                      -                      -
    Highlands  -                      -                      -
    Miami Dade -                      -                      -
    Palm Beach -                      -                      -
    Seminole   -                      -                      -
    St. Johns  -                      -                      -

.. code-block::
   :emphasize-lines: 1,1

    $ mlr --from data/flins.csv --icsv --oxtab stats2 -a corr,linreg-ols,r2 -f tiv_2011,tiv_2012
    tiv_2011_tiv_2012_corr  0.935363
    tiv_2011_tiv_2012_ols_m 1.089091
    tiv_2011_tiv_2012_ols_b 103095.523356
    tiv_2011_tiv_2012_ols_n 8
    tiv_2011_tiv_2012_r2    0.874904

.. code-block::
   :emphasize-lines: 1,1

    $ mlr --from data/flins.csv --icsv --opprint stats2 -a corr,linreg-ols,r2 -f tiv_2011,tiv_2012 -g county
    county     tiv_2011_tiv_2012_corr tiv_2011_tiv_2012_ols_m tiv_2011_tiv_2012_ols_b tiv_2011_tiv_2012_ols_n tiv_2011_tiv_2012_r2
    Seminole   -                      -                       -                       1                       -
    Miami Dade 1.000000               0.930643                -2311.154328            2                       1.000000
    Palm Beach -                      -                       -                       1                       -
    Highlands  1.000000               1.055693                -4529.793939            2                       1.000000
    Duval      -                      -                       -                       1                       -
    St. Johns  -                      -                       -                       1                       -

Color/shape data
----------------------------------------------------------------

The `colored-shapes.dkvp <https://github.com/johnkerl/miller/blob/master/docs/data/colored-shapes.dkvp>`_ file is some sample data produced by the `mkdat2 <https://github.com/johnkerl/miller/blob/master/doc/datagen/mkdat2>`_ script. The idea is:

* Produce some data with known distributions and correlations, and verify that Miller recovers those properties empirically.
* Each record is labeled with one of a few colors and one of a few shapes.
* The ``flag`` field is 0 or 1, with probability dependent on color
* The ``u`` field is plain uniform on the unit interval.
* The ``v`` field is the same, except tightly correlated with ``u`` for red circles.
* The ``w`` field is autocorrelated for each color/shape pair.
* The ``x`` field is boring Gaussian with mean 5 and standard deviation about 1.2, with no dependence on color or shape.

Peek at the data:

.. code-block::
   :emphasize-lines: 1,1

    $ wc -l data/colored-shapes.dkvp
       10078 data/colored-shapes.dkvp

.. code-block::
   :emphasize-lines: 1,1

    $ head -n 6 data/colored-shapes.dkvp | mlr --opprint cat
    color  shape    flag i  u                   v                    w                   x
    yellow triangle 1    11 0.6321695890307647  0.9887207810889004   0.4364983936735774  5.7981881667050565
    red    square   1    15 0.21966833570651523 0.001257332190235938 0.7927778364718627  2.944117399716207
    red    circle   1    16 0.20901671281497636 0.29005231936593445  0.13810280912907674 5.065034003400998
    red    square   0    48 0.9562743938458542  0.7467203085342884   0.7755423050923582  7.117831369597269
    purple triangle 0    51 0.4355354501763202  0.8591292672156728   0.8122903963006748  5.753094629505863
    red    square   0    64 0.2015510269821953  0.9531098083420033   0.7719912015786777  5.612050466474166

Look at uncategorized stats (using `creach <https://github.com/johnkerl/scripts/blob/master/fundam/creach>`_ for spacing).

Here it looks reasonable that ``u`` is unit-uniform; something's up with ``v`` but we can't yet see what:

.. code-block::
   :emphasize-lines: 1,1

    $ mlr --oxtab stats1 -a min,mean,max -f flag,u,v data/colored-shapes.dkvp | creach 3
    flag_min  0
    flag_mean 0.398889
    flag_max  1
    
    u_min     0.000044
    u_mean    0.498326
    u_max     0.999969
    
    v_min     -0.092709
    v_mean    0.497787
    v_max     1.072500

The histogram shows the different distribution of 0/1 flags:

.. code-block::
   :emphasize-lines: 1,1

    $ mlr --opprint histogram -f flag,u,v --lo -0.1 --hi 1.1 --nbins 12 data/colored-shapes.dkvp
    bin_lo    bin_hi   flag_count u_count v_count
    -0.100000 0.000000 6058       0       36
    0.000000  0.100000 0          1062    988
    0.100000  0.200000 0          985     1003
    0.200000  0.300000 0          1024    1014
    0.300000  0.400000 0          1002    991
    0.400000  0.500000 0          989     1041
    0.500000  0.600000 0          1001    1016
    0.600000  0.700000 0          972     962
    0.700000  0.800000 0          1035    1070
    0.800000  0.900000 0          995     993
    0.900000  1.000000 4020       1013    939
    1.000000  1.100000 0          0       25

Look at univariate stats by color and shape. In particular, color-dependent flag probabilities pop out, aligning with their original Bernoulli probablities from the data-generator script:

.. code-block::
   :emphasize-lines: 1,1

    $ mlr --opprint stats1 -a min,mean,max -f flag,u,v -g color then sort -f color data/colored-shapes.dkvp
    color  flag_min flag_mean flag_max u_min    u_mean   u_max    v_min     v_mean   v_max
    blue   0        0.584354  1        0.000044 0.517717 0.999969 0.001489  0.491056 0.999576
    green  0        0.209197  1        0.000488 0.504861 0.999936 0.000501  0.499085 0.999676
    orange 0        0.521452  1        0.001235 0.490532 0.998885 0.002449  0.487764 0.998475
    purple 0        0.090193  1        0.000266 0.494005 0.999647 0.000364  0.497051 0.999975
    red    0        0.303167  1        0.000671 0.492560 0.999882 -0.092709 0.496535 1.072500
    yellow 0        0.892427  1        0.001300 0.497129 0.999923 0.000711  0.510627 0.999919

.. code-block::
   :emphasize-lines: 1,1

    $ mlr --opprint stats1 -a min,mean,max -f flag,u,v -g shape then sort -f shape data/colored-shapes.dkvp
    shape    flag_min flag_mean flag_max u_min    u_mean   u_max    v_min     v_mean   v_max
    circle   0        0.399846  1        0.000044 0.498555 0.999923 -0.092709 0.495524 1.072500
    square   0        0.396112  1        0.000188 0.499385 0.999969 0.000089  0.496538 0.999975
    triangle 0        0.401542  1        0.000881 0.496859 0.999661 0.000717  0.501050 0.999995

Look at bivariate stats by color and shape. In particular, ``u,v`` pairwise correlation for red circles pops out:

.. code-block::
   :emphasize-lines: 1,1

    $ mlr --opprint --right stats2 -a corr -f u,v,w,x data/colored-shapes.dkvp
    u_v_corr  w_x_corr
    0.133418 -0.011320

.. code-block::
   :emphasize-lines: 1,1

    $ mlr --opprint --right stats2 -a corr -f u,v,w,x -g color,shape then sort -nr u_v_corr data/colored-shapes.dkvp
     color    shape  u_v_corr  w_x_corr
       red   circle  0.980798 -0.018565
    orange   square  0.176858 -0.071044
     green   circle  0.057644  0.011795
       red   square  0.055745 -0.000680
    yellow triangle  0.044573  0.024605
    yellow   square  0.043792 -0.044623
    purple   circle  0.035874  0.134112
      blue   square  0.032412 -0.053508
      blue triangle  0.015356 -0.000608
    orange   circle  0.010519 -0.162795
       red triangle  0.008098  0.012486
    purple triangle  0.005155 -0.045058
    purple   square -0.025680  0.057694
     green   square -0.025776 -0.003265
    orange triangle -0.030457 -0.131870
    yellow   circle -0.064773  0.073695
      blue   circle -0.102348 -0.030529
     green triangle -0.109018 -0.048488
