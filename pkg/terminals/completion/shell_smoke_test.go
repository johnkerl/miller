//go:build !windows

// Shell-level smoke tests for the generated bash completion script.
//
// These guard the shell glue itself -- in particular the bash 3.2 (macOS)
// array-slicing defect that caused every candidate to be mashed into a single
// COMPREPLY entry and dumped onto the command line. The Go-level engine tests
// cannot catch that class of bug, since it lives entirely in the emitted shell.
//
// The test sources the real generated script under /bin/bash, drives
// _mlr_complete with synthetic COMP_WORDS/COMP_CWORD, and inspects COMPREPLY.
// It builds the mlr binary because the shim shells back out to it.

package completion

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
)

var (
	builtMlrOnce sync.Once
	builtMlrPath string
	builtMlrErr  error
)

// buildMlr builds the mlr binary once per test process and returns its path.
func buildMlr(t *testing.T) string {
	t.Helper()
	builtMlrOnce.Do(func() {
		dir, err := os.MkdirTemp("", "mlr-completion-smoke")
		if err != nil {
			builtMlrErr = err
			return
		}
		bin := filepath.Join(dir, "mlr")
		cmd := exec.Command("go", "build", "-o", bin, "github.com/johnkerl/miller/v6/cmd/mlr")
		if out, err := cmd.CombinedOutput(); err != nil {
			builtMlrErr = err
			t.Logf("go build output:\n%s", out)
			return
		}
		builtMlrPath = bin
	})
	if builtMlrErr != nil {
		t.Fatalf("building mlr: %v", builtMlrErr)
	}
	return builtMlrPath
}

func shellSingleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// runBashCompletion sources the script, drives _mlr_complete for the given
// command-line words (words[0] is replaced with the real mlr binary), and
// returns the COMPREPLY entries. dir, if non-empty, is the working directory
// (so file-completion behavior is deterministic).
func runBashCompletion(t *testing.T, scriptPath, mlrBin, dir string, words []string, cword int) []string {
	t.Helper()

	quoted := make([]string, len(words))
	for i, w := range words {
		if i == 0 {
			quoted[i] = shellSingleQuote(mlrBin)
		} else {
			quoted[i] = shellSingleQuote(w)
		}
	}

	driver := strings.Join([]string{
		"source " + shellSingleQuote(scriptPath),
		"COMP_WORDS=(" + strings.Join(quoted, " ") + ")",
		"COMP_CWORD=" + strconv.Itoa(cword),
		"COMPREPLY=()",
		"_mlr_complete",
		`for c in "${COMPREPLY[@]}"; do printf 'ENTRY:%s\n' "$c"; done`,
	}, "\n")

	cmd := exec.Command("/bin/bash", "-c", driver)
	if dir != "" {
		cmd.Dir = dir
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("bash driver failed: %v\n%s", err, out)
	}

	var entries []string
	for _, line := range strings.Split(string(out), "\n") {
		if rest, ok := strings.CutPrefix(line, "ENTRY:"); ok {
			entries = append(entries, rest)
		}
	}
	return entries
}

