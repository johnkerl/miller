#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"

// ----------------------------------------------------------------
char* write_temp_file_or_die(char* contents) {
	char* path = mktemp(mlr_strdup_or_die("/tmp/mlr.ut.XXXXXXXX"));
	FILE* fp = fopen(path, "w");
	int len = strlen(contents);
	int rc = fwrite(contents, 1, len, fp);
	if (rc != len) {
		perror("fwrite");
		fprintf(stderr, "%s: fwrite (%d) to \"%s\" failed.\n",
		MLR_GLOBALS.argv0, len, path);
		exit(1);
	}
	fclose(fp);
	return path;
}

// ----------------------------------------------------------------
void unlink_file_or_die(char* path) {
	int rc = unlink(path);
	if (rc != 0) {
		perror("unlink");
		fprintf(stderr, "unlink of \"%s\" failed.\n", path);
		exit(1);
	}
}
