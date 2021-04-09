// ================================================================
// TOOO
// ================================================================

package regtest

import (
	"container/list"
	"errors"
	"fmt"
	"os"
	"strings"

	"miller/src/lib"
	"miller/src/platform"
)

const DefaultPath = "./reg-test/cases"
const CommandSuffix = ".cmd"
const EnvSuffix = ".env"
const ExpectedStdoutSuffix = ".expout"
const ExpectedStderrSuffix = ".experr"
const ShouldFailSuffix = ".should-fail"

const MajorSeparator = "================================================================"
const MinorSeparator = "----------------------------------------------------------------"

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

	failDirNames      *list.List
	failCaseNames     *list.List
	firstNFailsToShow int
}

// ----------------------------------------------------------------
func NewRegTester(
	exeName string,
	verbName string,
	doPopulate bool,
	verbosityLevel int,
	firstNFailsToShow int,
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
		failDirNames:       list.New(),
		failCaseNames:      list.New(),
		firstNFailsToShow:  firstNFailsToShow,
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

	if this.failCaseNames.Len() > 0 && this.firstNFailsToShow > 0 {
		fmt.Println()
		fmt.Println("RERUNS OF FIRST FAILED CASE FILES:")
		verbosityLevel := 3
		i := 0
		for e := this.failCaseNames.Front(); e != nil; e = e.Next() {
			this.executeSingleCmdFile(e.Value.(string), verbosityLevel)
			i++
			if i >= this.firstNFailsToShow {
				break
			}
		}
	}

	if this.failDirNames.Len() > 0 {
		fmt.Println()
		fmt.Println("FAILED CASE DIRECTORIES:")
		for e := this.failDirNames.Front(); e != nil; e = e.Next() {
			fmt.Printf("  %s/\n", e.Value.(string))
		}
	}

	fmt.Println()
	fmt.Printf("NUMBER OF CASES            PASSED %d\n", this.casePassCount)
	fmt.Printf("NUMBER OF CASES            FAILED %d\n", this.caseFailCount)
	fmt.Printf("NUMBER OF CASE-DIRECTORIES PASSED %d\n", this.directoryPassCount)
	fmt.Printf("NUMBER OF CASE-DIRECTORIES FAILED %d\n", this.directoryFailCount)
	fmt.Println()

	if this.casePassCount > 0 && this.caseFailCount == 0 {
		platform.PrintHiGreen("PASS")
		fmt.Printf(" overall\n")
		return true
	} else {
		platform.PrintHiRed("FAIL")
		fmt.Printf(" overall\n")
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
			this.failDirNames.PushBack(path)
		}
		return passed
	} else if mode.IsRegular() {
		if strings.HasSuffix(path, CommandSuffix) {
			if this.doPopulate {
				this.populateSingleCmdFile(path)
				return true
			} else {
				passed := this.executeSingleCmdFile(path, this.verbosityLevel)
				if passed {
					this.casePassCount++
				} else {
					this.caseFailCount++
					this.failCaseNames.PushBack(path)
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
		fmt.Printf("%s BEGIN %s/\n", MajorSeparator, dirName)
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
				platform.PrintHiGreen("PASS")
				fmt.Printf(" %s/\n", dirName)
			} else {
				platform.PrintHiRed("FAIL")
				fmt.Printf(" %s/\n", dirName)
			}
		}
	}

	if hasDirectEntries && this.verbosityLevel >= 1 {
		fmt.Printf("%s END   %s/\n", MajorSeparator, dirName)
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
	envFileName := this.changeExtension(cmdFileName, CommandSuffix, EnvSuffix)

	cmd, err := this.loadFile(cmdFileName)
	if err != nil {
		fmt.Printf("%s: %v\n", cmdFileName, err)
		return
	}

	if this.verbosityLevel >= 2 {
		fmt.Println("Command:")
		fmt.Println(cmd)
	}

	// The .env needn't exist (most test cases don't have one) in which case
	// the envKeyValuePairs map will be empty.
	envKeyValuePairs, err := this.loadEnvFile(envFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set any case-specific environment variables before running the case.
	for key, value := range envKeyValuePairs {
		if this.verbosityLevel >= 3 {
			fmt.Printf("SETENV %s=%s\n", key, value)
		}
		os.Setenv(key, value)
	}

	actualStdout, actualStderr, actualExitCode, err := RunMillerCommand(this.exeName, cmd)

	// Unset any case-specific environment variables after running the case.
	// This is important since the setenv is done in the current process,
	// and we don't want to affect subsequent test cases.
	for key, _ := range envKeyValuePairs {
		if this.verbosityLevel >= 3 {
			fmt.Printf("UNSETENV %s\n", key)
		}
		os.Setenv(key, "")
	}

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
	verbosityLevel int,
) bool {

	if verbosityLevel >= 2 {
		fmt.Printf("%s begin %s\n", MinorSeparator, cmdFileName)
		defer fmt.Printf("%s end   %s\n", MinorSeparator, cmdFileName)
	}

	expectedStdoutFileName := this.changeExtension(cmdFileName, CommandSuffix, ExpectedStdoutSuffix)
	expectedStderrFileName := this.changeExtension(cmdFileName, CommandSuffix, ExpectedStderrSuffix)
	expectFailFileName := this.changeExtension(cmdFileName, CommandSuffix, ShouldFailSuffix)
	envFileName := this.changeExtension(cmdFileName, CommandSuffix, EnvSuffix)

	cmd, err := this.loadFile(cmdFileName)
	if err != nil {
		if verbosityLevel >= 2 {
			fmt.Printf("%s: %v\n", cmdFileName, err)
		}
		return false
	}

	if verbosityLevel >= 2 {
		fmt.Println("Command:")
		fmt.Println(cmd)
	}

	// The .env needn't exist (most test cases don't have one) in which case
	// the envKeyValuePairs map will be empty.
	envKeyValuePairs, err := this.loadEnvFile(envFileName)
	if err != nil {
		fmt.Println(err)
		return false
	}

	expectedStdout, err := this.loadFile(expectedStdoutFileName)
	if err != nil {
		if verbosityLevel >= 2 {
			fmt.Printf("%s: %v\n", expectedStdoutFileName, err)
		}
		return false
	}
	expectedStderr, err := this.loadFile(expectedStderrFileName)
	if err != nil {
		if verbosityLevel >= 2 {
			fmt.Printf("%s: %v\n", expectedStderrFileName, err)
		}
		return false
	}
	expectedExitCode := 0
	if this.FileExists(expectFailFileName) {
		expectedExitCode = 1
	}

	passed := true

	// Set any case-specific environment variables before running the case.
	for key, value := range envKeyValuePairs {
		if verbosityLevel >= 3 {
			fmt.Printf("SETENV %s=%s\n", key, value)
		}
		os.Setenv(key, value)
	}

	actualStdout, actualStderr, actualExitCode, err := RunMillerCommand(this.exeName, cmd)

	// Unset any case-specific environment variables after running the case.
	// This is important since the setenv is done in the current process,
	// and we don't want to affect subsequent test cases.
	for key, _ := range envKeyValuePairs {
		if verbosityLevel >= 3 {
			fmt.Printf("UNSETENV %s\n", key)
		}
		os.Setenv(key, "")
	}

	if verbosityLevel >= 3 {
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
		if verbosityLevel >= 2 {
			fmt.Printf(
				"%s: stdout does not match expected %s\n",
				cmdFileName,
				expectedStdoutFileName,
			)
		}
		passed = false
	}

	if actualStderr != expectedStderr {
		if verbosityLevel >= 2 {
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
		if verbosityLevel >= 2 {
			fmt.Printf(
				"%s: exit code %d does not match expected %d\n",
				cmdFileName,
				actualExitCode, expectedExitCode,
			)
		}
		passed = false
	}

	if verbosityLevel >= 1 {
		if passed {
			platform.PrintHiGreen("pass")
			fmt.Printf(" %s\n", cmdFileName)
		} else {
			platform.PrintHiRed("fail")
			fmt.Printf(" %s\n", cmdFileName)
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

// ----------------------------------------------------------------
func (this *RegTester) loadEnvFile(
	envFileName string,
) (map[string]string, error) {
	// If the file doesn't exist that's the normal case -- most cases do not
	// have a .env file.
	_, err := os.Stat(envFileName)
	if os.IsNotExist(err) {
		return nil, nil
	}

	// If the file does exist and isn't loadable, that's an error.
	contents, err := this.loadFile(envFileName)
	if err != nil {
		return nil, err
	}

	keyValuePairs := make(map[string]string)
	lines := strings.Split(contents, "\n")
	for _, line := range lines {
		line = strings.TrimSuffix(line, "\r")
		if line == "" {
			continue
		}
		fields := strings.SplitN(line, "=", 2)
		if len(fields) != 2 {
			return nil, errors.New(
				fmt.Sprintf(
					"%s: could not parse env line \"%s\" from file \"%s\".\n",
					lib.MlrExeName(), line, envFileName,
				),
			)
		}
		keyValuePairs[fields[0]] = fields[1]
	}
	return keyValuePairs, nil
}
