..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Cookbook part 2: Random things, and some math
================================================================

Randomly selecting words from a list
----------------------------------------------------------------

Given this `word list <./data/english-words.txt>`_, first take a look to see what the first few lines look like:

.. code-block:: none
   :emphasize-lines: 1-1

    head data/english-words.txt
    a
    aa
    aal
    aalii
    aam
    aardvark
    aardwolf
    aba
    abac
    abaca

Then the following will randomly sample ten words with four to eight characters in them:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --from data/english-words.txt --nidx filter -S 'n=strlen($1);4<=n&&n<=8' then sample -k 10
    thionine
    birchman
    mildewy
    avigate
    addedly
    abaze
    askant
    aiming
    insulant
    coinmate

Randomly generating jabberwocky words
----------------------------------------------------------------

These are simple *n*-grams as `described here <http://johnkerl.org/randspell/randspell-slides-ts.pdf>`_. Some common functions are `located here <https://github.com/johnkerl/miller/blob/master/docs/ngrams/ngfuncs.mlr.txt>`_. Then here are scripts for `1-grams <https://github.com/johnkerl/miller/blob/master/docs/ngrams/ng1.mlr.txt>`_ `2-grams <https://github.com/johnkerl/miller/blob/master/docs/ngrams/ng2.mlr.txt>`_ `3-grams <https://github.com/johnkerl/miller/blob/master/docs/ngrams/ng3.mlr.txt>`_ `4-grams <https://github.com/johnkerl/miller/blob/master/docs/ngrams/ng4.mlr.txt>`_, and `5-grams <https://github.com/johnkerl/miller/blob/master/docs/ngrams/ng5.mlr.txt>`_.

The idea is that words from the input file are consumed, then taken apart and pasted back together in ways which imitate the letter-to-letter transitions found in the word list -- giving us automatically generated words in the same vein as *bromance* and *spork*:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --nidx --from ./ngrams/gsl-2000.txt put -q -f ./ngrams/ngfuncs.mlr -f ./ngrams/ng5.mlr
    beard
    plastinguish
    politicially
    noise
    loan
    country
    controductionary
    suppery
    lose
    lessors
    dollar
    judge
    rottendence
    lessenger
    diffendant
    suggestional

Program timing
----------------------------------------------------------------

This admittedly artificial example demonstrates using Miller time and stats functions to introspectively acquire some information about Miller's own runtime. The ``delta`` function computes the difference between successive timestamps.

.. code-block:: none

    $ ruby -e '10000.times{|i|puts "i=#{i+1}"}' > lines.txt
    
    $ head -n 5 lines.txt
    i=1
    i=2
    i=3
    i=4
    i=5
    
    mlr --ofmt '%.9le' --opprint put '$t=systime()' then step -a delta -f t lines.txt | head -n 7
    i     t                 t_delta
    1     1430603027.018016 1.430603027e+09
    2     1430603027.018043 2.694129944e-05
    3     1430603027.018048 5.006790161e-06
    4     1430603027.018052 4.053115845e-06
    5     1430603027.018055 2.861022949e-06
    6     1430603027.018058 3.099441528e-06
    
    mlr --ofmt '%.9le' --oxtab \
      put '$t=systime()' then \
      step -a delta -f t then \
      filter '$i>1' then \
      stats1 -a min,mean,max -f t_delta \
      lines.txt
    t_delta_min  2.861022949e-06
    t_delta_mean 4.077508505e-06
    t_delta_max  5.388259888e-05

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

Generating random numbers from various distributions
----------------------------------------------------------------

Here we can chain together a few simple building blocks:

.. code-block:: none
   :emphasize-lines: 1-1

    cat expo-sample.sh
    # Generate 100,000 pairs of independent and identically distributed
    # exponentially distributed random variables with the same rate parameter
    # (namely, 2.5). Then compute histograms of one of them, along with
    # histograms for their sum and their product.
    #
    # See also https://en.wikipedia.org/wiki/Exponential_distribution
    #
    # Here I'm using a specified random-number seed so this example always
    # produces the same output for this web document: in everyday practice we
    # wouldn't do that.
    
    mlr -n \
      --seed 0 \
      --opprint \
      seqgen --stop 100000 \
      then put '
        # https://en.wikipedia.org/wiki/Inverse_transform_sampling
        func expo_sample(lambda) {
          return -log(1-urand())/lambda
        }
        $u = expo_sample(2.5);
        $v = expo_sample(2.5);
        $s = $u + $v;
        $p = $u * $v;
      ' \
      then histogram -f u,s,p --lo 0 --hi 2 --nbins 50 \
      then bar -f u_count,s_count,p_count --auto -w 20

