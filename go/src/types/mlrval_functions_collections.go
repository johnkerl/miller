package types

import (
	"bytes"
	"strconv"
	"strings"

	"miller/src/lib"
)

// ================================================================
// Map/array count. Scalars (including strings) have length 1.
func MlrvalLength(ma *Mlrval) Mlrval {
	switch ma.mvtype {
	case MT_ERROR:
		return MlrvalFromInt(0)
		break
	case MT_ABSENT:
		return MlrvalFromInt(0)
		break
	case MT_ARRAY:
		return MlrvalFromInt(int(len(ma.arrayval)))
		break
	case MT_MAP:
		return MlrvalFromInt(int(ma.mapval.FieldCount))
		break
	}
	return MlrvalFromInt(1)
}

// ================================================================
func depth_from_array(ma *Mlrval) Mlrval {
	maxChildDepth := 0
	for _, child := range ma.arrayval {
		childDepth := MlrvalDepth(&child)
		iChildDepth := int(childDepth.intval)
		if iChildDepth > maxChildDepth {
			maxChildDepth = iChildDepth
		}
	}
	return MlrvalFromInt(int(1 + maxChildDepth))
}

func depth_from_map(ma *Mlrval) Mlrval {
	maxChildDepth := 0
	for pe := ma.mapval.Head; pe != nil; pe = pe.Next {
		child := pe.Value
		childDepth := MlrvalDepth(child)
		iChildDepth := int(childDepth.intval)
		if iChildDepth > maxChildDepth {
			maxChildDepth = iChildDepth
		}
	}
	return MlrvalFromInt(int(1 + maxChildDepth))
}

func depth_from_scalar(ma *Mlrval) Mlrval {
	return MlrvalFromInt(0)
}

// We get a Golang "initialization loop" due to recursive depth computation
// if this is defined statically. So, we use a "package init" function.
var depth_dispositions = [MT_DIM]UnaryFunc{}

func init() {
	depth_dispositions = [MT_DIM]UnaryFunc{
		/*ERROR  */ _erro1,
		/*ABSENT */ _absn1,
		/*VOID   */ depth_from_scalar,
		/*STRING */ depth_from_scalar,
		/*INT    */ depth_from_scalar,
		/*FLOAT  */ depth_from_scalar,
		/*BOOL   */ depth_from_scalar,
		/*ARRAY  */ depth_from_array,
		/*MAP    */ depth_from_map,
	}
}

func MlrvalDepth(ma *Mlrval) Mlrval {
	return depth_dispositions[ma.mvtype](ma)
}

// ================================================================
func leafcount_from_array(ma *Mlrval) Mlrval {
	sumChildLeafCount := 0
	for _, child := range ma.arrayval {
		// Golang initialization loop if we do this :(
		// childLeafCount := MlrvalLeafCount(&child)

		childLeafCount := MlrvalFromInt(1)
		if child.mvtype == MT_ARRAY {
			childLeafCount = leafcount_from_array(&child)
		} else if child.mvtype == MT_MAP {
			childLeafCount = leafcount_from_map(&child)
		}

		iChildLeafCount := int(childLeafCount.intval)
		sumChildLeafCount += iChildLeafCount
	}
	return MlrvalFromInt(int(sumChildLeafCount))
}

func leafcount_from_map(ma *Mlrval) Mlrval {
	sumChildLeafCount := 0
	for pe := ma.mapval.Head; pe != nil; pe = pe.Next {
		child := pe.Value

		// Golang initialization loop if we do this :(
		// childLeafCount := MlrvalLeafCount(child)

		childLeafCount := MlrvalFromInt(1)
		if child.mvtype == MT_ARRAY {
			childLeafCount = leafcount_from_array(child)
		} else if child.mvtype == MT_MAP {
			childLeafCount = leafcount_from_map(child)
		}

		iChildLeafCount := int(childLeafCount.intval)
		sumChildLeafCount += iChildLeafCount
	}
	return MlrvalFromInt(int(sumChildLeafCount))
}

