// ================================================================
// TOOO: comment
// ================================================================

package regtest

import (
	"container/list"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"miller/src/lib"
	"miller/src/platform"
)

const DefaultPath = "./regtest/cases"
const CmdName = "cmd"
const EnvName = "env"
const PreCopyName = "precopy"
const ExpectedStdoutName = "expout"
const ExpectedStderrName = "experr"
const PostCompareName = "postcmp"
const ShouldFailName = "should-fail"

const MajorSeparator = "================================================================"
const MinorSeparator = "----------------------------------------------------------------"

type stringPair struct {
	first  string
	second string
}

// ----------------------------------------------------------------
type RegTester struct {
	exeName        string
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
	doPopulate bool,
	verbosityLevel int,
	firstNFailsToShow int,
) *RegTester {
	return &RegTester{
		exeName:            exeName,
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

	// Directory count may be zero if we were invoked with all paths on the
	// command line being .cmd files.
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
		basename := filepath.Base(path)
		if basename == CmdName {
			passed := this.executeSingleCmdFile(path, this.verbosityLevel)
			if passed {
				this.casePassCount++
			} else {
				this.caseFailCount++
				this.failCaseNames.PushBack(path)
			}
			return passed
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
	// TODO: comment
	isCaseDirectory := this.isCaseDirectory(dirName)

	if !isCaseDirectory && this.verbosityLevel >= 1 {
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
		if !isCaseDirectory {
			if passed {
				platform.PrintHiGreen("PASS")
				fmt.Printf(" %s/\n", dirName)
			} else {
				platform.PrintHiRed("FAIL")
				fmt.Printf(" %s/\n", dirName)
			}
		}
	}

	if !isCaseDirectory && this.verbosityLevel >= 1 {
		fmt.Printf("%s END   %s/\n", MajorSeparator, dirName)
		fmt.Println()
	}

	return passed
}

// ----------------------------------------------------------------
// Sees if a directory contains a single test case. If so, we don't want to
// print a banner at start and end at the default verbosity level.  In order to
// manage print volume, we only want to print for directories-of-cases (next
// level up).

// TODO: don't print container-of-containers, via entries.any.isCaseDirectory
// or somesuch.

func (this *RegTester) isCaseDirectory(
	dirName string,
) bool {
	cmdFilePath := dirName + string(filepath.Separator) + CmdName
	return this.FileExists(cmdFilePath)
}

// ----------------------------------------------------------------
// This is the main regression-test logic for a single .cmd file (a single mlr
// invocation) and its associated supporting files.
func (this *RegTester) executeSingleCmdFile(
	cmdFilePath string,
	verbosityLevel int,
) bool {

	if verbosityLevel >= 2 {
		fmt.Printf("%s begin %s\n", MinorSeparator, cmdFilePath)
		defer fmt.Printf("%s end   %s\n", MinorSeparator, cmdFilePath)
	}

	// Given 'regtest/cases/foo/0038.cmd', get 'regtest/cases/foo' and '0038'.
	// Various support files use syntax ${CASEDIR} within them so they're
	// relocatable, but we need to expand those in order to execute the test
	// case.
	caseDir := filepath.Dir(cmdFilePath)

	cmd, err := this.loadFile(cmdFilePath, caseDir)
	if err != nil {
		if verbosityLevel >= 2 {
			fmt.Printf("%s: %v\n", cmdFilePath, err)
		}
		return false
	}

	slash := string(filepath.Separator) // Or backslash on Windows ... although modern Windows versions handle slashes fine.
	envFileName := caseDir + slash + EnvName
	preCopyFileName := caseDir + slash + PreCopyName
	expectedStdoutFileName := caseDir + slash + ExpectedStdoutName
	expectedStderrFileName := caseDir + slash + ExpectedStderrName
	expectFailFileName := caseDir + slash + ShouldFailName
	postCompareFileName := caseDir + slash + PostCompareName

	cmd, err = this.loadFile(cmdFilePath, caseDir)
	if err != nil {
		if verbosityLevel >= 2 {
			fmt.Printf("%s: %v\n", cmdFilePath, err)
		}
		return false
	}

	if verbosityLevel >= 2 {
		fmt.Println("Command:")
		fmt.Println(cmd)
	}

	// The .env needn't exist (most test cases don't have one) in which case
	// the returned map will be empty.
	envKeyValuePairs, err := this.loadEnvFile(envFileName, caseDir)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// The .precopy needn't exist (most test cases don't have one) in which case
	// the returned map will be empty.
	preCopySrcDestPairs, err := this.loadStringPairFile(preCopyFileName, caseDir)
	if err != nil {
		fmt.Println(err)
		return false
	}

	passed := true

	// Set any case-specific environment variables before running the case.
	for pe := envKeyValuePairs.Head; pe != nil; pe = pe.Next {
		key := pe.Key
		value := pe.Value.(string)
		if verbosityLevel >= 3 {
			fmt.Printf("SETENV %s=%s\n", key, value)
		}
		os.Setenv(key, value)
	}

	// Copy any files requested by the test. (Most don't; some do, e.g. those
	// which test the write-in-place logic of mlr -I.)
	for pe := preCopySrcDestPairs.Front(); pe != nil; pe = pe.Next() {
		pair := pe.Value.(stringPair)
		src := pair.first
		dst := pair.second
		if verbosityLevel >= 3 {
			fmt.Printf("%s: copy %s to %s\n", cmdFilePath, src, dst)
		}
		err := this.copyFile(src, dst, caseDir)
		if err != nil {
			fmt.Printf("%s: %v\n", dst, err)
			passed = false
		}
	}

	// ****************************************************************
	// HERE IS WHERE WE RUN THE MILLER COMMAND LINE FOR THE TEST CASE
	actualStdout, actualStderr, actualExitCode, err := RunMillerCommand(this.exeName, cmd)
	// ****************************************************************

	// Unset any case-specific environment variables after running the case.
	// This is important since the setenv is done in the current process,
	// and we don't want to affect subsequent test cases.
	for pe := envKeyValuePairs.Head; pe != nil; pe = pe.Next {
		key := pe.Key
		if verbosityLevel >= 3 {
			fmt.Printf("UNSETENV %s\n", key)
		}
		os.Setenv(key, "")
	}

	// The .postcmp needn't exist (most test cases don't have one) in which case
	// the returned map will be empty.
	postCompareExpectedActualPairs, err := this.loadStringPairFile(postCompareFileName, caseDir)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if this.doPopulate {
		// Populate mode: write out the actual stdout/stderr/exit-code to disk
		// as expected values for subsequent runs.

		// TODO: temp replace-all for CR/LF to LF. Will need re-work once auto-detect is ported.
		actualStdout = strings.ReplaceAll(actualStdout, "\r\n", "\n")
		actualStderr = strings.ReplaceAll(actualStderr, "\r\n", "\n")

		// Write the .expout file
		err = this.storeFile(expectedStdoutFileName, actualStdout)
		if err != nil {
			fmt.Printf("%s: %v\n", expectedStdoutFileName, err)
			passed = false
		} else {
			if this.verbosityLevel >= 1 {
				fmt.Printf("wrote %s\n", expectedStdoutFileName)
			}
		}

		// Write the .experr file
		err = this.storeFile(expectedStderrFileName, actualStderr)
		if err != nil {
			fmt.Printf("%s: %v\n", expectedStderrFileName, err)
			passed = false
		} else {
			if this.verbosityLevel >= 1 {
				fmt.Printf("wrote %s\n", expectedStdoutFileName)
			}
		}

		// Write the .should-fail file
		if actualExitCode != 0 {
			err = this.storeFile(expectFailFileName, "")
			if err != nil {
				fmt.Printf("%s: %v\n", expectedStderrFileName, err)
				passed = false
			} else {
				if this.verbosityLevel >= 1 {
					fmt.Printf("wrote %s\n", expectedStdoutFileName)
				}
			}

		}

		for pe := postCompareExpectedActualPairs.Front(); pe != nil; pe = pe.Next() {
			pair := pe.Value.(stringPair)
			expectedFileName := pair.first
			actualFileName := pair.second

			err := this.copyFile(actualFileName, expectedFileName, caseDir)
			if err != nil {
				fmt.Printf("Could not copy %s to %s: %v\n", actualFileName, expectedFileName, err)
				passed = false
			}
			if verbosityLevel >= 3 {
				fmt.Printf("Copied %s to %s: %v\n", actualFileName, expectedFileName, err)
			}
		}

	} else {
		// Verify mode: check actuals against expecteds

		// Load the .expout file
		expectedStdout, err := this.loadFile(expectedStdoutFileName, caseDir)
		if err != nil {
			if verbosityLevel >= 2 {
				fmt.Printf("%s: %v\n", expectedStdoutFileName, err)
			}
			return false
		}

		// Load the .experr file
		expectedStderr, err := this.loadFile(expectedStderrFileName, caseDir)
		if err != nil {
			if verbosityLevel >= 2 {
				fmt.Printf("%s: %v\n", expectedStderrFileName, err)
			}
			return false
		}

		// Load the .should-fail file
		expectedExitCode := 0
		if this.FileExists(expectFailFileName) {
			expectedExitCode = 1
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

		// Compare stdout to .expout
		if actualStdout != expectedStdout {
			if verbosityLevel >= 2 {
				fmt.Printf(
					"%s: stdout does not match expected %s\n",
					cmdFilePath,
					expectedStdoutFileName,
				)
			}
			passed = false
		}

		// Compare stderr to .experr
		if actualStderr != expectedStderr {
			if verbosityLevel >= 2 {
				fmt.Printf(
					"%s: stderr does not match expected %s\n",
					cmdFilePath,
					expectedStderrFileName,
				)
			}
			// TODO: needs normalization of os.Args[0] -> "mlr" throughout the codebase,
			// else we get spurious mismatch between expected strings like 'mlr: ...'
			// and actuals like 'C:\miller\go\mlr.exe: ...'

			// passed = false
		}

		// Compare exit code
		if actualExitCode != expectedExitCode {
			if verbosityLevel >= 2 {
				fmt.Printf(
					"%s: exit code %d does not match expected %d\n",
					cmdFilePath,
					actualExitCode, expectedExitCode,
				)
			}
			passed = false
		}

		// Compare any additional output files. Most test cases don't have
		// these (just stdout/stderr), but some do: for example, those which
		// test the tee verb/function.
		for pe := postCompareExpectedActualPairs.Front(); pe != nil; pe = pe.Next() {
			pair := pe.Value.(stringPair)
			expectedFileName := pair.first
			actualFileName := pair.second
			ok, err := this.compareFiles(expectedFileName, actualFileName, caseDir)
			if err != nil {
				fmt.Printf("%s: %v\n", cmdFilePath, err)
				passed = false
			} else if !ok {
				if verbosityLevel >= 2 {
					fmt.Printf(
						"%s: %s does not match %s\n",
						cmdFilePath, expectedFileName, actualFileName,
					)
				}
				// TODO: if verbosityLevel >= 3, print the contents of both files
				passed = false
			} else {
				if verbosityLevel >= 2 {
					fmt.Printf(
						"%s: %s matches %s\n",
						cmdFilePath, expectedFileName, actualFileName,
					)
				}
			}
		}

		// Clean up any requested file-copies so that we're git-clean after the regression-test run.
		for pe := preCopySrcDestPairs.Front(); pe != nil; pe = pe.Next() {
			pair := pe.Value.(stringPair)
			dst := pair.second
			os.Remove(dst)
			if verbosityLevel >= 3 {
				fmt.Printf("%s: clean up %s\n", cmdFilePath, dst)
			}
		}

		// Clean up any extra output files so that we're git-clean after the regression-test run.
		for pe := postCompareExpectedActualPairs.Front(); pe != nil; pe = pe.Next() {
			pair := pe.Value.(stringPair)
			actualFileName := pair.second
			os.Remove(actualFileName)
			if verbosityLevel >= 3 {
				fmt.Printf("%s: clean up %s\n", cmdFilePath, actualFileName)
			}
		}
	}

	if verbosityLevel >= 1 {
		if passed {
			platform.PrintHiGreen("pass")
			fmt.Printf(" %s\n", cmdFilePath)
		} else {
			platform.PrintHiRed("fail")
			fmt.Printf(" %s\n", cmdFilePath)
		}
	}

	return passed
}

// ----------------------------------------------------------------
func (this *RegTester) FileExists(fileName string) bool {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}

func (this *RegTester) loadFile(
	fileName string,
	caseDir string,
) (string, error) {
	byteContents, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	contents := string(byteContents)
	contents = strings.ReplaceAll(contents, "${CASEDIR}", caseDir)
	return contents, nil
}

func (this *RegTester) storeFile(
	fileName string,
	contents string,
) error {
	err := os.WriteFile(fileName, []byte(contents), 0666)
	if err != nil {
		return err
	}
	return nil
}

func (this *RegTester) copyFile(
	src string,
	dst string,
	caseDir string,
) error {
	contents, err := this.loadFile(src, caseDir)
	if err != nil {
		return err
	}
	err = this.storeFile(dst, contents)
	if err != nil {
		return err
	}
	return nil
}

func (this *RegTester) compareFiles(
	expectedFileName string,
	actualFileName string,
	caseDir string,
) (bool, error) {
	expectedContents, err := this.loadFile(expectedFileName, caseDir)
	if err != nil {
		return false, err
	}
	actualContents, err := this.loadFile(actualFileName, caseDir)
	if err != nil {
		return false, err
	}
	// TODO: maybe rethink later with autoterm
	expectedContents = strings.ReplaceAll(expectedContents, "\r\n", "\n")
	actualContents = strings.ReplaceAll(actualContents, "\r\n", "\n")

	return expectedContents == actualContents, nil
}

// ----------------------------------------------------------------
func (this *RegTester) loadEnvFile(
	filename string,
	caseDir string,
) (*lib.OrderedMap, error) {
	// If the file doesn't exist that's the normal case -- most cases do not
	// have a .env file.
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return lib.NewOrderedMap(), nil
	}

	// If the file does exist and isn't loadable, that's an error.
	contents, err := this.loadFile(filename, caseDir)
	if err != nil {
		return nil, err
	}

	keyValuePairs := lib.NewOrderedMap()
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
					"%s: could not parse line \"%s\" from file \"%s\".\n",
					lib.MlrExeName(), line, filename,
				),
			)
		}
		keyValuePairs.Put(fields[0], fields[1])
	}
	return keyValuePairs, nil
}

// ----------------------------------------------------------------
func (this *RegTester) loadStringPairFile(
	filename string,
	caseDir string,
) (*list.List, error) {
	// If the file doesn't exist that's the normal case -- most cases do not
	// have a .precopy or .postcmp file.
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return list.New(), nil
	}

	// If the file does exist and isn't loadable, that's an error.
	contents, err := this.loadFile(filename, caseDir)
	if err != nil {
		return nil, err
	}

	pairs := list.New()
	lines := strings.Split(contents, "\n")
	for _, line := range lines {
		line = strings.TrimSuffix(line, "\r")
		if line == "" {
			continue
		}
		fields := strings.SplitN(line, " ", 2) // TODO: split on multi-space
		if len(fields) != 2 {
			return nil, errors.New(
				fmt.Sprintf(
					"%s: could not parse line \"%s\" from file \"%s\".\n",
					lib.MlrExeName(), line, filename,
				),
			)
		}
		pair := stringPair{first: fields[0], second: fields[1]}
		pairs.PushBack(pair)
	}
	return pairs, nil
}
