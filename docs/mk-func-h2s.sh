#!/bin/bash

# In case the user running this has a .mlrrc
export MLRRC=__none__

mlr -F | grep -v '^[a-zA-Z]' | uniq | while read funcname; do
  displayname=$funcname
  if [ "$funcname" = '+' ]; then
    displayname='\+'
  elif [ "$funcname" = '-' ]; then
    displayname='\-'
  elif [ "$funcname" = '*' ]; then
    displayname='\*'
  elif [ "$funcname" = '|' ]; then
    displayname='\|'
  fi

  echo ""
  echo ".. _\"$funcname\":"
  echo ""
  echo "$displayname"
  echo "^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^"
  echo ""
  echo '::'
  echo ''
  mlr --help-function "$funcname" | sed 's/^/    /'
  echo ''
  echo ''
done

mlr -F | grep '^[a-zA-Z]' | sort -u | while read funcname; do
  echo ""
  echo ".. _\"$funcname\":"
  echo ""
  echo "$displayname"
  echo "^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^"
  echo ""
  echo '::'
  echo ''
  mlr --help-function "$funcname" | sed 's/^/    /'
  echo ''
  echo ''
done

