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
// This is an example of using Miller as a library.
package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"

	"github.com/johnkerl/miller/pkg/bifs"
	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/input"
	"github.com/johnkerl/miller/pkg/output"
	"github.com/johnkerl/miller/pkg/types"
)

// Put your record-processing logic here.
func custom_record_processor(irac *types.RecordAndContext) (*types.RecordAndContext, error) {
	irec := irac.Record

	v := irec.Get("i")
	if v == nil {
		return nil, fmt.Errorf("did not find key \"i\" at filename %s record number %d",
			irac.Context.FILENAME, irac.Context.FNR,
		)
	}
	v2 := bifs.BIF_times(v, v)
	irec.PutReference("i2", v2)

	return irac, nil
}

// Put your various options here.
func custom_options() *cli.TOptions {
	return &cli.TOptions{
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
}

// This function you don't need to modify.
func run_custom_processor(
	fileNames []string,
	options *cli.TOptions,
	record_processor func (irac *types.RecordAndContext) (*types.RecordAndContext, error),
) error {
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

	// Set up the channels for the record-reader.
	readerChannel := make(chan *list.List, 2) // list of *types.RecordAndContext
	inputErrorChannel := make(chan error, 1)
	// Not needed in this example
	readerDownstreamDoneChannel := make(chan bool, 1)

	// Instantiate the record-writer
	recordWriter, err := output.Create(&options.WriterOptions)
	if err != nil {
		return err
	}
	bufferedOutputStream := bufio.NewWriter(outputStream)

	// Start the record-reader.
	go recordReader.Read(
		fileNames, *initialContext, readerChannel, inputErrorChannel, readerDownstreamDoneChannel)

	// Loop through the record stream.
	var retval error
	done := false
	for !done {
		select {

		case ierr := &lt;-inputErrorChannel:
			retval = ierr
			break

		case iracs := &lt;-readerChannel:
			// Handle the record batch
			for e := iracs.Front(); e != nil; e = e.Next() {
				irac := e.Value.(*types.RecordAndContext)
				if irac.Record != nil {
					orac, err := record_processor(irac)
					if err != nil {
						retval = err
						done = true
						break
					}
					recordWriter.Write(orac.Record, bufferedOutputStream, outputIsStdout)
				}
				if irac.OutputString != "" {
					fmt.Fprintln(bufferedOutputStream, irac.OutputString)
				}
				if irac.EndOfStream {
					done = true
				}
			}
			break

		}
	}

	bufferedOutputStream.Flush()

	return retval
}

func main() {
	options := custom_options()
	err := run_custom_processor(os.Args[1:], options, custom_record_processor)
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
{"a": "pan", "b": "pan", "i": 1, "x": 0.3467901443380824, "y": 0.7268028627434533, "i2": 1}
{"a": "eks", "b": "pan", "i": 2, "x": 0.7586799647899636, "y": 0.5221511083334797, "i2": 4}
{"a": "wye", "b": "wye", "i": 3, "x": 0.20460330576630303, "y": 0.33831852551664776, "i2": 9}
{"a": "eks", "b": "wye", "i": 4, "x": 0.38139939387114097, "y": 0.13418874328430463, "i2": 16}
{"a": "wye", "b": "pan", "i": 5, "x": 0.5732889198020006, "y": 0.8636244699032729, "i2": 25}$ ./main2 data/small.csv
```
