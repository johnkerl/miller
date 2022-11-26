package mlrval

// TODO: comment about mvtype; deferrence; copying of deferrence.
func (mv *Mlrval) Copy() *Mlrval {
	other := *mv
	if mv.mvtype == MT_MAP {
		other.x = &mlrvalExtended{
			mapval: mv.x.mapval.Copy(),
		}
	} else if mv.mvtype == MT_ARRAY {
		other.x = &mlrvalExtended{
			arrayval: CopyMlrvalArray(mv.x.arrayval),
		}
	} else if mv.mvtype == MT_FUNC {
		other.x = &mlrvalExtended{
			funcval: mv.x.funcval,
		}
	}
	return &other
}
