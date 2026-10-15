package colorizer

import "strings"

// passThroughTerminalEscapes disables sanitizing ESC/BEL (and other C0
// controls) in record keys/values written to a TTY. Default is false so
// untrusted field data cannot drive the terminal (OSC 52 clipboard writes,
// title spoofing, screen clears, etc.). See #2255.
var passThroughTerminalEscapes bool

// SetPassThroughTerminalEscapes controls whether terminal-active bytes in
// record payloads are left intact when writing to a TTY.
func SetPassThroughTerminalEscapes(v bool) {
	passThroughTerminalEscapes = v
}

// EscapeTerminalActive replaces ESC, BEL, and other C0 controls (except tab,
// LF, and CR) with caret notation, matching the spirit of `cat -v`. Valid
// UTF-8 payload outside those bytes is left unchanged. This breaks CSI/OSC
// sequences that would otherwise be interpreted by the terminal.
func EscapeTerminalActive(s string) string {
	if !needsTerminalEscape(s) {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) + 8)
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '\t' || c == '\n' || c == '\r':
			b.WriteByte(c)
		case c == 0x7f: // DEL
			b.WriteString("^?")
		case c < 0x20: // C0 including ESC (0x1b) and BEL (0x07)
			b.WriteByte('^')
			b.WriteByte(c + 0x40) // ESC -> ^[, BEL -> ^G
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}

func needsTerminalEscape(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\t' || c == '\n' || c == '\r' {
			continue
		}
		if c < 0x20 || c == 0x7f {
			return true
		}
	}
	return false
}

// sanitizePayloadForOutput escapes terminal-active bytes when writing record
// keys/values to a TTY, unless the user opted out.
// SanitizePayloadForOutput escapes terminal-active bytes in record payloads
// written to a TTY. Exported for writers that bypass MaybeColorize* (CSV).
func SanitizePayloadForOutput(text string, outputIsStdout bool) string {
	return sanitizePayloadForOutput(text, outputIsStdout)
}

func sanitizePayloadForOutput(text string, outputIsStdout bool) string {
	if passThroughTerminalEscapes {
		return text
	}
	if !(outputIsStdout && stdoutIsATTY) {
		return text
	}
	return EscapeTerminalActive(text)
}

// SetStdoutIsATTYForTest overrides TTY detection in unit tests.
func SetStdoutIsATTYForTest(v bool) {
	stdoutIsATTY = v
}
