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

package runtime

import (
	"fmt"

	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

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
}

func NewStackVariable(name string) *StackVariable {
	return NewStackVariableAux(name, true)
}

// TODO: comment re function literals
func NewStackVariableAux(name string, cacheable bool) *StackVariable {
	return &StackVariable{
		name: name,
	}
}

func (sv *StackVariable) GetName() string {
	return sv.name
}

// STACK METHODS

type Stack struct {
	// Save/restore stack of framesets, one pushed per user-defined
	// function/subroutine call. The CURRENT frameset is the tail element
	// (stackFrameSets[len-1]); pushing appends and popping truncates, so neither
	// allocates a new slice once capacity is established. (Order among the saved
	// sets is irrelevant: all get/set go through the cached head.)
	stackFrameSets []*StackFrameSet

	// Invariant: equal to the tail of the stackFrameSets list. This is cached
	// since all sets/gets in between frameset-push and frameset-pop will all
	// and only be operating on the head.
	head *StackFrameSet

	// pool retains popped framesets for reuse, so repeated function calls do not
	// each allocate a fresh StackFrameSet (and its initial StackFrame).
	pool []*StackFrameSet
}

func NewStack() *Stack {
	head := newStackFrameSet()
	stackFrameSets := []*StackFrameSet{head}
	return &Stack{
		stackFrameSets: stackFrameSets,
		head:           head,
	}
}

// For when a user-defined function/subroutine is being entered
func (stack *Stack) PushStackFrameSet() {
	var frameset *StackFrameSet
	n := len(stack.pool)
	if n > 0 {
		frameset = stack.pool[n-1]
		stack.pool = stack.pool[:n-1]
		frameset.reset()
	} else {
		frameset = newStackFrameSet()
	}
	stack.stackFrameSets = append(stack.stackFrameSets, frameset)
	stack.head = frameset
}

// For when a user-defined function/subroutine is being exited
func (stack *Stack) PopStackFrameSet() {
	n := len(stack.stackFrameSets)
	popped := stack.stackFrameSets[n-1]
	stack.stackFrameSets = stack.stackFrameSets[0 : n-1]
	stack.pool = append(stack.pool, popped)
	stack.head = stack.stackFrameSets[len(stack.stackFrameSets)-1]
}

// All of these are simply delegations to the head frameset

// For when an if/for/etc block is being entered
func (stack *Stack) PushStackFrame() {
	stack.head.pushStackFrame()
}

// For when an if/for/etc block is being exited
func (stack *Stack) PopStackFrame() {
	stack.head.popStackFrame()
}

// Returns nil on no-such
func (stack *Stack) Get(
	stackVariable *StackVariable,
) *mlrval.Mlrval {
	return stack.head.get(stackVariable)
}

// For 'num a = 2', setting a variable at the current frame regardless of outer
// scope.  It's an error to define it again in the same scope, whether the type
// is the same or not.
func (stack *Stack) DefineTypedAtScope(
	stackVariable *StackVariable,
	typeName string,
	mlrval *mlrval.Mlrval,
) error {
	return stack.head.defineTypedAtScope(stackVariable, typeName, mlrval)
}

// For untyped declarations at the current scope -- these are in binds of
// for-loop variables, except for triple-for.
// E.g. 'for (k, v in $*)' uses SetAtScope.
// E.g. 'for (int i = 0; i < 10; i += 1)' uses DefineTypedAtScope
// E.g. 'for (i = 0; i < 10; i += 1)' uses Set.
func (stack *Stack) SetAtScope(
	stackVariable *StackVariable,
	mlrval *mlrval.Mlrval,
) error {
	return stack.head.setAtScope(stackVariable, mlrval)
}

// For 'a = 2', checking for outer-scoped to maybe reuse, else insert new in
// current frame. If the variable is entirely new it's set in the current frame
// with no type-checking. If it's not new the assignment is subject to
// type-checking for wherever the variable was defined. E.g. if it was
// previously defined with 'str a = "hello"' then this Set returns an error.
// However if it waa previously assigned untyped with 'a = "hello"' then the
// assignment is OK.
func (stack *Stack) Set(
	stackVariable *StackVariable,
	mlrval *mlrval.Mlrval,
) error {
	return stack.head.set(stackVariable, mlrval)
}

// E.g. 'x[1] = 2' where the variable x may or may not have been already set.
func (stack *Stack) SetIndexed(
	stackVariable *StackVariable,
	indices []*mlrval.Mlrval,
	mlrval *mlrval.Mlrval,
) error {
	return stack.head.setIndexed(stackVariable, indices, mlrval)
}

// E.g. 'unset x'
func (stack *Stack) Unset(
	stackVariable *StackVariable,
) {
	stack.head.unset(stackVariable)
}

