package dsl

import (
	"fmt"

	"github.com/johnkerl/miller/pkg/parsing/token"
)

// TokenToLocationInfo is used to track runtime errors back to source-code locations in DSL
// expressions, so we can have more informative error messages.
func TokenToLocationInfo(sourceToken *token.Token) string {
	if sourceToken == nil {
		return ""
	} else {
		return fmt.Sprintf(" at DSL expression line %d column %d", sourceToken.Pos.Line, sourceToken.Pos.Column)
	}
}
