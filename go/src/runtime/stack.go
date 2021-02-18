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
// STACK VARIABLE

// StackVariable is an opaque handle which a callsite can hold onto, which
// keeps stack-offset information in it that is private to us.
type StackVariable struct {
	name string

	// Type like "int" or "num" or "var" is stored in the stack itself.  A
	// StackVariable can appear in the CST (concrete syntax tree) on either the
	// left-hand side or right-hande side of an assignment -- in the latter
	// case the callsite won't know the type until the value is read off the
	// stack.

	// TODO: comment
	frameSetIndex int
	indexInFrame  int
}

// TODO: be sure to invalidate slot 0 for struct uninit
func NewStackVariable(name string) *StackVariable {
	return &StackVariable{
		name:          name,
		frameSetIndex: 0,
		indexInFrame:  0,
	}
}

// ================================================================
// STACK METHODS

type Stack struct {
	// list of *StackFrameSet
	stackFrameSets *list.List

	// Invariant: equal to the head of the stackFrameSets list. This is cached
	// since all sets/gets in between frameset-push and frameset-pop will all
	// and only be operating on the head.
	head *StackFrameSet
}

func NewStack() *Stack {
	stackFrameSets := list.New()
	head := newStackFrameSet()
	stackFrameSets.PushFront(head)
	return &Stack{
		stackFrameSets: stackFrameSets,
		head:           head,
	}
}

// For when a user-defined function/subroutine is being entered
func (this *Stack) PushStackFrameSet() {
	this.head = newStackFrameSet()
	this.stackFrameSets.PushFront(this.head)
}

// For when a user-defined function/subroutine is being exited
func (this *Stack) PopStackFrameSet() {
	this.stackFrameSets.Remove(this.stackFrameSets.Front())
	this.head = this.stackFrameSets.Front().Value.(*StackFrameSet)
}

// ----------------------------------------------------------------
// Delegations to topmost frameset

// For when an if/for/etc block is being entered
func (this *Stack) PushStackFrame() {
	this.head.pushStackFrame()
}

// For when an if/for/etc block is being exited
func (this *Stack) PopStackFrame() {
	this.head.popStackFrame()
}

// For 'num a = 2', setting a variable at the current frame regardless of outer
// scope.  It's an error to define it again in the same scope, whether the type
// is the same or not.
func (this *Stack) DefineTypedAtScope(
	stackVariable *StackVariable,
	typeName string,
	mlrval *types.Mlrval,
) error {
	return this.head.defineTypedAtScope(stackVariable, typeName, mlrval)
}

// For untyped declarations at the current scope -- these are in binds of
// for-loop variables, except for triple-for.
// E.g. 'for (k, v in $*)' uses SetAtScope.
// E.g. 'for (int i = 0; i < 10; i += 1)' uses DefineTypedAtScope
// E.g. 'for (i = 0; i < 10; i += 1)' uses Set.
func (this *Stack) SetAtScope(
	stackVariable *StackVariable,
	mlrval *types.Mlrval,
) error {
	return this.head.setAtScope(stackVariable, mlrval)
}

// For 'a = 2', checking for outer-scoped to maybe reuse, else insert new in
// current frame. If the variable is entirely new it's set in the current frame
// with no type-checking. If it's not new the assignment is subject to
// type-checking for wherever the variable was defined. E.g. if it was
// previously defined with 'str a = "hello"' then this Set returns an error.
// However if it waa previously assigned untyped with 'a = "hello"' then the
// assignment is OK.
func (this *Stack) Set(
	stackVariable *StackVariable,
	mlrval *types.Mlrval,
) error {
	return this.head.set(stackVariable, mlrval)
}

// E.g. 'x[1] = 2' where the variable x may or may not have been already set.
func (this *Stack) SetIndexed(
	stackVariable *StackVariable,
	indices []*types.Mlrval,
	mlrval *types.Mlrval,
) error {
	return this.head.setIndexed(stackVariable, indices, mlrval)
}

// E.g. 'unset x'
func (this *Stack) Unset(
	stackVariable *StackVariable,
) {
	this.head.unset(stackVariable)
}

// E.g. 'unset x[1]'
func (this *Stack) UnsetIndexed(
	stackVariable *StackVariable,
	indices []*types.Mlrval,
) {
	this.head.unsetIndexed(stackVariable, indices)
}

// Returns nil on no-such
func (this *Stack) Get(
	stackVariable *StackVariable,
) *types.Mlrval {
	return this.head.get(stackVariable)
}

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
	// TODO: to array
	stackFrames := list.New()
	stackFrames.PushFront(newStackFrame())
	return &StackFrameSet{
		stackFrames: stackFrames,
	}
}

func (this *StackFrameSet) pushStackFrame() {
	// TODO: to array
	this.stackFrames.PushFront(newStackFrame())
}

