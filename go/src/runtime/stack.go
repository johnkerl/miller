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
	stackFrameSets.PushFront(newStackFrameSet())
	return &Stack{
		stackFrameSets: stackFrameSets,
	}
}

// For when a user-defined function/subroutine is being entered
func (this *Stack) PushStackFrameSet() {
	this.stackFrameSets.PushFront(newStackFrameSet())
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
	head.pushStackFrame()
}

// For when an if/for/etc block is being exited
func (this *Stack) PopStackFrame() {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	head.popStackFrame()
}

// For 'var a = 2', setting a variable at the current frame regardless of outer
// scope.  It's an error to define it again in the same scope, whether the type
// is the same or not.
func (this *Stack) DefineTypedAtScope(
	variableName string,
	typeName string,
	mlrval *types.Mlrval,
) error {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	return head.defineTypedAtScope(variableName, typeName, mlrval)
}

// TODO: comment -- for-loops (but not triple-fors) ...
func (this *Stack) SetAtScope(
	variableName string,
	mlrval *types.Mlrval,
) error {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	return head.setAtScope(variableName, mlrval)
}

func (this *Stack) DefineTypedAtScopeIndexed(
	variableName string,
	typeName string,
	indices []*types.Mlrval,
	mlrval *types.Mlrval,
) error {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	return head.defineTypedAtScopeIndexed(variableName, typeName, indices, mlrval)
}

// For 'a = 2', checking for outer-scoped to maybe reuse, else insert new in
// current frame.
func (this *Stack) Set(
	variableName string,
	mlrval *types.Mlrval,
) error {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	return head.set(variableName, mlrval)
}

func (this *Stack) SetIndexed(
	variableName string,
	indices []*types.Mlrval,
	mlrval *types.Mlrval,
) error {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	return head.setIndexed(variableName, indices, mlrval)
}

func (this *Stack) Unset(variableName string) {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	head.unset(variableName)
}

func (this *Stack) UnsetIndexed(
	variableName string,
	indices []*types.Mlrval,
) {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	head.unsetIndexed(variableName, indices)
}

// Returns nil on no-such
func (this *Stack) Get(variableName string) *types.Mlrval {
	head := this.stackFrameSets.Front().Value.(*StackFrameSet)
	return head.get(variableName)
}

// ----------------------------------------------------------------
func (this *Stack) Dump() {
	fmt.Printf("STACK FRAMESETS (count %d):\n", this.stackFrameSets.Len())
	for entry := this.stackFrameSets.Front(); entry != nil; entry = entry.Next() {
		stackFrameSet := entry.Value.(*StackFrameSet)
		stackFrameSet.dump()
	}
}

// ================================================================
// STACKFRAMESET METHODS

type StackFrameSet struct {
	stackFrames *list.List // list of *StackFrame
}

func newStackFrameSet() *StackFrameSet {
	stackFrames := list.New()
	stackFrames.PushFront(newStackFrame())
	return &StackFrameSet{
		stackFrames: stackFrames,
	}
}

func (this *StackFrameSet) pushStackFrame() {
	this.stackFrames.PushFront(newStackFrame())
}

func (this *StackFrameSet) popStackFrame() {
	this.stackFrames.Remove(this.stackFrames.Front())
}

