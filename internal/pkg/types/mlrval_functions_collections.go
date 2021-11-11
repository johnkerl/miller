package types

import (
	"bytes"
	"strconv"
	"strings"

	"mlr/internal/pkg/lib"
)

// ================================================================
// Map/array count. Scalars (including strings) have length 1.
func BIF_length(input1 *Mlrval) *Mlrval {
	switch input1.mvtype {
	case MT_ERROR:
		return MlrvalFromInt(0)
		break
	case MT_ABSENT:
		return MlrvalFromInt(0)
		break
	case MT_ARRAY:
		return MlrvalFromInt(int(len(input1.arrayval)))
		break
	case MT_MAP:
		return MlrvalFromInt(int(input1.mapval.FieldCount))
		break
	}
	return MlrvalFromInt(1)
}

// ================================================================
func depth_from_array(input1 *Mlrval) *Mlrval {
	maxChildDepth := 0
	for _, child := range input1.arrayval {
		childDepth := BIF_depth(&child)
		lib.InternalCodingErrorIf(childDepth.mvtype != MT_INT)
		iChildDepth := int(childDepth.intval)
		if iChildDepth > maxChildDepth {
			maxChildDepth = iChildDepth
		}
	}
	return MlrvalFromInt(int(1 + maxChildDepth))
}

func depth_from_map(input1 *Mlrval) *Mlrval {
	maxChildDepth := 0
	for pe := input1.mapval.Head; pe != nil; pe = pe.Next {
		child := pe.Value
		childDepth := BIF_depth(child)
		lib.InternalCodingErrorIf(childDepth.mvtype != MT_INT)
		iChildDepth := int(childDepth.intval)
		if iChildDepth > maxChildDepth {
			maxChildDepth = iChildDepth
		}
	}
	return MlrvalFromInt(int(1 + maxChildDepth))
}

func depth_from_scalar(input1 *Mlrval) *Mlrval {
	return MlrvalFromInt(0)
}

// We get a Golang "initialization loop" due to recursive depth computation
// if this is defined statically. So, we use a "package init" function.
var depth_dispositions = [MT_DIM]UnaryFunc{}

func init() {
	depth_dispositions = [MT_DIM]UnaryFunc{
		/*ERROR  */ _erro1,
		/*ABSENT */ _absn1,
		/*NULL   */ _zero1,
		/*VOID   */ depth_from_scalar,
		/*STRING */ depth_from_scalar,
		/*INT    */ depth_from_scalar,
		/*FLOAT  */ depth_from_scalar,
		/*BOOL   */ depth_from_scalar,
		/*ARRAY  */ depth_from_array,
		/*MAP    */ depth_from_map,
		/*FUNC   */ _erro1,
	}
}

func BIF_depth(input1 *Mlrval) *Mlrval {
	return depth_dispositions[input1.mvtype](input1)
}

// ================================================================
func leafcount_from_array(input1 *Mlrval) *Mlrval {
	sumChildLeafCount := 0
	for _, child := range input1.arrayval {
		// Golang initialization loop if we do this :(
		// childLeafCount := BIF_leafcount(&child)

		childLeafCount := MlrvalFromInt(1)
		if child.mvtype == MT_ARRAY {
			childLeafCount = leafcount_from_array(&child)
		} else if child.mvtype == MT_MAP {
			childLeafCount = leafcount_from_map(&child)
		}

		lib.InternalCodingErrorIf(childLeafCount.mvtype != MT_INT)
		iChildLeafCount := int(childLeafCount.intval)
		sumChildLeafCount += iChildLeafCount
	}
	return MlrvalFromInt(int(sumChildLeafCount))
}

func leafcount_from_map(input1 *Mlrval) *Mlrval {
	sumChildLeafCount := 0
	for pe := input1.mapval.Head; pe != nil; pe = pe.Next {
		child := pe.Value

		// Golang initialization loop if we do this :(
		// childLeafCount := BIF_leafcount(child)

		childLeafCount := MlrvalFromInt(1)
		if child.mvtype == MT_ARRAY {
			childLeafCount = leafcount_from_array(child)
		} else if child.mvtype == MT_MAP {
			childLeafCount = leafcount_from_map(child)
		}

		lib.InternalCodingErrorIf(childLeafCount.mvtype != MT_INT)
		iChildLeafCount := int(childLeafCount.intval)
		sumChildLeafCount += iChildLeafCount
	}
	return MlrvalFromInt(int(sumChildLeafCount))
}