Namely:

* Set the Miller random-number seed so this webdoc looks the same every time I regenerate it.
* Use pretty-printed tabular output.
* Use pretty-printed tabular output.
* Use ``seqgen`` to produce 100,000 records ``i=0``, ``i=1``, etc.
* Send those to a ``put`` step which defines an inverse-transform-sampling function and calls it twice, then computes the sum and product of samples.
* Send those to a histogram, and from there to a bar-plotter. This is just for visualization; you could just as well output CSV and send that off to your own plotting tool, etc.

The output is as follows:

.. code-block:: none
   :emphasize-lines: 1-1

    sh expo-sample.sh
    bin_lo bin_hi u_count                        s_count                         p_count
    0      0.04   [64]*******************#[9554] [326]#...................[3703] [19]*******************#[39809]
    0.04   0.08   [64]*****************...[9554] [326]*****...............[3703] [19]*******.............[39809]
    0.08   0.12   [64]****************....[9554] [326]*********...........[3703] [19]****................[39809]
    0.12   0.16   [64]**************......[9554] [326]************........[3703] [19]***.................[39809]
    0.16   0.2    [64]*************.......[9554] [326]**************......[3703] [19]**..................[39809]
    0.2    0.24   [64]************........[9554] [326]*****************...[3703] [19]*...................[39809]
    0.24   0.28   [64]**********..........[9554] [326]******************..[3703] [19]*...................[39809]
    0.28   0.32   [64]*********...........[9554] [326]******************..[3703] [19]*...................[39809]
    0.32   0.36   [64]********............[9554] [326]*******************.[3703] [19]#...................[39809]
    0.36   0.4    [64]*******.............[9554] [326]*******************#[3703] [19]#...................[39809]
    0.4    0.44   [64]*******.............[9554] [326]*******************.[3703] [19]#...................[39809]
    0.44   0.48   [64]******..............[9554] [326]*******************.[3703] [19]#...................[39809]
    0.48   0.52   [64]*****...............[9554] [326]******************..[3703] [19]#...................[39809]
    0.52   0.56   [64]*****...............[9554] [326]******************..[3703] [19]#...................[39809]
    0.56   0.6    [64]****................[9554] [326]*****************...[3703] [19]#...................[39809]
    0.6    0.64   [64]****................[9554] [326]******************..[3703] [19]#...................[39809]
    0.64   0.68   [64]***.................[9554] [326]****************....[3703] [19]#...................[39809]
    0.68   0.72   [64]***.................[9554] [326]****************....[3703] [19]#...................[39809]
    0.72   0.76   [64]***.................[9554] [326]***************.....[3703] [19]#...................[39809]
    0.76   0.8    [64]**..................[9554] [326]**************......[3703] [19]#...................[39809]
    0.8    0.84   [64]**..................[9554] [326]*************.......[3703] [19]#...................[39809]
    0.84   0.88   [64]**..................[9554] [326]************........[3703] [19]#...................[39809]
    0.88   0.92   [64]**..................[9554] [326]************........[3703] [19]#...................[39809]
    0.92   0.96   [64]*...................[9554] [326]***********.........[3703] [19]#...................[39809]
    0.96   1      [64]*...................[9554] [326]**********..........[3703] [19]#...................[39809]
    1      1.04   [64]*...................[9554] [326]*********...........[3703] [19]#...................[39809]
    1.04   1.08   [64]*...................[9554] [326]********............[3703] [19]#...................[39809]
    1.08   1.12   [64]*...................[9554] [326]********............[3703] [19]#...................[39809]
    1.12   1.16   [64]*...................[9554] [326]********............[3703] [19]#...................[39809]
    1.16   1.2    [64]*...................[9554] [326]*******.............[3703] [19]#...................[39809]
    1.2    1.24   [64]#...................[9554] [326]******..............[3703] [19]#...................[39809]
    1.24   1.28   [64]#...................[9554] [326]*****...............[3703] [19]#...................[39809]
    1.28   1.32   [64]#...................[9554] [326]*****...............[3703] [19]#...................[39809]
    1.32   1.36   [64]#...................[9554] [326]****................[3703] [19]#...................[39809]
    1.36   1.4    [64]#...................[9554] [326]****................[3703] [19]#...................[39809]
    1.4    1.44   [64]#...................[9554] [326]****................[3703] [19]#...................[39809]
    1.44   1.48   [64]#...................[9554] [326]***.................[3703] [19]#...................[39809]
    1.48   1.52   [64]#...................[9554] [326]***.................[3703] [19]#...................[39809]
    1.52   1.56   [64]#...................[9554] [326]***.................[3703] [19]#...................[39809]
    1.56   1.6    [64]#...................[9554] [326]**..................[3703] [19]#...................[39809]
    1.6    1.64   [64]#...................[9554] [326]**..................[3703] [19]#...................[39809]
    1.64   1.68   [64]#...................[9554] [326]**..................[3703] [19]#...................[39809]
    1.68   1.72   [64]#...................[9554] [326]*...................[3703] [19]#...................[39809]
    1.72   1.76   [64]#...................[9554] [326]*...................[3703] [19]#...................[39809]
    1.76   1.8    [64]#...................[9554] [326]*...................[3703] [19]#...................[39809]
    1.8    1.84   [64]#...................[9554] [326]#...................[3703] [19]#...................[39809]
    1.84   1.88   [64]#...................[9554] [326]#...................[3703] [19]#...................[39809]
    1.88   1.92   [64]#...................[9554] [326]#...................[3703] [19]#...................[39809]
    1.92   1.96   [64]#...................[9554] [326]#...................[3703] [19]#...................[39809]
    1.96   2      [64]#...................[9554] [326]#...................[3703] [19]#...................[39809]

