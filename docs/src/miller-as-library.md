<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-main-flag-list/index.html">Flags</a>
&nbsp;
<a class="quicklink" href="../reference-verbs/index.html">Verbs</a>
&nbsp;
<a class="quicklink" href="../reference-dsl-builtin-functions/index.html">Functions</a>
&nbsp;
<a class="quicklink" href="../glossary/index.html">Glossary</a>
&nbsp;
<a class="quicklink" href="../release-docs/index.html">Release docs</a>
</span>
</div>
# Miller as a library

Very initially and experimentally, as of Miller 6.9.1, Go developers will be able to access Miller source
code --- moved from `internal/pkg/` to `pkg/` --- within their own Go projects.

Caveat emptor: Miller's backward-compatibility guarantees are at the CLI level; API is not guaranteed stable.
For this reason, please be careful with your version pins.

I'm happy to discuss this new area further at the [discussions page](https://github.com/johnkerl/miller/discussions).

## Setup

```
$ mkdir use-mlr

$ cd cd use-mlr

$ go mod init github.com/johnkerl/miller-library-example
go: creating new go.mod: module github.com/johnkerl/miller-library-example

# One of:
$ go get github.com/johnkerl/miller
$ go get github.com/johnkerl/miller@0f27a39a9f92d4c633dd29d99ad203e95a484dd3
# Etc.

$ go mod tidy
```

## One example use

<pre class="pre-non-highlight-non-pair">
package main

import (
	"fmt"

	"github.com/johnkerl/miller/pkg/bifs"
	"github.com/johnkerl/miller/pkg/mlrval"
)

func main() {
	a := mlrval.FromInt(2)
	b := mlrval.FromInt(60)
	c := bifs.BIF_pow(a, b)
	fmt.Println(c.String())
}
</pre>

```
$ go build main1.go
$ ./main1
1152921504606846976
```

Or simply:
```
$ go run main1.go
1152921504606846976
```

## Another example use

<pre class="pre-non-highlight-non-pair">
package main

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"os"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/input"
	"github.com/johnkerl/miller/pkg/output"
	"github.com/johnkerl/miller/pkg/transformers"
	"github.com/johnkerl/miller/pkg/types"
)

func convert_csv_to_json(fileNames []string) error {
	options := &cli.TOptions{
		ReaderOptions: cli.TReaderOptions{
			InputFileFormat: "csv",
			IFS:             ",",
			IRS:             "\n",
			RecordsPerBatch: 1,
		},
		WriterOptions: cli.TWriterOptions{
			OutputFileFormat: "json",
		},
	}
	outputStream := os.Stdout
	outputIsStdout := true

	// Since Go is concurrent, the context struct needs to be duplicated and
	// passed through the channels along with each record.
	initialContext := types.NewContext()

	// Instantiate the record-reader.
	// RecordsPerBatch is tracked separately from ReaderOptions since join/repl
	// may use batch size of 1.
	recordReader, err := input.Create(&options.ReaderOptions, options.ReaderOptions.RecordsPerBatch)
	if err != nil {
		return err
	}

	// Instantiate the record-writer
	recordWriter, err := output.Create(&options.WriterOptions)
	if err != nil {
		return err
	}

	cat, err := transformers.NewTransformerCat(
		false, // doCounters bool,
		"",    // counterFieldName string,
		nil,   // groupByFieldNames []string,
		false, // doFileName bool,
		false, // doFileNum bool,
	)
	if err != nil {
		return err
	}
	recordTransformers := []transformers.IRecordTransformer{cat}

	// Set up the reader-to-transformer and transformer-to-writer channels.
	readerChannel := make(chan *list.List, 2) // list of *types.RecordAndContext
	writerChannel := make(chan *list.List, 1) // list of *types.RecordAndContext

	// We're done when a fatal error is registered on input (file not found,
	// etc) or when the record-writer has written all its output. We use
	// channels to communicate both of these conditions.
	inputErrorChannel := make(chan error, 1)
	doneWritingChannel := make(chan bool, 1)
	dataProcessingErrorChannel := make(chan bool, 1)

	readerDownstreamDoneChannel := make(chan bool, 1)

	// Start the reader, transformer, and writer. Let them run until fatal input
	// error or end-of-processing happens.
	bufferedOutputStream := bufio.NewWriter(outputStream)

	go recordReader.Read(fileNames, *initialContext, readerChannel, inputErrorChannel, readerDownstreamDoneChannel)
	go transformers.ChainTransformer(readerChannel, readerDownstreamDoneChannel, recordTransformers,
		writerChannel, options)
	go output.ChannelWriter(writerChannel, recordWriter, &options.WriterOptions, doneWritingChannel,
		dataProcessingErrorChannel, bufferedOutputStream, outputIsStdout)

	var retval error
	done := false
	for !done {
		select {
		case ierr := &lt;-inputErrorChannel:
			retval = ierr
			break
		case _ = &lt;-dataProcessingErrorChannel:
			retval = errors.New("exiting due to data error") // details already printed
			break
		case _ = &lt;-doneWritingChannel:
			done = true
			break
		}
	}

	bufferedOutputStream.Flush()

	return retval
}

func main() {
	err := convert_csv_to_json(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
</pre>

<pre class="pre-non-highlight-non-pair">
host,status
apoapsis.east.our.org,up
nadir.west.our.org,down
</pre>

```
$ go build main2.go
$ ./main2 data/hostnames.csv
{"host": "apoapsis.east.our.org", "status": "up"}
{"host": "nadir.west.our.org", "status": "down"}
```