func leafcount_from_scalar(input1 *Mlrval) *Mlrval {
	return MlrvalFromInt(1)
}

var leafcount_dispositions = [MT_DIM]UnaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*NULL   */ _zero1,
	/*VOID   */ leafcount_from_scalar,
	/*STRING */ leafcount_from_scalar,
	/*INT    */ leafcount_from_scalar,
	/*FLOAT  */ leafcount_from_scalar,
	/*BOOL   */ leafcount_from_scalar,
	/*ARRAY  */ leafcount_from_array,
	/*MAP    */ leafcount_from_map,
	/*FUNC   */ _erro1,
}

func BIF_leafcount(input1 *Mlrval) *Mlrval {
	return leafcount_dispositions[input1.mvtype](input1)
}

// ----------------------------------------------------------------
func has_key_in_array(input1, input2 *Mlrval) *Mlrval {
	if input2.mvtype == MT_STRING {
		return MLRVAL_FALSE
	}
	if input2.mvtype != MT_INT {
		return MLRVAL_ERROR
	}
	_, ok := UnaliasArrayIndex(&input1.arrayval, input2.intval)
	return MlrvalFromBool(ok)
}

func has_key_in_map(input1, input2 *Mlrval) *Mlrval {
	if input2.mvtype == MT_STRING || input2.mvtype == MT_INT {
		return MlrvalFromBool(input1.mapval.Has(input2.String()))
	} else {
		return MLRVAL_ERROR
	}
}

func BIF_haskey(input1, input2 *Mlrval) *Mlrval {
	if input1.mvtype == MT_ARRAY {
		return has_key_in_array(input1, input2)
	} else if input1.mvtype == MT_MAP {
		return has_key_in_map(input1, input2)
	} else {
		return MLRVAL_ERROR
	}
}

