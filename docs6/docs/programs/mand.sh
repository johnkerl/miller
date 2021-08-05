#!/bin/bash

iheight=$(stty size | mlr --nidx --fs space cut -f 1)
iwidth=$(stty size | mlr --nidx --fs space cut -f 2)
if [ "$1" = "1" ]; then
  echo "rcorn=-1.787582,icorn=-0.000002,side=0.000004,maxits=1000,iheight=$iheight,iwidth=$iwidth"
elif [ "$1" = "2" ]; then
  echo "rcorn=-0.162950,icorn=+1.026100,side=0.000200,maxits=100000,iheight=$iheight,iwidth=$iwidth"
elif [ "$1" = "3" ]; then
  echo "rcorn=-1.755350,icorn=+0.014230,side=0.000020,maxits=10000,iheight=$iheight,iwidth=$iwidth"
elif [ "$1" =  4 ]; then
  echo "do_julia=true,jr= 0.35,ji=0.35,maxits=1000,iheight=$iheight,iwidth=$iwidth"
elif [ "$1" =  5 ]; then
  echo "do_julia=true,jr= 0.0,ji=0.63567,iheight=$iheight,iwidth=$iwidth"
elif [ "$1" =  6 ]; then
  echo "do_julia=true,jr= 0.4,ji=0.34745,maxits=1000,iheight=$iheight,iwidth=$iwidth"
elif [ "$1" =  7 ]; then
  echo "do_julia=true,jr= 0.36,ji=0.36,maxits=80,iheight=$iheight,iwidth=$iwidth"
elif [ "$1" =  8 ]; then
  echo "do_julia=true,jr=-0.55,ji=0.55,maxits=100,iheight=$iheight,iwidth=$iwidth"
elif [ "$1" =  9 ]; then
  echo "do_julia=true,jr=-0.51,ji=0.51,maxits=1000,iheight=$iheight,iwidth=$iwidth"
elif [ "$1" = 10 ]; then
  echo "do_julia=true,jr=-1.26,ji=-0.03,rcorn=-0.3,icorn=-0.3,side=0.6,iheight=$iheight,iwidth=$iwidth"
elif [ "$1" = 11 ]; then
  echo "do_julia=true,jr=-1.26,ji=-0.03,rcorn=-0.6,icorn=-0.6,side=.2,iheight=$iheight,iwidth=$iwidth"
elif [ "$1" = 12 ]; then
  echo "do_julia=true,jr=-1.26,ji=-0.03,rcorn=-0.75,icorn=-0.03125,side=.0625,iheight=$iheight,iwidth=$iwidth"
elif [ "$1" = 13 ]; then
  echo "do_julia=true,jr=-1.26,ji=-0.03,rcorn=-0.75,icorn=-0.01,side=.02,iheight=$iheight,iwidth=$iwidth"
else
  echo "iheight=$iheight,iwidth=$iwidth"
fi | mlr put -f programs/mand.mlr 
