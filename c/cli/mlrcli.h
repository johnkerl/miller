// ================================================================
// Miller command-line parsing
// ================================================================

#ifndef MLRCLI_H
#define MLRCLI_H

#include "containers/slls.h"
#include "containers/sllv.h"
#include "cli/quoting.h"
#include "containers/lhmsll.h"
#include "containers/lhmss.h"
#include "input/lrec_reader.h"
#include "output/lrec_writer.h"

// ----------------------------------------------------------------
typedef struct _cli_reader_opts_t {

	char* ifile_fmt;
	char* irs;
	char* ifs;
	char* ips;
	char* input_json_flatten_separator;

	int   allow_repeat_ifs;
	int   allow_repeat_ips;
	int   use_implicit_csv_header;
	int   use_mmap_for_read;

	// Command for popen on input, e.g. "zcat -cf <". Can be null in which case
	// files are read directly rather than through a pipe.
	char*  prepipe;

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

	lrec_reader_t* plrec_reader;
	sllv_t*        pmapper_list;
	lrec_writer_t* plrec_writer;
	slls_t* filenames;

	char* ofmt;
	long long nr_progress_mod;

} cli_opts_t;

// ----------------------------------------------------------------
cli_opts_t* parse_command_line(int argc, char** argv);

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
void cli_opts_free(cli_opts_t* popts, context_t* pctx);

// The caller can unconditionally free the return value
char* cli_sep_from_arg(char* arg);

#endif // MLRCLI_H
