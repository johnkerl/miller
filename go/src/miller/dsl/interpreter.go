package dsl

import (
	"errors"
	"strconv"

	"miller/containers"
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
		value, err := this.evaluateNode(rhsNode, inrec, context)
		// xxx temp undefined-handling
		if err != nil {
			return nil, err
		}
		inrec.Put(&fieldName, &value)
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

func (this *Interpreter) evaluateNode(
	node *ASTNode,
	inrec *containers.Lrec,
	context *runtime.Context,
) (string, error) {
	var sval = ""
	if node.Token != nil {
		sval = string(node.Token.Lit)
	}

	switch node.NodeType {

	case NodeTypeStringLiteral:
		return sval[1 : len(sval)-1], nil // xxx temp -- fix this in the grammar or ast-insert?
	case NodeTypeNumberLiteral:
		return sval, nil
	case NodeTypeBooleanLiteral:
		return sval, nil

	case NodeTypeDirectFieldName:
		fieldName := sval[1:] // xxx temp -- fix this in the grammar or ast-insert?
		fieldValue := inrec.Get(&fieldName)
		if fieldValue == nil {
			return "", errors.New("unhandled")
		} else {
			return *fieldValue, nil
		}
		break
	case NodeTypeIndirectFieldName:
		return "", errors.New("unhandled")
		break

	case NodeTypeStatementBlock:
		return "", errors.New("unhandled")
		break
	case NodeTypeAssignment:
		return "", errors.New("unhandled")
		break
	case NodeTypeOperator:
		this.checkArity(node, 2) // xxx temp -- binary-only for now
		return "", errors.New("unhandled")
		break
	case NodeTypeContextVariable:
		return this.evaluateContextVariableNode(node, context)
		break

	}
	return "", errors.New("unhandled")
}

func (this *Interpreter) evaluateContextVariableNode(
	node *ASTNode,
	context *runtime.Context,
) (string, error) {
	if node.Token == nil {
		return "", errors.New("internal coding error") // xxx libify
	}
	sval := string(node.Token.Lit)
	switch sval {
	case "FILENAME":
		return context.FILENAME, nil
		break
	case "FILENUM":
		return strconv.FormatInt(context.FILENUM, 10), nil
		break
	case "NF":
		return strconv.FormatInt(context.NF, 10), nil
		break
	case "NR":
		return strconv.FormatInt(context.NR, 10), nil
		break
	case "FNR":
		return strconv.FormatInt(context.FNR, 10), nil
		break
	}
	return "", errors.New("internal coding error") // xxx libify
}
