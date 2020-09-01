// ================================================================
// Miller command-line parsing
// ================================================================

#ifndef MLRCLI_H
#define MLRCLI_H

#include "lib/context.h"
#include "containers/slls.h"
#include "containers/sllv.h"
#include "cli/quoting.h"
#include "cli/comment_handling.h"
#include "cli/json_array_ingest.h"
#include "containers/lhmsll.h"
#include "containers/lhmss.h"
#include <unistd.h>

// ----------------------------------------------------------------
typedef struct _generator_opts_t {
	char* field_name;
	// xxx to do: convert to mv_t
	long long start;
	long long stop;
	long long step;
} generator_opts_t;

typedef struct _cli_reader_opts_t {

	char* ifile_fmt;
	char* irs;
	char* ifs;
	char* ips;
	char* input_json_flatten_separator;
	json_array_ingest_t  json_array_ingest;

	int   allow_repeat_ifs;
	int   allow_repeat_ips;
	int   use_implicit_csv_header;
	int   allow_ragged_csv_input;

	// Command for popen on input, e.g. "zcat -cf <". Can be null in which case
	// files are read directly rather than through a pipe.
	char* prepipe;

	comment_handling_t comment_handling;
	char* comment_string;

	// Fake internal-data-generator 'reader'
	generator_opts_t generator_opts;

} cli_reader_opts_t;

// ----------------------------------------------------------------
typedef struct _cli_writer_opts_t {

	char* ofile_fmt;
	char* ors;
	char* ofs;
	char* ops;

	int   headerless_csv_output;
	int   right_justify_xtab_value;
	int   right_align_pprint;
	int   pprint_barred;
	int   stack_json_output_vertically;
	int   wrap_json_output_in_outer_list;
	int   json_quote_int_keys;
	int   json_quote_non_string_values;
	char* output_json_flatten_separator;
	char* oosvar_flatten_separator;

	quoting_t oquoting;

} cli_writer_opts_t;

// ----------------------------------------------------------------
typedef struct _cli_opts_t {
	cli_reader_opts_t reader_opts;
	cli_writer_opts_t writer_opts;

	// These are used to construct the mapper list. In particular, for in-place mode
	// they're reconstructed for each file.  We make copies since each pass through a
	// CLI-parser operates destructively, principally by running strtok over
	// comma-delimited field-name lists.

	char**  original_argv;
	char**  non_in_place_argv;
	int     argc;
	int     mapper_argb;

	slls_t* filenames;

	char* ofmt;
	long long nr_progress_mod;

	int do_in_place;

	int no_input;
	int have_rand_seed;
	unsigned rand_seed;

} cli_opts_t;

// ----------------------------------------------------------------
cli_opts_t* parse_command_line(int argc, char** argv, sllv_t** ppmapper_list);

// See stream.c. The idea is that the mapper-chain is constructed once for normal stream-over-all-files
// mode, but per-file for in-place mode.
sllv_t* cli_parse_mappers(char** argv, int* pargi, int argc, cli_opts_t* popts);

int cli_handle_reader_options(char** argv, int argc, int *pargi, cli_reader_opts_t* preader_opts);
int cli_handle_writer_options(char** argv, int argc, int *pargi, cli_writer_opts_t* pwriter_opts);
int cli_handle_reader_writer_options(char** argv, int argc, int *pargi,
	cli_reader_opts_t* preader_opts, cli_writer_opts_t* pwriter_opts);

void cli_opts_init(cli_opts_t* popts);
void cli_reader_opts_init(cli_reader_opts_t* preader_opts);
void cli_writer_opts_init(cli_writer_opts_t* pwriter_opts);

void cli_apply_defaults(cli_opts_t* popts);
void cli_apply_reader_defaults(cli_reader_opts_t* preader_opts);
void cli_apply_writer_defaults(cli_writer_opts_t* pwriter_opts);

// For mapper join which has its separate input-format overrides:
void cli_merge_reader_opts(cli_reader_opts_t* pfunc_opts, cli_reader_opts_t* pmain_opts);

// For mapper tee & mapper put which have their separate output-format overrides:
void cli_merge_writer_opts(cli_writer_opts_t* pfunc_opts, cli_writer_opts_t* pmain_opts);

// Stream context is for lrec-writer drain on tee et al. when using aggregated
// output.  E.g. pretty-print output has column widths which are only computable
// after all output records have been retained. The free methods are used as
// drain triggers.
void cli_opts_free(cli_opts_t* popts);

// The caller can unconditionally free the return value
char* cli_sep_from_arg(char* arg);

#endif // MLRCLI_H
