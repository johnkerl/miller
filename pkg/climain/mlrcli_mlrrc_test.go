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
	if err := loadMlrrc(options, ""); err != nil {
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
	if err := loadMlrrc(options, ""); err != nil {
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
	if err := loadMlrrc(options, "j"); err != nil {
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
	if err := loadMlrrc(options, "j"); err != nil {
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
	if err := loadMlrrc(options, "p"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if options.WriterOptions.OutputFileFormat != "pprint" {
		t.Errorf("output format: got %s, want pprint", options.WriterOptions.OutputFileFormat)
	}
}

func TestMlrrcMissingProfileIsError(t *testing.T) {
	writeTempMlrrc(t, testMlrrcContents)
	options := cli.DefaultOptions()
	err := loadMlrrc(options, "nonesuch")
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
	err := loadMlrrc(options, "j")
	if err == nil {
		t.Fatal("expected error for profile with MLRRC=__none__, got nil")
	}
}

func TestMlrrcProfileWithNoMlrrcFileIsError(t *testing.T) {
	// Point MLRRC at a nonexistent file: an unopenable file is silently
	// skipped (the normal no-.mlrrc case), and HOME is remapped to an empty
	// temp directory, so no .mlrrc file is processed at all.
	t.Setenv("MLRRC", filepath.Join(t.TempDir(), "nonexistent"))
	t.Setenv("HOME", t.TempDir())
	options := cli.DefaultOptions()
	err := loadMlrrc(options, "j")
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
	if err := loadMlrrc(options, "j"); err != nil {
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
		err := loadMlrrc(options, "")
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
	if err := loadMlrrc(options, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// But it is a parse error when the section is selected.
	options = cli.DefaultOptions()
	if err := loadMlrrc(options, "j"); err == nil {
		t.Fatal("expected parse error for selected profile with bad line, got nil")
	}
}

func TestMlrrcGlobalParseErrorsStillFatal(t *testing.T) {
	writeTempMlrrc(t, "this-is-not-a-flag\n")
	options := cli.DefaultOptions()
	if err := loadMlrrc(options, ""); err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestMlrrcPrepipeStillDisallowed(t *testing.T) {
	// Code-execution flags are disallowed in .mlrrc -- inside profiles too.
	writeTempMlrrc(t, "[j]\nprepipe zcat\n")
	options := cli.DefaultOptions()
	if err := loadMlrrc(options, "j"); err == nil {
		t.Fatal("expected error for --prepipe within a profile, got nil")
	}
}
