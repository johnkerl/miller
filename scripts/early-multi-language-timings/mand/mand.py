#/usr/bin/python

from __future__ import division;
import sys

rcorn   = -2.0;
icorn   = -2.0;
side    =  4.0;
iheight =  500;
iwidth  = 1000;
maxits  =  100;

for ii in range(0, iheight+1):
	for ir in range(0, iwidth+1):
		cr = rcorn + (ir/iwidth) * side;
		ci = icorn + (ii/iheight) * side;

		zr = 0.0;
		zi = 0.0;

		# z := z^2 + c
		iti = 0;
		escaped = False;
		for iti in range(0, maxits):
			mag = zr*zr + zi+zi;
			if mag > 4.0:
					escaped = True;
					break;
			zt = zr*zr - zi*zi + cr;
			zi = 2*zr*zi + ci;
			zr = zt;
		if (escaped):
			sys.stdout.write("o");
		else:
			sys.stdout.write(".");
	sys.stdout.write("\n");
