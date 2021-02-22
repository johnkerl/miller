package types

import (
	"bytes"
	"strconv"
	"strings"

	"miller/src/lib"
)

// ================================================================
// Map/array count. Scalars (including strings) have length 1.
func MlrvalLength(output, input1 *Mlrval) {
	switch input1.mvtype {
	case MT_ERROR:
		output.SetFromInt(0)
		break
	case MT_ABSENT:
		output.SetFromInt(0)
		break
	case MT_ARRAY:
		output.SetFromInt(int(len(input1.arrayval)))
		break
	case MT_MAP:
		output.SetFromInt(int(input1.mapval.FieldCount))
		break
	default:
		output.SetFromInt(1)
		break
	}
}

// ================================================================
func depth_from_array(output, input1 *Mlrval) {
	maxChildDepth := 0
	for _, child := range input1.arrayval {
		childDepth := MlrvalFromAbsent()
		MlrvalDepth(&childDepth, &child)
		lib.InternalCodingErrorIf(childDepth.mvtype != MT_INT)
		iChildDepth := int(childDepth.intval)
		if iChildDepth > maxChildDepth {
			maxChildDepth = iChildDepth
		}
	}
	output.SetFromInt(int(1 + maxChildDepth))
}

func depth_from_map(output, input1 *Mlrval) {
	maxChildDepth := 0
	for pe := input1.mapval.Head; pe != nil; pe = pe.Next {
		child := pe.Value
		childDepth := MlrvalFromAbsent()
		MlrvalDepth(&childDepth, child)
		lib.InternalCodingErrorIf(childDepth.mvtype != MT_INT)
		iChildDepth := int(childDepth.intval)
		if iChildDepth > maxChildDepth {
			maxChildDepth = iChildDepth
		}
	}
	output.SetFromInt(int(1 + maxChildDepth))
}

func depth_from_scalar(output, input1 *Mlrval) {
	output.SetFromInt(0)
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

func MlrvalDepth(output, input1 *Mlrval) {
	depth_dispositions[input1.mvtype](output, input1)
}

// ================================================================
func leafcount_from_array(output, input1 *Mlrval) {
	sumChildLeafCount := 0
	for _, child := range input1.arrayval {
		// Golang initialization loop if we do this :(
		// childLeafCount := MlrvalLeafCount(&child)

		childLeafCount := MlrvalFromInt(1)
		if child.mvtype == MT_ARRAY {
			leafcount_from_array(&childLeafCount, &child)
		} else if child.mvtype == MT_MAP {
			leafcount_from_map(&childLeafCount, &child)
		}

		lib.InternalCodingErrorIf(childLeafCount.mvtype != MT_INT)
		iChildLeafCount := int(childLeafCount.intval)
		sumChildLeafCount += iChildLeafCount
	}
	output.SetFromInt(int(sumChildLeafCount))
}

func leafcount_from_map(output, input1 *Mlrval) {
	sumChildLeafCount := 0
	for pe := input1.mapval.Head; pe != nil; pe = pe.Next {
		child := pe.Value

		// Golang initialization loop if we do this :(
		// childLeafCount := MlrvalLeafCount(child)

		childLeafCount := MlrvalFromInt(1)
		if child.mvtype == MT_ARRAY {
			leafcount_from_array(&childLeafCount, child)
		} else if child.mvtype == MT_MAP {
			leafcount_from_map(&childLeafCount, child)
		}

		lib.InternalCodingErrorIf(childLeafCount.mvtype != MT_INT)
		iChildLeafCount := int(childLeafCount.intval)
		sumChildLeafCount += iChildLeafCount
	}
	output.SetFromInt(int(sumChildLeafCount))
}

func leafcount_from_scalar(output, input1 *Mlrval) {
	output.SetFromInt(1)
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

func MlrvalLeafCount(output, input1 *Mlrval) {
	leafcount_dispositions[input1.mvtype](output, input1)
}

// ----------------------------------------------------------------
func has_key_in_array(output, input1, input2 *Mlrval) {
	if input2.mvtype == MT_STRING {
		output.SetFromFalse()
		return
	}
	if input2.mvtype != MT_INT {
		output.SetFromError()
		return
	}
	_, ok := UnaliasArrayIndex(&input1.arrayval, input2.intval)
	output.SetFromBool(ok)
}

func has_key_in_map(output, input1, input2 *Mlrval) {
	if input2.mvtype == MT_STRING || input2.mvtype == MT_INT {
		output.SetFromBool(input1.mapval.Has(input2.String()))
	} else {
		output.SetFromError()
	}
}

func MlrvalHasKey(output, input1, input2 *Mlrval) {
	if input1.mvtype == MT_ARRAY {
		has_key_in_array(output, input1, input2)
	} else if input1.mvtype == MT_MAP {
		has_key_in_map(output, input1, input2)
	} else {
		output.SetFromError()
	}
}

// ================================================================
func MlrvalMapSelect(output *Mlrval, mlrvals []*Mlrval) {
	if len(mlrvals) < 1 {
		output.SetFromError()
		return
	}
	if mlrvals[0].mvtype != MT_MAP {
		output.SetFromError()
		return
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
					output.SetFromError()
					return
				}
			}
		} else {
			output.SetFromError()
			return
		}
	}

	for pe := oldmap.Head; pe != nil; pe = pe.Next {
		oldKey := pe.Key
		_, present := newKeys[oldKey]
		if present {
			newMap.PutCopy(oldKey, oldmap.Get(oldKey))
		}
	}

	output.SetFromMap(newMap)
}

