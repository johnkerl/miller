// ================================================================
// Just playing around -- nothing serious here.
// ================================================================

package repl

import (
	"fmt"
	"os"

	"golang.org/x/term"

	"miller/version"
)

func getInputIsTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

func getPrompt1() string {
	prompt1 := os.Getenv("MLR_REPL_PS1")
	if prompt1 == "" {
		prompt1 = "[mlr] "
	}
	return prompt1
}

func getPrompt2() string {
	prompt2 := os.Getenv("MLR_REPL_PS2")
	if prompt2 == "" {
		prompt2 = "[mlr] "
	}
	return prompt2
}

func (this *Repl) printStartupBanner() {
	if this.inputIsTerminal {
		fmt.Printf("Miller %s\n", version.STRING) // TODO: inhibit if mlr repl -q
		fmt.Printf("Type ':help' for on-line help.\n")
	}
}

func (this *Repl) printPrompt1() {
	if this.inputIsTerminal {
		fmt.Print(this.prompt1)
	}
}

func (this *Repl) printPrompt2() {
	if this.inputIsTerminal {
		fmt.Print(this.prompt2)
	}
}
