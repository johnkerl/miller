#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "lib/mlrutil.h"
#include "lib/mtrand.h"
#include "containers/slls.h"
#include "containers/lhmss.h"
#include "containers/lhmsi.h"
#include "input/lrec_readers.h"
#include "mapping/mappers.h"
#include "mapping/lrec_evaluators.h"
#include "output/lrec_writers.h"
#include "cli/mlrcli.h"
#include "cli/quoting.h"
#include "cli/argparse.h"

#ifdef HAVE_CONFIG_H
#include "config.h"
#else
#include "mlrvers.h"
#endif

// ----------------------------------------------------------------
static mapper_setup_t* mapper_lookup_table[] = {
	&mapper_bar_setup,
	&mapper_cat_setup,
	&mapper_check_setup,
	&mapper_count_distinct_setup,
	&mapper_cut_setup,
	&mapper_decimate_setup,
	&mapper_filter_setup,
	&mapper_grep_setup,
	&mapper_group_by_setup,
	&mapper_group_like_setup,
	&mapper_having_fields_setup,
	&mapper_head_setup,
	&mapper_histogram_setup,
	&mapper_join_setup,
	&mapper_label_setup,
	&mapper_put_setup,
	&mapper_regularize_setup,
	&mapper_rename_setup,
	&mapper_reorder_setup,
	&mapper_sample_setup,
	&mapper_sec2gmt_setup,
	&mapper_sort_setup,
	&mapper_stats1_setup,
	&mapper_stats2_setup,
	&mapper_step_setup,
	&mapper_tac_setup,
	&mapper_tail_setup,
	&mapper_top_setup,
	&mapper_uniq_setup,
};
static int mapper_lookup_table_length = sizeof(mapper_lookup_table) / sizeof(mapper_lookup_table[0]);

// ----------------------------------------------------------------
static lhmss_t* singleton_pdesc_to_chars_map = NULL;
static lhmss_t* get_desc_to_chars_map() {
	if (singleton_pdesc_to_chars_map == NULL) {
		singleton_pdesc_to_chars_map = lhmss_alloc();
		lhmss_put(singleton_pdesc_to_chars_map, "cr",        "\r");
		lhmss_put(singleton_pdesc_to_chars_map, "crcr",      "\r\r");
		lhmss_put(singleton_pdesc_to_chars_map, "newline",   "\n");
		lhmss_put(singleton_pdesc_to_chars_map, "lf",        "\n");
		lhmss_put(singleton_pdesc_to_chars_map, "lflf",      "\n\n");
		lhmss_put(singleton_pdesc_to_chars_map, "crlf",      "\r\n");
		lhmss_put(singleton_pdesc_to_chars_map, "crlfcrlf",  "\r\n\r\n");
		lhmss_put(singleton_pdesc_to_chars_map, "tab",       "\t");
		lhmss_put(singleton_pdesc_to_chars_map, "space",     " ");
		lhmss_put(singleton_pdesc_to_chars_map, "comma",     ",");
		lhmss_put(singleton_pdesc_to_chars_map, "newline",   "\n");
		lhmss_put(singleton_pdesc_to_chars_map, "pipe",      "|");
		lhmss_put(singleton_pdesc_to_chars_map, "slash",     "/");
		lhmss_put(singleton_pdesc_to_chars_map, "colon",     ":");
		lhmss_put(singleton_pdesc_to_chars_map, "semicolon", ";");
		lhmss_put(singleton_pdesc_to_chars_map, "equals",    "=");
	}
	return singleton_pdesc_to_chars_map;
}
static char* sep_from_arg(char* arg, char* argv0) {
	char* chars = lhmss_get(get_desc_to_chars_map(), arg);
	if (chars != NULL) // E.g. crlf
		return chars;
	else // E.g. '\r\n'
		return mlr_unbackslash(arg);
}

// ----------------------------------------------------------------
static lhmss_t* singleton_default_rses = NULL;
static lhmss_t* singleton_default_fses = NULL;
static lhmss_t* singleton_default_pses = NULL;
static lhmsi_t* singleton_default_repeat_ifses = NULL;
static lhmsi_t* singleton_default_repeat_ipses = NULL;

static lhmss_t* get_default_rses() {
	if (singleton_default_rses == NULL) {
		singleton_default_rses = lhmss_alloc();
		lhmss_put(singleton_default_rses, "dkvp",    "\n");
		lhmss_put(singleton_default_rses, "nidx",    "\n");
		lhmss_put(singleton_default_rses, "csv",     "\r\n");
		lhmss_put(singleton_default_rses, "csvlite", "\n");
		lhmss_put(singleton_default_rses, "pprint",  "\n");
		lhmss_put(singleton_default_rses, "xtab",    "(N/A)");
	}
	return singleton_default_rses;
}

static lhmss_t* get_default_fses() {
	if (singleton_default_fses == NULL) {
		singleton_default_fses = lhmss_alloc();
		lhmss_put(singleton_default_fses, "dkvp",    ",");
		lhmss_put(singleton_default_fses, "nidx",    " ");
		lhmss_put(singleton_default_fses, "csv",     ",");
		lhmss_put(singleton_default_fses, "csvlite", ",");
		lhmss_put(singleton_default_fses, "pprint",  " ");
		lhmss_put(singleton_default_fses, "xtab",    "\n");
	}
	return singleton_default_fses;
}

