package parser

import (
	"errors"
	"strconv"
)

import errs "github.com/goccmack/gocc/internal/frontend/errors"

import "github.com/goccmack/gocc/internal/frontend/token"

// ParserUTab

type ParserUTab struct {
	actTab     *ActionTabU
	canRecover *CanRecover
	gotoTab    *GotoTabU
	prodTab    *ProdTabU
	stack      *stack
	nextToken  *token.Token
	pos        token.Position
	tokenMap   *token.TokenMap
}

func NewParserUTab(tm *token.TokenMap) *ParserUTab {
	p := &ParserUTab{
		actTab:     getActionTableUncompressed(),
		canRecover: getCanRecoverTableUncompressed(),
		gotoTab:    getGotoTableUncompressed(),
		prodTab:    getProductionsTableUncompressed(),
		stack:      NewStack(),
		tokenMap:   tm,
	}
	p.stack.Push(0, nil) //TODO: which attribute should be pushed here?
	return p
}

func (P *ParserUTab) Reset() {
	P.stack.reset()
	P.stack.Push(0, nil) //TODO: which attribute should be pushed here?
}

func (P *ParserUTab) Error(err error, scanner Scanner) (recovered bool, errorAttrib *errs.Error) {
	errorAttrib = &errs.Error{
		Err:            err,
		ErrorToken:     P.nextToken,
		ErrorPos:       P.pos,
		ErrorSymbols:   P.popNonRecoveryStates(),
		ExpectedTokens: make([]string, 0, 8),
	}
	for t, act := range P.actTab[P.stack.Top()] {
		if act != nil {
			errorAttrib.ExpectedTokens = append(errorAttrib.ExpectedTokens, P.tokenMap.TokenString(token.Type(t)))
		}
	}

	action := P.actTab[P.stack.Top()][P.tokenMap.Type("error")]
	if action == nil {
		return
	}
	P.stack.Push(State(action.(Shift)), errorAttrib) // action can only be shift

	for P.actTab[P.stack.Top()][P.nextToken.Type] == nil && P.nextToken.Type != token.EOF {
		P.nextToken, P.pos = scanner.Scan()
	}

	return
}

func (P *ParserUTab) popNonRecoveryStates() (removedAttribs []errs.ErrorSymbol) {
	if rs, ok := P.firstRecoveryState(); ok {
		errorSymbols := P.stack.PopN(int(P.stack.TopIndex() - rs))
		removedAttribs = make([]errs.ErrorSymbol, len(errorSymbols))
		for i, e := range errorSymbols {
			removedAttribs[i] = e
		}
	} else {
		removedAttribs = []errs.ErrorSymbol{}
	}
	return
}

// recoveryState points to the highest state on the stack, which can recover
func (P *ParserUTab) firstRecoveryState() (recoveryState int, canRecover bool) {
	recoveryState, canRecover = P.stack.TopIndex(), P.canRecover[P.stack.Top()]
	for recoveryState > 0 && !canRecover {
		recoveryState--
		canRecover = P.canRecover[P.stack.Peek(recoveryState)]
	}
	return
}

func (P *ParserUTab) newError(err error) error {
	errmsg := "Error: " + P.TokString(P.nextToken) + " @ " + P.pos.String()
	if err != nil {
		errmsg += " " + err.Error()
	} else {
		errmsg += ", expected one of: "
		actRow := P.actTab[P.stack.Top()]
		i := 0
		for t, act := range actRow {
			if act != nil {
				errmsg += P.tokenMap.TokenString(token.Type(t))
				if i < len(actRow)-1 {
					errmsg += " "
				}
				i++
			}
		}
	}
	return errors.New(errmsg)
}

func (P *ParserUTab) TokString(tok *token.Token) string {
	msg := P.tokenMap.TokenString(tok.Type) + "(" + strconv.Itoa(int(tok.Type)) + ")"
	msg += " " + string(tok.Lit)
	return msg
}

func (this *ParserUTab) Parse(scanner Scanner) (res interface{}, err error) {
	this.Reset()
	this.nextToken, this.pos = scanner.Scan()
	for acc := false; !acc; {
		action := this.actTab[this.stack.Top()][this.nextToken.Type]
		if action == nil {
			if recovered, errAttrib := this.Error(nil, scanner); !recovered {
				this.nextToken, this.pos = errAttrib.ErrorToken, errAttrib.ErrorPos
				return nil, this.newError(nil)
			}
			if action = this.actTab[this.stack.Top()][this.nextToken.Type]; action == nil {
				panic("Error recover led to invalid action")
			}
		}
		// fmt.Printf("S%d %s %s\n", this.stack.Top(), this.nextToken, action)
		switch act := action.(type) {
		case Accept:
			res = this.stack.PopN(1)[0]
			acc = true
		case Shift:
			this.stack.Push(State(act), this.nextToken)
			this.nextToken, this.pos = scanner.Scan()
		case Reduce:
			prod := this.prodTab[int(act)]
			attrib, err := prod.ReduceFunc(this.stack.PopN(prod.NumSymbols))
			if err != nil {
				return nil, this.newError(err)
			} else {
				this.stack.Push(this.gotoTab[this.stack.Top()][prod.HeadIndex], attrib)
			}
		default:
			panic("unknown action: " + action.String())
		}
	}
	return res, nil
}
