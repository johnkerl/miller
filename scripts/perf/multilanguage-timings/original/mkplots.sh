#!/bin/sh

experiments=`mlr cut -f experiment < tn.out|sort -u|sed 's/experiment=//'`
for experiment in $experiments; do
  mlr --onidx --ofs ' ' filter '$experiment=="'$experiment'"' then cut -x -f experiment tn.out | pgr -nc -title $experiment -xlabel nlines -ylabel seconds -legend 'a b' -lop &
done