static lhmss_t* get_default_pses() {
	if (singleton_default_pses == NULL) {
		singleton_default_pses = lhmss_alloc();
		lhmss_put(singleton_default_pses, "dkvp",    "=");
		lhmss_put(singleton_default_pses, "nidx",    "(N/A)");
		lhmss_put(singleton_default_pses, "csv",     "(N/A)");
		lhmss_put(singleton_default_pses, "csvlite", "(N/A)");
		lhmss_put(singleton_default_pses, "pprint",  "(N/A)");
		lhmss_put(singleton_default_pses, "xtab",    " ");
	}
	return singleton_default_pses;
}

static lhmsi_t* get_default_repeat_ifses() {
	if (singleton_default_repeat_ifses == NULL) {
		singleton_default_repeat_ifses = lhmsi_alloc();
		lhmsi_put(singleton_default_repeat_ifses, "dkvp",    FALSE);
		lhmsi_put(singleton_default_repeat_ifses, "csv",     FALSE);
		lhmsi_put(singleton_default_repeat_ifses, "csvlite", FALSE);
		lhmsi_put(singleton_default_repeat_ifses, "nidx",    FALSE);
		lhmsi_put(singleton_default_repeat_ifses, "xtab",    FALSE);
		lhmsi_put(singleton_default_repeat_ifses, "pprint",  TRUE);
	}
	return singleton_default_repeat_ifses;
}

static lhmsi_t* get_default_repeat_ipses() {
	if (singleton_default_repeat_ipses == NULL) {
		singleton_default_repeat_ipses = lhmsi_alloc();
		lhmsi_put(singleton_default_repeat_ipses, "dkvp",    FALSE);
		lhmsi_put(singleton_default_repeat_ipses, "csv",     FALSE);
		lhmsi_put(singleton_default_repeat_ipses, "csvlite", FALSE);
		lhmsi_put(singleton_default_repeat_ipses, "nidx",    FALSE);
		lhmsi_put(singleton_default_repeat_ipses, "xtab",    TRUE);
		lhmsi_put(singleton_default_repeat_ipses, "pprint",  FALSE);
	}
	return singleton_default_repeat_ipses;
}

// For displaying the default separators in on-line help
static char* rebackslash(char* sep) {
	if (streq(sep, "\r"))
		return "\\r";
	else if (streq(sep, "\n"))
		return "\\n";
	else if (streq(sep, "\r\n"))
		return "\\r\\n";
	else if (streq(sep, "\t"))
		return "\\t";
	else if (streq(sep, " "))
		return "space";
	else
		return sep;
}

// ----------------------------------------------------------------
#define DEFAULT_OFMT "%lf"

#define DEFAULT_OQUOTING QUOTE_MINIMAL

// ----------------------------------------------------------------
// The main_usage() function is split out into subroutines in support of the
// manpage autogenerator.

static void main_usage_synopsis(FILE* o, char* argv0) {
	fprintf(o, "Usage: %s [I/O options] {verb} [verb-dependent options ...] {zero or more file names}\n", argv0);
}

static void main_usage_examples(FILE* o, char* argv0, char* leader) {
	fprintf(o, "%s%s --csv --rs lf --fs tab cut -f hostname,uptime file1.tsv file2.tsv\n", leader, argv0);
	fprintf(o, "%s%s --csv cut -f hostname,uptime mydata.csv\n", leader, argv0);
	fprintf(o, "%s%s --csv filter '$status != \"down\" && $upsec >= 10000' *.csv\n", leader, argv0);
	fprintf(o, "%s%s --nidx put '$sum = $7 + 2.1*$8' *.dat\n", leader, argv0);
	fprintf(o, "%sgrep -v '^#' /etc/group | %s --ifs : --nidx --opprint label group,pass,gid,member then sort -f group\n", leader, argv0);
	fprintf(o, "%s%s join -j account_id -f accounts.dat then group-by account_name balances.dat\n", leader, argv0);
	fprintf(o, "%s%s put '$attr = sub($attr, \"([0-9]+)_([0-9]+)_.*\", \"\\1:\\2\")' data/*\n", leader, argv0);
	fprintf(o, "%s%s stats1 -a min,mean,max,p10,p50,p90 -f flag,u,v data/*\n", leader, argv0);
	fprintf(o, "%s%s stats2 -a linreg-pca -f u,v -g shape data/*\n", leader, argv0);
}

static void list_all_verbs_raw(FILE* o) {
	for (int i = 0; i < mapper_lookup_table_length; i++) {
		fprintf(o, "%s\n", mapper_lookup_table[i]->verb);
	}
}

static void list_all_verbs(FILE* o, char* leader) {
	char* separator = " ";
	int leaderlen = strlen(leader);
	int separatorlen = strlen(separator);
	int linelen = leaderlen;
	int j = 0;
	for (int i = 0; i < mapper_lookup_table_length; i++) {
		char* verb = mapper_lookup_table[i]->verb;
		int verblen = strlen(verb);
		linelen += separatorlen + verblen;
		if (linelen >= 80) {
			fprintf(o, "\n");
			linelen = leaderlen + separatorlen + verblen;
			j = 0;
		}
		if (j == 0)
			fprintf(o, "%s", leader);
		fprintf(o, "%s%s", separator, verb);
		j++;
	}
	fprintf(o, "\n");
}

static void main_usage_help_options(FILE* o, char* argv0) {
	fprintf(o, "  -h or --help Show this message.\n");
	fprintf(o, "  --version              Show the software version.\n");
	fprintf(o, "  {verb name} --help     Show verb-specific help.\n");
	fprintf(o, "  --list-all-verbs or -l List only verb names.\n");
	fprintf(o, "  --help-all-verbs       Show help on all verbs.\n");
}

