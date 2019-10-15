#!/bin/bash

counter=0
echo '<table border=1>'
mlr --list-all-functions-as-table | while read name class nargs; do
  counter=$[counter+1]

  if [ $counter -eq 1 ]; then
    echo '<tr class="mlrbg">'
    echo "<th>$name</th> <th>$class</th> <th>$nargs</th>"
  else
    echo '<tr>'
    echo "<td><a href="#$name">$name</a></td> <td>$class</td> <td>$nargs</td>"
  fi
  echo '</tr>'
done
echo '</table>'
