// ================================================================
// Constructors
// ================================================================

package mlrval

import (
	"errors"
	"fmt"

	"github.com/johnkerl/miller/internal/pkg/lib"
)

// TODO: comment for JSON-scanner context.
func FromPending() *Mlrval {
	return &Mlrval{
		mvtype:        MT_PENDING,
		printrep:      "(bug-if-you-see-this:case-1)",
		printrepValid: false,
	}
}

// TODO: comment JIT context. Some things we already know are typed -- DSL
// things, or JSON contents.  Others are deferred, e.g. items from any file
// format except JSON.
// TODO: comment re inferBool.
func FromDeferredType(input string) *Mlrval {
	return &Mlrval{
		mvtype:        MT_PENDING,
		printrep:      input,
		printrepValid: true,
	}
}

func FromError(err error) *Mlrval {
	return &Mlrval{
		mvtype:        MT_ERROR,
		err:           err,
		printrep:      ERROR_PRINTREP,
		printrepValid: true,
	}
}

func FromErrorString(err string) *Mlrval {
	return &Mlrval{
		mvtype:        MT_ERROR,
		err:           errors.New(err),
		printrep:      ERROR_PRINTREP,
		printrepValid: true,
	}
}

func FromAnonymousError() *Mlrval {
	return &Mlrval{
		mvtype:        MT_ERROR,
		printrep:      ERROR_PRINTREP,
		printrepValid: true,
	}
}

func FromTypeErrorUnary(funcname string, input1 *Mlrval) *Mlrval {
	return FromError(
		fmt.Errorf(
			"%s: unacceptable type %s with value %s",
			funcname,
			input1.GetTypeName(),
			input1.StringMaybeQuoted(),
		),
	)
}

func FromTypeErrorBinary(funcname string, input1, input2 *Mlrval) *Mlrval {
	return FromError(
		fmt.Errorf(
			"%s: unacceptable types %s, %s with values %s, %s",
			funcname,
			input1.GetTypeName(),
			input2.GetTypeName(),
			input1.StringMaybeQuoted(),
			input2.StringMaybeQuoted(),
		),
	)
}

func FromTypeErrorTernary(funcname string, input1, input2, input3 *Mlrval) *Mlrval {
	return FromError(
		fmt.Errorf(
			"%s: unacceptable types %s, %s, %s with values %s, %s, %s",
			funcname,
			input1.GetTypeName(),
			input2.GetTypeName(),
			input3.GetTypeName(),
			input1.StringMaybeQuoted(),
			input2.StringMaybeQuoted(),
			input3.StringMaybeQuoted(),
		),
	)
}

// TODO: comment non-JIT context like mlr put -s.
// TODO: comment re inferBool.
func FromInferredType(input string) *Mlrval {
	mv := &Mlrval{
		mvtype:        MT_PENDING,
		printrep:      input,
		printrepValid: true,
	}
	// TODO: comment re data files vs literals context -- this is for the latter
	if input == "true" {
		return TRUE
	} else if input == "false" {
		return FALSE
	} else {
		packageLevelInferrer(mv)
		return mv
	}
}

func FromString(input string) *Mlrval {
	if input == "" {
		return VOID
	}
	return &Mlrval{
		mvtype:        MT_STRING,
		printrep:      input,
		printrepValid: true,
	}
}

func (mv *Mlrval) SetFromString(input string) *Mlrval {
	mv.printrep = input
	mv.printrepValid = true
	if input == "" {
		mv.mvtype = MT_VOID
	} else {
		mv.mvtype = MT_STRING
	}
	return mv
}

func (mv *Mlrval) SetFromVoid() *Mlrval {
	mv.printrep = ""
	mv.printrepValid = true
	mv.mvtype = MT_VOID
	return mv
}

func FromInt(input int64) *Mlrval {
	return &Mlrval{
		mvtype:        MT_INT,
		printrepValid: false,
		intf:          input,
	}
}

// TryFromIntString is used by the mlrval Formatter (fmtnum DSL function,
// format-values verb, etc).  Each mlrval has printrep and a printrepValid for
// its original string, then a type-code like MT_INT or MT_FLOAT, and
// type-specific storage like intval or floatval.
//
// If the user has taken a mlrval with original string "314" and formatted it
// with "0x%04x" then its printrep will be "0x013a" but its type should still
// be MT_INT.
//
// If on the other hand the user has formatted the same mlrval with
// "[[%0x04x]]" then its printrep will be "[[0x013a]]" and it will be
// MT_STRING.  This function supports that.
func TryFromIntString(input string) *Mlrval {
	intval, ok := lib.TryIntFromString(input)
	if ok {
		return FromPrevalidatedIntString(input, intval)
	} else {
		return FromString(input)
	}
}

