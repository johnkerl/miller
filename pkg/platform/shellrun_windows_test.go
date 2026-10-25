package platform

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Under Git Bash or MSYS2 the shell must be runnable by a native Windows
// program, for system() and exec.Command as well as for os.StartProcess
// (used for --prepipe and piped redirects), which needs an absolute path.
func TestGetShellRunArrayUnderMSYS(t *testing.T) {
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("sh is not on the PATH")
	}
	t.Setenv("MSYSTEM", "MINGW64")

	shellRunArray := GetShellRunArray("echo hello")
	output, err := exec.Command(shellRunArray[0], shellRunArray[1:]...).Output()
	if err != nil {
		t.Fatalf("running %q: %v", shellRunArray, err)
	}
	if strings.TrimSpace(string(output)) != "hello" {
		t.Fatalf("expected \"hello\", got %q", output)
	}
	if !filepath.IsAbs(shellRunArray[0]) {
		t.Fatalf("expected an absolute path to the shell, got %q", shellRunArray[0])
	}
}