func leafcount_from_scalar(ma *Mlrval) Mlrval {
	return MlrvalFromInt(1)
}

var leafcount_dispositions = [MT_DIM]UnaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*VOID   */ leafcount_from_scalar,
	/*STRING */ leafcount_from_scalar,
	/*INT    */ leafcount_from_scalar,
	/*FLOAT  */ leafcount_from_scalar,
	/*BOOL   */ leafcount_from_scalar,
	/*ARRAY  */ leafcount_from_array,
	/*MAP    */ leafcount_from_map,
}

func MlrvalLeafCount(ma *Mlrval) Mlrval {
	return leafcount_dispositions[ma.mvtype](ma)
}

// ----------------------------------------------------------------
func has_key_in_array(ma, mb *Mlrval) Mlrval {
	if mb.mvtype == MT_STRING {
		return MlrvalFromFalse()
	}
	if mb.mvtype != MT_INT {
		return MlrvalFromError()
	}
	_, ok := UnaliasArrayIndex(&ma.arrayval, mb.intval)
	return MlrvalFromBool(ok)
}

func has_key_in_map(ma, mb *Mlrval) Mlrval {
	if mb.mvtype == MT_STRING || mb.mvtype == MT_INT {
		return MlrvalFromBool(ma.mapval.Has(mb.String()))
	} else {
		return MlrvalFromError()
	}
}

func MlrvalHasKey(ma, mb *Mlrval) Mlrval {
	if ma.mvtype == MT_ARRAY {
		return has_key_in_array(ma, mb)
	} else if ma.mvtype == MT_MAP {
		return has_key_in_map(ma, mb)
	} else {
		return MlrvalFromError()
	}
}

// ================================================================
func MlrvalMapSelect(mlrvals []*Mlrval) Mlrval {
	if len(mlrvals) < 1 {
		return MlrvalFromError()
	}
	if mlrvals[0].mvtype != MT_MAP {
		return MlrvalFromError()
	}
	oldmap := mlrvals[0].mapval
	newMap := NewMlrmap()

	newKeys := make(map[string]bool)
	for _, selectArg := range mlrvals[1:] {
		if selectArg.mvtype == MT_STRING {
			newKeys[selectArg.printrep] = true
		} else if selectArg.mvtype == MT_ARRAY {
			for _, element := range selectArg.arrayval {
				if element.mvtype == MT_STRING {
					newKeys[element.printrep] = true
				} else {
					return MlrvalFromError()
				}
			}
		} else {
			return MlrvalFromError()
		}
	}

	for pe := oldmap.Head; pe != nil; pe = pe.Next {
		oldKey := pe.Key
		_, present := newKeys[oldKey]
		if present {
			newMap.PutCopy(oldKey, oldmap.Get(oldKey))
		}
	}

	return MlrvalFromMap(newMap)
}

// ----------------------------------------------------------------
func MlrvalMapExcept(mlrvals []*Mlrval) Mlrval {
	if len(mlrvals) < 1 {
		return MlrvalFromError()
	}
	if mlrvals[0].mvtype != MT_MAP {
		return MlrvalFromError()
	}
	newMap := mlrvals[0].mapval.Copy()

	for _, exceptArg := range mlrvals[1:] {
		if exceptArg.mvtype == MT_STRING {
			newMap.Remove(exceptArg.printrep)
		} else if exceptArg.mvtype == MT_ARRAY {
			for _, element := range exceptArg.arrayval {
				if element.mvtype == MT_STRING {
					newMap.Remove(element.printrep)
				} else {
					return MlrvalFromError()
				}
			}
		} else {
			return MlrvalFromError()
		}
	}

	return MlrvalFromMap(newMap)
}

