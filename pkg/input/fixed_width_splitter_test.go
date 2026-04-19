package input

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestReferenceSplitter(t *testing.T) {

	type testCase struct {
		scenario            string
		referenceRow        string
		line                string
		expectedLeftSingle  []string
		expectedLeftMulti   []string
		expectedRightSingle []string
		expectedRightMulti  []string
	}

	testCases := []testCase{
		{
			scenario:            "base left",
			referenceRow:        "Name      Place    Thing",
			line:                "Name      Place    Thing",
			expectedLeftSingle:  []string{"Name      ", "Place    ", "Thing"},
			expectedLeftMulti:   []string{"Name      ", "Place    ", "Thing"},
			expectedRightSingle: []string{"Name", "      Place", "    Thing"},
			expectedRightMulti:  []string{"Name", "      Place", "    Thing"},
		},
		{
			scenario:            "base right",
			referenceRow:        "    Name      Place    Thing",
			line:                "    Name      Place    Thing",
			expectedLeftSingle:  []string{"    ", "Name      ", "Place    ", "Thing"},
			expectedLeftMulti:   []string{"    ", "Name      ", "Place    ", "Thing"},
			expectedRightSingle: []string{"    Name", "      Place", "    Thing"},
			expectedRightMulti:  []string{"    Name", "      Place", "    Thing"},
		},
		{
			scenario:            "left with data row",
			referenceRow:        "Name      Place    Thing",
			line:                "JohnDoe   Nyc      Bottle",
			expectedLeftSingle:  []string{"JohnDoe   ", "Nyc      ", "Bottle"},
			expectedLeftMulti:   []string{"JohnDoe   ", "Nyc      ", "Bottle"},
			expectedRightSingle: []string{"John", "Doe   Nyc  ", "    Bottle"},
			expectedRightMulti:  []string{"John", "Doe   Nyc  ", "    Bottle"},
		},
		{
			scenario:            "left multi word",
			referenceRow:        "Name      Last Seen     Thing",
			line:                "Name      Last Seen     Thing",
			expectedLeftSingle:  []string{"Name      ", "Last ", "Seen     ", "Thing"},
			expectedLeftMulti:   []string{"Name      ", "Last Seen     ", "Thing"},
			expectedRightSingle: []string{"Name", "      Last", " Seen", "     Thing"},
			expectedRightMulti:  []string{"Name", "      Last Seen", "     Thing"},
		},
		{
			scenario:            "right multi word",
			referenceRow:        "    Name     Last Seen   Thing",
			line:                "    Name     Last Seen   Thing",
			expectedLeftSingle:  []string{"    ", "Name     ", "Last ", "Seen   ", "Thing"},
			expectedLeftMulti:   []string{"    ", "Name     ", "Last ", "Seen   ", "Thing"},
			expectedRightSingle: []string{"    Name", "     Last", " Seen", "   Thing"},
			expectedRightMulti:  []string{"    Name", "     Last Seen", "   Thing"},
		},
		{
			scenario:            "left multi word data row",
			referenceRow:        "Name      Last Seen     Thing",
			line:                "Max       two days ago  Bottle",
			expectedLeftSingle:  []string{"Max       ", "two d", "ays ago  ", "Bottle"},
			expectedLeftMulti:   []string{"Max       ", "two days ago  ", "Bottle"},
			expectedRightSingle: []string{"Max ", "      two ", "days ", "ago  Bottle"},
			expectedRightMulti:  []string{"Max ", "      two days ", "ago  Bottle"},
		},
		{
			scenario:            "right with data row",
			referenceRow:        "    Name    Place   Thing",
			line:                " JohnDoe  NewYork  Bottle",
			expectedLeftSingle:  []string{" Joh", "nDoe  Ne", "wYork  B", "ottle"},
			expectedLeftMulti:   []string{" Joh", "nDoe  Ne", "wYork  B", "ottle"},
			expectedRightSingle: []string{" JohnDoe", "  NewYork", "  Bottle"},
			expectedRightMulti:  []string{" JohnDoe", "  NewYork", "  Bottle"},
		},
		{
			scenario:            "right multi word data row",
			referenceRow:        "   Name     Last Seen   Thing",
			line:                "JohnDoe  two days ago  Bottle",
			expectedLeftSingle:  []string{"Joh", "nDoe  two", " days", " ago  B", "ottle"},
			expectedLeftMulti:   []string{"Joh", "nDoe  two", " days ago  B", "ottle"},
			expectedRightSingle: []string{"JohnDoe", "  two day", "s ago", "  Bottle"},
			expectedRightMulti:  []string{"JohnDoe", "  two days ago", "  Bottle"},
		},
		{
			scenario:            "single column",
			referenceRow:        "Name",
			line:                "Blah Blah Blah",
			expectedLeftSingle:  []string{"Blah Blah Blah"},
			expectedLeftMulti:   []string{"Blah Blah Blah"},
			expectedRightSingle: []string{"Blah Blah Blah"},
			expectedRightMulti:  []string{"Blah Blah Blah"},
		},
		{
			scenario:            "empty row has one column",
			referenceRow:        "",
			line:                "anything",
			expectedLeftSingle:  []string{"anything"},
			expectedLeftMulti:   []string{"anything"},
			expectedRightSingle: []string{"anything"},
			expectedRightMulti:  []string{"anything"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			{
				indexes := parseLeftAlign(tc.referenceRow, false)
				fields := split(tc.line, indexes)
				assert.Equal(t, tc.expectedLeftSingle, fields)
				assert.Equal(t, tc.line, strings.Join(fields, ""))
			}
			{
				indexes := parseLeftAlign(tc.referenceRow, true)
				fields := split(tc.line, indexes)
				assert.Equal(t, tc.expectedLeftMulti, fields)
				assert.Equal(t, tc.line, strings.Join(fields, ""))
			}
			{
				indexes := parseRightAlign(tc.referenceRow, false)
				fields := split(tc.line, indexes)
				assert.Equal(t, tc.expectedRightSingle, fields)
				assert.Equal(t, tc.line, strings.Join(fields, ""))
			}
			{
				indexes := parseRightAlign(tc.referenceRow, true)
				fields := split(tc.line, indexes)
				assert.Equal(t, tc.expectedRightMulti, fields)
				assert.Equal(t, tc.line, strings.Join(fields, ""))
			}
		})
	}
}

