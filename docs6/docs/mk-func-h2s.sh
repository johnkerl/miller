#!/bin/bash

# In case the user running this has a .mlrrc
export MLRRC=__none__

mlr help list-functions | grep -v '^[a-zA-Z]' | uniq | while read funcname; do
  # Displayname is for HTML rendered in the document body. Linkname is for
  # #-style links which don't do so well with special characters.
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
  elif [ "$funcname" = '&' ]; then
    displayname='&'
    linkname='bitwise-and'
  elif [ "$funcname" = '^' ]; then
    displayname='^'
    linkname='bitwise-xor'

  elif [ "$funcname" = '||' ]; then
    displayname='\|\|'
    linkname='logical-or'
  elif [ "$funcname" = '&&' ]; then
    displayname='&&'
    linkname='logical-and'
  elif [ "$funcname" = '^^' ]; then
    displayname='^^'
    linkname='logical-xor'

  elif [ "$funcname" = '?' ]; then
    displayname='?'
    linkname='question-mark'
  elif [ "$funcname" = ':' ]; then
    displayname='\:'
    linkname='colon'
  elif [ "$funcname" = '?:' ]; then
    displayname='?:'
    linkname='question-mark-colon'
  elif [ "$funcname" = '? :' ]; then
    displayname='?:'
    linkname='question-mark-colon'
  elif [ "$funcname" = '??' ]; then
    displayname='??'
    linkname='absent-coalesce'
  elif [ "$funcname" = '???' ]; then
    displayname='???'
    linkname='absent-empty-coalesce'

  elif [ "$funcname" = '!' ]; then
    displayname='\!'
    linkname='exclamation-point'

  elif [ "$funcname" = '=~' ]; then
    displayname='=~'
    linkname='regmatch'
  elif [ "$funcname" = '!=~' ]; then
    displayname='!=~'
    linkname='regnotmatch'
  elif [ "$funcname" = '>>' ]; then
    displayname='\>\>'
    linkname='srsh'
  elif [ "$funcname" = '>>>' ]; then
    displayname='\>\>\>'
    linkname='ursh'
  elif [ "$funcname" = '>>>=' ]; then
    displayname='\>\>\>='
    linkname='ursheq'
  fi

  # TODO: fix section-links for mnkdocs
  #echo ''
  #echo ".. _reference-dsl-${linkname}:"
  echo ''
  if [ "$linkname" = "$displayname" ]; then
    echo "## $displayname"
  else
    echo "<a id=$linkname />"
    echo "## $displayname"
  fi
  echo ''
  echo '<pre class="pre-non-highlight-non-pair">'
  mlr help function "$funcname"
  echo '</pre>'
  echo ''
done

mlr help list-functions | grep '^[a-zA-Z]' | sort -u | while read funcname; do
  displayname=$funcname
  linkname="$funcname"

  # TODO: fix section-links for mnkdocs
  #echo ''
  #echo ".. _reference-dsl-${linkname}:"
  echo ''
  if [ "$linkname" = "$displayname" ]; then
    echo "## $displayname"
  else
    echo "<a id=$linkname />"
    echo "## $displayname"
  fi
  echo ''
  echo '<pre class="pre-non-highlight-non-pair">'
  mlr help function "$funcname"
  echo '</pre>'
  echo ''
done
