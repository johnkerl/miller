#ifndef STREAM_H
#define STREAM_H

#include "cli/mlrcli.h"
#include "lib/context.h"
#include "containers/slls.h"
#include "containers/sllv.h"
#include "input/lrec_readers.h"
#include "mapping/mappers.h"
#include "output/lrec_writers.h"

int do_stream_chained(context_t* pctx, sllv_t* pmapper_list, cli_opts_t* popts);

#endif // STREAM_H