Sieve of Eratosthenes
----------------------------------------------------------------

The `Sieve of Eratosthenes <http://en.wikipedia.org/wiki/Sieve_of_Eratosthenes>`_ is a standard introductory programming topic. The idea is to find all primes up to some *N* by making a list of the numbers 1 to *N*, then striking out all multiples of 2 except 2 itself, all multiples of 3 except 3 itself, all multiples of 4 except 4 itself, and so on. Whatever survives that without getting marked is a prime. This is easy enough in Miller. Notice that here all the work is in ``begin`` and ``end`` statements; there is no file input (so we use ``mlr -n`` to keep Miller from waiting for input data).

.. code-block:: none
   :emphasize-lines: 1-1

    cat programs/sieve.mlr
    # ================================================================
    # Sieve of Eratosthenes: simple example of Miller DSL as programming language.
    # ================================================================
    
    # Put this in a begin-block so we can do either
    #   mlr -n put -q -f name-of-this-file.mlr
    # or
    #   mlr -n put -q -f name-of-this-file.mlr -e '@n = 200'
    # i.e. 100 is the default upper limit, and another can be specified using -e.
    begin {
      @n = 100;
    }
    
    end {
      for (int i = 0; i <= @n; i += 1) {
        @s[i] = true;
      }
      @s[0] = false; # 0 is neither prime nor composite
      @s[1] = false; # 1 is neither prime nor composite
      # Strike out multiples
      for (int i = 2; i <= @n; i += 1) {
        for (int j = i+i; j <= @n; j += i) {
          @s[j] = false;
        }
      }
      # Print survivors
      for (int i = 0; i <= @n; i += 1) {
        if (@s[i]) {
          print i;
        }
      }
    }

.. code-block:: none
   :emphasize-lines: 1-1

    mlr -n put -f programs/sieve.mlr
    2
    3
    5
    7
    11
    13
    17
    19
    23
    29
    31
    37
    41
    43
    47
    53
    59
    61
    67
    71
    73
    79
    83
    89
    97

Mandelbrot-set generator
----------------------------------------------------------------

The `Mandelbrot set <http://en.wikipedia.org/wiki/Mandelbrot_set>`_ is also easily expressed. This isn't an important case of data-processing in the vein for which Miller was designed, but it is an example of Miller as a general-purpose programming language -- a test case for the expressiveness of the language.

The (approximate) computation of points in the complex plane which are and aren't members is just a few lines of complex arithmetic (see the Wikipedia article); how to render them is another task.  Using graphics libraries you can create PNG or JPEG files, but another fun way to do this is by printing various characters to the screen:

