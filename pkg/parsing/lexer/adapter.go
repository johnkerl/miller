// adapter.go provides NewLexer([]byte) for compatibility with pkg/dsl/cst/root.go.
// The PGPG-generated lexer uses NewMlrLexer(io.Reader); this adapter accepts []byte.

package lexer

import (
	"bytes"

	liblexers "github.com/johnkerl/pgpg/go/lib/pkg/lexers"
)

// NewLexer creates a lexer from a byte slice. This matches the API expected by
// pkg/dsl/cst/root.go (same signature as the former PGPG lexer).
func NewLexer(src []byte) liblexers.AbstractLexer {
	return NewMlrLexer(bytes.NewReader(src))
}
