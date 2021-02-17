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

package runtime

import (
	"container/list"
	"errors"
	"fmt"

	"miller/src/lib"
	"miller/src/types"
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
func (this *Stack) SetAtScope(
	variableName string,
	typeName string,
	mlrval *types.Mlrval,
) error {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	return head.SetAtScope(variableName, typeName, mlrval)
}

func (this *Stack) SetAtScopeIndexed(
	variableName string,
	typeName string,
	indices []*types.Mlrval,
	mlrval *types.Mlrval,
) error {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	return head.SetAtScopeIndexed(variableName, typeName, indices, mlrval)
}

// For 'a = 2', checking for outer-scoped to maybe reuse, else insert new in current frame.
func (this *Stack) Set(
	variableName string,
	typeName string,
	mlrval *types.Mlrval,
) error {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	return head.Set(variableName, typeName, mlrval)
}

func (this *Stack) SetIndexed(
	variableName string,
	typeName string,
	indices []*types.Mlrval,
	mlrval *types.Mlrval,
) error {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	return head.SetIndexed(variableName, typeName, indices, mlrval)
}

func (this *Stack) Unset(variableName string) {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	head.Unset(variableName)
}

func (this *Stack) UnsetIndexed(
	variableName string,
	indices []*types.Mlrval,
) {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	head.UnsetIndexed(variableName, indices)
}

// Returns nil on no-such
func (this *Stack) Get(variableName string) *types.Mlrval {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	return head.Get(variableName)
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
			fmt.Printf("      %-16s %s\n", k, v.ValueString())
		}
	}
}

// ----------------------------------------------------------------
// Sets the variable at the current frame whether it's defined outer from there
// or not.
//
// OK to use SetAtScope:
//
//   k = 1                 <-- top-level -frame, k=1
//   for (k in $*) { ... } <-- another k is bound in the loop
//   $k = k                <-- k is still 1
//
// Not OK to use SetAtScope:
//
//   z = 1         <-- top-level frame, z=1
//   if (NR < 2) {
//     z = 2       <-- this should adjust top-level z, not bind within if-block
//   } else {
//     z = 3       <-- this should adjust top-level z, not bind within else-block
//   }
//   $z = z        <-- z should be 2 or 3, not 1

func (this *StackFrameSet) SetAtScope(
	variableName string,
	typeName string,
	mlrval *types.Mlrval,
) error {
	return this.stackFrames.Front().Value.(*StackFrame).Set(variableName, typeName, mlrval)
}

func (this *StackFrameSet) SetAtScopeIndexed(
	variableName string,
	typeName string,
	indices []*types.Mlrval,
	mlrval *types.Mlrval,
) error {
	return this.stackFrames.Front().Value.(*StackFrame).SetIndexed(variableName, typeName, indices, mlrval)
}

// Used for the above SetAtScope example where we look for outer-scope names,
// then set a new one only if not found in an outer scope.
func (this *StackFrameSet) Set(
	variableName string,
	typeName string,
	mlrval *types.Mlrval,
) error {
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		if stackFrame.Has(variableName) {
			return stackFrame.Set(variableName, typeName, mlrval)
		}
	}
	return this.SetAtScope(variableName, typeName, mlrval)
}

func (this *StackFrameSet) SetIndexed(
	variableName string,
	typeName string,
	indices []*types.Mlrval,
	mlrval *types.Mlrval,
) error {
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		if stackFrame.Has(variableName) {
			return stackFrame.SetIndexed(variableName, typeName, indices, mlrval)
		}
	}
	return this.SetAtScopeIndexed(variableName, typeName, indices, mlrval)
}

// ----------------------------------------------------------------
func (this *StackFrameSet) Unset(variableName string) {
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		if stackFrame.Has(variableName) {
			stackFrame.Unset(variableName)
			return
		}
	}
}

func (this *StackFrameSet) UnsetIndexed(
	variableName string,
	indices []*types.Mlrval,
) {
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		if stackFrame.Has(variableName) {
			stackFrame.UnsetIndexed(variableName, indices)
			return
		}
	}
}

// ----------------------------------------------------------------
// Returns nil on no-such
func (this *StackFrameSet) Get(variableName string) *types.Mlrval {
	// Scope-walk
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		mlrval := stackFrame.Get(variableName)
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
	vars map[string]*types.TypeGatedMlrvalVariable
}

func NewStackFrame() *StackFrame {
	return &StackFrame{
		vars: make(map[string]*types.TypeGatedMlrvalVariable),
	}
}

// Returns nil on no such
func (this *StackFrame) Get(variableName string) *types.Mlrval {
	slot := this.vars[variableName]
	if slot == nil {
		return nil
	} else {
		return slot.GetValue()
	}
}

// Returns nil on no such
func (this *StackFrame) Has(variableName string) bool {
	return this.vars[variableName] != nil
}

func (this *StackFrame) Clear() {
	this.vars = make(map[string]*types.TypeGatedMlrvalVariable)
}

// TODO: audit for honor of error-return at callsites
func (this *StackFrame) Set(
	variableName string,
	typeName string,
	mlrval *types.Mlrval,
) error {
	slot := this.vars[variableName]
	if slot == nil {
		slot, err := types.NewTypeGatedMlrvalVariable(variableName, typeName, mlrval)
		if err != nil {
			return err
		} else {
			this.vars[variableName] = slot
			return nil
		}
		this.vars[variableName] = slot
		return nil
	} else {
		return slot.Assign(mlrval.Copy())
	}
}

func (this *StackFrame) Unset(variableName string) {
	slot := this.vars[variableName]
	if slot != nil {
		slot.Unassign()
	}
}

// TODO: audit for honor of error-return at callsites
func (this *StackFrame) SetIndexed(
	variableName string,
	typeName string,
	indices []*types.Mlrval,
	mlrval *types.Mlrval,
) error {
	value := this.Get(variableName)
	if value == nil {
		lib.InternalCodingErrorIf(len(indices) < 1)
		leadingIndex := indices[0]
		if leadingIndex.IsString() || leadingIndex.IsInt() {
			newval := types.MlrvalEmptyMap()
			newval.PutIndexed(indices, mlrval)
			return this.Set(variableName, typeName, &newval)
		} else {
			return errors.New(
				fmt.Sprintf(
					"%s: map indices must be int or string; got %s.\n",
					lib.MlrExeName(), leadingIndex.GetTypeName(),
				),
			)
		}
	} else {
		// TODO: propagate error return.
		// For example maybe the variable exists and is an array but
		// the leading index is a string.
		return value.PutIndexed(indices, mlrval)
	}
}

func (this *StackFrame) UnsetIndexed(
	variableName string,
	indices []*types.Mlrval,
) {
	value := this.Get(variableName)
	if value == nil {
		return
	}
	value.RemoveIndexed(indices)
}
