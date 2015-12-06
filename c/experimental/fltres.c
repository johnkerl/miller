#include <stdio.h>

static void show(long long i) {
	double d = (double) i;
	long long j = (long long) d;
	printf("%016llx %lld %lf %016llx\n", i, i, d, j);
}

int main() {
	show(0x7ffffffffffff9ff);
	show(0x7ffffffffffffa00);
	show(0x7ffffffffffffbff);
	show(0x7ffffffffffffc00);
	show(0x7ffffffffffffdff);
	show(0x7ffffffffffffe00);
	show(0x7ffffffffffffffe);
	show(0x7fffffffffffffff);

	return 0;
}
