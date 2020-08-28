// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A scanner for Go source text. Takes a []byte as source which can
// then be tokenized through repeated calls to the Scan function.
// For a sample use of a scanner, see the implementation of Tokenize.
//
package scanner

import (
	"bytes"
	"strconv"
	"unicode"
	"unicode/utf8"
)
import "github.com/goccmack/gocc/internal/frontend/token"

// A Scanner holds the scanner's internal state while processing
// a given text.  It can be allocated as part of another data
// structure but must be initialized via Init before use. For
// a sample use, see the implementation of Tokenize.
//
type Scanner struct {
	// immutable state
	src      []byte // source
	tokenMap *token.TokenMap

	// scanning state
	pos    token.Position // previous reading position (position before ch)
	offset int            // current reading offset (position after ch)
	ch     rune           // one char look-ahead

	// public state - ok to modify
	ErrorCount int // number of errors encountered
}

// Read the next Unicode char into S.ch.
// S.ch < 0 means end-of-file.
//
func (S *Scanner) next() {
	if S.offset < len(S.src) {
		S.pos.Offset = S.offset
		S.pos.Column++
		if S.ch == '\n' {
			// next character starts a new line
			S.pos.Line++
			S.pos.Column = 1
		}
		r, w := rune(S.src[S.offset]), 1
		if r == 0 {
			S.error(S.pos, "illegal character NUL")
		} else if r >= 80 {
			// not ASCII
			r, w = utf8.DecodeRune(S.src[S.offset:])
			if r == utf8.RuneError && w == 1 {
				S.error(S.pos, "illegal UTF-8 encoding")
			}
		}
		S.offset += w
		S.ch = r
	} else {
		S.pos.Offset = len(S.src)
		S.ch = -1 // eof
	}
}

// The mode parameter to the Init function is a set of flags (or 0).
// They control scanner behavior.
//
const (
	ScanComments      = 1 << iota // return comments as COMMENT tokens
	AllowIllegalChars             // do not report an error for illegal chars
	InsertSemis                   // automatically insert semicolons
)

// Init prepares the scanner S to tokenize the text src. Calls to Scan
// will use the error handler err if they encounter a syntax error and
// err is not nil. Also, for each error encountered, the Scanner field
// ErrorCount is incremented by one. The filename parameter is used as
// filename in the token.Position returned by Scan for each token. The
// mode parameter determines how comments and illegal characters are
// handled.
//
func (S *Scanner) Init(src []byte, tokenMap *token.TokenMap) {
	// Explicitly initialize all fields since a scanner may be reused.
	S.src = src
	S.tokenMap = tokenMap
	S.pos = token.Position{Offset: 0, Line: 1, Column: 0}
	S.offset = 0
	S.ErrorCount = 0
	S.next()
}

func charString(ch rune) string {
	var s string
	switch ch {
	case -1:
		return "EOF"
	case '\a':
		s = "\\a"
	case '\b':
		s = "\\b"
	case '\f':
		s = "\\f"
	case '\n':
		s = "\\n"
	case '\r':
		s = "\\r"
	case '\t':
		s = "\\t"
	case '\v':
		s = "\\v"
	case '\\':
		s = "\\\\"
	case '\'':
		s = "\\'"
	default:
		s = string(ch)
	}
	return "'" + s + "' (U+" + strconv.FormatInt(int64(ch), 16) + ")"
}

func (S *Scanner) error(pos token.Position, msg string) {
	S.ErrorCount++
}

func (S *Scanner) expect(ch rune) {
	if S.ch != ch {
		S.error(S.pos, "expected "+charString(ch)+", found "+charString(S.ch))
	}
	S.next() // always make progress
}

var prefix = []byte("line ")

func (S *Scanner) scanComment(pos token.Position) {
	// first '/' already consumed

	if S.ch == '/' {
		//-style comment
		for S.ch >= 0 {
			S.next()
			if S.ch == '\n' {
				// '\n' is not part of the comment for purposes of scanning
				// (the comment ends on the same line where it started)
				if pos.Column == 1 {
					text := S.src[pos.Offset+2 : S.pos.Offset]
					if bytes.HasPrefix(text, prefix) {
						// comment starts at beginning of line with "//line ";
						// get filename and line number, if any
						i := bytes.Index(text, []byte{':'})
						if i >= 0 {
							if line, err := strconv.Atoi(string(text[i+1:])); err == nil && line > 0 {
								// valid //line filename:line comment;
								// update scanner position
								S.pos.Line = line - 1 // -1 since the '\n' has not been consumed yet
							}
						}
					}
				}
				return
			}
		}

	} else {
		// / *-style comment * /
		S.expect('*')
		for S.ch >= 0 {
			ch := S.ch
			S.next()
			if ch == '*' && S.ch == '/' {
				S.next()
				return
			}
		}
	}

	S.error(pos, "comment not terminated")
}

