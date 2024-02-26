#!/usr/bin/env python

import sys

nrow = 2
ncol = 100
if len(sys.argv) == 2:
    ncol = int(sys.argv[1])
if len(sys.argv) == 3:
    nrow = int(sys.argv[1])
    ncol = int(sys.argv[2])

prefix = "k"
for i in range(nrow):
    for j in range(ncol):
        if j == 0:
            sys.stdout.write("%s%07d" % (prefix, j))
        else:
                sys.stdout.write("\t%s%07d" % (prefix, j))
    sys.stdout.write("\n")
    prefix = "v"
