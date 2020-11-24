package types

import (
	"bytes"
	"strconv"
	"strings"

	"miller/lib"
)

// ================================================================
// Map/array count. Scalars (including strings) have length 1.
func MlrvalLength(ma *Mlrval) Mlrval {
	switch ma.mvtype {
	case MT_ERROR:
		return MlrvalFromInt64(0)
		break
	case MT_ABSENT:
		return MlrvalFromInt64(0)
		break
	case MT_ARRAY:
		return MlrvalFromInt64(int64(len(ma.arrayval)))
		break
	case MT_MAP:
		return MlrvalFromInt64(int64(ma.mapval.FieldCount))
		break
	}
	return MlrvalFromInt64(1)
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
	return MlrvalFromInt64(int64(1 + maxChildDepth))
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
	return MlrvalFromInt64(int64(1 + maxChildDepth))
}

func depth_from_scalar(ma *Mlrval) Mlrval {
	return MlrvalFromInt64(0)
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

		childLeafCount := MlrvalFromInt64(1)
		if child.mvtype == MT_ARRAY {
			childLeafCount = leafcount_from_array(&child)
		} else if child.mvtype == MT_MAP {
			childLeafCount = leafcount_from_map(&child)
		}

		iChildLeafCount := int(childLeafCount.intval)
		sumChildLeafCount += iChildLeafCount
	}
	return MlrvalFromInt64(int64(sumChildLeafCount))
}

func leafcount_from_map(ma *Mlrval) Mlrval {
	sumChildLeafCount := 0
	for pe := ma.mapval.Head; pe != nil; pe = pe.Next {
		child := pe.Value

		// Golang initialization loop if we do this :(
		// childLeafCount := MlrvalLeafCount(child)

		childLeafCount := MlrvalFromInt64(1)
		if child.mvtype == MT_ARRAY {
			childLeafCount = leafcount_from_array(child)
		} else if child.mvtype == MT_MAP {
			childLeafCount = leafcount_from_map(child)
		}

		iChildLeafCount := int(childLeafCount.intval)
		sumChildLeafCount += iChildLeafCount
	}
	return MlrvalFromInt64(int64(sumChildLeafCount))
}

func leafcount_from_scalar(ma *Mlrval) Mlrval {
	return MlrvalFromInt64(1)
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
	_, ok := unaliasArrayIndex(&ma.arrayval, mb.intval)
	return MlrvalFromBool(ok)
}

func has_key_in_map(ma, mb *Mlrval) Mlrval {
	if mb.mvtype == MT_INT {
		return MlrvalFromFalse()
	}
	if mb.mvtype != MT_STRING {
		return MlrvalFromError()
	}
	return MlrvalFromBool(ma.mapval.Has(&mb.printrep))
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
		oldKey := *pe.Key
		_, present := newKeys[oldKey]
		if present {
			newMap.PutCopy(&oldKey, oldmap.Get(&oldKey))
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
			newMap.Remove(&exceptArg.printrep)
		} else if exceptArg.mvtype == MT_ARRAY {
			for _, element := range exceptArg.arrayval {
				if element.mvtype == MT_STRING {
					newMap.Remove(&element.printrep)
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
			buffer.WriteString(*pe.Key)
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
			buffer.WriteString(*pe.Key)
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
			retval.mapval.PutReference(&key, &value)
		} else if len(pair) == 2 {
			key := pair[0]
			value := MlrvalFromInferredType(pair[1])
			retval.mapval.PutReference(&key, &value)
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
			retval.mapval.PutReference(&key, &value)
		} else if len(pair) == 2 {
			key := pair[0]
			value := MlrvalFromString(pair[1])
			retval.mapval.PutReference(&key, &value)
		} else {
			lib.InternalCodingErrorIf(true)
		}
	}

	return retval
}

// ----------------------------------------------------------------
// splitnv("a=3,b=4,c=5", "=", ",") -> {"1":3,"2":4,"3":5}
func MlrvalSplitNV(ma, mb, mc *Mlrval) Mlrval {
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
		key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
		if len(pair) == 1 {
			value := MlrvalFromInferredType(pair[0])
			retval.mapval.PutReference(&key, &value)
		} else if len(pair) == 2 {
			value := MlrvalFromInferredType(pair[1])
			retval.mapval.PutReference(&key, &value)
		} else {
			lib.InternalCodingErrorIf(true)
		}
	}

	return retval
}

