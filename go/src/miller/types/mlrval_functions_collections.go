package types

import ()

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
	if mb.mvtype != MT_INT {
		return MlrvalFromError()
	}
	_, ok := unaliasArrayIndex(&ma.arrayval, mb.intval)
	return MlrvalFromBool(ok)
}

func has_key_in_map(ma, mb *Mlrval) Mlrval {
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
