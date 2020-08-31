package runtime

import (
	"miller/containers"
)

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
}

func (this *Context) UpdateForInputRecord(inrec *containers.Lrec) {
	if inrec != nil { // do not count the end-of-stream marker which is a nil record pointer
		this.NF = inrec.FieldCount
		this.NR++
		this.FNR++
	}
}
