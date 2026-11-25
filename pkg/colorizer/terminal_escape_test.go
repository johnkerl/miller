package colorizer

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestEscapeTerminalActive_OSC52(t *testing.T) {
	payload := base64.StdEncoding.EncodeToString([]byte("TUI-INJECTION-POC"))
	raw := "\x1b]52;c;" + payload + "\x07"
	got := EscapeTerminalActive(raw)
	if strings.Contains(got, "\x1b") || strings.Contains(got, "\x07") {
		t.Fatalf("escape left terminal-active bytes: %q", got)
	}
	if !strings.HasPrefix(got, "^[]52;c;") {
		t.Fatalf("expected caret-escaped ESC prefix, got %q", got)
	}
	if !strings.HasSuffix(got, "^G") {
		t.Fatalf("expected caret-escaped BEL suffix, got %q", got)
	}
}

func TestEscapeTerminalActive_PreservesPrintableAndCommonWS(t *testing.T) {
	in := "hello\tworld\nline\r"
	if got := EscapeTerminalActive(in); got != in {
		t.Fatalf("printable/common ws changed: got %q want %q", got, in)
	}
}

func TestSanitizePayloadForOutput_TTYOnly(t *testing.T) {
	prevTTY := stdoutIsATTY
	prevPass := passThroughTerminalEscapes
	t.Cleanup(func() {
		stdoutIsATTY = prevTTY
		passThroughTerminalEscapes = prevPass
	})

	raw := "x\x1b]0;title\x07y"
	stdoutIsATTY = true
	passThroughTerminalEscapes = false
	if got := sanitizePayloadForOutput(raw, true); strings.Contains(got, "\x1b") {
		t.Fatalf("TTY stdout should sanitize, got %q", got)
	}
	if got := sanitizePayloadForOutput(raw, false); got != raw {
		t.Fatalf("non-stdout should not sanitize, got %q", got)
	}

	stdoutIsATTY = false
	if got := sanitizePayloadForOutput(raw, true); got != raw {
		t.Fatalf("non-TTY should not sanitize, got %q", got)
	}

	stdoutIsATTY = true
	passThroughTerminalEscapes = true
	if got := sanitizePayloadForOutput(raw, true); got != raw {
		t.Fatalf("pass-through should leave bytes intact, got %q", got)
	}
}

func TestMaybeColorizeValue_SanitizesOnTTY(t *testing.T) {
	prevTTY := stdoutIsATTY
	prevPass := passThroughTerminalEscapes
	prevColor := colorization
	t.Cleanup(func() {
		stdoutIsATTY = prevTTY
		passThroughTerminalEscapes = prevPass
		colorization = prevColor
	})

	stdoutIsATTY = true
	passThroughTerminalEscapes = false
	colorization = ColorizeOutputNever
	raw := "a\x1bb"
	got := MaybeColorizeValue(raw, true)
	if strings.Contains(got, "\x1b") {
		t.Fatalf("MaybeColorizeValue left ESC on TTY: %q", got)
	}
	if got != "a^[b" {
		t.Fatalf("got %q want %q", got, "a^[b")
	}
}
