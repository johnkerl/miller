#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main(int argc, char** argv) {
	double rcorn     = -2.0;
	double icorn     = -2.0;
	double side      =  4.0;
	int    iheight   =  500;
	int    iwidth    = 1000;
	int    maxits    =  100;
	int    levelstep =  100;
	double jr        = 0.0;
	double ji        = 0.0;
	int    do_julia  = 0;
	char*  chars     = "@X*o-.";
	int    nchars = strlen(chars);

	for (int argi = 1; argi < argc; argi++) {
		if (sscanf(argv[argi], "rcorn=%lf", &rcorn)) {
		} else if (sscanf(argv[argi], "icorn=%lf", &icorn)) {
		} else if (sscanf(argv[argi], "side=%lf", &side)) {
		} else if (sscanf(argv[argi], "iheight=%d", &iheight)) {
		} else if (sscanf(argv[argi], "iwidth=%d", &iwidth)) {
		} else if (sscanf(argv[argi], "dims=%d %d", &iheight, &iwidth)) {
		} else if (sscanf(argv[argi], "maxits=%d", &maxits)) {
		} else if (sscanf(argv[argi], "levelstep=%d", &levelstep)) {
		} else if (sscanf(argv[argi], "julia=%lf,%lf", &jr,&ji)) {
			do_julia = 1;
		} else {
			fprintf(stderr, "b04k!\n");
			exit(1);
		}
	}

	putchar('+');
	for (int ir = 1; ir < iwidth-1; ir += 1) {
		putchar('-');
	}
	putchar('+');
	putchar('\n');

	for (int ii = iheight-1; ii > 0; ii -= 1) {
		double pi = icorn + ((double)ii/(double)iheight) * side;
		putchar('|');
		for (int ir = 1; ir < iwidth-1; ir += 1) {
			double pr = rcorn + ((double)ir/(double)iwidth) * side;

			double zr = 0.0;
			double zi = 0.0;
			double cr = 0.0;
			double ci = 0.0;

			if (!do_julia) {
				zr = 0.0;
				zi = 0.0;
				cr = pr;
				ci = pi;
			} else {
				zr = pr;
				zi = pi;
				cr = jr;
				ci = ji;
			}

			// z := z^2 + c
			int iti = 0;
			int escaped = 0;
			for (iti = 0; iti < maxits; iti += 1) {
				double mag = zr*zr + zi+zi;
				if (mag > 4.0) {
						escaped = 1;
						break;
				}
				double zt = zr*zr - zi*zi + cr;
				zi = 2*zr*zi + ci;
				zr = zt;
			}

			if (!escaped) {
				putchar(' ');
			} else {
				int level = (iti / levelstep) % nchars;
				putchar(chars[level]);
			}
		}
		putchar('|');
		putchar('\n');
	}

	putchar('+');
	for (int ir = 1; ir < iwidth-1; ir += 1) {
		putchar('-');
	}
	putchar('+');
	putchar('\n');

	return 0;
}
