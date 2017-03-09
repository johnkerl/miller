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

static void hex_print(char* pre, char* string, char* post) {
	printf("%s", pre);
	for (char* p = string; *p; p++) {
		if (p > string)
			printf(" ");
		printf("%02x", *p);
	}
	printf("%s", post);
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
	hex_print("\nCase start: [", contents, "]\n");

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
	hex_print("\nCase start: [", contents, "]\n");

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
	hex_print("\nCase start: [", contents, "]\n");

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
	hex_print("\nCase start: [", contents, "]\n");

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
	hex_print("\nCase start: [", contents, "]\n");

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
	hex_print("\nCase start: [", contents, "]\n");

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
static char* test_mlr_alloc_read_line_single_delimiter_with_autodetect() {
	char   delimiter = '\n';
	char*  path = NULL;
	FILE*  fp = NULL;
	char*  contents = NULL;
	char*  line = NULL;
	size_t linelen = 0;
	size_t linecap = 0;
	int    do_auto_line_term = TRUE;
	context_t ctx;
	context_t* pctx = &ctx;

	context_init_from_first_file_name(pctx, "fake-file-name");

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	contents = "";
	path = write_temp_file_or_die(contents);
	fp = fopen_or_die(path);
	linelen = linecap = 4;
	hex_print("\nCase start: [", contents, "]\n");

	// Read line
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line == NULL);
	mu_assert_lf(pctx->auto_line_term_detected == TRUE);
	mu_assert_lf(streq(pctx->auto_line_term, "\n"));

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
	hex_print("\nCase start: [", contents, "]\n");

	// Read line
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line != NULL);
	mu_assert_lf(streq(line, ""));
	mu_assert_lf(linelen == strlen(""));
	mu_assert_lf(linecap > linelen);
	mu_assert_lf(pctx->auto_line_term_detected == TRUE);
	mu_assert_lf(streq(pctx->auto_line_term, "\n"));

	// Read to EOF
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line == NULL);
	mu_assert_lf(linecap > linelen);
	mu_assert_lf(pctx->auto_line_term_detected == TRUE);
	mu_assert_lf(streq(pctx->auto_line_term, "\n"));

	fclose(fp);
	unlink_file_or_die(path);

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	contents = "abc\n";
	path = write_temp_file_or_die(contents);
	fp = fopen_or_die(path);
	linelen = linecap = 4;
	hex_print("\nCase start: [", contents, "]\n");

	// Read line
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line != NULL);
	mu_assert_lf(streq(line, "abc"));
	mu_assert_lf(linelen == strlen("abc"));
	mu_assert_lf(linecap > linelen);
	mu_assert_lf(pctx->auto_line_term_detected == TRUE);
	mu_assert_lf(streq(pctx->auto_line_term, "\n"));

	// Read to EOF
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line == NULL);
	mu_assert_lf(linecap > linelen);
	mu_assert_lf(pctx->auto_line_term_detected == TRUE);
	mu_assert_lf(streq(pctx->auto_line_term, "\n"));

	fclose(fp);
	unlink_file_or_die(path);

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	contents = "abc\n\n";
	path = write_temp_file_or_die(contents);
	fp = fopen_or_die(path);
	linelen = linecap = 4;
	hex_print("\nCase start: [", contents, "]\n");

	// Read line 1
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line != NULL);
	mu_assert_lf(streq(line, "abc"));
	mu_assert_lf(linelen == strlen("abc"));
	mu_assert_lf(linecap > linelen);
	mu_assert_lf(pctx->auto_line_term_detected == TRUE);
	mu_assert_lf(streq(pctx->auto_line_term, "\n"));

	// Read line 2
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line != NULL);
	mu_assert_lf(streq(line, ""));
	mu_assert_lf(linelen == strlen(""));
	mu_assert_lf(linecap > linelen);
	mu_assert_lf(pctx->auto_line_term_detected == TRUE);
	mu_assert_lf(streq(pctx->auto_line_term, "\n"));

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
	hex_print("\nCase start: [", contents, "]\n");

	// Read line 1
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line != NULL);
	mu_assert_lf(streq(line, ""));
	mu_assert_lf(linelen == strlen(""));
	mu_assert_lf(linecap > linelen);
	mu_assert_lf(pctx->auto_line_term_detected == TRUE);
	mu_assert_lf(streq(pctx->auto_line_term, "\n"));

	// Read line 2
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line != NULL);
	mu_assert_lf(streq(line, "abc"));
	mu_assert_lf(linelen == strlen("abc"));
	mu_assert_lf(linecap > linelen);
	mu_assert_lf(pctx->auto_line_term_detected == TRUE);
	mu_assert_lf(streq(pctx->auto_line_term, "\n"));

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
	hex_print("\nCase start: [", contents, "]\n");

	// Read line 1
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line != NULL);
	mu_assert_lf(streq(line, ""));
	mu_assert_lf(linelen == strlen(""));
	mu_assert_lf(linecap > linelen);
	mu_assert_lf(pctx->auto_line_term_detected == TRUE);
	mu_assert_lf(streq(pctx->auto_line_term, "\n"));

	// Read line 2
	line = mlr_alloc_read_line_single_delimiter(fp, delimiter, &linelen, &linecap, do_auto_line_term, pctx);
	printf("linelen=%d linecap=%d line=\"%s\"\n", (int)linelen, (int)linecap, line);
	mu_assert_lf(line != NULL);
	mu_assert_lf(streq(line, "abc"));
	mu_assert_lf(linelen == strlen("abc"));
	mu_assert_lf(linecap > linelen);
	mu_assert_lf(pctx->auto_line_term_detected == TRUE);
	mu_assert_lf(streq(pctx->auto_line_term, "\n"));

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

// ================================================================
static char * run_all_tests() {
	mu_run_test(test_mlr_alloc_read_line_single_delimiter);
	mu_run_test(test_mlr_alloc_read_line_single_delimiter_with_autodetect);
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
