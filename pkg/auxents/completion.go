package auxents

import "fmt"

func genCompletion(args []string) int {
	verb := args[1]
	args = args[2:]

	printUsage := func() {
		fmt.Printf("Usage: mlr %s SHELL\n", verb)
		fmt.Println("Supported shells: bash")
		fmt.Println()
		fmt.Println("Add below to your bashrc to enable completion")
		fmt.Println("source <(mlr completion bash)")
	}

	if len(args) != 1 {
		printUsage()
		return 1
	}

	if args[0] == "-h" || args[0] == "--help" {
		printUsage()
		return 0
	}

	if args[0] != "bash" {
		fmt.Printlf("Unsupported shell: %s\n", args[0])
		printUsage()
		return 1
	}

	fmt.Println(`complete -o nospace -o nosort -C "mlr _complete_bash" mlr`)
	return 0
}
