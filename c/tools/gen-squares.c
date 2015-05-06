#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <unistd.h>

// Samples unit-uniform x,y pairs, conditioned on being in two subboxes the size
// of which is described by the "chomp" parameter.
//
// 0 < chomp < 1:
//
// +----------------oooooooo
// |                oooooooo
// |                oooooooo
// |                oooooooo <-- chomp
// |                       |
// |                       |
// oooooooo                | <-- 1-chomp
// oooooooo                |
// oooooooo                |
// oooooooo----------------+
//        ^        ^
//        |        |
//    1-chomp    chomp
//
// chomp = 0 means output all x,y pairs
//
// -1 < chomp < 0:
//    1+chomp    -chomp
//        |        |
//        v        v
// oooooooo----------------+
// oooooooo                |
// oooooooo                |
// oooooooo                | <-- -chomp
// |                       |
// |                       |
// |                oooooooo <-- 1+chomp
// |                oooooooo
// |                oooooooo
// +----------------oooooooo

static void usage(char* argv0) {
	fprintf(stderr, "Usage: %s [chomp [n]]\n", argv0);
	exit(1);
}

int main(int argc, char** argv) {
	double chomp = 0.5;
	int n = 100000;
	if (argc == 2) {
		if (sscanf(argv[1], "%lf", &chomp) != 1)
			usage(argv[0]);
	} else if (argc == 3) {
		if (sscanf(argv[1], "%lf", &chomp) != 1)
			usage(argv[0]);
		if (sscanf(argv[2], "%d", &n) != 1)
			usage(argv[0]);
	} else {
		usage(argv[0]);
	}

	srand(time(0) ^ getpid());
	if (chomp > 0.0) {
		double lo = 1.0 - chomp;
		double hi = chomp;
		for (int i = 0; i < n; /* increment in loop */) {
			double x = (double)rand() / (double)RAND_MAX;
			double y = (double)rand() / (double)RAND_MAX;
			if ((x < lo && y < lo) || (x > hi && y > hi)) {
				i++;
				printf("x=%.18lf,y=%.18lf\n", x, y);
			}
		}
	} else {
		double lo = 1+chomp;
		double hi = -chomp;
		for (int i = 0; i < n; /* increment in loop */) {
			double x = (double)rand() / (double)RAND_MAX;
			double y = (double)rand() / (double)RAND_MAX;
			if ((x < lo && y > hi) || (x > hi && y < lo)) {
				i++;
				printf("x=%.18lf,y=%.18lf\n", x, y);
			}
		}
	}
	return 0;
}
