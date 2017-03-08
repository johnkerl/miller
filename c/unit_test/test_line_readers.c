#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include "lib/mlr_arch.h"
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "lib/minunit.h"
#include "lib/mlr_test_util.h"
#include "input/line_readers.h"

int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

typedef ssize_t getdelim_t(char** restrict pline, size_t* restrict plinecap, int delimiter, FILE* restrict stream);

// ----------------------------------------------------------------
static FILE* fopen_or_die(char* filename) {
	FILE* fp = fopen(filename, "r");
	if (fp == NULL) {
		perror("fopen");
		fprintf(stderr, "Couldn't open \"%s\" for read; exiting.\n", filename);
		exit(1);
	}
	return fp;
}

// ----------------------------------------------------------------
static char* test_getdelim_impl(getdelim_t* pgetdelim) {
	char delimiter = '\n';
	char* contents = NULL;
	char* path = NULL;
	FILE* fp = NULL;
	char* line = NULL;
	size_t linecap = 0;
	int rc = 0;

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	contents = "";
	path = write_temp_file_or_die(contents);
	fp = fopen_or_die(path);

	// Read line
	line = NULL;
	linecap = 0;
	rc = (*pgetdelim)(&line, &linecap, delimiter, fp);
	mu_assert_lf(rc == -1);
	mu_assert_lf(streq(line, ""));
	mu_assert_lf(linecap >= 1+strlen(contents));

	// Read past EOF
	line = NULL;
	linecap = 0;
	rc = (*pgetdelim)(&line, &linecap, delimiter, fp);
	mu_assert_lf(rc == -1);
	mu_assert_lf(streq(line, ""));

	fclose(fp);
	unlink_file_or_die(path);


	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	contents = "\n";
	path = write_temp_file_or_die(contents);
	fp = fopen_or_die(path);

	// Read line
	line = NULL;
	linecap = 0;
	rc = (*pgetdelim)(&line, &linecap, delimiter, fp);
	mu_assert_lf(rc == 1);
	mu_assert_lf(streq(line, "\n"));
	mu_assert_lf(linecap >= 1+strlen(contents));

	// Read past EOF
	line = NULL;
	linecap = 0;
	rc = (*pgetdelim)(&line, &linecap, delimiter, fp);
	mu_assert_lf(rc == -1);


	fclose(fp);
	unlink_file_or_die(path);

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	contents = "\r\n";
	path = write_temp_file_or_die(contents);
	fp = fopen_or_die(path);

	// Read line
	line = NULL;
	linecap = 0;
	rc = (*pgetdelim)(&line, &linecap, delimiter, fp);
	mu_assert_lf(rc == 2);
	mu_assert_lf(streq(line, "\r\n"));
	mu_assert_lf(linecap >= 1+strlen(contents));

	// Read past EOF
	line = NULL;
	linecap = 0;
	rc = (*pgetdelim)(&line, &linecap, delimiter, fp);
	mu_assert_lf(rc == -1);


	fclose(fp);
	unlink_file_or_die(path);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_mlr_alloc_read_line_single_delimiter() {
	char   delimiter = '\n';
	char*  path = NULL;
	FILE*  fp = NULL;
	char*  contents = NULL;
	char*  line = NULL;
	size_t linelen = 0;
	size_t linecap = 0;
	int    do_auto_line_term = FALSE;
	context_t* pctx = NULL;

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	contents = "";
	path = write_temp_file_or_die(contents);
	fp = fopen_or_die(path);
	linelen = linecap = 4;
	printf("\n");
	printf("Case start[%s]\n", contents);

	// Read line
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line == NULL);

	// Read past EOF
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	mu_assert_lf(line == NULL);

	fclose(fp);
	unlink_file_or_die(path);

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	contents = "\n";
	path = write_temp_file_or_die(contents);
	fp = fopen_or_die(path);
	linelen = linecap = 4;
	printf("\n");
	printf("Case start[%s]\n", contents);

	// Read line
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line != NULL);
	mu_assert_lf(streq(line, ""));
	mu_assert_lf(linelen == strlen(""));
	mu_assert_lf(linecap > linelen);

	// Read to EOF
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line == NULL);
	mu_assert_lf(linecap > linelen);

	fclose(fp);
	unlink_file_or_die(path);

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	contents = "abc\n";
	path = write_temp_file_or_die(contents);
	fp = fopen_or_die(path);
	linelen = linecap = 4;
	printf("\n");
	printf("Case start[%s]\n", contents);

	// Read line
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line != NULL);
	mu_assert_lf(streq(line, "abc"));
	mu_assert_lf(linelen == strlen("abc"));
	mu_assert_lf(linecap > linelen);

	// Read to EOF
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line == NULL);
	mu_assert_lf(linecap > linelen);

	fclose(fp);
	unlink_file_or_die(path);

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	contents = "abc\n\n";
	path = write_temp_file_or_die(contents);
	fp = fopen_or_die(path);
	linelen = linecap = 4;
	printf("\n");
	printf("Case start[%s]\n", contents);

	// Read line 1
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line != NULL);
	mu_assert_lf(streq(line, "abc"));
	mu_assert_lf(linelen == strlen("abc"));
	mu_assert_lf(linecap > linelen);

	// Read line 2
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line != NULL);
	mu_assert_lf(streq(line, ""));
	mu_assert_lf(linelen == strlen(""));
	mu_assert_lf(linecap > linelen);

	// Read to EOF
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line == NULL);
	mu_assert_lf(linecap > linelen);

	fclose(fp);
	unlink_file_or_die(path);

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	contents = "\nabc\n";
	path = write_temp_file_or_die(contents);
	fp = fopen_or_die(path);
	linelen = linecap = 4;
	printf("\n");
	printf("Case start[%s]\n", contents);

	// Read line 1
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line != NULL);
	mu_assert_lf(streq(line, ""));
	mu_assert_lf(linelen == strlen(""));
	mu_assert_lf(linecap > linelen);

	// Read line 2
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line != NULL);
	mu_assert_lf(streq(line, "abc"));
	mu_assert_lf(linelen == strlen("abc"));
	mu_assert_lf(linecap > linelen);

	// Read to EOF
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line == NULL);
	mu_assert_lf(linecap > linelen);

	fclose(fp);
	unlink_file_or_die(path);

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	contents = "\nabc";
	path = write_temp_file_or_die(contents);
	fp = fopen_or_die(path);
	linelen = linecap = 4;
	printf("\n");
	printf("Case start[%s]\n", contents);

	// Read line 1
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line != NULL);
	mu_assert_lf(streq(line, ""));
	mu_assert_lf(linelen == strlen(""));
	mu_assert_lf(linecap > linelen);

	// Read line 2
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line != NULL);
	mu_assert_lf(streq(line, "abc"));
	mu_assert_lf(linelen == strlen("abc"));
	mu_assert_lf(linecap > linelen);

	// Read to EOF
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line == NULL);
	mu_assert_lf(linecap > linelen);

	fclose(fp);
	unlink_file_or_die(path);


	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	return NULL;
}

// ----------------------------------------------------------------
// On unixish platforms this is testing the system getdelim() which is correct by definition.
// (At least, we verify it's behaving as we expect.) On Windows, this tests our local_getdelim
// which is a radundant test -- but again, it confirms behavior is as expected.
static char* test_getdelim() {
	return test_getdelim_impl(&mlr_arch_getdelim);
}

// ----------------------------------------------------------------
// This tests our homemade getdelim replacement, for running on Windows which lacks getdelim.
// xxx WIP
//static char* test_local_getdelim() {
//	return test_getdelim_impl(&local_getdelim);
//}

// ================================================================
static char * run_all_tests() {
	mu_run_test(test_getdelim);
	mu_run_test(test_mlr_alloc_read_line_single_delimiter);
	return 0;
}

int main(int argc, char **argv) {
	mlr_global_init(argv[0], NULL);
	printf("TEST_LINE_READERS ENTER\n");
	char *result = run_all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_LINE_READERS: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
