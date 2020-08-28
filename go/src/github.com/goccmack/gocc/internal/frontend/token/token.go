package token

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type Token struct {
	Type Type
	Lit  []byte
}

func NewToken(typ Type, lit []byte) *Token {
	return &Token{typ, lit}
}

func (this *Token) Equals(that *Token) bool {
	if this == nil || that == nil {
		return this == that
	}

	if this.Type != that.Type {
		return false
	}

	return bytes.Equal(this.Lit, that.Lit)
}

func (this *Token) String() string {
	str := ""
	if this.Type == EOF {
		str += "\"$\""
	} else {
		str += "\"" + string(this.Lit) + "\""
	}
	str += "(" + strconv.Itoa(int(this.Type)) + ")"
	return str
}

type Type int

const (
	ILLEGAL Type = iota - 1
	EOF
)

func (T Type) String() string {
	return strconv.Itoa(int(T))
}

// Position describes an arbitrary source position
// including the file, line, and column location.
// A Position is valid if the line number is > 0.
//
type Position struct {
	Offset int // offset, starting at 0
	Line   int // line number, starting at 1
	Column int // column number, starting at 1 (character count)
}

// IsValid returns true if the position is valid.
func (pos *Position) IsValid() bool { return pos.Line > 0 }

// String returns a string in one of several forms:
//
//	file:line:column    valid position with file name
//	line:column         valid position without file name
//	file                invalid position with file name
//	-                   invalid position without file name
//
func (pos Position) String() string {
	s := ""
	if pos.IsValid() {
		s += fmt.Sprintf("%d:%d", pos.Line, pos.Column)
	}
	if s == "" {
		s = "-"
	}
	return s
}

func (T *Token) IntValue() (int64, error) {
	return strconv.ParseInt(string(T.Lit), 10, 64)
}

func (T *Token) UintValue() (uint64, error) {
	return strconv.ParseUint(string(T.Lit), 10, 64)
}

func (T *Token) SDTVal() string {
	sdt := string(T.Lit)
	rex, err := regexp.Compile("\\$[0-9]+")
	if err != nil {
		panic(err)
	}
	idx := rex.FindAllStringIndex(sdt, -1)
	res := ""
	if len(idx) <= 0 {
		res = sdt
	} else {
		for i, loc := range idx {
			if loc[0] > 0 {
				if i > 0 {
					res += sdt[idx[i-1][1]:loc[0]]
				} else {
					res += sdt[0:loc[0]]
				}
			}
			res += "X["
			res += sdt[loc[0]+1 : loc[1]]
			res += "]"
		}
		if idx[len(idx)-1][1] < len(sdt) {
			res += sdt[idx[len(idx)-1][1]:]
		}
	}
	return strings.TrimSpace(res[2 : len(res)-2])
}

// Tokenmap

type TokenMap struct {
	tokenMap  []string
	stringMap map[string]Type
}

func NewMap() *TokenMap {
	tm := &TokenMap{make([]string, 0, 10), make(map[string]Type)}
	tm.AddToken("$")
	// tm.AddToken("ε")
	return tm
}

func (this *TokenMap) AddToken(str string) {
	if _, exists := this.stringMap[str]; exists {
		return
	}
	this.stringMap[str] = Type(len(this.tokenMap))
	this.tokenMap = append(this.tokenMap, str)
}

func NewMapFromFile(file string) (*TokenMap, error) {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return NewMapFromString(string(src)), nil
}

func NewMapFromStrings(input []string) *TokenMap {
	tm := NewMap()
	for _, s := range input {
		tm.AddToken(s)
	}
	return tm
}

func NewMapFromString(input string) *TokenMap {
	tokens := strings.Fields(input)
	return NewMapFromStrings(tokens)
}

func (this *TokenMap) Len() int {
	return len(this.tokenMap)
}

func (this *TokenMap) Type(key string) Type {
	tok, ok := this.stringMap[key]
	if !ok {
		return ILLEGAL
	}
	return tok
}

func (this *TokenMap) TokenString(typ Type) string {
	tok := int(typ)
	if tok < 0 || tok >= len(this.tokenMap) {
		return "illegal " + strconv.Itoa(tok)
	}
	return this.tokenMap[tok]
}

func (this *TokenMap) String() string {
	res := ""
	for str, tok := range this.stringMap {
		res += str + " : " + strconv.Itoa(int(tok)) + "\n"
	}
	return res
}

func (this *TokenMap) Strings() []string {
	return this.tokenMap[1:]
}

func (this *TokenMap) Equals(that *TokenMap) bool {
	if this == nil || that == nil {
		return false
	}

	if len(this.stringMap) != len(that.stringMap) ||
		len(this.tokenMap) != len(that.tokenMap) {
		return false
	}

	for str, tok := range this.stringMap {
		if tok1, ok := that.stringMap[str]; !ok || tok1 != tok {
			return false
		}
	}

	return true
}

func (this *TokenMap) Tokens() []*Token {
	res := make([]*Token, 0, len(this.stringMap))
	for typ, str := range this.tokenMap {
		res = append(res, &Token{Type(typ), []byte(str)})
	}
	return res
}

func (this *TokenMap) WriteFile(file string) error {
	out := ""
	for i := 1; i < len(this.tokenMap); i++ {
		out += this.TokenString(Type(i)) + "\n"
	}
	return ioutil.WriteFile(file, []byte(out), 0644)
}
