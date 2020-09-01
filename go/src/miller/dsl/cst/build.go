package cst

import (
	"errors"

	"miller/dsl"
)

func Build(ast *dsl.AST) (*Root, error) {
	if ast.Root == nil {
		return nil, errors.New("Cannot build CST from nil AST root")
	}
	if ast.Root.NodeType != dsl.NodeTypeStatementBlock {
		return nil, errors.New("Non-statement-block AST root node unhandled")
	}
	astChildren := ast.Root.Children

	// For this very early stub, only process statement nodes (which is all the
	// grammar produces at this point ...)
	for _, astChild := range astChildren {
		if astChild.NodeType != dsl.NodeTypeAssignment {
			return nil, errors.New("Non-assignment AST node unhandled")
		}

		statement, err := newSrecAssignment(astChild)
		if err != nil {
			return nil, err
		}

		// have a helper method for DirectSrecFieldAssignment from AST NodeTypeAssignment node
	}

	return Root, nil
}

// ----------------------------------------------------------------
// xxx make into an ASTNode method?
func checkASTNodeArity(
	node *ASTNode,
	arity int,
) error {
	if len(node.Children) != arity {
		return errors.New("AST node arity malformed") // xxx libify
	} else {
		return nil
	}
}

// ----------------------------------------------------------------
func newSrecAssignment(astNode *ASTNode) (*DirectSrecFieldAssignment, error) {

	err := this.checkArity(astChild, 2)
	if err != nil {
		return nil, err
	}

	lhsASTNode := astChild.Children[0]
	rhsASTNode := astChild.Children[1]

	// strip off leading '$'. TODO: move into the AST-builder
	lhsFieldName := string(lhsNode.Token.Lit)[1:]
	rhs = nil // xxx newEvaluable(rhsASTNode)

	//		mvalue, err := this.evaluateNode(rhsNode, inrec, context)
	//		if err != nil {
	//			return nil, err
	//		} else {
	//			if !mvalue.IsAbsent() {
	//				inrec.Put(&fieldName, &mvalue)
	//			}
	//		}

	return NewDirectSrecFieldAssignment(lhsFieldName, rhs)
}

// ----------------------------------------------------------------
//func (this *Interpreter) evaluateNode(
//	node *ASTNode,
//) (lib.Mlrval, error) {
//	var sval = ""
//	if node.Token != nil {
//		sval = string(node.Token.Lit)
//	}
//
//	switch node.NodeType {
//
//	case NodeTypeStringLiteral:
//		// xxx temp "..." strip -- fix this in the grammar or ast-insert
//		return lib.MlrvalFromString(sval[1 : len(sval)-1]), nil
//	case NodeTypeIntLiteral:
//		return lib.MlrvalFromInt64String(sval), nil
//	case NodeTypeFloatLiteral:
//		return lib.MlrvalFromFloat64String(sval), nil
//	case NodeTypeBoolLiteral:
//		return lib.MlrvalFromBoolString(sval), nil
//
//	case NodeTypeDirectFieldName:
//		fieldName := sval[1:] // xxx temp -- fix this in the grammar or ast-insert?
//		fieldValue := inrec.Get(&fieldName)
//		if fieldValue == nil {
//			return lib.MlrvalFromAbsent(), nil
//		} else {
//			return *fieldValue, nil
//		}
//		break
//	case NodeTypeIndirectFieldName:
//		return lib.MlrvalFromError(), errors.New("unhandled1")
//		break
//
//	case NodeTypeStatementBlock:
//		return lib.MlrvalFromError(), errors.New("unhandled2")
//		break
//	case NodeTypeAssignment:
//		return lib.MlrvalFromError(), errors.New("unhandled3")
//		break
//	case NodeTypeOperator:
//		this.checkArity(node, 2) // xxx temp -- binary-only for now
//		return this.evaluateBinaryOperatorNode(node, node.Children[0], node.Children[1],
//			inrec, context)
//		break
//	case NodeTypeContextVariable:
//		return this.evaluateContextVariableNode(node, context)
//		break
//
//	}
//	return lib.MlrvalFromError(), errors.New("unhandled4")
//}

//func (this *Interpreter) evaluateContextVariableNode(
//	node *ASTNode,
//	context *containers.Context,
//) (lib.Mlrval, error) {
//	if node.Token == nil {
//		return lib.MlrvalFromError(), errors.New("internal coding error") // xxx libify
//	}
//	sval := string(node.Token.Lit)
//	switch sval {
//	case "FILENAME":
//		return lib.MlrvalFromString(context.FILENAME), nil
//		break
//	case "FILENUM":
//		return lib.MlrvalFromInt64(context.FILENUM), nil
//		break
//	case "NF":
//		return lib.MlrvalFromInt64(context.NF), nil
//		break
//	case "NR":
//		return lib.MlrvalFromInt64(context.NR), nil
//		break
//	case "FNR":
//		return lib.MlrvalFromInt64(context.FNR), nil
//		break
//
//	case "IPS":
//		return lib.MlrvalFromString(context.IPS), nil
//		break
//	case "IFS":
//		return lib.MlrvalFromString(context.IFS), nil
//		break
//	case "IRS":
//		return lib.MlrvalFromString(context.IRS), nil
//		break
//
//	case "OPS":
//		return lib.MlrvalFromString(context.OPS), nil
//		break
//	case "OFS":
//		return lib.MlrvalFromString(context.OFS), nil
//		break
//	case "ORS":
//		return lib.MlrvalFromString(context.ORS), nil
//		break
//
//		break
//	}
//	return lib.MlrvalFromError(), errors.New("internal coding error") // xxx libify
//}

//func (this *Interpreter) evaluateBinaryOperatorNode(
//	node *ASTNode,
//	leftChild *ASTNode,
//	rightChild *ASTNode,
//	inrec *containers.Lrec,
//	context *containers.Context,
//) (lib.Mlrval, error) {
//	sop := string(node.Token.Lit)
//
//	leftValue, leftErr := this.evaluateNode(leftChild, inrec, context)
//	if leftErr != nil {
//		return lib.MlrvalFromError(), leftErr
//	}
//	rightValue, rightErr := this.evaluateNode(rightChild, inrec, context)
//	if rightErr != nil {
//		return lib.MlrvalFromError(), rightErr
//	}
//
//	switch sop {
//	case ".":
//		return lib.MlrvalDot(&leftValue, &rightValue), nil
//		break
//	case "+":
//		return lib.MlrvalPlus(&leftValue, &rightValue), nil
//		break
//	case "-":
//		return lib.MlrvalMinus(&leftValue, &rightValue), nil
//		break
//	case "*":
//		return lib.MlrvalTimes(&leftValue, &rightValue), nil
//		break
//	case "/":
//		return lib.MlrvalDivide(&leftValue, &rightValue), nil
//		break
//	case "//":
//		return lib.MlrvalIntDivide(&leftValue, &rightValue), nil
//		break
//		// xxx continue ...
//	}
//
//	return lib.MlrvalFromError(), errors.New("internal coding error") // xxx libify
//}
