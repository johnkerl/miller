package bifs

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

// ================================================================
// Map/array count. Scalars (including strings) have length 1; strlen is for string length.
func BIF_length(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	switch input1.Type() {
	case mlrval.MT_ERROR:
		return mlrval.FromInt(0)
		break
	case mlrval.MT_ABSENT:
		return mlrval.FromInt(0)
		break
	case mlrval.MT_ARRAY:
		arrayval := input1.AcquireArrayValue()
		return mlrval.FromInt(int(len(arrayval)))
		break
	case mlrval.MT_MAP:
		mapval := input1.AcquireMapValue()
		return mlrval.FromInt(int(mapval.FieldCount))
		break
	}
	return mlrval.FromInt(1)
}

// ================================================================
func depth_from_array(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	maxChildDepth := 0
	arrayval := input1.AcquireArrayValue()
	for _, child := range arrayval {
		childDepth := BIF_depth(child)
		lib.InternalCodingErrorIf(!childDepth.IsInt())
		iChildDepth := int(childDepth.AcquireIntValue())
		if iChildDepth > maxChildDepth {
			maxChildDepth = iChildDepth
		}
	}
	return mlrval.FromInt(int(1 + maxChildDepth))
}

func depth_from_map(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	maxChildDepth := 0
	mapval := input1.AcquireMapValue()
	for pe := mapval.Head; pe != nil; pe = pe.Next {
		child := pe.Value
		childDepth := BIF_depth(child)
		lib.InternalCodingErrorIf(!childDepth.IsInt())
		iChildDepth := int(childDepth.AcquireIntValue())
		if iChildDepth > maxChildDepth {
			maxChildDepth = iChildDepth
		}
	}
	return mlrval.FromInt(int(1 + maxChildDepth))
}

func depth_from_scalar(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(0)
}

// We get a Golang "initialization loop" due to recursive depth computation
// if this is defined statically. So, we use a "package init" function.
var depth_dispositions = [mlrval.MT_DIM]UnaryFunc{}

func init() {
	depth_dispositions = [mlrval.MT_DIM]UnaryFunc{
		/*INT    */ depth_from_scalar,
		/*FLOAT  */ depth_from_scalar,
		/*BOOL   */ depth_from_scalar,
		/*VOID   */ depth_from_scalar,
		/*STRING */ depth_from_scalar,
		/*ARRAY  */ depth_from_array,
		/*MAP    */ depth_from_map,
		/*FUNC   */ _erro1,
		/*ERROR  */ _erro1,
		/*NULL   */ _zero1,
		/*ABSENT */ _absn1,
	}
}

func BIF_depth(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return depth_dispositions[input1.Type()](input1)
}

// ================================================================
func leafcount_from_array(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	sumChildLeafCount := 0
	arrayval := input1.AcquireArrayValue()
	for _, child := range arrayval {
		// Golang initialization loop if we do this :(
		// childLeafCount := BIF_leafcount(&child)

		childLeafCount := mlrval.FromInt(1)
		if child.IsArray() {
			childLeafCount = leafcount_from_array(child)
		} else if child.IsMap() {
			childLeafCount = leafcount_from_map(child)
		}

		lib.InternalCodingErrorIf(!childLeafCount.IsInt())
		iChildLeafCount := int(childLeafCount.AcquireIntValue())
		sumChildLeafCount += iChildLeafCount
	}
	return mlrval.FromInt(int(sumChildLeafCount))
}

func leafcount_from_map(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	sumChildLeafCount := 0
	mapval := input1.AcquireMapValue()
	for pe := mapval.Head; pe != nil; pe = pe.Next {
		child := pe.Value

		// Golang initialization loop if we do this :(
		// childLeafCount := BIF_leafcount(child)

		childLeafCount := mlrval.FromInt(1)
		if child.IsArray() {
			childLeafCount = leafcount_from_array(child)
		} else if child.IsMap() {
			childLeafCount = leafcount_from_map(child)
		}

		lib.InternalCodingErrorIf(!childLeafCount.IsInt())
		iChildLeafCount := int(childLeafCount.AcquireIntValue())
		sumChildLeafCount += iChildLeafCount
	}
	return mlrval.FromInt(int(sumChildLeafCount))
}

func leafcount_from_scalar(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(1)
}

