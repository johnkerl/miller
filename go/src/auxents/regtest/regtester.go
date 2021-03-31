// ================================================================
// TOOO
// ================================================================

package regtest

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

const DefaultPath = "./reg-test/cases"
const CommandSuffix = ".cmd"
const ExpectedStdoutSuffix = ".expout"
const ExpectedStderrSuffix = ".experr"
const ShouldFailSuffix = ".should-fail"

const MajorSeparator = "================================================================"
const MinorSeparator = "----------------------------------------------------------------"

var PASS = color.HiGreenString("PASS")
var pass = color.HiGreenString("pass")
var FAIL = color.HiRedString("FAIL")
var fail = color.HiRedString("fail")

// ----------------------------------------------------------------
type RegTester struct {
	exeName        string
	verbName       string
	verbosityLevel int
	doPopulate     bool

	directoryPassCount int
	directoryFailCount int
	casePassCount      int
	caseFailCount      int
}

// ----------------------------------------------------------------
func NewRegTester(
	exeName string,
	verbName string,
	doPopulate bool,
	verbosityLevel int,
) *RegTester {
	return &RegTester{
		exeName:            exeName,
		verbName:           verbName,
		doPopulate:         doPopulate,
		verbosityLevel:     verbosityLevel,
		directoryPassCount: 0,
		directoryFailCount: 0,
		casePassCount:      0,
		caseFailCount:      0,
	}
}

func (this *RegTester) resetCounts() {
	this.directoryPassCount = 0
	this.directoryFailCount = 0
	this.casePassCount = 0
	this.caseFailCount = 0
}

// ----------------------------------------------------------------
// Top-level entrypoint for the regtester. See the usage function in entry.go
// for semantics.

func (this *RegTester) Execute(
	paths []string,
) bool {

	this.resetCounts()

	if len(paths) == 0 {
		paths = []string{DefaultPath}
	}

	fmt.Println("REGRESSION TEST:")
	for _, path := range paths {
		fmt.Printf("  %s\n", path)
	}
	fmt.Printf("Using executable: %s\n", this.exeName)
	fmt.Println()

	for _, path := range paths {
		this.executeSinglePath(path)
	}

	fmt.Println()
	fmt.Printf("NUMBER OF CASE-DIRECTORIES PASSED %d\n", this.directoryPassCount)
	fmt.Printf("NUMBER OF CASE-DIRECTORIES FAILED %d\n", this.directoryFailCount)
	fmt.Printf("NUMBER OF CASES            PASSED %d\n", this.casePassCount)
	fmt.Printf("NUMBER OF CASES            FAILED %d\n", this.caseFailCount)
	fmt.Println()

	if this.casePassCount > 0 && this.caseFailCount == 0 {
		fmt.Printf("%s overall\n", PASS)
		return true
	} else {
		fmt.Printf("%s overall\n", FAIL)
		return false
	}
}

// ----------------------------------------------------------------
// Recursively invoked routine to process either a single .cmd file, or a
// directory of such, or a directory of directories.

func (this *RegTester) executeSinglePath(
	path string,
) bool {
	handle, err := os.Stat(path)
	if err != nil {
		fmt.Printf("%s: %v\n", path, err)
		return false
	}
	mode := handle.Mode()
	if mode.IsDir() {
		passed := this.executeSingleDirectory(path)
		if passed {
			this.directoryPassCount++
		} else {
			this.directoryFailCount++
		}
		return passed
	} else if mode.IsRegular() {
		if strings.HasSuffix(path, CommandSuffix) {
			if this.doPopulate {
				this.populateSingleCmdFile(path)
				return true
			} else {
				passed := this.executeSingleCmdFile(path)
				if passed {
					this.casePassCount++
				} else {
					this.caseFailCount++
				}
				return passed
			}
		}
		return true // No .cmd files directly inside
	}

	fmt.Printf("%s: neither directory nor regular file.\n", path)
	return false // fall-through
}

