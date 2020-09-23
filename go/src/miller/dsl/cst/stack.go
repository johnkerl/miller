package cst

import (
	"container/list"
	"fmt"

	"miller/types"
)

// ================================================================
// Stack frames for begin/end/if/for/function blocks
// ================================================================

// ----------------------------------------------------------------
type Stack struct {
	stackFrames *list.List // list of *StackFrame
}

func NewStack() *Stack {
	return &Stack{
		stackFrames: list.New(),
	}
}

func (this *Stack) PushStackFrame() {
	this.stackFrames.PushFront(NewStackFrame())
}

func (this *Stack) PopStackFrame() {
	this.stackFrames.Remove(this.stackFrames.Front())
}

func (this *Stack) BindVariable(name string, mlrval *types.Mlrval) {
	this.stackFrames.Front().Value.(*StackFrame).Bind(name, mlrval)
}

// Returns nil on no-such
func (this *Stack) ReadVariable(name string) *types.Mlrval {

	// Scope-walk
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		mlrval := stackFrame.ReadVariable(name)
		if mlrval != nil {
			return mlrval
		}
	}
	return nil
}

// Returns nil on no-such
func (this *Stack) Dump() {
	fmt.Printf("STACK FRAMES (count %d):\n", this.stackFrames.Len())
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		fmt.Printf("  VARIABLES (count %d):\n", len(stackFrame.vars))
		for k, v := range stackFrame.vars {
			fmt.Printf("    %-16s %s\n", k, v.String())
		}
	}
}

// ----------------------------------------------------------------
type StackFrame struct {
	// TODO: just a map for now. In the C impl, pre-computation of
	// name-to-array-slot indices was an important optimization, especially for
	// compute-intensive scenarios.
	vars map[string]*types.Mlrval
}

func NewStackFrame() *StackFrame {
	return &StackFrame{
		vars: make(map[string]*types.Mlrval),
	}
}

func (this *StackFrame) Clear() {
	this.vars = make(map[string]*types.Mlrval)
}

func (this *StackFrame) Bind(name string, mlrval *types.Mlrval) {
	this.vars[name] = mlrval.Copy()
}

// Returns nil on no such
func (this *StackFrame) ReadVariable(name string) *types.Mlrval {
	return this.vars[name]
}
