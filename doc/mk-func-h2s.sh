#!/bin/bash

mlr -F | grep -v '^[a-zA-Z]' | uniq | while read funcname; do
  echo ""
  echo "<a id=$funcname/>"
  echo "<h2>$funcname</h2>"
  echo ""
  echo "<p/>"
  echo '<div class="pokipanel">'
  echo '<pre>'
  mlr --help-function "$funcname"
  echo '</pre>'
  echo '</div>'
  echo ""
done

mlr -F | grep '^[a-zA-Z]' | sort -u | while read funcname; do
  echo ""
  echo "<a id=$funcname/>"
  echo "<h2>$funcname</h2>"
  echo ""
  echo "<p/>"
  echo '<div class="pokipanel">'
  echo '<pre>'
  mlr --help-function "$funcname"
  echo '</pre>'
  echo '</div>'
  echo ""
done

