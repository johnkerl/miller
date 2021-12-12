#!/bin/bash

cat "$@" | while read line; do
  millisleep 100
  echo $line
done
