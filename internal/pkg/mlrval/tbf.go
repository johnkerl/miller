package mlrval

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/johnkerl/miller/internal/pkg/lib"
)

// things to be filed

// NewMlrvalForAutoDeepen is for auto-deepen of nested maps in things like
//
//   $foo[1]["a"][2]["b"] = 3
//
// Autocreated levels are maps.  Array levels can be explicitly created e.g.
//
//   $foo[1]["a"] ??= []
//   $foo[1]["a"][2]["b"] = 3
func NewMlrvalForAutoDeepen(mvtype MVType) (*Mlrval, error) {
	if mvtype == MT_STRING || mvtype == MT_INT {
		empty := FromEmptyMap()
		return empty, nil
	} else {
		return nil, errors.New(
			"mlr: indices must be string, int, or array thereof; got " + GetTypeName(mvtype),
		)
	}
}

func (mv *Mlrval) Arrayify() *Mlrval {
	if mv.IsMap() {
		if mv.mapval.IsEmpty() {
			return mv
		}

		convertible := true
		i := 0
		for pe := mv.mapval.Head; pe != nil; pe = pe.Next {
			sval := strconv.Itoa(i + 1) // Miller user-space indices are 1-up
			i++
			if pe.Key != sval {
				convertible = false
			}
			pe.Value = pe.Value.Arrayify()
		}

		if convertible {
			arrayval := make([]Mlrval, mv.mapval.FieldCount)
			i := 0
			for pe := mv.mapval.Head; pe != nil; pe = pe.Next {
				arrayval[i] = *pe.Value.Copy()
				i++
			}
			return FromArray(arrayval)

		} else {
			return mv
		}

	} else if mv.IsArray() {
		// TODO: comment (or rethink) that this modifies its inputs!!
		output := mv.Copy()
		for i := range mv.arrayval {
			output.arrayval[i] = *output.arrayval[i].Arrayify()
		}
		return output

	} else {
		return mv
	}
}

//func MlrvalFromEmptyMap() *Mlrval {
//	return &Mlrval{
//		mvtype:        MT_MAP,
//		printrep:      "(bug-if-you-see-this-map-type)",
//		printrepValid: false,
//		mapval:        NewMlrmap(),
//	}
//}
//
//func MlrvalFromMap(mlrmap *Mlrmap) *Mlrval {
//	mv := MlrvalFromEmptyMap()
//	if mlrmap == nil {
//		// TODO maybe return 2nd-arg error in the API
//		return ERROR
//	}
//
//	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
//		mv.mapval.PutCopy(pe.Key, pe.Value)
//	}
//	return mv
//}
//
//// Like previous but doesn't copy. Only safe when the argument's sole purpose
//// is to be passed into here.
//func MlrvalFromMapReferenced(mlrmap *Mlrmap) *Mlrval {
//	mv := MlrvalFromEmptyMap()
//	if mlrmap == nil {
//		// xxx maybe return 2nd-arg error in the API
//		return ERROR
//	}
//
//	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
//		mv.mapval.PutReference(pe.Key, pe.Value)
//	}
//	return mv
//}
//
//// Does not copy the data. We can make a MlrvalFromArrayLiteralCopy if needed,
//// using values.CopyMlrvalArray().
//func MlrvalEmptyArray() Mlrval {
//	return Mlrval{
//		mvtype:        MT_ARRAY,
//		printrep:      "(bug-if-you-see-this-array-type)",
//		printrepValid: false,
//		intval:        0,
//		floatval:      0.0,
//		boolval:       false,
//		arrayval:      make([]Mlrval, 0, 10),
//		mapval:        nil,
//	}
//}
//
//// Users can do things like '$new[1][2][3] = 4' even if '$new' isn't already
//// allocated. This function supports that.
//func NewSizedMlrvalArray(length int) *Mlrval {
//	arrayval := make([]Mlrval, length, 2*length)
//
//	for i := 0; i < int(length); i++ {
//		arrayval[i] = *VOID
//	}
//
//	return &Mlrval{
//		mvtype:        MT_ARRAY,
//		printrep:      "(bug-if-you-see-this-array-type)",
//		printrepValid: false,
//		intval:        0,
//		floatval:      0.0,
//		boolval:       false,
//		arrayval:      arrayval,
//		mapval:        nil,
//	}
//}
//
//// Does not copy the data. We can make a SetFromArrayLiteralCopy if needed
//// using values.CopyMlrvalArray().
//func MlrvalFromArrayReference(input []Mlrval) *Mlrval {
//	return &Mlrval{
//		mvtype:        MT_ARRAY,
//		printrepValid: false,
//		arrayval:      input,
//	}
//}

func LengthenMlrvalArray(array *[]Mlrval, newLength64 int) {
	newLength := int(newLength64)
	lib.InternalCodingErrorIf(newLength <= len(*array))

	if newLength <= cap(*array) {
		newArray := (*array)[:newLength]
		for zindex := len(*array); zindex < newLength; zindex++ {
			// TODO: comment why not MT_ABSENT or MT_VOID
			newArray[zindex] = *NULL
		}
		*array = newArray
	} else {
		newArray := make([]Mlrval, newLength, 2*newLength)
		zindex := 0
		for zindex = 0; zindex < len(*array); zindex++ {
			newArray[zindex] = (*array)[zindex]
		}
		for zindex = len(*array); zindex < newLength; zindex++ {
			// TODO: comment why not MT_ABSENT or MT_VOID
			newArray[zindex] = *NULL
		}
		*array = newArray
	}
}

//// NewMlrvalForAutoDeepen is for auto-deepen of nested maps in things like
////
////   $foo[1]["a"][2]["b"] = 3
////
//// Autocreated levels are maps.  Array levels can be explicitly created e.g.
////
////   $foo[1]["a"] ??= []
////   $foo[1]["a"][2]["b"] = 3
//func NewMlrvalForAutoDeepen(mvtype MVType) (*Mlrval, error) {
//	if mvtype == MT_STRING || mvtype == MT_INT {
//		empty := MlrvalFromEmptyMap()
//		return empty, nil
//	} else {
//		return nil, errors.New(
//			"mlr: indices must be string, int, or array thereof; got " + GetTypeName(mvtype),
//		)
//	}
//}

// TODO: FILE
//
//// MlrvalFromInferredTypeForDataFiles is for parsing field values directly from
//// data files (except JSON, which is typed -- "true" and true are distinct).
//// Mostly the same as MlrvalFromInferredType, except it doesn't auto-infer
//// true/false to bool; don't auto-infer NaN/Inf to float; etc.
//func MlrvalFromInferredTypeForDataFiles(input string) *Mlrval {
//	return inferrer(input, false)
//}
//
//// MlrvalFromInferredType is for parsing field values not directly from data
//// files.  Mostly the same as MlrvalFromInferredTypeForDataFiles, except it
//// auto-infers true/false to bool; don't auto-infer NaN/Inf to float; etc.
//func MlrvalFromInferredType(input string) *Mlrval {
//	return inferrer(input, true)
//}

func (mv *Mlrval) GetNumericToFloatValueOrDie() (floatValue float64) {
	floatValue, ok := mv.GetNumericToFloatValue()
	if !ok {
		fmt.Fprintf(
			os.Stderr,
			"%s: couldn't parse \"%s\" as number.",
			"mlr", mv.String(),
		)
		os.Exit(1)
	}
	return floatValue
}

func (mv *Mlrval) AssertNumeric() {
	_ = mv.GetNumericToFloatValueOrDie()
}