// ----------------------------------------------------------------
func MlrvalMapSum(mlrvals []*Mlrval) Mlrval {
	if len(mlrvals) == 0 {
		return MlrvalEmptyMap()
	}
	if len(mlrvals) == 1 {
		return *mlrvals[0]
	}
	if mlrvals[0].mvtype != MT_MAP {
		return MlrvalFromError()
	}
	newMap := mlrvals[0].mapval.Copy()

	for _, otherMapArg := range mlrvals[1:] {
		if otherMapArg.mvtype != MT_MAP {
			return MlrvalFromError()
		}

		for pe := otherMapArg.mapval.Head; pe != nil; pe = pe.Next {
			newMap.PutCopy(pe.Key, pe.Value)
		}
	}

	return MlrvalFromMap(newMap)
}

// ----------------------------------------------------------------
func MlrvalMapDiff(mlrvals []*Mlrval) Mlrval {
	if len(mlrvals) == 0 {
		return MlrvalEmptyMap()
	}
	if len(mlrvals) == 1 {
		return *mlrvals[0]
	}
	if mlrvals[0].mvtype != MT_MAP {
		return MlrvalFromError()
	}
	newMap := mlrvals[0].mapval.Copy()

	for _, otherMapArg := range mlrvals[1:] {
		if otherMapArg.mvtype != MT_MAP {
			return MlrvalFromError()
		}

		for pe := otherMapArg.mapval.Head; pe != nil; pe = pe.Next {
			newMap.Remove(pe.Key)
		}
	}

	return MlrvalFromMap(newMap)
}

// ================================================================
// joink([1,2,3], ",") -> "1,2,3"
// joink({"a":3,"b":4,"c":5}, ",") -> "a,b,c"
func MlrvalJoinK(ma, mb *Mlrval) Mlrval {
	if mb.mvtype != MT_STRING {
		return MlrvalFromError()
	}
	fieldSeparator := mb.printrep
	if ma.mvtype == MT_MAP {
		var buffer bytes.Buffer

		for pe := ma.mapval.Head; pe != nil; pe = pe.Next {
			buffer.WriteString(pe.Key)
			if pe.Next != nil {
				buffer.WriteString(fieldSeparator)
			}
		}

		return MlrvalFromString(buffer.String())
	} else if ma.mvtype == MT_ARRAY {
		var buffer bytes.Buffer

		for i, _ := range ma.arrayval {
			if i > 0 {
				buffer.WriteString(fieldSeparator)
			}
			// Miller userspace array indices are 1-up
			buffer.WriteString(strconv.Itoa(i + 1))
		}

		return MlrvalFromString(buffer.String())
	} else {
		return MlrvalFromError()
	}
}

// ----------------------------------------------------------------
// joinv([3,4,5], ",") -> "3,4,5"
// joinv({"a":3,"b":4,"c":5}, ",") -> "3,4,5"
func MlrvalJoinV(ma, mb *Mlrval) Mlrval {
	if mb.mvtype != MT_STRING {
		return MlrvalFromError()
	}
	fieldSeparator := mb.printrep

	if ma.mvtype == MT_MAP {
		var buffer bytes.Buffer

		for pe := ma.mapval.Head; pe != nil; pe = pe.Next {
			buffer.WriteString(pe.Value.String())
			if pe.Next != nil {
				buffer.WriteString(fieldSeparator)
			}
		}

		return MlrvalFromString(buffer.String())
	} else if ma.mvtype == MT_ARRAY {
		var buffer bytes.Buffer

		for i, element := range ma.arrayval {
			if i > 0 {
				buffer.WriteString(fieldSeparator)
			}
			buffer.WriteString(element.String())
		}

		return MlrvalFromString(buffer.String())
	} else {
		return MlrvalFromError()
	}
}

