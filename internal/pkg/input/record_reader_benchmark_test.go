package input

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnkerl/miller/internal/pkg/cli"
)

// go test -run=nonesuch -bench=. github.com/johnkerl/miller/internal/pkg/input/...

func BenchmarkDKVPParse(b *testing.B) {
	readerOptions := &cli.TReaderOptions{
		InputFileFormat: "dkvp",
		IFS:             ",",
		IPS:             "=",
		IRS:             "\n",
	}
	reader, err := NewRecordReaderDKVP(readerOptions, 1)
	assert.Nil(b, err)

	for i := 0; i < b.N; i++ {
		_, _ = recordFromDKVPLine(
			reader,
			"color=yellow,shape=triangle,flag=true,k=1,index=11,quantity=43.6498,rate=9.8870",
		)
	}
}

func BenchmarkNIDXParse(b *testing.B) {
	readerOptions := &cli.TReaderOptions{
		InputFileFormat: "nidx",
		IFS:             " ",
		AllowRepeatIFS:  true,
		IRS:             "\n",
	}
	reader, err := NewRecordReaderNIDX(readerOptions, 1)
	assert.Nil(b, err)

	for i := 0; i < b.N; i++ {
		_, _ = recordFromDKVPLine(
			reader,
			"yellow triangle true 1 11 43.6498 9.8870",
		)
	}
}

func BenchmarkXTABParse(b *testing.B) {
	readerOptions := &cli.TReaderOptions{
		InputFileFormat: "xtab",
		IPS:             " ",
		IFS:             "\n",
		IRS:             "\n",
	}
	reader, err := NewRecordReaderXTAB(readerOptions, 1)
	assert.Nil(b, err)

	stanza := newStanza()
	stanza.dataLines.PushBack("color    yellow")
	stanza.dataLines.PushBack("shape    triangle")
	stanza.dataLines.PushBack("flag     true")
	stanza.dataLines.PushBack("k        1")
	stanza.dataLines.PushBack("index    11")
	stanza.dataLines.PushBack("quantity 43.6498")
	stanza.dataLines.PushBack("rate     9.8870")

	for i := 0; i < b.N; i++ {
		_, _ = reader.recordFromXTABLines(stanza.dataLines)
	}
}
