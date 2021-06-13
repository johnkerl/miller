#!/bin/bash

# In case the user running this has a .mlrrc
export MLRRC=__none__

mlr -F | grep -v '^[a-zA-Z]' | uniq | while read funcname; do
  displayname=$funcname
  linkname="$funcname"
  if [ "$funcname" = '+' ]; then
    displayname='\+'
    linkname='plus'
  elif [ "$funcname" = '-' ]; then
    displayname='\-'
    linkname='minus'
  elif [ "$funcname" = '*' ]; then
    displayname='\*'
    linkname='times'
  elif [ "$funcname" = '**' ]; then
    displayname='\**'
    linkname='exponentiation'
  elif [ "$funcname" = '|' ]; then
    displayname='\|'
    linkname='bitwise-or'
  elif [ "$funcname" = '?' ]; then
    displayname='\?'
    linkname='question-mark'
  elif [ "$funcname" = ':' ]; then
    displayname='\:'
    linkname='colon'
  elif [ "$funcname" = '? :' ]; then
    displayname='\?'
    linkname='question-mark-colon'
  elif [ "$funcname" = '?:' ]; then
    displayname='\?'
    linkname='question-mark-colon'
  fi

  echo ""
  echo ".. _reference-dsl-${linkname}:"
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
  displayname=$funcname
  linkname="$funcname"
  if [ "$funcname" = '+' ]; then
    displayname='\+'
    linkname='plus'
  elif [ "$funcname" = '-' ]; then
    displayname='\-'
    linkname='minus'
  elif [ "$funcname" = '*' ]; then
    displayname='\*'
    linkname='times'
  elif [ "$funcname" = '**' ]; then
    displayname='\**'
    linkname='exponentiation'
  elif [ "$funcname" = '|' ]; then
    displayname='\|'
    linkname='bitwise-or'
  elif [ "$funcname" = '?' ]; then
    displayname='\?'
    linkname='question-mark'
  elif [ "$funcname" = ':' ]; then
    displayname='\:'
    linkname='colon'
  elif [ "$funcname" = '? :' ]; then
    displayname='\?'
    linkname='question-mark-colon'
  fi

  echo ""
  echo ".. _reference-dsl-${linkname}:"
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
