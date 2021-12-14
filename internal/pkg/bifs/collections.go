package types

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ================================================================
// Map/array count. Scalars (including strings) have length 1.
func BIF_length(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	switch input1.Type() {
	case mlrval.MT_ERROR:
		return mlrval.FromInt(0)
		break
	case mlrval.MT_ABSENT:
		return mlrval.FromInt(0)
		break
	case mlrval.MT_ARRAY:
		return mlrval.FromInt(int(len(input1.AcquireArrayValue().([]mlrval.Mlrval))))
		break
	case mlrval.MT_MAP:
		return mlrval.FromInt(int(input1.AcquireArrayValue().(types.Mlrmap).FieldCount))
		break
	}
	return mlrval.FromInt(1)
}

// ================================================================
func depth_from_array(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	maxChildDepth := 0
	for _, child := range input1.AcquireArrayValue().([]mlrval.Mlrval) {
		childDepth := BIF_depth(&child)
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
	for pe := input1.AcquireArrayValue().(types.Mlrmap).Head; pe != nil; pe = pe.Next {
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

func BIF_depth(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return depth_dispositions[input1.Type()](input1)
}

// ================================================================
func leafcount_from_array(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	sumChildLeafCount := 0
	for _, child := range input1.AcquireArrayValue().([]mlrval.Mlrval) {
		// Golang initialization loop if we do this :(
		// childLeafCount := BIF_leafcount(&child)

		childLeafCount := mlrval.FromInt(1)
		if child.IsArray() {
			childLeafCount = leafcount_from_array(&child)
		} else if child.IsMap() {
			childLeafCount = leafcount_from_map(&child)
		}

		lib.InternalCodingErrorIf(!childLeafCount.IsInt())
		iChildLeafCount := int(childLeafCount.AcquireIntValue())
		sumChildLeafCount += iChildLeafCount
	}
	return mlrval.FromInt(int(sumChildLeafCount))
}

func leafcount_from_map(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	sumChildLeafCount := 0
	for pe := input1.AcquireArrayValue().(types.Mlrmap).Head; pe != nil; pe = pe.Next {
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
	_, ok := UnaliasArrayIndex(&input1.AcquireArrayValue().([]mlrval.Mlrval), input2.AcquireIntValue())
	return mlrval.FromBool(ok)
}

func has_key_in_map(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if input2.IsString() || input2.IsInt() {
		return mlrval.FromBool(input1.AcquireArrayValue().(types.Mlrmap).Has(input2.String()))
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
func BIF_mapselect(mlrvals []*mlrval.Mlrval) *mlrval.Mlrval {
	if len(mlrvals) < 1 {
		return mlrval.ERROR
	}
	if !mlrvals[0].IsMap() {
		return mlrval.ERROR
	}
	oldmap := mlrvals[0].AcquireArrayValue().(types.Mlrmap)
	newMap := NewMlrmap()

	newKeys := make(map[string]bool)
	for _, selectArg := range mlrvals[1:] {
		if selectArg.IsString() {
			newKeys[selectArg.AcquireStringValue()] = true
		} else if selectArg.IsInt() {
			newKeys[selectArg.String()] = true
		} else if selectArg.IsArray() {
			for _, element := range selectArg.AcquireArrayValue().([]mlrval.Mlrval) {
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
	newMap := mlrvals[0].AcquireArrayValue().(types.Mlrmap).Copy()

	for _, exceptArg := range mlrvals[1:] {
		if exceptArg.IsString() {
			newMap.Remove(exceptArg.AcquireStringValue())
		} else if exceptArg.IsInt() {
			newMap.Remove(exceptArg.String())
		} else if exceptArg.IsArray() {
			for _, element := range exceptArg.AcquireArrayValue().([]mlrval.Mlrval) {
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
	if mlrvals[0].Type() != MT_MAP {
		return mlrval.ERROR
	}
	newMap := mlrvals[0].AcquireArrayValue().(types.Mlrmap).Copy()

	for _, otherMapArg := range mlrvals[1:] {
		if otherMapArg.Type() != MT_MAP {
			return mlrval.ERROR
		}

		for pe := otherMapArg.AcquireArrayValue().(types.Mlrmap).Head; pe != nil; pe = pe.Next {
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
	if mlrvals[0].Type() != MT_MAP {
		return mlrval.ERROR
	}
	newMap := mlrvals[0].AcquireArrayValue().(types.Mlrmap).Copy()

	for _, otherMapArg := range mlrvals[1:] {
		if otherMapArg.Type() != MT_MAP {
			return mlrval.ERROR
		}

		for pe := otherMapArg.AcquireArrayValue().(types.Mlrmap).Head; pe != nil; pe = pe.Next {
			newMap.Remove(pe.Key)
		}
	}

	return mlrval.FromMap(newMap)
}

// ================================================================
// joink([1,2,3], ",") -> "1,2,3"
// joink({"a":3,"b":4,"c":5}, ",") -> "a,b,c"
func BIF_joink(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if input2.Type() != MT_STRING {
		return mlrval.ERROR
	}
	fieldSeparator := input2.AcquireStringValue()
	if input1.IsMap() {
		var buffer bytes.Buffer

		for pe := input1.AcquireArrayValue().(types.Mlrmap).Head; pe != nil; pe = pe.Next {
			buffer.WriteString(pe.Key)
			if pe.Next != nil {
				buffer.WriteString(fieldSeparator)
			}
		}

		return mlrval.FromString(buffer.String())
	} else if input1.IsArray() {
		var buffer bytes.Buffer

		for i := range input1.AcquireArrayValue().([]mlrval.Mlrval) {
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
	if input2.Type() != MT_STRING {
		return mlrval.ERROR
	}
	fieldSeparator := input2.AcquireStringValue()

	if input1.IsMap() {
		var buffer bytes.Buffer

		for pe := input1.AcquireArrayValue().(types.Mlrmap).Head; pe != nil; pe = pe.Next {
			buffer.WriteString(pe.Value.String())
			if pe.Next != nil {
				buffer.WriteString(fieldSeparator)
			}
		}

		return mlrval.FromString(buffer.String())
	} else if input1.IsArray() {
		var buffer bytes.Buffer

		for i, element := range input1.AcquireArrayValue().([]mlrval.Mlrval) {
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
	if input2.Type() != MT_STRING {
		return mlrval.ERROR
	}
	pairSeparator := input2.AcquireStringValue()
	if input3.Type() != MT_STRING {
		return mlrval.ERROR
	}
	fieldSeparator := input3.AcquireStringValue()

	if input1.IsMap() {
		var buffer bytes.Buffer

		for pe := input1.AcquireArrayValue().(types.Mlrmap).Head; pe != nil; pe = pe.Next {
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

		for i, element := range input1.AcquireArrayValue().([]mlrval.Mlrval) {
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
	if input2.Type() != MT_STRING {
		return mlrval.ERROR
	}
	pairSeparator := input2.AcquireStringValue()
	if input3.Type() != MT_STRING {
		return mlrval.ERROR
	}
	fieldSeparator := input3.AcquireStringValue()

	output := mlrval.FromEmptyMap()

	fields := lib.SplitString(input1.AcquireStringValue(), fieldSeparator)
	for i, field := range fields {
		pair := strings.SplitN(field, pairSeparator, 2)
		if len(pair) == 1 {
			key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
			value := mlrval.FromInferredType(pair[0])
			output.AcquireArrayValue().(types.Mlrmap).PutReference(key, value)
		} else if len(pair) == 2 {
			key := pair[0]
			value := mlrval.FromInferredType(pair[1])
			output.AcquireArrayValue().(types.Mlrmap).PutReference(key, value)
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
	if input2.Type() != MT_STRING {
		return mlrval.ERROR
	}
	pairSeparator := input2.AcquireStringValue()
	if input3.Type() != MT_STRING {
		return mlrval.ERROR
	}
	fieldSeparator := input3.AcquireStringValue()

	output := mlrval.FromEmptyMap()

	fields := lib.SplitString(input1.AcquireStringValue(), fieldSeparator)
	for i, field := range fields {
		pair := strings.SplitN(field, pairSeparator, 2)
		if len(pair) == 1 {
			key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
			value := mlrval.FromString(pair[0])
			output.AcquireArrayValue().(types.Mlrmap).PutReference(key, value)
		} else if len(pair) == 2 {
			key := pair[0]
			value := mlrval.FromString(pair[1])
			output.AcquireArrayValue().(types.Mlrmap).PutReference(key, value)
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
	if input2.Type() != MT_STRING {
		return mlrval.ERROR
	}

	output := mlrval.FromEmptyMap()

	fields := lib.SplitString(input1.AcquireStringValue(), input2.AcquireStringValue())
	for i, field := range fields {
		key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
		value := mlrval.FromInferredType(field)
		output.AcquireArrayValue().(types.Mlrmap).PutReference(key, value)
	}

	return output
}

// ----------------------------------------------------------------
// splitnvx("3,4,5", ",") -> {"1":"3","2":"4","3":"5"}
func BIF_splitnvx(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsStringOrVoid() {
		return mlrval.ERROR
	}
	if input2.Type() != MT_STRING {
		return mlrval.ERROR
	}

	output := mlrval.FromEmptyMap()

	fields := lib.SplitString(input1.AcquireStringValue(), input2.AcquireStringValue())
	for i, field := range fields {
		key := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
		value := mlrval.FromString(field)
		output.AcquireArrayValue().(types.Mlrmap).PutReference(key, value)
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

	output := NewSizedMlrvalArray(int(len(fields)))

	for i, field := range fields {
		value := mlrval.FromInferredType(field)
		output.AcquireArrayValue().([]mlrval.Mlrval)[i] = *value
	}

	return output
}

// ----------------------------------------------------------------
// BIF_splitax splits a string to an array, without type-inference:
// e.g. splitax("3,4,5", ",") -> ["3","4","5"]
func BIF_splitax(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsStringOrVoid() {
		return mlrval.ERROR
	}
	if input2.Type() != MT_STRING {
		return mlrval.ERROR
	}
	input := input1.AcquireStringValue()
	fieldSeparator := input2.AcquireStringValue()

	return mlrvalSplitAXHelper(input, fieldSeparator)
}

// mlrvalSplitAXHelper is Split out for the benefit of BIF_splitax and
// BIF_unflatten.
func mlrvalSplitAXHelper(input string, separator string) *mlrval.Mlrval {
	fields := lib.SplitString(input, separator)

	output := NewSizedMlrvalArray(int(len(fields)))

	for i, field := range fields {
		output.AcquireArrayValue().([]mlrval.Mlrval)[i] = *mlrval.FromString(field)
	}

	return output
}

// ----------------------------------------------------------------
func BIF_get_keys(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsMap() {
		// TODO: make a ReferenceFrom with comments
		output := NewSizedMlrvalArray(input1.AcquireArrayValue().(types.Mlrmap).FieldCount)
		i := 0
		for pe := input1.AcquireArrayValue().(types.Mlrmap).Head; pe != nil; pe = pe.Next {
			output.AcquireArrayValue().([]mlrval.Mlrval)[i] = *mlrval.FromString(pe.Key)
			i++
		}
		return output

	} else if input1.IsArray() {
		output := NewSizedMlrvalArray(int(len(input1.AcquireArrayValue().([]mlrval.Mlrval))))
		for i := range input1.AcquireArrayValue().([]mlrval.Mlrval) {
			output.AcquireArrayValue().([]mlrval.Mlrval)[i] = *mlrval.FromInt(int(i + 1)) // Miller user-space indices are 1-up
		}
		return output

	} else {
		return mlrval.ERROR
	}
}

// ----------------------------------------------------------------
func BIF_get_values(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsMap() {
		// TODO: make a ReferenceFrom with commenbs
		output := NewSizedMlrvalArray(input1.AcquireArrayValue().(types.Mlrmap).FieldCount)
		i := 0
		for pe := input1.AcquireArrayValue().(types.Mlrmap).Head; pe != nil; pe = pe.Next {
			output.AcquireArrayValue().([]mlrval.Mlrval)[i] = *pe.Value.Copy()
			i++
		}
		return output

	} else if input1.IsArray() {
		output := NewSizedMlrvalArray(int(len(input1.AcquireArrayValue().([]mlrval.Mlrval))))
		for i, value := range input1.AcquireArrayValue().([]mlrval.Mlrval) {
			output.AcquireArrayValue().([]mlrval.Mlrval)[i] = *value.Copy()
		}
		return output

	} else {
		return mlrval.ERROR
	}
}

// ----------------------------------------------------------------
func BIF_append(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.Type() != MT_ARRAY {
		return mlrval.ERROR
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
func BIF_flatten(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	if input3.IsMap() || input3.IsArray() {
		if input1.Type() != MT_STRING && input1.Type() != MT_VOID {
			return mlrval.ERROR
		}
		prefix := input1.AcquireStringValue()
		if input2.Type() != MT_STRING {
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
	if input2.Type() != MT_STRING {
		return mlrval.ERROR
	}
	if input1.Type() != MT_MAP {
		return input1
	}
	oldmap := input1.AcquireArrayValue().(types.Mlrmap)
	separator := input2.AcquireStringValue()
	newmap := oldmap.CopyUnflattened(separator)
	return mlrval.FromMapReferenced(newmap)
}

// ----------------------------------------------------------------
// Converts maps with "1", "2", ... keys into arrays. Recurses nested data structures.
func BIF_arrayify(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsMap() {
		if input1.AcquireArrayValue().(types.Mlrmap).IsEmpty() {
			return input1
		}

		convertible := true
		i := 0
		for pe := input1.AcquireArrayValue().(types.Mlrmap).Head; pe != nil; pe = pe.Next {
			sval := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
			i++
			if pe.Key != sval {
				convertible = false
			}
			pe.Value = BIF_arrayify(pe.Value)
		}

		if convertible {
			AcquireArrayValue().([]mlrval.Mlrval) := make([]mlrval.Mlrval, input1.AcquireArrayValue().(types.Mlrmap).FieldCount)
			i := 0
			for pe := input1.AcquireArrayValue().(types.Mlrmap).Head; pe != nil; pe = pe.Next {
				AcquireArrayValue().([]mlrval.Mlrval)[i] = *pe.Value.Copy()
				i++
			}
			return mlrval.FromArrayReference(AcquireArrayValue().([]mlrval.Mlrval))

		} else {
			return input1
		}

	} else if input1.IsArray() {
		// TODO: comment (or rethink) that this modifies its inputs!!
		output := input1.Copy()
		for i := range input1.AcquireArrayValue().([]mlrval.Mlrval) {
			output.AcquireArrayValue().([]mlrval.Mlrval)[i] = *BIF_arrayify(&output.AcquireArrayValue().([]mlrval.Mlrval)[i])
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
	} else if input1.Type() != MT_STRING {
		return mlrval.ERROR
	} else {
		output := mlrval.FromPending()
		err := output.UnmarshalJSON([]byte(input1.AcquireStringValue()))
		if err != nil {
			return mlrval.ERROR
		}
		return &output
	}
}

func BIF_json_stringify_unary(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	outputBytes, err := input1.MarshalJSON(JSON_SINGLE_LINE, false)
	if err != nil {
		return mlrval.ERROR
	} else {
		return mlrval.FromString(string(outputBytes))
	}
}

func BIF_json_stringify_binary(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	var jsonFormatting TJSONFormatting = JSON_SINGLE_LINE
	useMultiline, ok := input2.GetBoolValue()
	if !ok {
		return mlrval.ERROR
	}
	if useMultiline {
		jsonFormatting = JSON_MULTILINE
	}

	outputBytes, err := input1.MarshalJSON(jsonFormatting, false)
	if err != nil {
		return mlrval.ERROR
	} else {
		return mlrval.FromString(string(outputBytes))
	}
}
