#!/bin/bash

while read funcname; do
  echo "<h2>$funcname</h2>"
  echo ""
  echo "<p/>"
  echo '<div class="pokipanel">'
  echo '<pre>'
  mlr --help-function $funcname | fmt -80
  echo '</pre>'
  echo '</div>'
  echo ""
done
