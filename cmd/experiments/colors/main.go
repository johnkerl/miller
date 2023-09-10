// This is the entry point for the mlr executable.
package main

import (
	"fmt"
	"github.com/johnkerl/miller/pkg/colorizer"
)

const boldString = "\u001b[1m"
const underlineString = "\u001b[4m"
const reversedString = "\u001b[7m"
const redString = "\u001b[1;31m"
const blueString = "\u001b[1;34m"
const defaultString = "\u001b[0m"

func main() {
	fmt.Printf("Hello, world!\n")

	fmt.Printf("1. before %s during %s after\n", "", "")
	fmt.Printf("4. before %s during %s after\n", boldString, defaultString)
	fmt.Printf("2. before %s during %s after\n", redString, defaultString)
	fmt.Printf("3. before %s during %s after\n", blueString, defaultString)
	fmt.Printf("5. before %s during %s after\n", redString+boldString, defaultString)
	fmt.Printf("7. before %s during %s after\n", blueString, defaultString)
	fmt.Printf("8. before %s during %s after\n", boldString, defaultString)
	fmt.Printf("9. before %s during %s after\n", boldString, defaultString)
	fmt.Println()

	names := []string{
		"plain",
		"red",
		"bold",
		"bold-red",
		"red-bold",
		"blue-underline",
		"208",
		"reversed-208",
	}
	for _, name := range names {
		colorizer.SetKeyColor(name)
		fmt.Printf("testing: [%s]\n", colorizer.MaybeColorizeKey(name, true))
	}
}
