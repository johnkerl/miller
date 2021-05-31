..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Cookbook part 3: Stats with and without out-of-stream variables
================================================================

Overview
----------------------------------------------------------------

One of Miller's strengths is its compact notation: for example, given input of the form

.. code-block:: none
   :emphasize-lines: 1,1

    $ head -n 5 ../data/medium
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729

you can simply do

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --oxtab stats1 -a sum -f x ../data/medium
    x_sum 4986.019682

or

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint stats1 -a sum -f x -g b ../data/medium
    b   x_sum
    pan 965.763670
    wye 1023.548470
    zee 979.742016
    eks 1016.772857
    hat 1000.192668

rather than the more tedious

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --oxtab put -q '
      @x_sum += $x;
      end {
        emit @x_sum
      }
    ' data/medium
    x_sum 4986.019682

or

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint put -q '
      @x_sum[$b] += $x;
      end {
        emit @x_sum, "b"
      }
    ' data/medium
    b   x_sum
    pan 965.763670
    wye 1023.548470
    zee 979.742016
    eks 1016.772857
    hat 1000.192668

The former (``mlr stats1`` et al.) has the advantages of being easier to type, being less error-prone to type, and running faster.

Nonetheless, out-of-stream variables (which I whimsically call *oosvars*), begin/end blocks, and emit statements give you the ability to implement logic -- if you wish to do so -- which isn't present in other Miller verbs.  (If you find yourself often using the same out-of-stream-variable logic over and over, please file a request at https://github.com/johnkerl/miller/issues to get it implemented directly in C as a Miller verb of its own.)

The following examples compute some things using oosvars which are already computable using Miller verbs, by way of providing food for thought.

Mean without/with oosvars
----------------------------------------------------------------

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint stats1 -a mean -f x data/medium
    x_mean
    0.498602

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint put -q '
      @x_sum += $x;
      @x_count += 1;
      end {
        @x_mean = @x_sum / @x_count;
        emit @x_mean
      }
    ' data/medium
    x_mean
    0.498602

Keyed mean without/with oosvars
----------------------------------------------------------------

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint stats1 -a mean -f x -g a,b data/medium
    a   b   x_mean
    pan pan 0.513314
    eks pan 0.485076
    wye wye 0.491501
    eks wye 0.483895
    wye pan 0.499612
    zee pan 0.519830
    eks zee 0.495463
    zee wye 0.514267
    hat wye 0.493813
    pan wye 0.502362
    zee eks 0.488393
    hat zee 0.509999
    hat eks 0.485879
    wye hat 0.497730
    pan eks 0.503672
    eks eks 0.522799
    hat hat 0.479931
    hat pan 0.464336
    zee zee 0.512756
    pan hat 0.492141
    pan zee 0.496604
    zee hat 0.467726
    wye zee 0.505907
    eks hat 0.500679
    wye eks 0.530604

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint put -q '
      @x_sum[$a][$b] += $x;
      @x_count[$a][$b] += 1;
      end{
        for ((a, b), v in @x_sum) {
          @x_mean[a][b] = @x_sum[a][b] / @x_count[a][b];
        }
        emit @x_mean, "a", "b"
      }
    ' data/medium
    a   b   x_mean
    pan pan 0.513314
    pan wye 0.502362
    pan eks 0.503672
    pan hat 0.492141
    pan zee 0.496604
    eks pan 0.485076
    eks wye 0.483895
    eks zee 0.495463
    eks eks 0.522799
    eks hat 0.500679
    wye wye 0.491501
    wye pan 0.499612
    wye hat 0.497730
    wye zee 0.505907
    wye eks 0.530604
    zee pan 0.519830
    zee wye 0.514267
    zee eks 0.488393
    zee zee 0.512756
    zee hat 0.467726
    hat wye 0.493813
    hat zee 0.509999
    hat eks 0.485879
    hat hat 0.479931
    hat pan 0.464336

Variance and standard deviation without/with oosvars
----------------------------------------------------------------

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --oxtab stats1 -a count,sum,mean,var,stddev -f x data/medium
    x_count  10000
    x_sum    4986.019682
    x_mean   0.498602
    x_var    0.084270
    x_stddev 0.290293

.. code-block:: none
   :emphasize-lines: 1,1

    $ cat variance.mlr
    @n += 1;
    @sumx += $x;
    @sumx2 += $x**2;
    end {
      @mean = @sumx / @n;
      @var = (@sumx2 - @mean * (2 * @sumx - @n * @mean)) / (@n - 1);
      @stddev = sqrt(@var);
      emitf @n, @sumx, @sumx2, @mean, @var, @stddev
    }

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --oxtab put -q -f variance.mlr data/medium
    n      10000
    sumx   4986.019682
    sumx2  3328.652400
    mean   0.498602
    var    0.084270
    stddev 0.290293