// ----------------------------------------------------------------
// splitnvx("a=3,b=4,c=5", "=", ",") -> {"1":"3","2":"4","3":"5"}
func MlrvalSplitNVX(ma, mb, mc *Mlrval) Mlrval {
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
		key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
		if len(pair) == 1 {
			value := MlrvalFromString(pair[0])
			retval.mapval.PutReference(&key, &value)
		} else if len(pair) == 2 {
			value := MlrvalFromString(pair[1])
			retval.mapval.PutReference(&key, &value)
		} else {
			lib.InternalCodingErrorIf(true)
		}
	}

	return retval
}

// ----------------------------------------------------------------
// splitak("a=3,b=4,c=5", "=", ",") -> ["a","b","c"]
func MlrvalSplitAK(ma, mb, mc *Mlrval) Mlrval {
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

	fields := lib.SplitString(ma.printrep, fieldSeparator)

	retval := NewSizedMlrvalArray(int64(len(fields)))

	for i, field := range fields {
		pair := strings.SplitN(field, pairSeparator, 2)

		if len(pair) == 1 {
			key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
			retval.arrayval[i] = MlrvalFromString(key)
		} else if len(pair) == 2 {
			key := MlrvalFromInferredType(pair[0])
			retval.arrayval[i] = key
		} else {
			lib.InternalCodingErrorIf(true)
		}
	}

	return *retval
}

// ----------------------------------------------------------------
// splitav("a=3,b=4,c=5", "=", ",") -> [3,4,5]
func MlrvalSplitAV(ma, mb, mc *Mlrval) Mlrval {
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

	fields := lib.SplitString(ma.printrep, fieldSeparator)

	retval := NewSizedMlrvalArray(int64(len(fields)))

	for i, field := range fields {
		pair := strings.SplitN(field, pairSeparator, 2)
		if len(pair) == 1 {
			value := MlrvalFromInferredType(pair[0])
			retval.arrayval[i] = value
		} else if len(pair) == 2 {
			value := MlrvalFromInferredType(pair[1])
			retval.arrayval[i] = value
		} else {
			lib.InternalCodingErrorIf(true)
		}
	}

	return *retval
}

// ----------------------------------------------------------------
// splitav("a=3,b=4,c=5", "=", ",") -> ["3","4","5"]
func MlrvalSplitAVX(ma, mb, mc *Mlrval) Mlrval {
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

	fields := lib.SplitString(ma.printrep, fieldSeparator)

	retval := NewSizedMlrvalArray(int64(len(fields)))

	for i, field := range fields {
		pair := strings.SplitN(field, pairSeparator, 2)
		if len(pair) == 1 {
			value := MlrvalFromString(pair[0])
			retval.arrayval[i] = value
		} else if len(pair) == 2 {
			value := MlrvalFromString(pair[1])
			retval.arrayval[i] = value
		} else {
			lib.InternalCodingErrorIf(true)
		}
	}

	return *retval
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

	retval := NewSizedMlrvalArray(int64(len(fields)))

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
	fieldSeparator := mb.printrep

	fields := lib.SplitString(ma.printrep, fieldSeparator)

	retval := NewSizedMlrvalArray(int64(len(fields)))

	for i, field := range fields {
		value := MlrvalFromString(field)
		retval.arrayval[i] = value
	}

	return *retval
}

// ----------------------------------------------------------------
func MlrvalKeys(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_MAP {
		retval := NewSizedMlrvalArray(ma.mapval.FieldCount)
		i := 0
		for pe := ma.mapval.Head; pe != nil; pe = pe.Next {
			retval.arrayval[i] = MlrvalFromString(*pe.Key)
			i++
		}
		return *retval

	} else if ma.mvtype == MT_ARRAY {
		retval := NewSizedMlrvalArray(int64(len(ma.arrayval)))
		for i, _ := range ma.arrayval {
			retval.arrayval[i] = MlrvalFromInt64(int64(i + 1)) // Miller user-space indices are 1-up
		}
		return *retval

	} else {
		return MlrvalFromError()
	}
}

// ----------------------------------------------------------------
func MlrvalValues(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_MAP {
		retval := NewSizedMlrvalArray(ma.mapval.FieldCount)
		i := 0
		for pe := ma.mapval.Head; pe != nil; pe = pe.Next {
			retval.arrayval[i] = *pe.Value.Copy()
			i++
		}
		return *retval

	} else if ma.mvtype == MT_ARRAY {
		retval := NewSizedMlrvalArray(int64(len(ma.arrayval)))
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