var leafcount_dispositions = [mlrval.MT_DIM]UnaryFunc{
	/*INT    */ leafcount_from_scalar,
	/*FLOAT  */ leafcount_from_scalar,
	/*BOOL   */ leafcount_from_scalar,
	/*VOID   */ leafcount_from_scalar,
	/*STRING */ leafcount_from_scalar,
	/*ARRAY  */ leafcount_from_array,
	/*MAP    */ leafcount_from_map,
	/*FUNC   */ _erro1,
	/*ERROR  */ _erro1,
	/*NULL   */ _zero1,
	/*ABSENT */ _absn1,
}

func BIF_leafcount(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return leafcount_dispositions[input1.Type()](input1)
}

// ----------------------------------------------------------------
func has_key_in_array(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if input2.IsString() {
		return mlrval.FALSE
	}
	if !input2.IsInt() {
		return mlrval.ERROR
	}
	arrayval := input1.AcquireArrayValue()
	_, ok := unaliasArrayIndex(&arrayval, input2.AcquireIntValue())
	return mlrval.FromBool(ok)
}

func has_key_in_map(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if input2.IsString() || input2.IsInt() {
		return mlrval.FromBool(input1.AcquireMapValue().Has(input2.String()))
	} else {
		return mlrval.ERROR
	}
}

func BIF_haskey(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsArray() {
		return has_key_in_array(input1, input2)
	} else if input1.IsMap() {
		return has_key_in_map(input1, input2)
	} else {
		return mlrval.ERROR
	}
}

// ================================================================
func BIF_concat(mlrvals []*mlrval.Mlrval) *mlrval.Mlrval {
	output := mlrval.FromEmptyArray()

	for _, arg := range mlrvals {
		argArray := arg.GetArray()
		if argArray == nil { // not an array
			output.ArrayAppend(arg.Copy())
		} else {
			for i := range argArray {
				output.ArrayAppend(argArray[i].Copy())
			}
		}
	}

	return output
}

// ================================================================
func BIF_mapselect(mlrvals []*mlrval.Mlrval) *mlrval.Mlrval {
	if len(mlrvals) < 1 {
		return mlrval.ERROR
	}
	if !mlrvals[0].IsMap() {
		return mlrval.ERROR
	}
	oldmap := mlrvals[0].AcquireMapValue()
	newMap := mlrval.NewMlrmap()

	newKeys := make(map[string]bool)
	for _, selectArg := range mlrvals[1:] {
		if selectArg.IsString() {
			newKeys[selectArg.AcquireStringValue()] = true
		} else if selectArg.IsInt() {
			newKeys[selectArg.String()] = true
		} else if selectArg.IsArray() {
			for _, element := range selectArg.AcquireArrayValue() {
				if element.IsString() {
					newKeys[element.AcquireStringValue()] = true
				} else {
					return mlrval.ERROR
				}
			}
		} else {
			return mlrval.ERROR
		}
	}

	for pe := oldmap.Head; pe != nil; pe = pe.Next {
		oldKey := pe.Key
		_, present := newKeys[oldKey]
		if present {
			newMap.PutCopy(oldKey, oldmap.Get(oldKey))
		}
	}

	return mlrval.FromMap(newMap)
}

// ----------------------------------------------------------------
func BIF_mapexcept(mlrvals []*mlrval.Mlrval) *mlrval.Mlrval {
	if len(mlrvals) < 1 {
		return mlrval.ERROR
	}
	if !mlrvals[0].IsMap() {
		return mlrval.ERROR
	}
	newMap := mlrvals[0].AcquireMapValue().Copy()

	for _, exceptArg := range mlrvals[1:] {
		if exceptArg.IsString() {
			newMap.Remove(exceptArg.AcquireStringValue())
		} else if exceptArg.IsInt() {
			newMap.Remove(exceptArg.String())
		} else if exceptArg.IsArray() {
			for _, element := range exceptArg.AcquireArrayValue() {
				if element.IsString() {
					newMap.Remove(element.AcquireStringValue())
				} else {
					return mlrval.ERROR
				}
			}
		} else {
			return mlrval.ERROR
		}
	}

	return mlrval.FromMap(newMap)
}

