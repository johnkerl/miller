package input

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnkerl/miller/v6/pkg/cli"
)

// go test -run=nonesuch -bench=. github.com/johnkerl/miller/v6/pkg/input/...

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
	stanza.dataLines = append(stanza.dataLines, "color    yellow")
	stanza.dataLines = append(stanza.dataLines, "shape    triangle")
	stanza.dataLines = append(stanza.dataLines, "flag     true")
	stanza.dataLines = append(stanza.dataLines, "k        1")
	stanza.dataLines = append(stanza.dataLines, "index    11")
	stanza.dataLines = append(stanza.dataLines, "quantity 43.6498")
	stanza.dataLines = append(stanza.dataLines, "rate     9.8870")

	for i := 0; i < b.N; i++ {
		_, _ = reader.recordFromXTABLines(stanza.dataLines)
	}
}
