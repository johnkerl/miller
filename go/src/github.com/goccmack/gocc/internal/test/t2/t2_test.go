package t2

import (
	"os/exec"
	"strings"
	"testing"
)

func TestEmptyKeyword(t *testing.T) {
	var buf = new(strings.Builder)
	cmd := exec.Command("gocc", "t2.bnf")
	cmd.Stdout = buf
	err := cmd.Run()

	e, ok := err.(*exec.ExitError)
	if !ok || e.Success() {
		t.Fatalf("gocc t2.bnf should return with exit error, but get %v.", err)
	}
	expectedStderr := `empty production alternative: Maybe you are missing the "empty" keyword in "B : \t<<  >>"`
	actual := buf.String()
	if !strings.Contains(actual, expectedStderr) {
		t.Fatalf("%q should contains %q, but not", actual, expectedStderr)
	}
}
