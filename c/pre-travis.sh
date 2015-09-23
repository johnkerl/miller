#!/bin/sh -e
cd ..
make distclean
autoreconf -fiv
./configure
make check