static void main_usage_functions(FILE* o, char* argv0, char* leader) {
	lrec_evaluator_list_functions(o, leader);
	fprintf(o, "Please use \"%s --help-function {function name}\" for function-specific help.\n", argv0);
	fprintf(o, "Please use \"%s --help-all-functions\" or \"%s -f\" for help on all functions.\n", argv0, argv0);
}

static void main_usage_data_format_examples(FILE* o, char* argv0) {
	fprintf(o,
		"  DKVP: delimited key-value pairs (Miller default format)\n"
		"  +---------------------+\n"
		"  | apple=1,bat=2,cog=3 |  Record 1: \"apple\" => \"1\", \"bat\" => \"2\", \"cog\" => \"3\"\n"
		"  | dish=7,egg=8,flint  |  Record 2: \"dish\" => \"7\", \"egg\" => \"8\", \"3\" => \"flint\"\n"
		"  +---------------------+\n"
		"\n"
		"  NIDX: implicitly numerically indexed (Unix-toolkit style)\n"
		"  +---------------------+\n"
		"  | the quick brown     | Record 1: \"1\" => \"the\", \"2\" => \"quick\", \"3\" => \"brown\"\n"
		"  | fox jumped          | Record 2: \"1\" => \"fox\", \"2\" => \"jumped\"\n"
		"  +---------------------+\n"
		"\n"
		"  CSV/CSV-lite: comma-separated values with separate header line\n"
		"  +---------------------+\n"
		"  | apple,bat,cog       |\n"
		"  | 1,2,3               | Record 1: \"apple => \"1\", \"bat\" => \"2\", \"cog\" => \"3\"\n"
		"  | 4,5,6               | Record 2: \"apple\" => \"4\", \"bat\" => \"5\", \"cog\" => \"6\"\n"
		"  +---------------------+\n"
		"\n"
		"  PPRINT: pretty-printed tabular\n"
		"  +---------------------+\n"
		"  | apple bat cog       |\n"
		"  | 1     2   3         | Record 1: \"apple => \"1\", \"bat\" => \"2\", \"cog\" => \"3\"\n"
		"  | 4     5   6         | Record 2: \"apple\" => \"4\", \"bat\" => \"5\", \"cog\" => \"6\"\n"
		"  +---------------------+\n"
		"\n"
		"  XTAB: pretty-printed transposed tabular\n"
		"  +---------------------+\n"
		"  | apple 1             | Record 1: \"apple\" => \"1\", \"bat\" => \"2\", \"cog\" => \"3\"\n"
		"  | bat   2             |\n"
		"  | cog   3             |\n"
		"  |                     |\n"
		"  | dish 7              | Record 2: \"dish\" => \"7\", \"egg\" => \"8\"\n"
		"  | egg  8              |\n"
		"  +---------------------+\n");
}

static void main_usage_data_format_options(FILE* o, char* argv0) {
	fprintf(o, "  --idkvp   --odkvp   --dkvp            Delimited key-value pairs, e.g \"a=1,b=2\"\n");
	fprintf(o, "                                        (default)\n");
	fprintf(o, "  --inidx   --onidx   --nidx            Implicitly-integer-indexed fields\n");
	fprintf(o, "                                        (Unix-toolkit style)\n");
	fprintf(o, "  --icsv    --ocsv    --csv             Comma-separated value (or tab-separated\n");
	fprintf(o, "                                        with --fs tab, etc.)\n");
	fprintf(o, "  --ipprint --opprint --pprint --right  Pretty-printed tabular (produces no\n");
	fprintf(o, "                                        output until all input is in)\n");
	fprintf(o, "  --ixtab   --oxtab   --xtab --xvright  Pretty-printed vertical-tabular\n");
	fprintf(o, "  The --right option right-justifies all fields for PPRINT output format.\n");
	fprintf(o, "  The --xvright option right-justifies values for XTAB format.\n");
	fprintf(o, "  -p is a keystroke-saver for --nidx --fs space --repifs\n");
	fprintf(o, "  Examples: --csv for CSV-formatted input and output; --idkvp --opprint for\n");
	fprintf(o, "  DKVP-formatted input and pretty-printed output.\n");
}

static void main_usage_compressed_data_options(FILE* o, char* argv0) {
	fprintf(o, "  --prepipe {command} This allows Miller to handle compressed inputs. You can do\n");
	fprintf(o, "  without this for single input files, e.g. \"gunzip < myfile.csv.gz | %s ...\".\n",
		argv0);
	fprintf(o, "  However, when multiple input files are present, between-file separations are\n");
	fprintf(o, "  lost; also, the FILENAME variable doesn't iterate. Using --prepipe you can\n");
	fprintf(o, "  specify an action to be taken on each input file. This pre-pipe command must\n");
	fprintf(o, "  be able to read from standard input; it will be invoked with\n");
	fprintf(o, "    {command} < {filename}.\n");
	fprintf(o, "  Examples:\n");
	fprintf(o, "    %s --prepipe 'gunzip'\n", argv0);
	fprintf(o, "    %s --prepipe 'zcat -cf'\n", argv0);
	fprintf(o, "    %s --prepipe 'xz -cd'\n", argv0);
	fprintf(o, "    %s --prepipe cat\n", argv0);
	fprintf(o, "  Note that this feature is quite general and is not limited to decompression\n");
	fprintf(o, "  utilities. You can use it to apply per-file filters of your choice.\n");
	fprintf(o, "  For output compression (or other) utilities, simply pipe the output:\n");
	fprintf(o, "    %s ... | {your compression command}\n", argv0);
}

