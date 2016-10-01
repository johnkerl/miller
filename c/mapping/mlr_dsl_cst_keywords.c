#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/hss.h"
#include "mlr_dsl_cst.h"
#include "context_flags.h"

// ----------------------------------------------------------------
typedef void keyword_usage_func_t(FILE* ostream);

// E
// ENV
// FILENAME
// FILENUM
// FNR
// NF
// NR
// PI

// Should match the keywords defined in dsls/mlr_dsl_lexer.l.

static keyword_usage_func_t mlr_dsl_all_keyword_usage;
static keyword_usage_func_t mlr_dsl_begin_keyword_usage;
static keyword_usage_func_t mlr_dsl_break_keyword_usage;
static keyword_usage_func_t mlr_dsl_call_keyword_usage;
static keyword_usage_func_t mlr_dsl_continue_keyword_usage;
static keyword_usage_func_t mlr_dsl_do_keyword_usage;
static keyword_usage_func_t mlr_dsl_dump_keyword_usage;
static keyword_usage_func_t mlr_dsl_edump_keyword_usage;
static keyword_usage_func_t mlr_dsl_elif_keyword_usage;
static keyword_usage_func_t mlr_dsl_else_keyword_usage;
static keyword_usage_func_t mlr_dsl_emit_keyword_usage;
static keyword_usage_func_t mlr_dsl_emitf_keyword_usage;
static keyword_usage_func_t mlr_dsl_emitp_keyword_usage;
static keyword_usage_func_t mlr_dsl_end_keyword_usage;
static keyword_usage_func_t mlr_dsl_eprint_keyword_usage;
static keyword_usage_func_t mlr_dsl_eprintn_keyword_usage;
static keyword_usage_func_t mlr_dsl_filter_keyword_usage;
static keyword_usage_func_t mlr_dsl_for_keyword_usage;
static keyword_usage_func_t mlr_dsl_func_keyword_usage;
static keyword_usage_func_t mlr_dsl_if_keyword_usage;
static keyword_usage_func_t mlr_dsl_in_keyword_usage;
static keyword_usage_func_t mlr_dsl_local_keyword_usage;
static keyword_usage_func_t mlr_dsl_print_keyword_usage;
static keyword_usage_func_t mlr_dsl_printn_keyword_usage;
static keyword_usage_func_t mlr_dsl_return_keyword_usage;
static keyword_usage_func_t mlr_dsl_stderr_keyword_usage;
static keyword_usage_func_t mlr_dsl_stdout_keyword_usage;
static keyword_usage_func_t mlr_dsl_subr_keyword_usage;
static keyword_usage_func_t mlr_dsl_tee_keyword_usage;
static keyword_usage_func_t mlr_dsl_unset_keyword_usage;
static keyword_usage_func_t mlr_dsl_while_keyword_usage;

// ----------------------------------------------------------------
typedef struct _keyword_usage_entry_t {
	char* name;
	keyword_usage_func_t* pusage_func;
} keyword_usage_entry_t;

static keyword_usage_entry_t KEYWORD_USAGE_TABLE[] = {

	{ "all",      mlr_dsl_all_keyword_usage      },
	{ "begin",    mlr_dsl_begin_keyword_usage    },
	{ "break",    mlr_dsl_break_keyword_usage    },
	{ "call",     mlr_dsl_call_keyword_usage     },
	{ "continue", mlr_dsl_continue_keyword_usage },
	{ "do",       mlr_dsl_do_keyword_usage       },
	{ "dump",    mlr_dsl_dump_keyword_usage    },
	{ "edump",   mlr_dsl_edump_keyword_usage   },
	{ "elif",     mlr_dsl_elif_keyword_usage     },
	{ "else",     mlr_dsl_else_keyword_usage     },
	{ "emit",    mlr_dsl_emit_keyword_usage    },
	{ "emitf",   mlr_dsl_emitf_keyword_usage   },
	{ "emitp",   mlr_dsl_emitp_keyword_usage   },
	{ "end",      mlr_dsl_end_keyword_usage      },
	{ "eprint",  mlr_dsl_eprint_keyword_usage  },
	{ "eprintn", mlr_dsl_eprintn_keyword_usage },
	{ "filter",  mlr_dsl_filter_keyword_usage  },
	{ "for",      mlr_dsl_for_keyword_usage      },
	{ "func",     mlr_dsl_func_keyword_usage     },
	{ "if",       mlr_dsl_if_keyword_usage       },
	{ "in",       mlr_dsl_in_keyword_usage       },
	{ "local",    mlr_dsl_local_keyword_usage    },
	{ "print",   mlr_dsl_print_keyword_usage   },
	{ "printn",  mlr_dsl_printn_keyword_usage  },
	{ "return",   mlr_dsl_return_keyword_usage   },
	{ "stderr",  mlr_dsl_stderr_keyword_usage  },
	{ "stdout",  mlr_dsl_stdout_keyword_usage  },
	{ "subr",     mlr_dsl_subr_keyword_usage     },
	{ "tee",     mlr_dsl_tee_keyword_usage     },
	{ "unset",   mlr_dsl_unset_keyword_usage   },
	{ "while",    mlr_dsl_while_keyword_usage    },

};
static int KEYWORD_USAGE_TABLE_SIZE = sizeof(KEYWORD_USAGE_TABLE)/sizeof(KEYWORD_USAGE_TABLE[0]);

