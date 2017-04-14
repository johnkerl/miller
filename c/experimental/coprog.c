#include <stdio.h>
#include <stdlib.h>

int main(int argc, char** argv) {
	char* cmd = argv[1];
	FILE* pfp = popen(cmd, "r+");
	if (pfp == NULL) {
		perror("popen");
		exit(1);
	}
	if (setlinebuf(pfp) < 0) {
		perror("setlinebuf");
		exit(1);
	}
	char* line1 = NULL;
	size_t linecap1 = 0;
	ssize_t linelen1;

	while ((linelen1 = getline(&line1, &linecap1, stdin)) > 0) {
		if (linelen1 > 0 && line1[linelen1-1] == '\n') {
			line1[linelen1-1] = 0;
			linelen1--;
		}
		printf("ILINE: %s\n", line1);
		int rc2 = fprintf(pfp, "%s\n", line1);
		if (rc2 < 0) {
			perror("fprintf");
			exit(1);
		}
		if (fflush(pfp) < 0) {
			perror("fflush");
			exit(1);
		}
		char* line2 = NULL;
		size_t linecap2 = 0;
		ssize_t linelen2;
		linelen2 = getline(&line2, &linecap2, pfp);
		if (linelen2 <= 0) {
			perror("pipe read");
			exit(1);
		}
		if (linelen2 > 0 && line2[linelen2-1] == '\n') {
			line2[linelen2-1] = 0;
			linelen2--;
		}
		printf("PLINE: %s\n", line2);
	}
	fclose(pfp);

	return 0;
}
