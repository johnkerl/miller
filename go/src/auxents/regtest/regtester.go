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

	directoryPassCount int
	directoryFailCount int
	casePassCount      int
	caseFailCount      int
}

// ----------------------------------------------------------------
func NewRegTester(
	exeName string,
	verbName string,
	verbosityLevel int,
) *RegTester {
	return &RegTester{
		exeName:            exeName,
		verbName:           verbName,
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
// TODO: comment

func (this *RegTester) Execute(
	paths []string,
) bool {

	this.resetCounts()
	// TODO: print using-exe

	if len(paths) == 0 {
		paths = []string{DefaultPath}
	}

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
// TODO: comment

func (this *RegTester) executeSinglePath(
	path string,
) bool {
	handle, err := os.Stat(path)
	if err != nil {
		fmt.Printf("%s: %v\n", path, err)
		// TODO: attribs ...
		return false
	}
	mode := handle.Mode()
	if mode.IsDir() {
		ok := this.executeSingleDirectory(path)
		if ok {
			this.directoryPassCount++
		} else {
			this.directoryFailCount++
		}
	} else if mode.IsRegular() {
		if strings.HasSuffix(path, CommandSuffix) {
			ok := this.executeSingleCmdFile(path)
			if ok {
				this.casePassCount++
			} else {
				this.caseFailCount++
			}
		}
	}
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
	}

	return passed
}

// ----------------------------------------------------------------
// TODO: comment
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
func (this *RegTester) executeSingleCmdFile(
	cmdFileName string,
) bool {

	if this.verbosityLevel >= 2 {
		fmt.Printf("%s begin %s\n", MinorSeparator, cmdFileName)
		defer fmt.Printf("%s end   %s\n", MinorSeparator, cmdFileName)
	}

	expectedStdoutFileName := this.changeExtension(cmdFileName, CommandSuffix, ExpectedStdoutSuffix)
	expectedStderrFileName := this.changeExtension(cmdFileName, CommandSuffix, ExpectedStderrSuffix)

	cmd, err := this.loadFile(cmdFileName)
	if err != nil {
		if this.verbosityLevel >= 3 {
			fmt.Printf("%s: %v\n", cmdFileName, err)
		}
		return false
	}
	expectedStdout, err := this.loadFile(expectedStdoutFileName)
	if err != nil {
		if this.verbosityLevel >= 3 {
			fmt.Printf("%s: %v\n", expectedStdoutFileName, err)
		}
		return false
	}
	expectedStderr, err := this.loadFile(expectedStderrFileName)
	if err != nil {
		if this.verbosityLevel >= 3 {
			fmt.Printf("%s: %v\n", expectedStderrFileName, err)
		}
		return false
	}

	passed := true

	actualStdout, actualStderr, actualExitCode, err := RunMillerCommand(this.exeName, cmd)
	// xxx temp
	expectedExitCode := 0

	if this.verbosityLevel >= 3 {

		fmt.Println("actualStdout:")
		fmt.Println(actualStdout)

		fmt.Println("expectedStdout:")
		fmt.Println(expectedStdout)

		fmt.Println("actualStderr:")
		fmt.Println(actualStderr)

		fmt.Println("expectedStderr:")
		fmt.Println(expectedStderr)

		fmt.Println("actualExitCode:")
		fmt.Println(actualExitCode)

		fmt.Println("expectedExitCode:")
		fmt.Println(expectedExitCode)
	}

	// xxx Windows CR/LF <-> LF handling
	if actualStdout != expectedStdout {
		if this.verbosityLevel >= 2 {
			fmt.Printf(
				"stdout does not match expected %s\n",
				expectedStdoutFileName,
			)
		}
		passed = false
	}

	// xxx Windows CR/LF <-> LF handling
	if actualStderr != expectedStderr {
		if this.verbosityLevel >= 2 {
			fmt.Printf(
				"stderr does not match expected %s\n",
				expectedStderrFileName,
			)
		}
		passed = false
	}

	if actualExitCode != expectedExitCode {
		if this.verbosityLevel >= 2 {
			fmt.Printf(
				"Exit code %d does not match expected %d\n",
				actualExitCode, expectedExitCode,
			)
		}
		passed = false
	}

	if this.verbosityLevel >= 2 {
		if passed {
			fmt.Printf("%s %s\n", pass, cmdFileName)
		} else {
			fmt.Printf("%s %s\n", fail, cmdFileName)
		}
	}

	return false
}

// ----------------------------------------------------------------
func (this *RegTester) changeExtension(
	fileName string,
	oldExtension string,
	newExtension string,
) string {
	return strings.TrimSuffix(fileName, oldExtension) + newExtension
}

func (this *RegTester) loadFile(
	fileName string,
) (string, error) {
	byteContents, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("%s: %v\n", fileName, err)
		return "", err // xxx or nah
	}
	return string(byteContents), nil
}
