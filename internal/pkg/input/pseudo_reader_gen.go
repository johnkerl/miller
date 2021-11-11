package input

import (
	"errors"
	"fmt"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/types"
)

type PseudoReaderGen struct {
	readerOptions *cli.TReaderOptions
}

func NewPseudoReaderGen(readerOptions *cli.TReaderOptions) (*PseudoReaderGen, error) {
	return &PseudoReaderGen{
		readerOptions: readerOptions,
	}, nil
}

func (reader *PseudoReaderGen) Read(
	filenames []string, // ignored
	context types.Context,
	inputChannel chan<- *types.RecordAndContext,
	errorChannel chan error,
	downstreamDoneChannel <-chan bool, // for mlr head
) {
	reader.process(&context, inputChannel, errorChannel, downstreamDoneChannel)
	inputChannel <- types.NewEndOfStreamMarker(&context)
}

func (reader *PseudoReaderGen) process(
	context *types.Context,
	inputChannel chan<- *types.RecordAndContext,
	errorChannel chan error,
	downstreamDoneChannel <-chan bool, // for mlr head
) {
	context.UpdateForStartOfFile("(gen-pseudo-reader)")

	start, err := reader.tryParse("start", reader.readerOptions.GeneratorOptions.StartAsString)
	if err != nil {
		errorChannel <- err
		return
	}

	step, err := reader.tryParse("step", reader.readerOptions.GeneratorOptions.StepAsString)
	if err != nil {
		errorChannel <- err
		return
	}

	stop, err := reader.tryParse("stop", reader.readerOptions.GeneratorOptions.StopAsString)
	if err != nil {
		errorChannel <- err
		return
	}

	var doneComparator types.BinaryFunc = types.BIF_greater_than
	if step.GetNumericNegativeorDie() {
		doneComparator = types.BIF_less_than
	}

	key := reader.readerOptions.GeneratorOptions.FieldName
	value := start.Copy()

	eof := false
	for !eof {

		// See if downstream processors will be ignoring further data (e.g. mlr
		// head).  If so, stop reading. This makes 'mlr head hugefile' exit
		// quickly, as it should.
		eof := false
		select {
		case _ = <-downstreamDoneChannel:
			eof = true
			break
		default:
			break
		}
		if eof {
			break
		}

		mdone := doneComparator(value, stop)
		done, _ := mdone.GetBoolValue()
		if done {
			break
		}

		record := types.NewMlrmap()
		record.PutCopy(key, value)

		context.UpdateForInputRecord()
		inputChannel <- types.NewRecordAndContext(
			record,
			context,
		)

		value = types.BIF_plus_binary(value, step)
	}
}

func (reader *PseudoReaderGen) tryParse(
	name string,
	svalue string,
) (*types.Mlrval, error) {
	mvalue := types.MlrvalFromInferredType(svalue)
	if mvalue == nil || !mvalue.IsNumeric() {
		return nil, errors.New(
			fmt.Sprintf("mlr: gen: %s \"%s\" is not parseable as number", name, svalue),
		)
	}
	return mvalue, nil
}

//#include <stdio.h>
//#include <stdlib.h>
//#include "lib/mlr_globals.h"
//#include "lib/mlrutil.h"
//#include "input/lrec_readers.h"
//
//typedef struct _lrec_reader_gen_state_t {
//	char* field_name;
//	unsigned long long start;
//	unsigned long long stop;
//	unsigned long long step;
//	unsigned long long current_value;
//} lrec_reader_gen_state_t;
//
//static void    lrec_reader_gen_free(lrec_reader_t* preader);
//static void*   lrec_reader_gen_open(void* pvstate, char* prepipe, char* filename);
//static void    lrec_reader_gen_close(void* pvstate, void* pvhandle, char* prepipe);
//static void    lrec_reader_gen_sof(void* pvstate, void* pvhandle);
//static lrec_t* lrec_reader_gen_process(void* pvstate, void* pvhandle, context_t* pctx);
//
//// ----------------------------------------------------------------
//lrec_reader_t* lrec_reader_gen_alloc(char* field_name,
//	unsigned long long start, unsigned long long stop, unsigned long long step)
//{
//	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));
//
//	lrec_reader_gen_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_gen_state_t));
//	pstate->field_name    = field_name;
//	pstate->start         = start;
//	pstate->stop          = stop;
//	pstate->step          = step;
//	pstate->current_value = start;
//
//	plrec_reader->pvstate       = (void*)pstate;
//	plrec_reader->popen_func    = lrec_reader_gen_open;
//	plrec_reader->pclose_func   = lrec_reader_gen_close;
//	plrec_reader->pprocess_func = lrec_reader_gen_process;
//	plrec_reader->psof_func     = lrec_reader_gen_sof;
//	plrec_reader->pfree_func    = lrec_reader_gen_free;
//
//	return plrec_reader;
//}
//
//static void* lrec_reader_gen_open(void* pvstate, char* prepipe, char* filename) {
//	return NULL;
//}
//
//static void lrec_reader_gen_close(void* pvstate, void* pvhandle, char* prepipe) {
//}
//
//static void lrec_reader_gen_free(lrec_reader_t* preader) {
//	free(preader->pvstate);
//	free(preader);
//}
//
//static void lrec_reader_gen_sof(void* pvstate, void* pvhandle) {
//	lrec_reader_gen_state_t* pstate = pvstate;
//	pstate->current_value = pstate->start;
//}
//
//// ----------------------------------------------------------------
//static lrec_t* lrec_reader_gen_process(void* pvstate, void* pvhandle, context_t* pctx) {
//	lrec_reader_gen_state_t* pstate = pvstate;
//	if (pstate->current_value > pstate->stop) {
//		return NULL;
//	}
//
//	lrec_t* prec = lrec_unbacked_alloc();
//	char* key = pstate->field_name;
//	char* value = mlr_alloc_string_from_ll(pstate->current_value);
//	pstate->current_value += pstate->step;
//
//	lrec_put(prec, key, value, FREE_ENTRY_VALUE);
//
//	return prec;
//}
