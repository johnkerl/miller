package runtime

import (
	"miller/containers"
)

// ----------------------------------------------------------------
type LrecAndContext struct {
	Lrec    *containers.Lrec
	Context Context
}

func NewLrecAndContext(
	lrec *containers.Lrec,
	context *Context,
) *LrecAndContext {
	return &LrecAndContext{
		lrec,
		*context,
	}
}

// ----------------------------------------------------------------
// xxx comment about who writes this and who reads this
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

func (this *Context) UpdateForStartOfFile(filename string) {
	this.FILENAME = filename
	this.FILENUM++
	this.FNR = 0
}

func (this *Context) UpdateForInputRecord(inrec *containers.Lrec) {
	if inrec != nil { // do not count the end-of-stream marker which is a nil record pointer
		this.NF = inrec.FieldCount
		this.NR++
		this.FNR++
	}
}
