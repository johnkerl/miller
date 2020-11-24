// ================================================================
// Stack frames for begin/end/if/for/function blocks
//
// A Miller DSL stack has two levels of nesting:
// * A Stack contains a list of StackFrameSet, one per function or Miller outermost statement block
// * A StackFrameSet contains a list of StackFrame, one per if/for/etc within a function
//
// This is because of the following.
//
// (1) a = 1              <-- outer stack frame in same frameset
//     if (condition) {   <-- inner stack frame in same frameset
//       a = 2            <-- this should update the outer 'a', not create new inner 'a'
//     }
//
// (2) a = 1              <-- outer stack frame in same frameset
//     if (condition) {   <-- inner stack frame in same frameset
//       var a = 2        <-- this should create new inner 'a', not update the outer 'a'
//     }
//
// (3) a = 1              <-- outer stack frame
//     func f() {         <-- stack frame in a new frameset
//       a = 2            <-- this should create new inner 'a', not update the outer 'a'
//     }
// ================================================================

package cst

import (
	"container/list"
	"fmt"

	"miller/lib"
	"miller/types"
)

// ================================================================
// STACK METHODS

type Stack struct {
	stackFrameSets *list.List // list of *StackFrameSet
}

func NewStack() *Stack {
	stackFrameSets := list.New()
	stackFrameSets.PushFront(NewStackFrameSet())
	return &Stack{
		stackFrameSets: stackFrameSets,
	}
}

// For when a user-defined function/subroutine is being entered
func (this *Stack) PushStackFrameSet() {
	this.stackFrameSets.PushFront(NewStackFrameSet())
}

// For when a user-defined function/subroutine is being exited
func (this *Stack) PopStackFrameSet() {
	this.stackFrameSets.Remove(this.stackFrameSets.Front())
}

// ----------------------------------------------------------------
// Delegations to topmost frameset

// For when an if/for/etc block is being entered
func (this *Stack) PushStackFrame() {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	head.PushStackFrame()
}

// For when an if/for/etc block is being exited
func (this *Stack) PopStackFrame() {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	head.PopStackFrame()
}

// For 'var a = 2', setting a variable at the current frame regardless of outer scope.
func (this *Stack) BindVariable(
	name string,
	mlrval *types.Mlrval,
) {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	head.BindVariable(name, mlrval)
}

func (this *Stack) BindVariableIndexed(
	name string,
	indices []*types.Mlrval,
	mlrval *types.Mlrval,
) {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	head.BindVariableIndexed(name, indices, mlrval)
}

// For 'a = 2', checking for outer-scoped to maybe reuse, else insert new in current frame.
func (this *Stack) SetVariable(name string, mlrval *types.Mlrval) {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	head.SetVariable(name, mlrval)
}

func (this *Stack) SetVariableIndexed(
	name string,
	indices []*types.Mlrval,
	mlrval *types.Mlrval,
) {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	head.SetVariableIndexed(name, indices, mlrval)
}

func (this *Stack) UnsetVariable(name string) {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	head.UnsetVariable(name)
}

func (this *Stack) UnsetVariableIndexed(
	name string,
	indices []*types.Mlrval,
) {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	head.UnsetVariableIndexed(name, indices)
}

// Returns nil on no-such
func (this *Stack) ReadVariable(name string) *types.Mlrval {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	return head.ReadVariable(name)
}

// ----------------------------------------------------------------
func (this *Stack) Dump() {
	fmt.Printf("STACK FRAMESETS (count %d):\n", this.stackFrameSets.Len())
	for entry := this.stackFrameSets.Front(); entry != nil; entry = entry.Next() {
		stackFrameSet := entry.Value.(*StackFrameSet)
		stackFrameSet.Dump()
	}
}

// ================================================================
// STACKFRAMESET METHODS

type StackFrameSet struct {
	stackFrames *list.List // list of *StackFrame
}

func NewStackFrameSet() *StackFrameSet {
	stackFrames := list.New()
	stackFrames.PushFront(NewStackFrame())
	return &StackFrameSet{
		stackFrames: stackFrames,
	}
}

func (this *StackFrameSet) PushStackFrame() {
	this.stackFrames.PushFront(NewStackFrame())
}

func (this *StackFrameSet) PopStackFrame() {
	this.stackFrames.Remove(this.stackFrames.Front())
}

