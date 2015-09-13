// ================================================================
// Data structures for mlr join
// ================================================================

#ifndef JOIN_BUCKET_KEEPER_H
#define JOIN_BUCKET_KEEPER_H

#include "containers/lrec.h"
#include "containers/slls.h"
#include "containers/sllv.h"
#include "input/lrec_reader.h"
#include "mapping/context.h"

typedef struct _join_bucket_t {
	slls_t* pleft_field_values;
	sllv_t* precords;
	int     was_paired;
} join_bucket_t;

typedef struct _join_bucket_keeper_t {
	lrec_reader_t* plrec_reader;
	void*          pvhandle;
	context_t*     pctx;

	slls_t*        pleft_field_names;

	join_bucket_t* pbucket;

	lrec_t*        prec_peek;
	int            leof;
	int            state;
} join_bucket_keeper_t;

join_bucket_keeper_t* join_bucket_keeper_alloc(
	char* left_file_name,
	char* input_file_format,
	int   use_mmap_for_read,
	char* irs,
	char* ifs,
	int   allow_repeat_ifs,
	char* ips,
	int   allow_repeat_ips,
	slls_t* pleft_field_names);

join_bucket_keeper_t* join_bucket_keeper_alloc_from_reader(
	lrec_reader_t* plrec_reader,
	char*          left_file_name,
	slls_t*        pleft_field_names);

void join_bucket_keeper_free(join_bucket_keeper_t* pkeeper);

// *pprecords_paired should not be freed by the caller.
// *pprecords_left_unpaired should be freed by the caller.
void join_bucket_keeper_emit(
	join_bucket_keeper_t* pkeeper,
	slls_t*               pright_field_values,
	sllv_t**              pprecords_paired,
	sllv_t**              pprecords_left_unpaired);

void join_bucket_print(join_bucket_t* pbucket, char* indent);
void join_bucket_keeper_print(join_bucket_keeper_t* pkeeper);

#endif // JOIN_BUCKET_KEEPER_H
