// ================================================================
// Constants/singletons of various types
// ================================================================

package mlrval

// MlrvalFromPending is designed solely for the JSON API, for something
// intended to be mutated after construction once its type is (later) known.
// Whereas ERROR, ABSENT, etc are all singletons, this one
// must be mutable and therefor non-singleton.

func MlrvalFromPending() Mlrval {
	return Mlrval{
		mvtype:   MT_PENDING,
		printrep: "(bug-if-you-see-this:case-3)", // INVALID_PRINTREP,
	}
}

// These are made singletons as part of a copy-reduction effort.  They're not
// marked const (I haven't figured out the right way to get that to compile;
// just using `const` isn't enough) but the gentelpersons' agreement is that
// the caller should never modify these.

var TRUE = &Mlrval{
	mvtype:        MT_BOOL,
	printrep:      "true",
	printrepValid: true,
	boolval:       true,
}

var FALSE = &Mlrval{
	mvtype:        MT_BOOL,
	printrep:      "false",
	printrepValid: true,
}

var VOID = &Mlrval{
	mvtype:        MT_VOID,
	printrep:      "",
	printrepValid: true,
}

var ERROR = &Mlrval{
	mvtype:        MT_ERROR,
	printrep:      ERROR_PRINTREP,
	printrepValid: true,
}

var ABSENT = &Mlrval{
	mvtype:        MT_ABSENT,
	printrep:      ABSENT_PRINTREP,
	printrepValid: true,
}

var NULL = &Mlrval{
	mvtype:        MT_NULL,
	printrep:      "null",
	printrepValid: true,
}