// ----------------------------------------------------------------
// joinkv([3,4,5], "=", ",") -> "1=3,2=4,3=5"
// joinkv({"a":3,"b":4,"c":5}, "=", ",") -> "a=3,b=4,c=5"
func MlrvalJoinKV(ma, mb, mc *Mlrval) Mlrval {
	if mb.mvtype != MT_STRING {
		return MlrvalFromError()
	}
	pairSeparator := mb.printrep
	if mc.mvtype != MT_STRING {
		return MlrvalFromError()
	}
	fieldSeparator := mc.printrep

	if ma.mvtype == MT_MAP {
		var buffer bytes.Buffer

		for pe := ma.mapval.Head; pe != nil; pe = pe.Next {
			buffer.WriteString(pe.Key)
			buffer.WriteString(pairSeparator)
			buffer.WriteString(pe.Value.String())
			if pe.Next != nil {
				buffer.WriteString(fieldSeparator)
			}
		}

		return MlrvalFromString(buffer.String())
	} else if ma.mvtype == MT_ARRAY {
		var buffer bytes.Buffer

		for i, element := range ma.arrayval {
			if i > 0 {
				buffer.WriteString(fieldSeparator)
			}
			// Miller userspace array indices are 1-up
			buffer.WriteString(strconv.Itoa(i + 1))
			buffer.WriteString(pairSeparator)
			buffer.WriteString(element.String())
		}

		return MlrvalFromString(buffer.String())
	} else {
		return MlrvalFromError()
	}
}

// ================================================================
// splitkv("a=3,b=4,c=5", "=", ",") -> {"a":3,"b":4,"c":5}
func MlrvalSplitKV(ma, mb, mc *Mlrval) Mlrval {
	if ma.mvtype != MT_STRING {
		return MlrvalFromError()
	}
	if mb.mvtype != MT_STRING {
		return MlrvalFromError()
	}
	pairSeparator := mb.printrep
	if mc.mvtype != MT_STRING {
		return MlrvalFromError()
	}
	fieldSeparator := mc.printrep

	retval := MlrvalEmptyMap()

	fields := lib.SplitString(ma.printrep, fieldSeparator)
	for i, field := range fields {
		pair := strings.SplitN(field, pairSeparator, 2)
		if len(pair) == 1 {
			key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
			value := MlrvalFromInferredType(pair[0])
			retval.mapval.PutReference(key, &value)
		} else if len(pair) == 2 {
			key := pair[0]
			value := MlrvalFromInferredType(pair[1])
			retval.mapval.PutReference(key, &value)
		} else {
			lib.InternalCodingErrorIf(true)
		}
	}

	return retval
}

// ----------------------------------------------------------------
// splitkvx("a=3,b=4,c=5", "=", ",") -> {"a":"3","b":"4","c":"5"}
func MlrvalSplitKVX(ma, mb, mc *Mlrval) Mlrval {
	if ma.mvtype != MT_STRING {
		return MlrvalFromError()
	}
	if mb.mvtype != MT_STRING {
		return MlrvalFromError()
	}
	pairSeparator := mb.printrep
	if mc.mvtype != MT_STRING {
		return MlrvalFromError()
	}
	fieldSeparator := mc.printrep

	retval := MlrvalEmptyMap()

	fields := lib.SplitString(ma.printrep, fieldSeparator)
	for i, field := range fields {
		pair := strings.SplitN(field, pairSeparator, 2)
		if len(pair) == 1 {
			key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
			value := MlrvalFromString(pair[0])
			retval.mapval.PutReference(key, &value)
		} else if len(pair) == 2 {
			key := pair[0]
			value := MlrvalFromString(pair[1])
			retval.mapval.PutReference(key, &value)
		} else {
			lib.InternalCodingErrorIf(true)
		}
	}

	return retval
}

// ----------------------------------------------------------------
// splitnv("a,b,c", ",") -> {"1":"a","2":"b","3":"c"}
func MlrvalSplitNV(ma, mb *Mlrval) Mlrval {
	if ma.mvtype != MT_STRING {
		return MlrvalFromError()
	}
	if mb.mvtype != MT_STRING {
		return MlrvalFromError()
	}

	retval := MlrvalEmptyMap()

	fields := lib.SplitString(ma.printrep, mb.printrep)
	for i, field := range fields {
		key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
		value := MlrvalFromInferredType(field)
		retval.mapval.PutReference(key, &value)
	}

	return retval
}

