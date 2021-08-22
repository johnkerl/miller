#!/bin/bash

set -euo pipefail

iheight=$(stty size | mlr --nidx --fs space cut -f 1)
iwidth=$(stty size | mlr --nidx --fs space cut -f 2)

rcorn=-2.0
icorn=-2.0
side=4.0
iheight=50
iwidth=100
maxits=100
do_julia=false
jr=0.0
ji=0.0

if [ "$1" = "1" ]; then
  rcorn=-1.787582;icorn=-0.000002;side=0.000004;maxits=1000
elif [ "$1" = "2" ]; then
  rcorn=-0.162950;icorn=1.026100;side=0.000200;maxits=100000
elif [ "$1" = "3" ]; then
  rcorn=-1.755350;icorn=0.014230;side=0.000020;maxits=10000
elif [ "$1" =  4 ]; then
  do_julia=true;jr=0.35;ji=0.35;maxits=1000
elif [ "$1" =  5 ]; then
  do_julia=true;jr=0.0;ji=0.63567
elif [ "$1" =  6 ]; then
  do_julia=true;jr=0.4;ji=0.34745;maxits=1000
elif [ "$1" =  7 ]; then
  do_julia=true;jr=0.36;ji=0.36;maxits=80
elif [ "$1" =  8 ]; then
  do_julia=true;jr=-0.55;ji=0.55;maxits=100
elif [ "$1" =  9 ]; then
  do_julia=true;jr=-0.51;ji=0.51;maxits=1000
elif [ "$1" = 10 ]; then
  do_julia=true;jr=-1.26;ji=-0.03;rcorn=-0.3;icorn=-0.3;side=0.6
elif [ "$1" = 11 ]; then
  do_julia=true;jr=-1.26;ji=-0.03;rcorn=-0.6;icorn=-0.6;side=0.2
elif [ "$1" = 12 ]; then
  do_julia=true;jr=-1.26;ji=-0.03;rcorn=-0.75;icorn=-0.03125;side=.0625
elif [ "$1" = 13 ]; then
  do_julia=true;jr=-1.26;ji=-0.03;rcorn=-0.75;icorn=-0.01;side=.02
fi

mlr -n put \
  -s rcorn=$rcorn \
  -s icorn=$icorn \
  -s side=$side \
  -s iheight=$iheight \
  -s iwidth=$iwidth \
  -s maxits=$maxits \
  -s do_julia=$do_julia \
  -s jr=$jr \
  -s ji=$ji \
  -f programs/mand.mlr