static void main_usage_separator_options(FILE* o, char* argv0) {
	fprintf(o, "  --rs     --irs     --ors              Record separators, e.g. 'lf' or '\\r\\n'\n");
	fprintf(o, "  --fs     --ifs     --ofs  --repifs    Field separators, e.g. comma\n");
	fprintf(o, "  --ps     --ips     --ops              Pair separators, e.g. equals sign\n");
	fprintf(o, "  Notes:\n");
	fprintf(o, "  * IPS/OPS are only used for DKVP and XTAB formats, since only in these formats\n");
	fprintf(o, "    do key-value pairs appear juxtaposed.\n");
	fprintf(o, "  * IRS/ORS are ignored for XTAB format. Nominally IFS and OFS are newlines;\n");
	fprintf(o, "    XTAB records are separated by two or more consecutive IFS/OFS -- i.e.\n");
	fprintf(o, "    a blank line.\n");
	fprintf(o, "  * OFS must be single-character for PPRINT format. This is because it is used\n");
	fprintf(o, "    with repetition for alignment; multi-character separators would make\n");
	fprintf(o, "    alignment impossible.\n");
	fprintf(o, "  * OPS may be multi-character for XTAB format, in which case alignment is\n");
	fprintf(o, "    disabled.\n");
	fprintf(o, "  * DKVP, NIDX, CSVLITE, PPRINT, and XTAB formats are intended to handle\n");
	fprintf(o, "    platform-native text data. In particular, this means LF line-terminators\n");
	fprintf(o, "    by default on Linux/OSX. You can use \"--dkvp --rs crlf\" for\n");
	fprintf(o, "    CRLF-terminated DKVP files, and so on.\n");
	fprintf(o, "  * CSV is intended to handle RFC-4180-compliant data. In particular, this means\n");
	fprintf(o, "    it uses CRLF line-terminators by default. You can use \"--csv --rs lf\" for\n");
	fprintf(o, "    Linux-native CSV files.\n");
	fprintf(o, "  * You can specify separators in any of the following ways, shown by example:\n");
	fprintf(o, "    - Type them out, quoting as necessary for shell escapes, e.g.\n");
	fprintf(o, "      \"--fs '|' --ips :\"\n");
	fprintf(o, "    - C-style escape sequences, e.g. \"--rs '\\r\\n' --fs '\\t'\".\n");
	fprintf(o, "    - To avoid backslashing, you can use any of the following names:\n");
	fprintf(o, "     ");
	lhmss_t* pmap = get_desc_to_chars_map();
	for (lhmsse_t* pe = pmap->phead; pe != NULL; pe = pe->pnext) {
		fprintf(o, " %s", pe->key);
	}
	fprintf(o, "\n");
	fprintf(o, "  * Default separators by format:\n");
	fprintf(o, "      %-12s %-8s %-8s %s\n", "File format", "RS", "FS", "PS");
	lhmss_t* default_rses = get_default_rses();
	lhmss_t* default_fses = get_default_fses();
	lhmss_t* default_pses = get_default_pses();
	for (lhmsse_t* pe = default_rses->phead; pe != NULL; pe = pe->pnext) {
		char* filefmt = pe->key;
		char* rs = pe->value;
		char* fs = lhmss_get(default_fses, filefmt);
		char* ps = lhmss_get(default_pses, filefmt);
		fprintf(o, "      %-12s %-8s %-8s %s\n", filefmt, rebackslash(rs), rebackslash(fs), rebackslash(ps));
	}
}

static void main_usage_csv_options(FILE* o, char* argv0) {
	fprintf(o, "  --implicit-csv-header Use 1,2,3,... as field labels, rather than from line 1\n");
	fprintf(o, "                     of input files. Tip: combine with \"label\" to recreate\n");
	fprintf(o, "                     missing headers.\n");
	fprintf(o, "  --headerless-csv-output   Print only CSV data lines.\n");
}

static void main_usage_double_quoting(FILE* o, char* argv0) {
	fprintf(o, "  --quote-all        Wrap all fields in double quotes\n");
	fprintf(o, "  --quote-none       Do not wrap any fields in double quotes, even if they have\n");
	fprintf(o, "                     OFS or ORS in them\n");
	fprintf(o, "  --quote-minimal    Wrap fields in double quotes only if they have OFS or ORS\n");
	fprintf(o, "                     in them (default)\n");
	fprintf(o, "  --quote-numeric    Wrap fields in double quotes only if they have numbers\n");
	fprintf(o, "                     in them\n");
}

static void main_usage_numerical_formatting(FILE* o, char* argv0) {
	fprintf(o, "  --ofmt {format}    E.g. %%.18lf, %%.0lf. Please use sprintf-style codes for\n");
	fprintf(o, "                     double-precision. Applies to verbs which compute new\n");
	fprintf(o, "                     values, e.g. put, stats1, stats2. See also the fmtnum\n");
	fprintf(o, "                     function within mlr put (mlr --help-all-functions).\n");
	fprintf(o, "                     Defaults to %s.\n", DEFAULT_OFMT);
}

static void main_usage_other_options(FILE* o, char* argv0) {
	fprintf(o, "  --seed {n} with n of the form 12345678 or 0xcafefeed. For put/filter\n");
	fprintf(o, "                     urand()/urandint()/urand32().\n");
}