You can also do this keyed, of course, imitating the keyed-mean example above.

Min/max without/with oosvars
----------------------------------------------------------------

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --oxtab stats1 -a min,max -f x data/medium
    x_min 0.000045
    x_max 0.999953

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --oxtab put -q '@x_min = min(@x_min, $x); @x_max = max(@x_max, $x); end{emitf @x_min, @x_max}' data/medium
    x_min 0.000045
    x_max 0.999953

Keyed min/max without/with oosvars
----------------------------------------------------------------

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint stats1 -a min,max -f x -g a data/medium
    a   x_min    x_max
    pan 0.000204 0.999403
    eks 0.000692 0.998811
    wye 0.000187 0.999823
    zee 0.000549 0.999490
    hat 0.000045 0.999953

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint --from data/medium put -q '
      @min[$a] = min(@min[$a], $x);
      @max[$a] = max(@max[$a], $x);
      end{
        emit (@min, @max), "a";
      }
    '
    a   min      max
    pan 0.000204 0.999403
    eks 0.000692 0.998811
    wye 0.000187 0.999823
    zee 0.000549 0.999490
    hat 0.000045 0.999953

Delta without/with oosvars
----------------------------------------------------------------

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint step -a delta -f x data/small
    a   b   i x                   y                   x_delta
    pan pan 1 0.3467901443380824  0.7268028627434533  0
    eks pan 2 0.7586799647899636  0.5221511083334797  0.411890
    wye wye 3 0.20460330576630303 0.33831852551664776 -0.554077
    eks wye 4 0.38139939387114097 0.13418874328430463 0.176796
    wye pan 5 0.5732889198020006  0.8636244699032729  0.191890

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint put '$x_delta = is_present(@last) ? $x - @last : 0; @last = $x' data/small
    a   b   i x                   y                   x_delta
    pan pan 1 0.3467901443380824  0.7268028627434533  0
    eks pan 2 0.7586799647899636  0.5221511083334797  0.411890
    wye wye 3 0.20460330576630303 0.33831852551664776 -0.554077
    eks wye 4 0.38139939387114097 0.13418874328430463 0.176796
    wye pan 5 0.5732889198020006  0.8636244699032729  0.191890

Keyed delta without/with oosvars
----------------------------------------------------------------

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint step -a delta -f x -g a data/small
    a   b   i x                   y                   x_delta
    pan pan 1 0.3467901443380824  0.7268028627434533  0
    eks pan 2 0.7586799647899636  0.5221511083334797  0
    wye wye 3 0.20460330576630303 0.33831852551664776 0
    eks wye 4 0.38139939387114097 0.13418874328430463 -0.377281
    wye pan 5 0.5732889198020006  0.8636244699032729  0.368686

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint put '$x_delta = is_present(@last[$a]) ? $x - @last[$a] : 0; @last[$a]=$x' data/small
    a   b   i x                   y                   x_delta
    pan pan 1 0.3467901443380824  0.7268028627434533  0
    eks pan 2 0.7586799647899636  0.5221511083334797  0
    wye wye 3 0.20460330576630303 0.33831852551664776 0
    eks wye 4 0.38139939387114097 0.13418874328430463 -0.377281
    wye pan 5 0.5732889198020006  0.8636244699032729  0.368686

Exponentially weighted moving averages without/with oosvars
----------------------------------------------------------------

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint step -a ewma -d 0.1 -f x data/small
    a   b   i x                   y                   x_ewma_0.1
    pan pan 1 0.3467901443380824  0.7268028627434533  0.346790
    eks pan 2 0.7586799647899636  0.5221511083334797  0.387979
    wye wye 3 0.20460330576630303 0.33831852551664776 0.369642
    eks wye 4 0.38139939387114097 0.13418874328430463 0.370817
    wye pan 5 0.5732889198020006  0.8636244699032729  0.391064

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint put '
      begin{ @a=0.1 };
      $e = NR==1 ? $x : @a * $x + (1 - @a) * @e;
      @e=$e
    ' data/small
    a   b   i x                   y                   e
    pan pan 1 0.3467901443380824  0.7268028627434533  0.346790
    eks pan 2 0.7586799647899636  0.5221511083334797  0.387979
    wye wye 3 0.20460330576630303 0.33831852551664776 0.369642
    eks wye 4 0.38139939387114097 0.13418874328430463 0.370817
    wye pan 5 0.5732889198020006  0.8636244699032729  0.391064
