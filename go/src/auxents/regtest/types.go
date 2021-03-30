// ================================================================
// TOOO
// ================================================================

package regtest

import (
	"fmt"
	"os"
	"strings"

	"miller/regression/support"
)

const DefaultPath = "./reg-test/cases"
const CommandSuffix = ".cmd"
const ExpectedStdoutSuffix = ".expout"
const ExpectedStderrSuffix = ".experr"

const MajorSeparator = "================================================================"
const MinorSeparator = "----------------------------------------------------------------"

// ----------------------------------------------------------------
type RegTester struct {
	exeName        string
	verbName       string
	verbosityLevel int
}

// ----------------------------------------------------------------
func NewRegTester(
	exeName string,
	verbName string,
	verbosityLevel int,
) *RegTester {
	return &RegTester{
		exeName:        exeName,
		verbName:       verbName,
		verbosityLevel: verbosityLevel,
	}
}

// ----------------------------------------------------------------
// TODO: comment
// TODO: []error since we want to report all ...
// TODO: maybe just a bool since report-all-to-terminal API instead?

func (this *RegTester) Execute(
	paths []string,
) error {
	if len(paths) == 0 {
		paths = []string{DefaultPath}
	}

	for _, path := range paths {
		this.executeSinglePath(path)
	}

	return nil
}

// ----------------------------------------------------------------
// TODO: comment
// TODO: []error since we want to report all ...

func (this *RegTester) executeSinglePath(
	path string,
) error {

	handle, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		return err // xxx or nah
	}
	switch mode := handle.Mode(); {
	case mode.IsDir():
		return this.executeSingleDirectory(path)
	case mode.IsRegular():
		return this.executeSingleFile(path)
	}

	return nil
}

// ----------------------------------------------------------------
func (this *RegTester) executeSingleDirectory(
	dirName string,
) error {
	fmt.Printf("%s BEGIN %s\n", MajorSeparator, dirName) // TODO: depending on verbosity level
	entries, err := os.ReadDir(dirName)
	if err != nil {
		fmt.Println(err)
		return err // xxx or nah
	}

	for i := range entries {
		entry := &entries[i]
		err := this.executeSinglePath(dirName + "/" + (*entry).Name())
		if err != nil {
			fmt.Println(err)
			return err // xxx or nah
		}
	}

	fmt.Printf("%s END   %s\n", MajorSeparator, dirName) // TODO: depending on verbosity level
	return nil
}

// ----------------------------------------------------------------
func (this *RegTester) executeSingleFile(
	fileName string,
) error {
	if !strings.HasSuffix(fileName, CommandSuffix) {
		return nil
	}

	fmt.Printf("%s BEGIN %s\n", MinorSeparator, fileName) // TODO: depending on verbosity level

	byteContents, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return err // xxx or nah
	}
	contents := string(byteContents)

    stdout, stderr, exitCode, err := support.RunMillerCommand(this.exeName, contents)
	// xxx expected exitCode
	fmt.Println("stdout <<", stdout, ">>")
	fmt.Println("stderr <<", stderr, ">>")
	fmt.Println("exitCode", exitCode)

	fmt.Printf("%s END   %s\n", MinorSeparator, fileName) // TODO: depending on verbosity level
	return nil
}