// ----------------------------------------------------------------
func MlrvalMapExcept(output *Mlrval, mlrvals []*Mlrval) {
	if len(mlrvals) < 1 {
		output.SetFromError()
		return
	}
	if mlrvals[0].mvtype != MT_MAP {
		output.SetFromError()
		return
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
					output.SetFromError()
					return
				}
			}
		} else {
			output.SetFromError()
			return
		}
	}

	output.SetFromMap(newMap)
}

// ----------------------------------------------------------------
func MlrvalMapSum(output *Mlrval, mlrvals []*Mlrval) {
	if len(mlrvals) == 0 {
		output.SetFromEmptyMap()
		return
	}
	if len(mlrvals) == 1 {
		output.CopyFrom(mlrvals[0])
		return
	}
	if mlrvals[0].mvtype != MT_MAP {
		output.SetFromError()
		return
	}
	newMap := mlrvals[0].mapval.Copy()

	for _, otherMapArg := range mlrvals[1:] {
		if otherMapArg.mvtype != MT_MAP {
			output.SetFromError()
			return
		}

		for pe := otherMapArg.mapval.Head; pe != nil; pe = pe.Next {
			newMap.PutCopy(pe.Key, pe.Value)
		}
	}

	output.SetFromMap(newMap)
}

// ----------------------------------------------------------------
func MlrvalMapDiff(output *Mlrval, mlrvals []*Mlrval) {
	if len(mlrvals) == 0 {
		output.SetFromEmptyMap()
		return
	}
	if len(mlrvals) == 1 {
		output.CopyFrom(mlrvals[0])
		return
	}
	if mlrvals[0].mvtype != MT_MAP {
		output.SetFromError()
		return
	}
	newMap := mlrvals[0].mapval.Copy()

	for _, otherMapArg := range mlrvals[1:] {
		if otherMapArg.mvtype != MT_MAP {
			output.SetFromError()
			return
		}

		for pe := otherMapArg.mapval.Head; pe != nil; pe = pe.Next {
			newMap.Remove(pe.Key)
		}
	}

	output.SetFromMap(newMap)
}

// ================================================================
// joink([1,2,3], ",") -> "1,2,3"
// joink({"a":3,"b":4,"c":5}, ",") -> "a,b,c"
func MlrvalJoinK(output, input1, input2 *Mlrval) {
	if input2.mvtype != MT_STRING {
		output.SetFromError()
		return
	}
	fieldSeparator := input2.printrep
	if input1.mvtype == MT_MAP {
		var buffer bytes.Buffer

		for pe := input1.mapval.Head; pe != nil; pe = pe.Next {
			buffer.WriteString(pe.Key)
			if pe.Next != nil {
				buffer.WriteString(fieldSeparator)
			}
		}

		output.SetFromString(buffer.String())
	} else if input1.mvtype == MT_ARRAY {
		var buffer bytes.Buffer

		for i, _ := range input1.arrayval {
			if i > 0 {
				buffer.WriteString(fieldSeparator)
			}
			// Miller userspace array indices are 1-up
			buffer.WriteString(strconv.Itoa(i + 1))
		}

		output.SetFromString(buffer.String())
	} else {
		output.SetFromError()
	}
}

