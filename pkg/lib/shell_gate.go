// Package-level gate for shell-outs, controlled by the --no-shell main flag
// or a truthy MLR_NO_SHELL environment variable.
//
// When disabled, every place Miller can execute an external command in the
// data path -- the DSL system() and exec() functions, piped output redirects
// (e.g. `tee | "command"`), and --prepipe/--prepipex -- fails cleanly instead
// of running the command. This exists so that automation driving Miller (in
// particular the `mlr mcp` server, which sets MLR_NO_SHELL on the subprocesses
// it spawns) can run agent-constructed command lines without also granting
// arbitrary command execution.
//
// The gate is one-way by design: it can be disabled during startup, but never
// re-enabled, so an argv-injected flag cannot override an environment-level
// opt-out.

package lib

var shellOutEnabled = true

// DisableShellOut turns off all shell-out capability for the remainder of the
// process lifetime.
func DisableShellOut() {
	shellOutEnabled = false
}

// ShellOutEnabled reports whether shell-outs are permitted.
func ShellOutEnabled() bool {
	return shellOutEnabled
}

// IsTruthyEnvValue returns true for non-empty strings commonly used as
// boolean env-var truthy values: "1", "true", "yes" (case-insensitive).
func IsTruthyEnvValue(v string) bool {
	switch v {
	case "1", "true", "True", "TRUE", "yes", "Yes", "YES":
		return true
	}
	return false
}