// ----------------------------------------------------------------
func (this *RegTester) executeSingleDirectory(
	dirName string,
) bool {
	passed := true
	hasDirectEntries := this.directoryHasDirectEntries(dirName)

	if hasDirectEntries && this.verbosityLevel >= 1 {
		fmt.Printf("%s BEGIN %s\n", MajorSeparator, dirName)
	}

	entries, err := os.ReadDir(dirName)
	if err != nil {
		fmt.Printf("%s: %v\n", dirName, err)
		passed = false
	} else {

		for i := range entries {
			entry := &entries[i]
			path := dirName + "/" + (*entry).Name()

			ok := this.executeSinglePath(path)
			if !ok {
				passed = false
			}
		}

		// Only print if there are .cmd files directly in this directory.
		// Otherwise it's just a directory-of-directories and we don't need to
		// multiply announce.
		if hasDirectEntries {
			if passed {
				fmt.Printf("%s %s\n", PASS, dirName)
			} else {
				fmt.Printf("%s %s\n", FAIL, dirName)
			}
		}
	}

	if hasDirectEntries && this.verbosityLevel >= 1 {
		fmt.Printf("%s END   %s\n", MajorSeparator, dirName)
		fmt.Println()
	}

	return passed
}

// ----------------------------------------------------------------
// Sees if a directory has .cmd files directly in it (vs in a subdirectory).
// If so, we want to print a banner at start and end.
func (this *RegTester) directoryHasDirectEntries(
	dirName string,
) bool {

	entries, err := os.ReadDir(dirName)
	if err != nil {
		fmt.Printf("%s: %v\n", dirName, err)
		return false
	}

	for i := range entries {
		entry := &entries[i]
		entryName := (*entry).Name()

		if strings.HasSuffix(entryName, CommandSuffix) {
			return true
		}
	}
	return false
}

// ----------------------------------------------------------------
// TODO: comment
func (this *RegTester) populateSingleCmdFile(
	cmdFileName string,
) {

	if this.verbosityLevel >= 1 {
		fmt.Printf("%s begin %s\n", MinorSeparator, cmdFileName)
		defer fmt.Printf("%s end   %s\n", MinorSeparator, cmdFileName)
	}

	expectedStdoutFileName := this.changeExtension(cmdFileName, CommandSuffix, ExpectedStdoutSuffix)
	expectedStderrFileName := this.changeExtension(cmdFileName, CommandSuffix, ExpectedStderrSuffix)

	cmd, err := this.loadFile(cmdFileName)
	if err != nil {
		fmt.Printf("%s: %v\n", cmdFileName, err)
		return
	}

	if this.verbosityLevel >= 2 {
		fmt.Println("Command:")
		fmt.Println(cmd)
	}

	actualStdout, actualStderr, actualExitCode, err := RunMillerCommand(this.exeName, cmd)

	if this.verbosityLevel >= 3 {

		fmt.Printf("actualStdout [%d]:\n", len(actualStdout))
		fmt.Println(actualStdout)

		fmt.Printf("actualStderr [%d]:\n", len(actualStderr))
		fmt.Println(actualStderr)

		fmt.Println("actualExitCode:")
		fmt.Println(actualExitCode)

		fmt.Println()
	}

	// TODO: temp replace-all for CR/LF to LF. Will need re-work once auto-detect is ported.
	actualStdout = strings.ReplaceAll(actualStdout, "\r\n", "\n")
	actualStderr = strings.ReplaceAll(actualStderr, "\r\n", "\n")

	err = this.storeFile(expectedStdoutFileName, actualStdout)
	if err != nil {
		fmt.Printf("%s: %v\n", expectedStdoutFileName, err)
		return
	}
	err = this.storeFile(expectedStderrFileName, actualStderr)
	if err != nil {
		fmt.Printf("%s: %v\n", expectedStderrFileName, err)
		return
	}

	if this.verbosityLevel >= 1 {
		fmt.Printf("wrote %s\n", cmdFileName)
	}
}

