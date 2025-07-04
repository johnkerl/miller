#!/bin/bash

for x in *.md.in; do
    sed -i .emd 's/  *$//' $x
    rm $x.emd
done