// ----------------------------------------------------------------
// splitnvx("3,4,5", ",") -> {"1":"3","2":"4","3":"5"}
func MlrvalSplitNVX(ma, mb *Mlrval) Mlrval {
	if ma.mvtype != MT_STRING {
		return MlrvalFromError()
	}
	if mb.mvtype != MT_STRING {
		return MlrvalFromError()
	}

	retval := MlrvalEmptyMap()

	fields := lib.SplitString(ma.printrep, mb.printrep)
	for i, field := range fields {
		key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
		value := MlrvalFromString(field)
		retval.mapval.PutReference(key, &value)
	}

	return retval
}

// ----------------------------------------------------------------
// splita("3,4,5", ",") -> [3,4,5]
func MlrvalSplitA(ma, mb *Mlrval) Mlrval {
	if ma.mvtype != MT_STRING {
		return MlrvalFromError()
	}
	if mb.mvtype != MT_STRING {
		return MlrvalFromError()
	}
	fieldSeparator := mb.printrep

	fields := lib.SplitString(ma.printrep, fieldSeparator)

	retval := NewSizedMlrvalArray(int(len(fields)))

	for i, field := range fields {
		value := MlrvalFromInferredType(field)
		retval.arrayval[i] = value
	}

	return *retval
}

// ----------------------------------------------------------------
// splitax("3,4,5", ",") -> ["3","4","5"]

func MlrvalSplitAX(ma, mb *Mlrval) Mlrval {
	if ma.mvtype != MT_STRING {
		return MlrvalFromError()
	}
	if mb.mvtype != MT_STRING {
		return MlrvalFromError()
	}
	input := ma.printrep
	fieldSeparator := mb.printrep

	return *mlrvalSplitAXHelper(input, fieldSeparator)
}

// Split out for MlrvalSplitAX and MlrvalUnflatten
func mlrvalSplitAXHelper(input string, separator string) *Mlrval {
	fields := lib.SplitString(input, separator)

	retval := NewSizedMlrvalArray(int(len(fields)))

	for i, field := range fields {
		value := MlrvalFromString(field)
		retval.arrayval[i] = value
	}

	return retval
}

// ----------------------------------------------------------------
func MlrvalGetKeys(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_MAP {
		retval := NewSizedMlrvalArray(ma.mapval.FieldCount)
		i := 0
		for pe := ma.mapval.Head; pe != nil; pe = pe.Next {
			retval.arrayval[i] = MlrvalFromString(pe.Key)
			i++
		}
		return *retval

	} else if ma.mvtype == MT_ARRAY {
		retval := NewSizedMlrvalArray(int(len(ma.arrayval)))
		for i, _ := range ma.arrayval {
			retval.arrayval[i] = MlrvalFromInt(int(i + 1)) // Miller user-space indices are 1-up
		}
		return *retval

	} else {
		return MlrvalFromError()
	}
}

// ----------------------------------------------------------------
func MlrvalGetValues(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_MAP {
		retval := NewSizedMlrvalArray(ma.mapval.FieldCount)
		i := 0
		for pe := ma.mapval.Head; pe != nil; pe = pe.Next {
			retval.arrayval[i] = *pe.Value.Copy()
			i++
		}
		return *retval

	} else if ma.mvtype == MT_ARRAY {
		retval := NewSizedMlrvalArray(int(len(ma.arrayval)))
		for i, value := range ma.arrayval {
			retval.arrayval[i] = *value.Copy()
		}
		return *retval

	} else {
		return MlrvalFromError()
	}
}

// ----------------------------------------------------------------
func MlrvalAppend(ma, mb *Mlrval) Mlrval {
	if ma.mvtype != MT_ARRAY {
		return MlrvalFromError()
	}

	macopy := ma.Copy()
	mbcopy := mb.Copy()
	macopy.ArrayAppend(mbcopy)
	return *macopy
}