// ================================================================
func BIF_mapselect(mlrvals []*Mlrval) *Mlrval {
	if len(mlrvals) < 1 {
		return MLRVAL_ERROR
	}
	if mlrvals[0].mvtype != MT_MAP {
		return MLRVAL_ERROR
	}
	oldmap := mlrvals[0].mapval
	newMap := NewMlrmap()

	newKeys := make(map[string]bool)
	for _, selectArg := range mlrvals[1:] {
		if selectArg.mvtype == MT_STRING {
			newKeys[selectArg.printrep] = true
		} else if selectArg.mvtype == MT_INT {
			newKeys[selectArg.String()] = true
		} else if selectArg.mvtype == MT_ARRAY {
			for _, element := range selectArg.arrayval {
				if element.mvtype == MT_STRING {
					newKeys[element.printrep] = true
				} else {
					return MLRVAL_ERROR
				}
			}
		} else {
			return MLRVAL_ERROR
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
func BIF_mapexcept(mlrvals []*Mlrval) *Mlrval {
	if len(mlrvals) < 1 {
		return MLRVAL_ERROR
	}
	if mlrvals[0].mvtype != MT_MAP {
		return MLRVAL_ERROR
	}
	newMap := mlrvals[0].mapval.Copy()

	for _, exceptArg := range mlrvals[1:] {
		if exceptArg.mvtype == MT_STRING {
			newMap.Remove(exceptArg.printrep)
		} else if exceptArg.mvtype == MT_INT {
			newMap.Remove(exceptArg.String())
		} else if exceptArg.mvtype == MT_ARRAY {
			for _, element := range exceptArg.arrayval {
				if element.mvtype == MT_STRING {
					newMap.Remove(element.printrep)
				} else {
					return MLRVAL_ERROR
				}
			}
		} else {
			return MLRVAL_ERROR
		}
	}

	return MlrvalFromMap(newMap)
}

// ----------------------------------------------------------------
func BIF_mapsum(mlrvals []*Mlrval) *Mlrval {
	if len(mlrvals) == 0 {
		return MlrvalFromEmptyMap()
	}
	if len(mlrvals) == 1 {
		return mlrvals[0]
	}
	if mlrvals[0].mvtype != MT_MAP {
		return MLRVAL_ERROR
	}
	newMap := mlrvals[0].mapval.Copy()

	for _, otherMapArg := range mlrvals[1:] {
		if otherMapArg.mvtype != MT_MAP {
			return MLRVAL_ERROR
		}

		for pe := otherMapArg.mapval.Head; pe != nil; pe = pe.Next {
			newMap.PutCopy(pe.Key, pe.Value)
		}
	}

	return MlrvalFromMap(newMap)
}

// ----------------------------------------------------------------
func BIF_mapdiff(mlrvals []*Mlrval) *Mlrval {
	if len(mlrvals) == 0 {
		return MlrvalFromEmptyMap()
	}
	if len(mlrvals) == 1 {
		return mlrvals[0]
	}
	if mlrvals[0].mvtype != MT_MAP {
		return MLRVAL_ERROR
	}
	newMap := mlrvals[0].mapval.Copy()

	for _, otherMapArg := range mlrvals[1:] {
		if otherMapArg.mvtype != MT_MAP {
			return MLRVAL_ERROR
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
func BIF_joink(input1, input2 *Mlrval) *Mlrval {
	if input2.mvtype != MT_STRING {
		return MLRVAL_ERROR
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

		return MlrvalFromString(buffer.String())
	} else if input1.mvtype == MT_ARRAY {
		var buffer bytes.Buffer

		for i := range input1.arrayval {
			if i > 0 {
				buffer.WriteString(fieldSeparator)
			}
			// Miller userspace array indices are 1-up
			buffer.WriteString(strconv.Itoa(i + 1))
		}

		return MlrvalFromString(buffer.String())
	} else {
		return MLRVAL_ERROR
	}
}

// ----------------------------------------------------------------
// joinv([3,4,5], ",") -> "3,4,5"
// joinv({"a":3,"b":4,"c":5}, ",") -> "3,4,5"
func BIF_joinv(input1, input2 *Mlrval) *Mlrval {
	if input2.mvtype != MT_STRING {
		return MLRVAL_ERROR
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

		return MlrvalFromString(buffer.String())
	} else if input1.mvtype == MT_ARRAY {
		var buffer bytes.Buffer

		for i, element := range input1.arrayval {
			if i > 0 {
				buffer.WriteString(fieldSeparator)
			}
			buffer.WriteString(element.String())
		}

		return MlrvalFromString(buffer.String())
	} else {
		return MLRVAL_ERROR
	}
}

// ----------------------------------------------------------------
// joinkv([3,4,5], "=", ",") -> "1=3,2=4,3=5"
// joinkv({"a":3,"b":4,"c":5}, "=", ",") -> "a=3,b=4,c=5"
func BIF_joinkv(input1, input2, input3 *Mlrval) *Mlrval {
	if input2.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}
	pairSeparator := input2.printrep
	if input3.mvtype != MT_STRING {
		return MLRVAL_ERROR
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

		return MlrvalFromString(buffer.String())
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

		return MlrvalFromString(buffer.String())
	} else {
		return MLRVAL_ERROR
	}
}

// ================================================================
// splitkv("a=3,b=4,c=5", "=", ",") -> {"a":3,"b":4,"c":5}
func BIF_splitkv(input1, input2, input3 *Mlrval) *Mlrval {
	if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR
	}
	if input2.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}
	pairSeparator := input2.printrep
	if input3.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}
	fieldSeparator := input3.printrep

	output := MlrvalFromEmptyMap()

	fields := lib.SplitString(input1.printrep, fieldSeparator)
	for i, field := range fields {
		pair := strings.SplitN(field, pairSeparator, 2)
		if len(pair) == 1 {
			key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
			value := MlrvalFromInferredType(pair[0])
			output.mapval.PutReference(key, value)
		} else if len(pair) == 2 {
			key := pair[0]
			value := MlrvalFromInferredType(pair[1])
			output.mapval.PutReference(key, value)
		} else {
			lib.InternalCodingErrorIf(true)
		}
	}
	return output
}

// ----------------------------------------------------------------
// splitkvx("a=3,b=4,c=5", "=", ",") -> {"a":"3","b":"4","c":"5"}
func BIF_splitkvx(input1, input2, input3 *Mlrval) *Mlrval {
	if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR
	}
	if input2.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}
	pairSeparator := input2.printrep
	if input3.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}
	fieldSeparator := input3.printrep

	output := MlrvalFromEmptyMap()

	fields := lib.SplitString(input1.printrep, fieldSeparator)
	for i, field := range fields {
		pair := strings.SplitN(field, pairSeparator, 2)
		if len(pair) == 1 {
			key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
			value := MlrvalFromString(pair[0])
			output.mapval.PutReference(key, value)
		} else if len(pair) == 2 {
			key := pair[0]
			value := MlrvalFromString(pair[1])
			output.mapval.PutReference(key, value)
		} else {
			lib.InternalCodingErrorIf(true)
		}
	}

	return output
}

