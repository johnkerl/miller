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
    x_sum 4986.019681679581

or

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint stats1 -a sum -f x -g b ../data/medium
    b   x_sum
    pan 965.7636699425815
    wye 1023.5484702619565
    zee 979.7420161495838
    eks 1016.7728571314786
    hat 1000.192668193983

rather than the more tedious

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --oxtab put -q '
      @x_sum += $x;
      end {
        emit @x_sum
      }
    ' data/medium
    x_sum 4986.019681679581

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
    pan 965.7636699425815
    wye 1023.5484702619565
    zee 979.7420161495838
    eks 1016.7728571314786
    hat 1000.192668193983

The former (``mlr stats1`` et al.) has the advantages of being easier to type, being less error-prone to type, and running faster.

Nonetheless, out-of-stream variables (which I whimsically call *oosvars*), begin/end blocks, and emit statements give you the ability to implement logic -- if you wish to do so -- which isn't present in other Miller verbs.  (If you find yourself often using the same out-of-stream-variable logic over and over, please file a request at https://github.com/johnkerl/miller/issues to get it implemented directly in Go as a Miller verb of its own.)

The following examples compute some things using oosvars which are already computable using Miller verbs, by way of providing food for thought.

Mean without/with oosvars
----------------------------------------------------------------

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint stats1 -a mean -f x data/medium
    x_mean
    0.49860196816795804

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
    0.49860196816795804

Keyed mean without/with oosvars
----------------------------------------------------------------

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint stats1 -a mean -f x -g a,b data/medium
    a   b   x_mean
    pan pan 0.5133141190437597
    eks pan 0.48507555383425127
    wye wye 0.49150092785839306
    eks wye 0.4838950517724162
    wye pan 0.4996119901034838
    zee pan 0.5198298297816007
    eks zee 0.49546320772681596
    zee wye 0.5142667998230479
    hat wye 0.49381326184632596
    pan wye 0.5023618498923658
    zee eks 0.4883932942792647
    hat zee 0.5099985721987774
    hat eks 0.48587864619953547
    wye hat 0.4977304763723314
    pan eks 0.5036718595143479
    eks eks 0.5227992666570941
    hat hat 0.47993053101017374
    hat pan 0.4643355557376876
    zee zee 0.5127559183726382
    pan hat 0.492140950155604
    pan zee 0.4966041598627583
    zee hat 0.46772617655014515
    wye zee 0.5059066170573692
    eks hat 0.5006790659966355
    wye eks 0.5306035254809106

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
    pan pan 0.5133141190437597
    pan wye 0.5023618498923658
    pan eks 0.5036718595143479
    pan hat 0.492140950155604
    pan zee 0.4966041598627583
    eks pan 0.48507555383425127
    eks wye 0.4838950517724162
    eks zee 0.49546320772681596
    eks eks 0.5227992666570941
    eks hat 0.5006790659966355
    wye wye 0.49150092785839306
    wye pan 0.4996119901034838
    wye hat 0.4977304763723314
    wye zee 0.5059066170573692
    wye eks 0.5306035254809106
    zee pan 0.5198298297816007
    zee wye 0.5142667998230479
    zee eks 0.4883932942792647
    zee zee 0.5127559183726382
    zee hat 0.46772617655014515
    hat wye 0.49381326184632596
    hat zee 0.5099985721987774
    hat eks 0.48587864619953547
    hat hat 0.47993053101017374
    hat pan 0.4643355557376876

Variance and standard deviation without/with oosvars
----------------------------------------------------------------

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --oxtab stats1 -a count,sum,mean,var,stddev -f x data/medium
    x_count  10000
    x_sum    4986.019681679581
    x_mean   0.49860196816795804
    x_var    0.08426974433144456
    x_stddev 0.2902925151144007

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
    sumx   4986.019681679581
    sumx2  3328.652400179729
    mean   0.49860196816795804
    var    0.08426974433144456
    stddev 0.2902925151144007

You can also do this keyed, of course, imitating the keyed-mean example above.

Min/max without/with oosvars
----------------------------------------------------------------

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --oxtab stats1 -a min,max -f x data/medium
    x_min 4.509679127584487e-05
    x_max 0.999952670371898

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --oxtab put -q '@x_min = min(@x_min, $x); @x_max = max(@x_max, $x); end{emitf @x_min, @x_max}' data/medium
    x_min 4.509679127584487e-05
    x_max 0.999952670371898

