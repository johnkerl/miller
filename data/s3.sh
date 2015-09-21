#!/bin/sh
mlr --opprint "$@" step -a rsum,delta,counter -f x,y -g a ../data/small
echo
mlr --opprint "$@" step -a rsum,delta,counter -f x,y ../data/small
