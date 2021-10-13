package lib

import (
	"fmt"
	"os"
	"path"
	"runtime"
)

// InternalCodingErrorIf is a lookalike for C's __FILE__ and __LINE__ printing,
// with exit 1 if the condition is true.
func InternalCodingErrorIf(condition bool) {
	if !condition {
		return
	}
	_, fileName, fileLine, ok := runtime.Caller(1)
	if ok {
		fmt.Fprintf(
			os.Stderr,
			"Internal coding error detected at file %s line %d\n",
			// Full path preferred but breaks diffs on regression-test actual vs expected
			// stderr comparison on expect-fail cases.
			path.Base(fileName),
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
	// panic("eek")
	os.Exit(1)
}

// InternalCodingErrorWithMessageIf is a lookalike for C's __FILE__ and
// __LINE__ printing, with exit 1 if the condition is true.
func InternalCodingErrorWithMessageIf(condition bool, message string) {
	if !condition {
		return
	}
	_, fileName, fileLine, ok := runtime.Caller(1)
	if ok {
		fmt.Fprintf(
			os.Stderr,
			"Internal coding error detected at file %s line %d: %s\n",
			path.Base(fileName),
			fileLine,
			message,
		)
	} else {
		fmt.Fprintf(
			os.Stderr,
			"Internal coding error detected at file %s line %s: %s\n",
			"(unknown)",
			"(unknown)",
			message,
		)
	}
	// Uncomment this and re-run if you want to get a stack trace to get the
	// call-tree that led to the indicated file/line:
	// panic("eek")
	os.Exit(1)
}

// InternalCodingErrorPanic is like InternalCodingErrorIf, expect that it
// panics the process (for stack trace, which is usually not desired), and that
// it requires the if-test to be at the caller.
func InternalCodingErrorPanic(message string) {
	_, fileName, fileLine, ok := runtime.Caller(1)
	if ok {
		panic(
			fmt.Sprintf(
				"Internal coding error detected at file %s line %d: %s\n",
				path.Base(fileName),
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

// WhereAreWe shows a stack trace from the current callsite.
func WhereAreWe() {
	// Start at 1, not 0, since this function itself is not of interest.
	for i := 1; i < 20; i++ {
		_, fileName, fileLine, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fmt.Printf("  %s %d\n", fileName, fileLine)
	}
}
