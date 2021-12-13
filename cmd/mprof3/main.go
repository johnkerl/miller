// Experiments in performance/profiling.
package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strconv"
	"strings"

	"github.com/pkg/profile" // for trace.out

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/input"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/output"
	"github.com/johnkerl/miller/internal/pkg/types"
)

func main() {

	// Respect env $GOMAXPROCS, if provided, else set default.
	haveSetGoMaxProcs := false
	goMaxProcsString := os.Getenv("GOMAXPROCS")
	if goMaxProcsString != "" {
		goMaxProcs, err := strconv.Atoi(goMaxProcsString)
		if err != nil {
			runtime.GOMAXPROCS(goMaxProcs)
			haveSetGoMaxProcs = true
		}
	}
	if !haveSetGoMaxProcs {
		// As of Go 1.16 this is the default anyway. For 1.15 and below we need
		// to explicitly set this.
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	debug.SetGCPercent(500) // Empirical: See README-profiling.md

	if os.Getenv("MPROF_PPROF") != "" {
		// profiling with cpu.pprof and go tool pprof -http=:8080 cpu.pprof
		profFilename := "cpu.pprof"
		handle, err := os.Create(profFilename)
		if err != nil {
			fmt.Fprintln(os.Stderr, os.Args[0], ": ", "Could not start CPU profile: ", err)
			return
		}
		defer handle.Close()

		if err := pprof.StartCPUProfile(handle); err != nil {
			fmt.Fprintln(os.Stderr, os.Args[0], ": ", "Could not start CPU profile: ", err)
			return
		}
		defer pprof.StopCPUProfile()

		fmt.Fprintf(os.Stderr, "CPU profile started.\n")
		fmt.Fprintf(os.Stderr, "go tool pprof -http=:8080 cpu.pprof\n")
		defer fmt.Fprintf(os.Stderr, "CPU profile finished.\n")
	}

	if os.Getenv("MPROF_TRACE") != "" {
		// tracing with trace.out and go tool trace trace.out
		fmt.Fprintf(os.Stderr, "go tool trace trace.out\n")
		defer profile.Start(profile.TraceProfile, profile.ProfilePath(".")).Stop()
	}

	options := cli.DefaultOptions()
	types.SetInferrerStringOnly()

	filenames := os.Args[1:]
	lib.InternalCodingErrorIf(len(filenames) != 1)
	filename := filenames[0]

	err := Stream(filename, options, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mlr: %v.\n", err)
		os.Exit(1)
	}
}

func getBatchSize() int {
	return 1000
}

// ================================================================
type IRecordReader interface {
	Read(ioChannel chan<- *list.List) error
}

func Stream(
	filename string,
	options *cli.TOptions,
	outputStream io.WriteCloser,
) error {
	initialContext := types.NewContext()

	// Instantiate the record-reader
	recordReader, err := NewRecordReaderDKVPChanPipelined(&options.ReaderOptions, filename, initialContext)
	if err != nil {
		return err
	}

	// Instantiate the record-writer
	recordWriter, err := output.NewRecordWriterDKVP(&options.WriterOptions)
	if err != nil {
		return err
	}

	bufferedOutputStream := bufio.NewWriter(os.Stdout)
	defer bufferedOutputStream.Flush()

	ioChannel := make(chan *list.List, 1)
	errorChannel := make(chan error, 1)
	doneWritingChannel := make(chan bool, 1)

	go recordReader.Read(ioChannel)
	go output.ChannelWriter(ioChannel, recordWriter, &options.WriterOptions, doneWritingChannel,
		bufferedOutputStream, true)

	done := false
	for !done {
		select {
		case err := <-errorChannel:
			////fmt.Fprintf(os.Stderr, "ECHAN READ\n")
			fmt.Fprintln(os.Stderr, "mlr", ": ", err)
			os.Exit(1)
		case _ = <-doneWritingChannel:
			////fmt.Fprintf(os.Stderr, "ZCHAN READ\n")
			done = true
			break
		}
	}

	return nil
}

// ================================================================

type RecordReaderDKVPChanPipelined struct {
	readerOptions  *cli.TReaderOptions
	filename       string
	initialContext *types.Context
}

func NewRecordReaderDKVPChanPipelined(
	readerOptions *cli.TReaderOptions,
	filename string,
	initialContext *types.Context,
) (*RecordReaderDKVPChanPipelined, error) {
	return &RecordReaderDKVPChanPipelined{
		readerOptions:  readerOptions,
		filename:       filename,
		initialContext: initialContext,
	}, nil
}

func (reader *RecordReaderDKVPChanPipelined) Read(
	readerChannel chan<- *list.List,
) error {
	handle, err := lib.OpenFileForRead(
		reader.filename,
		reader.readerOptions.Prepipe,
		reader.readerOptions.PrepipeIsRaw,
		reader.readerOptions.FileInputEncoding,
	)
	if err != nil {
		return err
	} else {
		reader.processHandle(handle, reader.filename, reader.initialContext, readerChannel)
		handle.Close()
	}

	eom := types.NewEndOfStreamMarker(reader.initialContext)
	leom := list.New()
	leom.PushBack(eom)
	readerChannel <- leom
	////fmt.Fprintf(os.Stderr, "IOCHAN WRITE EOM\n")
	return nil
}

func provideChannelizedLines(
	lineScanner *bufio.Scanner,
	linesChannel chan<- string,
) {
	for lineScanner.Scan() {
		linesChannel <- lineScanner.Text()
	}
	close(linesChannel) // end-of-stream marker
}

func (reader *RecordReaderDKVPChanPipelined) processHandle(
	handle io.Reader,
	filename string,
	context *types.Context,
	readerChannel chan<- *list.List,
) {
	context.UpdateForStartOfFile(filename)
	m := getBatchSize()

	lineScanner := input.NewLineScanner(handle, reader.readerOptions.IRS)
	linesChannel := make(chan string, m)
	go provideChannelizedLines(lineScanner, linesChannel)

	eof := false
	for !eof {
		var recordsAndContexts *list.List
		recordsAndContexts, eof = reader.getRecordBatch(linesChannel, m, context)
		//fmt.Fprintf(os.Stderr, "GOT RECORD BATCH OF LENGTH %d\n", recordsAndContexts.Len())
		readerChannel <- recordsAndContexts
	}
}

// TODO: comment copiously we're trying to handle slow/fast/short/long
// reads: tail -f, smallfile, bigfile.
func (reader *RecordReaderDKVPChanPipelined) getRecordBatch(
	linesChannel <-chan string,
	maxBatchSize int,
	context *types.Context,
) (
	recordsAndContexts *list.List,
	eof bool,
) {
	//fmt.Printf("GRB ENTER\n")
	recordsAndContexts = list.New()
	eof = false

	for i := 0; i < maxBatchSize; i++ {
		//fmt.Fprintf(os.Stderr, "-- %d/%d %d/%d\n", i, maxBatchSize, len(linesChannel), cap(linesChannel))
		if len(linesChannel) == 0 && i > 0 {
			//fmt.Println(" .. BREAK")
			break
		}
		//fmt.Println(" .. B:BLOCK")
		line, more := <-linesChannel
		//fmt.Printf(" .. E:BLOCK <<%s>> %v\n", line, more)
		if !more {
			eof = true
			break
		}

		// Check for comments-in-data feature
		// TODO: funcptr this away
		if reader.readerOptions.CommentHandling != cli.CommentsAreData {
			if strings.HasPrefix(line, reader.readerOptions.CommentString) {
				if reader.readerOptions.CommentHandling == cli.PassComments {
					recordsAndContexts.PushBack(types.NewOutputStringList(line+"\n", context))
					continue
				} else if reader.readerOptions.CommentHandling == cli.SkipComments {
					continue
				}
				// else comments are data
			}
		}

		record := reader.recordFromDKVPLine(line)
		context.UpdateForInputRecord()
		recordAndContext := types.NewRecordAndContext(record, context)
		recordsAndContexts.PushBack(recordAndContext)
	}

	//fmt.Printf("GRB EXIT\n")
	return recordsAndContexts, eof
}

func (reader *RecordReaderDKVPChanPipelined) recordFromDKVPLine(
	line string,
) *types.Mlrmap {
	record := types.NewMlrmapAsRecord()

	var pairs []string
	if reader.readerOptions.IFSRegex == nil { // e.g. --no-ifs-regex
		pairs = lib.SplitString(line, reader.readerOptions.IFS)
	} else {
		pairs = lib.RegexSplitString(reader.readerOptions.IFSRegex, line, -1)
	}

	for i, pair := range pairs {
		var kv []string
		if reader.readerOptions.IPSRegex == nil { // e.g. --no-ips-regex
			kv = strings.SplitN(line, reader.readerOptions.IPS, 2)
		} else {
			kv = lib.RegexSplitString(reader.readerOptions.IPSRegex, pair, 2)
		}

		if len(kv) == 0 {
			// Ignore. This is expected when splitting with repeated IFS.
		} else if len(kv) == 1 {
			// E.g the pair has no equals sign: "a" rather than "a=1" or
			// "a=".  Here we use the positional index as the key. This way
			// DKVP is a generalization of NIDX.
			key := strconv.Itoa(i + 1) // Miller userspace indices are 1-up
			value := types.MlrvalFromInferredTypeForDataFiles(kv[0])
			record.PutReference(key, value)
		} else {
			key := kv[0]
			value := types.MlrvalFromInferredTypeForDataFiles(kv[1])
			record.PutReference(key, value)
		}
	}
	return record
}
