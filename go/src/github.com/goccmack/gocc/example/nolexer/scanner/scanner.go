package scanner

import (
	"github.com/goccmack/gocc/example/nolexer/token"
)

type Scanner struct {
	src []byte
	pos int
}

func NewString(s string) *Scanner {
	return &Scanner{[]byte(s), 0}
}

func isWhiteSpace(c byte) bool {
	return c == ' ' ||
		c == '\t' ||
		c == '\n' ||
		c == '\r'
}

func (S *Scanner) skipWhiteSpace() {
	for S.pos < len(S.src) && isWhiteSpace(S.src[S.pos]) {
		S.pos++
	}
}

func (S *Scanner) scanId() string {
	pos := S.pos
	for S.pos < len(S.src) && !isWhiteSpace(S.src[S.pos]) {
		S.pos++
	}
	S.pos++
	return string(S.src[pos : S.pos-1])
}

func (S *Scanner) Scan() (tok *token.Token) {
	S.skipWhiteSpace()

	if S.pos >= len(S.src) {
		return &token.Token{Type: token.EOF}
	}

	pos := S.pos

	lit := S.scanId()
	switch lit {
	case "hiya":
		return &token.Token{Type: token.TokMap.Type("hiya"),
			Lit: []byte("hiya"),
			Pos: token.Pos{Offset: pos}}
	case "hello":
		return &token.Token{Type: token.TokMap.Type("hello"),
			Lit: []byte("hello"),
			Pos: token.Pos{Offset: pos}}
	default:
		return &token.Token{Type: token.TokMap.Type("name"),
			Lit: []byte(lit),
			Pos: token.Pos{Offset: pos}}
	}
}
