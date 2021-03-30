// ================================================================
// FOO STUB
// ================================================================

package regression

import (
	"os"
	"testing"

	"miller/regression/support"
)

func TestWorkingDirectory(t *testing.T) {
	wd, _ := os.Getwd()
	t.Log("pwd is:", wd)
}

func TestFoo(t *testing.T) {
	stdout, stderr, exitCode, err := support.RunMillerCommand(os.Args[0], "cat testdata/abixy")
	if err != nil {
		t.Fatal(err)
	}
	// TODO: really want multiple verbosity levels ...
	t.Log("stdout <<", stdout, ">>")
	t.Log("stderr <<", stderr, ">>")
	t.Log("exitCode", exitCode)
	if exitCode != 0 {
		t.Fatal()
	}
}

func TestBar(t *testing.T) {
	stdout, stderr, exitCode, err := support.RunMillerCommand(os.Args[0], "cxt testdata/abixy")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("stdout <<", stdout, ">>")
	t.Log("stderr <<", stderr, ">>")
	t.Log("exitCode", exitCode)
	if exitCode != 1 {
		t.Fatal()
	}
}