// ----------------------------------------------------------------


// ================================================================
// Pass function_name == NULL to get usage for all keywords.
// Note keywords are defined in dsls/mlr_dsl_lexer.l.
void mlr_dsl_keyword_usage(FILE* ostream, char* keyword) {
	if (keyword == NULL) {
		for (int i = 0; i < KEYWORD_USAGE_TABLE_SIZE; i++) {
			if (i > 0) {
				fprintf(ostream, "\n");
			}
			KEYWORD_USAGE_TABLE[i].pusage_func(ostream);
		}
	} else {

		int found = FALSE;
		for (int i = 0; i < KEYWORD_USAGE_TABLE_SIZE; i++) {
			if (streq(keyword, KEYWORD_USAGE_TABLE[i].name)) {
				KEYWORD_USAGE_TABLE[i].pusage_func(ostream);
				found = TRUE;
				break;
			}
		}
		if (!found) {
			fprintf(ostream, "%s: unrecognized keyword \"%s\".\n", MLR_GLOBALS.bargv0, keyword);
		}
	}
}

void mlr_dsl_list_all_keywords_raw(FILE* ostream) {
	for (int i = 0; i < KEYWORD_USAGE_TABLE_SIZE; i++) {
		printf("%s\n", KEYWORD_USAGE_TABLE[i].name);
	}
}

// ----------------------------------------------------------------

static void mlr_dsl_all_keyword_usage(FILE* ostream) {
	fprintf(ostream, "all: xxx temp\n");
}

static void mlr_dsl_begin_keyword_usage(FILE* ostream) {
	fprintf(ostream, "begin: xxx temp\n");
}

static void mlr_dsl_break_keyword_usage(FILE* ostream) {
	fprintf(ostream, "break: xxx temp\n");
}

static void mlr_dsl_call_keyword_usage(FILE* ostream) {
	fprintf(ostream, "call: xxx temp\n");
}

static void mlr_dsl_continue_keyword_usage(FILE* ostream) {
	fprintf(ostream, "continue: xxx temp\n");
}

static void mlr_dsl_do_keyword_usage(FILE* ostream) {
	fprintf(ostream, "do: xxx temp\n");
}

// ----------------------------------------------------------------
static void mlr_dsl_dump_keyword_usage(FILE* ostream) {
    fprintf(ostream,
		"dump: prints all currently defined out-of-stream variables immediately\n"
		"  to stdout as JSON.\n"
		"\n"
		"  With >, >>, or |, the data do not become part of the output record stream but\n"
		"  are instead redirected.\n"
		"\n"
		"  The > and >> are for write and append, as in the shell, but (as with awk) the\n"
		"  file-overwrite for > is on first write, not per record. The | is for piping to\n"
		"  a process which will process the data. There will be one open file for each\n"
		"  distinct file name (for > and >>) or one subordinate process for each distinct\n"
		"  value of the piped-to command (for |). Output-formatting flags are taken from\n"
		"  the main command line.\n"
		"\n"
		"  Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump }'\n"
		"  Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump >  \"mytap.dat\"}'\n"
		"  Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump >> \"mytap.dat\"}'\n"
		"  Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump | \"jq .[]\"}'\n");
}

