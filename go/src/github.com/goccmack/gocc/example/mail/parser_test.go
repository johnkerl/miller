package mail

import (
	"testing"

	"github.com/goccmack/gocc/example/mail/lexer"
	"github.com/goccmack/gocc/example/mail/token"
)

var testData1 = map[string]bool{
	"mymail@google.com":          true,
	"@google.com":                false,
	`"quoted string"@mymail.com`: true,
	`"unclosed quote@mymail.com`: false,
}

func Test1(t *testing.T) {
	for input, ok := range testData1 {
		l := lexer.NewLexer([]byte(input))
		tok := l.Scan()
		switch {
		case tok.Type == token.INVALID:
			if ok {
				t.Errorf("%s", input)
			}
		case tok.Type == token.TokMap.Type("addrspec"):
			if !ok {
				t.Errorf("%s", input)
			}
		default:
			t.Fatalf("This must not happen")
		}
	}
}

var checkData2 = []string{
	"addr1@gmail.com",
	"addr2@gmail.com",
	"addr3@gmail.com",
}

var testData2 = `
	addr1@gmail.com
	addr2@gmail.com
	addr3@gmail.com
`

func Test2(t *testing.T) {
	l := lexer.NewLexer([]byte(testData2))
	num := 0
	for tok := l.Scan(); tok.Type == token.TokMap.Type("addrspec"); tok = l.Scan() {
		if string(tok.Lit) != checkData2[num] {
			t.Errorf("%s != %s", string(tok.Lit), checkData2[num])
		}
		num++
	}
	if num != len(checkData2) {
		t.Fatalf("%d addresses parsed", num)
	}
}
