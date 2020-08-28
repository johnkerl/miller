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

package ast

import (
	"fmt"
)

func walkLexNode(node LexNode, visitor LexNodeVisitor) LexNodeVisitor {
	switch n := node.(type) {
	case *LexAlt:
		return n.Walk(visitor)
	case *LexCharLit:
		return n.Walk(visitor)
	case *LexCharRange:
		return n.Walk(visitor)
	case *LexDot:
		return n.Walk(visitor)
	case *LexGroupPattern:
		return n.Walk(visitor)
	case *LexIgnoredTokDef:
		return n.Walk(visitor)
	case *LexImports:
		return n.Walk(visitor)
	case *LexOptPattern:
		return n.Walk(visitor)
	case *LexPattern:
		return n.Walk(visitor)
	case *LexProductions:
		return n.Walk(visitor)
	case *LexRegDef:
		return n.Walk(visitor)
	case *LexRegDefId:
		return n.Walk(visitor)
	case *LexRepPattern:
		return n.Walk(visitor)
	case *LexTokDef:
		return n.Walk(visitor)
	}
	panic(fmt.Sprintf("Unsupported LexNode type %T", node))
}

func (this *LexAlt) Walk(visitor LexNodeVisitor) LexNodeVisitor {
	for _, term := range this.Terms {
		if v := walkLexNode(term, visitor); v == nil {
			return nil
		}
	}
	return visitor.Visit(this)
}

func (this *LexCharLit) Walk(visitor LexNodeVisitor) LexNodeVisitor {
	return visitor.Visit(this)
}

func (this *LexCharRange) Walk(visitor LexNodeVisitor) LexNodeVisitor {
	return visitor.Visit(this)
}

func (this *LexDot) Walk(visitor LexNodeVisitor) LexNodeVisitor {
	return visitor.Visit(this)
}

func (this *LexGroupPattern) Walk(visitor LexNodeVisitor) LexNodeVisitor {
	for _, term := range this.LexPattern.Alternatives {
		if v := walkLexNode(term, visitor); v == nil {
			return nil
		}
	}
	return visitor.Visit(this)
}

func (this *LexIgnoredTokDef) Walk(visitor LexNodeVisitor) LexNodeVisitor {
	return visitor.Visit(this)
}

func (this *LexImports) Walk(visitor LexNodeVisitor) LexNodeVisitor {
	return visitor.Visit(this)
}

func (this *LexOptPattern) Walk(visitor LexNodeVisitor) LexNodeVisitor {
	for _, term := range this.LexPattern.Alternatives {
		if v := walkLexNode(term, visitor); v == nil {
			return nil
		}
	}
	return visitor.Visit(this)
}

func (this *LexPattern) Walk(visitor LexNodeVisitor) LexNodeVisitor {
	for _, term := range this.Alternatives {
		if v := walkLexNode(term, visitor); v == nil {
			return nil
		}
	}
	return visitor.Visit(this)
}

func (this *LexProductions) Walk(visitor LexNodeVisitor) LexNodeVisitor {
	for _, term := range this.Productions {
		if v := walkLexNode(term, visitor); v == nil {
			return nil
		}
	}
	return visitor.Visit(this)
}

func (this *LexRegDef) Walk(visitor LexNodeVisitor) LexNodeVisitor {
	return visitor.Visit(this)
}

func (this *LexRegDefId) Walk(visitor LexNodeVisitor) LexNodeVisitor {
	return visitor.Visit(this)
}

func (this *LexRepPattern) Walk(visitor LexNodeVisitor) LexNodeVisitor {
	for _, term := range this.LexPattern.Alternatives {
		if v := walkLexNode(term, visitor); v == nil {
			return nil
		}
	}
	return visitor.Visit(this)
}

func (this *LexTokDef) Walk(visitor LexNodeVisitor) LexNodeVisitor {
	return visitor.Visit(this)
}