// ----------------------------------------------------------------
// splitnv("a,b,c", ",") -> {"1":"a","2":"b","3":"c"}
func BIF_splitnv(input1, input2 *Mlrval) *Mlrval {
	if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR
	}
	if input2.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}

	output := MlrvalFromEmptyMap()

	fields := lib.SplitString(input1.printrep, input2.printrep)
	for i, field := range fields {
		key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
		value := MlrvalFromInferredType(field)
		output.mapval.PutReference(key, value)
	}

	return output
}

// ----------------------------------------------------------------
// splitnvx("3,4,5", ",") -> {"1":"3","2":"4","3":"5"}
func BIF_splitnvx(input1, input2 *Mlrval) *Mlrval {
	if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR
	}
	if input2.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}

	output := MlrvalFromEmptyMap()

	fields := lib.SplitString(input1.printrep, input2.printrep)
	for i, field := range fields {
		key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
		value := MlrvalFromString(field)
		output.mapval.PutReference(key, value)
	}

	return output
}

// ----------------------------------------------------------------
// splita("3,4,5", ",") -> [3,4,5]
func BIF_splita(input1, input2 *Mlrval) *Mlrval {
	if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR
	}
	if !input2.IsString() {
		return MLRVAL_ERROR
	}
	fieldSeparator := input2.printrep

	fields := lib.SplitString(input1.printrep, fieldSeparator)

	output := NewSizedMlrvalArray(int(len(fields)))

	for i, field := range fields {
		value := MlrvalFromInferredType(field)
		output.arrayval[i] = *value
	}

	return output
}

// ----------------------------------------------------------------
// BIF_splitax splits a string to an array, without type-inference:
// e.g. splitax("3,4,5", ",") -> ["3","4","5"]
func BIF_splitax(input1, input2 *Mlrval) *Mlrval {
	if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR
	}
	if input2.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}
	input := input1.printrep
	fieldSeparator := input2.printrep

	return mlrvalSplitAXHelper(input, fieldSeparator)
}

// mlrvalSplitAXHelper is Split out for the benefit of BIF_splitax and
// BIF_unflatten.
func mlrvalSplitAXHelper(input string, separator string) *Mlrval {
	fields := lib.SplitString(input, separator)

	output := NewSizedMlrvalArray(int(len(fields)))

	for i, field := range fields {
		output.arrayval[i] = *MlrvalFromString(field)
	}

	return output
}

// ----------------------------------------------------------------
func BIF_get_keys(input1 *Mlrval) *Mlrval {
	if input1.mvtype == MT_MAP {
		// TODO: make a ReferenceFrom with comments
		output := NewSizedMlrvalArray(input1.mapval.FieldCount)
		i := 0
		for pe := input1.mapval.Head; pe != nil; pe = pe.Next {
			output.arrayval[i] = *MlrvalFromString(pe.Key)
			i++
		}
		return output

	} else if input1.mvtype == MT_ARRAY {
		output := NewSizedMlrvalArray(int(len(input1.arrayval)))
		for i := range input1.arrayval {
			output.arrayval[i] = *MlrvalFromInt(int(i + 1)) // Miller user-space indices are 1-up
		}
		return output

	} else {
		return MLRVAL_ERROR
	}
}

// ----------------------------------------------------------------
func BIF_get_values(input1 *Mlrval) *Mlrval {
	if input1.mvtype == MT_MAP {
		// TODO: make a ReferenceFrom with commenbs
		output := NewSizedMlrvalArray(input1.mapval.FieldCount)
		i := 0
		for pe := input1.mapval.Head; pe != nil; pe = pe.Next {
			output.arrayval[i] = *pe.Value.Copy()
			i++
		}
		return output

	} else if input1.mvtype == MT_ARRAY {
		output := NewSizedMlrvalArray(int(len(input1.arrayval)))
		for i, value := range input1.arrayval {
			output.arrayval[i] = *value.Copy()
		}
		return output

	} else {
		return MLRVAL_ERROR
	}
}

