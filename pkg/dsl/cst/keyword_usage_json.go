// Machine-readable (JSON) accessors over the keyword catalog, for
// `mlr help --as-json` and similar tooling. The per-keyword usage functions print
// their bodies directly to stdout (they predate any structured-help need), so
// here we capture that output by temporarily redirecting os.Stdout.

package cst

import (
	"bytes"
	"io"
	"os"
	"strings"
)

// KeywordInfoForJSON is the structured view of a single DSL keyword.
type KeywordInfoForJSON struct {
	Name string `json:"name"`
	Help string `json:"help"`
}

// captureStdout runs f with os.Stdout redirected to a pipe and returns whatever
// f printed. The keyword usage functions write via fmt.Println/Printf, which
// resolve os.Stdout at call time, so swapping it here captures their output.
// Help generation is single-threaded and one-shot, so the global swap is safe.
func captureStdout(f func()) string {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return ""
	}
	os.Stdout = w

	done := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()

	f()

	_ = w.Close()
	os.Stdout = old
	s := <-done
	_ = r.Close()
	return s
}

func makeKeywordInfoForJSON(entry *tKeywordUsageEntry) *KeywordInfoForJSON {
	return &KeywordInfoForJSON{
		Name: entry.name,
		Help: strings.TrimRight(captureStdout(entry.usageFunc), "\n"),
	}
}

// GetKeywordInfosForJSON returns the full keyword catalog in source-table order.
func GetKeywordInfosForJSON() []*KeywordInfoForJSON {
	infos := make([]*KeywordInfoForJSON, 0, len(KEYWORD_USAGE_TABLE))
	for i := range KEYWORD_USAGE_TABLE {
		infos = append(infos, makeKeywordInfoForJSON(&KEYWORD_USAGE_TABLE[i]))
	}
	return infos
}

// GetKeywordInfoForJSON returns the structured view of a single keyword, or nil
// if there is no such keyword.
func GetKeywordInfoForJSON(name string) *KeywordInfoForJSON {
	for i := range KEYWORD_USAGE_TABLE {
		if KEYWORD_USAGE_TABLE[i].name == name {
			return makeKeywordInfoForJSON(&KEYWORD_USAGE_TABLE[i])
		}
	}
	return nil
}
