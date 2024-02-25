// This file contains the interface for file-format-specific record-readers, as
// well as a collection of utility functions.

package input

import (
	"regexp"
	"strings"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/lib"
)

// IPairSplitter splits a string into left and right, e.g. for IPS.
// This helps us reuse code for splitting by IPS string, or IPS regex.
type iPairSplitter interface {
	Split(input string) []string
}

func newPairSplitter(options *cli.TReaderOptions) iPairSplitter {
	if options.IPSRegex == nil {
		return &tIPSSplitter{ips: options.IPS}
	} else {
		return &tIPSRegexSplitter{ipsRegex: options.IPSRegex}
	}
}

type tIPSSplitter struct {
	ips string
}

func (s *tIPSSplitter) Split(input string) []string {
	return strings.SplitN(input, s.ips, 2)
}

type tIPSRegexSplitter struct {
	ipsRegex *regexp.Regexp
}

func (s *tIPSRegexSplitter) Split(input string) []string {
	return lib.RegexCompiledSplitString(s.ipsRegex, input, 2)
}

// IFieldSplitter splits a string into pieces, e.g. for IFS.
// This helps us reuse code for splitting by IFS string, or IFS regex.
type iFieldSplitter interface {
	Split(input string) []string
}

func newFieldSplitter(options *cli.TReaderOptions) iFieldSplitter {
	if options.IFSRegex == nil {
		return &tIFSSplitter{ifs: options.IFS, allowRepeatIFS: options.AllowRepeatIFS}
	} else {
		return &tIFSRegexSplitter{ifsRegex: options.IFSRegex}
	}
}

type tIFSSplitter struct {
	ifs            string
	allowRepeatIFS bool
}

func (s *tIFSSplitter) Split(input string) []string {
	fields := lib.SplitString(input, s.ifs)
	if s.allowRepeatIFS {
		fields = lib.StripEmpties(fields) // left/right trim
	}
	return fields
}

type tIFSRegexSplitter struct {
	ifsRegex *regexp.Regexp
}

func (s *tIFSRegexSplitter) Split(input string) []string {
	return lib.RegexCompiledSplitString(s.ifsRegex, input, -1)
}