// ----------------------------------------------------------------
// joinv([3,4,5], ",") -> "3,4,5"
// joinv({"a":3,"b":4,"c":5}, ",") -> "3,4,5"
func MlrvalJoinV(output, input1, input2 *Mlrval) {
	if input2.mvtype != MT_STRING {
		output.SetFromError()
		return
	}
	fieldSeparator := input2.printrep

	if input1.mvtype == MT_MAP {
		var buffer bytes.Buffer

		for pe := input1.mapval.Head; pe != nil; pe = pe.Next {
			buffer.WriteString(pe.Value.String())
			if pe.Next != nil {
				buffer.WriteString(fieldSeparator)
			}
		}

		output.SetFromString(buffer.String())
	} else if input1.mvtype == MT_ARRAY {
		var buffer bytes.Buffer

		for i, element := range input1.arrayval {
			if i > 0 {
				buffer.WriteString(fieldSeparator)
			}
			buffer.WriteString(element.String())
		}

		output.SetFromString(buffer.String())
	} else {
		output.SetFromError()
	}
}

// ----------------------------------------------------------------
// joinkv([3,4,5], "=", ",") -> "1=3,2=4,3=5"
// joinkv({"a":3,"b":4,"c":5}, "=", ",") -> "a=3,b=4,c=5"
func MlrvalJoinKV(output, input1, input2, input3 *Mlrval) {
	if input2.mvtype != MT_STRING {
		output.SetFromError()
		return
	}
	pairSeparator := input2.printrep
	if input3.mvtype != MT_STRING {
		output.SetFromError()
		return
	}
	fieldSeparator := input3.printrep

	if input1.mvtype == MT_MAP {
		var buffer bytes.Buffer

		for pe := input1.mapval.Head; pe != nil; pe = pe.Next {
			buffer.WriteString(pe.Key)
			buffer.WriteString(pairSeparator)
			buffer.WriteString(pe.Value.String())
			if pe.Next != nil {
				buffer.WriteString(fieldSeparator)
			}
		}

		output.SetFromString(buffer.String())
	} else if input1.mvtype == MT_ARRAY {
		var buffer bytes.Buffer

		for i, element := range input1.arrayval {
			if i > 0 {
				buffer.WriteString(fieldSeparator)
			}
			// Miller userspace array indices are 1-up
			buffer.WriteString(strconv.Itoa(i + 1))
			buffer.WriteString(pairSeparator)
			buffer.WriteString(element.String())
		}

		output.SetFromString(buffer.String())
	} else {
		output.SetFromError()
	}
}

// ================================================================
// splitkv("a=3,b=4,c=5", "=", ",") -> {"a":3,"b":4,"c":5}
func MlrvalSplitKV(output, input1, input2, input3 *Mlrval) {
	if input1.mvtype != MT_STRING {
		output.SetFromError()
		return
	}
	if input2.mvtype != MT_STRING {
		output.SetFromError()
		return
	}
	pairSeparator := input2.printrep
	if input3.mvtype != MT_STRING {
		output.SetFromError()
		return
	}
	fieldSeparator := input3.printrep

	*output = MlrvalEmptyMap()

	fields := lib.SplitString(input1.printrep, fieldSeparator)
	for i, field := range fields {
		pair := strings.SplitN(field, pairSeparator, 2)
		if len(pair) == 1 {
			key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
			value := MlrvalFromInferredType(pair[0])
			output.mapval.PutReference(key, &value)
		} else if len(pair) == 2 {
			key := pair[0]
			value := MlrvalFromInferredType(pair[1])
			output.mapval.PutReference(key, &value)
		} else {
			lib.InternalCodingErrorIf(true)
		}
	}
}

// ----------------------------------------------------------------
// splitkvx("a=3,b=4,c=5", "=", ",") -> {"a":"3","b":"4","c":"5"}
func MlrvalSplitKVX(output, input1, input2, input3 *Mlrval) {
	if input1.mvtype != MT_STRING {
		output.SetFromError()
		return
	}
	if input2.mvtype != MT_STRING {
		output.SetFromError()
		return
	}
	pairSeparator := input2.printrep
	if input3.mvtype != MT_STRING {
		output.SetFromError()
		return
	}
	fieldSeparator := input3.printrep

	*output = MlrvalEmptyMap()

	fields := lib.SplitString(input1.printrep, fieldSeparator)
	for i, field := range fields {
		pair := strings.SplitN(field, pairSeparator, 2)
		if len(pair) == 1 {
			key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
			value := MlrvalFromString(pair[0])
			output.mapval.PutReference(key, &value)
		} else if len(pair) == 2 {
			key := pair[0]
			value := MlrvalFromString(pair[1])
			output.mapval.PutReference(key, &value)
		} else {
			lib.InternalCodingErrorIf(true)
		}
	}
}

