// ================================================================
// FOO STUB
// ================================================================

package regression

import (
	"os"
	"testing"
)

func TestWorkingDirectory(t *testing.T) {
         wd, _ := os.Getwd()
         t.Log(wd)
}