func (S *Scanner) findNewline(pos token.Position) bool {
	// first '/' already consumed; assume S.ch == '/' || S.ch == '*'

	// read ahead until a newline or non-comment token is found
	newline := false
	for pos1 := pos; S.ch >= 0; {
		if S.ch == '/' {
			//-style comment always contains a newline
			newline = true
			break
		}
		S.scanComment(pos1)
		if pos1.Line < S.pos.Line {
			// / *-style comment contained a newline * /
			newline = true
			break
		}
		S.skipWhitespace() // S.insertSemi is set
		if S.ch == '\n' {
			newline = true
			break
		}
		if S.ch != '/' {
			// non-comment token
			break
		}
		pos1 = S.pos
		S.next()
		if S.ch != '/' && S.ch != '*' {
			// non-comment token
			break
		}
	}

	// reset position to where it was upon calling findNewline
	S.pos = pos
	S.offset = pos.Offset + 1
	S.next()

	return newline
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' ||
		ch >= 0x80 && unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9' || ch >= 0x80 && unicode.IsDigit(ch)
}

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch) - '0'
	case 'a' <= ch && ch <= 'f':
		return int(ch) - 'a' + 10
	case 'A' <= ch && ch <= 'F':
		return int(ch) - 'A' + 10
	}
	return 16 // larger than any legal digit val
}

func (S *Scanner) scanEscape(quote rune) {
	pos := S.pos

	var i, base, max uint32
	switch S.ch {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
		S.next()
		return
	case '0', '1', '2', '3', '4', '5', '6', '7':
		i, base, max = 3, 8, 255
	case 'x':
		S.next()
		i, base, max = 2, 16, 255
	case 'u':
		S.next()
		i, base, max = 4, 16, unicode.MaxRune
	case 'U':
		S.next()
		i, base, max = 8, 16, unicode.MaxRune
	default:
		S.next() // always make progress
		S.error(pos, "unknown escape sequence")
		return
	}

	var x uint32
	for ; i > 0; i-- {
		d := uint32(digitVal(S.ch))
		if d > base {
			S.error(S.pos, "illegal character in escape sequence")
			return
		}
		x = x*base + d
		S.next()
	}
	if x > max || 0xd800 <= x && x < 0xe000 {
		S.error(pos, "escape sequence is invalid Unicode code point")
	}
}

func (S *Scanner) scanChar(pos token.Position) {
	// '\'' already consumed

	n := 0
	for S.ch != '\'' {
		ch := S.ch
		n++
		S.next()
		if ch == '\n' || ch < 0 {
			S.error(pos, "character literal not terminated")
			n = 1
			break
		}
		if ch == '\\' {
			S.scanEscape('\'')
		}
	}

	S.next()

	if n != 1 {
		S.error(pos, "illegal character literal")
	}
}

func (S *Scanner) scanIdentifier(pos token.Position) token.Type {
	ch0 := S.ch
	for isLetter(S.ch) || isDigit(S.ch) || S.ch == '!' {
		S.next()
	}
	switch {
	case string(S.src[pos.Offset:S.pos.Offset]) == "import":
		return S.tokenMap.Type("import")
	case ch0 == '!':
		return S.tokenMap.Type("ignoredTokId")
	case ch0 == '_':
		return S.tokenMap.Type("regDefId")
	case unicode.IsUpper(ch0):
		return S.tokenMap.Type("prodId")
	default:
		return S.tokenMap.Type("tokId")
	}
}

var (
	rw_CharLit   = []byte("char_lit")
	rw_StringLit = []byte("string_lit")
)

func (S *Scanner) isReservedWord(lit []byte) bool {
	return bytes.Equal(lit, rw_CharLit) || bytes.Equal(lit, rw_StringLit)
}

// func (S *Scanner) scanInteger() {
// 	if !isDigit(S.ch) { S.error(S.pos, "integer index expected") }
//
//
// 	for S.next(); isDigit(S.ch); S.next() {}
// }

func (S *Scanner) scanNumber() token.Type {
	if S.ch == '-' {
		S.next()
	}
	for isDigit(S.ch) {
		S.next()
	}
	return S.tokenMap.Type("int_lit")
}

func (S *Scanner) scanSDTLit(pos token.Position) {
	// '<' already consumed
	S.next() // consume second <
	for cmp := false; !cmp; {
		if S.ch < 0 {
			S.error(pos, "SDT not terminated")
			break
		}
		if S.ch == '>' {
			S.next()
			if S.ch == '>' {
				break
			}
		}
		S.next()
	}
	S.next()
}

func (S *Scanner) scanString(pos token.Position) {
	// '"' already consumed

	for S.ch != '"' {
		ch := S.ch
		S.next()
		if ch == '\n' || ch < 0 {
			S.error(pos, "string not terminated")
			break
		}
		if ch == '\\' {
			S.scanEscape('"')
		}
	}

	S.next()
}