// ----------------------------------------------------------------
// splitnv("a,b,c", ",") -> {"1":"a","2":"b","3":"c"}
func MlrvalSplitNV(output, input1, input2 *Mlrval) {
	if input1.mvtype != MT_STRING {
		output.SetFromError()
		return
	}
	if input2.mvtype != MT_STRING {
		output.SetFromError()
		return
	}

	*output = MlrvalEmptyMap()

	fields := lib.SplitString(input1.printrep, input2.printrep)
	for i, field := range fields {
		key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
		value := MlrvalFromInferredType(field)
		output.mapval.PutReference(key, &value)
	}
}

// ----------------------------------------------------------------
// splitnvx("3,4,5", ",") -> {"1":"3","2":"4","3":"5"}
func MlrvalSplitNVX(output, input1, input2 *Mlrval) {
	if input1.mvtype != MT_STRING {
		output.SetFromError()
		return
	}
	if input2.mvtype != MT_STRING {
		output.SetFromError()
		return
	}

	*output = MlrvalEmptyMap()

	fields := lib.SplitString(input1.printrep, input2.printrep)
	for i, field := range fields {
		key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
		value := MlrvalFromString(field)
		output.mapval.PutReference(key, &value)
	}
}

// ----------------------------------------------------------------
// splita("3,4,5", ",") -> [3,4,5]
func MlrvalSplitA(output, input1, input2 *Mlrval) {
	if input1.mvtype != MT_STRING {
		output.SetFromError()
		return
	}
	if input2.mvtype != MT_STRING {
		output.SetFromError()
		return
	}
	fieldSeparator := input2.printrep

	fields := lib.SplitString(input1.printrep, fieldSeparator)

	*output = *NewSizedMlrvalArray(int(len(fields)))

	for i, field := range fields {
		value := MlrvalFromInferredType(field)
		output.arrayval[i] = value
	}
}

// ----------------------------------------------------------------
// splitax("3,4,5", ",") -> ["3","4","5"]

func MlrvalSplitAX(output, input1, input2 *Mlrval) {
	if input1.mvtype != MT_STRING {
		output.SetFromError()
		return
	}
	if input2.mvtype != MT_STRING {
		output.SetFromError()
		return
	}
	input := input1.printrep
	fieldSeparator := input2.printrep

	mlrvalSplitAXHelper(output, input, fieldSeparator)
}

// Split out for MlrvalSplitAX and MlrvalUnflatten
func mlrvalSplitAXHelper(output *Mlrval, input string, separator string) {
	fields := lib.SplitString(input, separator)

	*output = *NewSizedMlrvalArray(int(len(fields)))

	for i, field := range fields {
		value := MlrvalFromString(field)
		output.arrayval[i] = value
	}
}

// ----------------------------------------------------------------
func MlrvalGetKeys(output, input1 *Mlrval) {
	if input1.mvtype == MT_MAP {
		// TODO: make a ReferenceFrom with commenbs
		*output = *NewSizedMlrvalArray(input1.mapval.FieldCount)
		i := 0
		for pe := input1.mapval.Head; pe != nil; pe = pe.Next {
			output.arrayval[i] = MlrvalFromString(pe.Key)
			i++
		}

	} else if input1.mvtype == MT_ARRAY {
		*output = *NewSizedMlrvalArray(int(len(input1.arrayval)))
		for i, _ := range input1.arrayval {
			output.arrayval[i] = MlrvalFromInt(int(i + 1)) // Miller user-space indices are 1-up
		}

	} else {
		output.SetFromError()
	}
}

// ----------------------------------------------------------------
func MlrvalGetValues(output, input1 *Mlrval) {
	if input1.mvtype == MT_MAP {
		// TODO: make a ReferenceFrom with commenbs
		*output = *NewSizedMlrvalArray(input1.mapval.FieldCount)
		i := 0
		for pe := input1.mapval.Head; pe != nil; pe = pe.Next {
			output.arrayval[i] = *pe.Value.Copy()
			i++
		}

	} else if input1.mvtype == MT_ARRAY {
		*output = *NewSizedMlrvalArray(int(len(input1.arrayval)))
		for i, value := range input1.arrayval {
			output.arrayval[i] = *value.Copy()
		}

	} else {
		output.SetFromError()
	}
}

// ----------------------------------------------------------------
func MlrvalAppend(output, input1, input2 *Mlrval) {
	if input1.mvtype != MT_ARRAY {
		output.SetFromError()
		return
	}

	*output = *input1.Copy()
	output.ArrayAppend(input2.Copy())
}

