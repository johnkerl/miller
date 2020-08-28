//Copyright 2013 Vastech SA (PTY) LTD
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package util

type Stack struct {
	stack []interface{}
}

func NewStack(capacity int) *Stack {
	return &Stack{make([]interface{}, 0, capacity)}
}

/*
Returns the number of items on the stack or 0 if the stack is empty
*/
func (this *Stack) Len() int {
	return len(this.stack)
}

/*
Returns the item at index in the stack. The index of the bottom of the stack is 0.
Returns nil if index is higher than stack.Len()-1
*/
func (this *Stack) Peek(index int) interface{} {
	if index > len(this.stack)-1 {
		return nil
	}
	return this.stack[index]
}

/*
Removes and returns the last item pushed or nil if the stack is empty
*/
func (this *Stack) Pop() (item interface{}) {
	if len(this.stack) == 0 {
		return nil
	}
	item = this.stack[len(this.stack)-1]
	this.stack = this.stack[:len(this.stack)-1]
	return
}

/*
Push a new item onto the stack. item may not be nil
*/
func (this *Stack) Push(items ...interface{}) *Stack {
	for _, item := range items {
		if item == nil {
			panic("nil item may not be pushed")
		}
		this.stack = append(this.stack, item)
	}
	return this
}

/*
Returns the top element without popping the stack
*/
func (this *Stack) Top() interface{} {
	return this.stack[len(this.stack)-1]
}
