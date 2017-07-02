#!/bin/bash

echo ================================================================
#### xxx temp
#### sed 's/-lm/-lm -lpcreposix/' c/Makefile.am > temp; mv temp c/Makefile.am
#### xxx temp
#### sed 's/-lm/-lm -lpcreposix/' c/unit_test/Makefile.am > temp; mv temp c/unit_test/Makefile.am
sed 's/#undef MLR_ON_MSYS2/#define MLR_ON_MSYS2/' c/lib/mlr_arch.h > temp; mv temp c/lib/mlr_arch.h

echo ================================================================
./configure
make -C c/parsing lemon.exe
make
