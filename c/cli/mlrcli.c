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
#include "mapping/rval_evaluators.h"
#include "mapping/mlr_dsl_cst.h"
#include "output/lrec_writers.h"
#include "cli/mlrcli.h"
#include "cli/quoting.h"
#include "cli/argparse.h"

#ifdef HAVE_CONFIG_H
#include "config.h"
#define VERSION_STRING PACKAGE_VERSION
#else
#include "mlrvers.h"
#define VERSION_STRING MLR_VERSION
#endif

// ----------------------------------------------------------------
#define DEFAULT_OFMT                     "%lf"
#define DEFAULT_OQUOTING                 QUOTE_MINIMAL
#define DEFAULT_JSON_FLATTEN_SEPARATOR   ":"
#define DEFAULT_OOSVAR_FLATTEN_SEPARATOR ":"

// ----------------------------------------------------------------
static mapper_setup_t* mapper_lookup_table[] = {
	&mapper_bar_setup,
	&mapper_bootstrap_setup,
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
	&mapper_merge_fields_setup,
	&mapper_nest_setup,
	&mapper_nothing_setup,
	&mapper_put_setup,
	&mapper_regularize_setup,
	&mapper_rename_setup,
	&mapper_reorder_setup,
	&mapper_repeat_setup,
	&mapper_reshape_setup,
	&mapper_sample_setup,
	&mapper_sec2gmt_setup,
	&mapper_sec2gmtdate_setup,
	&mapper_shuffle_setup,
	&mapper_sort_setup,
	&mapper_stats1_setup,
	&mapper_stats2_setup,
	&mapper_step_setup,
	&mapper_tac_setup,
	&mapper_tail_setup,
	&mapper_tee_setup,
	&mapper_top_setup,
	&mapper_uniq_setup,
};
static int mapper_lookup_table_length = sizeof(mapper_lookup_table) / sizeof(mapper_lookup_table[0]);

// ----------------------------------------------------------------
static lhmss_t* singleton_pdesc_to_chars_map = NULL;
static lhmss_t* get_desc_to_chars_map() {
	if (singleton_pdesc_to_chars_map == NULL) {
		singleton_pdesc_to_chars_map = lhmss_alloc();
		lhmss_put(singleton_pdesc_to_chars_map, "cr",        "\r",       NO_FREE);
		lhmss_put(singleton_pdesc_to_chars_map, "crcr",      "\r\r",     NO_FREE);
		lhmss_put(singleton_pdesc_to_chars_map, "newline",   "\n",       NO_FREE);
		lhmss_put(singleton_pdesc_to_chars_map, "lf",        "\n",       NO_FREE);
		lhmss_put(singleton_pdesc_to_chars_map, "lflf",      "\n\n",     NO_FREE);
		lhmss_put(singleton_pdesc_to_chars_map, "crlf",      "\r\n",     NO_FREE);
		lhmss_put(singleton_pdesc_to_chars_map, "crlfcrlf",  "\r\n\r\n", NO_FREE);
		lhmss_put(singleton_pdesc_to_chars_map, "tab",       "\t",       NO_FREE);
		lhmss_put(singleton_pdesc_to_chars_map, "space",     " ",        NO_FREE);
		lhmss_put(singleton_pdesc_to_chars_map, "comma",     ",",        NO_FREE);
		lhmss_put(singleton_pdesc_to_chars_map, "newline",   "\n",       NO_FREE);
		lhmss_put(singleton_pdesc_to_chars_map, "pipe",      "|",        NO_FREE);
		lhmss_put(singleton_pdesc_to_chars_map, "slash",     "/",        NO_FREE);
		lhmss_put(singleton_pdesc_to_chars_map, "colon",     ":",        NO_FREE);
		lhmss_put(singleton_pdesc_to_chars_map, "semicolon", ";",        NO_FREE);
		lhmss_put(singleton_pdesc_to_chars_map, "equals",    "=",        NO_FREE);
	}
	return singleton_pdesc_to_chars_map;
}
// Always strdup so the caller can unconditionally free our return value
char* cli_sep_from_arg(char* arg) {
	char* chars = lhmss_get(get_desc_to_chars_map(), arg);
	if (chars != NULL) // E.g. crlf
		return mlr_strdup_or_die(chars);
	else // E.g. '\r\n'
		return mlr_alloc_unbackslash(arg);
}

// ----------------------------------------------------------------
static lhmss_t* singleton_default_rses = NULL;
static lhmss_t* singleton_default_fses = NULL;
static lhmss_t* singleton_default_pses = NULL;
static lhmsi_t* singleton_default_repeat_ifses = NULL;
static lhmsi_t* singleton_default_repeat_ipses = NULL;

lhmss_t* get_default_rses() {
	if (singleton_default_rses == NULL) {
		singleton_default_rses = lhmss_alloc();
		lhmss_put(singleton_default_rses, "dkvp",     "\n",    NO_FREE);
		lhmss_put(singleton_default_rses, "json",     "(N/A)", NO_FREE);
		lhmss_put(singleton_default_rses, "nidx",     "\n",    NO_FREE);
		lhmss_put(singleton_default_rses, "csv",      "\r\n",  NO_FREE);

		char* csv_rs = "\r\n";
		char* env_default = getenv("MLR_CSV_DEFAULT_RS");
		if (env_default != NULL && !streq(env_default, ""))
			csv_rs = cli_sep_from_arg(env_default);
		lhmss_put(singleton_default_rses, "csv", csv_rs, NO_FREE);

		lhmss_put(singleton_default_rses, "csvlite",  "\n",    NO_FREE);
		lhmss_put(singleton_default_rses, "markdown", "\n",    NO_FREE);
		lhmss_put(singleton_default_rses, "pprint",   "\n",    NO_FREE);
		lhmss_put(singleton_default_rses, "xtab",     "(N/A)", NO_FREE);
	}
	return singleton_default_rses;
}