// ----------------------------------------------------------------
func (this *RegTester) executeSingleCmdFile(
	cmdFileName string,
) bool {

	if this.verbosityLevel >= 1 {
		fmt.Printf("%s begin %s\n", MinorSeparator, cmdFileName)
		defer fmt.Printf("%s end   %s\n", MinorSeparator, cmdFileName)
	}

	expectedStdoutFileName := this.changeExtension(cmdFileName, CommandSuffix, ExpectedStdoutSuffix)
	expectedStderrFileName := this.changeExtension(cmdFileName, CommandSuffix, ExpectedStderrSuffix)
	expectFailFileName := this.changeExtension(cmdFileName, CommandSuffix, ShouldFailSuffix)

	cmd, err := this.loadFile(cmdFileName)
	if err != nil {
		if this.verbosityLevel >= 2 {
			fmt.Printf("%s: %v\n", cmdFileName, err)
		}
		return false
	}

	if this.verbosityLevel >= 2 {
		fmt.Println("Command:")
		fmt.Println(cmd)
	}

	expectedStdout, err := this.loadFile(expectedStdoutFileName)
	if err != nil {
		if this.verbosityLevel >= 2 {
			fmt.Printf("%s: %v\n", expectedStdoutFileName, err)
		}
		return false
	}
	expectedStderr, err := this.loadFile(expectedStderrFileName)
	if err != nil {
		if this.verbosityLevel >= 2 {
			fmt.Printf("%s: %v\n", expectedStderrFileName, err)
		}
		return false
	}
	expectedExitCode := 0
	if this.FileExists(expectFailFileName) {
		expectedExitCode = 1
	}

	passed := true

	actualStdout, actualStderr, actualExitCode, err := RunMillerCommand(this.exeName, cmd)

	if this.verbosityLevel >= 3 {

		fmt.Printf("actualStdout [%d]:\n", len(actualStdout))
		fmt.Println(actualStdout)

		fmt.Printf("expectedStdout [%d]:\n", len(expectedStdout))
		fmt.Println(expectedStdout)

		fmt.Printf("actualStderr [%d]:\n", len(actualStderr))
		fmt.Println(actualStderr)

		fmt.Printf("expectedStderr [%d]:\n", len(expectedStderr))
		fmt.Println(expectedStderr)

		fmt.Println("actualExitCode:")
		fmt.Println(actualExitCode)

		fmt.Println("expectedExitCode:")
		fmt.Println(expectedExitCode)

		fmt.Println()
	}

	// TODO: temp replace-all for CR/LF to LF. Will need re-work once auto-detect is ported.
	actualStdout = strings.ReplaceAll(actualStdout, "\r\n", "\n")
	actualStderr = strings.ReplaceAll(actualStderr, "\r\n", "\n")
	expectedStdout = strings.ReplaceAll(expectedStdout, "\r\n", "\n")
	expectedStderr = strings.ReplaceAll(expectedStderr, "\r\n", "\n")

	if actualStdout != expectedStdout {
		if this.verbosityLevel >= 2 {
			fmt.Printf(
				"%s: stdout does not match expected %s\n",
				cmdFileName,
				expectedStdoutFileName,
			)
		}
		passed = false
	}

	if actualStderr != expectedStderr {
		if this.verbosityLevel >= 2 {
			fmt.Printf(
				"%s: stderr does not match expected %s\n",
				cmdFileName,
				expectedStderrFileName,
			)
		}
		// TODO: needs normalization of os.Args[0] -> "mlr" throughout the codebase,
		// else we get spurious mismatch between expected strings like 'mlr: ...'
		// and actuals like 'C:\miller\go\mlr.exe: ...'

		// passed = false
	}

	if actualExitCode != expectedExitCode {
		if this.verbosityLevel >= 2 {
			fmt.Printf(
				"%s: exit code %d does not match expected %d\n",
				cmdFileName,
				actualExitCode, expectedExitCode,
			)
		}
		passed = false
	}

	if this.verbosityLevel >= 1 {
		if passed {
			fmt.Printf("%s %s\n", pass, cmdFileName)
		} else {
			fmt.Printf("%s %s\n", fail, cmdFileName)
		}
	}

	return passed
}

// ----------------------------------------------------------------
func (this *RegTester) changeExtension(
	fileName string,
	oldExtension string,
	newExtension string,
) string {
	return strings.TrimSuffix(fileName, oldExtension) + newExtension
}

func (this *RegTester) FileExists(fileName string) bool {
	fileInfo, err := os.Stat(fileName)
	if err != nil { // TODO: neither true nor false; throw & abend the entire suite maybe
		return false
	}
	return !fileInfo.IsDir()
}

func (this *RegTester) loadFile(
	fileName string,
) (string, error) {
	byteContents, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("%s: %v\n", fileName, err)
		return "", err
	}
	return string(byteContents), nil
}

func (this *RegTester) storeFile(
	fileName string,
	contents string,
) error {
	err := os.WriteFile(fileName, []byte(contents), 0666)
	if err != nil {
		fmt.Printf("%s: %v\n", fileName, err)
		return err
	}
	return nil
}
