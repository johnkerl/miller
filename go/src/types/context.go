package types

import (
	"bytes"
	"strconv"

	"miller/src/cliutil"
)

// Since Go is concurrent, the context struct (AWK-like variables such as
// FILENAME, NF, NR, FNR, etc.) needs to be duplicated and passed through the
// channels along with each record.
//
// Strings to be printed from put/filter DSL print/dump/etc statements are
// passed along to the output channel via this OutputString rather than
// fmt.Println directly in the put/filter handlers since we want all print
// statements and record-output to be in the same goroutine, for deterministic
// output ordering.

type RecordAndContext struct {
	Record       *Mlrmap
	Context      Context
	OutputString string
	EndOfStream  bool
}

func NewRecordAndContext(
	record *Mlrmap,
	context *Context,
) *RecordAndContext {
	return &RecordAndContext{
		Record: record,
		// Since Go is concurrent, the context struct needs to be duplicated and
		// passed through the channels along with each record. Here is where
		// the copy happens, via the '*' in *context.
		Context:      *context,
		OutputString: "",
		EndOfStream:  false,
	}
}

// For the record-readers to update their initial context as each new record is read.
func (this *RecordAndContext) Copy() *RecordAndContext {
	if this == nil {
		return nil
	}
	recordCopy := this.Record.Copy()
	contextCopy := this.Context
	return &RecordAndContext{
		Record:       recordCopy,
		Context:      contextCopy,
		OutputString: "",
		EndOfStream:  false,
	}
}

// For print/dump/etc to insert strings sequenced into the record-output stream.
func NewOutputString(
	outputString string,
	context *Context,
) *RecordAndContext {
	return &RecordAndContext{
		Record:       nil,
		Context:      *context,
		OutputString: outputString,
		EndOfStream:  false,
	}
}

// For the record-readers to update their initial context as each new record is read.
func NewEndOfStreamMarker(context *Context) *RecordAndContext {
	return &RecordAndContext{
		Record:       nil,
		Context:      *context,
		OutputString: "",
		EndOfStream:  true,
	}
}

// ----------------------------------------------------------------
type Context struct {
	FILENAME string
	FILENUM  int

	// This is computed dynammically from the current record's field-count
	// NF int
	NR  int
	FNR int

	IPS string
	IFS string
	IRS string

	OPS      string
	OFS      string
	ORS      string
	OFLATSEP string
}

func NewContext(options *cliutil.TOptions) *Context {
	context := &Context{
		FILENAME: "(stdin)",
		FILENUM:  0,

		NR:  0,
		FNR: 0,

		IPS: "=",
		IFS: ",",
		IRS: "\n",

		OPS:      "=",
		OFS:      ",",
		ORS:      "\n",
		OFLATSEP: ".",
	}

	// Remember command-line values to pass along to CST evaluators.  The
	// options struct-pointer can be nil when invoked by non-DSL verbs such as
	// join or seqgen.
	//
	// TODO: FILENAME/FILENUM/NR/FNR should be in one struct, and the rest in
	// another. The former vary per record; the latter are command-line-driven
	// and do not vary per record. All they have in common is they are
	// awk-like context-variables.
	if options != nil {
		context.IPS = options.ReaderOptions.IPS
		context.IFS = options.ReaderOptions.IFS
		context.IRS = options.ReaderOptions.IRS

		context.OPS = options.WriterOptions.OPS
		context.OFS = options.WriterOptions.OFS
		context.ORS = options.WriterOptions.ORS
		context.OFLATSEP = options.WriterOptions.OFLATSEP
	}

	return context
}

// For the record-readers to update their initial context as each new file is opened.
func (this *Context) UpdateForStartOfFile(filename string) {
	this.FILENAME = filename
	this.FILENUM++
	this.FNR = 0
}

// For the record-readers to update their initial context as each new record is read.
func (this *Context) UpdateForInputRecord() {
	this.NR++
	this.FNR++
}

func (this *Context) Copy() *Context {
	that := *this
	return &that
}

func (this *Context) GetStatusString() string {

	var buffer bytes.Buffer // 5x faster than fmt.Print() separately
	buffer.WriteString("FILENAME=\"")
	buffer.WriteString(this.FILENAME)

	buffer.WriteString("\",FILENUM=")
	buffer.WriteString(strconv.Itoa(this.FILENUM))

	buffer.WriteString(",NR=")
	buffer.WriteString(strconv.Itoa(this.NR))

	buffer.WriteString(",FNR=")
	buffer.WriteString(strconv.Itoa(this.FNR))

	return buffer.String()
}
