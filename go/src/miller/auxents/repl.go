// ================================================================
// Just playing around -- nothing serious here.
// ================================================================

package auxents

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func replUsage(verbName string, o *os.File, exitCode int) {
	fmt.Fprintf(o, "Usage: %s %s with no arguments\n", mlrExeName(), verbName)
	os.Exit(exitCode)
}

func replMain(args []string) int {
	repl := NewRepl()
	repl.Handle(os.Stdin, os.Stdout)
	return 0
}

// ================================================================
type Repl struct {
}

func NewRepl() *Repl {
	return &Repl{
	}
}

func (this *Repl) Handle(istream *os.File, ostream *os.File) {
	lineReader := bufio.NewReader(istream)

	for {
		line, err := lineReader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			// TODO: lib.MlrExeName()
			fmt.Fprintln(os.Stderr, "mlr repl:", err)
			os.Exit(1)
		}

		// This is how to do a chomp:
		line = strings.TrimRight(line, "\n")

		fmt.Println(line)
	}
}
