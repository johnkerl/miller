package types

import (
	"bytes"
	"strconv"
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
func (rac *RecordAndContext) Copy() *RecordAndContext {
	if rac == nil {
		return nil
	}
	recordCopy := rac.Record.Copy()
	contextCopy := rac.Context
	return &RecordAndContext{
		Record:       recordCopy,
		Context:      contextCopy,
		OutputString: "",
		EndOfStream:  false,
	}
}

// For print/dump/etc to insert strings sequenced into the record-output
// stream.  This avoids race conditions between different goroutines printing
// to stdout: we have a single designated goroutine printing to stdout. This
// makes output more predictable and intuitive for users; it also makes our
// regression tests run reliably the same each time.
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

	OPS     string
	OFS     string
	ORS     string
	FLATSEP string
}

// TODO: comment: Remember command-line values to pass along to CST evaluators.
// The options struct-pointer can be nil when invoked by non-DSL verbs such as
// join or seqgen.
func NewContext(
	IPS string,
	IFS string,
	IRS string,

	OPS string,
	OFS string,
	ORS string,
	FLATSEP string,
) *Context {
	context := &Context{
		FILENAME: "(stdin)",
		FILENUM:  0,

		NR:  0,
		FNR: 0,

		IPS: "=",
		IFS: ",",
		IRS: "\n",

		OPS:     "=",
		OFS:     ",",
		ORS:     "\n",
		FLATSEP: ".",
	}

	// TODO: FILENAME/FILENUM/NR/FNR should be in one struct, and the rest in
	// another. The former vary per record; the latter are command-line-driven
	// and do not vary per record. All they have in common is they are
	// awk-like context-variables.
	context.IPS = IPS
	context.IFS = IFS
	context.IRS = IRS

	context.OPS = OPS
	context.OFS = OFS
	context.ORS = ORS
	context.FLATSEP = FLATSEP

	return context
}

// TODO: comment: Remember command-line values to pass along to CST evaluators.
// The options struct-pointer can be nil when invoked by non-DSL verbs such as
// join or seqgen.
func NewNilContext() *Context { // TODO: rename
	context := &Context{
		FILENAME: "(stdin)",
		FILENUM:  0,

		NR:  0,
		FNR: 0,

		IPS: "=",
		IFS: ",",
		IRS: "\n",

		OPS:     "=",
		OFS:     ",",
		ORS:     "\n",
		FLATSEP: ".",
	}

	return context
}

// For the record-readers to update their initial context as each new file is opened.
func (context *Context) UpdateForStartOfFile(filename string) {
	context.FILENAME = filename
	context.FILENUM++
	context.FNR = 0
}

// For the record-readers to update their initial context as each new record is read.
func (context *Context) UpdateForInputRecord() {
	context.NR++
	context.FNR++
}

func (context *Context) Copy() *Context {
	other := *context
	return &other
}

func (context *Context) GetStatusString() string {

	var buffer bytes.Buffer // stdio is non-buffered in Go, so buffer for speed increase
	buffer.WriteString("FILENAME=\"")
	buffer.WriteString(context.FILENAME)

	buffer.WriteString("\",FILENUM=")
	buffer.WriteString(strconv.Itoa(context.FILENUM))

	buffer.WriteString(",NR=")
	buffer.WriteString(strconv.Itoa(context.NR))

	buffer.WriteString(",FNR=")
	buffer.WriteString(strconv.Itoa(context.FNR))

	return buffer.String()
}
