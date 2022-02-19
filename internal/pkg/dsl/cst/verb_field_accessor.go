// ================================================================
// This is an experimental technique for doing things like
//
//   mlr cut -f 'item,[[[3]]],item.name,foo["bar"]
//
// where the expressions aren't simple strings but rather correspond to DSL expressions like
//
//   $item  $[[[3]]]  $item.name  $foo["bar"]
//
// ================================================================

package cst

import (
	"fmt"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/dsl"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/runtime"
)

//type State struct {
//	Inrec                    *mlrval.Mlrmap
//	Context                  *types.Context
//	Oosvars                  *mlrval.Mlrmap
//	FilterExpression         *mlrval.Mlrval
//	Stack                    *Stack
//	OutputRecordsAndContexts *list.List // list of *types.RecordAndContext
//
//	// For holding "\0".."\9" between where they are set via things like
//	// '$x =~ "(..)_(...)"', and interpolated via things like '$y = "\2:\1"'.
//	RegexCaptures []string
//	Options       *cli.TOptions
//}

type VerbFieldAccessor struct {
	cstRootNode  *RootNode
	runtimeState *runtime.State
}

// NodeTypeDirectFieldValue
// NodeTypeIndirectFieldValue
// (BracedFieldValue is DirectFieldValue)
// Indexed with .
// Indexed with []
// NodeTypePositionalFieldName
// NodeTypePositionalFieldValue

// mlr -n put -v '$item; ${item}; $["item"]; $item.name; $item["name"]; $[[3]]; $[[[3]]]'
// AST:
// * statement block
//     * bare boolean
//         * direct field value "item"
//     * bare boolean
//         * direct field value "item"
//     * bare boolean
//         * indirect field value "$[]"
//             * string literal "item"
//     * bare boolean
//         * dot operator "."
//             * direct field value "item"
//             * local variable "name"
//     * bare boolean
//         * array or map index access "[]"
//             * direct field value "item"
//             * string literal "name"
//     * bare boolean
//         * positionally-indexed field name "$[]"
//             * int literal "3"
//     * bare boolean
//         * positionally-indexed field value "$[]"
//             * int literal "3"

func verbFieldAccessorASTValidator(dslString string, astNode *dsl.AST) error {
	// TODO: flesh this out
	err := fmt.Errorf("malformed field-selector syntax: \"%s\"", dslString)

	if astNode.RootNode.Type != dsl.NodeTypeStatementBlock {
		return err
	}
	if len(astNode.RootNode.Children) != 1 {
		return err
	}
	if astNode.RootNode.Children[0].Type != dsl.NodeTypeBareBoolean {
		return err
	}
	if len(astNode.RootNode.Children[0].Children) != 1 {
		return err
	}

	return nil
}

func NewVerbFieldAccessor(input string) (*VerbFieldAccessor, error) {
	cstRootNode := NewEmptyRoot(nil, DSLInstanceTypeVerbFieldAccessor)
	err := cstRootNode.Build(
		[]string{"$" + input},          // dslStrings []string
		DSLInstanceTypeVerbFieldAccessor, // dslInstanceType DSLInstanceType
		false,                          // isReplImmediate bool
		false,                          // doWarnings bool
		false,                          // warningsAreFatal bool
		verbFieldAccessorASTValidator,    // astBuildVisitorFunc ASTBuildVisitorFunc
	)
	if err != nil {
		return nil, err
	}

	options := cli.DefaultOptions()
	runtimeState := runtime.NewEmptyState(options)

	return &VerbFieldAccessor{
		cstRootNode,
		runtimeState,
	}, nil
}

func (g *VerbFieldAccessor) Get(record *mlrval.Mlrmap) *mlrval.Mlrval {
	// TODO: rework all the CST stuff to not have so much extra.
	// This is just a POC for now.
	g.runtimeState.Inrec = record
	node := g.cstRootNode.mainBlock.executables[0].(*BareBooleanStatementNode)
	return node.Evaluate(g.runtimeState)
}
