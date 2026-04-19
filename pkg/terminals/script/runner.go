package script

import (
	"bufio"
	"fmt"
	"os"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/dsl/cst"
	"github.com/johnkerl/miller/v6/pkg/input"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/output"
	"github.com/johnkerl/miller/v6/pkg/runtime"
	"github.com/johnkerl/miller/v6/pkg/types"
)

func NewScript(
	options *cli.TOptions,
	doWarnings bool,
	strictMode bool,
	dslStrings []string,
) (*Script, error) {
	recordReader, err := input.Create(&options.ReaderOptions, 1)
	if err != nil {
		return nil, err
	}

	recordWriter, err := output.Create(&options.WriterOptions)
	if err != nil {
		return nil, err
	}

	context := types.NewContext()
	runtimeState := runtime.NewEmptyState(options, strictMode)
	runtimeState.NoExitOnFunctionNotFound = true
	runtimeState.Update(nil, context)
	runtimeState.FilterExpression = mlrval.NULL

	cstRootNode := cst.NewEmptyRoot(
		&options.WriterOptions, cst.DSLInstanceTypeScript,
	).WithRedefinableUDFUDS().WithStrictMode(strictMode)

	// Collect all DSL: preloads first, then main script
	allDSLStrings := []string{}
	for _, filename := range options.DSLPreloadFileNames {
		theseStrings, err := lib.LoadStringsFromFileOrDir(filename, ".mlr")
		if err != nil {
			return nil, fmt.Errorf("cannot load from \"%s\": %w", filename, err)
		}
		allDSLStrings = append(allDSLStrings, theseStrings...)
	}
	allDSLStrings = append(allDSLStrings, dslStrings...)

	scr := &Script{
		name:         "script",
		doWarnings:   doWarnings,
		cstRootNode:  cstRootNode,
		options:      options,
		recordReader: recordReader,
		recordWriter: recordWriter,
		runtimeState: runtimeState,
	}

	scr.recordOutputFileName = "(stdout)"
	scr.recordOutputStream = os.Stdout
	scr.bufferedRecordOutputStream = bufio.NewWriter(os.Stdout)

	_, err = cstRootNode.Build(
		allDSLStrings,
		cst.DSLInstanceTypeScript,
		false,
		doWarnings,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return scr, nil
}

func (scr *Script) openFiles(filenames []string) {
	scr.options.FileNames = filenames
	scr.readerChannel = make(chan []*types.RecordAndContext, 2)
	scr.errorChannel = make(chan error, 1)
	scr.downstreamDoneChannel = make(chan bool, 1)

	go scr.recordReader.Read(
		filenames,
		*scr.runtimeState.Context,
		scr.readerChannel,
		scr.errorChannel,
		scr.downstreamDoneChannel,
	)
}

func (scr *Script) run() error {
	// Wire NextRecordFunc before running
	readerChannel := scr.readerChannel
	errorChannel := scr.errorChannel
	bufferedOutput := scr.bufferedRecordOutputStream

	scr.runtimeState.NextRecordFunc = func() (*mlrval.Mlrmap, *types.Context, bool) {
		for {
			var recordsAndContexts []*types.RecordAndContext
			var err error
			select {
			case recordsAndContexts = <-readerChannel:
			case err = <-errorChannel:
			}

			if err != nil {
				return nil, scr.runtimeState.Context, false
			}

			if recordsAndContexts == nil {
				return nil, scr.runtimeState.Context, false
			}

			lib.InternalCodingErrorIf(len(recordsAndContexts) != 1)
			rac := recordsAndContexts[0]

			if rac.EndOfStream {
				return nil, &rac.Context, false
			}

			if rac.Record == nil {
				// Output string from print/dump etc
				bufferedOutput.WriteString(rac.OutputString)
				bufferedOutput.Flush()
				continue
			}

			// Actual record
			return rac.Record, &rac.Context, true
		}
	}

	// Execute begin blocks
	err := scr.cstRootNode.ExecuteBeginBlocks(scr.runtimeState)
	if err != nil {
		return err
	}

	// Execute main block (script drives itself via next())
	_, err = scr.cstRootNode.ExecuteMainBlock(scr.runtimeState)
	if err != nil {
		return err
	}

	// Execute end blocks
	err = scr.cstRootNode.ExecuteEndBlocks(scr.runtimeState)
	if err != nil {
		return err
	}

	return nil
}

func (scr *Script) closeBufferedOutputStream() error {
	if scr.recordOutputStream != os.Stdout {
		err := scr.recordOutputStream.Close()
		if err != nil {
			return fmt.Errorf("error on close: %w", err)
		}
	}
	return nil
}