func TestSplit(t *testing.T) {

	type testCase struct {
		scenario string
		indexes  []int
		line     string
		expected []string
	}

	testCases := []testCase{
		{
			scenario: "normal split",
			indexes:  []int{4, 8, 13},
			line:     "abc123defghij",
			expected: []string{"abc1", "23de", "fghij"},
		},
		{
			scenario: "split with rest",
			indexes:  []int{4, 8},
			line:     "abc123defghij",
			expected: []string{"abc1", "23de", "fghij"},
		},
		{
			scenario: "split shorter string",
			indexes:  []int{4, 8, 13},
			line:     "abc123",
			expected: []string{"abc1", "23"},
		},
		{
			scenario: "empty indexes returns everything as one entry",
			indexes:  []int{},
			line:     "abc123defghij",
			expected: []string{"abc123defghij"},
		},
		{
			scenario: "empty line with valid indexes returns nil",
			indexes:  []int{2, 4},
			line:     "",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			result := split(tc.line, tc.indexes)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestWidthSplitter(t *testing.T) {

	type testCase struct {
		scenario    string
		widthsStr   string
		line        string
		expected    []string
		expectedErr string
	}

	testCases := []testCase{
		{
			scenario:    "empty widths spec",
			widthsStr:   "",
			line:        "blah blah",
			expected:    []string{"blah blah"},
			expectedErr: "",
		},
		{
			scenario:    "simple",
			widthsStr:   "5,5,6",
			line:        "blah blah hello",
			expected:    []string{"blah ", "blah ", "hello"},
			expectedErr: "",
		},
		{
			scenario:    "short line",
			widthsStr:   "5,5,6",
			line:        "blah blah hi",
			expected:    []string{"blah ", "blah ", "hi"},
			expectedErr: "",
		},
		{
			scenario:    "long line",
			widthsStr:   "5,5,2",
			line:        "blah blah hello",
			expected:    []string{"blah ", "blah ", "he", "llo"},
			expectedErr: "",
		},
		{
			scenario:    "non numeric width",
			widthsStr:   "a,b,c",
			line:        "blah blah",
			expected:    nil,
			expectedErr: `invalid width: a, error: strconv.Atoi: parsing "a": invalid syntax`,
		},
		{
			scenario:    "zero width",
			widthsStr:   "0,1,1",
			line:        "blah blah",
			expected:    nil,
			expectedErr: `not a positive width: 0`,
		},
		{
			scenario:    "negative width",
			widthsStr:   "1,-1,1",
			line:        "blah blah",
			expected:    nil,
			expectedErr: `not a positive width: -1`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			indexes, err := parseWidths(tc.widthsStr)
			if err != nil {
				assert.Equal(t, tc.expectedErr, err.Error())
			} else {
				fields := split(tc.line, indexes)
				assert.Equal(t, tc.expected, fields)
				assert.Equal(t, tc.line, strings.Join(fields, ""))
			}
		})
	}
}