Keyed min/max without/with oosvars
----------------------------------------------------------------

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint stats1 -a min,max -f x -g a data/medium
    a   x_min                  x_max
    pan 0.00020390740306253097 0.9994029107062516
    eks 0.0006917972627396018  0.9988110946859143
    wye 0.0001874794831505655  0.9998228522652893
    zee 0.0005486114815762555  0.9994904324789629
    hat 4.509679127584487e-05  0.999952670371898

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint --from data/medium put -q '
      @min[$a] = min(@min[$a], $x);
      @max[$a] = max(@max[$a], $x);
      end{
        emit (@min, @max), "a";
      }
    '
    a   min                    max
    pan 0.00020390740306253097 0.9994029107062516
    eks 0.0006917972627396018  0.9988110946859143
    wye 0.0001874794831505655  0.9998228522652893
    zee 0.0005486114815762555  0.9994904324789629
    hat 4.509679127584487e-05  0.999952670371898

Delta without/with oosvars
----------------------------------------------------------------

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint step -a delta -f x data/small
    a   b   i x                   y                   x_delta
    pan pan 1 0.3467901443380824  0.7268028627434533  0
    eks pan 2 0.7586799647899636  0.5221511083334797  0.41188982045188116
    wye wye 3 0.20460330576630303 0.33831852551664776 -0.5540766590236605
    eks wye 4 0.38139939387114097 0.13418874328430463 0.17679608810483793
    wye pan 5 0.5732889198020006  0.8636244699032729  0.19188952593085962

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint put '$x_delta = is_present(@last) ? $x - @last : 0; @last = $x' data/small
    a   b   i x                   y                   x_delta
    pan pan 1 0.3467901443380824  0.7268028627434533  0
    eks pan 2 0.7586799647899636  0.5221511083334797  0.41188982045188116
    wye wye 3 0.20460330576630303 0.33831852551664776 -0.5540766590236605
    eks wye 4 0.38139939387114097 0.13418874328430463 0.17679608810483793
    wye pan 5 0.5732889198020006  0.8636244699032729  0.19188952593085962

Keyed delta without/with oosvars
----------------------------------------------------------------

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint step -a delta -f x -g a data/small
    a   b   i x                   y                   x_delta
    pan pan 1 0.3467901443380824  0.7268028627434533  0
    eks pan 2 0.7586799647899636  0.5221511083334797  0
    wye wye 3 0.20460330576630303 0.33831852551664776 0
    eks wye 4 0.38139939387114097 0.13418874328430463 -0.3772805709188226
    wye pan 5 0.5732889198020006  0.8636244699032729  0.36868561403569755

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint put '$x_delta = is_present(@last[$a]) ? $x - @last[$a] : 0; @last[$a]=$x' data/small
    a   b   i x                   y                   x_delta
    pan pan 1 0.3467901443380824  0.7268028627434533  0
    eks pan 2 0.7586799647899636  0.5221511083334797  0
    wye wye 3 0.20460330576630303 0.33831852551664776 0
    eks wye 4 0.38139939387114097 0.13418874328430463 -0.3772805709188226
    wye pan 5 0.5732889198020006  0.8636244699032729  0.36868561403569755

Exponentially weighted moving averages without/with oosvars
----------------------------------------------------------------

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint step -a ewma -d 0.1 -f x data/small
    a   b   i x                   y                   x_ewma_0.1
    pan pan 1 0.3467901443380824  0.7268028627434533  0.3467901443380824
    eks pan 2 0.7586799647899636  0.5221511083334797  0.3879791263832706
    wye wye 3 0.20460330576630303 0.33831852551664776 0.36964154432157387
    eks wye 4 0.38139939387114097 0.13418874328430463 0.37081732927653055
    wye pan 5 0.5732889198020006  0.8636244699032729  0.3910644883290776

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint put '
      begin{ @a=0.1 };
      $e = NR==1 ? $x : @a * $x + (1 - @a) * @e;
      @e=$e
    ' data/small
    a   b   i x                   y                   e
    pan pan 1 0.3467901443380824  0.7268028627434533  0.3467901443380824
    eks pan 2 0.7586799647899636  0.5221511083334797  0.3879791263832706
    wye wye 3 0.20460330576630303 0.33831852551664776 0.36964154432157387
    eks wye 4 0.38139939387114097 0.13418874328430463 0.37081732927653055
    wye pan 5 0.5732889198020006  0.8636244699032729  0.3910644883290776
