package input

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnkerl/miller/internal/pkg/cli"
)

func TestRecordFromDKVPLine(t *testing.T) {
	readerOptions := cli.DefaultReaderOptions()
	cli.FinalizeReaderOptions(&readerOptions) // compute IPS, IFS -> IPSRegex, IFSRegex
	reader, err := NewRecordReaderDKVP(&readerOptions, 1)
	assert.NotNil(t, reader)
	assert.Nil(t, err)

	line := ""
	record, err := recordFromDKVPLine(reader, line)
	assert.NotNil(t, record)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), record.FieldCount)

	line = "a=1,b=2,c=3"
	record, err = recordFromDKVPLine(reader, line)
	assert.NotNil(t, record)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), record.FieldCount)

	assert.NotNil(t, record.Head)
	assert.NotNil(t, record.Head.Next)
	assert.NotNil(t, record.Head.Next.Next)
	assert.Nil(t, record.Head.Next.Next.Next)
	assert.Equal(t, record.Head.Key, "a")
	assert.Equal(t, record.Head.Next.Key, "b")
	assert.Equal(t, record.Head.Next.Next.Key, "c")

	// Default is to dedupe to a=1,b=2,b_2=3
	line = "a=1,b=2,b=3"
	record, err = recordFromDKVPLine(reader, line)
	assert.NotNil(t, record)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), record.FieldCount)

	assert.NotNil(t, record.Head)
	assert.NotNil(t, record.Head.Next)
	assert.NotNil(t, record.Head.Next.Next)
	assert.Nil(t, record.Head.Next.Next.Next)
	assert.Equal(t, record.Head.Key, "a")
	assert.Equal(t, record.Head.Next.Key, "b")
	assert.Equal(t, record.Head.Next.Next.Key, "b_2")

	line = "a,b,c"
	record, err = recordFromDKVPLine(reader, line)
	assert.NotNil(t, record)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), record.FieldCount)

	assert.NotNil(t, record.Head)
	assert.NotNil(t, record.Head.Next)
	assert.NotNil(t, record.Head.Next.Next)
	assert.Nil(t, record.Head.Next.Next.Next)
	assert.Equal(t, record.Head.Key, "1")
	assert.Equal(t, record.Head.Next.Key, "2")
	assert.Equal(t, record.Head.Next.Next.Key, "3")
}