// TODO: comment
func (mv *Mlrval) SetFromPrevalidatedIntString(input string, intval int64) *Mlrval {
	mv.printrep = input
	mv.printrepValid = true
	mv.intf = intval
	mv.mvtype = MT_INT
	return mv
}

// TODO: comment
func FromPrevalidatedIntString(input string, intval int64) *Mlrval {
	mv := &Mlrval{}
	mv.SetFromPrevalidatedIntString(input, intval)
	return mv
}

func FromFloat(input float64) *Mlrval {
	return &Mlrval{
		mvtype:        MT_FLOAT,
		printrepValid: false,
		intf:          input,
	}
}

// TryFromFloatString is used by the mlrval Formatter (fmtnum DSL function,
// format-values verb, etc).  Each mlrval has printrep and a printrepValid for
// its original string, then a type-code like MT_INT or MT_FLOAT, and
// type-specific storage like intval or floatval.
//
// If the user has taken a mlrval with original string "3.14" and formatted it
// with "%.4f" then its printrep will be "3.1400" but its type should still be
// MT_FLOAT.
//
// If on the other hand the user has formatted the same mlrval with "[[%.4f]]"
// then its printrep will be "[[3.1400]]" and it will be MT_STRING.  This
// function supports that.
func TryFromFloatString(input string) *Mlrval {
	floatval, ok := lib.TryFloatFromString(input)
	if ok {
		return FromPrevalidatedFloatString(input, floatval)
	} else {
		return FromString(input)
	}
}

// TODO: comment
func (mv *Mlrval) SetFromPrevalidatedFloatString(input string, floatval float64) *Mlrval {
	mv.printrep = input
	mv.printrepValid = true
	mv.intf = floatval
	mv.mvtype = MT_FLOAT
	return mv
}

// TODO: comment
func FromPrevalidatedFloatString(input string, floatval float64) *Mlrval {
	mv := &Mlrval{}
	mv.SetFromPrevalidatedFloatString(input, floatval)
	return mv
}

func FromBool(input bool) *Mlrval {
	if input == true {
		return TRUE
	} else {
		return FALSE
	}
}

func FromBoolString(input string) *Mlrval {
	if input == "true" {
		return TRUE
	} else if input == "false" {
		return FALSE
	} else {
		lib.InternalCodingErrorIf(true)
		return nil // not reached
	}
}

// TODO: comment
func (mv *Mlrval) SetFromPrevalidatedBoolString(input string, boolval bool) *Mlrval {
	mv.printrep = input
	mv.printrepValid = true
	mv.intf = boolval
	mv.mvtype = MT_BOOL
	return mv
}

// The user-defined function is of type 'interface{}' here to avoid what would
// otherwise be a package-dependency cycle between this package and
// github.com/johnkerl/miller/internal/pkg/dsl/cst.
//
// Nominally the name argument is the user-specified name if `func f(a, b) {
// ... }`, or some autogenerated UUID like `fl0052` if `func (a, b) { ... }`.

func FromFunction(funcval interface{}, name string) *Mlrval {
	return &Mlrval{
		mvtype:        MT_FUNC,
		printrep:      name,
		printrepValid: true,
		intf:          funcval,
	}
}

func FromArray(arrayval []*Mlrval) *Mlrval {
	return &Mlrval{
		mvtype:        MT_ARRAY,
		printrep:      "(bug-if-you-see-this:case-4)", // INVALID_PRINTREP,
		printrepValid: false,
		intf:          CopyMlrvalArray(arrayval),
	}
}

func FromSingletonArray(element *Mlrval) *Mlrval {
	a := make([]*Mlrval, 1)
	a[0] = element
	return FromArray(a)
}

func FromEmptyArray() *Mlrval {
	return FromArray(make([]*Mlrval, 0))
}

func FromMap(mapval *Mlrmap) *Mlrval {
	return &Mlrval{
		mvtype:        MT_MAP,
		printrep:      "(bug-if-you-see-this:case-5)", // INVALID_PRINTREP,
		printrepValid: false,
		intf:          mapval.Copy(),
	}
}

func FromEmptyMap() *Mlrval {
	return FromMap(NewMlrmap())
}
