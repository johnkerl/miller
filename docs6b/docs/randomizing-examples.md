<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# Randomizing examples

## Generating random numbers from various distributions

Here we can chain together a few simple building blocks:

<pre>
<b>cat expo-sample.sh</b>
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
</pre>

Namely:

* Set the Miller random-number seed so this webdoc looks the same every time I regenerate it.
* Use pretty-printed tabular output.
* Use pretty-printed tabular output.
* Use ``seqgen`` to produce 100,000 records ``i=0``, ``i=1``, etc.
* Send those to a ``put`` step which defines an inverse-transform-sampling function and calls it twice, then computes the sum and product of samples.
* Send those to a histogram, and from there to a bar-plotter. This is just for visualization; you could just as well output CSV and send that off to your own plotting tool, etc.

The output is as follows:

<pre>
<b>sh expo-sample.sh</b>
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
</pre>

## Randomly selecting words from a list

Given this `word list <./data/english-words.txt>`_, first take a look to see what the first few lines look like:

<pre>
<b>head data/english-words.txt</b>
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
</pre>

Then the following will randomly sample ten words with four to eight characters in them:

<pre>
<b>mlr --from data/english-words.txt --nidx filter -S 'n=strlen($1);4<=n&&n<=8' then sample -k 10</b>
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
</pre>

## Randomly generating jabberwocky words

These are simple *n*-grams as `described here <http://johnkerl.org/randspell/randspell-slides-ts.pdf>`_. Some common functions are `located here <https://github.com/johnkerl/miller/blob/master/docs/ngrams/ngfuncs.mlr.txt>`_. Then here are scripts for `1-grams <https://github.com/johnkerl/miller/blob/master/docs/ngrams/ng1.mlr.txt>`_ `2-grams <https://github.com/johnkerl/miller/blob/master/docs/ngrams/ng2.mlr.txt>`_ `3-grams <https://github.com/johnkerl/miller/blob/master/docs/ngrams/ng3.mlr.txt>`_ `4-grams <https://github.com/johnkerl/miller/blob/master/docs/ngrams/ng4.mlr.txt>`_, and `5-grams <https://github.com/johnkerl/miller/blob/master/docs/ngrams/ng5.mlr.txt>`_.

The idea is that words from the input file are consumed, then taken apart and pasted back together in ways which imitate the letter-to-letter transitions found in the word list -- giving us automatically generated words in the same vein as *bromance* and *spork*:

<pre>
<b>mlr --nidx --from ./ngrams/gsl-2000.txt put -q -f ./ngrams/ngfuncs.mlr -f ./ngrams/ng5.mlr</b>
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
</pre>