func (this *StackFrameSet) Dump() {
	fmt.Printf("  STACK FRAMES (count %d):\n", this.stackFrames.Len())
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		fmt.Printf("    VARIABLES (count %d):\n", len(stackFrame.vars))
		for k, v := range stackFrame.vars {
			fmt.Printf("      %-16s %s\n", k, v.String())
		}
	}
}

// ----------------------------------------------------------------
// Sets the variable at the current frame whether it's defined outer from there
// or not.
//
// OK to use BindVariable:
//
//   k = 1                 <-- top-level -frame, k=1
//   for (k in $*) { ... } <-- another k is bound in the loop
//   $k = k                <-- k is still 1
//
// Not OK to use BindVariable:
//
//   z = 1         <-- top-level frame, z=1
//   if (NR < 2) {
//     z = 2       <-- this should adjust top-level z, not bind within if-block
//   } else {
//     z = 3       <-- this should adjust top-level z, not bind within else-block
//   }
//   $z = z        <-- z should be 2 or 3, not 1

func (this *StackFrameSet) BindVariable(
	name string,
	mlrval *types.Mlrval,
) {
	this.stackFrames.Front().Value.(*StackFrame).Set(name, mlrval)
}

func (this *StackFrameSet) BindVariableIndexed(
	name string,
	indices []*types.Mlrval,
	mlrval *types.Mlrval,
) {
	this.stackFrames.Front().Value.(*StackFrame).SetIndexed(name, indices, mlrval)
}

// Used for the above BindVariable example where we look for outer-scope names,
// then set a new one only if not found in an outer scope.
func (this *StackFrameSet) SetVariable(name string, mlrval *types.Mlrval) {
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		if stackFrame.Has(name) {
			stackFrame.Set(name, mlrval)
			return
		}
	}
	this.BindVariable(name, mlrval)
}

func (this *StackFrameSet) SetVariableIndexed(
	name string,
	indices []*types.Mlrval,
	mlrval *types.Mlrval,
) {
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		if stackFrame.Has(name) {
			stackFrame.SetIndexed(name, indices, mlrval)
			return
		}
	}
	this.BindVariableIndexed(name, indices, mlrval)
}

// ----------------------------------------------------------------
func (this *StackFrameSet) UnsetVariable(name string) {
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		if stackFrame.Has(name) {
			stackFrame.Unset(name)
			return
		}
	}
}

func (this *StackFrameSet) UnsetVariableIndexed(
	name string,
	indices []*types.Mlrval,
) {
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		if stackFrame.Has(name) {
			stackFrame.UnsetIndexed(name, indices)
			return
		}
	}
}

// ----------------------------------------------------------------
// Returns nil on no-such
func (this *StackFrameSet) ReadVariable(name string) *types.Mlrval {
	// Scope-walk
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		mlrval := stackFrame.Get(name)
		if mlrval != nil {
			return mlrval
		}
	}
	return nil
}

// ================================================================
// STACKFRAME METHODS

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

// Returns nil on no such
func (this *StackFrame) Get(name string) *types.Mlrval {
	return this.vars[name]
}

// Returns nil on no such
func (this *StackFrame) Has(name string) bool {
	return this.vars[name] != nil
}

func (this *StackFrame) Clear() {
	this.vars = make(map[string]*types.Mlrval)
}

func (this *StackFrame) Set(name string, mlrval *types.Mlrval) {
	this.vars[name] = mlrval.Copy()
}

func (this *StackFrame) Unset(name string) {
	value := types.MlrvalFromAbsent()
	this.vars[name] = &value
}

func (this *StackFrame) SetIndexed(
	name string,
	indices []*types.Mlrval,
	mlrval *types.Mlrval,
) {
	value := this.Get(name)
	if value == nil {
		lib.InternalCodingErrorIf(len(indices) < 1)
		leadingIndex := indices[0]
		if leadingIndex.IsString() {
			newval := types.MlrvalEmptyMap()
			newval.PutIndexed(indices, mlrval)
			this.Set(name, &newval)
		} else if leadingIndex.IsInt() {
			newval := types.MlrvalEmptyArray()
			newval.PutIndexed(indices, mlrval)
			this.Set(name, &newval)
		} else {
			// TODO:
			// return errors.New("...");
		}
	} else {
		// TODO: propagate error return.
		// For example maybe the variable exists and is an array but
		// the leading index is a string.
		_ = value.PutIndexed(indices, mlrval)
	}
}

func (this *StackFrame) UnsetIndexed(
	name string,
	indices []*types.Mlrval,
) {
	value := this.Get(name)
	if value == nil {
		return
	}
	value.UnsetIndexed(indices)
}