func (this *StackFrameSet) dump() {
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
// OK to use DefineTypedAtScope:
//
//   k = 1                 <-- top-level -frame, k=1
//   for (k in $*) { ... } <-- another k is bound in the loop
//   $k = k                <-- k is still 1
//
// Not OK to use DefineTypedAtScope:
//
//   z = 1         <-- top-level frame, z=1
//   if (NR < 2) {
//     z = 2       <-- this should adjust top-level z, not bind within if-block
//   } else {
//     z = 3       <-- this should adjust top-level z, not bind within else-block
//   }
//   $z = z        <-- z should be 2 or 3, not 1

func (this *StackFrameSet) defineTypedAtScope(
	variableName string,
	typeName string,
	mlrval *types.Mlrval,
) error {
	return this.stackFrames.Front().Value.(*StackFrame).defineTyped(variableName, typeName, mlrval)
}

func (this *StackFrameSet) defineTypedAtScopeIndexed(
	variableName string,
	typeName string,
	indices []*types.Mlrval,
	mlrval *types.Mlrval,
) error {
	// TODO WTF DOES THIS DO
	return this.stackFrames.Front().Value.(*StackFrame).setIndexed(variableName, indices, mlrval)
}

func (this *StackFrameSet) setAtScopeIndexed(
	variableName string,
	indices []*types.Mlrval,
	mlrval *types.Mlrval,
) error {
	// TODO WTF DOES THIS DO
	return this.stackFrames.Front().Value.(*StackFrame).setIndexed(variableName, indices, mlrval)
}

// TODO: comment
func (this *StackFrameSet) setAtScope(
	variableName string,
	mlrval *types.Mlrval,
) error {
	return this.stackFrames.Front().Value.(*StackFrame).set(variableName, mlrval)
}

// Used for the above DefineTypedAtScope example where we look for outer-scope names,
// then set a new one only if not found in an outer scope.
func (this *StackFrameSet) set(
	variableName string,
	mlrval *types.Mlrval,
) error {
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		if stackFrame.has(variableName) {
			return stackFrame.set(variableName, mlrval)
		}
	}
	return this.setAtScope(variableName, mlrval)
}

func (this *StackFrameSet) setIndexed(
	variableName string,
	indices []*types.Mlrval,
	mlrval *types.Mlrval,
) error {
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		if stackFrame.has(variableName) {
			return stackFrame.setIndexed(variableName, indices, mlrval)
		}
	}
	return this.setAtScopeIndexed(variableName, indices, mlrval)
}

// ----------------------------------------------------------------
func (this *StackFrameSet) unset(variableName string) {
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		if stackFrame.has(variableName) {
			stackFrame.unset(variableName)
			return
		}
	}
}

func (this *StackFrameSet) unsetIndexed(
	variableName string,
	indices []*types.Mlrval,
) {
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		if stackFrame.has(variableName) {
			stackFrame.unsetIndexed(variableName, indices)
			return
		}
	}
}

// ----------------------------------------------------------------
// Returns nil on no-such
func (this *StackFrameSet) get(variableName string) *types.Mlrval {
	// Scope-walk
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		mlrval := stackFrame.get(variableName)
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

func newStackFrame() *StackFrame {
	return &StackFrame{
		vars: make(map[string]*types.TypeGatedMlrvalVariable),
	}
}

// Returns nil on no such
func (this *StackFrame) get(variableName string) *types.Mlrval {
	slot := this.vars[variableName]
	if slot == nil {
		return nil
	} else {
		return slot.GetValue()
	}
}

// Returns nil on no such
func (this *StackFrame) has(variableName string) bool {
	return this.vars[variableName] != nil
}

func (this *StackFrame) clear() {
	this.vars = make(map[string]*types.TypeGatedMlrvalVariable)
}

// TODO: audit for honor of error-return at callsites
func (this *StackFrame) set(
	variableName string,
	mlrval *types.Mlrval,
) error {
	slot := this.vars[variableName]
	if slot == nil {
		slot, err := types.NewTypeGatedMlrvalVariable(variableName, "any", mlrval)
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

// TODO: audit for honor of error-return at callsites
func (this *StackFrame) defineTyped(
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
		return errors.New(
			fmt.Sprintf(
				"%s: variable %s has already been defined in the same scope.",
				lib.MlrExeName(), variableName,
			),
		)
	}
}

func (this *StackFrame) unset(variableName string) {
	slot := this.vars[variableName]
	if slot != nil {
		slot.Unassign()
	}
}

// TODO: audit for honor of error-return at callsites
func (this *StackFrame) setIndexed(
	variableName string,
	indices []*types.Mlrval,
	mlrval *types.Mlrval,
) error {
	value := this.get(variableName)
	if value == nil {
		lib.InternalCodingErrorIf(len(indices) < 1)
		leadingIndex := indices[0]
		if leadingIndex.IsString() || leadingIndex.IsInt() {
			newval := types.MlrvalEmptyMap()
			newval.PutIndexed(indices, mlrval)
			return this.set(variableName, &newval)
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

func (this *StackFrame) unsetIndexed(
	variableName string,
	indices []*types.Mlrval,
) {
	value := this.get(variableName)
	if value == nil {
		return
	}
	value.RemoveIndexed(indices)
}