// E.g. 'unset x[1]'
func (stack *Stack) UnsetIndexed(
	stackVariable *StackVariable,
	indices []*mlrval.Mlrval,
) {
	stack.head.unsetIndexed(stackVariable, indices)
}

// STACKFRAMESET METHODS

const stackFrameSetInitCap = 6

type StackFrameSet struct {
	stackFrames []*StackFrame
	// pool retains popped frames for reuse. Push/pop is strictly LIFO and a
	// StackFrameSet is reused across all records (it lives on the persistent
	// runtime.State), so without pooling each record's block entry/exit would
	// allocate and discard a StackFrame (a slice + a map). Pooling makes
	// per-record block execution allocation-free after the first record.
	pool []*StackFrame
}

func newStackFrameSet() *StackFrameSet {
	stackFrames := make([]*StackFrame, 1, stackFrameSetInitCap)
	stackFrames[0] = newStackFrame()
	return &StackFrameSet{
		stackFrames: stackFrames,
	}
}

// reset returns a pooled frameset to its freshly-constructed state: exactly one
// (cleared) base frame. Any extra frames are kept in the per-frameset frame
// pool for reuse. At a balanced PopStackFrameSet the set is already at depth 1,
// so this is normally just a clear of the base frame.
func (frameset *StackFrameSet) reset() {
	for len(frameset.stackFrames) > 1 {
		frameset.popStackFrame()
	}
	frameset.stackFrames[0].clear()
}

func (frameset *StackFrameSet) pushStackFrame() {
	n := len(frameset.pool)
	if n > 0 {
		frame := frameset.pool[n-1]
		frameset.pool = frameset.pool[:n-1]
		frame.clear()
		frameset.stackFrames = append(frameset.stackFrames, frame)
	} else {
		frameset.stackFrames = append(frameset.stackFrames, newStackFrame())
	}
}

func (frameset *StackFrameSet) popStackFrame() {
	n := len(frameset.stackFrames)
	frame := frameset.stackFrames[n-1]
	frameset.stackFrames = frameset.stackFrames[0 : n-1]
	frameset.pool = append(frameset.pool, frame)
}

// Returns nil on no-such
func (frameset *StackFrameSet) get(
	stackVariable *StackVariable,
) *mlrval.Mlrval {
	// Scope-walk
	numStackFrames := len(frameset.stackFrames)
	for offset := numStackFrames - 1; offset >= 0; offset-- {
		stackFrame := frameset.stackFrames[offset]
		mlrval := stackFrame.get(stackVariable)
		if mlrval != nil {
			return mlrval
		}
	}
	return nil
}

// See Stack.DefineTypedAtScope comments above
func (frameset *StackFrameSet) defineTypedAtScope(
	stackVariable *StackVariable,
	typeName string,
	mlrval *mlrval.Mlrval,
) error {
	offset := len(frameset.stackFrames) - 1
	// TODO: comment
	return frameset.stackFrames[offset].defineTyped(
		stackVariable, typeName, mlrval,
	)
}

// See Stack.SetAtScope comments above
func (frameset *StackFrameSet) setAtScope(
	stackVariable *StackVariable,
	mlrval *mlrval.Mlrval,
) error {
	offset := len(frameset.stackFrames) - 1
	return frameset.stackFrames[offset].set(stackVariable, mlrval)
}

// See Stack.Set comments above
func (frameset *StackFrameSet) set(
	stackVariable *StackVariable,
	mlrval *mlrval.Mlrval,
) error {
	// Scope-walk
	numStackFrames := len(frameset.stackFrames)
	for offset := numStackFrames - 1; offset >= 0; offset-- {
		stackFrame := frameset.stackFrames[offset]
		if stackFrame.has(stackVariable) {
			return stackFrame.set(stackVariable, mlrval)
		}
	}
	return frameset.setAtScope(stackVariable, mlrval)
}

// See Stack.SetIndexed comments above
func (frameset *StackFrameSet) setIndexed(
	stackVariable *StackVariable,
	indices []*mlrval.Mlrval,
	mlrval *mlrval.Mlrval,
) error {
	// Scope-walk
	numStackFrames := len(frameset.stackFrames)
	for offset := numStackFrames - 1; offset >= 0; offset-- {
		stackFrame := frameset.stackFrames[offset]
		if stackFrame.has(stackVariable) {
			return stackFrame.setIndexed(stackVariable, indices, mlrval)
		}
	}
	offset := numStackFrames - 1
	return frameset.stackFrames[offset].setIndexed(stackVariable, indices, mlrval)
}

// See Stack.Unset comments above
func (frameset *StackFrameSet) unset(
	stackVariable *StackVariable,
) {
	// Scope-walk
	numStackFrames := len(frameset.stackFrames)
	for offset := numStackFrames - 1; offset >= 0; offset-- {
		stackFrame := frameset.stackFrames[offset]
		if stackFrame.has(stackVariable) {
			stackFrame.unset(stackVariable)
			return
		}
	}
}