func (this *StackFrameSet) popStackFrame() {
	// TODO: to array
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

// See Stack.DefineTypedAtScope comments above
func (this *StackFrameSet) defineTypedAtScope(
	stackVariable *StackVariable,
	typeName string,
	mlrval *types.Mlrval,
) error {
	// TODO: to array
	return this.stackFrames.Front().Value.(*StackFrame).defineTyped(
		stackVariable, typeName, mlrval,
	)
}

// See Stack.SetAtScope comments above
func (this *StackFrameSet) setAtScope(
	stackVariable *StackVariable,
	mlrval *types.Mlrval,
) error {
	// TODO: to array
	return this.stackFrames.Front().Value.(*StackFrame).set(stackVariable, mlrval)
}

// See Stack.Set comments above
func (this *StackFrameSet) set(
	stackVariable *StackVariable,
	mlrval *types.Mlrval,
) error {
	// TODO: to array
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		if stackFrame.has(stackVariable) {
			return stackFrame.set(stackVariable, mlrval)
		}
	}
	return this.setAtScope(stackVariable, mlrval)
}

// See Stack.SetIndexed comments above
func (this *StackFrameSet) setIndexed(
	stackVariable *StackVariable,
	indices []*types.Mlrval,
	mlrval *types.Mlrval,
) error {
	// TODO: to array
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		if stackFrame.has(stackVariable) {
			return stackFrame.setIndexed(stackVariable, indices, mlrval)
		}
	}
	return this.stackFrames.Front().Value.(*StackFrame).setIndexed(stackVariable, indices, mlrval)
}

// See Stack.Unset comments above
func (this *StackFrameSet) unset(
	stackVariable *StackVariable,
) {
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		if stackFrame.has(stackVariable) {
			stackFrame.unset(stackVariable)
			return
		}
	}
}

// See Stack.UnsetIndexed comments above
func (this *StackFrameSet) unsetIndexed(
	stackVariable *StackVariable,
	indices []*types.Mlrval,
) {
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		if stackFrame.has(stackVariable) {
			stackFrame.unsetIndexed(stackVariable, indices)
			return
		}
	}
}

// Returns nil on no-such
func (this *StackFrameSet) get(
	stackVariable *StackVariable,
) *types.Mlrval {
	// Scope-walk
	for entry := this.stackFrames.Front(); entry != nil; entry = entry.Next() {
		stackFrame := entry.Value.(*StackFrame)
		mlrval := stackFrame.get(stackVariable)
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
func (this *StackFrame) get(
	stackVariable *StackVariable,
) *types.Mlrval {
	slot := this.vars[stackVariable.name]
	if slot == nil {
		return nil
	} else {
		return slot.GetValue()
	}
}

func (this *StackFrame) has(
	stackVariable *StackVariable,
) bool {
	return this.vars[stackVariable.name] != nil
}

func (this *StackFrame) clear() {
	this.vars = make(map[string]*types.TypeGatedMlrvalVariable)
}

// TODO: audit for honor of error-return at callsites
func (this *StackFrame) set(
	stackVariable *StackVariable,
	mlrval *types.Mlrval,
) error {
	slot := this.vars[stackVariable.name]
	if slot == nil {
		slot, err := types.NewTypeGatedMlrvalVariable(stackVariable.name, "any", mlrval)
		if err != nil {
			return err
		} else {
			this.vars[stackVariable.name] = slot
			return nil
		}
		this.vars[stackVariable.name] = slot
		return nil
	} else {
		return slot.Assign(mlrval)
	}
}

// TODO: audit for honor of error-return at callsites
func (this *StackFrame) defineTyped(
	stackVariable *StackVariable,
	typeName string,
	mlrval *types.Mlrval,
) error {
	slot := this.vars[stackVariable.name]
	if slot == nil {
		slot, err := types.NewTypeGatedMlrvalVariable(stackVariable.name, typeName, mlrval)
		if err != nil {
			return err
		} else {
			this.vars[stackVariable.name] = slot
			return nil
		}
		this.vars[stackVariable.name] = slot
		return nil
	} else {
		return errors.New(
			fmt.Sprintf(
				"%s: variable %s has already been defined in the same scope.",
				lib.MlrExeName(), stackVariable.name,
			),
		)
	}
}

func (this *StackFrame) unset(
	stackVariable *StackVariable,
) {
	slot := this.vars[stackVariable.name]
	if slot != nil {
		slot.Unassign()
	}
}

// TODO: audit for honor of error-return at callsites
func (this *StackFrame) setIndexed(
	stackVariable *StackVariable,
	indices []*types.Mlrval,
	mlrval *types.Mlrval,
) error {
	value := this.get(stackVariable)
	if value == nil {
		lib.InternalCodingErrorIf(len(indices) < 1)
		leadingIndex := indices[0]
		if leadingIndex.IsString() || leadingIndex.IsInt() {
			newval := types.MlrvalEmptyMap()
			newval.PutIndexed(indices, mlrval)
			return this.set(stackVariable, &newval)
		} else {
			return errors.New(
				fmt.Sprintf(
					"%s: map indices must be int or string; got %s.\n",
					lib.MlrExeName(), leadingIndex.GetTypeName(),
				),
			)
		}
	} else {
		// For example maybe the variable exists and is an array but the
		// leading index is a string.
		return value.PutIndexed(indices, mlrval)
	}
}

func (this *StackFrame) unsetIndexed(
	stackVariable *StackVariable,
	indices []*types.Mlrval,
) {
	value := this.get(stackVariable)
	if value == nil {
		return
	}
	value.RemoveIndexed(indices)
}