.. code-block:: none
   :emphasize-lines: 1-1

    cat programs/mand.mlr
    # Mandelbrot set generator: simple example of Miller DSL as programming language.
    begin {
      # Set defaults
      @rcorn     = -2.0;
      @icorn     = -2.0;
      @side      = 4.0;
      @iheight   = 50;
      @iwidth    = 100;
      @maxits    = 100;
      @levelstep = 5;
      @chars     = "@X*o-."; # Palette of characters to print to the screen.
      @verbose   = false;
      @do_julia  = false;
      @jr        = 0.0;      # Real part of Julia point, if any
      @ji        = 0.0;      # Imaginary part of Julia point, if any
    }
    
    # Here, we can override defaults from an input file (if any).  In Miller's
    # put/filter DSL, absent-null right-hand sides result in no assignment so we
    # can simply put @rcorn = $rcorn: if there is a field in the input like
    # 'rcorn = -1.847' we'll read and use it, else we'll keep the default.
    @rcorn     = $rcorn;
    @icorn     = $icorn;
    @side      = $side;
    @iheight   = $iheight;
    @iwidth    = $iwidth;
    @maxits    = $maxits;
    @levelstep = $levelstep;
    @chars     = $chars;
    @verbose   = $verbose;
    @do_julia  = $do_julia;
    @jr        = $jr;
    @ji        = $ji;
    
    end {
      if (@verbose) {
        print "RCORN     = ".@rcorn;
        print "ICORN     = ".@icorn;
        print "SIDE      = ".@side;
        print "IHEIGHT   = ".@iheight;
        print "IWIDTH    = ".@iwidth;
        print "MAXITS    = ".@maxits;
        print "LEVELSTEP = ".@levelstep;
        print "CHARS     = ".@chars;
      }
    
      # Iterate over a matrix of rows and columns, printing one character for each cell.
      for (int ii = @iheight-1; ii >= 0; ii -= 1) {
        num pi = @icorn + (ii/@iheight) * @side;
        for (int ir = 0; ir < @iwidth; ir += 1) {
          num pr = @rcorn + (ir/@iwidth) * @side;
          printn get_point_plot(pr, pi, @maxits, @do_julia, @jr, @ji);
        }
        print;
      }
    }
    
    # This is a function to approximate membership in the Mandelbrot set (or Julia
    # set for a given Julia point if do_julia == true) for a given point in the
    # complex plane.
    func get_point_plot(pr, pi, maxits, do_julia, jr, ji) {
      num zr = 0.0;
      num zi = 0.0;
      num cr = 0.0;
      num ci = 0.0;
    
      if (!do_julia) {
        zr = 0.0;
        zi = 0.0;
        cr = pr;
        ci = pi;
      } else {
        zr = pr;
        zi = pi;
        cr = jr;
        ci = ji;
      }
    
      int iti = 0;
      bool escaped = false;
      num zt = 0;
      for (iti = 0; iti < maxits; iti += 1) {
        num mag = zr*zr + zi+zi;
        if (mag > 4.0) {
            escaped = true;
            break;
        }
        # z := z^2 + c
        zt = zr*zr - zi*zi + cr;
        zi = 2*zr*zi + ci;
        zr = zt;
      }
      if (!escaped) {
        return ".";
      } else {
        # The // operator is Miller's (pythonic) integer-division operator
        int level = (iti // @levelstep) % strlen(@chars);
        return substr(@chars, level, level);
      }
    }

At standard resolution this makes a nice little ASCII plot:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr -n put -f ./programs/mand.mlr
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXX.XXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXXXXooXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXX**o..*XXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXXX*-....-oXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXX@XXXXXXXXXX*......o*XXXXXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@XXXXXXXXXX**oo*-.-........oo.XXXXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@XXXXXXXXXXXXX....................X..o-XXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@XXXXXXXXXXXXXXX*oo......................oXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@XXX*XXXXXXXXXXXX**o........................*X*X@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@XXXXXXooo***o*.*XX**X..........................o-XX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@XXXXXXXX*-.......-***.............................oXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@XXXXXXXX*@..........Xo............................*XX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@XXXX@XXXXXXXX*o@oX...........@...........................oXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    .........................................................o*XXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@XXXXXXXXX*-.oX...........@...........................oXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@XXXXXXXXXX**@..........*o............................*XXXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@XXXXXXXXXXXXX-........***.............................oXXXXXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@XXXXXXXXXXXXoo****o*.XX***@..........................o-XXXXXXXXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@XXXXX*XXXX*XXXXXXX**-........................***XXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@XXXXXXXXXXXXX*o*.....................@o*XXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@XXXXXXXXXXXX*....................*..o-XX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXX*ooo*-.o........oo.X*XXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXXXXXX**@.....*XXXXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXXXXXX*o....-o*XXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXXXXXXXo*o..*XXXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXXXXXXXXXX*o*XXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXXXXXXXXX@XXXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXXXXXX@@XXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

But using a very small font size (as small as my Mac will let me go), and by choosing the coordinates to zoom in on a particular part of the complex plane, we can get a nice little picture:

.. code-block:: none

    #!/bin/bash
    # Get the number of rows and columns from the terminal window dimensions
    iheight=$(stty size | mlr --nidx --fs space cut -f 1)
    iwidth=$(stty size | mlr --nidx --fs space cut -f 2)
    echo "rcorn=-1.755350,icorn=+0.014230,side=0.000020,maxits=10000,iheight=$iheight,iwidth=$iwidth" \
    | mlr put -f programs/mand.mlr

.. image:: pix/mand.png
