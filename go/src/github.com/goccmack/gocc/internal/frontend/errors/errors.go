package errors

import "github.com/goccmack/gocc/internal/frontend/token"

type ErrorSymbol interface {
}

type Error struct {
	Err            error
	ErrorToken     *token.Token
	ErrorPos       token.Position
	ErrorSymbols   []ErrorSymbol
	ExpectedTokens []string
}

func (E *Error) String() string {
	errmsg := "Got " + E.ErrorToken.String() + " @ " + E.ErrorPos.String()
	if E.Err != nil {
		errmsg += " " + E.Err.Error()
	} else {
		errmsg += ", expected one of: "
		for _, t := range E.ExpectedTokens {
			errmsg += t + " "
		}
	}
	return errmsg
}
