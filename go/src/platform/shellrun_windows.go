// ================================================================
// Wraps 'sh -c foo bar' or 'cmd /c foo bar', nominally for regression-testing.
// ================================================================

//go:build windows
// +build windows

package platform

import (
	"os"
)

// PowerShell or CMD?
//
// * Either PowerShell or CMD is fine for everything except the '...' in inline
//   put/filter statements.
//
// * PowerShell allows fractionally more, e.g. in mlr put '$flag = $x > 10' CMD
//   sees the ">" as a file-redirect (before our main() is ever entered) and
//   PowerShell lets it be operator it is intended to be.
//
// * Neither of them allows multi-line '...' inputs.
//
// * Both of the previous mean that a large number of regression tests need to
//   use mlr -f .../cases/.../0047.mlr .../cases/.../0047.input and PowerShell
//   doesn't help us eliminate this.
//
// * PowerShell has about a 2-second startup time per invocation so if we
//   use it for regression testing, we'd have to move away from one shell
//   invocation per Miller invocation and back to batching lots of Miller
//   statements into a file, doing diff/findstr to check status. That was the
//   case for Miller's original bash/file regtest framework and it was hard to
//   debug.
//
// In conclusion, we stick with CMD because it's faster, and PowerShell while
// more powerful isn't *sufficiently* more powerful to justify the
// batching-complexity to overcome its startup-latency overhead.

func GetShellRunArray(command string) []string {
	if os.Getenv("MSYSTEM") != "" {
		// Running inside MSYS2; sufficiently Unix-like already.
		return []string{"/bin/sh", "-c", command}
	} else {
		cmd := os.Getenv("COMSPEC")
		if cmd == "" {
			cmd = "C:\\Windows\\System32\\cmd.exe"
		}
		return []string{cmd, "/c", command}
	}
}