// ----------------------------------------------------------------
// First argumemnt is prefix.
// Second argument is delimiter.
// Third argument is map or array.
// flatten("a", ".", {"b": { "c": 4 }}) is {"a.b.c" : 4}.
// flatten("", ".", {"a": { "b": 3 }}) is {"a.b" : 3}.
func MlrvalFlatten(output, input1, input2, input3 *Mlrval) {
	if input3.mvtype == MT_MAP || input3.mvtype == MT_ARRAY {
		if input1.mvtype != MT_STRING && input1.mvtype != MT_VOID {
			output.SetFromError()
			return
		}
		prefix := input1.printrep
		if input2.mvtype != MT_STRING {
			output.SetFromError()
			return
		}
		delimiter := input2.printrep

		temp := input3.FlattenToMap(prefix, delimiter)
		output.CopyFrom(&temp)
	} else {
		output.CopyFrom(input3)
	}
}

// flatten($*, ".") is the same as flatten("", ".", $*)
func MlrvalFlattenBinary(output, input1, input2 *Mlrval) {
	MlrvalFlatten(output, MlrvalPointerFromVoid(), input2, input1)
}

// ----------------------------------------------------------------
// First argument is a map.
// Second argument is a delimiter string.
// unflatten({"a.b.c", ".") is {"a": { "b": { "c": 4}}}.
func MlrvalUnflatten(output, input1, input2 *Mlrval) {
	if input2.mvtype != MT_STRING {
		output.SetFromError()
		return
	}
	if input1.mvtype != MT_MAP {
		output.CopyFrom(input1)
		return
	}
	oldmap := input1.mapval
	separator := input2.printrep
	newmap := NewMlrmap()

	for pe := oldmap.Head; pe != nil; pe = pe.Next {
		// TODO: factor out a shared helper function bewteen here and MlrvalSplitAX.
		arrayOfIndices := MlrvalFromError()
		mlrvalSplitAXHelper(&arrayOfIndices, pe.Key, separator)
		newmap.PutIndexed(MakePointerArray(arrayOfIndices.arrayval), pe.Value.Copy())
	}
	output.SetFromMapReferenced(newmap)
}

// ----------------------------------------------------------------
// Converts maps with "1", "2", ... keys into arrays. Recurses nested data structures.
func MlrvalArrayify(output, input1 *Mlrval) {
	if input1.mvtype == MT_MAP {
		convertible := true
		i := 0
		for pe := input1.mapval.Head; pe != nil; pe = pe.Next {
			sval := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
			i++
			if pe.Key != sval {
				convertible = false
			}
			MlrvalArrayify(pe.Value, pe.Value)
		}

		if convertible {
			arrayval := make([]Mlrval, input1.mapval.FieldCount)
			i := 0
			for pe := input1.mapval.Head; pe != nil; pe = pe.Next {
				arrayval[i] = *pe.Value.Copy()
				i++
			}
			output.SetFromArrayLiteralReference(arrayval)

		} else {
			output.CopyFrom(input1)
		}

	} else if input1.mvtype == MT_ARRAY {
		// TODO: comment (or rethink) that this modifies its inputs!!
		for i, _ := range input1.arrayval {
			MlrvalArrayify(&input1.arrayval[i], &input1.arrayval[i])
		}
		output.CopyFrom(input1)

	} else {
		output.CopyFrom(input1)
	}
}

// ----------------------------------------------------------------
func MlrvalJSONParse(output, input1 *Mlrval) {
	if input1.mvtype == MT_VOID {
		output.CopyFrom(input1)
	} else if input1.mvtype != MT_STRING {
		output.SetFromError()
	} else {
		*output = MlrvalFromPending()
		err := output.UnmarshalJSON([]byte(input1.printrep))
		if err != nil {
			output.SetFromError()
		}
	}
}

func MlrvalJSONStringifyUnary(output, input1 *Mlrval) {
	outputBytes, err := input1.MarshalJSON(JSON_SINGLE_LINE)
	if err != nil {
		output.SetFromError()
	} else {
		output.SetFromString(string(outputBytes))
	}
}

func MlrvalJSONStringifyBinary(output, input1, input2 *Mlrval) {
	var jsonFormatting TJSONFormatting = JSON_SINGLE_LINE
	useMultiline, ok := input2.GetBoolValue()
	if !ok {
		output.SetFromError()
		return
	}
	if useMultiline {
		jsonFormatting = JSON_MULTILINE
	}

	outputBytes, err := input1.MarshalJSON(jsonFormatting)
	if err != nil {
		output.SetFromError()
	} else {
		output.SetFromString(string(outputBytes))
	}
}
