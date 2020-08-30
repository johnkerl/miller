package dsl

import (
	"miller/containers"
	"miller/runtime"
)

// Just a very temporary CST-free, AST-only interpreter to get me executing
// some DSL code with a minimum of keystroking, while I work out other issues
// including mlrval-valued lrecs.

type Interpreter struct {
}

func NewInterpreter() *Interpreter {
	return &Interpreter {
	}
}

func (this *Interpreter) InterpretOnInputRecord(
	inrec *containers.Lrec,
	context* runtime.Context,
) (outrec *containers.Lrec) {
	k := "foo"
	v := "bar"
	inrec.Put(&k, &v)
	return inrec
}
