package types

// Since Go is concurrent, the context struct (AWK-like variables such as
// FILENAME, NF, NF, FNR, etc.) needs to be duplicated and passed through the
// channels along with each record.
type RecordAndContext struct {
	Record  *Mlrmap
	Context Context
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
		Context: *context,
	}
}

// ----------------------------------------------------------------
type Context struct {
	FILENAME string
	FILENUM  int64

	NF  int64
	NR  int64
	FNR int64

	IPS string
	IFS string
	IRS string

	OPS string
	OFS string
	ORS string
}

func NewContext() *Context {
	return &Context{
		FILENAME: "(stdin)",
		FILENUM:  0,

		NF:  0,
		NR:  0,
		FNR: 0,

		IPS: "=",
		IFS: ",",
		IRS: "\n",

		OPS: "=",
		OFS: ",",
		ORS: "\n",
	}
}

// For the record-readers to update their initial context as each new file is opened.
func (this *Context) UpdateForStartOfFile(filename string) {
	this.FILENAME = filename
	this.FILENUM++
	this.FNR = 0
}

// For the record-readers to update their initial context as each new record is read.
func (this *Context) UpdateForInputRecord(inrec *Mlrmap) {
	if inrec != nil { // do not count the end-of-stream marker which is a nil record pointer
		this.NF = int64(inrec.FieldCount)
		this.NR++
		this.FNR++
	}
}
