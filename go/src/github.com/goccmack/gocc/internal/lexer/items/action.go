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

package items

import (
	"fmt"
)

type (
	Action interface {
		action()
		String() string
	}
	Accept string
	Ignore string
	Error  int
)

func (Accept) action() {}
func (Error) action()  {}
func (Ignore) action() {}

func (this Accept) String() string {
	return fmt.Sprintf("Accept(\"%s\")", string(this))
}

func (this Ignore) String() string {
	return fmt.Sprintf("Ignore(\"%s\")", string(this))
}

func (this Error) String() string {
	return fmt.Sprintf("Error(%d)", int(this))
}