// See Stack.UnsetIndexed comments above
func (frameset *StackFrameSet) unsetIndexed(
	stackVariable *StackVariable,
	indices []*mlrval.Mlrval,
) {
	// Scope-walk
	numStackFrames := len(frameset.stackFrames)
	for offset := numStackFrames - 1; offset >= 0; offset-- {
		stackFrame := frameset.stackFrames[offset]
		if stackFrame.has(stackVariable) {
			stackFrame.unsetIndexed(stackVariable, indices)
			return
		}
	}
}

// STACKFRAME METHODS

const stackFrameInitCap = 10

type StackFrame struct {
	// TODO: just a map for now. In the C impl, pre-computation of
	// name-to-array-slot indices was an important optimization, especially for
	// compute-intensive scenarios.
	//vars map[string]*types.TypeGatedMlrvalVariable

	// TODO: comment
	vars           []*types.TypeGatedMlrvalVariable
	namesToOffsets map[string]int
}

func newStackFrame() *StackFrame {
	vars := make([]*types.TypeGatedMlrvalVariable, 0, stackFrameInitCap)
	namesToOffsets := make(map[string]int)
	return &StackFrame{
		vars:           vars,
		namesToOffsets: namesToOffsets,
	}
}

// clear resets a frame for reuse from the pool, retaining its backing slice and
// map allocations. The vars elements are nilled so reuse does not pin the
// previous scope's variable values.
func (frame *StackFrame) clear() {
	for i := range frame.vars {
		frame.vars[i] = nil
	}
	frame.vars = frame.vars[:0]
	clear(frame.namesToOffsets)
}

// Returns nil on no such
func (frame *StackFrame) get(
	stackVariable *StackVariable,
) *mlrval.Mlrval {
	offset, ok := frame.namesToOffsets[stackVariable.name]
	if ok {
		return frame.vars[offset].GetValue()
	}
	return nil
}

func (frame *StackFrame) has(
	stackVariable *StackVariable,
) bool {
	_, ok := frame.namesToOffsets[stackVariable.name]
	return ok
}

// TODO: audit for honor of error-return at callsites
func (frame *StackFrame) set(
	stackVariable *StackVariable,
	mlrval *mlrval.Mlrval,
) error {
	offset, ok := frame.namesToOffsets[stackVariable.name]
	if !ok {
		slot, err := types.NewTypeGatedMlrvalVariable(stackVariable.name, "any", mlrval)
		if err != nil {
			return err
		}
		frame.vars = append(frame.vars, slot)
		offsetInFrame := len(frame.vars) - 1
		frame.namesToOffsets[stackVariable.name] = offsetInFrame
		return nil
	}
	return frame.vars[offset].Assign(mlrval)
}

// TODO: audit for honor of error-return at callsites
func (frame *StackFrame) defineTyped(
	stackVariable *StackVariable,
	typeName string,
	mlrval *mlrval.Mlrval,
) error {
	_, ok := frame.namesToOffsets[stackVariable.name]
	if !ok {
		slot, err := types.NewTypeGatedMlrvalVariable(stackVariable.name, typeName, mlrval)
		if err != nil {
			return err
		}
		frame.vars = append(frame.vars, slot)
		offsetInFrame := len(frame.vars) - 1
		frame.namesToOffsets[stackVariable.name] = offsetInFrame
		return nil
	}
	return fmt.Errorf(
		"%s: variable %s has already been defined in the same scope",
		"mlr", stackVariable.name,
	)
}

// TODO: audit for honor of error-return at callsites
func (frame *StackFrame) setIndexed(
	stackVariable *StackVariable,
	indices []*mlrval.Mlrval,
	mv *mlrval.Mlrval,
) error {
	value := frame.get(stackVariable)
	if value == nil {
		lib.InternalCodingErrorIf(len(indices) < 1)
		leadingIndex := indices[0]
		if leadingIndex.IsString() || leadingIndex.IsInt() {
			newval := mlrval.FromMap(mlrval.NewMlrmap())
			newval.PutIndexed(indices, mv)
			return frame.set(stackVariable, newval)
		}
		return fmt.Errorf(
			"%s: map indices must be int or string; got %s",
			"mlr", leadingIndex.GetTypeName(),
		)
	} else {
		// For example maybe the variable exists and is an array but the
		// leading index is a string.
		return value.PutIndexed(indices, mv)
	}
}

func (frame *StackFrame) unset(
	stackVariable *StackVariable,
) {
	offset, ok := frame.namesToOffsets[stackVariable.name]
	if ok {
		frame.vars[offset].Unassign()
	}
}

func (frame *StackFrame) unsetIndexed(
	stackVariable *StackVariable,
	indices []*mlrval.Mlrval,
) {
	value := frame.get(stackVariable)
	if value == nil {
		return
	}
	value.RemoveIndexed(indices)
}
