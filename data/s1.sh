#!/bin/sh
mlr --opprint "$@" stats1 -a mean,sum,count,min,max -f i,x,y -g a,b ../data/medium
echo
mlr --opprint "$@" stats1 -a mean,sum,count,min,max -f i,x,y ../data/medium