// ----------------------------------------------------------------
func BIF_mapsum(mlrvals []*mlrval.Mlrval) *mlrval.Mlrval {
	if len(mlrvals) == 0 {
		return mlrval.FromEmptyMap()
	}
	if len(mlrvals) == 1 {
		return mlrvals[0]
	}
	if mlrvals[0].Type() != mlrval.MT_MAP {
		return mlrval.ERROR
	}
	newMap := mlrvals[0].AcquireMapValue().Copy()

	for _, otherMapArg := range mlrvals[1:] {
		if otherMapArg.Type() != mlrval.MT_MAP {
			return mlrval.ERROR
		}

		for pe := otherMapArg.AcquireMapValue().Head; pe != nil; pe = pe.Next {
			newMap.PutCopy(pe.Key, pe.Value)
		}
	}

	return mlrval.FromMap(newMap)
}

// ----------------------------------------------------------------
func BIF_mapdiff(mlrvals []*mlrval.Mlrval) *mlrval.Mlrval {
	if len(mlrvals) == 0 {
		return mlrval.FromEmptyMap()
	}
	if len(mlrvals) == 1 {
		return mlrvals[0]
	}
	if !mlrvals[0].IsMap() {
		return mlrval.ERROR
	}
	newMap := mlrvals[0].AcquireMapValue().Copy()

	for _, otherMapArg := range mlrvals[1:] {
		if !otherMapArg.IsMap() {
			return mlrval.ERROR
		}

		for pe := otherMapArg.AcquireMapValue().Head; pe != nil; pe = pe.Next {
			newMap.Remove(pe.Key)
		}
	}

	return mlrval.FromMap(newMap)
}

// ================================================================
// joink([1,2,3], ",") -> "1,2,3"
// joink({"a":3,"b":4,"c":5}, ",") -> "a,b,c"
func BIF_joink(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input2.IsString() {
		return mlrval.ERROR
	}
	fieldSeparator := input2.AcquireStringValue()
	if input1.IsMap() {
		var buffer bytes.Buffer

		for pe := input1.AcquireMapValue().Head; pe != nil; pe = pe.Next {
			buffer.WriteString(pe.Key)
			if pe.Next != nil {
				buffer.WriteString(fieldSeparator)
			}
		}

		return mlrval.FromString(buffer.String())
	} else if input1.IsArray() {
		var buffer bytes.Buffer

		for i := range input1.AcquireArrayValue() {
			if i > 0 {
				buffer.WriteString(fieldSeparator)
			}
			// Miller userspace array indices are 1-up
			buffer.WriteString(strconv.Itoa(i + 1))
		}

		return mlrval.FromString(buffer.String())
	} else {
		return mlrval.ERROR
	}
}

// ----------------------------------------------------------------
// joinv([3,4,5], ",") -> "3,4,5"
// joinv({"a":3,"b":4,"c":5}, ",") -> "3,4,5"
func BIF_joinv(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input2.IsString() {
		return mlrval.ERROR
	}
	fieldSeparator := input2.AcquireStringValue()

	if input1.IsMap() {
		var buffer bytes.Buffer

		for pe := input1.AcquireMapValue().Head; pe != nil; pe = pe.Next {
			buffer.WriteString(pe.Value.String())
			if pe.Next != nil {
				buffer.WriteString(fieldSeparator)
			}
		}

		return mlrval.FromString(buffer.String())
	} else if input1.IsArray() {
		var buffer bytes.Buffer

		for i, element := range input1.AcquireArrayValue() {
			if i > 0 {
				buffer.WriteString(fieldSeparator)
			}
			buffer.WriteString(element.String())
		}

		return mlrval.FromString(buffer.String())
	} else {
		return mlrval.ERROR
	}
}

// ----------------------------------------------------------------
// joinkv([3,4,5], "=", ",") -> "1=3,2=4,3=5"
// joinkv({"a":3,"b":4,"c":5}, "=", ",") -> "a=3,b=4,c=5"
func BIF_joinkv(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input2.IsString() {
		return mlrval.ERROR
	}
	pairSeparator := input2.AcquireStringValue()
	if !input3.IsString() {
		return mlrval.ERROR
	}
	fieldSeparator := input3.AcquireStringValue()

	if input1.IsMap() {
		var buffer bytes.Buffer

		for pe := input1.AcquireMapValue().Head; pe != nil; pe = pe.Next {
			buffer.WriteString(pe.Key)
			buffer.WriteString(pairSeparator)
			buffer.WriteString(pe.Value.String())
			if pe.Next != nil {
				buffer.WriteString(fieldSeparator)
			}
		}

		return mlrval.FromString(buffer.String())
	} else if input1.IsArray() {
		var buffer bytes.Buffer

		for i, element := range input1.AcquireArrayValue() {
			if i > 0 {
				buffer.WriteString(fieldSeparator)
			}
			// Miller userspace array indices are 1-up
			buffer.WriteString(strconv.Itoa(i + 1))
			buffer.WriteString(pairSeparator)
			buffer.WriteString(element.String())
		}

		return mlrval.FromString(buffer.String())
	} else {
		return mlrval.ERROR
	}
}