static void mlr_dsl_edump_keyword_usage(FILE* ostream) {
    fprintf(ostream,
		"edump: prints all currently defined out-of-stream variables immediately\n"
		"  to stderr as JSON.\n"
		"\n"
		"  Example: mlr --from f.dat put -q '@v[NR]=$*; end { edump }'\n");
}

static void mlr_dsl_elif_keyword_usage(FILE* ostream) {
	fprintf(ostream, "elif: xxx temp\n");
}

static void mlr_dsl_else_keyword_usage(FILE* ostream) {
	fprintf(ostream, "else: xxx temp\n");
}

static void mlr_dsl_emit_keyword_usage(FILE* ostream) {
    fprintf(ostream,
		"emit: inserts an out-of-stream variable into the output record stream. Hashmap\n"
		"  indices present in the data but not slotted by emit arguments are not output.\n"
		"\n"
		"  With >, >>, or |, the data do not become part of the output record stream but\n"
		"  are instead redirected.\n"
		"\n"
		"  The > and >> are for write and append, as in the shell, but (as with awk) the\n"
		"  file-overwrite for > is on first write, not per record. The | is for piping to\n"
		"  a process which will process the data. There will be one open file for each\n"
		"  distinct file name (for > and >>) or one subordinate process for each distinct\n"
		"  value of the piped-to command (for |). Output-formatting flags are taken from\n"
		"  the main command line.\n"
		"\n"
		"  You can use any of the output-format command-line flags, e.g. --ocsv, --ofs,\n"
		"  etc., to control the format of the output if the output is redirected. See also %s -h.\n"
		"\n"
		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit @sums'\n"
		"  Example: mlr --from f.dat put --ojson '@sums[$a][$b]+=$x; emit > \"tap-\".$a.$b.\".dat\", @sums'\n"
		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit @sums, \"index1\", \"index2\"'\n"
		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit @*, \"index1\", \"index2\"'\n"
		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit >  \"mytap.dat\", @*, \"index1\", \"index2\"'\n"
		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit >> \"mytap.dat\", @*, \"index1\", \"index2\"'\n"
		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit | \"gzip > mytap.dat.gz\", @*, \"index1\", \"index2\"'\n"
		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit > stderr, @*, \"index1\", \"index2\"'\n"
		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit | \"grep somepattern\", @*, \"index1\", \"index2\"'\n"
		"\n"
		"  Please see http://johnkerl.org/miller/doc for more information.\n",
		MLR_GLOBALS.bargv0);
}

static void mlr_dsl_emitf_keyword_usage(FILE* ostream) {
    fprintf(ostream,
		"emitf: inserts non-indexed out-of-stream variable(s) side-by-side into the\n"
		"  output record stream.\n"
		"\n"
		"  With >, >>, or |, the data do not become part of the output record stream but\n"
		"  are instead redirected.\n"
		"\n"
		"  The > and >> are for write and append, as in the shell, but (as with awk) the\n"
		"  file-overwrite for > is on first write, not per record. The | is for piping to\n"
		"  a process which will process the data. There will be one open file for each\n"
		"  distinct file name (for > and >>) or one subordinate process for each distinct\n"
		"  value of the piped-to command (for |). Output-formatting flags are taken from\n"
		"  the main command line.\n"
		"\n"
		"  You can use any of the output-format command-line flags, e.g. --ocsv, --ofs,\n"
		"  etc., to control the format of the output if the output is redirected. See also %s -h.\n"
		"\n"
		"  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf @a'\n"
		"  Example: mlr --from f.dat put --oxtab '@a=$i;@b+=$x;@c+=$y; emitf > \"tap-\".$i.\".dat\", @a'\n"
		"  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf @a, @b, @c'\n"
		"  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf > \"mytap.dat\", @a, @b, @c'\n"
		"  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf >> \"mytap.dat\", @a, @b, @c'\n"
		"  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf > stderr, @a, @b, @c'\n"
		"  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf | \"grep somepattern\", @a, @b, @c'\n"
		"  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf | \"grep somepattern > mytap.dat\", @a, @b, @c'\n"
		"\n"
		"  Please see http://johnkerl.org/miller/doc for more information.\n",
		MLR_GLOBALS.bargv0);
}

