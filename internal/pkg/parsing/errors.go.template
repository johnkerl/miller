// ================================================================
// During the Miller build, after GOCC codegen and Go compile, this is copied
// over the top of GOCC codegen so that we can customize handling of error
// messages.
//
// Source:       src/parsing/errors.go.template
// Destination:
// src/parsing/errors/errors.go
// ================================================================

package errors

import (
	"fmt"
	"strings"

	"mlr/src/parsing/token"
)

type ErrorSymbol interface {
}

type Error struct {
	Err            error
	ErrorToken     *token.Token
	ErrorSymbols   []ErrorSymbol
	ExpectedTokens []string
	StackTop       int
}

func (e *Error) String() string {
	w := new(strings.Builder)
	fmt.Fprintf(w, "Error")
	if e.Err != nil {
		fmt.Fprintf(w, " %s\n", e.Err)
	} else {
		fmt.Fprintf(w, "\n")
	}
	fmt.Fprintf(w, "Token: type=%d, lit=%s\n", e.ErrorToken.Type, e.ErrorToken.Lit)
	fmt.Fprintf(w, "Pos: offset=%d, line=%d, column=%d\n", e.ErrorToken.Pos.Offset, e.ErrorToken.Pos.Line, e.ErrorToken.Pos.Column)
	fmt.Fprintf(w, "Expected one of: ")
	for _, sym := range e.ExpectedTokens {
		fmt.Fprintf(w, "%s ", sym)
	}
	fmt.Fprintf(w, "ErrorSymbol:\n")
	for _, sym := range e.ErrorSymbols {
		fmt.Fprintf(w, "%v\n", sym)
	}
	return w.String()
}

func (e *Error) Error() string {
	w := new(strings.Builder)
	fmt.Fprintf(
		w,
		"Parse error on token \"%s\" at line %d columnn %d.\n",
		string(e.ErrorToken.Lit),
		e.ErrorToken.Pos.Line,
		e.ErrorToken.Pos.Column,
	)
	if e.Err != nil {
		fmt.Fprintf(w, "%+v\n", e.Err)
	} else {
		suggestSemicolons := false
		for _, expected := range e.ExpectedTokens {
			if expected == ";" {
				suggestSemicolons = true
				break
			}
		}

		if suggestSemicolons {
			fmt.Fprintf(w, "Please check for missing semicolon.\n")
		}

		fmt.Fprintf(w, "Expected one of:\n")

		//for _, expected := range e.ExpectedTokens {
		//  fmt.Fprintf(w, "%s ", expected)
		//}

		// Print a carriage return every so often, in case there are many possible
		// next tokens.
		line := ""
		for _, expected := range e.ExpectedTokens {
			if line != "" {
				line += " "
			}
			line += expected
			if len(line) > 70 {
				fmt.Fprintf(w, "  %s\n", line)
				line = ""
			}
		}
		if line != "" {
			fmt.Fprintf(w, "  %s\n", line)
		}

	}
	return w.String()
}
