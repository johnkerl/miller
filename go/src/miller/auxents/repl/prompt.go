// ================================================================
// Handling for default and customized banners/prompts for the Miller REPL.
// ================================================================

package repl

import (
	"fmt"
	"os"

	"golang.org/x/term"

	"miller/version"
)

const ENV_PRIMARY_PROMPT = "MLR_REPL_PS1"
const ENV_SECONDARY_PROMPT = "MLR_REPL_PS2"
const DEFAULT_PRIMARY_PROMPT = "[mlr] "
const DEFAULT_SECONDARY_PROMPT = ""


func getInputIsTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

func getPrompt1() string {
	prompt1 := os.Getenv(ENV_PRIMARY_PROMPT)
	if prompt1 == "" {
		prompt1 = DEFAULT_PRIMARY_PROMPT
	}
	return prompt1
}

func getPrompt2() string {
	prompt2 := os.Getenv(ENV_SECONDARY_PROMPT)
	if prompt2 == "" {
		prompt2 = DEFAULT_SECONDARY_PROMPT
	}
	return prompt2
}

func (this *Repl) printStartupBanner() {
	if this.inputIsTerminal {
		fmt.Printf("Miller %s\n", version.STRING) // TODO: inhibit if mlr repl -q
		fmt.Printf("Type ':help' for on-line help; ':quit' to quit.\n")
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