lhmss_t* get_default_fses() {
	if (singleton_default_fses == NULL) {
		singleton_default_fses = lhmss_alloc();
		lhmss_put(singleton_default_fses, "dkvp",     ",",      NO_FREE);
		lhmss_put(singleton_default_fses, "json",     "(N/A)",  NO_FREE);
		lhmss_put(singleton_default_fses, "nidx",     " ",      NO_FREE);
		lhmss_put(singleton_default_fses, "csv",      ",",      NO_FREE);
		lhmss_put(singleton_default_fses, "csvlite",  ",",      NO_FREE);
		lhmss_put(singleton_default_fses, "markdown", "(N/A)",  NO_FREE);
		lhmss_put(singleton_default_fses, "pprint",   " ",      NO_FREE);
		lhmss_put(singleton_default_fses, "xtab",     "\n",     NO_FREE);
	}
	return singleton_default_fses;
}

lhmss_t* get_default_pses() {
	if (singleton_default_pses == NULL) {
		singleton_default_pses = lhmss_alloc();
		lhmss_put(singleton_default_pses, "dkvp",     "=",     NO_FREE);
		lhmss_put(singleton_default_pses, "json",     "(N/A)", NO_FREE);
		lhmss_put(singleton_default_pses, "nidx",     "(N/A)", NO_FREE);
		lhmss_put(singleton_default_pses, "csv",      "(N/A)", NO_FREE);
		lhmss_put(singleton_default_pses, "csvlite",  "(N/A)", NO_FREE);
		lhmss_put(singleton_default_pses, "markdown", "(N/A)", NO_FREE);
		lhmss_put(singleton_default_pses, "pprint",   "(N/A)", NO_FREE);
		lhmss_put(singleton_default_pses, "xtab",     " ",     NO_FREE);
	}
	return singleton_default_pses;
}

lhmsi_t* get_default_repeat_ifses() {
	if (singleton_default_repeat_ifses == NULL) {
		singleton_default_repeat_ifses = lhmsi_alloc();
		lhmsi_put(singleton_default_repeat_ifses, "dkvp",     FALSE, NO_FREE);
		lhmsi_put(singleton_default_repeat_ifses, "json",     FALSE, NO_FREE);
		lhmsi_put(singleton_default_repeat_ifses, "csv",      FALSE, NO_FREE);
		lhmsi_put(singleton_default_repeat_ifses, "csvlite",  FALSE, NO_FREE);
		lhmsi_put(singleton_default_repeat_ifses, "markdown", FALSE, NO_FREE);
		lhmsi_put(singleton_default_repeat_ifses, "nidx",     FALSE, NO_FREE);
		lhmsi_put(singleton_default_repeat_ifses, "xtab",     FALSE, NO_FREE);
		lhmsi_put(singleton_default_repeat_ifses, "pprint",   TRUE,  NO_FREE);
	}
	return singleton_default_repeat_ifses;
}

lhmsi_t* get_default_repeat_ipses() {
	if (singleton_default_repeat_ipses == NULL) {
		singleton_default_repeat_ipses = lhmsi_alloc();
		lhmsi_put(singleton_default_repeat_ipses, "dkvp",     FALSE, NO_FREE);
		lhmsi_put(singleton_default_repeat_ipses, "json",     FALSE, NO_FREE);
		lhmsi_put(singleton_default_repeat_ipses, "csv",      FALSE, NO_FREE);
		lhmsi_put(singleton_default_repeat_ipses, "csvlite",  FALSE, NO_FREE);
		lhmsi_put(singleton_default_repeat_ipses, "markdown", FALSE, NO_FREE);
		lhmsi_put(singleton_default_repeat_ipses, "nidx",     FALSE, NO_FREE);
		lhmsi_put(singleton_default_repeat_ipses, "xtab",     TRUE,  NO_FREE);
		lhmsi_put(singleton_default_repeat_ipses, "pprint",   FALSE, NO_FREE);
	}
	return singleton_default_repeat_ipses;
}

