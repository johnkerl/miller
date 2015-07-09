#ifndef JOIN_BUCKET_KEEPER_H
#define JOIN_BUCKET_KEEPER_H

#include "containers/lrec.h"
#include "containers/slls.h"
#include "containers/sllv.h"
#include "input/lrec_reader.h"
#include "mapping/context.h"

typedef struct _join_bucket_keeper_t {
	lrec_reader_t* plrec_reader;
	void*          pvhandle;
	context_t*     pctx;

	int            state;
	slls_t*        pleft_field_values;
	sllv_t*        precords;
	lrec_t*        prec_peek;
	int            leof;

} join_bucket_keeper_t;

join_bucket_keeper_t* join_bucket_keeper_alloc(
	char* left_file_name,
	char* input_file_format,
	int   use_mmap_for_read,
	char  irs,
	char  ifs,
	int   allow_repeat_ifs,
	char  ips,
	int   allow_repeat_ips
);

void join_bucket_keeper_free(join_bucket_keeper_t* pkeeper);

void join_bucket_keeper_emit(join_bucket_keeper_t* pkeeper,
	slls_t* pleft_field_names, slls_t* pright_field_values,
	sllv_t** ppbucket_paired, sllv_t** ppbucket_left_unpaired);


#endif // JOIN_BUCKET_KEEPER_H
