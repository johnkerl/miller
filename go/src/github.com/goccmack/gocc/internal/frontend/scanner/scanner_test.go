package scanner

import (
	"fmt"
	"testing"

	"github.com/goccmack/gocc/internal/frontend/token"
)

type testRecord struct {
	src    string
	typ    token.Type
	tokLit string
}

var testData = []testRecord{
	{"tokId", token.FRONTENDTokens.Type("tokId"), "tokId"},
	{"!whitespace", token.FRONTENDTokens.Type("ignoredTokId"), "!whitespace"},
	{":", token.FRONTENDTokens.Type(":"), ":"},
	{";", token.FRONTENDTokens.Type(";"), ";"},
	{"_regDefId", token.FRONTENDTokens.Type("regDefId"), "_regDefId"},
	{"|", token.FRONTENDTokens.Type("|"), "|"},
	{`'\u0011'`, token.FRONTENDTokens.Type("char_lit"), `'\u0011'`},
	{"-", token.FRONTENDTokens.Type("-"), "-"},
	{"(", token.FRONTENDTokens.Type("("), "("},
	{")", token.FRONTENDTokens.Type(")"), ")"},
	{"[", token.FRONTENDTokens.Type("["), "["},
	{"]", token.FRONTENDTokens.Type("]"), "]"},
	{"{", token.FRONTENDTokens.Type("{"), "{"},
	{"}", token.FRONTENDTokens.Type("}"), "}"},
	{"<< sdt lit >>", token.FRONTENDTokens.Type("g_sdt_lit"), "<< sdt lit >>"},
	{"ProdId", token.FRONTENDTokens.Type("prodId"), "ProdId"},
	{`"string lit"`, token.FRONTENDTokens.Type("string_lit"), `"string lit"`},
}

func Test1(tst *testing.T) {
	s := &Scanner{}
	for _, t := range testData {
		s.Init([]byte(t.src), token.FRONTENDTokens)
		tok, _ := s.Scan()
		if tok.Type != t.typ {
			tst.Error(fmt.Sprintf("src: %s, type: %d -- got type: %d\n", t.src, t.typ, tok.Type))
		}
		if string(tok.Lit) != t.tokLit {
			tst.Error(fmt.Sprintf("src: %s, expected lit: %s, got: %s\n", t.src, t.tokLit, string(tok.Lit)))
		}
	}
}

func Test2(t *testing.T) {
	s := &Scanner{}
	lit := "The SDT Lit"
	s.Init([]byte(fmt.Sprintf("<< %s >>", lit)), token.FRONTENDTokens)
	tok, _ := s.Scan()
	if tok.Type != token.FRONTENDTokens.Type("g_sdt_lit") {
		t.Error(fmt.Sprintf("Expected tok type: g_sdt_lit, got: %s", token.FRONTENDTokens.TokenString(tok.Type)))
	}
	if tok.SDTVal() != lit {
		t.Error(fmt.Sprintf("Expected SDTVal: %s, got: %s\n", lit, tok.SDTVal()))
	}
}