static void main_usage_then_chaining(FILE* o, char* argv0) {
	fprintf(o, "Output of one verb may be chained as input to another using \"then\", e.g.\n");
	fprintf(o, "  %s stats1 -a min,mean,max -f flag,u,v -g color then sort -f color\n", argv0);
}

static void main_usage_see_also(FILE* o, char* argv0) {
	fprintf(o, "For more information please see http://johnkerl.org/miller/doc and/or\n");
	fprintf(o, "http://github.com/johnkerl/miller.");
#ifdef HAVE_CONFIG_H
	fprintf(o, " This is Miller version %s.\n", PACKAGE_VERSION);
#else
	fprintf(o, " This is Miller version %s.\n", MLR_VERSION);
#endif // HAVE_CONFIG_H
}

// ----------------------------------------------------------------
static void main_usage(FILE* o, char* argv0) {
	main_usage_synopsis(o, argv0);
	fprintf(o, "\n");

	fprintf(o, "Command-line-syntax examples:\n");
	main_usage_examples(o, argv0, "  ");
	fprintf(o, "\n");

	fprintf(o, "Data-format examples:\n");
	main_usage_data_format_examples(o, argv0);
	fprintf(o, "\n");

	fprintf(o, "Help options:\n");
	main_usage_help_options(o, argv0);
	fprintf(o, "\n");

	fprintf(o, "Verbs:\n");
	list_all_verbs(o, "  ");
	fprintf(o, "\n");

	fprintf(o, "Functions for the filter and put verbs:\n");
	main_usage_functions(o, argv0, "  ");
	fprintf(o, "\n");

	fprintf(o, "Data-format options, for input, output, or both:\n");
	main_usage_data_format_options(o, argv0);
	fprintf(o, "\n");

	fprintf(o, "Compressed-data options:\n");
	main_usage_compressed_data_options(o, argv0);
	fprintf(o, "\n");

	fprintf(o, "Separator options, for input, output, or both:\n");
	main_usage_separator_options(o, argv0);
	fprintf(o, "\n");

	fprintf(o, "Relevant to CSV/CSV-lite input only:\n");
	main_usage_csv_options(o, argv0);
	fprintf(o, "\n");

	fprintf(o, "Double-quoting for CSV output:\n");
	main_usage_double_quoting(o, argv0);
	fprintf(o, "\n");

	fprintf(o, "Numerical formatting:\n");
	main_usage_numerical_formatting(o, argv0);
	fprintf(o, "\n");

	fprintf(o, "Other options:\n");
	main_usage_other_options(o, argv0);
	fprintf(o, "\n");

	fprintf(o, "Then-chaining:\n");
	main_usage_then_chaining(o, argv0);
	fprintf(o, "\n");

	main_usage_see_also(o, argv0);
}

// ----------------------------------------------------------------
static void usage_all_verbs(char* argv0) {
	char* separator = "================================================================";

	for (int i = 0; i < mapper_lookup_table_length; i++) {
		fprintf(stdout, "%s\n", separator);
		mapper_lookup_table[i]->pusage_func(stdout, argv0, mapper_lookup_table[i]->verb);
		fprintf(stdout, "\n");
	}
	fprintf(stdout, "%s\n", separator);
	exit(0);
}

static void usage_unrecognized_verb(char* argv0, char* arg) {
	fprintf(stderr, "%s: option \"%s\" not recognized.\n", argv0, arg);
	fprintf(stderr, "\n");
	main_usage(stderr, argv0);
	exit(1);
}

static void check_arg_count(char** argv, int argi, int argc, int n) {
	if ((argc - argi) < n) {
		main_usage(stderr, argv[0]);
		exit(1);
	}
}

static mapper_setup_t* look_up_mapper_setup(char* verb) {
	mapper_setup_t* pmapper_setup = NULL;
	for (int i = 0; i < mapper_lookup_table_length; i++) {
		if (streq(mapper_lookup_table[i]->verb, verb))
			return mapper_lookup_table[i];
	}

	return pmapper_setup;
}

