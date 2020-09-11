package lib

import (
	"fmt"
	"os"
	"runtime"
)

// Lookalike for C's __FILE__ and __LINE__ printing.

func InternalCodingErrorIf(condition bool) {
	if !condition {
		return
	}
	_, fileName, fileLine, ok := runtime.Caller(1)
	if ok {
		fmt.Fprintf(
			os.Stderr,
			"Internal coding error detected at file %s line %d\n",
			fileName,
			fileLine,
		)
	} else {
		fmt.Fprintf(
			os.Stderr,
			"Internal coding error detected at file %s line %s\n",
			"(unknown)",
			"(unknown)",
		)
	}
	// Uncomment this and re-run if you want to get a stack trace to get the
	// call-tree that led to the indicated file/line:
	panic("eek")
	os.Exit(1)
}

func InternalCodingErrorPanic(message string) {
	_, fileName, fileLine, ok := runtime.Caller(1)
	if ok {
		panic(
			fmt.Sprintf(
				"Internal coding error detected at file %s line %d: %s\n",
				fileName,
				fileLine,
				message,
			),
		)
	} else {
		panic(
			fmt.Sprintf(
				"Internal coding error detected at file %s line %s: %s\n",
				"(unknown)",
				"(unknown)",
				message,
			),
		)
	}
	os.Exit(1)
}
