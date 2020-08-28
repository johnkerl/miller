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

package action

import (
	"fmt"
)

type Action interface {
	Equal(Action) bool
	ResolveConflict(that Action) Action
	String() string
}

type (
	Accept bool
	Error  bool
	Reduce int
	Shift  int
)

const (
	ACCEPT = Accept(true)
	ERROR  = Error(true)
)

func (Accept) Equal(act Action) bool {
	if _, ok := act.(Accept); ok {
		return true
	}
	return false
}

func (Error) Equal(act Action) bool {
	if _, ok := act.(Error); ok {
		return true
	}
	return false
}

func (this Reduce) Equal(act Action) bool {
	if that, ok := act.(Reduce); ok {
		return this == that
	}
	return false
}

func (this Shift) Equal(act Action) bool {
	if that, ok := act.(Shift); ok {
		return this == that
	}
	return false
}

func (this Accept) ResolveConflict(that Action) Action {
	if _, ok := that.(Error); ok {
		return this
	}
	panic(fmt.Sprintf("Cannot have LR1 conflict with Accept."))
}

func (Error) ResolveConflict(that Action) Action {
	return that
}

func (this Shift) ResolveConflict(that Action) Action {
	switch that := that.(type) {
	case Accept:
		panic(fmt.Sprintf("Impossible conflict: Shift(%d)/Accept", int(this)))
	case Shift:
		panic(fmt.Sprintf("Cannot have Shift(%d)/Shift(%d)", int(this), int(that)))
	case Error:
		return this
	case Reduce:
		return this
	}
	panic(fmt.Sprintf("Conflict with unknown type of action: %T", that))
}

func (this Reduce) ResolveConflict(that Action) Action {
	switch that := that.(type) {
	case Accept:
		panic(fmt.Sprintf("Impossible conflict: Shift(%d)/Accept", int(this)))
	case Shift:
		return that
	case Error:
		return this
	case Reduce:
		if this < that {
			return this
		}
		return that
	}
	panic(fmt.Sprintf("Conflict with unknown type of action: %T", that))
}

func (this Accept) String() string {
	return "accept"
}

func (this Error) String() string {
	return "error"
}

func (this Reduce) String() string {
	return fmt.Sprintf("Reduce(%d)", this)
}

func (this Shift) String() string {
	return fmt.Sprintf("Shift(%d)", this)
}
