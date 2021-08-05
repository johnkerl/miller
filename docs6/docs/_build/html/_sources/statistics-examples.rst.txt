..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Statistics examples
====================

Computing interquartile ranges
----------------------------------------------------------------

For one or more specified field names, simply compute p25 and p75, then write the IQR as the difference of p75 and p25:

.. code-block:: none
   :emphasize-lines: 1-3

    mlr --oxtab stats1 -f x -a p25,p75 \
        then put '$x_iqr = $x_p75 - $x_p25' \
        data/medium 
    x_p25 0.24667037823231752
    x_p75 0.7481860062358446
    x_iqr 0.5015156280035271

For wildcarded field names, first compute p25 and p75, then loop over field names with ``p25`` in them:

.. code-block:: none
   :emphasize-lines: 1-7

    mlr --oxtab stats1 --fr '[i-z]' -a p25,p75 \
        then put 'for (k,v in $*) {
          if (k =~ "(.*)_p25") {
            $["\1_iqr"] = $["\1_p75"] - $["\1_p25"]
          }
        }' \
        data/medium 

Computing weighted means
----------------------------------------------------------------

This might be more elegantly implemented as an option within the ``stats1`` verb. Meanwhile, it's expressible within the DSL:

.. code-block:: none
   :emphasize-lines: 1-24

    mlr --from data/medium put -q '
      # Using the y field for weighting in this example
      weight = $y;
    
      # Using the a field for weighted aggregation in this example
      @sumwx[$a] += weight * $i;
      @sumw[$a] += weight;
    
      @sumx[$a] += $i;
      @sumn[$a] += 1;
    
      end {
        map wmean = {};
        map mean  = {};
        for (a in @sumwx) {
          wmean[a] = @sumwx[a] / @sumw[a]
        }
        for (a in @sumx) {
          mean[a] = @sumx[a] / @sumn[a]
        }
        #emit wmean, "a";
        #emit mean, "a";
        emit (wmean, mean), "a";
      }'
    a=pan,wmean=4979.563722208067,mean=5028.259010091302
    a=eks,wmean=4890.3815931472145,mean=4956.2900763358775
    a=wye,wmean=4946.987746229947,mean=4920.001017293998
    a=zee,wmean=5164.719684856538,mean=5123.092330239375
    a=hat,wmean=4925.533162478552,mean=4967.743946419371