func TestBashShim(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping shell smoke test in -short mode")
	}
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go toolchain not available to build mlr")
	}
	bashPath := "/bin/bash"
	if _, err := os.Stat(bashPath); err != nil {
		t.Skip("no /bin/bash available")
	}

	mlrBin := buildMlr(t)

	scriptDir := t.TempDir()
	scriptPath := filepath.Join(scriptDir, "mlr_completion.bash")
	if err := os.WriteFile(scriptPath, []byte(bashScript), 0o644); err != nil {
		t.Fatal(err)
	}

	containsAll := func(t *testing.T, got, want []string) {
		t.Helper()
		set := make(map[string]bool, len(got))
		for _, g := range got {
			set[g] = true
		}
		for _, w := range want {
			if !set[w] {
				t.Errorf("missing %q in COMPREPLY: %v", w, got)
			}
		}
	}

	// The core regression: multiple candidates must come back as SEPARATE
	// COMPREPLY entries, not one space-joined blob (the bash 3.2 array-slicing
	// bug). This is what caused all candidates to be inserted onto the line.
	t.Run("candidates are split, not joined", func(t *testing.T) {
		got := runBashCompletion(t, scriptPath, mlrBin, "", []string{"mlr", "--m"}, 1)
		if len(got) < 2 {
			t.Fatalf("expected multiple candidates for '--m', got %d: %v", len(got), got)
		}
		for _, c := range got {
			if strings.ContainsAny(c, " \t") {
				t.Errorf("candidate contains whitespace (slicing regression?): %q", c)
			}
			if !strings.HasPrefix(c, "--m") {
				t.Errorf("unexpected candidate for '--m': %q", c)
			}
		}
	})

	t.Run("verb names before first verb", func(t *testing.T) {
		got := runBashCompletion(t, scriptPath, mlrBin, "", []string{"mlr", ""}, 1)
		containsAll(t, got, []string{"cat", "sort", "put"})
	})

	t.Run("terminal subcommands and flags", func(t *testing.T) {
		got := runBashCompletion(t, scriptPath, mlrBin, "", []string{"mlr", ""}, 1)
		containsAll(t, got, []string{"help", "version", "repl"})

		gotFlags := runBashCompletion(t, scriptPath, mlrBin, "", []string{"mlr", "-"}, 1)
		containsAll(t, gotFlags, []string{"-h", "--help", "--version"})
	})

	t.Run("help topics and topic arguments", func(t *testing.T) {
		topics := runBashCompletion(t, scriptPath, mlrBin, "", []string{"mlr", "help", ""}, 2)
		containsAll(t, topics, []string{"flags", "verb", "function"})

		verbs := runBashCompletion(t, scriptPath, mlrBin, "", []string{"mlr", "help", "verb", ""}, 3)
		containsAll(t, verbs, []string{"cat", "sort"})
	})

	t.Run("verb flags inside a verb", func(t *testing.T) {
		got := runBashCompletion(t, scriptPath, mlrBin, "", []string{"mlr", "cat", "-"}, 2)
		containsAll(t, got, []string{"-n", "--filename"})
	})

	t.Run("main flags on bare dash include all flags", func(t *testing.T) {
		got := runBashCompletion(t, scriptPath, mlrBin, "", []string{"mlr", "-"}, 1)
		// Includes ordinary flags as well as the format-conversion
		// keystroke-savers (--c2j, --m2j, ...).
		containsAll(t, got, []string{"--icsv", "--ojson", "--m2j", "--c2p"})
	})

	t.Run("conversion matrix narrows by prefix", func(t *testing.T) {
		got := runBashCompletion(t, scriptPath, mlrBin, "", []string{"mlr", "--m2"}, 1)
		containsAll(t, got, []string{"--m2j", "--m2p"})
	})

	// File-completion directives (the bash 3.2 'no compopt' path uses
	// compgen -f). --from takes a filename argument, so it defers to file
	// completion. Run in a temp dir with known files for determinism.
	t.Run("file completion after filename flag", func(t *testing.T) {
		dataDir := t.TempDir()
		for _, name := range []string{"alpha.csv", "beta.csv"} {
			if err := os.WriteFile(filepath.Join(dataDir, name), nil, 0o644); err != nil {
				t.Fatal(err)
			}
		}
		got := runBashCompletion(t, scriptPath, mlrBin, dataDir, []string{"mlr", "--from", ""}, 2)
		containsAll(t, got, []string{"alpha.csv", "beta.csv"})
	})

	// Enum-value completion: an arg-taking flag whose values are a known set
	// (here file formats) offers those values rather than filenames.
	t.Run("enum value completion for format flag", func(t *testing.T) {
		got := runBashCompletion(t, scriptPath, mlrBin, "", []string{"mlr", "-i", ""}, 2)
		containsAll(t, got, []string{"csv", "json", "tsv"})
	})

	t.Run("then plus files inside verb", func(t *testing.T) {
		dataDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dataDir, "gamma.csv"), nil, 0o644); err != nil {
			t.Fatal(err)
		}
		got := runBashCompletion(t, scriptPath, mlrBin, dataDir, []string{"mlr", "cat", ""}, 2)
		containsAll(t, got, []string{"then", "gamma.csv"})
	})
}

func TestZshShim(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping shell smoke test in -short mode")
	}
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go toolchain not available to build mlr")
	}
	zshPath, err := exec.LookPath("zsh")
	if err != nil {
		t.Skip("no zsh available")
	}

	mlrBin := buildMlr(t)

	scriptPath := filepath.Join(t.TempDir(), "_mlr")
	if err := os.WriteFile(scriptPath, []byte(zshScript), 0o644); err != nil {
		t.Fatal(err)
	}

	// The script must source cleanly even when compinit has never run (zsh -f
	// loads no startup files). This guards the 'command not found: compdef'
	// regression: the script self-initializes the completion system.
	t.Run("sources cleanly and registers without precompinit", func(t *testing.T) {
		driver := strings.Join([]string{
			"source " + shellSingleQuote(scriptPath),
			"(( $+functions[_mlr] )) || { print -r FAIL_NO_FUNC; exit 3 }",
			"(( $+_comps[mlr] )) || { print -r FAIL_NO_COMPDEF; exit 4 }",
			"print -r OK",
		}, "\n")
		cmd := exec.Command(zshPath, "-f", "-c", driver)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("zsh sourcing failed: %v\n%s", err, out)
		}
		if !strings.Contains(string(out), "OK") {
			t.Fatalf("expected OK, got:\n%s", out)
		}
	})

	// zsh-side parsing: the ${(@f)} newline-split plus the [2,-1] slice must
	// yield each candidate as a separate array element (the zsh analogue of the
	// bash array-slicing regression).
	t.Run("splits candidates into separate elements", func(t *testing.T) {
		driver := fmt.Sprintf(`response=("${(@f)$(%s completion complete 1 mlr --m 2>/dev/null)}")
candidates=(${response[2,-1]})
for c in $candidates; do print -r -- "ENTRY:$c"; done`, shellSingleQuote(mlrBin))

		cmd := exec.Command(zshPath, "-f", "-c", driver)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("zsh driver failed: %v\n%s", err, out)
		}
		var entries []string
		for _, line := range strings.Split(string(out), "\n") {
			if rest, ok := strings.CutPrefix(line, "ENTRY:"); ok {
				entries = append(entries, rest)
			}
		}
		if len(entries) < 2 {
			t.Fatalf("expected multiple candidates for '--m', got %d: %v", len(entries), entries)
		}
		for _, c := range entries {
			if strings.ContainsAny(c, " \t") {
				t.Errorf("candidate contains whitespace (slicing regression?): %q", c)
			}
			if !strings.HasPrefix(c, "--m") {
				t.Errorf("unexpected candidate for '--m': %q", c)
			}
		}
	})
}
