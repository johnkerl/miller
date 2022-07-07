
//go:build windows
// +build windows

package platform

import (
	"syscall"

	sequences "github.com/nine-lives-later/go-windows-terminal-sequences"
)

func EnableAnsiEscapeSequences() {
	sequences.EnableVirtualTerminalProcessing(syscall.Stdout, true)
	sequences.EnableVirtualTerminalProcessing(syscall.Stderr, true)
}