// ================================================================
// splitkv("a=3,b=4,c=5", "=", ",") -> {"a":3,"b":4,"c":5}
func BIF_splitkv(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsStringOrVoid() {
		return mlrval.ERROR
	}
	if !input2.IsString() {
		return mlrval.ERROR
	}
	pairSeparator := input2.AcquireStringValue()
	if !input3.IsString() {
		return mlrval.ERROR
	}
	fieldSeparator := input3.AcquireStringValue()

	output := mlrval.FromMap(mlrval.NewMlrmap())

	fields := lib.SplitString(input1.AcquireStringValue(), fieldSeparator)
	for i, field := range fields {
		pair := strings.SplitN(field, pairSeparator, 2)
		if len(pair) == 1 {
			key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
			value := mlrval.FromInferredType(pair[0])
			output.AcquireMapValue().PutReference(key, value)
		} else if len(pair) == 2 {
			key := pair[0]
			value := mlrval.FromInferredType(pair[1])
			output.AcquireMapValue().PutReference(key, value)
		} else {
			lib.InternalCodingErrorIf(true)
		}
	}
	return output
}

// ----------------------------------------------------------------
// splitkvx("a=3,b=4,c=5", "=", ",") -> {"a":"3","b":"4","c":"5"}
func BIF_splitkvx(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsStringOrVoid() {
		return mlrval.ERROR
	}
	if !input2.IsString() {
		return mlrval.ERROR
	}
	pairSeparator := input2.AcquireStringValue()
	if !input3.IsString() {
		return mlrval.ERROR
	}
	fieldSeparator := input3.AcquireStringValue()

	output := mlrval.FromMap(mlrval.NewMlrmap())

	fields := lib.SplitString(input1.AcquireStringValue(), fieldSeparator)
	for i, field := range fields {
		pair := strings.SplitN(field, pairSeparator, 2)
		if len(pair) == 1 {
			key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
			value := mlrval.FromString(pair[0])
			output.AcquireMapValue().PutReference(key, value)
		} else if len(pair) == 2 {
			key := pair[0]
			value := mlrval.FromString(pair[1])
			output.AcquireMapValue().PutReference(key, value)
		} else {
			lib.InternalCodingErrorIf(true)
		}
	}

	return output
}

// ----------------------------------------------------------------
// splitnv("a,b,c", ",") -> {"1":"a","2":"b","3":"c"}
func BIF_splitnv(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsStringOrVoid() {
		return mlrval.ERROR
	}
	if !input2.IsString() {
		return mlrval.ERROR
	}

	output := mlrval.FromMap(mlrval.NewMlrmap())

	fields := lib.SplitString(input1.AcquireStringValue(), input2.AcquireStringValue())
	for i, field := range fields {
		key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
		value := mlrval.FromInferredType(field)
		output.AcquireMapValue().PutReference(key, value)
	}

	return output
}

// ----------------------------------------------------------------
// splitnvx("3,4,5", ",") -> {"1":"3","2":"4","3":"5"}
func BIF_splitnvx(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsStringOrVoid() {
		return mlrval.ERROR
	}
	if !input2.IsString() {
		return mlrval.ERROR
	}

	output := mlrval.FromMap(mlrval.NewMlrmap())

	fields := lib.SplitString(input1.AcquireStringValue(), input2.AcquireStringValue())
	for i, field := range fields {
		key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
		value := mlrval.FromString(field)
		output.AcquireMapValue().PutReference(key, value)
	}

	return output
}

// ----------------------------------------------------------------
// splita("3,4,5", ",") -> [3,4,5]
func BIF_splita(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsStringOrVoid() {
		return mlrval.ERROR
	}
	if !input2.IsString() {
		return mlrval.ERROR
	}
	fieldSeparator := input2.AcquireStringValue()

	fields := lib.SplitString(input1.AcquireStringValue(), fieldSeparator)

	arrayval := make([]*mlrval.Mlrval, len(fields))

	for i, field := range fields {
		value := mlrval.FromInferredType(field)
		arrayval[i] = value
	}

	return mlrval.FromArray(arrayval)
}

