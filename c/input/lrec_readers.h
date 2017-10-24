#ifndef LREC_READERS_H
#define LREC_READERS_H
#include "cli/mlrcli.h"
#include "input/lrec_reader.h"

// ----------------------------------------------------------------
// Primary entry points

lrec_reader_t*  lrec_reader_alloc(cli_reader_opts_t* popts);
lrec_reader_t*  lrec_reader_alloc_or_die(cli_reader_opts_t* popts);

lrec_reader_t* lrec_reader_gen_alloc(char* field_name, unsigned long long start, unsigned long long stop, unsigned long long step);
lrec_reader_t* lrec_reader_stdio_csvlite_alloc(char* irs, char* ifs, int allow_repeat_ifs, int use_implicit_header,
	char* comment_string);
lrec_reader_t* lrec_reader_stdio_csv_alloc(char* irs, char* ifs, int use_implicit_header, char* comment_string);
lrec_reader_t* lrec_reader_stdio_dkvp_alloc(char* irs, char* ifs, char* ips, int allow_repeat_ifs,
	char* comment_string);
lrec_reader_t* lrec_reader_stdio_nidx_alloc(char* irs, char* ifs, int allow_repeat_ifs, char* comment_string);
lrec_reader_t* lrec_reader_stdio_xtab_alloc(char* ifs, char* ips, int allow_repeat_ips, char* comment_string);
lrec_reader_t* lrec_reader_stdio_json_alloc(char* input_json_flatten_separator, json_array_ingest_t json_array_ingest,
	char* line_term);

lrec_reader_t* lrec_reader_mmap_csv_alloc(char* irs, char* ifs, int use_implicit_header, char* comment_string);
lrec_reader_t* lrec_reader_mmap_csvlite_alloc(char* irs, char* ifs, int allow_repeat_ifs, int use_implicit_header,
	char* comment_string);
lrec_reader_t* lrec_reader_mmap_dkvp_alloc(char* irs, char* ifs, char* ips, int allow_repeat_ifs, char* comment_string);
lrec_reader_t* lrec_reader_mmap_nidx_alloc(char* irs, char* ifs, int allow_repeat_ifs, char* comment_string);
lrec_reader_t* lrec_reader_mmap_xtab_alloc(char* ifs, char* ips, int allow_repeat_ips, char* comment_string);
lrec_reader_t* lrec_reader_mmap_json_alloc(char* input_json_flatten_separator, json_array_ingest_t json_array_ingest,
	char* line_term);

lrec_reader_t* lrec_reader_in_memory_alloc(sllv_t* precords);

// ----------------------------------------------------------------
// These entry points are made public for unit test

lrec_t* lrec_parse_stdio_nidx_single_sep(char* line, char ifs, int allow_repeat_ifs);
lrec_t* lrec_parse_stdio_nidx_multi_sep(char* line, char* ifs, int ifslen, int allow_repeat_ifs);

lrec_t* lrec_parse_stdio_dkvp_single_sep(char* line, char ifs, char ips, int allow_repeat_ifs);
lrec_t* lrec_parse_stdio_dkvp_multi_sep(char* line, char* ifs, char* ips, int ifslen, int ipslen, int allow_repeat_ifs);

slls_t* split_csv_header_line(char* line, char ifs, int allow_repeat_ifs);

slls_t* split_csvlite_header_line_single_ifs(char* line, char ifs, int allow_repeat_ifs);
slls_t* split_csvlite_header_line_multi_ifs(char* line, char* ifs, int ifslen, int allow_repeat_ifs);

lrec_t* lrec_parse_stdio_csvlite_data_line_single_ifs(header_keeper_t* pheader_keeper, char* filename, long long ilno,
	char* data_line, char ifs, int allow_repeat_ifs);
lrec_t* lrec_parse_stdio_csvlite_data_line_multi_ifs(header_keeper_t* pheader_keeper, char* filename, long long ilno,
	char* data_line, char* ifs, int ifslen, int allow_repeat_ifs);
lrec_t* lrec_parse_stdio_csvlite_data_line_single_ifs_implicit_header(header_keeper_t* pheader_keeper, char* filename, long long ilno,
	char* data_line, char ifs, int allow_repeat_ifs);
lrec_t* lrec_parse_stdio_csvlite_data_line_multi_ifs_implicit_header(header_keeper_t* pheader_keeper, char* filename, long long ilno,
	char* data_line, char* ifs, int ifslen, int allow_repeat_ifs);

lrec_t* lrec_parse_stdio_xtab_single_ips(slls_t* pxtab_lines, char ips, int allow_repeat_ips);
lrec_t* lrec_parse_stdio_xtab_multi_ips(slls_t* pxtab_lines, char* ips, int ipslen, int allow_repeat_ips);

lrec_t* lrec_parse_mmap_nidx_single_irs_single_ifs(file_reader_mmap_state_t *phandle,
	char irs, char ifs, int allow_repeat_ifs, int do_auto_line_term, context_t* pctx);
lrec_t* lrec_parse_mmap_nidx_single_irs_multi_ifs(file_reader_mmap_state_t *phandle,
	char irs, char* ifs, int ifslen, int allow_repeat_ifs, int do_auto_line_term, context_t* pctx);
lrec_t* lrec_parse_mmap_nidx_multi_irs_single_ifs(file_reader_mmap_state_t *phandle,
	char* irs, char ifs, int irslen, int allow_repeat_ifs);
lrec_t* lrec_parse_mmap_nidx_multi_irs_multi_ifs(file_reader_mmap_state_t *phandle,
	char* irs, char* ifs, int irslen, int ifslen, int allow_repeat_ifs);

lrec_t* lrec_parse_mmap_xtab_single_ifs_single_ips(file_reader_mmap_state_t* phandle, char ifs, char ips, int allow_repeat_ips,
	int do_auto_line_term, context_t* pctx);
lrec_t* lrec_parse_mmap_xtab_single_ifs_multi_ips(file_reader_mmap_state_t* phandle, char ifs, char* ips, int ipslen, int allow_repeat_ips,
	int do_auto_line_term, context_t* pctx);
lrec_t* lrec_parse_mmap_xtab_multi_ifs_single_ips(file_reader_mmap_state_t* phandle, char* ifs, char ips, int ifslen, int allow_repeat_ips);
lrec_t* lrec_parse_mmap_xtab_multi_ifs_multi_ips(file_reader_mmap_state_t* phandle, char* ifs, char* ips, int ipslen, int ifslen, int allow_repeat_ips);

#endif // LREC_READERS_H
