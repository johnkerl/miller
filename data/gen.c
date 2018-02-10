#include <stdio.h>

int main(int argc, char** argv) {
	long long n = 100;
	if (argc == 2) {
		(void)sscanf(argv[1], "%lld", &n);
	}
	for (long long i = 0; i < n; i++) {
		printf("i=%lld\n", i);
	}
	return 0;
}
