package climain

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/johnkerl/miller/v6/pkg/cli"
)

// writeTempMlrrc writes contents to a file in a temp directory and points the
// MLRRC environment variable at it, so loadMlrrc reads it and only it.
func writeTempMlrrc(t *testing.T, contents string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "mlrrc")
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("MLRRC", path)
}

const testMlrrcContents = `# A global setting, applied always:
icsv

[j]
ojson # A comment after a setting
jvstack

[p]
opprint

[j]
no-auto-flatten
`

func TestMlrrcGlobalOnlyBackCompat(t *testing.T) {
	// No sections at all: everything applies -- this is the historical
	// .mlrrc format.
	writeTempMlrrc(t, "icsv\nojson\n")
	options := cli.DefaultOptions()
	if err := loadMlrrcFiles(options, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if options.ReaderOptions.InputFileFormat != "csv" {
		t.Errorf("input format: got %s, want csv", options.ReaderOptions.InputFileFormat)
	}
	if options.WriterOptions.OutputFileFormat != "json" {
		t.Errorf("output format: got %s, want json", options.WriterOptions.OutputFileFormat)
	}
}

func TestMlrrcSectionsIgnoredWithoutProfile(t *testing.T) {
	writeTempMlrrc(t, testMlrrcContents)
	options := cli.DefaultOptions()
	if err := loadMlrrcFiles(options, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if options.ReaderOptions.InputFileFormat != "csv" {
		t.Errorf("input format: got %s, want csv", options.ReaderOptions.InputFileFormat)
	}
	// The [j] and [p] sections must not be applied.
	if options.WriterOptions.OutputFileFormat != "dkvp" {
		t.Errorf("output format: got %s, want dkvp", options.WriterOptions.OutputFileFormat)
	}
}

func TestMlrrcProfileSelection(t *testing.T) {
	writeTempMlrrc(t, testMlrrcContents)
	options := cli.DefaultOptions()
	if err := loadMlrrcFiles(options, "j"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Global setting applies first ...
	if options.ReaderOptions.InputFileFormat != "csv" {
		t.Errorf("input format: got %s, want csv", options.ReaderOptions.InputFileFormat)
	}
	// ... then the [j] settings.
	if options.WriterOptions.OutputFileFormat != "json" {
		t.Errorf("output format: got %s, want json", options.WriterOptions.OutputFileFormat)
	}
}

func TestMlrrcRepeatedSectionsAccumulate(t *testing.T) {
	writeTempMlrrc(t, testMlrrcContents)
	options := cli.DefaultOptions()
	if err := loadMlrrcFiles(options, "j"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// no-auto-flatten is in the second [j] block.
	if options.WriterOptions.AutoFlatten {
		t.Errorf("auto-flatten: got true, want false via second [j] block")
	}
}

func TestMlrrcOtherProfileNotApplied(t *testing.T) {
	writeTempMlrrc(t, testMlrrcContents)
	options := cli.DefaultOptions()
	if err := loadMlrrcFiles(options, "p"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if options.WriterOptions.OutputFileFormat != "pprint" {
		t.Errorf("output format: got %s, want pprint", options.WriterOptions.OutputFileFormat)
	}
}

func TestMlrrcMissingProfileIsError(t *testing.T) {
	writeTempMlrrc(t, testMlrrcContents)
	options := cli.DefaultOptions()
	err := loadMlrrcFiles(options, "nonesuch")
	if err == nil {
		t.Fatal("expected error for missing profile, got nil")
	}
	if !strings.Contains(err.Error(), "nonesuch") {
		t.Errorf("error should name the missing profile: %v", err)
	}
}

func TestMlrrcProfileWithMlrrcNoneIsError(t *testing.T) {
	t.Setenv("MLRRC", "__none__")
	options := cli.DefaultOptions()
	err := loadMlrrcFiles(options, "j")
	if err == nil {
		t.Fatal("expected error for profile with MLRRC=__none__, got nil")
	}
}

func TestMlrrcProfileWithNoMlrrcFileIsError(t *testing.T) {
	// Point MLRRC at a nonexistent file: an unopenable file is silently
	// skipped (the normal no-.mlrrc case), and HOME is remapped to an empty
	// temp directory, so no .mlrrc file is processed at all. XDG_CONFIG_HOME
	// is cleared so a developer's real config dir doesn't leak in.
	t.Setenv("MLRRC", filepath.Join(t.TempDir(), "nonexistent"))
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", "")
	options := cli.DefaultOptions()
	err := loadMlrrcFiles(options, "j")
	if err == nil {
		t.Fatal("expected error for profile with no .mlrrc file, got nil")
	}
	if !strings.Contains(err.Error(), "no .mlrrc file was found") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestMlrrcWhitespaceAndCommentsAroundSectionHeaders(t *testing.T) {
	writeTempMlrrc(t, "icsv\n\n  [ j ]  # a comment\nojson\n")
	options := cli.DefaultOptions()
	if err := loadMlrrcFiles(options, "j"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if options.WriterOptions.OutputFileFormat != "json" {
		t.Errorf("output format: got %s, want json", options.WriterOptions.OutputFileFormat)
	}
}

func TestMlrrcMalformedSectionHeadersAreErrors(t *testing.T) {
	for _, header := range []string{"[j", "[]", "[ ]", "[a[b]"} {
		writeTempMlrrc(t, header+"\n")
		options := cli.DefaultOptions()
		err := loadMlrrcFiles(options, "")
		if err == nil {
			t.Errorf("expected parse error for header %q, got nil", header)
		} else if !strings.Contains(err.Error(), "parse error") {
			t.Errorf("expected parse error for header %q, got: %v", header, err)
		}
	}
}

func TestMlrrcUnusedProfileLinesNotValidated(t *testing.T) {
	// A typo inside a section which isn't selected must not break other
	// invocations of mlr.
	writeTempMlrrc(t, "icsv\n[j]\nthis-is-not-a-flag\n")
	options := cli.DefaultOptions()
	if err := loadMlrrcFiles(options, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// But it is a parse error when the section is selected.
	options = cli.DefaultOptions()
	if err := loadMlrrcFiles(options, "j"); err == nil {
		t.Fatal("expected parse error for selected profile with bad line, got nil")
	}
}

func TestMlrrcGlobalParseErrorsStillFatal(t *testing.T) {
	writeTempMlrrc(t, "this-is-not-a-flag\n")
	options := cli.DefaultOptions()
	if err := loadMlrrcFiles(options, ""); err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestMlrrcPrepipeStillDisallowed(t *testing.T) {
	// Code-execution flags are disallowed in .mlrrc -- inside profiles too.
	writeTempMlrrc(t, "[j]\nprepipe zcat\n")
	options := cli.DefaultOptions()
	if err := loadMlrrcFiles(options, "j"); err == nil {
		t.Fatal("expected error for --prepipe within a profile, got nil")
	}
}

func TestMlrrcXdgConfigHomeIsLoaded(t *testing.T) {
	// No $HOME/.mlrrc, no ./.mlrrc, no $MLRRC: only
	// $XDG_CONFIG_HOME/miller/mlrrc is processed.
	xdgConfigHome := t.TempDir()
	if err := os.MkdirAll(filepath.Join(xdgConfigHome, "miller"), 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(xdgConfigHome, "miller", "mlrrc")
	if err := os.WriteFile(path, []byte("ojson\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("MLRRC", "")
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", xdgConfigHome)
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldwd) }()

	options := cli.DefaultOptions()
	if err := loadMlrrcFiles(options, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if options.WriterOptions.OutputFileFormat != "json" {
		t.Errorf("output format: got %s, want json", options.WriterOptions.OutputFileFormat)
	}
}

func TestMlrrcXdgConfigHomeDefaultsToHomeConfig(t *testing.T) {
	// When $XDG_CONFIG_HOME is unset, $HOME/.config/miller/mlrrc is used.
	home := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, ".config", "miller"), 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(home, ".config", "miller", "mlrrc")
	if err := os.WriteFile(path, []byte("ojson\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("MLRRC", "")
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", "")
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldwd) }()

	options := cli.DefaultOptions()
	if err := loadMlrrcFiles(options, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if options.WriterOptions.OutputFileFormat != "json" {
		t.Errorf("output format: got %s, want json", options.WriterOptions.OutputFileFormat)
	}
}

func TestMlrrcXdgConfigHomeStacksAfterHomeMlrrc(t *testing.T) {
	// Both $HOME/.mlrrc and $XDG_CONFIG_HOME/miller/mlrrc are processed, with
	// $HOME/.mlrrc applied first: a setting in the XDG file overrides one
	// from $HOME/.mlrrc.
	home := t.TempDir()
	if err := os.WriteFile(filepath.Join(home, ".mlrrc"), []byte("icsv\nojson\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	xdgConfigHome := t.TempDir()
	if err := os.MkdirAll(filepath.Join(xdgConfigHome, "miller"), 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(xdgConfigHome, "miller", "mlrrc")
	if err := os.WriteFile(path, []byte("otsv\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("MLRRC", "")
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", xdgConfigHome)
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldwd) }()

	options := cli.DefaultOptions()
	if err := loadMlrrcFiles(options, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if options.ReaderOptions.InputFileFormat != "csv" {
		t.Errorf("input format: got %s, want csv", options.ReaderOptions.InputFileFormat)
	}
	if options.WriterOptions.OutputFileFormat != "tsv" {
		t.Errorf("output format: got %s, want tsv (from XDG config, overriding json from $HOME/.mlrrc)", options.WriterOptions.OutputFileFormat)
	}
}

func TestMlrrcProfileFlagDisallowedWithinMlrrc(t *testing.T) {
	// Profiles are selected on the mlr command line, not from within a
	// .mlrrc file: --profile / -P there is a parse error, like --prepipe.
	for _, line := range []string{"profile j", "--profile j", "-P j"} {
		writeTempMlrrc(t, line+"\n")
		options := cli.DefaultOptions()
		err := loadMlrrcFiles(options, "")
		if err == nil {
			t.Errorf("expected parse error for %q within .mlrrc, got nil", line)
		} else if !strings.Contains(err.Error(), "parse error") {
			t.Errorf("expected parse error for %q within .mlrrc, got: %v", line, err)
		}
	}
	// Inside a selected profile section, too.
	writeTempMlrrc(t, "[j]\nprofile p\n")
	options := cli.DefaultOptions()
	if err := loadMlrrcFiles(options, "j"); err == nil {
		t.Fatal("expected parse error for --profile within a profile section, got nil")
	}
}
