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
  --seed 0.25 \
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
