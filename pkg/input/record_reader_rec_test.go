package input

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnkerl/miller/v6/pkg/cli"
)

func newTestRECReader(t *testing.T) *RecordReaderREC {
	readerOptions := cli.DefaultReaderOptions()
	err := cli.FinalizeReaderOptions(&readerOptions)
	assert.Nil(t, err)
	reader, err := NewRecordReaderREC(&readerOptions, 1)
	assert.NotNil(t, reader)
	assert.Nil(t, err)
	return reader
}

func TestRecordFromRECLinesBasic(t *testing.T) {
	reader := newTestRECReader(t)

	record, err := reader.recordFromRECLines([]string{
		"name: John Doe",
		"email: jdoe@example.com",
	})
	assert.Nil(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, int64(2), record.FieldCount)
	assert.Equal(t, "name", record.Head.Key)
	assert.Equal(t, "John Doe", record.Head.Value.String())
	assert.Equal(t, "email", record.Head.Next.Key)
	assert.Equal(t, "jdoe@example.com", record.Head.Next.Value.String())
}

func TestRecordFromRECLinesPlusContinuation(t *testing.T) {
	reader := newTestRECReader(t)

	record, err := reader.recordFromRECLines([]string{
		"notes: line one",
		"+ line two",
		"+line three",
	})
	assert.Nil(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, int64(1), record.FieldCount)
	assert.Equal(t, "notes", record.Head.Key)
	assert.Equal(t, "line one\nline two\nline three", record.Head.Value.String())
}

func TestRecordFromRECLinesBarePlusContinuation(t *testing.T) {
	reader := newTestRECReader(t)

	record, err := reader.recordFromRECLines([]string{
		"notes: line one",
		"+",
		"+ line three",
	})
	assert.Nil(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, "line one\n\nline three", record.Head.Value.String())
}

func TestRecordFromRECLinesBackslashContinuation(t *testing.T) {
	reader := newTestRECReader(t)

	record, err := reader.recordFromRECLines([]string{
		`name: hello \`,
		`world`,
	})
	assert.Nil(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, int64(1), record.FieldCount)
	assert.Equal(t, "hello world", record.Head.Value.String())
}

func TestRecordFromRECLinesColonInValue(t *testing.T) {
	reader := newTestRECReader(t)

	record, err := reader.recordFromRECLines([]string{
		"url: http://example.com: see notes",
	})
	assert.Nil(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, "url", record.Head.Key)
	assert.Equal(t, "http://example.com: see notes", record.Head.Value.String())
}

func TestRecordFromRECLinesDedupeFieldNames(t *testing.T) {
	reader := newTestRECReader(t)

	record, err := reader.recordFromRECLines([]string{
		"a: 1",
		"b: 2",
		"b: 3",
	})
	assert.Nil(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, int64(3), record.FieldCount)
	assert.Equal(t, "a", record.Head.Key)
	assert.Equal(t, "b", record.Head.Next.Key)
	assert.Equal(t, "b_2", record.Head.Next.Next.Key)
}

func TestRecordFromRECLinesOrphanPlusIsError(t *testing.T) {
	reader := newTestRECReader(t)

	record, err := reader.recordFromRECLines([]string{
		"+ oops",
		"name: x",
	})
	assert.Nil(t, record)
	assert.NotNil(t, err)
}

func TestRecordFromRECLinesMissingColonSpaceIsError(t *testing.T) {
	reader := newTestRECReader(t)

	record, err := reader.recordFromRECLines([]string{
		"this line has no colon-space separator",
	})
	assert.Nil(t, record)
	assert.NotNil(t, err)
}

func TestRecordFromRECLinesMissingColonSpaceBareColonIsError(t *testing.T) {
	reader := newTestRECReader(t)

	// Per the recutils spec, a non-empty value must be separated from the
	// field name by ": " (colon-space) -- a colon directly followed by
	// value text with no separating space is malformed. (A colon followed
	// by nothing at all is a different, valid case: an empty value -- see
	// TestRecordFromRECLinesBareColonEmptyValue.)
	record, err := reader.recordFromRECLines([]string{
		"Foo:bar",
	})
	assert.Nil(t, record)
	assert.NotNil(t, err)
}

func TestJoinRECBackslashContinuations(t *testing.T) {
	assert.Equal(t,
		[]string{"ab"},
		joinRECBackslashContinuations([]string{`a\`, "b"}),
	)
	assert.Equal(t,
		[]string{"a", "b"},
		joinRECBackslashContinuations([]string{"a", "b"}),
	)
	// A trailing backslash with nothing left to join to keeps its
	// (backslash-stripped) content as-is.
	assert.Equal(t,
		[]string{"a"},
		joinRECBackslashContinuations([]string{`a\`}),
	)
}

func TestParseRECFields(t *testing.T) {
	fields, err := parseRECFields([]string{"a: 1", "+ 2", "+3"})
	assert.Nil(t, err)
	assert.Equal(t, []tRECField{{key: "a", value: "1\n2\n3"}}, fields)

	_, err = parseRECFields([]string{"+ orphan"})
	assert.NotNil(t, err)
}

func TestRecordFromRECLinesBareColonEmptyValue(t *testing.T) {
	reader := newTestRECReader(t)

	// Per the recutils spec, a field's value may be empty: a colon with
	// nothing after it (not even a space) is valid, distinct from the
	// "Foo:bar" case above which is missing the required separator space.
	record, err := reader.recordFromRECLines([]string{
		"notes:",
	})
	assert.Nil(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, int64(1), record.FieldCount)
	assert.Equal(t, "notes", record.Head.Key)
	assert.Equal(t, "", record.Head.Value.String())
}

func TestRecordFromRECLinesBareColonPlusContinuation(t *testing.T) {
	reader := newTestRECReader(t)

	// This is the idiomatic recutils "%doc:" pattern from the official
	// manual's tutorial example: an empty-valued field whose actual value
	// is supplied entirely by the "+"-continuation lines that follow. The
	// continuation becomes the value outright, with no spurious leading
	// "\n" from folding onto an empty base value.
	record, err := reader.recordFromRECLines([]string{
		"%doc:",
		"+ A book in my personal collection.",
	})
	assert.Nil(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, int64(1), record.FieldCount)
	assert.Equal(t, "%doc", record.Head.Key)
	assert.Equal(t, "A book in my personal collection.", record.Head.Value.String())
}