// ----------------------------------------------------------------
func BIF_append(input1, input2 *Mlrval) *Mlrval {
	if input1.mvtype != MT_ARRAY {
		return MLRVAL_ERROR
	}

	output := input1.Copy()
	output.ArrayAppend(input2.Copy())
	return output
}

// ----------------------------------------------------------------
// First argumemnt is prefix.
// Second argument is delimiter.
// Third argument is map or array.
// flatten("a", ".", {"b": { "c": 4 }}) is {"a.b.c" : 4}.
// flatten("", ".", {"a": { "b": 3 }}) is {"a.b" : 3}.
func BIF_flatten(input1, input2, input3 *Mlrval) *Mlrval {
	if input3.mvtype == MT_MAP || input3.mvtype == MT_ARRAY {
		if input1.mvtype != MT_STRING && input1.mvtype != MT_VOID {
			return MLRVAL_ERROR
		}
		prefix := input1.printrep
		if input2.mvtype != MT_STRING {
			return MLRVAL_ERROR
		}
		delimiter := input2.printrep

		retval := input3.FlattenToMap(prefix, delimiter)
		return &retval
	} else {
		return input3
	}
}

// flatten($*, ".") is the same as flatten("", ".", $*)
func BIF_flatten_binary(input1, input2 *Mlrval) *Mlrval {
	return BIF_flatten(MLRVAL_VOID, input2, input1)
}

// ----------------------------------------------------------------
// First argument is a map.
// Second argument is a delimiter string.
// unflatten({"a.b.c", ".") is {"a": { "b": { "c": 4}}}.
func BIF_unflatten(input1, input2 *Mlrval) *Mlrval {
	if input2.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}
	if input1.mvtype != MT_MAP {
		return input1
	}
	oldmap := input1.mapval
	separator := input2.printrep
	newmap := oldmap.CopyUnflattened(separator)
	return MlrvalFromMapReferenced(newmap)
}

// ----------------------------------------------------------------
// Converts maps with "1", "2", ... keys into arrays. Recurses nested data structures.
func BIF_arrayify(input1 *Mlrval) *Mlrval {
	if input1.mvtype == MT_MAP {
		if input1.mapval.IsEmpty() {
			return input1
		}

		convertible := true
		i := 0
		for pe := input1.mapval.Head; pe != nil; pe = pe.Next {
			sval := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
			i++
			if pe.Key != sval {
				convertible = false
			}
			pe.Value = BIF_arrayify(pe.Value)
		}

		if convertible {
			arrayval := make([]Mlrval, input1.mapval.FieldCount)
			i := 0
			for pe := input1.mapval.Head; pe != nil; pe = pe.Next {
				arrayval[i] = *pe.Value.Copy()
				i++
			}
			return MlrvalFromArrayReference(arrayval)

		} else {
			return input1
		}

	} else if input1.mvtype == MT_ARRAY {
		// TODO: comment (or rethink) that this modifies its inputs!!
		output := input1.Copy()
		for i := range input1.arrayval {
			output.arrayval[i] = *BIF_arrayify(&output.arrayval[i])
		}
		return output

	} else {
		return input1
	}
}

// ----------------------------------------------------------------
func BIF_json_parse(input1 *Mlrval) *Mlrval {
	if input1.mvtype == MT_VOID {
		return input1
	} else if input1.mvtype != MT_STRING {
		return MLRVAL_ERROR
	} else {
		output := MlrvalFromPending()
		err := output.UnmarshalJSON([]byte(input1.printrep))
		if err != nil {
			return MLRVAL_ERROR
		}
		return &output
	}
}

func BIF_json_stringify_unary(input1 *Mlrval) *Mlrval {
	outputBytes, err := input1.MarshalJSON(JSON_SINGLE_LINE, false)
	if err != nil {
		return MLRVAL_ERROR
	} else {
		return MlrvalFromString(string(outputBytes))
	}
}

func BIF_json_stringify_binary(input1, input2 *Mlrval) *Mlrval {
	var jsonFormatting TJSONFormatting = JSON_SINGLE_LINE
	useMultiline, ok := input2.GetBoolValue()
	if !ok {
		return MLRVAL_ERROR
	}
	if useMultiline {
		jsonFormatting = JSON_MULTILINE
	}

	outputBytes, err := input1.MarshalJSON(jsonFormatting, false)
	if err != nil {
		return MLRVAL_ERROR
	} else {
		return MlrvalFromString(string(outputBytes))
	}
}