// ----------------------------------------------------------------
// First argumemnt is prefix.
// Second argument is delimiter.
// Third argument is map or array.
// flatten("a", ".", {"b": { "c": 4 }}) is {"a.b.c" : 4}.
// flatten("", ".", {"a": { "b": 3 }}) is {"a.b" : 3}.
func MlrvalFlatten(ma, mb, mc *Mlrval) Mlrval {
	if mc.mvtype == MT_MAP || mc.mvtype == MT_ARRAY {
		if ma.mvtype != MT_STRING && ma.mvtype != MT_VOID {
			return MlrvalFromError()
		}
		prefix := ma.printrep
		if mb.mvtype != MT_STRING {
			return MlrvalFromError()
		}
		delimiter := mb.printrep

		return mc.FlattenToMap(prefix, delimiter)
	} else {
		return *mc
	}
}

// flatten($*, ".") is the same as flatten("", ".", $*)
func MlrvalFlattenBinary(ma, mb *Mlrval) Mlrval {
	return MlrvalFlatten(MlrvalPointerFromVoid(), mb, ma)
}

// ----------------------------------------------------------------
// First argument is a map.
// Second argument is a delimiter string.
// unflatten({"a.b.c", ".") is {"a": { "b": { "c": 4}}}.
func MlrvalUnflatten(ma, mb *Mlrval) Mlrval {
	if mb.mvtype != MT_STRING {
		return MlrvalFromError()
	}
	if ma.mvtype != MT_MAP {
		return *ma
	}
	oldmap := ma.mapval
	separator := mb.printrep
	newmap := NewMlrmap()

	for pe := oldmap.Head; pe != nil; pe = pe.Next {
		// TODO: factor out a shared helper function bewteen here and MlrvalSplitAX.
		arrayOfIndices := mlrvalSplitAXHelper(pe.Key, separator)
		newmap.PutIndexed(MakePointerArray(arrayOfIndices.arrayval), pe.Value.Copy())
	}
	return MlrvalFromMapReferenced(newmap)
}

// ----------------------------------------------------------------
// Converts maps with "1", "2", ... keys into arrays. Recurses nested data structures.
func MlrvalArrayify(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_MAP {
		convertible := true
		i := 0
		for pe := ma.mapval.Head; pe != nil; pe = pe.Next {
			sval := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
			i++
			if pe.Key != sval {
				convertible = false
			}
			temp := MlrvalArrayify(pe.Value)
			pe.Value = &temp
		}

		if convertible {
			arrayval := make([]Mlrval, ma.mapval.FieldCount)
			i := 0
			for pe := ma.mapval.Head; pe != nil; pe = pe.Next {
				arrayval[i] = *pe.Value.Copy()
				i++
			}
			return MlrvalFromArrayLiteralReference(arrayval)

		} else {
			return *ma
		}

	} else if ma.mvtype == MT_ARRAY {
		for i, _ := range ma.arrayval {
			ma.arrayval[i] = MlrvalArrayify(&ma.arrayval[i])
		}
		return *ma

	} else {
		return *ma
	}
}

// ----------------------------------------------------------------
func MlrvalJSONParse(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_VOID {
		return *ma
	} else if ma.mvtype != MT_STRING {
		return MlrvalFromError()
	} else {
		output := MlrvalFromPending()
		err := output.UnmarshalJSON([]byte(ma.printrep))
		if err == nil {
			return output
		} else {
			return MlrvalFromError()
		}
	}
}

func MlrvalJSONStringifyUnary(ma *Mlrval) Mlrval {
	outputBytes, err := ma.MarshalJSON(JSON_SINGLE_LINE)
	if err != nil {
		return MlrvalFromError()
	} else {
		return MlrvalFromString(string(outputBytes))
	}
}

func MlrvalJSONStringifyBinary(ma, mb *Mlrval) Mlrval {
	var jsonFormatting TJSONFormatting = JSON_SINGLE_LINE
	useMultiline, ok := mb.GetBoolValue()
	if !ok {
		return MlrvalFromError()
	}
	if useMultiline {
		jsonFormatting = JSON_MULTILINE
	}

	outputBytes, err := ma.MarshalJSON(jsonFormatting)
	if err != nil {
		return MlrvalFromError()
	} else {
		return MlrvalFromString(string(outputBytes))
	}
}
