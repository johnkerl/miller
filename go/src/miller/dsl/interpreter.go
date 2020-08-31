package dsl

import (
	"errors"
	"strconv"

	"miller/containers"
	"miller/lib"
	"miller/runtime"
)

// Just a very temporary CST-free, AST-only interpreter to get me executing
// some DSL code with a minimum of keystroking, while I work out other issues
// including mlrval-valued lrecs.
type Interpreter struct {
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (this *Interpreter) InterpretOnInputRecord(
	inrec *containers.Lrec,
	context *runtime.Context,
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
		value, defined, err := this.evaluateNode(rhsNode, inrec, context)
		if err != nil {
			return nil, err
		}
		if defined {
			inrec.Put(&fieldName, &value)
		}
	}

	return inrec, nil
}

// xxx make into ASTNode method
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

// xxx needs null/undefined/string/error. then, string->mlrval.
func (this *Interpreter) evaluateNode(
	node *ASTNode,
	inrec *containers.Lrec,
	context *runtime.Context,
) (string, bool, error) {
	var sval = ""
	if node.Token != nil {
		sval = string(node.Token.Lit)
	}

	switch node.NodeType {

	case NodeTypeStringLiteral:
		// xxx temp -- fix this in the grammar or ast-insert?
		return sval[1 : len(sval)-1], true, nil
	case NodeTypeNumberLiteral:
		return sval, true, nil // xxx temp -- to mlrval
	case NodeTypeBooleanLiteral:
		return sval, true, nil // xxx temp -- to mlrval

	case NodeTypeDirectFieldName:
		fieldName := sval[1:] // xxx temp -- fix this in the grammar or ast-insert?
		fieldValue := inrec.Get(&fieldName)
		if fieldValue == nil {
			return "", false, nil
		} else {
			return *fieldValue, true, nil
		}
		break
	case NodeTypeIndirectFieldName:
		return "", true, errors.New("unhandled")
		break

	case NodeTypeStatementBlock:
		return "", true, errors.New("unhandled")
		break
	case NodeTypeAssignment:
		return "", true, errors.New("unhandled")
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
	return "", true, errors.New("unhandled")
}

func (this *Interpreter) evaluateContextVariableNode(
	node *ASTNode,
	context *runtime.Context,
) (string, bool, error) {
	if node.Token == nil {
		return "", true, errors.New("internal coding error") // xxx libify
	}
	sval := string(node.Token.Lit)
	switch sval {
	case "FILENAME":
		return context.FILENAME, true, nil
		break
	case "FILENUM":
		return strconv.FormatInt(context.FILENUM, 10), true, nil
		break
	case "NF":
		return strconv.FormatInt(context.NF, 10), true, nil
		break
	case "NR":
		return strconv.FormatInt(context.NR, 10), true, nil
		break
	case "FNR":
		return strconv.FormatInt(context.FNR, 10), true, nil
		break

	case "IPS":
		return context.IPS, true, nil
		break
	case "IFS":
		return context.IFS, true, nil
		break
	case "IRS":
		return context.IRS, true, nil
		break
	case "OPS":
		return context.OPS, true, nil
		break
	case "OFS":
		return context.OFS, true, nil
		break
	case "ORS":
		return context.ORS, true, nil
		break
	}
	return "", true, errors.New("internal coding error") // xxx libify
}

func (this *Interpreter) evaluateBinaryOperatorNode(
	node *ASTNode,
	leftChild *ASTNode,
	rightChild *ASTNode,
	inrec *containers.Lrec,
	context *runtime.Context,
) (string, bool, error) {
	sval := string(node.Token.Lit)

	leftValue, leftDefined, leftErr := this.evaluateNode(leftChild, inrec, context)
	if leftErr != nil {
		return "", true, leftErr
	}
	if !leftDefined {
		return "", false, nil
	}
	rightValue, rightDefined, rightErr := this.evaluateNode(rightChild, inrec, context)
	if rightErr != nil {
		return "", true, rightErr
	}
	if !rightDefined {
		return "", false, nil
	}

	switch sval {
	case ".":
		return leftValue + rightValue, true, nil
		break
	}

	switch sval {
	case "+":
		// xxx make a lib method -- Itoa64
		//return lib.Itoa64(leftInt + rightInt), true, nil
		a := lib.MlrvalFromInt64String(leftValue)
		b := lib.MlrvalFromInt64String(rightValue)
		c := lib.MlrvalPlus(&a, &b)
		return c.String(), true, nil
		break
	}

	// make a helper method for int-pairings
	leftInt, lerr := strconv.ParseInt(leftValue, 10, 0)
	if lerr != nil {
		// to do: consider error-propagation through the AST evaluator, with
		// null/undefined/error in the binop matrices etc.
		//
		// need to separate internal coding errors, from data-dependent ones
		//
		//return "", true, lerr
		return "(error)", true, nil
	}
	rightInt, rerr := strconv.ParseInt(rightValue, 10, 0)
	if rerr != nil {
		//return "", true, rerr
		return "(error)", true, nil
	}

	switch sval {
	case "-":
		return lib.Itoa64(leftInt - rightInt), true, nil
		break
	case "*":
		return lib.Itoa64(leftInt * rightInt), true, nil
		break
	case "/":
		return lib.Itoa64(leftInt / rightInt), true, nil
		break
	case "^":
		return lib.Itoa64(leftInt ^ rightInt), true, nil
		break
	case "&":
		return lib.Itoa64(leftInt & rightInt), true, nil
		break
	case "|":
		return lib.Itoa64(leftInt | rightInt), true, nil
		break
	case "//":
		return "", true, errors.New("unhandled")
		break
	}

	return "", true, errors.New("internal coding error") // xxx libify
}
