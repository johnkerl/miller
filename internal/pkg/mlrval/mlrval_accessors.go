package mlrval

import (
	"strconv"

	"github.com/johnkerl/miller/internal/pkg/lib"
)

func (mv *Mlrval) GetArrayLength() (int, bool) {
	if mv.IsArray() {
		return len(mv.intf.([]*Mlrval)), true
	} else {
		return -999, false
	}
}

func CopyMlrvalArray(input []*Mlrval) []*Mlrval {
	output := make([]*Mlrval, len(input))
	for i, element := range input {
		if element == nil {
			output[i] = nil
		} else {
			output[i] = element.Copy()
		}
	}
	return output
}

// ---------------------------------------------------------------
// For the flatten verb and DSL function.

func (mv *Mlrval) FlattenToMap(prefix string, delimiter string) Mlrval {
	retval := NewMlrmap()

	if mv.IsMap() {
		// Without this, the for-loop below is zero-pass and fields with "{}"
		// values would disappear entirely in a JSON-to-CSV conversion.
		if mv.intf.(*Mlrmap).IsEmpty() {
			if prefix != "" {
				retval.PutCopy(prefix, FromString("{}"))
			}
		}

		for pe := mv.intf.(*Mlrmap).Head; pe != nil; pe = pe.Next {
			nextPrefix := pe.Key
			if prefix != "" {
				nextPrefix = prefix + delimiter + nextPrefix
			}
			if pe.Value.IsMap() || pe.Value.IsArray() {
				nextResult := pe.Value.FlattenToMap(nextPrefix, delimiter)
				lib.InternalCodingErrorIf(nextResult.mvtype != MT_MAP)
				for pf := nextResult.intf.(*Mlrmap).Head; pf != nil; pf = pf.Next {
					retval.PutCopy(pf.Key, pf.Value.Copy())
				}
			} else {
				retval.PutCopy(nextPrefix, pe.Value.Copy())
			}
		}

	} else if mv.IsArray() {
		// Without this, the for-loop below is zero-pass and fields with "[]"
		// values would disappear entirely in a JSON-to-CSV conversion.
		if len(mv.intf.([]*Mlrval)) == 0 {
			if prefix != "" {
				retval.PutCopy(prefix, FromString("[]"))
			}
		}

		for zindex, value := range mv.intf.([]*Mlrval) {
			nextPrefix := strconv.Itoa(zindex + 1) // Miller user-space indices are 1-up
			if prefix != "" {
				nextPrefix = prefix + delimiter + nextPrefix
			}
			if value.IsMap() || value.IsArray() {
				nextResult := value.FlattenToMap(nextPrefix, delimiter)
				lib.InternalCodingErrorIf(nextResult.mvtype != MT_MAP)
				for pf := nextResult.intf.(*Mlrmap).Head; pf != nil; pf = pf.Next {
					retval.PutCopy(pf.Key, pf.Value.Copy())
				}
			} else {
				retval.PutCopy(nextPrefix, value.Copy())
			}
		}

	} else {
		retval.PutCopy(prefix, mv.Copy())
	}

	return *FromMap(retval)
}

// Increment is used by stats1.
func (mv *Mlrval) Increment() {
	if mv.mvtype == MT_INT {
		mv.intval++
	} else if mv.mvtype == MT_FLOAT {
		mv.intf = mv.intf.(float64) + 1.0
	}
}
