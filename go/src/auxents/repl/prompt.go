// ================================================================
// Handling for default and customized banners/prompts for the Miller REPL.
// ================================================================

package repl

import (
	"fmt"
	"os"
	"runtime"

	"golang.org/x/term"

	"miller/src/version"
)

const ENV_PRIMARY_PROMPT = "MLR_REPL_PS1"
const ENV_SECONDARY_PROMPT = "MLR_REPL_PS2"
const DEFAULT_PRIMARY_PROMPT = "[mlr] "
const DEFAULT_SECONDARY_PROMPT = ""

func getInputIsTerminal() bool {
	if runtime.GOOS == "windows" && os.Getenv("MSYSTEM") != "" {
		// Sadly, term.IsTerminal doesn't work inside MSYS2 but does work
		// outside MSYS2. So in that case we simply return true so that the mlr
		// repl has a prompt.
		return true
	} else {
		return term.IsTerminal(int(os.Stdin.Fd()))
	}
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
		// TODO: inhibit if mlr repl -q
		fmt.Printf("Miller %s for %s:%s:%s\n", version.STRING, runtime.GOOS, runtime.GOARCH, runtime.Version())
		fmt.Printf("Type ':help' for on-line help; ':quit' to quit.\n")
	}
}

func (this *Repl) printPrompt1() {
	if this.inputIsTerminal && this.showPrompts {
		fmt.Print(this.prompt1)
	}
}

func (this *Repl) printPrompt2() {
	if this.inputIsTerminal && this.showPrompts {
		fmt.Print(this.prompt2)
	}
}
