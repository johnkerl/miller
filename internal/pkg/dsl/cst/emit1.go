// ================================================================
// The other emit variants (emit, emitp, emitf) need to take only oosvars, etc.
// -- not arbitrary expressions which *evaluate* to map. Emit1, by contrast,
// takes any expression which evaluates to a map. So you can do 'emit1
// mapsum({"id": $id}, $some_map_valued_field})'.
//
// The reason for this is LR1 shift-reduce conflicts. When I originally
// implemented emit/emitp/emitf, I permitted a lot of options for lashing
// together multiple oosvars, indexing, redirection, etc. When we try to let emit (not
// emit1) take arbitrary Rvalue as argument, we get LR1 conflicts since the
// parse can't disambiguate between all the possibilities for commas and
// parentheses for emit-lashing and emit-indexing, and all the possibilities
// for commas and parentheses for the Rvalue expression itself.
//
// So, we have emit/emitp which permit grammatical complexity in the
// lashing/indexing/redirection, and emit1 which permits grammatical complexity
// in the emittable.
// ================================================================

package cst

import (
	"fmt"

	"mlr/internal/pkg/dsl"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/runtime"
	"mlr/internal/pkg/types"
)

type Emit1StatementNode struct {
	evaluable IEvaluable
}

func (root *RootNode) BuildEmit1StatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeEmit1Statement)
	return root.buildEmit1StatementNode(astNode, false)
}

func (root *RootNode) buildEmit1StatementNode(
	astNode *dsl.ASTNode,
	isEmitP bool,
) (IExecutable, error) {
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)

	evaluable, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}

	return &Emit1StatementNode{
		evaluable: evaluable,
	}, nil
}

func (node *Emit1StatementNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	value := node.evaluable.Evaluate(state)
	if value.IsAbsent() {
		return nil, nil
	}

	// TODO: in a to-be-developed strict mode, fatal here.
	valueAsMap := value.GetMap() // nil if not a map
	if valueAsMap == nil {
		return nil, nil
	}

	if state.OutputChannel != nil {
		state.OutputChannel <- types.NewRecordAndContext(valueAsMap, state.Context)
	} else {
		fmt.Println(valueAsMap.String())
	}

	return nil, nil
}
