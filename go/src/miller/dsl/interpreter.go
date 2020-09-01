package dsl

import (
	"errors"

	"miller/containers"
	"miller/lib"
)

// ----------------------------------------------------------------
// Just a very temporary CST-free, AST-only interpreter to get me executing
// some DSL code with a minimum of keystroking, while I work out other issues
// including mlrval-valued lrecs, and port of mvfuncs from C to Go.
type Interpreter struct {
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

// ----------------------------------------------------------------
func (this *Interpreter) InterpretOnInputRecord(
	inrec *containers.Lrec,
	context *containers.Context,
	ast *AST,
) (outrec *containers.Lrec, err error) {
	root := ast.Root
	if root == nil {
		return nil, errors.New("internal coding error") // xxx libify
	}
	if root.NodeType != NodeTypeStatementBlock {
		return nil, errors.New("non-statement-block root node unhandled")
	}
	children := root.Children
	if children == nil {
		return nil, errors.New("internal coding error") // xxx libify
	}

	// For this very early stub, only process statement nodes (which is all the
	// grammar produces at this point ...)
	for _, child := range children {
		if child.NodeType != NodeTypeAssignment {
			return nil, errors.New("non-assignment node unhandled")
		}

		err := this.checkArity(child, 2)
		if err != nil {
			return nil, err
		}
		lhsNode := child.Children[0]
		rhsNode := child.Children[1]

		fieldName := string(lhsNode.Token.Lit)[1:] // strip off leading '$'
		mvalue, err := this.evaluateNode(rhsNode, inrec, context)
		if err != nil {
			return nil, err
		} else {
			if !mvalue.IsAbsent() {
				inrec.Put(&fieldName, &mvalue)
			}
		}
	}

	return inrec, nil
}

// ----------------------------------------------------------------
// xxx make into ASTNode method?
func (this *Interpreter) checkArity(
	node *ASTNode,
	arity int,
) error {
	if len(node.Children) != arity {
		return errors.New("internal coding error") // xxx libify
	} else {
		return nil
	}
}

// ----------------------------------------------------------------
func (this *Interpreter) evaluateNode(
	node *ASTNode,
	inrec *containers.Lrec,
	context *containers.Context,
) (lib.Mlrval, error) {
	var sval = ""
	if node.Token != nil {
		sval = string(node.Token.Lit)
	}

	switch node.NodeType {

	case NodeTypeStringLiteral:
		// xxx temp "..." strip -- fix this in the grammar or ast-insert
		return lib.MlrvalFromString(sval[1 : len(sval)-1]), nil
	case NodeTypeIntLiteral:
		return lib.MlrvalFromInt64String(sval), nil
	case NodeTypeFloatLiteral:
		return lib.MlrvalFromFloat64String(sval), nil
	case NodeTypeBoolLiteral:
		return lib.MlrvalFromBoolString(sval), nil

	case NodeTypeDirectFieldName:
		fieldName := sval[1:] // xxx temp -- fix this in the grammar or ast-insert?
		fieldValue := inrec.Get(&fieldName)
		if fieldValue == nil {
			return lib.MlrvalFromAbsent(), nil
		} else {
			return *fieldValue, nil
		}
		break
	case NodeTypeIndirectFieldName:
		return lib.MlrvalFromError(), errors.New("unhandled1")
		break

	case NodeTypeStatementBlock:
		return lib.MlrvalFromError(), errors.New("unhandled2")
		break
	case NodeTypeAssignment:
		return lib.MlrvalFromError(), errors.New("unhandled3")
		break
	case NodeTypeOperator:
		this.checkArity(node, 2) // xxx temp -- binary-only for now
		return this.evaluateBinaryOperatorNode(node, node.Children[0], node.Children[1],
			inrec, context)
		break
	case NodeTypeContextVariable:
		return this.evaluateContextVariableNode(node, context)
		break

	}
	return lib.MlrvalFromError(), errors.New("unhandled4")
}

func (this *Interpreter) evaluateContextVariableNode(
	node *ASTNode,
	context *containers.Context,
) (lib.Mlrval, error) {
	if node.Token == nil {
		return lib.MlrvalFromError(), errors.New("internal coding error") // xxx libify
	}
	sval := string(node.Token.Lit)
	switch sval {
	case "FILENAME":
		return lib.MlrvalFromString(context.FILENAME), nil
		break
	case "FILENUM":
		return lib.MlrvalFromInt64(context.FILENUM), nil
		break
	case "NF":
		return lib.MlrvalFromInt64(context.NF), nil
		break
	case "NR":
		return lib.MlrvalFromInt64(context.NR), nil
		break
	case "FNR":
		return lib.MlrvalFromInt64(context.FNR), nil
		break

	case "IPS":
		return lib.MlrvalFromString(context.IPS), nil
		break
	case "IFS":
		return lib.MlrvalFromString(context.IFS), nil
		break
	case "IRS":
		return lib.MlrvalFromString(context.IRS), nil
		break

	case "OPS":
		return lib.MlrvalFromString(context.OPS), nil
		break
	case "OFS":
		return lib.MlrvalFromString(context.OFS), nil
		break
	case "ORS":
		return lib.MlrvalFromString(context.ORS), nil
		break

		break
	}
	return lib.MlrvalFromError(), errors.New("internal coding error") // xxx libify
}

func (this *Interpreter) evaluateBinaryOperatorNode(
	node *ASTNode,
	leftChild *ASTNode,
	rightChild *ASTNode,
	inrec *containers.Lrec,
	context *containers.Context,
) (lib.Mlrval, error) {
	sop := string(node.Token.Lit)

	leftValue, leftErr := this.evaluateNode(leftChild, inrec, context)
	if leftErr != nil {
		return lib.MlrvalFromError(), leftErr
	}
	rightValue, rightErr := this.evaluateNode(rightChild, inrec, context)
	if rightErr != nil {
		return lib.MlrvalFromError(), rightErr
	}

	switch sop {
	case "+":
		return lib.MlrvalPlus(&leftValue, &rightValue), nil
		break
	case "-":
		return lib.MlrvalMinus(&leftValue, &rightValue), nil
		break
	case "*":
		return lib.MlrvalTimes(&leftValue, &rightValue), nil
		break
	case "/":
		return lib.MlrvalDivide(&leftValue, &rightValue), nil
		break
	case "//":
		return lib.MlrvalIntDivide(&leftValue, &rightValue), nil
		break
		// xxx continue ...
	}

	return lib.MlrvalFromError(), errors.New("internal coding error") // xxx libify
}
