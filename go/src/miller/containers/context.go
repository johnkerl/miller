package containers

// Since Go is concurrent, the context struct (AWK-like variables such as
// FILENAME, NF, NF, FNR, etc.) needs to be duplicated and passed through the
// channels along with each record.
type LrecAndContext struct {
	Lrec    *Lrec
	Context Context
}

func NewLrecAndContext(
	lrec *Lrec,
	context *Context,
) *LrecAndContext {
	return &LrecAndContext{
		Lrec: lrec,
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
		"(stdin)",
		0,

		0,
		0,
		0,

		"=",
		",",
		"\n",

		"=",
		",",
		"\n",
	}
}

// For the record-readers to update their initial context as each new file is opened.
func (this *Context) UpdateForStartOfFile(filename string) {
	this.FILENAME = filename
	this.FILENUM++
	this.FNR = 0
}

// For the record-readers to update their initial context as each new record is read.
func (this *Context) UpdateForInputRecord(inrec *Lrec) {
	if inrec != nil { // do not count the end-of-stream marker which is a nil record pointer
		this.NF = inrec.FieldCount
		this.NR++
		this.FNR++
	}
}