static void free_opt_singletons() {
	lhmss_free(singleton_pdesc_to_chars_map);
	lhmss_free(singleton_default_rses);
	lhmss_free(singleton_default_fses);
	lhmss_free(singleton_default_pses);
	lhmsi_free(singleton_default_repeat_ifses);
	lhmsi_free(singleton_default_repeat_ipses);
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
// The main_usage() function is split out into subroutines in support of the
// manpage autogenerator.

static void main_usage_synopsis(FILE* o, char* argv0) {
	fprintf(o, "Usage: %s [I/O options] {verb} [verb-dependent options ...] {zero or more file names}\n", argv0);
}

static void main_usage_examples(FILE* o, char* argv0, char* leader) {

	fprintf(o, "%s%s --csv cut -f hostname,uptime mydata.csv\n", leader, argv0);
	fprintf(o, "%s%s --tsv --rs lf filter '$status != \"down\" && $upsec >= 10000' *.tsv\n", leader, argv0);
	fprintf(o, "%s%s --nidx put '$sum = $7 < 0.0 ? 3.5 : $7 + 2.1*$8' *.dat\n", leader, argv0);
	fprintf(o, "%sgrep -v '^#' /etc/group | %s --ifs : --nidx --opprint label group,pass,gid,member then sort -f group\n", leader, argv0);
	fprintf(o, "%s%s join -j account_id -f accounts.dat then group-by account_name balances.dat\n", leader, argv0);
	fprintf(o, "%s%s --json put '$attr = sub($attr, \"([0-9]+)_([0-9]+)_.*\", \"\\1:\\2\")' data/*.json\n", leader, argv0);
	fprintf(o, "%s%s stats1 -a min,mean,max,p10,p50,p90 -f flag,u,v data/*\n", leader, argv0);
	fprintf(o, "%s%s stats2 -a linreg-pca -f u,v -g shape data/*\n", leader, argv0);
	fprintf(o, "%s%s put -q '@sum[$a][$b] += $x; end {emit @sum, \"a\", \"b\"}' data/*\n", leader, argv0);
	fprintf(o, "%s%s --from estimates.tbl put '\n", leader, argv0);
	fprintf(o, "  for (k,v in $*) {\n");
	fprintf(o, "    if (isnumeric(v) && k =~ \"^[t-z].*$\") {\n");
	fprintf(o, "      $sum += v; $count += 1\n");
	fprintf(o, "    }\n");
	fprintf(o, "  }\n");
	fprintf(o, "  $mean = $sum / $count # no assignment if count unset'\n");
	fprintf(o, "%s%s --from infile.dat put -f analyze.mlr\n", leader, argv0);
	fprintf(o, "%s%s --from infile.dat put 'tee > \"./taps/data-\".$a.\"-\".$b, $*'\n", leader, argv0);
	fprintf(o, "%s%s --from infile.dat put 'tee | \"gzip > ./taps/data-\".$a.\"-\".$b.\".gz\", $*'\n", leader, argv0);
	fprintf(o, "%s%s --from infile.dat put -q '@v=$*; dump | \"jq .[]\"'\n", leader, argv0);
	fprintf(o, "%s%s --from infile.dat put  '(NR %% 1000 == 0) { print > stderr, \"Checkpoint \".NR}'\n",
		leader, argv0);
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
	rval_evaluator_list_functions(o, leader);
	fprintf(o, "Please use \"%s --help-function {function name}\" for function-specific help.\n", argv0);
	fprintf(o, "Please use \"%s --help-all-functions\" or \"%s -f\" for help on all functions.\n", argv0, argv0);
	fprintf(o, "Please use \"%s --help-all-keywords\" or \"%s -k\" for help on all keywords.\n", argv0, argv0);
}

static void main_usage_data_format_examples(FILE* o, char* argv0) {
	fprintf(o,
		"  DKVP: delimited key-value pairs (Miller default format)\n"
		"  +---------------------+\n"
		"  | apple=1,bat=2,cog=3 | Record 1: \"apple\" => \"1\", \"bat\" => \"2\", \"cog\" => \"3\"\n"
		"  | dish=7,egg=8,flint  | Record 2: \"dish\" => \"7\", \"egg\" => \"8\", \"3\" => \"flint\"\n"
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
		"  Tabular JSON: nested objects are supported, although arrays within them are not:\n"
		"  +---------------------+\n"
		"  | {                   |\n"
		"  |  \"apple\": 1,        | Record 1: \"apple\" => \"1\", \"bat\" => \"2\", \"cog\" => \"3\"\n"
		"  |  \"bat\": 2,          |\n"
		"  |  \"cog\": 3           |\n"
		"  | }                   |\n"
		"  | {                   |\n"
		"  |   \"dish\": {         | Record 2: \"dish:egg\" => \"7\", \"dish:flint\" => \"8\", \"garlic\" => \"\"\n"
		"  |     \"egg\": 7,       |\n"
		"  |     \"flint\": 8      |\n"
		"  |   },                |\n"
		"  |   \"garlic\": \"\"      |\n"
		"  | }                   |\n"
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
		"  +---------------------+\n"
		"\n"
		"  Markdown tabular (supported for output only):\n"
		"  +-----------------------+\n"
		"  | | apple | bat | cog | |\n"
		"  | | ---   | --- | --- | |\n"
		"  | | 1     | 2   | 3   | | Record 1: \"apple => \"1\", \"bat\" => \"2\", \"cog\" => \"3\"\n"
		"  | | 4     | 5   | 6   | | Record 2: \"apple\" => \"4\", \"bat\" => \"5\", \"cog\" => \"6\"\n"
		"  +-----------------------+\n");
}

static void main_usage_data_format_options(FILE* o, char* argv0) {
	fprintf(o, "  --idkvp   --odkvp   --dkvp      Delimited key-value pairs, e.g \"a=1,b=2\"\n");
	fprintf(o, "                                  (this is Miller's default format).\n");
	fprintf(o, "\n");
	fprintf(o, "  --inidx   --onidx   --nidx      Implicitly-integer-indexed fields\n");
	fprintf(o, "                                  (Unix-toolkit style).\n");
	fprintf(o, "\n");
	fprintf(o, "  --icsv    --ocsv    --csv       Comma-separated value (or tab-separated\n");
	fprintf(o, "                                  with --fs tab, etc.)\n");
	fprintf(o, "\n");
	fprintf(o, "  --itsv    --otsv    --tsv       Keystroke-savers for \"--icsv --ifs tab\",\n");
	fprintf(o, "                                  \"--ocsv --ofs tab\", \"--csv --fs tab\".\n");
	fprintf(o, "\n");
	fprintf(o, "  --ipprint --opprint --pprint    Pretty-printed tabular (produces no\n");
	fprintf(o, "                                  output until all input is in).\n");
	fprintf(o, "                      --right     Right-justifies all fields for PPRINT output.\n");
	fprintf(o, "\n");
	fprintf(o, "            --omd                 Markdown-tabular (only available for output).\n");
	fprintf(o, "\n");
	fprintf(o, "  --ixtab   --oxtab   --xtab      Pretty-printed vertical-tabular.\n");
	fprintf(o, "                      --xvright   Right-justifies values for XTAB format.\n");
	fprintf(o, "\n");
	fprintf(o, "  --ijson   --ojson   --json      JSON tabular: sequence or list of one-level\n");
	fprintf(o, "                                  maps: {...}{...} or [{...},{...}].\n");
	fprintf(o, "                      --jvstack   Put one key-value pair per line for JSON\n");
	fprintf(o, "                                  output.\n");
	fprintf(o, "                      --jlistwrap Wrap JSON output in outermost [ ].\n");
	fprintf(o, "                      --jquoteall Quote map keys in JSON output, even if they're\n");
	fprintf(o, "                                  numeric.\n");
	fprintf(o, "              --jflatsep {string} Separator for flattening multi-level JSON keys,\n");
	fprintf(o, "                                  e.g. '{\"a\":{\"b\":3}}' becomes a:b => 3 for\n");
	fprintf(o, "                                  non-JSON formats. Defaults to %s.\n",
		DEFAULT_JSON_FLATTEN_SEPARATOR);
	fprintf(o, "\n");
	fprintf(o, "  -p is a keystroke-saver for --nidx --fs space --repifs\n");
	fprintf(o, "\n");
	fprintf(o, "  Examples: --csv for CSV-formatted input and output; --idkvp --opprint for\n");
	fprintf(o, "  DKVP-formatted input and pretty-printed output.\n");
	fprintf(o, "\n");
	fprintf(o, "  PLEASE USE \"%s --csv --rs lf\" FOR NATIVE UN*X (LINEFEED-TERMINATED) CSV FILES.\n", argv0);
	fprintf(o, "  You can also have MLR_CSV_DEFAULT_RS=lf in your shell environment, e.g.\n");
	fprintf(o, "  \"export MLR_CSV_DEFAULT_RS=lf\" or \"setenv MLR_CSV_DEFAULT_RS lf\" depending on\n");
	fprintf(o, "  which shell you use.\n");
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
	fprintf(o, "    Linux-native CSV files.  You can also have \"MLR_CSV_DEFAULT_RS=lf\" in your\n");
	fprintf(o, "    shell environment, e.g.  \"export MLR_CSV_DEFAULT_RS=lf\" or \"setenv\n");
	fprintf(o, "    MLR_CSV_DEFAULT_RS lf\" depending on which shell you use.\n");
	fprintf(o, "  * TSV is simply CSV using tab as field separator (\"--fs tab\").\n");
	fprintf(o, "  * FS/PS are ignored for markdown format; RS is used.\n");
	fprintf(o, "  * All RS/FS/PS options are ignored for JSON format: JSON doesn't allow\n");
	fprintf(o, "    changing these.\n");
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
	fprintf(o, "  --quote-original   Wrap fields in double quotes if and only if they were\n");
	fprintf(o, "                     quoted on input. This isn't sticky for computed fields:\n");
	fprintf(o, "                     e.g. if fields a and b were quoted on input and you do\n");
	fprintf(o, "                     \"put '$c = $a . $b'\" then field c won't inherit a or b's\n");
	fprintf(o, "                     was-quoted-on-input flag.\n");
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
	fprintf(o, "  --nr-progress-mod {m}, with m a positive integer: print filename and record\n");
	fprintf(o, "                     count to stderr every m input records.\n");
	fprintf(o, "  --from {filename}  Use this to specify an input file before the verb(s),\n");
	fprintf(o, "                     rather than after. May be used more than once. Example:\n");
	fprintf(o, "                     \"%s --from a.dat --from b.dat cat\" is the same as\n", argv0);
	fprintf(o, "                     \"%s cat a.dat b.dat\".\n", argv0);
	fprintf(o, "  -n                 Process no input files, nor standard input either. Useful\n");
	fprintf(o, "                     for %s put with begin/end statements only. (Same as --from\n", argv0);
	fprintf(o, "                     /dev/null.) Also useful in \"%s -n put -v '...'\" for\n", argv0);
	fprintf(o, "                     analyzing abstract syntax trees (if that's your thing).\n");
}

static void main_usage_then_chaining(FILE* o, char* argv0) {
	fprintf(o, "Output of one verb may be chained as input to another using \"then\", e.g.\n");
	fprintf(o, "  %s stats1 -a min,mean,max -f flag,u,v -g color then sort -f color\n", argv0);
}

static void main_usage_see_also(FILE* o, char* argv0) {
	fprintf(o, "For more information please see http://johnkerl.org/miller/doc and/or\n");
	fprintf(o, "http://github.com/johnkerl/miller.");
	fprintf(o, " This is Miller version %s.\n", VERSION_STRING);
}

static void print_type_arithmetic_info(FILE* o, char* argv0) {
	for (int i = -2; i < MT_DIM; i++) {
		mv_t a = (mv_t) {.type = i, .free_flags = NO_FREE, .u.intv = 0};
		if (i == -2)
			printf("%-6s |", "(+)");
		else if (i == -1)
			printf("%-6s +", "------");
		else
			printf("%-6s |", mt_describe_type_simple(a.type));

		for (int j = 0; j < MT_DIM; j++) {
			mv_t b = (mv_t) {.type = j, .free_flags = NO_FREE, .u.intv = 0};
			if (i == -2) {
				printf(" %-6s", mt_describe_type_simple(b.type));
			} else if (i == -1) {
				printf(" %-6s", "------");
			} else {
				mv_t c = x_xx_plus_func(&a, &b);
				printf(" %-6s", mt_describe_type_simple(c.type));
			}
		}

		fprintf(o, "\n");
	}
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
	fprintf(stderr, "Please run \"%s --help\" for usage information.\n", argv0);
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
static void cli_set_reader_defaults(cli_reader_opts_t* preader_opts) {
	preader_opts->ifile_fmt                      = "dkvp";
	preader_opts->irs                            = NULL;
	preader_opts->ifs                            = NULL;
	preader_opts->ips                            = NULL;
	preader_opts->input_json_flatten_separator   = DEFAULT_JSON_FLATTEN_SEPARATOR;

	preader_opts->allow_repeat_ifs               = NEITHER_TRUE_NOR_FALSE;
	preader_opts->allow_repeat_ips               = NEITHER_TRUE_NOR_FALSE;
	preader_opts->use_implicit_csv_header        = FALSE;
	preader_opts->use_mmap_for_read              = TRUE;

	preader_opts->prepipe                        = NULL;
}

static void cli_set_writer_defaults(cli_writer_opts_t* pwriter_opts) {
	pwriter_opts->ofile_fmt                      = "dkvp";
	pwriter_opts->ors                            = NULL;
	pwriter_opts->ofs                            = NULL;
	pwriter_opts->ops                            = NULL;

	pwriter_opts->headerless_csv_output          = FALSE;
	pwriter_opts->right_justify_xtab_value       = FALSE;
	pwriter_opts->left_align_pprint              = TRUE;
	pwriter_opts->stack_json_output_vertically   = FALSE;
	pwriter_opts->wrap_json_output_in_outer_list = FALSE;
	pwriter_opts->quote_json_values_always       = FALSE;
	pwriter_opts->output_json_flatten_separator  = DEFAULT_JSON_FLATTEN_SEPARATOR;
	pwriter_opts->oosvar_flatten_separator       = DEFAULT_OOSVAR_FLATTEN_SEPARATOR;

	pwriter_opts->oquoting                       = DEFAULT_OQUOTING;
}

static void cli_set_defaults(cli_opts_t* popts) {
	memset(popts, 0, sizeof(*popts));

	cli_set_reader_defaults(&popts->reader_opts);
	cli_set_writer_defaults(&popts->writer_opts);

	popts->plrec_reader      = NULL;
	popts->pmapper_list      = sllv_alloc();
	popts->plrec_writer      = NULL;
	popts->filenames         = slls_alloc();

	popts->ofmt              = DEFAULT_OFMT;
	popts->nr_progress_mod   = 0LL;
}

// ----------------------------------------------------------------
static int handle_terminal_usage(char** argv, int argc, int argi)
{
	if (streq(argv[argi], "--version")) {
		printf("Miller %s\n", VERSION_STRING);
		return TRUE;
	} else if (streq(argv[argi], "-h")) {
		main_usage(stdout, argv[0]);
		return TRUE;
	} else if (streq(argv[argi], "--help")) {
		main_usage(stdout, argv[0]);
		return TRUE;
	} else if (streq(argv[argi], "--print-type-arithmetic-info")) {
		print_type_arithmetic_info(stdout, argv[0]);
		return TRUE;

	} else if (streq(argv[argi], "--help-all-verbs")) {
		usage_all_verbs(argv[0]);
	} else if (streq(argv[argi], "--list-all-verbs") || streq(argv[argi], "-l")) {
		list_all_verbs(stdout, "");
		return TRUE;
	} else if (streq(argv[argi], "--list-all-verbs-raw") || streq(argv[argi], "-L")) {
		list_all_verbs_raw(stdout);
		return TRUE;

	} else if (streq(argv[argi], "--list-all-functions-raw")) {
		rval_evaluator_list_all_functions_raw(stdout);
		return TRUE;
	} else if (streq(argv[argi], "--help-all-functions") || streq(argv[argi], "-f")) {
		rval_evaluator_function_usage(stdout, NULL);
		return TRUE;
	} else if (streq(argv[argi], "--help-function") || streq(argv[argi], "--hf")) {
		check_arg_count(argv, argi, argc, 2);
		rval_evaluator_function_usage(stdout, argv[argi+1]);
		return TRUE;

	} else if (streq(argv[argi], "--list-all-keywords-raw")) {
		mlr_dsl_list_all_keywords_raw(stdout);
		return TRUE;
	} else if (streq(argv[argi], "--help-all-keywords") || streq(argv[argi], "-k")) {
		mlr_dsl_keyword_usage(stdout, NULL);
		return TRUE;
	} else if (streq(argv[argi], "--help-keyword") || streq(argv[argi], "--hk")) {
		check_arg_count(argv, argi, argc, 2);
		mlr_dsl_keyword_usage(stdout, argv[argi+1]);
		return TRUE;

	// main-usage subsections, individually accessible for the benefit of
	// the manpage-autogenerator
	} else if (streq(argv[argi], "--usage-synopsis")) {
		main_usage_synopsis(stdout, argv[0]);
		return TRUE;
	} else if (streq(argv[argi], "--usage-examples")) {
		main_usage_examples(stdout, argv[0], "");
		return TRUE;
	} else if (streq(argv[argi], "--usage-list-all-verbs")) {
		list_all_verbs(stdout, "");
		return TRUE;
	} else if (streq(argv[argi], "--usage-help-options")) {
		main_usage_help_options(stdout, argv[0]);
		return TRUE;
	} else if (streq(argv[argi], "--usage-functions")) {
		main_usage_functions(stdout, argv[0], "");
		return TRUE;
	} else if (streq(argv[argi], "--usage-data-format-examples")) {
		main_usage_data_format_examples(stdout, argv[0]);
		return TRUE;
	} else if (streq(argv[argi], "--usage-data-format-options")) {
		main_usage_data_format_options(stdout, argv[0]);
		return TRUE;
	} else if (streq(argv[argi], "--usage-compressed-data-options")) {
		main_usage_compressed_data_options(stdout, argv[0]);
		return TRUE;
	} else if (streq(argv[argi], "--usage-separator-options")) {
		main_usage_separator_options(stdout, argv[0]);
		return TRUE;
	} else if (streq(argv[argi], "--usage-csv-options")) {
		main_usage_csv_options(stdout, argv[0]);
		return TRUE;
	} else if (streq(argv[argi], "--usage-double-quoting")) {
		main_usage_double_quoting(stdout, argv[0]);
		return TRUE;
	} else if (streq(argv[argi], "--usage-numerical-formatting")) {
		main_usage_numerical_formatting(stdout, argv[0]);
		return TRUE;
	} else if (streq(argv[argi], "--usage-other-options")) {
		main_usage_other_options(stdout, argv[0]);
		return TRUE;
	} else if (streq(argv[argi], "--usage-then-chaining")) {
		main_usage_then_chaining(stdout, argv[0]);
		return TRUE;
	} else if (streq(argv[argi], "--usage-see-also")) {
		main_usage_see_also(stdout, argv[0]);
		return TRUE;
	}
	return FALSE;
}

// Returns TRUE if the current flag was handled.
static int handle_reader_options(char** argv, int argc, int *pargi, cli_reader_opts_t* preader_opts)
{
	int argi = *pargi;
	int oargi = argi;

	if (streq(argv[argi], "--irs")) {
		check_arg_count(argv, argi, argc, 2);
		preader_opts->irs = cli_sep_from_arg(argv[argi+1]);
		argi += 2;

	} else if (streq(argv[argi], "--ifs")) {
		check_arg_count(argv, argi, argc, 2);
		preader_opts->ifs = cli_sep_from_arg(argv[argi+1]);
		argi += 2;

	} else if (streq(argv[argi], "--repifs")) {
		preader_opts->allow_repeat_ifs = TRUE;
		argi += 1;

	} else if (streq(argv[argi], "--implicit-csv-header")) {
		preader_opts->use_implicit_csv_header = TRUE;
		argi += 1;

	} else if (streq(argv[argi], "--ips")) {
		check_arg_count(argv, argi, argc, 2);
		preader_opts->ips = cli_sep_from_arg(argv[argi+1]);
		argi += 2;

	} else if (streq(argv[argi], "--icsv")) {
		preader_opts->ifile_fmt = "csv";
		argi += 1;

	} else if (streq(argv[argi], "--icsvlite")) {
		preader_opts->ifile_fmt = "csvlite";
		argi += 1;

	} else if (streq(argv[argi], "--itsv")) {
		preader_opts->ifile_fmt = "csv";
		preader_opts->ifs = "\t";
		argi += 1;

	} else if (streq(argv[argi], "--itsvlite")) {
		preader_opts->ifile_fmt = "csvlite";
		preader_opts->ifs = "\t";
		argi += 1;

	} else if (streq(argv[argi], "--idkvp")) {
		preader_opts->ifile_fmt = "dkvp";
		argi += 1;

	} else if (streq(argv[argi], "--ijson")) {
		preader_opts->ifile_fmt = "json";
		argi += 1;

	} else if (streq(argv[argi], "--inidx")) {
		preader_opts->ifile_fmt = "nidx";
		argi += 1;

	} else if (streq(argv[argi], "--ixtab")) {
		preader_opts->ifile_fmt = "xtab";
		argi += 1;

	} else if (streq(argv[argi], "--ipprint")) {
		preader_opts->ifile_fmt        = "csvlite";
		preader_opts->ifs              = " ";
		preader_opts->allow_repeat_ifs = TRUE;
		argi += 1;

	} else if (streq(argv[argi], "--mmap")) {
		preader_opts->use_mmap_for_read = TRUE;
		argi += 1;

	} else if (streq(argv[argi], "--no-mmap")) {
		preader_opts->use_mmap_for_read = FALSE;
		argi += 1;

	} else if (streq(argv[argi], "--prepipe")) {
		check_arg_count(argv, argi, argc, 2);
		preader_opts->prepipe = argv[argi+1];
		preader_opts->use_mmap_for_read = FALSE;
		argi += 2;

	}
	*pargi = argi;
	return argi != oargi;
}

// Returns TRUE if the current flag was handled.
static int handle_writer_options(char** argv, int argc, int *pargi, cli_writer_opts_t* pwriter_opts)
{
	int argi = *pargi;
	int oargi = argi;

	if (streq(argv[argi], "--ors")) {
		check_arg_count(argv, argi, argc, 2);
		pwriter_opts->ors = cli_sep_from_arg(argv[argi+1]);
		argi += 2;

	} else if (streq(argv[argi], "--ofs")) {
		check_arg_count(argv, argi, argc, 2);
		pwriter_opts->ofs = cli_sep_from_arg(argv[argi+1]);
		argi += 2;

	} else if (streq(argv[argi], "--headerless-csv-output")) {
		pwriter_opts->headerless_csv_output = TRUE;
		argi += 1;

	} else if (streq(argv[argi], "--ops")) {
		check_arg_count(argv, argi, argc, 2);
		pwriter_opts->ops = cli_sep_from_arg(argv[argi+1]);
		argi += 2;

	} else if (streq(argv[argi], "--xvright")) {
		pwriter_opts->right_justify_xtab_value = TRUE;
		argi += 1;

	} else if (streq(argv[argi], "--jvstack")) {
		pwriter_opts->stack_json_output_vertically = TRUE;
		argi += 1;

	} else if (streq(argv[argi], "--jlistwrap")) {
		pwriter_opts->wrap_json_output_in_outer_list = TRUE;
		argi += 1;

	} else if (streq(argv[argi], "--jquoteall")) {
		pwriter_opts->quote_json_values_always = TRUE;
		argi += 1;

	} else if (streq(argv[argi], "--vflatsep")) {
		check_arg_count(argv, argi, argc, 2);
		pwriter_opts->oosvar_flatten_separator = cli_sep_from_arg(argv[argi+1]);
		argi += 2;

	} else if (streq(argv[argi], "--ocsv")) {
		pwriter_opts->ofile_fmt = "csv";
		argi += 1;

	} else if (streq(argv[argi], "--ocsvlite")) {
		pwriter_opts->ofile_fmt = "csvlite";
		argi += 1;

	} else if (streq(argv[argi], "--otsv")) {
		pwriter_opts->ofile_fmt = "csv";
		pwriter_opts->ofs = "\t";
		argi += 1;

	} else if (streq(argv[argi], "--otsvlite")) {
		pwriter_opts->ofile_fmt = "csvlite";
		pwriter_opts->ofs = "\t";
		argi += 1;

	} else if (streq(argv[argi], "--omd")) {
		pwriter_opts->ofile_fmt = "markdown";
		argi += 1;

	} else if (streq(argv[argi], "--odkvp")) {
		pwriter_opts->ofile_fmt = "dkvp";
		argi += 1;

	} else if (streq(argv[argi], "--ojson")) {
		pwriter_opts->ofile_fmt = "json";
		argi += 1;

	} else if (streq(argv[argi], "--onidx")) {
		pwriter_opts->ofile_fmt = "nidx";
		argi += 1;

	} else if (streq(argv[argi], "--oxtab")) {
		pwriter_opts->ofile_fmt = "xtab";
		argi += 1;

	} else if (streq(argv[argi], "--opprint")) {
		pwriter_opts->ofile_fmt = "pprint";
		argi += 1;

	} else if (streq(argv[argi], "--right")) {
		pwriter_opts->left_align_pprint = FALSE;
		argi += 1;

	} else if (streq(argv[argi], "--quote-all")) {
		pwriter_opts->oquoting = QUOTE_ALL;
		argi += 1;

	} else if (streq(argv[argi], "--quote-none")) {
		pwriter_opts->oquoting = QUOTE_NONE;
		argi += 1;

	} else if (streq(argv[argi], "--quote-minimal")) {
		pwriter_opts->oquoting = QUOTE_MINIMAL;
		argi += 1;

	} else if (streq(argv[argi], "--quote-numeric")) {
		pwriter_opts->oquoting = QUOTE_NUMERIC;
		argi += 1;

	} else if (streq(argv[argi], "--quote-original")) {
		pwriter_opts->oquoting = QUOTE_ORIGINAL;
		argi += 1;

	}
	*pargi = argi;
	return argi != oargi;
}

// Returns TRUE if the current flag was handled.
static int handle_reader_writer_options(char** argv, int argc, int *pargi,
	cli_reader_opts_t* preader_opts, cli_writer_opts_t* pwriter_opts)
{
	int argi = *pargi;
	int oargi = argi;

	if (streq(argv[argi], "--rs")) {
		check_arg_count(argv, argi, argc, 2);
		preader_opts->irs = cli_sep_from_arg(argv[argi+1]);
		pwriter_opts->ors = cli_sep_from_arg(argv[argi+1]);
		argi += 2;

	} else if (streq(argv[argi], "--fs")) {
		check_arg_count(argv, argi, argc, 2);
		preader_opts->ifs = cli_sep_from_arg(argv[argi+1]);
		pwriter_opts->ofs = cli_sep_from_arg(argv[argi+1]);
		argi += 2;

	} else if (streq(argv[argi], "-p")) {
		preader_opts->ifile_fmt = "nidx";
		pwriter_opts->ofile_fmt = "nidx";
		preader_opts->ifs = " ";
		pwriter_opts->ofs = " ";
		preader_opts->allow_repeat_ifs = TRUE;
		argi += 1;

	} else if (streq(argv[argi], "--ps")) {
		check_arg_count(argv, argi, argc, 2);
		preader_opts->ips = cli_sep_from_arg(argv[argi+1]);
		pwriter_opts->ops = cli_sep_from_arg(argv[argi+1]);
		argi += 2;

	} else if (streq(argv[argi], "--jflatsep")) {
		check_arg_count(argv, argi, argc, 2);
		preader_opts->input_json_flatten_separator  = cli_sep_from_arg(argv[argi+1]);
		pwriter_opts->output_json_flatten_separator = cli_sep_from_arg(argv[argi+1]);
		argi += 2;

	} else if (streq(argv[argi], "--csv")) {
		preader_opts->ifile_fmt = "csv";
		pwriter_opts->ofile_fmt = "csv";
		argi += 1;

	} else if (streq(argv[argi], "--csvlite")) {
		preader_opts->ifile_fmt = "csvlite";
		pwriter_opts->ofile_fmt = "csvlite";
		argi += 1;

	} else if (streq(argv[argi], "--tsv")) {
		preader_opts->ifile_fmt = pwriter_opts->ofile_fmt = "csv";
		preader_opts->ifs = "\t";
		pwriter_opts->ofs = "\t";
		argi += 1;

	} else if (streq(argv[argi], "--tsvlite")) {
		preader_opts->ifile_fmt = pwriter_opts->ofile_fmt = "csvlite";
		preader_opts->ifs = "\t";
		pwriter_opts->ofs = "\t";
		argi += 1;

	} else if (streq(argv[argi], "--dkvp")) {
		preader_opts->ifile_fmt = "dkvp";
		pwriter_opts->ofile_fmt = "dkvp";
		argi += 1;

	} else if (streq(argv[argi], "--json")) {
		preader_opts->ifile_fmt = "json";
		pwriter_opts->ofile_fmt = "json";
		argi += 1;

	} else if (streq(argv[argi], "--nidx")) {
		preader_opts->ifile_fmt = "nidx";
		pwriter_opts->ofile_fmt = "nidx";
		argi += 1;

	} else if (streq(argv[argi], "--xtab")) {
		preader_opts->ifile_fmt = "xtab";
		pwriter_opts->ofile_fmt = "xtab";
		argi += 1;

	} else if (streq(argv[argi], "--pprint")) {
		preader_opts->ifile_fmt        = "csvlite";
		preader_opts->ifs              = " ";
		preader_opts->allow_repeat_ifs = TRUE;
		pwriter_opts->ofile_fmt        = "pprint";
		argi += 1;

	}
	*pargi = argi;
	return argi != oargi;
}

// ----------------------------------------------------------------
cli_opts_t* parse_command_line(int argc, char** argv)
{
	cli_opts_t* popts = mlr_malloc_or_die(sizeof(cli_opts_t));

	cli_set_defaults(popts);

	int no_input             = FALSE;
	int have_rand_seed       = FALSE;
	unsigned rand_seed       = 0;


	int argi = 1;
	for (; argi < argc; /* variable increment: 1 or 2 depending on flag */) {

		if (argv[argi][0] != '-') {
			break; // No more flag options to process
		} else if (handle_terminal_usage(argv, argc, argi)) {
			exit(0);
		} else if (handle_reader_options(argv, argc, &argi, &popts->reader_opts)) {
			// handled
		} else if (handle_writer_options(argv, argc, &argi, &popts->writer_opts)) {
			// handled
		} else if (handle_reader_writer_options(argv, argc, &argi, &popts->reader_opts, &popts->writer_opts)) {
			// handled

		} else if (streq(argv[argi], "-n")) {
			no_input = TRUE;
			argi += 1;

		} else if (streq(argv[argi], "--from")) {
			check_arg_count(argv, argi, argc, 2);
			slls_append(popts->filenames, argv[argi+1], NO_FREE);
			argi += 2;

		} else if (streq(argv[argi], "--ofmt")) {
			check_arg_count(argv, argi, argc, 2);
			popts->ofmt = argv[argi+1];
			argi += 2;

		} else if (streq(argv[argi], "--nr-progress-mod")) {
			check_arg_count(argv, argi, argc, 2);
			if (sscanf(argv[argi+1], "%lld", &popts->nr_progress_mod) != 1) {
				main_usage(stderr, argv[0]);
				exit(1);
			}
			if (popts->nr_progress_mod <= 0) {
				main_usage(stderr, argv[0]);
				exit(1);
			}
			argi += 2;

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
			argi += 2;

		//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
		} else {
			usage_unrecognized_verb(argv[0], argv[argi]);
		}
	}

	lhmss_t* default_rses = get_default_rses();
	lhmss_t* default_fses = get_default_fses();
	lhmss_t* default_pses = get_default_pses();
	lhmsi_t* default_repeat_ifses = get_default_repeat_ifses();
	lhmsi_t* default_repeat_ipses = get_default_repeat_ipses();

	if (popts->reader_opts.irs == NULL)
		popts->reader_opts.irs = lhmss_get(default_rses, popts->reader_opts.ifile_fmt);
	if (popts->reader_opts.ifs == NULL)
		popts->reader_opts.ifs = lhmss_get(default_fses, popts->reader_opts.ifile_fmt);
	if (popts->reader_opts.ips == NULL)
		popts->reader_opts.ips = lhmss_get(default_pses, popts->reader_opts.ifile_fmt);

	if (popts->reader_opts.allow_repeat_ifs == NEITHER_TRUE_NOR_FALSE)
		popts->reader_opts.allow_repeat_ifs = lhmsi_get(default_repeat_ifses, popts->reader_opts.ifile_fmt);
	if (popts->reader_opts.allow_repeat_ips == NEITHER_TRUE_NOR_FALSE)
		popts->reader_opts.allow_repeat_ips = lhmsi_get(default_repeat_ipses, popts->reader_opts.ifile_fmt);

	if (popts->writer_opts.ors == NULL)
		popts->writer_opts.ors = lhmss_get(default_rses, popts->writer_opts.ofile_fmt);
	if (popts->writer_opts.ofs == NULL)
		popts->writer_opts.ofs = lhmss_get(default_fses, popts->writer_opts.ofile_fmt);
	if (popts->writer_opts.ops == NULL)
		popts->writer_opts.ops = lhmss_get(default_pses, popts->writer_opts.ofile_fmt);

	// xxx fold into get-default-or-die methods
	if (popts->reader_opts.irs == NULL) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n", argv[0], __FILE__, __LINE__);
		exit(1);
	}
	if (popts->reader_opts.ifs == NULL) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n", argv[0], __FILE__, __LINE__);
		exit(1);
	}
	if (popts->reader_opts.ips == NULL) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n", argv[0], __FILE__, __LINE__);
		exit(1);
	}

	if (popts->reader_opts.allow_repeat_ifs == NEITHER_TRUE_NOR_FALSE) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n", argv[0], __FILE__, __LINE__);
		exit(1);
	}
	if (popts->reader_opts.allow_repeat_ips == NEITHER_TRUE_NOR_FALSE) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n", argv[0], __FILE__, __LINE__);
		exit(1);
	}

	if (popts->writer_opts.ors == NULL) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n", argv[0], __FILE__, __LINE__);
		exit(1);
	}
	if (popts->writer_opts.ofs == NULL) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n", argv[0], __FILE__, __LINE__);
		exit(1);
	}
	if (popts->writer_opts.ops == NULL) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n", argv[0], __FILE__, __LINE__);
		exit(1);
	}

	if (streq(popts->writer_opts.ofile_fmt, "pprint") && strlen(popts->writer_opts.ofs) != 1) {
		fprintf(stderr, "%s: OFS for PPRINT format must be single-character; got \"%s\".\n",
			argv[0], popts->writer_opts.ofs);
		return NULL;
	}

	// xxx fold into alloc-or-die methods
	popts->plrec_writer = lrec_writer_alloc(&popts->writer_opts);
	if (popts->plrec_writer == NULL) {
		main_usage(stderr, argv[0]);
		exit(1);
	}

	// xxx make method
	if ((argc - argi) < 1) {
		main_usage(stderr, argv[0]);
		exit(1);
	}
	while (TRUE) {
		check_arg_count(argv, argi, argc, 1);
		char* verb = argv[argi];

		// xxx fold into look-up-or-die methods
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
		sllv_append(popts->pmapper_list, pmapper);

		if (argi >= argc || !streq(argv[argi], "then"))
			break;
		argi++;
	}

	for ( ; argi < argc; argi++) {
		slls_append(popts->filenames, argv[argi], NO_FREE);
	}

	if (no_input) {
		slls_free(popts->filenames);
		popts->filenames = NULL;
	} else if (popts->filenames->length == 0) {
		// No filenames means read from standard input, and standard input cannot be mmapped.
		popts->reader_opts.use_mmap_for_read = FALSE;
	}

	popts->plrec_reader = lrec_reader_alloc(&popts->reader_opts);
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

	popts->plrec_reader->pfree_func(popts->plrec_reader);

	for (sllve_t* pe = popts->pmapper_list->phead; pe != NULL; pe = pe->pnext) {
		mapper_t* pmapper = pe->pvvalue;
		pmapper->pfree_func(pmapper);
	}
	sllv_free(popts->pmapper_list);

	popts->plrec_writer->pfree_func(popts->plrec_writer);

	slls_free(popts->filenames);

	free(popts);

	free_opt_singletons();
}