// ----------------------------------------------------------------
// BIF_splitax splits a string to an array, without type-inference:
// e.g. splitax("3,4,5", ",") -> ["3","4","5"]
func BIF_splitax(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsStringOrVoid() {
		return mlrval.ERROR
	}
	if !input2.IsString() {
		return mlrval.ERROR
	}
	input := input1.AcquireStringValue()
	fieldSeparator := input2.AcquireStringValue()

	return bif_splitax_helper(input, fieldSeparator)
}

// bif_splitax_helper is split out for the benefit of BIF_splitax and
// BIF_unflatten.
func bif_splitax_helper(input string, separator string) *mlrval.Mlrval {
	fields := lib.SplitString(input, separator)

	arrayval := make([]*mlrval.Mlrval, len(fields))

	for i, field := range fields {
		arrayval[i] = mlrval.FromString(field)
	}

	return mlrval.FromArray(arrayval)
}

// ----------------------------------------------------------------
func BIF_get_keys(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsMap() {
		// TODO: make a ReferenceFrom with comments
		mapval := input1.AcquireMapValue()
		arrayval := make([]*mlrval.Mlrval, mapval.FieldCount)
		i := 0
		for pe := mapval.Head; pe != nil; pe = pe.Next {
			arrayval[i] = mlrval.FromString(pe.Key)
			i++
		}
		return mlrval.FromArray(arrayval)

	} else if input1.IsArray() {
		inputarrayval := input1.AcquireArrayValue()
		arrayval := make([]*mlrval.Mlrval, len(inputarrayval))
		for i := range inputarrayval {
			arrayval[i] = mlrval.FromInt(int(i + 1)) // Miller user-space indices are 1-up
		}
		return mlrval.FromArray(arrayval)

	} else {
		return mlrval.ERROR
	}
}

// ----------------------------------------------------------------
func BIF_get_values(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsMap() {
		mapval := input1.AcquireMapValue()
		arrayval := make([]*mlrval.Mlrval, mapval.FieldCount)
		i := 0
		for pe := mapval.Head; pe != nil; pe = pe.Next {
			arrayval[i] = pe.Value.Copy()
			i++
		}
		return mlrval.FromArray(arrayval)

	} else if input1.IsArray() {
		inputarrayval := input1.AcquireArrayValue()
		arrayval := make([]*mlrval.Mlrval, len(inputarrayval))
		for i, value := range inputarrayval {
			arrayval[i] = value.Copy()
		}
		return mlrval.FromArray(arrayval)

	} else {
		return mlrval.ERROR
	}
}

// ----------------------------------------------------------------
func BIF_append(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsArray() {
		return mlrval.ERROR
	}

	output := input1.Copy()
	output.ArrayAppend(input2.Copy())
	return output
}

// ----------------------------------------------------------------
// First argument is prefix.
// Second argument is delimiter.
// Third argument is map or array.
// flatten("a", ".", {"b": { "c": 4 }}) is {"a.b.c" : 4}.
// flatten("", ".", {"a": { "b": 3 }}) is {"a.b" : 3}.
func BIF_flatten(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	if input3.IsMap() || input3.IsArray() {
		if !input1.IsString() && input1.Type() != mlrval.MT_VOID {
			return mlrval.ERROR
		}
		prefix := input1.AcquireStringValue()
		if !input2.IsString() {
			return mlrval.ERROR
		}
		delimiter := input2.AcquireStringValue()

		retval := input3.FlattenToMap(prefix, delimiter)
		return &retval
	} else {
		return input3
	}
}

// flatten($*, ".") is the same as flatten("", ".", $*)
func BIF_flatten_binary(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return BIF_flatten(mlrval.VOID, input2, input1)
}

// ----------------------------------------------------------------
// First argument is a map.
// Second argument is a delimiter string.
// unflatten({"a.b.c", ".") is {"a": { "b": { "c": 4}}}.
func BIF_unflatten(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input2.IsString() {
		return mlrval.ERROR
	}
	if input1.Type() != mlrval.MT_MAP {
		return input1
	}
	oldmap := input1.AcquireMapValue()
	separator := input2.AcquireStringValue()
	newmap := oldmap.CopyUnflattened(separator)
	return mlrval.FromMap(newmap)
}

