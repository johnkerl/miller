#!/bin/bash

cat "$@" | while read line; do
  sleep 1
  echo $line
done