static void mlr_dsl_emitp_keyword_usage(FILE* ostream) {
    fprintf(ostream,
		"emitp: inserts an out-of-stream variable into the output record stream.\n"
		"  Hashmap indices present in the data but not slotted by emitp arguments are\n"
		"  output concatenated with \":\".\n"
		"\n"
		"  With >, >>, or |, the data do not become part of the output record stream but\n"
		"  are instead redirected.\n"
		"\n"
		"  The > and >> are for write and append, as in the shell, but (as with awk) the\n"
		"  file-overwrite for > is on first write, not per record. The | is for piping to\n"
		"  a process which will process the data. There will be one open file for each\n"
		"  distinct file name (for > and >>) or one subordinate process for each distinct\n"
		"  value of the piped-to command (for |). Output-formatting flags are taken from\n"
		"  the main command line.\n"
		"\n"
		"  You can use any of the output-format command-line flags, e.g. --ocsv, --ofs,\n"
		"  etc., to control the format of the output if the output is redirected. See also %s -h.\n"
		"\n"
		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp @sums'\n"
		"  Example: mlr --from f.dat put --opprint '@sums[$a][$b]+=$x; emitp > \"tap-\".$a.$b.\".dat\", @sums'\n"
		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp @sums, \"index1\", \"index2\"'\n"
		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp @*, \"index1\", \"index2\"'\n"
		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp >  \"mytap.dat\", @*, \"index1\", \"index2\"'\n"
		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp >> \"mytap.dat\", @*, \"index1\", \"index2\"'\n"
		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp | \"gzip > mytap.dat.gz\", @*, \"index1\", \"index2\"'\n"
		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp > stderr, @*, \"index1\", \"index2\"'\n"
		"  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp | \"grep somepattern\", @*, \"index1\", \"index2\"'\n"
		"\n"
		"  Please see http://johnkerl.org/miller/doc for more information.\n",
		MLR_GLOBALS.bargv0);
}

static void mlr_dsl_end_keyword_usage(FILE* ostream) {
	fprintf(ostream, "end: xxx temp\n");
}

static void mlr_dsl_eprint_keyword_usage(FILE* ostream) {
    fprintf(ostream,
		"eprint: prints expression immediately to stderr.\n"
		"  Example: mlr --from f.dat put -q 'eprint \"The sum of x and y is \".($x+$y)'\n"
		"  Example: mlr --from f.dat put -q 'for (k, v in $*) { eprint k . \" => \" . v }'\n"
		"  Example: mlr --from f.dat put  '(NR %% 1000 == 0) { eprint \"Checkpoint \".NR}'\n");
}

static void mlr_dsl_eprintn_keyword_usage(FILE* ostream) {
    fprintf(ostream,
		"eprintn: prints expression immediately to stderr, without trailing newline.\n"
		"  Example: mlr --from f.dat put -q 'eprintn \"The sum of x and y is \".($x+$y); eprint \"\"'\n");
}

// ----------------------------------------------------------------
static void mlr_dsl_filter_keyword_usage(FILE* ostream) {
    fprintf(ostream,
		"filter: includes/excludes the record in the output record stream.\n"
		"\n"
		"  Example: mlr --from f.dat put 'filter (NR == 2 || $x > 5.4)'\n"
		"\n"
		"  Instead of put with 'filter false' you can simply use put -q.  The following\n"
		"  uses the input record to accumulate data but only prints the running sum\n"
		"  without printing the input record:\n"
		"\n"
		"  Example: mlr --from f.dat put -q '@running_sum += $x * $y; emit @running_sum'\n");
}

static void mlr_dsl_for_keyword_usage(FILE* ostream) {
	fprintf(ostream, "for: xxx temp\n");
}

static void mlr_dsl_func_keyword_usage(FILE* ostream) {
	fprintf(ostream, "func: xxx temp\n");
}

static void mlr_dsl_if_keyword_usage(FILE* ostream) {
	fprintf(ostream, "if: xxx temp\n");
}

static void mlr_dsl_in_keyword_usage(FILE* ostream) {
	fprintf(ostream, "in: xxx temp\n");
}

static void mlr_dsl_local_keyword_usage(FILE* ostream) {
	fprintf(ostream, "local: xxx temp\n");
}

