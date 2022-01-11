#!/bin/bash

# ================================================================
# See for context
# https://miller.readthedocs.io/en/latest/randomizing-examples/#randomly-generating-jabberwocky-words
# ================================================================

# ----------------------------------------------------------------
ourdir=$(dirname $0)
us=$(basename $0)

default_n=5
default_ocount=40
default_verbose="false"

usage() {
  echo "Usage: $us [options] {word-list files}" 1>&2
  echo "Options:" 1>&2
  echo "-n {n} The n for n-grams; default $default_n." 1>&2
  echo "-o {o} Number of words to produce; default $default_ocount." 1>&2
  echo "-v     Verbose processing; default off." 1>&2
  echo "If no wordlists are provided, stdin i=s read." 1>&2
  exit 1
}

# ----------------------------------------------------------------
n=$default_n
ocount=$default_ocount
verbose=$default_verbose
wordlist=$default_wordlist

while getopts n:o:vh? f
do
  case $f in
    n)  n="$OPTARG";      continue;;
    o)  ocount="$OPTARG"; continue;;
    v)  verbose="true";   continue;;
    h)  usage;;
    \?) usage;;
  esac
done
shift $(($OPTIND-1))

if [ $n -lt 2 ]; then
  echo "${us}: n must be >= 2." 1>&2
fi

wordlist="$@"

mlr --nidx put -q \
  -s n=$n -s ocount=$ocount -s verbose=$verbose \
  -f $ourdir/ngfuncs.mlr -f $ourdir/ngrams.mlr \
  $wordlist