func (S *Scanner) scanRawString(pos token.Position) {
	// '\140' already consumed

	for S.ch != '\140' {
		ch := S.ch
		S.next()
		if ch < 0 {
			S.error(pos, "string not terminated")
			break
		}
	}

	S.next()
}

func (S *Scanner) skipWhitespace() {
	for S.ch == ' ' || S.ch == '\t' || S.ch == '\n' || S.ch == '\r' {
		S.next()
	}
}

// Helper functions for scanning multi-byte tokens such as >> += >>= .
// Different routines recognize different length tok_i based on matches
// of ch_i. If a token ends in '=', the result is tok1 or tok3
// respectively. Otherwise, the result is tok0 if there was no other
// matching character, or tok2 if the matching character was ch2.

func (S *Scanner) switch2(tok0, tok1 token.Type) token.Type {
	if S.ch == '=' {
		S.next()
		return tok1
	}
	return tok0
}

func (S *Scanner) switch3(tok0, tok1 token.Type, ch2 rune, tok2 token.Type) token.Type {
	if S.ch == '=' {
		S.next()
		return tok1
	}
	if S.ch == ch2 {
		S.next()
		return tok2
	}
	return tok0
}

func (S *Scanner) switch4(tok0, tok1 token.Type, ch2 rune, tok2, tok3 token.Type) token.Type {
	if S.ch == '=' {
		S.next()
		return tok1
	}
	if S.ch == ch2 {
		S.next()
		if S.ch == '=' {
			S.next()
			return tok3
		}
		return tok2
	}
	return tok0
}

var semicolon = []byte{';'}

// Scan scans the next token and returns the token position pos,
// the token tok, and the literal text lit corresponding to the
// token. The source end is indicated by token.EOF.
//
// For more tolerant parsing, Scan will return a valid token if
// possible even if a syntax error was encountered. Thus, even
// if the resulting token sequence contains no illegal tokens,
// a client may not assume that no error occurred. Instead it
// must check the scanner's ErrorCount or the number of calls
// of the error handler, if there was one installed.
//
func (S *Scanner) Scan() (*token.Token, token.Position) {
scanAgain:
	S.skipWhitespace()

	// current token start
	pos, tok := S.pos, token.ILLEGAL

	// determine token value
	switch ch := S.ch; {
	case ch == '!' || isLetter(ch):
		tok = S.scanIdentifier(pos)
	default:
		S.next() // always make progress
		switch ch {
		case -1:
			tok = S.tokenMap.Type("$")
		case '"':
			tok = S.tokenMap.Type("string_lit")
			S.scanString(pos)
		case '\'':
			tok = S.tokenMap.Type("char_lit")
			S.scanChar(pos)
		case '\140':
			tok = S.tokenMap.Type("string_lit")
			S.scanRawString(pos)
		case '-':
			tok = S.tokenMap.Type("-")
		case '{':
			tok = S.tokenMap.Type("{")
		case '}':
			tok = S.tokenMap.Type("}")
		case ':':
			tok = S.tokenMap.Type(":")
		case ';':
			tok = S.tokenMap.Type(";")
		case ',':
			tok = S.tokenMap.Type(",")
		case '[':
			tok = S.tokenMap.Type("[")
		case ']':
			tok = S.tokenMap.Type("]")
		case '(':
			tok = S.tokenMap.Type("(")
		case ')':
			tok = S.tokenMap.Type(")")
		case '|':
			tok = S.tokenMap.Type("|")
		case '/':
			if S.ch == '/' || S.ch == '*' {
				// comment
				S.scanComment(pos)
				goto scanAgain
			} else {
				tok = S.tokenMap.Type("/")
			}
		case '<':
			switch S.ch {
			case '<':
				tok = S.tokenMap.Type("g_sdt_lit")
				S.scanSDTLit(pos)
			case '=':
				tok = S.tokenMap.Type("<=")
				S.next()
			default:
				tok = S.tokenMap.Type("<")
			}
		case '.':
			tok = S.tokenMap.Type(".")
		default:
			S.error(pos, "illegal character "+charString(ch))
		}
	}
	return token.NewToken(tok, S.src[pos.Offset:S.pos.Offset]), pos
}

// An implementation of an ErrorHandler may be provided to the Scanner.
// If a syntax error is encountered and a handler was installed, Error
// is called with a position and an error message. The position points
// to the beginning of the offending token.
//
type ErrorHandler interface {
	Error(pos token.Position, msg string)
}

// Within ErrorVector, an error is represented by an Error node. The
// position Pos, if valid, points to the beginning of the offending
// token, and the error condition is described by Msg.
//
type Error struct {
	Pos token.Position
	Msg string
}

func (e *Error) String() string {
	if e.Pos.IsValid() {
		// don't print "<unknown position>"
		return e.Pos.String() + ": " + e.Msg
	}
	return e.Msg
}