static void mlr_dsl_print_keyword_usage(FILE* ostream) {
    fprintf(ostream,
		"print: prints expression immediately to stdout.\n"
		"  Example: mlr --from f.dat put -q 'print \"The sum of x and y is \".($x+$y)'\n"
		"  Example: mlr --from f.dat put -q 'for (k, v in $*) { print k . \" => \" . v }'\n"
		"  Example: mlr --from f.dat put  '(NR %% 1000 == 0) { print > stderr, \"Checkpoint \".NR}'\n");
}

static void mlr_dsl_printn_keyword_usage(FILE* ostream) {
    fprintf(ostream,
		"printn: prints expression immediately to stdout, without trailing newline.\n"
		"  Example: mlr --from f.dat put -q 'printn \".\"; end { print \"\" }'\n");
}

static void mlr_dsl_return_keyword_usage(FILE* ostream) {
	fprintf(ostream, "return: xxx temp\n");
}

static void mlr_dsl_stderr_keyword_usage(FILE* ostream) {
    fprintf(ostream,
		"stderr: Used for tee, emit, emitf, emitp, print, and dump in place of filename\n"
		"  to print to standard error.\n");
}

static void mlr_dsl_stdout_keyword_usage(FILE* ostream) {
    fprintf(ostream,
		"stdout: Used for tee, emit, emitf, emitp, print, and dump in place of filename\n"
		"  to print to standard output.\n");
}

static void mlr_dsl_subr_keyword_usage(FILE* ostream) {
	fprintf(ostream, "subr: xxx temp\n");
}

static void mlr_dsl_tee_keyword_usage(FILE* ostream) {
    fprintf(ostream,
		"tee: prints the current record to specified file.\n"
		"  This is an immediate print to the specified file (except for pprint format\n"
		"  which of course waits until the end of the input stream to format all output).\n"
		"\n"
		"  The > and >> are for write and append, as in the shell, but (as with awk) the\n"
		"  file-overwrite for > is on first write, not per record. The | is for piping to\n"
		"  a process which will process the data. There will be one open file for each\n"
		"  distinct file name (for > and >>) or one subordinate process for each distinct\n"
		"  value of the piped-to command (for |). Output-formatting flags are taken from\n"
		"  the main command line.\n"
		"\n"
		"  You can use any of the output-format command-line flags, e.g. --ocsv, --ofs,\n"
		"  etc., to control the format of the output. See also %s -h.\n"
		"\n"
		"  Example: mlr --from f.dat put 'tee >  \"/tmp/data-\".$a, $*'\n"
		"  Example: mlr --from f.dat put 'tee >> \"/tmp/data-\".$a.$b, $*'\n"
		"  Example: mlr --from f.dat put 'tee >  stderr, $*'\n"
		"  Example: mlr --from f.dat put -q 'tee | \"tr \[a-z\\] \[A-Z\\]\", $*'\n"
		"  Example: mlr --from f.dat put -q 'tee | \"tr \[a-z\\] \[A-Z\\] > /tmp/data-\".$a, $*'\n"
		"  Example: mlr --from f.dat put -q 'tee | \"gzip > /tmp/data-\".$a.\".gz\", $*'\n"
		"  Example: mlr --from f.dat put -q --ojson 'tee | \"gzip > /tmp/data-\".$a.\".gz\", $*'\n",
		MLR_GLOBALS.bargv0);
}

static void mlr_dsl_unset_keyword_usage(FILE* ostream) {
    fprintf(ostream,
		"unset: clears field(s) from the current record, or an out-of-stream variable.\n"
		"\n"
		"  Example: mlr --from f.dat put 'unset $x'\n"
		"  Example: mlr --from f.dat put 'unset $*'\n"
		"  Example: mlr --from f.dat put 'for (k, v in $*) { if (k =~ \"a.*\") { unset $[k] } }'\n"
		"  Example: mlr --from f.dat put '...; unset @sums'\n"
		"  Example: mlr --from f.dat put '...; unset @sums[\"green\"]'\n"
		"  Example: mlr --from f.dat put '...; unset @*'\n");
}

static void mlr_dsl_while_keyword_usage(FILE* ostream) {
	fprintf(ostream, "while: xxx temp\n");
}
