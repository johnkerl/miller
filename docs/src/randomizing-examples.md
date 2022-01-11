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
# Randomizing examples

## Generating random numbers from various distributions

Here we can chain together a few simple building blocks:

<pre class="pre-highlight-in-pair">
<b>cat expo-sample.sh</b>
</pre>
<pre class="pre-non-highlight-in-pair">
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
  ' \
  then histogram -f u,s --lo 0 --hi 2 --nbins 50 \
  then bar -f u_count,s_count --auto -w 20
</pre>

Namely:

* Set the Miller random-number seed so this webdoc looks the same every time I regenerate it.
* Use pretty-printed tabular output.
* Use `seqgen` to produce 100,000 records `i=0`, `i=1`, etc.
* Send those to a `put` step which defines an inverse-transform-sampling function and calls it twice, then computes the sum and product of samples.
* Send those to a histogram, and from there to a bar-plotter. This is just for visualization; you could just as well output CSV and send that off to your own plotting tool, etc.

The output is as follows:

<pre class="pre-highlight-in-pair">
<b>sh expo-sample.sh</b>
</pre>
<pre class="pre-non-highlight-in-pair">
bin_lo bin_hi u_count                        s_count
0      0.04   [64]*******************#[9554] [326]#...................[3703]
0.04   0.08   [64]*****************...[9554] [326]*****...............[3703]
0.08   0.12   [64]****************....[9554] [326]*********...........[3703]
0.12   0.16   [64]**************......[9554] [326]************........[3703]
0.16   0.2    [64]*************.......[9554] [326]**************......[3703]
0.2    0.24   [64]************........[9554] [326]*****************...[3703]
0.24   0.28   [64]**********..........[9554] [326]******************..[3703]
0.28   0.32   [64]*********...........[9554] [326]******************..[3703]
0.32   0.36   [64]********............[9554] [326]*******************.[3703]
0.36   0.4    [64]*******.............[9554] [326]*******************#[3703]
0.4    0.44   [64]*******.............[9554] [326]*******************.[3703]
0.44   0.48   [64]******..............[9554] [326]*******************.[3703]
0.48   0.52   [64]*****...............[9554] [326]******************..[3703]
0.52   0.56   [64]*****...............[9554] [326]******************..[3703]
0.56   0.6    [64]****................[9554] [326]*****************...[3703]
0.6    0.64   [64]****................[9554] [326]******************..[3703]
0.64   0.68   [64]***.................[9554] [326]****************....[3703]
0.68   0.72   [64]***.................[9554] [326]****************....[3703]
0.72   0.76   [64]***.................[9554] [326]***************.....[3703]
0.76   0.8    [64]**..................[9554] [326]**************......[3703]
0.8    0.84   [64]**..................[9554] [326]*************.......[3703]
0.84   0.88   [64]**..................[9554] [326]************........[3703]
0.88   0.92   [64]**..................[9554] [326]************........[3703]
0.92   0.96   [64]*...................[9554] [326]***********.........[3703]
0.96   1      [64]*...................[9554] [326]**********..........[3703]
1      1.04   [64]*...................[9554] [326]*********...........[3703]
1.04   1.08   [64]*...................[9554] [326]********............[3703]
1.08   1.12   [64]*...................[9554] [326]********............[3703]
1.12   1.16   [64]*...................[9554] [326]********............[3703]
1.16   1.2    [64]*...................[9554] [326]*******.............[3703]
1.2    1.24   [64]#...................[9554] [326]******..............[3703]
1.24   1.28   [64]#...................[9554] [326]*****...............[3703]
1.28   1.32   [64]#...................[9554] [326]*****...............[3703]
1.32   1.36   [64]#...................[9554] [326]****................[3703]
1.36   1.4    [64]#...................[9554] [326]****................[3703]
1.4    1.44   [64]#...................[9554] [326]****................[3703]
1.44   1.48   [64]#...................[9554] [326]***.................[3703]
1.48   1.52   [64]#...................[9554] [326]***.................[3703]
1.52   1.56   [64]#...................[9554] [326]***.................[3703]
1.56   1.6    [64]#...................[9554] [326]**..................[3703]
1.6    1.64   [64]#...................[9554] [326]**..................[3703]
1.64   1.68   [64]#...................[9554] [326]**..................[3703]
1.68   1.72   [64]#...................[9554] [326]*...................[3703]
1.72   1.76   [64]#...................[9554] [326]*...................[3703]
1.76   1.8    [64]#...................[9554] [326]*...................[3703]
1.8    1.84   [64]#...................[9554] [326]#...................[3703]
1.84   1.88   [64]#...................[9554] [326]#...................[3703]
1.88   1.92   [64]#...................[9554] [326]#...................[3703]
1.92   1.96   [64]#...................[9554] [326]#...................[3703]
1.96   2      [64]#...................[9554] [326]#...................[3703]
</pre>

## Randomly selecting words from a list

Given this [word list](./data/english-words.txt), first take a look to see what the first few lines look like:

<pre class="pre-highlight-in-pair">
<b>head data/english-words.txt</b>
</pre>
<pre class="pre-non-highlight-in-pair">
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

<pre class="pre-highlight-in-pair">
<b>mlr --from data/english-words.txt --nidx filter -S 'n=strlen($1);4<=n&&n<=8' then sample -k 10</b>
</pre>
<pre class="pre-non-highlight-in-pair">
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

These are simple *n*-grams, adapted from a previous version [described here](http://johnkerl.org/randspell/randspell-slides-ts.pdf). Some common functions are [located here](https://github.com/johnkerl/miller/blob/master/docs/src/ngrams/ngfuncs.mlr) with main Miller script [here](https://github.com/johnkerl/miller/blob/master/docs/src/ngrams/ngrams.mlr) and wrapper script [here](https://github.com/johnkerl/miller/blob/master/docs/src/ngrams/ngrams.sh).

The idea is that words from the input file are consumed, then taken apart and pasted back together in ways which imitate the letter-to-letter transitions found in the word list -- giving us automatically generated words in the same vein as *bromance* and *spork*:

<pre class="pre-highlight-in-pair">
<b>mlr --nidx --from ./ngrams/gsl-2000.txt put -q -f ./ngrams/ngfuncs.mlr -f ./ngrams/ngrams.mlr</b>
</pre>
<pre class="pre-non-highlight-in-pair">
burse
serious
land
seasure
clainst
tray
wherhoose
stry
jourt
strue
partist
ornear
devel
praction
roup
</pre>