// ----------------------------------------------------------------
// Converts maps with "1", "2", ... keys into arrays. Recurses nested data structures.
func BIF_arrayify(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsMap() {
		if input1.AcquireMapValue().IsEmpty() {
			return input1
		}

		convertible := true
		i := 0
		for pe := input1.AcquireMapValue().Head; pe != nil; pe = pe.Next {
			sval := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
			i++
			if pe.Key != sval {
				convertible = false
			}
			pe.Value = BIF_arrayify(pe.Value)
		}

		if convertible {
			mapval := input1.AcquireMapValue()
			arrayval := make([]*mlrval.Mlrval, input1.AcquireMapValue().FieldCount)
			i := 0
			for pe := mapval.Head; pe != nil; pe = pe.Next {
				arrayval[i] = pe.Value.Copy()
				i++
			}
			return mlrval.FromArray(arrayval)

		} else {
			return input1
		}

	} else if input1.IsArray() {
		// TODO: comment (or rethink) that this modifies its inputs!!
		output := input1.Copy()
		arrayval := output.AcquireArrayValue()
		for i := range input1.AcquireArrayValue() {
			arrayval[i] = BIF_arrayify(arrayval[i])
		}
		return output

	} else {
		return input1
	}
}

// ----------------------------------------------------------------
func BIF_json_parse(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsVoid() {
		return input1
	} else if !input1.IsString() {
		return mlrval.ERROR
	} else {
		output := mlrval.FromPending()
		err := output.UnmarshalJSON([]byte(input1.AcquireStringValue()))
		if err != nil {
			return mlrval.ERROR
		}
		return output
	}
}

func BIF_json_stringify_unary(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	outputBytes, err := input1.MarshalJSON(mlrval.JSON_SINGLE_LINE, false)
	if err != nil {
		return mlrval.ERROR
	} else {
		return mlrval.FromString(string(outputBytes))
	}
}

func BIF_json_stringify_binary(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	var jsonFormatting mlrval.TJSONFormatting = mlrval.JSON_SINGLE_LINE
	useMultiline, ok := input2.GetBoolValue()
	if !ok {
		return mlrval.ERROR
	}
	if useMultiline {
		jsonFormatting = mlrval.JSON_MULTILINE
	}

	outputBytes, err := input1.MarshalJSON(jsonFormatting, false)
	if err != nil {
		return mlrval.ERROR
	} else {
		return mlrval.FromString(string(outputBytes))
	}
}

func unaliasArrayIndex(array *[]*mlrval.Mlrval, mindex int) (int, bool) {
	n := int(len(*array))
	return unaliasArrayLengthIndex(n, mindex)
}

// Input "mindex" is a Miller DSL array index. These are 1-up, so 1..n where n
// is the length of the array. Also, -n..-1 are aliases to 1..n. 0 is never a
// valid index.
//
// Output "zindex" is a Golang array index. These are 0-up, so 0..(n-1).
//
// The second return value indicates whether the Miller index is in-bounds.
// Even if it's out of bounds, while the second return value is false, the
// first return is correctly de-aliased. E.g. if the array has length 5 and the
// mindex is 8, zindex is 7 and valid=false. This is so in array-slice
// operations like 'v = myarray[2:8]' the callsite can hand back slots 2-5 of
// the array (which is the same way Python handles beyond-the-end indexing).

// Examples with n = 5:
//
// mindex zindex ok
// -7    -2      false
// -6    -1      false
// -5     0      true
// -4     1      true
// -3     2      true
// -2     3      true
// -1     4      true
//  0    -1      false
//  1     0      true
//  2     1      true
//  3     2      true
//  4     3      true
//  5     4      true
//  6     5      false
//  7     6      false

func unaliasArrayLengthIndex(n int, mindex int) (int, bool) {
	if 1 <= mindex {
		zindex := mindex - 1
		if mindex <= n { // in bounds
			return zindex, true
		} else { // out of bounds
			return zindex, false
		}
	} else if mindex <= -1 {
		zindex := mindex + n
		if -n <= mindex { // in bounds
			return zindex, true
		} else { // out of bounds
			return zindex, false
		}
	} else {
		// mindex is 0
		return -1, false
	}
}
