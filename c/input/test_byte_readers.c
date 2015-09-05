#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include "lib/mlrutil.h"
#include "lib/minunit.h"
#include "input/byte_readers.h"

#ifdef __TEST_BYTE_READERS_MAIN__
int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

// ----------------------------------------------------------------
static char* test_string_byte_reader() {
	byte_reader_t* pbr = string_byte_reader_alloc();

	int ok = pbr->popen_func(pbr, "");
	mu_assert_lf(ok == TRUE);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	pbr->pclose_func(pbr);

	ok = pbr->popen_func(pbr, "a");
	mu_assert_lf(ok == TRUE);
	mu_assert_lf(pbr->pread_func(pbr) == 'a');
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	pbr->pclose_func(pbr);

	ok = pbr->popen_func(pbr, "abc");
	mu_assert_lf(ok == TRUE);
	mu_assert_lf(pbr->pread_func(pbr) == 'a');
	mu_assert_lf(pbr->pread_func(pbr) == 'b');
	mu_assert_lf(pbr->pread_func(pbr) == 'c');
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	pbr->pclose_func(pbr);

	return NULL;
}

#if 0
// ----------------------------------------------------------------
// xxx mkstemp
// xxx tmpfile
// xxx pop from buf
// xxx take dirname from argv[1]?

// ----------------------------------------------------------------
static FILE* make_temp_file(char* contents) {
	xxx

	int fd = mkstemp("/tmp/mlr-ut-XXXXXXXX");

}
#endif
static char* test_temp() {
	printf("hello\n");
	char* path = mktemp(strdup("/tmp/mlr.ut.XXXXXXXX"));
	printf("path=%s\n", path);
	FILE* fp = fopen(path, "w");
	char* buf = "a=1,b=2\nc=3";
	int len = strlen(buf);
	int rc = fwrite(buf, 1, len, fp);
	if (rc != len) {
		perror("fwrite");
		fprintf(stderr, "fwrite (%d) to \"%s\" failed.\n", len, path);
		exit(1);
	}
	fclose(fp);
	rc = unlink(path);
	if (rc != 0) {
		perror("unlink");
		fprintf(stderr, "unlink of \"%s\" failed.\n", path);
		exit(1);
	}

	return NULL;
}

// ================================================================
static char * run_all_tests() {
	mu_run_test(test_string_byte_reader);
	mu_run_test(test_temp);
	return 0;
}

int main(int argc, char **argv) {
	char *result = run_all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_BYTE_READERS: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
#endif // __TEST_BYTE_READERS_MAIN__