// ----------------------------------------------------------------
cli_opts_t* parse_command_line(int argc, char** argv) {
	cli_opts_t* popts = mlr_malloc_or_die(sizeof(cli_opts_t));
	memset(popts, 0, sizeof(*popts));

	popts->irs               = NULL;
	popts->ifs               = NULL;
	popts->ips               = NULL;
	popts->allow_repeat_ifs  = NEITHER_TRUE_NOR_FALSE;
	popts->allow_repeat_ips  = NEITHER_TRUE_NOR_FALSE;
	popts->use_implicit_csv_header = FALSE;
	popts->headerless_csv_output   = FALSE;

	popts->ors               = NULL;
	popts->ofs               = NULL;
	popts->ops               = NULL;
	popts->right_justify_xtab_value = FALSE;
	popts->ofmt              = DEFAULT_OFMT;
	popts->oquoting          = DEFAULT_OQUOTING;

	popts->plrec_reader      = NULL;
	popts->plrec_writer      = NULL;

	popts->prepipe           = NULL;
	popts->filenames         = NULL;

	popts->ifile_fmt         = "dkvp";
	popts->ofile_fmt         = "dkvp";

	popts->use_mmap_for_read = TRUE;
	int left_align_pprint    = TRUE;

	int have_rand_seed       = FALSE;
	unsigned rand_seed       = 0;

	int argi = 1;
	for (; argi < argc; argi++) {
		if (argv[argi][0] != '-') {
			break;

		} else if (streq(argv[argi], "--version")) {
#ifdef HAVE_CONFIG_H
			printf("Miller %s\n", PACKAGE_VERSION);
#else
			printf("Miller %s\n", MLR_VERSION);
#endif // HAVE_CONFIG_H
			exit(0);
		} else if (streq(argv[argi], "-h")) {
			main_usage(stdout, argv[0]);
			exit(0);
		} else if (streq(argv[argi], "--help")) {
			main_usage(stdout, argv[0]);
			exit(0);
		} else if (streq(argv[argi], "--help-all-verbs")) {
			usage_all_verbs(argv[0]);
		} else if (streq(argv[argi], "--list-all-verbs") || streq(argv[argi], "-l")) {
			list_all_verbs(stdout, "");
			exit(0);
		} else if (streq(argv[argi], "--list-all-verbs-raw")) {
			list_all_verbs_raw(stdout);
			exit(0);
		} else if (streq(argv[argi], "--list-all-functions-raw")) {
			lrec_evaluator_list_all_functions_raw(stdout);
			exit(0);
		} else if (streq(argv[argi], "--help-all-functions") || streq(argv[argi], "-f")) {
			lrec_evaluator_function_usage(stdout, NULL);
			exit(0);

		} else if (streq(argv[argi], "--help-function") || streq(argv[argi], "--hf")) {
			check_arg_count(argv, argi, argc, 2);
			lrec_evaluator_function_usage(stdout, argv[argi+1]);
			exit(0);

		// main-usage subsections, individually accessible for the benefit of
		// the manpage-autogenerator
		} else if (streq(argv[argi], "--usage-synopsis")) {
			main_usage_synopsis(stdout, argv[0]);
			exit(0);
		} else if (streq(argv[argi], "--usage-examples")) {
			main_usage_examples(stdout, argv[0], "");
			exit(0);
		} else if (streq(argv[argi], "--usage-list-all-verbs")) {
			list_all_verbs(stdout, "");
			exit(0);
		} else if (streq(argv[argi], "--usage-help-options")) {
			main_usage_help_options(stdout, argv[0]);
			exit(0);
		} else if (streq(argv[argi], "--usage-functions")) {
			main_usage_functions(stdout, argv[0], "");
			exit(0);
		} else if (streq(argv[argi], "--usage-data-format-examples")) {
			main_usage_data_format_examples(stdout, argv[0]);
			exit(0);
		} else if (streq(argv[argi], "--usage-data-format-options")) {
			main_usage_data_format_options(stdout, argv[0]);
			exit(0);
		} else if (streq(argv[argi], "--usage-compressed-data-options")) {
			main_usage_compressed_data_options(stdout, argv[0]);
			exit(0);
		} else if (streq(argv[argi], "--usage-separator-options")) {
			main_usage_separator_options(stdout, argv[0]);
			exit(0);
		} else if (streq(argv[argi], "--usage-csv-options")) {
			main_usage_csv_options(stdout, argv[0]);
			exit(0);
		} else if (streq(argv[argi], "--usage-double-quoting")) {
			main_usage_double_quoting(stdout, argv[0]);
			exit(0);
		} else if (streq(argv[argi], "--usage-numerical-formatting")) {
			main_usage_numerical_formatting(stdout, argv[0]);
			exit(0);
		} else if (streq(argv[argi], "--usage-other-options")) {
			main_usage_other_options(stdout, argv[0]);
			exit(0);
		} else if (streq(argv[argi], "--usage-then-chaining")) {
			main_usage_then_chaining(stdout, argv[0]);
			exit(0);
		} else if (streq(argv[argi], "--usage-see-also")) {
			main_usage_see_also(stdout, argv[0]);
			exit(0);

		} else if (streq(argv[argi], "--rs")) {
			check_arg_count(argv, argi, argc, 2);
			popts->ors = sep_from_arg(argv[argi+1], argv[0]);
			popts->irs = sep_from_arg(argv[argi+1], argv[0]);
			argi++;
		} else if (streq(argv[argi], "--irs")) {
			check_arg_count(argv, argi, argc, 2);
			popts->irs = sep_from_arg(argv[argi+1], argv[0]);
			argi++;
		} else if (streq(argv[argi], "--ors")) {
			check_arg_count(argv, argi, argc, 2);
			popts->ors = sep_from_arg(argv[argi+1], argv[0]);
			argi++;

		} else if (streq(argv[argi], "--fs")) {
			check_arg_count(argv, argi, argc, 2);
			popts->ofs = sep_from_arg(argv[argi+1], argv[0]);
			popts->ifs = sep_from_arg(argv[argi+1], argv[0]);
			argi++;
		} else if (streq(argv[argi], "--ifs")) {
			check_arg_count(argv, argi, argc, 2);
			popts->ifs = sep_from_arg(argv[argi+1], argv[0]);
			argi++;
		} else if (streq(argv[argi], "--ofs")) {
			check_arg_count(argv, argi, argc, 2);
			popts->ofs = sep_from_arg(argv[argi+1], argv[0]);
			argi++;
		} else if (streq(argv[argi], "--repifs")) {
			popts->allow_repeat_ifs = TRUE;
		} else if (streq(argv[argi], "--implicit-csv-header")) {
			popts->use_implicit_csv_header = TRUE;
		} else if (streq(argv[argi], "--headerless-csv-output")) {
			popts->headerless_csv_output = TRUE;

		} else if (streq(argv[argi], "-p")) {
			popts->ifile_fmt = "nidx";
			popts->ofile_fmt = "nidx";
			popts->ifs = " ";
			popts->ofs = " ";
			popts->allow_repeat_ifs = TRUE;

		} else if (streq(argv[argi], "--ps")) {
			check_arg_count(argv, argi, argc, 2);
			popts->ops = sep_from_arg(argv[argi+1], argv[0]);
			popts->ips = sep_from_arg(argv[argi+1], argv[0]);
			argi++;
		} else if (streq(argv[argi], "--ips")) {
			check_arg_count(argv, argi, argc, 2);
			popts->ips = sep_from_arg(argv[argi+1], argv[0]);
			argi++;
		} else if (streq(argv[argi], "--ops")) {
			check_arg_count(argv, argi, argc, 2);
			popts->ops = sep_from_arg(argv[argi+1], argv[0]);
			argi++;

		} else if (streq(argv[argi], "--xvright")) {
			popts->right_justify_xtab_value = TRUE;

		} else if (streq(argv[argi], "--csv"))      { popts->ifile_fmt = popts->ofile_fmt = "csv";
		} else if (streq(argv[argi], "--icsv"))     { popts->ifile_fmt = "csv";
		} else if (streq(argv[argi], "--ocsv"))     { popts->ofile_fmt = "csv";

		} else if (streq(argv[argi], "--csvlite"))  { popts->ifile_fmt = popts->ofile_fmt = "csvlite";
		} else if (streq(argv[argi], "--icsvlite")) { popts->ifile_fmt = "csvlite";
		} else if (streq(argv[argi], "--ocsvlite")) { popts->ofile_fmt = "csvlite";

		} else if (streq(argv[argi], "--dkvp"))     { popts->ifile_fmt = popts->ofile_fmt = "dkvp";
		} else if (streq(argv[argi], "--idkvp"))    { popts->ifile_fmt = "dkvp";
		} else if (streq(argv[argi], "--odkvp"))    { popts->ofile_fmt = "dkvp";

		} else if (streq(argv[argi], "--nidx"))     { popts->ifile_fmt = popts->ofile_fmt = "nidx";
		} else if (streq(argv[argi], "--inidx"))    { popts->ifile_fmt = "nidx";
		} else if (streq(argv[argi], "--onidx"))    { popts->ofile_fmt = "nidx";

		} else if (streq(argv[argi], "--xtab"))     { popts->ifile_fmt = popts->ofile_fmt = "xtab";
		} else if (streq(argv[argi], "--ixtab"))    { popts->ifile_fmt = "xtab";
		} else if (streq(argv[argi], "--oxtab"))    { popts->ofile_fmt = "xtab";

		} else if (streq(argv[argi], "--ipprint")) {
			popts->ifile_fmt        = "csvlite";
			popts->ifs              = " ";
			popts->allow_repeat_ifs = TRUE;

		} else if (streq(argv[argi], "--opprint")) {
			popts->ofile_fmt = "pprint";
		} else if (streq(argv[argi], "--pprint")) {
			popts->ifile_fmt        = "csvlite";
			popts->ifs              = " ";
			popts->allow_repeat_ifs = TRUE;
			popts->ofile_fmt        = "pprint";
		} else if (streq(argv[argi], "--right"))   {
			left_align_pprint = FALSE;

		} else if (streq(argv[argi], "--ofmt")) {
			check_arg_count(argv, argi, argc, 2);
			popts->ofmt = argv[argi+1];
			argi++;

		} else if (streq(argv[argi], "--quote-all"))     { popts->oquoting = QUOTE_ALL;
		} else if (streq(argv[argi], "--quote-none"))    { popts->oquoting = QUOTE_NONE;
		} else if (streq(argv[argi], "--quote-minimal")) { popts->oquoting = QUOTE_MINIMAL;
		} else if (streq(argv[argi], "--quote-numeric")) { popts->oquoting = QUOTE_NUMERIC;

		} else if (streq(argv[argi], "--mmap")) {
			popts->use_mmap_for_read = TRUE;
		} else if (streq(argv[argi], "--no-mmap")) {
			popts->use_mmap_for_read = FALSE;

		} else if (streq(argv[argi], "--seed")) {
			check_arg_count(argv, argi, argc, 2);
			if (sscanf(argv[argi+1], "0x%x", &rand_seed) == 1) {
				have_rand_seed = TRUE;
			} else if (sscanf(argv[argi+1], "%u", &rand_seed) == 1) {
				have_rand_seed = TRUE;
			} else {
				main_usage(stderr, argv[0]);
				exit(1);
			}
			argi++;

		} else if (streq(argv[argi], "--prepipe")) {
			check_arg_count(argv, argi, argc, 2);
			popts->prepipe = argv[argi+1];
			popts->use_mmap_for_read = FALSE;
			argi++;

		} else {
			usage_unrecognized_verb(argv[0], argv[argi]);
		}
	}

	lhmss_t* default_rses = get_default_rses();
	lhmss_t* default_fses = get_default_fses();
	lhmss_t* default_pses = get_default_pses();
	lhmsi_t* default_repeat_ifses = get_default_repeat_ifses();
	lhmsi_t* default_repeat_ipses = get_default_repeat_ipses();

	if (popts->irs == NULL)
		popts->irs = lhmss_get(default_rses, popts->ifile_fmt);
	if (popts->ifs == NULL)
		popts->ifs = lhmss_get(default_fses, popts->ifile_fmt);
	if (popts->ips == NULL)
		popts->ips = lhmss_get(default_pses, popts->ifile_fmt);

	if (popts->allow_repeat_ifs == NEITHER_TRUE_NOR_FALSE)
		popts->allow_repeat_ifs = lhmsi_get(default_repeat_ifses, popts->ifile_fmt);
	if (popts->allow_repeat_ips == NEITHER_TRUE_NOR_FALSE)
		popts->allow_repeat_ips = lhmsi_get(default_repeat_ipses, popts->ifile_fmt);

	if (popts->ors == NULL)
		popts->ors = lhmss_get(default_rses, popts->ofile_fmt);
	if (popts->ofs == NULL)
		popts->ofs = lhmss_get(default_fses, popts->ofile_fmt);
	if (popts->ops == NULL)
		popts->ops = lhmss_get(default_pses, popts->ofile_fmt);

	if (popts->irs == NULL) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n", argv[0], __FILE__, __LINE__);
		exit(1);
	}
	if (popts->ifs == NULL) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n", argv[0], __FILE__, __LINE__);
		exit(1);
	}
	if (popts->ips == NULL) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n", argv[0], __FILE__, __LINE__);
		exit(1);
	}

	if (popts->allow_repeat_ifs == NEITHER_TRUE_NOR_FALSE) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n", argv[0], __FILE__, __LINE__);
		exit(1);
	}
	if (popts->allow_repeat_ips == NEITHER_TRUE_NOR_FALSE) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n", argv[0], __FILE__, __LINE__);
		exit(1);
	}

	if (popts->ors == NULL) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n", argv[0], __FILE__, __LINE__);
		exit(1);
	}
	if (popts->ofs == NULL) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n", argv[0], __FILE__, __LINE__);
		exit(1);
	}
	if (popts->ops == NULL) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n", argv[0], __FILE__, __LINE__);
		exit(1);
	}

	if (streq(popts->ofile_fmt, "pprint") && strlen(popts->ofs) != 1) {
		fprintf(stderr, "%s: OFS for PPRINT format must be single-character; got \"%s\".\n",
			argv[0], popts->ofs);
		return NULL;
	}
	if      (streq(popts->ofile_fmt, "dkvp"))
		popts->plrec_writer = lrec_writer_dkvp_alloc(popts->ors, popts->ofs, popts->ops);
	else if (streq(popts->ofile_fmt, "csv"))
		popts->plrec_writer = lrec_writer_csv_alloc(popts->ors, popts->ofs, popts->oquoting,
			popts->headerless_csv_output);
	else if (streq(popts->ofile_fmt, "csvlite"))
		popts->plrec_writer = lrec_writer_csvlite_alloc(popts->ors, popts->ofs, popts->headerless_csv_output);
	else if (streq(popts->ofile_fmt, "nidx"))
		popts->plrec_writer = lrec_writer_nidx_alloc(popts->ors, popts->ofs);
	else if (streq(popts->ofile_fmt, "xtab"))
		popts->plrec_writer = lrec_writer_xtab_alloc(popts->ofs, popts->ops, popts->right_justify_xtab_value);
	else if (streq(popts->ofile_fmt, "pprint"))
		popts->plrec_writer = lrec_writer_pprint_alloc(popts->ors, popts->ofs[0], left_align_pprint);
	else {
		main_usage(stderr, argv[0]);
		exit(1);
	}

	if ((argc - argi) < 1) {
		main_usage(stderr, argv[0]);
		exit(1);
	}

	popts->pmapper_list = sllv_alloc();
	while (TRUE) {
		check_arg_count(argv, argi, argc, 1);
		char* verb = argv[argi];

		mapper_setup_t* pmapper_setup = look_up_mapper_setup(verb);
		if (pmapper_setup == NULL) {
			fprintf(stderr, "%s: verb \"%s\" not found. Please use \"%s --help\" for a list.\n",
				argv[0], verb, argv[0]);
			exit(1);
		}

		if ((argc - argi) >= 2) {
			if (streq(argv[argi+1], "-h") || streq(argv[argi+1], "--help")) {
				pmapper_setup->pusage_func(stdout, argv[0], verb);
				exit(0);
			}
		}

		// It's up to the parse func to print its usage on CLI-parse failure.
		mapper_t* pmapper = pmapper_setup->pparse_func(&argi, argc, argv);
		if (pmapper == NULL) {
			exit(1);
		}
		sllv_add(popts->pmapper_list, pmapper);

		if (argi >= argc || !streq(argv[argi], "then"))
			break;
		argi++;
	}

	popts->filenames = &argv[argi];

	// No filenames means read from standard input, and standard input cannot be mmapped.
	if (argi == argc)
		popts->use_mmap_for_read = FALSE;

	popts->plrec_reader = lrec_reader_alloc(popts->ifile_fmt, popts->use_mmap_for_read,
		popts->irs, popts->ifs, popts->allow_repeat_ifs, popts->ips, popts->allow_repeat_ips,
		popts->use_implicit_csv_header);
	if (popts->plrec_reader == NULL) {
		main_usage(stderr, argv[0]);
		exit(1);
	}

	if (have_rand_seed) {
		mtrand_init(rand_seed);
	} else {
		mtrand_init_default();
	}

	return popts;
}

// ----------------------------------------------------------------
void cli_opts_free(cli_opts_t* popts) {
	if (popts == NULL)
		return;

	popts->plrec_reader->pfree_func(popts->plrec_reader, popts->plrec_reader->pvstate);

	for (sllve_t* pe = popts->pmapper_list->phead; pe != NULL; pe = pe->pnext) {
		mapper_t* pmapper = pe->pvdata;
		pmapper->pfree_func(pmapper->pvstate);
	}
	sllv_free(popts->pmapper_list);

	popts->plrec_writer->pfree_func(popts->plrec_writer->pvstate);
	free(popts);
}
