package input

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/types"
)

func TestNewRecordReaderDKVPX(t *testing.T) {
	readerOptions := cli.DefaultReaderOptions()
	readerOptions.InputFileFormat = "dkvpx"
	assert.NoError(t, cli.FinalizeReaderOptions(&readerOptions))

	reader, err := NewRecordReaderDKVPX(&readerOptions, 1)
	assert.NotNil(t, reader)
	assert.NoError(t, err)
}

func TestRecordReaderDKVPX_ReadStdin(t *testing.T) {
	readerOptions := cli.DefaultReaderOptions()
	readerOptions.InputFileFormat = "dkvpx"
	assert.NoError(t, cli.FinalizeReaderOptions(&readerOptions))

	reader, err := NewRecordReaderDKVPX(&readerOptions, 1)
	assert.NoError(t, err)

	ctx := types.Context{}
	readerChannel := make(chan []*types.RecordAndContext, 4)
	errorChannel := make(chan error, 1)

	input := strings.NewReader("x=1,y=2,z=3\n")
	go reader.processHandle(input, "(test)", &ctx, readerChannel, errorChannel, nil)

	records := <-readerChannel
	assert.Len(t, records, 1)
	assert.False(t, records[0].EndOfStream)
	assert.Equal(t, "x", records[0].Record.Head.Key)
	assert.Equal(t, "1", records[0].Record.Head.Value.String())
}
